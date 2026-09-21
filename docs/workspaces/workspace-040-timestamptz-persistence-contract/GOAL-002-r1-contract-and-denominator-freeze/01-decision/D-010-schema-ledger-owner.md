---
id: D-010-schema-ledger-owner
doc: decision-entry
status: accepted
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-010 · schema_migrations owner 承接

承接 Root D-011：`core.persistence` 追加 `schema_migrations.applied_at` conversion；Store runner 只执行 catalog/ledger，不直接改历史版本；authsession v1 仅是历史创建来源。C2/C3 需把该 ModuleID 与 v73+ version allocation 写入 append-only plan。
