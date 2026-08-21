---
id: E-002
goal: GOAL-010-r3-s04-scheduled-tasks
date: 2026-08-14
status: recorded
parent: GOAL-010-r3-s04-scheduled-tasks
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# E-002 · S2 实现完成

## 事实

- 2026-08-14：S2 实现完成，覆盖 D-002 §1–§5：
  - **cron** store/cron.go：自研 5 字段解析（*、数字、*/n 步进、a,b 列表、a-b 范围）；Next 含当前分钟槽（inclusive），5 年窗口搜索。
  - **migration 0021**（admin.scheduled-tasks）：scheduled_tasks + task_runs（级联）；checksum 076c8fa3…；**0022**（core.operationlog）：CHECK + scheduled-tasks.create/update/delete；checksum cb64c05e…。
  - **store**：任务 CRUD + runs（任务级/全局）+ enabled 扫描。
  - **调度器** scheduler.go：30s tick 循环、分钟槽去重（lastRun 内存表）、system.noop 处理器、Execute 手动/调度共用、运行记录写入。
  - **handler** scheduledtasks.go：任务资源工厂（cron 校验 INVALID_CRON、TASK_KEY_TAKEN）+ 手动触发 + 单任务/全局运行历史；审计事件。
  - **模块** provider（五面贡献）+ schema ×2（任务页 + 运行历史页，navigate 复用）+ fragment；装配（profile/composition/testsupport/compiled）。
  - **web**：i18n zh/en；fixture（admin + sha 重钉 c61a3b3e…）；smoke admin 页面集 + scheduled-tasks。
  - **测试**：cron 单测（合法/非法/边界/闰年）、调度器（到点执行、分钟槽去重、手动触发）、handler（生命周期/门禁/审计）、provider（注册面 + 端到端）。
- 实现期修正：Next 语义从「严格大于」改为「含当前槽」（否则当前分钟槽永不执行）+ 分钟槽去重防双跑；全局 runs 路径从 /api/scheduled-tasks/runs 改为 /api/task-runs（避免与 {id} 路由冲突）。
