---
id: workspace-032-rate-limiter-atomic-port
title: 限流器端口原子化工作区（架构分支 · AllowRecord）
status: active
root_goal: GOAL-001-rate-limiter-atomic-port
canonical_scope: docs/workspaces/workspace-032-rate-limiter-atomic-port/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-032-rate-limiter-atomic-port
primary_plan: VP-032-rate-limiter-atomic-port
created: 2026-09-03
updated: 2026-09-03
version: 0.3.0
parent: null
---

# 工作区上下文 · 限流器端口原子化

本工作区是 [VP-032-rate-limiter-atomic-port](../../vision/plans/VP-032-rate-limiter-atomic-port.md)（**`active`** v0.2.0 · 2026-09-03 用户指令激活）的唯一 lead delivery workspace。**架构分支**（承接 VP-027 residual R-007 · 不重开 VP-027）：在 `kernel.RateLimiter` 新增原子 `AllowRecord`，迁移冻结 14 处 Allow→Record 使用点，消除 TOCTOU；内存供应商实现；Redis 仍 RT-Q05 trigger-gated。

- **Root** `GOAL-001-rate-limiter-atomic-port`：**`active`** · **1/3**（R1 合同落盘已关门 → R2 14 处生产迁移进行中 → R3 证据与关门），纲领见 Root `00-meta.md`。
- 激活门禁已满足（2026-09-03）：[VRev-073](../../vision/reviews/VRev-073-vp032-rate-limiter-atomic-port-activation.md) self `pass`（0 required；I-032-001/002 已冻结）；**架构类轻量 freshness PASS**（`42036a3c` → `b1c03acd`：协议 pin / 依赖锁 / Profile 默认集 / provenance 零变更；区间代码 = VP-030 已审结目）不暂挂 `go`。
- 不改变 Charter `primary_workspace`（仍为 workspace-001）。
- **消费基线**：VP-027 已 closed 的端口（Allow/Record/Clear/RetryAfterSeconds + 内存供应商 + Redis 接缝声明）· VP-030 三桶限流调用点（分母 #12–#14）· 停机合同 VP-021（本波不新增后台协程）。
- **红线（激活即生效）**：不重开 VP-027；不实现 Redis / 不消耗 RT-Q05 trigger；不改 Profile 默认集 / 模块矩阵 / Manifest（VP-008 `go`）；不改其它内核端口；`Allow`/`Record` 保留兼容。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-032-rate-limiter-atomic-port` | 与本区目标及资料引用的 `workspace_id` 一致 |
| Root Goal | `GOAL-001-rate-limiter-atomic-port` | `parent: null`；**active** · 1/3 |
| canonical 范围 | `docs/workspaces/workspace-032-rate-limiter-atomic-port/` | 本区唯一目标状态范围 |
| 共享资料目录 | `none` | 暂无固定共享资料 |
| 愿景角色 | `delivery` | VP-032 lead（active）；不改变 Charter primary workspace |
| 规划对齐 | `primary_plan` = `VP-032-rate-limiter-atomic-port`（`active` v0.2.0） | 2026-09-03 激活/开区（VRev-073 self `pass`；arch 类 freshness PASS `42036a3c`→`b1c03acd`） |

## 愿景对齐

Charter：`schema-ui-core-admin-foundation@0.4.0`（2026-08-31 strategic：同进程基座 · 成功边界 #6 · H-002）。
VP-032：限流器端口原子化（vision_ref @0.4.0）——五条方向级退出判据 = 原子性 / 行为等价（14 处）/ 兼容 / 边界保持 / required=0 闭合；红线 = 不重开 VP-027、不实现 Redis、不改 Profile 默认集。

## 纲领阶段（Root 路线图指针）

| 阶段 | 内容 | 状态 |
|------|------|------|
| R1 | **合同落盘**（GOAL-002）：D-002 冻结 + kernel.AllowRecord + Memory 单锁实现 + 合同级测试 | **已关门**（GOAL-002 done · A-003 关门） |
| R2 | **14 处迁移 + handler 回归**（GOAL-003 · 判据 2/3）：按立即消费 / 失败预算两口径迁生产调用点 | **进行中**（GOAL-003 已立项） |
| R3 | **证据与关门**（判据 4/5）：证据矩阵 / 越界核账 / 审计闭合 | 待 R2 关门 |

## 固定共享资料引用

> `shared-materials/index.json` 只能提供候选路径与摘要。缺完整引用字段的行无效。

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | — | none | — |
