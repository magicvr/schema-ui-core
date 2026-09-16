---
id: D-002-open-r1-scope-freeze
doc: decision
goal_id: GOAL-001-admin-workflow-continuity
status: accepted
created: 2026-09-17
updated: 2026-09-17
parent: GOAL-001-admin-workflow-continuity
version: 0.1.0
---

# D-002 · 建立 R1 分母与语义冻结子目标

## 决策

在 Root 之下按纲领路线图开设平铺子目标 `GOAL-002-r1-scope-semantics-freeze`，承载 R1 的代码事实盘点、可核对矩阵、用户裁决前的候选方案和阶段审计。

## 原因与边界

- 当前 VP-037 的 R2～R4 方案依赖四项 required 信息；直接进入实现会把未确认的持久化、dirty-state 和反馈语义混入代码。
- 子目标只推进 R1 的信息就绪与方案冻结，不提前创建 R2～R4 的实现目标，也不改变 Root `active · 0/5`。
- Saved View 持久化等会改变实现成本的方案仍须单独取得用户确认；本决策不替用户选择。

## 证据

- `docs/workspaces/workspace-037-admin-workflow-continuity/GOAL-002-r1-scope-semantics-freeze/00-meta.md`
- `docs/workspaces/workspace-037-admin-workflow-continuity/goal-tree.md`
