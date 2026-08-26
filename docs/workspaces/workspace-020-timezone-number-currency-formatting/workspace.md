---
id: workspace-020-timezone-number-currency-formatting
title: 时区 / 数字 / 货币格式语义工作区
status: active
root_goal: GOAL-001-timezone-number-currency-formatting
canonical_scope: docs/workspaces/workspace-020-timezone-number-currency-formatting/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-020-timezone-number-currency-formatting
primary_plan: VP-020-timezone-number-currency-formatting
created: 2026-08-26
updated: 2026-08-26
version: 0.5.0
parent: null
---

# 工作区上下文 · 时区 / 数字 / 货币格式语义

本工作区是 [VP-020-timezone-number-currency-formatting](../../vision/plans/VP-020-timezone-number-currency-formatting.md)（**`active`** · 2026-08-26 激活；VRev-044 self `pass`）的唯一 lead delivery workspace。

- **Root** `GOAL-001-timezone-number-currency-formatting`：**`active`** · 3/4（R1～R3 已关门；R4 待立项）。
- 不改变 Charter `primary_workspace`（仍为 workspace-001）。
- 消费已交付基架：VP-007 locale 运行时（closed v0.3.0）、VP-005 设计系统（closed）、VP-011 用户/角色边界（closed）。
- 边界：不承接汇率/换算/计费（业务域）、DB `timestamptz` 持久化合同（架构 RT-T03，仍 `registered`）、翻译中心（VP-007）、改 Profile 默认集 / 模块矩阵 / Manifest 装配语义。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-020-timezone-number-currency-formatting` | 与本区目标及资料引用的 `workspace_id` 一致 |
| Root Goal | `GOAL-001-timezone-number-currency-formatting` | `parent: null`；**active** · 3/4（R1～R3 已关门；R4 待立项） |
| canonical 范围 | `docs/workspaces/workspace-020-timezone-number-currency-formatting/` | 本区唯一目标状态范围 |
| 共享资料目录 | `none` | 暂无固定共享资料 |
| 愿景角色 | `delivery` | VP-020 lead；不改变 Charter primary workspace |
| 规划对齐 | `primary_plan` = `VP-020-timezone-number-currency-formatting`（`active`） | 2026-08-26 激活/开区（VRev-044 self `pass`；Admin 类 freshness PASS） |

## 愿景对齐

Charter：`schema-ui-core-admin-foundation@0.2.0`。
VP-020：时区 / 数字 / 货币格式语义（Admin 功能分支基架能力剩余 #5）；消费 VP-007 locale 运行时与 VP-005 设计系统；DB 时区持久化合同仍归架构 RT-T03。
与 VP-009/VP-010 正交：格式相关安全/符合性 gap 不进本区。
VP-008 `go` 消费有效性：Admin 类 freshness **PASS**（`66f5fd1f` → `c6fda691`，2026-08-26；VRev-044）。

## 纲领阶段（Root 路线图指针）

| 阶段 | 内容 | 状态 |
|------|------|------|
| R1 | 合同冻结：时区来源（I-020-001）、数字/货币落点（I-020-002）、设置归属（I-020-005）；I-003/I-004 保持 VP 冻结投影 | **已关门**（GOAL-002 done 3/3 · A-001 self pass；合同正文 = GOAL-002 D-001） |
| R2 | 时区语义：会话/用户级解析与展示（IANA / offset / auto） | 依赖 R1 | **已关门**（GOAL-003 done 5/5 · A-001 self pass；timezone.ts + 头部时区选择 + 站点默认接入 + 统一语义） |
| R3 | 数字 / 货币语义：locale 驱动格式与输入解析合同 | 依赖 R2 | **已关门**（GOAL-004 done 6/6 · A-001 self pass → A-002 grok fail(2 required) → fixed → A-004 grok 复审 pass；money.ts + defaultCurrency 端到端） |
| R4 | 证据与关门：快测 + 双 locale 范例；无越界；required = 0 | 依赖 R3 | **已立项**（GOAL-005 active 3/4 · 证据矩阵/越界核账/F-007 加严已完成；Root 关门审计待执行） |

## 固定共享资料引用

> `shared-materials/index.json` 只能提供候选路径与摘要。缺完整引用字段的行无效。

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | — | none | — |