---
id: GOAL-001-store-dialects
doc: audit-entry
record_id: A-001
source: independent
scope: Root close-out · 代码实现对照 VP-013 退出判据 1–6（不以治理文档为通过依据）
verdict: pass
status: recorded
parent: null
created: 2026-08-21
updated: 2026-08-21
version: 0.1.0
---

# A-001 · 独立关门审计：代码实现是否支撑 Root 闭门（2026-08-21）

- **source**：independent
- **auditor**：grok-4.6（Grok Build `/audit` · 思考强度 high）
- **类型**：close-out
- **scope**：`workspace-013-store-dialects` / `GOAL-001-store-dialects` 根目标闭门；对照 VP-013 方向级退出判据与 Root 成功标准，**以本轮代码阅读 + 本机复跑 + HEAD CI 为准**。治理五件套 / 子目标审计 / execution 台账只作范围参考，**不得作为本条通过依据**。
- **verdict**：pass

## 范围与区间

- **workspace**：`workspace-013-store-dialects`（`workspace.md` `root_goal: GOAL-001-store-dialects`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`primary_plan: VP-013-store-dialects`）
- **covered**：
  1. 内核 `kernel.Store` / `kernel.Tx` 端口与生产公共面是否仍暴露 `*sql.Tx` / 驱动类型
  2. PostgreSQL 实现（pgx v5 stdlib、配置面、`Open`/`Ping`/`WasFresh`、`readyz` 路径）
  3. compiled 台账双方言 apply + checksum fail-closed
  4. SQLite 是否仍为默认；无 PG 时开发/快测路径
  5. 生产向验收：迁移、Start+Ready/`readyz`、跨模块共事务、备份/恢复合同
  6. 无 ORM；未进 Admin/业务域；升级路径（in-place vs dump/fresh）
- **excluded**：其他工作区上下文；愿景层 VRev（应走 `/vision-audit`）；未改 Charter 的逐字审阅；对象存储 / Redis / 队列（VP 非目标）
- **今日独立证据窗口**：2026-08-21；HEAD `54722036a2d95bb201bb8002f6911ab868f8d915`

## 成果（有证据）

本轮**自己**核对，不是转述子目标 A 条。

| 主张 | 本轮证据 |
|------|----------|
| 内核端口存在且 composition 注入 `kernel.Store` | `apps/api/internal/kernel/store.go`（`Store`/`Tx`）；`apps/api/internal/store/open.go` `Open(...) (kernel.Store, error)`；`composition.openStore` 返回 `kernel.Store` |
| 模块仓库公共面是 `TxRunner.Run(..., func(kernel.Tx))`，不是 `*sql.Tx` | 全仓 `type TxRunner interface` 共 12 处，均为 `Run(context.Context, func(kernel.Tx) error)`；handler `Register*` 签名为 `st kernel.Store` |
| PostgreSQL 实现 = `database/sql` + `pgx/v5/stdlib`，无 ORM | `apps/api/go.mod` require `github.com/jackc/pgx/v5 v5.10.0` + `modernc.org/sqlite`；全仓无 gorm/ent/sqlx/AutoMigrate |
| 默认方言 SQLite；未知方言 fail-closed | `apps/api/configs/config.yaml` `db.dialect: sqlite`；`config.go` 空值归一 sqlite、未知值 `LoadError`；`compose.yaml` 仍 `DB_PATH` 文件库，无 postgres service |
| compiled catalog = 48 条、版本 1–48 连续 | 本轮 helper/测试打印 `catalog=48`；`compiled/persistence.go` 收集 14 个 migration provider |
| 全量 PG fresh bootstrap + 幂等再 migrate + Unix 时间列为 bigint | **本轮** `go test ./internal/store -run TestFullCatalogPostgresBootstrapIntegration -v` → **PASS**（1.88s，未 skip） |
| checksum drift fail-closed | `TestPostgresMigrateRunnerIntegration` **本轮 PASS** |
| 跨模块共事务 commit / rollback | **本轮** `TestPostgresCrossModuleSharedTx` **PASS**（credential + operation_log 同 tx） |
| SQLite→PG 逻辑拷贝最小原型 | **本轮** `TestPostgresDataMigrationPrototype` **PASS**（测例级，非产品 CLI） |
| PG 完整启动（catalog apply + seed + Start/Ready） | **本轮** `TestCompositionPostgresStartup` + `TestCompositionPostgresConfigDriven` **PASS** |
| `/readyz` 走 `kernel.Store.Ping` + 模块门禁 | `handler/health.go` `readyz(st kernel.Store, ready)`；composition `gate.setReady()` 在 `runtime.Start`+`Ready` 之后；auth-session Ready 含 `Ping` + `SystemDataReady` |
| SQLite 回归（本机） | 本轮 `go test ./internal/config,kernel,store,composition,handler` 全绿（handler 115.9s，默认 sqlite 路径） |
| CI 双路径（HEAD） | GitHub `r6-basic-matrix` run `32432792181` conclusion **success**；job `api + postgres` `go test ./...` success（`PG_TEST_*` + `postgres:15-alpine`） |
| 备份合同 `pg_dump -F c` → `pg_restore` **本轮独立 round-trip** | 本轮对 fresh-bootstrap 库 `r13indaudit`：48 迁移 / 35 张 public 表；checksum 聚合 `4f21236ecaeb6af0dfe1b2194d9f157f` 与 restore 库 `r13indrest` **一致**；dump 65824 bytes；服务器 `15.4`；客户端 `pg_dump 17.11` |
| 生产代码无 `instr(` | 全仓 `instr(` 仅测试注释（`postgres_test.go`） |

