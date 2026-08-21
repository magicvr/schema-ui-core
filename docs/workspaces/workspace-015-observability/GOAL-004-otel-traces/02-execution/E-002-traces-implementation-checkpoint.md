---
id: GOAL-004-otel-traces
doc: execution-entry
record_id: E-002
status: recorded
parent: GOAL-001-observability
created: 2026-08-21
updated: 2026-08-22
version: 1.0.0
---

## E-002 · R3 OTel traces 接入实施与验证（checkpoint）

### 事实

按 GOAL-004 D-001 实施全切片（checkpoint `2ab4ec4`，11 files +652）：

- **`internal/config`**：`Config` 新增 `TracesEnabled` / `TracesEndpoint` / `TracesSampleRatio`（默认 1.0）；`yamlFile` 增 `observability.traces.*`；env `OBSERVABILITY_TRACES_ENABLED|ENDPOINT|SAMPLE_RATIO`；`validateTraces` fail-closed（enabled 无 endpoint 拒绝 / 非 http(s) URL 拒绝 / disabled 死 endpoint 拒绝 / ratio 出 (0,1] 拒绝——经 Load + ValidateProd 双侧）。两份 YAML 补 `traces:` 段。
- **`internal/obs`（新 tracing.go）**：`Tracing` 类型——disabled 纯 no-op（无 provider/exporter/后台 goroutine）；enabled 走 `otlptracehttp` + BatchSpanProcessor + `ParentBased(TraceIDRatioBased)`；Resource = service.name/version + deployment.environment；`otel.SetErrorHandler` → slog warn；`serverSpan`（W3C TraceContext+Baggage 提取、span name = `<METHOD> <route>`、attrs method/route/url.path 无 query）+ `finishServerSpan`（status 属性、≥500 → codes.Error）。
- **`observer.go`**：`SetTracing` + `Wrap` 扩展——同一拦截点同时产出 metrics 与（启用时）server span。
- **`internal/composition`**：`newTracing` provider、`newObserver` 注入 tracing、`registerLifecycle` OnStop join `tracing.Shutdown(ctx)`（flush 后退出）。
- **依赖**：`go.opentelemetry.io/otel` v1.45.0 + `otel/sdk` + `otlptracehttp`（go.mod/go.sum）。

### 验证

- 新增测试：config 10 子测试（缺省 no-op / YAML 启用 / 缺 endpoint / 非 http 协议 / 死 endpoint / ratio 越界 / env 覆盖 / 零值 ValidateProd / 手工 Config 复检 / 无效 env 保默认）；obs 4 组（disabled no-op / span 形状与 503→Error / traceparent join / **httptest OTLP sink 实收 POST**）；composition 2 组（wiring no-op 缺省 / enabled provider）。
- `go vet`（obs/composition/config）exit 0；`go build ./...` 通过；全仓 `go test ./...` 无 FAIL。
- **live 冒烟（负路径）**：真实二进制 `OBSERVABILITY_TRACES_ENABLED=true` + endpoint `http://127.0.0.1:1`（不可达）→ 启动成功、`/healthz` 200、metrics 200；连续请求后 stdout 出现 2 次 `otel export issue` WARN（BSP 周期重试，`traces export: Post ...connection refused`）——**证明 span 已创建并异步导出、失败仅告警绝不致命**；进程正常停止。冒烟产物已清理。

### Git checkpoint

| hash | scope |
|------|-------|
| `2ab4ec4` | apps/api：config（三键+校验+测试）、internal/obs（tracing.go/tracing_test.go/observer 扩展）、composition（newTracing 接线 + 测试）、go.mod/go.sum、两份 YAML |