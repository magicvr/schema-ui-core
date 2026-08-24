---
id: GOAL-005-r4-evidence
title: R4 证据与关门就绪
status: active
parent: GOAL-001-account-email-identity
created: 2026-08-24
updated: 2026-08-24
version: 0.1.0
progress: 0/3
plan_refs:
  - VP-018-account-email-identity
primary_plan: VP-018-account-email-identity
serves_summary: 承接 Root R4：端到端证据（绑定→校验信经 017 默认渠道出站记录可取→完成校验）、唯一性 fail-closed 可核对、无 IAM/邀请/密码策略越界声明、N-1 边界说明；self 审计后关门并交 Root 关门环。
---

# GOAL-005 · R4 证据与关门就绪

## 概述

本目标承接 Root 纲领阶段 **R4**：把 R1～R3 的交付固化为可核对证据——核心是一条**经真实渠道适配器**的端到端流（BindEmail → `mail.OutboxSink` 出站记录 → 从记录取码 → VerifyEmail → verified），加唯一性 fail-closed 与边界声明。

对齐递归：GOAL-005 → Root GOAL-001（R4）→ VP-018 → Charter @0.2.0。

## 检查点（progress 来源）

| # | 检查点 | 证据 |
|---|--------|------|
| C1 | 端到端证据测试落盘且绿：校验信经 OutboxSink（017 默认渠道适配器）写入出站记录并可取码完成校验 | **完成**：TestR4EndToEndBindVerifyThroughMockChannel PASS；两阶段派发修正（E-001） |
| C2 | 唯一性 fail-closed + 无邮箱账号不受影响 + N-1 边界声明（附件） | **完成**：e2e 同链断言 + `attachments/r4-evidence.md` |
| C3 | self 审计闭合（无开放 required） | **完成**：A-001 self pass |

`progress` = 已完成检查点 / 3。当前 **3/3**（已关门）。

## 边界

- 不改产品代码；证据缺口回 R2/R3 整改而非在本目标补丁。
- 审计模式：self（Root D-001 映射；Root 关门另走 independent 环）。
