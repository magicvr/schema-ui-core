---
id: GOAL-015-w14-user-perspective-review
doc: decision
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-08-17
updated: 2026-08-17
version: 0.3.0
---

# 决策记录 · GOAL-015

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|
| I-001 | required（对整改波次） | F-01～F-14 的 in-scope / defer / 优先级 | 未来整改波次 | 未来整改波次 S2 | 用户裁决（P-004） | **deferred** | D-002：本波只审视+落盘，不实施整改；延期触发＝用户指示开始整改（A-002 核验三字段齐备） |
| I-002 | non-blocking | F-01 handler 目录暴露方式 | 未来整改波次（F-01） | 未来整改波次 S3 | as-built + 方案 | **collecting** | 见 D-001 §3 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-17 | 开波：真实用户视角审视结论 + 改进项台账（F-01～F-14）+ 建议实施顺序 | accepted（范围+台账） | `01-decision/D-001-w14-scope-and-findings.md` |
| D-002 | 2026-08-17 | 波次交付边界：本波 = 审视 + 台账落盘 + 审计会签 + 台账同步；F-01～F-14 整改 deferred（待用户裁决，I-001 门禁移至未来整改波次） | accepted（范围） | `01-decision/D-002-w14-delivery-boundary.md` |

## 待决问题（P-004）

- **本波**：无阻断待决。S1/S2 只读完成（A-001 self pass）；S3 独立审计（A-002 independent pass，无 required）；S4 响应三条 non-blocking 全 fixed 并关门（A-003 self pass）。GOAL-015 `done`（4/4）。
- **面向未来整改波次（deferred，非本波）**：F-01～F-14 的 in-scope / defer / 优先级需用户裁决；其中 F-01 是否新增「handler 目录」端点/静态选项、F-04 通知本地化方案（存 messageKey vs 成品文案）、F-08 调试框移除 vs「开发者模式」开关，三个方案选择一并留待用户。