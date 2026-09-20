---
id: GOAL-006-r3-wire-formatter-and-unit-family-matrix
title: R3 · 公共 wire fixed-6 formatter 与单位族/时区矩阵（R3-A/B）
status: active
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
progress: 0/3
plan_refs:
  - VP-040-timestamptz-persistence-contract
primary_plan: VP-040-timestamptz-persistence-contract
serves_summary: 承接 Root R3 的 R3-A/B：落码共享 fixed-6 UTC wire formatter 并替换全部 inline 布局、同步 Go/Web fixture、补齐输入兼容矩阵（0/3/6/9 位与 +00:00、拒绝无时区），并以单位族矩阵（秒/毫秒/可空/sentinel 各至少一个 endpoint）与 VP-020 会话时区展示 round-trip 收口 I-040-004。
---

# GOAL-006 · R3 · 公共 wire fixed-6 formatter 与单位族/时区矩阵

## 概述

承接 `GOAL-001-timestamptz-persistence-contract` 的 **R3 阶段 R3-A/B 检查点**（Root `D-018`）。R1（`GOAL-002`）与 R2（`GOAL-003`/`004`/`005`）均已关门；本目标只做**公共 wire 面与回归矩阵**，不改 Store/迁移/描述符。

**子目标编号与 slug 经用户 2026-09-20 确认**（AGENTS §11）：`GOAL-006-r3-wire-formatter-and-unit-family-matrix`。

**`I-041-007` 载体裁决（用户 2026-09-20）**：Go 单测锁 wire 形状与 parser 兼容；Web 用组件/单测锁会话时区展示 round-trip；**不**引入浏览器 e2e。

## 范围

| # | 交付 | 依据 |
|--:|------|------|
| 1 | 一个 shared fixed-6 UTC formatter（`YYYY-MM-DDTHH:MM:SS.ffffffZ`），替换 `rfc3339Milli` 与 inventory 列出的全部 inline 布局 | `D-003`；`D-009`；inventory §C2/R3 要求 1 |
| 2 | 输入 parser 兼容矩阵：0/3/6/9 位小数与 `+00:00` 等等价 offset → UTC；**拒绝无时区/模糊本地** | `D-005`；要求 2 |
| 3 | Go 与 Web fixture 同步（inventory 逐项），保留非时间字段的 wire 形状 | 要求 3 |
| 4 | **单位族矩阵**：秒 / 毫秒 / 可空 / sentinel 转换 各至少一个 endpoint 的 wire 输出断言 | 要求 4 |
| 5 | **VP-020 会话时区展示 round-trip**：会话/用户时区设置下的展示与 UTC 存储一致（Web 组件/单测载体） | `I-040-004`；`I-041-007` |
| 6 | `I-041-008`（wire 输出的破坏性/兼容期）证据收集与判定 | `D-018` §5 |

## 非目标（本子目标**不**做）

- **不**改 Store 物理类型、codec、迁移描述符、canonical SQL/checksum（v1–v87 不可变）。
- **不**改 C3 的 Port/产物语义。
- **不**做 R3-C（PG 跨版本矩阵 + 备份有界核对）与 R3-D（退出矩阵/关门）——待本目标完成后再立项。
- **不**新增浏览器 e2e 依赖（用户 `I-041-007` 裁决）。
- **不**引入 ORM/第三库/Redis/MQ/多实例；**不**把驱动类型泄漏进公共契约。
- **不**扩展 VP-020 的时区能力本身（只做回归矩阵）。

## 红线

- 不把「非 DB 文本」（邮件正文、audit `detail`、任意 payload 字段）纳入固定 6 位输出（`D-009`）。
- 不以 `pgtype`/驱动时间类型作为 wire 值。
- 不以「测试通过」代替独立审计结论；required finding 未合法闭合不得关门。

## 成功标准（本子目标检查点，用于 `progress` 派生）

| 检查点 | 判据 | 状态 |
|--------|------|------|
| **A** | shared fixed-6 formatter 落码并替换 inventory 的 Go 面；输入兼容矩阵（0/3/6/9 位、`+00:00`、拒绝无时区）有可执行测试；Go/Web fixture 同步；`I-041-008` 判定 | pending |
| **B** | 单位族矩阵（秒/毫秒/可空/sentinel 各至少一个 endpoint）通过；VP-020 会话时区展示 round-trip 通过（Go + Web 组件/单测）；`I-040-004` 关闭 | pending |
| **C** | self + grok independent 审计落盘、required 合法闭合 → 静默关门 | pending |

`progress: 0/3` 由 A～C 等权派生；**不**放行阶段、**不**关闭 finding、**不**推导 `done`。

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 | 影响门禁 | 状态 | 证据 |
|----|------|----------|----------|------|------|
| I-040-004 | required（继承） | VP-020 展示/输入与 UTC 存储的回归矩阵 | B | open（本目标 B 关闭） | `D-018` §5；载体裁决 `I-041-007` |
| I-041-007 | required（继承） | VP-020 矩阵的载体与验收口径 | B | **verified（用户 2026-09-20）** | `D-018` §5 |
| I-041-008 | required（本目标新增） | wire 输出的破坏性：前端/现网是否依赖 3 位小数；是否需兼容期 | A | **collecting** | `D-018` §5；inventory §Web consumer |

## 父目标

- `GOAL-001-timestamptz-persistence-contract`（R3 检查点）

## 关门条件

**A～C 全部完成**、`03-audit` 的 self 与 independent 意见落盘、required 合法闭合后才可静默关门；关门**不**等于 Root 关门（Root 关门仍须 R3-C/D 与用户确认）。
