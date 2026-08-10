---
id: workspace-009-production-hardening
title: 生产加固工作区
status: active
root_goal: GOAL-001-production-hardening
canonical_scope: docs/workspace-009-production-hardening/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-009-production-hardening
primary_plan: VP-009-production-hardening
created: 2026-08-10
updated: 2026-08-10
version: 0.1.0
parent: null
---

# 工作区上下文 · 生产加固

本工作区承接 [VP-009-production-hardening](../vision/plans/VP-009-production-hardening.md)（`active`）的实现层治理。它承担 VP-008 §`go` 消费有效性规则下的**共享基架安全/健壮性缺陷重验证**：修复 2026-08-10 代码审查发现的缺陷（输入 `raw/audit-20260810-api-web-bug-review.md`，gitignored 临时记录），恢复 VP-008 `go` 的消费有效性。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-009-production-hardening` | 与本工作区所有目标及共享资料引用的 `workspace_id` 一致。 |
| Root Goal | `GOAL-001-production-hardening` | Root 固定为 `GOAL-001`，且 `parent: null`。 |
| canonical 范围 | `docs/workspace-009-production-hardening/` | 当前工作区唯一的目标状态范围。 |
| 共享资料目录 | `none` | 暂无固定共享资料；缺少完整引用字段的资料不得作为证据。 |
| 愿景角色 | `delivery` | 本区是 VP-009 的 lead delivery workspace，不改变 Charter 的 primary workspace。 |
| 规划对齐 | `plan_refs` / `primary_plan` = `VP-009-production-hardening` | 指向 [VP-009](../vision/plans/VP-009-production-hardening.md)。 |

## 愿景对齐

Charter 唯一来源为 `schema-ui-core-admin-foundation@0.2.0`。VP-009 已为 `active`，本区是其唯一 lead delivery workspace；VP-008 的 `go` 因共享基架安全/健壮性缺陷（2026-08-10 代码审查）按 VP-008 §`go` 消费有效性规则暂挂，由本区完成修复与重验证后恢复。

`I-READINESS-005` 的 independent provider（workspace-008 D-002）延续适用：**grok build（模型 `grok-4.5`、思考强度 high）执行 `audit` 命令**，审计模式 `cross`；provider 选择本身不是已完成审计证据。

## 固定共享资料引用

> `shared-materials/index.json` 只能提供候选路径与摘要。缺 `material_id`、`source`、`version`、64 位十六进制 `sha256` 或匹配 `workspace_id` 的行无效，不能作为事实、证据或跨工作区上下文来源。

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | — | — | — |
