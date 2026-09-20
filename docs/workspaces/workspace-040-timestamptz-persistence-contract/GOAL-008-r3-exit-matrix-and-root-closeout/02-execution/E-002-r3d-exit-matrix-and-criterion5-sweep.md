---
id: E-002-r3d-exit-matrix-and-criterion5-sweep
doc: execution-entry
status: active
parent: GOAL-008-r3-exit-matrix-and-root-closeout
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# E-002 · 检查点 A：退出判据矩阵与判据 5 反向核验

## 事实

- **退出矩阵落盘**：`attachments/r3d-root-exit-criteria-matrix-v0.1.md` —— Root 六条成功标准逐条落成「判据 → 主张 → 证据 → 结论」，每条指向已关门目标的可核对产物。
  - 判据 1（合同冻结）：`GOAL-002` `D-002`/`D-003`、Root `D-003`/`D-005`/`D-009`/`D-015`、22 份冻结附件、`internal/temporalcontract`（`Count = 90` / 44 表）+ `internal/temporal` codec。
  - 判据 2（双方言迁移 + checksum + 非时间 `INTEGER` 排除）：15 个 v73–v87 descriptor；**v73–v87 的 15/15 checksum 冻结在 `internal/store/migrate_test.go`**（本轮实测 15 of 15 present）；真实迁移边界矩阵 `internal/w040contracttest/migration_boundaries_test.go`。
  - 判据 3（写入读回一致 / VP-020 不漂移 / 含真实 PG 路径）：写入截断 + 读侧 fail-closed 适配器 + R3-A/B wire formatter、parser、单位族矩阵、VP-020 round-trip；**真实 PG 路径本轮实测非 skip**。
  - 判据 4（可升级后恢复）：SQLite/PG restore harness + **升级路径** `TestC3RecoveryAnchorsOnPostgresUpgrade` + R3-C 跨版本矩阵 18 个 supported 格（含 15→16/17）；residual 台账 `F-I-005` 已 `fixed`（`GOAL-002/A-048`），不得重记为 open。
  - 判据 5（反向核验）：本轮亲自扫描，见下。
  - 判据 6（独立意见 + 开放 required + 用户确认）：六个已关门目标关门向 independent 裁决与开放 required = 0 汇总；**用户确认仍待**（`I-041-010`）。
- **判据 5 反向核验（本轮执行）**：
  - 依赖清单在 workspace-040 全程**零变更**（`go.mod`/`go.sum`/`package.json`/lockfile 均无 diff）。
  - `go.mod` 无 ORM / Redis / MQ / 第三数据库依赖字符串；`apps/web/package.json` 无同类依赖。
  - 驱动类型：`pgx.`/`modernc.org/sqlite` 仅出现在 `cmd/` 工具与适配器空导入；`*sql.Tx` 仅 `internal/store` 适配器内部与 `internal/testsupport`；`kernel` 中只出现在**注释**里。
  - `sql.NullTime` 属 `database/sql`（非驱动类型），用于 R2 已记录的 NULL 表示（`authsession.User.LockedUntil`）；handler 仅据此派生布尔，从不送上 wire。
  - 业务域/维护提示：本工作区只触碰 `kernel/backup.go`（C3 Port 表面），未向 kernel 新增业务词汇。
  - 多实例/分布式协调：108 个非测试生产文件全部落在 temporal/store/migration/handler-wire/backup/jobs 范围内。
- **真实路径实测（判据 3/4）**：`TestPGRestoreToNewDB`（13.51s）、`TestSQLiteRestoreToNewDB`（0.74s）、`TestC3RecoveryAnchorsOnPostgresUpgrade`（19.59s）、`TestCompositionPostgresStartup`（14.13s）**全部 PASS、无一 skip**。

## 证据

- `attachments/r3d-root-exit-criteria-matrix-v0.1.md`（矩阵本体，含 §5 反向核验的命令与结果、§6 跨目标开放 required 汇总、§7 待用户裁决项）。
- 各已关门目标的 `03-audit` 索引（`GOAL-002`～`GOAL-007`）。

## 进度评估

检查点 A 完成 → `progress: 1/3`。**待办**：检查点 B（self + grok independent 关门审计）与检查点 C（**用户确认关门**，`I-041-010` required）。
