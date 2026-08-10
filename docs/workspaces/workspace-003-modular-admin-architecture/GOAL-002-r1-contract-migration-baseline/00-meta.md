---
id: GOAL-002-r1-contract-migration-baseline
title: R1 · 契约与迁移基线冻结
status: done
parent: GOAL-001-modular-admin-architecture
created: 2026-08-04
updated: 2026-08-04
version: 0.5.0
progress: 4/4
plan_refs:
  - VP-003-modular-admin-architecture
primary_plan: VP-003-modular-admin-architecture
serves_summary: 为 Root 的 R1 方案冻结建立可追溯的模块、迁移、生命周期错误与协议继承基线，不提前承诺 R2 实现或 VP 终态。
---

# GOAL-002 · R1 契约与迁移基线冻结

## 概述

本子目标承接 [Root GOAL-001](../GOAL-001-modular-admin-architecture/00-meta.md) 的 R1 阶段。它把 R1 方案冻结前的四项 Root required 信息门禁收集为可核对的盘点、矩阵与决策证据：模块/注册盘点（Root I-001）、迁移与 seed 所有权（Root I-002）、Fx 与生命周期契约（Root I-003）、协议继承一致性（Root I-007）。

Root `00-meta.md` 的信息表仍是这些门禁的 canonical 状态源；本子目标不复制或改写 Root 的 `open` 状态。

## 愿景对齐

| 字段 | 值 |
|------|----|
| `parent` | `GOAL-001-modular-admin-architecture` |
| `plan_refs` | `VP-003-modular-admin-architecture` |
| `primary_plan` | `VP-003-modular-admin-architecture` |
| Charter | `schema-ui-core-admin-foundation@0.2.0`（经 VP-003 间接对齐） |

本目标只服务 VP-003 R1 的实施边界冻结，与 Root 的模块化单体方向一致；不改变 Charter、VP status 或 Root 成功边界。

## 成功标准

- [x] **C1 / Root I-001**：形成 API、Web、Shell、路由、导航、Schema、权限、中央注册路径、模块候选与跨模块依赖的可追溯清单，并产出 `mvp`/`admin`（及已识别 custom，如有）的 Profile **候选模块集与依赖闭包矩阵**；不冻结 I-004 要求的精确集合或配置覆盖顺序。证据见 `attachments/r1-c1-module-profile-inventory.md`。
- [x] **C2 / Root I-002**：形成 `0001` 起迁移链、seed 所有权、checksum、快照/恢复、回滚、tombstone 与 system-data reconcile 的现状和目标边界记录。证据见 `attachments/r1-c2-migration-seed-boundary.md`。
- [x] **C3 / Root I-003**：记录 Fx/Go 兼容候选、框架无关模块 API、模块核心六项必须能力与按需能力边界、capability 协商及 fail-closed 语义、启动/就绪/停止/失败清理与错误分类的 R1 冻结决策，并区分现状事实与 R2 实施工作；按需能力不得覆盖核心六项，实现仍属于 R2。证据见 `attachments/r1-c3-lifecycle-contract.md`。
- [x] **C4 / Root I-007**：对照 VP-003 继承的 [I-PROTO-001 v0.1.3 覆盖表](../../workspace-001-mvp-admin-foundation/GOAL-001-mvp-admin-foundation/attachments/I-PROTO-001-coverage-draft.md) 完成 R1 模块候选与 `include`/`include-partial`/`exclude` 矩阵，并固定范围扩大所需的新决策与覆盖表升版门槛；只读取该固定协议范围，不读取其所属工作区过程状态。证据见 `attachments/r1-c4-protocol-matrix.md`。

四个检查点等权形成 `progress: 4/4`；`4/4` 只表示 C1-C4 证据收集完成，不自动等于 Root I-* `verified`、R1 方案冻结或阶段放行。本目标已取得独立阶段审计并完成 `/govern` 响应；Root I-001/I-002/I-003/I-007 已由 Root D-004 verified、R1 已推进至 `1/6`。该 Root close-out 仍不代表 R2 实现或 VP 关门。

## 范围与非目标

范围包括现有实现的事实盘点、R1 方案/兼容边界决策、可核对证据路径和阶段审计准备。当前阶段不以文档或测试存在宣称模块化实现完成。

非目标：

- 不冻结 R2 `mvp`/`admin` Profile 精确模块集合或配置覆盖顺序（Root I-004）。
- 不实现 Fx 组合根、模块依赖图、Manifest 聚合 API 或 R3 试点；这些属于后续纲领阶段。
- 不迁移 users、roles、settings、activity 或其他一方能力，不删除旧路径。
- 不扩大 `I-PROTO-001 v0.1.3`，不改变 `D-UPLOAD` 的 `exclude`。

## 阶段计划

1. 收集并核对 C1/C2 的现状清单与证据。
2. 基于清单收敛 C3 的 R1 契约候选和错误语义边界。
3. 对照 C4 协议基线，记录范围一致性和变更门槛。
4. 形成 R1 冻结候选，完成本目标 self 阶段复盘，并在高影响门禁前取得可核对的 independent 审计意见。

R1 的阶段审计模式按 Root 预置建议采用至少 `independent`；Grok Build provider 已通过一次只读设计审计，命令、身份和输出见 [A-001](03-audit/A-001-grok-r1-design-review.md) 及附件。后续 provider 失败或无可核对输出时，不降级冒充 independent。

## 相关 Root 信息门禁

| Root 信息项 | 本目标承接 | Root 当前状态 |
|------------|------------|---------------|
| I-001 | C1 模块/注册/依赖盘点 | open |
| I-002 | C2 迁移/seed/恢复边界 | open |
| I-003 | C3 Fx/API/生命周期错误语义 | open |
| I-007 | C4 协议继承一致性 | open |

## 台账布局

本目标使用平铺的 `01-decision/`、`02-execution/`、`03-audit/` 目录；索引只保留摘要和链接，独立记录按 `D-NNN-*`、`E-NNN-*`、`A-NNN-*` 单调编号。
