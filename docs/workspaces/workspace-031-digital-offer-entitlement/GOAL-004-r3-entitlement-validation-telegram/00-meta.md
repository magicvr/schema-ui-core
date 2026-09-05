---
id: GOAL-004-r3-entitlement-validation-telegram
title: R3 权益核验/消耗 + 可选 Telegram 注册
status: active
parent: GOAL-001-digital-offer-entitlement
created: 2026-09-05
updated: 2026-09-05
version: 0.1.0
progress: 2/3
plan_refs:
  - VP-031-digital-offer-entitlement
primary_plan: VP-031-digital-offer-entitlement
serves_summary: 承载 VP-031 R3 实施（分母 = GOAL-002 D-002 v1.2.0）：权益统一核验 API（§5.1 聚合规则）、count 型原子消耗（§5.2 跨方言算法）、BIZOFFER_ENTITLEMENT_* 错误码登记、Telegram price/buy/entitlements 命令 Register（§6）与 price/entitlements 查询限流桶（§8）。
---

# GOAL-004 · R3 权益核验/消耗 + 可选 Telegram 注册

## 概述

执行 Root 纲领 **R3**：按 GOAL-002 D-002 v1.2.0 合同实施权益统一核验（Check，§5.1 聚合优先级）与 count 型原子消耗（Consume，§5.2 跨方言算法与 void 线性化）、`BIZOFFER_ENTITLEMENT_INVALID` / `_INSUFFICIENT` 错误码登记、可选 Telegram 命令（`price` / `buy` / `entitlements`，经 `kernel.TelegramDispatcher` Register，未启用时 DisabledDispatcher 零感知）以及 price/entitlements 查询限流桶。承接 GOAL-003 A-001 F-001 遗留：purchase/price 桶在通道入口的接线。

对齐递归：GOAL-004 → Root GOAL-001（R3）→ VP-031（判据 3/5）→ Charter @0.4.0。

## 纲领检查点（P-001）

| 检查点 | 内容 | 状态 |
|--------|------|------|
| C1 | **核验/消耗落地**：Check（§5.1 聚合）+ Consume（§5.2 算法）+ 错误码登记 + 双库测试（含 RowsAffected 竞争、并发 void） | **已关门**（2026-09-05 · E-001） |
| C2 | **Telegram Register + 查询桶**：price/buy/entitlements 命令、`tg:` request_id 派生、price/entitlements 查询桶、未启用通道零依赖测试 | **已关门**（2026-09-05 · E-001） |
| C3 | **审视与关门**：self 审计 + codex independent 审计（判据 3/5 门禁）；意见响应；Root/goal-tree 回写 | 待开始 |

`progress` = 已关门检查点数 / 3。当前 **2/3**。

## 成功标准（方向级）

1. Check 以 (subjectID, offerID) 统一核验，聚合 reason 确定性（valid/no_entitlement/expired/exhausted/voided），判据 3 三态可测（表驱动）。
2. Consume 按 §5.2 算法原子扣减：全有或全无、RowsAffected 竞争重读、与 void 线性化；双数据库测试通过。
3. Telegram 命令在通道启用时可注册（DisabledDispatcher no-op 时模块测试不依赖 Bot API）；空 SubjectID fail-closed（判据 5）。
4. 关门前开放 required finding = 0（C3 审计闭合）。

## 信息就绪与未知项

R1/R2 已关闭全部信息项；R3 无新增 required 信息项。实现细节以 D-002 v1.2.0 条款为准。

## 父目标

- `GOAL-001-digital-offer-entitlement`（Root · 纲领 R3）

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺记账；索引文件在本目标 `01-decision.md` / `02-execution.md` / `03-audit.md`。

## 备注

- 审计模式：C3 关门按 AGENTS P-003 为 **independent**（判据 3/5 门禁 + 通道面），provider = 本地 codex（gpt-5.6-sol · medium）。
- 承接 GOAL-003 A-001 F-001 / A-003：purchase 桶已在 Service API 生效；本目标补 Telegram 入口的 purchase/price/entitlements 桶接线。
