---
id: GOAL-005-r4-repository-surface
title: R4 · 仓库公共面收口（*sql.Tx → kernel.Store/Tx + 运行时 SQL 债）
status: active
parent: GOAL-001-store-dialects
created: 2026-08-20
updated: 2026-08-20
version: 0.1.0
progress: —
plan_refs:
  - VP-013-store-dialects
primary_plan: VP-013-store-dialects
serves_summary: 交付 R4：handler/jobs/模块仓库公共契约去掉 `*sql.Tx` 与驱动类型，仓库改走 `kernel.Store`/`kernel.Tx`；改写运行时 SQL 方言债（INSERT OR IGNORE / LIKE / COLLATE NOCASE / RETURNING / 布尔）；接入 composition 的 postgres DSN 启动，完成 postgres 生产向应用路径。
---

# GOAL-005 · R4 · 仓库公共面收口

## 概述

Root 纲领 **R4**（依赖 R1；可与 R3 部分并行，现 R3 已 done）：把 Handler / jobs / 模块公共契约中的 **`*sql.Tx` 与驱动类型全部收口**，仓库改经内核端口（`kernel.Store`/`kernel.Tx`），并改写 R1 v1.4 §3 点名的**运行时 SQL 方言债**。R4 完成后，composition 可把 postgres DSN 接入完整应用启动（模块全部讲 kernel.Tx），为 R5 双路径验收铺路。

## 非目标

- 不做 R3 已完成的迁移对写；不做 R5 升级策略 / 备份合同（I-001/I-004）；不引 ORM；不改 Profile/Compose 默认（sqlite 仍是缺省开发路径）。

## 纲领路线图（P-001）

| 阶段 | 内容 | 状态 |
|------|------|------|
| S0 | 泄漏面扫描（I-003 补全）：`*sql.Tx`/驱动类型/方言 SQL 全量清单 | 待做（GOAL-002 E-001 已有部分） |
| S1 | 内核端口接缝：把各模块 `TxRunner`/`WithTx(ctx, func(*sql.Tx))` 接口改为 `kernel.Store`/`func(kernel.Tx)`（或 `Run`） | 待做 |
| S2 | 逐模块仓库迁移签名 + SQL 债改写（operationlog `INSERT OR IGNORE`→`ON CONFLICT DO NOTHING`；wallet/recyclebin `LIKE`→显式 `ILIKE`/校对决策；users/roles `ORDER BY … COLLATE NOCASE`→CITEXT/LOWER；插入取 id 用 `RETURNING`；布尔 `INTEGER` 0/1 保持并按 R1 落盘） | 待做 |
| S3 | jobs / handler 公共面收口（`CommitFunc` 等 `func(kernel.Tx)`） | 待做 |
| S4 | composition postgres 启动路由 + 运行证据（postgres DSN 启动、readyz 模块门禁全绿） | 待做 |
| S5 | 关门：sqlite 全量回归 + postgres 生产向运行验收；self + independent（compatibility/production 门禁） | 待做 |

## 成功标准

1. `apps/api` 缺省 sqlite 全量回归 0 FAIL（不改缺省行为）。
2. handler / jobs / 模块公共契约不再 import 具体驱动、不再出现 `*sql.Tx`（grep 可核对）。
3. 运行时方言债已改写且行为等价（大小写/布尔/upsert 语义与 sqlite 现状一致或书面落盘差异）。
4. **postgres DSN 完整启动**：compiled catalog apply + `readyz` 模块门禁全绿（live PG 运行证据）。
5. 未引 ORM；未改 Profile/Compose 默认；未重做 R3。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 全量 `*sql.Tx`/方言 SQL 泄漏面（handler/jobs/模块） | S0/S2 | S0 前补全 | 代码扫描（grep） | collecting | 责任人：本区编排 | GOAL-002 E-001 + R3 迁移期扫描；S0 补全到 handler/jobs |
| I-002 | required | 各模块运行时 SQL 债的具体改写决策（LIKE 大小写、COLLATE NOCASE 查询侧、INSERT OR IGNORE、RETURNING） | S2 每处 | 每处前 | 逐处核对 + 测试 | open | — | 待确认（S1/S2 落盘） |
| I-003 | non-blocking | postgres 生产向启动的运维面（连接池/SSL/超时键） | S4 验收 | S4 | OpenOptions 扩展 + 文档 | collecting | S4 前 | R2 已留字段；T3 用 ConnectTimeout |

> Root I-001（SQLite→PG 升级策略）与 Root I-004（PG 备份合同）为 R5，不构成本目标到期门禁。

## 父目标

- `GOAL-001-store-dialects`

## 台账布局

五件套 + `01-decision/`、`02-execution/`、`03-audit/`、`attachments/`。
