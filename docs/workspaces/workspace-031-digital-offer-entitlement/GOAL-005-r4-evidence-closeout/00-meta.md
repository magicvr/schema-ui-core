---
id: GOAL-005-r4-evidence-closeout
title: R4 证据矩阵/边界核账/关门审计
status: done
parent: GOAL-001-digital-offer-entitlement
created: 2026-09-05
updated: 2026-09-05
version: 0.1.0
progress: 2/2
plan_refs:
  - VP-031-digital-offer-entitlement
primary_plan: VP-031-digital-offer-entitlement
serves_summary: 承载 VP-031 R4（分母 = GOAL-002 D-002 v1.2.0）：VP-031 判据 1～8 的证据矩阵（判据 → 代码/测试/审计路径）、边界核账（红线逐条核对）与关门审计（self + codex independent），支撑 Root GOAL-001 关门。
---

# GOAL-005 · R4 证据矩阵/边界核账/关门审计

## 概述

执行 Root 纲领 **R4**：汇总 R1～R3 的合同、实施与审计证据，构建 VP-031 方向级判据 1～8 的证据矩阵，逐条边界核账，完成关门审计并支撑 Root 关门。

对齐递归：GOAL-005 → Root GOAL-001（R4）→ VP-031（判据 6/7/8 + 全量）→ Charter @0.4.0。

## 纲领检查点（P-001）

| 检查点 | 内容 | 状态 |
|--------|------|------|
| C1 | **证据矩阵 + 边界核账**：判据 1～8 逐条映射证据路径；红线/Profile/subject-only/ Charter 未改核对 | **已关门**（2026-09-05 · E-001） |
| C2 | **关门审计与 Root 关门**：self + codex independent 关门审计；意见响应；Root/goal-tree/VP-031 投影同步 | **已关门**（第 2 次关门：A-006 fail 4 required → A-009 补齐 → A-010 independent `pass` 0 required） |

`progress` = 已关门检查点数 / 2。当前 **0/2**。

## 成功标准（方向级）

1. 判据 1～8 每条都有可核对证据路径（代码/测试/审计台账），无「证据不足」项。
2. 边界核账：无类目/SKU/税/库存/物流订单；未进默认 Profile；未解禁通用 Entitlement 接缝；购买与权益只挂 subject_id；未改 Charter。
3. 关门审计：开放 required = 0（fixed / accepted-residual / user-overruled）。

## 信息就绪与未知项

R1～R3 已关闭全部信息项；R4 无新增 required 信息项。

## 父目标

- `GOAL-001-digital-offer-entitlement`（Root · 纲领 R4）

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺记账；索引文件在本目标 `01-decision.md` / `02-execution.md` / `03-audit.md`。

## 备注

- 关门审计模式：**independent**（release/close-out 门禁），provider = 本地 codex（gpt-5.6-sol · medium）。
