---
id: E-003
goal: GOAL-006-r2-f04-notification-center
title: S2 实现 + S3 验证（通知模块全量落地 + 回归 + 冒烟）
date: 2026-08-14
status: recorded
parent: GOAL-006-r2-f04-notification-center
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# E-003 · S2 实现 + S3 验证

## S2 事实（checkpoint `a065288`，34 files / +1489）

| 项 | 落地 |
|----|------|
| 迁移 0016/0017 | `notifications` 表（FK→users、事件 CHECK、索引）+ `users.notifications_enabled` |
| 端点 | list（分页/unreadOnly）/ read（owner 幂等）/ read-all / unread-count / settings |
| 事件钩子 | `auth.OnLockOpened`（锁开窗）+ disable/unlock（users_state）+ 改密（自助+管理员）→ `NotifyAccountEvent`（best-effort，主开关门控） |
| 裁剪 | 每用户 500 条（两段式：计数 → 删最旧已读，未读不丢） |
| 页面 | `notifications`（schema：设置表单 + table + 行 read + read-all） |
| 铃铛 | App.tsx 顶栏 `NotificationBell`（徽标/下拉 5 条/查看全部；失败静默） |
| 装配 | mvp/admin/demo += `admin.notifications`（内容扩展）；smoke 页面集；i18n 键 |

## S3 验证事实（2026-08-14）

| 项 | 结果 |
|----|------|
| `go test ./... -count=1` | ✅ 全绿（事件钩子 4 类、已读/全部已读/未读计数、越权 404、设置开关抑制、裁剪保未读） |
| `npm test` | ✅ 896/896（铃铛 3 用例：徽标/下拉/fail-open） |
| 错误契约 | +2 码（INVALID_SETTINGS_BODY / NOTIFICATION_NOT_FOUND）入 frozen 集 + catalog |
| 迁移账本 | 17 条全绿（0016/0017 追加） |

## 门禁结论

S2/S3 完成。进入 S4。
