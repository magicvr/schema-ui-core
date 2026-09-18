---
id: D-013-open-typecheck-evidence-goal
doc: decision-entry
status: accepted
goal_id: GOAL-001-admin-workflow-continuity
created: 2026-09-18
updated: 2026-09-18
parent: null
version: 1.0.0
---

# D-013 · 用户裁决 F-005 处置路径并开设整改子目标 GOAL-008

## 决定

用户在 R6 C6 审计 `A-002` 后按 P-004 裁决两项：

1. **F-005 处置路径 = 方案 A（立独立目标系统性处置）**：未选「只改约定不追溯」（B）与「accepted-residual」（C）。
2. **R6 关门时机 = 等 F-005 裁决落定后再关门**：R6 保持 `active · 7/8`，不在本轮投影 `done`。

据此开设整改子目标 `GOAL-008-typecheck-evidence-convention`（`active · 2/4`），承接 F-005 的跨工作区部分。它是 **Root 下的整改子目标，不是纲领阶段**，因此不改变 Root 六阶段分母与 `progress: 4/6`。

## 范围

- 确立 `tsc -b` 为唯一类型检查口径，提供 `npm run typecheck` 可发现入口，并确定防复发守卫形态。
- 只处理 `workspace-037-admin-workflow-continuity` 的 canonical 范围；其他工作区的同类历史条目只登记不代改（AGENTS §6c 禁止跨区写入）。
- 不重开 R6 的视觉范围，不修改 `tsconfig` 的 project references 结构。

## 理由

F-005 的性质是**证据有效性**缺陷而非本次 UI 变更的正确性缺陷：R6 的 C5/C7/C8 交付经独立复核属实且 `tsc -b` 全绿，但「类型检查通过」这一结论在流程上曾长期不成立。用户选择方案 A 的理由是系统性处置可一次性厘清跨阶段证据有效性，避免同类失实证据继续扩散。

将处置移出 R6 的原因：R6 的 canonical 范围不覆盖其他工作区，且把跨区流程问题作为 R6 阻断项会让 R6 的视觉交付结论长期悬置。独立子目标可与 R6 状态解耦。

## 未选方案

- 不在 R6 内直接改跨区台账：违反 AGENTS §6c。
- 不把 F-005 记为 accepted-residual：用户已明确选 A；残余会让历史失实证据保持无标注状态。
- 不修改 `tsconfig.json` 结构：solution-style + project references 是正当形态，`tsc -b` 工作正常。

## 门禁

`GOAL-008` 的 C3 守卫形态（`I-008-003`）若存在实质歧义，按 P-004 询问用户。本目标不关闭 R5-I-004、R6、Root 或 VP-037；R6 的完成投影待 `GOAL-008` 处置落定后单独执行。
