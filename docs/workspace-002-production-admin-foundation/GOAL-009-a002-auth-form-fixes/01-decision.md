---
title: 决策 · A-002 缺陷修复
status: active
created: 2026-08-03
updated: 2026-08-03
parent: GOAL-001-production-admin-foundation
version: 0.1.0
---

# 决策 · GOAL-009

## D-001 · 承接 Root A-002 两条 required 缺陷修复（F-002-002/003）

- **日期**：2026-08-03
- **状态**：accepted
- **用户裁决**（P-004）：A-002 三条 required 均走 `fixed` 关闭路径；F-002-002/003 由本目标承载；F-002-001 归 `GOAL-010-a002-schema-adapter`（Root D-014）；recommended F-002-004~006 为本目标可选加分；A-002 同 scope self 审计延后至修复完成后随关门补（P-004 §3.1）。
- **决定**：以本目标承载 F-002-002（表单校验错误阻断提交）与 F-002-003（认证失效状态清理）的修复与回归证据。
- **理由**：F-002-002/003 是确定性小缺陷，关闭路径由审计明确；独立成目标便于验收、回归与 Root finding 的关闭证据复核。
- **未选方案**：
  - **与 GOAL-010 合并单目标**：混合小修复与大改造，验收口径混杂，P-001 路线图难以对齐。
  - **不在 Goal 层跟踪**：破坏本仓按波次建子目标的惯例，Root finding 关闭缺少载体。
- **影响**：Root A-002 F-002-002/003 保持 `open`，直至本目标 S1/S2 实施完成并有关门证据后按 `fixed` 闭合。
- **后续**：S1 → S2 → S3 回归 → S4 关门审计；`/audit` finding-closure 复审建议在 S4 前请求。
