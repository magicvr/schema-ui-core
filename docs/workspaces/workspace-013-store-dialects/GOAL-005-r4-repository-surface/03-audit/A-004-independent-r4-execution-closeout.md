---
id: A-004
doc: audit-entry
record_id: A-004
goal: GOAL-005-r4-repository-surface
source: independent
scope: R4 实施（S0 扫描；S1 kernel.Store/kernel.Tx 接缝；S2 全模块仓库迁移 + 运行时 SQL 债改写 INSERT OR IGNORE→ON CONFLICT DO NOTHING、LIKE/COLLATE NOCASE→LOWER 等价；S3 jobs/handler/auth/composition D 回调链收口；S4 composition 公共面 kernel.Store + postgres DSN 完整启动 + readyz 门禁）
verdict: conditional
status: recorded
parent: GOAL-001-store-dialects
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# A-004 · R4 实施独立交叉审计（2026-08-20）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：`execution-facts` + `close-out-readiness`（compatibility/production 门禁）· GOAL-005 R4 实施 S0–S4
- **verdict**：conditional
- **工作区**：`workspace-013-store-dialects`（Root `GOAL-001-store-dialects`；canonical `docs/workspaces/workspace-013-store-dialects/`；`shared_materials_catalog: none`）

## 范围与区间

