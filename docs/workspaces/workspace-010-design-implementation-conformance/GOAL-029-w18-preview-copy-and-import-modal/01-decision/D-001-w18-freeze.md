---
id: D-001
goal: GOAL-029-w18-preview-copy-and-import-modal
title: W18 方案冻结：预览弹窗/复制链接与导入模态模板
status: accepted
created: 2026-08-18
updated: 2026-08-18
version: 0.1.0
parent: GOAL-001-design-implementation-conformance
---

# D-001 · W18 方案冻结

## 1. 触发

用户点名 `/govern` 处理 GOAL-024 A-007 F-001 / F-002（recommended open）。GOAL-024 保持 `done`。

## 2. 决定

### 2.1 F-001 · 预览与复制

- **预览**：在用户手势内同步 `window.open("about:blank")`，`await` 鉴权拉取 blob 后再 `location.replace(objectUrl)`。弹窗被拦则失败返回，不事后 `window.open`。
- **释放**：预览成功后 `setTimeout` 撤销 object URL（约 60s）；弹窗失败立刻 revoke。
- **复制**：写入 `new URL(downloadPath, location.origin).href`，不再复制 `blob:`。该 URL 仍受 Bearer 门禁，**不是**对外免登分享链。
- **不做**：Lightbox；带签名/query token 的公开下载。

### 2.2 F-002 · 导入模态

- 用户导入表单 `file` 字段挂 `afterComponent: "import-template-download"`（沿用 W17 本地扩展）。
- 移除页面级 `import-template-block`。
- `!response.ok` 显示可见错误，不再静默返回。
- 补 vitest：导入 200 `fieldErrors` 渲染 `data-import-error-rows`。

## 3. 未选方案

| 方案 | 未选理由 |
|------|----------|
| 继续复制 blob URL | A-007 明确不可分享 |
| 新增签名下载端点 | 超出 recommended 波次 |
| Lightbox | 原冻结有，但 A-007 不阻断；本波不扩 |

## 4. 后续

S2 实施；S3 Web 定向 + tsc。
