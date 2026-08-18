---
id: D-003
doc: decision-entry
goal: GOAL-001-admin-functional-modules
status: accepted
created: 2026-08-18
updated: 2026-08-18
version: 1.0.0
---

# D-003 · VP-011 有界关门与四档地图上提

## 背景

用户确认：剩余能力逻辑上都应由未来 VP/工作区承载，四档能力地图属于组合层路线图，不应继续留在 VP-011/workspace-011 的交付范围内。标准 Admin 功能模块波次已完成。

## 决策

1. 四档能力地图上提至 `docs/vision/roadmap.md` 作为组合层后续方向权威登记；workspace-011 的 `I-011-002` 保留为历史证据。
2. VP-011 有界关门：`status: active → closed`；关门记录列出 residual（S-05/S-06/S-07/S-08/S-13 与 B-01～B-11 reclassify 到未来 VP/工作区）。
3. workspace-011 Root `GOAL-001-admin-functional-modules` 置 `done`；子目标已全部 done。
4. 当前 active 交付 VP 调整为 VP-012；VP-009/VP-010 继续为持续程序。
5. 不修改 Charter，不改变 `vision_id@version`。

## 审计模式

文档收口与状态变更（低风险、可逆，但涉及 VP/Root 状态）：**self** 已通过 VRev-027 与本次记录；未触发 security/data/migration 门禁。

## 未选方案

- 保持 VP-011 active 作为 backlog 容器：会让“已完成交付”与“未来路线图”混淆。
- 不 reclassify S-05～S-08/S-13 直接关门：会造成“常用档未交付”的模糊 residual，关门语义不干净。
