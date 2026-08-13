@echo off
setlocal enabledelayedexpansion
rem =====================================================================
rem dev.cmd - schema-ui-core quick local dev launcher (start/stop)
rem
rem Corresponds to QUICKSTART "Path B - local two-process" (dev default):
rem   API  -> apps/api   go run ./cmd/server   (:25080, APP_ENV=development)
rem   Web  -> apps/web   npm run dev           (:25173, /api proxied to :25080)
rem
rem Startup order: API starts first; Web is launched only after the API port
rem   is listening AND /readyz returns HTTP 200 (DB ping + module graph
rem   Start+Ready), so the Vite /api proxy never races an uninitialized backend.
rem
rem Stop precision: each started window records its PID to %TEMP%
rem   (schema-ui-dev-api-<port>.pid / schema-ui-dev-web-<port>.pid). `stop`
rem   kills ONLY those recorded processes (plus the listeners on the two
rem   dev ports as a fallback) - it never matches window titles, so unrelated
rem   console windows (Claude Code, other CLIs) are never touched.
rem
rem Usage:
rem   dev.cmd start [profile] [options]   start API+Web (default profile: admin)
rem   dev.cmd stop                        stop API+Web
rem   dev.cmd status                      show listening state
rem   dev.cmd help                        this help
rem
rem profile:  mvp | admin | demo | custom  (default admin)
rem   custom requires --modules "core.xxx,module.yyy" (full explicit module closure)
rem options:
rem   --profile <name>      explicit profile
rem   --modules <a,b,c>     APP_MODULES_ENABLED override (REQUIRED for custom)
rem   --browser             open web UI in the default browser (default)
rem   --no-browser          do not open a browser
rem env:
rem   API_PORT (default 25080) / WEB_PORT (default 25173)  port overrides
rem   ADMIN_INITIAL_PASSWORD / AUTH_JWT_SECRET pass through to the API if exported
rem     (dev mode: built-in dev secret, seed password defaults to "admin")
rem
rem NOTE: keep this file ASCII-only. Non-ASCII bytes in .cmd are re-parsed by
rem cmd.exe in the OEM codepage and can corrupt the whole script parse.
rem =====================================================================

set "ROOT=%~dp0"
set "SELF=%~f0"
set "API_DIR=%ROOT%apps\api"
set "WEB_DIR=%ROOT%apps\web"

rem ---- defaults ----
set "PROFILE=admin"
set "MODULES="
set "OPEN_BROWSER=1"
if "%API_PORT%"=="" set "API_PORT=25080"
if "%WEB_PORT%"=="" set "WEB_PORT=25173"

set "CMD=%~1"
if "%CMD%"=="" goto :usage
shift

if /i "%CMD%"=="help"   goto :usage
if /i "%CMD%"=="status" goto :status
if /i "%CMD%"=="stop"   goto :stop
if /i "%CMD%"=="start"  goto :parse_start

echo ERROR: unknown command '%CMD%'
goto :usage

rem =====================================================================
rem start argument parsing
rem =====================================================================
:parse_start
if "%~1"=="" goto :start_validated
if /i "%~1"=="--profile"    ( set "PROFILE=%~2"      & shift & shift & goto :parse_start )
if /i "%~1"=="--modules"    ( set "MODULES=%~2"      & shift & shift & goto :parse_start )
if /i "%~1"=="--browser"    ( set "OPEN_BROWSER=1"   & shift & goto :parse_start )
if /i "%~1"=="--no-browser" ( set "OPEN_BROWSER=0"   & shift & goto :parse_start )
if /i "%~1"=="--help"       goto :usage
set "PROFILE=%~1"
shift
goto :parse_start

:start_validated
rem normalize case for known profiles
if /i "%PROFILE%"=="mvp"    set "PROFILE=mvp"
if /i "%PROFILE%"=="admin"  set "PROFILE=admin"
if /i "%PROFILE%"=="demo"   set "PROFILE=demo"
if /i "%PROFILE%"=="custom" set "PROFILE=custom"
if not "%PROFILE%"=="mvp" if not "%PROFILE%"=="admin" if not "%PROFILE%"=="demo" if not "%PROFILE%"=="custom" (
  echo ERROR: unknown profile '%PROFILE%' -- expected mvp, admin, demo or custom
  goto :usage
)
if /i "%PROFILE%"=="custom" if "%MODULES%"=="" (
  echo ERROR: profile=custom requires --modules "core.x,module.y,..." to enumerate the full module closure
  goto :usage
)

rem ---- tool check ----
set "TOOLS_MISSING="
where go  >nul 2>&1 || set "TOOLS_MISSING=go"
where npm >nul 2>&1 || set "TOOLS_MISSING=!TOOLS_MISSING! npm"
where curl >nul 2>&1 || set "TOOLS_MISSING=!TOOLS_MISSING! curl"
if defined TOOLS_MISSING (
  echo ERROR: missing required tools: !TOOLS_MISSING!
  exit /b 2
)

