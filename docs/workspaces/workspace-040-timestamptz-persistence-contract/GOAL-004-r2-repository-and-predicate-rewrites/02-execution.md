---
id: GOAL-004-r2-repository-and-predicate-rewrites
doc: execution
status: active
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# 执行记录 · GOAL-004（R2 M3）

> 只记**事实**（做了什么、产物路径、进度评估）。计划单独标注。

## E-001 · 子目标立项与 M3 载体冻结

- **来源**：Root `D-017`（用户 2026-09-20 P-004 裁决）：M2 落码后按 `D-018` 无过渡期使全仓测试红，故 **M2 不单独提交**，M3 转绿后合并提交；用户确认 slug `GOAL-004-r2-repository-and-predicate-rewrites`。
- **产物**：本目标五件套 + `01-decision/D-001-temporal-binding-and-predicate-scope.md`（`I-041-005` 定稿：域类型 `time.Time`/`sql.NullTime`，形态转换只在 store 适配器；逐列 transform 分档；回归载体）。
- **同时落盘**：Root `D-017`；`D-016` §5 `I-041-003` → `verified`（常驻 PostgreSQL 15.4 实测执行 + 三条用户约束）。
- **事实（已完成的依赖改造，M2 范围内）**：
  - `internal/store`：SQLite `sqlTx` / PG `pgTx` 增加**域时间参数归一化**（canonical 字符串 / `time.Time`），使仓储层无需按方言分支（`D-001` §1）。
  - `internal/store`：ledger 写入改为**同事务探测列形状**（`sqliteLedgerAppliedAt` / `postgresLedgerAppliedAt`）。**实测发现**：M2 若按冻结附件 §6 字面统一绑 `time.Time`/canonical 串，会让 v1–v72 期间写入 INTEGER 列的 ledger 行落成文本，v73 重建时 `strftime` 解析失败 → `NOT NULL constraint failed: schema_migrations.applied_at`。该偏差已记录，待 independent 复审。
- **进度评估**：本目标 A 检查点开始；`progress: 0/3`（A/B/C 均未完成）。
- **未提交状态**：M2 + M3 均在工作树中；按 `D-017` §1，转绿后一并提交。

## E-002 · M3 落码（读写/谓词/边界测试）

- **D-001 §1 的读侧更正（实测）**：`database/sql` **不会**把 string 解析成 `*time.Time`（Go 1.26 `convert.go` 实测报 `unsupported Scan, storing driver.Value type string into type *time.Time`）。原 D-001 §1 表格的「canonical TEXT 直接解析」写法**错误**；已改为由 `internal/store/scan.go` 承担读侧归一化，并同步修订 D-001。
- **store 适配器（本目标 A 检查点产物）**：
  - `internal/store/scan.go`：`*time.Time` / `sql.NullTime` / `**time.Time` 目标经临时值转换；SQLite 文本必须与 canonical fixed-6 往返一致，否则 fail closed；其余目标直通。
  - `internal/store/store.go` / `postgres.go`：`sqlTx`/`pgTx` 的 `Query`/`QueryRow` 返回包装结果；写侧 `bindSQLiteArgs`（→ canonical 串）/ `bindPostgresArgs`（→ `time.Time`）。
  - `internal/store/scan_test.go`：读适配器单测（canonical、NULL、非规范文本拒绝、直通）。
