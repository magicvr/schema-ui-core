---
doc_type: vision-roadmap
title: 愿景组合编排
status: active
created: 2026-07-31
updated: 2026-08-06
parent: null
version: 0.7.2
---

# 组合编排 · Schema UI Core Admin 基架

本文件索引已落盘的 VP 与用户确认的后续方向；它不是 Goal 路线图，也不汇总 progress%。

## 已落盘意图

| 顺序 | VP | 意图 | 前置 | 状态 |
|------|----|------|------|------|
| 1 | [VP-001-mvp-admin-foundation](plans/VP-001-mvp-admin-foundation.md) | 初始化 React + Go Admin MVP，覆盖固定协议来源、核心账号权限与协议范例验证。 | 无 | **closed**（2026-08-01；lead: workspace-001-mvp-admin-foundation；三条退出判据经 R6 工作区 Q2 证据满足，用户确认关门） |
| 2 | [VP-002-production-admin-foundation](plans/VP-002-production-admin-foundation.md) | 在 I-PROTO-001 冻结子集之上，交付可直接 fork 使用的生产级 Schema 驱动 Admin 基架：Renderer、真实认证、持久化权限、CRUD 与工程化启动。 | 继承 VP-001 协议验证基线 | **closed**（2026-08-04；lead: workspace-002-production-admin-foundation；七条产品级成功标准经工作区 Q2 证据满足——Root `GOAL-001` `done / 5/5`（A-007 self close-out `pass`）、GOAL-002～013 全部 `done`、Root 03-audit 开放 required=0——用户确认关门） |
| 3 | [VP-003-modular-admin-architecture](plans/VP-003-modular-admin-architecture.md) | 将生产级 Admin 基架收敛为单主线模块化单体：薄内核、框架无关模块契约、Fx 组合根、Profile、后端聚合 Manifest 与完整数据/运维闭环。 | 继承 VP-002 产品基线；strategic re-align 已由 VRev-006 核对 | **closed**（2026-08-06；lead: workspace-003-modular-admin-architecture；七条方向级退出判据经工作区 Q2 证据满足——Root `GOAL-001` `done / 6/6`（A-018 self + A-019 independent + A-020 response）、A-021 独立动态代码复审 `pass`、Root 03-audit 开放 required=0、Vision Review 0 open required——用户确认关门；有界 residual R4-I004 点名 workspace-003/GOAL-006） |

## 已确认但尚未纳入新 VP 的后续方向

| 顺序 | 方向 | 与前序关系 | 建立 VP 前的约束 |
|------|------|------------|------------------|
| 4 | 订单、钱包、类目、通知等业务能力 | 以 VP-003 的单主线模块架构为默认技术承载边界。 | 当前 Charter 不把具体业务产品列为成功条件；建 VP 前须由 `/vision` 复核是否需要 strategic 修订，并明确独立退出判据。不得用业务模块倒逼恢复长期双线。 |

VP-003 是下一个明确 VP。剩余业务能力在建立对应 `VP-00N-*.md` 前不是可引用的 `primary_plan`，也不属于 VP-003 的架构退出证据。

## 单主线模块化策略

未来 fork 起点统一由同一代码主线、模块候选集与启动时 Profile 表达，权威见 [module-architecture.md](../architecture/module-architecture.md) 和 VP-003。原 [dual-track-contract.md](dual-track-contract.md) 已转为历史记录；本次只替换未来意图，不宣称执行过不存在的 MVP/Admin 分支删除或合并。
