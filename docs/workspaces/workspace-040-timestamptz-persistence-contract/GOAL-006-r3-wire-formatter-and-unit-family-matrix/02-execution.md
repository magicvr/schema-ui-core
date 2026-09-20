---
id: GOAL-006-r3-wire-formatter-and-unit-family-matrix
doc: execution
status: active
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.3.0
---

# 执行记录 · GOAL-006-r3-wire-formatter-and-unit-family-matrix

## 执行索引

| E-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|--------|------|
| E-001 | 2026-09-20 | R3 首个阶段子目标立项（用户确认 slug；`I-041-007` 载体裁决） | recorded | `02-execution/E-001-goal-created.md` |
| E-002 | 2026-09-20 | R3-A 第一段：Go 侧单一 fixed-6 formatter 落码与全量替换（含 parser 统一与 fixture 同步） | recorded | `02-execution/E-002-r3a-go-wire-sweep.md` |
| E-003 | 2026-09-20 | R3-A 第二段：Web fixture 同步 + 三个漏网公共 wire 时间字段（healthz/readyz、mail outbox、mail config）修复 | recorded | `02-execution/E-003-r3a-web-fixtures-and-wire-field-gap.md` |
| E-004 | 2026-09-20 | R3-B：单位族矩阵（秒/毫秒/可空/sentinel 四端点）+ VP-020 会话时区展示 round-trip | recorded | `02-execution/E-004-r3b-unit-family-matrix-and-tz-roundtrip.md` |

## 事实边界

> R3 边界由 Root `D-018` 冻结（R3-A formatter+fixture、R3-B 单位族/时区矩阵、R3-C 跨版本矩阵+备份核对、R3-D 退出矩阵+用户确认）。本目标承担 **R3-A/B**。**已完成**：Go 侧单一 fixed-6 formatter + 全量输出替换 + parser 统一 + Go fixture 同步（E-002）；Web fixture 同步与三处漏网字段修复（E-003）；`I-041-008` 经用户裁决无破坏性/不需兼容期（Root `D-019`）；单位族四端点矩阵与 VP-020 会话时区展示 round-trip（E-004），`I-040-004` 关闭。**未完成**：检查点 C（self + grok independent 审计、required 合法闭合、静默关门）。`progress: 2/3`。
