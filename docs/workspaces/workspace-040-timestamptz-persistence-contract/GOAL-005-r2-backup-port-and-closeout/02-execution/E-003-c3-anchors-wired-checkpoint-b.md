---
id: E-003-c3-anchors-wired-checkpoint-b
doc: execution-entry
status: active
parent: GOAL-005-r2-backup-port-and-closeout
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-003 · C3 §4.2 调用点全部接线（检查点 B）

## 事实

- **用户裁决落盘**：`01-decision/D-001-partial-upgrade-and-c3-anchor-scope.md` —— `I-041-006` **accepted**（保留每-descriptor 事务，可续跑 + A/C 回滚；不做批级原子化）；C3 §4.2 **全部接线**。
- **SQLite 侧**（`internal/store/migrate.go`、`recovery.go`）：
  - 批前 A：`snapshotBeforeBatch`（`VACUUM INTO` → `<db>.batch-rollback-<ts>.sqlite` + 完整性校验）；`pending == 0` 或 fresh/内存库时跳过；失败 → 批次不开始。
  - 批次内 C：既有 `snapshotBeforePending` 保留（`<db>.pre-v%04d-<ts>.sqlite`）。
  - 批次后 B：`createRecoveryPoint` 调 `kernel.RecoveryPointPort`，成功后写 sidecar marker `<db>.recovery-point.json`（含 ID / ArtifactRef / CatalogVersion / ChecksumSet / VerifiedAt）。
  - 形状校验：`sqliteConvertedShape`（90 列 TEXT + `records` 缺席）。
- **PG 侧**（`internal/store/postgres.go`、`recovery.go`）：
  - 批前 A / 批次内 C：经注入的 `store.RollbackArtifactCreator`（由 `backup.PgProvider` 实现，`CreateRollbackArtifact` → `pg_dump -F c`）落到 `ArtifactDir`；未配置时记 note 并跳过（不静默）。
  - 批次后 B：同 SQLite，但 marker 写为 `COMMENT ON DATABASE`（`vp040-recovery-point:<base64url-json>`）——**无 schema 变更**。
  - **新增** `verifyIntegrityPG` + `postgresConvertedShape`（补 C3 §4.2 点名的「PG 侧无 `verifyIntegrity`」缺口）。
- **§4.3 可检测性与有界重试**：`dbIdentity.RecoveryPointMissing` + 新动作 `actionVerifyRecoveryPoint`；`planStartup` 在「catalog 到 head 且无 marker」时返回该动作（不再 `actionNoop`）；动作内**每次启动最多补创一次**；B 失败**不阻断启动**但记 `RecoveryNote` 且**不写 marker**。
- **依赖方向**：store 仅依赖窄接口 `store.RollbackArtifactCreator` + `kernel.RecoveryPointPort`，**不导入** `internal/backup`（provider 由调用方注入），避免包环；`internal/backup` 新增 `PgProvider.CreateRollbackArtifact`。
- **共享分母**：90 列分母抽到新包 `internal/temporalcontract`（由 `backup` 与 `store` 共用，避免 store↔backup 循环）。

## 证据（可执行）

- `TestC3RecoveryAnchorsOnUpgrade`（SQLite，注入 fake port/creator）：A 恰好 1 个 `batch-rollback-`、C 至少 1 个 `pre-v`；B 调用 1 次且 marker 写入（`CatalogVersion` 正确）；**有 marker 再开 = noop（B 不再调用）**；删 marker 再开 = **恰好 1 次**补创并重写 marker。
- `TestC3RecoveryPointFailureNeverBlocksStartup`：B 失败 → 启动成功、`RecoveryNote` 记录、**无 marker**、`RecoveryPointState().HasPoint == false`（不假装门禁满足）。
- `TestC3RecoveryAnchorsOnPostgresUpgrade`（**真实 PG 15.4 + 容器 `pg_dump`**）：v86→v87 升级 → A=1、C=1、B 真实生成（`pg_dump` → 新库 `pg_restore` → 校验）→ DB 注释 marker 写入；重开 = noop；清 marker 后重开 = 显式动作补创成功。
- `go build ./...` exit 0；`go vet ./internal/store/ ./internal/backup/` 无输出。
- 全仓 `go test -count=1 ./...`：见 `02-execution.md` 事实边界（本轮已复跑）。

## 未完成（诚实边界）

- **检查点 C（R2 关门审计）未开始**：需 self + independent（grok build）关门审计落盘、required 合法闭合。
- 本文件**不**声称 R2 已放行。
