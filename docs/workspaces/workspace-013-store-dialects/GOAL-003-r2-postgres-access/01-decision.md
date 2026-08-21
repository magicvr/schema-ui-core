---
id: GOAL-003-r2-postgres-access
doc: decision
status: active
parent: GOAL-001-store-dialects
created: 2026-08-20
updated: 2026-08-20
version: 0.1.0
---

# 决策记录 · GOAL-003

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | R2 证据边界（Open+Ping vs 现行 `/readyz`） | S1 方案 | S1 前 | R1 合同 v1.3/v1.4 §2 | verified | 2026-08-20 | v1.4「R2 证据边界」 |
| I-002 | required | 本目标审计模式 | S5 关门 | S5 前 | Root D-001 第 5 条 | verified | 2026-08-20 | self 先行；R2 实现后 independent；R3/R5 independent |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-20 | R2 方案：访问层边界与实施切面 | accepted | [D-001-r2-plan.md](01-decision/D-001-r2-plan.md) |
