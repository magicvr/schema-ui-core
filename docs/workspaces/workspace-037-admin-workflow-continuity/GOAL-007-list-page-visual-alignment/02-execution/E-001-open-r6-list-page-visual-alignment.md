---
id: E-001-open-r6-list-page-visual-alignment
doc: execution-entry
status: recorded
parent: GOAL-001-admin-workflow-continuity
created: 2026-09-18
updated: 2026-09-18
version: 1.0.0
goal_id: GOAL-007-list-page-visual-alignment
---

# E-001 · 开设 R6 并完成范例/现状基线

## 事实

2026-09-18，按用户明确指令在 `[workspace-037-admin-workflow-continuity]` 下新增 `GOAL-007-list-page-visual-alignment`，父目标为 `GOAL-001-admin-workflow-continuity`。Root/VP 保持 `active`，R5 用户书面关门确认仍开放。

已完成只读基线盘点并形成 D-001：

- 范例页是静态 HTML，位于 `raw/new-table/schema_ui_core_2/code.html`，视觉说明位于 `raw/new-table/monochrome_technical/DESIGN.md`。
- 通用列表调用链为 `App.tsx` → `renderer/schema-table.tsx` → `components/data-table.tsx`；search form/query bridge 位于 `renderer/render.tsx`。
- 现有 token 可覆盖范例的 surface/outline 语义，无需新增全局 token。
- 查询/重置逻辑与单页分页旧行为已定位到实现和测试；本轮只调整展示条件与布局。
- 顶部功能栏与左侧导航边界已登记，R6 不修改 shell 结构。

## 当前状态

- R6：`active · 1/4`；C1 完成。
- 实现尚未开始；C2/C3/C4 仍未完成。
- 多选批量功能未纳入 R6 必做范围。