- **covered**：(1) 全仓公共面是否已无 `*sql.Tx` / 驱动类型；(2) 运行时 SQL 债改写后 sqlite/PG 行为等价（D-001 点名项 + 本审补扫到的 sqlite 专用函数）；(3) postgres 完整启动测试证据（`TestCompositionPostgresStartup` / `TestFullCatalogPostgresBootstrapIntegration`）；(4) sqlite 全量回归；(5) 未引 ORM / 未改 Profile/Compose 默认。对照 GOAL-005 成功标准 1–5、D-001、S0–S4 路线图、self A-001～A-003。
- **excluded**：R5 升级策略 / 备份合同（Root I-001 / I-004）；`kernel.ErrNoRows` 独立 sentinel 与 `store.Store.WithTx` 删除（D-001 明确 R5 前）；HTTP 业务 handler 的领域行为；不改 `status` / `progress` / 方案正文 / goal-tree；不读取或比较其他工作区。
- **P-005**：GOAL-005 I-001 `00-meta` 标 **verified**（S0 `*sql.Tx` 四类清单）；I-002 **open**（required，最晚 S2 每处前）；I-003 non-blocking `collecting`（S4 运维面）。本审 scope 含 S2 完成主张与 close-out-readiness，故 **到期未关的 I-002 与本审新发现的 sqlite 专用 `instr()` 影响本 scope**。I-003 不阻断关门。
- **资料引用**：无（`shared_materials_catalog: none`）。
- **对照 self**：A-001（S0/S1/S2 首批 `pass`）、A-002（S2/S3 `conditional`，F-001 required + F-002 recommended）、A-003（响应 A-002，`pass`，声称 F-001/F-002 `fixed` 且成功标准 1–4 已满足）。

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 工作区绑定合格 | `workspace.md`：id / Root / canonical / `primary_plan: VP-013-store-dialects` 一致；资料目录 `none` |
| S1 接缝：模块 `TxRunner.Run(ctx, func(kernel.Tx))` | 各模块 `store/repository.go`、`authsession/repository.go`、`jobs/repository.go`、`settings`、`operationlog`、`systemdata`：接口均为 `Run(context.Context, func(kernel.Tx) error)` |
| S3 D 链：jobs/auth/handler/composition 公共回调为 `kernel.Tx` | `jobs.CommitFunc` / `CompleteWithCommit`；`auth.ServiceCredentialUseTxRecorder`；`operationlog.RecordOperationTx(kernel.Tx, …)`；`composition.newMuxWithExtraProviders` 传入 `func(tx kernel.Tx, …)` |
| 生产公共面无 `func(*sql.Tx)` | 本审 ripgrep `func(tx *sql.Tx)`：命中仅测试 helper + `store.Store.WithTx`（sqlite 实现，D-001 留 R5）。handler / jobs / 模块生产文件 **0** 处 |
| 公共面无具体驱动类型 | `modernc.org/sqlite` 仅 `store/store.go`（+ `auth_test.go` blank）；`pgx/v5/stdlib` 仅 `store/postgres.go`。handler/jobs/modules 生产包不 import 驱动 |
| S4：`openStore` 返回 `kernel.Store`，不再 type-assert `*store.Store` | `composition.go` `openStore` → `store.Open(...)`；`newJobRuntime` / `newAuthSessionRepository` / `newMux` / `registerLifecycle` / `handler.Register*` / `systemmonitoring.Provider` 均消费 `kernel.Store` |
| 点名 SQL 债：`INSERT OR IGNORE` 已清 | 全仓 `INSERT OR IGNORE` **0** 命中。`operationlog/retention.go` 为 `INSERT … ON CONFLICT (…) DO NOTHING` |
| 点名 SQL 债：运行时 `LIKE` / `COLLATE NOCASE` 已按 D-001 改写 | wallet/recyclebin：`LOWER(col) LIKE LOWER(?)`；`usersSortSQL` / `rolesSortSQL`：`ORDER BY LOWER(col)`（不再 `COLLATE NOCASE`）。运行时 `COLLATE NOCASE` 仅注释 + 迁移 sqlite DDL |
| postgres 完整启动（本审独立复跑，非 skip） | `SCHEMA_UI_R2_PG_DSN` → 容器 `r2-pg-probe` postgres:17-alpine。`TestCompositionPostgresStartup` **PASS**（0.97s；`NewApp.Start`：catalog apply + admin 种子 + reconcile + 模块 Start/Ready，随后 `gate.setReady`）。`TestFullCatalogPostgresBootstrapIntegration` **PASS**（0.81s；含 logincaptcha 仓库在 bootstrapped PG 上 SetEnabled/Create/Consume） |
| sqlite 全量回归 + 构建（本审独立复跑） | `apps/api`：`go build ./...` PASS；`go test -count=1 ./...` **0 FAIL**（DSN 已设，含 live PG；handler 151s 等均 ok） |
| 未引 ORM；未改 Profile/Compose 默认 | `go.mod` 无 gorm/ent/sqlx/bun。`git diff 299c7dc^..HEAD -- compose.yaml apps/api/configs/config.yaml apps/api/internal/config/config.default.yaml` **空**。`config.default.yaml` `profile: mvp`、`dialect: sqlite`；`compose.yaml` 不设 `DB_DIALECT`/`DB_DSN` |

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| 1. `apps/api` 缺省 sqlite 全量回归 0 FAIL | **满足** | 本审 `go test -count=1 ./...` 0 FAIL |
| 2. handler / jobs / 模块公共契约不再 import 具体驱动、不再出现 `*sql.Tx` | **满足**（R5 残留见 F-004） | grep 生产公共签名 0 命中；`store.WithTx(*sql.Tx)` 按 D-001 留 R5 |
| 3. 运行时方言债已改写且行为等价（大小写/布尔/upsert 与 sqlite 现状一致或书面落盘差异） | **部分** | 点名三项已改写；`LOWER(…) LIKE LOWER(?)` 在 live PG 上可执行（本审 `psql` 断言 `t`）。**sqlite 专用 `instr()` 仍在 9 处运行时检索**，live PG `function instr does not exist`（F-001）。I-002 台账仍 `open`（F-002） |
| 4. postgres DSN 完整启动：compiled catalog apply + `readyz` 模块门禁全绿 | **满足（Start+Ready 等价路径）** | 本审 `TestCompositionPostgresStartup` PASS。测试未 HTTP GET `/readyz`，但 `registerLifecycle` 在 `runtime.Start`+`Ready`（含 `Ping` + `SystemDataReady`）成功后 `gate.setReady()`，与 `readyz(st, ready)` 同一门（F-003） |
| 5. 未引 ORM；未改 Profile/Compose 默认；未重做 R3 | **满足** | `go.mod`；compose/config diff 空；迁移台账未在本目标重写 |

## Findings

### F-001 · 运行时检索仍用 sqlite 专用 `instr()`，postgres 生产路径会在带 `q` 的列表查询上失败（required）

| 字段 | 值 |
|------|-----|
| 严重度 | high |
| 建议 | required |
| 状态 | open |
| evidence | 运行时 9 处：`authsession/users_repository.go:417`（`ListUsers`）；`authsession/roles_repository.go:424`（`ListRoles`）；`authsession/notifications_repository.go:146`；`datadictionary/store/repository.go:93` / `:251`；`scheduledtasks/store/repository.go:86` / `:300`；`operationlog/repository.go:317`；`wallet/store/repository.go:601`（`ListEntries`）。本审对 `r2-pg-probe`：`SELECT instr('Hello','e')` → `ERROR: function instr(unknown, unknown) does not exist`。`TestCompositionPostgresStartup` / catalog bootstrap **不调用**这些检索，故 S4 全绿不能覆盖本缺陷 |
| closure | 无 |
| 关联 I-00N | I-001（S0「方言 SQL 全量清单」未收入）；I-002（S2 运行时 SQL 债；点名集未含 `instr`） |

