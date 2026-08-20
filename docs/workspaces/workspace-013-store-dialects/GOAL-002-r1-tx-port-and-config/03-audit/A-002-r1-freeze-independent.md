---
id: A-002
doc: audit-entry
record_id: A-002
source: independent
scope: GOAL-002 R1 方案冻结（D-001 + attachments/r1-tx-port-and-config-freeze.md）是否足以作为 R2 实施合同
verdict: conditional
status: recorded
parent: GOAL-002-r1-tx-port-and-config
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# A-002 · R1 方案冻结独立交叉审计（2026-08-20）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：`design-plan` · GOAL-002 R1 方案冻结是否合理（D-001、冻结合同附件、E-001 扫描、A-001 自审主张）；不含 PG 驱动实现、台账对写、模块签名迁移的实施事实
- **verdict**：conditional
- **工作区**：`workspace-013-store-dialects`（Root `GOAL-001-store-dialects`；canonical `docs/workspaces/workspace-013-store-dialects/`；`shared_materials_catalog: none`）

## 范围与区间

- **covered**：冻结合同附件 §1–§6；D-001 取舍；GOAL-002 成功标准 1–3 与 I-001/I-002；对照 VP-013 / RT-P03 内核端口声明与 Root R1 字面范围；抽查现行 `store`/`config`/`composition`/`kernel.MigrationContribution` 是否支撑冻结主张。
- **excluded**：不审 R2 代码（尚无）；不改 `status`/`progress`/方案正文；不读取或比较其他工作区目标台账。
- **P-005**：GOAL-002 I-001/I-002 均标 `verified`、无到期未关 required 信息项。缺口是**未登记**的方言能力与配置语义，不是已到期 I-00N。Root I-002（驱动）仍 open、最晚 R2，不构成本冻结失败。
- **资料引用**：无。

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| Tx 公共面离开 `*sql.Tx`，`Run` 一事务、panic rollback、禁止嵌套 `Run` / `LastInsertId` / unwrap | 附件 §2–§3；现行 `store.WithTx` 确为 `func(*sql.Tx)`（`apps/api/internal/store/store.go`）；仓库内无 `LastInsertId` 调用 |
| 占位符 `?` + 实现 rebind，模块禁止 `switch` 方言 | 附件 §1、§3；与 VP-013「逻辑 schema 一份、物理 SQL 可成对」兼容 |
| 配置键 `db.dialect`/`db.path`/`db.dsn` + env；缺省 sqlite；不读 `DATABASE_URL` | 附件 §5；VP-013 配置面只要求 R1 冻结键名，未预先指定键名 |
| sqlite 下 DSN fail-closed、未知 dialect 拒绝、不改 Compose 默认、无 ORM | 附件 §5–§6、D-001 未选方案；现行 `config.default.yaml` 仅 `db.path` |
| 打开点与事务面扫描属实（范围内） | E-001 对齐 `OpenWithCatalog(cfg.DBPath)`、`yamlFile.DB` 仅 `Path`、`MigrationContribution.Apply func(*sql.Tx)` |
| 本目标未改运行时（文档主张） | E-002；本审未把 git 状态当冻结质量证明 |
| 工作区绑定合格 | `workspace.md`：id / Root / canonical / `primary_plan: VP-013-store-dialects` 一致 |

方向判断：**Tx 端口形状 + 配置键名这一刀是对的**，足以作为 R1 骨架，不能当作已经与 RT-P03 内核端口同构的完整实施合同。

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| 1. Tx 端口合同落盘：公共类型不含 `*sql.Tx` / 驱动类型；`Run` 为一事务 | **部分** | 附件给出接口；但 `kernel.Store` 混入 `WasFresh` 等非 Tx 方法且语义未中立（F-003） |
| 2. 配置键名与缺省/fail-closed 与 VP-013 配置面一致 | **部分** | 键名与缺省 sqlite 一致；postgres 下 `db.path`「不使用」与现行文件根耦合未闭合（F-002） |
| 3. 未修改应用代码 | **本 scope 不作为冻结合理性主证** | E-002 主张；与合同是否可实施正交 |

