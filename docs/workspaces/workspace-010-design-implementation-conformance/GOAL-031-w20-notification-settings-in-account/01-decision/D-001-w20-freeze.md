---
id: D-001
goal: GOAL-031-w20-notification-settings-in-account
title: W20 方案冻结：通知设置迁入个人中心
status: accepted
created: 2026-08-18
updated: 2026-08-18
version: 0.1.0
parent: GOAL-001-design-implementation-conformance
---

# D-001 · W20 方案冻结

## 1. 触发

用户确认：通知列表页不应有「通知设置」；若要给人操作，应放个人中心。

## 2. 决定

1. **列表页**：删除 `notification-settings` 表单与 `saveSettings` action。只留搜索/筛选 + `notification-center`。去掉仅设置表单需要的 `form.record.load`。
2. **个人中心**：在资料与安全之间新增 Tab「通知」。表单回读 `GET /api/notifications/settings`，提交 `PATCH` 同一端点。字段用 `switch`（布尔），不再用 select 字符串。
3. **API**：路径与权限不变。GET/PATCH 的 `enabled` 改为 JSON bool（PATCH 仍兼容 `"true"`/`"false"` 字符串）。
4. **文案**：Tab 用 `schema.account.tab.notifications`；表单标签沿用 `schema.notifications.settings.*`。

## 3. 未选方案

| 方案 | 未选理由 |
|------|----------|
| 塞进「安全」Tab | 这是偏好，不是改密/MFA |
| 只从列表删掉、个人中心不放 | 用户仍需要操作入口 |
| 列表留「去设置」链接 | 本波不做导航交叉；个人中心已有入口 |

## 4. 后续

S2 改 schema + i18n + 契约测试。go 不暂挂。
