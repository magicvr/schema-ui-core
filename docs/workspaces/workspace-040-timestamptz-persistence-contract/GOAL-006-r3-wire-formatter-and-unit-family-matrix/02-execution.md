---
id: GOAL-006-r3-wire-formatter-and-unit-family-matrix
doc: execution
status: active
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.2.0
---

# 执行记录 · GOAL-006-r3-wire-formatter-and-unit-family-matrix

## 执行索引

| E-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|--------|------|
| E-001 | 2026-09-20 | R3 首个阶段子目标立项（用户确认 slug；`I-041-007` 载体裁决） | recorded | `02-execution/E-001-goal-created.md` |
| E-002 | 2026-09-20 | R3-A 第一段：Go 侧单一 fixed-6 formatter 落码与全量替换（含 parser 统一与 fixture 同步） | recorded | `02-execution/E-002-r3a-go-wire-sweep.md` |

## 事实边界

> R3 边界由 Root `D-018` 冻结（R3-A formatter+fixture、R3-B 单位族/时区矩阵、R3-C 跨版本矩阵+备份核对、R3-D 退出矩阵+用户确认）。本目标承担 **R3-A/B**。**已完成**：Go 侧单一 fixed-6 formatter + 全量输出替换 + parser 统一 + Go fixture 同步（E-002），全仓 `go test -count=1 ./...` 64/64 包 ok。**未完成**：Web fixture 同步与 `I-041-008` 判定 → 检查点 A 未完成；R3-B（单位族矩阵 + VP-020 时区 round-trip）未开始。