---
id: A-004
doc: audit-entry
record_id: A-004
source: independent
scope: GOAL-002 R1 方案冻结 v1.1（D-001 + D-002 + attachments/r1-tx-port-and-config-freeze.md）是否合理、可否作 R2 实施合同
verdict: conditional
status: recorded
parent: GOAL-002-r1-tx-port-and-config
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# A-004 · R1 冻结合同 v1.1 独立交叉审计（2026-08-20）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：`design-plan` · GOAL-002 R1 方案冻结是否合理（现行权威 = 附件 **v1.1.0** + D-001 + D-002）；含 A-003 对 A-002 的闭合是否经得起核对。不含 PG 驱动实现、台账对写、模块签名迁移的实施事实
- **verdict**：conditional
- **工作区**：`workspace-013-store-dialects`（Root `GOAL-001-store-dialects`；canonical `docs/workspaces/workspace-013-store-dialects/`；`shared_materials_catalog: none`）

## 范围与区间

- **covered**：冻结合同附件 v1.1.0 §1–§6；D-001 / D-002 取舍；GOAL-002 成功标准 1–3 与 I-001/I-002；对照 VP-013 / RT-P03 内核端口声明与 Root R1/R2 字面范围；抽查现行 `store` / `config` / `composition` / `jobs` / `operationlog` / `authsession` 迁移是否支撑 v1.1 主张。
- **excluded**：不审 R2 代码（尚无）；不改 `status` / `progress` / 方案正文 / goal-tree；不读取或比较其他工作区目标台账；不把 GOAL-002 已标 `done` 当作冻结质量证据。
- **P-005**：GOAL-002 I-001 / I-002 均标 `verified`、无到期未关 required 信息项。本条缺口是合同文本与现行代码不一致、以及 R2 `Open` 与 catalog 对写的时序未钉死，不是已到期 I-00N。Root I-002（驱动）仍 open、最晚 R2 方案冻结，不构成本冻结失败。
- **资料引用**：无。

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| Tx 公共面离开 `*sql.Tx`，`Run` 一事务、panic rollback、禁止嵌套 `Run` / `LastInsertId` / unwrap；Tx 仅回调内有效 | 附件 v1.1 §2–§3；现行 `store.WithTx` 确为 `func(*sql.Tx) error`（`apps/api/internal/store/store.go`）；仓库内无 `LastInsertId` 调用 |
| 占位符 `?` + 实现 rebind；upsert 允许 `ON CONFLICT(cols) DO UPDATE/NOTHING`，禁止 `INSERT OR IGNORE` 等 | 附件 §3；现行多数 upsert 已是 `ON CONFLICT(...)`（如 `authsession/systemdata`、`settings/repository`）；`operationlog/retention.go` 的 `INSERT OR IGNORE` 已点名为 R3/R4 改写 |
| 配置键 `db.dialect` / `db.path` / `db.dsn`；缺省 sqlite；sqlite 下 DSN fail-closed；不读 `DATABASE_URL` | 附件 §5；现行 `yamlFile.DB` 仅 `Path`（`config.go`）；`composition.go` 用 `filepath.Dir(cfg.DBPath)` 派生 uploads / brand-assets / avatars，并把 path 传给 system-monitoring |
| postgres 下 `db.path` 不作 SQL、仍经 `filepath.Dir` 作文件根；省略时缺省 `./data/schema-ui.db`；监控 postgres 下 `DBSizeBytes=0` | 附件 §5 条 3–6；D-002 决定 2；现行 `handler/systemmonitoring.go` 对 path 无条件 `os.Stat`（R2 须按方言改，合同已写） |
| `WasFresh` 留在 `kernel.Store`；求值在 migrate 前；sqlite 对齐现行 `databaseIsEmpty` | 附件 §2；`store.go` `open()`：先 `databaseIsEmpty` 再 `migrate` |
| `ErrNoRows` 双 `errors.Is`；`Dialect()` 禁止模块按之分支 SQL；`OpenOptions` 可增字段 | 附件 §1 / §2 / §3；D-002 决定 4–6 |
| A-002 骨架判断仍成立 | Tx ≠ `*sql.Tx`、缺省 sqlite、无 ORM、不改 Compose 默认：本审同意 |
| 工作区绑定合格 | `workspace.md`：id / Root / canonical / `primary_plan: VP-013-store-dialects` 一致；资料目录 `none` |

