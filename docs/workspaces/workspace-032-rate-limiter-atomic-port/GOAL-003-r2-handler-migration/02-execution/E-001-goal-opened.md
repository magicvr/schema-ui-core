---
doc_type: goal-execution
id: E-001-goal-opened
parent: GOAL-003-r2-handler-migration
date: 2026-09-03
status: active
version: 0.1.0
---

# E-001 · R2 立项与迁移准备

## 事实时间线

- 2026-09-03：R1 合同冻结目标 `GOAL-002-r1-contract-freeze` 完成 A-001 + A-002 审计响应并正式关门（A-003 `status: done`）。
- 2026-09-03：Root `GOAL-001` 纲领 R2 阶段启动，立项 `GOAL-003-r2-handler-migration`。
- 2026-09-03：五件套与三个 ledger 目录建齐，D-001 锁定 14 处生产迁移分母与口径。
- 2026-09-03：核对既有测试基线，准备执行 14 处调用点改造。

## 产物

- `docs/workspaces/workspace-032-rate-limiter-atomic-port/GOAL-003-r2-handler-migration/`
- `01-decision/D-001-inherit-r1-contract-and-migration-scope.md`

## 下一步（计划）

- C1：执行 14 处生产调用点重构（立即消费 4 处 + 失败预算 10 处）。
- C2：运行 handler 与 channel 限流回归测试。
