---
id: A-001
doc: audit-entry
record_id: A-001
goal: GOAL-006-r5-dual-path-acceptance
source: independent
scope: R5 实施（U0 sqlite 回归 + live PG boot/启动基线；U1 SQLite→PG 升级策略 I-001；U2 PG 备份/恢复合同 I-004；U3 跨模块共事务 live 验收）+ close-out-readiness（production 门禁；VP-013 退出判据 1–6；Root 5/5）
verdict: conditional
status: recorded
parent: GOAL-001-store-dialects
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# A-001 · R5 实施独立交叉审计（2026-08-20）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：`execution-facts` + `close-out-readiness`（production 门禁）· GOAL-006 R5 U0–U3 实施 + U4/Root 5/5 关门就绪
- **verdict**：conditional
- **工作区**：`workspace-013-store-dialects`（Root `GOAL-001-store-dialects`；canonical `docs/workspaces/workspace-013-store-dialects/`；`shared_materials_catalog: none`）

## 范围与区间

- **covered**：(1) 双路径证据是否可复核（sqlite `go test ./...` 0 FAIL、PG boot、PG 完整启动、跨模块共事务 commit/rollback）；(2) Root I-001 / GOAL-006 RT-I-001 是否真正 verified（决策 + 可执行证据）；(3) Root I-004 / RT-I-004 备份合同 `pg_dump`/`pg_restore` round-trip 是否成立；(4) R5 是否可支撑工作区根目标 5/5 关门（VP-013 退出判据 1–6）。对照 GOAL-006 成功标准 1–5、D-001/D-002、E-001/E-002、U0–U4 路线图。
- **excluded**：不改 `status` / `progress` / 方案正文 / goal-tree；不读取或比较其他工作区；不重开 R1–R4 已闭合 required（GOAL-002..GOAL-005 仅作前提指针）；不把本审临时探测库当产品产物。
- **P-005**：GOAL-006 `00-meta` / `01-decision.md` 将 RT-I-001 / RT-I-004 标 **verified**（D-002，最晚 U1/U2）。权威 Root 台账 `GOAL-001/00-meta.md` 与 `GOAL-001/01-decision.md` 仍将 I-001 / I-004 标 **open**、证据「待确认」，最晚阶段「R5 开始前 / R5」。本审 scope 含 U1/U2 完成主张与 Root 关门就绪，故 **Root 到期未关的 I-001/I-004 影响本 scope**。资料目录 `none`。
- **对照 self**：本目标 `03-audit/` 在本条写入前 **无 A 条目**（无 self）。无法做逐条 finding 对照。D-001 U4 与项目独立审计路径均要求 self 之后 independent。

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 工作区绑定合格 | `workspace.md`：id / Root / canonical / `primary_plan: VP-013-store-dialects` 一致；`shared_materials_catalog: none`。Charter `schema-ui-core-admin-foundation@0.2.0` 与 VP-013 `vision_ref` 一致 |
| U0 sqlite 全量回归 0 FAIL（本审独立复跑） | `apps/api`：`SCHEMA_UI_R2_PG_DSN=postgres://probe:probe@127.0.0.1:5432/probe?sslmode=disable` → `go test -count=1 -timeout 8m ./...` **exit 0**（handler 150.6s 等均 `ok`，0 FAIL）。DSN 已设，live PG 门控测试未 skip |
| U0 PG 全量 boot（本审独立复跑） | 容器 `r2-pg-probe` postgres:17-alpine。`TestFullCatalogPostgresBootstrapIntegration` **PASS**（0.68s；scratch `r3full`；48 迁移 fresh bootstrap + 幂等 + checksum drift fail-closed） |
| U0 PG 完整启动 / readyz 等价路径（本审独立复跑） | `TestCompositionPostgresStartup` **PASS**（0.85s；`NewApp.Start`：catalog apply + 种子 + reconcile + 模块 Start+Ready，随后 `gate.setReady`）。测试未 HTTP GET `/readyz`（F-007） |
| U3 跨模块共事务（本审独立复跑） | `apps/api/internal/store/postgres_test.go` `TestPostgresCrossModuleSharedTx` **PASS**（0.91s 与 0.87s 各一次）。`CreateServiceCredential` 经 `authsession.Repository.withTx` → `kernel.Store.Run`；同事务 `operationlog.RecordOperationTx`：commit 路径 cred+op = 1/1；审计强制失败路径 0/0 |
| U1 书面策略存在 | D-002：in-place 跨引擎不可行；推荐 fresh bootstrap；逻辑数据迁移 / 模块导入列为支持路径；自动 in-place 转换不提供。fresh bootstrap 有 live 证据（上列 T3 测试） |
| U2 合同文本 + 工具链 | D-002：逻辑备份 `pg_dump -F c` → `pg_restore`；物理/PITR `pg_basebackup` 任选。SQLite `VACUUM INTO` 仍仅 `store/migrate.go` snapshot（sqlite 路径，应保留） |
| U2 原主张玩具 round-trip 可复现 | 遗留库 `r5u2` / `r5u2rest`：仅 `public.users(id text, name text, created_at bigint)`，各 2 行（u1 Alice / u2 Bob）。本审另建 `r5indsrc`→ dump → `r5inddst` restore，`count(*)=2`，随后已 DROP |
| U2 **本审** catalog 级 round-trip **成立** | 对 48 迁移 / 35 张 public 表的 bootstrapped 库 `pg_dump -F c`（66257 bytes）→ `CREATE DATABASE r5auditrest` → `pg_restore` exit 0。`schema_migrations` 48=48；台账 checksum 聚合 `md5(string_agg(version\|\|':'\|\|checksum))` 双方均为 `4f21236ecaeb6af0dfe1b2194d9f157f`。遗留库 `r5auditrest` 仍在 probe 容器（独立证据，非产品代码） |
| 未引 ORM；未改缺省 sqlite | `apps/api/go.mod` 无 gorm/ent/sqlx/bun；有 `pgx/v5` + `modernc.org/sqlite`。`config.yaml` / `config.default.yaml`：`dialect: sqlite`。`compose.yaml` 无 `DB_DIALECT`/`DB_DSN` |
| 生产公共面 `func(*sql.Tx)` | 本审 ripgrep：生产导出签名仅 `store.Store.WithTx`（sqlite 实现；R4 D-001 划「R5 前」）。handler / jobs / 模块生产文件 0 处 |
| 运行时 `instr()` | 本审仅命中 `postgres_test.go` 注释；R4 A-005 改写仍在 |

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| 1. sqlite 缺省路径全量回归 0 FAIL；PG 双路径（bootstrap + 完整启动）live 全绿 | **满足** | 本审 `go test ./...` 0 FAIL；两门控测试 PASS（非 skip） |
| 2. I-001：SQLite→PG 升级策略有书面结论与证据（in-place 不可行则 dump/restore/rebootstrap + 范围） | **部分** | 书面结论有（D-002）；rebootstrap 有 live 证据。**SQLite→PG 数据搬运无抽样原型**；D-002 把「逻辑数据迁移」写成支持路径，但工具「由 fork/运维选型」。Root I-001 仍 `open` /「待确认」（F-001、F-002） |
| 3. I-004：PG 备份/恢复合同明确（替换 `VACUUM INTO`），有可执行验证 | **实质满足 / 台账未关** | 合同明确。原落盘证据是玩具 `users` 表（非台账）。**本审**对 48 迁移 catalog 的 dump/restore 成立。Root I-004 仍 `open`（F-001）；原「账台账计数」核对面与实跑不一致（F-004） |
| 4. 跨模块共事务在 PG 上验收通过；`readyz` 生产向就绪 | **满足（Start+Ready 等价路径）** | U3 测试 PASS。`readyz` 实现为 `kernel.Store.Ping` + 可选 `ready()` 门；composition 启动走同一门。无 HTTP GET `/readyz` 探针（F-007） |
| 5. VP-013 退出判据 1–6 全链可核对；无 open required；Root 5/5 关门 | **不满足** | 判据 1/3/5 可核对；2/4 策略与工具证据存在但 Root 信息项未关；6 = 本审开放 required ≠ 0。U4 无 self。 **不同意 Root 5/5 关门** |

