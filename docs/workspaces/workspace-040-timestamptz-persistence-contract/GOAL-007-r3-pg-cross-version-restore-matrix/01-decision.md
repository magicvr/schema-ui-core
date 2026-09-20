---
id: GOAL-007-r3-pg-cross-version-restore-matrix
doc: decision
status: active
parent: null
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# 决策记录 · GOAL-007-r3-pg-cross-version-restore-matrix（R3-C）

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 | 影响门禁 | 状态 | 证据 / 决策 |
|----|------|----------|----------|------|-------------|
| I-041-004 | non-blocking（继承） | PG 15/16/17 跨版本 `pg_restore` 兼容矩阵 | B | open | `D-018` §5；关闭或 residual 见本表更新 |
| I-041-009 | required（本目标新增） | 「supported / unsupported」判定口径与组合边界 | A | **verified（2026-09-21）** | `D-001` §2–§3（先于矩阵本体冻结）；前置实测 `attachments/r3c-pg-tool-compatibility-probe-v0.1.md` |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|--------|------|
| D-001 | 2026-09-21 | R3-C 矩阵定义与 supported/unsupported 判定口径（检查点 A；关闭 `I-041-009`） | accepted | `01-decision/D-001-r3c-matrix-definition-and-criterion.md` |

> 新决策从 `01-decision/D-NNN-<slug>.md` 写入；编号在本目标内单调不复用。

## 已继承的边界事实（不得在本目标重开）

- **证据定位**：Docker 临时容器的跨版本结果只作 **CI/reproducibility**，不是生产就绪证据（`D-017` §3 约束②，用户裁决）。
- **破坏性动作范围**：只可作用于一次性/专用测试 database（`D-017`，用户裁决）。
- **不重开**：R1/R2 的 canonical SQL/checksum（v1–v87 不可变）、Store codec、C3 Port 与三类产物区分、R3-A/B 的 wire 合同。
- **Root 关门（判据 6）**仍须用户确认，属 R3-D。
