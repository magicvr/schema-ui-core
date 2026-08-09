---
id: E-009-row-action-no-drawer
title: 行 Edit/Delete 不再打开 recordView Drawer
date: 2026-08-09
status: recorded
parent: GOAL-001-design-system-and-ui-experience
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
---

# E-009 · 行 action 与详情 Drawer 解耦

## 触发

用户驳回 closeout：users 列表点 Edit 等 action 会打开 record details 侧栏，不符合使用逻辑；要求一并修所有列表。

## 根因

1. `invokeAction` / `openModal` 调用 `setSelectedRow(row)`，而 selection 驱动 `recordView` Drawer。
2. 次要：action 按钮点击可能冒泡到 `tr.onRowClick`（已用 stopPropagation + interactive 判定加固）。

## 修复

| 位置 | 改动 |
|------|------|
| `render.tsx` | `invokeAction` / `openModal` **不再** `setSelectedRow`；modal 预填仍用 `activeModal.row` / `modalRow` |
| `data-table.tsx` | 行点击忽略 button/input/`[data-row-click-ignore]`；actions/selection 列 stopPropagation |
| `schema-table.tsx` | action 按钮与 checkbox stopPropagation |
| 测试 | `data-table.test.tsx` + `visual-fidelity.test.tsx` Edit 不打开 drawer |

## 验证

vitest 全绿（含 schema-crud Edit 预填仍通过）。
