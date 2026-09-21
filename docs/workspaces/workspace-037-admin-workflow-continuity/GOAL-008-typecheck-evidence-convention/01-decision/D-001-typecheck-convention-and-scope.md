---
id: D-001-typecheck-convention-and-scope
doc: decision-entry
status: accepted
goal_id: GOAL-008-typecheck-evidence-convention
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# D-001 · 承接 F-005：以 `tsc -b` 为唯一类型检查口径

## 决定

1. **承接来源**：本目标承接 R6 `GOAL-007` 审计 `A-002 F-005`（high required）的**跨工作区部分**。用户于 2026-09-18 按 P-004 裁决处置路径为**方案 A：立独立目标系统性处置**（未选「只改约定不追溯」与「accepted-residual」）。
2. **唯一正确口径**：本仓 Web 类型检查一律使用 `tsc -b`（等价入口 `npm run build` 的前半段）。裸 `tsc --noEmit`（不带 `-b` / `-p`）在 `apps/web` **不构成类型检查**，不得作为验证证据。
3. **范围**：只处理本工作区 `workspace-037-admin-workflow-continuity` 的 canonical 范围。其他工作区的同类历史条目（已登记于 `E-001`）**只登记、不代改**，符合 AGENTS §6c「禁止跨区写入」。
4. **不重开 R6**：R6 的视觉范围与 C1～C8 交付事实不受本目标影响；R6 保持 `active · 7/8`，其 C6 的 F-005 处置已移交本目标。

## 理由

`apps/web/tsconfig.json` 自仓库脚手架提交 `a3e1e5ad` 起即为 solution-style（`{"files": [], "references": [...]}`），本身不含源文件。TypeScript 在非 `-b` 模式下按该文件编译时没有任何输入，因此恒 exit 0。这不是配置错误（solution-style 是 project references 的正当形态，`build` 脚本用 `tsc -b` 正是为此），问题在于**验证命令选错**：把「无输入」误读为「无错误」。

该缺陷的性质是**证据有效性**，不是本次 UI 变更的正确性。因此它不适合作为 R6 的阻断项留在 R6 内处理（R6 的视觉交付已达标），也不适合降级为残余（会让同类失实证据继续扩散）。立独立目标可以在不污染 R6 分母的前提下系统性处置，并与 R6 的 `active` 状态解耦。

## 未选方案

- **不记为 accepted-residual**：用户已明确选择方案 A；且残余会让历史失实证据保持无标注状态，后续审计需重复解释。
- **不在 R6 内直接改跨区台账**：AGENTS §6c 明确禁止跨区写入；R6 的 canonical 范围只覆盖本工作区。
- **不修改 `tsconfig.json` 结构**：solution-style + project references 是 Vite/TS 模板的正当形态，`tsc -b` 正常工作；改结构会牵动 `tsconfig.app.json`/`tsconfig.node.json` 与构建链，超出本目标必要范围。
- **不只加文档说明而不提供脚本入口**：仅靠文档容易被绕过；提供 `npm run typecheck` 使正确命令成为最省力的默认选择。

## 门禁

C3 的防复发守卫形态（`I-008-003`）尚未决策，需在 C3 前确定；若形态选择存在实质歧义，按 P-004 询问用户。本目标不关闭 R5-I-004、R6、Root 或 VP-037。
