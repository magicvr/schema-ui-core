---
id: GOAL-006-channel-provider-contract
title: R5 渠道供应商合同冻结
status: active
parent: GOAL-001-outbound-mail
created: 2026-08-24
updated: 2026-08-24
version: 0.1.0
progress: 0/3
plan_refs:
  - VP-017-outbound-mail
primary_plan: VP-017-outbound-mail
serves_summary: 承接 Root R5：冻结具名渠道与 MailSender 的关系、mock 站内语义与持久化、SMTP 保留规则；关闭 I-011。不实施 Resend/设置页，不回退 R1～R4 代码。
---

# GOAL-006 · R5 渠道供应商合同冻结

## 概述

本目标承接 Root `GOAL-001-outbound-mail` 纲领阶段 **R5**：在写 Resend 适配器或设置页之前，把渠道供应商合同一次冻清——渠道标识、与 `kernel.MailSender` 的解析关系、mock 站内出站记录的持久化与取信面、SMTP 适配器保留为渠道而非删除。

Root D-006 已冻结 I-007（渠道集）、I-008（mock = 管理员记录）、I-012（邮件 tab 形状）。本目标关闭 **I-011**（mock 持久化）并写出可供 R6 实施的渠道注册形状。

对齐递归：GOAL-006 → Root GOAL-001（R5）→ VP-017 v0.4.0 → Charter @0.2.0。不进入账号 email、邀请、自助恢复、用户站内通知或业务域。不重开 GOAL-002～005。

## 检查点（progress 来源）

| # | 检查点 | 证据 |
|---|--------|------|
| C1 | 渠道合同冻结决策落盘（本目标 D-00N；Root I-011 → verified） | 决策文件 |
| C2 | 合同可被 R6 消费：渠道 id / 默认 mock / SMTP 保留 / 公共面无供应商类型，有可核对说明或最小类型骨架 | 决策附件或 kernel/composition 注释级骨架（若本目标选择只写决策则 C2=决策中的可实施条款） |
| C3 | 自审 A-001 闭合（无开放 required finding） | 本目标 `03-audit/` |

`progress` = 已完成检查点 / 3。当前 **0/3**。

## 边界

- 只冻结合同；Resend HTTP 适配器、mock 产品表落地、设置页、热切换、试发分别归 R6/R7。
- 不删除、不回退 SMTP / CaptureSink / `MailSender` 实施史。
- 不改 Profile / 模块矩阵；不重开 VP-007。
- 开项目标审计模式 **none**（文档 scaffold）；合同冻结落盘后走 **self**。

## 成功标准

1. I-011 已 verified 或合规 residual；mock 持久化与管理端取信面可被 R6 实施。
2. 渠道解析规则可核对：默认 mock；显式 Resend；SMTP 仍可被选为渠道；模块只见 `MailSender`。
3. 未实施 R6/R7 产品面；未解冻 018。
