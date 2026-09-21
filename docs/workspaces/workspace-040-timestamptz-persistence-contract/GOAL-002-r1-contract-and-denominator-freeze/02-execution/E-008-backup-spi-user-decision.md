---
id: E-008-backup-spi-user-decision
doc: execution-entry
status: recorded
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-008 · C3 Backup SPI/Service 用户裁决

用户选择统一 Backup SPI/Service 的最小合同：SQLite/PG 各自 native provider；统一 metadata、verification 与 restore-to-new-db；迁移失败优先 transaction rollback；不做调度、权限、远程存储、保留策略、KMS/TLS/UI。C3 仍待具体接口与失败边界设计。

证据：Root `D-006-r1-backup-spi-user-decision.md`、child `D-004-backup-spi-contract.md`。
