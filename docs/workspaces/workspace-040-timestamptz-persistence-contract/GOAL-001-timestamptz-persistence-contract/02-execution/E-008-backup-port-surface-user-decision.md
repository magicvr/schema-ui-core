---
id: E-008-backup-port-surface-user-decision
doc: execution-entry
status: recorded
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-008 · Backup/RecoveryPoint Port 公共边界裁决

用户选择将最小 Backup/RecoveryPoint Port 纳入 kernel 公共契约；BackupService orchestration、native providers、metadata/verification、restore-to-new-db 与完整运维/权限/调度能力留在 internal。Port 细节仍待 C3 设计与审计。

证据：Root `01-decision/D-007-backup-port-surface-user-decision.md`。
