---
id: E-001
doc: execution-entry
goal: GOAL-002-r1-tx-port-and-config
status: recorded
parent: GOAL-001-store-dialects
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# E-001 · 现行 db.path / WithTx 扫描（2026-08-20）

## 已发生事实

只读扫描，未改代码。

### 打开与配置

- YAML `db.path`；env `DB_PATH` 覆盖。`config.Config.DBPath` 缺省 `./data/schema-ui.db`。无 dialect/dsn 键。
- `composition` 调用 `store.OpenWithCatalog(cfg.DBPath, catalog)`。
- `store.Open` 使用 `modernc.org/sqlite`，`MaxOpenConns=1`。

### 事务面

- `Store.WithTx(ctx, func(*sql.Tx) error)` 是平台事务边界。
- 模块仓库（wallet、jobs 等）在公共结构里嵌入 `WithTx(context.Context, func(*sql.Tx) error)`。
- `kernel.MigrationContribution.Apply` / `Reconcile` 为 `func(*sql.Tx) error`。
- `jobs.CommitFunc` 为 `func(*sql.Tx) (json.RawMessage, error)`。

### 对 Root I-003 的收集（non-blocking，未完成）

至少这些公共签名泄漏 `*sql.Tx`：`store.Store.WithTx`、`jobs.Repository` runner、`wallet/store.Repository` runner、`kernel.MigrationContribution.Apply`、若干 handler 测试。完整清单仍属 R4。

## 证据

| 主张 | 路径 |
|------|------|
| db.path | `apps/api/internal/config/config.default.yaml`、`config.go` `yamlFile.DB` / `DB_PATH` |
| 打开 | `apps/api/internal/composition/composition.go`；`store/store.go` `OpenWithCatalog` |
| WithTx | `apps/api/internal/store/store.go` |
| 迁移 Apply | `apps/api/internal/kernel/contribution.go` |
