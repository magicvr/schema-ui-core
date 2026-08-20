---
id: A-008
doc: audit-entry
record_id: A-008
source: independent
scope: GOAL-002 R1 方案冻结 v1.3（D-001 + D-002 + D-003 + D-004 + attachments/r1-tx-port-and-config-freeze.md）是否合理、可否作 R2 实施合同
verdict: pass
status: recorded
parent: GOAL-002-r1-tx-port-and-config
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# A-008 · R1 冻结合同 v1.3 独立交叉审计（2026-08-20）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：`design-plan` · GOAL-002 R1 方案冻结是否合理（现行权威 = 附件 **v1.3.0** + D-001 + D-002 + D-003 + D-004）；含 A-007 对 A-006 的闭合是否经得起核对。不含 PG 驱动实现、台账对写、模块签名迁移的实施事实
- **verdict**：pass
- **工作区**：`workspace-013-store-dialects`（Root `GOAL-001-store-dialects`；canonical `docs/workspaces/workspace-013-store-dialects/`；`shared_materials_catalog: none`）

## 范围与区间

- **covered**：冻结合同附件 v1.3.0 §1–§6；D-001～D-004 取舍；GOAL-002 成功标准 1–3 与 I-001/I-002；对照 VP-013 / RT-P03 内核端口声明与 Root R1/R2 字面范围；抽查现行 `store` / `config` / `composition` / `handler/health.go` / `jobs` / `operationlog` / `wallet` / `authsession`（含 `COLLATE NOCASE`）/ `kernel.MigrationChecksum` 是否支撑 v1.3 主张。
- **excluded**：不审 R2 代码（尚无）；不改 `status` / `progress` / 方案正文 / goal-tree；不读取或比较其他工作区目标台账；不把 GOAL-002 已标 `done` 当作冻结质量证据。
- **P-005**：GOAL-002 I-001 / I-002 均标 `verified`、无到期未关 required 信息项。本条无新的到期 required 信息缺口。Root I-002（驱动）仍 open、最晚 R2 方案冻结，不构成本冻结失败。
- **资料引用**：无。

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| Tx 公共面离开 `*sql.Tx`，`Run` 一事务、panic rollback、禁止嵌套 `Run` / `LastInsertId` / unwrap；Tx 仅回调内有效 | 附件 v1.3 §2–§3；现行 `store.WithTx` 确为 `func(*sql.Tx) error`（`apps/api/internal/store/store.go`）；仓库内无 `LastInsertId` 调用 |
| 占位符 `?` + rebind；upsert 允许 `ON CONFLICT(cols) DO UPDATE/NOTHING`，禁止 `INSERT OR IGNORE` 等 | 附件 §3；`operationlog/retention.go` 的 `INSERT OR IGNORE` 仍点名为 R3/R4 改写 |
| 时间**单位**：INTEGER 秒或毫秒、按列沿用；抽样表与现行读写一致 | jobs / operation_log = `UnixMilli`；wallet 账户/流水、`schema_migrations`、logincaptcha `expires_at` = `Unix()`。与 A-004 F-001 所证冲突已消失 |
| 时间**宽度**：逻辑 UTC 整数与物理 SQL 拆开；sqlite 时间列保持 `INTEGER`；postgres Unix 时间列（秒与毫秒）= `BIGINT`/`int8`；禁止 int4、禁止把 sqlite `INTEGER` 字面当 postgres DDL | 附件 §3「时间存储」；D-004 决定 1。现行毫秒值约 1.7×10¹² 超出 int4 上限，合同已禁止按字面抄 `INTEGER` |
| R2 postgres `Open` = 连接 + `Ping` + `WasFresh`，不 apply 现行 compiled catalog；空/`nil` catalog 作探测；含 SQLite 专用 SQL 则 fail closed | 附件 §2「Open 与 catalog apply」；现行 `open()` 仍是 `databaseIsEmpty` 后立刻 `migrate`（sqlite 路径保留） |
| R2 证据边界 ≠ 现行 HTTP `/readyz` 模块门禁全绿 | 附件 §2「R2 证据边界」；`handler/health.go` `readyz` = `Store.Ping` + 模块图 `ready`；`SystemDataReady` 依赖 Reconcile/表。A-006 F-002 原文问题已写入合同 |
| 配置键 `db.dialect` / `db.path` / `db.dsn`；缺省 sqlite；sqlite 下 DSN fail-closed；postgres 下 path 经 `filepath.Dir` 作文件根 | 附件 §5；现行 `yamlFile.DB` 仅 `Path`；`composition.go` 用 `Dir(cfg.DBPath)` 派生 uploads / brand-assets / avatars |
| `WasFresh` 求值在 catalog apply 前；postgres 只计用户基表，并跟服务器解析后的 `search_path`（禁止字面 `"$user"`） | 附件 §2；sqlite 侧对齐 `databaseIsEmpty` |
| 非时间 INTEGER 宽度按列决定；已点名钱包 `balance_*` / `amount_delta` / `balance_after_*` | 附件 §3；`wallet/migration/migration.go` |
| `ErrNoRows` 双 `errors.Is`；`Dialect()` 禁止模块按之分支 SQL；`OpenOptions` 可增字段；`LIKE` / INTEGER 0/1 布尔列入 R3 物理 SQL | 附件 §1 / §2 / §3 |
| A-006 一条必改（时间宽度）的**原问题** | 本审同意 A-007 对 F-001 的文本修正可核对；**宽度半段现已 fixed** |
| 工作区绑定合格 | `workspace.md`：id / Root / canonical / `primary_plan: VP-013-store-dialects` 一致；资料目录 `none` |

