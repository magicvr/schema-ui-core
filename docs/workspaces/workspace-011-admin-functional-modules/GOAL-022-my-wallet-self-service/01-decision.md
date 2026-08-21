---
id: GOAL-022-my-wallet-self-service
doc: decision
status: active
parent: GOAL-001-admin-functional-modules
created: 2026-08-16
updated: 2026-08-16
version: 0.2.0
---

# 决策记录 · GOAL-022

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 只读 vs 有界自助操作 | S1 | S1 | 用户裁决 | verified | — | 用户 2026-08-16 裁决：只读自服务（D-002 §1） |
| I-002 | required | 路由与 get-or-create 时机 | S1 | S1 | 对照 GOAL-020 | verified | — | 用户 2026-08-16 裁决：/my-wallet + 惰性开户（D-002 §2） |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-16 | 开项目录：承接 W12 T-04 移交 | accepted | `01-decision/D-001-open-from-w12.md` |
| D-002 | 2026-08-16 | 方案冻结：只读自服务 + /my-wallet 惰性开户 | accepted | `01-decision/D-002-s1-plan-freeze.md` |
| D-003 | 2026-08-16 | S4 go 判定（只读加法面，无门禁语义变化） | accepted | `01-decision/D-003-s4-go-judgment.md` |