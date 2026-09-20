---
id: D-001-temporal-binding-and-predicate-scope
doc: decision-entry
status: accepted
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-001 · M3 绑定载体、谓词改造范围与回归载体

## 决定的来源

- `GOAL-001` Root `D-016`（R2 边界与 M1–M4 检查点）、`D-017`（M2/M3 次序 + `I-041-003` 关闭约束）、`D-018`（无过渡期）、`D-019`（F-5）、`D-020`（测试载体）、`D-021`（F-I-005 residual）。
- R1 冻结合同：`r1-c2-per-column-conversion-contract-v1.0-fc.md` §2 的 `read/write` 列 + `r1-c2-predicate-exact-sql-v1.0-fc.md` §1–§5 的 exact new SQL。
- 本轮实测（本目标 A 检查点前）：SQLite 绑 `time.Time` 会写入驱动自己的 36 字符布局 `2026-09-20 12:40:48.814963 +0000 UTC`，**违反 canonical 27 字符合同** → `I-041-005` 因此立项并在本决策定稿。

## 1. 绑定与扫描载体（`I-041-005` 定稿）

**规则：模块与仓储层只用 `time.Time` / `sql.NullTime` / `*time.Time`（域类型），不感知方言；形态转换只发生在 store 适配器。**

| 项 | 规则 | 载体 |
|----|------|------|
| 写（绑定） | 仓储直接把域时间值作为参数绑定；**不得**再手写 `.Unix()` / `.UnixMilli()` | 模块侧；`sqlTx`/`pgTx` 适配器 |
| SQLite 写 | 适配器把 `time.Time` / `*time.Time` / `sql.NullTime` / codec `Value`/`NullValue` 归一为 **canonical fixed-6 UTC 字符串**（nil 保持 NULL） | `internal/store/store.go` `bindSQLiteArgs` |
| PG 写 | 适配器把 `sql.NullTime` / codec `Value`/`NullValue` 归一为 `time.Time`（nil 保持 NULL）；`time.Time` 直通 | `internal/store/postgres.go` `bindPostgresArgs` |
| 读（扫描） | 扫入 `time.Time` / `sql.NullTime` / `*time.Time`；**由 `internal/store/scan.go` 的读适配器解析**（`database/sql` 不会把 string 解析成 `*time.Time`——Go 1.26 实测报 `unsupported Scan, storing driver.Value type string into type *time.Time`）；SQLite 的存量文本必须恰为 canonical fixed-6，否则 fail closed；PG 原生 `timestamptz` 直通 | 模块侧声明域类型；store 适配器 |
| 禁止 | 仓库层绑定 `int64` 秒/毫秒；仓库层按方言分支；把驱动时间类型泄漏进公共契约 | — |
| 无需排序/比较适配 | canonical fixed-6 定宽 → **词法序 = 时刻序**；TEXT 比较与 `ORDER BY` 逐字保留 | — |

- `internal/temporal` 的 `Value`/`NullValue` 保持 **无驱动耦合**（不实现 `driver.Valuer`）；仓储不需要它们，适配器兼容它们以便将来统一。
- **ledger 特例**（M2 已落码）：`schema_migrations.applied_at` 在 v73 之前是 INTEGER 秒、之后是 canonical TEXT / `timestamptz(6)`，故 runner 在**同一事务内**探测列形状后绑定（`sqliteLedgerAppliedAt` / `postgresLedgerAppliedAt`）——这是"无过渡期 + 单批内改形"的必然结果。

## 2. 逐列 transform 分档（必须逐条落实）

| # | 列 | 改造要点 |
|--:|----|----------|
| 5 | `users.locked_until` | 读：已锁定分支 `locked_until IS NOT NULL AND locked_until > ?`；未锁定分支 `locked_until IS NULL OR locked_until <= ?`；写：`time.Time` 或 NULL（清锁写 NULL，不再写 0） |
| 6 | `users.last_login_failure_at` | 读：`(last_login_failure_at IS NULL OR last_login_failure_at < ?)`；写：NULL 或时刻；**删除写 0 路径** |
| 20 | `login_failures.locked_until` | 插入绑 NULL（不再 `0`）；读入 `sql.NullTime` → `Valid && After(now)`；开锁绑 `time.Time` |
| 34 / 78 | `mail_config.updated_at` / `telegram_config.updated_at` | 读：`NULL` = 未初始化/未配置；写：NULL 或时刻（不再 0） |
| 61 | `task_runs.finished_at` | 去掉 `COALESCE(finished_at, 0)`，直接选列扫 `sql.NullTime`；写 NULL（不再 0）；`ORDER BY started_at DESC` 保留 |
| 72 / 73 | `vouchers.expires_at` / `redeemed_at` | 读：`if exp.Valid`（去掉 `> 0`）；写：负值 **fail closed**（Root `D-012`） |
| 4 / 11 | `users.updated_at` / `roles.updated_at` | 单调：`Truncate(µs) now` 与 `old + 1µs` 取大（`D-013`） |
| 其余 | 普通列 | 去 `.Unix()`/`.UnixMilli()`；`NULL` 语义逐字保留 |

## 3. 回归与测试载体

- **SQLite**：全仓 `go test ./...` 必须绿。
- **PG**：至少一条真实路径绿（`PG_TEST_*` 常驻 PostgreSQL 15.4，见 Root `D-017` §3）；**破坏性 migration/round-trip 只能作用于一次性/专用测试 database**（用户约束①）。
- **金额列断言拆分**：`wallet_accounts.balance_total` / `wallet_ledger_entries.amount_delta` 保持 `bigint`；时间列断言改 `timestamp with time zone`（精度 6）；leftover 21 列 + `postgres_test.go:312-316` 缺失七列补入（`D-021` residual ②③）。
- **边界测试重定向**（`D-020`）：`internal/w040contracttest/` 用例改指真实迁移。
- 测试侧若直接以 `int64` 读时间列（非仓储路径），同样改为 `time.Time`/`sql.NullTime`。

## 未选方案

- **在仓储层按方言分支绑定**：违反 `kernel.Tx` 合同（模块不得读 Dialect），且把方言判断散落到 15 个模块。**未采用**。
- **让 codec `Value` 实现 `driver.Valuer` 并在模块侧绑定**：`driver.Valuer` 只能返回 `driver.Value`，无法同时满足 SQLite 要字符串、PG 要 `time.Time`（除非引入方言感知），且会把 codec 与驱动耦合（违反 `D-001` 的包边界）。**未采用**。
- **给每个时间列建视图/触发器做转换**：引入第二套真相，且迁移不可逆点被隐藏。**未采用**。
