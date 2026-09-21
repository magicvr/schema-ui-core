---
id: D-012-reopen-r6-layout-revision
doc: decision-entry
status: accepted
goal_id: GOAL-001-admin-workflow-continuity
created: 2026-09-18
updated: 2026-09-18
parent: null
version: 1.0.0
---

# D-012 · 用户选择回开 R6 并追加列表布局修订

## 决定

用户在 R6 已完成后提出同一列表视觉范围的进一步修订，并选择方案 A：回开既有 `GOAL-007-list-page-visual-alignment`，保留原 C1～C4 历史事实，追加 C5/C6；不新建 GOAL-008。

Root 的六阶段分母不变，但 R6 检查点因修订回到未完成状态，因此 Root 从 `active · 5/6` 回投影为 `active · 4/6`。R5 `GOAL-006`、R5-I-004、Root 和 VP-037 均保持 active。

## 范围

- 页面标题右上角只显示视图切换/保存视图标签组。
- 含“列配置”的页面 actions 移到筛选栏下面、列表上面。
- 筛选 reset/展开收起与筛选控件同属网格，位于最后一行最右。
- 分页改为范例页列表内 footer 样式和布局，保留现有逻辑。

不得修改顶部功能栏、左侧导航、后端 schema、Saved View 存储格式、查询/重置/分页状态逻辑或未实现的多选能力。

## 理由

这是对 R6 同一交付分母和同一范例合同的局部纠偏；回开可保留原审计证据的历史边界，同时避免重复目标和新的 Root 阶段。
