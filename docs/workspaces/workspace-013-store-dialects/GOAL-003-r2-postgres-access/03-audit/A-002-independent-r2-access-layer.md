---
id: A-002
doc: audit-entry
record_id: A-002
goal: GOAL-003-r2-postgres-access
source: independent
scope: R2 实施切片——pgx v5 stdlib；kernel.Store/Tx/Result/Rows/Row/ErrNoRows；store.Open 方言分发；postgres 空/nil catalog 探测打开与非空 fail-closed；config db.dialect/db.dsn 校验；?→$n rebind；goroutine-local 嵌套 Run；测试与回归
verdict: pass
status: recorded
parent: GOAL-003-r2-postgres-access
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# A-002 · R2 访问层实施独立交叉审计（2026-08-20）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：`execution-facts` + `close-out-readiness` · GOAL-003 R2 访问层实施切片（S2–S5 / 用户 8 项）
- **verdict**：pass
- **工作区**：`workspace-013-store-dialects`（Root `GOAL-001-store-dialects`；canonical `docs/workspaces/workspace-013-store-dialects/`；`shared_materials_catalog: none`）

## 范围与区间

- **covered**：(1) pgx v5 stdlib 驱动引入；(2) `kernel.Store` / `kernel.Tx` / `Result` / `Rows` / `Row` / `ErrNoRows` 端口；(3) `store.Open(ctx, OpenOptions, catalog)` 方言分发；(4) postgres 空/`nil` catalog 探测打开 = 连接 + Ping + WasFresh，非空 catalog fail-closed；(5) config `db.dialect` / `db.dsn` 与启动校验（枚举、DSN 配对、`db.path` 文件路径形状 + 扩展名谓词）；(6) `'?'` → `$n` rebind；(7) goroutine-local 嵌套 `Run` 检测；(8) 测试与回归（本审独立 `go build ./...`、`go test ./...`）。对照 GOAL-003 成功标准 1–5、D-001、R1 冻结合同 v1.4 §2/§3/§5、Root D-001 第 5 条 / D-002 第 5 条（R2 实现后 independent）。
- **excluded**：R3 双方言台账对写；R4 模块公共签名 / `*sql.Tx` 收口；R5 升级/备份合同；现行 HTTP `/readyz` 模块门禁全绿（v1.4「R2 证据边界」明文排除）；不改 `status` / `progress` / 方案正文 / goal-tree；不读取或比较其他工作区。
- **P-005**：GOAL-003 I-001 / I-002 均为 `verified`，最晚阶段分别为 S1 / S5，无到期未关 required 信息项。Root I-001 / I-004 最晚 R5、I-003 最晚 R4，均未到期，不构成本切片失败。
- **资料引用**：无（`shared_materials_catalog: none`）。
- **对照 self**：A-001（2026-08-20，`source: self`，`verdict: pass`，0 required，F-001～F-003 recommended）。

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 工作区绑定合格 | `workspace.md`：id / Root / canonical / `primary_plan: VP-013-store-dialects` 一致；资料目录 `none` |
| pgx v5 stdlib，驱动名 `"pgx"`，仅 store 触达 | `apps/api/internal/store/postgres.go` blank import `github.com/jackc/pgx/v5/stdlib` + `sql.Open("pgx", dsn)`；`go list`/`go mod why`：store → stdlib → `github.com/jackc/pgx/v5 v5.10.0`。全仓 `.go` 仅 postgres.go 引用 jackc/pgx |
| `kernel` 端口形状与 `ErrNoRows` 别名 | `apps/api/internal/kernel/store.go`：`Store`/`Tx`/`Result`/`Rows`/`Row`；`var ErrNoRows = sql.ErrNoRows`；`Result` 仅 `RowsAffected`（无 `LastInsertId`） |
| `Open` 方言分发 | `apps/api/internal/store/open.go` + `options.go`：sqlite apply catalog；postgres → `openPostgres`；未知方言 fail-closed |
| postgres 非空 catalog **先于** `sql.Open` fail-closed | `postgres.go`：`len(catalog) > 0` 在 `sql.Open` 之前返回含 `R3` 的错误；`TestOpenPostgresFailsClosedOnNonEmptyCatalog` PASS（无 PG 也能过，支持「未连网」） |
| 空/`nil` catalog 探测路径存在 | `len(catalog) == 0` 走 DSN + Ping + `postgresWasFresh`；`TestOpenPostgresRequiresDSN` PASS；集成测试 env 门控 |
| sqlite 缺省路径仍 apply catalog | `store.go` `open()`：`databaseIsEmpty` 后 `migrate(catalog)`；`OpenWithCatalog` 保留为测试适配器 |
| composition 经 `store.Open`；postgres + compiled catalog fail-closed | `composition.go` `openStore`：把 `PersistenceCatalog()` 交给 `Open`；非 `*store.Store` 再关连接并报「not wired until R3」（R2 下因非空 catalog 通常在 Open 已失败） |
| config 字段 + 枚举 + DSN 配对 + path 谓词 | `config.go` `DBDialect`/`DBDSN`；YAML `db.dialect`/`db.dsn`；env `DB_DIALECT`/`DB_DSN`；Load 空=sqlite、未知 LoadError；`ValidateProd`→`validateDB`；谓词 1/2/3/6（尾分隔符 / 已存在目录 / Base ∈ {.,..} / 无扩展名）。全仓 `DATABASE_URL` 零引用 |
| `'?'` → `$n` rebind（不引 sqlx） | `rebind.go`；`pgTx` 的 Exec/Query/QueryRow 均 rebind；sqlite `sqlTx` 不 rebind |
| 嵌套 `Run` 为 goroutine 局部，非 Store 级 flag | `runmarker.go` `enterRun`/`leaveRun`（`runtime.Stack` goroutine id）；sqlite 与 postgres `Run` 共用 |
| 未引 ORM；未改模块公共签名；未写 postgres DDL 到模块仓库 | commit `1305754` 文件清单无 `internal/modules/**`；`kernel.MigrationContribution.Apply` 仍 `func(*sql.Tx)`；`Store.WithTx` 仍在；无 gorm/sqlx/ent |
| 构建 + 全量测试绿（本审独立复跑） | `apps/api`：`go build ./...` PASS；`go test ./...` PASS（2026-08-20）。新测：`TestKernelStore*` 5 PASS、`TestRebindPostgres` PASS、`TestSearchPathCandidates` PASS、`TestOpenPostgresRequiresDSN` PASS、`TestOpenPostgresFailsClosedOnNonEmptyCatalog` PASS、`TestOpenPostgresProbeIntegration` **SKIP**（`SCHEMA_UI_R2_PG_DSN` UNSET）。`TestDBDialectConfig` 4 子例 + `TestValidateDBPairs` 7 子例 PASS |

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| 1. 缺省 sqlite 构建 + 既有测试回归不破 | **满足** | `go build ./...`；`go test ./...` 全绿；sqlite `open()` 仍 apply catalog |
| 2. postgres `Open`：连 + Ping + WasFresh；空 catalog 探测、非空 fail-closed | **满足（离线路径）**；live 探测未跑 | 代码顺序 + fail-closed / DSN 单测 PASS；`TestOpenPostgresProbeIntegration` SKIP（D-001 第 7 条 / 成功标准 4 允许） |
| 3. config 校验：dialect 枚举 / DSN 配对 / path 扩展名谓词 | **满足** | `TestDBDialectConfig` + `TestValidateDBPairs` |
| 4. 内核接口 + rebind 单测；pg 运行时 env 门控 | **满足** | `TestKernelStore*`、`TestRebindPostgres`、门控 SKIP |
| 5. 未引 ORM、未改模块签名、未写 postgres DDL 到模块仓库 | **满足** | commit `1305754` diff scope；`Apply`/`WithTx` 签名未改 |

