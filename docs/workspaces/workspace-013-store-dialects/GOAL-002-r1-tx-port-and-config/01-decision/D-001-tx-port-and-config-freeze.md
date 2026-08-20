---
id: D-001-tx-port-and-config-freeze
doc: decision-entry
status: accepted
created: 2026-08-20
updated: 2026-08-20
parent: GOAL-002-r1-tx-port-and-config
version: 1.0.0
---

# D-001 · 冻结内核 Tx 端口与 db 配置键

- **日期**：2026-08-20
- **状态**：accepted
- **工作区**：`workspace-013-store-dialects`
- **决定**：采纳附件 [r1-tx-port-and-config-freeze.md](../attachments/r1-tx-port-and-config-freeze.md) 为 R1 合同：
  1. 公共持久化面是 `kernel.Store.Run` + `kernel.Tx`（Exec/Query/QueryRow），不是 `*sql.Tx`。
  2. SQL 占位符统一 `?`，由实现 rebind；禁止 `LastInsertId` 与嵌套 `Run`。
  3. 配置键：`db.dialect`（`DB_DIALECT`）、既有 `db.path`（`DB_PATH`）、`db.dsn`（`DB_DSN`）。缺省 sqlite。postgres 必须有 DSN；sqlite 下 DSN 必须为空。不读 `DATABASE_URL`。
  4. 本目标不写代码。驱动选型留 I-002/R2。
  5. 审计模式：**self**（契约冻结、无运行时、可逆文档）。independent 在 R2 实现后做。
- **理由**：现行唯一打开点是 `composition` → `store.OpenWithCatalog(cfg.DBPath)`，唯一事务面是 `WithTx(*sql.Tx)`，配置只有 `db.path`。把方言差收进内核实现、把业务留在模块 Repository，才能满足 VP-013 / RT-P03。`?` + rebind 避免模块 `switch` 方言。不用 `DATABASE_URL` 以免与现有 `DB_*` 双权威。sqlite 下拒绝残留 DSN，避免「以为在跑文件库其实配了 PG」。postgres 下允许残留 `DB_PATH`，避免现网 env 误杀。
- **未选方案**：
  - 内核上帝 Store（GetUser/SaveOrder）：破坏模块 Persistence。
  - 引入 GORM/ent/sqlx 当端口：与 RT-P03 冲突。
  - 模块直接写 `$1` / 方言 SQL 作为公共约定：泄漏方言。
  - `LastInsertId` 留在 Tx：Postgres 无此语义。
  - 嵌套事务/savepoint：R1 不需要，SQLite 单连接下易误用。
  - 默认 `dialect=postgres`：违反内嵌默认。
  - 用 `DATABASE_URL` 兼作 DSN：第二权威。
- **影响范围**：R2 起的 `kernel`/`store`/`config`；R3 迁移 Apply 签名；R4 jobs/模块 `WithTx`。不改 Profile、Manifest、模块矩阵。
- **关联信息项**：GOAL-002 I-001/I-002 → verified。Root I-002（驱动）仍 open，最晚 R2。Root I-003（泄漏清单）本回合开始收集，见 E-001。
- **后续**：R1 自审关门。合同已被 [D-002](D-002-a002-freeze-patch.md) 升为 v1.1、[D-003](D-003-a004-freeze-patch.md) 升为 v1.2、[D-004](D-004-a006-freeze-patch.md) 升为 **v1.3**；R2 按 v1.3 实施，勿按本条原文的 v1.0～v1.2 猜测时间宽度 / `Open` vs catalog。
