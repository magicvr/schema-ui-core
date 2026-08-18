---
id: E-002
goal: GOAL-031-w20-notification-settings-in-account
title: S2 实施设置迁入个人中心
status: completed
created: 2026-08-18
updated: 2026-08-18
version: 0.1.0
parent: GOAL-001-design-implementation-conformance
---

# E-002 · S2 实施（2026-08-18）

## 已发生事实

1. `notifications.json` 删除设置表单与 `saveSettings`；去掉 `form.record.load`。
2. `account.json` 增加「通知」Tab（资料之后）：`switch` + `GET/PATCH /api/notifications/settings`。
3. GET/PATCH 响应 `enabled` 改为 JSON bool；PATCH 仍兼容字符串。
4. 新增 `schema.account.tab.notifications`；dval 锁定「设置不在收件箱、在个人中心」。

## 下一步（计划）

S3 记录测试。
