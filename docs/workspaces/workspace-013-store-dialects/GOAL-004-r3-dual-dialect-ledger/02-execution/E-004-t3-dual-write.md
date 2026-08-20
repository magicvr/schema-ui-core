---
id: E-004
doc: execution-entry
goal: GOAL-004-r3-dual-dialect-ledger
status: recorded
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# E-004 · T3 对写实施（PostgresApply 管道 + 12/13 模块；operationlog 待续）

## 2026-08-20 · T3 主体落地

### 已发生事实

- `kernel.MigrationContribution` 增 `ApplyPostgres func(Tx) error`（可选；`postgres.migrate` 优先用 PostgresApply，否则回退可移植 Apply）。checksum 仍绑 sqlite/canonical 历史（v1.4 §4）。
- **已对写 12 个模块**（PostgresApply + PG DDL，时间列 BIGINT、money BIGINT、`COLLATE NOCASE`→CITEXT、partial index 保留）：authsession、notifications、jobs、datadictionary、datapermission、mfa、logincaptcha、corepersistence、settings、scheduledtasks、recyclebin、wallet。account（纯 ADD 列 ALTER）可移植无需对写。
- **惊喜发现**：未对写前，现行 compiled catalog 用可移植 `Apply` 已能在 live PG 上完整 fresh bootstrap（`INTEGER` 时间列作为 int4 可执行）——`TestFullCatalogPostgresBootstrapIntegration` 全量 apply + 台账 + 幂等通过。
- **BIGINT 合规断言**：上述 bootstrap 测试新增对 12 模块 ~20 个时间/金额列的 `data_type=bigint` 断言（全绿）。
- **待续（T3 收尾）**：`operationlog` 模块——`operation_log` 多次 rebuild（事件 CHECK 累积）与 `operation_log_archive`，需逐条成对 PG DDL（`created_at BIGINT`）。当前全量 bootstrap 仍可执行（int4），但 `operation_log.created_at` 在 PG 上非 BIGINT，不满足 v1.4 §3。

### 证据

| 主张 | 路径 / 命令 / commit |
|------|----------------------|
| PostgresApply 管道 | `apps/api/internal/kernel/contribution.go`、`apps/api/internal/store/postgres.go` |
| 12 模块对写 | `apps/api/internal/modules/*/migration/migration.go`（diff，commit `8baca38` + `ad3a876`） |
| 全量 PG bootstrap + BIGINT | `TestFullCatalogPostgresBootstrapIntegration`（live postgres:17-alpine，`-count=1` PASS） |
| 全量回归 | `apps/api`: `go test ./...` 0 FAIL；commit `8baca38`/`ad3a876` |
