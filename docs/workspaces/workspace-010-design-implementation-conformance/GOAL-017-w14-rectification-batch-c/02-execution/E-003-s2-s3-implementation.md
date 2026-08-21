---
id: GOAL-017-w14-rectification-batch-c
doc: execution
status: active
parent: GOAL-015-w14-user-perspective-review
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# E-003 · S2/S3 实施与回归

## 事实

- **2026-08-17**：F-08 实施——`App.tsx` 删除 `pageId` + `route` 调试框。
- **2026-08-17**：F-09 实施——`render.tsx` toast 去掉错误码前缀，增加 `data-feedback-code`/`title`；`render.tsx`/`schema-table.tsx` 硬编码英文反馈改为 i18n `error.*` 键；en/zh 新增键。
- **2026-08-17**：F-10 实施——`App.tsx` `PageSchemaErrorSurface` 改为友好标题 + 重新加载按钮 + 技术详情折叠；新增 `shell.pageSchemaError.*` i18n 键。
- **2026-08-17**：S3 回归——`npx tsc -b` 通过；Web 全量待最终跑（批 C 已跑相关 75 项测试）；如后续跑全量记录于 E-004。
