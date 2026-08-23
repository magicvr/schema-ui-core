---
id: GOAL-032-w21-startup-db-identity
doc: audit-entry
record_id: A-001
source: independent
scope: GOAL-032 S5 close-out（目标整体 · Identify/Plan/Execute 合同与双方言启动路径）
verdict: conditional
audit_type: close-out
status: recorded
parent: GOAL-001-design-implementation-conformance
created: 2026-08-22
updated: 2026-08-22
version: 0.1.0
---

# A-001 · 独立交叉审计 · W21 启动身份/迁移计划关门（2026-08-22）

- **source**：independent
- **auditor**：grok-4.6（Grok Build `/audit`）
- **类型** / **scope**：close-out · GOAL-032 全目标（S1 合同 + S2/S3 实现 + S4 回归证据）；migration/data 启动门禁
- **verdict**：**conditional**
- **完整意见**：本文件（未超 32 KiB，无附件）

## 范围与区间

| 项 | 值 |
|----|-----|
| 工作区 | `workspace-010-design-implementation-conformance` |
| Root | `GOAL-001-design-implementation-conformance` |
| canonical | `docs/workspaces/workspace-010-design-implementation-conformance/` |
| 资料目录 | `none`（无共享资料引用可核） |
| 日期 | 2026-08-22 |
| covered | D-001 身份/计划合同；`identity.go` classify/plan；sqlite/postgres Execute；E-001/E-002 主张；I-001/I-002；S4 单测与本轮复跑 |
| excluded | 未重跑活进程 `go run ./cmd/server` + `/healthz`（E-001 活测本轮未复验）；未读取其他工作区上下文（`00-meta` 对 store 方言权威的 Q2 路径仅作范围声明） |
| 信息项 | I-001 / I-002 在 S1 标为 verified；本意见不把它们重新标成开放 I-00N，但对 I-001「四表 = 完整」是否撑得起 restore **整表盖章** 提出 required finding（F-001） |

工作区绑定：`workspace.md` 的 `id` / `root_goal` / `canonical_scope` 与本目标路径一致。本波未改 Profile / 模块矩阵 / Manifest。

## 成果（有证据）

对照 D-001 与 as-built，下列主张**可核对**：

| 主张 | 证据 |
|------|------|
| 沿用 `schema_migrations` 当历史表，不新建第二张表 | `identity.go` ledger DDL；D-001「沿用」；I-002 verified |
| Identify → Plan → Execute 接入双方言 runner | `store/migrate.go` `Store.migrate`；`store/postgres.go` `postgres.migrate`；`open.go` 注释与 postgres 非空 catalog 走 migrate |
| 身份枚举与计划动作与 D-001 同名 | `identity.go` `dbIdentityKind` / `startupAction`；`TestClassifyIdentity` / `TestPlanStartup` **PASS**（本轮） |
| 外库 refuse，错误含 `identity=foreign` | `planStartup` foreign 分支；`TestPostgresMigrateRejectsConflictingUsersTable` **PASS**（本轮，未 skip） |
| ledger-less R2 `{users, refresh_tokens}` adopt | `TestMigrateExistingR2DB` **PASS**；`TestPostgresMigrateAdoptsLedgerlessR2` **PASS** |
| 当前 catalog 全量库丢 ledger 后 restore 盖章、不 42P07、保留用户行 | `TestPostgresMigrateRestoresLostLedger` **PASS**（先 bootstrap 全 catalog，再 `DROP TABLE schema_migrations`，再 Open） |
| 空 ledger 表 fail closed（与旧 runner 一致） | `appliedMigrations` / `appliedMigrationsPG`：存在且 0 行 → error `partial bootstrap`，probe 在 classify 前失败 |
| 未钉死 sqlite | `dev.cmd` 不设置 `DB_DIALECT`；注释写明跟随 `configs/.env` |
| 本轮回归（不含活服务器） | `go test ./internal/store/ ./internal/kernel/` **ok**；定向身份/PG 用例见上 |

`kernel.IsDuplicateObject` / `ExecIdempotentDDL` 已落地且单测绿；生产 Apply 路径**未调用**它们（见 F-007）。postgres 内联 `stampRemainingPG` 已不在 `postgres.go`（改由 `actionRestoreLedger`）。

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| S1 方案冻结（EF 对照 + 身份 + 计划） | 完成 | D-001 accepted；I-001/I-002 有书面答案 |
| S2 Identify + Plan 纯函数 + 单测 | 完成 | `identity.go` + `identity_test.go` 本轮 PASS |
| S3 双方言 runner 按计划执行 | **部分** | fresh / ledger pending / r2 adopt / foreign refuse / **当前全量** restore 已接线；`ours-partial-no-ledger` 与「完整」启发式见 F-001～F-003 |
| S4 回归 | **部分** | 声明的三条 PG 测试 + classify 本轮绿；sqlite 无 restore/foreign 集成（F-004）；无「四表但不齐 catalog」用例（F-001） |
| S5 independent 关门 | 本次 | 本条 **conditional**；有未闭合 high/med required → **不可无条件关门** |
| 不误伤外库 | 达成（已测路径） | PG 冲突 `users` refuse；classify `foreign other` |
| 已健康库不重放 DDL | 达成（ledger 在、或当前全量丢 ledger） | ours-ledger pending/noop；全量 restore 不跑 Apply。**不**覆盖「有 jobs 但缺 v43+ 对象」（F-001） |
| sqlite + postgres 同一套身份/计划枚举 | 枚举达成；**执行未对齐** | F-003 |
| 不改 Profile / 模块矩阵 / Manifest | 成立 | 仅 store 启动路径与 authsession v1 PG adopt |

