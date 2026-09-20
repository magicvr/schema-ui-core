---
id: E-014-v73-allocation-negative-truncation-decisions
doc: execution-entry
status: recorded
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-014 · v73 baseline 与负时间解释裁决

用户接受 v73–v87 allocation 作为 R2 baseline（未发布前可有记录调整/拆分，发布后 append-only），并冻结负 instant 的向零微秒截断解释：Go codec 处理 raw fractional time，整数 legacy PG conversion 不依赖 typmod rounding；0 sentinel only → NULL。

证据：Root D-014/D-015。
