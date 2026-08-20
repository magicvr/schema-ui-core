---
id: D-004-a006-freeze-patch
doc: decision-entry
status: accepted
created: 2026-08-20
updated: 2026-08-20
parent: GOAL-002-r1-tx-port-and-config
version: 1.0.0
---

# D-004 · 响应 A-006：冻结合同 v1.3（全部 required `fixed`）

- **日期**：2026-08-20
- **状态**：accepted
- **工作区**：`workspace-013-store-dialects`
- **用户确认**：本轮 `/govern` 输入「响应工作区12 GOAL-002 A-006」。工作区 12 的 GOAL-002 无 A-006；唯一待响应条目为 `[workspace-013] GOAL-002` A-006。A-006 要求必改走书面修正（非 residual / overruled）；合同写法在「一律 `BIGINT`」与「逻辑整数 + 物理成对」中择一。F-002～F-005 recommended 一并 `fixed`（与 A-002 / A-004 响应同口径）。

## 决定

修订 [r1-tx-port-and-config-freeze.md](../attachments/r1-tx-port-and-config-freeze.md) 至 **v1.3.0**。D-001 / D-002 / D-003 的 Tx 形状、占位符 `?`、upsert、时间**单位**（秒/毫秒按列）、配置键名、path 文件根、`WasFresh` 骨架、R2 postgres `Open` 不 apply catalog、无 ORM、缺省 sqlite 仍成立。本条只补 A-006 指出的时间**宽度**与 R2/R3 边界。GOAL-002 **保持 `done`**（不改检查点 2/2）。本回合不改 `apps/api`。

1. **F-001**（A-006 选项 1，并写明逻辑/物理拆开）：SQLite 时间列保持 `INTEGER`（64-bit affinity）；PostgreSQL 上 Unix 时间列（**秒与毫秒**）映射为 **`BIGINT` / `int8`**。禁止毫秒列（及本条范围内的秒列）使用 postgres `INTEGER`/`int4`。禁止把 SQLite `INTEGER` 字面当作两方言公共列类型或 postgres DDL。A-005 对 A-004 F-001 的**单位**修正仍然有效，不覆盖本条。
2. **F-002**：R2 postgres 本拍证据 = `Open` + `Ping`（及池）。不是 composition 全量 bootstrap，不是现行 HTTP `/readyz` 模块门禁全绿（`Store.Ping` + 模块图 `ready` + `SystemDataReady` 依赖 Reconcile/表）。
3. **F-003**：postgres 下 `db.path` 文件路径形状写可测谓词：拒绝尾部分隔符；已存在则拒绝 `Stat` 为目录；`Base` 非空且非 `.`/`..`；允许 cwd 相对文件（`schema-ui.db`，此时 `Dir` 为 `.` 合法）；禁止单靠 `Dir == "."` 拒绝。
4. **F-004**：postgres `WasFresh` 跟服务器实际解析后的 `search_path`，取第一个存在的用户 schema；禁止把字面 `"$user"` 当 schema 名；缺省落到 `public`。
5. **F-005**：R3 对写按列决定非时间 INTEGER 的 postgres 宽度；已点名钱包 `balance_*` / `amount_delta` / `balance_after_*` 用 `BIGINT`，不得只修时间列。

不另立 GOAL-002 I-00N：缺口是合同文本，不是新的未知信息。

## 理由

A-006 independent `conditional` 给出可核对反证：PostgreSQL `INTEGER` 是 int4，现行 jobs / operation_log 毫秒值约 1.7×10¹²，超出上限约 800×。v1.2 已纠正单位，但把公共列形态仍写成 `INTEGER`。用户要求响应本条；必改项可在附件写清。int4 上限是引擎事实，不构成「毫秒能否放进 postgres INTEGER」的 P-004 选边。

F-001 选选项 1（sqlite `INTEGER` + postgres `BIGINT`），并在合同里显式拆开逻辑类型与物理 SQL：实施者拿到的是可抄的 DDL 映射，同时满足选项 2 的分层。秒列一并 `BIGINT`：避免 2038 年 int4 溢出，也避免「秒用 int4、毫秒用 int8」两套规则。

F-002～F-005 与 A-004 响应同口径写入合同：都是 R2/R3 实施者会误读的边界，写进附件比只留审计建议便宜。

## 未选方案

- F-001 选项 2 只写「逻辑整数 + 物理成对」、不钉 postgres `BIGINT`：合法但 R3 仍可能把毫秒列对成 int4；本条把映射写死。
- 只禁止毫秒列用 int4、秒列仍用 postgres `INTEGER`：秒列至 2038 不溢出，但两套物理类型增加对写清单分叉，A-006 选项 1 已要求秒与毫秒都映射 `BIGINT`。
- 把 SQLite 时间列也改成显式 `BIGINT`：SQLite 无独立 `BIGINT` 存储类（`BIGINT` affinity 仍是 INTEGER）；改 DDL 字面不增加能力，且 R1 不改运行时/迁移文本。
- 把毫秒改成秒或 RFC3339 以迁就 int4：会改 jobs 租约约 1000×，或改逻辑类型与 checksum。
- `accepted-residual` / `user-overruled`：用户未选。
- 重开 GOAL-002 或改 `status`：与 A-002 / A-004 响应同口径；补丁记在本条 + 附件版本。GOAL-002 `done` 仍不证明冻结质量。

## 影响范围

- R2：`Open`/`Ping`/配置校验对照 **v1.3**；`db.path` 谓词；`WasFresh` `search_path` 解析。本拍不把现行 `/readyz` 当验收。
- R3：Unix 时间列 postgres DDL = `BIGINT`；钱包等非时间宽整数按列 `BIGINT`。F-001 闭合前不得把 `INTEGER` 时间列字面对写到 postgres。
- 不改 Profile、Manifest、模块矩阵、Compose 默认。
- Root I-002（驱动）仍 open，最晚 R2 方案冻结。

## 关联信息项

- GOAL-002 I-001：仍 verified（打开 + `WithTx`）；A-006 所指合同事实错误由本条写入 v1.3，不另立 I-00N。
- GOAL-002 I-002：仍 verified。
- Root I-002 / I-003 / I-001 / I-004：不变。F-002 不替代 Root I-002。

## 后续

R2 立项前先冻结 Root I-002（驱动）。实施合同已被 [D-005](D-005-a008-freeze-patch.md) 升为 **v1.4**（本条 v1.3 的时间宽度与 Open+Ping 证据边界仍成立；path 扩展名 / `COLLATE NOCASE` / checksum 输入 / 嵌套检测以 v1.4 为准）。A-008 已对 v1.3 复审为 pass。R2 方案冻结后再 `/audit`（independent）。
