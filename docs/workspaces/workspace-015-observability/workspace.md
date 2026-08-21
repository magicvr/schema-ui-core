---
id: workspace-015-observability
title: 可观测性工作区
status: active
root_goal: GOAL-001-observability
canonical_scope: docs/workspaces/workspace-015-observability/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-015-observability
primary_plan: VP-015-observability
created: 2026-08-21
updated: 2026-08-21
version: 0.1.0
parent: null
---

# 工作区上下文 · 可观测性

本工作区是 [VP-015-observability](../../vision/plans/VP-015-observability.md)（**`active`**，架构 A4）的唯一 lead delivery workspace。

- **Root** `GOAL-001-observability`：纲领 R1～R5；R1 已完成（[GOAL-002](GOAL-002-metrics-export-contract/00-meta.md) done 3/3）；R2 已完成（[GOAL-003](GOAL-003-metrics-scrape-endpoint/00-meta.md) done 4/4）；R3 已完成（[GOAL-004](GOAL-004-otel-traces/00-meta.md) done 4/4）。
- 不改变 Charter `primary_workspace`（仍为 workspace-001）。
- 不重开 `workspace-014-object-storage`。
- 不承接 Admin 功能或业务域；不重开 VP-012 / VP-013 / VP-014。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-015-observability` | 与本区目标及资料引用的 `workspace_id` 一致 |
| Root Goal | `GOAL-001-observability` | `parent: null`；交付目标 |
| canonical 范围 | `docs/workspaces/workspace-015-observability/` | 本区唯一目标状态范围 |
| 共享资料目录 | `none` | 暂无固定共享资料 |
| 愿景角色 | `delivery` | VP-015 lead；不改变 Charter primary workspace |
| 规划对齐 | `primary_plan` = `VP-015-observability`（`active`） | 架构 A4 指标导出 + OpenTelemetry |

## 愿景对齐

Charter：`schema-ui-core-admin-foundation@0.2.0`。  
VP-015：Prometheus 类指标导出 + OTLP traces；无收集器为内嵌默认；系列/诊断面须带 `module_id`、不得泄露秘密。  
与 VP-009/VP-010 正交：安全/符合性 gap 不进本区。

## 纲领阶段（Root 路线图指针）

| 阶段 | 内容 | 状态 |
|------|------|------|
| R1 | 导出合同与配置面冻结 | 已完成（GOAL-002 done 3/3） |
| R2 | 指标 scrape 接入 | 已完成（GOAL-003 done 4/4） |
| R3 | OpenTelemetry traces 接入 | 已完成（GOAL-004 done 4/4；I-002 已闭合） |
| R4 | 与现有 request-id / correlation 关联 | 待立项（I-005 须先闭合） |
| R5 | 默认无收集器 + 显式导出双路径证据 | 待立项 |

## 固定共享资料引用

> `shared-materials/index.json` 只能提供候选路径与摘要。缺完整引用字段的行无效。

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | — | — | — |
