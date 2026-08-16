---
id: GOAL-020-wallet-auto-account
doc: decision
status: active
parent: GOAL-001-admin-functional-modules
created: 2026-08-16
updated: 2026-08-16
version: 0.4.0
---

# 决策记录 · GOAL-020-wallet-auto-account

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 自动开户触发面与幂等 | 方案 | S1 | — | **verified** | — | D-001 §1 |
| I-002 | required | 手动创建边界与错误码 | 方案 | S1 | — | **verified** | — | D-001 §2 |
| I-003 | non-blocking | 前端选项调整副作用 | 方案 | S1 | — | **verified** | — | D-001 §3 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-16 | 方案冻结：自动开户 get-or-create + 手动边界 | accepted | 01-decision/D-001-auto-account-plan.md |
| D-002 | 2026-08-16 | S4 go 影响判定（内容扩展不触发失效，不暂挂） | accepted | 01-decision/D-002-s4-go-judgment.md |