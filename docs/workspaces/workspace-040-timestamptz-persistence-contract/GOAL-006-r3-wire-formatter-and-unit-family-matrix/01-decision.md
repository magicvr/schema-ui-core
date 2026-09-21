---
id: GOAL-006-r3-wire-formatter-and-unit-family-matrix
doc: decision
status: active
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.2.0
---

# 决策记录 · GOAL-006-r3-wire-formatter-and-unit-family-matrix（R3-A/B）

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 | 影响门禁 | 状态 | 证据 / 决策 |
|----|------|----------|----------|------|-------------|
| I-041-008 | required | wire 输出的破坏性/兼容期 | A | **verified（用户 2026-09-20 P-004）** | Root `D-019` §1：仓内无破坏性依赖、不需兼容期；证据表 6 条 |
| I-040-004 | required（继承） | VP-020 时区展示矩阵 | B | **verified（本目标 B）** | 载体 `I-041-007`；证据 `02-execution/E-004` §2–3 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|--------|------|
| — | — | 本目标无自有决策：两个门禁项分别由 Root `D-019`（`I-041-008`）与执行证据（`I-040-004`）收口 | — | — |

> 新决策从 `01-decision/D-NNN-<slug>.md` 写入；编号在本目标内单调不复用。

## 边界事实（供审计取证）

- R3-A 的实施范围由 `D-003`/`D-005`/`D-009` 与 inventory 冻结；E-003 §2 记录的「三个漏网公共 wire 时间字段」是**按类别补齐 `D-003` 合同**，不是新决策：未改字段语义、JSON 键名或 NULL 语义（`mail.PublicView` 保持 `*time.Time` + JSON null，遵循 `GOAL-004/A-003` 的 `user-overruled` 裁决）。
- R3-B 只做回归矩阵，不新增 VP-020 能力（`D-018` §3）。
