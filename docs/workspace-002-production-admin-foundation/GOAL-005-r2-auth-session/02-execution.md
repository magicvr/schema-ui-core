---
title: 执行记录 · R2 · 真实认证与请求级身份
status: active
created: 2026-08-02
updated: 2026-08-02
parent: null
version: 0.4.0
---

# 执行记录 · GOAL-005

## 2026-08-02 · 立项

- 用户确认创建 `GOAL-005-r2-auth-session`；按 Root `D-007`（短 JWT Access + Opaque Refresh + SQLite + 接受 JWT 库）立项。
- 建齐五件套与 `attachments/`；`goal-tree.md` 同步（树 + 表，parent = `GOAL-001-production-admin-foundation`）。
- 登记 `I-005-001`（安全配置参数）、`I-005-002`（前端令牌存储）、`I-005-003`（SQLite 表结构与驱动）为 required；`I-005-004`（CORS 最小假设）、`I-005-005`（R3 边界）为 non-blocking。
- 本回合**无产品代码变更**；Root 纲领检查点仍为 `1/5`（R2 未实施完成）。

> 本节只记录立项与文档落盘事实，不代表任何认证功能已经实施或通过验收。

## 2026-08-02 · 方案参数定稿（D-002 / D-003 / D-004）

- 用户裁决三项开放前提：前端令牌存储 = **access 内存 + refresh localStorage**；Go 依赖 = **golang-jwt/jwt v5 + modernc.org/sqlite（纯 Go 无 CGO）**；TTL/哈希 = **access 15m / refresh 30d / bcrypt cost 10**。
- 记录决策 D-002（前端存储，含 localStorage XSS 残余的书面缓解边界）、D-003（依赖 + SQLite 表结构 + 迁移 + admin 种子）、D-004（安全配置参数与 env 键集）。
- `I-005-001` / `I-005-002` / `I-005-003` → **verified**（已关闭）；GOAL-005 方案冻结完成，可进入后端认证端点实施。
- **未做**：尚无产品代码变更；前端登录页与令牌存储定稿后的实现、后端实施均未开始。
- **计划（非事实）**：进入后端认证端点实施（login/refresh/logout、请求级身份中间件、SQLite 存储、`/api/accounts/me` 与 records gate 改走请求身份）。

## 2026-08-02 · 后端认证端点实施（成功标准 1–4）

- **依赖（D-003）**：`go.mod` 引入 `github.com/golang-jwt/jwt/v5 v5.3.1`、`modernc.org/sqlite v1.55.0`（纯 Go 无 CGO）、`golang.org/x/crypto v0.54.0`。
- **配置（D-004）**：`internal/config` 新增 `AUTH_JWT_SECRET` / `AUTH_ACCESS_TTL` / `AUTH_REFRESH_TTL` / `DB_PATH` / `ADMIN_INITIAL_PASSWORD` / `AUTH_DEV_SESSION_ENABLED`；`.env.example` 同步；生产缺 secret / 种子密码 fail-closed，dev 使用文档化默认。
- **存储（D-003）**：新建 `internal/store`（modernc sqlite）：`users` + `refresh_tokens` 表（幂等建表），admin 种子（roles `admin`/`editor`），refresh token 以 SHA-256 哈希存储，支持撤销。
- **认证核心**：新建 `internal/auth`：JWT 签发/校验（HMAC-SHA256，`sub`=user id）、opaque refresh 生成/哈希、bcrypt 校验、请求身份上下文；`Middleware` 解析 `Authorization: Bearer`，失败 fail-closed 401；`AUTH_DEV_SESSION_ENABLED` 显式启用时以 `StaticDevSession` 兜底（M9）。
- **端点与路由改造**：`POST /api/auth/login|refresh|logout`（login 失败 401、缺字段 400；refresh 轮换并撤销旧 token；logout 幂等）；`GET /api/accounts/me` 与 records 写路由改走请求身份中间件（无/无效 token 401，非 admin 403）；records GET 只读保持公开（本阶段边界）。
- **测试**：新增 `internal/store`（种子幂等、token 生命周期）、`internal/auth`（login 成功/失败/枚举防护、refresh 轮换/撤销/过期、logout、JWT 过期/密钥错）、`internal/handler`（端点与 401/403、`/me`、dev-session 兜底）测试；`account_test` / `records_test` / `health_test` 改写为真实认证接线。
- **验证**：`go build` / `go vet` / `go test ./...` 全绿；`gofmt` 干净；web `npm test` 425/425、`npm run build` 通过；**真实服务器 smoke**（dev 种子 admin/admin）验证 login → `/me`(Bearer) → 无 token 401 → refresh 轮换 + 旧 token 撤销 → logout 撤销 → 写路由 admin 200 / 无 token 401。
- **E2E 边界**：`apps/web/e2e/shell.spec.ts` 已改写为真实认证链（login → /me → records），但**本机未能执行**：Windows 本机 `127.0.0.1:5173` 被禁止绑定（EACCES；8080/9999 可绑定，非代码问题）；CI `browser-e2e`（Linux）将覆盖。
- **未做**：前端登录页 / 会话恢复 / 令牌存储定稿后的实现（成功标准 5）；CORS 最小假设（`I-005-004` non-blocking，仍 open）；R3 身份模型（不实施）。
- **计划（非事实）**：下一拍实现前端认证闭环（登录页、refresh 续期、401 转登、localStorage 存储）并收口 CORS / R3 边界。

## 2026-08-02 · 前端认证闭环实施（成功标准 5–6，I-005-004/005 收口）

- **令牌存储（D-002）**：新增 `apps/web/src/account/tokens.ts` — access 仅内存，refresh 存 localStorage（`schema-ui.refreshToken`），`clearTokens`/`hasSession`。
- **认证客户端**：新增 `apps/web/src/account/auth-client.ts` — `login`/`logout`/`restoreSession`/`fetchMe` 与 **`authFetch`**（自动挂 `Authorization: Bearer`；401 时静默 refresh 一次并重试；refresh 失败 → 清会话并通知 auth-lost，UI 回登录页）。`AuthError` 稳定错误码（`INVALID_CREDENTIALS` 等）。
- **状态管理**：新增 `apps/web/src/account/AuthContext.tsx` — `AuthProvider`（启动 `restoreSession` 恢复、`login`/`logout` 切换状态、监听 auth-lost）+ `useAuth`。
- **登录页**：新增 `apps/web/src/app/LoginPage.tsx` — 未认证时整屏登录表单（用户名/密码、提交中禁用、错误显示）。
- **接线**：`main.tsx` 改为 `AuthProvider` + `AuthGate`（loading → 启动屏；unauthenticated → LoginPage；authenticated → App，`navigationContext.user` 取真实身份、`recordsFetcher={authFetch}`）；`App.tsx` 新增 `onLogout`/`currentUser` 与 Sign out 按钮。
- **收口**：记录 **D-005**（`I-005-004` → verified：R2 同源最小假设，不引入 CORS 头，跨源属 R5）、**D-006**（`I-005-005` → verified：`account.User` 形状为 `$context.user`/`/me` 契约，R3 同形状支撑）。
- **测试**：新增 `tokens.test.ts`（access 内存 / refresh localStorage / 清除）、`auth-client.test.ts`（login 成败、Bearer 注入、401 刷新重试、refresh 失败通知 auth-lost、restore、logout）、`LoginPage.test.tsx`（渲染/提交/错误/禁用）。
- **验证**：web `npm test` **441/441**（新增 16 项）、`npm run build`（tsc + vite）通过；api `go test` 不受影响。
- **E2E**：`apps/web/e2e/shell.spec.ts` 改写为登录门禁流（未认证 → 登录页 → admin/admin 登录 → shell → 独立 API 认证链断言）；**本机仍受 `127.0.0.1:5173` 绑定被禁限制未能执行**，CI `browser-e2e`（Linux）覆盖。
- **未做**：GOAL-005 尚未 close-out（需审计 + 用户确认置 `done`）；Root R2 检查点未勾选；R3 不实施。

## 计划（非事实）

- 下一拍：定稿 `I-005-001` / `I-005-002` / `I-005-003`（安全配置、前端存储、SQLite 表与驱动选型）并记录决策，再进入后端认证端点与中间件实施。