R2 证据边界（v1.4 §2）：本拍可核对 = `Open` + `Ping`（及 `database/sql` 池）+ 空 catalog 探测路径，**不是** composition 全量 bootstrap，**不是**现行 `/readyz` 200。`openStore` 在 postgres + compiled catalog 上 fail-closed，与合同一致。Root 纲领里的「readyz」按 v1.4 读作「可 Ping 的连接就绪」。

## Findings

### F-001 · live PostgreSQL 探测未在本机运行

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | recommended |
| 状态 | open |
| evidence | `apps/api/internal/store/postgres_test.go` `TestOpenPostgresProbeIntegration`（本审：`SCHEMA_UI_R2_PG_DSN` UNSET → SKIP）；本审未执行 live `Open`+`Ping`+`WasFresh`+rebind `Run` |
| closure | 无 |
| 关联 I-00N | 无（成功标准 4 / D-001 第 7 条把无 PG skip 定为合法） |

空 catalog 探测的**真实 PG** 路径（连接、Ping、WasFresh、`?`→`$n` 的 `Run`）没有运行时证据。fail-closed 与解析单测不能替代一次真实握手。**不构成** R2 代码验收必改（合同与成功标准 4 允许 skip），但 R3 对写开工前应在 compose/CI 设 `SCHEMA_UI_R2_PG_DSN` 跑通该门控测试。

