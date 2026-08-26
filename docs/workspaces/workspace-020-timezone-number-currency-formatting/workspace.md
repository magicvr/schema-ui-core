---
id: workspace-020-timezone-number-currency-formatting
title: 时区 / 数字 / 货币格式语义工作区
status: done
root_goal: GOAL-001-timezone-number-currency-formatting
canonical_scope: docs/workspaces/workspace-020-timezone-number-currency-formatting/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-020-timezone-number-currency-formatting
primary_plan: VP-020-timezone-number-currency-formatting
created: 2026-08-26
updated: 2026-08-27
version: 0.6.0
parent: null
---

# 工作区上下文 · 时区 / 数字 / 货币格式语义（已结项）

本工作区是 [VP-020-timezone-number-currency-formatting](../../vision/plans/VP-020-timezone-number-currency-formatting.md)（**`closed`** v0.3.0 · 2026-08-27 关门 · VRev-044 self `pass` + VRev-045 关门审查）的唯一 lead delivery workspace。**工作区已结项**（2026-08-27 用户书面确认）：历史绑定保留，默认不接新区。

- **Root** `GOAL-001-timezone-number-currency-formatting`：**`done`** · 4/4（R1～R4 全部关门；关门审计 A-001 self + A-002 grok independent 双 `pass`，0 required）。
- 不改变 Charter `primary_workspace`（仍为 workspace-001）。
- 消费已交付基架：VP-007 locale 运行时（closed v0.3.0）、VP-005 设计系统（closed）、VP-011 用户/角色边界（closed）。
- 边界（保持与开区一致）：不承接汇率/换算/计费（业务域）、DB `timestamptz` 持久化合同（架构 RT-T03，仍 `registered`）、翻译中心（VP-007）、改 Profile 默认集 / 模块矩阵 / Manifest 装配语义。
- 全部残余项随 R4 核账结清（用户书面接受）：`defaultCurrency` 句法三字母（非完整 ISO 目录）、分组位序容差、业务金额展示不接线（VP-020 冻结）；安全整数已加严 fixed。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-020-timezone-number-currency-formatting` | 与本区目标及资料引用的 `workspace_id` 一致 |
| Root Goal | `GOAL-001-timezone-number-currency-formatting` | `parent: null`；**done** · 4/4（2026-08-27 结项） |
| canonical 范围 | `docs/workspaces/workspace-020-timezone-number-currency-formatting/` | 本区唯一目标状态范围 |
| 共享资料目录 | `none` | 暂无固定共享资料 |
| 愿景角色 | `delivery` | VP-020 lead（closed 历史绑定）；不改变 Charter primary workspace |
| 规划对齐 | `primary_plan` = `VP-020-timezone-number-currency-formatting`（`closed` v0.3.0） | 2026-08-26 激活/开区（VRev-044 self `pass`；Admin 类 freshness PASS）；2026-08-27 关门（VRev-045） |

## 愿景对齐

Charter：`schema-ui-core-admin-foundation@0.2.0`。
VP-020：时区 / 数字 / 货币格式语义（Admin 功能分支基架能力剩余 #5）——**已 close**（v0.3.0，2026-08-27）：全套交付证据见本区 Goal 台账（R1 合同 → R2 时区 → R3 数字/货币 → R4 证据与关门）；DB 时区持久化合同仍归架构 RT-T03（`registered`）。
与 VP-009/VP-010 正交：格式相关安全/符合性 gap 归持续程序。

## 纲领阶段（Root 路线图指针 · 全部已关门）

| 阶段 | 内容 | 状态 |
|------|------|------|
| R1 | 合同冻结：时区来源（I-020-001）、数字/货币落点（I-020-002）、设置归属（I-020-005）；I-003/I-004 保持 VP 冻结投影 | **已关门**（GOAL-002 done 3/3 · 合同正文 = GOAL-002 D-001） |
| R2 | 时区语义：会话/用户级解析与展示（IANA / offset / auto） | 依赖 R1 | **已关门**（GOAL-003 done 5/5） |
| R3 | 数字 / 货币语义：locale 驱动格式与输入解析合同 | 依赖 R2 | **已关门**（GOAL-004 done 6/6 · grok 独立审闭环） |
| R4 | 证据与关门：快测 + 双 locale 范例；无越界；required = 0 | 依赖 R3 | **已关门**（GOAL-005 done 4/4 · Root 关门审计双 pass） |

## 结项记录

| 日期 | 结论 | 证据 | 残余 |
|------|------|------|------|
| 2026-08-27 | Root done 4/4 · 工作区结项（用户书面确认）；VP-020 closed v0.3.0 | Root `03-audit/A-001`（self pass）+ `A-002`（grok independent pass，0 required）；证据矩阵 `GOAL-005/attachments/r4-evidence-matrix.md`；全量回归（Go 全绿 / web 1181） | 无开放必改；残余均书面接受（分组位序、币种句法近似、业务展示不接线等，见 R4 核账表） |

## 固定共享资料引用

> `shared-materials/index.json` 只能提供候选路径与摘要。缺完整引用字段的行无效。

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | — | none | — |