---
id: D-003
title: 记录 C1/C2 现状基线与目标边界
status: accepted
created: 2026-08-04
updated: 2026-08-04
parent: GOAL-002-r1-contract-migration-baseline
version: 0.1.0
---

# D-003 · C1/C2 现状基线与目标边界

## 决定

1. 将 [R1 C1 evidence](../attachments/r1-c1-module-profile-inventory.md) 作为当前模块、注册路径、Profile 候选和依赖闭包的 baseline；其中 `present`、`candidate`、`absent` 三类不得混写。
2. 将 [R1 C2 evidence](../attachments/r1-c2-migration-seed-boundary.md) 作为迁移链、seed、ledger/checksum、快照/恢复和未实现回滚/tombstone/reconcile 边界的 baseline。
3. 将 C1、C2 标记为本子目标已完成的证据收集检查点，子目标进度推进至 `2/4`。Root I-001、I-002 仍保持 `open`，不把子目标检查点完成改写为 Root 信息验证或 R1 放行。

## 取舍与理由

- 当前集中式 `handler.Register`、静态 Manifest 消费和 Store 迁移链足以提供可追踪事实，但没有据此虚构 Fx 模块 Registry、Shell 实现、Profile 配置、tombstone 或独立 reconcile。
- `admin.users`、`admin.roles`、`admin.settings`、`admin.activity` 只能作为候选模块；I-004 要求的精确 Profile 集和 precedence 仍由后续 R2 冻结。
- 迁移前 snapshot 是恢复后备，不等于应用层 rollback。`seedRBAC` 的幂等 ensure 也不等于目标要求的 system-data reconcile；这两个缺口保留给 C3/R2/R4 的决策与实现验证。

## 约束

本决定不修改 Root `00-meta.md` 的 required 信息状态，不创建 R2 子目标，不改变 I-PROTO-001 v0.1.3 的范围。
