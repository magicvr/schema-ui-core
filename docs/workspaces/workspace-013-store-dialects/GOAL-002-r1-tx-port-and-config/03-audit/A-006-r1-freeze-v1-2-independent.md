---
id: A-006
doc: audit-entry
record_id: A-006
source: independent
scope: GOAL-002 R1 方案冻结 v1.2（D-001 + D-002 + D-003 + attachments/r1-tx-port-and-config-freeze.md）是否合理、可否作 R2 实施合同
verdict: conditional
status: recorded
parent: GOAL-002-r1-tx-port-and-config
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# A-006 · R1 冻结合同 v1.2 独立交叉审计（2026-08-20）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：`design-plan` · GOAL-002 R1 方案冻结是否合理（现行权威 = 附件 **v1.2.0** + D-001 + D-002 + D-003）；含 A-005 对 A-004 的闭合是否经得起核对。不含 PG 驱动实现、台账对写、模块签名迁移的实施事实
- **verdict**：conditional
- **工作区**：`workspace-013-store-dialects`（Root `GOAL-001-store-dialects`；canonical `docs/workspaces/workspace-013-store-dialects/`；`shared_materials_catalog: none`）

## 范围与区间

- **covered**：冻结合同附件 v1.2.0 §1–§6；D-001 / D-002 / D-003 取舍；GOAL-002 成功标准 1–3 与 I-001/I-002；对照 VP-013 / RT-P03 内核端口声明与 Root R1/R2 字面范围；抽查现行 `store` / `config` / `composition` / `jobs` / `operationlog` / `wallet` / `logincaptcha` / `authsession` 迁移是否支撑 v1.2 主张。
- **excluded**：不审 R2 代码（尚无）；不改 `status` / `progress` / 方案正文 / goal-tree；不读取或比较其他工作区目标台账；不把 GOAL-002 已标 `done` 当作冻结质量证据。
- **P-005**：GOAL-002 I-001 / I-002 均标 `verified`、无到期未关 required 信息项。本条缺口是时间**存储宽度**在两方言上不等价，不是已到期 I-00N。Root I-002（驱动）仍 open、最晚 R2 方案冻结，不构成本冻结失败。
- **资料引用**：无。

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| Tx 公共面离开 `*sql.Tx`，`Run` 一事务、panic rollback、禁止嵌套 `Run` / `LastInsertId` / unwrap；Tx 仅回调内有效 | 附件 v1.2 §2–§3；现行 `store.WithTx` 确为 `func(*sql.Tx) error`（`apps/api/internal/store/store.go`）；仓库内无 `LastInsertId` 调用 |
| 占位符 `?` + rebind；upsert 允许 `ON CONFLICT(cols) DO UPDATE/NOTHING`，禁止 `INSERT OR IGNORE` 等 | 附件 §3；`operationlog/retention.go` 的 `INSERT OR IGNORE` 仍点名为 R3/R4 改写 |
| 时间**单位**：INTEGER 秒或毫秒、按列沿用；抽样表与现行读写一致 | jobs / operation_log = `UnixMilli`；wallet 账户/流水、`schema_migrations`、logincaptcha `expires_at` = `Unix()`；RFC3339 出现在 handler JSON，不是 SQL 列。与 A-004 F-001 所证冲突已消失 |
| R2 postgres `Open` = 连接 + `Ping` + `WasFresh`，不 apply 现行 compiled catalog；空/`nil` catalog 作探测；含 SQLite 专用 SQL 则 fail closed | 附件 §2「Open 与 catalog apply」；现行 `open()` 仍是 `databaseIsEmpty` 后立刻 `migrate`（sqlite 路径保留）；`authsession/migration/migration.go` 仍含 `sqlite_master` / `PRAGMA`（R3 点名债，代码未改属合同预期） |
| 配置键 `db.dialect` / `db.path` / `db.dsn`；缺省 sqlite；sqlite 下 DSN fail-closed；postgres 下 path 经 `filepath.Dir` 作文件根，且须为文件路径形状 | 附件 §5；现行 `yamlFile.DB` 仅 `Path`；`composition.go` 用 `Dir(cfg.DBPath)` 派生 uploads / brand-assets / avatars；`systemmonitoring.go` 对 path 无条件 `os.Stat`（R2 须按方言改，合同已写 `DBSizeBytes=0`） |
| `WasFresh` 求值在 catalog apply 前；postgres 只计用户基表 | 附件 §2；sqlite 侧对齐 `databaseIsEmpty`（`type = 'table' AND name NOT LIKE 'sqlite_%'`） |
| `ErrNoRows` 双 `errors.Is`；`Dialect()` 禁止模块按之分支 SQL；`OpenOptions` 可增字段；`LIKE` / INTEGER 0/1 布尔列入 R3 物理 SQL | 附件 §1 / §2 / §3；D-002 / D-003 |
| A-004 两条必改的**原问题**（单位写错；`Open` 与 catalog 时序未钉） | 本审同意 A-005 对这两点的文本修正可核对 |
| 工作区绑定合格 | `workspace.md`：id / Root / canonical / `primary_plan: VP-013-store-dialects` 一致；资料目录 `none` |

