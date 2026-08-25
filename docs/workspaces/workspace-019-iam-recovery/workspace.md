---
id: workspace-019-iam-recovery
title: IAM 工作区（密码策略 / 邀请入职 / 自助恢复状态机）
status: active
root_goal: GOAL-001-iam-recovery
canonical_scope: docs/workspaces/workspace-019-iam-recovery/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-019-iam-recovery
primary_plan: VP-019-iam-recovery
created: 2026-08-25
updated: 2026-08-25
version: 0.1.0
parent: null
---

# 工作区上下文 · IAM（密码策略 / 邀请入职 / 自助恢复状态机）

本工作区是 [VP-019-iam-recovery](../../vision/plans/VP-019-iam-recovery.md)（**`active`** · 2026-08-25 激活；VRev-043 independent `pass`）的唯一 lead delivery workspace。

- **Root** `GOAL-001-iam-recovery`：**`active`** · 0/4（2026-08-25 开区；R1 合同冻结待启动）。
- 不改变 Charter `primary_workspace`（仍为 workspace-001）。
- 硬前置已完成：VP-017（运输，渠道模型 v0.5.0 closed）+ VP-018（邮箱身份 v1.0.0 closed）。
- 边界：不承接 SMS、模板中心、多邮箱、组织权限、OIDC、业务域；不改 Profile 默认集 / 模块矩阵 / Manifest 装配语义。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-019-iam-recovery` | 与本区目标及资料引用的 `workspace_id` 一致 |
| Root Goal | `GOAL-001-iam-recovery` | `parent: null`；**active**（纲领 R1～R4） |
| canonical 范围 | `docs/workspaces/workspace-019-iam-recovery/` | 本区唯一目标状态范围 |
| 共享资料目录 | `none` | 暂无固定共享资料 |
| 愿景角色 | `delivery` | VP-019 lead；不改变 Charter primary workspace |
| 规划对齐 | `primary_plan` = `VP-019-iam-recovery`（`active`） | 2026-08-25 激活/开区（VR-047；VRev-043 pass） |

## 愿景对齐

Charter：`schema-ui-core-admin-foundation@0.2.0`。  
VP-019：IAM 三件（密码策略 / 邀请入职 / 自助恢复状态机）；消费 VP-017 `MailSender` 与 VP-018 邮箱身份；证明依据 = 已绑定且已校验邮箱（2026-08-22 产品事实）。  
与 VP-009/VP-010 正交：安全/符合性 gap（枚举、重放、open-relay、邀请滥用）不进本区。  
VP-008 `go` 消费有效性：Admin 类 freshness **PASS**（`092bf37` → `66f5fd1f`，2026-08-25；VRev-043）。

## 纲领阶段（Root 路线图指针）

| 阶段 | 内容 | 状态 |
|------|------|------|
| R1 | 合同冻结：证明形态（I-019-001/6 位码候选）、TTL/冷却（I-019-002 → R2 前）、策略参数与 Profile 边界（I-019-003/007）、邀请形态（I-019-004/005）、MFA 与自助恢复（I-019-009） | **待启动**（Root 00-meta I-001～009） |
| R2 | 自助恢复全链：登录页发起 → 投递 → 校验 → 设新码 → 登录；会话/令牌语义 | 依赖 R1 |
| R3 | 密码策略（配置面 + 强制面）+ 邀请入职 | 依赖 R2 |
| R4 | 证据与关门：端到端经现行渠道取信、策略强制可核对、邀请全链、无越界 | 依赖 R3 |

## 固定共享资料引用

> `shared-materials/index.json` 只能提供候选路径与摘要。缺完整引用字段的行无效。

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | — | none | — |