## 对照成功标准

VP-013 方向级退出判据 / Root `00-meta` 成功标准（标准文本作对照清单；**判定只引用上表本轮证据**）。

| # | 标准 | 本轮判定 | 依据 |
|---|------|----------|------|
| 1 | 内核端口落地；handler / 模块公共契约无 `*sql.Tx` / 驱动类型 | **达成**（有残留，见 F-001） | 接口与模块 `TxRunner` 已收口；sqlite 具体类型仍留 `WithTx(*sql.Tx)`，测试仍调用 |
| 2 | 开区 compiled 台账两方言可 apply + checksum | **达成** | 本轮 PG 全量 boot PASS；sqlite handler/store 测试 PASS；checksum drift 测例 PASS |
| 3 | SQLite 仍为本地/Compose 默认；无 PG 仍能开发快测 | **达成** | `config.yaml` / `compose.yaml`；PG 测试无 env 则 skip（`pgtest.DSN()==""`） |
| 4 | 生产向验收以 PG 为准（迁移、readyz、共事务、备份/恢复之一可核对） | **达成** | 本轮四项均复跑：boot、composition Start/Ready、共事务、pg_dump/restore |
| 5 | 未引入 ORM；未进 Admin 功能/业务域 | **达成** | `go.mod`；本区代码为 store/config/migration/repository 端口，无新业务页 |
| 6 | 开放 required finding = 0（或已合法闭合） | **本条无 required** | 下列 findings 均为 recommended；不阻断闭门 |

### P-005 信息项（本轮按代码复核，不采信 meta 的 verified 字样）

| ID | 级别 | 本轮结论 |
|----|------|----------|
| I-001 | required · R5 | **有界 residual 与 VP 退出 2 一致**：代码中**没有** SQLite 文件库 in-place 升 PG 的产品路径；有测例级 repository 拷贝（本轮 PASS）。fresh bootstrap 本轮已证。 |
| I-002 | required · R2 | **成立**：`pgx/v5/stdlib`，驱动名 `pgx`。 |
| I-003 | non-blocking · R4 | **大体成立**：生产模块公共面无 `*sql.Tx`；sqlite `Store.WithTx` 与测试仍泄漏（F-001）。 |
| I-004 | required · R5 | **本轮独立成立**：`pg_dump -F c` → `pg_restore` 对 48/35 catalog 台账 checksum 一致。应用内**没有** `Backup()` API（`kernel.Store` 无此方法；sqlite 仍 `VACUUM INTO`）。 |

## Findings

### F-001 · sqlite 具体类型仍公开 `WithTx(*sql.Tx)`，测试继续走驱动事务

| 字段 | 值 |
|------|-----|
| level | recommended |
| 严重度 | med |
| status | open |
| evidence | `apps/api/internal/store/store.go` L72 注释仍写「until R4 moves their public signature」；L107 `func (s *Store) WithTx(ctx, fn func(*sql.Tx) error)`。生产 `*postgres` **没有**对应方法。调用点全部在 `_test.go` / `testsupport/store.go`（auth、composition、wallet、settings、handler/wallet 等）。模块 `TxRunner` 与 handler 签名已是 `kernel.Tx`/`kernel.Store`。 |
| closure | — |

描述：R4 声称收口的是「Handler / 模块公共契约」。生产装配已走 `kernel.Store`，**不**把 `WithTx` 注入模块。但 sqlite 实现类型仍把 `*sql.Tx` 留在可调用面上，且注释仍把 R4 写成未完成。新测试若复制 `testsupport` 模式会继续方言锁定。不否定退出判据 1 的生产契约，但是闭门后的卫生债。

### F-002 · 连接池旋钮未接到配置 / composition

| 字段 | 值 |
|------|-----|
| level | recommended |
| 严重度 | med |
| status | open |
| evidence | `store/options.go` 有 `PoolMaxOpenConns` / `PoolMaxIdleConns` / `ConnMaxLifetime`；`postgres.applyPostgresPoolOptions` 仅在 opts>0 时生效。`composition.openStore` 只传 `Dialect/Path/DSN`，**不传池参数**。`config.Config` / `config.yaml` **无**池键。`database/sql` 默认 `MaxOpenConns=0`（无上限）。 |
| closure | — |

描述：VP 清单写了「连接池」。实现上有 stdlib 池能力，但生产启动未设置上限。这不等于「没有连接」，也不阻断 fresh boot / readyz；作为生产权威方言，缺操作面限流。建议后续用配置键接到 `OpenOptions`，不需要重开 VP。