方向判断：**v1.1 作为 R1 Tx + 配置 + upsert 骨架是合理的**；时间单位合同与现行代码不符，且 R2 `Open(..., catalog)` 与 R3 台账对写的边界未钉死，**还不能当作无条件的 R2 实施合同**。

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| 1. Tx 端口合同落盘：公共类型不含 `*sql.Tx` / 驱动类型；`Run` 为一事务 | **部分** | 接口与事务语义已可实施；方言 SQL 面的时间单位主张与现行 INTEGER 列不一致（F-001） |
| 2. 配置键名与缺省/fail-closed；path 文件根闭合 | **满足（R2 映射层）** | 附件 §5 相对 A-002 F-002 已闭合；`Dir(path)` 与 composition 对齐 |
| 3. 未修改应用代码 | **本 scope 不作为冻结合理性主证** | E-003；与合同是否可实施正交 |

GOAL-002 自身成功标准仍窄于 VP-013 / RT-P03 全端口。D-002 已把「同构」降级为：连接/事务/占位符/upsert/时间 + 配置面；迁移 runner 对写与备份属 R3/R5。本审同意该降级。**时间一节写入了，但写错了现行事实**（F-001），不能再当作 A-002 F-001 已 `fixed`。

## Findings

### F-001 · 时间存储合同声称「与现行一致：INTEGER Unix 秒」，与现行列单位不符

| 字段 | 值 |
|------|-----|
| 严重度 | high |
| 建议 | required |
| 状态 | open |
| 影响门禁 | R1 冻结完整性（A-002 F-001 时间半段的闭合）；R3 台账对写 / R4 仓库时间列；误跟合同会改 jobs 租约与 operation_log 语义 |
| 关联 I-00N | 无（GOAL-002 I-001 仍只覆盖打开路径与 `WithTx`；属未登记的合同事实错误） |

附件 v1.1 §3「时间存储」写：**默认列形态与现行一致：INTEGER Unix 秒**（`time.Now().UTC().Unix()`）；亚秒则改 RFC3339 **文本列**。D-002 / A-003 据此把 A-002 F-001 标为 `fixed`。

现行不是单一 Unix 秒：

| 面 | 单位 | 证据 |
|----|------|------|
| jobs 时间列（`created_at` / `updated_at` / `lease_expires_at` / `expires_at` / `finished_at`） | **INTEGER 毫秒** | `apps/api/internal/jobs/model.go` `toMillis` / `fromMillis`（`UnixMilli`）；DDL `apps/api/internal/modules/jobs/migration/migration.go` 同为 INTEGER |
| operation_log `created_at` / retention cutoff | **INTEGER 毫秒** | `operationlog/retention.go` 注释「created_at is unix milliseconds」+ `UnixMilli()`；`repository.go` 读写同样毫秒 |
| 若干模块与 `schema_migrations` | INTEGER **秒** | 如 `wallet/store/repository.go` `now.Unix()`；`store/migrate.go` `time.Now().UTC().Unix()` |
| ID 前缀（非时间列） | 毫秒 hex | `jobs.NewID`、wallet 流水 id；与列单位无关，不得拿来证明「列都是秒」 |

按字面执行 v1.1 会把 jobs 租约/过期从毫秒改成秒（约 1000× 偏差）或把毫秒列改成 RFC3339 文本（改逻辑类型与 checksum）。亚秒路径指定 RFC3339 文本，也覆盖不了已经存在的 INTEGER 毫秒列。