### VP-013 退出判据（close-out）

| # | 判据 | 本审 |
|---|------|------|
| 1 | 内核端口落地；公共契约无 `*sql.Tx` / 驱动类型 | **满足**（R4 已收口；`WithTx` 为 sqlite 实现残留，见 F-006） |
| 2 | PG 可对开区 compiled 迁移 fresh bootstrap；既有 SQLite 升级路径证据 **或** 书面有界 residual | **条件满足**：fresh bootstrap 本审 PASS。升级路径被写成 dump/restore，但无 SQLite→PG 原型；残余未按 P-005 用户书面接受字段收口；Root I-001 未关 |
| 3 | SQLite 默认仍可用；两方言逻辑 schema 一致；新迁移双 apply + checksum | **满足**（sqlite 全量绿；R3 双方言台账 + 本审 PG boot） |
| 4 | 生产向验收以 PG 为准（迁移、`readyz`、跨模块共事务、备份/恢复合同之一可核对） | **实质满足**：四项均有本审可核对证据。Root I-004 台账未关 |
| 5 | 未引入 ORM；未改 Charter；未进 Admin 功能/业务域 | **满足** |
| 6 | 开放 required finding = 0（或已合法闭合） | **不满足**：本条 F-001～F-003 仍 open |

## Findings

