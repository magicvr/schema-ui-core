---
title: 目标树 · workspace-001-mvp-admin-foundation
status: active
created: 2026-07-31
updated: 2026-08-01
parent: null
version: 0.31.0
workspace_id: workspace-001-mvp-admin-foundation
---

# 目标树 · MVP Admin 基架

> 工作区：`workspace-001-mvp-admin-foundation`  
> canonical：`docs/workspace-001-mvp-admin-foundation/`  
> Root：`GOAL-001-mvp-admin-foundation`  
> primary_plan：`VP-001-mvp-admin-foundation`

## 树

```text
GOAL-001-mvp-admin-foundation  [done]  MVP Admin 基架  progress=6/6
├── GOAL-002-r1-repo-layout-conventions  [done]  R1 · 仓库布局与包管理约定
├── GOAL-003-r1-api-go-scaffold          [done]  R1 · Go API 工程骨架
├── GOAL-004-r1-web-react-scaffold       [done]  R1 · React Web 工程骨架
├── GOAL-005-r3-admin-shell-navigation   [done]  R3 · Admin 外壳与导航
├── GOAL-006-r4-account-permission       [done]  R4 · 核心账号与权限
├── GOAL-007-r5-examples-contract-verification [done]  R5 · 纳入域范例与契约验证
├── GOAL-008-r6-integration-acceptance-vp-evidence [done] R6 · 集成验收与 VP 证据
└── GOAL-009-mvp-bugfix-followup       [active] MVP 代码审视 bug 修正  progress=0/5
```

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-mvp-admin-foundation | MVP Admin 基架 | null | done | 6/6 | 2026-08-01 |
| GOAL-002-r1-repo-layout-conventions | R1 · 仓库布局与包管理约定 | GOAL-001-mvp-admin-foundation | done | — | 2026-07-31 |
| GOAL-003-r1-api-go-scaffold | R1 · Go API 工程骨架 | GOAL-001-mvp-admin-foundation | done | — | 2026-07-31 |
| GOAL-004-r1-web-react-scaffold | R1 · React Web 工程骨架 | GOAL-001-mvp-admin-foundation | done | — | 2026-07-31 |
| GOAL-005-r3-admin-shell-navigation | R3 · Admin 外壳与导航 | GOAL-001-mvp-admin-foundation | done | — | 2026-07-31 |
| GOAL-006-r4-account-permission | R4 · 核心账号与权限 | GOAL-001-mvp-admin-foundation | done | — | 2026-07-31 |
| GOAL-007-r5-examples-contract-verification | R5 · 纳入域范例与契约验证 | GOAL-001-mvp-admin-foundation | done | — | 2026-08-01 |
| GOAL-008-r6-integration-acceptance-vp-evidence | R6 · 集成验收与 VP 证据 | GOAL-001-mvp-admin-foundation | done | — | 2026-08-01 |
| GOAL-009-mvp-bugfix-followup | MVP 代码审视 bug 修正 | GOAL-001-mvp-admin-foundation | active | 0/5 | 2026-08-01 |

## 说明

- **Root 已关门（2026-08-01）**：`progress: 6/6` = 纲领路线图 R1–R6 全部完成（等权派生）；Root `status: done`（用户 `/govern` 授权，Root D-016）。VP-001 已 `closed`（`/vision`）；Root `done` 不等于完整协议支持。
- **GOAL-009（2026-08-01）**：Root 关门后的修正跟随子目标 `GOAL-009-mvp-bugfix-followup` → **`active`**（`progress: 0/5`）。承接代码审视 required findings F-009-001～005；长文 [audit-code-review-bugs-2026-08-01.md](GOAL-009-mvp-bugfix-followup/attachments/audit-code-review-bugs-2026-08-01.md)；A-001 independent conditional。不改 Root 纲领 6/6，不重开 VP 范围。
- R6：GOAL-008-r6-integration-acceptance-vp-evidence → **`done`**（2026-08-01）。A-005 self + A-006 independent close-out 均 pass、开放 required=0；五项 `I-008` required 均 verified；用户 `/govern` 响应 A-006 并授权关门（GOAL-008 D-006 / Root D-016）；goal-tree 同步。
- R5：GOAL-007-r5-examples-contract-verification → **`done`**（2026-08-01）。A-008 self pass + A-009 independent pass；`I-PROTO-003` verified；用户 `/govern` 响应 A-009 授权关门（D-012 / Root D-014）；`progress` **4/6 → 5/6**。
- 层级只由各目标 `00-meta.md` 的 `parent` 表达；本文件为区内索引，不是第二套状态源。
- R1：2026-07-31 `/govern` 阶段/关门自审 A-003（self）对 002/003/004 均为 **pass** → 三子目标 `done`；Root 纲领 R1 → 完成。Root A-001（independent）R1 证据复核 **pass** → A-002/D-006 响应：维持完成事实。
- R2：完成；用户书面确认按 v0.1.3 覆盖表正式冻结，D-009 / A-006 已留痕，`I-PROTO-001` = `verified`。冻结范围不等同于完整协议支持或 R3-R5 实施完成。
- R3：完成；GOAL-005 A-006 `source: self` / `verdict: pass` 已留痕，73 项测试、构建、固定 fixture 对照和 HTTP 入口证据已入账；A-007 independent 关门复审 `pass`，A-008 已响应其 recommended（F-001 索引表时点绑定、F-002 schema 等价性随 Root `I-PROTO-004` 跟进）。R4/R5 与完整协议支持仍在后续范围。
- R4：方案冻结完成（GOAL-006 D-004，`I-006-001` 与父目标 `I-PROTO-002` 均 `verified`）；R4 **实施完成**（2026-07-31）：Go 会话与 `/api/accounts/me`、Go 独立鉴权、Web `$context` 挂载、D-PERM 求值引擎与 17 例 fixture 对照（13 valid 求值 + 4 invalid 错误码；13 valid 中含 5 例 execution 时序断言）落地，`go test`/`go build`、web 94 项测试、`npm run build`、HTTP 运行时与代理联调证据入账；A-001（self）与 A-002（independent）实施阶段均 pass，A-003 已响应（F-001/F-002 fixed，F-003～F-005 跟踪）。R4 **关门完成**（2026-07-31）：A-004 关门自审（self）与 A-005 独立关门复审（independent）均 pass、开放 required=0；经用户 `/govern` 授权 GOAL-006 → `done`、Root 纲领 R4 检查点完成（`progress` → 4/6）、goal-tree 同步。F-002～F-004（Renderer 接线 / token 会话 / 双端一致性 oracle）为 recommended 跟踪项，随 R5 / 生产化 / `I-PROTO-004` 解决。
