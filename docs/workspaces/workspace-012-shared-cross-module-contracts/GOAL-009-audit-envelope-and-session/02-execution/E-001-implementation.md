---
id: E-001-implementation
goal: GOAL-009-audit-envelope-and-session
doc: execution-entry
record_id: E-001
status: recorded
parent: null
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
---

# E-001 · JWT sid、session 侧表、writer envelope

1. access JWT 增加 `sid`（refresh token id）；middleware / service-credential / dev-session 写入 `User.SessionID`。
2. 迁移 0048：`operation_log_session` + `operation_log_archive_session`；retention archive 同步拷贝 session。
3. 生产 mutation writer 改 `NewDetail` + `recordAudit`（含 users_state/MFA/wallet/roles/import/export/datapermission/captcha/dictionary/account）。字典类型级联删除改为 `auditDetail("delete", {entries})`，不再手写 JSON。
4. 读路径：`/api/operations` 与 CSV 导出增加 `sessionId`（与 `correlationId` 同形）。后台 Job 无登录会话，不写 session。
5. 测试：`go test ./internal/modules/operationlog ./internal/auth ./internal/store ./internal/handler ./internal/composition ./internal/modules/wallet -count=1` 通过。覆盖 0048 加法迁移、归档 session、auth/users/settings/service-credential session、envelope `ParseDetail`。
