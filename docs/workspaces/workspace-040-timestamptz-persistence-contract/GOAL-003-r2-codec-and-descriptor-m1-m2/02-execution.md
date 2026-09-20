---
id: GOAL-003-r2-codec-and-descriptor-m1-m2
doc: execution
status: active
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.3.0
---

# 执行记录 · GOAL-003

## 执行索引

| E-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|--------|------|
| E-001 | 2026-09-20 | R2 边界冻结与 GOAL-003 立项（用户确认 slug） | recorded | `02-execution/E-001-r2-boundary-and-goal-created.md` |
| E-002 | 2026-09-20 | 共享 codec 落码与检查点 A 完成（I-041-001 定稿） | recorded | `02-execution/E-002-codec-implemented-checkpoint-a.md` |
| E-003 | 2026-09-20 | 15 个 descriptor 落码与检查点 B/C（SQLite `Apply` + PG 显式 DDL + 目录/指纹断言） | recorded | `02-execution/E-003-descriptors-and-pg-landed-checkpoints-b-c.md` |
| E-004 | 2026-09-20 | 检查点 D：canonical SQL / 真实 checksum 落盘 + `D-021` residual 复审 + 独立审计修正 | recorded | `02-execution/E-004-checkpoint-d-checksums-and-residual-rereview.md` |

## 事实边界

> 只登记已经发生且有证据的事实。
>
> - **检查点 A 已完成**：`apps/api/internal/temporal` 落码 + 11 个单测全绿 + 依赖边界实测 `errors fmt time`。
> - **检查点 B/C 已完成**（E-003）：15 个 descriptor 的 SQLite `Apply` 与 PG `ApplyPostgres` 落码；冻结目录表记录真实 checksum；真实 PostgreSQL 15.4 上 `TestFullCatalogPostgresBootstrapIntegration` 通过；PG 类型/精度断言与金额列 `bigint` 断言就位。
> - **检查点 D 已完成**（E-004）：canonical SQL 与 15 个真实 `MigrationChecksum` 落盘于 `attachments/r2-v73-v87-generated-statements-v0.1.md` 并由冻结表锁定；`D-021` residual 的 independent 复审已落盘（A-002）并判可 `fixed` 闭合（编排器 A-003 `pass`）；生成器对同一历史字节级复现。本子目标据此静默关门（`done · 4/4`）。**本关门不放行 R2，也不等于 M3/M4 完成。**
> - M3（仓储与谓词改造）由 `GOAL-004-r2-repository-and-predicate-rewrites` 承担；按 Root `D-017`，M2 与 M3 的提交在同一绿态下进行。
