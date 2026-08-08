---
id: workspace-006-design-system-and-ui-experience
title: 设计系统与 Schema 驱动 UI/UX 体验工作区
status: done
root_goal: GOAL-001-design-system-and-ui-experience
canonical_scope: docs/workspace-006-design-system-and-ui-experience/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-005-design-system-and-ui-experience
primary_plan: VP-005-design-system-and-ui-experience
created: 2026-08-09
updated: 2026-08-09
version: 0.2.0
parent: null
---

# 工作区上下文 · 设计系统与 Schema 驱动 UI/UX

本工作区承接 [VP-005 · 现代设计系统与 Schema 驱动 UI/UX 体验产品化](../vision/plans/VP-005-design-system-and-ui-experience.md)（**`active`**，2026-08-09 用户确认激活 + 本轮 `/govern` 开区）的实现层治理。唯一目标状态只保存在本目录的 `goal-tree.md` 与平铺 Goal 五件套中；本文件不维护第二套 progress 或审计台账。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-006-design-system-and-ui-experience` | 用户 2026-08-09 经 `/govern` **书面确认** slug。 |
| Root Goal | `GOAL-001-design-system-and-ui-experience` | 本区唯一 `parent: null` Root；`status: done`（2026-08-09 用户书面确认关门，见 D-005）；纲领 S1–S5（`progress: 5/5`）。 |
| canonical 范围 | `docs/workspace-006-design-system-and-ui-experience/` | 当前工作区唯一的目标状态范围。 |
| 共享资料目录 | `none` | 未声明固定共享资料引用；不得把候选资料作为事实或证据。 |
| 愿景角色 | `delivery` | 不改变 Charter 的 primary 工作区 `workspace-001-mvp-admin-foundation`。 |
| 规划对齐 | `VP-005-design-system-and-ui-experience` | `plan_refs` 与 `primary_plan` 均指向 **active** VP-005。 |

## 愿景对齐

VP-005 的 `vision_ref` 精确匹配现行 Charter `schema-ui-core-admin-foundation@0.2.0`。本工作区为该 VP 的唯一 lead / delivery。视觉升级范围钉死在 [I-PROTO-FULL-001](../workspace-005-full-protocol-contract-v2-7-0/GOAL-001-full-protocol-contract-v2-7-0/attachments/I-PROTO-FULL-001-coverage-v2-7-0.md) include 面与 VP-005 真实 type 表（VRev-011 F-V018 fixed）；**禁止**借本区扩张协议 disposition。

**禁止**在 closed [workspace-003](../workspace-003-modular-admin-architecture/workspace.md) / [workspace-004](../workspace-004-module-contribution-readiness/workspace.md) / [workspace-005](../workspace-005-full-protocol-contract-v2-7-0/workspace.md) 吸收本意图。

建区本身 **不**勾选任何 S1–S5 检查点，也 **不**构成 Design Token / Shell / Renderer 视觉已产品化声明。

**视觉方向（D-004）**：Stitch 定稿已冻结为实施输入（I-005 closed）；仓库摘要见 Root `attachments/visual-direction-stitch-summary.md`；本地截图在 `raw/stitch-vp005-visual-refs/`（gitignore）。定稿 **不**勾选 S1–S5；开放 required finding 仍 **F-002**（S1 完成门禁）。

## 固定共享资料引用

`shared_materials_catalog: none`。当前没有固定共享资料引用，也不得以其他工作区的资料或记录替代本区的事实、信息门禁或审计证据。协议 pin 与覆盖表以 Charter + inventory + I-PROTO-FULL-001 为只读权威输入。

## 串行阶段说明

Root 的 S1–S5 纲领阶段原则上串行；同一阶段内可在信息门禁与计划就绪后建立并行子目标。跨区纲领阶段写在 `docs/vision/roadmap.md` 与 VP-005。
