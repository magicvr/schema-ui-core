---
id: E-004-r4-c3-cutover-routes
doc: execution-entry
goal: GOAL-009-r4-c3-users-roles-migration
source: orchestrator
date: 2026-08-05
status: recorded
---

# E-004 · R4-C3 composition 切换（路由切片）

## 已发生事实

按冻结包 §7 步骤 3 完成 users/roles 路由的 composition 切换（C3.3 路由切片）：

- `handler/health.go` `Register`：移除 `admin.users`/`admin.roles` 的中心
  `registerResource` 分支；保留 core auth/accounts/health/schema。新增
  `MountProviderRoutes`（同一资源工厂）供测试环境与兼容场景使用。
- `handler/testhelpers_test.go`：测试环境改走 `MountProviderRoutes` 挂载
  users/roles 路由（与 composition 一致）。
- `composition/newMux`：构建 `usersmodule.New(a, st)` / `rolesmodule.New(a, st)`
  provider，调用 `kernel.RegisterContributions(plan, providers)` 消费 finalize，
  挂载返回的路由；无永久双路径。
- `modules/{users,roles}/provider_test.go`：兼容比较测试改写为**鉴权行为验证**
  （匿名 401 / 管理员登录后 list 200 / 未知 detail 404），直接验证 C3.3 生产路径，
  并回应 Grok A-003 F-IND-C32-001（鉴权成功路径对比）。
- `kernel.StandardAdminCapabilities()` 导出，users/roles provider 复用（F-IND-C32-003
  fixed）。

## 验证

- `go build ./...`、`go test ./...`（apps/api 全量）、`go vet ./...` 通过。
- 生产 mux 由 provider finalize 挂载 users/roles 路由；中心 Register 不再挂载；
  无重复注册。settings/activity/manifest/schema 仍走中心路径（C4 范围）。

## C3.3 剩余

Schema owner map（`handler/schema.go` `schemaDocumentsForPlan` 的 users/roles gating）
与 Manifest `adminModules`（`manifest/manifest.go`）对 users/roles 的中心投影仍未
清除，需 schema 内容迁入 `modules/{users,roles}/schema` + manifest 消费 provider
fragment；该部分与 settings/activity（C4）的同类特例纠缠，登记为 C3.3 续作。
