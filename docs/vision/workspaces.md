---
doc_type: vision-workspaces
title: 工作区贡献图
status: active
created: 2026-07-31
updated: 2026-08-18
parent: null
version: 0.14.0
---

# 工作区贡献图

| workspace_id | canonical_scope | root_goal | role | primary_plan | status |
|--------------|-----------------|-----------|------|--------------|--------|
| workspace-001-mvp-admin-foundation | docs/workspaces/workspace-001-mvp-admin-foundation/ | GOAL-001-mvp-admin-foundation | primary | VP-001-mvp-admin-foundation | active |
| workspace-002-production-admin-foundation | docs/workspaces/workspace-002-production-admin-foundation/ | GOAL-001-production-admin-foundation | delivery | VP-002-production-admin-foundation | active |
| workspace-003-modular-admin-architecture | docs/workspaces/workspace-003-modular-admin-architecture/ | GOAL-001-modular-admin-architecture | delivery | VP-003-modular-admin-architecture | active |
| workspace-004-module-contribution-readiness | docs/workspaces/workspace-004-module-contribution-readiness/ | GOAL-001-module-contribution-readiness | delivery | VP-004-module-contribution-readiness | active |
| workspace-005-full-protocol-contract-v2-7-0 | docs/workspaces/workspace-005-full-protocol-contract-v2-7-0/ | GOAL-001-full-protocol-contract-v2-7-0 | delivery | VP-006-full-protocol-contract-v2-7-0 | active |
| workspace-006-design-system-and-ui-experience | docs/workspaces/workspace-006-design-system-and-ui-experience/ | GOAL-001-design-system-and-ui-experience | delivery | VP-005-design-system-and-ui-experience | active |
| workspace-007-localization-and-system-settings | docs/workspaces/workspace-007-localization-and-system-settings/ | GOAL-001-localization-and-system-settings | delivery | VP-007-localization-and-system-settings | active |
| workspace-008-admin-module-readiness | docs/workspaces/workspace-008-admin-module-readiness/ | GOAL-001-admin-module-readiness | delivery | VP-008-admin-module-readiness-and-foundation-convergence | active |
| workspace-009-production-hardening | docs/workspaces/workspace-009-production-hardening/ | GOAL-001-production-hardening | lead | VP-009-production-hardening | active（Root **active** 长期程序容器；波次 W1–W4 与 W6 done，W5 扫描 0 中高危未开子目标；2026-08-10 语义纠正） |
| workspace-010-design-implementation-conformance | docs/workspaces/workspace-010-design-implementation-conformance/ | GOAL-001-design-implementation-conformance | lead | VP-010-design-implementation-conformance | active（Root **active** 长期程序容器；2026-08-11 开区；波次 W1–W13 done，`go` 均无新暂挂） |
| workspace-011-admin-functional-modules | docs/workspaces/workspace-011-admin-functional-modules/ | GOAL-001-admin-functional-modules | lead | VP-011-admin-functional-modules | active（历史；VP-011 已于 2026-08-18 有界 closed；Root done；四档能力地图上提至 vision roadmap） |
| workspace-012-shared-cross-module-contracts | docs/workspaces/workspace-012-shared-cross-module-contracts/ | GOAL-001-shared-cross-module-contracts | lead | VP-012-shared-cross-module-contracts | active（2026-08-18 开区；VP-012 激活；首波 = 横切契约波） |

## 说明

