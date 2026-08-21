---
id: D-004
goal: GOAL-006-r2-f04-notification-center
title: A-003 响应 — F-001 required fixed + recommended 全落地（无 P-004 裁决冲突）
date: 2026-08-14
status: accepted
parent: GOAL-006-r2-f04-notification-center
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# D-004 · A-003（grok independent · conditional）响应

> P-004 检查：F-001 为 required 且 self（A-002）未发现。选择 **fixed**——设置页无法真正关闭开关是交付面缺陷，修复成本低、可核对。

## 处置

| finding | 处置 | 修复 |
|---------|------|------|
| F-001 required | **fixed** | `GET /api/notifications/settings`（表单面 `"true"/"false"` 字符串）+ PATCH 兼收 bool/字符串 + 页面 `recordSource` 回填 + `TestNotificationSettingsStringForm` |
| F-002 | **fixed** | disable/unlock 仅真实转移发（handler 预读）+ 用例 |
| F-003 | **fixed** | NotifyAccountEvent 失败 `slog.Error` |
| F-004 | **fixed** | 跨用户 404 / 开关抑制 / 锁内不重复 4 用例 |
| F-005 | **fixed** | 铃铛 aria-label 带未读数 + `shell.notifications.unavailable` |
| F-006 | 留痕 | D-002 `8：钩子不依赖模块启用（迁移全局）；行可能落库但 UI 面不可用 |

## 验证

`go test ./... -count=1` + `npm test` 复跑全绿后关门（E-004）。
