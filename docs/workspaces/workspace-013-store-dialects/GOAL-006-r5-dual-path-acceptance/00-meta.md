---
id: GOAL-006-r5-dual-path-acceptance
title: R5 · 双路径验收（升级策略 + 备份合同 + 全链证据 → Root 关门）
status: active
parent: GOAL-001-store-dialects
created: 2026-08-20
updated: 2026-08-20
version: 0.1.0
progress: 4/5
plan_refs:
  - VP-013-store-dialects
primary_plan: VP-013-store-dialects
serves_summary: 交付 R5：SQLite 缺省路径回归 + PostgreSQL 生产向验收（迁移、共事务、备份/恢复合同），收敛 Root I-001（SQLite→PG 升级策略）与 I-004（PG 备份合同），跑通全链双路径证据，达成 VP-013 退出判据并关闭工作区根目标。
---

# GOAL-006 · R5 · 双路径验收

## 概述

Root 纲领 **R5**（依赖 R3/R4）：把 VP-013 的「SQLite 内嵌缺省 + PostgreSQL 生产权威」在**两条完整路径**上验收并给出台账证据，收敛 Root 剩余 required 信息项，最后关闭工作区根目标。

## 非目标

- 不新增业务域 / Admin 页面；不引 ORM；不改 Profile/Compose 缺省（sqlite 仍是 dev/mvp/快测默认）；不重做 R1–R4。

## 纲领路线图（P-001）

| 阶段 | 内容 | 状态 |
|------|------|------|
| U0 | 双路径回归基线：sqlite 全量 `go test ./...` + live PG 全量 boot + 完整启动（既有证据固化） | ✅（sqlite 0 FAIL；`TestFullCatalogPostgresBootstrapIntegration` + `TestCompositionPostgresStartup` 全绿） |
| U1 | **I-001 升级策略**（SQLite→PG）：结论/残余（in-place vs dump/restore vs fresh bootstrap），书面落盘 + 抽样证据 | ✅ D-002 / E-002（fresh bootstrap + 逻辑迁移；in-place 跨引擎不可行书面 residual） |
| U2 | **I-004 备份/恢复合同**：替代 `VACUUM INTO` 的 PG 生产路径（pg_dump 或等价），落盘 + 可执行验证 | ✅ D-002 / E-002（pg_dump→pg_restore round-trip 实跑通过） |
| U3 | 跨模块共事务验收（job/审计/钱包等共用事务在 PG 上的一致性与回滚）+ `readyz` 生产向就绪 | ✅ E-002（`TestPostgresCrossModuleSharedTx` live PG：commit+rollback 双双验证） |
| U4 | 关门：VP-013 退出判据 1–6 全链核对 + self + independent（production 门禁）→ GOAL-006 done → **Root 5/5 关门** | 🔄 待 independent 后关门 |

## 成功标准

1. sqlite 缺省路径全量回归 0 FAIL；PG 双路径（bootstrap + 完整启动）live 全绿。
2. I-001：SQLite→PG 升级策略有书面结论与证据（in-place 不可行则 dump/restore/rebootstrap + 范围）。
3. I-004：PG 备份/恢复合同明确（替换 `VACUUM INTO`），有可执行验证。
4. 跨模块共事务在 PG 上验收通过（一致/回滚语义）；`readyz` 生产向就绪。
5. VP-013 退出判据 1–6 全链可核对；无 open required；Root 5/5 关门。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| RT-I-001 | required | 存量 SQLite→PG：in-place vs dump/restore vs fresh bootstrap | U1 / Exit 2 | U1 | 原型/抽样 + 决策 | **verified** | 2026-08-20（D-002） | fresh bootstrap + 逻辑迁移；in-place 跨引擎不可行（书面 residual 范围） |
| RT-I-004 | required | PG 备份/恢复合同（替代 `VACUUM INTO`） | U2 / Exit 4 | U2 | 设计 + 可执行验证 | **verified** | 2026-08-20（D-002） | pg_dump -F c → pg_restore round-trip 实跑通过（count=2）；pg_basebackup 供规模选型 |

> 编号沿用 Root 信息项 I-001 / I-004（R5 最晚需要阶段），此处作旅程登记。

## 父目标

- `GOAL-001-store-dialects`

## 台账布局

五件套 + `01-decision/`、`02-execution/`、`03-audit/`、`attachments/`。
