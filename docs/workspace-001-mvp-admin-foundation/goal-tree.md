---
title: 目标树 · workspace-001-mvp-admin-foundation
status: active
created: 2026-07-31
updated: 2026-07-31
parent: null
version: 0.8.0
workspace_id: workspace-001-mvp-admin-foundation
---

# 目标树 · MVP Admin 基架

> 工作区：`workspace-001-mvp-admin-foundation`  
> canonical：`docs/workspace-001-mvp-admin-foundation/`  
> Root：`GOAL-001-mvp-admin-foundation`  
> primary_plan：`VP-001-mvp-admin-foundation`

## 树

```text
GOAL-001-mvp-admin-foundation  [active]  MVP Admin 基架  progress=3/6
├── GOAL-002-r1-repo-layout-conventions  [done]  R1 · 仓库布局与包管理约定
├── GOAL-003-r1-api-go-scaffold          [done]  R1 · Go API 工程骨架
├── GOAL-004-r1-web-react-scaffold       [done]  R1 · React Web 工程骨架
├── GOAL-005-r3-admin-shell-navigation   [done]  R3 · Admin 外壳与导航
└── GOAL-006-r4-account-permission       [active]  R4 · 核心账号与权限
```

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-mvp-admin-foundation | MVP Admin 基架 | null | active | 3/6 | 2026-07-31 |
| GOAL-002-r1-repo-layout-conventions | R1 · 仓库布局与包管理约定 | GOAL-001-mvp-admin-foundation | done | — | 2026-07-31 |
| GOAL-003-r1-api-go-scaffold | R1 · Go API 工程骨架 | GOAL-001-mvp-admin-foundation | done | — | 2026-07-31 |
| GOAL-004-r1-web-react-scaffold | R1 · React Web 工程骨架 | GOAL-001-mvp-admin-foundation | done | — | 2026-07-31 |
| GOAL-005-r3-admin-shell-navigation | R3 · Admin 外壳与导航 | GOAL-001-mvp-admin-foundation | done | — | 2026-07-31 |
| GOAL-006-r4-account-permission | R4 · 核心账号与权限 | GOAL-001-mvp-admin-foundation | active | — | 2026-07-31 |

## 说明

- Root `progress: 3/6` = 纲领路线图 R1–R6 中 R1、R2、R3 已完成（等权派生）；**不**放行 R4-R6、不关闭未相关 finding、不推导 Root `done`。
- 层级只由各目标 `00-meta.md` 的 `parent` 表达；本文件为区内索引，不是第二套状态源。
- R1：2026-07-31 `/govern` 阶段/关门自审 A-003（self）对 002/003/004 均为 **pass** → 三子目标 `done`；Root 纲领 R1 → 完成。Root A-001（independent）R1 证据复核 **pass** → A-002/D-006 响应：维持完成事实。
- R2：完成；用户书面确认按 v0.1.3 覆盖表正式冻结，D-009 / A-006 已留痕，`I-PROTO-001` = `verified`。冻结范围不等同于完整协议支持或 R3-R5 实施完成。
- R3：完成；GOAL-005 A-006 `source: self` / `verdict: pass` 已留痕，73 项测试、构建、固定 fixture 对照和 HTTP 入口证据已入账；A-007 independent 关门复审 `pass`，A-008 已响应其 recommended（F-001 索引表时点绑定、F-002 schema 等价性随 Root `I-PROTO-004` 跟进）。R4/R5 与完整协议支持仍在后续范围。
- R4：规划中；GOAL-006 已立项为 `active`，`I-006-001` 在方案冻结前验证，父目标 `I-PROTO-002` 为 R4 **实施**门禁。本条目反映规划事实，不放行 R4 实施。
