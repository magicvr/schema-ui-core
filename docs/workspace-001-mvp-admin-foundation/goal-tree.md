---
title: 目标树 · workspace-001-mvp-admin-foundation
status: active
created: 2026-07-31
updated: 2026-08-01
parent: null
version: 0.19.0
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
| GOAL-007-r5-examples-contract-verification | R5 · 纳入域范例与契约验证 | GOAL-001-mvp-admin-foundation | active | — | 2026-08-01 |

## 说明

- Root `progress: 4/6` = 纲领路线图 R1–R6 中 R1、R2、R3、R4 已完成（等权派生）；**不**放行 R5/R6、不关闭未相关 finding、不推导 Root `done`。
- R5：GOAL-007-r5-examples-contract-verification（`active`）。阶段 1–2 完成（见历史条目）。**A-006（independent, pass）响应 + 阶段 3 完成**（2026-08-01）：用户裁决「不需要自审，直接推进」；`I-PROTO-004`=**vendor**（Root D-012 / GOAL-007 D-008）→ verified；schemas/fixtures vendor+SHA pin；Ajv 结构校验 + conformance 行为对照（`npm test` **326** 项 / build / go test 全绿）；登记表 v0.6.0。**阶段 4（验收/关门）未开始**；`I-PROTO-003` 仍 open。Root `progress` 仍 **4/6**（阶段 3 不抬升纲领检查点）。
- 层级只由各目标 `00-meta.md` 的 `parent` 表达；本文件为区内索引，不是第二套状态源。
- R1：2026-07-31 `/govern` 阶段/关门自审 A-003（self）对 002/003/004 均为 **pass** → 三子目标 `done`；Root 纲领 R1 → 完成。Root A-001（independent）R1 证据复核 **pass** → A-002/D-006 响应：维持完成事实。
- R2：完成；用户书面确认按 v0.1.3 覆盖表正式冻结，D-009 / A-006 已留痕，`I-PROTO-001` = `verified`。冻结范围不等同于完整协议支持或 R3-R5 实施完成。
- R3：完成；GOAL-005 A-006 `source: self` / `verdict: pass` 已留痕，73 项测试、构建、固定 fixture 对照和 HTTP 入口证据已入账；A-007 independent 关门复审 `pass`，A-008 已响应其 recommended（F-001 索引表时点绑定、F-002 schema 等价性随 Root `I-PROTO-004` 跟进）。R4/R5 与完整协议支持仍在后续范围。
- R4：方案冻结完成（GOAL-006 D-004，`I-006-001` 与父目标 `I-PROTO-002` 均 `verified`）；R4 **实施完成**（2026-07-31）：Go 会话与 `/api/accounts/me`、Go 独立鉴权、Web `$context` 挂载、D-PERM 求值引擎与 17 例 fixture 对照（13 valid 求值 + 4 invalid 错误码；13 valid 中含 5 例 execution 时序断言）落地，`go test`/`go build`、web 94 项测试、`npm run build`、HTTP 运行时与代理联调证据入账；A-001（self）与 A-002（independent）实施阶段均 pass，A-003 已响应（F-001/F-002 fixed，F-003～F-005 跟踪）。R4 **关门完成**（2026-07-31）：A-004 关门自审（self）与 A-005 独立关门复审（independent）均 pass、开放 required=0；经用户 `/govern` 授权 GOAL-006 → `done`、Root 纲领 R4 检查点完成（`progress` → 4/6）、goal-tree 同步。F-002～F-004（Renderer 接线 / token 会话 / 双端一致性 oracle）为 recommended 跟踪项，随 R5 / 生产化 / `I-PROTO-004` 解决。
