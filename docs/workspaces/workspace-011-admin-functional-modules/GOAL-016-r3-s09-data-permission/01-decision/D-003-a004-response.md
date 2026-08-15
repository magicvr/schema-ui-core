---
id: D-003
goal: GOAL-016-r3-s09-data-permission
title: A-004 响应：S1 独立审计 required 全 fixed
date: 2026-08-15
status: accepted
parent: GOAL-016-r3-s09-data-permission
created: 2026-08-15
updated: 2026-08-15
version: 1.0.0
---

# D-003 · A-004 响应（2026-08-15）

## 审计结论

A-004（grok build · grok-4.6 · high · independent）verdict **conditional**，1 required：

- **F-001（required med）**：范围只注入 list()；detail/update/delete/batchDelete 按 id 放行构成 IDOR（一旦登记资源）；导出不经该 list；「零代码接入」与实体 List 不消费 Scope 矛盾；组合根应为 24→26。

## 响应（全 fixed）

1. **行访问全覆盖**：D-002 §2 扩展——Get/Update/Delete 按 owner 约束（self 不属本人 → 404 不泄露存在性）、BatchDelete 仅删本人行、Create 时 self 资源 owner_column=actor。
2. **导出面必办**：登记规范要求登记资源时同步评估 data-transfer 导出路径并施加同约束（v1 无登记资源无暴露面，接入点登记为后续目标必办）。
3. **ScopeAware 契约**：已登记资源实体必须实现 ScopeAware（工厂登记时校验，未实现拒绝登记——配置面 fail-closed）；修正「零代码接入」表述为两项动作。
4. **default_scope 必填**：PATCH policies 省略 → 400 INVALID_PATCH_FIELD（无隐式默认）。
5. **组合根修正**：admin 权限 24→26、导航 12→13（composition_test.go L465 核对）。

闭合路径：**fixed**（D-002 方案修正，S2 实施前可核对）。复审：A-005 grok reaudit。
