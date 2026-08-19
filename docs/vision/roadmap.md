---
doc_type: vision-roadmap
title: 愿景组合编排
status: active
created: 2026-07-31
updated: 2026-08-19
parent: null
version: 0.23.0
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
| 5 | [VP-006-full-protocol-contract-v2-7-0](plans/VP-006-full-protocol-contract-v2-7-0.md) | **`schema-ui-docs@v2.7.0` 整份契约**可验证兼容：覆盖表升版、Renderer/后端实现、范例与验证；纠正「长期停留在 MVP 子集」的组合焦点。 | 继承 VP-003/004；以 inventory + 上游 pin 为权威；`I-PROTO-001 v0.1.3` 仅作升版起点 | **closed**（2026-08-08 用户书面确认；lead: workspace-005-full-protocol-contract-v2-7-0；`I-PROTO-FULL-001` v1.0.1 = 12/12 include、318 executed + 2 local adapter excluded） |
| 6 | [VP-005-design-system-and-ui-experience](plans/VP-005-design-system-and-ui-experience.md) | Design Token、shadcn/ui 风格、Renderer/Shell 视觉与状态工效产品化。 | 继承 VP-003/004 + **VP-006 已 closed 的整份协议面**；VRev-011 `F-V018`/`F-V019`/`F-V020` → **fixed**（v0.3.0） | **closed**（2026-08-09 用户书面确认；v0.5.0；lead: `workspace-006-design-system-and-ui-experience`；Root `GOAL-001-design-system-and-ui-experience` `done 5/5`） |
| 7 | [VP-007-localization-and-system-settings](plans/VP-007-localization-and-system-settings.md) | 建立 `zh-CN` / `en-US` 多语种运行时与 `auto` 解析，并把既有 Settings 产品化为 General / Branding / Localization / Appearance 四类系统设置。 | 继承 VP-003/004 模块边界、VP-005 设计系统与 VP-006 完整协议面；不改变双 Profile 的 Settings 可见性边界 | **closed**（2026-08-09 用户书面确认；lead: `workspace-007-localization-and-system-settings`，Root done 6/6） |
| 8 | [VP-008-admin-module-readiness-and-foundation-convergence](plans/VP-008-admin-module-readiness-and-foundation-convergence.md) | 在正式业务模块开发前，对当前代码主线执行全基架准入：现状扫描、代码/功能/治理缺口、UI 协议判断、阻断整改与 `go`/`no-go`。 | 继承 VP-003/004 模块架构与贡献契约、VP-005 设计系统、VP-006 完整协议面、VP-007 locale/settings；不重开历史 VP | **closed**（2026-08-10 用户书面确认；候选 `ed99e88` clean，S0–S5 完成、open required = 0、`go` 签发；lead: workspace-008-admin-module-readiness，Root `GOAL-001-admin-module-readiness` done 6/6） |
| 9 | [VP-009-production-hardening](plans/VP-009-production-hardening.md) | 生产加固：**共享基架持续安全与健壮性程序**（周期扫描、波次修复、与 VP-008 `go` 消费有效性接口）；具体 finding 由工作区波次子目标承接。 | 继承 VP-003/004/005/006/007 + **VP-008 `go` 消费有效性**；共享基架缺陷可暂挂/恢复 `go` | **active**（2026-08-10 语义纠正为长期程序；曾误 `closed` 已撤销；lead: workspace-009-production-hardening；Root **active** 程序容器；波次 W1–W4 与 W6 均 done，W5 扫描 0 中高危未开子目标） |
| 10 | [VP-010-design-implementation-conformance](plans/VP-010-design-implementation-conformance.md) | 设计意图与实现符合性：**共享基架持续对齐程序**（周期对照 as-designed / as-built、conformance gap 分流、波次整改、与 VP-008 `go` 消费有效性接口、与 VP-009 正交）；具体 gap 由工作区波次子目标承接。 | 继承 VP-003/004/005/006/007/008 + **VP-008 `go` 消费有效性**；与 **VP-009** 正交（安全 vs 符合性） | **active**（2026-08-11 用户确认类 VP-009 长期程序；lead: workspace-010-design-implementation-conformance；Root **active** 程序容器；波次 W1–W13 均 done，`go` 均无新暂挂） |
| 11 | [VP-011-admin-functional-modules](plans/VP-011-admin-functional-modules.md) | 标准 Admin 功能模块（通用模块 + 常用业务领域）分档交付：有界调研 → 三档分档 → 分波实现；一等公民 / 常用 / 增补。 | 继承 VP-008 `go` 消费有效性（freshness review **PASS**，候选 `f14ab9d`）+ VP-009/010 无开放阻断；VP-001～008 已固化协议/架构/设计/locale 基线 | **closed**（2026-08-18 有界关门；lead: workspace-011-admin-functional-modules；Root done；四档能力地图上提至本 roadmap） |
| 12 | [VP-012-shared-cross-module-contracts](plans/VP-012-shared-cross-module-contracts.md) | 共享横切契约与平台基架：correlation、审计模型、并发/幂等、异步 Job、maintenance 门控、API Token；不承载业务领域。 | 继承 VP-011 的 R5 四档能力地图；与 VP-009/VP-010 正交分流；不改变 Charter 边界 | **closed**（2026-08-19 完整关门 · 首波；lead: workspace-012-shared-cross-module-contracts；Root done 6/6；后续 session/effective actor、保留/归档、其余 writer envelope 移交本文件 Tier A） |

