---
id: E-002
goal: GOAL-025-w16-rectification-batch-a
title: 批 A S2 实施（F01/F07/F08）
status: completed
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# E-002 · 批 A S2 实施

## 1. 执行事实

- **日期**：2026-08-17
- **动作**：
  1. **F01 首次强制改密**：
     - `users` 表新增 `must_change_password`（migration 0038）。
     - `authsession.User` / `account.User` 增加字段，登录与 `/me` 用户快照携带 `mustChangePassword`。
     - 认证中间件对 `must_change_password=1` 用户仅放行 `POST /api/account/password`、`GET /api/account/profile`，其余业务 API 返回 `403 MUST_CHANGE_PASSWORD`。
     - 新建/导入/管理员重置密码默认置 `true`；自助改密成功后清标记；强制改密场景返回新 token pair。
  2. **F07 一键下线其他设备**：
     - 新增 `POST /api/account/sessions/revoke-others`。
     - 后端 `BumpTokenVersionAndRevokeAll` + `IssueTokensFor` 为当前设备重签令牌；前端个人中心会话页新增 `account-session-toolbar` 自定义组件。
  3. **F08 验证码刷新 + MFA 备份**：
     - 登录页验证码区域新增“换一题”按钮。
     - MFA 绑定弹窗新增“复制密钥”与“下载恢复码 (.txt)”按钮。
- **测试**：
  - Go：新增 `TestForcedPasswordChangeGateAndReissue`、`TestRevokeOthersReissuesTokensAndRevokesOtherSessions`；修正 error contract 与 notifications 测试以适配新语义；`go build ./...` 通过，相关 handler 定向测试通过。
  - Web：新增 `ForcePasswordChange`、`AccountSessionToolbar` 测试，扩展 `LoginPage`/`MfaManager` 测试；`npx tsc -b` 通过，相关 vitest 22/22 通过。

## 2. 证据

| 主张 | 路径 / 证据 |
|------|-------------|
| F01 migration | `apps/api/internal/modules/authsession/migration/migration.go`（0038） |
| F01 gate | `apps/api/internal/auth/auth.go` `isMustChangePasswordAllowed` |
| F07 endpoint | `apps/api/internal/handler/account_self.go` `revokeOthers` |
| F08 UI | `apps/web/src/app/LoginPage.tsx`、`apps/web/src/components/mfa-manager.tsx` |
| Go tests | `apps/api/internal/handler/w16_batch_a_test.go` |
| Web tests | `apps/web/src/components/force-password-change.test.tsx`、`account-session-toolbar.test.tsx` |
