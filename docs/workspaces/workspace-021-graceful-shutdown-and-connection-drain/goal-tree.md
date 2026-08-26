---
title: 目标树 · workspace-021-graceful-shutdown-and-connection-drain
status: active
created: 2026-08-27
updated: 2026-08-27
parent: null
version: 0.1.0
workspace_id: workspace-021-graceful-shutdown-and-connection-drain
---

# 目标树 · 优雅停机 / 连接排空合同

> 工作区：`workspace-021-graceful-shutdown-and-connection-drain`（**active** · 2026-08-27 开区）
> canonical：`docs/workspaces/workspace-021-graceful-shutdown-and-connection-drain/`
> Root：`GOAL-001-graceful-shutdown-and-connection-drain`（**active** · 0/3）
> primary_plan：`VP-021-graceful-shutdown-and-connection-drain`（**active** v0.2.0 · VRev-046 self pass + 架构类 freshness PASS）

## 树

```text
GOAL-001-graceful-shutdown-and-connection-drain [active 2/3]  · 优雅停机 / 连接排空合同
└── GOAL-002-r1-contract-freeze [done 3/3]  · R1 合同冻结（停机顺序 / 超时与配置键 / Job 停机语义 / Store 排空）
    （R2 / R3 阶段待立项）
```

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-graceful-shutdown-and-connection-drain | 优雅停机 / 连接排空合同 | null | active | 2/3 | 2026-08-27 |
| GOAL-002-r1-contract-freeze | R1 合同冻结（停机顺序 / 超时与配置键 / Job 语义 / Store 排空） | GOAL-001-graceful-shutdown-and-connection-drain | done | 3/3 | 2026-08-27 |