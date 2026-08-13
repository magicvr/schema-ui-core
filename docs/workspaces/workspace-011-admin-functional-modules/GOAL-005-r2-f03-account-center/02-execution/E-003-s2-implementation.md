---
id: E-003
goal: GOAL-005-r2-f03-account-center
title: S2 · 实现（admin.account 模块全量落地）
date: 2026-08-14
status: recorded
parent: GOAL-005-r2-f03-account-center
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# E-003 · S2 · 实现

## 事实（checkpoint `5a50524`，40 files / +1590）

| 项 | 落地 | 证据 |
|----|------|------|
| 迁移 0013 | `users.enabled INTEGER NOT NULL DEFAULT 1`（owner admin.account） | apps/api/internal/modules/account/migration |
| 迁移 0014 | operation_log 事件 CHECK 扩展（5 个新事件，表重建） | apps/api/internal/modules/operationlog/migration |
| 持久化 | `User.Enabled` + 全列清单更新 + `ListRefreshTokensForUser` / `RevokeRefreshTokenIfOwned` / `SetUserEnabled`（self/last-admin 守卫）/ `UnlockUser` | authsession/account_operations.go |
| auth | Login 403 `ACCOUNT_DISABLED`；Refresh/Middleware 对停用账号 401 fail-closed | auth/auth.go、handler/auth.go |
| 端点 | profile / password / sessions / revoke / enable / disable / unlock（8 路由） | handler/account_self.go、users_state.go |
| 模块 | `admin.account` provider + fragment（navigation.user）+ schema（account.json） | modules/account/* |
| users 页 | enabled/locked 列 + enable/disable/unlock 行操作（权限键表达式 + disabledWhen 行状态） | modules/users/schema/users.json |
| 装配 | mvp+admin+demo 默认集 + adminFunctionalOrder 尾部 + BuiltinModules + composition | kernel/profile.go、composition.go |
| i18n | en-US/zh-CN 键（account 页 + users 动作），zh-CN 补齐 | i18n/messages |
| smoke | SM-007 页面集 mvp/admin/demo 增加 account | scripts/smoke.sh |
| 错误契约 | `INVALID_PASSWORD` / `INVALID_PASSWORD_BODY` / `SESSION_NOT_FOUND` / `ACCOUNT_DISABLED` 入 frozen 集 + catalog | error_contract_test.go、errorcatalog |

## 偏差

无。实现与 D-002 冻结一致（路由/语义/守卫/权限键逐项对照）。
