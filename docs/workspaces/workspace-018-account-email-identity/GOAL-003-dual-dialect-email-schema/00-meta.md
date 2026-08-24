---
id: GOAL-003-dual-dialect-email-schema
title: R2 双方言 schema + 唯一性
status: active
parent: GOAL-001-account-email-identity
created: 2026-08-24
updated: 2026-08-24
version: 0.1.0
progress: 0/4
plan_refs:
  - VP-018-account-email-identity
primary_plan: VP-018-account-email-identity
serves_summary: 承接 Root R2：迁移 0054 为 users 加 email（可空）与 email_status（pending/verified CHECK）列 + lower(email) 表达式唯一索引（双方言同构、可移植 DDL）；同步五处黄金断言与专项测试；independent 审计后关门。不实施绑定流（R3）、不写仓库层语义。
---

# GOAL-003 · R2 双方言 schema + 唯一性

## 概述

本目标承接 Root `GOAL-001-account-email-identity` 纲领阶段 **R2**：按 R1 冻结的身份合同（GOAL-002 D-001 §1/§2/§3/§6），以模块迁移贡献方式为 `core.auth-session` 增加 **0054 account_email_identity**：

- `users.email TEXT`（可空；NULL = 未绑定，多行 NULL 不冲突）
- `users.email_status TEXT NULL CHECK (IN ('pending','verified'))`（仅 email 非空时有意义）
- `CREATE UNIQUE INDEX idx_users_email_lower ON users(lower(email))`（绑定即占槽的物理承载；SQLite/PostgreSQL 同构表达式索引，NULL 互异语义双方言一致）

对齐递归：GOAL-003 → Root GOAL-001（R2）→ VP-018 → Charter @0.2.0。不进入绑定/校验流、验证码数值（I-005 归 R3 门）、模板或业务域。

## 检查点（progress 来源）

| # | 检查点 | 证据 |
|---|--------|------|
| C1 | 迁移 0054 双方言落地（DDL 可移植；ApplyPostgres nil 先例一致），checksum 台账参与 | 待完成 |
| C2 | 五处黄金断言同步（identity.go head / identity_test lockedHead / migrate_test ×2 / operations_test / restart_test） | 待完成 |
| C3 | 专项测试通过（升级路径 + 全新库 + 大小写唯一冲突 + 多 NULL 共存 + CHECK 拒绝） | 待完成 |
| C4 | independent 审计（grok build · grok-4.6 · high）意见落盘且开放 required = 0 | 待完成 |

`progress` = 已完成检查点 / 4。当前 **0/4**。

## 边界

- 只做 schema 与约束；绑定/校验流、验证码时效数值（I-005）、管理员代填（I-006）归 R3。
- 不改既有迁移 DDL 字符串（checksum 不可变）；不改 `fingerprintR2` 收紧口径（无 ledger 遗留库 fail-closed 保持）。
- 不动 Profile 默认集、不动 Web 面。

## 成功标准

1. 全新库与升级路径均干净应用；`go build ./...` 与 store/authsession 包测试绿。
2. 大小写折叠唯一可核对：`A@x` 与 `a@x` 第二次插入被索引拒绝；多 NULL 共存。
3. independent 审计 required = 0；Root progress 推进至 2/4。
