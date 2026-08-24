---
id: workspace-018-account-email-identity
title: 账号邮箱身份工作区
status: active
root_goal: GOAL-001-account-email-identity
canonical_scope: docs/workspaces/workspace-018-account-email-identity/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-018-account-email-identity
primary_plan: VP-018-account-email-identity
created: 2026-08-24
updated: 2026-08-24
version: 0.3.0
parent: null
---

# 工作区上下文 · 账号邮箱身份

本工作区是 [VP-018-account-email-identity](../../vision/plans/VP-018-account-email-identity.md)（**`active` · 已解冻**）的唯一 lead delivery workspace。

- **Root** `GOAL-001-account-email-identity`：**`active`** · 1/4。2026-08-24 解冻（D-003）：VP-017 已按现行渠道分母再次 `closed`（v0.5.0），用户确认恢复推进；R1 身份合同已冻结（GOAL-002 D-001 · A-001 self pass）。
- 不改变 Charter `primary_workspace`（仍为 workspace-001）。
- 不承接自助恢复状态机、邀请、密码策略、SMS、模板产品或业务域。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-018-account-email-identity` | 与本区目标及资料引用的 `workspace_id` 一致 |
| Root Goal | `GOAL-001-account-email-identity` | `parent: null`；**active**（D-003 解冻） |
| canonical 范围 | `docs/workspaces/workspace-018-account-email-identity/` | 本区唯一目标状态范围 |
| 共享资料目录 | `none` | 暂无固定共享资料 |
| 愿景角色 | `delivery` | VP-018 lead；不改变 Charter primary workspace |
| 规划对齐 | `primary_plan` = `VP-018-account-email-identity`（`active` · 已解冻） | 解冻已于 2026-08-24 生效（VP-017 v0.5.0 现行分母再关门） |

## 愿景对齐

Charter：`schema-ui-core-admin-foundation@0.2.0`。  
VP-018：账号邮箱身份面；消费 `MailSender`；2026-08-24 解冻（VP-017 按现行渠道分母再关门）。  
与 VP-009/VP-010 正交：安全/符合性 gap 不进本区。

## 纲领阶段（Root 路线图指针）

| 阶段 | 内容 | 状态 |
|------|------|------|
| R1 | 身份合同冻结（唯一性细则、校验形态、可空、换绑） | **已完成**（GOAL-002 done · A-001 pass） |
| R2 | 双方言 schema + 唯一性约束 | **待启动** |
| R3 | 绑定/校验消费 `MailSender` | **待启动** |
| R4 | 证据 | **待启动** |

## 固定共享资料引用

> `shared-materials/index.json` 只能提供候选路径与摘要。缺完整引用字段的行无效。

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | — | none | — |
