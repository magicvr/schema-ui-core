---
id: E-003
doc: execution-entry
goal: GOAL-004-r3-dual-dialect-ledger
status: recorded
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# E-003 · T2a 实施：postgres 迁移运行器（live PG 证明）

## 2026-08-20 · T2a 完成

### 已发生事实

- 新增 `(*postgres).migrate(catalog)`（`internal/store/postgres.go`）：镜像 sqlite 运行器——
  - 复用 `normalizeCatalog` / `validateApplied` / `pendingMigrations`；
  - `appliedMigrationsPG`（`to_regclass('schema_migrations')` 探测 + 顺序读台账）；
  - `applyMigrationPG` = 一次 `Run`（一个迁移一个事务），经 `kernel.Tx` 执行 Apply + 台账插入（`?`→`$n` rebind 生效）；
  - 台账 checksum 仍绑 sqlite/canonical 历史文本（R1 v1.4 §4）。
- 集成测试 `TestPostgresMigrateRunnerIntegration`（`SCHEMA_UI_R2_PG_DSN` 门控，live postgres:17-alpine）：
  - 便携 scratch catalog（bootstrap 建 `schema_migrations`+`r3_users`；seed 用 `?` rebind 插 2 行）→ apply 成功，台账 2 行、`r3_users` 2 行、`WasFresh` 翻转正确；
  - 重开 + 再 migrate = 幂等；checksum drift fail-closed（`drift`）。
  - 为共用 live DB 增加了确定性 clean-slate（seed drop + 专用连接 teardown drop），与 `TestOpenPostgresProbeIntegration` 解耦，`-count=2` 复跑绿、测试后库净。
- **计划细化（相对 D-001 §4/§5）**：`openPostgres` 的「非空 catalog 解闸 + composition postgres 路由」**并入 T3**——因为现行 compiled catalog 仍是 sqlite-only（含 `sqlite_master`/`PRAGMA`），在真实双方言 catalog 落地前解闸只会得到 sqlite SQL 语法错误而非干净行为。本拍交付可点亮 T3 的运行器本身。
- 验证：`go build ./...` 绿；`go test ./...` **0 FAIL**（含 live PG 三测试）。

### 证据

| 主张 | 路径 / 命令 / commit |
|------|----------------------|
| 运行器 | `apps/api/internal/store/postgres.go`（`migrate`/`applyMigrationPG`/`appliedMigrationsPG`） |
| live 证明 | `TestPostgresMigrateRunnerIntegration`（`go test -run TestPostgresMigrateRunnerIntegration ./internal/store/` PASS；postgres:17-alpine） |
| 全量回归 | `apps/api`: `go test ./...` 0 FAIL；commit `4359c7f` |