A-003 对 A-002 F-001 的 `fixed` **在时间半段不成立**（upsert 半段仍成立）。本条不是新发明时间规则，是指出 v1.1 写了错误现行事实。

**要补什么（须书面，择一即可）**：

1. 改冻结合同：INTEGER 时间列允许 **秒或毫秒**，**按表/列沿用现行单位**，禁止无证据换单位；亚秒继续用 INTEGER 毫秒，不把 RFC3339 文本当现行默认。或
2. 列一张现行时间列单位表（秒 / 毫秒 / 文本）作为 R3 对照，并写明「默认秒」只约束**新列**。

不得再用「与现行一致：一律 Unix 秒」关闭本条。

### F-002 · `Open(..., catalog)` 把 catalog apply 冻进打开路径，未钉 R2 postgres 与 R3 对写的边界

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | required |
| 状态 | open |
| 影响门禁 | R2 实施方案（`store.Open` / `WasFresh` / 是否 apply catalog） |
| 关联 I-00N | Root I-002（驱动，最晚 R2 方案冻结）不替代本条；本条是 Open 语义 |

现行 `OpenWithCatalog`：空库探测 → **立刻** `migrate(catalog)`（`store.go` `open()`）。v1.1 把打开点冻成 `store.Open(ctx, OpenOptions, catalog)`，并把 `WasFresh` 定义在「Open 成功、catalog apply 之前」。D-002 把 `store.Open` / `WasFresh` 列为 R2 工作。

Root / 工作区纲领：R2 = 驱动、连接池、`readyz`；R3 = compiled 台账双方言对写。当前 catalog 仍含 SQLite 专用 SQL，例如 `authsession/migration/migration.go` 的 `sqlite_master` 与 `PRAGMA table_info` / `PRAGMA foreign_key_list`（在 Apply 路径上，不是模块运行时 SQL）。R2 若按字面「postgres `Open` 即 apply 现有 catalog」，会在 sqlite-only SQL 上失败，或被迫提前改写迁移（侵入 R3）。R2 若跳过 apply，又与冻结合同里的 `Open(..., catalog)` 及 `WasFresh` 时机不一致。

合同未写下列任一者：postgres 在 R3 完成前只连 + `Ping` / `WasFresh`、对 sqlite-only catalog **fail closed**、或允许空 catalog 的探测打开。R2 实施者无法唯一选择。

**要补什么**：在附件或 D 条目写明 R2 postgres `Open` 相对 catalog 的行为（三选一，禁止默示）。

### F-003 · 模块迁移中的 `sqlite_master` / `PRAGMA` 未像 `INSERT OR IGNORE` 一样点名为 R3 改写债

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |

§1 已禁止模块按方言分支 SQL，并声明仅 store 实现与迁移对写文件可分方言。§3 点名了 `operationlog/retention.go` 的 `INSERT OR IGNORE`。同属 compiled Apply、且不能靠 `?` rebind 解决的还有 `authsession/migration/migration.go`（`sqlite_master`、`PRAGMA`）。不点名则 R3 清单会漏。不阻断 R2 驱动选型；应在冻结合同 §3/§6 或 R3 方案前补进已知方言 SQL 债。

### F-004 · 表头把 `db.path` 称作「数据根」，语义实际是 `filepath.Dir(path)`

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |

§5 表格写 postgres 下 `db.path`「仍作数据目录根」，正文第 3/5 条正确写为 `filepath.Dir(db.path)`。缺省 `./data/schema-ui.db` → Dir = `./data`，与现行 composition 一致。若运维把 postgres 的 path 配成目录（如 `./data`），`Dir` 变成 `.`，文件根错位。D-002 已否决新键 `db.data_dir`。建议补一句：**postgres 下 path 仍必须是文件路径形状**（与 sqlite 缺省同形），禁止配成目录。

### F-005 · `WasFresh` postgres「用户表」未排除视图/序列；LIKE 大小写与 INTEGER 0/1 布尔未冻

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |

