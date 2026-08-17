---
id: E-002
goal: GOAL-027-w16-rectification-batch-c
title: 批 C S2 实施（F05/F06/F09/F10）
status: completed
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# E-002 · 批 C S2 实施

## 1. 执行事实

- **日期**：2026-08-17
- **动作**：
  1. **F05 Cron 预览**：
     - 新增 `POST /api/scheduled-tasks/cron/preview`，返回 `{ description, nextRuns }`。
     - 前端新增 `cron-preview` custom component 并挂到定时任务页。
  2. **F06 监控自动刷新**：
     - 前端新增 `monitoring-auto-refresh` custom component，提供关闭/5s/10s/30s 并调用 `crud.reloadList()`。
  3. **F09 字典 Badge**：
     - `dict_entries` 新增 `badge_style`（migration 0039），CRUD schema/API 支持 `badgeStyle`。
     - `SchemaTableColumnSpec` 新增 `badgeStyleField`，字典标签列按样式渲染 Badge。
  4. **F10 页脚版权/备案号**：
     - `site_settings` 新增 `copyright_text`/`icp_number`（migration 0040），settings schema/API 支持。
     - `/api/branding` 与前端 `Branding` 透出页脚字段，Shell 渲染页脚。
- **测试**：
  - Go 新增 `TestCronPreviewEndpoint`、`TestDictEntryBadgeStylePersists`、`TestSettingsFooterFieldsPersist`。
  - Web 全量 vitest 1056/1056 通过，`tsc -b` 通过。

## 2. 证据

| 主张 | 路径 / 证据 |
|------|-------------|
| F05 API/UI | `handler/scheduledtasks.go`、`components/cron-preview.tsx` |
| F06 UI | `components/monitoring-auto-refresh.tsx` |
| F09 migration/UI | `datadictionary/migration`、`dictionary-entries.json`、`schema-table.tsx` |
| F10 migration/UI | `settings/migration`、`settings.go`、`App.tsx`、`branding.ts` |
