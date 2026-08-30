---
title: 目标树 · workspace-021-graceful-shutdown-and-connection-drain
status: done
created: 2026-08-27
updated: 2026-08-30
parent: null
version: 0.2.1
workspace_id: workspace-021-graceful-shutdown-and-connection-drain
---

# 目标树 · 优雅停机 / 连接排空合同（已结项）

> 工作区：`workspace-021-graceful-shutdown-and-connection-drain`（**done** · 2026-08-27 结项）
> canonical：`docs/workspaces/workspace-021-graceful-shutdown-and-connection-drain/`
> Root：`GOAL-001-graceful-shutdown-and-connection-drain`（**done** · 3/3）
> primary_plan：`VP-021-graceful-shutdown-and-connection-drain`（closed v0.3.0 · 2026-08-27 用户指令授权 · VRev-047）

## 树

```text
GOAL-001-graceful-shutdown-and-connection-drain [done 3/3]  · 优雅停机 / 连接排空合同
├── GOAL-002-r1-contract-freeze [done 3/3]  · R1 合同冻结（停机顺序 / 超时与配置键 / Job 停机语义 / Store 排空）
├── GOAL-003-r2-impl-and-test [done 3/3]  · R2 实现与测试（shutdown_timeout 配置键 / main 接线 / compose 对齐 / 测试锁）
└── GOAL-004-r3-evidence-closeout [done 3/3]  · R3 证据与关门（合同 §8 harness · 双方言 · 双审）
```

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-graceful-shutdown-and-connection-drain | 优雅停机 / 连接排空合同 | null | done | 3/3 | 2026-08-27 |
| GOAL-002-r1-contract-freeze | R1 合同冻结（停机顺序 / 超时与配置键 / Job 语义 / Store 排空） | GOAL-001-graceful-shutdown-and-connection-drain | done | 3/3 | 2026-08-27 |
| GOAL-003-r2-impl-and-test | R2 实现与测试（shutdown_timeout 配置键 / main 接线 / compose 对齐 / 测试锁） | GOAL-001-graceful-shutdown-and-connection-drain | done | 3/3 | 2026-08-27 |
| GOAL-004-r3-evidence-closeout | R3 证据与关门（合同 §8 harness · 双方言 · Root 关门审计） | GOAL-001-graceful-shutdown-and-connection-drain | done | 3/3 | 2026-08-27 |