---
id: GOAL-002-navigation-group-polish
title: 导航分组语义样式与默认折叠增量
status: done
parent: GOAL-001-nav-group-collapsible
created: 2026-09-07
updated: 2026-09-07
version: 0.1.0
progress: 100%
plan_refs:
  - VP-034-nav-group-collapsible
primary_plan: VP-034-nav-group-collapsible
serves_summary: 在已完成的导航分组契约上补充工程化语义样式，并将非 active 分组默认关闭、active 直达/深链自动展开。
---

# GOAL-002 · 导航分组语义样式与默认折叠增量

## 概述

这是 workspace-034 在 Root `GOAL-001-nav-group-collapsible` 结项后的有界增量目标，承接同一 VP-034，不改写父目标已完成的历史事实。目标聚焦两个用户明确要求：

1. 优化导航分组样式，保持现有 Tailwind/shadcn/Linear/Vercel 的语义化视觉体系，新增清晰的分组层级、状态与交互样式；
2. 分组默认关闭；只有当前访问分组内的直接页面或已登记内页/动态深链时自动展开，手动折叠仍可保持。

父目标 `GOAL-001-nav-group-collapsible` 继续保持 `done 5/5`；本目标是后续增量，不把父目标重新标记为 active。

## 愿景对齐

- Charter：`schema-ui-core-admin-foundation@0.4.0`
- VP：`VP-034-nav-group-collapsible`（`active` · v0.3.0）
- 工作区：`workspace-034-nav-group-collapsible`（follow-up delivery increment）
- 父目标：`GOAL-001-nav-group-collapsible`（`done 5/5`）

## 范围

### 纳入

- sidebar/mobile drawer 分组标题的工程化语义样式：层级、边界、hover/focus、active 组与展开指示。
- 继续使用现有语义颜色/间距/圆角体系，不引入新的主题变量或业务模块分支。
- 非 active 分组默认关闭；当前组内直接页面、已登记内页/动态深链自动展开。
- `sessionStorage` 继续保存手动 open/closed 偏好；storage 不可用/损坏时默认关闭，active 深链仍优先自动展开。
- 分组样式、默认折叠、键盘操作、sessionStorage 与深链回归。

### 不纳入

- 修改 VP-034 五组成员、组序、API group 契约或 Manifest schema。
- 将 top/user slot 或 Dashboard 单例改造成分组。
- localStorage、服务端持久化、多级分组、拖拽排序或新主题色体系。
- 重开父目标或修改 VP-034 status。

## 纲领路线图

以下 2 个检查点是本目标 progress 的唯一来源，默认等权：

| 检查点 | 目的 | 状态 |
|---|---|---|
| P1 | 导航分组语义样式在 desktop/mobile/sidebar 层级与 active/focus 状态上完成 | completed |
| P2 | 默认关闭、active 直达/深链自动展开、手动状态保持与回归完成 | completed |

`progress: 100%` = 2/2 个检查点完成；不替代 status、审计或验收。

## 信息需求（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚阶段 | 验证 / 收集动作 | 状态 | 证据 |
|---|---|---|---|---|---|---|---|
| I-002-001 | non-blocking | 新增分组样式应复用哪些现有视觉语义，而不引入新主题变量？ | P1 | P1 | 对照现有 NavigationLink、accent/muted/border/focus 类并补样式回归 | verified（用户要求 + 实现/测试） | D-001；E-001；`apps/web/src/app/nav-groups.test.tsx` |
| I-002-002 | required | 非 active 分组默认是否关闭、active 直达/深链是否自动展开？ | P2 | P2 | 用户确认 + App/navigation 测试 | verified（用户决策 + 测试） | D-001；E-001 |
| I-002-003 | non-blocking | sessionStorage 偏好与 active 深链优先级是否保持既有 D-005 语义？ | P2 | P2 | 复用现有 storage namespace 并补优先级测试 | verified（继承决策 + 测试） | D-005；D-001；E-001 |

## 父级关系

- `parent: GOAL-001-nav-group-collapsible`；父目标完整 id 与当前 workspace 相同。
- 父目标已结项；本增量只在本目标五件套与当前 workspace goal-tree 中保存新增状态。

## 台账布局

本目标使用五件套与 `01-decision/`、`02-execution/`、`03-audit/` 三个平铺 ledger 目录。
