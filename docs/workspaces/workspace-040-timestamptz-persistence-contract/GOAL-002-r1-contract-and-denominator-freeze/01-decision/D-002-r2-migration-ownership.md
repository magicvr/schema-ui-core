---
id: D-002-r2-migration-ownership
doc: decision-entry
status: accepted
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-002 · R2 migration 归属承接

承接 Root D-004 用户裁决：R2 按模块追加 migration；各模块拥有自己的 conversion version/Apply/ApplyPostgres；公共 codec/test contract 可共享；历史 v1–v72 canonical SQL/checksum 不改，转换从 v73 之后 append-only。

R1 C2/C3 方案必须把此不变量写入：模块 owner、版本分配、checksum、SQLite table rebuild、PG `USING`、CHECK/partial index/predicate 顺序、备份/回滚与公共 wire 影响。