### F-001 · Root I-001 / I-004 仍 open，「verified」只写在 GOAL-006 旅程表（required）

| 字段 | 值 |
|------|----|
| 严重度 | high |
| 建议 | required |
| 状态 | open |
| evidence | `GOAL-001-store-dialects/00-meta.md` I-001/I-004 = **open**，证据「待确认」，最晚「R5 开始前 / R5」。同文件 `01-decision.md` 信息表同为 open。`GOAL-001/03-audit.md` 信息就绪仍写 I-001/I-004 open。GOAL-006 `00-meta.md` / `01-decision.md` 将 RT-I-001/RT-I-004 标 verified（D-002）。GOAL-006 `03-audit.md` 索引（本条写入前）仍写「open（U1/U2 门禁）」。`01-decision.md` 决策索引 **未列入 D-002**（文件存在于 `01-decision/D-002-upgrade-backup-contract.md`） |
| closure | 无 |
| 关联 I-00N | Root I-001（required，R5 / 退出 2）；Root I-004（required，R5 / 退出 4） |

描述：P-005 以 **Root 信息项** 为工作区根目标关门门禁。子目标旅程登记不能替代 Root 台账。VP-013 退出判据 2/4/6 与 Root 成功标准 4 都挂在 I-001/I-004 上。当前出现「GOAL-006 已 verified、Root 仍待确认」双账。E-013 已写明 R5 是这两项的实施窗口，但 Root 从未被回写。

**阻断** Root 5/5 与 VP-013 退出判据 6。不否定 D-002 文本与本审复跑到的 bootstrap/备份工具证据。

`/govern` 应在响应本条时：把 D-002 列入 GOAL-006 决策索引；按 F-002 的实质结论回写 Root I-001/I-004（`verified` 或 `accepted-residual` + 范围/复审触发）；同步 Root `00-meta` / `01-decision` / `03-audit` 信息表。本审不改这些文件。

### F-002 · I-001 把未实证的「SQLite→PG 逻辑数据迁移」写成支持路径（required）

| 字段 | 值 |
|------|----|
| 严重度 | med |
| 建议 | required |
| 状态 | open |
| evidence | D-001：U1「落盘 Decision + **抽样原型证据**」；不足则以书面 residual 收口。D-002 §U1.2 将 **逻辑数据迁移（dump/restore 或模块级导出/导入）** 列为支持路径；§U1.4 残余只覆盖「自动 in-place 不提供」，并写「dump/restore 或重放工具由 fork/运维选型」。E-002 未给出任何 SQLite 表导出并导入 PG 的命令/行数/抽样。`VACUUM INTO` 仍只服务 sqlite 升级快照（`store/migrate.go`），不能产生 PG 可 restore 的档案。本审未发现 SQLite→PG 迁移器、脚本或原型测试 |
| closure | 无（须 `/govern` 二选一：抽样原型 **或** 用户书面 residual） |
| 关联 I-00N | Root I-001 / RT-I-001；VP-013 退出判据 2 |