描述：PostgreSQL 无 `instr()`（有 `strpos` / `position`）。用户/角色/通知/字典/定时任务/操作日志/钱包流水在 **非空搜索词** 下会 500。这不是 Unicode 大小写残余，而是 SQL 无法在生产方言解析。

S0 把 I-001 闭合为 144 处 `*sql.Tx` 四类清单，未把 `instr()` 记入方言 SQL 债。D-001 / A-002 F-001 点名的是 `LIKE` / `COLLATE NOCASE` / `INSERT OR IGNORE`，这三项本审确认已改。**不否定** kernel.Tx 收口与点名改写；**阻断** R4 无条件关门 / 把「运行时方言债已改写且行为等价」写成已满足。

可移植改写（供 `/govern` 实施，本审不改代码）：与 wallet/recyclebin 已采用的 `LOWER(col) LIKE …` 对齐，或对用户输入做 `%`/`_` 转义后使用 `LIKE`；避免引入仅 PG 的 `strpos`。改写后须在 sqlite **与** live PG 上覆盖带 `q` 的列表。

### F-002 · I-002 仍 open：点名改写已发生，台账未闭合；S0 方言清单不完整（required）

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | required |
| 状态 | open |
| evidence | `00-meta.md` / `01-decision.md` I-002 status=`open`、证据「待确认」、最晚「S2 每处前」。`03-audit.md` 信息就绪表仍写 I-001 collecting、I-002 open。P-005：影响方案/实施/关门的 required 信息项须在对应阶段前由证据关闭。本审可从代码还原点名项结论（见下），但治理台账未 `verified`，且未登记 `instr()` |
| closure | 无（须 `/govern` 把 I-002 标 `verified` 并纳入 F-001，或用户书面 residual） |
| 关联 I-00N | I-002（required，S2）；旁注 I-001 00-meta=verified 但清单不完整 |

描述：P-005 不允许带着到期 required 信息项做本阶段关门。S2 点名改写**事实成立**，I-002 仍是「待确认」。本审还原的点名结论（可供编排器闭合，**不能**当作 00-meta 已更新）：

- **`INSERT OR IGNORE` → `INSERT … ON CONFLICT (…) DO NOTHING`**：`operationlog/retention.go` 三处 archive；全仓 0 残留。
- **`LIKE` → `LOWER(col) LIKE LOWER(?)`**：wallet `ListAccounts`、recyclebin `List`（未选 D-001 的 `ILIKE`）。
- **`ORDER BY … COLLATE NOCASE` → `LOWER(col)`**：`usersSortSQL` / `rolesSortSQL`（未依赖查询侧 CITEXT）。
- **插入取 id / `RETURNING`**：运行时无 `LastInsertId`；无 `RETURNING`；插入走显式 id。
- **布尔**：模块运行时 SQL 无 PG `BOOLEAN`；仍 INTEGER 0/1。

**额外**：I-002 点名集未含 `instr()`；闭合 I-002 时必须把 F-001 登记为仍开放的运行时债，或与 F-001 一并改写后 verified。

### F-003 · S4 / 大小写等价的运行证据偏窄（recommended）

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | recommended |
| 状态 | open |
| evidence | `postgres_startup_test.go` 只 `app.Start`，注释写「readyz-equivalent」，无 HTTP `GET /readyz`。wallet/recyclebin/authsession 无 `LOWER`/`q` 大小写对测；`datadictionary/store` 无测试文件。本审 live PG 证明 `LOWER('Hello') LIKE LOWER('%ELL%')` = `t`，但未跑用户/钱包检索的 PG 端到端 |
| closure | 无 |

描述：Start+Ready 与 `readyz` 共用 Ping + 模块门 + `gate.setReady`，本审**不**主张 readyz 逻辑未接线。缺口是：(1) 字面「readyz 全绿」无 HTTP 探针；(2) 点名 `LIKE`/`COLLATE` 改写没有「`Foo` 能搜到 `foo`」的双方言测试。不构成 F-001 以外的必改；建议 S5 补一条 HTTP `/readyz` 与一条 ASCII 大小写检索/排序（sqlite + live PG）。

SQLite `LOWER` 默认仅 ASCII；PG `LOWER` 跟 locale。非 ASCII 大小写可能仍有差。若产品检索只保证 ASCII，应在 I-002 书面落盘；否则属未声明差异。

