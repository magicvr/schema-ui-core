---
id: E-004-r6-layout-revision-implementation
doc: execution-entry
status: recorded
goal_id: GOAL-007-list-page-visual-alignment
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# E-004 · 完成 C5 布局修订与回归

2026-09-18，完成 R6 C5 的通用列表布局修订：

- `apps/web/src/renderer/schema-table.tsx` 将视图切换/“保存为视图”保留在标题区右上角，将包含“列配置”的页面 actions、Saved View 管理操作和其他页面级操作放回筛选栏下方、列表上方；顶部功能栏、左侧导航、Saved View 存储格式和未实现的多选操作未改动。
- `apps/web/src/components/list-filter-panel.tsx` 将筛选操作单元放入筛选网格最后一行最右侧；查询按钮由 `apps/web/src/renderer/render.tsx` 移入该操作单元，与重置及展开/收起额外筛选项共同保持现有提交/重置逻辑。
- `apps/web/src/components/data-table.tsx` 与 `schema-table.tsx` 将分页 footer 放入列表 surface 内，并保持单页前后页、当前页、跳页和每页数量的既有状态逻辑；有效列表响应即使只有一页仍显示分页控件。
- 相关回归更新于 `data-table.test.tsx`、`schema-table.test.tsx`、`search-form-filters.test.tsx`；Saved View 集成行为继续由既有测试覆盖。

验证事实：

- `npm exec -- tsc -b --pretty false`：通过（**更正**：原记 `tsc --noEmit` 属空转，见 A-002 F-005；`apps/web/tsconfig.json` 为 solution-style `{"files": []}`，裸 `tsc --noEmit` 不检查任何文件）。
- 定向 Vitest：4 个文件、58 个测试通过；查询动作修复后的定向回归另通过 3 个文件、91 个测试。
- 全量 Web Vitest：111 个文件、1413 个测试通过。
- `git diff --check`：通过。

C5 已完成；C6 修订审计与 Root/VP 的 R6 完成投影待后续审计记录。R5-I-004、Root 和 VP 仍保持开放。
