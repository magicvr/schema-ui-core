---
id: GOAL-004-r3-dual-dialect-ledger
doc: decision
status: active
parent: GOAL-001-store-dialects
created: 2026-08-20
updated: 2026-08-20
version: 0.1.0
---

# 决策记录 · GOAL-004

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 逐迁移 SQL 方言债完整清单 | T0 / T3 | T0 前 | 扫描 | collecting | — | R2 已扫部分（见 E-001） |
| I-002 | required | 时间单位 / 非时间宽度 / 布尔 逐列证据 | T3 对写 | 每迁移前 | 核对 | open | T3 | 待确认 |
| I-003 | required | catalog 分列 vs 成对（checksum 绑定） | T0 | T0 前 | v1.4 §4 | open | — | D-001 内裁 |
| I-004 | non-blocking | PG fresh bootstrap 语义 | T4 | T4 前 | live 对比 | collecting | T4 | 待确认 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-20 | R3 方案：Apply 形状、catalog 形态、checksum 绑定与对写规则 | accepted | [D-001-r3-plan.md](01-decision/D-001-r3-plan.md) |
