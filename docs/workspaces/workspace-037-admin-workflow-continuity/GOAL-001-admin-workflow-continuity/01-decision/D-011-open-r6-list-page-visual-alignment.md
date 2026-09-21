---
id: D-011-open-r6-list-page-visual-alignment
doc: decision-entry
status: accepted
parent: GOAL-001-admin-workflow-continuity
created: 2026-09-18
updated: 2026-09-18
version: 1.0.0
goal_id: GOAL-001-admin-workflow-continuity
---

# D-011 · 保持 R5/Root/VP 开放并开设 R6 列表页视觉收敛

## 决定

尊重用户“还不能关门根目标”的明确指令：不把 R5-I-004 用户书面确认推断为通过，不关闭 R5、Root 或 VP-037。在同一 Root/VP/workspace 下新增 `GOAL-007-list-page-visual-alignment`，承接参考 `raw/new-table` 的通用列表页视觉与筛选布局改造。

Root 的显式路线图从 5 个检查点扩展为 6 个：R1～R4 保持已完成，R5 仍为 `active · 3/4`，新增 R6 为 `active · 1/4`；Root 派生进度同步为 `active · 4/6`。R6 不改变 Charter 方向、不新建 VP/workspace，也不替代 R5 的关门审计与用户确认。

## 范围

- R6 只覆盖通用列表页的视觉层、筛选折叠、页面级按钮/视图位置与语义文案、列配置顺序、单页分页展示和对应回归。
- 顶部功能栏、左侧导航、后端 schema、查询/重置逻辑、Saved View 存储格式和未实现的多选批量动作保持边界不变。
- R6 的具体实施合同、信息需求和 C1 基线由 `GOAL-007-list-page-visual-alignment` 承接。

## 未选方案

- 不在 R5 的组合验收目标内偷偷加入产品实现：这会混淆 R5 已冻结的关门范围。
- 不创建新的 VP 或 workspace：R6 仍是 VP-037 Admin 工作流连续性方向内的实现层增量。
- 不关闭 Root/VP 后再以新目标“重开”：用户已明确当前不可关门，且 R5-I-004 尚未有书面确认。

## 后续

先完成 GOAL-007 C2/C3 实施与回归，再按 R6 C4 做 self 审计；R5-I-004 仍按原 `GOAL-006` 台账等待用户明确裁决。
