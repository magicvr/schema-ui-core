---
id: workspace-027-rate-limiter-port
title: 通用限流器端口工作区（架构分支 · 三端口第二个）
status: active
root_goal: GOAL-001-rate-limiter-port
canonical_scope: docs/workspaces/workspace-027-rate-limiter-port/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-027-rate-limiter-port
primary_plan: VP-027-rate-limiter-port
created: 2026-09-01
updated: 2026-09-01
version: 0.1.0
parent: null
---

# 工作区上下文 · 通用限流器端口

本工作区是 [VP-027-rate-limiter-port](../../vision/plans/VP-027-rate-limiter-port.md)（**`active`** v0.2.0 · 2026-09-01 用户指令激活）的唯一 lead delivery workspace。**架构分支**（H-002 同进程基座基础设施端口早期化 · Charter 0.4.0 成功边界 #6 · 承接 RT-Q05）：交付通用限流器端口——RateLimiter 端口（Allow/Record/Reset/RetryAfter + key 寻址 + 供应商无关）+ 滑动窗口内存供应商（演进既有 `loginRateLimiter`）+ 7 处使用点完整迁移 + Redis 供应商接缝声明（不实现）。

- **Root** `GOAL-001-rate-limiter-port`：**`active`** · **0/4**（R1 合同冻结 → R2 内存供应商+使用点迁移 → R3 接缝与共享约定 → R4 证据与关门），纲领见 Root `00-meta.md`。
- 激活门禁已满足（2026-09-01）：[VRev-062](../../vision/reviews/VRev-062-vp027-rate-limiter-port-activation.md) self `pass`（0 required；VRev-058/059 全部 findings 已闭合）；**架构类轻量 freshness PASS**（`54fb57e7` → `5744868d`：协议 pin / 依赖锁 / 迁移台账 / Profile 装配 / provenance 五域零变更；区间代码全部为 VP-026 已审结目交付）不暂挂 `go`。
- 不改变 Charter `primary_workspace`（仍为 workspace-001）。
- **消费基线**：VP-003 模块契约（kernel Provider/Registrar）· 内核基础设施端口形态参照 Cache / Store / ObjectStore / Mail（VP-026/013/014/017 先例）· 停机合同 VP-021（后台协程须声明 SIGTERM 排空）· 使用点迁移须保持 D-001 P1 防暴破防护与 W12 D-002 窗口常量；GOAL-014 账号分层锁定（DB 行锁）显式排除，不纳入端口。
- **红线（激活即生效）**：不预制 Redis 实现（不引入客户端依赖 / **不消耗 RT-Q05 trigger**）；不改 Profile 默认集 / 模块矩阵 / Manifest 装配（VP-008 `go` 消费有效性）；Redis 轨道共享约定（key 前缀/命名空间/连接管理/测试 harness）继承 owner 文档 `docs/architecture/cache-redis-seam-and-track.md`（单一所有者，不跨区绑 Goal D-001）；缓存/事件语义归 VP-026/028 独立交付。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-027-rate-limiter-port` | 与本区目标及资料引用的 `workspace_id` 一致 |
| Root Goal | `GOAL-001-rate-limiter-port` | `parent: null`；**active** · 0/4 |
| canonical 范围 | `docs/workspaces/workspace-027-rate-limiter-port/` | 本区唯一目标状态范围 |
| 共享资料目录 | `none` | 暂无固定共享资料 |
| 愿景角色 | `delivery` | VP-027 lead（active）；不改变 Charter primary workspace |
| 规划对齐 | `primary_plan` = `VP-027-rate-limiter-port`（`active` v0.2.0） | 2026-09-01 激活/开区（VRev-062 self `pass`；arch 类 freshness PASS `54fb57e7`→`5744868d`） |

## 愿景对齐

Charter：`schema-ui-core-admin-foundation@0.4.0`（2026-08-31 strategic：同进程基座 · 成功边界 #6 · H-002）。
VP-027：通用限流器端口（vision_ref @0.4.0）——七条方向级退出判据 = 端口契约 / 内存供应商 / 使用点迁移不回归（7 处构造点）/ Redis 接缝不引入客户端 / 共享约定单一所有者 / 边界保持 / required=0 闭合；红线 = 不预制 Redis（不消耗 RT-Q05 trigger）、不改 Profile 默认集/Manifest、W12 D-002 窗口常量保持。

## 纲领阶段（Root 路线图指针）

| 阶段 | 内容 | 状态 |
|------|------|------|
| R1 | **合同冻结**（判据 #1 + I-027-001/003/004 裁决）：RateLimiter 端口 API 形态 · RetryAfter 语义 · 窗口语义（滑动窗口保持）· key 维度 · 供应商无关面 | **待启动**（I-027-001 required 前置裁决） |
| R2 | **内存供应商 + 使用点迁移**（判据 #2/#3 + I-027-002）：演进 `loginRateLimiter` + 7 处构造点接入 + 回归（D-001 P1 / W12 D-002 常量保持） | **待启动** |
| R3 | **接缝与共享约定**（判据 #4/#5）：Redis 接缝声明 + 轨道约定继承登记（owner = cache-redis-seam-and-track.md） | **待启动** |
| R4 | **证据与关门**（判据 #6/#7）：证据矩阵 / 越界核账 / 审计闭合 | **待启动** |

## 固定共享资料引用

> `shared-materials/index.json` 只能提供候选路径与摘要。缺完整引用字段的行无效。

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | — | none | — |