rem ---- install web deps if missing ----
if not exist "%WEB_DIR%\node_modules" (
  echo Installing web dependencies ^(npm ci^) ...
  pushd "%WEB_DIR%"
  call npm ci
  popd
  if errorlevel 1 (
    echo ERROR: npm ci failed
    exit /b 2
  )
)

echo.
echo == schema-ui-core dev ^| profile=%PROFILE% ==
echo   API  :%API_PORT%   %API_DIR%
echo   Web  :%WEB_PORT%   %WEB_DIR%
if defined MODULES echo   APP_MODULES_ENABLED=%MODULES%
echo.

rem ---- launch API ----
set "API_CMD=set APP_ENV=development && set APP_PROFILE=%PROFILE% && set HTTP_ADDR=:%API_PORT%"
if defined MODULES set "API_CMD=!API_CMD! && set APP_MODULES_ENABLED=%MODULES%"
set "API_CMD=!API_CMD! && go run ./cmd/server"

set "API_UP="
netstat -ano | findstr /R /C:":%API_PORT% " | findstr "LISTENING" >nul 2>&1 && set "API_UP=1"
if defined API_UP (
  echo [skip] API already listening on :%API_PORT%
) else (
  echo [start] API ... go run ./cmd/server
  set "DEV_CMD=!API_CMD!"
  set "DEV_DIR=%API_DIR%"
  set "DEV_PIDFILE=%TEMP%\schema-ui-dev-api-%API_PORT%.pid"
  powershell -NoProfile -Command "$p = Start-Process cmd.exe -ArgumentList '/c',$env:DEV_CMD -WorkingDirectory $env:DEV_DIR -PassThru; $p.Id | Out-File -Encoding ascii $env:DEV_PIDFILE"
)

rem ---- gate: Web must NOT start until API is up AND fully ready ----
echo.
echo Waiting for API to come up ...
call :wait_listening "%API_PORT%" "API"
if errorlevel 1 (
  echo ERROR: API not listening on :%API_PORT% within ~60s -- Web NOT started.
  exit /b 1
)

call :wait_ready "%API_PORT%"
if errorlevel 1 (
  echo ERROR: API /readyz did not return 200 within ~60s -- Web NOT started.
  echo   Check the API window for errors ^(migrations / module Start+Ready^).
  exit /b 1
)

rem ---- launch Web (only after API is fully ready) ----
set "WEB_CMD=set WEB_PORT=%WEB_PORT% && npm run dev"
set "WEB_UP="
netstat -ano | findstr /R /C:":%WEB_PORT% " | findstr "LISTENING" >nul 2>&1 && set "WEB_UP=1"
if defined WEB_UP (
  echo [skip] Web already listening on :%WEB_PORT%
) else (
  echo [start] Web ... npm run dev
  set "DEV_CMD=!WEB_CMD!"
  set "DEV_DIR=%WEB_DIR%"
  set "DEV_PIDFILE=%TEMP%\schema-ui-dev-web-%WEB_PORT%.pid"
  powershell -NoProfile -Command "$p = Start-Process cmd.exe -ArgumentList '/c',$env:DEV_CMD -WorkingDirectory $env:DEV_DIR -PassThru; $p.Id | Out-File -Encoding ascii $env:DEV_PIDFILE"
)

rem ---- wait for Web + open browser ----
echo.
echo Waiting for Web to come up ...
call :wait_listening "%WEB_PORT%" "Web"

if "%OPEN_BROWSER%"=="1" (
  echo Opening browser: http://localhost:%WEB_PORT%
  start "" "http://localhost:%WEB_PORT%"
) else (
  echo Browser disabled. Open manually: http://localhost:%WEB_PORT%
)
if not defined ADMIN_INITIAL_PASSWORD (
  echo Login: admin / admin   ^(dev default; export ADMIN_INITIAL_PASSWORD to override^)
)
echo Stop:  %SELF% stop
echo.
exit /b 0

rem =====================================================================
rem wait until a port starts listening (max ~60s)
rem exit code: 0 = listening, 1 = timed out
rem usage: call :wait_listening <port> <label>
rem =====================================================================
:wait_listening
set /a WL_N=0
:wait_listening_loop
netstat -ano | findstr /R /C:":%~1 " | findstr "LISTENING" >nul 2>&1
if not errorlevel 1 (
  echo   ok   %~2 listening on :%~1
  exit /b 0
)
set /a WL_N+=1
if %WL_N% geq 60 (
  echo   error %~2 not listening on :%~1 after ~60s ^(check its window for errors^)
  exit /b 1
)
timeout /t 1 /nobreak >nul 2>&1 || ping -n 2 127.0.0.1 >nul 2>&1
goto :wait_listening_loop

