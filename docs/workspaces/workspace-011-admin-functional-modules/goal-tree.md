---
title: 目标树 · workspace-011-admin-functional-modules
status: active
created: 2026-08-14
updated: 2026-08-14
parent: null
version: 0.2.5
workspace_id: workspace-011-admin-functional-modules
---

# 目标树 · 标准 Admin 功能模块（通用 + 常用业务领域）

> 工作区：`workspace-011-admin-functional-modules`
> canonical：`docs/workspaces/workspace-011-admin-functional-modules/`
> Root：`GOAL-001-admin-functional-modules`（**交付目标 · active**）
> primary_plan：`VP-011-admin-functional-modules`（**active**）

## 树

```text
GOAL-001-admin-functional-modules [active]  · 标准 Admin 功能模块（分档交付）
├── GOAL-002-r1-bounded-research [done]       · R1 有界调研：候选池收集 + 三档分档（5/5）
├── GOAL-003-r2-f01-dashboard [active]        · R2-F01 仪表盘/控制台（生产 home）
├── GOAL-004-r2-f02-data-import-export [active] · R2-F02 数据导入/导出（共享能力）（4/5 · S4 就绪，关门待 independent）
├── GOAL-005-r2-f03-account-center [done]      · R2-F03 个人中心与账户安全 + 账号启停（5/5）
└── GOAL-006-r2-f04-notification-center [active] · R2-F04 通知中心（站内通知）
```

Root 于 2026-08-14 开区（VP-011 激活 + freshness review PASS，候选 `f14ab9d`）。首阶段 R1 = 有界调研；分档产出后按 Root 路线图逐波立项（R2 一等公民 / R3 常用 / R4 增补 backlog）。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-admin-functional-modules | 标准 Admin 功能模块（通用 + 常用业务领域 · 分档交付） | null | active | —（纲领路线图就位） | 2026-08-14 |
| GOAL-002-r1-bounded-research | R1 · 有界调研：候选池收集与三档分档 | GOAL-001-admin-functional-modules | done | 5/5 | 2026-08-14 |
| GOAL-003-r2-f01-dashboard | R2-F01 · 仪表盘/控制台（生产 Profile home） | GOAL-001-admin-functional-modules | active | 0/5 | 2026-08-14 |
| GOAL-004-r2-f02-data-import-export | R2-F02 · 数据导入/导出（schema 驱动 · 共享能力） | GOAL-001-admin-functional-modules | active | 4/5 | 2026-08-14 |
| GOAL-005-r2-f03-account-center | R2-F03 · 个人中心与账户安全 + 账号启停 | GOAL-001-admin-functional-modules | done | 5/5 | 2026-08-14 |
| GOAL-006-r2-f04-notification-center | R2-F04 · 通知中心（站内通知） | GOAL-001-admin-functional-modules | active | 0/5 | 2026-08-14 |

## 维护说明

- 层级唯一来源是目标 `00-meta.md` 的 `parent`。
- 阶段子目标按 Root 纲领路线图立项；progress 只写在子目标。