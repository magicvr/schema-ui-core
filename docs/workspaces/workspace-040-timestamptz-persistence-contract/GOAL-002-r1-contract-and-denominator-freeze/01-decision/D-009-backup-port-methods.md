---
id: D-009-backup-port-methods
doc: decision-entry
status: accepted
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-009 · Backup Port 最小方法承接

承接 Root D-010：kernel Port 仅暴露 `CreateRecoveryPoint`，成功返回必须完成最低验证；Verify/RestoreTo/provider/tool orchestration 留在 internal BackupService。C3 需冻结 request/recovery metadata、后置验证条件与错误语义。
