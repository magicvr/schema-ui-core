---
id: GOAL-002-identity-contract-freeze
title: R1 身份合同冻结
status: done
parent: GOAL-001-account-email-identity
created: 2026-08-24
updated: 2026-08-24
version: 1.0.0
progress: 3/3
plan_refs:
  - VP-018-account-email-identity
primary_plan: VP-018-account-email-identity
serves_summary: 承接 Root R1：关闭 I-001（非空邮箱唯一性细则）与 I-002（校验投递形态），冻结账号邮箱身份合同（可空、绑定占槽语义、唯一性失败语义、换绑同合同、状态机三态边界），供 R2 schema 与 R3 绑定流实施。不写 DDL、不改应用代码。
---

# GOAL-002 · R1 身份合同冻结

## 概述

本目标承接 Root `GOAL-001-account-email-identity` 纲领阶段 **R1**（D-003 解冻后待启动 → 进行中）：在写任何 schema 或绑定流代码之前，把账号邮箱身份合同一次冻清——唯一性细则（含未校验地址是否占槽、大小写/规范化）、校验投递形态（验证码 vs 魔术链接）、可空与换绑（VP 已冻结，投影落款）、状态机三态（未绑定 / 待校验 / 已校验）边界。

VP-018 已冻结项：`users` email 可空（I-003）、换绑进本波（I-004）。本目标裁决并关闭 **I-001 / I-002**。

对齐递归：GOAL-002 → Root GOAL-001（R1）→ VP-018 → Charter @0.2.0。不进入自助恢复、邀请、密码策略、SMS、模板或业务域；不重开 VP-012～017。

## 检查点（progress 来源）

| # | 检查点 | 证据 |
|---|--------|------|
| C1 | I-001 / I-002 关闭（verified；用户书面裁决留痕） | **完成**：会话裁决答复（i002_form / i001_slot / i001_norm）；Root 镜像表已同步 |
| C2 | 身份合同冻结决策落盘（本目标 D-NNN），条款可供 R2 schema 与 R3 绑定流直接实施 | **完成**：D-001 §1～§7 七条款 |
| C3 | 自审 A-001 闭合（self；无开放 required finding） | **完成**：A-001 self pass（0 required；F-1/F-2 note 已留痕/移交） |

`progress` = 已完成检查点 / 3。当前 **3/3**（已关门）。

## 边界

- 只冻结合同；双方言 schema 归 R2；绑定/校验流归 R3；证据归 R4。
- 不写 `users` DDL、不改应用代码、不改 Profile 默认集。
- 不删除、不回退 VP-017 运输面实施史；运输验收对齐其现行默认渠道。
- 审计模式：开题 **none**；合同冻结落盘后 **self**；schema 迁移与绑定实施按 Root D-001 走 **independent**（grok build · grok-4.6 · `/audit`）。

## 成功标准

1. I-001 / I-002 均 verified 或经用户书面接受残余；Root 镜像表同步。
2. 合同条款可核对：可空性、占槽语义、唯一性失败语义、规范化规则、投递形态、换绑路径、状态机三态——R2/R3 无需再猜。
3. 未实施 R2/R3 产品面；未越界进入 IAM 恢复或模板中心。