## Findings

### F-001 · restore-ledger 四表启发式整表盖章，会静默跳过 jobs 之后的 catalog

- 严重度：**high**
- 建议：**required**
- 关联：I-001（S1 已 verified 的「完整最小表集合」）；D-001 `restore-ledger`
- 状态：**open**

D-001：`ours-complete-no-ledger` → **不重放 Apply**，对**当前 catalog 整表盖章**。判定完整的实现是：

```89:91:apps/api/internal/store/identity.go
func lostLedgerLooksComplete(tables []string) bool {
	have := tableNameSet(contentTables(tables))
	return have["users"] && have["refresh_tokens"] && have["operation_log"] && have["jobs"]
}
```

`jobs` 由 `core.jobs` **version 42**（`async_jobs`）创建。本轮 `TestMigrateFreshDB` 断言 compiled catalog 为 **1..48**，v42 之后至少还有：43 `operation_log_wallet_jobs`、**44 `service_credentials`**、45、46、47 `operation_log_archive`、48 `operation_log_session`。

`stampCatalog` 按传入 catalog **逐条 INSERT**，不检查这些对象是否存在。之后再启动走 `ours-ledger` + `noop`，缺口被 checksum 盖住，不会再 pending。

`TestPostgresMigrateRestoresLostLedger` 只覆盖「先用**当前** catalog 全量 bootstrap，再丢 ledger」——此时四表与 v43–v48 对象同时存在，测不出该洞。

I-001 在 catalog 头已是 v48 的同一天把「完整」答成这四张表。该答案可以识别**现场那次**全量库丢 ledger，但**撑不起**「盖章 = 当前 catalog 已全部 Apply」的执行语义。丢 ledger 的旧完整库（已有 jobs、尚未有 service_credentials）会以成功启动的形式少跑安全相关 DDL。

建议（任一即可，需决策留痕）：restore 前核验当前 catalog 应有对象是否都在，缺则 fail closed 或只盖章已存在版本再 pending；或把「完整」指纹改成与 catalog 头绑定的对象集，而不是固定四表。

### F-002 · `adopt-then-pending` 对已有 post-v1 对象非幂等，兑现不了 D-001 的 partial 计划

- 严重度：**med**
- 建议：**required**
- 关联：D-001 `ours-partial-no-ledger` / `adopt-then-pending`
- 状态：**open**

双方言 Execute 对 `adopt-r2` 与 `adopt-then-pending` 同一条臂：先 Apply `catalog[0]`（v1），再 `applyPending` 其余版本（`migrate.go` 59–70 行；`postgres.go` 101–112 行）。

v1 PG 只补 ledger / `refresh_tokens`（`migrateBaselinePG` 探针）。v2+（如 `migrateRBACPG`）仍是裸 `CREATE TABLE`。无 ledger 且已有 `roles` / `operation_log` 等、但尚未凑齐四表时，分类为 `ours-partial-no-ledger`，pending 会在已有对象上 42P07 失败。因在独立迁移事务里失败，表现为 **fail closed**（不会继续重建 `operation_log`），但 D-001 写的是「v1 只补缺、再 pending」，不是「partial 直接失败」。

分类器不区分「多出来的是我们 catalog 表」还是「无关表」。因此该一等身份只有在「我方 `users`、且没有任何 post-v1 catalog 表」时才可能跑完。

### F-003 · sqlite `ours-partial` 仍受 v1 精确二表指纹约束，与 D-001 / postgres 分叉

- 严重度：**med**
- 建议：**required**
- 关联：D-001；既有 V-MIG-03 `TestMigrateFailClosedPartialBaseline`
- 状态：**open**

D-001：无 ledger、我方 `users`、非 R2 精确集、非 complete → `adopt-then-pending`；并写明 sqlite 与 postgres **同一套**身份/计划枚举。

sqlite v1 `fingerprintR2` 仍要求表集合**恰好** `{users, refresh_tokens}`。因此：

- sqlite **users-only**（我方形态）：classify → partial → 计划 adopt-then-pending → v1 指纹失败。`TestMigrateFailClosedPartialBaseline` 本轮 **PASS**，明确要求此路径 fail closed 且不留下 ledger。
- postgres **users-only**：`migrateBaselinePG` 会建 ledger 并补 `refresh_tokens`（与 sqlite 相反）。无对等 PG 单测。

