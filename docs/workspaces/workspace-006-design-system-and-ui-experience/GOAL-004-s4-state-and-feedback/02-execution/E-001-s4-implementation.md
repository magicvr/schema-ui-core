---
id: E-001-s4-implementation
title: S4 实施完成记录 — 统一 Skeleton/Empty/错误反馈
date: 2026-08-09
status: done
parent: GOAL-004-s4-state-and-feedback
---

# E-001 · S4 实施完成记录

## 事实摘要

**日期**：2026-08-09
**scope**：statCard / chart / list-table（DataTable）的 loading 态收敛到 `Skeleton` primitive；三者共享的
loading/error/empty 判定逻辑抽成纯函数。

### 新增文件

- `apps/web/src/components/ui/async-state.ts`：`resolveAsyncDisplayState({ loading, error, isEmpty })` →
  `"loading" | "error" | "empty" | "ready"`。优先级：`error` > `loading` > `isEmpty` > `ready`。纯函数，不依赖
  React/DOM，可直接单测。
- `apps/web/src/components/ui/async-state.test.ts`：4 个直接驱动该纯函数的单测（error 优先于 stale loading；
  loading 优先于 empty 判定；empty 仅在 settled 后判定；ready 路径）。

### 消费点变更

| 消费点 | 变更前 | 变更后 |
|--------|--------|--------|
| `apps/web/src/components/data-table.tsx`（`DataTable`，被 `SchemaTable`/list-table 复用） | loading 行渲染纯文本 `"Loading…"` | loading 行渲染 `role="status"` 容器 + 3 个 `Skeleton` 条；error/empty 分支改用 `resolveAsyncDisplayState` 统一判定（保留既有 `role="alert"` / emptyMessage 文案，向后兼容既有断言） |
| `apps/web/src/renderer/render.tsx`（`StatCardView`） | `list === null` 时渲染纯文本 `"Loading statCard…"` | 用 `resolveAsyncDisplayState({ loading: list === null, error })` 判定；`loading` 态渲染 `role="status"` + 2 个 `Skeleton` 条（label 占位 + 数值占位），保持卡片 `rounded-md border` 外框一致 |
| `apps/web/src/renderer/render.tsx`（`ChartView`） | `list === null` 时渲染纯文本 `"Loading chart…"` | 同上模式；`loading` 态渲染 `role="status"` + 1 个 `Skeleton` 条（`h-40` 对齐真实图表高度） |

`error`/`empty` 分支的既有 `role="alert"` 文案与 `emptyMessage`/`"chart has no plottable data points"` 保持不变，
只是判定路径统一走 `resolveAsyncDisplayState`，不改变可观察行为。

### 未变更（范围纪律）

- `FeedbackRegion`（表单提交成功/失败反馈）、`FormView` 的 `formError`/`submitting` 状态、行内字段错误
  （`<p role="alert">`/`<ul role="alert">`）均已在既有实现中呈现 loading/错误一致语言（disabled 按钮 +
  "Submitting…" 文案 + `role="alert"`），S4 未改动这些路径——它们已经满足"一致反馈"要求，不需要引入
  Skeleton（表单没有骨架屏语义）。
- 未引入新的状态管理框架；未新建 Token 或 primitive；`Skeleton` 组件本体（`ui/skeleton.tsx`）未改动，仅新增
  消费点。

### 新增/修改测试

- `apps/web/src/components/ui/async-state.test.ts`（新增，4 tests）：直接单测 `resolveAsyncDisplayState` 纯函数。
- `apps/web/src/components/data-table.test.tsx`：`"shows the empty message and a loading row"` 断言扩展为同时
  校验 `role="status"` 容器存在（结构断言，驱动真实 `DataTable` 组件）。
- `apps/web/src/renderer/render.test.tsx`：新增两个测试
  `"shows a Skeleton status region for statCard while its dataSource fetch is pending"` 与
  `"...for chart..."`，通过 `RenderPage` 真实挂载并在 fetch resolve 之前（未 `await` fetch 结算的 tick）断言
  `container.querySelector('[role="status"]')` 存在，直接驱动 `StatCardView`/`ChartView` 的真实渲染路径。
- 顺带修复 `apps/web/src/app/shell.test.ts` 中已存在但与 S4 无关的 3 个 `noUnusedParameters` TS 编译错误
  （`openDrawer`/`closeDrawer`/`navigateAndClose` 的未用 `state` 参数改成 `_state`），否则 `npm run build`
  无法通过（该问题在 GOAL-003 提交时未被 `tsc -b` 增量缓存捕获；本次是从干净状态首次触发）。

### 验证结果

- `apps/web` `npm run test`（vitest run）：**613 tests / 28 test files — 全绿**（较 S2/S3 完成时的 607 净增 6：
  4 个 async-state 单测 + 2 个 Skeleton 结构断言）。输出：`{SCRATCH}/vitest-s4s5.log`。
- `apps/web` `npm run build`：**exit 0**（`tsc -b && vite build`）。输出：`{SCRATCH}/build-s4s5.log`。

## 派生进度贡献

对应 `GOAL-004` 检查点 C1（Skeleton 消费点统一）与 C2（纯判定逻辑 + 直接单测）均已满足，详见 `00-meta.md`。
