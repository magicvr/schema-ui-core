---
id: GOAL-004-s4-state-and-feedback
title: S4 · 状态与反馈一致性（Loading / Empty / 错误 / 表单异步反馈）
status: done
parent: GOAL-001-design-system-and-ui-experience
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
progress: 3/3
---

# GOAL-004 · S4 · 状态与反馈一致性

## 概述

承接 Root [GOAL-001-design-system-and-ui-experience](../GOAL-001-design-system-and-ui-experience/00-meta.md) 的 **S4 阶段**：让 `apps/web` 主范例路径（statCard / chart / list-table / form 提交异步反馈）的 Loading / Empty / 错误反馈可观察地一致，复用既有 `Skeleton` primitive 而不是散落的纯文本占位（如 `"Loading statCard…"` / `"Loading chart…"` / `"Loading…"` 文本行）。

**实施依据**：Root D-002/D-003/D-004（Token/视觉方向已 accepted）；不引入新的状态管理框架，只做"收敛到一致"。

## 现状盘点（实施前）

| 消费点 | 原始状态 | 问题 |
|--------|----------|------|
| `data-table.tsx`（list-table） | `loading` 时渲染纯文本 `"Loading…"` 单元格 | 未用 Skeleton |
| `render.tsx` StatCardView | `list === null` 时渲染纯文本 `"Loading statCard…"` | 未用 Skeleton；loading/error 判定逻辑内联、未抽出可单测纯函数 |
| `render.tsx` ChartView | `list === null` 时渲染纯文本 `"Loading chart…"` | 同上 |
| 错误态（table/statCard/chart/form） | 已一致使用 `role="alert"` + `text-destructive` | 已达标，本阶段保留 |
| Empty 态（table） | `emptyMessage` 纯文本，无 Skeleton 需求（空态本身不需要骨架屏） | 已达标，本阶段保留 |
| Form 提交异步反馈（`render.tsx` FormView） | `submitting` 态按钮文案变 "Submitting…" + disabled | 已是一致模式，本阶段保留 |

## 方案

- 新增共享纯函数 `resolveAsyncDisplayState({ loading, error, isEmpty })` → `"loading" | "error" | "empty" | "ready"`（`apps/web/src/components/ui/async-state.ts`），把"是否 loading / 是否 empty / 错误优先"的判定逻辑从渲染 JSX 中分离，可直接单测。
- `DataTable`、`StatCardView`、`ChartView` 三个消费点统一改用该纯函数决定分支，并在 `loading` 分支渲染 `Skeleton` primitive（`role="status"` + `aria-label`）取代纯文本占位。
- 不改变错误态（`role="alert"`）与 empty 态（纯文本 `emptyMessage`）的既有可观察行为，只统一 loading 态。
- Form 提交异步反馈（`submitting` → disabled + "Submitting…"）已符合一致性要求，不改动。

## 成功标准（可验收 · 等权检查点 · 共 3 项）

- [x] **C1**：新增 `resolveAsyncDisplayState` 纯函数并有直接单测覆盖 loading/error/empty/ready 四种优先级组合（`async-state.test.ts`）。
- [x] **C2**：`DataTable`（list-table）、`StatCardView`、`ChartView` 的 loading 态均改用 `Skeleton` primitive + `role="status"`，不再渲染纯文本 "Loading…" 占位；有直接驱动这些渲染路径的单测断言 Skeleton/status 存在。
- [x] **C3**：`apps/web` vitest 全量 + build 通过，无回归（既有错误态/空态断言保持不变）。

## 派生进度展示

`progress: 3/3` 由上方 3 个等权检查点派生；仅为展示，不放行阶段、不关闭 finding、不推导 Root `status: done`。

## 信息就绪与未知项

无独立信息项；继承 Root I-001/I-002/I-005（均 closed）。

## 父目标

[GOAL-001-design-system-and-ui-experience](../GOAL-001-design-system-and-ui-experience/00-meta.md)（Root；本目标为 S4 阶段子目标）

## 台账布局

本目标使用 ledger 目录：`01-decision/`、`02-execution/`、`03-audit/`。索引文件保留 frontmatter 与条目表；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*`。
