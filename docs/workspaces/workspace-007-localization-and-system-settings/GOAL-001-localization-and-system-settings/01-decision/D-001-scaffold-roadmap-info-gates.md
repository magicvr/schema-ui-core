---
id: D-001
doc: decision
title: S0 · scaffold 工作区与 Root、纲领路线图 S0–S5、信息门禁登记
status: accepted
parent: GOAL-001-localization-and-system-settings
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
---

# D-001 · scaffold 工作区与 Root、纲领路线图、信息门禁登记

## 决定

1. **建区**：按 `/vision` 轮次用户确认（2026-08-09），创建 `docs/workspaces/workspace-007-localization-and-system-settings/`（`vision_role: delivery`，`plan_refs`/`primary_plan` = `VP-007-localization-and-system-settings`，`shared_materials_catalog: none`）与 Root `GOAL-001-localization-and-system-settings`（`parent: null`，`status: active`）。
2. **纲领路线图（P-001）**：按 VP-007「建议实现阶段」建立 S0 → S1 → S2 → S3 → S4 → S5 六个串行纲领阶段；同一阶段内可并行子目标。`progress: 0/6` 由 6 个阶段检查点等权派生。
3. **信息门禁（P-005）**：`I-L10N-001`～`005` 全部登记为 `required` / `open`，与 VP-007 同 id 对齐；S0 契约冻结前必须全部关闭（否则阻断 S0→S1 放行）。尚未到期的 S1/S3/S4 门禁不阻断本轮 scaffold。
4. **审计模式**：本轮 scaffold 为低风险、可逆、无门禁语义变化的治理文档操作 → 模式 `none`；S0 契约冻结、S1 实施起点与关门时按 P-002 另行审视（self 或 independent）。

## 未选方案

- 不把本波次吸收进 closed workspace-003/004/005/006（VP-007 Non-goals 禁止）。
- 不按信息门禁机械创建两个信息收集子目标（P-005 §按规模拆分）：S0 作为纲领阶段在 Root 内闭环收集与冻结。

## 影响

- S0 阶段子目标创建、方案冻结、放行均受 `I-L10N-001`～`005` 门禁约束；任何 residual 须用户书面接受（P-004）。
