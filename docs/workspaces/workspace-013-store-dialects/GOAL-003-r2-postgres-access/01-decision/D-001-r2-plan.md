---
id: D-001
doc: decision-entry
goal: GOAL-003-r2-postgres-access
status: accepted
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# D-001 · R2 方案：访问层边界与实施切面

## 触发

- R1 合同 v1.4（GOAL-002 D-005 / A-009）把 R2 边界钉死；Root I-002 经 Root D-002 已 `verified`（pgx v5 stdlib）。
- 本目标进入 R2 实施：驱动、连接池、`Open`/`Ping`/`WasFresh`（探测打开）、config 校验。R3 才 apply 双方言 catalog。

## 决定

1. **依赖**：`github.com/jackc/pgx/v5/stdlib`（驱动名 `"pgx"`），仅 `internal/store` 触达；内核 `database/sql` 形状端口不引 pgx 类型。
2. **接口**（`internal/kernel`，新增 `store.go`）：`Dialect`（`sqlite`/`postgres`）、`Store`（`Dialect/Run/Ping/Close/WasFresh/MarkSystemDataReady/SystemDataReady`）、`Tx`/`Result`/`Rows`/`Row`、`ErrNoRows = sql.ErrNoRows`（别名变量，R4 切断前两路 `Is` 同真）。
3. **打开**（`internal/store`）：新增 `OpenOptions{Dialect,Path,DSN,PoolMaxOpenConns,PoolMaxIdleConns,ConnMaxLifetime,ConnectTimeout}` 与 `Open(ctx, opts, catalog) (kernel.Store, error)` 方言分发。应用代码仍走 `Open`；`OpenWithCatalog` 保留为 sqlite 测试适配器（R5 前删除或降 helper）。
   - sqlite：现行路径 apply catalog（缺省开发不变）。
   - postgres：`sql.Open("pgx", dsn)` → `Ping` → `WasFresh`(search_path 感知)；**空/`nil` catalog = 探测打开**；**非空 catalog fail closed**（含 sqlite 专用 SQL，不得半执行）。
4. **占位符**：统一 `?`，postgres 在 store 层 rebind `?`→`$1..$n`（内联实现，不引 sqlx）。本拍 rebind 与测试即占位符合同走通。
5. **config**：`DBDialect` / `DBDSN` 字段 + YAML `db.dialect`/`db.dsn` + env `DB_DIALECT`/`DB_DSN`。校验并入 `ValidateProd()`（启动 fail-closed，全环境）：
   - dialect ∈ {空, `sqlite`, `postgres`}；空=sqlite。
   - sqlite：`dsn` 必须空；`path` 非空。
   - postgres：`dsn` 必须非空；`path` 用缺省 `./data/schema-ui.db`、必须文件路径形状（含扩展名谓词：尾分隔符 / 已存在目录 / Base ∈ {.,..} / 无扩展名 → 拒绝）。
   - 不读 `DATABASE_URL`。
6. **composition.openStore 路由**：经 `store.Open(ctx, opts, catalog)` 分发；sqlite 行为逐字节不变，postgres 下对非空 catalog **fail closed**（明确错误，符合 v1.4 R2 边界）。
7. **测试**：sqlite 运行时（`Run`/`Tx`/回滚）+ rebind 单测 + config 校验单测 + postgres 探测集成测试（env `SCHEMA_UI_R2_PG_DSN` 门控，缺省 skip）。
8. **不改**：模块 / jobs / handler 公共签名（`WithTx(*sql.Tx)` 与 `MigrationContribution.Apply func(*sql.Tx)`），全部留 R4/R3。

## 为什么

- pgx stdlib 与 `database/sql` 形状端口天然兼容；生产向连接池由 `database/sql` 池语义承载。
- postgres 非空 catalog fail closed 是 v1.4 明文 R2 行为，防止把 sqlite 专用 SQL 半执行到 PG。
- 不立刻切模块签名，避免把 R4 的 `*sql.Tx` 迁移混进 R2，破坏「缺省开发不能断」。
- config 校验进 `ValidateProd`：现有 main 唯一启动校验点，全环境生效且无需改 main。
- pg 集成测试用 env 门控：无 PG 的本地/CI 不因连不上而红，符合「没有 Postgres 仍能开发快测」。

## 未选方案

- **pgx 原生驱动直接当端口**（pgxpool + pgx.Tx）：泄漏驱动类型到内核面，违反 R1 §1。排除。
- **本回合就切 `MigrationContribution.Apply` 为 `func(kernel.Tx)`**：全模块迁移面一次大改，属 R3 双方言对写的既有契约；R2 保持 `*sql.Tx` 避免范围爆炸。排除。
- **postgres 本拍就 apply 现行 catalog**：明文禁止（v1.4）。排除。
- **config 校验放独立方法且不接 ValidateProd**：main 只调 ValidateProd；另接=漏。纳入后者。
- **默认 dialect=postgres**：违反 RT-P03 内嵌默认。排除。

## 影响范围

- `apps/api`：`go.mod`/`go.sum`；`internal/kernel`（新 `store.go`）；`internal/store`（`open.go`/`options.go`/`rebind.go`/`postgres.go` + `store.go` 增方法）；`internal/config`（字段 + 校验 + default.yaml/config.yaml）；`internal/composition`（openStore 走 `Open`）。模块与 handler 不改。
- 用户可感知行为：显式配置 `db.dialect=postgres` + `db.dsn` 的启动会 fail closed（clear error）而非静默开 sqlite 文件；缺省/`sqlite` 行为不变。

## 后续

- S2–S5 实施 + self 审计（A-001）；R2 实现后的独立审计（项目默认 grok build）在实现切片完成并自审后另行执行；关门前补 A 条目与证据。
