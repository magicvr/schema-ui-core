---
doc_type: goal-execution
id: E-002-contract-frozen
parent: GOAL-002-r1-contract-freeze
date: 2026-09-03
status: active
version: 0.1.0
---

# E-002 · C2 合同冻结与端口落地

## 事实时间线

- 2026-09-03：D-002 v0.1.0 冻结 AllowRecord 单锁等价、兼容、14 处分母、R1/R2 切分。
- 2026-09-03：`apps/api/kernel/ratelimit.go` 接口新增 `AllowRecord`；`RetryAfterSeconds` godoc 允许在 AllowRecord false 后调用。
- 2026-09-03：`apps/api/kernel/ratelimit_test.go` stub 补 `AllowRecord`。
- 2026-09-03：`apps/api/internal/ratelimit/memory.go` 抽 `allowLocked`/`recordLocked`；`AllowRecord` 同一把锁内 Allow-then-Record。
- 2026-09-03：`memory_test.go` 增顺序等价、拒绝不登记、并发预算；既有并发测试混入 AllowRecord。
- 2026-09-03：验证（`apps/api`）：`go test ./kernel/ ./internal/ratelimit/` 绿；`go test -race ./internal/ratelimit/` 绿；`go build ./...` 通过。生产 14 处调用点未改。
- 2026-09-03：Git checkpoint `bdfe925f`（owned paths：vision VP-032 激活台账 + workspace-032 + kernel/Memory AllowRecord）。

## 产物

- `01-decision/D-002-allowrecord-port-contract.md`
- `apps/api/kernel/ratelimit.go`
- `apps/api/kernel/ratelimit_test.go`
- `apps/api/internal/ratelimit/memory.go`
- `apps/api/internal/ratelimit/memory_test.go`

## 下一步（计划）

- C3：R1 阶段关门 self（A-001）；通过后 Root progress → 1/3，立项 R2 迁移。