与 self A-001 F-001：**同意**事实与级别。关门校准见「与既有意见的异同」。

### F-002 · 嵌套 `Run` 检测是 goroutine-id 启发式；注释仍写 ctx；缺并发 `Run` 回归

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |
| evidence | `apps/api/internal/store/runmarker.go`；`store.go` L74–76 与 `postgres.go` L73–75 注释写「via the ctx passed to fn」，实现未把 marker 放入 ctx；仅 `TestKernelStoreNestedRunFailsClosed`（同 goroutine） |
| closure | 无 |

实现按 v1.4 允许的「goroutine 局部」做：`runtime.Stack` 取 goroutine id + 进程级 map 计数，**不是** Store 级 `inRun`，因此 postgres 池上不同 goroutine 的并发 `Run` 不会被误杀。合同允许该等价机制。缺口：

1. 同 goroutine 回调内再 `Run` 能拦；回调里另起 goroutine 再 `Run` 视为并发、不拦（调用方误用）。
2. sqlite/postgres `Run` 注释声称 ctx 传递，与代码不符，R4 接入时易误改。
3. 没有「两 goroutine 并发 `Run` 必须成功」的回归，挡不住以后改成 Store 级互斥。

R2 模块仍走 `WithTx`，`Run` 尚未成为公共入口。R3/R4 接入前复核。与 self A-001 F-002：**同意**启发式判断；本条补注释漂移与并发回归缺口。

### F-003 · WasFresh 的 search_path 为手解析，缺 live 验证，且不跳过显式系统 schema

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |
| evidence | `postgres.go` `searchPathCandidates` / `firstExistingSchema` / `postgresWasFresh`；`TestSearchPathCandidates` 仅解析层；无 `current_schemas(false)`、无 live 非默认 search_path |
| closure | 无 |

`$user` 经 `current_user` 解析，**没有**把字面 `"$user"` 拿去查 `information_schema`，满足 v1.4 禁令的字面。但实现是 `SHOW search_path` 手解析，不是服务器已解析列表。`firstExistingSchema` 不跳过 `pg_catalog` / `information_schema`：若 search_path **显式**把系统 schema 放在最前，随后 `count` 因 `NOT IN (pg_catalog, …)` 得到 0，可能把后面 `public` 里已有用户基表的库报成 fresh。默认 `"$user", public` 不受影响。建议 R3 前改为 `current_schemas(false)`（或跳过系统 schema）并并入 F-001 的 live probe。

与 self A-001 F-003：**同意**缺 live；本条多一条系统 schema 短路径风险。

### F-004 · `Run` 在 `fn` panic 时未按合同 rollback 再 re-panic

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | recommended |
| 状态 | open |
| evidence | R1 冻结 v1.4 §2「`fn` panic：rollback，再 re-panic」；`store.go` `(*Store).Run` 与 `postgres.go` `(*postgres).Run` 均无 `recover`；`TestKernelStoreRunCommitAndRollback` 只覆盖 error 回滚，不覆盖 panic |
| closure | 无 |

sqlite 与 postgres 的新 `Run` 在 `fn` 返回 error 时 rollback（已测）；`fn` panic 时只靠 `defer leaveRun()` 清嵌套标记，事务连接依赖 `*sql.Tx` 最终器/GC 才可能回滚，postgres 池上会多占连接。v1.4 写明 R2 起必须遵守该语义。

