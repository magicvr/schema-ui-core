---
id: E-002-backup-port-and-providers-checkpoint-a
doc: execution-entry
status: active
parent: GOAL-005-r2-backup-port-and-closeout
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-002 · Backup Port + provider + restore harness（检查点 A，含 PG 路径）

## 事实

- **kernel Port（Root `D-007` + `D-010` 最小面）**：`apps/api/kernel/backup.go`
  - `RecoveryPointPort` 只导出 `CreateRecoveryPoint`；
  - 中性类型 `RecoveryPointRequest` / `RecoveryPoint` / `VerificationSummary`（`Passed()` 四项全真）；
  - 常量 `TimeContractVP040Timestamptz`、`ContractShapeConverted` / `ContractShapeLegacy`；
  - **无驱动类型**（`ArtifactRef` 不透明；`Dialect` 用既有 `kernel.Dialect`）。
- **internal/backup**：
  - `backup.go`：C3 §5.1 的**冻结错误分类**（`TimeContractMismatch` / `TemporalColumnSetIncomplete` / `ChecksumMismatch` / `ArtifactNotFound` / `ArtifactUnreadable` / `ToolFailure` / `InvalidRequest` / `SampleMismatch`）+ `KindOf` 沿 wrap 链取分类 + ledger 指纹（`sha256` over 排序后的 `version:name:checksum`）。
  - `temporal_columns.go`：**冻结 90 列分母**，由 `r1-time-column-inventory-v0.3.md` **机械解析生成**（90 行全部解析成功，无手抄）。
  - `verify.go`：SQLite 实测（`PRAGMA integrity_check` / `foreign_key_check` / 逐表 `table_info` 的 90 列类型 / `records` 缺席 / ledger 指纹 / batch version）+ 样本校验（sentinel 0→NULL、canonical fixed-6 逐值回环）+ PG 实测（`information_schema` 的 `timestamp with time zone` 且 `datetime_precision = 6`、`records` 缺席、ledger 指纹）；分类顺序按 C3 §5.1（缺列 → `TemporalColumnSetIncomplete`；形状不符 → `TimeContractMismatch`）。
  - `provider_sqlite.go`：`VACUUM INTO` 生成 artifact（类 B 只有在转换提交并校验后才成立，provider 自身不作此声明）、`Restore` 到**新文件**（`O_EXCL`，不做就地 restore）+ `integrity_check`。
  - `provider_pg.go`：`pg_dump -F c --no-owner` / `pg_restore --exit-on-error --no-owner`，客户端工具由**固定版本容器镜像**提供（Root `D-017` §3 约束③；默认 `postgres:15-alpine`，与常驻 server 15.x 同主版本）；restore 到**新建数据库**并返回 `DROP DATABASE … WITH (FORCE)` 清理；错误信息**脱敏** DSN 口令。
  - `service.go`：provider 选择 → 生成 → **restore 到新目标** → 校验 → 只在 `VerificationSummary.Passed()` 时返回 `RecoveryPoint`；否则返回分类错误且**不返回** RecoveryPoint（Root `D-010` 后置条件）。`var _ kernel.RecoveryPointPort = (*Service)(nil)` 编译期断言。
- **restore harness（C3 §5 / §5.1 / §6）**：
  - `TestSQLiteRestoreToNewDB`：新建转换库 → 生成 recovery point → restore 到新文件 → 独立复测（integrity ok、FK 0、90/90 列 TEXT、ledger 指纹一致）+ 秒/毫秒样本逐一核对。
  - `TestLegacyArtifactMustFail`：**A 类**（预转换库）→ 校验必须以 `TimeContractMismatch`（或 `TemporalColumnSetIncomplete`）失败，并**显式断言不属于** `ArtifactNotFound` / `ArtifactUnreadable` / `ToolFailure`；同时断言 Service 不会对预转换源返回 recovery point。
  - `TestPerMigrationSnapshotIsNotARecoveryPointSource`：**C 类**（批次中途：v1–v73 已应用 → 混合形状）同样被拒。
  - `TestMissingArtifactIsNotAContractFailure`：缺文件必须归 `ArtifactNotFound`（分类边界反向锁定）。
  - `TestPgProviderCommandConstruction`：工具调用参数含 `-F c` / `--no-owner` / 挂载点，且工具失败信息**不含**口令。
  - `TestPGRestoreToNewDB`（真实 PG 15.4 + 容器客户端 `pg_dump 15.19`）：pg_dump → `pg_restore` 到新库 → 90 列 `timestamptz(6)` + ledger 指纹 + batch version 全部通过；客户端版本被记录（C3 §5「记录组合，不默认通过」）。
  - `TestPGLegacyArtifactMustFail`：**真实**预转换 dump 被以 `TimeContractMismatch` 拒绝，错误文本逐列列出 90 个 `bigint` 列（可核对证据）。

## 验证

- `go build ./...` exit 0；`go vet ./internal/backup/ ./kernel/` 无输出。
- `go test -count=1 ./internal/backup/` 全绿（含真实 PostgreSQL 15.4 + Docker 固定版本客户端路径，~19s）。

## 未完成（诚实边界）

- **C3 §4.2 的 before/after 调用点尚未接线**：批前 A（双方言）、批次成功后 B（`CreateRecoveryPoint`）、PG 侧对称的 snapshot/verifyIntegrity 缺口仍存在；§4.3 的「B 缺失补创」动作未落码。是否在本子目标内接线，须编排器与用户按 P-004 定（列入 `01-decision` 待决）。
- **`I-041-006`（部分升级语义）未裁决**：descriptor 各自事务下 v85 类失败会留下 v73–v84 已提交；runner 的现状是「下次启动从部分状态继续」（可续跑），回滚依赖 A/C artifact。
- 本文件**不**声称 R2 已放行；检查点 C（R2 关门审计）未开始。
