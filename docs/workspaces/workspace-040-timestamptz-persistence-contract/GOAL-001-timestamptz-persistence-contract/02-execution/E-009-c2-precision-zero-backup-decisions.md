---
id: E-009-c2-precision-zero-backup-decisions
doc: execution-entry
status: recorded
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-009 · C2 精度、D0 与 PG backup provider 用户裁决

用户选择：纳秒余数统一截断到微秒；mail/telegram config legacy 0 → NULL 并去掉 default 0；PG provider 固定 `pg_dump -F c` + `pg_restore`。具体逐列 mapping、Port API、metadata/restore script 与测试仍待 C2/C3 设计。

证据：Root `01-decision/D-008-c2-precision-zero-backup-decisions.md`。
