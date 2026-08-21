---
id: E-008-r6-c62-authsession-repository
doc: execution-entry
goal: GOAL-013-r6-old-path-removal
source: orchestrator
date: 2026-08-06
status: recorded
---

# E-008 · R6 C6.2 auth-session/RBAC repository 迁出（切片 5）

## 已发生事实

- `apps/api/internal/modules/authsession/` 新增 module-owned Repository，拥有 account、
  refresh token、用户管理与 roles/RBAC 的类型、错误和 SQL。Repository 只消费结构化
  `TxRunner.WithTx(context.Context, func(*sql.Tx) error)`，不导入具体 `store` 包，也不
  取得底层 `*sql.DB`。
- composition 由 Fx 构造并共享一个 auth-session Repository；Authenticator、
  admin.users Provider 与 admin.roles Provider 均消费该实例。生产 `internal/auth`、
  `modules/users`、`modules/roles` 已无领域化 `store.User` / `store.Role` / Store CRUD
  依赖。
- users/roles handler 改为窄 Repository 接口，并将 module-owned sentinel 显式映射到
  既有 HTTP 错误契约。认证 login/refresh/logout、request identity、角色委派与
  last-admin/self-protection 行为保持由真实 Provider/composition 路径覆盖。
- operation-log writer 暂由独立 `OperationRecorder` 接口注入，但 payload/实现仍位于
  `store`；旧 `store.go`/`users.go`/`roles.go` 领域实现也仍存在。因此本切片只是生产
  接线迁移，不构成 A-010 F-001 的完整 `fixed` 或 C6.2 放行。

## 验证

- `apps/api: go test ./...`：通过。
- `apps/api: go vet ./...`：通过。
- `git diff --check`：通过。
- 生产依赖扫描（排除 `*_test.go`）：`internal/auth`、`modules/users`、`modules/roles`
  对 `internal/store`、`*store.Store` 及 store-owned account/RBAC 类型与 sentinel 为零
  命中。
- 首轮定向测试曾发现 module `ErrNotFound` 未翻译导致资源 404 误报 500；补齐显式
  映射后，handler/users/roles/composition 定向测试与全量测试均通过。

## Git checkpoint

- `f4fb882` · `feat(api): move auth and RBAC repositories into module`
- owned paths：auth-session repository、Authenticator、users/roles handlers/providers、
  composition 接线及相应测试；未暂存既有 store 换行状态噪音。

## 仍开放

- A-010 F-001 / F-C62-004：迁出 admin.settings 与 core.operationlog repository，删除
  `store` 中全部领域 SQL/types/errors/test ownership，使其只保留 DB 生命周期、事务、
  migration runner/ledger/snapshot 与 readiness；完成后才可请求 C6.2 independent audit。