## 组合门闩（用户 2026-08-08）

1. **协议优先于视觉**：在 VP-006 未 `closed` 前，**不得**将 VP-005 作为 `primary_plan` 推进实现，不得启动视觉优化波次。  
2. **MVP 子集不是终态成功条件**：`I-PROTO-001 v0.1.3` 是历史 MVP 冻结切片；整份 v2.7.0 契约由 VP-006 收口。  
3. 已关闭 VP-001～004 的历史证据与 status **不重写**。

> **协议 pin 注记（2026-08-14 · VR-020 · editorial）**：协议覆盖权威由 `v2.7.0` 升至 `v2.8.0`（additive 超集；`I-PROTO-FULL-001` v1.0.1 仍为 v2.7.0 历史分母、已被 v2.8.0 覆盖）。身份权威见 `apps/web/src/protocol/upstream/provenance-v2.8.json`。

## 已确认但尚未纳入新 VP 的后续方向

> 2026-08-14：本表原有方向（订单/钱包/类目/通知等业务能力、业务能力后续波次）曾由 **VP-011** 承接（标准 Admin 通用模块 + 常用业务领域一起入候选池再分档）。**2026-08-18：VP-011 已 `closed`**，四档能力地图已上提至本文件；后续各领域/波次按四档地图在组合层登记，触发时新建 VP/工作区。

| 顺序 | 方向 | 与前序关系 | 建立 VP 前的约束 |
|------|------|------------|------------------|
| 9 | 订单、钱包、类目、通知等业务能力 | 默认承载：VP-003 架构 + VP-004 playbook + **VP-006 协议面** + VP-005 设计系统 + VP-007 locale/settings，并消费 VP-008 的准入结论 | 建 VP 前须 `/vision` 复核；VP-008 已 `closed` 且 `go` 已签发（候选 `ed99e88`），后续业务 VP 可实现；每个后续业务 VP 激活前必须针对拟消费候选与 scope 完成并记录 freshness review，复核失败或证据不可用时暂停 `go` 并回流 VP-008 重验证或 P-004 裁决；单领域问题留在该业务 VP 的 Root/Goal 台账，共享基架或 `go` 语义问题由 `/vision` 决定重开 VP-008 或新建准入 VP；不得用业务模块倒逼恢复长期双线、跳过协议覆盖或私增协议语义 |
| 10 | 业务能力后续波次 | 在 VP-008 `go` 有效且无未恢复的共享基架 Critical/High 阻断时推进 | 消费 VP-008 `go`；共享基架问题回流 **active VP-009** 波次，不默认新开一次性加固 VP | 建 VP 前须 `/vision` 复核 |
| 11 | 共享基架横切能力、扩展接缝与后续业务域（四档能力地图） | 横切契约**首波**已由 **VP-012** `closed` 交付；lead 区增量 R7 保留/归档、R8 session+envelope 已 done。剩余 Tier A（impersonation/effective actor 产品化若再出现）与扩展接缝、Tier D 按本表登记，按触发条件独立立项 | 共享基架安全/符合性问题回流 **active VP-009/VP-010**；横切增量不默认重开 VP-012，须 `/vision` 复核后新建 VP；Tier D 在业务成立后再独立立项 |

### 四档能力地图（组合层后续方向 · 2026-08-18 上提）

> 本表为组合层权威登记，替代 VP-011 R5 作为“未来方向”的归属。详细历史证据仍保留在 `workspace-011` 的 `I-011-002`。

| 档 | 含义 | 当前/未来承载 |
|----|------|----------------|
| **Tier A · Admin 基架规划** | 身份/安全/审计/并发/运维/配置迁移等通用基架能力 | **VP-012** 首波已 `closed`。保留/归档 = GOAL-008（done）。session + D-003 外 envelope = GOAL-009（done）。effective actor 冻结为当前 actor；不做 impersonation |
| **Tier B · 扩展接缝** | 领域事件、通知 Transport、SSO、Approval、Entitlement、多组织上下文等 | 预留接口；出现真实触发条件后新建 VP/工作区 |
| **Tier C · Admin 体验增强** | 全局搜索、Saved Views、批量任务结果、统一错误恢复等 | 后续体验波次；按用户/产品触发新建 VP/工作区 |
| **Tier D · 真实业务领域** | Catalog、库存、订单/支付/退款、物流、营销、订阅、工单、CMS 等 | 业务成立后独立立项为业务 VP/工作区 |

**当前组合焦点**：无 active 交付 VP。持续程序 = **VP-009 `active`**（共享基架安全与健壮性；lead workspace-009；波次 W1–W4 与 W6 均 done，W5 扫描 0 中高危未开子目标）与 **VP-010 `active`**（设计意图—实现符合性；lead workspace-010；波次 W1–W13 均 done）。VP-001～008 仍为历史 `closed`；**VP-011 已于 2026-08-18 有界 `closed`**；**VP-012 已于 2026-08-19 完整 `closed`**（首波横切契约；lead workspace-012 Root done 6/6）。VP-008 `go` 消费有效性在无新的共享基架阻断时保持可消费。协议覆盖权威 `I-PROTO-FULL-001`（v2.7.0 历史分母，被 v2.8.0 覆盖）。

## 单主线模块化策略

未来 fork 起点统一由同一代码主线、模块候选集与启动时 Profile 表达，权威见 [module-architecture.md](../architecture/module-architecture.md) 和 VP-003。原 [dual-track-contract.md](dual-track-contract.md) 已转为历史记录。
