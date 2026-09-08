---
title: 目标树 · workspace-034-nav-group-collapsible
status: done
created: 2026-09-07
updated: 2026-09-08
parent: null
version: 0.1.0
workspace_id: workspace-034-nav-group-collapsible
---

# 目标树 · Admin 导航分组折叠体验

> 工作区：`workspace-034-nav-group-collapsible`
> canonical：`docs/workspaces/workspace-034-nav-group-collapsible/`
> Root：`GOAL-001-nav-group-collapsible`（原始有界交付容器 · done 5/5）
> primary_plan：`VP-034-nav-group-collapsible`（active · v0.3.0）

## 目标树

```text
GOAL-001-nav-group-collapsible [done · 5/5]
├── GOAL-002-navigation-group-polish [done · 2/2]（结项后增量）
├── GOAL-003-sidebar-engine-navigation [done · 3/3]（Sidebar Engine 与通用详情抽屉视觉增量）
└── GOAL-004-sidebar-active-indicator [done · 2/2]（组间距与页面 active 指示微调）
```

## 纲领路线图

```text
R1 导航清单 / Profile-slot 矩阵 / 分组 IA 冻结 [completed]
 → R2 模块注册与跨模块聚合契约 [completed]
 → R3 Shell 折叠展开 / 可访问性 / 直接 URL 自动展开 [completed]
 → R4 当前 sidebar 导航全量迁移与 default/optional/custom/demo 回归 [completed]
 → R5 证据矩阵 / 全量回归 / Goal 审计 / 关门准备 [completed]
```

## 状态表

| id | title | status | progress | parent | notes |
|----|-------|--------|----------|--------|-------|
| GOAL-001-nav-group-collapsible | Admin 导航分组折叠体验交付 | **done** | 5/5 | null | 原始 R1-R5 交付已完成；后续样式/默认折叠增量由 GOAL-002 承载；VP-034 仍 active，另走 `/vision` |
| GOAL-002-navigation-group-polish | 导航分组语义样式与默认折叠增量 | **done** | 2/2 | GOAL-001-nav-group-collapsible | P1/P2 完成；A-001 self `pass`；默认关闭、active 深链与语义样式已验证；Git checkpoint `fcd6fe9b` 已创建 |
| GOAL-003-sidebar-engine-navigation | Sidebar Engine 导航与通用详情抽屉视觉优化 | **done** | 3/3 | GOAL-001-nav-group-collapsible | P1/P2/P3 完成；A-001 self `pass`；API/Web 全量回归与 scope-specific browser checks 通过；checkpoint `de71fff0` |
| GOAL-004-sidebar-active-indicator | Sidebar 分组间距与页面激活指示微调 | **done** | 2/2 | GOAL-001-nav-group-collapsible | P1/P2 完成；A-001 self `pass`；space-y-2、左侧 active 竖线、secondary 替换光点已验证；全量 Web 回归与构建通过 |

## 维护说明

- Root `progress: 100%` 由 GOAL-001 `00-meta.md` 的 5 个显式检查点派生；GOAL-002 `progress: 100%` 由其 P1/P2 检查点派生；GOAL-003 `progress: 100%` 由其 P1/P2/P3 检查点派生；GOAL-004 `progress: 100%` 由其 P1/P2 检查点派生。
- 新建阶段子目标前，先在 Root 决策/路线图中冻结阶段边界；目标文件夹在本工作区根平铺。
- status/progress/parent 或新增子目标发生变化时，必须同步本文件树与状态表。