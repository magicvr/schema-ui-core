---
id: workspace-036-admin-command-palette
title: Admin 全局检索与 Command Palette 工作区
status: active
root_goal: GOAL-001-admin-command-palette
canonical_scope: docs/workspaces/workspace-036-admin-command-palette/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-036-admin-command-palette
primary_plan: VP-036-admin-command-palette
created: 2026-09-14
updated: 2026-09-14
version: 0.1.0
parent: null
---

# 工作区上下文 · Admin 全局检索与 Command Palette

本工作区是 [VP-036-admin-command-palette](../../vision/plans/VP-036-admin-command-palette.md)（**`active` · v0.2.0**）的唯一 delivery workspace，承接已注册页面、导航项与声明式动作的权限安全检索、Command Palette、跨模块 provider 聚合、键盘可访问性与导航分组联动。不重开 VP-034，不替代 VP-009/VP-010，不承载实体全文搜索或新的业务域。

- Root `GOAL-001-admin-command-palette`：**`done · 4/4`**；2026-09-14 用户书面确认关门，VP-036 `closed` 由 `/vision` 另行执行。
- 激活门禁已满足（2026-09-14）：[VRev-092](../../vision/reviews/VRev-092-vp036-admin-command-palette-activation.md) self `pass`；Admin freshness `5c341ec7` → `97aefe8c` PASS；I-036-006 verified。
- 用户已确认工作区 slug = `workspace-036-admin-command-palette`；Root = `GOAL-001-admin-command-palette`。
- Vision open required：0；V-F123 为 inherited recommended，由 R1 的可机器核对矩阵承接。
- 红线：不实现或解除 `RT-X01` / `RT-X02`、Redis、MQ、多实例、跨进程索引；不把实体记录全文搜索、Saved Views、批量结果中心、未保存保护或 Toast 全局重做混入首波。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-036-admin-command-palette` | 与本区目标及资料引用的 `workspace_id` 一致；当前无固定共享资料 |
| Root Goal | `GOAL-001-admin-command-palette` | `parent: null`；**done · 4/4** |
| canonical 范围 | `docs/workspaces/workspace-036-admin-command-palette/` | 本区唯一目标状态范围 |
| 共享资料目录 | `none` | 本区暂无固定共享资料；不得声明共享资料引用 |
| 愿景角色 | `delivery` | VP-036 唯一 delivery workspace；不改变 Charter primary workspace |
| 规划对齐 | `primary_plan` = `VP-036-admin-command-palette` | `plan_refs` 必填且已精确绑定 |

## 愿景对齐

- Charter：`schema-ui-core-admin-foundation@0.4.0`
- VP：`VP-036-admin-command-palette`（**`active` · v0.2.0**）
- 计划审视：[VRev-090](../../vision/reviews/VRev-090-vp036-admin-command-palette-planned.md) self `pass`
- 激活审视：[VRev-092](../../vision/reviews/VRev-092-vp036-admin-command-palette-activation.md) self `pass`
- Vision open required：0；V-F123 recommended 不阻断开区，但须在 R1 矩阵承接

## 纲领阶段

| 阶段 | 目的 | 状态 |
|------|------|------|
| R1 | 页面/导航/动作分母、Profile/权限语义、排序/去重、快捷键与实体搜索排除冻结 | completed；矩阵与 item-level oracle 已落盘，A-001 self `pass` + A-002 grok `conditional` 已由 A-003 响应，I-036-001～003 verified |
| R2 | `SearchableItem` 契约与跨模块 provider 聚合 | completed；provider v1、Manifest projection、冲突/匹配/cap 与测试已落盘，A-004 self `pass` |
| R3 | Admin Shell Command Palette、键盘可访问、i18n/theme、直接路由与分组联动 | completed；A-005 self `pass` + A-006 grok independent `pass`，A-007 已响应，mvp/admin SQLite/Postgres smoke 已通过 |
| R4 | Profile×权限×路由回归、证据矩阵、边界复核与关门 | completed；E-005/A-008 self `pass` + A-009 grok independent `pass`，A-010 闭合 2 条 recommended，四 Profile matrix 与 mvp/admin SQLite/Postgres smoke 已通过；**2026-09-14 用户确认 Root `done`** |

纲领阶段按 R1 → R2 → R3 → R4 串行推进；同一阶段内的细粒度子目标须在 R1 边界冻结后按证据与并行价值创建。工作区建立本身不代表任何实现阶段完成。

## 固定共享资料引用

> `shared-materials/index.json` 只能提供候选路径与摘要。本区 `shared_materials_catalog: none`，当前无可作为事实依据的共享资料引用。

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | — | — | none |

## 备注

本工作区只保存 VP-036 的实现层目标、决策、执行事实与 Goal 审计；不得将 Vision Review 或 VP 状态复制成第二套目标状态源。愿景组合编排仍以 `docs/vision/` 为准；本区 `goal-tree.md` 与目标五件套是实现层状态真相源。
