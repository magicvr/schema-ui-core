---
id: D-004-backup-spi-contract
doc: decision-entry
status: accepted
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-004 · C3 Backup SPI/Service 合同承接

承接 Root D-006：SQLite/PG 各自 native provider；统一 backup metadata、verification 与 restore-to-new-db contract；迁移失败优先事务 rollback；不做 scheduler、权限、远程存储、保留策略、KMS/TLS 或 UI。

C3 需明确：

1. provider 输入/输出与 dialect ownership；
2. metadata 至少包含 dialect、catalog/schema version、checksum、source/target、时间合同版本；
3. restore verification 包含 schema type、90-column temporal shape、代表性 seconds/milliseconds→UTC instant 样本、NULL/sentinel normalization；
4. native tool/version mismatch、unknown GUC、partial restore 与事务 rollback 的失败边界；
5. 与 VP-013 既有 `pg_dump`/`pg_restore` 与 SQLite snapshot 证据的复用/升级关系。
