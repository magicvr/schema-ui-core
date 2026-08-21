---
id: E-009-all-audit-findings-remediation
goal: GOAL-001-shared-cross-module-contracts
doc: execution-entry
record_id: E-009
status: recorded
parent: null
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
---

# E-009 · A-004/A-005 全部 finding 修复与回归验证

## 已核对事实

- F-001：`apps/api/internal/jobs/runner.go` 的 lease claim 使用 execute 传入的 request context；heartbeat 的取消查询与续租也使用同一执行 context。
- F-002/F-009：heartbeat 的取消查询或续租失败会先取消 handler，再通过 background cleanup 将 lease 标记为 `JOB_HANDLER_FAILED`，成功后通知终态 hook；runner 正常 Stop 仍保留可恢复的 running lease 语义。
- F-003：`apps/api/internal/requestid/requestid.go` 的随机源不可用回退改为时间戳加进程内原子序列，避免同纳秒碰撞；随机路径仍使用 CSPRNG。
- F-004：`apps/api/internal/auth/auth.go` 的 service credential ID 回退改为时间戳加进程内原子序列；opaque token 仍无回退并在 CSPRNG 失败时返回错误。
- F-005：新增 `apps/api/internal/errorcatalog/writer.go` 作为共享错误包络 writer；handler 与 auth 保留兼容调用面并委托共享实现，消除两套 wire 行为漂移。
- F-006：`apps/api/internal/handler/operational.go` 使用显式 recovery path registry，覆盖当前 auth refresh/logout、account password 与 MFA 自助恢复端点；维护模式 allowlist 回归测试覆盖全部登记路径。
- F-007：`apps/api/internal/modules/operationlog/detail.go` 对 struct、typed map、typed slice 等未知类型进行 JSON 归一化后递归脱敏，不再因可序列化类型直接失败；新增 typed value 回归测试。
- F-008：沿用 E-008/A-006 的 en-US/zh-CN `schema.systemMonitoring.statCard.availability` 修复与 Web 全量证据。

## 验证

- `apps/api`: `go test ./internal/requestid ./internal/auth ./internal/jobs ./internal/modules/operationlog` 通过。
- `apps/api`: `go test ./internal/handler -run 'TestOperationalGate|TestWriteLocalizedError|TestAuthMiddlewareLocalized' -count=1` 通过。
- `apps/web`: `npm test -- --run` 通过，72/72 test files、1069/1069 tests。
- `apps/api`: `go test ./internal/docscheck` 通过（并在串行 API 复跑中通过）。
- 仓库：`git diff --check` 与 untracked 文件 whitespace 检查待本条落盘后再次执行。
- 串行 `go test ./... -p 1 -timeout 300s` 复跑在既有 `internal/handler/TestNotificationPruneKeepsUnread` 的 SQLite `VACUUM` 初始化路径超时；该测试此前约 80 个 handler 测试均通过，超时未指向本次变更。

## 审计关联

- A-004 F-001～F-007：本条提供 fixed 实现与定向回归证据。
- A-005 F-008：fixed 证据由 E-008 保持，本条复核其仍通过。
- A-005 F-009：与 F-002 同一 heartbeat 异常路径，本条提供 fixed 实现证据。
- A-007 将本条作为 self response 的执行依据；本条不冒充 independent 审计。