描述：VP-013 退出判据 2 **允许**「书面记录不可自动升级、需 dump/restore 的有界 residual」。fresh bootstrap 本审已证，in-place 跨引擎不可行也成立。缺口是：GOAL-006 把 I-001 标 **verified**（暗示支持路径已落地），同时把「逻辑数据迁移」写成产品支持项，却 **零抽样**。pg_dump 是 **PG→PG**（I-004），不能冒充 SQLite→PG。

建议闭合（本审不改决策正文）：

1. **preferred residual**：明确「本 VP **不提供** SQLite→PG 数据搬运器；存量路径 = fresh bootstrap + 运维自备导出/回放」。范围、复审触发、用户书面接受（P-005 `accepted-residual`）。然后 I-001 可关。
2. **fixed**：给一份最小原型（例如从 sqlite `users`/`schema_migrations` 抽样导入已 bootstrap 的 PG）并落盘命令与行数。

不要在 F-002 开放时把 I-001 继续写成无条件 verified，也不要静默当作已关。

### F-003 · U4 缺 self；生产关门链不完整（required）

| 字段 | 值 |
|------|----|
| 严重度 | med |
| 建议 | required |
| 状态 | open |
| evidence | D-001 决定 5：U4 = VP-013 退出 1–6 核对 + **self + independent**（production 门禁）。项目级 `docs/architecture/independent-audit-execution.md`：先 self，再 grok-build `/audit`。GOAL-006 `03-audit/` 在本条前为空；`00-meta` U4 仍 🔄。本条为 A-001 independent，不能冒充 self |
| closure | 无 |
| 关联 I-00N | 无（过程门禁）；影响 U4 / 退出 6 |

描述：交叉审计可以在 self 缺失时仍然出意见（本条即是）。**不能**据此把 GOAL-006 / Root 标 `done`。编排器须先有覆盖 U0–U3 + 退出 1–6 的 self，再汇总 self + 本 A-001，按 P-003 闭合 required 后才能放行 U4。

### F-004 · I-004 原落盘证据是玩具表，不是声称的「台账计数」（recommended）

| 字段 | 值 |
|------|----|
| 严重度 | med |
| 建议 | recommended |
| 状态 | open |
| evidence | D-002：「本 VP 以**账台账计数 + 代表性行计数**为核对面」；实跑 = 库 `r5u2` 单表 `users` 2 行，无 `schema_migrations`。E-002 只记 `count=2`。本审 catalog 级 dump/restore 已证明合同在真实 48 迁移库上成立（见成果表） |
| closure | 无 |

描述：**不否定** `pg_dump -F c` → `pg_restore` 合同，也 **不**要求实现应用内备份 API（R1 冻结合同的 `kernel.Store` 无 `Backup()`；sqlite `VACUUM INTO` 应继续留给文件库升级）。缺口是证据包装：用玩具表关闭「台账核对」过宽。`/govern` 把本审 catalog 级数字写入 E 条或附件即可把本条降为 hygiene。

### F-005 · 台账索引与 Root 附属项漂移（recommended）

| 字段 | 值 |
|------|----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |
| evidence | GOAL-006 `01-decision.md` 决策索引只有 D-001，漏 D-002。`goal-tree.md` GOAL-006 progress 为 `—`，`00-meta` 为 `4/5`（本审不改）。Root I-003 仍 `collecting`，E-013 已写 verified。`workspace.md` 纲领表 R2–R5 仍写「依赖 R1」未勾完成 |
| closure | 无 |

描述：不单独阻断 U0–U3 事实。Root 5/5 前应把索引与 I-003（non-blocking）一并刷到与 E-013/D-002 一致，避免下一审计再踩双账。

### F-006 · R4 声明的 R5 前残留：`WithTx(*sql.Tx)` 与 `kernel.ErrNoRows` 别名（recommended）

| 字段 | 值 |
|------|----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |
| evidence | R4 D-001 决定 1/4：「`store.Store.WithTx` … R5 前」删除或降为测试 helper；`kernel.ErrNoRows` 改为独立 sentinel（R5 前）。现状：`store/store.go:107` 仍导出 `WithTx(ctx, func(*sql.Tx))`；`kernel/store.go:25` `var ErrNoRows = sql.ErrNoRows`。GOAL-006 D-001：「R5 不改公共契约」。测试仍调用 `WithTx` |
| closure | 无 |

