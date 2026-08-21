---
id: D-002-independent-audit-provider-grok-build
doc: decision-entry
goal: GOAL-001-admin-module-readiness
status: accepted
supersedes: D-001（仅 provider 记录；开区绑定/命名/Root/单工作区不变）
created: 2026-08-10
updated: 2026-08-10
version: 1.0.0
---

# D-002 · independent 交叉审计执行路径改为 grok build（grok 4.5 · 思考强度 high · `audit`）

## 用户指令

用户在目标级指令（2026-08-10）中明确指定本 Root 的独立交叉审计入口：

> 需要交叉审计时候，可以调用 grok build（模型 grok 4.5，思考强度 high）进行独立审计（它有跟你一样的 skills，可以直接调用 `audit` 命令）。

按直接用户指令优先原则（先于关门修正 [D-001](D-001-workspace-s0-bindings.md) 中记录的 provider），本决策 **更新 `I-READINESS-005` 的 independent provider**：

| 字段 | 旧值（D-001） | 新值（本决策） |
|------|----------------|----------------|
| provider | GitHub Copilot · `/audit` | **grok build · 模型 `grok-4.5` · 思考强度 high · 执行 `audit` 命令** |
| 执行方式 | 会话内指定 /audit | 与本次编排同 Skills 的独立 grok 4.5 会话，读取 `audit` skill（`.grok/skills/audit/SKILL.md`）并产出可核对的独立审计意见 |
| 审计模式 | `cross`（self + independent） | 保持 `cross`（self + independent）不变 |

## 范围与不变项

- 仅更新 `I-READINESS-005` 的 provider 与执行路径记录；**不改变**单工作区绑定、Root 命名、`GOAL-001-admin-module-readiness`、`parent: null`、审计模式 `cross`、审计覆盖 scope（compatibility / data / migration / production/release / 跨边界治理语义）。
- D-001 的历史记录保留原文（不得改写历史决策）；本决策以 `supersedes` 声明 provider 段被替换，`I-READINESS-005` 的证据与闭合以本决策为准。
- provider 选择本身不是已完成的审计证据：S5 必须由该 provider 的**独立会话**产出可核对的 Goal 审计意见（落盘 `03-audit/A-NNN-*.md`，`source: independent`）。independent 会话不可用或无可核对输出时，独立门禁保持未满足，不得由 self 或编排器冒充。

## 未选方案

- **维持 GitHub Copilot `/audit`**：与用户 2026-08-10 目标级直接指令冲突；直接用户指令优先。
- **不落盘、仅口头改**：P-003/P-004 要求书面留痕；provider 与审计路径记录必须版本化落盘。
