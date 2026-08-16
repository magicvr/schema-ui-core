---
id: E-013
doc: execution
status: recorded
parent: GOAL-014-w13-settings-tabs-and-topbar
created: 2026-08-16
updated: 2026-08-16
version: 0.1.0
---

# E-013 · S2 T-07 列表筛选即时生效实施

- `apps/web/src/renderer/render.tsx`（FormInner，搜索模式）：
  - 新增 `handleFieldChange`：字段变更更新本地 values；搜索模式下非 `input` 类型控件（select/日期等）→ **立即** `crud.searchFormSubmit(node, next)`——列表马上按新条件重新获取（page 重置 1）；`input` 文本框只更新草稿值（点击配对的搜索按键才提交）。
  - `activeFilters`（筛选记录 chips）改为从**已提交的目标表查询**反推（`crud.tableQuery(targetTable)`：q → query.q；其余字段 → query.filters[id]）——未提交的输入不显示为筛选记录，chips 与列表实际生效条件严格一致。
- 保留：搜索按键（submit）提交、重置按钮、chip 逐条移除（均即时提交）；表格级 `props.filters` 本就即时生效，未改动。
- 无协议/schema/Go 变更。
