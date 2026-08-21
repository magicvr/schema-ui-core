---
id: E-003-r3-canonical-layout-correction
doc: execution
goal: GOAL-004-r3-bounded-pilot
date: 2026-08-05
status: recorded
---

# E-003 · 修正 GOAL-004 canonical 布局

## 事实

发现仓库根存在重复的 `GOAL-004-r3-bounded-pilot/`，而 workspace-003
canonical 目录同时已有该目标的 `00-meta.md`、决策索引和部分记录。保留
canonical 索引与状态真相，将根目录版本中缺失的 `D-003`、`E-001/E-002`、
`03-audit.md`、`A-001` 和两个证据附件移入 canonical；两份重复决策中的
唯一边界约束合并进 canonical `D-001/D-002`，没有覆盖其较完整索引。

根目录重复文件与空目录随后移除。GOAL-004 的 `status` 和 `progress` 未因
布局修正改变，C1、I-006 和 R3 门闩仍按原台账保持未完成状态。

## 校验

- canonical 目标现具备 `00-meta.md`、三个稳定索引、三个 ledger 目录和
  `attachments/`。
- `D-003`、`E-001/E-002`、`A-001` 以及附件链接均落在 canonical 相对路径。
- 仓库根不再存在第二个 `GOAL-004-r3-bounded-pilot/` 目录。
