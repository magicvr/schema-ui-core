---
title: 审计台账 · A-002 缺陷修复
status: active
created: 2026-08-03
updated: 2026-08-03
parent: GOAL-001-production-admin-foundation
version: 0.2.1
---

# 审计台账 · GOAL-009

## 正式意见索引

| 编号 | source | 日期 | scope | verdict | 状态 |
|------|--------|------|-------|---------|------|
| A-001 | self | 2026-08-03 | S1～S4 关门（F-002-002/003 修复与回归） | pass | 已响应（closed） |
| A-002 | independent | 2026-08-03 | S1～S4 关门及 Root F-002-002/003 关闭复核 | conditional | 已响应（F-001 **fixed**；R-001 handled；维持 `done / 4/4`） |

> 本目标为 A-002 响应载体；Root 层面的 A-002 意见与响应记录见 [Root 03-audit](../GOAL-001-production-admin-foundation/03-audit.md)。

## A-001 · self · close-out（2026-08-03）

- **scope**：GOAL-009 S1～S4——F-002-002（表单提交门禁）与 F-002-003（认证失效清理）的修复、回归与关闭证据。
- **source**：self；**verdict**：pass。
- **背景**：P-004 §3.1 自审裁决于 D-001 留痕（延后至修复完成后随关门补）；本轮按用户指令实施关门审计。

### 复核结论

- **S1（F-002-002）**：`apps/web/src/renderer/render.tsx` `FormView` 存在 `hasBlockingErrors`（gate/reaction 任一非空），`handleSubmit` 开头早退，提交按钮 `disabled={submitting || hasBlockingErrors}`；`render.test.tsx` 新增 3 条回归（字段 gate 错误、reaction 错误 → 按钮禁用 + fetch spy 未调用；正对照 → POST 恰好 1 次），「错误显示后请求未发出」断言成立。与 finding 建议关闭路径逐字对应。
- **S2（F-002-003）**：`apps/web/src/account/auth-client.ts` `authFetch` 在 refresh 成功但重试仍 401 时 `clearTokens()` + `onAuthLost?.()`；`login` 的 `/me` 失败回滚 token 并以原 `AuthError` 拒绝（`AuthContext.login` 拒绝 → `LoginPage` 呈现登录失败，不再进入 `authenticated`）；`auth-client.test.ts` 新增 3 条回归（重试 401 清凭据 + lost 1 次；`/me` 500 回滚；`/me` 401 + refresh 401 回滚）。与 finding 建议关闭路径逐字对应。
- **S3（回归与构建）**：web `vitest run` 23 文件 / **464/464**（含 S1/S2 共 6 条新增）；`tsc -b` + `vite build` 干净；`apps/api` `go test ./... -count=1` 全绿（7 包）+ `go vet ./...` 干净（基线确认；S1/S2 为纯 web 改动）。
- **信息与意见**：本目标无 I-00N；本区无其他开放审计意见。
- **结论**：成功标准 S1～S4 证据链成立（02-execution 时间线可指回）；**建议 Root A-002 F-002-002/003 按 `fixed` 合法闭合**（用户裁决 D-014 已定 fixed 路径）。无 required 残留，可关门。

### 响应（编排器 · 2026-08-03）

- 采纳 A-001 `pass`：GOAL-009 四检查点全勾选（`4/4`），置 `done`；Root 03-audit A-002 关闭证据表更新 F-002-002/003 → `fixed`。独立 `/audit` finding-closure 复审为可选加固，不阻断本目标关门（D-014 fixed 路径 + A-001 self 复核已闭环）。

## A-002 · independent · close-out（2026-08-03）

- **source**：independent
- **auditor**：Codex · `$audit`
- **类型 / scope**：close-out；GOAL-009 S1～S4，含 F-002-002（表单提交门禁）、F-002-003（认证失效清理）及 Root A-002 的 `fixed` 关闭记录；不审 GOAL-010、Root/VP-002 关门或 recommended F-002-004～006 的实施取舍。
- **verdict**：conditional

### 范围与证据边界

- 工作区绑定已核对：`workspace-002-production-admin-foundation` 的 canonical scope、Root、`primary_plan: VP-002-production-admin-foundation` 与 Charter 链一致；本目标 `parent`、`done / 4/4` 在 `00-meta.md` 与 `goal-tree.md` 一致。
- 代码与文档证据位于当前未提交工作树（基准 HEAD `5e084893ed2eb3604ab640d75c8a689d8efd6219`）；本意见只确认该快照，不把它扩写为 clean revision、CI 或发布验收。
- 本目标声明无 I-00N；未发现影响本目标关门的 required 信息项。Root F-002-001 归 GOAL-010，继续阻断 Root/VP-002 关门，但不否定本目标 S1/S2 的实现事实。

### 成果与成功标准

