---
id: workspace-014-object-storage
title: 对象存储适配器工作区
status: active
root_goal: GOAL-001-object-storage
canonical_scope: docs/workspaces/workspace-014-object-storage/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-014-object-storage
primary_plan: VP-014-object-storage
created: 2026-08-21
updated: 2026-08-21
version: 0.1.0
parent: null
---

# 工作区上下文 · 对象存储适配器

本工作区是 [VP-014-object-storage](../../vision/plans/VP-014-object-storage.md)（**`active`**，架构 A2）的唯一 lead delivery workspace。

- **Root** `GOAL-001-object-storage`：纲领 R1～R5；本回合只 scaffold，不创建子目标、不改代码。
- 不改变 Charter `primary_workspace`（仍为 workspace-001）。
- 不重开 `workspace-013-store-dialects`。
- 不承接 Admin 功能或业务域；不重开 VP-012 / VP-013。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-014-object-storage` | 与本区目标及资料引用的 `workspace_id` 一致 |
| Root Goal | `GOAL-001-object-storage` | `parent: null`；交付目标 |
| canonical 范围 | `docs/workspaces/workspace-014-object-storage/` | 本区唯一目标状态范围 |
| 共享资料目录 | `none` | 暂无固定共享资料 |
| 愿景角色 | `delivery` | VP-014 lead；不改变 Charter primary workspace |
| 规划对齐 | `primary_plan` = `VP-014-object-storage`（`active`） | 架构 A2 对象存储适配器 |

## 愿景对齐

Charter：`schema-ui-core-admin-foundation@0.2.0`。  
VP-014：内核对象存储端口；S3 兼容为生产验收权威；本地盘内嵌默认且合同平等。  
与 VP-009/VP-010 正交：安全/符合性 gap 不进本区。

## 纲领阶段（Root 路线图指针）

| 阶段 | 内容 | 状态 |
|------|------|------|
| R1 | 内核对象存储端口与配置面冻结 | 待立项 |
| R2 | S3 兼容接入（驱动、凭证、`readyz`） | 待立项 |
| R3 | 三类落盘收口（avatars / brand-assets / uploads） | 待立项 |
| R4 | 公共面去本地路径 / `os.File` | 待立项 |
| R5 | 双路径证据（本地盘默认 + S3 生产验收） | 待立项 |

## 固定共享资料引用

> `shared-materials/index.json` 只能提供候选路径与摘要。缺完整引用字段的行无效。

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | — | — | — |
