---
id: E-017-v73-negative-truncation-decision
doc: execution-entry
status: recorded
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-017 · v73 allocation / negative truncation policy

用户接受 v73–v87 为未发布 R2 baseline（可记录调整/拆分），并冻结负 fractional time 的 Go codec 向零截断、legacy integer PG conversion 不依赖 typmod round。

证据：Root D-014/D-015、child D-012。
