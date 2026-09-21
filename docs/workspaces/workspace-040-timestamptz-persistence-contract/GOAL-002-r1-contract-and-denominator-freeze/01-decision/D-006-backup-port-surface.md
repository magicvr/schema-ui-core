---
id: D-006-backup-port-surface
doc: decision-entry
status: accepted
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-006 · 最小 Backup/RecoveryPoint Port 承接

承接 Root D-007：只有最小 Backup/RecoveryPoint Port 进入 `kernel`；BackupService orchestration、native providers、metadata/verification、restore-to-new-db 与运维边界留在 `apps/api/internal`。不把 provider/driver 类型带进模块公共契约。

C3 必须冻结 Port 的最小方法、metadata/artifact/verification 类型、dialect/contract version、失败语义与调用者；完整备份产品能力排除。
