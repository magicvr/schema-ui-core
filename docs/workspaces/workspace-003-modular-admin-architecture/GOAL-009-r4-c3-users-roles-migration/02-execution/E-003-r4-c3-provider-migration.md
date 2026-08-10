---
id: E-003-r4-c3-provider-migration
doc: execution-entry
goal: GOAL-009-r4-c3-users-roles-migration
source: orchestrator
date: 2026-08-05
status: recorded
---

# E-003 · R4-C3 Users/Roles provider 化（C3.2）

## 已发生事实

按冻结包 §7 步骤 1-2，为 `admin.users`/`admin.roles` 建立 provider 并与中心输出
兼容比较（测试内，未切换生产）：

- `handler`：`registerResource` 重构为复用 `resourceRoutes`（返回
  `[]kernel.RouteContribution`）；导出 `UsersResource`/`RolesResource`/
  `ResourceRoutes`。provider 生成的 HTTP surface 与中心 `registerResource` 挂载的
  路由逐字节一致（同一工厂、同一 middleware/权限门禁）。
- `kernel`：`BuiltinModules` 的 `admin.users`/`admin.roles` descriptor 补充
  Routes + Fragments 声明（与 provider descriptor 全匹配，满足 C2
  `descriptorsMatch`）。
- `modules/users/provider.go`：`Provider`（Descriptor + `Register` 写 HTTP 5 路由/
  Schema page `users`/permissions `users.read`/`users.write`/navigation `menu_users`/
  Manifest fragment `users`；`CompiledPersistence` 返回 nil）。
- `modules/roles/provider.go`：同上（HTTP 5 路由/Schema `roles`/permissions
  `roles.read`/`roles.write`/`roles.assign`/navigation `menu_roles`/fragment `roles`）。
- 测试：`modules/{users,roles}/provider_test.go` — 表面注册断言 + 与中心 mux 的
  请求级兼容比较（匿名 GET/POST/PATCH/DELETE 返回相同 status）。

## 验证

- `go build ./...`、`go test ./...`（apps/api 全量）通过；`go vet ./...` 干净。
- fx 边界：`modules/users`、`modules/roles` 均无 `go.uber.org/fx` import。
- 生产 mux 未切换：composition 仍走中心 `handler.Register`；provider 仅测试内
  生成 surface 并与中心比较（无永久双注册）。

## 边界

C3.3（composition 消费 `RegisterContributions`、移除中心 users/roles 分支、Schema
owner map、Manifest `adminModules`）与 C3.4（行为矩阵 + 双 Profile + operationlog
失败注入 + 复审）待实施。Schema 文档内容仍由 handler fixture 提供，provider 仅注册
page 元数据（PageContribution）。
