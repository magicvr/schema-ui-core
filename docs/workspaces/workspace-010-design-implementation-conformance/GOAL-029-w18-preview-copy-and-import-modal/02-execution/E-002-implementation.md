---
id: E-002
goal: GOAL-029-w18-preview-copy-and-import-modal
title: S2 实施预览/复制与导入模态模板
status: completed
created: 2026-08-18
updated: 2026-08-18
version: 0.1.0
parent: GOAL-001-design-implementation-conformance
---

# E-002 · S2 实施（2026-08-18）

## 已发生事实

1. `library.preview`：手势内 `window.open("about:blank")`，鉴权拉 blob 后 `location.replace`；60s 后 `revokeObjectURL`；弹窗被拦返回 `POPUP_BLOCKED`。
2. `library.copyLink`：写入源站绝对 download URL，不再复制 `blob:`。
3. `users.json` 导入表单 `file.afterComponent = import-template-download`；移除页面 `import-template-block`。
4. 模板下载失败展示 `data-import-template-error`。
5. 补定向测试：预览开窗、复制绝对 URL、200 `fieldErrors` 行列表、模板失败提示。

## 下一步（计划）

S3 记录测试结果。
