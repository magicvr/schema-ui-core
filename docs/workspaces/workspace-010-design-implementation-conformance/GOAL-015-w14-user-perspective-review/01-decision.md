---
id: GOAL-015-w14-user-perspective-review
doc: decision
status: done
parent: GOAL-001-design-implementation-conformance
created: 2026-08-17
updated: 2026-08-17
version: 0.6.0
---

# 决策记录 · GOAL-015

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|
| I-001 | required → **已裁决** | F-01～F-14 的 in-scope / defer / 优先级 | 整改子目标范围（R1～R4） | R1～R4 | 用户裁决（P-004） | **closed** | 用户 2026-08-17 书面裁决（D-003）：全部 in-scope、分批 A→C→D→B 作为 GOAL-015 子目标；三方案选择冻结。裁决为整改输入；GOAL-015 待整改完成后关门 |
| I-002 | non-blocking | F-01 handler 目录暴露方式 | 未来整改波次（F-01） | 未来整改波次 S3 | as-built + 方案 | **collecting** | 见 D-001 §3 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-17 | 开波：真实用户视角审视结论 + 改进项台账（F-01～F-14）+ 建议实施顺序 | accepted（范围+台账） | `01-decision/D-001-w14-scope-and-findings.md` |
| D-002 | 2026-08-17 | 波次交付边界：本波 = 审视 + 台账落盘 + 审计会签 + 台账同步；F-01～F-14 整改 deferred（待用户裁决，I-001 门禁移至未来整改波次）——**「移至未来波次并把 S4 视为可放行关门」部分已被用户修正（E-003/A-004）** | accepted（范围，部分被修正） | `01-decision/D-002-w14-delivery-boundary.md` |
| D-003 | 2026-08-17 | **I-001 用户书面裁决**：F-01～F-14 全部 in-scope（分批 A→C→D→B，另起整改波次）；F-01 新增端点、F-04 存 messageKey、F-08 直接移除 | accepted（用户裁决） | `01-decision/D-003-i001-user-adjudication.md` |

## 待决问题（P-004）

- **本波（active · 4/8，整改承接中）**：S1/S2 只读完成（A-001 self pass）；S3 独立审计（A-002 independent pass，无 required）；S4 审计响应 + 台账同步 + I-001 用户书面裁决（D-003）：F-01～F-14 全部 in-scope（分批 A→C→D→B）+ 三方案选择（F-01 新增端点 / F-04 存 messageKey / F-08 直接移除）。**用户结构裁决（E-005/A-006）：整改完成前 GOAL-015 不得 done；整改子目标为 GOAL-015 下级（渐进添加）。** 曾两次关门尝试（E-002/A-003、E-004/A-005）均被用户否决/修正。
- **整改实施（GOAL-015 子目标分批）**：批次 A（F-01～F-04，GOAL-016）→ 批次 C（F-08～F-10）→ 批次 D（F-11～F-14）→ 批次 B（F-05～F-07），作为 GOAL-015 下级子目标渐进添加；全部完成 + S5 终审后 GOAL-015 方可关门。