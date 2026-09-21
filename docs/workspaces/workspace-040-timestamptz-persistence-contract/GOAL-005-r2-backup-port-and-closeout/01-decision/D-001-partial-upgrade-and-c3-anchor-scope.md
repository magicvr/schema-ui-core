---
id: D-001-partial-upgrade-and-c3-anchor-scope
doc: decision-entry
status: accepted
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-001 · 部分升级语义与 C3 调用点范围（用户 P-004 裁决）

## 决定的来源

- **用户 2026-09-20 裁决**（`/govern` 询问，两问同轮）：
  1. `I-041-006`（部分升级/批级快照）：**接受现有语义** —— descriptor 各自事务；批次失败后下次启动**从部分状态继续（可续跑）**，回滚依靠 C3 的 A/C artifact。**不**引入批级单事务原子性。
  2. C3 §4.2 调用点：**全部接线** —— 批前 A、批次成功后 B（`CreateRecoveryPoint`）、§4.3「无 B 则有限重试」、PG 侧对称的 snapshot 与形状校验入口。
- 触发依据：`GOAL-004/03-audit/A-002-independent-m3-implementation.md` 待复审事项 7；Root `D-016` §2 第 9 项；冻结 C3 边界 `r1-c3-backup-recovery-boundary-v1.0-fc.md` §4.2/§4.3。

## 1. `I-041-006` 关闭（`accepted`）

- **语义（冻结为事实）**：`applyPending` / `applyPendingPG` 每条迁移一个事务；批次中途失败时，**已提交的迁移保持已提交**，ledger 记录前缀，**下次启动从该前缀继续**（`planStartup` → `actionApplyPending`）。
- **回滚路径**：A（批前旧合同 artifact）与 C（per-migration 快照）承担回滚；B 只在批次提交并校验后存在。
- **不做**：批级单事务原子化（会改变事务边界、使 per-migration 快照失去意义、在 PG 上把整批 DDL 锁持有到提交）。
- 运维含义必须与 C3 恢复边界一起读取：部分升级**不是**损坏状态，但**不得**被当作「转换已完成」；`verifyIntegrity`/`verifyIntegrityPG` 与形状探测负责区分。

## 2. C3 §4.2 调用点（全部接线）

| 调用点 | SQLite | PG |
|--------|--------|----|
| 1 批前 A | `Store.snapshotBeforeBatch`（`VACUUM INTO`，文件名 `<db>.batch-rollback-<ts>.sqlite`；失败 → 批次不开始） | `postgres.snapshotBeforeBatchPG`（经注入的 `RollbackArtifactCreator` = `PgProvider`，`<db>.batch-rollback-<ts>.dump`；未配置 creator/dir 时记 note 并跳过） |
| 2 批次内 C | `snapshotBeforePending`（既有，`<db>.pre-v%04d-<ts>.sqlite`） | `snapshotBeforePendingPG`（`<db>.pre-v%04d.dump`），fresh 库跳过 |
| 3 批次后 B | `Store.createRecoveryPoint` → `kernel.RecoveryPointPort.CreateRecoveryPoint`；成功后写 sidecar marker `<db>.recovery-point.json` | 同左，但 marker 写为 `COMMENT ON DATABASE`（`vp040-recovery-point:<base64url-json>`，无需 schema 变更） |
| 4 迁移事务内 | **禁止**（未接线） | **禁止**（未接线） |
| 形状校验 | `verifyIntegrity`（既有）+ `sqliteConvertedShape` | **新增** `verifyIntegrityPG` + `postgresConvertedShape`（补 C3 §4.2 点名的 PG 缺口） |
| §4.3 无 B | `actionVerifyRecoveryPoint`：形状已转换 → **每次启动最多补创一次**；形状未转换 → 显式报错（回滚路径） | 同左（PG 分支） |
| §4.3 不阻断启动 | B 创建失败只记 `RecoveryNote`，启动继续；**不写 marker**，`RecoveryPointState().HasPoint` 保持 false | 同左 |

- **「无 B」不是「已满足门禁」**：`planStartup` 在 catalog 到 head 且无 marker 时返回**显式动作** `verify-recovery-point`，不再返回 `actionNoop`。
- **依赖方向**：store 只依赖窄接口 `store.RollbackArtifactCreator` + `kernel.RecoveryPointPort`，**不**导入 `internal/backup`（provider 由 composition/测试注入），避免包环。

## 未选方案

- **批级单事务**：见 §1「不做」，且触及 Root 红线（不修改历史 DDL/Apply 语义需另开决策）。**未采用**。
- **只接 B + 重试**：会把 C3 §4.2 的 R2 义务留成半成品，独立审计按现有证据很可能判 required。**未采用**（用户也选择全接）。
- **用新表存 marker**：需要新增迁移（catalog 88），会改变指纹/丢账本表集合并牵动 R1 冻结面。**未采用**：SQLite 用 sidecar 文件、PG 用数据库注释，均无 schema 变更。

## 影响与边界

- 本决策**不**放行 R2 关门；检查点 C（self + independent 关门审计）仍未开始。
- `I-041-004`（PG 15/16/17 跨版本 dump/restore 矩阵）仍为 R3 前复核的 `non-blocking` 项；本批实测组合 = server 15.4 + 容器客户端 15.19（同主版本），已记录。
