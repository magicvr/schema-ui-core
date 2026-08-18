---
id: A-006
goal: GOAL-024-w16-user-perspective-improvements
title: 编排响应 · A-004/A-005 冲突裁决与 A-005 required 闭合
source: self
date: 2026-08-18
verdict: pass
scope: 响应 A-004（independent · pass）与 A-005（independent · fail）的冲突与 findings
---

# A-006 · 编排响应 · A-004/A-005 冲突裁决与 A-005 required 闭合

## 1. 冲突说明

- **A-004**（independent · gemini-3.7-flash-high · **pass**）：主张 10 项全部落地、可关门。
- **A-005**（independent · grok-build grok-4.6 · **fail**）：主张 F02 预览鉴权不可用、F03 导入前端未落地，2 条 required（F-001/F-002）未闭合，不同意 A-003/A-004 关门结论。
- **P-004 裁决**：两条 independent 意见在关门结论上冲突。已展示冲突并给建议，用户书面选择 **「采纳 A-005，先修正再关门」**（2026-08-18）。

## 2. 响应对象

- A-004 F-001/F-002（recommended · 文档残留）
- A-005 F-001（required · F02 预览/复制链接）
- A-005 F-002（required · F03 模板下载 + 逐行错误）
- A-005 F-003（recommended · F04 调账二次确认）
- A-005 F-004（recommended · F05 Cron 防抖预览）
- A-005 F-005（recommended · 父目标台账收口）

## 3. 关闭证据表

| finding | 状态 | 证据 |
|---------|------|------|
| A-005 F-001 | **fixed** | `render.tsx` `library.preview`/`library.copyLink` 改为经 `authFetch` 拉取 blob 后打开/复制 object URL；新增 `download-behavior.test.tsx` `library.preview` 用例。 |
| A-005 F-002 | **fixed** | 新增 `import-template-download` 组件并挂到 users 页（`users.json`）；`render.tsx` 支持在 200 响应的 `fieldErrors` 上渲染行号/字段/原因列表并保持模态打开；新增 `importErrors.title` i18n。 |
| A-005 F-003 | **fixed** | `render.tsx` `FormView` 对 `adjust-wallet-form` 的负数/大额 `amountDelta` 增加 `window.confirm` 二次确认；新增 `schema.wallet.adjustConfirm` i18n。 |
| A-005 F-004 | **fixed** | `cron-preview.tsx` 增加 400ms 防抖自动预览；保留手动按钮。 |
| A-005 F-005 | **fixed** | `GOAL-024/00-meta.md` 子目标表 A 行改为 done 4/4；`01/02/03.md` 索引 frontmatter status 置 done；`02-execution.md` 补齐 E-003（R1～R3/S5）。 |
| A-004 F-001/F-002 | **fixed** | 上述 F-005 修复一并覆盖（00-meta 子表 + 索引 status）。 |

## 4. 回归证据

- Go：本轮未改 Go 业务代码；此前 `go test ./...` 通过，users/schema/docscheck 定向复跑通过。
- Web：`tsc -b` 通过；vitest **1057/1057** 通过（含新增 `library.preview` 用例）。

## 5. 结论

A-004 与 A-005 的关门冲突已由用户裁决「采纳 A-005，先修正再关门」；A-005 两条 required 及三条 recommended 均按上表 fixed；父目标台账收口完成。`GOAL-024` 维持 **done · 8/8**。
