---
id: GOAL-010-w9-branding-asset-upload
doc: execution
status: recorded
parent: GOAL-001-design-implementation-conformance
created: 2026-08-15
updated: 2026-08-15
version: 0.1.0
---

# E-006 · 关门后用户跟进：上传字段图片预览

2026-08-15 用户跟进（关门后）：设置页上传 logo 后应显示图片预览，而非「Value: /api/branding/assets/…」文本。

- UploadField（form-controls.tsx）：URL 形态的值（同源路径 / http(s)）渲染缩略图预览；加载失败（onError）自动回退原文本；值变更时重置失败态；移除按钮保持不变。
- 非 URL 值（裸 id 等）与多选字段仍为文本形态；javascript:/data: 不进入预览。
- 测试：form-controls.upload.test.tsx 6 例（预览 + 移除 + onError 回退 + 非 URL 文本 + 空值/多选/readOnly）；representative-pages 上传流断言改为预览 img。
- 验证：tsc -b 0 错误；vitest 全量 969/969。
