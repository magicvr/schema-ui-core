---
id: E-013-backup-port-methods-decision
doc: execution-entry
status: recorded
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-013 · Backup Port 最小方法裁决

用户选择 kernel Port 仅暴露 `CreateRecoveryPoint`，成功返回必须已经完成最低验证；Verify/RestoreTo/provider/tool orchestration 留在 internal BackupService。

证据：Root D-010、child D-009。
