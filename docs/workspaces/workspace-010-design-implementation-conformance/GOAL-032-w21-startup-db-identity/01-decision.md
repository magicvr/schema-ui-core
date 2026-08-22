---
id: GOAL-032-w21-startup-db-identity
doc: decision
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-08-22
updated: 2026-08-22
version: 0.2.0
---

# 决策记录 · GOAL-032 · W21

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 无 ledger 时「完整我方库」最小表集合 | S1 | S1 | 现场故障表清单 + catalog | **verified** | — | D-001 四表；D-002 收紧为含 catalog 头表 |
| I-002 | required | 是否新建历史表还是沿用 `schema_migrations` | S1 | S1 | EF 对照 | **verified** | — | D-001：沿用 |

## 决策索引

| 编号 | 标题 | 日期 | 状态 | 摘要 |
|------|------|------|------|------|
| [D-001](01-decision/D-001-identity-and-plan-freeze.md) | 启动身份/计划合同（含 EF 历史表对照） | 2026-08-22 | accepted | ledger 权威 + 无 ledger 身份探针；restore/partial 以 D-002 为准 |
| [D-002](01-decision/D-002-a001-response.md) | 响应 A-001 F-001～F-003 | 2026-08-22 | accepted | restore 绑定 catalog 头；unsafe refuse；sqlite V-MIG-03 保留 |
