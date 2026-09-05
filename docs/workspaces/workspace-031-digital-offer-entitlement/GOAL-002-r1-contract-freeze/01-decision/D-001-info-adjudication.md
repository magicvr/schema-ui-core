---
doc_type: goal-decision
id: D-001-info-adjudication
parent: GOAL-002-r1-contract-freeze
date: 2026-09-05
status: accepted
version: 1.0.0
---

# D-001 · R1 信息裁决（I-031-001～005）

## 决定

2026-09-05 用户在 P-004 裁决点书面裁决 R1 三个 required 信息项，全部采纳编排器建议项：

| 信息项 | 问题 | 裁决 | 级别 |
|--------|------|------|------|
| I-031-001 | 首波权益形态 | **二者并存**：Offer 建档时二选一（`duration` → `expires_at`；`count` → `remaining_count`），一 Offer 固定一种 | required |
| I-031-002 | 购买状态最小子集 | **同步一拍 fulfilled**：同进程同步购买，freeze → deduct_frozen → 购买凭证 + 权益同事务落库即为 `fulfilled`；失败整体回滚（或等价补偿），不落购买凭证；**无 `pending` 态** | required |
| I-031-003 | Admin 人工发放/撤销 | **只读 + 作废**：Admin 可查询权益、可作废（独立权限键 + 审计）；**不开放人工发放**，发放只走购买路径，资金路径不可绕过 | required |

随 R1 合同一并冻结的 non-blocking 默认项（lead 建议、用户可随时否决，否决走 02 决策）：

| 信息项 | 问题 | 默认裁决 | 级别 |
|--------|------|----------|------|
| I-031-004 | Telegram 命令清单 | `price`（价目）/ `buy`（购买）/ `entitlements`（我的权益）三命令；仅 `channel.telegram` 启用时 Register | non-blocking |
| I-031-005 | 模块 id | `biz.digital-offer`（对齐 `admin.*` / `channel.*` 前缀惯例，业务域独立命名空间） | non-blocking |

## 未选方案

- I-031-001 仅时长 / 仅次数：按次套餐或会员期二者必有其无法表达，后续补形态需迁移。
- I-031-002 引入 pending：首波无异步履约场景，pending 带来孤儿补偿与超时清理复杂度；同步单事务是更强的 fail-closed。
- I-031-003 开放人工发放：绕过支付路径，扩大滥用面；完全只读则误发/滥用无纠错手段。
- I-031-005 复用 `admin.*` 前缀：会与 Admin 面模块混淆；业务域应独立命名空间。

## 影响与回写

- 本裁决关闭 Root `GOAL-001-digital-offer-entitlement` 信息台账 I-031-001～003（required → verified）与 I-031-004～005（non-blocking → verified）；Root/VP-031 台账同步回写。
- 合同细节（字段、API、错误码、桶阈值）由 D-002 在本裁决框架内细化；偏离裁决框架需新的 02 决策。
