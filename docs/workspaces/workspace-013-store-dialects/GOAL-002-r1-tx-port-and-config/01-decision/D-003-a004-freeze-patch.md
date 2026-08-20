---
id: D-003-a004-freeze-patch
doc: decision-entry
status: accepted
created: 2026-08-20
updated: 2026-08-20
parent: GOAL-002-r1-tx-port-and-config
version: 1.0.0
---

# D-003 · 响应 A-004：冻结合同 v1.2（全部 required `fixed`）

- **日期**：2026-08-20
- **状态**：accepted
- **工作区**：`workspace-013-store-dialects`
- **用户确认**：本轮 `/govern` 输入「响应 GOAL-002 A-004」。A-004 要求必改走书面修正（非 residual / overruled）；F-001/F-002 写法由「与现行一致」与 Root R2/R3 分工唯一判定。F-003～F-005 recommended 一并 `fixed`（与 A-002 响应同口径）。

## 决定

修订 [r1-tx-port-and-config-freeze.md](../attachments/r1-tx-port-and-config-freeze.md) 至 **v1.2.0**。D-001 / D-002 的 Tx 形状、占位符 `?`、upsert 形态、配置键名、path 文件根、`WasFresh` 骨架、无 ORM、缺省 sqlite 仍成立。本条只补 A-004 指出的合同事实错误与 R2/R3 时序。GOAL-002 **保持 `done`**（不改检查点 2/2）。本回合不改 `apps/api`。

1. **F-001**（A-004 选项 1）：INTEGER 时间列允许 **秒或毫秒**，**按表/列沿用现行单位**；禁止无证据换单位。亚秒现行列继续 INTEGER 毫秒；RFC3339 文本不是现行 SQL 默认。附件列出抽样单位表（jobs / operation_log = 毫秒；wallet / `schema_migrations` / 若干模块 = 秒）。新 INTEGER 时间列必须在 R3 对写条目声明单位，禁止再写「一律 Unix 秒」当缺省。A-003 对 A-002 F-001 **时间半段**的 `fixed` 以本条重闭合；upsert 半段仍依 D-002。
2. **F-002**：R2 postgres `Open` = 连接 + `Ping` + 求值 `WasFresh`，**不** apply 现行 compiled catalog。允许空/`nil` catalog 作探测打开。postgres 上 apply 含 SQLite 专用 SQL 的 catalog → fail closed。sqlite 方言 R2 期间仍 apply。R3 完成后两方言均在 `WasFresh` 之后 apply 双方言 catalog。
3. **F-003**：点名 `authsession/migration/migration.go` 的 `sqlite_master` / `PRAGMA` 为 R3 改写债（与 `INSERT OR IGNORE` 并列）。
4. **F-004**：postgres 下 `db.path` 必须是文件路径形状（与缺省 `./data/schema-ui.db` 同形）；禁止配成目录。表头不再把 path 本身称作「数据根」；语义仍是 `filepath.Dir(path)`。
5. **F-005**：postgres `WasFresh` 只计用户**基表**（`BASE TABLE` / `relkind = 'r'`），不计视图/序列。`LIKE` 大小写与 INTEGER 0/1 布尔写入 §3，属 R3 物理 SQL，不能靠 rebind。

不另立 GOAL-002 I-00N：缺口是合同文本，不是新的未知信息。

## 理由

A-004 independent `conditional` 给出可核对反证：v1.1 把现行时间写成一律 Unix 秒，与 jobs / operation_log 的 INTEGER 毫秒冲突；`Open(..., catalog)` 把 catalog apply 冻进打开路径，但 Root R2 只含驱动/池/`readyz`、R3 才对写台账。用户要求响应本条；必改项均可在附件写清。

F-001 选「按列沿用」而不是「默认秒只约束新列」：原合同意图是「与现行一致」，纠正方式是承认现行单位，而不是给新列另立与现行部分列冲突的缺省。抽样表使该规则可核对，但不冒充全库穷尽。

F-002 选「只连 + Ping + WasFresh」：由 Root 纲领 R2≠R3 唯一决定。空 catalog 允许探测打开；sqlite-only catalog 在 postgres 上 fail closed，避免半执行。

## 未选方案

- F-001 选项 2（「默认秒」只约束新列、另附全量表）：新列缺省仍会与毫秒列并存冲突；全量穷尽属 R3 对写清单，不是本冻结能诚实完成的。
- F-001 把毫秒列改成秒或 RFC3339 文本：会改 jobs 租约/过期约 1000×，或改逻辑类型与 checksum。
- F-002 R2 postgres 即 apply 现有 catalog：现行 catalog 含 sqlite-only SQL，会失败或被迫提前改写迁移（侵入 R3）。
- F-002 跳过 apply 却不写进合同：实施者无法唯一选择。
- `accepted-residual` / `user-overruled`：用户未选。
- 重开 GOAL-002 或改 `status`：与 A-002 响应同口径；补丁记在本条 + 附件版本。GOAL-002 `done` 仍不证明冻结质量。

## 影响范围

- R2：postgres `Open` 不得 apply sqlite-only catalog；`WasFresh` 基表判定；`db.path` 文件路径形状校验。
- R3/R4：时间列按现行单位对写；点名 `sqlite_master`/`PRAGMA`；`LIKE` / 布尔。
- 不改 Profile、Manifest、模块矩阵、Compose 默认。
- Root I-002（驱动）仍 open，最晚 R2 方案冻结。

## 关联信息项

- GOAL-002 I-001：仍 verified（打开 + `WithTx`）；A-004 所指合同事实错误由本条写入 v1.2，不另立 I-00N。
- GOAL-002 I-002：仍 verified。
- Root I-002 / I-003 / I-001 / I-004：不变。F-002 不替代 Root I-002。

## 后续

R2 立项前先冻结 Root I-002（驱动）。实施合同已被 [D-004](D-004-a006-freeze-patch.md) 升为 **v1.3**（本条 v1.2 的时间单位与 `Open` 不 apply catalog 仍成立；时间宽度以 v1.3 为准）。建议 R2 方案冻结后再 `/audit`（independent）。
