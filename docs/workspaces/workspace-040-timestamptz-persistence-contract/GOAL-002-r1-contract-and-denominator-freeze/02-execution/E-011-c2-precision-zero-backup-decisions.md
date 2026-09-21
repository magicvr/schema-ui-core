---
id: E-011-c2-precision-zero-backup-decisions
doc: execution-entry
status: recorded
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-011 · C2 精度、D0 与 PG backup provider 用户裁决

用户已裁决：纳秒余数截断到微秒；mail/telegram config 0 → NULL；PG backup provider 固定 `pg_dump -F c` + `pg_restore`。这些是 C2/C3 方向，具体实现证据尚未落盘。

证据：Root D-008、child D-007。
