---
id: workspace-017-outbound-mail
title: 出站邮件工作区
status: active
root_goal: GOAL-001-outbound-mail
canonical_scope: docs/workspaces/workspace-017-outbound-mail/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-017-outbound-mail
primary_plan: VP-017-outbound-mail
created: 2026-08-22
updated: 2026-08-22
version: 0.1.0
parent: null
---

# 工作区上下文 · 出站邮件

本工作区是 [VP-017-outbound-mail](../../vision/plans/VP-017-outbound-mail.md)（**`active`**，架构 A6）的唯一 lead delivery workspace。

- **Root** `GOAL-001-outbound-mail`：纲领 R1～R4（端口冻结 → SMTP 接入 → 默认 sink / 公共面 → 显式路径 + `readyz`）。
- 不改变 Charter `primary_workspace`（仍为 workspace-001）。
- 不重开 `workspace-016-key-rotation-and-backup`。
- 不承接账号 email、邀请、自助恢复、Admin 功能或业务域；不重开 VP-012～016。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-017-outbound-mail` | 与本区目标及资料引用的 `workspace_id` 一致 |
| Root Goal | `GOAL-001-outbound-mail` | `parent: null`；交付目标 |
| canonical 范围 | `docs/workspaces/workspace-017-outbound-mail/` | 本区唯一目标状态范围 |
| 共享资料目录 | `none` | 暂无固定共享资料 |
| 愿景角色 | `delivery` | VP-017 lead；不改变 Charter primary workspace |
| 规划对齐 | `primary_plan` = `VP-017-outbound-mail`（`active`） | 架构 A6 内核发送端口 + SMTP |

## 愿景对齐

Charter：`schema-ui-core-admin-foundation@0.2.0`。  
VP-017：内核出站邮件发送端口；SMTP 为生产实现；capture/log sink 为内嵌默认；重启生效。  
与 VP-009/VP-010 正交：安全/符合性 gap 不进本区。

## 纲领阶段（Root 路线图指针）

| 阶段 | 内容 | 状态 |
|------|------|------|
| R1 | 端口与发送合同冻结 | **已完成**（GOAL-002；D-002） |
| R2 | SMTP 接入与配置面 | **已完成**（GOAL-003；D-003，隐式 TLS 465 唯一路径） |
| R3 | 默认 sink + 公共面去 SMTP 客户端类型 | **已完成**（GOAL-004；D-004） |
| R4 | 显式路径证据 + `readyz` | **已完成**（GOAL-005；D-005） |

Root progress 4/4；关门审计进行中（03-audit.md）。

## 固定共享资料引用

> `shared-materials/index.json` 只能提供候选路径与摘要。缺完整引用字段的行无效。

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | — | — | — |
