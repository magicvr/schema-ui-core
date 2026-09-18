---
id: E-005-r6-c7-control-placement-implementation
doc: execution-entry
status: recorded
goal_id: GOAL-007-list-page-visual-alignment
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# E-005 · 完成 C7 控制位纠偏与回归

2026-09-18，按 D-003 完成 R6 C7 的三项 UI 纠偏：

1. **图标视觉强调**（`apps/web/src/renderer/schema-table.tsx`）：“列配置”触发器（`data-saved-view-columns-trigger`）加列图标，“保存为视图”按钮加加号图标；两者图标均 `aria-hidden="true"`，可见文案不变，可访问名称未受影响。
2. **搜索配对回调**（`apps/web/src/renderer/render.tsx`）：恢复 `searchButtonSlot`，搜索提交按钮重新渲染在关键词输入所属网格单元内，并恢复 `-ml-px` 重叠与 `rounded-l-none`；筛选网格末位操作单元（`data-filter-actions`）只保留重置与展开/收起。`form-controls.tsx` 的 `paired`/`rounded-r-none` 机制在 C5 期间未被改动，本轮直接复用。
3. **视图表单归属**（`apps/web/src/renderer/schema-table.tsx`）：新增 `data-saved-view-management` 容器，把“更新视图”“删除视图”与新建视图表单移入视图 surface（`data-saved-views-surface`）内部、标签组（`data-saved-view-tabs`）下方；页面级 actions 行（`data-list-page-actions`）不再承载视图管理操作，仍只保留列配置与 schema toolbar 触发器。查询、重置、Saved View 存储格式与分页状态逻辑均未改动。

回归更新于 `search-form-filters.test.tsx`（配对规则断言回到同格 + `-ml-px`）、`saved-views.ui.test.tsx`（新增视图表单位置与图标两项）与 `schema-table.test.tsx`（筛选操作单元不再持有提交按钮）。

验证事实：

- `npm exec -- tsc --noEmit --pretty false`：通过。
- 定向 Vitest：4 个文件、60 个测试通过。
- 全量 Web Vitest：111 个文件、1415 个测试通过。
- `git diff --check`：通过。
- 真实 Chromium 布局核对（临时预览页，核对后已删除）：视图表单距“保存为视图”按钮 4px，且不在页面级 actions 行内；搜索按钮与其输入水平重叠 1px、垂直偏移 0；两个触发器均含图标；无 console 错误。
- Git checkpoint：C5 轮次的已验证快照以 `0cc18418` 提交（保留原轮次边界），本轮 C7 以 `e4b8c5fa` 提交；两者均只暂存本轮 owned paths，未使用 `git add -A`。

C7 实现与回归已完成；C6 修订审计与 Root/VP 的 R6 完成投影待后续审计记录。R5-I-004、Root 和 VP 仍保持开放。
