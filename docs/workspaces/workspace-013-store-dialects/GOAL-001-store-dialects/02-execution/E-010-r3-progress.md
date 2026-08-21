---
id: E-010
doc: execution-entry
goal: GOAL-001-store-dialects
status: recorded
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# E-010 · R3（GOAL-004）主体实施：T1/T2a/T3（12/13 模块）

## 2026-08-20 · R3 推进

### 已发生事实

- R3（GOAL-004，progress 3/5）：
  - **T1**：`MigrationContribution.Apply/Reconcile` → `kernel.Tx`（14 迁移包；commit `8932148`）。
  - **T2a**：postgres 迁移运行器 live 证明（commit `4359c7f`）。
  - **T3**：`ApplyPostgres` 管道 + **12/13 模块**对写（BIGINT 时间/金额、CITEXT、partial index）；**全量 compiled catalog 在 live PG fresh bootstrap + 台账 + 幂等 + 20 列 BIGINT 断言全绿**（commit `8baca38` + `ad3a876`）。`operationlog` 对写待 T3 收尾（A-003 F-001 required）。
- sqlite 回归 0 FAIL（sqlite Apply/checksum 未改）。全量 `go test ./...` 0 FAIL（含 live PG）。

### 证据

| 主张 | 路径 / commit |
|------|---------------|
| R3 T1/T2a/T3 主体 | `docs/workspaces/workspace-013-store-dialects/GOAL-004-*/02-execution/E-00{2,3,4}-*.md`、`03-audit/A-00{1,2,3}-*.md`；commits `8932148`/`4359c7f`/`8baca38`/`ad3a876` |
| Root R3 行 | `GOAL-001-store-dialects/00-meta.md` |
