---
id: E-007-backup-spi-user-decision
doc: execution-entry
status: recorded
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-007 · Backup SPI/Service 用户裁决

用户选择在 VP-040 C3 建立统一 Backup SPI/Service 合同，由 SQLite/PG provider 使用各自原生备份机制，统一元数据、验证与 restore-to-new-db；迁移失败优先事务 rollback；不实现调度、权限、远程存储、保留策略、KMS/TLS 或完整备份产品。

证据：Root `01-decision/D-006-r1-backup-spi-user-decision.md`。
