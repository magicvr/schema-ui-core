---
doc_type: vision-roadmap
title: 愿景组合编排
status: active
created: 2026-07-31
updated: 2026-08-09
parent: null
version: 0.13.0
---

# 组合编排 · Schema UI Core Admin 基架

本文件索引已落盘的 VP 与用户确认的后续方向；它不是 Goal 路线图，也不汇总 progress%。

## 已落盘意图

| 顺序 | VP | 意图 | 前置 | 状态 |
|------|----|------|------|------|
| 1 | [VP-001-mvp-admin-foundation](plans/VP-001-mvp-admin-foundation.md) | 初始化 React + Go Admin MVP，覆盖固定协议来源、核心账号权限与协议范例验证。 | 无 | **closed**（2026-08-01；lead: workspace-001-mvp-admin-foundation） |
| 2 | [VP-002-production-admin-foundation](plans/VP-002-production-admin-foundation.md) | 在 I-PROTO-001 冻结子集之上，交付可直接 fork 使用的生产级 Schema 驱动 Admin 基架。 | 继承 VP-001 协议验证基线 | **closed**（2026-08-04；lead: workspace-002-production-admin-foundation） |
| 3 | [VP-003-modular-admin-architecture](plans/VP-003-modular-admin-architecture.md) | 单主线模块化单体：薄内核、模块契约、Fx、Profile、后端聚合 Manifest。 | 继承 VP-002；strategic re-align 见 VRev-006 | **closed**（2026-08-06；lead: workspace-003-modular-admin-architecture） |
| 4 | [VP-004-module-contribution-readiness](plans/VP-004-module-contribution-readiness.md) | 一方模块贡献 playbook 与 Core vs 模块归属方法论。 | 继承 VP-003 | **closed**（2026-08-06；lead: workspace-004-module-contribution-readiness） |
| 5 | [VP-006-full-protocol-contract-v2-7-0](plans/VP-006-full-protocol-contract-v2-7-0.md) | **`schema-ui-docs@v2.7.0` 整份契约**可验证兼容：覆盖表升版、Renderer/后端实现、范例与验证；纠正「长期停留在 MVP 子集」的组合焦点。 | 继承 VP-003/004；以 inventory + 上游 pin 为权威；`I-PROTO-001 v0.1.3` 仅作升版起点 | **closed**（2026-08-08 用户书面确认；lead: workspace-005-full-protocol-contract-v2-7-0；`I-PROTO-FULL-001` 12/12 include 冻结） |
| 6 | [VP-005-design-system-and-ui-experience](plans/VP-005-design-system-and-ui-experience.md) | Design Token、shadcn/ui 风格、Renderer/Shell 视觉与状态工效产品化。 | 继承 VP-003/004 + **VP-006 已 closed 的整份协议面**；VRev-011 `F-V018`/`F-V019`/`F-V020` → **fixed**（v0.3.0） | **active**（2026-08-09 用户确认激活；v0.4.0；lead 区待 `/govern` scaffold） |

## 组合门闩（用户 2026-08-08）

1. **协议优先于视觉**：在 VP-006 未 `closed` 前，**不得**将 VP-005 作为 `primary_plan` 推进实现，不得启动视觉优化波次。  
2. **MVP 子集不是终态成功条件**：`I-PROTO-001 v0.1.3` 是历史 MVP 冻结切片；整份 v2.7.0 契约由 VP-006 收口。  
3. 已关闭 VP-001～004 的历史证据与 status **不重写**。

## 已确认但尚未纳入新 VP 的后续方向

| 顺序 | 方向 | 与前序关系 | 建立 VP 前的约束 |
|------|------|------------|------------------|
| 7 | 订单、钱包、类目、通知等业务能力 | 默认承载：VP-003 架构 + VP-004 playbook + **VP-006 协议面** +（若已 closed）VP-005 设计系统 | 建 VP 前须 `/vision` 复核；不得用业务模块倒逼恢复长期双线或跳过协议覆盖 |

**当前交付意图**：**VP-005** **`active`**（2026-08-09 用户确认；设计系统与 Schema 驱动 UI/UX；lead 工作区待 `/govern`）。VP-001～004、**VP-006** 均 **closed**（覆盖权威 `I-PROTO-FULL-001`）。激活 ≠ 视觉产品化已交付。

## 单主线模块化策略

未来 fork 起点统一由同一代码主线、模块候选集与启动时 Profile 表达，权威见 [module-architecture.md](../architecture/module-architecture.md) 和 VP-003。原 [dual-track-contract.md](dual-track-contract.md) 已转为历史记录。
