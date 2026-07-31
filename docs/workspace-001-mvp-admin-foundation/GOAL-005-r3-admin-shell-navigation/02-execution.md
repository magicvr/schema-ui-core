---
id: GOAL-005-r3-admin-shell-navigation
doc: execution
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-07-31
version: 0.1.0
---

# 执行记录 · GOAL-005

## 时间线

### 2026-07-31 · R3 目标立项与范围登记

- `/govern` 复核当前显式工作区、Charter/VP 对齐、Root 路线图、R2 D-009/A-006 以及 R1 Web 边界。
- 创建本目标五件套和 `attachments/` 目录，并将 `GOAL-005-r3-admin-shell-navigation` 挂到 `GOAL-001-mvp-admin-foundation`。
- 将 R3 范围记录为 App manifest 装载、Admin shell、导航入口和路由语义；明确排除 R4 权限、R5 Renderer/业务范例及完整协议支持。
- 登记 `I-005-001` 至 `I-005-005` 为 required/open，记录方案冻结前的验证动作；当前没有把任何未知写成 verified。
- 同步工作区 `goal-tree.md`，Root `status` 仍为 `active`、`progress` 仍为 `2/6`。
- 本次没有修改 `apps/web`，没有产生 manifest loader、router、navigation 或 shell 的实现证据；父目标 `I-PROTO-002` / `I-PROTO-003` 未改变。

## 待办（计划 · 非完成事实）

1. 对照固定上游资料和 fixture，解析 `I-005-001` 至 `I-005-005`，补齐 manifest、导航和 shell 的未决契约。
2. 在 required 信息项可核对后，记录 R3 的 manifest 最小子集、路由映射、默认/fallback/active-route 和 shell 边界决策。
3. 按冻结方案在 `apps/web` 实施，并为无效 manifest、未知路由、fallback 和 active-route 建立可核对测试/运行时证据。
4. 执行阶段自审，响应全部相关意见，确认无开放 required finding 后再申请 R3 关门。

## 进度评估

R3 当前为 `active` 的规划阶段；目标五件套、范围、路线图和信息台账已建立，代码实现为未开始，`I-005-001` 至 `I-005-005` 均为 `open`。本次立项不改变 Root `progress: 2/6`，也不构成 R3 实施或验收完成。
