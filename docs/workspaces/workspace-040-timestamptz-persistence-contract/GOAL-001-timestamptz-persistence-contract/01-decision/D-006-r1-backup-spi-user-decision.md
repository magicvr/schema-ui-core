---
id: D-006-r1-backup-spi-user-decision
doc: decision-entry
status: accepted
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-006 · C3 统一 Backup SPI/Service 裁决

## 用户裁决

C3 不直接扩展为完整备份产品，但本 VP 建立统一 **Backup SPI/Service 合同**：

- SQLite 与 PostgreSQL 各自由 provider 使用原生备份机制；不假定两方言备份/恢复语义相同。
- 统一备份元数据、验证与 restore-to-new-db 合同，至少能核对目标方言、schema/catalog version、checksum、类型/代表性数据抽样与 UTC 时间值。
- 迁移失败优先依赖同一事务 rollback；备份是更高一级恢复保障，不替代事务回滚。
- 本 VP 不实现调度、用户权限、远程存储、保留策略、KMS/TLS 或完整备份管理 UI；这些留待后续独立 VP。

## 门禁

C3 必须冻结 SPI/Service 的最小接口、SQLite/PG provider 责任、元数据 schema、restore-to-new-db 验证与失败边界；随后由 self + grok independent 审计。该裁决不放行实现代码，也不关闭 I-040-003。
