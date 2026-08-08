---
id: workspace-005-full-protocol-contract-v2-7-0
title: schema-ui-docs@v2.7.0 整份契约可验证兼容工作区
status: active
root_goal: GOAL-001-full-protocol-contract-v2-7-0
canonical_scope: docs/workspace-005-full-protocol-contract-v2-7-0/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-006-full-protocol-contract-v2-7-0
primary_plan: VP-006-full-protocol-contract-v2-7-0
created: 2026-08-08
updated: 2026-08-08
version: 0.1.0
parent: null
---

# 工作区上下文 · 整份 v2.7.0 契约可验证兼容

本工作区承接 [VP-006 · schema-ui-docs@v2.7.0 整份契约可验证兼容](../vision/plans/VP-006-full-protocol-contract-v2-7-0.md) 的实现层治理（**已 closed**，2026-08-08 用户书面确认；历史绑定保留）。唯一目标状态只保存在本目录的 `goal-tree.md` 与平铺 Goal 五件套中；本文件不维护第二套 progress 或审计台账。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-005-full-protocol-contract-v2-7-0` | 当前 delivery 工作区的稳定标识（用户 2026-08-08 经 `/vision` + `/govern` 确认开区）。 |
| Root Goal | `GOAL-001-full-protocol-contract-v2-7-0` | 本区唯一 `parent: null` Root；**终态 `done / 6/6`**（VP-006 2026-08-08 用户书面确认关门）。 |
| canonical 范围 | `docs/workspace-005-full-protocol-contract-v2-7-0/` | 当前工作区唯一的目标状态范围。 |
| 共享资料目录 | `none` | 未声明固定共享资料引用；不得把候选资料作为事实或证据。 |
| 愿景角色 | `delivery` | 不改变 Charter 的 primary 工作区 `workspace-001-mvp-admin-foundation`。 |
| 规划对齐 | `VP-006-full-protocol-contract-v2-7-0` | `plan_refs` 与 `primary_plan` 均指向 **closed** VP-006（历史绑定保留）。 |

## 愿景对齐

VP-006 的 `vision_ref` 精确匹配现行 Charter `schema-ui-core-admin-foundation@0.2.0`。本工作区为该 VP 的唯一 lead / delivery。协议 pin 与全量清单见 [protocol-inventory-v2.7.0.md](../vision/protocol-inventory-v2.7.0.md)；历史 MVP 覆盖表 `I-PROTO-001 v0.1.3` 仅作升版起点与回归对照，**不是**本区退出上界。

**禁止**在 closed [workspace-003-modular-admin-architecture](../workspace-003-modular-admin-architecture/workspace.md) 或 [workspace-004-module-contribution-readiness](../workspace-004-module-contribution-readiness/workspace.md) 吸收本意图。  
**禁止**在本 VP closed 前启动 [VP-005](../vision/plans/VP-005-design-system-and-ui-experience.md) 视觉实施。

建区本身不勾选任何 S0–S5 检查点，也不构成「已完整支持 v2.7.0」声明。

## 固定共享资料引用

`shared_materials_catalog: none`。当前没有固定共享资料引用，也不得以其他工作区的资料或记录替代本区的事实、信息门禁或审计证据。上游协议以 Charter pin + inventory 为只读权威输入。

## 串行阶段说明

Root 的 S0–S5 纲领阶段原则上串行；同一阶段内可在信息门禁与计划就绪后建立并行子目标。跨区纲领阶段写在 `docs/vision/roadmap.md` 与 VP-006。
