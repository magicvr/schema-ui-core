---
id: GOAL-005-requestid-correlation
doc: execution-entry
record_id: E-002
status: recorded
parent: GOAL-001-observability
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

## E-002 · R4 request-id 关联实施与验证（checkpoint）

### 事实

按 GOAL-005 D-001 实施（checkpoint `bc5e196`，3 files +85）：

- `internal/obs/tracing.go`：`serverSpan` 增 `requestID` 参数——非空时 span 打属性 `correlation.request_id`，并 `baggage.NewMember("request-id", …)` + `SetMember` 注入 W3C baggage（失败静默跳过，旁路语义）。
- `internal/obs/observer.go`：`Wrap` 读取 `requestid.FromContext(r.Context())` 传入（requestid 中间件在 mux 外层运行，此处必已就绪；obs → requestid 为叶包依赖，无环）。

### 验证

- 新增 3 组测试（`tracing_test.go`）：真实 `requestid.Middleware` 链（X-Request-ID 头 → span 属性相等）；`serverSpan` 直接断言 baggage `request-id` 键；无 request id 时属性/baggage 均跳过。
- `go vet ./internal/obs/` exit 0；`go build ./...` 通过；全仓 `go test ./...` 无 FAIL。

### Git checkpoint

| hash | scope |
|------|-------|
| `bc5e196` | `apps/api/internal/obs/{tracing.go,observer.go,tracing_test.go}` |

### 备注

- 关联判据（D-001 §1 核对式）已由「X-Request-ID → 属性相等」测试锁定，可直接作为 R5 trace 侧的核对样例。