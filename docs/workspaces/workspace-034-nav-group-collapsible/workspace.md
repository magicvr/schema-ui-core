---
id: workspace-034-nav-group-collapsible
title: Admin 导航分组折叠体验工作区
status: active
root_goal: GOAL-001-nav-group-collapsible
canonical_scope: docs/workspaces/workspace-034-nav-group-collapsible/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-034-nav-group-collapsible
primary_plan: VP-034-nav-group-collapsible
created: 2026-09-07
updated: 2026-09-07
version: 0.1.0
parent: null
---

# 工作区上下文 · Admin 导航分组折叠体验

本工作区是 [VP-034-nav-group-collapsible](../../vision/plans/VP-034-nav-group-collapsible.md)（`active` · v0.3.0）的唯一 lead delivery workspace，承接 Admin Shell 左侧导航分组、模块注册协议、既有导航迁移与回归验证。

- 本区不属于 workspace-010；VP-010 继续承担长期设计意图—实现符合性程序，本区是 Admin 功能分支的有界交付。
- 本区必须覆盖当前已注册 sidebar 导航，不以“既有模块”作为排除理由。
- `menu_dashboard` 保持顶层主入口单例；top/user slot 导航保持原 slot 语义并做兼容回归。
- 组与模块解耦：模块可将多个导航条目注册到同一组，也可在明确理由下保留顶层单例或不分组。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-034-nav-group-collapsible` | 与本区目标及资料引用的 `workspace_id` 一致 |
| Root Goal | `GOAL-001-nav-group-collapsible` | `parent: null`；长期 VP-034 有界交付 Root |
| canonical 范围 | `docs/workspaces/workspace-034-nav-group-collapsible/` | 本区唯一目标状态范围 |
| 共享资料目录 | `none` | 本区暂无固定共享资料 |
| 愿景角色 | `delivery` | VP-034 lead delivery；不改变 Charter primary workspace |
| 规划对齐 | `primary_plan` = `VP-034-nav-group-collapsible` | `plan_refs` 必填且已精确绑定 |

## 愿景对齐

- Charter：`schema-ui-core-admin-foundation@0.4.0`
- VP：`VP-034-nav-group-collapsible`（`active` · v0.3.0）
- 激活审视：[VRev-083](../../vision/reviews/VRev-083-vp034-nav-group-collapsible-activation.md) self `pass`
- Scope 修正审视：[VRev-084](../../vision/reviews/VRev-084-vp034-existing-navigation-scope-correction.md) self `pass`
- Vision open required：0；V-F121 为 inherited recommended，不阻断

## 初始分组基线（R1 待冻结）

| group key | 初始标题 | sidebar 节点 |
|-----------|----------|--------------|
| `identity-access` | Identity & access | `menu_users`、`menu_roles`、`menu_data_permission` |
| `content-data` | Content & data | `menu_files`、`menu_dictionary` |
| `operations` | Operations | `menu_activity`、`menu_monitoring`、`menu_scheduled_tasks`、`menu_recycle_bin` |
| `communications` | Communications | `menu_mail`、`menu_mail_outbox`、`menu_telegram` |
| `commerce` | Commerce | `menu_wallet`、`menu_wallet_vouchers`、`menu_digitaloffer_offers`、`menu_digitaloffer_entitlements`、`menu_digitaloffer_purchases` |
| 顶层单例 | Dashboard | `menu_dashboard`（有意不包进单项组） |

`dev.examples` 的现有 `Examples` 组保留并纳入 demo 回归。`menu_settings`、`menu_account`、`menu_wallet_self` 和通知面保留在 user slot，不搬到 sidebar。

## 纲领阶段

| 阶段 | 目的 | 状态 |
|------|------|------|
| R1 | 现有导航清单、Profile/slot 矩阵、分组信息架构与组 key/顺序冻结 | completed |
| R2 | 模块 NavigationContribution 与 Manifest 聚合契约：可选 `group`、跨模块共组、向后兼容 | completed |
| R3 | Shell 左侧分组折叠/展开、键盘可访问性、激活态自动展开与状态保持 | planned |
| R4 | 当前已注册 sidebar 导航全量迁移，覆盖默认/optional/custom/demo 组合并保持 top/user slot | planned |
| R5 | 证据矩阵、全量回归、self/independent 审计与 VP 关门准备 | planned |

纲领阶段按 R1 → R2 → R3 → R4 → R5 串行推进；同一阶段内的细粒度子目标需在 R1 冻结后按证据与并行价值创建。

## 固定共享资料引用

> `shared-materials/index.json` 只能提供候选路径与摘要。本区 `shared_materials_catalog: none`，当前无可作为事实依据的共享资料引用。

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | — | — | none |

## 备注

本工作区只保存 VP-034 的实现层目标、决策、执行事实与 Goal 审计；不得将 Vision Review 或 VP 状态复制成第二套目标状态源。