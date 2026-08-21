---
id: GOAL-017-w14-rectification-batch-c
doc: decision
status: active
parent: GOAL-015-w14-user-perspective-review
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# D-001 · GOAL-017 S1 方案冻结（F-08～F-10）

## 决策

### F-08 · 移除 pageId/route 技术信息框

- 直接删除 `App.tsx` 页面标题右侧的 `pageId` + `route` 调试框（D-003 已冻结）。
- 不做替代展示；页面标题区域保留面包屑 + 标题。

### F-09 · 反馈 toast 去错误码前缀 + 本地化

- Toast 不再显示 `ERROR_CODE: ` 前缀；错误码通过 `data-feedback-code` 属性与 `title` 保留为可访问/调试信息。
- renderer 中用户可见的硬编码英文反馈改为 i18n `error.*` 键，包含：
  - `error.actionNotFound` / `error.actionNotRequest` / `error.actionNotExecuted` / `error.batchActionNotExecuted`
  - `error.emptySelection` / `error.customHandlerNotFound` / `error.customHandlerMissingRowId`
  - `error.invalidNavigateUrl` / `error.rowNavigationFailed` / `error.navigateNoUrl`
  - `error.uploadNotUploadAction` / `error.uploadRequiresAction`
  - `error.tableColumnsRequired` / `error.tableDataSourceInvalid`
- 中英文键均落盘。

### F-10 · Schema 加载失败友好化

- `PageSchemaErrorSurface` 主标题不再显示 `error.code`，改为友好标题 + 说明 + “重新加载”按钮。
- 原始 `error.code`、`error.url`、validation issues 折叠在 `<details>` “技术详情”中，保留调试可访问性。
- 新增 i18n：`shell.pageSchemaError.title` / `.description` / `.reload` / `.technicalDetails`。

## 信息项更新

| ID | 状态 | 说明 |
|----|------|------|
| I-001 | **closed** | F-09 错误码保留在 `data-feedback-code`/`title`，不再作为 toast 可见前缀 |
| I-002 | **closed** | F-10 恢复动作 = “重新加载页面”，技术信息折叠保留 |
