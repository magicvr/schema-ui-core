---
id: D-001-s4-baseline-and-plan
title: S4 状态一致性盘点与方案
date: 2026-08-09
status: accepted
parent: GOAL-004-s4-state-and-feedback
---

# D-001 · S4 状态一致性盘点与方案

## 现状盘点（实施前）

| 路径 | Loading | Empty | 错误 |
|------|---------|-------|------|
| `components/data-table.tsx`（statCard/chart 之外的 list-table 主渲染） | 纯文本 `Loading…`（`<td>` 文案） | 纯文本 `emptyMessage` | 纯文本 `text-destructive`（无 `role="alert"`） |
| `render.tsx` `StatCardView` | 纯文本 `Loading statCard…` | N/A（无空态分支） | `role="alert"` `<p>` |
| `render.tsx` `ChartView` | 纯文本 `Loading chart…` | 纯文本 `chart has no plottable data points` | `role="alert"` `<p>` |
| `renderer/form-controls.tsx` 表单提交异步反馈 | 按钮文案 `Submitting…` + disabled | N/A | `role="alert"` `<span>`（upload 字段） |
| `render.tsx` `FormView` 提交错误 | 同上 | N/A | `role="alert"` `<p>`（`formError`） |

**结论**：三条主渲染路径（list-table / statCard / chart）各自用临时纯文本表达 loading，未复用已存在的 `Skeleton` primitive（`components/ui/skeleton.tsx`），且 DataTable 的错误态缺 `role="alert"`。表单提交异步反馈（disabled + `Submitting…` + `role="alert"` 错误）已经内部一致，本次不改。

## 方案（accepted）

1. 新增 `components/ui/async-state.ts`：纯函数 `resolveAsyncDisplayState({ loading, error, isEmpty }) → "loading" | "error" | "empty" | "ready"`，优先级 error > loading > empty > ready。三条主渲染路径统一改为调用该函数决定分支，不各自发明判定顺序。
2. `DataTable`（服务于 list-table 主范例）：loading 分支渲染 `role="status"` + 3 条 `Skeleton` 行；error 分支加 `role="alert"`；empty 分支不变（仍是 emptyMessage 文案，非视觉回归）。
3. `StatCardView` / `ChartView`：loading 分支从纯文本换成 `role="status"` + `Skeleton` 块（statCard 两条骨架条模拟 label+value；chart 一条骨架块模拟图形区域）；error/empty 分支保留既有 `role="alert"` / 提示文案不变（不引入新的空态语义）。
4. 不改动表单提交异步反馈（已一致），不引入新的状态管理库或 Suspense 边界。

## 未选方案

- 引入 React Suspense + `<Skeleton>` fallback：需要把 fetch 改成支持 Suspense 的数据源，改动面过大、超出 S4 范围（收敛现状 vs 重构数据层），放弃。
- 为每个组件单独写局部 loading 判定：会重复三份几乎相同的 if/else，且不可单测（判定逻辑内嵌在 JSX 里），放弃，改为抽纯函数集中单测。
