---
id: E-001-retention-implementation
goal: GOAL-008-audit-log-retention-settings
doc: execution-entry
record_id: E-001
status: recorded
parent: null
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
---

# E-001 · 设置字段、归档表、sweeper 与定向测试

## 已完成事实

1. 迁移 0046：`site_settings` 增加 `operation_log_retention_days`（默认 90）与 `operation_log_expiration_action`（默认 archive）。
2. 迁移 0047：`operation_log_archive` + `operation_log_archive_correlation`。
3. Settings GET/PATCH/reset 与 Audit log 页签（inputNumber + select）已接线；i18n zh-CN/en-US 已补。
4. `ApplyRetention`：archive 先拷后删热表；delete 只删。Sweeper 每小时读当前设置，不硬编码。
5. 定向测试：settings repository/handler、operationlog retention、store catalog、web schema-keys / representative-pages。

S2 审计尚未做。
