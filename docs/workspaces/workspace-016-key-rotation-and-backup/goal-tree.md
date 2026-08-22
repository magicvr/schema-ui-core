---
title: 目标树 · workspace-016-key-rotation-and-backup
status: active
created: 2026-08-22
updated: 2026-08-22
parent: null
version: 0.3.0
workspace_id: workspace-016-key-rotation-and-backup
---

# 目标树 · 密钥轮换与备份恢复

> 工作区：`workspace-016-key-rotation-and-backup`
> canonical：`docs/workspaces/workspace-016-key-rotation-and-backup/`
> Root：`GOAL-001-key-rotation-and-backup`（**active** · 2/5）
> primary_plan：`VP-016-key-rotation-and-backup`（**active** · 架构 A5）

## 树

```text
GOAL-001-key-rotation-and-backup [active 2/5]   · 密钥轮换与备份恢复合同（JWT + 轮换后恢复）
├── GOAL-002-rotation-contract-freeze [done 4/4]   · R1 轮换合同冻结与配置面落地（A-001 self pass）
├── GOAL-003-dual-key-jwt [done 4/4]               · R2 JWT 双密钥实现（A-001 self + A-002 independent 双 pass；F-001/2/3 fixed）
└── （R3 子目标待开：先关 I-004 恢复剧本决策，再两方言取证）
```

R1、R2 已关门。R2 审计：self + independent（grok build）双 pass，recommended findings 全部 fixed。下一阶段 R3（轮换后恢复证据，依赖 R2 ✓）。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-key-rotation-and-backup | 密钥轮换与备份恢复合同（JWT + 轮换后恢复） | null | active | 2/5 | 2026-08-22 |
| GOAL-002-rotation-contract-freeze | R1 轮换合同冻结与配置面落地 | GOAL-001-key-rotation-and-backup | done | 4/4 | 2026-08-22 |
| GOAL-003-dual-key-jwt | R2 JWT 双密钥实现（重叠窗验签） | GOAL-001-key-rotation-and-backup | done | 4/4 | 2026-08-22 |

