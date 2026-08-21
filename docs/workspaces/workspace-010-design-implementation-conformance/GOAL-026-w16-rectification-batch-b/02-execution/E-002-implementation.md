---
id: E-002
goal: GOAL-026-w16-rectification-batch-b
title: 批 B S2 实施（F02/F03/F04）
status: completed
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# E-002 · 批 B S2 实施

## 1. 执行事实

- **日期**：2026-08-17
- **动作**：
  1. **F02 文件预览/复制**：
     - 文件库列表/详情新增 `downloadUrl` 字段。
     - schema 增加 `preview` / `copyLink` 自定义 action；renderer 增加 `library.preview`（新标签页打开）与 `library.copyLink`（剪贴板写入）。
  2. **F03 导入模板与错误明细**：
     - 新增 `GET /api/import/{resource}/template`（users 返回 CSV 模板）。
     - 导入响应新增 `fieldErrors`（`rowNumber/field/reason`），同时保留 `errors` 兼容。
  3. **F04 金额格式化与调账警示**：
     - `SchemaTableColumnSpec` 增加 `format: "currency"`，钱包余额列使用分转元格式化。
     - 调账 `amountDelta` 输入为负数或绝对值 > 100000 时显示警示文案。
- **测试**：
  - Go 新增 `TestImportTemplateUsers`、`TestImportFieldErrors`、`TestFileLibraryRowsExposeDownloadUrl`。
  - Web 新增 `schema-table` 货币格式化测试；全量 vitest 1056/1056 通过。

## 2. 证据

| 主张 | 路径 / 证据 |
|------|-------------|
| F02 downloadUrl | `apps/api/internal/handler/filelibrary.go`、`file-library.json`、`render.tsx` |
| F03 template/fieldErrors | `apps/api/internal/handler/import.go`、`w16_batch_b_test.go` |
| F04 currency/warning | `apps/web/src/renderer/schema-table.tsx`、`form-controls.tsx`、`wallet.json` |
