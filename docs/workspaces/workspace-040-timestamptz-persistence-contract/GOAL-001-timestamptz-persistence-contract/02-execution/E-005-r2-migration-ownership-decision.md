---
id: E-005-r2-migration-ownership-decision
doc: execution-entry
status: recorded
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-005 · R2 migration 归属用户裁决

用户选择 R2 按模块追加 migration：每个模块拥有自己的 conversion migration version/Apply/ApplyPostgres；公共 codec/test contract 可共享；v1–v72 canonical SQL/checksum 不改，conversion 从 v73 之后 append-only。该裁决尚不等于 C2 设计已完成；逐模块 migration 顺序、codec、predicate/index rebuild、backup/rollback 仍由 R1 C2/C3 证据冻结。

证据：Root `01-decision/D-004-r2-migration-ownership-user-decision.md`。
