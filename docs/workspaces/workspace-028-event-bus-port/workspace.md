---
id: workspace-028-event-bus-port
title: 进程内事件总线端口工作区（架构分支 · 三端口第三个）
status: active
root_goal: GOAL-001-event-bus-port
canonical_scope: docs/workspaces/workspace-028-event-bus-port/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-028-event-bus-port
primary_plan: VP-028-event-bus-port
created: 2026-09-01
updated: 2026-09-01
version: 0.1.0
parent: null
---

# 工作区上下文 · 进程内事件总线端口

本工作区是 [VP-028-event-bus-port](../../vision/plans/VP-028-event-bus-port.md)（**`active`** v0.2.0 · 2026-09-01 用户指令激活）的唯一 lead delivery workspace。**架构分支**（H-002 同进程基座基础设施端口早期化 · Charter 0.4.0 成功边界 #6 · 承接 RT-Q02 运输端口前置）：交付进程内事件总线**运输端口**——EventBus 端口（类型化 Publish/Subscribe/Unsubscribe + 订阅生命周期 + 错误语义）+ 进程内 channel 实现 + outbox/MQ 接缝声明（不实现）。

- **Root** `GOAL-001-event-bus-port`：`active` · **1/4**（R1 ✅ 契约冻结 → R2 进程内实现 → R3 接缝与对齐 → R4 证据与关门），纲领见 Root `00-meta.md`。
- 激活门禁已满足（2026-09-01）：[VRev-064](../../vision/reviews/VRev-064-vp028-event-bus-port-activation.md) self `pass`（0 required；VRev-058/059 全部 findings 已闭合）；**架构类轻量 freshness PASS**（`5744868d` → `29727510`：协议 pin / 依赖锁 / 迁移台账 / Profile 装配 / provenance 五域零变更；区间代码全部为 VP-027 已审结目交付）不暂挂 `go`。
- 不改变 Charter `primary_workspace`（仍为 workspace-001）。
- **消费基线**：VP-003 模块契约（kernel Provider/Registrar）· 内核基础设施端口形态参照 Cache / RateLimiter / Store / ObjectStore / Mail（VP-026/027/013/014/017 先例）· 停机合同 VP-021（若选异步 channel 投递须声明 SIGTERM 取消订阅/排空）。
- **红线（激活即生效）**：不预制 outbox / 外部 broker（不引入客户端依赖 / 不预裁 RT-Q06 表结构 / **不消耗 RT-Q02 trigger**）；不改 Profile 默认集 / 模块矩阵 / Manifest 装配（VP-008 `go` 消费有效性）；**不解除** Admin 功能分支 typed domain event 扩展接缝的 trigger-gated；EventBus **不是** Job 端口替代（持久化/重试/定时工作仍走 Job）；不属 Redis 轨道（owner 仍为 `docs/architecture/cache-redis-seam-and-track.md`）；缓存/限流语义归 VP-026/027 独立交付。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-028-event-bus-port` | 与本区目标及资料引用的 `workspace_id` 一致 |
| Root Goal | `GOAL-001-event-bus-port` | `parent: null`；active · 0/4 |
| canonical 范围 | `docs/workspaces/workspace-028-event-bus-port/` | 本区唯一目标状态范围 |
| 共享资料目录 | `none` | 暂无固定共享资料 |
| 愿景角色 | `delivery` | VP-028 lead（active）；不改变 Charter primary workspace |
| 规划对齐 | `primary_plan` = `VP-028-event-bus-port`（`active` v0.2.0） | 2026-09-01 激活/开区（VRev-064 self `pass`；arch 类 freshness PASS `5744868d`→`29727510`） |

## 愿景对齐

Charter：`schema-ui-core-admin-foundation@0.4.0`（2026-08-31 strategic：同进程基座 · 成功边界 #6 · H-002）。
VP-028：进程内事件总线运输端口（vision_ref @0.4.0）——八条方向级退出判据 = 端口契约 / 进程内实现 / 接缝声明 / 对齐登记不解除 Admin gated / 共享约定 / 停机语义 / 边界保持 / required=0 闭合；红线 = 不预制 outbox/broker（不消耗 RT-Q02 trigger）、不改 Profile 默认集/Manifest、不解除 Admin typed domain event gated、EventBus ≠ Job。

## 纲领阶段（Root 路线图指针）

| 阶段 | 内容 | 状态 |
|------|------|------|
| R1 | **契约冻结**（判据 #1/#6 + I-028-001/002/003 裁决）：EventBus 端口 API 形态 · 类型化机制 · 投递语义（同步 vs 异步 + 缓冲满最小语义）· handler 错误语义 · 停机语义（异步须 SIGTERM 排空，否则同步） | **已关门**（2026-09-01 · GOAL-002 `done` 3/3：三信息项用户裁决 + D-002 v0.1.0 + kernel/eventbus.go + 双审 pass（A-001 self + A-002 grok independent · 0 required）） |
| R2 | **进程内实现**（判据 #2）：channel 分发 + 订阅管理 + 错误语义实现与测试（发布/订阅/退订、并发、顺序、handler panic 隔离） | 待 R1 |
| R3 | **接缝与对齐**（判据 #3/#4/#5）：outbox/MQ 运输接缝声明 + Admin typed domain event gated 对齐登记 + topic/订阅命名与契约测试 harness | 待 R2 |
| R4 | **证据与关门**（判据 #7/#8）：证据矩阵 / 越界核账 / 审计闭合 | 待 R1–R3 |

## 固定共享资料引用

> `shared-materials/index.json` 只能提供候选路径与摘要。缺完整引用字段的行无效。

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | — | none | — |
