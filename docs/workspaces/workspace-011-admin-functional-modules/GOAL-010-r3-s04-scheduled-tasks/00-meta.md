---
id: GOAL-010-r3-s04-scheduled-tasks
title: R3-S04 · 定时任务管理（cron 后台）
status: active
parent: GOAL-001-admin-functional-modules
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
progress: 0/5
---

# GOAL-010-r3-s04-scheduled-tasks · 定时任务管理（cron 后台）

## 概述

常用档 S-04（I-011-001 §4）：新建标准 Admin 模块（`admin.scheduled-tasks` 候选名），提供定时任务管理——任务定义（cron 表达式/启用/说明）、运行历史、手动触发、调度器执行循环。平台型后台高频能力，基架无。

## 当前边界

- 任务定义 CRUD + cron 校验 + 启用/停用 + 手动触发 + 最近运行历史
- 进程内调度器（best-effort 单实例语义文档化）+ 新表（任务/运行记录）
- 权限键 + 审计（任务事件）+ admin 默认集

## 成功标准与路线图（P-001）

- [ ] **S1 · 方案冻结**：任务模型（cron 校验/运行语义）、调度器边界（单实例/持久化）、权限/审计、Profile 归属；方案级 self 审视
- [ ] **S2 · 实现**：迁移 + 任务仓库 + 调度器 + 端点 + schema 页 + 测试
- [ ] **S3 · 验证**：单元/集成 + 调度实测 + 全量回归
- [ ] **S4 · go 影响判定 + 自审**
- [ ] **S5 · 关门**：独立审计（grok）+ required 闭合 + goal-tree 同步

progress: 0/5 由五个等权检查点派生（S1 完成后更新）。

## 审计策略

独立审计沿用 grok build（用户书面偏好，S-01 关门确认）；S-11/S-12 为 security/data 门禁必须独立。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 |
|----|------|-----------------|----------|--------------|-----------------|------|
| I-001 | required | cron 表达式校验器（依赖/自研）与运行语义（重叠/补跑） | S1 方案 | 评估 robfig/cron 依赖 or 自研校验 | open |
| I-002 | required | 调度器边界：进程内 vs 持久化任务队列；重启语义 | S1 方案 | 单实例假设文档化 | open |
| I-003 | required | Profile 归属：admin 默认集 | S1 方案 | 先例一致 | open |

## 父目标

- [GOAL-001-admin-functional-modules](../GOAL-001-admin-functional-modules/00-meta.md)

## 台账布局

本目标从首条记录起使用 `01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；索引与目录条目共同构成正式记录。