描述：不把这项记为 R5 必改——与 GOAL-006 D-001「不改公共契约」冲突时，以 R5 方案为准。VP-013 退出判据 1 的「handler / 模块公共契约」已收口。若 `/govern` 要在本目标消化 R4 残留，应另写决策，不要 silently 当 U4 成功标准。

### F-007 · postgres 路径无 HTTP `/readyz` 探针；注释仍写「trivial SQLite read」（recommended）

| 字段 | 值 |
|------|----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |
| evidence | `postgres_startup_test.go` 只 `app.Start`，注释「readyz gate path」。`handler/health.go:71` 注释「trivial SQLite read」；实现已是 `st.Ping` + `ready()`。`health_test.go` 的 `/readyz` 走 sqlite 测试环境 |
| closure | 无 |

描述：Start+Ready 与 `readyz` 共用 Ping + 模块门 + `gate.setReady`。本审 **不**主张 readyz 未接线。缺口是字面 HTTP 探针与过时注释。不阻断退出判据 4。

## 必改项汇总（required）

1. **F-001**：回写 Root I-001 / I-004（及 Root `01-decision` / `03-audit` 信息表）；GOAL-006 决策索引补 D-002。Root 台账仍 open 时 **禁止** Root 5/5 / VP-013 退出 6。
2. **F-002**：I-001 按「抽样原型」或「用户书面 residual（无 SQLite→PG 搬运器；fresh bootstrap + 运维自备）」三路径之一合法闭合；禁止继续用「逻辑数据迁移已支持」空主张支撑 verified。
3. **F-003**：补 self（U0–U3 + 退出 1–6），再与本 A-001 一并响应；不得仅凭 independent 关门。

无 self 可与本审形成 P-004「一要一否」冲突。F-002 的 residual / overruled 须用户书面（P-004）；本审建议 **accepted-residual**（不自研迁移器）或 **fixed**（最小原型），不建议把未实证路径写成已验证。

## 与既有意见的异同

本目标无 self / 无既有 independent。对照同工作区 R4 independent（GOAL-005 A-004，**不**把该目标状态当本条证据）：

| 项 | R4 independent A-004 / 本区前提 | 本审（GOAL-006 A-001） |
|----|----------------------------------|------------------------|
| sqlite 全量 0 FAIL | 当时 PASS | **同意并复跑**：仍 0 FAIL |
| PG boot + composition Start | 当时 PASS | **同意并复跑**：仍 PASS |
| `instr()` | A-004 F-001 required；A-005 称 fixed | **同意已固定**：运行时检索 0 处 `instr(` |
| I-001/I-004 划给 R5 | A-004 excluded | **本审范围**：策略文本有、Root 未关、SQLite→PG 数据路径无原型 |
| `WithTx` / `ErrNoRows` | A-004 F-004 recommended，划 R5 | **仍在**；本审 F-006 保持 recommended（R5 D-001 冻结契约） |
| HTTP `/readyz` | A-004 F-003 recommended | **仍缺**；F-007 |
| 本目标 self | （R4 有 self 链） | **缺失**（F-003） |

无 P-004 意见冲突。

## 结论 + 建议给编排器/用户的下一步

**verdict = conditional。不同意现在做 Root 5/5 关门。**

U0 / U3 的双路径与共事务主张 **名实相符**（本审独立复跑，非 skip）。I-004 的 `pg_dump`/`pg_restore` round-trip **成立**：原证据是玩具表，本审把 48 迁移 catalog 跑通且台账 checksum 一致。I-001 **不能**按当前写法无条件 verified：fresh bootstrap 已证，SQLite→PG 数据路径未证，Root 台账未关。

建议 `/govern`：

1. 补 GOAL-006 **self**（F-003）。
2. 按 F-002 裁决 I-001（residual 或原型）并 **回写 Root I-001/I-004**（F-001）。
3. 把 catalog 级 dump/restore 数字记入执行台账（F-004，建议）。
4. required 三路径闭合后，再核对 VP-013 退出 1–6，考虑 GOAL-006 `done` → Root 5/5。不要在 F-001～F-003 开放时把 Root 或 VP-013 标关闭。

## 声明

本意见不修改 status/progress；响应由 /govern 处理。