- **ledger 写入的 PG 预编译语句类型问题（实测）**：v73 之前 ledger 行绑 int64、之后绑 `time.Time`，而 pgx 按 SQL 文本缓存 prepared statement → v73 之后仍按 `bigint` 解析参数，报 `22P02 invalid input syntax for type bigint: "2026-09-20 12:57:15.914522 +0000 UTC"`。修复：`ledgerWrite` 按探测到的形状给出**不同文本**（legacy `CAST(? AS bigint)` / canonical `CAST(? AS timestamptz)`）；SQLite 侧共用一条文本（动态类型）。
- **模块侧改造（4 个并行子代理，按 D-001 逐列分档）**：
  - `modules/authsession` + `systemdata`：`User.LockedUntil` → `sql.NullTime`；#5/#6/#20 sentinel 谓词与写 NULL；#4/#11 微秒单调；其余列域类型化。新增 `temporal_sentinel_test.go`。
  - `internal/jobs`、`internal/mail`、`internal/channel/telegram`、`modules/channel/telegram`：ms 族（jobs/operation_log/mail_outbox）绑瞬时；#34/#78 `NULL = 未初始化`；删除 `toMillis`。
  - `modules/recyclebin`、`scheduledtasks`、`notifications`、`datadictionary`、`datapermission`、`logincaptcha`、`mfa`、`settings`：#61 去 `COALESCE(finished_at,0)`；`site_settings` 运行时 upsert 绑瞬时（v7 seeder 保持 Unix 秒，见下）。
  - `modules/wallet`、`digitaloffer`、`operationlog`：`operation_log.created_at` 毫秒族绑瞬时；#72/#73 `Valid` 判存在 + 负值 fail closed（写路径同时守卫 `redeemed_at`）；钱包账本 `ORDER BY created_at, id` 逐字保留。
- **边界测试重定向（`D-020` §2，本目标 C 检查点产物）**：新增 `apps/api/internal/w040contracttest/migration_boundaries_test.go` —— 同一批边界值改为**跑真实 v73–v87 迁移**：v1–v72 建库 → 以 legacy 整数形状播种 → 应用到 head → 用 store 读适配器断言 sentinel 0→NULL（#5/#6/#34/#78/#57/#61/#72/#73）、999 ms 无进位、负毫秒 floor（-1/-999/-1000/-1001）、公元 9999、fixed-6 词法序 = 时刻序、voucher 负值 fail closed（v85 事务回滚 + 列仍为 INTEGER）、普通负值为合法 instant。原 `contract_boundaries_test.go` 的冻结文本用例保留为文本层记录。
- **测试侧修正**（时间列由测试直接读写者）：`migrate_telegram_upgrade_test.go`（改走 `st.Run` 适配器 + 瞬时比较）、`postgres_telegram_test.go`（`lastMessageAt` 瞬时）、`postgres_test.go`（PG seed 由 `now.Unix()` 改绑 `now`）。
- **store 断言拆分（`D-021` residual ②③）**：`postgres_test.go` 时间列断言改 `timestamp with time zone`；`wallet_accounts.balance_total` / `wallet_ledger_entries.amount_delta` 单列断言保持 `bigint`；leftover 名单补齐至 **21 名**（新增 `last_login_failure_at`/`sent_at`/`consumed_at`/`last_sent_at`/`redeemed_at`/`last_message_at`/`received_at`）并新增 `datetime_precision = 6` 断言。
- **偏差（须 independent 复审，见 `03-audit.md`）**：
  1. `modules/settings/migration/migration.go` 的 v7 seeder **保持绑定 Unix 秒**（冻结附件 §2.12 要求「与 v84 同批改目标格式」）：v7 在 v84 之前执行，此时列仍是 `INTEGER`；若改为 canonical 串，v84 的 `strftime(updated_at,'unixepoch')` 会得 NULL → `NOT NULL constraint failed`（与 ledger 实测同型错误）。
  2. PG ledger 写入改为两条**显式 CAST** 语句（见上）——属「Go 控制流/绑定」层，不进 checksum。
- **证据（本目标 A/B 检查点）**：
  - `go build ./...` exit 0；`go vet ./internal/store/ ./internal/temporal/ ./internal/temporalmigrate/ ./internal/w040contracttest/` 无输出。
  - `go test ./internal/w040contracttest/` 全绿（真实迁移边界矩阵，含 PG 侧不变量）。**注意**：该包当前只覆盖 SQLite 侧真实迁移；PG 侧由 `internal/store/postgres_test.go` 的目录/精度断言承担。
  - `go test ./internal/store/ -run TestFullCatalogPostgresBootstrapIntegration` **ok**（真实 PostgreSQL 15.4 上 v1–v87 全量应用）。
  - 各模块包分别报 `ok`（4 个子代理的收口记录见其报告）。
