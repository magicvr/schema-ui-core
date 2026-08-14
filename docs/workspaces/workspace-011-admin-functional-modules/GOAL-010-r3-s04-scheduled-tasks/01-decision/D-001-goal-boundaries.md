---
id: D-001
goal: GOAL-010-r3-s04-scheduled-tasks
title: 立项边界：模块身份、Profile 归属与审计策略
date: 2026-08-14
status: accepted
parent: GOAL-010-r3-s04-scheduled-tasks
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# D-001 · 立项边界（S-04 定时任务）

## 决定

1. **模块身份**：`admin.scheduled-tasks`（标准 Admin 功能模块）；Descriptor 依赖 core.auth-session / core.schema-render / core.navigation-capability / core.operationlog。
2. **Profile 归属（I-003 闭合）**：进入 **admin 默认集**（内容扩展，先例一致）；mvp/demo 不启用。
3. **审计策略**：调度器 + 运行记录为 data 面 → 独立审计（grok，用户书面偏好统一沿用）。
4. **迁移**：0021 = scheduled_tasks + task_runs（归属 admin.scheduled-tasks）；0022 = operation_log CHECK + 3 个 scheduled-tasks 事件（归属 core.operationlog）。
5. **cron 校验（I-001 闭合）**：**自研 5 字段校验器**（minute hour dom month dow；支持 *、数字、步进 */n、列表 a,b；范围与语义校验）——go.mod 无 cron 依赖，引入外部库与仓库依赖纪律不符。
