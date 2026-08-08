---
doc_type: vision-workspaces
title: 工作区贡献图
status: active
created: 2026-07-31
updated: 2026-08-08
parent: null
version: 0.6.0
---

# 工作区贡献图

| workspace_id | canonical_scope | root_goal | role | primary_plan | status |
|--------------|-----------------|-----------|------|--------------|--------|
| workspace-001-mvp-admin-foundation | docs/workspace-001-mvp-admin-foundation/ | GOAL-001-mvp-admin-foundation | primary | VP-001-mvp-admin-foundation | active |
| workspace-002-production-admin-foundation | docs/workspace-002-production-admin-foundation/ | GOAL-001-production-admin-foundation | delivery | VP-002-production-admin-foundation | active |
| workspace-003-modular-admin-architecture | docs/workspace-003-modular-admin-architecture/ | GOAL-001-modular-admin-architecture | delivery | VP-003-modular-admin-architecture | active |
| workspace-004-module-contribution-readiness | docs/workspace-004-module-contribution-readiness/ | GOAL-001-module-contribution-readiness | delivery | VP-004-module-contribution-readiness | active |
| workspace-005-full-protocol-contract-v2-7-0 | docs/workspace-005-full-protocol-contract-v2-7-0/ | GOAL-001-full-protocol-contract-v2-7-0 | delivery | VP-006-full-protocol-contract-v2-7-0 | active |

## 说明

- 首个工作区由 `/govern` 于 2026-07-31 开区；与 Charter `primary_workspace`、工作区 `workspace.md` 的 `vision_role: primary` 一致。
- 第二个工作区由用户于 2026-08-01 确认，经 `/vision` 完成 VP-002 激活与绑定、由 `/govern` 建立实现层；它是 VP-002 当前唯一 lead workspace，角色为 `delivery`。
- 新 delivery 工作区不改变 Charter 的 `primary_workspace`，也不重开 VP-001 或旧 Root。
- **VP-002 已于 2026-08-04 经 `/vision` 用户确认关门（`closed`）**：workspace-002 与 Root `GOAL-001-production-admin-foundation`（`done / 5/5`）的历史绑定保留，默认不接新区（reopen 须用户确认）。
- 2026-08-04 strategic re-align 后，工作区及其已完成 Root 均精确对齐 Charter `@0.2.0`，但不改变 VP/Goal 的历史状态、progress 或证据。
- `VP-003-modular-admin-architecture` 已于 2026-08-04 由用户确认激活，并绑定 `workspace-003-modular-admin-architecture` 为当前唯一 lead / delivery 工作区；`/govern` 已建立对应 Root。建区不构成 VP 实现或关门证据。
- **VP-003 已于 2026-08-06 经 `/vision` 用户确认关门（`closed`）**：workspace-003 与 Root `GOAL-001-modular-admin-architecture`（`done / 6/6`）的历史绑定保留，默认不接新区（reopen 须用户确认）；关门证据链见 VP-003 关门记录。
- **VP-004 已于 2026-08-06 经 `/vision` 用户确认关门（`closed`）**：workspace-004 与 Root `GOAL-001-module-contribution-readiness`（`done / 4/4`）的历史绑定保留，默认不接新区（reopen 须用户确认）；关门证据链见 VP-004 关门记录。不改变 Charter `primary_workspace`。
- **VP-006 已于 2026-08-08 经 `/vision` 用户确认激活（`active`）**：唯一 lead / delivery = `workspace-005-full-protocol-contract-v2-7-0`；Root `GOAL-001-full-protocol-contract-v2-7-0` 由 `/govern` 同日 scaffold。激活与建区 **不**等于覆盖表 `I-PROTO-FULL-001` 已冻结或全量协议已兼容。**禁止**在 closed workspace-003/004 吸收本意图。不改变 Charter `primary_workspace`。
- **VP-006 已于 2026-08-08 经用户书面确认关门（`closed`）**：workspace-005 与 Root `GOAL-001-full-protocol-contract-v2-7-0`（`done / 6/6`）的历史绑定保留；整份 v2.7.0 契约覆盖权威 = `I-PROTO-FULL-001`（12/12 include，0 exclude），实现/验证/关门证据链见该区五件套（E-001～E-004、A-001/A-002）。VP-005 实施解冻与否由用户另行决策（`F-V018` 仍 open required，未闭合前不得放行）。
- 目标生命周期与 progress 以工作区内 `goal-tree.md` / 五件套为准；本文件不是第二套状态源。
