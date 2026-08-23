---
id: GOAL-032-w21-startup-db-identity
doc: execution-entry
record_id: E-004
status: recorded
parent: GOAL-001-design-implementation-conformance
created: 2026-08-22
updated: 2026-08-22
version: 0.1.0
---

# E-004 · 响应 A-003 与关门

## 2026-08-22 · 剩余意见处理 + S5

### 已发生事实

1. 用户确认并要求响应 A-003：F-001～F-003 **fixed**，处理其余开放意见后关门。
2. `TestCompleteFingerprintTracksCatalogHead` 改为 `max == completeFingerprintCatalogHead`，并强制 `lockedHeadExtraTables[head]` 表名落在 `completeLostLedgerTables`。
3. `postV1CatalogTables` 扩大到 roles/permissions/menu/operation_log/jobs/credentials/notifications/dict/wallet/recycle/captcha/tasks/mfa 等。
4. 新增 `TestMigrateRestoresLostLedgerSQLite`（OpenSeeded → DROP `schema_migrations` → 再 Open，ledger 行数 = catalog）。
5. 删除 `kernel.ExecIdempotentDDL`。
6. `go test ./internal/store/ ./internal/kernel/ -count=1` **ok**。
7. F-006 按 D-003 **accepted-residual**（双探针）。

### 证据

| 主张 | 路径 |
|------|------|
| 指纹锁表名 | `identity_test.go` `lockedHeadExtraTables` |
| post-v1 扩大 | `identity.go` `postV1CatalogTables` |
| sqlite restore | `migrate_test.go` `TestMigrateRestoresLostLedgerSQLite` |
| 删除死 helper | `kernel/duplicate_object.go` 无 `ExecIdempotentDDL` |
| 合同 | D-003 |
| 关门自审 | A-004 |
