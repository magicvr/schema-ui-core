---
doc_type: vision-roadmap
title: 愿景组合编排
status: active
created: 2026-07-31
updated: 2026-08-11
parent: null
version: 0.18.0
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
| 9 | [VP-009-production-hardening](plans/VP-009-production-hardening.md) | 生产加固：**共享基架持续安全与健壮性程序**（周期扫描、波次修复、与 VP-008 `go` 消费有效性接口）；具体 finding 由工作区波次子目标承接。 | 继承 VP-003/004/005/006/007 + **VP-008 `go` 消费有效性**；共享基架缺陷可暂挂/恢复 `go` | **active**（2026-08-10 语义纠正为长期程序；曾误 `closed` 已撤销；lead: workspace-009-production-hardening；Root **active** 程序容器；波次 W1 GOAL-002 done、W2 GOAL-003 done） |
| 10 | [VP-010-design-implementation-conformance](plans/VP-010-design-implementation-conformance.md) | 设计意图与实现符合性：**共享基架持续对齐程序**（周期对照 as-designed / as-built、conformance gap 分流、波次整改、与 VP-008 `go` 消费有效性接口、与 VP-009 正交）；具体 gap 由工作区波次子目标承接。 | 继承 VP-003/004/005/006/007/008 + **VP-008 `go` 消费有效性**；与 **VP-009** 正交（安全 vs 符合性） | **active**（2026-08-11 用户确认类 VP-009 长期程序；lead: workspace-010-design-implementation-conformance；Root **active** 程序容器；首波 W1 = GOAL-002 范例面可选化） |

## 组合门闩（用户 2026-08-08）

1. **协议优先于视觉**：在 VP-006 未 `closed` 前，**不得**将 VP-005 作为 `primary_plan` 推进实现，不得启动视觉优化波次。  
2. **MVP 子集不是终态成功条件**：`I-PROTO-001 v0.1.3` 是历史 MVP 冻结切片；整份 v2.7.0 契约由 VP-006 收口。  
3. 已关闭 VP-001～004 的历史证据与 status **不重写**。

## 已确认但尚未纳入新 VP 的后续方向

| 顺序 | 方向 | 与前序关系 | 建立 VP 前的约束 |
|------|------|------------|------------------|
| 9 | 订单、钱包、类目、通知等业务能力 | 默认承载：VP-003 架构 + VP-004 playbook + **VP-006 协议面** + VP-005 设计系统 + VP-007 locale/settings，并消费 VP-008 的准入结论 | 建 VP 前须 `/vision` 复核；VP-008 已 `closed` 且 `go` 已签发（候选 `ed99e88`），后续业务 VP 可实现；每个后续业务 VP 激活前必须针对拟消费候选与 scope 完成并记录 freshness review，复核失败或证据不可用时暂停 `go` 并回流 VP-008 重验证或 P-004 裁决；单领域问题留在该业务 VP 的 Root/Goal 台账，共享基架或 `go` 语义问题由 `/vision` 决定重开 VP-008 或新建准入 VP；不得用业务模块倒逼恢复长期双线、跳过协议覆盖或私增协议语义 |
| 10 | 业务能力后续波次 | 在 VP-008 `go` 有效且无未恢复的共享基架 Critical/High 阻断时推进 | 消费 VP-008 `go`；共享基架问题回流 **active VP-009** 波次，不默认新开一次性加固 VP | 建 VP 前须 `/vision` 复核 |
+ **VP-010 `active`**（**设计意图—实现符合性**持续对齐程序；lead workspace-010；Root 长期容器；首波 W1 为 GOAL-002 范例面可选化）。VP-001～008 仍为历史 `closed`；**VP-008 `go` 消费有效性在无新的共享基架阻断时保持可消费**（W1 曾落盘恢复证据；新阻断由 VP-009 波次处理；改变 Profile/模块矩阵/Manifest 装配语义的符合性缺口由 VP-010 波次按规则暂挂/重验证）。协议覆盖权威 `I-PROTO-FULL-001`。订单/钱包/类目/通知仍为后续独立业务 VP 候选，激活前须 freshness review，并受 `go` 失效/暂挂与 VP-009 开放 Critical/High、VP-010 开放 blocker/major gap
**当前组合焦点**：**VP-009 `active`**（共享基架**持续**安全与健壮性程序；lead workspace-009；Root 长期容器；已完成波次 W1/W2 为档案）。VP-001～008 仍为历史 `closed`；**VP-008 `go` 消费有效性在无新的共享基架阻断时保持可消费**（W1 曾落盘恢复证据；新阻断由 VP-009 波次处理）。协议覆盖权威 `I-PROTO-FULL-001`。订单/钱包/类目/通知仍为后续独立业务 VP 候选，激活前须 freshness review，并受 `go` 失效/暂挂与 VP-009 开放 Critical/High 约束。

## 单主线模块化策略

未来 fork 起点统一由同一代码主线、模块候选集与启动时 Profile 表达，权威见 [module-architecture.md](../architecture/module-architecture.md) 和 VP-003。原 [dual-track-contract.md](dual-track-contract.md) 已转为历史记录。
