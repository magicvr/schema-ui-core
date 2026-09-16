---
id: E-004-r4-f001-transport-classification
doc: execution-entry
status: recorded
parent: GOAL-005-r4-unified-feedback-recovery
created: 2026-09-17
updated: 2026-09-17
version: 1.0.0
goal_id: GOAL-005-r4-unified-feedback-recovery
---

# E-004 · 响应 A-002 F-001：保留 transport 分类

## 事实

2026-09-17，针对 A-002 independent 的 required F-001，已按 `fixed` 路径修正 `apps/web/src/renderer/render.tsx`：

- `runRequest`、`runBatchRequest`、custom preview/download 和 recordSource 的 transport catch 不再把原始错误统一压成 `REQUEST_FAILED`；统一通过 `transportFailureResult` 保留 `AbortError` → `feedback.timeout` 与网络 `TypeError`/`Failed to fetch` → `feedback.offline` 的分类。
- FormInner 的防御性 submit catch 继续直接交给 `feedbackFromError`；写入路径不提供 retry，recordSource 读取仍只提供用户触发的显式 retry。
- 新增表单/action AbortError 超时（保留值、无 retry）和 recordSource AbortError/离线分类及单次显式 retry 回归。

## 验证

- 受影响的 4 个测试文件、92 项通过：`render.test.tsx`、`feedback-policy.test.ts`、`feedback.test.tsx`、`schema-table.test.tsx`。
- `tsc -p tsconfig.app.json --noEmit` 通过。
- 本条只记录 F-001 响应事实，不修改 A-002 的原始 verdict；A-003 independent recheck 待执行。
