---
id: GOAL-004-r2-repository-and-predicate-rewrites
title: R2 · 仓储读写与谓词改造（M3）
status: done
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 1.0.0
progress: 3/3
plan_refs:
  - VP-040-timestamptz-persistence-contract
primary_plan: VP-040-timestamptz-persistence-contract
serves_summary: 承接 Root R2 的 M3：把 v73–v87 已改成 TEXT/timestamptz(6) 的时间列的全部读写路径与谓词改造成 codec 语义（无过渡期），完成双方言回归（含真实 PG 路径）、边界测试重定向、金额列断言拆分与 leftover 21 列补入，使全仓测试转绿。
---

# GOAL-004 · R2 · 仓储读写与谓词改造（M3）

## 概述

本子目标承接 `GOAL-001-timestamptz-persistence-contract` 的 **R2 阶段 M3 检查点**（`D-016` §4）。M1/M2 由 `GOAL-003` 承担；按 Root `D-017`（用户 2026-09-20 裁决），M2 的 schema 落码**不单独提交**（`D-018` 无过渡期 → 仓储层未改造前全仓测试红，AGENTS §6b fail closed），**M3 转绿后与 M2 一并提交**。

**子目标编号与 slug 经用户 2026-09-20 确认**（AGENTS §11 硬约束：禁止静默默认 slug）：`GOAL-004-r2-repository-and-predicate-rewrites`。

## 范围

| # | 交付 | 判据来源 |
|--:|------|----------|
| 1 | 各模块时间列 `int64` → `time.Time` / `sql.NullTime`（或 codec `Value`/`NullValue`）读写改造；**无过渡期** | `D-018`；逐列合同 `read/write` 列 |
| 2 | 谓词按 `r1-c2-predicate-exact-sql-v1.0-fc.md` 的 exact new SQL 改造（含 `#5/#6/#20` 的 `IS NULL` 分支、`#61` 去 `COALESCE(…,0)`、`#72/#73` 去 `> 0` 判定、`ORDER BY` 保留） | 谓词附件 §1–§5 |
| 3 | **金额列断言拆分**：`wallet_accounts.balance_total` / `wallet_ledger_entries.amount_delta` 保持 `bigint`；时间列断言改 `timestamp with time zone`（精度 6） | `D-021` residual ②；`postgres_test.go:291-307` |
| 4 | leftover 21 列名补入 PG 断言集合；`postgres_test.go:312-316` 缺失七列补齐 | `D-021` residual ③；`r1-v73-owner-allocation-draft-v0.1.md` L37–L38 |
| 5 | **边界测试重定向**：`apps/api/internal/w040contracttest/` 用例改指真实迁移（`D-020`） | `D-020` |
| 6 | **双方言回归**：含至少一条真实 PG 路径；破坏性 migration/round-trip **只能**作用于一次性/专用测试 database（Root `D-017` §3 约束①） | `D-016` M3；`D-017` |
| 7 | 全仓 `go test ./...` **转绿**（M2+M3 合并提交的前置） | `D-017` §1 |

## 非目标（本子目标**不**做）

- **不**改 v1–v87 的 canonical SQL / checksum / identity / Apply 语义（M2 已冻结的 descriptor 只读）。
- **不**做 Backup Port 类型表面与 provider（归 R2 **M4 前**，另立子目标）。
- **不**做公共 wire formatter 实施（**用户裁决归 R3**，`D-016`）。
- **不**做 VP-020 展示/输入时区回归矩阵与其用例 ID（属 **R3**）。
- **不**引入 ORM / 第三库 / Redis / MQ / 多实例；**不**把驱动类型泄漏进 handler/模块公共契约。
- **不**以 Docker 隔离 PG / 版本矩阵作为本子目标的完成条件（Root `D-017` §3 约束②）。

## 红线

- 不把 M2 已冻结的 descriptor/canonical SQL 静默改写；发现偏离必须回到 `GOAL-002`/`GOAL-003` 决策。
- 破坏性 migration / round-trip **不得**作用于共享业务 schema（Root `D-017` §3 约束①）。
- 不以「本地临时容器验证」冒充生产就绪证据（`D-021`）。
- 不以「测试通过」代替独立审计结论。

## 成功标准（本子目标检查点，用于 `progress` 派生）

| 检查点 | 判据 | 状态 |
|--------|------|------|
| **A** | 全部时间列读写/谓词改造落码（含 `internal/**` 与 `modules/**` 的 store/repository 层） | **completed**（E-002：store 读写适配器 + 4 路并行模块改造；sentinel 谓词按 D-001 §2 逐条落实） |
| **B** | 双方言回归：SQLite 全绿 + 至少一条真实 PG 路径绿；金额列断言拆分与 leftover 21 列补入完成 | **completed**（E-002：`postgres_test.go` 时间列改 `timestamp with time zone` + `datetime_precision = 6`、金额列保持 `bigint`、21 名清单；`TestFullCatalogPostgresBootstrapIntegration` 与 `TestPurchasePostgresAcceptance` 在真实 PG 15.4 上通过） |
| **C** | 边界测试重定向到真实迁移；全仓 `go test ./...` 绿 → 与 M2 一并提交 | **completed**（E-002：`migration_boundaries_test.go` 真实迁移矩阵全绿；全仓 `go test -count=1 ./...` 63/63 包 ok） |

`progress: 3/3` 由 A～C 等权派生；**不**放行阶段、**不**关闭 finding、**不**推导 `done`。

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 | 影响门禁 | 状态 | 证据 |
|----|------|----------|----------|------|------|
| I-041-003 | required（继承） | PG 侧可执行验证环境 | A/B/C | **verified（用户 2026-09-20，Root `D-017` §3）** | 常驻 PostgreSQL 15.4 实测执行；约束：破坏性 migration 只可作用于一次性/专用测试 database |
| I-040-001 | required（继承） | 逐列编解码/精度合同 | A/B | **verified**（R1 关门时） | `r1-c2-per-column-conversion-contract-v1.0-fc.md` |
| I-040-003 | required（继承） | 双方言原地转换、失败恢复、备份依赖 | B/C | **verified**（R1 关门时） | `r1-c3-backup-recovery-boundary-v1.0-fc.md`；`D-019` |
| I-041-005 | required（本目标新增） | 双方言绑定形态与统一入口 | A | **verified（本目标 `D-001` §1 定稿 + 实测）** | `D-001` §1：域类型 + store 适配器归一化（写：`bindSQLiteArgs`/`bindPostgresArgs`；读：`scan.go`）；待 independent 复审确认 | 本轮实测：SQLite 绑 `time.Time` 会写 36 字符驱动格式（`2026-09-20 12:40:48.814963 +0000 UTC`），违反 canonical 27 字符合同 |

## 父目标

- `GOAL-001-timestamptz-persistence-contract`（R2 检查点）

## 关门条件

本子目标在 **A～C 全部完成**、`03-audit` 的 self 与 independent 意见落盘、required finding 合法闭合、`I-041-005` 由证据（或用户裁决）关闭后才可关门；关门**不**放行 M4 的 Backup Port 及其门禁。
