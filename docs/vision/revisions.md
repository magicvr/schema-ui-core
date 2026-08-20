---
doc_type: vision-revisions
title: Charter 修订台账
status: active
created: 2026-07-31
updated: 2026-08-20
parent: null
version: 0.4.11
---

# Charter 修订台账

`VR-*` 记录 Charter 的初建与后续修订；它不同于 [reviews.md](reviews.md) 中的 `VRev-*` 愿景审视。

| id | date | class | scope | summary |
|----|------|-------|-------|---------|
| VR-001 | 2026-07-31 | initial | Charter 初建 | 用户确认唯一愿景、外部协议来源 `schema-ui-docs@v2.7.0`、React + Go 技术方向、三条非目标，以及首个 MVP VP。当前无既有 VP 或工作区，因此无 re-align 影响。 |
| VR-002 | 2026-07-31 | editorial | H-001 状态措辞 | 响应 VRev-003 `F-V006`：将 H-001 状态分列「清单提取 verified / 覆盖子集冻结 open」，避免读成覆盖已可冻结。`vision_id@version` 仍为 `schema-ui-core-admin-foundation@0.1.0`；**无** strategic、**无** re-align。 |
| VR-003 | 2026-07-31 | editorial | H-001 覆盖冻结事实 | Root D-009 已按 `I-PROTO-001` v0.1.3 完成 R2 冻结，更新 H-001 的第二个分列状态为 `verified`；`vision_id@version` 仍为 `schema-ui-core-admin-foundation@0.1.0`，不改 Charter 目的、边界或非目标，**无** strategic、**无** re-align。 |
| VR-004 | 2026-08-04 | strategic | 单主线模块化终态 | 用户接受架构评议的全部建议并明确：VP 必须表达完整最终意图，Activity/Settings 等试点只作为迭代路线图。Charter 升至 `schema-ui-core-admin-foundation@0.2.0`，以单主线 + Profile 替代双线长期维护；建立 planned `VP-003-modular-admin-architecture` 与架构契约；双线契约退役为历史；VP-001/002、协议清单、workspace-001/002 与两棵 Root 仅更新精确对齐引用，不重开已关闭交付、不改变 Goal status/progress。VRev-006 self review 记录战略 re-align 结果。 |
| VR-005 | 2026-08-06 | editorial | VP 关系指针 | VP-003 已 closed；「与工作区/VP 的关系」改为索引 closed VP-001～003，并将可挂接下一个意图指向 planned `VP-004-module-contribution-readiness`。不改目的/边界/非目标；`vision_id@version` 仍为 `@0.2.0`；**无** re-align。 |
| VR-006 | 2026-08-06 | editorial | VP 关系指针 | VP-004 用户确认激活；Charter 关系节改为 `active` + lead `workspace-004-module-contribution-readiness`（区待 scaffold）。不改目的/边界/非目标；`vision_id@version` 仍为 `@0.2.0`；**无** re-align。 |
| VR-007 | 2026-08-08 | editorial | 协议目标澄清 + 组合指针 | 用户确认：目标为 `schema-ui-docs@v2.7.0` **整份契约**可验证兼容（非长期停在 MVP 子集）。Charter 协议来源/H-001/关系节澄清 `I-PROTO-001 v0.1.3` 仅为 MVP 切片；组合焦点 → planned VP-006，VP-005 实施硬冻结至 VP-006 closed。不改目的/成功边界编号/非目标正文；`vision_id@version` 仍为 `@0.2.0`；**无** strategic、**无** re-align 宽阻断。 |
| VR-008 | 2026-08-08 | editorial | VP 关系指针 | VP-006 用户确认激活；Charter 关系节改为 active 交付 VP-006 + lead `workspace-005-full-protocol-contract-v2-7-0`；VP-005 仍实施冻结。不改目的/边界/非目标；`vision_id@version` 仍为 `@0.2.0`；**无** re-align。 |
| VR-009 | 2026-08-09 | editorial | VP 关系指针 + 组合焦点 | VRev-011 F-V018/019/020 fixed（VP-005 v0.3.0）；用户确认激活 VP-005（v0.4.0 `active`）。Charter 关系节改为当前交付 VP-005；协议语义注 VP-006 closed。不改目的/边界/非目标；`vision_id@version` 仍为 `@0.2.0`；**无** re-align。 |
| VR-010 | 2026-08-09 | editorial | VP 关系指针 | `/govern` 开区：VP-005 `lead_workspace` = `workspace-006-design-system-and-ui-experience`（Root `GOAL-001-design-system-and-ui-experience`）。Charter 关系节同步 lead 路径。不改目的/边界/非目标；`@0.2.0`；**无** re-align。 |
| VR-011 | 2026-08-09 | editorial | VP 关系指针 + 组合焦点 | VRev-015 `F-V027`/`F-V028` → fixed（用户书面「确认关门」）。VP-005 `active → closed`（v0.5.0，关门记录含 exit↔证据映射 + residual 点名）；Charter 关系节改为「无 active 交付 VP」。不改目的/边界/非目标；`@0.2.0`；**无** re-align。 |
| VR-012 | 2026-08-09 | editorial | 协议来源与 VP 关系指针 | 响应 VRev-016 `F-V031`：删除协议来源段残留的“VP-005 active”旧指针，明确 VP-005 closed、当前无 active，并索引 planned VP-007 v0.1.1（0 区）。不改目的/成功边界/非目标或 `vision_id@version`（仍 `@0.2.0`）；**无** re-align。 |
| VR-013 | 2026-08-09 | editorial | VP 关系指针 + 组合焦点 | 用户确认激活 VP-007（v0.2.0 `active`；lead delivery `workspace-007-localization-and-system-settings`，Root 待 `/govern` scaffold）；Charter 关系节改为当前交付 VP-007；roadmap / workspaces / VP-007 绑定表同步。不改目的/边界/非目标；`@0.2.0`；**无** re-align。 |
| VR-014 | 2026-08-10 | editorial | VP-007 关门投影 + VP-008 组合指针 | 修正 Charter 与 reviews 摘要残留的“VP-007 active”旧投影：VP-007 已于 2026-08-09 用户书面确认 `closed`，当前无 active 交付 VP。用户确认新建 planned `VP-008-admin-module-readiness-and-foundation-convergence`，作为正式业务模块前的全基架准入波次；未激活、未建区。不改目的/成功边界/非目标或 `vision_id@version`（仍 `@0.2.0`）；**无** strategic、**无** re-align。 |
| VR-015 | 2026-08-10 | editorial | VRev-026 响应 + VP-008 激活指针 | 用户采纳 VRev-026 `pass`、V-F053 → `fixed`，并确认激活 VP-008。VP-008 v0.10.0 改为 `active`、0 workspace；Charter/roadmap/reviews 投影同步，空区按 §5.1 从 2026-08-10 起进入 14 日宽限，最迟 2026-08-24 复核。不改目的/成功边界/非目标或 `vision_id@version`（仍 `@0.2.0`）；**无** strategic、**无** re-align。 |
| VR-016 | 2026-08-10 | editorial | VP-008 lead workspace + Root scaffold | 用户确认单工作区、lead `workspace-008-admin-module-readiness`、Root `GOAL-001-admin-module-readiness` 与 independent provider `GitHub Copilot · /audit`；`/govern` 完成物理 scaffold。VP-008 保持 `active`，workspace/Root/Goal 对齐 Charter `@0.2.0`，尚未产生 `go`；不改 Charter 目的、成功边界、非目标或 `vision_id@version`，**无** strategic、**无** re-align。 |
| VR-017 | 2026-08-10 | editorial | I-PROTO-FULL-001 执行分母勘误 | 响应 workspace-008 F-001：覆盖权威升至 v1.0.1，12/12 域、24/24 registry、16/16 suite include 不变；将 `320/320 全绿` 修正为 320 total = 318 executed + 2 local adapter excluded。同步 Charter H-001、roadmap、workspaces 与 VP-006 关门投影；不改 Charter 目的/成功边界/非目标或 `vision_id@version`（仍 `@0.2.0`），不重开 VP-006，**无** strategic、**无** re-align。 |
| VR-018 | 2026-08-10 | editorial | VP-008 关门投影 | 用户确认关闭 VP-008：status `active → closed`（v0.13.0）；Charter 目标语义、组合焦点与关系节、roadmap、workspaces、reviews 摘要同步为 closed 投影（`go` 候选 `ed99e88`、clean，解锁后续标准业务模块，每个激活前须 freshness review）。不改 Charter 目的/成功边界/非目标或 `vision_id@version`（仍 `@0.2.0`），不重开历史 VP，**无** strategic、**无** re-align。 |
| VR-019 | 2026-08-14 | editorial | VP 关系指针 + 组合编排 | 用户确认结构选型并新建 planned `VP-011-admin-functional-modules`（标准 Admin 通用模块 + 常用业务领域分档交付；0 区）；调研回写 Root 纲领路线图而非 VP，VP 只留意图 + 三档方法论；roadmap 已落盘意图 + 后续方向同步；激活前须完成 VP-008 `go` 消费前 freshness review。不改 Charter 目的/成功边界/非目标或 `vision_id@version`（仍 `@0.2.0`），**无** strategic、**无** re-align。 |
| VR-020 | 2026-08-14 | editorial | 协议 pin bump（v2.7.0 → v2.8.0） | 用户裁决 A（editorial）：`schema-ui-docs@v2.7.0` → `v2.8.0`（additive 超集，v2.7.0 机器契约保留，新增 Host/App 互操作层）；Charter 协议来源/目标语义/成功边界 1/H-001 同步；`I-PROTO-FULL-001` v1.0.1 保留为 v2.7.0 历史分母、被 v2.8.0 覆盖；身份权威 = provenance-v2.8.json（VP-010 W3 固定）。不改 `vision_id@version`（仍 `@0.2.0`）、不重开历史 VP，**无** strategic、**无** re-align。 |
| VR-021 | 2026-08-14 | editorial | VP 关系指针 + 组合编排 | 用户确认激活 VP-011（v0.2.0 `active`，lead workspace-011-admin-functional-modules）：消费前 freshness review 第二轮 **PASS**（候选 `f14ab9d`，V-001～V-008 全绿；F-1a/b/c 由 VP-010 W6 修复；VR-020 pin bump 已对齐）；`/govern` 当日 scaffold 工作区 + Root（首阶段 = 有界调研）。不改 Charter 目的/成功边界/非目标或 `vision_id@version`（仍 `@0.2.0`），**无** strategic、**无** re-align。 |
| VR-022 | 2026-08-18 | editorial | VP 关系指针 + 组合编排 | 用户确认新建并激活 VP-012（v0.1.0 `active`，lead workspace-012-shared-cross-module-contracts）：承载共享横切契约与平台基架（correlation/审计模型/并发幂等/异步 Job + maintenance 门控/API Token）；承接 VP-011 R5 中“横切契约”部分，Tier D 业务域仍按触发条件；不改变 Charter 边界或 `vision_id@version`（仍 `@0.2.0`），**无** strategic、**无** re-align。 |
| VR-023 | 2026-08-18 | editorial | VP 关系指针 + 组合编排 + 关门投影 | 用户确认四档能力地图上提至 vision roadmap，并关闭 VP-011（`active → closed` v0.4.0）与 workspace-011 Root（`active → done`）：标准 Admin 功能模块波次完成；S-05/S-06/S-07/S-08/S-13 及 B-01～B-11 reclassify 到四档地图未来 VP/工作区；当前 active 交付 VP = VP-012；不改变 Charter 边界或 `vision_id@version`（仍 `@0.2.0`），**无** strategic、**无** re-align。 |
| VR-024 | 2026-08-19 | editorial | VP-012 完整关门投影 | 用户确认关闭 VP-012（v0.2.0 `active → closed`，完整 · 首波）：lead workspace-012 Root `done 6/6`；VRev-028 V-F057 → fixed；session/effective actor、保留/归档、D-003 外 writer envelope 移交 roadmap Tier A，不作为本 VP residual。Charter 关系节改为无 active 交付 VP；持续程序仍为 VP-009/VP-010。不改 Charter 目的/成功边界/非目标或 `vision_id@version`（仍 `@0.2.0`），**无** strategic、**无** re-align。 |
| VR-025 | 2026-08-20 | editorial | 组合编排双分支 | 用户确认后续方向改为可并行的架构 / 业务两条轨道。落盘名称为架构分支（运行时平台横切清单）与产品分支（原四档地图；业务领域 = Tier D）。不改 Charter 目的/成功边界/非目标或 `vision_id@version`（仍 `@0.2.0`），不建 VP、不开区，**无** strategic、**无** re-align。同日被 **VR-026** 三分取代，本行保留为过程记录。 |
| VR-026 | 2026-08-20 | editorial | 组合编排三分支 | 用户确认将后续方向改为 **架构 / Admin 功能 / 业务域** 三条可并行轨道（取代 VR-025 二分）。原四档 A 剩余+B+C → Admin 功能；Tier D → 业务域；运行时平台清单 → 架构。不改 Charter 目的/成功边界/非目标或 `vision_id@version`（仍 `@0.2.0`），不建 VP、不开区，**无** strategic、**无** re-align。 |
| VR-027 | 2026-08-20 | editorial | 架构分支 Store 双方言 | 用户确认：不引入 ORM；自持 PostgreSQL + SQLite 两个 Store 方言；内核 Store 为持久化端口而非业务仓库；逻辑 schema 一份、物理 SQL 可成对；PostgreSQL 为生产验收权威，SQLite 为 dev/mvp/快测默认且合同平等（不得残缺）。写入 roadmap RT-P03。不改 Charter 目的/成功边界/非目标或 `vision_id@version`（仍 `@0.2.0`），不建 VP、不开区，**无** strategic、**无** re-align。 |
