---
id: E-003
goal: GOAL-003-r2-f01-dashboard
title: S2 实现 + S3 验证（dashboard 模块 + home 装配 + 回归 + 本地冒烟）
date: 2026-08-14
status: recorded
parent: GOAL-003-r2-f01-dashboard
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# E-003 · S2 实现 + S3 验证

## S2 事实（checkpoint `a9642f4`，13 files / +251）

| 项 | 落地 |
|----|------|
| 模块 | `modules/dashboard/`（provider + fragment + schema：text + grid × statCard×2） |
| home 装配 | `adminFunctionalOrder` 头部插 `admin.dashboard` → mvp/admin home=dashboard（demo 仍 overview） |
| Profile | mvp/admin/demo 默认集 += `admin.dashboard`（内容扩展声明 D-002 `3） |
| 导航 | sidebar `menu_dashboard`（Order 0；PolicyAdminEditorViewer；图标 dashboard） |
| 指标 | statCard × 2（`/api/users` / `/api/roles` envelope total）——无新端点 |

**观察（留痕）**：sidebar 视觉顺序按模块 id 字母序聚合（account/activity/dashboard/roles/settings/users）——既有聚合约定，不改排序语义；功能 home 顺序（dashboard 首位）由 `adminFunctionalOrder` 保证。

## S3 验证事实（2026-08-14）

| 项 | 结果 |
|----|------|
| `go test ./... -count=1` | ✅ 全绿（home 推导 dashboard、demo overview 不变、permissions/nav 计数、kernel 列表） |
| `npm test` | ✅ 892/892（schema 结构校验含 dashboard 页；i18n 键全） |
| 本地冒烟（admin profile） | ✅ manifest pages 含 dashboard、`home=dashboard`、schema/dashboard 200、users/roles 计数可读 |

## 门禁结论

S2/S3 完成。进入 S4。
