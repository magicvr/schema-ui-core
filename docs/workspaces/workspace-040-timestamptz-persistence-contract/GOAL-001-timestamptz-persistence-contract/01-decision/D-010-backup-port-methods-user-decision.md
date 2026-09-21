---
id: D-010-backup-port-methods-user-decision
doc: decision-entry
status: accepted
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-010 · Backup Port 最小方法裁决

用户选择 kernel Backup/RecoveryPoint Port 仅暴露 `CreateRecoveryPoint`。成功返回的后置条件必须包含最低验证已完成，不能返回未验证 artifact。Verify、RestoreTo 与 provider/tool orchestration 留在 internal BackupService；未来只有出现真实 kernel 消费者时才扩展 Port。
