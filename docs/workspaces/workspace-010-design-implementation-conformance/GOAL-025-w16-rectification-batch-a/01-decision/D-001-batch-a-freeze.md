---
id: D-001
goal: GOAL-025-w16-rectification-batch-a
title: W16 批 A 方案冻结（F01 / F07 / F08）
status: approved
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# D-001 · W16 批 A 方案冻结

## 1. 触发

GOAL-024 D-002/D-003 已冻结技术方案与分批规划；本决策将批 A 三项目（F01/F07/F08）的端点、存储、前端交互细化为可实施设计，并关闭本子目标 I-001/I-002。

## 2. W16-F01 · 首次登录强制修改初始密码

### 2.1 存储

- `users` 表新增 `must_change_password INTEGER NOT NULL DEFAULT 0`（迁移贡献，authsession module）。
- 种子 `admin`、管理员新建用户、导入用户、管理员重置密码时，`must_change_password = 1`。
- 用户自改密成功后将 `must_change_password = 0`，并沿用现有 `token_version + 1` 与 refresh-token 吊销语义。

### 2.2 API

- `POST /api/auth/login` 成功响应增加 `mustChangePassword: boolean`（`account.User` 或登录响应 map 增加字段）。
- `POST /api/account/password` 保持 `{ currentPassword, newPassword }`；当 `must_change_password=1` 时用初始密码作为 `currentPassword`，成功后清标记。
- 强制改密门禁：在认证中间件或账号端点入口增加白名单：
  - 放行：`POST /api/auth/login`、`POST /api/auth/refresh`、`POST /api/auth/logout`、`POST /api/account/password`、`GET /api/account/profile`、`GET /api/accounts/me`、验证码/MFA 相关必要端点。
  - 其余业务 API 在 `must_change_password=1` 时返回 `403 MUST_CHANGE_PASSWORD`。

### 2.3 前端

- `AuthContext` 解析登录响应 `mustChangePassword`；为 `true` 时跳转强制改密页，阻止进入业务页面。
- 强制改密页复用现有改密表单与错误反馈；成功后刷新会话（更新 access/refresh token 与用户状态）进入系统。

## 3. W16-F07 · 个人中心一键下线其他所有设备

### 3.1 API

- 新增 `POST /api/account/sessions/revoke-others`（`admin.account` 模块贡献路由）。
- 行为：
  1. 基于当前认证身份调用 `BumpTokenVersionAndRevokeAll(userID, now)`：吊销该用户全部 refresh token 并 bump `token_version`（使其他设备已签发的 access token 立即失效）。
  2. 调用 `auth.IssueTokensFor(userID, now)` 为当前设备签发新的 access/refresh 对并返回 `{ accessToken, refreshToken, user }`。
- 响应：`200 OK`，携带新令牌；当前前端用返回令牌替换本地存储，保持当前设备在线。

### 3.2 前端

- 个人中心会话卡片头部增加“下线其他设备”按钮。
- 点击后确认对话框；成功后替换 token、刷新会话列表并 Toast。

## 4. W16-F08 · 登录算术验证码主动刷新 + MFA 密钥复制与恢复码 txt 下载

### 4.1 登录验证码换一题

- 前端登录页验证码区域增加“换一题”按钮，重新请求现有验证码接口并更新 `captchaId`/图片/答案 session；无后端新契约。

### 4.2 MFA 密钥复制与恢复码下载

- MFA 绑定弹窗增加“复制密钥”按钮：`navigator.clipboard.writeText(secret)` + Toast。
- 增加“下载恢复码 (.txt)”按钮：将恢复码数组生成文本 Blob 并触发下载；不新增 API。

## 5. 信息项关闭

| ID | 级别 | 结论 |
|----|------|------|
| I-001 | required | 强制改密白名单按 §2.2 冻结：登录/刷新/登出/改密/资料 + 验证码/MFA 必要端点，其余 403。 |
| I-002 | non-blocking | F07 采用「bump token_version + 吊销全部旧 refresh + 为当前设备重签新令牌」方案，当前设备不中断。 |

## 6. 未选方案

- F07 不采用「仅吊销 refresh token、不 bump token_version」：其他设备已签发的 access token 会残留至过期，不符合“一键下线”即时语义。
- F01 不采用纯前端拦截：缺少后端强制门禁，无法满足安全合规。
