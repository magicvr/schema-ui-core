---
id: E-008
doc: execution
status: recorded
parent: GOAL-014-w13-settings-tabs-and-topbar
created: 2026-08-16
updated: 2026-08-16
version: 0.1.0
---

# E-008 · 用户缺陷反馈修复：顶栏头像即时刷新

## 问题

2026-08-16 用户反馈：「个人中心页确实上传了用户头像，但顶部工具栏显示用户头像的地方并没有正确显示头像」。根因：资料保存（PATCH /api/account/profile）只刷新了页面数据（reloadList），**会话用户快照（/me）不会更新**——顶栏用户菜单的 `currentUser` 来自登录/恢复时的会话，头像/显示名在保存后仍是旧值，必须整页刷新或重新登录才会出现新头像。

## 修复

- **Go**（`handler/account_self.go`）：`PATCH /api/account/profile` 成功响应增加 `X-Schema-UI-Config-Changed: account.profile` 头——复用设置页品牌（settings.branding）既有的「响应头 → 宿主事件」通道，泛型 Renderer 无需感知产品端点。
- **Web**（`app/config-events.ts`）：新增 `ACCOUNT_PROFILE_NAMESPACE = "account.profile"`。
- **Web**（`account/AuthContext.tsx`）：新增 `refreshSession()`（best-effort 重解析 /me 并 setSession，失败保持当前会话）；AuthProvider 挂载期订阅 `account.profile` 事件 → 自动刷新会话。保存资料后顶栏头像/显示名**立即更新，无需刷新页面**（顺带修复改名后顶栏显示名不更新的同类问题）。

## 测试

- Go：`TestAccountProfileUpdateName` 断言 PATCH 响应携带 `X-Schema-UI-Config-Changed: account.profile`。
- Web：`account/auth-context.test.tsx`（新，3 例）——account.profile 事件触发 /me 刷新并更新会话用户；其他命名空间（settings.branding）不触发；刷新失败保持当前会话不抛错。
- e2e：`shell.spec.ts` 增加**纯 UI 驱动**的头像替换流（用户菜单 → 个人中心 → `#field-avatarUrl` 上传 → 保存）——断言顶栏头像 img 在**不刷新页面**的情况下切换为新头像 URL（即用户报告场景的端到端复现与验证）。
