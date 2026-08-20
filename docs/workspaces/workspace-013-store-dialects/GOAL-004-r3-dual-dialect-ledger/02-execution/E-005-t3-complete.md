---
id: E-005
doc: execution-entry
goal: GOAL-004-r3-dual-dialect-ledger
status: recorded
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# E-005 · T3 完成（operationlog 对写 + postgres open 解闸 + 合规扫尾）

## 2026-08-20 · T3 落地完成

### 已发生事实

- **operationlog 对写**（commit `5e0341e`）：18 个迁移以 `pgTimeDDL` 从 canonical sqlite 字面量派生 PG DDL（`created_at`/`archived_at` → `BIGINT`；事件 CHECK 列表与 checksum 输入锁步）；`pgExecDDL`/`pgRebuild`/`pgRebuildWithCorrelation` 三个 PG apply 构造器；0043/0045 保留 correlation 备份/重建/恢复（PG 上 FK 依赖 operation_log，重建 rename 会打断）。`operation_log`/`operation_log_archive` 的 created_at/archived_at 现为 `bigint`（断言）。
- **postgres open 解闸**（commit `e7fd924`）：`openPostgres` 对非空 catalog 改为经 `postgres.migrate` 执行（R2 期 fail-closed 守卫依 R3 具备双写能力而移除）；`TestOpenPostgresAppliesNonEmptyCatalogIntegration` 证明 Open 时 bootstrap。
- **合规扫尾**：全量 bootstrap 测试新增系统级检查——public schema 中所有疑似 Unix 时间列（created_at/updated_at/expires_at/applied_at/archived_at/lease_expires_at/finished_at/started_at/read_at/revoked_at/last_used_at/restored_at/deleted_at）`data_type='integer'` 计数必须为 0 → **PASS**（0 残留）。
- 全量 compiled catalog（**48 迁移全部双写**）在 live PG fresh bootstrap + 台账 + 幂等 + BIGINT 断言全绿；`go test ./...` 0 FAIL。
- **边界注记**：composition 层的 postgres DSN 启动/路由仍未接入——仓库公共签名仍是 `*store.Store`/`WithTx(*sql.Tx)`，模块须先迁移到 `kernel.Store`/`kernel.Tx`；该工作属 **R4 仓库公共面收口**（Root R4），不在 R3 交付面内。

### 证据

| 主张 | 路径 / commit |
|------|---------------|
| operationlog 双写 | `apps/api/internal/modules/operationlog/migration/migration.go`（`5e0341e`） |
| open 解闸 | `apps/api/internal/store/postgres.go`（`e7fd924`） |
| 全量 boot + 无 int 时间残留 | `TestFullCatalogPostgresBootstrapIntegration`（postgres:17-alpine；系统级检查 0 残留） |
| 全量回归 | `apps/api`: `go test ./...` 0 FAIL |
