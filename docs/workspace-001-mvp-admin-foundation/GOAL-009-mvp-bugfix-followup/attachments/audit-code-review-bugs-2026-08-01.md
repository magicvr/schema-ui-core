---
title: 代码审视 · apps/api + apps/web 对照 VP-001（bug 与一致性）
status: active
doc_type: audit-attachment
source: independent
created: 2026-08-01
updated: 2026-08-01
parent: GOAL-009-mvp-bugfix-followup
version: 0.1.0
related_opinion: A-001
scope: apps/api + apps/web 实现 vs VP-001 意图；真实 bug / 集成失真 / 文档与稳健性
---

# 代码审视附件 · 2026-08-01

> **性质**：独立代码审视长文（`source: independent`）。  
> **正式台账索引**：本目标 [03-audit.md](../03-audit.md) **A-001**。  
> **实施后状态**：本附件 Findings 为 A-001 **时点基线**快照；F-009-001～007 实施后的闭合状态以 [03-audit.md](../03-audit.md) 闭合留痕与 **A-002** 关门复审为准（勿以附件 `状态: open` 误判未修）。  
> **范围**：对照 [VP-001](../../../../vision/plans/VP-001-mvp-admin-foundation.md) 已声明 MVP 边界与冻结表 I-PROTO-001 v0.1.3；**不是**完整协议宿主或生产 IAM 验收。  
> **验证基线（审视时）**：`go test ./...` 通过；Web Vitest **395/395** 通过。

## 1. 总评（对照 VP-001）

| 维度 | 结论 |
|------|------|
| 是否符合 VP-001 意图 | **基本符合**：可 fork 的 Admin 协议基架 + 范例驱动，非完整业务 Admin / 全量 schema 宿主 |
| 是否达到 VP 退出判据 | **达到已声明 MVP 边界**；与 VP 关门 residual 一致 |
| 是否有 bug | **有少量真实功能/一致性问题**；另有有意边界与文档滞后 |

一句话：与已关闭 VP-001 对齐度高；按「生产可用、schema 全驱动、后端强制鉴权」扩展标准则有明确缺口。本目标**只承接应修正的 bug / 失真 / 文档与最小稳健性**，不扩下一波次架构。

## 2. 对照 VP 三条退出判据（摘要）

1. **可运行可 fork + 固定协议边界**：符合（Go stdlib；Web 固定 2.7 manifest；pin `ca9e5fe…`；无未声明业务域）。
2. **纳入域有实现 + 范例 + 验证**：大体符合（范例驱动）；`schemaUrl` 通用 page 管线未做（有意/后续波次，**不在本目标必修**）。
3. **账号权限前后端集成**：最小闭环有；list-edit 演示失真、Go `Allow` 未挂写路由；浏览器拒绝路径弱（VP residual 已承认）。

## 3. Findings 全表（本目标范围）

编号 `F-009-00N` 供本目标实施与审计引用。级别：`required` = 本目标关门前应修；`recommended` = 建议修但不单独阻断若用户降级。

### F-009-001 · PATCH 不更新 `updatedAt`（required）

| 字段 | 值 |
|------|-----|
| **严重度** | med |
| **位置** | [`apps/api/internal/handler/records.go`](../../../../../apps/api/internal/handler/records.go) `update()` |
| **描述** | 成功 PATCH name/status/owner 后不刷新 `UpdatedAt`；`sort=updatedAt` 与“最近修改”语义错误 |
| **关闭要求** | 更新成功时写入当前 UTC 时间（或等价单调时间）；单测断言 PATCH 后 `updatedAt` 变化 |
| **状态** | open |

### F-009-002 · list-edit 权限门禁几乎恒过 + 硬编码 context（required）

| 字段 | 值 |
|------|-----|
| **严重度** | med |
| **位置** | [`apps/web/src/app/examples/list-edit-lifecycle-page.tsx`](../../../../../apps/web/src/app/examples/list-edit-lifecycle-page.tsx) |
| **描述** | (1) `context` 硬编码 `{ roles: ["admin"] }`，不用 boot 的 `/me`；(2) `PAGE_DOCUMENT` 仅有 `permissionIntent`，无 `permissions`/`permissionCascade`；(3) 引擎无 cascade/local 时 `effectivePermission` 默认 `true` → Edit/Delete 几乎总 `EXECUTED`，无法观察拒绝 |
| **关闭要求** | 接入真实 `navigationContext`（或从 `/me` 同源 snapshot）；在 page 文档上挂可失败的权限表达式；至少一条可核对的拒绝路径（单元或组件测） |
| **状态** | open |

