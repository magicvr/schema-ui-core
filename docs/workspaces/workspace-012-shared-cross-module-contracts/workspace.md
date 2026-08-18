---
id: workspace-012-shared-cross-module-contracts
title: 共享横切契约与平台基架工作区
status: active
root_goal: GOAL-001-shared-cross-module-contracts
canonical_scope: docs/workspaces/workspace-012-shared-cross-module-contracts/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-012-shared-cross-module-contracts
primary_plan: VP-012-shared-cross-module-contracts
created: 2026-08-18
updated: 2026-08-18
version: 0.1.1
parent: null
---

# 工作区上下文 · 共享横切契约与平台基架

本工作区是 [VP-012-shared-cross-module-contracts](../../vision/plans/VP-012-shared-cross-module-contracts.md)（`active`）的唯一 lead delivery workspace。

- **Root** 为交付目标（默认 `active`），首波 = **横切契约波**。
- **子目标** = 各契约的独立交付波次。
- 不改变 Charter `primary_workspace`（仍为 workspace-001）。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-012-shared-cross-module-contracts` | 与本区目标及资料引用的 `workspace_id` 一致 |
| Root Goal | `GOAL-001-shared-cross-module-contracts` | `parent: null`；交付目标 |
| canonical 范围 | `docs/workspaces/workspace-012-shared-cross-module-contracts/` | 本区唯一目标状态范围 |
| 共享资料目录 | `none` | 暂无固定共享资料 |
| 愿景角色 | `delivery` | VP-012 lead；不改变 Charter primary workspace |
| 规划对齐 | `primary_plan` = `VP-012-shared-cross-module-contracts` | 共享横切契约与平台基架 |

## 愿景对齐

Charter：`schema-ui-core-admin-foundation@0.2.0`。  
VP-012：共享横切契约与平台基架（correlation / 审计模型 / 并发幂等 / 异步 Job / maintenance 门控 / API Token）；不承载 Tier D 业务域。  
与 VP-009/VP-010 正交：安全威胁面归 009、符合性 gap 归 010；本区交付可被两者消费的契约。

## 纲领阶段（Root 路线图指针）

| 阶段 | 内容 | 状态 |
|------|------|------|
| R1 | correlation / request-id / 错误恢复契约 | ✅ 已完成（GOAL-002-r1-correlation-error-contract） |
| R2 | 审计事件模型增强（结构化 diff / 脱敏 / correlation） | ✅ 已完成（GOAL-003-r2-audit-event-model；A-006 pass / A-007 closed） |
| R3 | 乐观并发 + 幂等契约 | ✅ 已完成（GOAL-004-r3-optimistic-concurrency-idempotency；A-004 pass / A-005 closed） |
| R4 | 异步 Job / 长操作契约 | 未开始 |
| R5 | maintenance / degraded / read-only 门控 | 未开始 |
| R6 | API Token / Service Credential | 未开始 |

## 固定共享资料引用

> `shared-materials/index.json` 只能提供候选路径与摘要。缺完整引用字段的行无效。

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | — | — | — |
