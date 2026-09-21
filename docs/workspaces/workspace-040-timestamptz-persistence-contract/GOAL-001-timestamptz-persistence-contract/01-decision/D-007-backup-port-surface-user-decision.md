---
id: D-007-backup-port-surface-user-decision
doc: decision-entry
status: accepted
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-007 · Backup/RecoveryPoint Port 公共边界裁决

用户选择：

- 只把**最小 Backup/RecoveryPoint Port** 纳入 `kernel` 公共契约，供 composition/运维边界调用。
- `BackupService` orchestration、SQLite/PG provider、metadata/verification、restore-to-new-db、调度、权限、远程存储、保留策略、KMS/TLS 与具体 native 机制全部留在 `apps/api/internal`。
- 不把整个 Backup 子系统公共化；不增加模块对 provider/driver 的依赖。
- Port 的最小接口、metadata、artifact/restore verification 与错误/失败语义仍需 C3 设计并经 self + grok independent 审计。