### F-009-003 · Account 加载失败静默丢弃（required）

| 字段 | 值 |
|------|-----|
| **严重度** | low–med |
| **位置** | [`apps/web/src/main.tsx`](../../../../../apps/web/src/main.tsx) |
| **描述** | `loadAccountContext()` 的 `error` 被丢弃；API 失败时 shell 仍渲染且无提示；无 permissions 的导航仍全显 |
| **关闭要求** | 失败时至少有可观察 UI/日志面（不必全站阻断）；调用链保留 error；测覆盖失败路径展示或 prop |
| **状态** | open |

### F-009-004 · `sessionProvider` 注释与 nil 行为不符（required）

| 字段 | 值 |
|------|-----|
| **严重度** | low |
| **位置** | [`apps/api/internal/handler/account.go`](../../../../../apps/api/internal/handler/account.go) |
| **描述** | 注释写 “nil means static dev session”，实际 `h.sessionProvider()` 在 nil 时 panic |
| **关闭要求** | 修正注释与/或 nil 时回落 `StaticDevSession`（二选一须一致）；防回归测可选 |
| **状态** | open |

### F-009-005 · API/Web README 严重滞后（required）

| 字段 | 值 |
|------|-----|
| **严重度** | low |
| **位置** | [`apps/api/README.md`](../../../../../apps/api/README.md)、[`apps/web/README.md`](../../../../../apps/web/README.md)（及明显过时的 renderer/host 说明若存在） |
| **描述** | API README 仍写「仅 /healthz、无鉴权」；Web 侧仍暗示 R4/R5 未落地，与现状不符 |
| **关闭要求** | 更新为当前端点、session 模型、范例页与测试命令；标明 MVP 非生产鉴权边界 |
| **状态** | open |

### F-009-006 · records 写路由未调用 `account.Allow`（recommended）

| 字段 | 值 |
|------|-----|
| **严重度** | med（误用时 high） |
| **位置** | [`permission.go`](../../../../../apps/api/internal/account/permission.go) vs [`records.go`](../../../../../apps/api/internal/handler/records.go) |
| **描述** | Go 鉴权库存在但 PATCH/DELETE/GET 匿名可写；与「后端独立鉴权」叙事易混淆。VP/oracle 未强制真实 login 越权 HTTP |
| **关闭要求（若纳入）** | 至少对写操作挂 fail-closed 检查（dev session + 表达式或角色门槛）+ 测试；或书面 residual：明确“演示 API 无路由鉴权” |
| **状态** | open（本目标默认 **建议修**；用户可 residual） |

### F-009-007 · 无 body 大小限制 / pageSize 上限（recommended）

| 字段 | 值 |
|------|-----|
| **严重度** | low |
| **位置** | records handler |
| **描述** | 无 `MaxBytesReader`；`pageSize` 无上限（当前 8 条 seed 无害） |
| **关闭要求（若纳入）** | 合理上限 + 非法参数 400 测试 |
| **状态** | open |

## 4. 明确不在本目标范围（有意边界 / 下波次）

| 项 | 理由 |
|----|------|
| `schemaUrl` → 通用 page 加载管线 | 架构能力，非本波 bug；需新 VP/阶段 |
| `host/` 空占位补全 | 同上 |
| 真实 login/logout/token/IAM | R4 可选；非 VP residual 必改 |
| D-UPLOAD / batch action / multi-round `$deps` | I-PROTO-001 排除 |
| 浏览器多浏览器 e2e 矩阵 / 视觉回归 | 增强项 |
| Root/VP 重开或改 VP status | 本目标不改 VP；Root 保持 `done`，本子目标为关门后修正跟随 |

## 5. 建议实施顺序（供 GOAL-009）

| 序 | Finding | 建议 |
|----|---------|------|
| 1 | F-009-001 | API PATCH `updatedAt` + test |
| 2 | F-009-002 | list-edit 真实 context + 权限表达式 + 拒绝路径测 |
| 3 | F-009-003 | account error 可观察 |
| 4 | F-009-004 | nil/注释一致 |
| 5 | F-009-005 | README |
| 6 | F-009-006/007 | 按用户是否纳入 recommended |

## 6. 证据与方法说明

- 源码通读：`apps/api`（cmd/internal/pkg）、`apps/web/src`（account/app/protocol/renderer）。
- 对照文档：VP-001、I-PROTO-001 v0.1.3、goal-tree（R1–R6 done）。
- 运行验证：审视当日 `go test ./...`；`npm test -- --run` → 395 passed。
- 本附件**不**改任何目标 `status`/`progress`；响应与实施归 `/govern` 与本目标执行记录。
