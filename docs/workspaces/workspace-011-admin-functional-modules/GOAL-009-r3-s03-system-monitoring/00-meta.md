---
id: GOAL-009-r3-s03-system-monitoring
title: R3-S03 · 系统监控与错误日志查看（health/指标/错误日志 UI）
status: active
parent: GOAL-001-admin-functional-modules
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
progress: 0/5
---

# GOAL-009-r3-s03-system-monitoring · 系统监控与错误日志查看（health/指标/错误日志 UI）

## 概述

常用档 S-03（I-011-001 §4）：新建标准 Admin 模块（`admin.system-monitoring` 候选名），提供系统监控与错误日志查看 UI——健康状态（复用 /healthz /readyz）、基础指标（进程/存储/请求面）、错误日志视图（复用 operationlog 事件面 + 内核日志）。运维型后台高频能力，基架无产品面。

## 当前边界

- 只读监控面：health/ready 状态、模块装配摘要、存储信息、错误日志查询（best-effort）
- 权限键与导航（admin 默认集，Profile 内容扩展先例）
- 错误日志面与 operationlog 事件面的切分（S1 冻结）

## 成功标准与路线图（P-001）

- [ ] **S1 · 方案冻结**：边界（指标集/日志切分）、权限键、协议对照、Profile 归属；方案级 self 审视
- [ ] **S2 · 实现**：模块 provider + 端点 + schema 页 + 测试
- [ ] **S3 · 验证**：单元/集成 + 全量回归
- [ ] **S4 · go 影响判定 + 自审**
- [ ] **S5 · 关门**：独立审计（grok，沿用用户偏好）+ required 闭合 + goal-tree 同步

progress: 0/5 由五个等权检查点派生（S1 完成后更新）。

## 审计策略

独立审计沿用 grok build（用户书面偏好，S-01 关门确认）；S-11/S-12 为 security/data 门禁必须独立。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 |
|----|------|-----------------|----------|--------------|-----------------|------|
| I-001 | required | 监控指标集与错误日志切分（与 operationlog/activity 的边界） | S1 方案 | 对照既有 activity 页与 healthz/readyz | open |
| I-002 | required | Profile 归属：进入 admin 默认集？ | S1 方案 | S-02/S-01 先例（内容扩展） | open |

## 父目标

- [GOAL-001-admin-functional-modules](../GOAL-001-admin-functional-modules/00-meta.md)

## 台账布局

本目标从首条记录起使用 `01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；索引与目录条目共同构成正式记录。
