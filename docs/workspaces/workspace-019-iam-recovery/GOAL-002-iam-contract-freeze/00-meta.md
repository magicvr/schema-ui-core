---
id: GOAL-002-iam-contract-freeze
title: R1 IAM 合同冻结（恢复 / 策略 / 邀请）
status: done
parent: GOAL-001-iam-recovery
created: 2026-08-25
updated: 2026-08-25
version: 1.0.0
progress: 3/3
plan_refs:
  - VP-019-iam-recovery
primary_plan: VP-019-iam-recovery
serves_summary: 承接 Root R1：在已 verified 的 I-001/I-002/I-009（Root D-002）之上，继续裁决 I-003/I-004/I-005 并投影确认 I-007/I-008 与 VP 已冻结项，把 IAM 三件（自助恢复 / 密码策略 / 邀请入职）的实施合同一次冻清，供 R2 全链与 R3 实施。不写 DDL、不改应用代码。
---

# GOAL-002 · R1 IAM 合同冻结（恢复 / 策略 / 邀请）

## 概述

本目标承接 Root `GOAL-001-iam-recovery` 纲领阶段 **R1**：在写任何 schema 或流代码之前，把 IAM 波次合同一次冻清——

1. **自助恢复状态机**：证明形态与时效数值已由 Root D-002 冻结输入（6 位码 / 10 min TTL / 60 s 冷却 / 错 5 次作废 / MFA 第二因子门）；本目标把它写成可实施的合同条款（发起资格、挑战生命周期、完成动作、错误语义、审计事件）。
2. **密码策略**：VP-019 已冻结配置边界（UI 仅 `admin.settings` tab 扩展、强制面 `core.auth-session` 全 Profile 生效、禁止为策略改 mvp 默认集）；默认参数（I-003）与既有账号生效边界（I-007 投影确认）在本目标冻结。
3. **邀请入职**：形态与建号方式（I-004）、有效期/一次性/撤销语义（I-005）在本目标裁决并冻结。
4. **会话语义**：I-008（non-blocking，最晚 R4）可顺带投影现行改密行为提前确认。

对齐递归：GOAL-002 → Root GOAL-001（R1）→ VP-019 → Charter @0.2.0。不进入 SMS、模板中心、多邮箱、组织权限、OIDC 或业务域；不重开 VP-012～018。

## 检查点（progress 来源）

| # | 检查点 | 证据 |
|---|--------|------|
| C1 | I-001 / I-002 / I-009 关闭（verified；用户书面裁决留痕），立项解锁 | **完成**：Root D-002 + 结构化裁决会话（i001_proof_form / i002_ttl_cooldown / i009_mfa_boundary，均取推荐项）；Root 镜像表已同步 |
| C2 | IAM 合同冻结决策落盘（本目标 D-NNN）：恢复状态机条款 + 策略默认参数与生效边界 + 邀请形态与生命周期 + 会话语义投影；I-003/I-004/I-005 verified，I-007/I-008 投影留痕 | **完成**：D-001 §1～§5 五节条款；第二轮裁决会话五项均取推荐项；Root 镜像表 I-003～I-008 已同步 |
| C3 | 自审 A-001 闭合（self；无开放 required finding） | **完成**：A-001 self `pass`（0 required；N-1～N-3 notes 移交 R2/R3 设计，不阻断关门） |

`progress` = 已完成检查点 / 3。当前 **3/3**（已关门 · 2026-08-25）。

## 边界

- 只冻结合同；恢复全链实施归 R2；策略/邀请实施归 R3；端到端证据归 R4。
- 不写 DDL、不改应用代码、不改 Profile 默认集 / 模块矩阵 / Manifest 装配语义。
- 不动管理员重置既有特权路径（`must_change_password`）与 admin.mfa 既有面。
- 审计模式：开题 **none**；合同冻结落盘后 **self**；R2 schema 迁移与各实施切片按 Root 决策走 **independent**（grok build · grok-4.6 · high · `/audit`）。

## 成功标准

1. I-003 / I-004 / I-005 均 verified（用户书面裁决留痕）；I-007 / I-008 投影确认留痕；Root 镜像表同步。
2. 合同条款可核对：恢复流程每一步的前置、动作与失败语义；策略参数、配置面与强制面归属；邀请全链生命周期——R2/R3 无需再猜。
3. 未实施 R2/R3 产品面；未越界进入 SMS / 模板中心 / 多邮箱 / 业务域；未重开相邻 VP 已关门契约。
