---
id: GOAL-003-r2-kernel-composition-root
title: R2 · 内核与组合根基础
status: done
parent: GOAL-001-modular-admin-architecture
created: 2026-08-04
updated: 2026-08-05
version: 0.4.0
progress: 5/5
plan_refs:
  - VP-003-modular-admin-architecture
primary_plan: VP-003-modular-admin-architecture
serves_summary: 在 R1 冻结的框架无关模块语义、Profile 候选边界和迁移/协议约束上，建立可验证的薄内核、组合根、Profile 解析、确定性依赖图和 Manifest/迁移聚合骨架。
---

# GOAL-003 · R2 内核与组合根基础

## 概述

本子目标承接 Root R2 阶段。它实现 R1 已冻结的边界，不重新打开 R1 的协议范围或 Profile 候选盘点；精确 `mvp`/`admin` Profile 集与配置覆盖顺序在本目标的 C1/I-004 中冻结。Root I-004、I-005 仍由 Root canonical 台账维护，子目标只承接证据。

## 愿景对齐

| 字段 | 值 |
|------|----|
| `parent` | `GOAL-001-modular-admin-architecture` |
| `plan_refs` | `VP-003-modular-admin-architecture` |
| `primary_plan` | `VP-003-modular-admin-architecture` |
| Charter | `schema-ui-core-admin-foundation@0.2.0`（经 VP-003 间接对齐） |

## 成功标准

- [x] **C1 / Root I-004**：冻结 `mvp`、`admin`（以及需要时 custom fork）精确模块集合、依赖闭包和配置覆盖顺序，并用正反例证明 unknown/conflict/missing dependency fail closed。
- [x] **C2 / R2 platform**：实现薄内核与框架无关模块公共 API；Fx 只在组合根使用，模块 API 不暴露 Fx 类型；Go/Fx 兼容范围与稳定错误分类可核对。
- [x] **C3 / R2 graph**：实现确定性模块注册/依赖图/贡献冲突校验、capability negotiation 和启动前 fail-closed 行为，覆盖拓扑启动、反向停止与失败清理语义。
- [x] **C4 / Root I-005**：形成全局迁移收集与 Manifest/Schema/Navigation/permission/config 聚合骨架，并从 `/.well-known/schema-ui/app-manifest.json` 提供确定性、登录前可读的端点；静态生产 Manifest 不得继续作为静默兜底。
- [x] **C5 / verification**：完成 API/Web focused tests、双 Profile contract tests、startup/readiness/stop failure tests、migration/manifest aggregation fixtures 和独立 R2 gate audit；required findings 已合法闭合，Root R2 stage response 仍由 Root 台账单独承接。

五个检查点等权形成 `progress: 5/5`；`5/5` 只表示本子目标检查点完成，不等于 Root R2 放行、VP exit #2/#4 已完成或 Root done。I-004/I-005 已由 Root D-006/E-006/A-005 与 child A-005 response verified；I-006 仍开放。

## 阶段计划

1. 先收敛 C1 的精确 Profile 配置与 I-004 反例；依赖闭包必须由同一解析器产出。
2. 设计并实现 C2/C3 的薄内核、框架无关模块契约、组合根和确定性图校验。
3. 在 C2/C3 通过 focused tests 后实现 C4 的 migration collector 与 Manifest aggregation/代理骨架。
4. 运行 C5 verification，执行 self review，并在 R2 gate 前用 Grok Build 做 independent audit；根据 required findings 回流修正。

## 范围与非目标

范围：Go API 内核/组合根、配置/Profile 解析、模块注册与依赖图、聚合契约、迁移收集接口、Manifest 登录前代理和对应测试。

非目标：不迁移 users/roles/settings/activity 实际业务模块，不完成 R3 试点，不关闭 I-006，不删除旧中央注册/静态 Manifest（除非 C4 只删除其生产静默兜底所需的最小路径并有明确证据），不扩大 I-PROTO-001 v0.1.3。

## 依赖与门禁

| Root 信息项 | 本目标承接 | 当前状态 | 最晚需要 |
|-------------|------------|----------|----------|
| I-004 | C1 exact Profile set/precedence | verified | 2026-08-05 Root D-006/E-006/A-005 and child A-004 response |
| I-005 | C4 aggregation input/order/conflict/cache/auth boundary | verified | 2026-08-05 Root D-006/E-006/A-005 and child A-004 response |
| I-006 | 仅记录不提前关闭 | open | R3 方案冻结前 |

R1 child `GOAL-002` 已 `done`；本目标不能复用其 `progress` 或 audit status。R2 implementation risk is `independent`-eligible because it touches production startup, migration and compatibility boundaries; provider failure remains fail-closed.