- **workspace-011（2026-08-14 开区）**：VP-011（标准 Admin 功能模块分档交付）唯一 lead delivery 工作区；消费前 freshness review **PASS**（候选 `f14ab9d`；VR-020 pin bump、W5/W6 归档、F-1a/b/c fixed）；Root `GOAL-001-admin-functional-modules` 纲领路线图 R1 有界调研 → R2 一等公民 → R3 常用 → R4 增补 backlog → R5 四档能力地图与跨模块路线图登记。**2026-08-18：VP-011 有界 closed，Root done，四档能力地图上提至 vision roadmap**。不改变 Charter `primary_workspace`。
- **workspace-012（2026-08-18 开区）**：VP-012（共享横切契约与平台基架）唯一 lead delivery 工作区；首波 = 横切契约波（correlation/审计模型/并发幂等/异步 Job + maintenance 门控/API Token）；不承载 Tier D 业务域；与 VP-009/VP-010 正交分流。
- 首个工作区由 `/govern` 于 2026-07-31 开区；与 Charter `primary_workspace`、工作区 `workspace.md` 的 `vision_role: primary` 一致。
- 第二个工作区由用户于 2026-08-01 确认，经 `/vision` 完成 VP-002 激活与绑定、由 `/govern` 建立实现层；它是 VP-002 当前唯一 lead workspace，角色为 `delivery`。
- 新 delivery 工作区不改变 Charter 的 `primary_workspace`，也不重开 VP-001 或旧 Root。
- **VP-002 已于 2026-08-04 经 `/vision` 用户确认关门（`closed`）**：workspace-002 与 Root `GOAL-001-production-admin-foundation`（`done / 5/5`）的历史绑定保留，默认不接新区（reopen 须用户确认）。
- 2026-08-04 strategic re-align 后，工作区及其已完成 Root 均精确对齐 Charter `@0.2.0`，但不改变 VP/Goal 的历史状态、progress 或证据。
- `VP-003-modular-admin-architecture` 已于 2026-08-04 由用户确认激活，并绑定 `workspace-003-modular-admin-architecture` 为当前唯一 lead / delivery 工作区；`/govern` 已建立对应 Root。建区不构成 VP 实现或关门证据。
- **VP-003 已于 2026-08-06 经 `/vision` 用户确认关门（`closed`）**：workspace-003 与 Root `GOAL-001-modular-admin-architecture`（`done / 6/6`）的历史绑定保留，默认不接新区（reopen 须用户确认）；关门证据链见 VP-003 关门记录。
- **VP-004 已于 2026-08-06 经 `/vision` 用户确认关门（`closed`）**：workspace-004 与 Root `GOAL-001-module-contribution-readiness`（`done / 4/4`）的历史绑定保留，默认不接新区（reopen 须用户确认）；关门证据链见 VP-004 关门记录。不改变 Charter `primary_workspace`。
- **VP-006 已于 2026-08-08 经 `/vision` 用户确认激活（`active`）**：唯一 lead / delivery = `workspace-005-full-protocol-contract-v2-7-0`；Root `GOAL-001-full-protocol-contract-v2-7-0` 由 `/govern` 同日 scaffold。激活与建区 **不**等于覆盖表 `I-PROTO-FULL-001` 已冻结或全量协议已兼容。**禁止**在 closed workspace-003/004 吸收本意图。不改变 Charter `primary_workspace`。
- **VP-006 已于 2026-08-08 经用户书面确认关门（`closed`）**：workspace-005 与 Root `GOAL-001-full-protocol-contract-v2-7-0`（`done / 6/6`）的历史绑定保留；整份 v2.7.0 契约覆盖权威 = `I-PROTO-FULL-001` v1.0.1（12/12 include，320 case = 318 executed + 2 local adapter excluded），实现/验证/关门与勘误证据链见该区五件套（E-001～E-007、A-001～A-003）。
- **VP-005 已于 2026-08-09 经 `/vision` 用户书面确认关门（`closed`）**：workspace-006 与 Root `GOAL-001-design-system-and-ui-experience`（`done / 5/5`）的历史绑定保留，默认不接新区（reopen 须用户确认）；S1–S5 交付与回归证据（vitest 616/616 + e2e 2/2）见该区五件套（D-008 / A-012 / E-010；GOAL-005 E-001/E-002）；residual：F-VUI-007/010/011 `accepted-residual`、I-004 open non-blocking。不改变 Charter `primary_workspace`。
- **VP-007 已于 2026-08-09 经 `/vision` 用户确认激活（`active`）**：唯一 lead / delivery = `workspace-007-localization-and-system-settings`；`/govern` 已于 2026-08-09 scaffold：Root `GOAL-001-localization-and-system-settings`（S0–S5 纲领路线图）。
- **VP-007 已于 2026-08-09 经用户书面确认关门（`closed`）**：workspace-007 与 Root `GOAL-001-localization-and-system-settings`（`done / 6/6`）的历史绑定保留，默认不接新区（reopen 须用户确认）；S0–S5 证据与 independent 审计闭环见该区五件套（GOAL-006 A-001/A-002、D-002、E-001～E-003；Root `attachments/S5-evidence-matrix.md`）。不改变 Charter `primary_workspace`。
- **VP-008 已于 2026-08-10 经用户书面确认关门（`closed`）**：workspace-008 与 Root `GOAL-001-admin-module-readiness`（`done / 6/6`，S0–S5）的历史绑定保留，默认不接新区（reopen 须用户确认）；`go` 已于 2026-08-10 签发（候选 `ed99e88`，clean；D-001），解锁后续标准业务模块实现，每个后续业务 VP 激活前须完成消费前 freshness review。关门证据链见该区五件套（GOAL-007 D-001、A-001/A-002/A-003、S5-evidence-matrix）。不改变 Charter `primary_workspace`。
- **VP-009 已于 2026-08-10 落盘并激活**：初衷曾被写成「单次加固波次」；同日完成 W1（GOAL-002 16 项）后曾误记 `closed`。
- **VP-009 / workspace-009 语义纠正（2026-08-10 用户书面）**：VP-009 = **持续**共享基架安全与健壮性程序（`active`）；workspace-009 = 唯一 lead；Root `GOAL-001-production-hardening` = **长期程序容器（`active`）**；有界波次 = 子目标（W1–W4 与 W6 done；W5 扫描 0 中高危未开子目标）。单波完成或 `go` 恢复**不等于** VP/Root 关门。不改变 Charter `primary_workspace`。
- **VP-010 / workspace-010 已于 2026-08-11 落盘并激活**：长期**设计意图—实现符合性**程序（类 VP-009 持续程序语义，与 009 正交：安全 vs 符合性）；lead = workspace-010-design-implementation-conformance；Root `GOAL-001-design-implementation-conformance` = **长期程序容器（`active`）**；波次 W1–W13 均 done（W1 = GOAL-002 范例/演示产品面可选模块化）。改变 Profile 默认集 / 模块矩阵 / Manifest 装配语义的 gap 按 VP-008 `go` 消费有效性规则暂挂/重验证。单波完成**不等于** VP/Root 关门。不改变 Charter `primary_workspace`。
- 目标生命周期与 progress 以工作区内 `goal-tree.md` / 五件套为准；本文件不是第二套状态源。波次 progress 不得推导 Root/`VP` done。
