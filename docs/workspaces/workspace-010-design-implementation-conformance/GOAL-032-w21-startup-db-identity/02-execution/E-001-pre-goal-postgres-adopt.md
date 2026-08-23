---
id: GOAL-032-w21-startup-db-identity
doc: execution-entry
record_id: E-001
status: recorded
parent: GOAL-001-design-implementation-conformance
created: 2026-08-22
updated: 2026-08-22
version: 0.1.0
---

# E-001 · 立项前 postgres 启动应急修复

## 2026-08-22 · 立项前已落地（本波治理补记）

用户确认「指定 PG 不应失败」且否定钉死 sqlite 之后、本目标五件套创建之前，已经改过启动路径。本条只记事实，不把后续 Identify/Plan 抽取算进本条。

### 已发生事实

1. `migrateBaselinePG` 不再无条件 `CREATE TABLE`：空库全量建；`users` 像 schema-ui 则只补缺失对象（先探针，禁止在 PG 事务里吞 42P07）；外库 `users` fail closed。
2. postgres runner：打开时记下 preexisting 表；若无 ledger 且像完整安装（`users`+`refresh_tokens`+`operation_log`+`jobs`），v1 之后 stamp 剩余 catalog，避免重放 rebuild。
3. `appliedMigrationsPG` 改查 `current_schema()` 的基表；空 ledger 与 sqlite 一样 fail closed。
4. `backfillRoles` 改为先读完 `users` 再 `INSERT`（pgx 不允许 Rows 未关时同连接 Exec）。
5. `kernel.IsDuplicateObject` / `ExecIdempotentDDL` 落地，并注明 **不能** 用它在已失败的 PG 事务里继续 DDL（25P02）。
6. `dev.cmd`：**撤回** `DB_DIALECT=sqlite` 钉死；overlay 不再写死 sqlite。保留 API 崩溃 `pause` 与 `%TEMP%` 反斜杠 / `addr` 引号修复。
7. 活 PG（`configs/.env` 的 postgres）`go run ./cmd/server`：`/healthz` 与 `/readyz` HTTP 200。既有用户未重种，`admin/admin` 登录 401 是数据态，不是启动失败。

### 证据

| 主张 | 路径 / 命令 |
|------|-------------|
| PG 基线 adopt | `apps/api/internal/modules/authsession/migration/migration.go` `migrateBaselinePG` |
| 完整旧库 stamp | `apps/api/internal/store/postgres.go` `lostLedgerLooksComplete` / `stampRemainingPG`（E-002 起由 Plan 驱动） |
| backfill 安全 | 同 migration.go `backfillRoles` |
| 测试 | `TestPostgresMigrateAdoptsLedgerlessR2`、`TestPostgresMigrateRejectsConflictingUsersTable`、`TestPostgresMigrateRestoresLostLedger` |
| 启动活测 | `go run ./cmd/server` + `curl /healthz` `/readyz` → 200 |
