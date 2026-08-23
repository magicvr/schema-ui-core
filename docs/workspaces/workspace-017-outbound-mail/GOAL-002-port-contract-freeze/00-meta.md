---
id: GOAL-002-port-contract-freeze
title: R1 发送端口与合同冻结
status: done
parent: GOAL-001-outbound-mail
created: 2026-08-22
updated: 2026-08-22
version: 0.2.0
progress: 3/3
plan_refs:
  - VP-017-outbound-mail
primary_plan: VP-017-outbound-mail
serves_summary: 承接 Root R1：冻结内核同步 Send 合同（默认 sink 形态与取报文方式、To 基数、公共面禁 SMTP 客户端类型），关闭 I-001/I-002，并落地 kernel 邮件端口代码与单测。
---

# GOAL-002 · R1 发送端口与合同冻结

## 概述

本目标承接 Root `GOAL-001-outbound-mail` 纲领阶段 **R1**：在改任何 SMTP 接入代码之前，把内核出站邮件的发送合同一次冻清——端口形态、消息字段、默认 sink 形态与测试取报文方式、单次 `Send` 的 To 基数、公共面类型边界。合同以决策记录 + kernel 端口代码双重落盘，R2～R4 的实施只消费这份冻结件。

对齐递归：GOAL-002 → Root GOAL-001（R1）→ VP-017（RT-M01 / I-017-003 / I-017-004）→ Charter @0.2.0。不进入账号 email、邀请、自助恢复、模板产品或业务域。

## 检查点（progress 来源）

| # | 检查点 | 证据 |
|---|--------|------|
| C1 | 合同冻结决策落盘（Root D-002 + 本目标 D-001；Root 信息表 I-001/I-002 → verified） | 决策文件 |
| C2 | kernel 邮件端口代码 + 合同单测绿（`go test ./internal/kernel/`） | `apps/api/internal/kernel/mail.go` + 测试运行输出 |
| C3 | 自审 A-001 闭合（无开放 required finding） | 本目标 `03-audit/` |

`progress` = 已完成检查点 / 3。

## 边界

- 只冻结合同并落地端口类型；SMTP 适配器、配置键、capture sink 落地、`readyz` 分别归 R2/R3/R4 子目标。
- 不修改 composition 行为、不改 `readyz`、不改 handler / 模块公共契约（它们尚无邮件消费方）。