**不升 required** 的理由：R2 composition 仍用 `WithTx`（同样无 recover，且为既有生产路径）；`Run` 尚未被模块调用。必须在 **R4 把 `Run` 作成公共事务入口之前**补 `recover` + rollback + re-panic，并加单测。本目标关门不因本条阻断。

### F-005 · `pgx` 在 `go.mod` 标为 `// indirect`，与 store 实引不一致

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | recommended |
| 状态 | open |
| evidence | `apps/api/go.mod`：`github.com/jackc/pgx/v5 v5.10.0 // indirect`；`go mod why` 显示 `internal/store` → `pgx/v5/stdlib`；本审 `go build`/`go test` 仍通过 |
| closure | 无 |

依赖版本与驱动名正确，功能不破。`go mod tidy` 应将主模块直接 import 提升为 direct require，避免后续 tidy/工具链把「间接」读错。属卫生项，不阻断 R2。

## 必改项汇总

无。0 required。

开放 recommended：F-001（live PG）、F-002（嵌套检测启发式/注释/并发回归）、F-003（search_path 手解析 + 系统 schema）、F-004（`Run` panic rollback）、F-005（go.mod indirect）。

## 与既有意见的异同（self A-001）

| 点 | self A-001 | 本独立审 A-002 |
|----|------------|----------------|
| verdict | pass | **pass**（同意） |
| 必改 | 0 | **0**（同意） |
| live PG SKIP | F-001 recommended | **同意**；编号同主题为 F-001 |
| 嵌套检测启发式 | F-002 recommended | **同意**，并补注释写 ctx、缺并发回归（本 F-002） |
| search_path 缺 live | F-003 recommended | **同意**，并补系统 schema 短路径（本 F-003） |
| `Run` panic 语义 | 未写 | **新增** F-004 recommended（不升 required） |
| pgx `// indirect` | 未写 | **新增** F-005 recommended |
| 成功标准 1–5 离线项 | 均 ✅ | **同意**；本审独立复跑 `go test ./...` |
| 关门条件 | 「独立审 **与（或随后）live PG** 完成后才具备关门条件」 | **部分不同意**：独立审本条已落盘；live PG 按 D-001 第 7 条 / 成功标准 4 是 recommended 残余，**不是**本目标关门 required。建议 R3 开工前跑，不把 GOAL-003 `done` 绑死在本机是否有 PG |
| 现行 `/readyz` 200 | 正确未当验收 | **同意** |

无 verdict 冲突、无「一要一否」的必改冲突。P-004 冲突裁决点 **不触发**。

## 结论 + 建议给编排器/用户的下一步

R2 访问层实施切片在可离线核对的范围内 **名实相符**：驱动、端口、Open 分发、非空 catalog fail-closed、config 门禁、rebind、goroutine 局部嵌套检测、sqlite 回归均有代码与本审复跑测试支撑。scope 内无 high/med **required**，无到期 required 信息项。

**close-out-readiness**：从必改门禁看，GOAL-003 **可以进入关门流程**。本条即 Root D-001/D-002 要求的 independent。剩余 F-001～F-005 均为 recommended，不阻断 `done`。编排器不得把本 pass 当成已改 status。

建议 `/govern` 下一步：

1. 响应 A-001 + A-002（本条无需闭合必改；recommended 可保持 open 并指定复核点：F-001/F-003 → R3 前 live PG；F-002/F-004 → R4 `Run` 公共化前；F-005 → 下次改 `go.mod` 时 tidy）。
2. 更新 GOAL-003 路线图检查点（00-meta S2–S5 仍写「待做」，与 E-002/测试事实不一致——属编排器文档，本审未改）。
3. 用户确认后将 GOAL-003 标 `done`，再立项 R3 台账对写。
4. 有 PG 时设 `SCHEMA_UI_R2_PG_DSN` 跑 `TestOpenPostgresProbeIntegration`（建议，非本目标必改）。

## 声明

本意见不修改 status/progress/方案正文/goal-tree；响应由 `/govern` 处理。
