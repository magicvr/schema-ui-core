---
id: GOAL-010-r4-c4-schema-other-migration
title: R4-C4 · Schema 与其他能力迁移
status: done
parent: GOAL-005-r4-full-module-migration
created: 2026-08-05
updated: 2026-08-05
version: 0.1.0
progress: 4/4
plan_refs:
  - VP-003-modular-admin-architecture
primary_plan: VP-003-modular-admin-architecture
serves_summary: 承接 GOAL-005 C4：把仍存在的 Schema-driven Admin 能力（settings/activity 及横切 operationlog 读面）按 C1 冻结范围迁入 Provider 契约，Schema owner map 转贡献驱动，Records 保持 historical-only，并承接 C3 遗留门禁（Manifest secrecy、Ready 失败清理、PolicyID/Visibility/JSON 校验器、ledger drift/unknown）。
---

# GOAL-010 · R4-C4 Schema 与其他能力迁移

## 概述

本子目标是 `GOAL-005-r4-full-module-migration` 的 C4 子目标，按 C1 冻结范围把
settings/activity（及 operationlog 只读 activity 面）从中心 Register/Schema
owner map/Manifest adminModules 迁入 GOAL-008/009 的 Provider 契约。Records 保持
historical-only，不恢复任何产品 CRUD。承接 GOAL-008 E-004 登记给 C3/C5 的门禁中
与 C4 相关的项（Manifest secrecy、Ready 失败清理、PolicyID/Visibility/JSON 校验器、
ledger drift/unknown）。

## 愿景对齐

| 字段 | 值 |
|------|----|
| `parent` | `GOAL-005-r4-full-module-migration` |
| `plan_refs` | `VP-003-modular-admin-architecture` |
| `primary_plan` | `VP-003-modular-admin-architecture` |
| Charter | `schema-ui-core-admin-foundation@0.2.0`（经 VP-003 间接对齐） |
| 审计模式 | `independent`；迁移切片使用 Grok Build `grok-4.5` / `high` |

## 成功标准

- [x] **C4.1 / settings 迁移**：`admin.settings` provider 化（HTTP/Schema/Auth/Nav/
  Manifest 经 Registrar 贡献）；Schema 内容与 Manifest fragment 模块所有；中心
  settings 分支与 adminModules 特例清除。
- [x] **C4.2 / activity 迁移**：`admin.activity` provider 化（只读 operations 面）；
  operationlog writer 保持 core.operationlog 职责，Activity 只是查询/UI；Activity
  disabled 时 writer 仍工作。
- [x] **C4.3 / Schema owner map 转贡献驱动**：`schemaDocumentsForPlan` 的 owner map
  改为由 provider/schema 贡献驱动（解决 F-IND-C33-001 residual）；settings/activity
  与 users/roles 一致。
- [x] **C4.4 / C3 遗留门禁 + 关门**：Manifest secrecy 扫描、Ready 失败反向清理、
  PolicyID/Visibility/JSON 校验器、Records historical-only 保持（按 D-002 收窄；
  ledger drift/unknown 运行时 fail-closed 移交 C5 数据门禁）；self + Grok
  independent 无开放 required finding。

四个检查点等权；`progress: 4/4`（C4.1-C4.4 完成）。ledger drift/unknown 运行时
fail-closed 属 store/migration 路径改造，登记 C5 residual。完成本子目标只表示 C4
关闭，不关闭 GOAL-005、Root 或 VP-003，不自动放行 C5。

## 信息门禁

| 编号 | 级别 | 必须回答的问题 | 影响 | 最晚阶段 | 收集动作 | 状态 | 证据 |
|------|------|----------------|------|----------|----------|------|------|
| C4-I001 | required | settings/activity 当前中心注册、schema/manifest 投影与 operationlog 读面状态？ | C4.1/C4.2 | C4.1 | 全仓扫描（沿用 GOAL-009 E-002 模式） | verified | E-002：中心状态 + 行为矩阵 |
| C4-I002 | required | Schema owner map 转贡献驱动的语义（provider page 贡献 vs plan 门禁）是否固定？ | C4.3 | C4.3 | 设计 + 测试 | verified | schemaDocumentsForPlan 模块驱动 owner |
| C4-I003 | required | Manifest secrecy、Ready 失败清理、PolicyID/Visibility/JSON 校验器、ledger drift/unknown 的实现边界？ | C4.4 | C4.4 | GOAL-008 E-004 登记 + 实施 | verified | E-002；ledger drift 登记 C5 residual |
| C4-I004 | non-blocking | Records historical-only 是否保持（不恢复 CRUD）？ | C4.4 验收 | C4.4 | 负向断言 | open | D-003；GOAL-007 |

## 阶段路线图

1. settings/activity 迁移扫描 + 行为边界（C4.1 起点）。
2. settings/activity provider 化 + 中心特例清除（C4.1/C4.2）。
3. Schema owner map 转贡献驱动（C4.3）。
4. C3 遗留门禁实施 + 全量回归 + self + Grok independent 复审 + 关门（C4.4）。

## 范围与非目标

范围包括 settings/activity provider 化、Schema owner map 贡献驱动、Manifest
secrecy、Ready 失败清理、PolicyID/Visibility/JSON 校验器、ledger drift/unknown。
非目标包括 Records 恢复、R5 Profile 运维收敛、R6 旧路径终态删除。`0003`/`0006`
迁移账本与历史 operation-log 保留。
