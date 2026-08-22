---
id: workspace-016-key-rotation-and-backup
title: 密钥轮换与备份恢复工作区
status: active
root_goal: GOAL-001-key-rotation-and-backup
canonical_scope: docs/workspaces/workspace-016-key-rotation-and-backup/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-016-key-rotation-and-backup
primary_plan: VP-016-key-rotation-and-backup
created: 2026-08-22
updated: 2026-08-22
version: 0.2.0
parent: null
---

# 工作区上下文 · 密钥轮换与备份恢复

本工作区是 [VP-016-key-rotation-and-backup](../../vision/plans/VP-016-key-rotation-and-backup.md)（**`closed`**，2026-08-22 有界关门 · 架构 A5）的唯一 lead delivery workspace。历史绑定保留，默认不接新区。

- **Root** `GOAL-001-key-rotation-and-backup`：纲领 R1～R5 全部完成；Root 于 2026-08-22 **done 5/5**（GOAL-002～006 均 done）。
- 不改变 Charter `primary_workspace`（仍为 workspace-001）。
- 不重开 `workspace-015-observability`。
- 不承接 Admin 功能或业务域；不重开 VP-012 / VP-013 / VP-014 / VP-015。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-016-key-rotation-and-backup` | 与本区目标及资料引用的 `workspace_id` 一致 |
| Root Goal | `GOAL-001-key-rotation-and-backup` | `parent: null`；交付目标 |
| canonical 范围 | `docs/workspaces/workspace-016-key-rotation-and-backup/` | 本区唯一目标状态范围 |
| 共享资料目录 | `none` | 暂无固定共享资料 |
| 愿景角色 | `delivery` | VP-016 lead；不改变 Charter primary workspace |
| 规划对齐 | `primary_plan` = `VP-016-key-rotation-and-backup`（`closed`） | 架构 A5 JWT 轮换 + 轮换后恢复；历史绑定 |

## 愿景对齐

Charter：`schema-ui-core-admin-foundation@0.2.0`。  
VP-016：JWT current+previous 轮换合同；既有 SQLite `VACUUM INTO` 与 PG `pg_dump`/`pg_restore` 上的轮换后恢复；单密钥为内嵌默认；重启生效。  
与 VP-009/VP-010 正交：安全/符合性 gap 不进本区。

## 纲领阶段（Root 路线图指针）

| 阶段 | 内容 | 状态 |
|------|------|------|
| R1 | 轮换合同与配置面冻结 | **已完成**（D-002 + GOAL-002 done · A-001 self pass） |
| R2 | JWT 双密钥实现 | **已完成**（GOAL-003 done · self A-001 + independent A-002 双 pass · F-001/2/3 fixed） |
| R3 | 轮换后恢复证据 | **已完成**（GOAL-004 done · A-001 self pass · 双方言循环全绿） |
| R4 | 默认单密钥仍可用 | **已完成**（GOAL-005 done · A-001 self pass · 6/6 判据面成立） |
| R5 | 显式双密钥：轮换路径 **与** 恢复路径证据 | **已完成**（GOAL-006 done · 四项新鲜实跑全 PASS · Root 关门审计响应闭合） |

## 固定共享资料引用

> `shared-materials/index.json` 只能提供候选路径与摘要。缺完整引用字段的行无效。

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | — | — | — |
