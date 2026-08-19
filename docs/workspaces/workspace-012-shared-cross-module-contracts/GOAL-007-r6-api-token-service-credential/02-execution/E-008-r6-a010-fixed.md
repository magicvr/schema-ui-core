---
id: E-008-r6-a010-fixed
goal: GOAL-007-r6-api-token-service-credential
doc: execution-entry
record_id: E-008
status: recorded
parent: null
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
---

# E-008 · A-010 F-010 fixed 实施

## 已完成事实

- `apps/api/internal/auth/auth.go` 不再丢弃 service-credential 使用审计或 `last_used_at` 更新错误；任一错误均返回 503 `STORAGE_UNAVAILABLE`，并跳过 downstream handler。
- 生产 composition 改用 `MarkServiceCredentialUsedWithAudit` 调用方事务 seam，使用审计与 `last_used_at` 更新原子提交；审计失败会回滚元数据更新。兼容旧 seam 仍 fail closed。
- `apps/api/internal/auth/auth_test.go` 新增 recorder 故障与 metadata 故障回归，均核对 503、downstream 未调用及失败审计路径的状态。

## 验证

- `go test ./internal/auth -count=1` 通过。
- `go test ./internal/composition -count=1` 通过。
- `go test ./internal/modules/authsession -count=1` 通过。
- `go test ./internal/handler -count=1` 通过（264.787s）。
- `go test ./... -p 1 -count=1 -timeout 600s` 通过（API 全包串行回归）。

## 关联

- 响应 workspace-012 Root A-010 F-010；决策见 R6 D-004。
- 本条不冒充 independent 复审；A-010 原 verdict 与 finding 原文保留。
