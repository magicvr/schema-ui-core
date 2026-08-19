---
id: E-010-a008-recommended-remediation
goal: GOAL-001-shared-cross-module-contracts
doc: execution-entry
record_id: E-010
status: recorded
parent: null
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
---

# E-010 · A-008 recommended residual 修复与 runner 终态回归

## 已核对事实

- A-008 F-001：新增 `apps/api/internal/jobs/runner_failure_test.go`，直接验证 `abortLease` 在注入失败原因时清理 lease、写入 `JOB_HANDLER_FAILED` 并触发终态 hook；`go test ./internal/jobs` 通过。
- A-008 F-002：`apps/api/internal/jobs/runner.go` 的 `finish` 在初始取消查询失败、以及 Complete 失败后的二次取消查询失败时均调用 `abortLease`，不再静默依赖 lease 到期回收。
- 原 A-004/A-005 F-001～F-009 修复保持不变；A-008 独立复审确认的 Web/API 定向证据保持有效。

## 验证

- `apps/api`: `go test ./internal/jobs -run TestAbortLeasePersistsFailureAndNotifiesTerminalHook -count=1 -timeout 60s` 通过。
- `apps/api`: `go test ./internal/jobs -count=1 -timeout 60s` 通过。
- `apps/api`: `go test ./internal/handler -run 'TestOperationalGate|TestWriteLocalizedError|TestAuthMiddlewareLocalized' -count=1 -timeout 60s` 通过。
- `git diff --check` 与 tracked/untracked whitespace 检查在本条落盘后复核。

## 审计关联

- A-008 F-001 与 F-002：由本条按 `fixed` 路径闭合；无 residual 或 overrule。
- A-009 为本条对应的 self response；不改变 Root `status` / `progress`。
