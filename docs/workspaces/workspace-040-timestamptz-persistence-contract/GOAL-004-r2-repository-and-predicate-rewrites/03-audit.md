---
id: GOAL-004-r2-repository-and-predicate-rewrites
doc: audit
status: active
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# 审计台账 · GOAL-004（R2 M3）

> 本文件是唯一正式审计台账索引：`self` 与 `independent` **共用** `A-NNN` 序列。
> 每条意见的正文在 `03-audit/A-NNN-<slug>.md`；本文件登记条目头（`source` / 日期 / scope / `verdict`）。
> 独立审计默认只写意见，不修改 `status` / `progress` / 方案正文；响应归编排器。

## 意见索引

| A-ID | source | 日期 | scope | verdict | 摘要 | 文件 |
|------|--------|------|-------|---------|------|------|
| — | — | — | — | — | 尚无意见（本目标 A 检查点实施中） | — |

## 待复审事项（编排器登记，供独立审计取证）

| # | 事项 | 证据位置 | 说明 |
|--:|------|----------|------|
| 1 | ledger 写入的**列形状探测**偏离 M2 冻结附件 §6 字面（改为同事务探测） | `internal/store/identity.go`、`internal/store/migrate.go`、`internal/store/postgres.go` | 实测原字面写法会致 `NOT NULL constraint failed: schema_migrations.applied_at`；是否接受由 independent 判定 |
| 2 | v86 descriptor 的 `ModuleID` 由台账的 `admin.channel.telegram` 更正为模块实际的 `channel.telegram` | `modules/channel/telegram/migration/vp040_temporal.go`、`attachments/r2-v73-v87-generated-statements-v0.1.md` | 模块 ID 由 provider 机械决定；台账 §1 的 `ModuleID` 列不在 `D-014`/台账 §4 点名的「唯一允许输入」清单内 |
| 3 | v73–v87 语句由生成器从 **v72 live `sqlite_master`** 机械导出（含 m0/m4 查询语句进入 checksum 输入） | `attachments/r2-v73-v87-generated-statements-v0.1.md` | 逐列类型编辑的机械性与 m0/m4 进入哈希是否符合 `D-017` 由 independent 判定 |
| 4 | **PG ledger 写入使用两条显式 CAST 语句**（legacy `CAST(? AS bigint)` / canonical `CAST(? AS timestamptz)`），偏离「同一语句文本」的直觉写法 | `internal/store/identity.go` `postgresLedgerWrite` | 实测 pgx 按 SQL 文本缓存 prepared statement，跨 v73 复用同一文本会把参数类型钉死为 `bigint`（`22P02`）；是否接受由 independent 判定 |
| 5 | `modules/settings/migration/migration.go` 的 v7 seeder **保持 Unix 秒**，偏离冻结附件 §2.12「与 v84 同批改目标格式」 | 该文件 `migrate0007` / `migrate0007PG` | v7 在 v84 之前执行、列仍是 INTEGER；改绑 canonical 串会使 v84 的 `strftime` 得 NULL → `NOT NULL constraint failed`。是否接受该偏离由 independent 判定 |
| 6 | 读侧合同（`database/sql` 不解析 string → `*time.Time`）由 `internal/store/scan.go` 承担，D-001 §1 已同步更正 | `internal/store/scan.go`、`internal/store/scan_test.go` | 原 D-001 §1 表格写法经 Go 1.26 实测证伪；属**更正**而非新增设计，须记录 |
| 7 | 失败即部分升级：descriptor 各自事务，故 v85 失败时 v73–v84 已提交（非整批回滚） | `internal/store/migrate.go` `applyPending`；`internal/w040contracttest/migration_boundaries_test.go` | C3 恢复边界已覆盖；是否需要在 M4 前补「批级回滚/快照」由 `D-021` 复审一并判定 |
| 8 | 4 个子代理并行改造的一致性 | 各模块 `02-execution` 事实 + `git diff` | 判据：无方言分支、无驱动类型泄漏、`ORDER BY`/谓词逐字保留、无断言弱化 |
