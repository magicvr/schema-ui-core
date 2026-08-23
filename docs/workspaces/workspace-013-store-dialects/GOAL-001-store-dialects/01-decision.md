---
id: GOAL-001-store-dialects
doc: decision
status: active
parent: null
created: 2026-08-20
updated: 2026-08-23
version: 0.2.0
---

# 决策记录 · GOAL-001

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 存量 SQLite → PostgreSQL：in-place 或 dump/restore residual | R5 / 退出 2 | R5 前 | 原型或书面 residual | **accepted-residual** | 复审触发 = 出现 in-place/内置搬运需求时另立目标 | D-002 有界 residual（本 VP 不提供 SQLite→PG 搬运器；存量路径 = fresh bootstrap + 运维自备导出/回放，2026-08-20）+ `TestPostgresDataMigrationPrototype` live round-trip PASS；Root 独立关门审计 A-001 pass（2026-08-21）。2026-08-23 回写对齐 00-meta 口径（[workspace-010 GOAL-033](../../../workspace-010-design-implementation-conformance/GOAL-033-w22-residual-closeout/00-meta.md) B6 复核） |
| I-002 | required | PG 驱动选型（非 ORM） | R2 方案 | R2 实施前 | D-002（pgx v5 stdlib + 编译证据） | verified | — | D-002；`go get` 编译通过 |
| I-003 | non-blocking | `*sql.Tx` 公共面泄漏清单 | R4 范围 | R4 方案 | 扫描 | collecting | R4 前补全 | GOAL-002 E-001 部分清单 |
| I-004 | required | PG 备份/恢复合同 | R5 / 退出 4 | R5 前 | 设计 | **verified** | — | pg_dump -F c → pg_restore round-trip：catalog 48 迁移 / 35 表 checksum 一致（GOAL-006 D-002/E-002，2026-08-20；独立复核 2026-08-21）。2026-08-23 回写对齐 00-meta 口径 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-20 | 开区 scaffold 与 A1 纲领路线图 | accepted | [D-001-workspace-root-establishment.md](01-decision/D-001-workspace-root-establishment.md) |
| D-002 | 2026-08-20 | R2 PostgreSQL 驱动选型（I-002 → verified） | accepted | [D-002-postgres-driver-selection.md](01-decision/D-002-postgres-driver-selection.md) |
