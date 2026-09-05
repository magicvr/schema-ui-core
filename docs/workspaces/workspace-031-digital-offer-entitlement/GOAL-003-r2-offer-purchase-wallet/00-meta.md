---
id: GOAL-003-r2-offer-purchase-wallet
title: R2 Offer CRUD + 购买 + 钱包扣款
status: active
parent: GOAL-001-digital-offer-entitlement
created: 2026-09-05
updated: 2026-09-05
version: 0.2.0
progress: 3/4
plan_refs:
  - VP-031-digital-offer-entitlement
primary_plan: VP-031-digital-offer-entitlement
serves_summary: 承载 VP-031 R2 实施（分母 = GOAL-002 D-002 v1.0.0）：digitaloffer 模块骨架（迁移/store/provider/manifest/schema）、Offer CRUD 与 Admin 面、C 端上架列表、单事务购买链路（freeze→deduct_frozen→凭证+权益）与 §4.4 并发验收测试。
---

# GOAL-003 · R2 Offer CRUD + 购买 + 钱包扣款

## 概述

执行 Root 纲领 **R2**：按 GOAL-002 D-002 v1.0.0 合同实施 `biz.digital-offer` 模块的 Offer CRUD、Admin 协议页面与权限键、C 端上架列表、购买单事务链路与并发验收测试。权益核验/消耗与 Telegram 命令归 R3（GOAL-004）。

对齐递归：GOAL-003 → Root GOAL-001（R2）→ VP-031（判据 1/2）→ Charter @0.4.0。实施分母 = D-002 v1.0.0；偏离合同需先修合同再实施。

## 纲领检查点（P-001）

| 检查点 | 内容 | 状态 |
|--------|------|------|
| C1 | **切片方案**：文件/接口清单映射 D-002 条款（迁移 DDL、store API、provider 贡献、composition 装配、错误码注册） | **已关门**（2026-09-05 · E-001） |
| C2 | **模块骨架落地**：`apps/api/modules/digitaloffer/`（migration + store + offer service + provider + manifest + schema）+ composition 装配 + errorcatalog 注册；构建与既有测试绿 | **已关门**（2026-09-05 · E-001；commit 0e65f392） |
| C3 | **购买链路与验收**：单事务购买（§4.2/§4.4）+ 并发验收测试（双数据库）+ Admin 路由/权限/审计 fail-closed 测试 + C 端列表 | **已关门**（2026-09-05 · E-001；SQLite + 真 PG 矩阵全绿） |
| C4 | **审视与关门**：self 审计 + codex independent 审计（资金路径门禁）；意见响应；Root/goal-tree 回写 | 进行中（A-001 → A-002 fail → A-003 响应修复完成，待 A-004 closure 复审） |

`progress` = 已关门检查点数 / 4。当前 **3/4**。

## 成功标准（方向级）

1. Offer CRUD 具备 Admin 协议页面、权限键（digitaloffer.read / offer.manage）、审计（fail-closed）；C 端 `GET /api/biz/offers` 列出 on_sale 项（判据 1 前半）。
2. 购买路径满足余额不足拒绝、`freeze → deduct_frozen`、失败整体回滚、凭证与权益同事务、幂等重放；并发验收测试双数据库通过（判据 2）。
3. 不越界：无类目/SKU/税/库存；不进默认 Profile；购买只挂 subject_id；迁移占位符 `?`（D-002 §10）。
4. 关门前开放 required finding = 0（C4 审计闭合）。

## 信息就绪与未知项

R1 已关闭全部信息项（I-031-001～005 verified）；R2 无新增 required 信息项。实现细节一律以 D-002 v1.0.0 条款为准，无合同外未知。

## 父目标

- `GOAL-001-digital-offer-entitlement`（Root · 纲领 R2）

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺记账；索引文件在本目标 `01-decision.md` / `02-execution.md` / `03-audit.md`。

## 备注

- 审计模式：C4 关门按 AGENTS P-003 为 **independent**（资金路径 + 数据迁移），provider = 本地 codex（gpt-5.6-sol · medium）；C2/C3 期间 self 审计兜底。
- 实施过程发现合同缺口 → 先修 D-002（02 决策）再实施，禁止实现期静默偏离。