方向判断：**v1.3 作为 R1 Tx + 配置 + upsert + 时间（单位与宽度）+ R2 postgres `Open`/`Ping`/`WasFresh` 时序是合理的，可以作 R2 实施合同。** A-002 / A-004 / A-006 的 required 原缺口均已补上。剩余项不阻断 R2 立项（仍受 Root I-002 驱动选型约束），但 R2 配置校验与 R3 对写仍有应写清的 recommended 边界（F-001～F-004）。

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| 1. Tx 端口合同落盘：公共类型不含 `*sql.Tx` / 驱动类型；`Run` 为一事务 | **满足**（本目标范围内） | 接口、事务语义、upsert、时间单位与宽度均可实施；方言 SQL 全量可移植仍属 R3，附件已降级覆盖面 |
| 2. 配置键名与缺省/fail-closed；path 文件根闭合 | **满足（R2 映射层）** | 附件 §5；相对 A-002 F-002 已闭合。文件路径谓词有一处与自身示例不完全同构（F-001 recommended），缺省 `./data/schema-ui.db` 仍通过 |
| 3. 未修改应用代码 | **本 scope 不作为冻结合理性主证** | E-005；与合同是否可实施正交 |

GOAL-002 自身成功标准仍窄于 VP-013 / RT-P03 全端口。D-002 已把覆盖面降级为：连接/事务/占位符/upsert/时间 + 配置面；迁移 runner 对写与备份属 R3/R5。本审同意该降级，且 **时间一节单位 + 宽度现均可核对**。不得把本 pass 读成 RT-P03 全端口已冻，也不得把 GOAL-002 `done` 当冻结质量证明。

## Findings

### F-001 · postgres `db.path` 谓词仍放行「尚不存在的 `./data`」，与合同自己的禁例不完全同构

| 字段 | 值 |
|------|----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |
| 影响门禁 | R2 配置校验实现；不阻断 Open/Ping 骨架 |
| 关联 I-00N | 无 |

A-006 F-003 要求可测谓词；A-007 / D-004 宣称 `fixed`。v1.3 §5 条 5 **散文**仍禁止把 path 配成目录，并点名 `./data`（`filepath.Dir` → `.`，文件根错位）。编号谓词却是：

1. 拒绝尾部分隔符（`./data/` 会拦）
2. **仅当路径已存在**时拒绝 `Stat` 为目录
3. `Base` 非空且非 `.` / `..`
4. 允许 cwd 文件 `schema-ui.db`（此时 `Dir == "."` 合法）
5. 不要单靠 `Dir == "."` 拒绝

