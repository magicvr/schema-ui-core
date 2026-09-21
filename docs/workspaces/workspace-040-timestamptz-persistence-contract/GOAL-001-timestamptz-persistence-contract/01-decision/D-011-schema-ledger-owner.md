---
id: D-011-schema-ledger-owner
doc: decision-entry
status: accepted
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-011 · schema_migrations conversion owner 裁决

用户选择 `core.persistence` 作为 `schema_migrations.applied_at` conversion 的 platform owner。Store runner 仍只执行 compiled catalog/ledger，不直接改写历史 migration；authsession v1 的历史创建归属保留，但不承担 VP-040 conversion owner。
