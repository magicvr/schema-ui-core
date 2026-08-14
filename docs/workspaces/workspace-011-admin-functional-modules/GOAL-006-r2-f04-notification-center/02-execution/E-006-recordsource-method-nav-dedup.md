---
id: E-006
goal: GOAL-006-r2-f04-notification-center
date: 2026-08-14
status: recorded
parent: GOAL-006-r2-f04-notification-center
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# E-006 · recordSource.method 缺失 + 顶部「通知」入口重复修复

## 事实（用户报告，2026-08-14）

1. **recordSource construction failed (recordSource.method): MISSING_RECORD_SOURCE_METHOD**
   - 根因：notifications.json 与 account.json 的 `recordSource` 缺 `method` 字段（host 渲染器 fail-closed 要求 method；settings.json 基线有 GET）。
   - 修复：两处补 `"method": "GET"`；全量扫描所有 schema 的 recordSource，无其它遗漏。

2. **顶部导航「通知」重复**
   - 根因：shell 顶部有硬编码 NotificationBell 铃铛（F-04 设计：未读徽标 + 下拉 + view-all），同时 notifications 模块 manifest fragment 在 `navigation.user` 又声明了一个 "Notifications"（bell 图标）入口 → 顶部出现两个「通知」。
   - 修复：fragment `navigation.user` 清空（`user: []`）；admin/mvp fixture 同步移除 user slot 的 notifications 项；STATIC_MANIFEST_SHA256 重钉（cd7e5df7…）。唯一入口保留为 shell NotificationBell（view all → /notifications 页）。
- 验证：vitest 903/903（s5 渲染 notifications 页通过 recordSource 门禁）；go 相关模块绿；e2e admin localization+shell 绿。
