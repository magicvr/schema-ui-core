---
id: A-003-grok-r4-c3-provider-audit
doc: audit-entry
goal: GOAL-009-r4-c3-users-roles-migration
source: independent
auditor: Grok Build / grok-4.5
date: 2026-08-05
scope: R4-C3.2 Users/Roles provider 化 + 中心输出兼容验证（冻结包 §7 步骤 1-2）
audit_type: execution-facts
verdict: pass
---

# A-003 · Grok R4-C3.2 Users/Roles Provider 化独立交叉审计

## 声明

本意见 `source: independent`，只读审计。未修改任何文件、`status` / `progress` /
goal-tree / 台账。响应与关门由 `/govern` 处理。

## 结论摘要

C3.2 在冻结包 §7 步骤 1-2 与 `00-meta` 成功标准下成立：admin.users / admin.roles
已以 kernel.Provider 形态落地；RegisterContributions 路径下 Descriptor 与
BuiltinModules 全字段匹配；Register 写入已声明的 HTTP/Schema/Auth/Nav/Manifest；
CompiledPersistence 返回 nil；兼容比较仅在测试内；生产 composition 仍走中心
handler.Register，无永久双注册；模块包无 Fx import；复跑 go test ./... 与 go vet
均通过。未误将 C3.3/C3.4 宣称为完成。**verdict: pass，开放 required 0，可以继续
C3.3。**

## 核验矩阵（相对 §7 步骤 1-2）

1. Provider 元数据 + 无发布 contract tests：pass
2. typed Provider 生成 surface + 测试内兼容比较：pass（同工厂 `resourceRoutes`）
3. 无生产永久双注册：pass（composition.go:91 仍 handler.Register）
4. Descriptor ↔ BuiltinModules 全规范匹配（descriptorsMatch）：pass
5. Register 写 HTTP×5/Schema/权限/Nav/Manifest 且已声明：pass
6. 路由与中心 registerResource 一致（工厂/中间件/权限门禁）：pass
7. CompiledPersistence → nil：pass
8. 模块包无 Fx import：pass
9. go test ./... / go vet ./...：pass
10. 未超范围宣称 C3.3/C3.4：pass

## Findings（recommended，不阻断 C3.2）

### F-IND-C32-001 · 兼容比较深度偏窄（匿名 status-only）
- 影响 C3.4 行为矩阵强度；未覆盖鉴权成功路径、错误体/字段、Schema fixture diff。
- closure: C3.4 行为矩阵补鉴权后 CRUD status/关键字段对比。

### F-IND-C32-002 · Manifest fragment 为占位 JSON
- C3.3 若直接用 stub 替换中心 adminModules 可能造成 Manifest 语义漂移。
- closure: C3.3 切换前对齐真实 fragment payload / 与中心投影可比对证据。

### F-IND-C32-003 · sixCapabilities 与 standardAdminCapabilities 双份维护
- 当前集合一致且 descriptorsMatch + 测试 fail-closed；长期漂移风险。
- closure: 导出共享 helper（已由编排器导出 kernel.StandardAdminCapabilities 并替换）。

### F-IND-C32-004 · RouteContribution.Middleware/Public 未填元数据
- 运行时认证已由 a.Middleware 嵌入 Handler，行为等价成立；C3.3 需明确发布规则。
- closure: C3.3 文档化「Handler 已含中间件，Middleware/Public 元数据可选」。

## 信息门禁

C3-I001/C3-I002 verified、C3-I003 collecting（最晚 C3.4）、C3-I004 non-blocking。
本 scope 无到期未关 required 阻断 C3.2 成立或进入 C3.3。

**明确声明：本独立审计员未修改任何 status / progress / goal-tree / 文件内容。**
