---
id: GOAL-003-r2-handler-migration
title: R2 生产使用点迁移与 handler 回归
status: active
parent: GOAL-001-rate-limiter-atomic-port
created: 2026-09-03
updated: 2026-09-04
version: 0.2.0
---

# GOAL-003-r2-handler-migration · 03-audit 索引

| id | date | source | scope | verdict | open required | summary | file |
|----|------|--------|-------|---------|---------------|---------|------|
| [A-001](03-audit/A-001-r2-handler-migration-self-audit.md) | 2026-09-03 | self | GOAL-003 关门自审（14处迁移+回归） | pass | 0 | 14 处生产调用点全部迁移至 AllowRecord；并发无穿透与净状态等价测试通过；全量回归全绿；commit b08798d4（**行为等价主张被 A-002 证伪，按 E-003/A-003 修正**） | [A-001-r2-handler-migration-self-audit.md](03-audit/A-001-r2-handler-migration-self-audit.md) |
| [A-002](03-audit/A-002-r2-handler-migration-independent.md) | 2026-09-04 | independent | GOAL-003 R2 生产迁移与 handler 回归关门 | **fail** | **2** | 14 处静态迁移与现有测试可复现；失败预算的 AllowRecord+Clear 会丢失既有失败历史，recovery start no-path 不再累计，行为等价与关门证据不成立（F-001/F-002） | [A-002-r2-handler-migration-independent.md](03-audit/A-002-r2-handler-migration-independent.md) |
| [A-003](03-audit/A-003-tokenized-reservation-fix-self-audit.md) | 2026-09-04 | self | A-002 F-001/F-002 修复复审（令牌化 Reserve/Cancel） | pass | 0 | 用户裁决方案 A；Reserve/Cancel 落地 + 10 处逐路径语义冻结；no-path 累计 / 不清历史 / 混合序列回归全绿；commit 3bfe66c2 | [A-003-tokenized-reservation-fix-self-audit.md](03-audit/A-003-tokenized-reservation-fix-self-audit.md) |
| [A-004](03-audit/A-004-tokenized-reservation-fix-independent.md) | 2026-09-04 | independent（grok-build · grok-4.6 · reasoning high） | A-002 F-001/F-002 令牌化修复关门复审 | **pass** | **0** | 独立复跑合同与五条回归（含 -race）；十处语义与 D-002 §3 一致；F-001/F-002 可按 fixed 闭合；R-001（登录 500 分支占槽）/R-002（过期文案投影）recommended | [A-004-tokenized-reservation-fix-independent.md](03-audit/A-004-tokenized-reservation-fix-independent.md) |

## 信息就绪核对（按 scope）

| 核对项 | 状态 | 备注 |
|--------|------|------|
| I-032-001/002 | verified | 继承 R1 冻结；不阻断实施 |
| 到期 required | 无 | — |
| 资料引用 | 无 | `shared_materials_catalog: none` |

## 审计记录（ledger）

`03-audit/` 平铺；编号递增；意见必须落盘（self / independent 共用序列）。
