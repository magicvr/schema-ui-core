---
id: workspace-040-timestamptz-persistence-contract
title: DB 时间列 timestamptz 持久化合同工作区
status: active
root_goal: GOAL-001-timestamptz-persistence-contract
canonical_scope: docs/workspaces/workspace-040-timestamptz-persistence-contract/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-040-timestamptz-persistence-contract
primary_plan: VP-040-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
parent: null
---

# 工作区上下文 · DB 时间列 timestamptz 持久化合同

本工作区是 [VP-040-timestamptz-persistence-contract](../../vision/plans/VP-040-timestamptz-persistence-contract.md) 的唯一 `delivery` workspace，承接架构分支 C1：Store 时间列的 UTC 绝对时刻语义、PostgreSQL 生产权威物理类型、SQLite 合同平等物理类型、双方言迁移台账与读写编解码。

- VP-040 于 2026-09-20 经 `/vision` 激活为 **`active` v0.2.0**；激活就绪 self Review = [VRev-104](../../vision/reviews/VRev-104-vp040-timestamptz-persistence-contract-activation.md) `pass`，open required = 0。
- Root `[workspace-040-timestamptz-persistence-contract] GOAL-001-timestamptz-persistence-contract` 仍为 **`active · 0/3`**；R1 已渐进建立为 `[workspace-040-timestamptz-persistence-contract] GOAL-002-r1-contract-and-denominator-freeze`（`active · 0/4`），R2/R3 尚未创建。
- 用户已将 R1 合同改为 PostgreSQL `timestamptz(6)` + SQLite fixed-6 UTC RFC3339 `TEXT`；所有绝对时刻列纳入分母；双方言各自原地转换；sentinel `0` 按语义转为 `NULL`。
- 红线：不引入 ORM/第三库；不把 SQLite 假装成原生 `timestamptz`；不把驱动时间类型泄漏到 handler/模块公共契约；不消耗 Redis/MQ/多实例/A3 trigger；不重开 VP-013/VP-020。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-040-timestamptz-persistence-contract` | 本区唯一稳定标识 |
| Root Goal | `GOAL-001-timestamptz-persistence-contract` | `parent: null`；初始 `active · 0/3` |
| canonical 范围 | `docs/workspaces/workspace-040-timestamptz-persistence-contract/` | 本区唯一目标状态范围 |
| 共享资料目录 | `none` | 当前无固定共享资料；不得声明共享资料引用 |
| 愿景角色 | `delivery` | VP-040 唯一 delivery workspace；不改变 Charter primary workspace |
| 规划对齐 | `primary_plan` = `VP-040-timestamptz-persistence-contract` | `plan_refs` 必填且精确绑定 |

## 愿景对齐

- Charter：`schema-ui-core-admin-foundation@0.4.0`
- VP：[VP-040-timestamptz-persistence-contract](../../vision/plans/VP-040-timestamptz-persistence-contract.md)（`active` · v0.2.0）
- 计划阶段审视：[VRev-101](../../vision/reviews/VRev-101-vp039-vp040-planned.md) self `pass`
- 激活审视：[VRev-104](../../vision/reviews/VRev-104-vp040-timestamptz-persistence-contract-activation.md) self `pass`
- Vision open required：0；R1/R2/R3 信息门禁仍以 Root 与 VP-040 为准

## 纲领阶段

| 阶段 | 目的 | 状态 |
|------|------|------|
| R1 | 合同与分母冻结：方言物理类型、列清单、零值/NULL、备份 residual | **active**（GOAL-002 · 1/4；C1 completed） |
| R2 | 双方言迁移 + Store 编解码，沿用不可变 checksum 台账 | pending |
| R3 | 读写/时区回归 + 备份有界核对 + 证据与关门 | pending |

纲领阶段按 R1 → R2 → R3 串行；新建阶段子目标前必须先在 Root 决策/路线图中冻结边界与信息门禁。

## 固定共享资料引用

> `shared_materials_catalog: none`；当前无可作为事实依据的共享资料引用。

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | — | — | none |

## 备注

本工作区只保存 VP-040 的实现层目标、决策、执行事实与 Goal 审计；不得把 Vision Review 或 VP 状态复制成第二套目标状态源。愿景组合编排仍以 `docs/vision/` 为准；本区 `goal-tree.md` 与目标五件套是实现层状态真相源。
