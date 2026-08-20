---
id: A-005
doc: audit-entry
record_id: A-005
goal: GOAL-004-r3-dual-dialect-ledger
source: independent
scope: R3 实施（T1 kernel.Tx 形状；T2a postgres 迁移运行器；T3 48 迁移双方言对写 + 全量 compiled catalog live PG fresh bootstrap + 台账幂等 + BIGINT/无 int 时间列合规 + store 级 open 解闸）
verdict: conditional
status: recorded
parent: GOAL-004-r3-dual-dialect-ledger
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# A-005 · R3 实施独立交叉审计（2026-08-20）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：`execution-facts` + `close-out-readiness` · GOAL-004 R3 实施（T1 / T2a / T3；用户 6 项）
- **verdict**：conditional
- **工作区**：`workspace-013-store-dialects`（Root `GOAL-001-store-dialects`；canonical `docs/workspaces/workspace-013-store-dialects/`；`shared_materials_catalog: none`）

## 范围与区间

- **covered**：(1) 全量 48 迁移在 PG 上 apply + `schema_migrations` checksum 绑 sqlite 历史；(2) 时间/金额列 BIGINT 合规（系统级无 int 时间列）；(3) PostgresApply 优先 / 回退可移植 Apply；(4) `openPostgres` 非空 catalog 解闸后行为；(5) sqlite 回归；(6) composition 层 postgres 启动路由属 **R4**（模块公共签名迁移），不构成本目标交付缺口。对照 GOAL-004 成功标准 1–5、D-001、T1–T3 路线图、self A-001～A-004。
- **excluded**：T4 作为尚未在 `00-meta` 勾选的阶段标签（本审仍核对 T4 所需双路径证据是否已存在）；R4 仓库公共面 / `WithTx(*sql.Tx)` / 运行时 `LIKE`/`INSERT OR IGNORE`/`COLLATE NOCASE` ORDER BY；R5 升级路径与备份合同；Root I-001/I-004；HTTP JSON 编码层（非迁移 Apply 路径）；不改 `status` / `progress` / 方案正文 / goal-tree；不读取或比较其他工作区。
- **P-005**：GOAL-004 I-001 `00-meta` 标 **verified**（48 双写）；I-002 **open**（required，最晚 T3 对写 DDL）；I-003 `00-meta` 仍 **open** 但 D-001 已裁「分列 catalog」；I-004 non-blocking、最晚 T4。本审 scope 含 T3 完成主张与 close-out-readiness，故 **到期未关的 I-002 影响本 scope**。
- **资料引用**：无（`shared_materials_catalog: none`）。
- **对照 self**：A-001（T1 `pass`）、A-002（T2a `pass`）、A-003（T2b/T3 `conditional`，F-001 required + F-002 recommended）、A-004（响应 A-003，`pass`，声称 F-001/F-002 `fixed`）。

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 工作区绑定合格 | `workspace.md`：id / Root / canonical / `primary_plan: VP-013-store-dialects` 一致；资料目录 `none` |
| T1：`Apply`/`Reconcile` 为 `func(Tx) error`；sqlite 运行器 `sqlTx` 适配 | `kernel/contribution.go`；`store/migrate.go` `applyMigration`；模块 `**/migration/*.go` 无 `func (tx *sql.Tx)` |
| T2a：postgres 运行器镜像 sqlite 台账合同 | `store/postgres.go` `migrate` / `applyMigrationPG` / `appliedMigrationsPG`：`normalizeCatalog` → 空台账先 apply v1 → `validateApplied`（name+checksum）→ pending；一台账一行插入 `migration.Checksum` |
| checksum 绑 sqlite 历史，PG SQL 不进 digest | `kernel.MigrationChecksum(sqliteDDL, transformID)`；`applyMigrationPG` 写入同一 `Checksum`；`TestCompiledMigrationCatalogOwnership` 冻结 48 条 sqlite checksum；本审 `go test -run TestCompiledMigrationCatalogOwnership ./internal/store/` PASS |
| PostgresApply 优先 / nil 回退 Apply | `applyMigrationPG`：`apply := migration.Apply; if ApplyPostgres != nil { apply = ApplyPostgres }`。CREATE/时间/金额/CITEXT/rebuild 走 PostgresApply；可移植 ADD COLUMN / DROP / 无时间侧表（v6/10/11/13/17/35/37/38/39/40/41/46/48）回退 Apply |
| 48 迁移全覆盖 | `Version:` 在 `**/migration/migration.go` 恰好 48 条（1–48 连续）；`TestCompiledMigrationCatalogOwnership` `len(catalog)==len(want)==48` |
| authsession sqlite 方言债未抄进 PG Apply | `sqlite_master` / `PRAGMA table_info\|foreign_key_list` 仅 sqlite `migrateBaseline` 路径；`migrateBaselinePG` 只执行 `postgresBaselineDDL`。`COLLATE NOCASE` 仅 sqlite `serviceCredentialsDDL`；PG 为 `CREATE EXTENSION citext` + `CITEXT` |
| operationlog 对写（A-003 F-001 关闭核对） | 18 条带时间的 operationlog 迁移有 `ApplyPostgres`（`pgTimeDDL` / `pgRebuild` / `pgRebuildWithCorrelation`）；v41 correlation、v48 session 无时间列，可移植 Apply。本审 live：`TestFullCatalogPostgresBootstrapIntegration` 断言 `operation_log.created_at` 与 archive 时间为 `bigint` |
| store 级 open 解闸（A-003 F-002 / A-002 F-001 关闭核对） | `openPostgres`：`len(catalog)>0` → `st.migrate(catalog)`，R2 非空 fail-closed 已移除（`TestOpenPostgresFailsClosedOnNonEmptyCatalog` 已不存在）。`TestOpenPostgresAppliesNonEmptyCatalogIntegration` 本审 PASS；`WasFresh` 仍取 **migrate 前** 快照 |
| 全量 catalog live PG fresh bootstrap + 幂等 + 台账行数 | 本审：`SCHEMA_UI_R2_PG_DSN=postgres://probe:probe@127.0.0.1:5432/probe?sslmode=disable`（容器 `r2-pg-probe` postgres:17-alpine）→ `TestFullCatalogPostgresBootstrapIntegration` PASS（scratch `r3full`；re-migrate 走 `validateApplied` ⇒ 台账 checksum 与 catalog 一致；`migCount == len(catalog)`） |
| BIGINT 点名列 + 系统级无 int 时间列 | 同测试 ~20 列 `data_type=bigint`（含 wallet 金额与 operationlog 时间）+ `information_schema` 对固定 `*_at` 名单 `data_type='integer'` 计数 0。本审该测试 PASS |
| sqlite 回归 + 全量构建 | 本审独立：`apps/api` `go build ./...` PASS；`go test ./...` **0 FAIL**（DSN 已设，含 live PG 集成测试；handler 包 153s 等均 ok） |
| 未引 ORM；模块公共仓库签名未迁 | 无 gorm/ent/sqlx 作为 Store；`store.Store.WithTx(ctx, func(*sql.Tx))` 仍在（R4） |
| composition postgres 启动路由未接（R4 边界成立） | `composition.openStore` 仍 `kst.(*store.Store)`，非 sqlite 则 Close 并报 `not wired until R3`。本审按用户/E-005/A-004：**属 R4**，不记为本目标必改 |

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| 1. 缺省 sqlite 全量构建 + 测试绿 | **满足** | 本审 `go build ./...`；`go test ./...` 0 FAIL |
| 2. 同一 compiled catalog 在 PG fresh bootstrap 全量 apply；台账同名 + 同 checksum（绑 sqlite 文本） | **满足** | live `TestFullCatalogPostgresBootstrapIntegration` + 冻结 checksum 测试 + `applyMigrationPG` 插入 `migration.Checksum` |
| 3. 模块迁移 Apply/Reconcile 可在 PG 执行：无 sqlite_master/PRAGMA/COLLATE NOCASE 直抄；时间 BIGINT；命名宽度/布尔按 §3 **落盘** | **部分** | 代码与 live 断言满足技术面；I-002 台账仍 `open`/`待确认`，无一张逐列落盘表（见 F-001） |
| 4. 对 `SCHEMA_UI_R2_PG_DSN` 可执行 T2 apply / 等价 fresh bootstrap | **满足（本审独立复跑）** | 上列 live 测试 PASS（非 skip） |
| 5. 未引 ORM；未改模块公共仓库签名；未碰运行时 SQL 债 | **满足** | `WithTx(*sql.Tx)` 仍在；运行时 `LIKE`/`INSERT OR IGNORE`/`ORDER BY COLLATE NOCASE` 仍在 repository（R4） |

