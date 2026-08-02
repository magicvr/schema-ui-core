---
title: 执行记录 · R3 · 持久化 RBAC、菜单投影与版本迁移
status: active
created: 2026-08-02
updated: 2026-08-02
parent: GOAL-001-production-admin-foundation
version: 0.1.0
---

# 执行记录 · GOAL-006

## 2026-08-02 · 目标立项

- 用户在 Root R3 信息取舍中书面确认推荐方案 B、`features` 菜单投影、两步迁移、读写权限边界和恢复证据口径；Root 记录 D-009 并将 `I-003` 置为 `verified`。
- 从核心五件套约定建立本目标，设定 `parent: GOAL-001-production-admin-foundation`、`status: active` 与六个顺序成功检查点；同步更新工作区 `goal-tree.md`。
- 记录 D-001，选择一个端到端目标承载 R3 强耦合闭环。
- 登记 `I-006-001/002` 两个 required 实施细化项；当前均为 `collecting`，尚未到期但会阻断各自列明的实现门禁。
- **未做**：没有产品代码、数据库、API、Web manifest 或测试行为变更；当前进度为 `0/6`。

## 待办（计划，不是事实）

1. 形成版本化 DDL、约束、迁移编号与 seed key 计划，关闭 `I-006-001`。
2. 选择首个真实 `page_ref` / `feature_key` 及 admin/viewer 矩阵，关闭 `I-006-002`。
3. 按 S1 → S6 顺序实施并在每个检查点记录可复现证据。
