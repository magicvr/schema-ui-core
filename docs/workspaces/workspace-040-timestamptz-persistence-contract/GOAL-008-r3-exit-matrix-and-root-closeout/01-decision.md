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
| I-041-010 | required | Root 关门的用户确认（判据 6） | C | **open**（2026-09-21 用户未选任一选项，改为要求先修 `dev.cmd start`；缺陷已 `fixed`，等待重新裁决） | **必须由用户 P-004 裁决**；不得静默代替。参见 `03-audit/A-004` §4（含新增的 `.env` dev 库指向问题） |
| I-041-011 | non-blocking | 产品级支持组合是否写入发布说明/退出矩阵 | A | **verified（用户 2026-09-21 裁决）** | 用户选择「只保留在退出矩阵与 R3-C 附件」，**不**写入产品发布说明 |

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
