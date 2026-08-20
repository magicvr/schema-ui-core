---
id: D-001
doc: decision-entry
goal: GOAL-004-r3-dual-dialect-ledger
status: accepted
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# D-001 · R3 方案：Apply 形状、catalog 形态、checksum 绑定与对写规则

## 触发

- Root R2（GOAL-003）已 done：`store.Open(ctx, OpenOptions, catalog)` + `kernel.Store`/`kernel.Tx` 可用，postgres live 探测通过。R3 依赖项满足。
- R1 合同 v1.4 §3/§4/§6 已冻结 R3 的硬规则（时间 BIGINT、checksum 绑 sqlite 历史、`INSERT OR IGNORE`/`sqlite_master`/`PRAGMA`/`COLLATE NOCASE` 点名债、非时间 INTEGER 宽度按列）。

## 决定

1. **迁移贡献形状**（v1.4 §4）：`kernel.MigrationContribution.Apply` / `Reconcile` 由 `func(*sql.Tx) error` 改为 **`func(kernel.Tx) error`**；T1 一次性切全部模块迁移签名。store 迁移运行器（`migrate.go`）改用 `kernel.Store.Run` 驱动 apply + 台账插入；sqlite `WithTx` 保持仅给 R4 前的模块运行时用。
2. **catalog 形态（I-003）**：采用 **按方言分列的 compiled catalog**——同一版本号集合，`sqlite` 与 `postgres` 各持一份物理 SQL（版本对齐）；模块迁移贡献提供两个 Apply（或一个便携 Apply + 方言分文件）。checksum 一律绑 **sqlite/canonical 历史文本 + transform id**（v1.4 §4）；postgres 成对 SQL 不进 digest。禁止改 sqlite 历史文本对齐 PG。
3. **对写范围（本目标内）**：只处理**迁移 Apply 路径**。
   - authsession `migration.go`：`sqlite_master` 探测 / `PRAGMA table_info|foreign_key_list` 改为 PG 等价（`information_schema` / `pg_catalog`）；`COLLATE NOCASE` DDL 在 PG 上用 `CITEXT` 或 `LOWER()` 唯一索引（显式选并落盘）；时间列 PG `BIGINT`。
   - 时间列：全部 Unix 秒/毫秒列 PG `BIGINT`（v1.3 硬规则）；新迁移明确单位。
   - 非时间 INTEGER：按列决定 `INTEGER`/`BIGINT`（wallet `balance_*`、流水/归档 `amount_delta`/`balance_after_*` 等）；布尔列选两方言等价形态并落盘。
   - `?` 可 rebind 文本保持统一 `?`；不可 rebind 的方言差异成对文件。
   - **不处理**模块运行时 SQL（operationlog `INSERT OR IGNORE`、wallet/recyclebin `LIKE`、users/roles `COLLATE NOCASE` ORDER BY、布尔运行时读写）——全部归 **R4**（Root R4 收口）。
4. **postgres 迁移运行器（T2）**：postgres Store 上实现一方言 apply：`WasFresh`（已实现）→ 空库执行 postgres catalog 全部迁移 → `schema_migrations` 台账（version/name/checksum 与 sqlite 同值，checksum 绑 sqlite 文本）。`Run` 内一个迁移一个事务；非空 catalog 不再 fail-closed（R3 起 postgres 允许 apply 双方言 catalog；R2 的「非空 fail-closed」解除——需同步改 R2 的 `openPostgres` 守卫与 `composition.openStore` 路由）。
5. **composition 路由（T2 伴随）**：postgres + catalog 走 apply 路径（不再 fail-closed）；`/readyz` 仍以 Open+Ping+模块就绪为准，模块门禁在 Reconcile 后。
6. **审计模式**：迁移/数据门禁 → **self + independent**（项目默认 grok build），T3 对写完成且双路径证据后关门。
7. **信息门禁**：I-001/I-002 在每迁移对写前闭合（逐迁移核对清单 + 逐列证据）；I-003 由本 D-001 裁决为「分列 catalog」；I-004（T4 前）live 对比。

## 为什么

- 形状转 `kernel.Tx` 是 v1.4 冻结目标；R2 已把 `kernel.Tx`/`Run` 落地，T1 只改贡献签名与运行器，改动集中在迁移包。
- 分列 catalog 比「单一 catalog 内嵌两套 SQL」更干净地满足 checksum 绑定（sqlite 历史不动、PG 各自成对），且 v1.4 §4 明确允许。
- 运行时 SQL 债划归 R4：避免 R3 一次改编译迁移 + 全部模块 Repository 的爆炸范围；Root 路线图亦把「模块仓库公共面收口」列 R4。
- R2 的 postgres 非空 catalog fail-closed 是「对写完成前」的临时闸；R3 开始具备双写能力后必须解闸，否则 postgres 永远无法 apply。

## 未选方案

- **单一 catalog、成对 SQL 同放一份**：`Parallel` 形态使 checksum / 分发给包复杂化；不选。
- **R3 顺带收口模块运行时 SQL 债**：范围爆炸、与 R4 职责重复；不选。
- **改 checksum 算法或 PG 文本入 digest**：v1.4 §4 明文禁止；不选。
- **postgres 仍禁止 apply 非空 catalog**：R3 矛盾；不选。

## 影响范围

- `apps/api/internal/kernel`（MigrationContribution 字段类型）；`internal/migration`（contract 校验）；`internal/store`（migrate runner + postgres apply + Open 守卫）；全部模块 `**/migration/` 包（签名 + 成对 SQL）；`internal/composition`（postgres 路由解闸）。R4 前的模块运行时仓库不动。

## 后续

- T1 形状迁移（自审）→ T2 postgres runner → T3 逐迁移对写 → T4 双路径证据 + self/independent 审计 → 关门。逐迁移对写时按列闭合 I-002；I-004 于 T4 前 live 对比。
