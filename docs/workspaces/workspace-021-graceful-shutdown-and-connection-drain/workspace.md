---
id: workspace-021-graceful-shutdown-and-connection-drain
title: 优雅停机 / 连接排空合同工作区
status: active
root_goal: GOAL-001-graceful-shutdown-and-connection-drain
canonical_scope: docs/workspaces/workspace-021-graceful-shutdown-and-connection-drain/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-021-graceful-shutdown-and-connection-drain
primary_plan: VP-021-graceful-shutdown-and-connection-drain
created: 2026-08-27
updated: 2026-08-27
version: 0.1.0
parent: null
---

# 工作区上下文 · 优雅停机 / 连接排空合同

本工作区是 [VP-021-graceful-shutdown-and-connection-drain](../../vision/plans/VP-021-graceful-shutdown-and-connection-drain.md)（**`active`** v0.2.0 · 2026-08-27 激活 · 架构分支 RT-D02）的唯一 lead delivery workspace。

- **Root** `GOAL-001-graceful-shutdown-and-connection-drain`：**`active`** · 0/3（纲领 R1 合同冻结 → R2 实现与测试 → R3 证据与关门）。
- 激活门禁已满足（2026-08-27）：VRev-046（self）`pass`（0 required；V-F081/V-F082 → 激活事务内 fixed）；**架构类 freshness PASS**（`ed99e88` → `fddaf638`，不暂挂 `go`）；VP-009/VP-010 无开放阻断。
- 不改变 Charter `primary_workspace`（仍为 workspace-001）。
- 消费已交付基架：VP-012 Job 六态（closed 2026-08-19）、VP-013 双方言 Store（closed 2026-08-21，A1）、VP-015 结构化日志 / correlation（closed 2026-08-22，A4）、VP-002 Compose 一键（RT-D01 delivered）。
- 边界（与 VP 冻结一致）：**不进** A3 余项（多实例、就绪探针扩依赖、PG 锁 vs Redis vs 队列评估仍 trigger-gated）、RT-D03 进程分离、RT-Q04 Job 租约 / leader election、RT-Q02 外部队列、K8s/Helm/Operator、TLS 终止（RT-D05）；不改 Profile 默认集 / 模块矩阵 / Manifest 装配语义。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-021-graceful-shutdown-and-connection-drain` | 与本区目标及资料引用的 `workspace_id` 一致 |
| Root Goal | `GOAL-001-graceful-shutdown-and-connection-drain` | `parent: null`；**active** · 0/3（2026-08-27 开区） |
| canonical 范围 | `docs/workspaces/workspace-021-graceful-shutdown-and-connection-drain/` | 本区唯一目标状态范围 |
| 共享资料目录 | `none` | 暂无固定共享资料 |
| 愿景角色 | `delivery` | VP-021 lead（active）；不改变 Charter primary workspace |
| 规划对齐 | `primary_plan` = `VP-021-graceful-shutdown-and-connection-drain`（`active` v0.2.0） | 2026-08-27 激活/开区（VRev-046 self `pass`；架构类 freshness PASS） |

## 愿景对齐

Charter：`schema-ui-core-admin-foundation@0.2.0`。
VP-021：优雅停机 / 连接排空合同（架构分支 RT-D02，单进程 + Compose 基线）——`active`（v0.2.0，2026-08-27 激活）：把现行「进程生命周期有、但无明确 drain 合同」的后端收成可核对的停机顺序 / HTTP drain / 运行中 Job 语义 / 双方言 Store 排空合同。与 VP-009/VP-010 正交：停机相关安全/符合性 gap 归持续程序。

## 纲领阶段（Root 路线图指针）

| 阶段 | 内容 | 状态 |
|------|------|------|
| R1 | **合同冻结**：停机顺序 / grace 与超时默认值与配置键（I-002）/ 运行中 Job 停机语义（I-001，等完成 vs 中断标记重跑）；Store 排空与迁移窗口重叠语义（I-003 最晚 R2） | 待立项（GOAL-002；R1 方案冻结前关闭 I-001/I-002） |
| R2 | **实现与测试**：`http.Server` Shutdown 合同化、Job 停机行为实现、双方言连接关闭顺序 + 迁移窗口重叠语义 | 依赖 R1 |
| R3 | **证据与关门**：SIGTERM/SIGINT → 排空 → 退出码可核对（信号测试或等价 harness，单进程 + Compose）；双方言一致性；开放 required = 0 | 依赖 R2 |

## 固定共享资料引用

> `shared-materials/index.json` 只能提供候选路径与摘要。缺完整引用字段的行无效。

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | — | none | — |