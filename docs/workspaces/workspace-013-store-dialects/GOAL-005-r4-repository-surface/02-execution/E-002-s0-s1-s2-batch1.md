---
id: E-002
doc: execution-entry
goal: GOAL-005-r4-repository-surface
status: recorded
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# E-002 · S0 扫描 + S1 接缝 + S2 首批 6 模块迁移（kernel.Tx）

## 2026-08-20 · R4 首批落地

### 已发生事实

- **S0 泄漏面扫描（I-001 → verified）**：ripgrep `func(tx *sql.Tx)` 全仓 144 处（~30 生产文件 + 测试）。生产面分四类：
  - A. 模块 store 的本地 `TxRunner`（`r.runner.WithTx(ctx, func(*sql.Tx))`）：wallet/datapermission/datadictionary/mfa/scheduledtasks/recyclebin/logincaptcha；
  - B. `r.withTx(label, func(*sql.Tx))` helper：authsession（多文件）/operationlog/settings；
  - C. jobs `TxRunner`（14 处）+ systemdata `TxRunner`（bootstrap/reconcile）；
  - D. 公共回调传入 `*sql.Tx`：auth.go:632、handler/service_credentials.go:190（`CreateServiceCredential(…, func(tx *sql.Tx))`）、composition.go `RecordOperationTx`。
- **S1 接缝**：`kernel.Tx` 已是端口；各模块 `TxRunner` 接口从 `WithTx(ctx, func(*sql.Tx))` 改为 **`Run(ctx, func(kernel.Tx) error)`**（与 `kernel.Store`/`*store.Store.Run` 对齐；sqlite `WithTx(*sql.Tx)` 保留作测试适配器，R5 前清理）。
- **S2 首批 6 模块**（commit `299c7dc`）：**logincaptcha**（含测试 fake `failingRunner`→`Run`）、datapermission、datadictionary、mfa、scheduledtasks、recyclebin 全部迁移到 `Run(ctx, func(kernel.Tx))`；回调内 `tx.Exec/Query/QueryRow` 加 ctx、`sql.ErrNoRows`→`kernel.ErrNoRows`。
- **live PG 证据**：`TestFullCatalogPostgresBootstrapIntegration` 在全量 bootstrapped PG 上跑迁移后的 `logincaptchastore.NewRepository`（SetEnabled/Enabled/CreateChallenge/ConsumeChallenge）**端到端 PASS**——证明 kernel.Tx 端口在 postgres 上真正可用。
- 全量 `go test ./...` 0 FAIL；`go build ./...` 绿。

### 证据

| 主张 | 路径 / commit |
|------|---------------|
| S0 扫描 | 本会话 ripgrep：144 处 `func(tx *sql.Tx)` |
| S1 接缝 + 6 模块 | `apps/api/internal/modules/{logincaptcha,datapermission,datadictionary,mfa,scheduledtasks,recyclebin}/store/repository.go`（commit `299c7dc`） |
| 迁移后仓库 live PG | `apps/api/internal/store/postgres_test.go`（logincaptchastore 在 bootstrapped PG 上端到端） |
| 全量回归 | `apps/api`: `go test ./...` 0 FAIL |

## 待续（S2 剩余 + S3/S4）

- wallet（`rowQueryer`/`ReconcileOnceTx`/helper `*sql.Tx` 深层）；authsession（withTx 多文件）+ systemdata；jobs（14 回调）；operationlog + settings（withTx helper + `INSERT OR IGNORE` SQL 债）；D 类公共回调（auth/handler/composition `RecordOperationTx`）。
- S3 jobs/handler 收口；S4 composition postgres 启动；S5 关门。