`WasFresh` 的 sqlite 侧是 `sqlite_master` **type = 'table'**。postgres 只写「非系统 catalog 的用户表」，未钉 `BASE TABLE` vs 视图/序列。R2 实现宜写 `information_schema.tables` / `pg_class relkind = 'r'`，只计基表。

另：模块已用 `LIKE`（wallet / recyclebin）与 `boolInt` 0/1。SQLite `LIKE` ASCII 大小写不敏感、PG `LIKE` 敏感；布尔列 INTEGER vs BOOLEAN 属 R3 物理 SQL。不升 required；R3 方案应有一句，避免当成「只 rebind 占位符即可移植」。

## 必改项汇总

1. **F-001**：修正时间单位合同，使之可核对地覆盖现行 **秒与毫秒 INTEGER 列**（或列出每列单位）；禁止「一律 Unix 秒 / 亚秒改 RFC3339 文本」作为现行默认。A-003 对 A-002 F-001 时间半段的 `fixed` 不能再当闭合证据。
2. **F-002**：写明 R2 postgres `Open` 在 R3 台账对写完成前如何对待 catalog apply（只连+Ping / fail closed / 其它有界行为）。

## 与既有意见的异同

| 项 | A-001 self | A-002 independent | A-003 self 响应 | A-004 independent（本条） |
|----|------------|-------------------|-----------------|---------------------------|
| 对象合同 | v1.0 | v1.0 | 宣称 v1.1 闭合 A-002 | **v1.1 文本 + 现行代码** |
| verdict | pass | conditional | pass（响应范围） | **conditional** |
| Tx 骨架 / 配置键 / 无 ORM / path 文件根 | 同意 | 骨架同意；path 为必改 | path 已补 | **同意 path 已补**；配置面可作 R2 映射 |
| upsert | 未单列 | 必改 F-001 | 附件已写允许/禁止形态 | **同意 upsert 规则可实施** |
| 时间类型 | 未覆盖 | 与 upsert 同条必改 | 写成「现行 = Unix 秒」 | **时间半段未真正 fixed**（F-001） |
| `WasFresh` / `OpenOptions` / `ErrNoRows` / `Dialect()` | 不足 | F-003～F-006 | 已写入 v1.1 | 骨架同意；WasFresh 视图/序列仅 recommended |
| `Open` vs catalog / R2·R3 | 未审 | 未单列 | 未写 | **新必改 F-002** |
| 可否无条件进 R2 | 当时允许 | 否 | 允许对照 v1.1 | **否**；先闭合 F-001 / F-002 |

无 P-004 冲突需用户在「秒 vs 毫秒」上选边才能承认事实：毫秒列已存在。用户要选的是合同写法（沿用每列单位 vs 新列才默认秒），不是否认现行单位。

## 结论 + 建议给编排器/用户的下一步

**方案冻结作为 R1 Tx + 配置 + upsert 骨架仍然合理；作为 R2 唯一实施合同仍不合理。** 主因不是 A-002 已补的 path / `WasFresh` / upsert 形态，而是：(1) 时间合同与现行 INTEGER 秒/毫秒并存冲突；(2) postgres `Open` 与 catalog apply 相对 R3 未钉死。verdict **conditional**。

建议 `/govern`：

1. 响应 A-004；按 D 条目修正附件（F-001、F-002；F-003～F-005 建议一并写）。
2. 在 F-001 / F-002 合法闭合前，**不要**把现行 v1.1 当完整 R2 实施合同开工（驱动选型仍受 Root I-002 约束，且不替代本条）。
3. 不要改本独立意见中的 `status` / `progress` 来「证明」冻结已过关。GOAL-002 已标 `done`、A-003 已标 `fixed`，均不构成时间单位或 Open 时序已经正确。

## 声明

本意见 `source: independent`，不修改目标 `status` / 检查点 / 派生 `progress` / 方案正文 / goal-tree。响应与放行由 `/govern` 处理。
