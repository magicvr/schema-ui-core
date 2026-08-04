---
id: workspace-002-production-admin-foundation
title: 生产级可用 Admin 基架工作区
status: active
root_goal: GOAL-001-production-admin-foundation
canonical_scope: docs/workspace-002-production-admin-foundation/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-002-production-admin-foundation
primary_plan: VP-002-production-admin-foundation
created: 2026-08-01
updated: 2026-08-04
version: 0.1.1
parent: null
---

# workspace-002-production-admin-foundation · 工作区上下文

## 绑定

| 字段 | 值 |
|------|----|
| Root Goal | `GOAL-001-production-admin-foundation` |
| canonical root | `docs/workspace-002-production-admin-foundation/` |
| 愿景角色 | `delivery` |
| primary plan | `VP-002-production-admin-foundation` |
| 共享资料目录 | `none` |

本工作区承接 [VP-002 · 生产级可用 Admin 基架](../vision/plans/VP-002-production-admin-foundation.md)，是该 VP 唯一历史 lead workspace。仓库级 `primary_workspace` 仍为 `workspace-001-mvp-admin-foundation`；本工作区不会重开已关闭的 VP-001 或已完成的旧 Root。

现行 Charter 为 `schema-ui-core-admin-foundation@0.2.0`。VP-002 在 `@0.1.0` 下关闭，2026-08-04 仅做精确 re-align；planned VP-003 尚未绑定工作区，本区不自动承接其实施。

## 规范范围

- 目标状态真相仅存在于本目录的 `goal-tree.md` 与 `GOAL-*` 五件套。
- 目标均在本目录根平铺；层级只由各 `00-meta.md` 的 `parent` 字段表达。
- 跨工作区历史仅使用带路径的 Q2 引用；不得建立跨工作区 `parent`。
- `shared_materials_catalog: none` 表示当前没有可作为事实或 finding 关闭依据的固定共享资料引用。

## 继承边界

VP-002 继承 VP-001 已冻结的 `I-PROTO-001 v0.1.3` 覆盖基线及其边界。该继承只限定首阶段的实施范围，不等于新工作区已经完成 Renderer、运行时或产品验收；扩大协议范围必须另行记录决策、版本与验证证据。

## 推进约定

Root 的五个纲领阶段原则上串行；同一阶段内部可以按依赖拆分并行子目标。只有当前阶段的信息门禁和方案边界就绪后，才按需创建具体子目标。
