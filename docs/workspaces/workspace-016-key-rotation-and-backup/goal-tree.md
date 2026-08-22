---
title: 目标树 · workspace-016-key-rotation-and-backup
status: active
created: 2026-08-22
updated: 2026-08-22
parent: null
version: 0.2.0
workspace_id: workspace-016-key-rotation-and-backup
---

# 目标树 · 密钥轮换与备份恢复

> 工作区：`workspace-016-key-rotation-and-backup`
> canonical：`docs/workspaces/workspace-016-key-rotation-and-backup/`
> Root：`GOAL-001-key-rotation-and-backup`（**active** · 1/5）
> primary_plan：`VP-016-key-rotation-and-backup`（**active** · 架构 A5）

## 树

```text
GOAL-001-key-rotation-and-backup [active 1/5]   · 密钥轮换与备份恢复合同（JWT + 轮换后恢复）
├── GOAL-002-rotation-contract-freeze [done 4/4]   · R1 轮换合同冻结与配置面落地（A-001 self pass）
└── （R2 子目标待开：先关 I-003 决策，再实施双密钥验签；关门 independent/grok build）
```

R1 已关门：合同 D-002 冻结（I-001/I-002 verified）+ 配置面代码落地（`go test ./...` exit 0）。下一阶段 R2。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-key-rotation-and-backup | 密钥轮换与备份恢复合同（JWT + 轮换后恢复） | null | active | 1/5 | 2026-08-22 |
| GOAL-002-rotation-contract-freeze | R1 轮换合同冻结与配置面落地 | GOAL-001-key-rotation-and-backup | done | 4/4 | 2026-08-22 |
