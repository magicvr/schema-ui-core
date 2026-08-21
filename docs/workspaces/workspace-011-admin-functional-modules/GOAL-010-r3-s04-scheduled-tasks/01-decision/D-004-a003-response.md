---
id: D-004
goal: GOAL-010-r3-s04-scheduled-tasks
title: A-003 响应：1 required + 3 recommended 全 fixed
date: 2026-08-14
status: accepted
parent: GOAL-010-r3-s04-scheduled-tasks
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# D-004 · A-003 响应（关门）

## 结论

A-003（grok-build independent，security/data）verdict **conditional**、1 required + 3 recommended —— **全部 fixed**：

| finding | 级别 | 处置 | 证据 |
|---------|------|------|------|
| F-001 标量 n/step 忽略步进（0/5 只匹配第 0 分） | required | **fixed** | store/cron.go：非 * / 非范围的基础值带 /step 一律拒绝（400 INVALID_CRON）；测试补 0/5、5/2、0-59/0 拒绝 |
| F-002 5 年窗口无匹配不记录 | recommended | **fixed** | scheduler.recordUnschedule：status=failed + detail「unschedulable…」，每任务每日至多 1 条（unscheduled 去重表）；测试覆盖 |
| F-003 未知 handler 静默回退 noop | recommended | **fixed** | taskEntity.validateHandler：写入端校验 handler ∈ HandlerKeys()（400 INVALID_HANDLER，新目录化错误码 + i18n）；测试覆盖 |
| F-004 Stop 非幂等（二次 close panic） | recommended | **fixed** | scheduler.onceStop（sync.Once）；重启/多实例双跑为 D-002 §3 文档化残余 |

## 验证

- 修复后回归：go test ./... 全绿、vitest 900/900、e2e mvp/admin 8/8。
