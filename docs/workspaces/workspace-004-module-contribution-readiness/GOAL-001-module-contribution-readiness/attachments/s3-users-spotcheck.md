---
title: S3 抽检 · admin.users vs playbook
status: active
created: 2026-08-06
updated: 2026-08-06
parent: null
version: 0.2.0
---

# S3 抽检记录 · admin.users

权威 playbook：`docs/architecture/module-contribution-playbook.md`  
被检模块：`apps/api/internal/modules/users/provider.go`（及 schema/、manifest/）

## MUST

| 项 | 证据摘录 | 结论 |
|----|----------|------|
| M1 | `const ModuleID = "admin.users"`；Descriptor Version `2.0.0`；DependsOn `core.auth-session`, `core.navigation-capability`, `core.schema-render`, `core.operationlog` | pass |
| M2 HTTP | Register → `reg.HTTP` 多条 `/api/users*` | pass |
| M2 Schema | `reg.Schema` PageID `users`，Document from `usersschema` | pass |
| M2 Authorization | `users.read` / `users.write` PermissionContribution | pass |
| M2 Navigation | `menu_users` NavigationContribution | pass |
| M2 Manifest | `reg.Manifest` FragmentID `users` | pass |
| M2 Persistence | CompiledPersistence nil；依赖 auth-session 拥有账户迁移 | pass（归属清晰） |
| M3 | composition.go `plan.HasModule("admin.users")` → usersmodule.New | pass |
| M4 | profile.go ProfileMVP/Admin 含 `admin.users` | pass |
| M5 | compiled.PersistenceProviders 含 auth/history/operationlog/settings 迁移；全局收集 | pass |
| M6 | provider_test.go 存在并行使 Provider | pass |

## DO NOT

| 项 | 观察 | 结论 |
|----|------|------|
| D1 | 无强制修改 apps/web Renderer 中央注册以暴露 users | pass |
| D2 | fragment 由模块贡献；非 public 静态生产 Manifest | pass |
| D3 | 使用 auth.Authenticator / authsession.Repository，无平行 DB | pass |
| D4 | 标准 Admin：Register 覆盖 HTTP/Schema/Authorization/Navigation/Manifest；Persistence 语义经依赖 + `CompiledPersistence` 空返回声明（非「按需能力」永久缺省核心六项）。未以 Configuration/Lifecycle/Observability 替代 §2.1 | pass |
| D5 | 接入路径为组合根静态候选 + `plan.HasModule("admin.users")` 装配；模块/Provider 侧无 `.so` / 远程下载 / 运行中启停插件路径 | pass |

### D4 / D5 补抽说明（2026-08-06 · 响应 A-002 F-001）

- **D4**：对照 playbook §2 D4 与 architecture §2.1/§2.2；`admin.users` 作为标准 Admin 功能模块满配核心贡献面，Persistence 空返回有 auth-session 归属说明，不构成「按需豁免核心六项」。  
- **D5**：对照 playbook §2 D5 与 architecture §3；启用仅通过 Profile/`modules.enabled` 在已编译候选集中选择；本抽检未发现 users 模块引入热插拔加载器。

## 结论

`admin.users` 可作为 playbook 标准 Admin 功能模块正例；playbook 路径与现网一致，DO NOT D1–D5 均已逐行勾选，未发现需改 runtime 的阻断缺口。
