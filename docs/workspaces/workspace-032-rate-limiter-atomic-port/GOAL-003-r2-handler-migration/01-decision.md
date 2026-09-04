---
id: GOAL-003-r2-handler-migration
title: R2 生产使用点迁移与 handler 回归
status: active
parent: GOAL-001-rate-limiter-atomic-port
created: 2026-09-03
updated: 2026-09-04
version: 0.2.0
---

# GOAL-003-r2-handler-migration · 01-decision 索引

| id | date | scope | summary | status |
|----|------|-------|---------|--------|
| [D-001-inherit-r1-contract-and-migration-scope](01-decision/D-001-inherit-r1-contract-and-migration-scope.md) | 2026-09-03 | 迁移范围与口径 | 继承 R1 D-002 合同，锁定 14 处生产分母与立即消费/失败预算迁移口径 | accepted |
| [D-002-tokenized-reservation-failure-budget](01-decision/D-002-tokenized-reservation-failure-budget.md) | 2026-09-04 | A-002 响应 · 令牌化保留 | 用户裁决方案 A：内核 Reserve/Cancel 契约（I-032-003）+ 10 处失败预算逐路径语义冻结；取代 GOAL-002 D-002 §4 失败预算口径 | accepted |

## 决策记录（ledger）

`01-decision/` 平铺；编号递增；正文只写已发生或用户已确认的决策。
