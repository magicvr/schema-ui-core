---
id: GOAL-003-r2-handler-migration
title: R2 生产使用点迁移与 handler 回归
status: active
parent: GOAL-001-rate-limiter-atomic-port
created: 2026-09-03
updated: 2026-09-04
version: 0.2.0
---

# GOAL-003-r2-handler-migration · 02-execution 索引

| id | date | scope | summary | status |
|----|------|-------|---------|--------|
| [E-001-goal-opened](02-execution/E-001-goal-opened.md) | 2026-09-03 | 立项 | 立项 R2；继承 R1 成果与 14 处分母；准备执行 C1 迁移 | active |
| [E-002-handler-migration](02-execution/E-002-handler-migration.md) | 2026-09-03 | 14处迁移+回归 | 14 处生产调用点全部迁移至 AllowRecord；消除 TOCTOU；补齐并发无穿透与净状态等价测试；commit b08798d4（**初版，被 A-002 证伪后按 E-003 修正**） | completed |
| [E-003-tokenized-reservation-fix](02-execution/E-003-tokenized-reservation-fix.md) | 2026-09-04 | A-002 响应修复 | 响应 A-002 F-001/F-002：内核新增 Reserve/Cancel；10 处失败预算按 D-002 逐路径语义冻结改造；补齐 no-path 累计 / 不清历史 / 混合序列回归；全量 + race 全绿 | completed |

## 执行记录（ledger）

`02-execution/` 平铺；编号递增；时间线只记事实。
