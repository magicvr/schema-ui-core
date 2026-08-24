---
id: workspace-018-account-email-identity
title: 账号邮箱身份工作区
status: done
root_goal: GOAL-001-account-email-identity
canonical_scope: docs/workspaces/workspace-018-account-email-identity/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-018-account-email-identity
primary_plan: VP-018-account-email-identity
created: 2026-08-24
updated: 2026-08-24
version: 1.0.0
parent: null
---

# 工作区上下文 · 账号邮箱身份

本工作区是 [VP-018-account-email-identity](../../vision/plans/VP-018-account-email-identity.md)（**`closed` · 同日关门**）的唯一 lead delivery workspace。

- **Root** `GOAL-001-account-email-identity`：**`done`** · 4/4（2026-08-24 关门：A-001 self pass + A-002 independent conditional→F-001 fixed 后归零）。R1～R4 = GOAL-002～005 全 done。
- 不改变 Charter `primary_workspace`（仍为 workspace-001）。
- 不承接自助恢复状态机、邀请、密码策略、SMS、模板产品或业务域。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-018-account-email-identity` | 与本区目标及资料引用的 `workspace_id` 一致 |
| Root Goal | `GOAL-001-account-email-identity` | `parent: null`；**done**（2026-08-24 关门） |
| canonical 范围 | `docs/workspaces/workspace-018-account-email-identity/` | 本区唯一目标状态范围 |
| 共享资料目录 | `none` | 暂无固定共享资料 |
| 愿景角色 | `delivery` | VP-018 lead；不改变 Charter primary workspace |
| 规划对齐 | `primary_plan` = `VP-018-account-email-identity`（`closed` · 同日关门） | 解冻与关门均在 2026-08-24（VP-017 v0.5.0 再关门后恢复，当日交付完毕） |

## 愿景对齐

Charter：`schema-ui-core-admin-foundation@0.2.0`。  
VP-018：账号邮箱身份面；消费 `MailSender`；2026-08-24 解冻（VP-017 按现行渠道分母再关门）。  
与 VP-009/VP-010 正交：安全/符合性 gap 不进本区。

## 纲领阶段（Root 路线图指针）

| 阶段 | 内容 | 状态 |
|------|------|------|
| R1 | 身份合同冻结（唯一性细则、校验形态、可空、换绑） | **已完成**（GOAL-002 done · A-001 pass） |
| R2 | 双方言 schema + 唯一性约束 | **已完成**（GOAL-003 done · 迁移 0054 · A-001 independent pass） |
| R3 | 绑定/校验消费 `MailSender` | **已完成**（GOAL-004 done · 迁移 0055 · required 归零后关门） |
| R4 | 证据 | **已完成**（GOAL-005 done · 端到端经真实 mock 渠道适配器 · self pass）——Root 已关门（A-002 independent 归零） |

## 固定共享资料引用

> `shared-materials/index.json` 只能提供候选路径与摘要。缺完整引用字段的行无效。

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | — | none | — |
