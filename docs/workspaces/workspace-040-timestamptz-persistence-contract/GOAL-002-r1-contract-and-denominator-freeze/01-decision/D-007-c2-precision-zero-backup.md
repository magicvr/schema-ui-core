---
id: D-007-c2-precision-zero-backup
doc: decision-entry
status: accepted
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-007 · C2 精度、D0 与 PG backup provider 承接

承接 Root D-008：纳秒统一截断到微秒；mail/telegram config legacy 0 → NULL 并移除 default 0；PG backup 固定 `pg_dump -F c` + `pg_restore`。C2/C3 仍需把规则展开为逐列 mapping、Port 方法、metadata 与 restore verification。
