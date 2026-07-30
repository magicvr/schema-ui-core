# F-018 · PowerShell isolated install smoke test
# Run from anywhere:
#   powershell -NoProfile -ExecutionPolicy Bypass -File skills/tests/test_install_ps1_isolated.ps1
# Exit 0 = pass; non-zero = fail. Prints JSON-ish summary lines for logs.

$ErrorActionPreference = 'Stop'

$TestsDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$PackageRoot = Split-Path -Parent $TestsDir
$InstallPs1 = Join-Path $PackageRoot 'install.ps1'

if (-not (Test-Path -LiteralPath $InstallPs1 -PathType Leaf)) {
    Write-Error "install.ps1 not found: $InstallPs1"
    exit 2
}

$TempRoot = Join-Path ([System.IO.Path]::GetTempPath()) ("gg-skills-install-" + [guid]::NewGuid().ToString('N'))
New-Item -ItemType Directory -Path $TempRoot -Force | Out-Null
$SkillsDest = Join-Path $TempRoot 'skills'

Write-Host "F-018 isolated install"
Write-Host "  package: $PackageRoot"
Write-Host "  target:  $TempRoot"

try {
    Push-Location $TempRoot
    & $InstallPs1 -All -SkillsDir $SkillsDest `
        -InitWorkspace `
        -WorkspaceSlug 'pilot-app' `
        -RootSlug 'pilot-vision' `
        -RootTitle 'Pilot vision'
    if ($LASTEXITCODE -ne 0 -and $null -ne $LASTEXITCODE) {
        throw "install.ps1 exited with code $LASTEXITCODE"
    }
    Pop-Location

    $required = @(
        (Join-Path $TempRoot 'AGENTS.md'),
        (Join-Path $TempRoot '.claude\skills\govern\SKILL.md'),
        (Join-Path $TempRoot '.claude\skills\audit\SKILL.md'),
        (Join-Path $TempRoot '.claude\skills\vision\SKILL.md'),
        (Join-Path $TempRoot '.claude\skills\vision-audit\SKILL.md'),
        (Join-Path $TempRoot '.grok\skills\govern\SKILL.md'),
        (Join-Path $TempRoot '.grok\skills\audit\SKILL.md'),
        (Join-Path $TempRoot '.grok\skills\vision\SKILL.md'),
        (Join-Path $TempRoot '.grok\skills\vision-audit\SKILL.md'),
        (Join-Path $TempRoot '.agents\skills\govern\SKILL.md'),
        (Join-Path $TempRoot '.agents\skills\audit\SKILL.md'),
        (Join-Path $TempRoot '.agents\skills\vision\SKILL.md'),
        (Join-Path $TempRoot '.agents\skills\vision-audit\SKILL.md'),
        (Join-Path $TempRoot '.github\copilot-instructions.md'),
        (Join-Path $TempRoot '.github\prompts\govern.prompt.md'),
        (Join-Path $TempRoot '.github\prompts\audit.prompt.md'),
        (Join-Path $TempRoot '.github\prompts\vision.prompt.md'),
        (Join-Path $TempRoot '.github\prompts\vision-audit.prompt.md'),
        (Join-Path $SkillsDest 'prompts\00-govern-orchestrator.md'),
        (Join-Path $SkillsDest 'prompts\05-independent-audit.md'),
        (Join-Path $SkillsDest 'prompts\06-vision-orchestrator.md'),
        (Join-Path $SkillsDest 'prompts\07-independent-vision-review.md'),
        (Join-Path $SkillsDest 'templates\workspace-context.md'),
        (Join-Path $SkillsDest 'contracts\skills-consumer-contract.schema.json'),
        (Join-Path $SkillsDest 'contracts\skills-consumer-contract.json'),
        # GOAL-019 D-004: core methodology default install (from package core/, not skills dest)
        (Join-Path $TempRoot 'docs\README.md'),
        (Join-Path $TempRoot 'docs\architecture\principles.md'),
        (Join-Path $TempRoot 'docs\architecture\workspace-protocol.md'),
        (Join-Path $TempRoot 'docs\architecture\overview.md'),
        (Join-Path $TempRoot 'docs\architecture\directory-layout.md'),
        (Join-Path $TempRoot 'docs\templates\workspace-context.md'),
        (Join-Path $TempRoot 'docs\templates\goal-folder\00-meta.md'),
        # GOAL-019 phase C: --init-workspace skeleton
        (Join-Path $TempRoot 'docs\workspace-001-pilot-app\workspace.md'),
        (Join-Path $TempRoot 'docs\workspace-001-pilot-app\goal-tree.md')
    )

    $missing = @()
    foreach ($path in $required) {
        if (-not (Test-Path -LiteralPath $path -PathType Leaf)) {
            $missing += $path
        }
    }

    $forbidden = @(
        (Join-Path $TempRoot '.github\prompts\new-goal.prompt.md'),
        (Join-Path $TempRoot 'docs\architecture\tech-stack.md')
    )
    $leaked = @()
    foreach ($path in $forbidden) {
        if (Test-Path -LiteralPath $path -PathType Leaf) {
            $leaked += $path
        }
    }

    $governText = Get-Content -LiteralPath (Join-Path $TempRoot '.claude\skills\govern\SKILL.md') -Raw -Encoding UTF8
    $auditText = Get-Content -LiteralPath (Join-Path $TempRoot '.claude\skills\audit\SKILL.md') -Raw -Encoding UTF8
    $visionText = Get-Content -LiteralPath (Join-Path $TempRoot '.claude\skills\vision\SKILL.md') -Raw -Encoding UTF8
    $visionAuditText = Get-Content -LiteralPath (Join-Path $TempRoot '.claude\skills\vision-audit\SKILL.md') -Raw -Encoding UTF8

    $contentOk = $true
    if ($governText -notmatch '00-govern-orchestrator') {
        Write-Host 'FAIL: Claude govern skill missing 00-govern-orchestrator ref'
        $contentOk = $false
    }
    if ($auditText -notmatch '05-independent-audit') {
        Write-Host 'FAIL: Claude audit skill missing 05-independent-audit ref'
        $contentOk = $false
    }
    if ($visionText -notmatch '06-vision-orchestrator') {
        Write-Host 'FAIL: Claude vision skill missing 06-vision-orchestrator ref'
        $contentOk = $false
    }
    if ($visionAuditText -notmatch '07-independent-vision-review') {
        Write-Host 'FAIL: Claude vision-audit skill missing 07-independent-vision-review ref'
        $contentOk = $false
    }
    if ($governText -notmatch 'workspace-<NNN>-<slug>/workspace\.md') {
        Write-Host 'FAIL: Claude govern skill missing workspace context ref'
        $contentOk = $false
    }
    if ($auditText -notmatch 'workspace-<NNN>-<slug>/workspace\.md') {
        Write-Host 'FAIL: Claude audit skill missing workspace context ref'
        $contentOk = $false
    }

    $codexGovernText = Get-Content -LiteralPath (Join-Path $TempRoot '.agents\skills\govern\SKILL.md') -Raw -Encoding UTF8
    $codexAuditText = Get-Content -LiteralPath (Join-Path $TempRoot '.agents\skills\audit\SKILL.md') -Raw -Encoding UTF8
    $codexVisionText = Get-Content -LiteralPath (Join-Path $TempRoot '.agents\skills\vision\SKILL.md') -Raw -Encoding UTF8
    $codexVisionAuditText = Get-Content -LiteralPath (Join-Path $TempRoot '.agents\skills\vision-audit\SKILL.md') -Raw -Encoding UTF8
    if ($codexGovernText -notmatch '00-govern-orchestrator') {
        Write-Host 'FAIL: Codex govern skill missing 00-govern-orchestrator ref'
        $contentOk = $false
    }
    if ($codexAuditText -notmatch '05-independent-audit') {
        Write-Host 'FAIL: Codex audit skill missing 05-independent-audit ref'
        $contentOk = $false
    }
    if ($codexVisionText -notmatch '06-vision-orchestrator') {
        Write-Host 'FAIL: Codex vision skill missing 06-vision-orchestrator ref'
        $contentOk = $false
    }
    if ($codexVisionAuditText -notmatch '07-independent-vision-review') {
        Write-Host 'FAIL: Codex vision-audit skill missing 07-independent-vision-review ref'
        $contentOk = $false
    }
    if ($codexGovernText -notmatch 'host:\s*codex') {
        Write-Host 'FAIL: Codex govern skill missing host: codex metadata'
        $contentOk = $false
    }

    $wsPath = Join-Path $TempRoot 'docs\workspace-001-pilot-app\workspace.md'
    if (Test-Path -LiteralPath $wsPath -PathType Leaf) {
        $wsText = Get-Content -LiteralPath $wsPath -Raw -Encoding UTF8
        if ($wsText -notmatch 'root_goal:\s*GOAL-001-pilot-vision') {
            Write-Host 'FAIL: workspace.md missing bound root_goal GOAL-001-pilot-vision'
            $contentOk = $false
        }
        if ($wsText -notmatch 'canonical_scope:\s*docs/workspace-001-pilot-app/') {
            Write-Host 'FAIL: workspace.md missing canonical_scope'
            $contentOk = $false
        }
    }
    # Scaffold must NOT create Root five-pack
    if (Test-Path -LiteralPath (Join-Path $TempRoot 'docs\workspace-001-pilot-app\GOAL-001-pilot-vision') -PathType Container) {
        Write-Host 'FAIL: init-workspace must not create GOAL-* five-pack'
        $contentOk = $false
    }

    if ($missing.Count -gt 0) {
        Write-Host 'FAIL: missing required install outputs:'
        $missing | ForEach-Object { Write-Host "  - $_" }
        exit 1
    }
    if ($leaked.Count -gt 0) {
        Write-Host 'FAIL: advanced primitive prompts installed without -WithPrimitives:'
        $leaked | ForEach-Object { Write-Host "  - $_" }
        exit 1
    }
    if (-not $contentOk) {
        exit 1
    }

    Write-Host 'PASS: isolated -All + InitWorkspace produced /govern+/audit+core+workspace skeleton; no GOAL five-pack; no tech-stack.'
    Write-Host "  evidence_dir=$TempRoot"
    exit 0
}
catch {
    Write-Host "FAIL: $($_.Exception.Message)"
    if ((Get-Location).Path -eq $TempRoot) { Pop-Location }
    exit 1
}
finally {
    if (Test-Path -LiteralPath $TempRoot) {
        try {
            Remove-Item -LiteralPath $TempRoot -Recurse -Force -ErrorAction SilentlyContinue
        } catch {
            Write-Host "WARN: could not remove temp dir: $TempRoot"
        }
    }
}
