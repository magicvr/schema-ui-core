---
id: E-002
doc: execution-entry
goal: GOAL-004-r3-dual-dialect-ledger
status: recorded
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# E-002 · T1 实施：迁移贡献形状转 kernel.Tx

## 2026-08-20 · T1 完成（全量测试绿）

### 已发生事实

- `kernel.MigrationContribution.Apply` / `Reconcile` 由 `func(*sql.Tx) error` 改为 **`func(Tx) error`**（`internal/kernel/contribution.go`；去掉 `database/sql` import）。
- `internal/store/migrate.go` `applyMigration` 以 `sqlTx{tx}` 适配器调用 `Apply`（sqlite 保持 `?`、事务原子）。
- 14 个模块迁移包机械迁移：签名 `tx *sql.Tx` → `tx kernel.Tx`，`tx.Exec/Query/QueryRow` 加 `context.Background()` 首参；`database/sql` import 移除（authsession 保留——仍用 `sql.NullString`）。kernel 测试夹具同步。
- 验证：`go build ./...` 0 错；`go test ./...` **0 FAIL**（含 live PG 探测）；kernel/store/jobs-migration 定向测试绿。

### 证据

| 主张 | 路径 / 命令 / commit |
|------|----------------------|
| 形状变更 | `apps/api/internal/kernel/contribution.go` |
| store 运行器 | `apps/api/internal/store/migrate.go`（`sqlTx{tx}` 适配 Apply） |
| 14 迁移包 | `apps/api/internal/modules/*/migration/migration.go`（diff 143+/142-，机械替换） |
| 全量回归 | `apps/api`: `go build ./...` + `go test ./...` 0 FAIL；commit `8932148` |
