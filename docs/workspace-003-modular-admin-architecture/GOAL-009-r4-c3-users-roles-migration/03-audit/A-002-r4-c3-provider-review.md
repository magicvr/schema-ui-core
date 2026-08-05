---
id: A-002-r4-c3-provider-review
doc: audit-entry
goal: GOAL-009-r4-c3-users-roles-migration
source: self
date: 2026-08-05
scope: R4-C3 Users/Roles provider 化实施、兼容验证与 C3.2 就绪
verdict: conditional
---

# A-002 · R4-C3 Provider 化 self 审计

## Finding closure

| finding | closure | 证据路径 |
|---------|---------|----------|
| F-C3-001（C3-I001/I002 待证据） | `fixed`（C3.1） | E-002 扫描 + 行为矩阵 |

## 已核实成果

- `handler.resourceRoutes`/`ResourceRoutes`/`UsersResource`/`RolesResource`：provider
  生成的 HTTP surface 与中心 `registerResource` 逐字节一致（同一工厂/middleware/权限）。
- `modules/users`、`modules/roles` Provider：Descriptor 与 BuiltinModules 全匹配
  （含 Routes/Fragments 声明）；Register 写 HTTP 5 路由/Schema/权限/导航/Manifest
  fragment；CompiledPersistence 返回 nil。
- 测试：表面注册断言 + 与中心 mux 请求级兼容比较（匿名 5 方法返回相同 status）。
- 全量 `go test ./...`（apps/api）+ `go vet` 通过；fx 边界满足（modules 无 Fx import）。
- 生产 mux 未切换：composition 仍走中心 `handler.Register`；无永久双注册。

## Open required

无。C3.3（composition 切换 + 中心特例清除）与 C3.4（行为矩阵 + 双 Profile +
operationlog 失败注入 + 复审）为后续检查点，不阻断 C3.2 关门。

## Gate

C3 保持 `active 2/4`；C3.1/C3.2 勾选。Grok independent 复审未完成前不关门。
本意见不修改 status/progress。
