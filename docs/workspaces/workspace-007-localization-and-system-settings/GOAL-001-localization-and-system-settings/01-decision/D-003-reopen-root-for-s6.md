---
id: D-003
doc: decision
title: 暂时回退 Root 关门状态承接 S6 子目标（GOAL-007）
status: accepted
parent: GOAL-001-localization-and-system-settings
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
---

# D-003 · 暂时回退 Root 关门状态承接 S6（2026-08-09）

## 决策

用户 2026-08-09 书面指令：设置页表单/详情页改造（方案 A）在工作区 7 以子目标承接治理上下文，Root 暂时回退关门状态。执行：

1. **新增子目标** `GOAL-007-s6-settings-form-page`（parent = GOAL-001），承载 S6 阶段（recordSource 预填 + settings 页四类内联表单重构），见其 `00-meta.md` / `01-decision/D-001`。
2. **Root 状态回退**：`status: done` → `active`；`progress: 6/6` → `6/7`（7 个纲领阶段检查点等权派生，S6 开放；仅展示）。
3. **历史不重写**：S0–S5 关门记录（GOAL-006 D-002）保留原文；本记录只说明重新开根原因与范围，不推翻已关门事实。
4. **暂时语义**：GOAL-007 C4（证据 + 关门审计 + 用户书面确认）完成后，Root 恢复 `done`、`progress: 7/7`，并同步 goal-tree / workspace.md / 索引 frontmatter。
5. **Vision 边界**：VP-007 保持 `closed`——S6 为已关闭波次上的增量产品化（非新愿景波次、不 re-align、不扩张 Profile 可见性）。

## 为什么

- 新增子目标使 Root `done` 不再成立（Root 状态须反映未关门子目标的存在）。
- 用户指令明确「在工作区 7 添加子目标承载此项治理的上下文，并暂时回退根目标的关门状态」——在已关门波次上延续产品化，属用户裁决的合法 re-open（P-004 留痕）。
- 回退仅限 Root 状态与进度展示，不触碰已关闭阶段审计结论与历史证据。

## 未选方案

| 方案 | 未选原因 |
|------|----------|
| 新工作区 / 新 VP 承接 | 用户明确指定工作区 7 子目标；改动为已关闭波次上的增量，非新愿景波次 |
| Root 保持 done、S6 游离于树外 | 违反 goal-tree 一致性（Root 须反映全部未关门子目标） |