## Findings

### F-001 · I-002 仍 open：T3 对写已发生，但逐列落盘未闭合（required）

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | required |
| 状态 | open |
| evidence | `GOAL-004/00-meta.md` I-002 status=`open`、证据「待确认」、最晚阶段「T3 每迁移前」；`01-decision.md` 同表仍 open；D-001 §7 要求 I-001/I-002 在每迁移对写前闭合。成功标准 3 写「命名宽度/布尔已按 §3 落盘」。本审代码/测试已能回答该问题，但治理台账未闭合 |
| closure | 无（须 `/govern` 把 I-002 标 `verified` 并落盘逐列结论，或用户书面 residual） |
| 关联 I-00N | I-002（required，T3）；旁注 I-003 仍标 open 尽管 D-001 已裁 |

描述：P-005 不允许带着到期 required 信息项做本阶段关门。T3 实施事实（双写 + live BIGINT 断言）**成立**，但 I-002 仍是「待确认」。这不是「列类型未知」——独立审已从 PG DDL + 测试还原如下，可供编排器直接闭合，**不能**把本段当作 00-meta 已更新。

本审还原的 I-002 结论（证据 = 各模块 `*PGDDL` / `pgTimeDDL` + live 断言）：

- **Unix 时间 → PG `BIGINT`**：`schema_migrations.applied_at`；`users.{created_at,updated_at,locked_until}`；`refresh_tokens.{expires_at,revoked_at,created_at}`；roles/permissions/menu_items 的 created/updated；`system_data_reconcile.applied_at`；`service_credentials.{expires_at,revoked_at,last_used_at,created_at,updated_at}`；`records.updated_at`（v6 已 drop）；`operation_log.created_at`；`operation_log_archive.{created_at,archived_at}`；`site_settings.updated_at`；`notifications.{read_at,created_at}`；字典/任务/验证码/回收站/数据权限/MFA/wallet/jobs 的 `*_at` / `lease_expires_at`。
- **金额 → PG `BIGINT`**：`wallet_accounts.balance_*`；`wallet_ledger_entries.{amount_delta,balance_after_*}`。
- **布尔/开关 → 两方言 `INTEGER` 0/1**（等价形态，非 PG `BOOLEAN`）：`roles.system`、`menu_items.enabled`、`users.enabled` / `notifications_enabled` / `must_change_password`、`jobs.cancel_requested`。
- **非时间计数/版本仍 `INTEGER`**：`token_version`、`failed_login_count`、`sort_order`、jobs progress/attempt/lease_version、`mismatch_count`、`operation_log_retention_days`。

**不阻断**已完成的 T1/T2a/T3 代码使用；**阻断** GOAL-004 无条件关门 / 把 T4 标完成，直到 I-002 按上表（或等价指针）`verified`。

### F-002 · 「系统级无 int 时间列」oracle 漏掉 `locked_until`（recommended）

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | recommended |
| 状态 | open |
| evidence | `postgres_test.go` leftover 名单仅 `created_at`…`deleted_at` 等 `*_at`；**不含** `locked_until`。`postgresAccountLockDDL` 将 `users.locked_until` 写成 `BIGINT`（unix seconds，authsession 注释明示），但无 `assertPGType("users","locked_until","bigint")`，系统级计数也不会在该列残留 int4 时失败 |
| closure | 无 |
| 关联 I-00N | I-002 |

描述：本审 **不**主张当前 fresh bootstrap 上 `locked_until` 仍为 int4（PG DDL 为 BIGINT）。缺口是回归网：若有人删掉 v12 `ApplyPostgres`，leftover 检查与点名断言都不会红。建议把 `locked_until` 纳入 leftover 名单或加一条 `bigint` 断言。不构成 T3 代码必改。

### F-003 · 公开注释仍描述 R2「非空 catalog fail-closed / 直到 R3 才接线」（recommended）

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |
| evidence | `store/open.go` 包注释仍写 postgres 非空 catalog **FAILS CLOSED**（与 `openPostgres` 现行 `migrate` 矛盾）；`postgres.go` 类型注释仍写「compiled catalog is NOT applied here」；`composition.openStore` 注释仍写「postgres … fails closed（dual-dialect apply is R3）」且错误串 `not wired until R3`。`TestFullCatalogPostgresBootstrapIntegration` 仍留「operation_log is the known outstanding module」但同测试已断言其 `bigint` |
| closure | 无 |

描述：行为已解闸；文档/注释停在 R2。会误导后续 R4 实施者以为 `store.Open(postgres, compiled catalog)` 仍 fail-closed。属卫生项，不改本审对 T3 行为的判定。

### F-004 · D-001 §5 仍把 composition postgres 路由列为 T2 伴随；后文与用户边界划为 R4（recommended）

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |
| evidence | D-001 决定 4–5：解闸含 `composition.openStore` 路由。E-003 把解闸并入 T3；E-005/A-004 再把 **composition 路由移交 R4**。用户本审书面：composition 层 postgres 启动路由属 R4。代码：`openStore` 仍断言 `*store.Store` |
| closure | 无 |

描述：本审 **同意**用户/E-005/A-004：在模块公共契约仍是 `*store.Store` / `WithTx(*sql.Tx)` 时，composition 接 postgres 不是 R3 交付面。建议 `/govern` **补丁 D-001 §5**（accepted 决策与现行边界一致），避免下一轮独立审把未接线 composition 当成开放 R3 必改。这是决策-执行漂移，不是 T3 实现缺口。无 P-004 冲突需用户在「R3 是否包含 composition」上另裁——用户本轮已书面划为 R4。

## 必改项汇总

- **F-001（required）**：闭合 I-002（`verified` + 逐列结论落盘，可引用本审还原表 / live 断言 / PG DDL）。在此之前 **不得** 将 GOAL-004 标 `done`，也不得把 T4 当无条件关门。

无其他 required。F-002～F-004 为 recommended。

## 与既有意见的异同

| 条目 | self | 本独立审 |
|------|------|----------|
| A-001 T1 `pass` | 形状机械切换、sqlite 回归 | **同意**；本审无 `*sql.Tx` Apply 残留，全量测试绿 |
| A-002 T2a `pass`；F-001 recommended 生产未解闸 | 运行器 live 证明；解闸并入 T3 | **同意历史**；解闸现已在 store 级落地（见下） |
| A-003 F-001 operationlog 未对写（required） | open → A-004 `fixed` | **同意关闭**：18 条 PostgresApply + live `bigint` 断言本审复跑 PASS |
| A-003 F-002 / A-002 F-001 生产解闸 | A-004 `fixed（store 级）`；composition → R4 | **同意关闭（store 级）**；composition 未接 **按 R4 边界接受**（F-004 只要求补丁 D-001，不重开 A-003 F-002） |
| A-004「本目标无开放必改」 | T3 完成，待 T4 + independent | **不同意用于 GOAL 关门**：代码/测试面 T3 主张成立，但 I-002 到期未关（F-001）→ 本审 **conditional**。A-003/A-004 未审 I-002 台账 |
| 系统级无 int 时间列 | A-004 采信 leftover=0 | **同意当前 bootstrap 为 0**（本审复跑）；**补充** oracle 漏 `locked_until`（F-002） |
| checksum 绑 sqlite | self 主张 | **同意**，并加冻结 48 checksum + remigrate `validateApplied` 链 |
| T4 双路径证据 | A-004 列为下一步 | 本审认为 **证据已基本齐**（sqlite `go test ./...` + live 全量 PG boot）。T4 剩下的是勾选路线图、闭合 I-002、响应本意见——不是再发明一条测试路径 |

无 verdict 对撞需要 P-004 另裁（self 未把 I-002 标为 finding；本审新增，不否定其 T3 代码结论）。

## 结论 + 建议给编排器/用户的下一步

R3 **实施事实**（T1 形状、T2a 运行器、T3 48 对写、live 全量 PG boot、台账 checksum、store 级解闸、sqlite 回归）**经独立复跑成立**。A-003 两条 finding 的 `fixed` 关闭 **可核对**。composition postgres 启动路由 **正确划在 R4 外**。

因 I-002 仍为到期 required 信息项，close-out-readiness = **conditional**，不是 `pass`。不是「48 迁移没写完」。

建议 `/govern`：

1. **必做**：响应 F-001，把 I-002（及顺手 I-003）在 `00-meta`/`01-decision` 标 `verified`，逐列结论落盘（可引用本审表）。
2. **建议**：F-002 补 `locked_until` 断言或 leftover 名单；F-003 改 `open.go` / `postgres.go` / composition 过期注释与错误串（R3→R4）；F-004 补丁 D-001 §5。
3. I-002 闭合后，T4 可用本审已跑的双路径证据做 self 关门审计，再标 GOAL-004 `done`（progress 由编排器按检查点重算；本意见不改）。

## 声明

本意见 `source: independent`，不修改 `status` / `progress` / goal-tree / 方案正文；不 `git add` / `commit`。响应由 `/govern` 处理。
