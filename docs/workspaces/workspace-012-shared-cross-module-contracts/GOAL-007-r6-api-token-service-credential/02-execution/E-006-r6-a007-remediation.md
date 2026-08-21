---
id: E-006-r6-a007-remediation
goal: GOAL-007-r6-api-token-service-credential
doc: execution-entry
record_id: E-006
status: recorded
parent: GOAL-007-r6-api-token-service-credential
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
---

# E-006 · R6 A-007 整改与全量回归

## 已完成事实

1. 提交 `b6ebfec` 将 create 201 的一次性字段由 `token` 对齐为冻结契约字段 `secret`，并确认 metadata/list 不返回 `secret`、`token` 或 `tokenHash`。
2. 重名错误的 `fieldErrors.reason` 对齐为 `name already exists`；`scopes` 增加 1～64 个唯一 permission key 的数量上限。
3. `service-credentials.use` 审计 detail 增加 `scopeCount`，继续由 operation 行承载 actor 与 correlation，未记录 raw secret 或 hash。
4. HTTP 黑盒测试逐项覆盖 accounts/profile/avatar/MFA/notifications/wallet 六类 user-only 表面，机器 principal 均为 401；认证测试增加 expired service credential 的 401 用例。
5. `apps/api` 执行 `go test ./...` 全部通过；其中 `internal/handler` 258.366s、`internal/auth` 33.049s、`internal/composition` 33.776s、`internal/authsession` 43.477s、`internal/docscheck` 1.260s。

## 门禁状态

A-007 F-001～F-005 的实现与测试整改已完成；S3 仍等待 independent finding-closure 复核，不在本条执行记录中放行或关门。
