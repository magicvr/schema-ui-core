---
id: GOAL-004-r3-dual-dialect-ledger
title: R3 · 双方言台账对写（compiled catalog 两方言 apply + checksum）
status: active
parent: GOAL-001-store-dialects
created: 2026-08-20
updated: 2026-08-20
version: 0.2.0
progress: 4/5
plan_refs:
  - VP-013-store-dialects
primary_plan: VP-013-store-dialects
serves_summary: 交付 R3：把 compiled-global 迁移台账对写到 SQLite/PostgreSQL 两方言——Apply 形状转 kernel.Tx、逐迁移成对/可移植物理 SQL、store 双方言迁移运行器、checksum 仍绑 sqlite 历史文本；R4 前不缺省开发路径。
---

# GOAL-004 · R3 · 双方言台账对写

## 概述

Root 纲领 **R3**（依赖 R2 访问层）：现行 compiled-global 迁移台账（写作时最高版本 **0048**，开区时目录为准）在 **SQLite 与 PostgreSQL 两方言**上都能 apply，`checksum` 仍 fail-closed（R1 合同 v1.4 §4 / VP-013 退出 3）。

- `kernel.MigrationContribution.Apply / Reconcile` 形状转 **`func(kernel.Tx) error`**（§4）；store 迁移运行器经 `kernel.Store.Run` 驱动。
- 逻辑 schema 一份；物理 SQL 按方言**成对**或 `?` rebind 可移植文本（§1/§3）。
- **checksum 输入不变** = sqlite/canonical 历史 SQL + transform id；postgres 成对 SQL **不入 digest**（§4 禁止改 hash / 改 sqlite 历史）。
- 点名方言债（§3/§6）：`authsession/migration/migration.go` 的 `sqlite_master` / `PRAGMA table_info|foreign_key_list` / `COLLATE NOCASE` DDL；时间列 postgres `BIGINT`；非时间 INTEGER 宽度（wallet `balance_*`/`amount_delta`……）；布尔列等价形态；`LIKE` 大小写显式决策（模块运行时 `LIKE` 属 R4 收口，本目标只处理**迁移 Apply 路径**）。
- postgres 路径：fresh bootstrap apply 双方言 catalog；`WasFresh` 仍按 R1 v1.2/v1.3 语义。

## 非目标

- **不做**模块仓库公共面的 `*sql.Tx` 签名迁移 / 运行时 `INSERT OR IGNORE`、`LIKE`、`COLLATE NOCASE` 改写（归 **R4**）。
- **不做**存量 SQLite→PG in-place 升级策略或 dump/restore（归 **R5** / Root I-001）；不做 PG 备份合同（R5 / I-004）。
- 不引入 ORM；不改 Profile / Compose 默认（sqlite 仍为缺省开发路径）。

## 纲领路线图（P-001）

| 阶段 | 内容 | 状态 |
|------|------|------|
| T0 | R3 方案冻结：catalog 分列/成对形态、checksum 绑定、Apply→`kernel.Tx`、`LIKE`/`COLLATE`/时间/宽度/布尔落盘规则 | ✅ D-001（2026-08-20） |
| T1 | `kernel.MigrationContribution.Apply/Reconcile` → `func(kernel.Tx)`；store 迁移运行器走 `sqlTx` 适配（sqlite 保持绿） | ✅ E-002 / A-001（2026-08-20；全量测试 0 FAIL） |
| T2 | store postgres 迁移运行器：fresh bootstrap / apply postgres catalog / `schema_migrations` 台账（checksum 同 sqlite 绑定） | ✅ E-003 / A-002（T2a：运行器 live 证明；**生产解闸并入 T3**，随真实双方言 catalog） |
| T3 | 逐迁移对写：authsession `sqlite_master`/`PRAGMA`/`COLLATE NOCASE`、时间 `BIGINT`、非时间 INTEGER 宽度、布尔、可 rebind 文本；双 apply + checksum | ✅ E-004/E-005 + A-003/A-004：48 迁移全双写；全量 PG boot + 台账 + 幂等 + 系统级无 int 时间列检查；postgres open 解闸（store 级）；composition 路由移交 R4 |
| T4 | 证据与关门：sqlite 回归 + PG fresh bootstrap 全量 apply；self + independent（迁移/数据门禁） | 待做 |

## 成功标准

1. `apps/api` 缺省 sqlite 全量构建 + 测试绿（回归不破）。
2. 同一 compiled catalog 在 postgres fresh bootstrap 上可 apply 全部迁移，`schema_migrations` 含同名 + 同 checksum 记录（checksum 绑 sqlite 历史文本）。
3. 所有模块迁移的 `Apply/Reconcile` 可在 PG 上执行：无 `sqlite_master`/`PRAGMA`/`COLLATE NOCASE` 直抄；时间列 PG `BIGINT`；命名宽度/布尔已按 §3 落盘。
4. postgres 运行证据：对 `SCHEMA_UI_R2_PG_DSN` 指向的库可执行 T2 的迁移 apply 测试（或等价的 fresh bootstrap 验证）。
5. 未引 ORM；未改模块公共仓库签名、未碰运行时 SQL 债（留 R4）。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | compiled 台账当前逐迁移 SQL 方言债的**完整**清单（不是 R1 抽样） | T3 对写 | T3 前 | 迁移包代码扫描 | **verified** | 2026-08-20 E-004/E-005 | 48 迁移全部双写（BIGINT 时间/金额、CITEXT、partial index、operationlog rebuild）；系统级无 int 时间列检查 0 残留 |
| I-002 | required | 各迁移时间列单位、非时间 INTEGER 宽度、布尔列的**逐表/逐列**证据 | T3 对写 DDL | T3 每迁移前 | 代码/数据核对 | open | T3 按迁移进度复核 | 待确认 |
| I-003 | required | catalog 按方言分列 vs 单一 catalog 成对 SQL 的选择（checksum 绑定影响） | T0 方案 | T0 前 | R1 v1.4 §4 两种合法解取舍 | open | 责任人：本区编排 | 待确认（D-001 内裁） |
| I-004 | non-blocking | PG fresh bootstrap 的时间/备份语义是否与 sqlite 完全一致 | T4 证据 | T4 | live PG 对比 | collecting | T4 前 | 待确认 |

> Root I-001（SQLite→PG 升级策略）与 I-004（备份合同）为 R5，不构成本目标到期门禁。

## 父目标

- `GOAL-001-store-dialects`

## 台账布局

五件套 + `01-decision/`、`02-execution/`、`03-audit/`、`attachments/`。
