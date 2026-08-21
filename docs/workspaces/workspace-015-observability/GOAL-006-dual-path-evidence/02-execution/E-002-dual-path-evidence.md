---
id: GOAL-006-dual-path-evidence
doc: execution-entry
record_id: E-002
status: recorded
parent: GOAL-001-observability
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

## E-002 · R5 双路径证据收集（live，2026-08-22 07:44–07:45 +08:00）

### 缺省路径证据（VP 退出 3 / Root 成功 3：无收集器仍可开发快测）

命令序列（工作目录 `apps/api`；仅设 `APP_ENV=development`，**不设任何 OBSERVABILITY_\* 变量**）：

```powershell
go build -o .\bin\obs-smoke.exe ./cmd/server
$env:APP_ENV='development'
Start-Process .\bin\obs-smoke.exe -RedirectStandardOutput obs-r5-a-out.log ...
```

实测输出：

| 核对项 | 结果 |
|--------|------|
| `GET /healthz` | **200** |
| `GET /readyz` | **200** |
| `Get-NetTCPConnection -LocalPort 25081 -State Listen`（metrics 默认端口） | **无监听** |
| `Get-NetTCPConnection -LocalPort 4318 -State Listen`（collector 默认端口） | **无监听** |
| 启动日志中 `observability|metrics` 提及数 | **0** |
| 进程停止 | 正常 |

结论：缺省（无显式配置）下进程启动、服务可用、无任何额外端口或日志面——「无收集器仍能开发快测」成立。

### 显式双路径证据（VP 退出 4 / Root 成功 4）

命令序列：

```powershell
# ① 真实 OTLP sink（otlp-sink，R5 工具 `cf9df6c`）：
go run ./cmd/otlp-sink            # 监听 :4318
# ② 启用 metrics + traces 启动 api：
$env:OBSERVABILITY_METRICS_ENABLED='true'
$env:OBSERVABILITY_METRICS_ADDR='127.0.0.1:25099'
$env:OBSERVABILITY_TRACES_ENABLED='true'
$env:OBSERVABILITY_TRACES_ENDPOINT='http://127.0.0.1:4318'
Start-Process .\bin\obs-smoke.exe ...
# ③ 带关联 id 的请求：
Invoke-WebRequest http://127.0.0.1:25080/healthz -Headers @{'X-Request-ID'='r5-evidence-0001'}   # ×2，均 200，响应头回显同 id
# ④ scrape：
Invoke-WebRequest http://127.0.0.1:25099/metrics
```

**指标侧实测**（④返回 200）：

```text
suc_build_info{commit="unknown",go_version="go1.26.0",profile="mvp",version="0.1.0"} 1
suc_http_requests_total{method="GET",module_id="core",route="/healthz",status="200"} 3
suc_kernel_modules_enabled{module_id="admin.users"} 1
```

- `suc_build_info` 与模块 gauge ✓；HTTP 系列带 `module_id` 且 route 为注册 pattern ✓（GOAL-003 判据）。

**trace 侧实测**（约 9 秒后 BSP 批量导出；sink stderr）：

```text
2026/08/22 07:44:47 otlp-sink listening on :4318
2026/08/22 07:45:26 sink: POST / bytes=1037 (total posts=1 bytes=1037)
```

- 真实 OTLP/HTTP protobuf 导出（1037 字节）到达显式配置的 endpoint ✓。
- 关联判据（GOAL-005 锁定的 `correlation.request_id` == `X-Request-ID`）由单测锁定；live 侧请求/响应头回显一致。

**双路径同一次显式运行同时成立**；全部进程已停止、冒烟产物已清理。

### 备注

- sink 的 POST 日志走 stderr（Go `log` 默认），取证时两文件都要看。
- 命令序列可重复（otlp-sink 已入库）；任何协作者/CI 可重跑本 E 条目序列（N-004 固化达成）。