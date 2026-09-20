---
id: GOAL-002-r1-contract-and-denominator-freeze
doc: decision
status: active
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# 决策记录 · GOAL-002

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 状态 | 证据 / 决策 |
|----|------|-----------------|----------|------|-------------|
| I-040-001 | required | PG `timestamptz(6)` / SQLite fixed-6 RFC3339 TEXT 的精度、编解码、排序与 NULL 规则 | C2/R2 | collecting | Root D-002；待 C1/C2 |
| I-040-002 | required | 全部绝对时刻列与排除列分母 | C1/C2/R2 | collecting | Root D-002；待 inventory |
| I-040-003 | required | SQLite/PG 原地转换、失败恢复、备份依赖 | C3/R2/R3 | collecting | Root D-002；待设计 |
| I-040-004 | required | VP-020 展示/输入回归矩阵 | R3 | open | Root I-040-004；后续阶段 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-09-20 | R1 合同冻结（承接 Root 用户裁决） | accepted | `01-decision/D-001-r1-contract-freeze.md` |
| D-002 | 2026-09-20 | R2 migration 归属承接 | accepted | `01-decision/D-002-r2-migration-ownership.md` |
| D-003 | 2026-09-20 | C2 wire 输入兼容承接 | accepted | `01-decision/D-003-wire-input-compat.md` |
| D-004 | 2026-09-20 | C3 Backup SPI/Service 合同承接 | accepted | `01-decision/D-004-backup-spi-contract.md` |
| D-005 | 2026-09-20 | C2/C3 guardrails 草案 | proposed | `01-decision/D-005-c2-c3-guardrails-proposed.md` |
| D-006 | 2026-09-20 | 最小 Backup/RecoveryPoint Port 承接 | accepted | `01-decision/D-006-backup-port-surface.md` |

> 用户裁决原文与范围记录在 Root `D-002-r1-contract-freeze-user-decisions.md`；本子目标承接并将其转为可验证 C1～C4 交付物。
