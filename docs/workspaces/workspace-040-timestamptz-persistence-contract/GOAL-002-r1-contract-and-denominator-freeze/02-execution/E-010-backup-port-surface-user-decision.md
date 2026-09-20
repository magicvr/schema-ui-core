---
id: E-010-backup-port-surface-user-decision
doc: execution-entry
status: recorded
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-010 · Backup/RecoveryPoint Port 公共边界裁决

用户选择最小 Backup/RecoveryPoint Port 进入 kernel；BackupService orchestration、native providers、metadata/verification、restore-to-new-db、调度/权限/远程存储/保留策略留在 internal。C3 仍需冻结 Port API 与失败语义。

证据：Root `D-007-backup-port-surface-user-decision.md`、child `D-006-backup-port-surface.md`。
