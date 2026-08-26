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
GOAL-001-graceful-shutdown-and-connection-drain [active 0/3]  · 优雅停机 / 连接排空合同
├── GOAL-002-r1-contract-freeze [pending]  · R1 合同冻结（停机顺序 / 超时与配置键 / Job 停机语义）
├── GOAL-003-r2-impl-and-test [pending]  · R2 实现与测试（HTTP drain / Job 语义 / 双方言 Store 排空）
└── GOAL-004-r3-evidence-closeout [pending]  · R3 证据与关门（信号→排空→退出码 harness · 双方言一致）
```

> 子目标将在纲领阶段启动时按 P-001 逐阶段创建（R1 阶段先立 GOAL-002；其编号与 slug 在立项时确定，上表为计划投影，非已创建目标）。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-graceful-shutdown-and-connection-drain | 优雅停机 / 连接排空合同 | null | active | 0/3 | 2026-08-27 |