对**尚不存在**的 `./data` / `.\data`：无尾部分隔符、`Stat` 跳过、`Base == "data"` → **谓词通过**，随后 `filepath.Dir("./data") == "."`，正是散文禁止的错位。缺省 `./data/schema-ui.db` 仍通过，已存在的 `./data` 目录仍拒绝。

因此 A-007 对 A-006 F-003 的 `fixed` **在缺省路径与「已存在目录」半段成立**；**「禁止 `./data`」半段未由 1–3 实现**。不是 A-006 原 required，本条不升级为必改。

**要补什么（R2 配置校验，择一即可）**：加一条与缺省同构、且能拦住无扩展名目录分量的谓词。建议：`filepath.Ext(filepath.Base(path))` 必须非空（缺省 `.db` 通过；`./data` 失败；`schema-ui.db` / `./schema-ui.db` 通过）。不要再写「判定用 1–3」却把 `./data` 当禁例。

### F-002 · `COLLATE NOCASE` 是不能靠 rebind 的 SQLite 专用 SQL，尚未点名

| 字段 | 值 |
|------|----|
| 严重度 | med |
| 建议 | recommended |
| 状态 | open |
| 影响门禁 | R3 台账对写 / R4 仓库查询；不阻断 R2 Open |
| 关联 I-00N | 无 |

附件已点名 `INSERT OR IGNORE`、`sqlite_master` / `PRAGMA`，并把 `LIKE` 大小写列入 R3 物理 SQL。现行仍有 SQLite 校对：

- DDL：`authsession/migration/migration.go` `service_credentials.name … COLLATE NOCASE UNIQUE`
- 查询：`users_repository.go` / `roles_repository.go` 的 `ORDER BY … COLLATE NOCASE`

PostgreSQL 无名为 `NOCASE` 的 collation：按字面 apply 或排序会失败；丢掉 collation 则 `UNIQUE` / 排序语义与 sqlite 不等价（大小写视为不同）。与已列入的 `LIKE` 同类，且出现在 compiled Apply。R3 方案应点名并成对改写（`CITEXT` / `LOWER()` 唯一索引 / 显式 collation），不要只扫 `PRAGMA`。

本条不把「未穷尽全部 SQLite 专用 SQL」升为必改：附件已有类规则「不能靠 rebind 的 SQLite 专用 SQL，R3 必须成对改写」；缺口是**点名清单漏了会直接 apply 失败的一句**。

### F-003 · 「checksum 算法不变」未说明成对物理 SQL 的输入是哪份文本

| 字段 | 值 |
|------|----|
| 严重度 | med |
| 建议 | recommended |
| 状态 | open |
| 影响门禁 | R3 台账对写；**不阻断** R2（postgres 本拍不 apply catalog） |
| 关联 I-00N | 无 |

现行算法（`kernel.MigrationChecksum`）：规范化 SQL 文本 + transform id → SHA-256；`CollectPersistence` 要求 **checksum 全局唯一**，且已应用库按该值 fail-closed 防漂移。

v1.3 同时要求：(1) 物理 SQL 可成对（时间列 sqlite `INTEGER` / postgres `BIGINT`）；(2) checksum **算法不变**。两方言 DDL 文本不同 → 摘要不同。一个 compiled catalog 不能为同一 `version` 挂两条不同 checksum；若改写 sqlite 历史 stmts 去对齐 `BIGINT`，现网 sqlite 库会 checksum 漂移 fail-closed。

附件 §4 把 checksum 写进「类型冻结」，但没写 **digest 的输入是 sqlite/canonical 文本、postgres 文本，还是两套独立 catalog**。R3 实施者不能从本句唯一推出合法做法。

**要补什么（R3 方案冻结，不改 hash 函数）**：写明 sqlite 历史 stmts + transform id 仍是该 version 的 checksum 权威（现网库可核对）；postgres 成对 SQL 要么不进入该 digest（另有对写证据），要么走**按方言分列的 compiled catalog**（版本对齐、checksum 各算）。禁止把「算法不变」读成「两方言 DDL 必须 byte-identical」。

