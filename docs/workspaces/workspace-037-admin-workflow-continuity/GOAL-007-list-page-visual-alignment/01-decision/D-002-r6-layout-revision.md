---
id: D-002-r6-layout-revision
doc: decision-entry
status: accepted
goal_id: GOAL-007-list-page-visual-alignment
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# D-002 · 回开 R6 并修订范例页布局合同

## 决定

用户选择方案 A：回开既有 `GOAL-007-list-page-visual-alignment`，保留 C1～C4 的历史事实，追加 C5/C6，不新建 GOAL-008。

C5 的布局合同为：

1. 页面标题区右上角只承载视图切换与“保存为视图”分段标签组；Saved View 的选择/保存/更新/删除状态与存储合同保持不变。
2. 页面 actions（包括“列配置”）位于筛选栏下面、列表上面；列配置仍是页面 actions 的最左/第一个操作项。
3. 筛选栏默认折叠且只显示第一行字段；查询、重置、展开/收起操作归入同一筛选网格的操作单元，放在可见筛选项的最后一行最右侧。查询与重置继续使用当前 handler。
4. 分页采用范例页的列表内 footer 布局：左侧记录统计/页码，右侧每页数量、页码导航、跳页；单页仍显示 footer，操作按现有边界禁用。
5. 顶部功能栏、左侧导航、后端 schema、全局语义 token 名称、Saved View 持久化格式、查询/重置/分页状态逻辑和未实现的多选能力均不扩展或重写。

## 理由

本轮是对已交付 R6 的布局边界校正，用户反馈明确指出原实现把整个 page actions 组错误地放进标题右上角。问题仍属于 R6 的同一分母和同一范例对齐范围；回开目标比新建高度重叠的修订目标更易追踪，也保留原 A-001 作为历史版本证据。

## 未选方案

- 不保留 GOAL-007 `done` 并新建 GOAL-008：会把同一列表视觉合同拆成两个重叠目标，增加 Root 阶段和台账重复成本。
- 不把页面 actions 继续 portal 到标题区：这违反范例页“筛选栏 → 页面 actions → 列表”的层级。
- 不通过重写 query/controller 实现视觉调整：本轮只修订 DOM、布局和样式，避免触碰已验证的状态语义。

## 门禁

C5 完成前不沿用旧 A-001 作为本轮交付结论；C6 需在新实现、定向回归和自审证据完成后追加审计意见。R5-I-004 仍开放，R6 不得关闭 R5、Root 或 VP。
