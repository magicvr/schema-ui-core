---
id: GOAL-021-wallet-deduct-frozen
doc: decision
status: active
parent: GOAL-001-admin-functional-modules
created: 2026-08-16
updated: 2026-08-16
version: 0.1.0
---

# 决策记录 · GOAL-021-wallet-deduct-frozen

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | deduct_frozen 语义与拒绝条件 | 方案 | S1 | — | **verified** | — | D-001 §1 |
| I-002 | required | 迁移重建策略 | 方案 | S1 | — | **verified** | — | D-001 §3 |
| I-003 | non-blocking | 幂等比对兼容影响 | 方案 | S1 | — | **verified** | — | D-001 §2 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-16 | 方案冻结：deduct_frozen + 幂等修复 + 演进登记 | accepted | 01-decision/D-001-deduct-frozen-plan.md |
| D-002 | 2026-08-16 | S4 go 影响判定（加法路由不触发失效，不暂挂） | accepted | 01-decision/D-002-s4-go-judgment.md |