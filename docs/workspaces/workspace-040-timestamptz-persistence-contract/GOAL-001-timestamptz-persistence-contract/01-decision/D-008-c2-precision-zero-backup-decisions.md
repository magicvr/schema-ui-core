---
id: D-008-c2-precision-zero-backup-decisions
doc: decision-entry
status: accepted
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-008 · C2 精度、D0 与 PG backup provider 裁决

用户 P-004 选择：

1. **精度**：统一截断到微秒。`timestamptz(6)` 与 SQLite fixed-6 TEXT 写入/迁移时，纳秒余数向零截断；不得在不同模块静默采用 round。
2. **配置 D0**：`mail_config.updated_at`、`telegram_config.updated_at` 中 legacy `0` 转为 `NULL`；放宽 nullable、去掉 `DEFAULT 0`，读取层把 NULL 视为未初始化/无更新时间。
3. **PG backup provider**：固定使用 `pg_dump -F c` + `pg_restore`；provider metadata/verification 固定工具主版本约束与 restore-to-new-db 校验，不并行引入 `pg_basebackup` provider。

这些裁决冻结方向，仍需 C2/C3 形成逐列 mapping、Port 方法、metadata schema、restore script、predicate/constraint 重建与测试证据。
