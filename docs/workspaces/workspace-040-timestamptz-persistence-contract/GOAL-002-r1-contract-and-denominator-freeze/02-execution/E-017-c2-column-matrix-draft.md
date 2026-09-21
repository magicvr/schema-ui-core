---
id: E-017-c2-column-matrix-draft
doc: execution-entry
status: recorded
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-017 · C2 90 列 codec/NULL mapping matrix v0.2

已形成 `attachments/r1-c2-column-contract-matrix-v0.2.md`：以 inventory v0.3 的 90 行为 row id，按 seconds/milliseconds × NN/N/D0 六类给出 deterministic target codec、精度、sentinel/NULL 基线，并单列 `login_failures`、config D0、task_runs 与 voucher `0→NULL/negative fail closed` 特殊规则。

该矩阵仍为 proposed；逐 owner DDL/USING/rebuild、predicate/index、migration version/checksum、runtime codec 与 backup evidence 尚未完成。
