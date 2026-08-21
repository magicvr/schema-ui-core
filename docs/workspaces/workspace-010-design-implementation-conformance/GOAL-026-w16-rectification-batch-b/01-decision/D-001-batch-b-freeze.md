---
id: D-001
goal: GOAL-026-w16-rectification-batch-b
title: W16 批 B 方案冻结（F02 / F03 / F04）
status: approved
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# D-001 · W16 批 B 方案冻结

## 1. 触发

GOAL-024 D-002 已给出 W16 技术方案；批 A 已 done，按 D-003 渐进添加批 B。本决策将 F02/F03/F04 细化为可实施设计，并关闭本子目标 I-001/I-002。

## 2. W16-F02 · 文件库在线预览与一键复制直链

- **API**：`fileRow` 增加 `downloadUrl` 字段（`/api/library/files/{id}/download`），列表/详情直接返回；无需新增端点。
- **前端**：
  - 文件库表格操作列增加 `preview` 与 `copyLink` 两个 row action。
  - `preview`：图片类型打开 Lightbox（custom component `file-preview-lightbox`），PDF/文本 `window.open(downloadUrl, "_blank")`。
  - `copyLink`：`navigator.clipboard.writeText(downloadUrl)` + Toast。
- **Renderer**：通过 schema 表格 actions + custom component 实现，不新增协议。

## 3. W16-F03 · 数据导入 CSV 模板与逐行错误定位

- **API**：
  - 新增 `GET /api/import/{resource}/template`（当前仅 `users`）返回 CSV 模板，`Content-Type: text/csv`。
  - 导入响应在保留 `errors` 的同时新增 `fieldErrors`：`[{ rowNumber, field, reason }]`；`importRowError` 增加 `field` 可空字段，`errors` 继续兼容。
- **前端**：
  - 导入模态框增加“下载 CSV 模板”链接。
  - 校验失败后以表格/列表展示 `fieldErrors`（第几行、字段、原因）。

## 4. W16-F04 · 钱包金额“分转元”格式化与调账警示

- **API**：金额传输仍为分，不破坏契约。
- **前端**：
  - `SchemaTableColumnSpec` 增加可选 `format?: "currency"`；钱包 schema 的余额/冻结/流水金额列声明该格式，渲染时 `/100` + 千分位 + 两位小数。
  - 调账/扣款表单：当 `amountDelta` 绝对值超过阈值（如 100000 分）或为负数时，显示高亮警示与二次确认。

## 5. 信息项关闭

| ID | 级别 | 结论 |
|----|------|------|
| I-001 | required | F03 采用新增 `fieldErrors` + 保留 `errors` 的兼容升级，不破坏既有调用方。 |
| I-002 | non-blocking | F02 使用现有 `downloadUrl`（`/api/library/files/{id}/download`）即可，无需新增鉴权面。 |

## 6. 未选方案

- F02 不新增独立预览端点；下载端点已具备鉴权与文件读取能力。
- F03 不直接替换 `errors` 字段，避免破坏现有导入消费者。
