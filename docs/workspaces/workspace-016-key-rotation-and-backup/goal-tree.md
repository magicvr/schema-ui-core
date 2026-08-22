---
title: 目标树 · workspace-016-key-rotation-and-backup
status: done
created: 2026-08-22
updated: 2026-08-22
parent: null
version: 0.6.0
workspace_id: workspace-016-key-rotation-and-backup
---

# 目标树 · 密钥轮换与备份恢复

> 工作区：`workspace-016-key-rotation-and-backup`
> canonical：`docs/workspaces/workspace-016-key-rotation-and-backup/`
> Root：`GOAL-001-key-rotation-and-backup`（**done** · 5/5）
> primary_plan：`VP-016-key-rotation-and-backup`（`active` → 关门收尾走 `/vision`）

## 树

```text
GOAL-001-key-rotation-and-backup [done 5/5]     · 密钥轮换与备份恢复合同（JWT + 轮换后恢复）
├── GOAL-002-rotation-contract-freeze [done 4/4]   · R1 轮换合同冻结与配置面落地（A-001 self pass）
├── GOAL-003-dual-key-jwt [done 4/4]               · R2 JWT 双密钥实现（A-001 self + A-002 independent 双 pass；recommended 全 fixed）
├── GOAL-004-r3-recovery-evidence [done 4/4]       · R3 轮换后恢复证据（SQLite + PG 双循环全绿；A-001 self pass）
├── GOAL-005-default-single-key [done 3/3]         · R4 默认单密钥仍可用（6/6 判据面实跑成立；A-001 self pass）
└── GOAL-006-r5-dual-path-evidence [done 4/4]      · R5 双路径证据 + Root 关门驱动（A-002 F-001 required + recommended 全 fixed）
```

Root 关门：R1～R5 全部完成；Root 关门审计 = A-001 self pass + A-002 independent conditional，唯一 required F-001 与四条 recommended 均已按 `fixed` 路径闭合（见 Root `03-audit.md` 响应节）。开放 required finding = 0。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-key-rotation-and-backup | 密钥轮换与备份恢复合同（JWT + 轮换后恢复） | null | done | 5/5 | 2026-08-22 |
| GOAL-002-rotation-contract-freeze | R1 轮换合同冻结与配置面落地 | GOAL-001-key-rotation-and-backup | done | 4/4 | 2026-08-22 |
| GOAL-003-dual-key-jwt | R2 JWT 双密钥实现（重叠窗验签） | GOAL-001-key-rotation-and-backup | done | 4/4 | 2026-08-22 |
| GOAL-004-r3-recovery-evidence | R3 轮换后恢复证据（SQLite + PG） | GOAL-001-key-rotation-and-backup | done | 4/4 | 2026-08-22 |
| GOAL-005-default-single-key | R4 默认单密钥仍可用（证据整合） | GOAL-001-key-rotation-and-backup | done | 3/3 | 2026-08-22 |
| GOAL-006-r5-dual-path-evidence | R5 显式双密钥双路径证据与 Root 关门 | GOAL-001-key-rotation-and-backup | done | 4/4 | 2026-08-22 |
