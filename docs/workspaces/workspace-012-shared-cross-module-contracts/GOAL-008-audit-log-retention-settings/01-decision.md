---
id: GOAL-008-audit-log-retention-settings
doc: decision
status: done
parent: GOAL-001-shared-cross-module-contracts
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
---

# 决策记录 · GOAL-008

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 保留天数默认与过期删/归档 | S0 | S0 结束前 | 用户书面策略 | verified | — | D-001 |
| I-002 | required | 审计模式 | S2 | S1 实施前 | 按 data 生命周期分级 | verified | — | D-002：independent + grok-build；先 self |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-19 | 可配置保留：默认 90 天归档 | accepted | [D-001-configurable-retention.md](01-decision/D-001-configurable-retention.md) |
| D-002 | 2026-08-19 | 关门审计模式 independent | accepted | [D-002-audit-mode.md](01-decision/D-002-audit-mode.md) |
