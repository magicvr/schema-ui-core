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
| I-001 | **required（本波关门）** | F-01～F-14 的 in-scope / defer / 优先级 | 本波 S4 关门 | 本波 S4 关门 | 用户裁决（P-004） | **open** | 前次执行擅将本项 deferred 并关门 done（绕过 P-004）；用户 2026-08-17 裁决回退关门（A-004/E-003）。本波关门须先取得用户裁决 |
| I-002 | non-blocking | F-01 handler 目录暴露方式 | 未来整改波次（F-01） | 未来整改波次 S3 | as-built + 方案 | **collecting** | 见 D-001 §3 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-17 | 开波：真实用户视角审视结论 + 改进项台账（F-01～F-14）+ 建议实施顺序 | accepted（范围+台账） | `01-decision/D-001-w14-scope-and-findings.md` |
| D-002 | 2026-08-17 | 波次交付边界：本波 = 审视 + 台账落盘 + 审计会签 + 台账同步；F-01～F-14 整改 deferred（待用户裁决，I-001 门禁移至未来整改波次） | accepted（范围） | `01-decision/D-002-w14-delivery-boundary.md` |

## 待决问题（P-004）

- **本波（重新推进中）**：S1/S2 只读完成（A-001 self pass）；S3 独立审计（A-002 independent pass，无 required）。S4 曾于 2026-08-17 关门（A-003 self pass，`done` 4/4），但**被用户裁决回退**：关门未先取得 I-001（F-01～F-14 的 in-scope/defer/优先级，required）的用户书面裁决，属绕过 P-004 门禁（A-004/E-003）。GOAL-015 现 **active（3/4）**。
- **本波关门所需裁决（I-001）**：用户须对 F-01～F-14 逐项（或按建议批次）裁决 in-scope / defer / 优先级；其中 F-01 handler 目录暴露方式（新增端点 / 静态选项 / fork 扩展点）、F-04 通知本地化方案（存 messageKey vs 成品文案）、F-08 调试框移除 vs「开发者模式」开关，三个方案选择一并裁决。取得裁决后 S4 方可关门。