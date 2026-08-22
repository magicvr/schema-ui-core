---
id: GOAL-032-w21-startup-db-identity
doc: execution-entry
record_id: E-002
status: recorded
parent: GOAL-001-design-implementation-conformance
created: 2026-08-22
updated: 2026-08-22
version: 0.1.0
---

# E-002 · Identify / Plan 接入双方言 migrate

## 2026-08-22 · S1–S4 实施

### 已发生事实

1. D-001 冻结：沿用 `schema_migrations` 当 EF 式历史表；无 ledger 时做身份探针，不盲 Migrate。
2. 新增 `apps/api/internal/store/identity.go`：`classifyIdentity` / `planStartup`；sqlite/postgres `probeIdentity`；`restoreLedger` 用 `CREATE TABLE IF NOT EXISTS schema_migrations` 后按 catalog 整表盖章（不重放 Apply）。
3. sqlite `Store.migrate` 与 postgres `migrate` 均按计划执行：`refuse` / `noop` / `apply-pending` / `fresh` / `adopt-r2` / `adopt-then-pending` / `restore-ledger`。
4. postgres 内联 `stampRemainingPG` 特判删除，改由 `actionRestoreLedger` 驱动。

### 证据

| 主张 | 路径 |
|------|------|
| 合同 | `01-decision/D-001-identity-and-plan-freeze.md` |
| 代码 | `apps/api/internal/store/identity.go`、`migrate.go`、`postgres.go`、`open.go` |
| 单测 | `apps/api/internal/store/identity_test.go` |
| 集成 | 既有 `TestPostgresMigrateAdoptsLedgerlessR2` / `RejectsConflictingUsersTable` / `RestoresLostLedger` |