GOAL-002 自身范围写的是 Tx 方法面 + 配置键，**窄于** VP-013 / RT-P03「内核端口 = 连接/事务/占位符/**upsert/时间类型**/迁移 runner/备份与就绪」。A-001 把窄合同写成「与 VP-013 / RT-P03 同构」不成立（F-001）。

## Findings

### F-001 · 冻结合同未覆盖 upsert / 时间类型，却要求模块 SQL 可移植

| 字段 | 值 |
|------|-----|
| 严重度 | high |
| 建议 | required |
| 状态 | open |
| 影响门禁 | R1 冻结完整性；R2 不得把当前附件当完整内核端口开工；R3 台账对写 / R4 仓库 SQL 无对照合同 |
| 关联 I-00N | 无（Root 与 GOAL-002 均未登记；属 P-005 发现后应回流的未登记 required） |

冻结合同 §1 禁止模块 `switch` 方言、§3 只冻了 `?` rebind；§6 推迟表未列 upsert / 时间类型。RT-P03 与 VP-013 内核端口明确包含「占位符、**upsert、时间类型**」。现行代码已混用可移植 `ON CONFLICT` 与 SQLite 专有 `INSERT OR IGNORE`（`apps/api/internal/modules/operationlog/retention.go`），时间值多为 Unix/millis 绑参、并非统一 SQL 时间类型。R2 若只实现 `Run`/`Tx`，R3/R4 仍会把方言写进模块 SQL。A-001「合同与 VP-013 / RT-P03 同构」与附件覆盖面不符。

**要补什么（择一，须书面）**：把 upsert（允许的 SQL 形态 / 禁止 `INSERT OR IGNORE` 等）与时间存储（绑参 UTC、禁止 `datetime('now')` 一类方言函数等）写入本冻结；或登记 required I-00N，最晚不晚于 R3 方案冻结，并写明 R2 不得扩展模块 SQL 方言面。

### F-002 · postgres 下 `db.path`「不使用」与现行文件根耦合冲突

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | required |
| 状态 | open |
| 影响门禁 | R1 配置面冻结；R2 配置校验 / `Open` 映射 |
| 关联 I-00N | GOAL-002 I-001（关闭过窄：只证了打开路径与 `WithTx` 形状） |

附件 §5：postgres 时 `db.dsn` 必填、`db.path` **可残留但不使用**。现行 `composition` 用 `filepath.Dir(cfg.DBPath)` 派生 `uploads` / `brand-assets` / `avatars`，并把 `cfg.DBPath` 传入 `systemmonitoringmodule.New`（`apps/api/internal/composition/composition.go`）。E-001 未扫描这些非 SQL 消费点。R2 无法同时字面执行「不用 path」又保持现网文件目录。须在冻结中明确：path 在 postgres 下是否仍作数据目录根；若否，缺的配置键是什么；空 path + 默认 `./data/schema-ui.db` 在 postgres 下是否仍派生文件根。

### F-003 · `WasFresh()` 已冻进 `kernel.Store`，无方言中立语义

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | required |
| 状态 | open |
| 影响门禁 | R2 `store.Open` 实现 |
| 关联 I-00N | — |

附件 §2 把 `WasFresh` / `MarkSystemDataReady` / `SystemDataReady` 放进持久化端口。现行 `WasFresh` 来自 `sqlite_master` 空库探测（`store.go` `databaseIsEmpty`）。附件未定义 postgres 的「fresh」（无表 / 无迁移元表 / 无业务表）。`MarkSystemDataReady` 为进程内 atomic，可中立；`WasFresh` 不能靠「照抄 SQLite」实现。R2 会猜测。须补方言中立定义，或把 bootstrap 就绪从 `kernel.Store` 拆出并标明非 R1 端口。

### F-004 · `OpenOptions`、并发 `Run`、Tx 生命周期未写可扩展边界

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |

