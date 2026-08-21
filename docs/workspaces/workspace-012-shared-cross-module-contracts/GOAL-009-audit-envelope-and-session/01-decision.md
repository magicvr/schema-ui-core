---
id: GOAL-009-audit-envelope-and-session
doc: decision
status: done
parent: GOAL-001-shared-cross-module-contracts
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
---

# 决策记录 · GOAL-009

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 | 影响门禁 | 最晚需要阶段 | 状态 | 证据 |
|----|------|----------|----------|--------------|------|------|
| I-001 | required | session 语义 | S0 | S0 | verified | D-001：refresh token id via JWT sid；机器凭据 = credential id；effective actor = actor |
| I-002 | required | writer 范围 | S1 | S0 | verified | D-001：全部生产 mutation 写路径改 NewDetail |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-19 | session 与 envelope 范围 | accepted | [D-001-session-and-envelope.md](01-decision/D-001-session-and-envelope.md) |
