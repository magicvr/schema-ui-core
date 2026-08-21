---
id: D-004
doc: decision
status: accepted
parent: GOAL-014-w13-settings-tabs-and-topbar
created: 2026-08-16
updated: 2026-08-16
version: 0.1.0
---

# D-004 · T-07 列表筛选即时生效（用户点名遗留逻辑）

## 背景

2026-08-16 用户点名：所有列表页筛选——文本框输入或下拉选择后，只出现【筛选记录】chips，但列表不重新筛选，必须点击「搜索」才生效。期望：

1. 筛选项变动（新增或移除条件）→ 变更 chips 显示的同时**立即**重新输出筛选后的列表；
2. 文本框+搜索按键组合：输入时不筛选、也不显示 chips 变更，点击对应「搜索」按键后才变动筛选项并重新输出列表。

## 决策

### 双轨提交模型（即时控件 vs 提交式控件）

- **即时生效（非文本框控件）**：select / datePicker / dateRangePicker / inputNumber 等筛选控件变更（新增条件、改值、清空移除）→ 立即调用 `searchFormSubmit`（更新目标表 query → 列表重新获取，page 重置为 1），chips 同步更新。
- **提交式（文本框 input）**：关键词 `q` 输入只更新本地表单值——不触发查询、不显示 chips 变更；点击配对的「搜索」按键（form submit）后提交查询，chips 与列表同时生效。
- **chips 真相源**：筛选记录从**已提交的目标表 query**（`crud.tableQuery(targetTable)` 的 q / filters）反推，而不是本地草稿 values——未提交的输入永不显示为筛选记录；重置/移除 chip 仍即时提交（既有逻辑保留）。

### 实现

- `apps/web/src/renderer/render.tsx` FormInner：
  - `handleFieldChange`：更新本地 values；搜索模式下非 `input` 类型字段 → 立即 `crud.searchFormSubmit(node, next)`；
  - `activeFilters` 改从 `crud.tableQuery(targetTable)` 反推（q → query.q；其余字段 → query.filters[id]）。
- 不改协议 schema、不改 Go、不改表格级 `props.filters`（本就即时生效）。

## 影响

- **go 判定**：纯前端渲染层改动 → Profile 默认集 / 模块矩阵 / Manifest 装配语义不变 → **VP-008 go 无影响、不暂挂**。
- **测试影响**：`search-form-filters.test.tsx` 重写为新交互契约（select 即时请求、q 输入不请求、搜索键提交、chips 跟随已提交条件、chip 移除即时生效、reset 保留）；其余页面级测试回归。
