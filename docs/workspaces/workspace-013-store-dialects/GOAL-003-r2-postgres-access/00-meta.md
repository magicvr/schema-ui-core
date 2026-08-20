---
id: GOAL-003-r2-postgres-access
title: R2 · PostgreSQL 接入（驱动、连接池、Open/Ping/WasFresh）
status: done
parent: GOAL-001-store-dialects
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
progress: 6/6
plan_refs:
  - VP-013-store-dialects
primary_plan: VP-013-store-dialects
serves_summary: 落地 R2：pgx v5 stdlib 驱动的 postgres 访问层（连接 + Ping + WasFresh 探测打开；空 catalog），加 db.dialect / db.dsn 配置校验；sqlite 缺省路径不变，不 apply sqlite 专用 catalog 到 postgres。
---

# GOAL-003 · R2 · PostgreSQL 接入

## 概述

Root 纲领 **R2**（依赖 R1 合同 v1.4）：把 PostgreSQL 方言接入内核持久化端口——

- 依赖：**pgx v5 stdlib**（Root D-002 / I-002 `verified`）。
- 交付 `kernel.Store` / `kernel.Tx` 接口，`store.Open(ctx, OpenOptions, catalog)` 按方言分发。
- postgres 本拍**只连 + Ping + WasFresh**；空/`nil` catalog 用作探测打开；非空 catalog **fail closed**（现行 compiled catalog 含 SQLite 专用 SQL，R3 双方言对写前不得半执行）。
- 配置键 `db.dialect` / `db.dsn` + 启动 fail-closed 校验。
- SQLite 保持缺省与现行 `OpenWithCatalog` 路径（模块公共面 `WithTx(*sql.Tx)` 仍为 R4 工作）。

## 非目标

- 不对写迁移（R3）；不收口模块仓库签名 / `*sql.Tx` 泄漏（R4）；不接备份合同 / 升级策略（R5）；不改 Compose 默认依赖。
- 不引入 ORM；不重建业务 Repository。

## 纲领路线图（P-001）

| 阶段 | 内容 | 状态 |
|------|------|------|
| S0 | 驱动选型（I-002） | ✅ Root D-002 / I-002 verified |
| S1 | R2 方案（OpenOptions / dialect 分发 / config 校验 / 证据边界） | 本目标 D-001 |
| S2 | 驱动 + `kernel.Store`/`kernel.Tx` + `store.Open` + rebind | 待做（本回合） |
| S3 | config `db.dialect`/`db.dsn` 字段 + 校验谓词 | 待做（本回合） |
| S4 | postgres Store：连接 / Ping / WasFresh / 空 catalog 探测，非空 fail closed | 待做（本回合） |
| S5 | 测试（sqlite 运行时 + rebind + config + pg 探测 gate）+ self 审计 + independent 审计 | ✅ 2026-08-20（A-001 self pass；A-002 independent pass；A-003 关闭全部 recommended；live PG 探测通过） |

## 成功标准

1. `apps/api` 在缺省 sqlite 下构建并通过既有测试（回归不破）。
2. `store.Open` 支持 `Dialect=postgres`：pgx 驱动连接 + `Ping` + `WasFresh`；空 catalog 探测成功，非空 catalog **fail closed** 且不产生半执行。
3. `db.dialect` / `db.dsn` 配置与校验落地：dialect ∈ {空, sqlite, postgres}；sqlite 下 DSN 必须空、postgres 下 DSN 必须非空；`db.path` 保持文件路径形状（含扩展名谓词）。
4. `kernel.Tx`/`Store` 接口与 `?`→`$n` rebind 有单元测试；pg 运行时测试以 env 门控（无 PG 时 skip），不作为本地必跑失败项。
5. 未引入 ORM；未改模块公共签名 / 未写任何 postgres DDL 到模块仓库。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | R2 证据边界（Open+Ping vs 现行 `/readyz`） | S1 方案 | S1 前 | R1 合同 v1.3/v1.4 §2 | **verified** | 2026-08-20 | v1.4「R2 证据边界」：本拍 postgres 可核对 = Open+Ping（及连接池），非 composition 全量 bootstrap，非 `/readyz` 200 |
| I-002 | required | 本目标审计模式 | S5 关门 | S5 前 | D-001（Root）第 5 条 | **verified** | 2026-08-20 | Root D-001 / D-002：方案+实施 self 先行，R2 实现后 independent（项目默认 grok build）另做；R3/R5 仍为 independent 门禁 |

## 父目标

- `GOAL-001-store-dialects`

## 台账布局

五件套 + `01-decision/`、`02-execution/`、`03-audit/`、`attachments/`。
