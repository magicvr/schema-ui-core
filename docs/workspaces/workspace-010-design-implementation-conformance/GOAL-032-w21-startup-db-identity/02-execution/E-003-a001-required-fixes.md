---
id: GOAL-032-w21-startup-db-identity
doc: execution-entry
record_id: E-003
status: recorded
parent: GOAL-001-design-implementation-conformance
created: 2026-08-22
updated: 2026-08-22
version: 0.1.0
---

# E-003 · A-001 required 修复

## 2026-08-22 · 代码与测试

### 已发生事实

1. `lostLedgerLooksComplete` 改为要求 `completeLostLedgerTables`（含 `service_credentials`、`operation_log_session`）。
2. 新增 `identityLostLedgerUnsafe` + `hasPostV1CatalogTables`；plan 为 refuse，reason 含 `identity=lost-ledger-unsafe`。
3. `TestCompleteFingerprintTracksCatalogHead` 锁 catalog 头版本 48。
4. `TestClassifyIdentity` 补：四表无头对象 → unsafe；users-only → partial；完整表但 oursUsers=false → foreign；ledger 即使 oursUsers=false → ours-ledger。
5. `TestPostgresMigrateRefusesIncompleteLostLedger`：全量 boot → drop ledger + drop `service_credentials`/`operation_log_session` → Open 必须 fail 且不建 `schema_migrations`。
6. `TestMigrateFailClosedForeignSQLite`：`orders` 表 → `identity=foreign`，无 ledger。
7. `go test ./internal/store/ -count=1` **ok**（本轮）。

### 证据

| 主张 | 路径 |
|------|------|
| 指纹与 unsafe | `apps/api/internal/store/identity.go` |
| 单测 | `identity_test.go` |
| PG 反例 | `postgres_test.go` `TestPostgresMigrateRefusesIncompleteLostLedger` |
| sqlite foreign | `migrate_test.go` `TestMigrateFailClosedForeignSQLite` |
| 合同修正 | D-002 |

### 计划（非事实）

请 `/audit` 复审 A-001 F-001～F-003 关闭证据后再议 S5 关门。不把 GOAL-032 标 done。
