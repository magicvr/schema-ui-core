---
id: D-004
doc: decision
title: S6 关门，Root 恢复 done（解除临时回退）
status: accepted
parent: GOAL-001-localization-and-system-settings
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
---

# D-004 · S6 关门，Root 恢复 done（2026-08-09）

## 决策

GOAL-007（S6，设置页表单/详情页改造）C4 用户书面确认关门（GOAL-007 `D-002`）后，Root 恢复：

- `status: active` → **`done`**；`progress: 7/7`（S0–S5 + S6 七个纲领阶段检查点全部完成）。
- 解除 `D-003` 的**暂时回退关门**语义；S0–S6 全部完成，无未关门子目标。
- 同步 goal-tree / workspace.md / 索引 frontmatter（GOAL-001 `0[123]-*.md` → `done`）。

## 为什么

- D-003 §4「暂时语义」：GOAL-007 C4（证据 + 关门审计 + 用户书面确认）完成后 Root 恢复 `done`、`7/7`。
- S6 关门证据：E-001/E-002 + A-001（self pass）+ A-002（independent pass，required 0）+ 用户书面确认（GOAL-007 D-002）。

## 未选方案

| 方案 | 未选原因 |
|------|----------|
| Root 保持 active 到下一波次 | 用户已书面确认 S6 关门；无继续开放理由，恢复 `done` 与 S0–S5 历史一致 |
