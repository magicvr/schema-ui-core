---
id: GOAL-005-r2-f03-account-center
title: R2-F03 · 个人中心与账户安全 + 账号启停
status: active
parent: GOAL-001-admin-functional-modules
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
progress: 0/5
---

# GOAL 005-r2-f03-account-center · 个人中心与账户安全 + 账号启停

## 概述

一等公民 F-03（I-011-001 `3）：**自助个人中心**（个人资料、修改密码、会话列表/吊销）与**管理员账号启停**（启用/停用/手动解锁）。复用改密/吊销基建（C-10）与失败自动锁定（C-11）；新增会话列表/吊销端点与启停端点（权限键 users.enable/disable + 状态字段迁移）。

## 当前边界

- 自助：资料查看/修改、改密（吊销 access token）、会话列表/吊销
- 启停：enabled 状态字段 + 迁移、启用/停用/手动解锁端点 + 权限键
- 安全相关（authz/session）→ 关门须 independent 审计
- 不改变登录/锁定协议语义（C-11 保持）；共享基架问题回流 VP-009

## 成功标准与路线图（P-001）

- [ ] **S1 · 方案冻结**：方案冻结：会话模型（列表/吊销粒度、token_version 联动）+ 启停状态字段与迁移设计 + 权限键；方案级 self 审视
- [ ] **S2 · 实现**：实现：API（会话列表/吊销、启停、自助改密/资料）+ 页面（个人中心、用户启停操作）+ 迁移 + 测试
- [ ] **S3 · 验证**：验证：单元/集成（会话吊销后旧 token 失效、停用后登录拒绝）+ 全量回归
- [ ] **S4 · go 影响判定 + 自审**：go 影响判定 + self 审计
- [ ] **S5 · 关门**：关门：required 全闭合 + **independent 审计**（grok build） + goal-tree 同步

progress: 0/5 由五个等权检查点派生（S1 完成后更新）。

## 审计策略

self + **independent 关门**（security/authz 高影响门禁；provider 沿用 grok build，P-004 在关门时确认）。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 |
|----|------|-----------------|----------|--------------|-----------------|------|
| I-001 | required | 会话模型：列表/吊销粒度、与 token_version 联动语义 | S1 方案 | 对照 auth 模块 session/token 机制 | open |
| I-002 | required | 启停状态字段与迁移（enabled + 迁移版本）及权限键（users.enable/disable） | S1 方案 | 对照迁移账本 0001～0010 约定 | open |
| I-003 | non-blocking | 停用与锁定（C-11）的交互语义（停用是否重置失败计数等） | S1 方案 | 方案冻结时定 | open |

## 父目标

- [GOAL-001-admin-functional-modules](../GOAL-001-admin-functional-modules/00-meta.md)

## 台账布局

本目标从首条记录起使用 `01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；索引与目录条目共同构成正式记录。
