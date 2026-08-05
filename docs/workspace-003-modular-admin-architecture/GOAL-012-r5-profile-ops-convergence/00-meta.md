---
id: GOAL-012-r5-profile-ops-convergence
title: R5 · Profile 运维与数据收敛
status: active
parent: GOAL-001-modular-admin-architecture
created: 2026-08-05
updated: 2026-08-05
version: 0.1.0
progress: 3/4
plan_refs:
  - VP-003-modular-admin-architecture
primary_plan: VP-003-modular-admin-architecture
serves_summary: 承接 Root R5：Profile 运维/配置收敛（custom/MODULES_ENABLED 覆盖顺序、Configuration 运行时迁移）、fresh/reconcile 深度验证、readyz 真实模块图 readiness、双 Profile Start/Ready 失败矩阵、升级恢复、Docker/代理与 fork 文档，并承接 R4 residual 清单。
---

# GOAL-012 · R5 Profile 运维与数据收敛

## 概述

本子目标是 Root `GOAL-001-modular-admin-architecture` 的 R5 阶段：完成 Profile
运维/配置收敛与文档、fresh/reconcile、readyz/诊断、代理/容器、升级恢复与 fork
文档。承接 R4 residual 清单（Schema 完全贡献驱动、中心 RegisterSettings/Register
Activity 适配器终态删除、PolicyID/Visibility 深化、readyz 真实、双 Profile
Start/Ready 矩阵、Configuration 运行时迁移）。R5 **不**否定 R2 已冻结的精确 Profile
集，除非新决策书面改写。

## 愿景对齐

| 字段 | 值 |
|------|----|
| `parent` | `GOAL-001-modular-admin-architecture` |
| `plan_refs` | `VP-003-modular-admin-architecture` |
| `primary_plan` | `VP-003-modular-admin-architecture` |
| Charter | `schema-ui-core-admin-foundation@0.2.0`（经 VP-003 间接对齐） |
| 审计模式 | `self`；升级/恢复/容器放行倾向 `independent` |

## 成功标准

- [ ] **C5.1 / Profile 运维与配置收敛**：custom Profile + `MODULES_ENABLED` 覆盖
  顺序、Configuration 运行时迁移边界；R4 residual（Schema 贡献驱动、中心适配器
  终态删除、PolicyID/Visibility 深化、双 Profile Start/Ready 矩阵）闭合或 residual。
- [x] **C5.2 / 数据生命周期**：fresh/bootstrap 与 versioned reconcile 分离、升级/
  恢复演练、checksum/ledger 深度 fail-closed；R4-I004 residual 复核。
- [x] **C5.3 / 可运维性**：readyz 真实模块图 readiness（迁移/reconcile/依赖/lifecycle）、
  /healthz 与 /readyz 分离、诊断语义。
- [x] **C5.4 / 可 fork 与文档**：Docker/生产代理、快速启动、升级恢复与 fork 文档
  反映新架构；CI/本地矩阵覆盖。

四个检查点等权；当前 `progress: 3/4`（C5.2/C5.3/C5.4 完成或 partial；C5.1 部分）。
versioned reconcile 与完整 fork/升级恢复文档登记 residual（E-003）。R4 residual 状态
见 E-002。完成本子目标表示 R5 关闭，不关闭 Root、VP-003、R6。

## 信息门禁

| 编号 | 级别 | 必须回答的问题 | 影响 | 最晚阶段 | 收集动作 | 状态 | 证据 |
|------|------|----------------|------|----------|----------|------|------|
| R5-I001 | required | Profile 运维/配置收敛与 R4 residual 的闭合边界？ | C5.1 | C5.1 | R4 residual 清单 + 设计 | collecting | GOAL-011 E-003 |
| R5-I002 | required | fresh/reconcile/升级恢复的深度验证边界？ | C5.2 | C5.2 | store 迁移 + 演练 | verified | E-003：fresh/upgrade/recovery 测试覆盖；versioned reconcile 登记 residual |
| R5-I003 | required | readyz 真实 readiness 的实现边界？ | C5.3 | C5.3 | 设计 + 实施 | verified | E-002：readinessGate + RegisterWithReadiness + 两态测试 |
| R5-I004 | non-blocking | hosted E2E/容器环境作为补充证据？ | C5.4 | C5.4 | 记录环境 | open | R4-I005 |

## 阶段路线图

1. Profile 运维/配置收敛 + R4 residual 闭合（C5.1）。
2. fresh/reconcile/升级恢复深度验证（C5.2）。
3. readyz 真实 readiness + 诊断（C5.3）。
4. Docker/代理/fork 文档 + 回归 + self + Grok（C5.4）。

## 范围与非目标

范围包括 Profile 配置收敛、数据生命周期、可运维性、可 fork 文档。非目标包括 R6
旧路径终态删除、Root/VP-003 关门、Records 恢复。R2 精确 Profile 集不改写。
