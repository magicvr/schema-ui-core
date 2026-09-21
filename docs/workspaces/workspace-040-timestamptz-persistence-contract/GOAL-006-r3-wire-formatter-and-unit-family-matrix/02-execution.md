---
id: GOAL-006-r3-wire-formatter-and-unit-family-matrix
doc: execution
status: done
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-21
version: 0.4.0
---

# 执行记录 · GOAL-006-r3-wire-formatter-and-unit-family-matrix

## 执行索引

| E-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|--------|------|
| E-001 | 2026-09-20 | R3 首个阶段子目标立项（用户确认 slug；`I-041-007` 载体裁决） | recorded | `02-execution/E-001-goal-created.md` |
| E-002 | 2026-09-20 | R3-A 第一段：Go 侧单一 fixed-6 formatter 落码与全量替换（含 parser 统一与 fixture 同步） | recorded | `02-execution/E-002-r3a-go-wire-sweep.md` |
| E-003 | 2026-09-20 | R3-A 第二段：Web fixture 同步 + 三个漏网公共 wire 时间字段（healthz/readyz、mail outbox、mail config）修复 | recorded | `02-execution/E-003-r3a-web-fixtures-and-wire-field-gap.md` |
| E-004 | 2026-09-20 | R3-B：单位族矩阵（秒/毫秒/可空/sentinel 四端点）+ VP-020 会话时区展示 round-trip | recorded | `02-execution/E-004-r3b-unit-family-matrix-and-tz-roundtrip.md` |
| E-005 | 2026-09-20 | 响应 independent A-002（`fail`）：3 条 required 修复 + 两条常驻 wire 守卫 + 族身份 provenance 检查 | recorded | `02-execution/E-005-a002-required-fixes.md` |

## 事实边界

> R3 边界由 Root `D-018` 冻结（R3-A formatter+fixture、R3-B 单位族/时区矩阵、R3-C 跨版本矩阵+备份核对、R3-D 退出矩阵+用户确认）。本目标承担 **R3-A/B**，**已关门**（2026-09-21，`done · 3/3`）。
>
> **完成**：Go 侧单一 fixed-6 formatter + 全量输出替换 + parser 统一 + Go fixture 同步（E-002）；Web fixture 同步、三处漏网字段修复与 `I-041-008` 用户裁决（E-003 + Root `D-019`）；单位族四端点矩阵与 VP-020 会话时区 round-trip（E-004），`I-040-004` 关闭；independent `A-002` 判 `fail` 后的 3 条 required 修复与常驻守卫（E-005）；`A-004` 定向复审 **pass / 开放 required = 0**；recommended `F-I-101`（PUT 回归改为冻结时钟的确定性断言）与 `F-I-102`（守卫扩到 B/D/E/F 四类并收回「闭集」措辞，A/C 两类书面记录）在 `A-005` 中列为 fixed。
>
> **审计事实**：`A-002`（independent · grok 4.6 · high）曾**不同意**「检查点 A 已完成」并给出 3 条 required；`A-004` 复核后确认全部 `fixed`、同意 A 成立且可进入检查点 C。关门依据 `A-005`。
>
> **移交**：R3-C（`GOAL-007-r3-pg-cross-version-restore-matrix`，`I-041-004` 仍 open）与 R3-D（`GOAL-008-r3-exit-matrix-and-root-closeout`，含用户确认关门）尚未立项；Root 仍 `active · 2/3`。
