---
id: GOAL-002-r1-contract-freeze
title: R1 合同冻结（Offer 字段 / 购买状态机 / 权益形态 / 命令清单 / 事务与限流边界）
status: active
parent: GOAL-001-digital-offer-entitlement
created: 2026-09-05
updated: 2026-09-05
version: 0.1.0
---

# GOAL-002-r1-contract-freeze · 03-audit 索引

| id | date | source | scope | verdict | open required | summary | file |
|----|------|--------|-------|---------|---------------|---------|------|
| （待 C3：A-001 self 合同自审；A-002 independent codex 审计） | | | | | | | |

## 信息就绪核对（按 scope）

| 核对项 | 状态 | 备注 |
|--------|------|------|
| I-031-001～003 | verified | 2026-09-05 用户书面裁决（D-001） |
| I-031-004～005 | verified | non-blocking 默认冻结（D-001），用户可否决 |
| V-F119 | 已纳入合同 | §8 冻结业务桶与禁 key-wide Clear |
| 到期 required | 无 | C1 已关门；C2/C3 无信息门禁 |

## 审计记录（ledger）

`03-audit/` 平铺；编号递增；意见必须落盘（self / independent 共用序列）。
