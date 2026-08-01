---
title: 决策 · 生产级可用 Admin 基架
status: active
created: 2026-08-01
updated: 2026-08-01
parent: null
version: 0.1.0
---

# 决策 · GOAL-001

## D-001 · 以新 delivery 工作区承接 VP-002

- **日期**：2026-08-01
- **决定**：建立 `workspace-002-production-admin-foundation` 与 Root `GOAL-001-production-admin-foundation`，将 VP-002 从 `planned` 激活为 `active`，并把本工作区设为其当前唯一 lead workspace；仓库 `primary_workspace` 保持 `workspace-001-mvp-admin-foundation`。
- **依据**：用户明确要求开启新工作区和 Root 承接 VP-002，并确认了工作区与 Root 命名。VP-001 与旧 Root 均已关闭/完成，新波次需要独立 canonical scope 和目标树。
- **边界**：不重开 VP-001 或旧 Root；不建立跨工作区 `parent`；只通过 Q2 路径引用已冻结历史基线。

### 未选方案

- **复用 workspace-001 并重开旧 Root**：会混合已关闭波次与新实施事实，破坏状态边界。
- **把新工作区设为 primary**：当前没有长期目的或仓库北极星换代，不应改写 Charter 的 `primary_workspace`。
- **把 VP-002 作为旧 Root 的子目标**：跨工作区 parent 被协议禁止，也无法提供独立交付树。

## D-002 · 采用五阶段串行纲领路线图

- **日期**：2026-08-01
- **决定**：Root 采用 `Renderer → 认证 → 持久化权限 → CRUD → 工程化与关门` 五个等权检查点，纲领阶段原则上串行，阶段内部可按依赖并行。
- **原因**：该顺序遵循 VP-002 的价值链，同时把原第三阶段拆成可独立验证的权限持久化、CRUD 和工程交付，避免一个阶段承载过多门禁。
- **执行约束**：只在当前阶段边界和 required 信息项就绪后创建具体子目标；`progress` 只按完成检查点数派生，不用于放行。

## D-003 · 先登记未知，再按阶段关闭

- **日期**：2026-08-01
- **决定**：将协议实施差量、认证机制、持久化模型、代表性 CRUD 实体及部署/fork 验收口径登记为阶段 required 信息项；操作日志保持 non-blocking。
- **原因**：这些未知不妨碍 Root 立项，但会改变对应阶段的方案或验收，必须在最晚门禁前以证据关闭或经用户书面接受有界 residual。
