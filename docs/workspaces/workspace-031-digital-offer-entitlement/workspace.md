---
id: workspace-031-digital-offer-entitlement
title: 数字 Offer 与权益工作区
status: active
root_goal: GOAL-001-digital-offer-entitlement
canonical_scope: docs/workspaces/workspace-031-digital-offer-entitlement/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-031-digital-offer-entitlement
primary_plan: VP-031-digital-offer-entitlement
created: 2026-09-05
updated: 2026-09-05
version: 0.1.0
parent: null
---

# 工作区上下文 · 数字 Offer 与权益

本工作区是 [VP-031-digital-offer-entitlement](../../vision/plans/VP-031-digital-offer-entitlement.md) 的唯一 lead delivery workspace，承接数字 Offer、薄购买凭证与本域权益校验。它是本仓库首个业务域 VP 的实现上下文，不扩展为 Catalog/SKU/税/库存/物流订单，不进入默认 Profile，不解禁通用 Entitlement/Approval Gate。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-031-digital-offer-entitlement` | 与 canonical 路径一致 |
| Root Goal | `GOAL-001-digital-offer-entitlement` | `parent: null`；**active · 0/4** |
| canonical 范围 | `docs/workspaces/workspace-031-digital-offer-entitlement/` | 本区唯一目标状态范围 |
| 共享资料目录 | `none` | 暂无固定共享资料；不得引用未固定材料 |
| 愿景角色 | `delivery` | VP-031 唯一 lead；不改变 Charter primary workspace |
| 规划对齐 | `primary_plan = VP-031-digital-offer-entitlement` | `vision_ref = schema-ui-core-admin-foundation@0.4.0` |

## 纲领阶段

| 阶段 | 内容 | 状态 |
|------|------|------|
| R1 | 合同冻结：Offer 字段、购买状态机、权益形态、命令清单、事务与限流边界 | 待开始；I-031-001～005 尚未处理；V-F119 待承接 |
| R2 | Offer CRUD、购买凭证、钱包 `freeze → deduct_frozen` / 失败 `unfreeze` | 待开始；依赖 R1 |
| R3 | 权益校验、过期/耗尽语义、可选 Telegram Register | 待开始；依赖 R2 |
| R4 | 证据矩阵、边界核账、审计闭合与关门 | 待开始；依赖 R1～R3 |

## 激活边界

- H-002：用户于 2026-09-05 书面确认采用同进程模块。
- VP-029 为硬前置且已关闭；VP-030 为软前置且已关闭。
- RT-Q03/RT-Q05：VRev-080 已评估为本波不需要 Redis；多实例或跨实例共享状态出现时复审。
- 建区不等于 R1 已完成，不等于业务实现、审计或 `go`。