### F-003 · 公共面仍导入 `database/sql`；jobs 仍用 `sql.ErrNoRows`

| 字段 | 值 |
|------|-----|
| level | recommended |
| 严重度 | low |
| status | open |
| evidence | `kernel/store.go` `var ErrNoRows = sql.ErrNoRows`，注释「Until the public surface stops importing database/sql (R4)」；`jobs/repository.go` L329、`jobs/model.go` L107 `errors.Is(err, sql.ErrNoRows)`。若干模块内部仍用 `sql.NullString`/`sql.NullInt64`（authsession、operationlog、wallet、recyclebin）。 |
| closure | — |

描述：不是 `*sql.Tx` 泄漏。R1 合同允许 `ErrNoRows` 别名过渡。闭门后可把模块改为 `kernel.ErrNoRows`，并把 `sql.Null*` 留在 store 适配层。

### F-004 · `pg_dump` 17 客户端对 PG 15.4 恢复出现 `transaction_timeout` 告警

| 字段 | 值 |
|------|-----|
| level | recommended |
| 严重度 | low |
| status | open |
| evidence | 本轮：服务器 `SHOW server_version` = `15.4`；客户端 `pg_dump (PostgreSQL) 17.11`。`pg_restore` 报 `unrecognized configuration parameter "transaction_timeout"`（errors ignored: 1）。**随后** `schema_migrations` 48=48、checksum `4f21236ecaeb6af0dfe1b2194d9f157f`、35 表双方一致。CI 用 `postgres:15-alpine`。 |
| closure | — |

描述：备份合同**功能上成立**。合同文本若只写 `pg_dump -F c` 而不钉客户端主版本，在 17 客户 / 15 服务组合下会有可忽略的 SET 噪声。建议合同补一句：dump 客户端主版本 ≤ 服务器，或恢复时允许该参数告警。

### F-005 · 先前 Root `status: done` 时 `03-audit` 无关门意见（程序缺口，本条补齐）

| 字段 | 值 |
|------|-----|
| level | recommended |
| 严重度 | low |
| status | open |
| evidence | 本审开始时 `GOAL-001-store-dialects/03-audit/` 为空；`03-audit.md` 索引「尚无」，但 `00-meta`/`goal-tree` 已 `done` 5/5。P-002 默认「关门审计后结项」。子目标有独立审，但**被审目标是 Root**，Root 台账当时没有 A 条。 |
| closure | — |

描述：这**不**否定本轮代码证据。它说明「把 Root 标 done」那次动作缺少 Root 级独立关门意见。本 A-001 补上该台账。`/govern` 可在响应里注明：闭门依据改为本独立审 + 代码/复跑，而不是空索引上的 status 字样。

## 必改项汇总

**无 required / 必改项。**

建议（非门禁）：F-001 删除或收窄 `Store.WithTx`；F-002 把池上限接到配置；F-003 模块改 `kernel.ErrNoRows`；F-004 备份合同钉客户端版本。

## 与既有意见的异同

子目标意见**只对照、不采信为通过**：

| 来源 | 异同 |
|------|------|
| GOAL-006 independent A-001（conditional：I-001 过满、备份曾用玩具表） | 本轮**不再**用玩具表：独立 dump/restore 为 48 迁移 / 35 表，checksum 与当时记录的 `4f21236ecaeb6af0dfe1b2194d9f157f` 相同，但是**本轮重新算出**。I-001 维持「无产品 in-place、仅测例原型」——与 VP 退出 2 的 residual 许可一致，**不开 required**。 |
| GOAL-005 independent A-004（`instr()` / 公共面） | 本轮全仓 `instr(` 仅剩测试注释；模块 `TxRunner` 已是 `kernel.Tx`。新发现是 sqlite `WithTx` 残留（F-001），R4 独立审未当作 Root 闭门项。 |
| Root 自述 E-014 / 子目标 self pass | **明确不作为本条通过依据。** 本条通过只建立在代码与 2026-08-21 复跑上。 |

## 结论 + 建议给编排器/用户的下一步

**verdict = pass。** 工作区 13 根目标在**代码实现**上已经满足 VP-013 退出判据 1–5；本轮独立复跑了 PG 全量台账、组合启动、跨模块共事务、以及 catalog 级 `pg_dump`/`pg_restore`。SQLite 默认路径仍可用。无 ORM。没有开放的 required finding。

先前把 Root 标成 `done` 时 Root `03-audit` 是空的（F-005）：那是程序不完整，不是实现缺失。本意见不改 `status`/`progress`/`goal-tree`。

建议 `/govern` 下一句：

```text
/govern 响应 GOAL-001 A-001（independent pass）：F-001～F-005 均为 recommended。请确认是否把 Root 闭门依据改为本独立审的代码+复跑证据；recommended 项可另立项或 accepted-residual。
```

## 声明

本意见 `source: independent`。不修改 `00-meta` status/progress、不改 goal-tree、不改方案正文。响应与是否维持 `done` 由 `/govern` 处理。
