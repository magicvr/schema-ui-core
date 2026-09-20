---
id: GOAL-003-r2-codec-and-descriptor-m1-m2
doc: execution
status: active
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.2.0
---

# 执行记录 · GOAL-003

## 执行索引

| E-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|--------|------|
| E-001 | 2026-09-20 | R2 边界冻结与 GOAL-003 立项（用户确认 slug） | recorded | `02-execution/E-001-r2-boundary-and-goal-created.md` |
| E-002 | 2026-09-20 | 共享 codec 落码与检查点 A 完成（I-041-001 定稿） | recorded | `02-execution/E-002-codec-implemented-checkpoint-a.md` |
| E-003 | 2026-09-20 | 15 个 descriptor 落码与检查点 B/C（SQLite `Apply` + PG 显式 DDL + 目录/指纹断言） | recorded | `02-execution/E-003-descriptors-and-pg-landed-checkpoints-b-c.md` |

## 事实边界

> 只登记已经发生且有证据的事实。
>
> - **检查点 A 已完成**：`apps/api/internal/temporal` 落码 + 11 个单测全绿 + 依赖边界实测 `errors fmt time`。
> - **检查点 B/C 已完成**（E-003）：15 个 descriptor 的 SQLite `Apply` 与 PG `ApplyPostgres` 落码；冻结目录表记录真实 checksum；真实 PostgreSQL 15.4 上 `TestFullCatalogPostgresBootstrapIntegration` 通过；PG 类型/精度断言与金额列 `bigint` 断言就位。
> - **检查点 D 尚未完成**：canonical SQL 清单附件已生成（`attachments/r2-v73-v87-generated-statements-v0.1.md`），但**`D-021` residual 的 independent 复审尚未发起/落盘**，故 D 仍为 pending；本批落码**不**放行任何 M3 之前/之后的 schema 变更门禁。
> - M3（仓储与谓词改造）由 `GOAL-004-r2-repository-and-predicate-rewrites` 承担；按 Root `D-017`，M2 与 M3 的提交在同一绿态下进行。
