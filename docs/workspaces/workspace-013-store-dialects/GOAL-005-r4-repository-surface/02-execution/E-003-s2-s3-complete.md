---
id: E-003
doc: execution-entry
goal: GOAL-005-r4-repository-surface
status: recorded
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# E-003 · S2 完成 + S3 收口（全仓公共面去掉 `*sql.Tx`）

## 2026-08-20 · R4 主体落地

### 已发生事实

- **S2 迁移完成**（续 batch1 的 6 模块）：**settings**、**operationlog**（repository + retention，含 `INSERT OR IGNORE` → `INSERT ... ON CONFLICT (...) DO NOTHING` 改写）、**authsession**（accounts/account_operations/users/roles/notifications/service_credentials + `withTx` helper）+ **systemdata**（bootstrap/reconcile）、**wallet**（store 深层 `rowQueryer`/`ReconcileOnceTx` + service + jobs.go）、**jobs**（`TxRunner`、`CompleteWithCommit`、`CommitFunc` → `func(kernel.Tx)`）。
- **D 回调链收口（S3）**：auth `ServiceCredentialUseTxRecorder`、authsession `ServiceCredentialAudit/RevokeAudit`、operationlog `TransactionalRecorder`、handler `ServiceCredentialOperations`、composition `RecordOperationTx` 全部 → `kernel.Tx`。
- **scanOperation/scanJob** 等本地 `sql.NullString`/`rowScanner` 读到内核/标准类型；`sql.ErrNoRows` → `kernel.ErrNoRows`。
- 测试同步（fake runner/CommitFunc/audit 回调）；`go test ./...` **0 FAIL**（含 live PG）。
- **残留（S2 语义债）**：运行时 `LIKE`（wallet/recyclebin）大小写与 `ORDER BY … COLLATE NOCASE` 查询侧尚未逐处改写为 PG 显式语义——两方言行为差异点已定位（PG LIKE 大小写敏感 / PG 上 `COLLATE NOCASE` 无此 collation；CITEXT 列已建）。该处改写为行为决策，作为 S2 收尾项（或并入 S5 验收前消除）。

### 证据

| 主张 | 路径 / commit |
|------|---------------|
| S2 全模块迁移（2 批） | `9d5a97d`（authsession/operationlog/settings + D 链）；`1cf5320`（wallet/jobs + handler test） |
| INSERT OR IGNORE 改写 | `modules/operationlog/retention.go`（ON CONFLICT DO NOTHING，两方言可移植） |
| 全量回归 | `apps/api`: `go test ./...` 0 FAIL（含 live PG） |
| 残留项 | wallet/recyclebin `LIKE`、users/roles `ORDER BY … COLLATE NOCASE`（运行时查询侧） |