`OpenOptions` 只点了 Dialect/Path/DSN；池参数属 R2。未写「可在不改 Tx 合同的前提下为 OpenOptions 增字段」。并发 `Run`（sqlite `MaxOpenConns=1` vs PG 池）、Tx 仅在 `Run` 回调内有效、回调返回后 Tx 失效，均未写。不阻断骨架，R2 实施前宜补一句，避免把 R1 形状误锁死或误允许回调外使用 Tx。

### F-005 · `kernel.ErrNoRows` 与 `sql.ErrNoRows` 关系未冻

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |

附件要求模块 `errors.Is(err, kernel.ErrNoRows)`。未规定是别名、wrap 使得两者皆 `Is`，还是切断 `sql.ErrNoRows`。模块改走 `kernel.Tx` 时（R4）若只换签名不换 sentinel，查询未命中会静默走错分支。建议冻：实现必须使 `errors.Is(err, kernel.ErrNoRows)` 与 `errors.Is(err, sql.ErrNoRows)` 同真，直至公共面不再 import `database/sql`。

### F-006 · `Dialect()` 在公共 Store 上未限制调用方

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |

§1 禁止模块 `switch` 方言，§2 又把 `Dialect()` 放在 `kernel.Store`。composition / `readyz` 需要它是合理的；不写清「模块仓库禁止按 `Dialect()` 分支 SQL」，R4 会把方言判断从 `*sql.Tx` 换一种泄漏方式。

### F-007 · A-001 把门禁审深放到 R2，不能替代本冻结的独立审

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |

Root D-001 书面选择 R1 冻结 `self`、迁移/生产实施 `independent`。GOAL-002 I-002 据此标 verified。P-003 对 data/compatibility 高影响门禁的默认是 `independent`。本条补上冻结层独立审。**不**要求仅因模式选择而改 GOAL-002 `status`；**要求**编排器不得以 A-001 `pass` 或「目标已 done」跳过 F-001～F-003。R2 实施独立审不能代替冻结合同修正。

## 必改项汇总

1. **F-001**：补冻 upsert + 时间类型，或登记最晚 R3 方案前的 required I-00N，并禁止 R2 把当前附件当完整端口。
2. **F-002**：写清 postgres 下 `db.path` 与 uploads/avatars/brand-assets/监控文件根的关系（沿用 / 新键 / 默认 `./data`）。
3. **F-003**：给 `WasFresh()` 方言中立语义，或将其移出 R1 `kernel.Store` 端口。

## 与既有意见的异同

| 项 | A-001 self | A-002 independent |
|----|------------|-------------------|
| verdict | pass | **conditional** |
| 开放 required | 0 | **F-001 / F-002 / F-003** |
| Tx ≠ `*sql.Tx`、配置键、无 ORM、缺省 sqlite | 同意 | 同意 |
| 「与 VP-013 / RT-P03 同构」 | 主张成立 | **不成立**（缺 upsert/时间；A-001 过称） |
| I-001 足以冻结 | 同意（打开 + WithTx） | 同意足以冻 **Tx 骨架**；不足以冻配置语义与方言 SQL 面 |
| 可无条件进 R2 | 允许按附件实现 | **否**；先闭合必改或用户书面 residual |

无 P-004 冲突需用户在「同构 vs 不同构」上选边：独立审给出可核对反证，编排器应修正合同或降级 A-001 主张，而不是两说并存放行 R2。

## 结论 + 建议给编排器/用户的下一步

**方案冻结作为 R1 Tx+配置骨架是合理的；作为 R2 唯一实施合同尚不合理。** verdict **conditional**。

建议 `/govern`：

1. 响应 A-002；闭合 F-001～F-003（修正附件 + D 条目，或用户书面 residual，范围须写清「R2 允许猜什么」）。
2. 在 required 合法闭合前 **不要立项/实施 R2** 把当前附件当完整端口。
3. 不要改本独立意见中的 `status`/`progress` 字段来「证明」冻结已过关。GOAL-002 已标 `done` 不构成冻结质量证据。

## 声明

本意见 `source: independent`，不修改目标 `status` / 检查点 / 派生 `progress` / 方案正文 / goal-tree。响应与放行由 `/govern` 处理。
