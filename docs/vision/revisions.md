---
doc_type: vision-revisions
title: Charter 修订台账
status: active
created: 2026-07-31
updated: 2026-08-10
parent: null
version: 0.3.8
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
