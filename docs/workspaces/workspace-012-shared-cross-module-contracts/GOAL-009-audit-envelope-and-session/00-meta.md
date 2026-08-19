---
id: GOAL-009-audit-envelope-and-session
title: R8 · 审计 envelope 全覆盖与 session 关联
status: done
parent: GOAL-001-shared-cross-module-contracts
created: 2026-08-19
updated: 2026-08-19
version: 0.1.1
progress: 100
plan_refs:
  - VP-012-shared-cross-module-contracts
primary_plan: VP-012-shared-cross-module-contracts
serves_summary: 把 D-003 外 mutation 审计改成结构化 envelope，并把 login session id 关联到 operation_log。
---

# GOAL-009 · R8 · 审计 envelope 全覆盖与 session 关联

## 概述

把仍手写 JSON 的审计写路径改为 `NewDetail`，并给每条新写入关联 session（用户 = refresh token id 经 access JWT `sid`；机器凭据 = credential id）。effective actor 冻结为当前 `actor_id`。

## 成功标准

1. 生产 mutation 写路径不再用手写 JSON 拼 detail；新写入可被 `ParseDetail` 解析（无 detail 的事件保持无 body）。
2. 用户态新审计行在已签发 sid 的请求上带 session_id；service-credential 行带 credential id。
3. 不改 Profile/模块矩阵/协议 pin；不做 impersonation。

## 纲领路线图

| 阶段 | 内容 | 状态 |
|------|------|------|
| S0 | 冻结 session 语义与 writer 范围 | ✅ 已完成（D-001） |
| S1 | 实现 + 测试 | ✅ 已完成（E-001） |
| S2 | 审计关门 | ✅ 已完成（A-001 self pass / A-002 independent pass / A-003 close） |
