---
id: GOAL-008-r3-exit-matrix-and-root-closeout
doc: decision
status: active
parent: null
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# 决策记录 · GOAL-008-r3-exit-matrix-and-root-closeout（R3-D）

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 | 影响门禁 | 状态 | 证据 / 决策 |
|----|------|----------|----------|------|-------------|
| I-041-010 | required | Root 关门的用户确认（判据 6） | C | open | **必须由用户 P-004 裁决**；不得静默代替 |
| I-041-011 | non-blocking | 产品级支持组合是否写入发布说明/退出矩阵 | A | open | R3-C `A-002` 移交的建议 P-004 点 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|--------|------|
| — | — | 尚无本目标决策（范围由 Root `D-018` §2 第 5 项冻结） | — | — |

> 新决策从 `01-decision/D-NNN-<slug>.md` 写入；编号在本目标内单调不复用。

## 已继承的边界事实（不得在本目标重开）

- **证据定位**：Docker 容器矩阵只作 CI/reproducibility（`GOAL-007/D-002` §4；`D-017` §3 约束②）。
- **residual 台账**：R1 `D-021` 的 `F-I-005` 已按 `fixed` 闭合（`GOAL-002/03-audit/A-048`），**不得**在退出矩阵里重新记成 open residual。
- **不重开**：R1/R2/R3-A/B/R3-C 的冻结决策与已关门结论；本目标只汇总证据与关门。
- **判据 6 是硬门禁**：Root `status: done` 需用户书面确认（本节 C 检查点）。
