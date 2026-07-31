---
id: GOAL-002-r1-repo-layout-conventions
doc: decision
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-07-31
version: 0.1.0
---

# 决策记录 · GOAL-002

## 信息需求与阶段门禁

本目标无独立开放 required 信息项。布局权威决策见父目标 [GOAL-001 `D-004`](../GOAL-001-mvp-admin-foundation/01-decision.md)。

## D-001 · 承接 Root D-004，本目标只落约定不写业务

**日期**：2026-07-31  
**状态**：accepted

**决定**：

1. 在本目标内把 monorepo 路径、包管理、边界写入可维护文档（优先扩展/新建 architecture 约定，并在根 README 给最短入口）。
2. 可创建空的 `apps/web`、`apps/api` 占位（含 `.gitkeep` 或最小 README），**不**在本目标实现 HTTP 业务或 UI 业务。
3. 可运行骨架分别由 GOAL-003 / GOAL-004 交付。

**为什么**：

- 分开「约定真相」与「可运行骨架」便于并行与验收。
- 避免在约定未落盘时两边 scaffold 路径漂移。

**未选方案**：

- **约定与双端 scaffold 挤进同一目标**：文档少，但失败面大、难并行关门。
