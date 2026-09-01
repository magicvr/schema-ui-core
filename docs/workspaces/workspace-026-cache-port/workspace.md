---
id: workspace-026-cache-port
title: 通用缓存端口工作区（架构分支 · 三端口第一个）
status: active
root_goal: GOAL-001-cache-port
canonical_scope: docs/workspaces/workspace-026-cache-port/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-026-cache-port
primary_plan: VP-026-cache-port
created: 2026-08-31
updated: 2026-08-31
version: 0.1.0
parent: null
---

# 工作区上下文 · 通用缓存端口

本工作区是 [VP-026-cache-port](../../vision/plans/VP-026-cache-port.md)（**`active`** v0.2.0 · 2026-08-31 用户指令激活）的唯一 lead delivery workspace。**架构分支**（H-002 同进程基座基础设施端口早期化 · Charter 0.4.0 成功边界 #6 · 承接 RT-Q03）：交付通用缓存端口——Cache 端口（Get/Set/Delete + TTL + 命名空间 + 并发安全）+ 绝对/滑动过期 + 可插拔策略接口 + 内存供应商（默认）+ Redis 供应商接缝声明（不实现）。

- **Root** `GOAL-001-cache-port`：**`active`** · 0/4（R1 契约冻结 · R2 内存供应商 · R3 接缝与共享约定 · R4 证据与关门），纲领见 Root `00-meta.md`。
- 激活门禁已满足（2026-08-31）：[VRev-060](../../vision/reviews/VRev-060-vp026-cache-port-activation.md) self `pass`（0 required；VRev-058/059 全部 findings 已闭合）；**架构类轻量 freshness PASS**（`055da2fd` → `54fb57e7`：协议 pin / 依赖锁 / 迁移台账 / Profile 装配 / provenance 五域零变更；区间代码变更全部可追溯至 VP-025 已审结目）不暂挂 `go`。
- 不改变 Charter `primary_workspace`（仍为 workspace-001）。
- **消费基线**：VP-003 模块契约（kernel Provider/Registrar）· 内核基础设施端口形态参照 Store / ObjectStore / Mail（VP-013/014/017 先例）· 停机合同 VP-021（后台协程须声明 SIGTERM 排空）。
- **红线（激活即生效）**：不预制 Redis 实现（不引入客户端依赖 / **不消耗 RT-Q03 trigger**）；不改 Profile 默认集 / 模块矩阵 / Manifest 装配（VP-008 `go` 消费有效性）；Redis 轨道共享约定（key 前缀/命名空间/连接管理/测试 harness）登记于架构短文或 owner VP 决策（单一所有者，不跨区绑 Goal D-001）；限流/消息归 VP-027/028 独立交付。
- **本 VP 是三端口第一个**：Redis 轨道约定（VP-026/027）的 owner 登记义务在本区落地；VP-027/028 激活时继承约定。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-026-cache-port` | 与本区目标及资料引用的 `workspace_id` 一致 |
| Root Goal | `GOAL-001-cache-port` | `parent: null`；**active** · 0/4 |
| canonical 范围 | `docs/workspaces/workspace-026-cache-port/` | 本区唯一目标状态范围 |
| 共享资料目录 | `none` | 暂无固定共享资料 |
| 愿景角色 | `delivery` | VP-026 lead（active）；不改变 Charter primary workspace |
| 规划对齐 | `primary_plan` = `VP-026-cache-port`（`active` v0.2.0） | 2026-08-31 激活/开区（VRev-060 self `pass`；arch 类 freshness PASS `055da2fd`→`54fb57e7`） |

## 愿景对齐

Charter：`schema-ui-core-admin-foundation@0.4.0`（2026-08-31 strategic：同进程基座 · 成功边界 #6 · H-002）。
VP-026：通用缓存端口（vision_ref @0.4.0）——八条方向级退出判据 = 端口契约 / 双策略+可插拔 / 内存供应商 / Redis 接缝不引入客户端 / 共享约定单一所有者 / 停机语义 / 边界保持 / required=0 闭合；红线 = 不预制 Redis（不消耗 RT-Q03 trigger）、不改 Profile 默认集/Manifest。

## 纲领阶段（Root 路线图指针）

| 阶段 | 内容 | 状态 |
|------|------|------|
| R1 | **合同冻结**（判据 #1/#6 + I-026-001/002 裁决）：Cache 端口 API 形态（泛型 vs []byte vs 结构化）· TTL/清理语义（惰性 vs 后台协程）· 命名空间形态（I-026-003）· 策略接口形态 | 计划（未开工） |
| R2 | **内存供应商**（判据 #3）：有界容量 + TTL 清理（按 R1 裁决）+ 驱逐语义 + 并发安全 + 双策略实装（判据 #2） | 计划 |
| R3 | **接缝与共享约定**（判据 #4/#5）：Redis 供应商接缝声明（端口不变/连接管理约定/不引入客户端）+ Redis 轨道约定（VP-026/027）落地为 owner 文档 + mail `cachedAdapter` 迁移评估（I-026-004） | 计划 |
| R4 | **证据与关门**（判据 #7/#8）：证据矩阵 / 越界核账 / 审计闭合 | 计划 |

## 固定共享资料引用

> `shared-materials/index.json` 只能提供候选路径与摘要。缺完整引用字段的行无效。

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | — | none | — |