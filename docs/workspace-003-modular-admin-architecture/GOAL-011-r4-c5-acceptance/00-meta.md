---
id: GOAL-011-r4-c5-acceptance
title: R4-C5 · 验收与关门
status: active
parent: GOAL-005-r4-full-module-migration
created: 2026-08-05
updated: 2026-08-05
version: 0.1.0
progress: 3/4
plan_refs:
  - VP-003-modular-admin-architecture
primary_plan: VP-003-modular-admin-architecture
serves_summary: 承接 GOAL-005 C5：R4 验收与关门——同一 Web 构建双 Profile、API/Schema/Manifest/授权/持久化与行为矩阵通过，闭合 C5 门禁（ledger drift/unknown、双 Profile 失败矩阵、校验器深化、中心适配器终态删除、Schema owner 完全贡献驱动），self + Grok independent 无开放 required finding，形成进入 R5 的结论。
---

# GOAL-011 · R4-C5 验收与关门

## 概述

本子目标是 `GOAL-005-r4-full-module-migration` 的 C5 检查点：R4 全量迁移后的验收
与关门。它验证同一 Web 构建在 `mvp`/`admin` 双 Profile 下工作，API/Schema/Manifest/
授权/持久化与行为矩阵通过，闭合 GOAL-010 E-003 登记的 C5 门禁，并经 self + Grok
independent 审计后形成进入 R5 的结论。不关闭 Root/VP-003/R5/R6。

## 愿景对齐

| 字段 | 值 |
|------|----|
| `parent` | `GOAL-005-r4-full-module-migration` |
| `plan_refs` | `VP-003-modular-admin-architecture` |
| `primary_plan` | `VP-003-modular-admin-architecture` |
| Charter | `schema-ui-core-admin-foundation@0.2.0`（经 VP-003 间接对齐） |
| 审计模式 | `cross`；关门使用 Grok Build `grok-4.5` / `high` |

## 成功标准

- [x] **C5.1 / 双 Profile 行为矩阵**：同一 Web 构建在 `mvp`/`admin` 双 Profile 下
  页面集/路由/Schema/Manifest/授权一致；mvp 禁用 settings/activity 时面消失但
  operationlog writer 仍工作。
- [x] **C5.2 / C5 数据门禁**：ledger drift/unknown 运行时 fail-closed
  （fresh/upgrade/reconcile 深度验证）；双 Profile register/conflict/Start/Ready
  失败清理矩阵。
- [x] **C5.3 / C5 收尾**：PolicyID/Visibility allowlist 深化、中心
  RegisterSettings/RegisterActivity 终态删除、Schema owner 完全 ContributionSet
  驱动、readyz 真实 readiness（或 residual）。
- [ ] **C5.4 / 关门审计**：self + Grok independent 无开放 required finding；
  形成进入 R5 的结论，向 GOAL-005 C5 close-out 提交。

四个检查点等权；当前 `progress: 3/4`（C5.1-C5.3 验证完成）。C5.3 收尾项以
accepted-residual 或文档化登记（详见 E-002）。完成本子目标表示 R4 关闭，不关闭
GOAL-005（done 需父级确认）、Root、VP-003、R5 或 R6。

## 信息门禁

| 编号 | 级别 | 必须回答的问题 | 影响 | 最晚阶段 | 收集动作 | 状态 | 证据 |
|------|------|----------------|------|----------|----------|------|------|
| C5-I001 | required | 同一 Web 构建在 mvp/admin 双 Profile 下行为矩阵是否通过？ | C5.1 | C5.1 | e2e + 集成矩阵 | verified | composition 双 Profile 测试 + Web integration |
| C5-I002 | required | ledger drift/unknown 与双 Profile 失败矩阵是否 fail-closed？ | C5.2 | C5.2 | store/migration + composition 矩阵 | verified | store migrate 测试 + TestDualProfileRegisterValidationFailClosed |
| C5-I003 | required | C5 收尾项（校验器深化/中心适配器删除/owner 贡献驱动/readyz）是否闭合？ | C5.3 | C5.3 | 实施 + 测试 | verified | E-002；residual 登记 |
| C5-I004 | required | R4 验收结论是否可形成（无开放 required，进入 R5 依据）？ | C5.4 | C5.4 | self + Grok | collecting | 待 C5.4 |

## 阶段路线图

1. 双 Profile 行为矩阵 + e2e（C5.1）。
2. ledger drift/unknown + 双 Profile 失败矩阵（C5.2）。
3. C5 收尾项（C5.3）。
4. self + Grok independent 关门审计 + R5 结论（C5.4）。

## 范围与非目标

范围包括双 Profile 验收、ledger/失败矩阵、C5 收尾、关门审计与 R5 结论。非目标包括
R5 运维/配置收敛、R6 终态删除、Root/VP-003 关门。Records historical-only 保持。
