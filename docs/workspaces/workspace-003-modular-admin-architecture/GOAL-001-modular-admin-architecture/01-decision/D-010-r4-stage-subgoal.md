---
id: D-010-r4-stage-subgoal
doc: decision-entry
goal: GOAL-001-modular-admin-architecture
date: 2026-08-05
status: accepted
---

# D-010 · 建立 R4 子目标与 C1 信息门禁

建立 [GOAL-005-r4-full-module-migration](../../GOAL-005-r4-full-module-migration/00-meta.md)
承接 Root R4，Root progress 保持 `3/6`。R4 先完成 C1 范围和信息冻结，再按
C2-C5 渐进实施；不得批量预建 R5/R6 或把 Users/Roles 当前代码当作模块迁移完成。

当前必须显式处理的冲突是 VP-003 的 `records/Schema CRUD` 与迁移 `0006
records_retire` 事实不一致。该事项成为 GOAL-005 R4-I003，需用户/规范决策
后才可进入 C2。
