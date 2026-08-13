---
id: workspace-010-design-implementation-conformance
title: 设计意图与实现符合性工作区
status: active
root_goal: GOAL-001-design-implementation-conformance
canonical_scope: docs/workspaces/workspace-010-design-implementation-conformance/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-010-design-implementation-conformance
primary_plan: VP-010-design-implementation-conformance
created: 2026-08-11
updated: 2026-08-12
version: 0.2.0
parent: null
---

# 工作区上下文 · 设计意图与实现符合性

本工作区是 [VP-010-design-implementation-conformance](../../vision/plans/VP-010-design-implementation-conformance.md)（`active` · **长期设计意图—实现符合性程序**）的唯一 lead delivery workspace。

- **Root** 为长期程序容器（默认 `active`）。  
- **子目标** 为有界符合性审视/整改波次（可 `done`）。  
- 不因单波完成而关闭本区或 VP；不改变 Charter `primary_workspace`。  
- 与 [workspace-009-production-hardening](../workspace-009-production-hardening/workspace.md) **正交**：009 = 安全与健壮性；本区 = 架构/产品意图与 as-built 对齐。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-010-design-implementation-conformance` | 与本区目标及资料引用的 `workspace_id` 一致 |
| Root Goal | `GOAL-001-design-implementation-conformance` | `parent: null`；长期容器 |
| canonical 范围 | `docs/workspaces/workspace-010-design-implementation-conformance/` | 本区唯一目标状态范围 |
| 共享资料目录 | `none` | 暂无固定共享资料 |
| 愿景角色 | `delivery` | VP-010 lead；不改变 Charter primary workspace |
| 规划对齐 | `primary_plan` = `VP-010-design-implementation-conformance` | 持续程序意图 |

## 愿景对齐

Charter：`schema-ui-core-admin-foundation@0.2.0`。  
VP-010 为设计意图—实现符合性持续程序；与 VP-008 `go` 消费有效性接口见该 VP。  
若本区波次改变 Profile 默认集 / 模块矩阵 / Manifest 装配语义，须按规则暂挂或重验证业务对 `go` 的消费。

## 波次（实现层指针）

| 波次 | 子目标 | status |
|------|--------|--------|
| W1 | GOAL-002-w1-examples-optional-module | **done**（6/6 · 2026-08-11 关门；go 已恢复） |
| W2 | GOAL-003-demo-profile | **done**（6/6 · 2026-08-11 关门；go 无影响不暂挂） |
| W3 | GOAL-004-w3-schema-host-protocol-conformance | **active**（2/6 · S1 候选目录 + S3 新协议到手；上游 v2.8.0 正式身份已纠偏重 pin（E-005），停止线解除） |

## 固定共享资料引用

> `shared-materials/index.json` 只能提供候选路径与摘要。缺完整引用字段的行无效。

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | — | — | — |
