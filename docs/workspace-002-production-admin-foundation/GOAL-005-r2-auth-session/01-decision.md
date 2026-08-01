---
title: 决策 · R2 · 真实认证与请求级身份
status: active
created: 2026-08-02
updated: 2026-08-02
parent: null
version: 0.2.0
---

# 决策 · GOAL-005

## D-001 · 承接 Root D-007，本目标定稿 R2 实施细节

- **日期**：2026-08-02
- **状态**：accepted
- **决定**：
  1. 本目标实现 Root `D-007` 冻结的认证方案（短 JWT Access + Opaque Refresh + SQLite + 接受 JWT 库），交付真实登录/登出/刷新/过期/撤销与请求级身份中间件。
  2. 实施细节（`I-005-001` 安全配置参数、`I-005-002` 前端令牌存储策略、`I-005-003` SQLite 表结构与驱动选型）在本目标**方案阶段定稿并记录决策**，不在立项时静默假设。
  3. 静态开发会话 `StaticDevSession` 仅作为显式本地 dev 兜底，生产默认不启用。
  4. 不实施 R3（用户—角色—菜单持久化与权限模型，Root `I-003`）。
- **理由**：用户确认创建 `GOAL-005-r2-auth-session`；方案机制已由 D-007 裁决，本目标负责把安全配置与存储细节以证据定稿，避免把假设写成既定事实。
- **未选方案**：
  - **拆分为前后端两个子目标**：范围过碎，登录/会话闭环需要端到端同验，一个目标更利于 M1–M14 一体验收。
  - **与 R3 合并为一个目标**：越界实施 `I-003`，且 R3 的持久化身份模型是独立信息门禁。
  - **本目标内重选 Cookie / 纯 JWT 方案**：D-007 已冻结 C+B 混合，不重开。

## D-002 · 前端令牌存储策略：access 内存 + refresh localStorage（I-005-002 裁决）

- **日期**：2026-08-02
- **状态**：accepted
- **决定**：
  1. **access token 仅存内存**（React state，不落任何持久介质）；**refresh token 存 `localStorage`**。
  2. 页面启动先尝试用 localStorage 中的 refresh 换新 access；401 时静默尝试 refresh 一次，失败转登录页；登出清除两者。
  3. **接受 localStorage 的 XSS 暴露面**（用户书面裁决）。
- **理由**：用户在「UX 持久性 vs XSS 暴露」取舍中选定跨会话持久方案。
- **缓解边界**（已书面确认的残余风险缓解，非新增决策）：access 15m 短时效限制窃取窗口；refresh 服务端可撤销（登出/旋转即作废）；凭据仅经 HTTPS；后续可加 CSP 加固（不在本目标强制）。
- **关联信息项**：`I-005-002` → `verified`（证据 = 本决策）。

## D-003 · Go 依赖选型：golang-jwt/jwt v5 + modernc.org/sqlite + bcrypt（I-005-003 裁决）

- **日期**：2026-08-02
- **状态**：accepted
- **决定**：
  1. JWT 库 = `github.com/golang-jwt/jwt/v5`；SQLite 驱动 = `modernc.org/sqlite`（**纯 Go、无 CGO**）；密码哈希 = `golang.org/x/crypto/bcrypt`（cost 10）。
  2. **SQLite 表结构（R2 最小形态）**：
     - `users`（`id` TEXT PK · `username` TEXT UNIQUE · `name` TEXT · `roles` TEXT JSON 数组 · `password_hash` TEXT · `created_at`/`updated_at`）
     - `refresh_tokens`（`id` TEXT PK · `user_id` FK → users · `token_hash` TEXT UNIQUE · `expires_at` INTEGER · `revoked_at` INTEGER NULL · `created_at`）
  3. **迁移**：启动时 `CREATE TABLE IF NOT EXISTS`（幂等，R2 最小）；正式迁移工具与用户—角色—菜单规范化属 R3（`I-003`）。
  4. **种子**：启动时若 `users` 为空，创建 admin（username=`admin`，密码 = `ADMIN_INITIAL_PASSWORD` env；dev 缺省 `admin`），`roles=["admin","editor"]`（对齐现有 `StaticDevSession` 覆盖的 D-PERM 场景）。
- **理由**：用户裁决；modernc 免 CGO，Windows 与本仓 CI 跨编译平滑；jwt v5 为事实标准；bcrypt cost 10 为标准安全。
- **关联信息项**：`I-005-003` → `verified`（证据 = 本决策）。

## D-004 · 安全配置参数（I-005-001 裁决）

- **日期**：2026-08-02
- **状态**：accepted
- **决定**：
  - **TTL**：access `15m`（`AUTH_ACCESS_TTL`）；refresh `30d`（`AUTH_REFRESH_TTL`，`720h`）。
  - **哈希**：bcrypt cost 10。
  - **env 键集**：`AUTH_JWT_SECRET`（签发密钥；生产必填，缺失 → 启动失败 fail-closed；dev 缺省本地开发密钥并打警告）、`AUTH_ACCESS_TTL`、`AUTH_REFRESH_TTL`、`DB_PATH`（默认 `./data/schema-ui.db`）、`ADMIN_INITIAL_PASSWORD`、`AUTH_DEV_SESSION_ENABLED`（默认 `false`；仅 dev 显式开启 `StaticDevSession` 兜底）。
  - `AUTH_JWT_SECRET` 与 `ADMIN_INITIAL_PASSWORD` **不落仓库**；`.env.example` 仅放占位说明。
- **理由**：用户裁决 profile；短 access + 可撤销 refresh 匹配 D-007；静态会话默认关闭满足验收 M9。
- **关联信息项**：`I-005-001` → `verified`（证据 = 本决策）。
