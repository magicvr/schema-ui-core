---
title: 目标树 · workspace-011-admin-functional-modules
status: active
created: 2026-08-14
updated: 2026-08-14
parent: null
version: 0.3.2
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
├── GOAL-003-r2-f01-dashboard [done]        · R2-F01 仪表盘/控制台（生产 home）（5/5 · 990daa8）
├── GOAL-004-r2-f02-data-import-export [done]    · R2-F02 数据导入/导出（共享能力）（5/5）
├── GOAL-005-r2-f03-account-center [done]      · R2-F03 个人中心与账户安全 + 账号启停（5/5）
├── GOAL-006-r2-f04-notification-center [done] · R2-F04 通知中心（站内通知）（5/5）
├── GOAL-007-r3-s02-file-library [done]    · R3-S02 文件/附件库（统一文件管理）（5/5）
├── GOAL-008-r3-s01-data-dictionary [done]  · R3-S01 数据字典（枚举/字典管理）（5/5）
├── GOAL-009-r3-s03-system-monitoring [done] · R3-S03 系统监控与错误日志（5/5）
├── GOAL-010-r3-s04-scheduled-tasks [done]  · R3-S04 定时任务管理（5/5）
├── GOAL-011-r3-s11-login-captcha [done]       · R3-S11 登录验证码（5/5）
├── GOAL-012-r3-s12-recycle-bin [done]      · R3-S12 回收站/软删除（5/5）
├── GOAL-014-form-experience [done]      · R4 表单体验：字段级校验/错误展示 + 弹窗布局（5/5）
└── GOAL-015-dict-inner-page-breadcrumb [active] · R4 数据字典内页 + 面包屑层级导航（0/5）
```

Root 于 2026-08-14 开区（VP-011 激活 + freshness review PASS，候选 `f14ab9d`）。首阶段 R1 = 有界调研；分档产出后按 Root 路线图逐波立项（R2 一等公民 / R3 常用 / R4 增补 backlog）。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-admin-functional-modules | 标准 Admin 功能模块（通用 + 常用业务领域 · 分档交付） | null | active | —（纲领路线图就位） | 2026-08-14 |
| GOAL-002-r1-bounded-research | R1 · 有界调研：候选池收集与三档分档 | GOAL-001-admin-functional-modules | done | 5/5 | 2026-08-14 |
| GOAL-003-r2-f01-dashboard | R2-F01 · 仪表盘/控制台（生产 Profile home） | GOAL-001-admin-functional-modules | done | 5/5 | 2026-08-14 |
| GOAL-004-r2-f02-data-import-export | R2-F02 · 数据导入/导出（schema 驱动 · 共享能力） | GOAL-001-admin-functional-modules | done | 5/5 | 2026-08-14 |
| GOAL-005-r2-f03-account-center | R2-F03 · 个人中心与账户安全 + 账号启停 | GOAL-001-admin-functional-modules | done | 5/5 | 2026-08-14 |
| GOAL-006-r2-f04-notification-center | R2-F04 · 通知中心（站内通知） | GOAL-001-admin-functional-modules | done | 5/5 | 2026-08-14 |
| GOAL-007-r3-s02-file-library | R3-S02 · 文件/附件库（统一文件管理、引用、清理） | GOAL-001-admin-functional-modules | done | 5/5 | 2026-08-14 |
| GOAL-008-r3-s01-data-dictionary | R3-S01 · 数据字典（枚举/字典管理） | GOAL-001-admin-functional-modules | done | 5/5 | 2026-08-14 |
| GOAL-009-r3-s03-system-monitoring | R3-S03 · 系统监控与错误日志查看 | GOAL-001-admin-functional-modules | done | 5/5 | 2026-08-14 |
| GOAL-010-r3-s04-scheduled-tasks | R3-S04 · 定时任务管理（cron 后台） | GOAL-001-admin-functional-modules | done | 5/5 | 2026-08-14 |
| GOAL-011-r3-s11-login-captcha | R3-S11 · 登录验证码 | GOAL-001-admin-functional-modules | done | 5/5 | 2026-08-14 |
| GOAL-012-r3-s12-recycle-bin | R3-S12 · 回收站/软删除管理 | GOAL-001-admin-functional-modules | done | 5/5 | 2026-08-14 |
| GOAL-013-nav-order-config | 导航顺序：默认清单 + 配置文件覆盖（方案 A） | GOAL-001-admin-functional-modules | done | 5/5 | 2026-08-14 |

## 维护说明

- 层级唯一来源是目标 `00-meta.md` 的 `parent`。
- 阶段子目标按 Root 纲领路线图立项；progress 只写在子目标。