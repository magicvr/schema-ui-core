---
id: GOAL-001-store-dialects
doc: decision
status: active
parent: null
created: 2026-08-20
updated: 2026-08-20
version: 0.1.0
---

# 决策记录 · GOAL-001

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 存量 SQLite → PostgreSQL：in-place 或 dump/restore residual | R5 / 退出 2 | R5 前 | 原型或书面 residual | open | R3 后复核 | 待确认 |
| I-002 | required | PG 驱动选型（非 ORM） | R2 方案 | R2 实施前 | D-002（pgx v5 stdlib + 编译证据） | verified | — | D-002；`go get` 编译通过 |
| I-003 | non-blocking | `*sql.Tx` 公共面泄漏清单 | R4 范围 | R4 方案 | 扫描 | collecting | R4 前补全 | GOAL-002 E-001 部分清单 |
| I-004 | required | PG 备份/恢复合同 | R5 / 退出 4 | R5 前 | 设计 | open | 可与 I-001 一并裁 | 待确认 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-20 | 开区 scaffold 与 A1 纲领路线图 | accepted | [D-001-workspace-root-establishment.md](01-decision/D-001-workspace-root-establishment.md) |
| D-002 | 2026-08-20 | R2 PostgreSQL 驱动选型（I-002 → verified） | accepted | [D-002-postgres-driver-selection.md](01-decision/D-002-postgres-driver-selection.md) |
