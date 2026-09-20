---
id: A-003-response-to-closeout-audit
doc: audit-entry
status: active
parent: GOAL-005-r2-backup-port-and-closeout
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
source: self
verdict: pass
---

# A-003（self · 编排器响应）· R2 关门审计 A-002

- **source**: self（编排器汇总响应）
- **日期**: 2026-09-20
- **scope**: 响应本目标 A-002（independent · grok-build grok-4.6 · high · `conditional` · open required = 3）的全部意见
- **verdict**: `pass`（3 条 required 与 3 条 recommended 均已处置；待 independent 复审确认）

## required 闭合

| finding | 级别 | 处置 | 证据 |
|---------|------|------|------|
| A-002 `F-I-001`：PG 类 C dump 固定文件名 → 重试会挡住「可续跑」 | high | **fixed** | `postgres.snapshotBeforePendingPG` 的文件名改为 `<db>.pre-v%04d-<ts(ms)>-<hex4>.dump`（`time.Now().UTC().Format("20060102T150405.000")` + `randomSuffix()`），与 SQLite 侧 D5 的毫秒精度+唯一化对齐；`PgProvider.Create` 的「拒绝覆盖」语义因此不再与「可续跑」冲突。证据：`TestC3RecoveryAnchorsOnPostgresUpgrade`（真实 PG，A=1/C=1/B=1）与 `internal/composition` 全绿 |
| A-002 `F-I-002`：composition 未注入 Port/creator → §4.3 在生产入口是空操作 | med | **fixed** | `composition.openStore` 新增 `recoveryWiring(dialect, cfg)`：构造 `backup.NewService(artifactDir)`，PG 方言注册 `backup.PgProvider{AdminDSN: cfg.DBDSN, ClientImage: backup.DefaultPgClientImage, WorkDir: artifactDir}` 并同时作为 `RollbackArtifacts`；artifact 目录 = `<DBPath 同级>/recovery`（sqlite）或用户缓存目录下按 DSN 派生的键（PG）。`OpenOptions` 因此收到 `RecoveryPoints` / `RollbackArtifacts` / `ArtifactDir`，生产启动不再静默跳过锚点。证据：`go test ./internal/composition/` 全绿（57s，含真实库启动路径） |
| A-002 `F-I-003`：PG `SampleVerified` 只是类型检查的别名 | med | **fixed** | 新增 `verifyPostgresSamples`：检查 legacy-0 sentinel 列（`mail_config.updated_at` / `telegram_config.updated_at` / `users.locked_until` / `login_failures.locked_until`）在恢复库中**全部为 NULL**，并逐列抽样要求回读值为**微秒精度**（`Nanosecond()%1000 == 0`）；service 的 PG 分支改为调用它并把 `SampleVerified` 置真。证据：`TestPGRestoreToNewDB`（真实 PG）通过；实现见 `internal/backup/verify.go` |

## recommended 处置

| finding | 处置 | 证据 |
|---------|------|------|
| `F-I-004`：`HasPoint` 只看 sidecar 是否存在，不解析内容 | **fixed** | `probeRecoveryState` 改为经 `readRecoveryMarker` 解析；解析失败时**不算**已满足（`HasPoint=false`）并把原因写入 `Detail`（`recovery marker is unreadable: …`） |
| `F-I-005`：`verifyIntegrityPG` 把阈值写死 87 | **fixed（改为数据驱动）** | 新增 `conversionCompletionVersion(catalog)`：以冻结 descriptor 台账的**最后一个转换描述符**（`vp040_temporal_digital_offer`）为「全量转换必须成立」的版本；被裁短的历史（v1–v72、v1–v86）不触发全量要求，因此**部分历史仍是合法状态**。证据：`TestPGLegacyArtifactMustFail`（catalog[:72] 需通过启动）与 `TestC3RecoveryAnchorsOnPostgresUpgrade`（v86→v87）同时绿 |
| `F-I-006`：PG 缺类 C 混合形状的 §6 反向断言 | **fixed** | 新增 `TestPGMidBatchArtifactMustFail`（真实 PG：v1–v73 混合形状 dump → service 以 `TimeContractMismatch` 拒绝，且显式排除 NotFound/Unreadable/ToolFailure） |

## 未闭合 / 保留

- 本响应**不**自证关门：required 闭合需 independent 复审确认（本目标检查点 C 的关门审计）。
- `I-041-004`（PG 15/16/17 跨版本矩阵）仍为 R3 前复核项；本批实测组合 server 15.4 + 客户端 15.19 已记录。
- A-002 未提出的其余自审偏差（见 A-001 §偏差）保持记录状态。
