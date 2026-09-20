---
id: D-012-v73-allocation-negative-truncation
doc: decision-entry
status: accepted
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-012 · v73 allocation / negative truncation 承接

承接 Root D-014/D-015：v73–v87 是 R2 当前 baseline，未发布前可有记录地拆分/调整；发布后 append-only。负 fractional `time.Time` 由 Go codec 向零截断；legacy integer sec/ms PG conversion 使用 date_trunc + integer interval，不依赖 typmod rounding。