方向判断：**v1.2 作为 R1 Tx + 配置 + upsert + 时间单位 + R2 `Open` 时序骨架是合理的**；A-004 必改的原缺口已补上。时间类型仍把 SQLite 的 `INTEGER`（可 64-bit）写成两方言公共列类型，**PostgreSQL `INTEGER` 是 int4**，装不下现行毫秒时间戳，**还不能当作无条件的方言 SQL / R3 对写合同**。R2 的 `Open`/`Ping`/配置校验可以按 v1.2 做，但不得把「INTEGER 时间列」字面抄进 postgres DDL。

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| 1. Tx 端口合同落盘：公共类型不含 `*sql.Tx` / 驱动类型；`Run` 为一事务 | **部分** | 接口与事务语义、upsert 形态、时间**单位**已可实施；时间**存储宽度**仍按 SQLite `INTEGER` 书写，与 postgres int4 冲突（F-001） |
| 2. 配置键名与缺省/fail-closed；path 文件根闭合 | **满足（R2 映射层）** | 附件 §5；相对 A-002 F-002 / A-004 F-004 已闭合到可实施。文件路径形状的判定谓词仍略松（F-003 recommended） |
| 3. 未修改应用代码 | **本 scope 不作为冻结合理性主证** | E-004；与合同是否可实施正交 |

GOAL-002 自身成功标准仍窄于 VP-013 / RT-P03 全端口。D-002 / D-003 已把覆盖面降级为：连接/事务/占位符/upsert/时间 + 配置面；迁移 runner 对写与备份属 R3/R5。本审同意该降级。**时间一节单位已纠正，宽度未写**（F-001），不能把「时间类型已冻」当成已完成。

## Findings

### F-001 · 「INTEGER 时间列」在 PostgreSQL 上不是方言中立类型；毫秒列当前值超出 int4

| 字段 | 值 |
|------|-----|
| 严重度 | high |
| 建议 | required |
| 状态 | open |
| 影响门禁 | R1 冻结完整性（时间类型，非单位）；R3 台账对写 DDL；误跟合同会在 postgres 上溢出/插入失败 |
| 关联 I-00N | 无（GOAL-002 I-001 仍只覆盖打开路径与 `WithTx`） |

附件 v1.2 §3 已正确改成：INTEGER 时间列允许 **秒或毫秒**、**按列沿用**。这闭合了 A-004 F-001 的**单位**错误。同一段仍把公共列形态写成 **`INTEGER`**，并写「亚秒现行列继续用 INTEGER 毫秒」。

两方言的 `INTEGER` 不是同一宽度：

| 引擎 | `INTEGER` 含义 | Unix 毫秒（约 1.7×10¹²） |
|------|----------------|---------------------------|
| SQLite | 动态整数，可 8 字节 | 可存 |
| PostgreSQL | `int4`，上限 2 147 483 647 | **超出约 800×** |

现行毫秒列（抽样已在附件中）今天就不能放进 postgres `INTEGER`：

- jobs：`created_at` / `updated_at` / `lease_expires_at` / `expires_at` / `finished_at` — `jobs/model.go` `toMillis` = `UnixMilli()`
- operation_log：`created_at` / archive `archived_at` — `operationlog/repository.go`、`retention.go` `UnixMilli()`

按字面把 DDL 的 `INTEGER` 抄到 PostgreSQL，这些列会溢出或插入失败。秒列目前仍落在 int4 内（至 2038），不能拿来证明毫秒列可移植。

A-005 对 A-004 F-001 的 `fixed` **在单位半段成立**；**宽度半段不是 A-004 原问题，本条新开**，不回溯否定 A-005 对「一律 Unix 秒」的修正。

**要补什么（须书面，择一即可）**：

1. 在冻结合同时间节写明：SQLite 侧保持 `INTEGER`（64-bit affinity）；**PostgreSQL 上 Unix 时间列（秒与毫秒）映射为 `BIGINT` / `int8`**，禁止毫秒列用 `INTEGER`/`int4`。或
2. 把「逻辑类型 = 绑参 UTC 整数秒或毫秒」与「物理 SQL 类型按方言成对」拆开，并在 R3 物理 SQL 清单把毫秒列 → `BIGINT` 列为**硬规则**（不是可选）。

不得再用「INTEGER 时间列按列沿用」当作 postgres 可执行的列类型合同。

### F-002 · R2 postgres 探测打开不能当成现行 composition `/readyz` 全门禁

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | recommended |
| 状态 | open |

v1.2 已钉死：R2 postgres `Open` 只连 + `Ping` + `WasFresh`，允许空 catalog。Root R2 文案仍是「驱动、连接池、`readyz`」。现行 HTTP `/readyz`（`handler/health.go`）= `Store.Ping` **加上** 模块图 `ready`；`core.auth-session` Ready 还要求 `SystemDataReady()`（`composition.go` `withLifecycleHooks`）。`SystemDataReady` 在 `Reconcile` 成功后才 mark，而 Reconcile 需要表。

因此：**默认进程用 postgres DSN 走现行 `openStore`（必带 compiled catalog）→ 按合同 fail closed；空 catalog 探测打开则过不了 SystemDataReady，现行 `/readyz` 不会 200。** 这不是合同自相矛盾（fail closed 已写），但 R2 实施者可能把 Root 的「`readyz`」读成「现网探针变绿」。R2 方案应写明：postgres 本拍证据是 `Open` + `Ping`（及池），不是 composition 全量 bootstrap / 现行 `/readyz` 模块门禁。

### F-003 · postgres 下 `db.path`「文件路径形状」缺少唯一判定谓词

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |

A-004 F-004 已补：禁止把 path 配成目录（例 `./data` → `Dir` 变成 `.`）。未写启动时如何判定：尾部分隔符、`os.Stat` 为目录、无扩展名、`filepath.Dir` 为 `.` 是否一律拒绝。缺省 `./data/schema-ui.db` 清楚；`schema-ui.db`（cwd）与 `./data/` 的处理仍可能分叉。R2 配置校验宜写一条可测谓词（建议：拒绝尾部分隔符；拒绝 `Stat` 为目录；至少要求 `filepath.Base` 非空且 `Dir` 不得因「path 本身是目录」而落到 `.`）。不阻断 Open 骨架。

### F-004 · postgres `WasFresh` 的 `search_path` 未写解析规则

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |

附件要求「当前连接 `search_path` 解析出的第一个用户 schema（缺省 `public`）」内基表为 0。PostgreSQL 默认常为 `"$user", public`，且会跳过不存在的 `$user` schema。未写「按服务器实际解析、跳过不存在项」。R2 实现宜跟驱动/服务器解析，而不是把字面 `"$user"` 当 schema 名。不阻断探测打开。

### F-005 · 非时间 `INTEGER`（钱包余额等）同样可能超出 postgres int4

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |

`wallet_accounts.balance_*` / `amount_delta` 等在 SQLite 也是 `INTEGER`（`wallet/migration/migration.go`）。本条不把余额升为必改（冻结合同时间节未主张余额类型）。R3 对写应按列决定 `BIGINT`，不要只修时间列、把钱包金额留在 int4。可与 F-001 同一句物理类型规则处理。

## 必改项汇总

1. **F-001**：把 Unix 时间列的 **PostgreSQL 物理类型**写成 `BIGINT`（或等价 64-bit），禁止毫秒列使用 postgres `INTEGER`/`int4`。A-005 对 A-004 F-001 的单位修正仍然有效，不覆盖本条。

## 与既有意见的异同

| 项 | A-004 independent | A-005 self 响应 | A-006 independent（本条） |
|----|-------------------|-----------------|---------------------------|
| 对象合同 | v1.1 | 宣称 v1.2 闭合 A-004 | **v1.2 文本 + 现行代码** |
| verdict | conditional | pass（响应范围） | **conditional** |
| 时间单位（秒 vs 毫秒） | 必改 F-001 | 抽样表 + 按列沿用 | **同意单位已 fixed** |
| 时间宽度（int4 vs int8） | 未单列 | 未写 | **新必改 F-001** |
| `Open` vs catalog / R2·R3 | 必改 F-002 | 只连+Ping+WasFresh | **同意已 fixed** |
| path 文件根 / 文件路径形状 | F-004 recommended | 已写入 | 骨架同意；判定谓词 recommended |
| `WasFresh` 基表 / LIKE / 布尔 | F-005 recommended | 已写入 | 同意已写入；search_path 解析仅 recommended |
| `sqlite_master` / `PRAGMA` 点名 | F-003 recommended | 已点名 | **同意** |
| 可否无条件进 R2 | 否 | 允许对照 v1.2 | **Open/配置/Ping 可以对照 v1.2**；**不得**把 INTEGER 时间列字面作为 postgres DDL；F-001 闭合前不要开始 R3 对写 |

无 P-004 冲突需用户在「毫秒能否放进 postgres INTEGER」上选边：int4 上限是引擎事实。用户要选的是合同写法（一律 `BIGINT` vs 逻辑整数 + 物理成对），不是否认溢出。

A-001「与 VP-013 / RT-P03 同构」仍不成立（备份/迁移对写未冻；时间宽度未冻）。D-002 已降级该主张；本条不要求再改 A-001 原文。

## 结论 + 建议给编排器/用户的下一步

**方案冻结作为 R1 Tx + 配置 + upsert + 时间单位 + R2 postgres `Open` 时序是合理的；作为完整「时间类型」方言合同仍不合理。** 主因不再是 A-004 的单位写错或 catalog 时序，而是 SQLite `INTEGER` ≠ PostgreSQL `INTEGER`，现行毫秒列无法按字面迁移。verdict **conditional**。

建议 `/govern`：

1. 响应 A-006；按 D 条目修正附件（F-001；F-002～F-005 建议一并写进 R2/R3 方案边界）。
2. F-001 合法闭合前：**可以**按 v1.2 立项 R2（驱动仍受 Root I-002 约束），**不可以**开始把现行 `INTEGER` 时间 DDL 对写到 postgres。
3. 不要改本独立意见中的 `status` / `progress` 来「证明」冻结已过关。GOAL-002 已标 `done`、A-005 已标 `fixed`，只证明 A-004 的单位与 Open 时序已改，不证明时间宽度可移植。

## 声明

本意见 `source: independent`，不修改目标 `status` / 检查点 / 派生 `progress` / 方案正文 / goal-tree。响应与放行由 `/govern` 处理。
