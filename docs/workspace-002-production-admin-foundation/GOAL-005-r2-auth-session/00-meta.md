---
id: GOAL-005-r2-auth-session
title: R2 · 真实认证与请求级身份
status: active
created: 2026-08-02
updated: 2026-08-02
parent: GOAL-001-production-admin-foundation
version: 0.4.0
progress: 6/6
---

# GOAL-005 · R2 · 真实认证与请求级身份

## 概述

承接 Root **D-007** 认证方案（短 JWT Access Token + Opaque Refresh Token + SQLite 存储 + 接受 JWT 库）。交付真实登录 / 登出 / 会话恢复 / 过期 / 撤销与**请求级身份中间件**，使请求身份不再由进程内 `StaticDevSession` 注入；静态开发会话仅保留为**显式本地 dev 兜底**（验收 M9）。

本目标不实施 R3（用户—角色—菜单持久化与权限模型，`I-003`）；SQLite 在 R2 仅为满足真实登录的最小用户/凭据 + 刷新令牌存储。

## 成功标准

- [x] 登录与登出：`POST /api/auth/login` 校验 SQLite 种子凭据（失败 → `401`，成功 → 返回短 JWT access + opaque refresh）；`POST /api/auth/logout` 撤销 refresh token，登出后请求失效 — `internal/handler/auth.go` + 端点测试（2026-08-02）
- [x] 会话刷新与过期：`POST /api/auth/refresh` 用有效 refresh token 换新 access；refresh 过期 / 已撤销 → `401`；access 短时效，过期后凭 refresh 续期（M5/M6） — `internal/auth` `Refresh` + 单测（2026-08-02）
- [x] 请求级身份中间件：业务路由经 `Authorization: Bearer` 解析请求身份；无 / 无效 / 过期 access → `401`；无权限 → `403`（复用 D-PERM `Allow`）；**不再**以 `StaticDevSession` 进程注入身份（M8/M10） — `internal/auth` `Middleware`、records write gate（2026-08-02）
- [x] SQLite 存储与依赖：刷新令牌（哈希）与最小种子用户/凭据落 SQLite（先支持 sqlite）；`go.mod` 引入 JWT 库与 SQLite 驱动并记录选型（M12） — `internal/store` + `go.mod`（2026-08-02）
- [x] 前端认证闭环：登录页、登出、会话恢复（refresh 自动续期 access）、`401` → 重登 / `403` 处理；前端令牌存储策略定稿并留痕（M4/M11） — `LoginPage` / `AuthProvider` / `authFetch` / `tokens.ts`（2026-08-02，D-002 存储已定稿）
- [x] 静态开发会话仅 dev 兜底 + 安全配置 env 化（M9/M13）：后端 `AUTH_DEV_SESSION_ENABLED` 默认 false 与 env 注入已完成；CORS 经 **D-005**（同源最小假设，R2 无需 CORS 头）收口 — 2026-08-02

## 派生进度

`progress: 6/6` 由上方 6 条成功标准等权派生（后端 1–4 + 前端 5–6 均完成，2026-08-02）。**尚未关门**：需 close-out 审计 + 用户确认后才置 `done`；**不**放行 Root R2 检查点（R2 勾选需本目标关门审计）。

## 信息需求

| ID | 问题 / 所需信息 | 级别 | 影响门禁 | 最晚阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据或结论 |
|----|-----------------|------|----------|----------|-----------------|------|-------------|------------|
| `I-005-001` | 安全配置参数最小集：JWT secret 注入、access/refresh TTL、env 键、密码哈希方案？ | required | 方案冻结与实施 | 方案冻结前 | 按 D-007 定稿并记录决策 | **verified** | 已关闭（D-004） | D-004：access 15m / refresh 30d / bcrypt cost 10；env 键集 `AUTH_JWT_SECRET`/`AUTH_ACCESS_TTL`/`AUTH_REFRESH_TTL`/`DB_PATH`/`ADMIN_INITIAL_PASSWORD`/`AUTH_DEV_SESSION_ENABLED`；secret 不落仓库，生产缺失 fail-closed |
| `I-005-002` | 前端令牌存储策略：access（内存）与 refresh（存储介质）的 XSS 边界取舍？ | required | 前端实施与验收 | 方案冻结前 | 定稿存储策略并留痕 | **verified** | 已关闭（D-002） | D-002：access 仅内存、refresh 存 localStorage（用户书面接受 XSS 暴露面）；缓解 = 短 access + 服务端撤销 + HTTPS；验收 M11 |
| `I-005-003` | SQLite 表结构与驱动选型：`refresh_tokens`/最小用户表、纯 Go vs CGO 驱动、迁移与种子方式？ | required | 实施路径 | 实施开始前 | 选型并记录决策 | **verified** | 已关闭（D-003） | D-003：`golang-jwt/jwt v5` + `modernc.org/sqlite`（纯 Go 无 CGO）+ `x/crypto/bcrypt`；`users` + `refresh_tokens` 表；启动幂等建表；admin 种子（env 密码）；R3 才规范化角色关系 |
| `I-005-004` | R2 最小 CORS / 同源托管假设是什么？ | non-blocking | R2 验收边界 | R2 验收前 | 记录 R2 所需最小假设；完整部署基线属 Root `I-005` | **verified** | 已关闭（D-005） | D-005：R2 采用同源最小假设（Vite `/api` 代理；Bearer 无 cookie 域耦合），不引入 CORS 头；跨源与 CORS 属 R5 |
| `I-005-005` | 与 R3 依赖边界：认证产物身份对象与 R3 用户—角色模型字段是否兼容？ | non-blocking | R3 平滑扩展 | R2 验收前 | 确认身份对象字段边界，不实施 R3 | **verified** | 已关闭（D-006） | D-006：`account.User {id,name,roles}` 为 `$context.user`/`/me` 契约形状，R3 以规范化持久化支撑同形状、不得破坏契约；`users.roles` JSON 列为 R3 占位 |

## 依赖与边界

| 项 | 说明 |
|----|------|
| 父目标 | `GOAL-001-production-admin-foundation`（Root `D-007` 已冻结机制 / 存储 / 依赖；`I-002` = verified） |
| 前置证据 | `I-002-auth-collection.md`（[../GOAL-001-production-admin-foundation/attachments/I-002-auth-collection.md](../GOAL-001-production-admin-foundation/attachments/I-002-auth-collection.md)，M1–M14 验收矩阵） |
| In | 登录/登出/刷新端点、请求级身份中间件、SQLite 刷新令牌 + 最小种子凭据、JWT 签发/校验、安全配置 env 化、前端登录页与令牌存储定稿、M1–M14 对应自动化测试 |
| Out | R3 用户—角色—菜单持久化与权限模型（`I-003`）；`D-UPLOAD`；覆盖表扩域；Cookie/同源会话方案（未选 A）；完整部署基线（R5 `I-005`） |

## 父目标

- [GOAL-001-production-admin-foundation](../GOAL-001-production-admin-foundation/00-meta.md)
