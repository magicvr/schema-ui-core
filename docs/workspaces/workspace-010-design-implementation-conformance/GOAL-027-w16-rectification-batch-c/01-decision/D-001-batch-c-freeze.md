---
id: D-001
goal: GOAL-027-w16-rectification-batch-c
title: W16 批 C 方案冻结（F05 / F06 / F09 / F10）
status: approved
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# D-001 · W16 批 C 方案冻结

## 1. 触发

GOAL-024 D-002 已给出 W16 技术方案；批 A/B 已 done，按 D-003 渐进添加批 C。本决策将 F05/F06/F09/F10 细化为可实施设计，并关闭本子目标 I-001/I-002。

## 2. W16-F05 · Cron 表达式自然语言解析与未来运行时间预览

- **API**：新增 `POST /api/scheduled-tasks/cron/preview`，请求 `{ "cron": "0 0 2 * * *" }`，响应 `{ "description": "...", "nextRuns": ["2026-08-18T02:00:00Z", ...] }`。
- **后端**：复用 `scheduledtasks/store` 的 5-field cron 解析，计算未来 3 次运行时间。
- **前端**：定时任务表单的 Cron 字段下方挂 custom component `cron-preview`，输入防抖调用该端点并展示描述与未来时间。

## 3. W16-F06 · 系统监控页面定时自动轮询刷新

- **前端**：监控页表格/指标节点增加可选 `refreshInterval?: number`（毫秒）；`SchemaTable` 支持该 prop 时按间隔静默 refetch。
- **UI**：监控页面顶部提供 `关闭 / 5秒 / 10秒 / 30秒` 下拉，写入页面局部状态并透传为节点 `refreshInterval`；不新增 API。

## 4. W16-F09 · 数据字典项 Badge 颜色/标签风格配置

- **存储/API**：`dict_entries` 表新增 `badge_style TEXT NOT NULL DEFAULT 'default'`；条目 CRUD schema/API 增加 `badgeStyle` 字段，取值 `default|success|warning|destructive|info`。
- **前端**：`SchemaTableColumnSpec` 增加可选 `badgeStyleField`；当列声明该字段时，单元格渲染为对应颜色的 Badge。

## 5. W16-F10 · 系统设置页脚版权文字与备案号

- **存储/API**：`settings` 资源增加 `copyrightText` 与 `icpNumber` 字段（通用设置 schema），`GET /api/settings` 与 `PATCH /api/settings/{id}` 自动覆盖。
- **前端**：Shell/`App.tsx` 底部统一读取设置并渲染页脚；未配置时隐藏。

## 6. 信息项关闭

| ID | 级别 | 结论 |
|----|------|------|
| I-001 | required | F05 响应结构冻结为 `{ description, nextRuns }`，前端 custom component 调用。 |
| I-002 | non-blocking | F10 使用现有 settings schema/API 增加两个字段，Shell 页脚读取。 |

## 7. 未选方案

- F06 不新增独立轮询 API；使用前端定时 refetch 现有端点。
- F09 不新增字典样式表/主题系统；仅扩展 `badgeStyle` 预设值。
