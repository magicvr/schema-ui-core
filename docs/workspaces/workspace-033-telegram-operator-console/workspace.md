---
id: workspace-033-telegram-operator-console
title: Telegram Bot 人工控制台工作区
status: active
root_goal: GOAL-001-telegram-operator-console
canonical_scope: docs/workspaces/workspace-033-telegram-operator-console/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-033-telegram-operator-console
primary_plan: VP-033-telegram-operator-console
created: 2026-09-04
updated: 2026-09-04
version: 0.1.0
parent: null
---

# 工作区上下文 · Telegram Bot 人工控制台

本工作区是 [VP-033-telegram-operator-console](../../vision/plans/VP-033-telegram-operator-console.md)（`active` v0.2.0）的唯一 lead delivery workspace，消费 VP-030 已交付的 Telegram runtime，在 Admin 功能分支交付连接状态、互斥 webhook/polling、业务占用位与未绑定人工文本控制台。

- **Root**：`GOAL-001-telegram-operator-console`，`active · 0/4`。
- **激活依据**：[VRev-075](../../vision/reviews/VRev-075-vp033-telegram-operator-console-activation.md) self `pass`，open required = 0；Admin freshness `42036a3c` → `dd1edade` PASS。
- **边界**：不重开 VP-030；不进入业务域；不进入 `mvp`/`admin` 默认 Profile；不解除 SSE/WebSocket、多 bot 或多实例 polling 门闩。
- **既存 residual**：workspace-030 `R-009` 仅保留其原有适用范围；若本区改动密钥存储或生产隔离，必须回流 P-004，不自动继承为本区风险接受。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-033-telegram-operator-console` | 与 canonical 路径一致 |
| Root Goal | `GOAL-001-telegram-operator-console` | `parent: null`；active 0/4 |
| canonical 范围 | `docs/workspaces/workspace-033-telegram-operator-console/` | 本区唯一目标状态范围 |
| 共享资料目录 | `none` | 暂无固定共享资料 |
| 愿景角色 | `delivery` | VP-033 唯一 lead；不改变 Charter primary workspace |
| 规划对齐 | `primary_plan = VP-033-telegram-operator-console` | `vision_ref = schema-ui-core-admin-foundation@0.4.0` |

## 纲领阶段

| 阶段 | 内容 | 状态 |
|------|------|------|
| R1 | 合同冻结：入站模式、占用位、heartbeat、发言权、显式公网 URL | 进行中（GOAL-002-r1-contract-freeze） |
| R2 | 连接、互斥热切换、占用位与设置页 | 待开始 |
| R3 | 会话落盘与未绑定人工 IM | 待开始 |
| R4 | 证据矩阵、审计闭合与关门 | 待开始 |

## 固定共享资料引用

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | — | none | — |