### F-004 · 嵌套 `Run` 的检测不能做成 Store 级互斥，否则与 PG 并发 `Run` 冲突

| 字段 | 值 |
|------|----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |
| 影响门禁 | R2 `Run` 实现；不阻断驱动选型 / Open+Ping |
| 关联 I-00N | 无 |

附件同时规定：禁止同一 Store 在 `fn` 内再 `Run`；Postgres 连接池允许**不同连接上并发 `Run`**。若用 Store 级 `inRun` 标志检测嵌套，并发 `Run` 会被误判为嵌套而 fail closed。R2 实现须按调用栈 / `ctx` 值检测「当前回调内」，不得用进程级或 Store 级互斥冒充嵌套门禁。SQLite `MaxOpenConns=1` 串行下 Store 级标志碰巧能用，不能当 postgres 合同。

## 必改项汇总

无。本条 **开放 required = 0**。

## 与既有意见的异同

| 项 | A-006 independent | A-007 self 响应 | A-008 independent（本条） |
|----|-------------------|-----------------|---------------------------|
| 对象合同 | v1.2 | 宣称 v1.3 闭合 A-006 | **v1.3 文本 + 现行代码** |
| verdict | conditional | pass（响应范围） | **pass** |
| 时间宽度（int4 vs int8） | 必改 F-001 | sqlite `INTEGER` + postgres `BIGINT` | **同意已 fixed** |
| R2 证据 ≠ 现行 `/readyz` | F-002 recommended | 已写入 | **同意已 fixed** |
| path 文件路径谓词 | F-003 recommended | 宣称已写入 | 骨架同意；**`./data` 禁例与谓词 1–3 不完全同构**（本条 F-001 recommended） |
| `WasFresh` `search_path` | F-004 recommended | 禁止字面 `"$user"` | **同意已写入** |
| 非时间 INTEGER 宽度 | F-005 recommended | 点名钱包 | **同意已写入** |
| `COLLATE NOCASE` / checksum 输入 / 嵌套检测 | 未单列 | 未写 | **新 recommended F-002～F-004** |
| 可否无条件进 R2 | Open/配置/Ping 对照 v1.2；不得抄 INTEGER 时间列 | 对照 **v1.3** | **可以对照 v1.3 立项 R2**（仍受 Root I-002）；R3 对写须 BIGINT + 处理 F-002/F-003 |

无 P-004 冲突。A-006 F-001 的引擎事实（int4 装不下毫秒）现已写入合同，本条不重新打开。

A-001「与 VP-013 / RT-P03 同构」仍不成立（备份/迁移对写未冻）。D-002 已降级该主张；本条不要求再改 A-001 原文。

## 结论 + 建议给编排器/用户的下一步

**方案冻结作为 R1 Tx + 配置 + upsert + 时间（单位与宽度）+ R2 postgres Open/Ping/WasFresh 合同是合理的。** 主因已从 A-006 的 int4 溢出，收成 recommended 级的配置谓词边角、R3 点名清单与 checksum 输入。verdict **pass**。

建议 `/govern`：

1. 响应 A-008。F-001～F-004 均为 recommended：可写入 v1.3 补丁、或登记到 R2/R3 方案边界；**不是**本冻结放行的必改门禁。
2. **可以**按 v1.3 立项 R2（驱动仍受 Root I-002 约束）。**不可以**对照 v1.2 把 `INTEGER` 时间列抄进 postgres DDL，**不可以**开始 R3 对写直到实施计划显式引用 v1.3 宽度规则，并处理 F-002/F-003。
3. 不要改本独立意见中的 `status` / `progress` 来「证明」冻结已过关。GOAL-002 已标 `done`、A-007 已标 A-006 `fixed`，只证明 A-006 的宽度与当时 recommended 已改写进合同，不证明 RT-P03 全端口或 R3 SQL 已可移植。

## 声明

本意见 `source: independent`，不修改目标 `status` / 检查点 / 派生 `progress` / 方案正文 / goal-tree。响应与放行由 `/govern` 处理。
