---
id: GOAL-003-metrics-scrape-endpoint
doc: execution-entry
record_id: E-002
status: recorded
parent: GOAL-001-observability
created: 2026-08-21
updated: 2026-08-21
version: 1.0.0
---

## E-002 · R2 scrape 接入实施与验证（checkpoint）

### 事实

按 GOAL-003 D-001 实施全切片（checkpoint `5ba04c5`，22 files +921）：

- **`internal/obs`（新包）**：
  - `observer.go`：私有 `prometheus.Registry` + 固定系列（`suc_build_info` / `suc_http_requests_total` / `suc_http_request_duration_seconds` DefBuckets / `suc_kernel_modules_enabled` + Go/process collectors）；`Wrap`（route/method 取注册 pattern，status 记录显式 WriteHeader、隐式 200）；`Handler(token)`（`subtle.ConstantTimeCompare` Bearer 守卫）；全部方法 nil-receiver 安全。
  - `instrumented_mux.go`：`InstrumentedMux` 内嵌 `*http.ServeMux`，覆写 `Handle`/`HandleFunc`；`Own(pattern, moduleID)` 所有权声明，缺省 `core`。
  - `server.go`：可选专用 listener；disabled 全 no-op；bind 失败 fail-closed；Serve 运行期错误仅 error 日志不杀进程。
- **`internal/handler`**：新增 `registrar.go`（`routeRegistrar` 接口 = Handle + HandleFunc）；9 个装配函数（`RegisterWithMFAProbes`、`authsHandler`、`accountsHandler`、`RegisterUpload`、`RegisterPublicBranding(Assets)`、`RegisterSchemas`、`RegisterManifest`、`RegisterBootstrapWithAvailability`）参数从 `*http.ServeMux` 放宽为接口——`*http.ServeMux` 与 `*obs.InstrumentedMux` 均满足，调用方零破坏。
- **`internal/composition`**：`newObserver`（build info + plan.IDs 逐模块 gauge）、`newMetricsServer`（Config→obs.Server 映射）、`newMuxWithExtraProviders` 改用 InstrumentedMux 并对 contributed routes 调 `Own(route.ModuleID)`、`registerLifecycle` OnStart 挂 metrics listener（失败回滚链完整）/OnStop join 停机；3 处既有测试调用点补 `nil` observer。
- **依赖**：`github.com/prometheus/client_golang v1.23.2`（go.mod/go.sum）。

### 验证

- `go vet ./internal/obs/ ./internal/composition/ ./internal/handler/` → exit 0。
- `go build ./...` → 通过；`go test ./...`（apps/api 全仓）→ 无 FAIL。
- 新增测试：obs 包 7 组（系列暴露/标签取 pattern 且 raw path 不泄露/隐式 200/nil 透传/Bearer 401+200/splitPattern/disabled no-op + 真实 listener scrape + bind 冲突 fail-closed + HandleFunc 拦截）；composition 3 组（probe 路由 `module_id="admin.probe"` + health `core` + 模块 gauge + build info / enabled 端到端 scrape / disabled 惰性）。
- **live 冒烟**：`go build` 真实二进制 + `OBSERVABILITY_METRICS_ENABLED=true`（127.0.0.1:25099）启动 → `/healthz` 200 → `GET /metrics` 200（10306 字节），实测含 `suc_build_info{profile="mvp"...}`、`suc_http_requests_total{module_id="core",route="/healthz",status="200"} 1`、`suc_kernel_modules_enabled{module_id="admin.users"}` 等全部系列；停机后端口释放。冒烟产物（exe/log/pid）已清理，不入库。

### Git checkpoint

| hash | scope |
|------|-------|
| `5ba04c5` | apps/api：go.mod/go.sum、internal/obs（新包 4 文件）、internal/handler（registrar + 9 签名放宽）、internal/composition（接线 + 4 测试文件） |
