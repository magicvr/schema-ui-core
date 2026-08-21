---
id: GOAL-005-r4-repository-surface
doc: decision
status: done
parent: GOAL-001-store-dialects
created: 2026-08-20
updated: 2026-08-20
version: 0.1.0
---

# 决策记录 · GOAL-005

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 全量 `*sql.Tx`/方言 SQL 泄漏面 | S0/S2 | S0 前补全 | 代码扫描 | collecting | — | GOAL-002 E-001 + S0 |
| I-002 | required | 运行时 SQL 债改写决策（LIKE/COLLATE/INSERT OR IGNORE/RETURNING/instr） | S2 | 每处前 | 逐处核对 + 测试 | verified | 2026-08-20（A-005） | ON CONFLICT / LOWER LIKE / CAST 形态；无 RETURNING 需求（R3 已处理） |
| I-003 | non-blocking | postgres 生产启动运维面 | S4 验收 | S4 | OpenOptions/文档 | collecting | S4 | R2/T3 已留字段 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-20 | R4 方案：收口形状、运行时 SQL 债规则、postgres 启动 | accepted | [D-001-r4-plan.md](01-decision/D-001-r4-plan.md) |
