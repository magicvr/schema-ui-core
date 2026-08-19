---
id: D-001-configurable-retention
doc: decision-entry
status: accepted
parent: GOAL-008-audit-log-retention-settings
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
---

# D-001 · 可配置审计日志保留

用户书面要求：保留天数和过期删除/归档做成设置；设置页给入口；给合适默认；管理员可改；不要硬编码数值。

## 冻结

| 项 | 值 |
|----|----|
| 默认天数 | **90**（常见管理端审计窗口；可改） |
| 默认过期动作 | **archive**（先冷存再离热表，避免默认销毁） |
| 天数范围 | 1–3650 |
| 动作 | `archive` \| `delete` |
| 存储 | `site_settings.operation_log_retention_days` / `operation_log_expiration_action` |
| 设置入口 | Settings → Audit log 页签 |
| 执行 | 启动时跑一次，之后每小时读当前设置再扫 |

未选：硬编码 90；默认直接 delete；0 = 永久保留（避免「关保留」无入口）。
