---
id: D-015-authorize-cross-workspace-typecheck-errata
doc: decision
status: accepted
goal_id: GOAL-001-admin-workflow-continuity
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# D-015 · 授权跨工作区类型检查空转条目追溯更正

## 决定

2026-09-18，用户就 `GOAL-008-typecheck-evidence-convention` 的 `I-008-004`（其他工作区历史条目是否追溯更正）作出书面裁决：**授权追溯更正**。

据此，`GOAL-008` 原先「跨区只登记、不代改」的边界（AGENTS §6c 禁止跨区写入）在本项上由用户明确授权放行，动作范围**限定为加勘误注记**：

1. 保留原命令行、原结论、原 verdict，只追加带日期的勘误注记，并注明「原记录保留不改」。
2. 跨区引用使用 Q2 路径（`docs/workspaces/workspace-037-admin-workflow-continuity/GOAL-008-typecheck-evidence-convention/`）。
3. 不改任何外部目标的 `status` / `progress` / verdict / finding 状态，不改其 `03-audit.md` 索引。
4. 追加在他人 `03-audit/A-*.md` 中的注记须写明「不改本条 verdict 与结论」，性质是事后勘误、不是第二条审计意见。

## 未选方案

- **保持 deferred**：用户已明确要求更正，不采纳；仅保留为备选（若后续认为跨区注记负担过重，可停止追加新注记，不影响已落盘的 11 处）。
- **改写原文命令与结论**：会破坏历史记录的可核对性与 P-003「审计意见不代改」原则，不采纳。
- **逐条考古全仓 `tsc` 简写（300+ 行）**：范围与收益不成比例，且无法从文本唯一确定当时的命令形态；不采纳，改为在 `A-002 F-002` 中登记为 recommended open。

## 边界与不变量

- 本决定不改变 Root 六阶段分母与 `progress`，不关闭 `R5-I-004`、Root 或 VP-037。
- 本决定不授权修改 `apps/web/src/typecheck-convention.guard.test.ts` 等已关门交付物——守卫缺口（`GOAL-008 A-002 F-001`）另行由用户裁决。
- 授权仅覆盖类型检查空转条目的勘误注记，不构成对 workspace-002/009/010/011 其他内容的修改许可。
