---
id: E-012-schema-ledger-owner
doc: execution-entry
status: recorded
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-012 · schema_migrations owner 用户裁决

用户选择 `core.persistence` 作为 `schema_migrations.applied_at` conversion owner；Store runner 继续执行 catalog/ledger，不直接改历史 migration；authsession v1 仅保留历史创建事实。

证据：Root D-011。
