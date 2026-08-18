---
id: D-004
goal: GOAL-024-w16-user-perspective-improvements
title: 响应 A-007：A-005 F-001/F-002 记 fixed；A-005 F-004 保持 recommended open
status: accepted
created: 2026-08-18
updated: 2026-08-18
version: 0.1.0
parent: GOAL-001-design-implementation-conformance
---

# D-004 · 响应 A-007：A-005 required 闭合与 F-004 处置

## 1. 触发

- **A-007**（independent · grok-build grok-4.6 · **pass**，2026-08-18）：复审 A-005 F-001～F-005 关闭证据。判定 A-005 F-001 / F-002 required **可按 `fixed` 闭合**；**不同意** A-006 将 A-005 F-004 标 `fixed`。
- 用户书面指令（`/govern`，2026-08-18）：**A-005 F-001/F-002 记 `fixed`；A-005 F-004 保持 recommended open 或书面 residual**。

## 2. 决定

1. **接受 A-007 `pass`**。A-005 两条 required（F-001 / F-002）按 `fixed` 合法闭合；闭合依据为 A-007 核对 + 本轮编排器对 `fbe7c40` 代码路径的复核。
2. **A-005 F-004 保持 `recommended` / `open`**。不写 `accepted-residual`：用户给出的是「open 或 residual」二选一授权，但未书面给出残余范围、期限与复审触发；按 P-003/P-004 不得编造 residual 条款。
3. **A-006 对 F-004 的 `fixed` 声明撤回**（见 A-008）。防抖增量仍成立，冻结方案「字段下方即时中文语义」未交付，不能整条闭合。
4. **A-007 自身 F-001～F-003**（均为 recommended）保持 open，不升格、不假装闭合。
5. **不重开 GOAL-024**。A-007 开放 required = 0，维持 `done · 8/8`。

## 3. 为什么

- A-005 F-001：`library.preview` / `library.copyLink` 已走鉴权 `fetcher` + blob object URL，不再裸 `window.open` 下载路径；`download-behavior.test.tsx` 有预览用例。原 401 + `attachment` 必改缺口已补。
- A-005 F-002：用户页有 `import-template-download`；导入 200 `fieldErrors` 会留在模态并以 `#行号 / 字段 / 原因` 列出。原「模板入口 + 逐行错误」required 已补。
- A-005 F-004：`cron-preview.tsx` 仍是独立输入；`scheduled-tasks.json` 仍挂页面块；`describeCron` 仍返回 `"every minute"` / `"every hour at minute N"` / `"cron schedule (5-field)"`。仅 400ms 防抖不足以标 fixed。

## 4. 未选方案

| 方案 | 未选理由 |
|------|----------|
| 把 A-005 F-004 标 `accepted-residual` | 用户未写残余范围、复审触发；默认走 recommended open，避免静默编造 residual |
| 当场补 Cron 字段绑定 + 中文描述并标 fixed | 超出本轮「响应 A-007」写入指令；F-004 非 required，不阻断维持关门 |
| 因 A-007 新 recommended 重开 GOAL-024 | A-007 明确不阻断维持关门；用户未要求重开 |

## 5. 影响

- 门禁：GOAL-024 无未闭合 required；不改变 `status` / `progress`。
- 台账：A-006 关闭表 F-004 行更正为 open；A-008 为正式响应节。
- 后续：A-005 F-004 / A-007 F-003 若要 residual 或整改，须另给书面范围。

## 6. 后续

- 落盘 A-008（self · response）与 E-004。
- 可选后续：`/audit` 复审本响应；或另立项补 Cron 字段绑定与中文 `describeCron`。
