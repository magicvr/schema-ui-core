---
id: GOAL-045-w33-list-actions-slot-and-roles-trigger
doc: decision
status: done
parent: GOAL-001-design-implementation-conformance
created: 2026-09-19
updated: 2026-09-19
version: 1.0.0
---

# 决策记录 · GOAL-045

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-045-001 | required | 插槽宿主由谁提供、何时可见、首帧是否闪跳 | C2 | C1 | 读 custom 分支与 `App.tsx` portal 先例 | **verified** | — | `D-001` §2 |
| I-045-002 | required | 目标表缺失时隐藏 vs 原地渲染 | C2 | C1 | 对照既有回落约定 | **verified** | — | `D-001` §3（fail-open） |
| I-045-003 | non-blocking | roles 页启用选择对既有测试的影响 | C3 | C2 | 跑 roles 与列表视觉用例 | **verified** | — | 全量 vitest + 双 profile e2e 全绿 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-09-19 | W33 列表页 actions 左侧插槽与 roles 触发面方案冻结（用户选型 A + 载体 workspace-010 W33） | done | [D-001-w33-slot-and-roles-trigger-freeze.md](01-decision/D-001-w33-slot-and-roles-trigger-freeze.md) |

## 承接输入

- 用户 2026-09-19 两个问题及 P-004 选型（位置方案 **A** / 载体 **workspace-010 W33**），原文见 `00-meta` §概述。
- `[workspace-038]` 已交付事实：导出分母 `users`/`roles`（R3 `D-001` §1）、custom 节点位置规律与 `data-list-page-actions` 行现状。VP-038 保持 `closed`。

## 冻结结论

1. **插槽声明**：`props.slot = "list-page-actions"` + `props.targetTable`；未声明或未知值 → 原地渲染。
2. **宿主与可见性**：表格提供左段宿主并经 CRUD seam 发布；`registerListActionsSlot` 决定该行是否渲染。
3. **回落**：fail-open 原地渲染（布局能力不得隐藏操作入口）。
4. **行布局**：`justify-between` 两段，左段在前；右段语义/样式不变。
5. **表格身份（实施中补齐）**：渲染器按 table 节点 id 作 key——不同表不得共用组件实例（否则跨页串列状态）。该条为实施中由 e2e 暴露的问题所补，记入 `E-001` §3 问题 5。