---
id: E-002
goal: GOAL-005-r2-f03-account-center
title: S1 · 方案冻结执行（I-001/I-002/I-003 关闭 + 必办-5 核对）
date: 2026-08-14
status: recorded
parent: GOAL-005-r2-f03-account-center
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# E-002 · S1 · 方案冻结

## 事实

- 产出 [D-002-s1-freeze.md](../01-decision/D-002-s1-freeze.md)（方案冻结：模块 `admin.account`、会话模型、启停语义、权限键、Profile 扩展声明、前端设计、未选方案）。
- 基架核对证据（2026-08-14，均来自 HEAD `77722d0`）：
  - 会话基建：`refresh_tokens` 表 + `token_version`（apps/api/internal/modules/authsession，W4 P0-3）——会话模型直接复用，无需新表。
  - 迁移账本：既有版本 1–12（authsession 1/2/9/11/12、corepersistence 3/6、operationlog 4/5/8、settings 7/10）→ **0013 空闲**，分配给 F-03。
  - 权限贡献机制：`kernel.PermissionContribution` + `systemdata.Reconcile`（PolicyAdmin 等）；模块注册范式：users/activity provider。
  - 前端 schema 能力：`form.recordSource` prefill（ADR-0021）、`disabledWhen {field,equals}` 行状态门控、`permissions.edit` 表达式门控（$context.user.permissions contains …）——account 页与 users 页扩展均可纯 schema 表达。
  - 事件表：operationlog 常量为追加式（无白名单），新增 5 个事件无迁移影响。
- **必办核对（I-011-001 `8）**：必办-5 ✅（D-002 `1/`3/`4）；必办-1/2/3/4 不适用（已在 D-002 `7 留痕）。

## 信息项关闭

| ID | 级别 | 结论 | 证据 |
|----|------|------|------|
| I-001 | required | 会话模型冻结：refresh_tokens 为会话表；单条粒度吊销；改密/停用走 token_version 联动；单条吊销为短窗残余（已文档化） | D-002 `2 |
| I-002 | required | 启停字段 `enabled`（默认 1）+ 迁移 0013 + 权限键 `users.enable/disable`（PolicyAdmin） | D-002 `1/`3 |
| I-003 | non-blocking | 停用与锁定正交：disable 不重置失败计数/锁定窗口；enable 不清锁；unlock 才清 | D-002 `3/`8 |

## 进度评估

S1 检查点完成：方案冻结 + 方案级 self 审视（A-001）就绪。**S1 done，进入 S2 实现**（D-002 `9 清单）。