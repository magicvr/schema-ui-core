---
id: GOAL-003-r2-codec-and-descriptor-m1-m2
title: R2 · 共享 codec 与 15 个 conversion descriptor（M1/M2）
status: active
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.3.0
progress: 3/4
plan_refs:
  - VP-040-timestamptz-persistence-contract
primary_plan: VP-040-timestamptz-persistence-contract
serves_summary: 承接 Root R2 的 M1/M2：落码共享时间 codec 与 15 个 conversion descriptor（SQLite `Apply` + PG `ApplyPostgres`），按 D-017 约定首次记录真实 MigrationChecksum，并保持 v1–v72 canonical SQL/checksum 不变。
---

# GOAL-003 · R2 · 共享 codec 与 15 个 conversion descriptor（M1/M2）

## 概述

本子目标承接 `GOAL-001-timestamptz-persistence-contract` 的 **R2 阶段 M1/M2 检查点**。R1（`GOAL-002`）已关门（`done · 4/4`，开放 required = 0），C2/C3 合同与边界均已冻结；本目标把其中**与 codec 和迁移描述符相关**的部分落成可执行代码。

**子目标 slug 与编号经用户 2026-09-20 确认**（AGENTS §11 硬约束：禁止静默默认 slug）。

## 范围

| # | 交付 | 对应检查点 |
|--:|------|-----------|
| 1 | **共享时间 codec**（`apps/api/internal/temporal`，包名候选）：`FromUnix` / `FromUnixMilli`、canonical fixed-6 UTC formatter/parser、sentinel / NULL helper、`UTC().Truncate(time.Microsecond)` | **M1** |
| 2 | codec **可执行单测**（fixed-6、负毫秒 floor、999 ms 无进位、公元 9999 年、sentinel 0、NULL、非法值 fail closed） | **M1** |
| 3 | **SQLite 侧 15 个 conversion descriptor**（v73–v87）按 `D-019` F-5 子女先行 / 裸四步；含子表重建与回填、索引最后建 | **M2** |
| 4 | **PG 侧 `ApplyPostgres` 显式 DDL**（**禁止** `pgTimeColRe` 派生）；毫秒族用**整数拆分式**（原式已弃用） | **M2** |
| 5 | **canonical SQL 切片 + 真实 `MigrationChecksum` 记录**（`D-017` 约定：单 checksum / SQLite DDL 切片 / `transform_id` 无方言后缀）→ **`D-021` residual 复审触发点** | **M2** |
| 6 | **`rebuildOperationLog` fail-closed 断言**（`D-019` §5 改动 1–3） | **M2** |
| 7 | v73+ 追加后的 `migrate_test.go` 目录断言（v1–v72 逐条冻结不变） | **M2** |

## 非目标（本子目标**不**做）

- **不做**仓储层读写/谓词改造与双方言回归（属 R2 **M3**，另立子目标）。
- **不做** `postgres_test.go` 的金额列断言拆分与 leftover 21 名补入（属 **M3**；但左侧「时间列断言」随 M2 的 v73+ 追加一并处理）。
- **不做**公共 wire formatter 实施（**用户裁决归 R3**，`D-016`）。
- **不做** Backup Port 类型表面与 provider（归 R2 **M4 前**，`D-016` 第 9 项，另立子目标）。
- **不做**边界测试重定向到真实迁移（属 **M3**，`D-020`）。
- **不**修改 v1–v72 的 canonical SQL / checksum / identity / Apply 语义。
- **不**引入 ORM / 第三库 / Redis / MQ / 多实例；**不**把驱动类型泄漏进公共契约。

## 红线

- 不把「R1 已冻结的合同」在落码时**静默改写**；任何偏离必须回到 `GOAL-002` 决策或新开决策。
- 不把 `D-021` 的 `accepted-residual` 读作「哈希已验证」——本目标正是**首次记录哈希**之处，记录后必须触发独立复审。
- 不以「测试通过」代替独立审计结论。

## 成功标准（本子目标检查点，用于 `progress` 派生）

| 检查点 | 判据 | 状态 |
|--------|------|------|
| **A** | 共享 codec 落码 + 单测全绿（含 `D-018` 的 Go 对拍用例） | **completed**（E-002：`apps/api/internal/temporal`，11 个单测全绿，`go vet`/`go build` 通过，依赖边界实测 `errors fmt time`） |
| **B** | 15 个 descriptor 的 SQLite `Apply` 落码，v73+ 目录断言通过，**v1–v72 逐条不变** | **completed**（E-003：15 个 descriptor 落码于 `modules/*/migration/vp040_temporal.go`（由 v72 live `sqlite_master` 机械导出）；`migrate_test.go` 目录表追加 v73–v87 真实 checksum；`TestMigrateFreshDB`/`TestCompiledMigrationCatalogOwnership`/`TestCompleteFingerprintTracksCatalogHead` 通过；v1–v72 行逐条不变） |
| **C** | PG `ApplyPostgres` 显式 DDL 落码（毫秒族整数拆分式），PG 侧类型断言就位 | **completed**（E-003：显式 `ALTER … TYPE timestamptz(6) USING …` 逐列 DDL；`TestFullCatalogPostgresBootstrapIntegration` 在真实 PostgreSQL 15.4 上 **ok**；`postgres_test.go` 时间列断言改 `timestamp with time zone` + `datetime_precision = 6`，金额列保持 `bigint`） |
| **D** | canonical SQL 与**真实 `MigrationChecksum` 落盘**，并按 `D-021` **发起 independent 复审**（residual 复审触发） | pending |

`progress: 3/4` 由 A～D 等权派生；**不**放行阶段、**不**关闭 finding、**不**推导 `done`。

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 | 影响门禁 | 状态 | 证据 |
|----|------|----------|----------|------|------|
| I-041-001 | required | Go codec 公共 API 形态（签名、错误分类、导出面）与 `D-018` Go 对拍验收一致 | A / M1 | **collecting（本目标内定稿）** | 定稿落盘于本目标 `01-decision/`；未经用户裁决的技术细节由独立审计复审 |
| I-040-001 | required（继承） | 逐列编解码/精度合同 | A/B/C | **verified**（R1 关门时） | `r1-c2-per-column-conversion-contract-v1.0-fc.md` |
| I-040-003 | required（继承） | 双方言原地转换、失败恢复、备份依赖 | B/C | **verified**（R1 关门时） | `r1-c3-backup-recovery-boundary-v1.0-fc.md`；`D-019` |

## 父目标

- `GOAL-001-timestamptz-persistence-contract`（R2 检查点）

## 关门条件

本子目标只有在 **A～D 全部完成**、`GOAL-003/03-audit` 的 self 与 independent 意见均已落盘、required finding 合法闭合后，才可由编排器静默关门；**`D-021` residual 的复审结论必须落盘**（不得以「哈希已记录」代替复审）。关门**不**放行 M3 之前的任何 schema 变更。