### F-004 · R5 已声明残留：`WithTx(*sql.Tx)` helper、`kernel.ErrNoRows` 别名、仓库内 `sql.Null*`（recommended）

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |
| evidence | `store/store.go:107` `WithTx(ctx, func(*sql.Tx))` 仍导出，测试仍调用。`kernel/store.go` `var ErrNoRows = sql.ErrNoRows`（注释写「直到 R4 公共面停止 import database/sql」——D-001 决定 4 把独立 sentinel 放到 **R5 前**）。jobs / wallet / recyclebin / operationlog / authsession 扫描仍用 `sql.NullString`/`NullInt64`（不出现在导出签名） |
| closure | 无（D-001 已划 R5） |

描述：不把这些记为 R4 必改。成功标准 2 的「公共契约」已收口。`database/sql` 仍被 kernel 别名与仓库内部 Null 类型拉进来，与 kernel 注释的「R4 切断」不完全一致；以 D-001 决定 4 为准。`/govern` 不必在本目标为 helper/sentinel 另开 required。

### F-005 · 过时注释仍描述 R4 未完成（recommended）

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |
| evidence | `store/store.go` `Run` 注释：「modules keep using WithTx(*sql.Tx) until R4」。`handler/health.go` `readyz` 注释：「trivial SQLite read」（实现已是 `kernel.Store.Ping`） |
| closure | 无 |

描述：行为已收口；注释停在迁移前。卫生项，不改变对本审成功标准 2/4 的判定。

## 必改项汇总（required）

1. **F-001**：把 9 处运行时 `instr(...)` 改成两方言可执行的检索（建议对齐已有 `LOWER(col) LIKE …`，并处理通配符语义），在 sqlite 与 live PG 上覆盖带 `q` 的列表。
2. **F-002**：闭合 I-002（点名项 verified + 登记/关闭 `instr()`）；I-001 清单补全或注明范围仅 `*sql.Tx`。

无 P-004 意见冲突需用户在「是否改写 `instr`」上另裁——本审建议 **fixed**（改写），不建议 residual（生产 PG 带搜索的列表会直接失败）。若用户书面 overruled / accepted-residual，须写明范围（哪些检索可以在 PG 上保持失败）与复审触发。

## 与既有意见的异同

| 项 | self | 本审（independent） |
|----|------|---------------------|
| A-001 S0/S1/S2 首批 | `pass`；6 模块 + live PG captcha | **同意**该切片。补充：S0「方言 SQL 全量」实际只扫了 `*sql.Tx`，漏 `instr()` |
| A-002 F-001 LIKE/COLLATE | required，open | **同意当时点名**；A-003 后代码侧 **fixed**（`LOWER`）。本审不重开该 finding 编号 |
| A-002 F-002 S4 | recommended | **同意已做**：`openStore` → `kernel.Store`；本审复跑 `TestCompositionPostgresStartup` PASS |
| A-003 「成功标准 1–4 已满足、0 open required」 | `pass` | **不同意无条件满足标准 3 / 0 required**。标准 1、2、5 满足；标准 4 在 Start+Ready 路径满足。标准 3 因 `instr()` 与 I-002 未闭合而不成立。A-003 对点名 LIKE/COLLATE/S4 的关闭证据 **成立**，对「全部运行时方言债」过宽 |
| 新 finding | 无 | F-001 required high（`instr`）；F-002 required med（I-002 台账）；F-003～F-005 recommended |

无 self/independent 对同一必改项一要一否的 P-004 冲突：self 未讨论 `instr()`；本审新开 F-001。

## 结论 + 建议给编排器/用户的下一步

R4 的 **kernel.Tx 公共面收口、点名 `INSERT OR IGNORE`/`LIKE`/`COLLATE NOCASE` 改写、composition postgres DSN 启动、sqlite 全量回归、未引 ORM / 未改默认** 均有本审可重复证据。A-003 把这些写成已完成 **就点名范围而言成立**。

**不能无条件关门**：生产方言上，带搜索词的用户/角色/通知/字典/任务/操作日志/钱包流水查询会因 `instr()` 失败（F-001，high required）；I-002 到期未闭合（F-002）。verdict = **conditional**。

建议 `/govern`：先改写 F-001 并补双方言检索测试 → 把 I-002 标 verified（含 `instr` 条目）→ 响应本 A-004 → 再考虑 S5 关门。可选补 HTTP `/readyz` 与大小写对测（F-003）。不要在 F-001 开放时把 GOAL-005 标 `done`。

## 声明

本意见不修改 status/progress；响应由 /govern 处理。
