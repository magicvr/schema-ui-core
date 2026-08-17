---
id: independent-audit-execution
doc: project-decision
title: 项目级决策 · 独立审计执行路径（grok build）
status: active
created: 2026-08-17
updated: 2026-08-17
parent: null
version: 1.0.0
---

# 独立审计执行路径（项目级决策）

> **项目级决策，跨工作区有效**：不绑定任何 workspace；新开工作区后本决策继续适用，无需重复落盘。

## 用户指令

用户项目级指令（2026-08-17）：

> 本项目在推进目标过程中，需要交叉审计的时候可以在自审计之后，直接调用本地的 grok build（模型 grok 4.6，思考强度 high）执行独立审计（grok build 可以使用本仓库所有的 skills，所以直接调用 `/audit` 命令即可），然后合并响应审计意见。

## 决策

1. **顺序**：推进目标过程中判定需要交叉审计（审计模式 `independent` / `cross`）时，**先自审（self）**；自审之后**直接调用本地 grok build** 执行独立审计；随后**合并响应**全部审计意见（P-003：编排器汇总 self + independent，驱动修正/复审/推进）。
2. **执行方式**：grok build（模型 **grok 4.6** · 思考强度 **high**）可使用本仓库全部 skills（`.grok/skills/` 含 `audit` / `govern` / `vision` / `vision-audit`），因此直接调用 **`/audit`** 命令即可（对应 `.grok/skills/audit/SKILL.md` 的独立交叉审计流程），无需手工搬运提示词。
3. **意见落盘**：grok 会话产出可核对的独立审计意见，按 P-003 落盘到被审目标 `03-audit/A-NNN-*.md`（`source: independent`，与 self 共用 A 序列）并更新 `03-audit.md` 索引；长证据可放 `attachments/`。写入方式：交叉工具直接追加，或代贴并保留 `source: independent`。
4. **合并响应**：编排器（/govern）汇总并响应与推进焦点相关的**全部**意见（self + independent）；存在未合法闭合的 required/必改 findings 时不得推进对应门禁或关门（三路径闭合：`fixed` / `accepted-residual` / `user-overruled`）；意见冲突按 P-004 展示并等用户裁决，未决不放行/不关门。
5. **适用范围**：本项目全部工作区与目标（含新开工作区）。单条目标需要其他 provider 时须另行记录用户指令或决策。

## 范围与不变项

- 本决策只确定 independent 的**默认 provider 与执行路径**（模型 / 思考强度 / 调用方式 / 顺序 / 合并响应）；**不改变** P-003 审计模式选择规则（按风险定 `none` / `self` / `independent` / `cross`）、意见落盘约定、P-004 用户裁决点、P-005 信息门禁。
- 模式为 `self` 时不需要 grok；`independent` / `cross` 已确定时按本路径执行。
- provider 选择本身不是审计证据：grok 会话必须产出可核对、可落盘的独立审计意见；grok 不可用或无可核对输出时，independent 门禁保持未满足，不得由 self 或编排器冒充。

## 关联与先例

- 延续 workspace-008 Root 决策 [D-002-independent-audit-provider-grok-build](../workspaces/workspace-008-admin-module-readiness/GOAL-001-admin-module-readiness/01-decision/D-002-independent-audit-provider-grok-build.md)（provider 记录 grok build，模型 grok-4.5 → 本决策升级 **grok-4.6**）。
- 既有实践：workspace-011 R2～R4 波次独立关门审计已全程 grok build（`auditor: grok-build (grok-4.6 · reasoning high)`，如 GOAL-019 A-007）。
- 规则依据：AGENTS.md §6b（P-003 交叉审计与意见响应 / P-004 用户裁决点）。

## 未选方案

- **继续使用其他 provider（如 GitHub Copilot `/audit`）**：workspace-008 D-002 已按用户指令切到 grok build；本决策延续该路径并升级模型版本，与既有实践一致。
- **不落盘、仅口头约定**：P-003/P-004 要求书面留痕；provider 与审计执行路径必须版本化落盘。
- **仅记录在某个工作区**：开新工作区后决策落空；故落盘于项目级 `docs/architecture/` 并在根 AGENTS.md 登记入口。