D-001 未记录对 V-MIG-03 的废止或「sqlite 保持 fail closed」例外。枚举对齐、执行未对齐。关门前必须二选一并留痕：改 sqlite v1 以兑现 adopt-then-pending，或收窄 D-001（partial 仅 PG / sqlite 维持 V-MIG-03）并改测试与文档。

### F-004 · sqlite 缺少 restore-ledger / foreign-refuse 集成测试

- 严重度：**low**
- 建议：**recommended**
- 状态：**open**

S4 证据列的是 PG 三条 + classify 单测。sqlite 有 R2 adopt 与 users-only fail-closed，没有「丢 ledger 的完整 sqlite 文件」或「orders 表 foreign refuse」的 Open 级测试。纯函数覆盖了 kind/action，未覆盖 sqlite probe + Execute。

### F-005 · classify 单测矩阵未锁若干边界

- 严重度：**low**
- 建议：**recommended**
- 状态：**open**

未覆盖：四表齐全但 `oursUsers=false` → 应为 `foreign`；`schema_migrations` 在表清单中但 `applied==nil`（空库仅系统表之外的 ledger 空表由 probe 先报错，classify 本身不处理）；ledger 行非 nil 时即使 `oursUsers=false` 仍 `ours-ledger`（D-001 以 ledger 为权威，建议测锁定，避免回归改掉）。

### F-006 · `users` 形态探针重复实现

- 严重度：**low**
- 建议：**recommended**
- 状态：**open**

`store.usersLooksLikePostgres` 与 `authsession/migration.usersLooksLikeSchemaUIPG` 列集相同（text `id/username/name/roles/password_hash`）。Identify 与 v1 Apply 双探针是防御深度，但漂移时会出现「runner 认为 ours、v1 认为不是」或相反。非阻断。

### F-007 · `ExecIdempotentDDL` 无调用点

- 严重度：**low**
- 建议：**recommended**
- 状态：**open**

E-001 记录 helper 落地。本轮检索生产调用仅为自身测试。注释已禁止在已失败的 PG 事务里用它吞 42P07（25P02）。保持未用可接受；若有人接到 v2+ Apply 上，会回到 D-001 明确拒绝的「IF NOT EXISTS / 吞 42P07 当幂等」。

## 必改项汇总

1. **F-001**：restore 不得在「仅有四表」时把当前 catalog（含 v43–v48，尤其 `service_credentials`）标成已应用；须核验对象或缩小盖章范围，并补「四表在、后继对象缺」测试。
2. **F-002**：明确 `ours-partial-no-ledger` 的可执行子集（仅 users / 仅 R2 缺口 vs 含 post-v1 表则 refuse），使计划动作与 Apply 幂等性一致。
3. **F-003**：消解 D-001 ↔ sqlite V-MIG-03 ↔ postgres users-only adopt 的合同冲突并留痕。

开放 required：**3**（F-001 high，F-002/F-003 med）。F-004～F-007 为 recommended。

## 信息门禁（P-005）

| ID | 级别 | 最晚阶段 | 台账状态 | 本轮 |
|----|------|----------|----------|------|
| I-001 | required | S1 | verified | 答案存在且被代码使用；**不足以**安全执行 stamp-all（F-001）。不把 I-001 改回 collecting（审计不改决策态） |
| I-002 | required | S1 | verified | 沿用 `schema_migrations`，与实现一致 |

无到期未答的 required 信息项被假装成事实。无用户书面 `accepted-residual`。共享资料：无。

## 与既有意见的异同

无既有 `03-audit/A-*`，无 self。本条为序列 **A-001**。不与历史 verdict 冲突。

## 结论 + 建议给编排器/用户的下一步

活故障主路径（指定 PG、外库 `users`、R2 无 ledger、**当前全量**丢 ledger）在代码与本轮 PG 测试中成立，且未钉 sqlite。这不是「没做 Identify/Plan」。

不能无条件 S5 关门：restore 的「完整」谓词弱于整表盖章（**high**），partial 计划在真实脏库上不可执行，sqlite/postgres/D-001 对 users-only 仍三岔。

建议 `/govern`：

1. 汇总本条 F-001～F-003；禁止在 required 未 `fixed` / 用户书面 `accepted-residual` / `user-overruled` 前把 GOAL-032 标 `done`。
2. 优先修 F-001（静默跳迁移），再定 partial/sqlite 合同（F-002/F-003）。
3. 修复后复跑本轮命令，并补 F-001 反例测试；可选补 F-004。
4. 不因本波改 VP-008 `go`（无 Profile/模块矩阵/Manifest 变化）。

## 声明

本意见不修改 `status` / 检查点 / 派生 `progress` / 方案正文 / goal-tree。响应由 `/govern` 处理。
