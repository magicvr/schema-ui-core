---
id: workspace-013-store-dialects
title: Store 双方言工作区
status: active
root_goal: GOAL-001-store-dialects
canonical_scope: docs/workspaces/workspace-013-store-dialects/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-013-store-dialects
primary_plan: VP-013-store-dialects
created: 2026-08-20
updated: 2026-08-21
version: 0.2.0
parent: null
---

# 工作区上下文 · Store 双方言

本工作区是 [VP-013-store-dialects](../../vision/plans/VP-013-store-dialects.md)（**`closed`**，2026-08-21 有界关门 · 架构 A1）的唯一 lead delivery workspace。历史绑定保留，默认不接新区。

- **Root** 已为 `done`（2026-08-20；R1～R5 5/5；A-001 independent / A-002 响应，2026-08-21 补齐闭门依据）。
- **子目标**按 Root 纲领阶段 R1～R5 全部 `done`。
- 不改变 Charter `primary_workspace`（仍为 workspace-001）。
- 不承接 Admin 功能或业务域；不重开 VP-012。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-013-store-dialects` | 与本区目标及资料引用的 `workspace_id` 一致 |
| Root Goal | `GOAL-001-store-dialects` | `parent: null`；交付目标 |
| canonical 范围 | `docs/workspaces/workspace-013-store-dialects/` | 本区唯一目标状态范围 |
| 共享资料目录 | `none` | 暂无固定共享资料 |
| 愿景角色 | `delivery` | VP-013 lead；不改变 Charter primary workspace |
| 规划对齐 | `primary_plan` = `VP-013-store-dialects`（`closed`） | 架构 A1 Store 双方言；历史绑定 |

## 愿景对齐

Charter：`schema-ui-core-admin-foundation@0.2.0`。  
VP-013：内核 Store 为持久化端口；PostgreSQL 生产权威；SQLite 内嵌默认且合同平等；无 ORM。  
与 VP-009/VP-010 正交：安全/符合性 gap 不进本区。

## 纲领阶段（Root 路线图指针）

| 阶段 | 内容 | 状态 |
|------|------|------|
| R1 | 内核持久化端口与配置面冻结（Tx ≠ `*sql.Tx`；SQLite 默认 / PG DSN；合同 v1.4） | ✅ GOAL-002 |
| R2 | PostgreSQL 接入（驱动、连接池、`readyz`） | ✅ GOAL-003 |
| R3 | 现有 compiled 台账双方言对写 + checksum | ✅ GOAL-004 |
| R4 | 模块仓库公共面收口（去掉 `*sql.Tx`） | ✅ GOAL-005 |
| R5 | 双路径证据（SQLite 默认 + PG 生产验收） | ✅ GOAL-006 |

## 固定共享资料引用

> `shared-materials/index.json` 只能提供候选路径与摘要。缺完整引用字段的行无效。

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | — | — | — |
