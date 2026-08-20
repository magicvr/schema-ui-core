---
id: GOAL-002-r1-tx-port-and-config
doc: decision
status: done
parent: GOAL-001-store-dialects
created: 2026-08-20
updated: 2026-08-20
version: 0.1.0
---

# 决策记录 · GOAL-002

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 现行打开路径与 WithTx 形状 | S1 冻结 | S1 前 | 读 store/config | verified | 2026-08-20 | E-001 |
| I-002 | required | 本目标审计模式 | S1 关门 | S1 前 | 风险分级 | verified | 2026-08-20 | D-001：self |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-20 | 冻结内核 Tx 端口与 db 配置键 | accepted | [D-001-tx-port-and-config-freeze.md](01-decision/D-001-tx-port-and-config-freeze.md) |
