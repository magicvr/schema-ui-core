---
id: D-001
title: 建立 R2 子目标与五项检查点
status: accepted
created: 2026-08-04
updated: 2026-08-04
parent: GOAL-003-r2-kernel-composition-root
version: 0.1.0
---

# D-001 · R2 阶段范围

## 决定

R2 采用 C1 exact Profile/I-004、C2 framework-agnostic kernel/module API、C3 graph/capability fail-closed、C4 migration/Manifest aggregation skeleton、C5 verification 五个等权检查点，依次推进；同一检查点内允许独立实现切片，但 Root R2 仍以本目标和审计结论为阶段边界。

## 依据

Root R1 已由 D-004 关闭并将 I-001/I-002/I-003/I-007 verified；VP-003 R2 要求薄内核、Fx 组合根、确定性图校验、全局迁移收集和 Manifest 聚合骨架。R2 不提前迁移一方业务模块或关闭 R3/R6 相关信息项。
