---
id: D-002
doc: decision-entry
goal: GOAL-001-admin-functional-modules
status: accepted
created: 2026-08-18
updated: 2026-08-18
version: 1.0.0
---

# D-002 · 四档能力地图与跨模块路线图登记（非立即实施）

## 背景

上一轮审视指出：workspace-011 的 `I-011-001` 已覆盖大量可见功能页，但遗漏三类跨模块能力——横切基架、未来扩展接缝、业务领域未列实体与流程。用户确认：不马上实施，但要把路线图记录下来，避免以后重复分析。

## 决策

1. **不**把三类能力继续堆进 workspace-011 的 S/B 模块编号，也不批量创建子目标。
2. 采用**四档能力地图**作为路线图登记框架：
   - Tier A：Admin 基架规划（通用能力）；
   - Tier B：扩展接缝（预留接口）；
   - Tier C：Admin 体验增强；
   - Tier D：真实业务领域。
3. 详细四档地图与推进顺序落盘到 Root 附件 `attachments/I-011-002-cross-module-roadmap.md`；Root 纲领路线图增加 **R5 登记阶段**。
4. VP-011 增加“能力分层边界”：三档方法论只用于功能模块；横切基架问题不自动进入 VP-011，按 VP-009/VP-010 或未来平台 VP 分流。
5. 组合层在 `docs/vision/roadmap.md` 登记“共享基架横切能力、扩展接缝与后续业务域（四档能力地图）”为后续方向。
6. 不修改 Charter，不改变任何 Goal status/progress，不创建新 Goal。

## 审计模式

文档路线图登记（低风险、可逆）：**none**；同步以 self Vision Review `VRev-027` 留痕。

## 未选方案

- 把所有遗漏能力直接扩为新的 S/B 模块：会导致通用 Admin 与交易/业务单体边界模糊。
- 立即创建多个“未来能力”子目标：用户明确“肯定不会马上就做”，应只登记路线图。
- 新建独立平台 VP：暂未到触发条件；先登记，待真实需求出现再 `/vision` 决定。
