---
id: A-001-self-checkpoints-b-c
doc: audit-entry
status: active
parent: GOAL-003-r2-codec-and-descriptor-m1-m2
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
source: self
verdict: conditional
---

# A-001（self）· GOAL-003 检查点 B/C

- **source**: self（编排器自审）
- **日期**: 2026-09-20
- **scope**: GOAL-003 检查点 B（15 个 descriptor 的 SQLite `Apply` + 目录断言）与 C（PG `ApplyPostgres` 显式 DDL + PG 侧断言）；对应实现 commit `c69ee93d`。
- **verdict**: `conditional`

## 成果（可核对）

| 判据 | 证据 | 结果 |
|------|------|------|
| 15 个 descriptor 落码（v73–v87） | `apps/api/modules/*/migration/vp040_temporal.go`（15 文件） | ✅ |
| v1–v72 canonical SQL / checksum 不变 | `TestCompiledMigrationCatalogOwnership`（冻结表前 72 行未改）、`TestMigrateFreshDB` 的 v1–v72 逐条断言 | ✅ |
| 真实 `MigrationChecksum` 已记录 | `internal/store/migrate_test.go` 冻结表 15 行 + `attachments/r2-v73-v87-generated-statements-v0.1.md`；两处由**独立路径**算出（runtime descriptor vs 生成器重算），值一致 | ✅ |
| PG `ApplyPostgres` 显式 DDL（无 `pgTimeColRe` 派生；毫秒族整数拆分式） | `vp040_temporal.go` 的 `*Postgres` 切片；`TestFullCatalogPostgresBootstrapIntegration` 在真实 PostgreSQL 15.4 通过 | ✅ |
| PG 侧类型/精度断言 | `postgres_test.go`：时间列 `timestamp with time zone` + `datetime_precision = 6`；21 名 leftover 集合；金额列 `bigint` | ✅ |
| SQLite 全量应用 v1–v87 | `TestMigrateFreshDB`（87 行）、各模块 provider 测试 | ✅ |
| `rebuildOperationLog` fail-closed 断言 | `sideTablesAbsent` + 方言探测查询；`D-019` §5 改动 1/2 落实，改动 3 走通用 F-5 执行器 | ✅（改动 3 的实现形态见偏差 1） |

## 偏差与开放项（编排器自评）

1. **改动 3 的实现形态**：`D-019` §5 改动 3 要求 v75「调用 `rebuildOperationLogWithSessions`（或与 F-5 同一通用 helper）；禁止 `pgRebuild`」。本轮 v75 走的是**通用字面语句执行器**（`temporalmigrate.Exec` + 与 F-5 相同的 14 步字面编排），**未**调用 `rebuildOperationLogWithSessions`；差别的实际后果是 `sideTablesAbsent` 断言对 v75 不生效（v75 自身在 rename 之前已 DROP 两张子表，结构上不可能命中该断言）。**待 independent 判定**是否等价。
2. **m0/m4 进入 checksum 输入**：冻结台账 §2 的 `m0 → m5` 顺序被理解为 checksum 覆盖 m0/m4 的查询语句；若审计认为 checksum 只应覆盖 DDL，则 15 个 checksum 需重算（记录在案，非静默假设）。
3. **生成器保留在树内**（`internal/store/vp040_generate_test.go`，`VP040_GENERATE=1` 才运行）。可复现性已实测（16/16 输出 SHA-256 不变）。是否允许「测试写源码」的工具留在仓库由审计判定。
4. **`F-I-005` residual**：本轮首次记录 v73+ 哈希，**已触发 `D-021` 的复审**；residual 是否可闭合取决于本轮 independent 复审，本自审**不**自证闭合。
5. **PG 侧无独立跨版本矩阵**：仅在 PG 15.4 上实测；`I-041-004`（15/16/17 跨版本）仍为 R3 前复核项（Root `D-017` §3 约束②）。

## 自审边界（诚实声明）

- 本自审以**可执行证据**（编译、测试、真实 PG、生成器复现）为主；未逐行复核 15 个生成文件与全部仓储 diff。
- M3 的模块改造由 4 个并行子代理完成，其**各自 scope 的自证报告**是二手证据；编排器只核对了聚合测试结果、`ORDER BY`/谓词保留的抽样、以及对若干高风险点（sentinel 谓词、扫描适配器、PG 语句类型）的直接验证。
- 以上偏差与边界**全部**提交 independent 复审；本文件不构成任何 required finding 的闭合。
