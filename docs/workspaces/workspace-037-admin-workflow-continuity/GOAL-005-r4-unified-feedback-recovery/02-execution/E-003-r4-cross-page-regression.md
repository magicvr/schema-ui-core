---
id: E-003-r4-cross-page-regression
doc: execution-entry
status: recorded
goal_id: GOAL-005-r4-unified-feedback-recovery
created: 2026-09-17
updated: 2026-09-17
parent: GOAL-005-r4-unified-feedback-recovery
version: 1.0.0
---

# E-003 · R4 跨页面恢复回归

## 已发生事实

- SchemaTable 读取失败显示分类后的不可用文案，并保留显式 `[data-table-retry]` 路径；statCard 失败后重试一次可恢复展示值，retry 不形成自动循环。
- FormInner 的 network/write failure 会重新启用提交按钮并保留 `role=alert`、诊断 code 与 dirty 值；没有自动写 retry。既有 error-localization、schema-crud 与 R3 dirty-state 回归继续通过。
- 既有代表性 Admin 页面 unreachable 回归已更新为 R4 catalog 文案，Host failure fixtures 与页面 schema/manifest recovery 测试继续通过，未把普通 feedback 合并进 HostFailureScreen。
- 受影响定向回归：8 个测试文件、136 个测试通过；前端全量回归：110 个测试文件、1401 个测试通过；`npx tsc -p tsconfig.app.json --noEmit` 通过。

## 证据路径

- `apps/web/src/components/ui/feedback.test.tsx`
- `apps/web/src/renderer/feedback-policy.test.ts`
- `apps/web/src/renderer/schema-table.test.tsx`
- `apps/web/src/renderer/render.test.tsx`
- `apps/web/src/app/representative-pages.integration.test.tsx`
- `apps/web/src/renderer/r3-dirty-state.ui.test.tsx`
