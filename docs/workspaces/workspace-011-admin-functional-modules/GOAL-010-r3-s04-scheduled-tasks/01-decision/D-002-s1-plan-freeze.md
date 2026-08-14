---
id: D-002
goal: GOAL-010-r3-s04-scheduled-tasks
title: 方案冻结：定时任务设计（S1）
date: 2026-08-14
status: accepted
parent: GOAL-010-r3-s04-scheduled-tasks
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# D-002 · 方案冻结（S-04 定时任务）

## 1. 数据模型（migration 0021，admin.scheduled-tasks）

`scheduled_tasks`：id TEXT PK ｜ key TEXT UNIQUE NOT NULL ｜ cron TEXT NOT NULL（5 字段） ｜ name TEXT NOT NULL ｜ enabled INTEGER 1 ｜ description TEXT ｜ handler TEXT NOT NULL DEFAULT 'system.noop'（处理器注册键） ｜ created_at / updated_at INTEGER NOT NULL

`task_runs`：id TEXT PK ｜ task_id TEXT NOT NULL REFERENCES scheduled_tasks(id) ON DELETE CASCADE ｜ status TEXT NOT NULL CHECK (status IN ('ran','failed')) ｜ started_at INTEGER NOT NULL ｜ finished_at INTEGER ｜ detail TEXT ｜ created_at INTEGER NOT NULL

## 2. cron 语法（自研，I-001）

- 5 字段：minute(0-59) hour(0-23) dom(1-31) month(1-12) dow(0-6, 0=Sunday)。
- 支持：`*`、数字、`*/n` 步进、`a,b` 列表；dom 与 dow 同时约束（AND 语义）。
- 校验失败 → 400 INVALID_CRON；空/非法字段拒绝。
- 下一次运行时间：从 now 起向后搜索最近 5 年内匹配时刻（步进 1 分钟）；搜索失败视为任务不可调度（记录 detail）。

## 3. 调度器（I-002 闭合）

- **进程内循环**（best-effort 单实例语义文档化）：每 30s tick 一次，扫描 enabled 任务，对 next_run <= now 的任务执行并记录 task_runs 行。
- **处理器注册点**：模块内 map[string]TaskHandler（v1 内置 `system.noop`：仅记录 ran 状态 + 起止时间）；业务处理器由后续目标接入（文档化）。
- 重启语义：无持久化 next_run（每次 tick 按 cron 计算）；错过窗口不补跑（文档化）。
- 手动触发：POST /api/scheduled-tasks/{id}/run（tasks.write）→ 立即执行一次并记录 run 行。

## 4. 端点与权限

| 端点 | 门禁 | 说明 |
|------|------|------|
| GET/POST/PATCH/DELETE /api/scheduled-tasks (+batch-delete) | 读 tasks.read / 写 tasks.write | 任务定义 CRUD（工厂） |
| POST /api/scheduled-tasks/{id}/run | tasks.write | 手动触发（立即执行 + run 行） |
| GET /api/scheduled-tasks/{id}/runs | tasks.read | 单任务运行历史（自定义只读端点，倒序分页） |
| GET /api/scheduled-tasks/runs | tasks.read | 全局运行历史（只读资源工厂） |

- 权限键：`tasks.read / tasks.write`（PolicyAdmin）；导航 `menu_scheduled_tasks`（PageID scheduled-tasks）。
- 审计事件：`scheduled-tasks.create / scheduled-tasks.update / scheduled-tasks.delete`（0022 CHECK 扩展）。

## 5. 页面与 Schema

- 页面 `scheduled-tasks`：任务列表（key/cron/name/enabled/handler/updatedAt）+ 工具栏新建 + 行操作（编辑/删除/触发/启停?——v1 启用字段在编辑表单）+ 行操作「Runs」navigate 到运行历史页。
- 页面 `task-runs`：全局运行历史表（taskKey/status/startedAt/finishedAt/detail，倒序）。
- i18n zh/en：manifest.title.scheduledTasks / taskRuns + schema.* 键。
- cron 字段在表单为文本输入（服务端校验 INVALID_CRON）。

## 6. 测试与验证

- cron 校验器单测（合法/非法/边界）；调度器测试（tick 执行、noop 记录、手动触发、启停）；CRUD + 门禁 + 审计事件。
- 组合根：admin 权限 18→20、导航 10→11；迁移 20→22。
- web：fixture/schema-keys/s5-denominator/e2e admin 导航；冒烟 SM-007 + scheduled-tasks。

## 7. 未选方案

- 外部 cron 库（robfig/cron）：新依赖 + 仓库纪律不符（D-001 §5）。
- 持久化 next_run + 补跑：单实例 best-effort 语义 + 文档化接受（D-002 §3）。
- 真实业务处理器：v1 提供注册点 + noop 内置（业务处理器后续目标接入）。
