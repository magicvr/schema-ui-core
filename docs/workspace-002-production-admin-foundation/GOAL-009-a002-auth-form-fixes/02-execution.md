---
title: 执行记录 · A-002 缺陷修复
status: active
created: 2026-08-03
updated: 2026-08-03
parent: GOAL-001-production-admin-foundation
version: 0.1.0
---

# 执行记录 · GOAL-009

## 2026-08-03 · 立项

- 用户按 P-004 裁决：A-002 三条 required 走 `fixed`；F-002-002/003 由本目标承载，F-002-001 归 `GOAL-010-a002-schema-adapter`（Root D-014）；recommended F-002-004~006 为本目标可选加分（是否纳入待用户决定）；self 审计延后至修复后随关门补。
- 建立五件套与 `attachments/`；goal-tree 已同步（Root A-002 响应波次）。
- 未修改任何产品代码；Root A-002 F-002-002/003 保持 `open`。
- **计划（非事实）**：实施 S1（表单提交门禁）→ S2（认证失效清理）→ S3（回归/构建证据）→ S4（关门审计与 Root finding 关闭）。

## 2026-08-03 · S1 已实施（表单提交门禁 · F-002-002）

- `apps/web/src/renderer/render.tsx` `FormView`：新增 `hasBlockingErrors`（`gate.errors` / `reaction.errors` 任一非空）；`handleSubmit` 开头拒绝提交（早于 search/default 分支）；提交按钮 `disabled={submitting || hasBlockingErrors}`。
- 回归测试 `apps/web/src/renderer/render.test.tsx`「RenderPage form submit gate（GOAL-009 S1 · F-002-002）」3 条：① 字段 gate 错误（textarea 缺 `form.controls.extended` → `FORM_CAPABILITY_REQUIRED`）→ 按钮禁用 + fetch spy 未调用；② reaction 错误（`$deps.*` 越界语法 → `REACTION_EXPRESSION_INVALID`）→ 按钮禁用 + fetch spy 未调用；③ 正对照（合法表单）→ 按钮可用，POST `/api/records` 恰好 1 次。
- 证据：`vitest run` 23 文件 / **461 用例全绿**（458 存量 + 3 新增）；`tsc -b` 干净；`vite build` 成功（5.81s）。S1 为纯 web 改动，`apps/api` 未涉及。
- Root A-002 F-002-002 仍 `open`（按计划于 S4 关门审计后以 `fixed` 闭合）。
- **计划（非事实）**：S2（认证失效清理）→ S3 回归（S1 已提前满足 web 侧）→ S4（关门审计与 Root finding 关闭）。

## 2026-08-03 · S2 已实施（认证失效状态清理 · F-002-003）

- `apps/web/src/account/auth-client.ts` `authFetch`：refresh 成功但**重试仍 401** → `clearTokens()` + 触发 `onAuthLost`（新 token 仍被拒 = 会话已失效，不再保留旋转后的凭据）。
- `login`：`/me` 失败（任意原因）→ **回滚已存 token**（`clearTokens()`）并以登录失败重新抛出，不再静默降级为「带 token + 空 features」的 `authenticated` 快照；`AuthContext.login` 拒绝 → `LoginPage` 呈现错误、不进入 shell。
- 回归测试 `apps/web/src/account/auth-client.test.ts` 3 条：① 重试 401 后 `lost` 恰好 1 次 + 双 token 已清；② `/me` 500 → `ME_FAILED` + 双 token 已清；③ `/me` 401 → refresh 401 → `ME_FAILED` + 双 token 已清。
- 证据：`vitest run` 23 文件 / **464 用例全绿**（461 + 3 新增）；`tsc -b` 干净；`vite build` 成功（3.72s）。纯 web 改动，`apps/api` 未涉及。
- Root A-002 F-002-003 仍 `open`（按计划于 S4 关门审计后以 `fixed` 闭合）。
- **计划（非事实）**：S3 回归（web 侧已全绿，API 未涉及）→ S4（关门审计与 Root finding 关闭）。

## 2026-08-03 · S3 已实施（回归与构建证据）

- web 全量回归（当前 revision）：`vitest run` 23 文件 / **464/464** 全绿（含 S1/S2 共 6 条新增用例）；`tsc -b` + `vite build` 干净（3.17s）。
- `apps/api` 基线确认：`go test ./... -count=1` 全绿（cmd/server、account、auth、config、handler、store 等 7 包；handler 6.888s / store 5.223s）；`go vet ./...` 干净。S1/S2 为纯 web 改动，API 无代码变更，本项为基线回归确认（「若涉及」口径：不涉及，但已确认基线未破坏）。
- Root A-002 F-002-002/003 仍 `open`。
- **计划（非事实）**：S4（关门审计与 Root finding 关闭）——建议先 `/audit` finding-closure 复审后再关门。

## 2026-08-03 · S4 已实施（关门审计与 Root finding 关闭）

- **A-001（self · close-out · pass）** 写入本目标 `03-audit.md`：S1（门禁代码 + 3 回归）、S2（清理代码 + 3 回归）、S3（464/464 + build + go test/vet）证据链逐一复核成立；无 I-00N、无本区开放意见；结论支持按 `fixed` 闭合 F-002-002/003（D-014 已裁决 fixed 路径）。
- **GOAL-009 置 `done`（4/4）**：成功标准 S1～S4 全勾选。
- **Root 03-audit A-002 关闭证据表更新**：F-002-002 → `fixed`（证据：GOAL-009 S1 + A-001）、F-002-003 → `fixed`（证据：GOAL-009 S2 + A-001）；F-002-001 仍 `open`（GOAL-010 承载）；recommended F-002-004~006 仍 open 非阻断。Root 与 VP-002 关门因 F-002-001 未闭合继续阻断。
- 独立 `/audit` finding-closure 复审：可选加固，未执行（关门检查清单以 self 关门审计满足；如需可后续补跑）。

## 2026-08-03 · 响应 A-002（independent · close-out · conditional）

- **A-002**（`$audit` · close-out · conditional）确认 S1～S3 成立、F-002-002/003 代码关闭证据充分；开 **F-001（required / medium）**——Root A-002 正式意见索引仍写「F-002-001~003 仍 open」，与关闭证据表矛盾；**R-001（recommended / low）**——补可复现 revision 身份。
- **响应（用户裁决走 fixed）**：Root `03-audit.md` 正式意见索引 A-002 行已同步为单义现状（F-002-002/003 `fixed`、F-002-001 `open`、recommended 非阻断）；**F-001 → `fixed`**（索引 ↔ 关闭证据表 ↔ goal-tree 注记三处一致）。**R-001 → handled**（HEAD `5e08489` + 9 个未提交修改已记录；不冒充 clean revision/CI；Root/VP-002 关门前须先 commit）。
- A-002 与 A-001（self · pass）在 S4 结论上趋同，**GOAL-009 维持 `done / 4/4`**；本目标 scope 无开放 required。Root F-002-001 仍 open（GOAL-010 载体），Root 与 VP-002 关门继续阻断。