- **S1 / F-002-002 成立**：`apps/web/src/renderer/render.tsx` 的 `FormView` 汇总 `gate.errors` / `reaction.errors` 为 `hasBlockingErrors`，在 search/default 分支前早退，并同步禁用提交按钮；`apps/web/src/renderer/render.test.tsx` 覆盖字段门禁错误、reaction 错误的「显示错误 + 禁用 + 不发请求」及合法 POST 正对照。
- **S2 / F-002-003 成立**：`apps/web/src/account/auth-client.ts` 在 refresh 成功但重试仍 401 时清除双 token 并触发一次 `onAuthLost`；login 后 `/me` 失败清 token 并拒绝登录；`apps/web/src/account/auth-client.test.ts` 对两类路径有直接回归。
- **S3 成立**：本轮独立重跑 `apps/web` `npm test -- --run`（23 files / **464 passed**）与 `npm run build`；`apps/api` `go test ./... -count=1` 与 `go vet ./...`；均 exit 0。聚焦测试亦分别通过 `render.test.tsx` 10/10、`auth-client.test.ts` 12/12。
- **S4 有条件成立**：Root A-002 关闭证据表已把 F-002-002/003 写为 `fixed`，后文也明确两项合法闭合；但同一 Root 审计台账的正式意见索引仍称 F-002-001～003 全部 open，当前状态表述不一致。

### Findings

#### F-001 · Root A-002 索引与关闭证据表状态矛盾

- **级别 / 严重度**：required / medium
- **影响门禁**：GOAL-009 close-out 的 S4「Root finding 按 `fixed` 合法闭合」及后续 Root 审计台账汇总。
- **证据**：[Root `03-audit.md`](../GOAL-001-production-admin-foundation/03-audit.md) 正式意见索引 A-002 行仍写「F-002-001～003 仍 open」；同文件「关闭证据表」和「后续」则写 F-002-002/003 已于 2026-08-03 `fixed`，仅 F-002-001 open。
- **要求**：由 `/govern` 将 Root A-002 索引状态同步为单义现状（F-002-001 open；F-002-002/003 fixed；recommended 仍非阻断），并复核索引、响应表与 goal-tree 摘要一致。不得改写 A-002 原始 finding 或把 F-002-001 误判为已闭合。

#### R-001 · 为当前修复补可复现 revision 身份

- **级别 / 严重度**：recommended / low（非阻断）
- **证据**：本轮开始与结束前的工作树均包含 9 个相关未提交修改；测试和构建可在当前工作树复现，但尚无 clean revision / commit 身份。
- **建议**：进入 Root/VP-002 关门前，由治理流程记录包含本修复的 revision 身份与 clean-tree/CI 边界；不要把本地工作树验证表述为 CI 或发布验收。

### 必改项汇总

- **F-001（required / medium）**：同步 Root A-002 正式意见索引与其关闭证据表；闭合前，本次独立意见不支持无条件维持 GOAL-009 close-out 结论。

### 与既有意见的异同

- 与 A-001（self · pass）同意 S1～S3 的实现与回归结论，也同意 F-002-002/003 的代码关闭证据充分。
- 与 A-001 的差异仅在 S4 台账一致性：A-001 支持立即关门，本意见因 Root A-002 索引仍显示两项 open 而判 `conditional`。该 required 差异须由 `/govern` 按 P-003/P-004 汇总响应，不由本审计静默覆盖既有状态。

### 结论与下一步

产品修复和当前回归均通过；关门记录仍有一处 required/medium 的单义性缺口。建议 `/govern` 走 `fixed`：同步 Root A-002 索引后核对 F-001 关闭证据，再决定是否维持 GOAL-009 `done / 4/4`；GOAL-010 与 Root/VP-002 关门继续按原门禁独立处理。

### 声明

本意见仅追加独立审计记录，不修改目标 `status`、检查点、`progress`、Root 审计响应、`goal-tree.md` 或 VP 状态；finding 响应与生命周期决策由 `/govern` 处理。

## A-002 响应（编排器 · 2026-08-03）

- **响应编号**：A-002（independent · close-out · conditional）
- **响应来源**：`/govern` 编排器（self 侧记录，非独立审）
- **用户裁决**（P-004）：按 `fixed` 同步 Root A-002 索引并复核 F-001 关闭证据（书面指令；同 scope 已有 A-001 self close-out 覆盖，P-004 §3.1 无需再补自审）。

### F-001（required / medium）→ fixed

- **修正**：Root `03-audit.md` 正式意见索引 A-002 行已同步为单义现状——**F-002-002/003 `fixed`（2026-08-03）**、**F-002-001 仍 open（载体 GOAL-010）**、F-002-004~006 recommended 非阻断。未改写 A-002 原始 finding，未把 F-002-001 误判为已闭合。
- **关闭证据复核**：索引 ↔ 同文件「关闭证据表」（F-002-002/003 fixed、F-002-001 open）↔ `goal-tree.md` 2026-08-03 关门注记（同口径）三处一致；Root 响应节「仍开放项/后续」亦已同步。**F-001 按 `fixed` 合法闭合。**

### R-001（recommended / low）→ handled

- revision 身份已记录：HEAD `5e08489`（docs(goal): 新增认证表单修复与 schema 适配目标并同步目标树）+ **9 个未提交修改**（`auth-client.ts`/`auth-client.test.ts`/`render.tsx`/`render.test.tsx` + 5 个治理文档）。本地工作树验证**不冒充** clean revision / CI / 发布验收；进入 Root/VP-002 关门流程前须先提交形成 clean revision（届时按 R-001 补记）。

### 结论

- F-001 闭合后，A-002（conditional）与 A-001（self · pass）在 S4 结论上趋同：S1～S3 实现与回归成立、F-002-002/003 关闭证据充分、台账单义——**GOAL-009 维持 `done / 4/4`**。无残留冲突、无开放 required（本目标 scope）。
- Root F-002-001 仍 open（GOAL-010 实施中），Root 与 VP-002 关门继续阻断；recommended F-002-004~006 仍待用户决定。
