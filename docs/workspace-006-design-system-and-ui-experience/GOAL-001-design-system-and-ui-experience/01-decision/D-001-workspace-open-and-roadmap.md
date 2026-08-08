---
id: GOAL-001-design-system-and-ui-experience
doc: decision-entry
record_id: D-001
status: accepted
parent: null
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
---

## D-001 · 开区与纲领路线图采纳

### 触发

- 用户 2026-08-09 经 `/vision`：闭合 VRev-011 `F-V018`/`F-V019`/`F-V020`，并确认 **激活 VP-005**（`active` v0.4.0）。
- 用户 2026-08-09 经 `/govern` 明确指令：为 active VP-005 开区，**确认** workspace slug = `workspace-006-design-system-and-ui-experience`，scaffold Root 与纲领路线图。
- Vision Review open required = 0；VP-006 **closed**（硬前置已满足）。

### 决定

1. 新建 delivery 工作区 **`workspace-006-design-system-and-ui-experience`**（编号 = 现有最大工作区 + 1；slug **用户书面确认**）。
2. Root **`GOAL-001-design-system-and-ui-experience`**，`parent: null`，`primary_plan` / `plan_refs` = `VP-005-design-system-and-ui-experience`，`vision_role: delivery`。
3. 采纳 VP-005 建议阶段 **S1–S5** 为本 Root 等权纲领检查点（P-001）；建区当日均未完成（`progress: 0/5`）。
4. 视觉范围分母 = VP-005 钉死 type 表 + `I-PROTO-FULL-001` include（只读引用 workspace-005 覆盖表）；**禁止**扩张协议 disposition；**禁止**杜撰 `Detail`/`Filter` Node 名。
5. WCAG AA / Cmd+K 默认 **不进** 退出分母（F-V019 路径 b）；对比度升格须用户书面裁决（I-004）。
6. **不**将激活/建区读成视觉产品化已交付；**禁止**在 closed workspace-003/004/005 吸收本意图。

### 为什么

- 结构选型：新纲领波次 → 已 active VP + 独立 delivery 工作区（P-006）。
- 协议硬前置与 Vision required 已齐；用户确认激活与 slug。
- 大目标先纲领路线图（P-001），再按阶段推进。

### 未选方案

| 方案 | 未选原因 |
|------|----------|
| 在 workspace-005 继续挂视觉实施 | 违反 VP-005/VP-006 禁止在 closed 协议区吸收视觉意图；区职责混淆 |
| 在 workspace-003/004 吸收 | closed 区默认不接新区；架构/playbook 已交付 |
| 跳过纲领直接大批量建子目标 | 违反 P-001；S1 Token 基线未就绪 |
