---
id: E-002
doc: execution-entry
goal: GOAL-003-r2-postgres-access
status: recorded
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# E-002 · R2 访问层实施（驱动 / 内核端口 / Open / config）

## 2026-08-20 · 实施完成并通过全量测试

### 已发生事实

- 依赖：`go get github.com/jackc/pgx/v5/stdlib@latest` → pgx **v5.10.0**（含 puddle/pgpassfile/pgservicefile 间接）。驱动名 `"pgx"`。
- 新增 `apps/api/internal/kernel/store.go`：`Dialect`（sqlite/postgres）、`Store`、`Tx`/`Result`/`Rows`/`Row`、`ErrNoRows = sql.ErrNoRows`。
- 新增 `internal/store/`：`options.go`（`OpenOptions`）、`open.go`（`Open` 方言分发）、`rebind.go`（`?`→`$n`）、`runmarker.go`（goroutine-local 嵌套 `Run` 检测）、`postgres.go`（连接 + Ping + WasFresh + pgTx；非空 catalog fail closed）。
- 修改 `internal/store/store.go`：`*Store` 增加 `Dialect()` 与 `Run(func(kernel.Tx))`（sqlite 实现，`sqlTx` 适配器），满足 `kernel.Store`；`WithTx(*sql.Tx)` 保持供 R4 前模块使用。
- 修改 `internal/composition/composition.go` `openStore`：经 `store.Open` 分发；sqlite 行为不变；postgres 下对 compiled catalog **fail closed**。
- 修改 `internal/config/config.go` + `config.default.yaml` + `configs/config.yaml`：`db.dialect`/`db.dsn` 字段 + env `DB_DIALECT`/`DB_DSN` + Load 规范化（空=sqlite、未知拒绝）+ `ValidateProd().validateDB()`（DSN/dialect 配对 + `db.path` 文件路径形状与扩展名谓词）。
- 测试：`kernel_store_test.go`（sqlite Dialect/WasFresh/Ping、Run commit+rollback、`kernel.ErrNoRows` 双 `errors.Is`、嵌套 Run fail-closed、Query scan）、`postgres_test.go`（rebind、search_path 解析、缺 DSN / 非空 catalog fail-closed、PG 探测集成 env 门控）、config_test 增 `TestDBDialectConfig` + `TestValidateDBPairs`。
- 验证：`go build ./...` 通过；`go test ./...` 全量通过（composition/kernel/jobs/cmd/server 均 ok）；新 R2 单测 10/10 PASS，PG 探测在无 `SCHEMA_UI_R2_PG_DSN` 时正确 SKIP。

### 证据

| 主张 | 路径 / 命令 / commit |
|------|----------------------|
| 依赖拉取 | `apps/api/go.mod`/`go.sum`：`github.com/jackc/pgx/v5 v5.10.0` |
| 内核端口 | `apps/api/internal/kernel/store.go` |
| store 实现 | `apps/api/internal/store/{options,open,rebind,runmarker,postgres,store}.go` |
| config 面 | `apps/api/internal/config/{config.go,config.default.yaml}`、`apps/api/configs/config.yaml` |
| composition 路由 | `apps/api/internal/composition/composition.go` |
| 测试证据 | `apps/api/internal/store/{kernel_store,postgres}_test.go`、`internal/config/config_test.go`；`go test ./...` 全绿 |
| checkpoint | commit `1305754`（38 files；governance + implementation） |