rem =====================================================================
rem wait until API /readyz returns HTTP 200 (SQLite ping + module graph
rem Start+Ready), so Web starts against a fully-initialized backend.
rem exit code: 0 = ready, 1 = timed out
rem usage: call :wait_ready <port>
rem =====================================================================
:wait_ready
set /a WR_N=0
set "WR_CODE=000"
:wait_ready_loop
for /f "delims=" %%c in ('curl.exe -s -o nul --max-time 2 -w "%%{http_code}" http://127.0.0.1:%~1/readyz 2^>nul') do set "WR_CODE=%%c"
if "%WR_CODE%"=="200" (
  echo   ok   API ready ^(/readyz 200^)
  exit /b 0
)
set /a WR_N+=1
if %WR_N% geq 60 (
  echo   error API /readyz not 200 within ~60s ^(last code %WR_CODE%^)
  exit /b 1
)
timeout /t 1 /nobreak >nul 2>&1 || ping -n 2 127.0.0.1 >nul 2>&1
goto :wait_ready_loop

rem =====================================================================
rem stop - kill only this script's own windows via recorded PID files,
rem then fall back to the two dev-port listeners (idempotent).
rem Window-title matching is deliberately NOT used: it can kill unrelated
rem console windows whose titles merely start with the same prefix.
rem =====================================================================
:stop
echo.
echo == Stopping schema-ui-core ==
set "FOUND="

rem 1) recorded PIDs (the exact windows this script spawned)
set "API_PIDFILE=%TEMP%\schema-ui-dev-api-%API_PORT%.pid"
set "WEB_PIDFILE=%TEMP%\schema-ui-dev-web-%WEB_PORT%.pid"
if exist "%API_PIDFILE%" (
  set /p API_PID=<"%API_PIDFILE%"
  if defined API_PID (
    taskkill /PID !API_PID! /T /F >nul 2>&1 && set "FOUND=1"
  )
  del "%API_PIDFILE%" >nul 2>&1
)
if exist "%WEB_PIDFILE%" (
  set /p WEB_PID=<"%WEB_PIDFILE%"
  if defined WEB_PID (
    taskkill /PID !WEB_PID! /T /F >nul 2>&1 && set "FOUND=1"
  )
  del "%WEB_PIDFILE%" >nul 2>&1
)

rem 2) port listeners (covers manually started services / stale pid files)
for /f "tokens=5" %%p in ('netstat -ano ^| findstr /R /C:":%API_PORT% " ^| findstr "LISTENING"') do (
  taskkill /PID %%p /T /F >nul 2>&1
  set "FOUND=1"
)
for /f "tokens=5" %%p in ('netstat -ano ^| findstr /R /C:":%WEB_PORT% " ^| findstr "LISTENING"') do (
  taskkill /PID %%p /T /F >nul 2>&1
  set "FOUND=1"
)

if defined FOUND (
  echo   stopped services ^(:%API_PORT% / :%WEB_PORT%^)
) else (
  echo   nothing running on :%API_PORT% or :%WEB_PORT%
)
exit /b 0

rem =====================================================================
rem status - show listening state
rem =====================================================================
:status
echo.
echo == schema-ui-core status ==
set "A="
netstat -ano | findstr /R /C:":%API_PORT% " | findstr "LISTENING" >nul 2>&1 && set "A=UP" || set "A=down"
set "W="
netstat -ano | findstr /R /C:":%WEB_PORT% " | findstr "LISTENING" >nul 2>&1 && set "W=UP" || set "W=down"
echo   API :%API_PORT%  %A%
echo   Web :%WEB_PORT%  %W%
exit /b 0

rem =====================================================================
rem usage
rem =====================================================================
:usage
echo.
echo schema-ui-core dev launcher (local two-process)
echo.
echo Usage:
echo   dev.cmd start [profile] [options]   Start API+Web locally (default profile: admin)
echo   dev.cmd stop                        Stop API+Web
echo   dev.cmd status                      Show listening state
echo   dev.cmd help                        This help
echo.
echo profile:   mvp ^| admin ^| demo ^| custom   (default admin)
echo options:
echo   --profile ^<name^>      explicit profile
echo   --modules ^<a,b,c^>     APP_MODULES_ENABLED override; REQUIRED for custom
echo   --browser               open web UI in browser (default)
echo   --no-browser            do not open browser
echo env:
echo   API_PORT / WEB_PORT     port overrides (defaults 25080 / 25173)
echo   ADMIN_INITIAL_PASSWORD  dev seed password default "admin"
echo.
echo examples:
echo   dev.cmd start
echo   dev.cmd start demo --no-browser
echo   dev.cmd start custom --modules core.server-registration,users
echo   dev.cmd stop
echo.
exit /b 1
