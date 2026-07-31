---
title: 目标树 · workspace-001-mvp-admin-foundation
status: active
created: 2026-07-31
updated: 2026-07-31
parent: null
version: 0.16.0
workspace_id: workspace-001-mvp-admin-foundation
---

# 目标树 · MVP Admin 基架

> 工作区：`workspace-001-mvp-admin-foundation`  
> canonical：`docs/workspace-001-mvp-admin-foundation/`  
> Root：`GOAL-001-mvp-admin-foundation`  
> primary_plan：`VP-001-mvp-admin-foundation`

## 树

```text
GOAL-001-mvp-admin-foundation  [active]  MVP Admin 基架  progress=4/6
├── GOAL-002-r1-repo-layout-conventions  [done]  R1 · 仓库布局与包管理约定
├── GOAL-003-r1-api-go-scaffold          [done]  R1 · Go API 工程骨架
├── GOAL-004-r1-web-react-scaffold       [done]  R1 · React Web 工程骨架
├── GOAL-005-r3-admin-shell-navigation   [done]  R3 · Admin 外壳与导航
├── GOAL-006-r4-account-permission       [done]  R4 · 核心账号与权限
└── GOAL-007-r5-examples-contract-verification [active]  R5 · 纳入域范例与契约验证
```

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-mvp-admin-foundation | MVP Admin 基架 | null | active | 4/6 | 2026-07-31 |
| GOAL-002-r1-repo-layout-conventions | R1 · 仓库布局与包管理约定 | GOAL-001-mvp-admin-foundation | done | — | 2026-07-31 |
| GOAL-003-r1-api-go-scaffold | R1 · Go API 工程骨架 | GOAL-001-mvp-admin-foundation | done | — | 2026-07-31 |
| GOAL-004-r1-web-react-scaffold | R1 · React Web 工程骨架 | GOAL-001-mvp-admin-foundation | done | — | 2026-07-31 |
| GOAL-005-r3-admin-shell-navigation | R3 · Admin 外壳与导航 | GOAL-001-mvp-admin-foundation | done | — | 2026-07-31 |
| GOAL-006-r4-account-permission | R4 · 核心账号与权限 | GOAL-001-mvp-admin-foundation | done | — | 2026-07-31 |
| GOAL-007-r5-examples-contract-verification | R5 · 纳入域范例与契约验证 | GOAL-001-mvp-admin-foundation | active | — | 2026-07-31 |

## 说明

- Root `progress: 4/6` = 纲领路线图 R1–R6 中 R1、R2、R3、R4 已完成（等权派生）；**不**放行 R5/R6、不关闭未相关 finding、不推导 Root `done`。
- R5：GOAL-007-r5-examples-contract-verification 已立项（`active`）。**阶段 1 契约发现与登记完成**（2026-07-31）：`I-007-001` 登记表落盘并核验 D-APP/D-PERM 复用验证命令，`I-007-001` → `verified`（登记层面）。**阶段 2 批次 2a 完成（D-DATA/D-TABLE）**（2026-07-31）：Go `GET /api/records`（list/detail）+ Web `records.ts`/`data-table.tsx` + `search-form-table`/`data-table` 范例页落地，`npm test` 114 项 / `go test` 21 项 / build / Edge 实测全绿；A-002 批次自审（self）pass。**阶段 2 批次 2b 完成（D-FORM/D-ACT）**（2026-07-31）：D-FORM §5 白名单控件表面 + D-ACT 非批量动作（复用 R4 `executeAction`）+ Go `PATCH`/`DELETE /api/records/{id}` + `form-controls`/`list-edit-lifecycle` 范例页落地，`npm test` 138 项 / `go test` 18 顶层 / build 全绿；A-003 批次自审（self）pass。批次 2c（D-EXPR/Renderer）与阶段 3/4 未开始。R5 **验收/关门**受父目标 `I-PROTO-003`（required，验收前闭合）门禁约束；阶段 1、批次 2a、批次 2b 均不抬升 progress。
- 层级只由各目标 `00-meta.md` 的 `parent` 表达；本文件为区内索引，不是第二套状态源。
- R1：2026-07-31 `/govern` 阶段/关门自审 A-003（self）对 002/003/004 均为 **pass** → 三子目标 `done`；Root 纲领 R1 → 完成。Root A-001（independent）R1 证据复核 **pass** → A-002/D-006 响应：维持完成事实。
- R2：完成；用户书面确认按 v0.1.3 覆盖表正式冻结，D-009 / A-006 已留痕，`I-PROTO-001` = `verified`。冻结范围不等同于完整协议支持或 R3-R5 实施完成。
- R3：完成；GOAL-005 A-006 `source: self` / `verdict: pass` 已留痕，73 项测试、构建、固定 fixture 对照和 HTTP 入口证据已入账；A-007 independent 关门复审 `pass`，A-008 已响应其 recommended（F-001 索引表时点绑定、F-002 schema 等价性随 Root `I-PROTO-004` 跟进）。R4/R5 与完整协议支持仍在后续范围。
- R4：方案冻结完成（GOAL-006 D-004，`I-006-001` 与父目标 `I-PROTO-002` 均 `verified`）；R4 **实施完成**（2026-07-31）：Go 会话与 `/api/accounts/me`、Go 独立鉴权、Web `$context` 挂载、D-PERM 求值引擎与 17 例 fixture 对照（13 valid 求值 + 4 invalid 错误码；13 valid 中含 5 例 execution 时序断言）落地，`go test`/`go build`、web 94 项测试、`npm run build`、HTTP 运行时与代理联调证据入账；A-001（self）与 A-002（independent）实施阶段均 pass，A-003 已响应（F-001/F-002 fixed，F-003～F-005 跟踪）。R4 **关门完成**（2026-07-31）：A-004 关门自审（self）与 A-005 独立关门复审（independent）均 pass、开放 required=0；经用户 `/govern` 授权 GOAL-006 → `done`、Root 纲领 R4 检查点完成（`progress` → 4/6）、goal-tree 同步。F-002～F-004（Renderer 接线 / token 会话 / 双端一致性 oracle）为 recommended 跟踪项，随 R5 / 生产化 / `I-PROTO-004` 解决。
