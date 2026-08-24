---
id: D-001
doc: decision-entry
goal: GOAL-007-mock-resend-delivery
status: accepted
created: 2026-08-24
updated: 2026-08-24
version: 1.0.0
---

# D-001 · 开设 R6 子目标

## 背景

Root R5 已完成：GOAL-006 D-002 冻结渠道合同，I-010 / I-011 verified，R6 门禁解除。P-001：按阶段开设子目标。

## 决定

1. 创建 `GOAL-007-mock-resend-delivery`，parent = `GOAL-001-outbound-mail`，status `active`。
2. 实施范围 = GOAL-006 D-002 条款 §1～§4 的代码落地（配置 / 迁移 + mock / API + 接线 / Resend）；**无新增合同裁决**——实施细节若与冻结节冲突，回退到 D-002 而不是改写它。
3. 实施内技术选型（非合同级）：迁移归属 corepersistence（core 表先例）；mock 记录 id 沿用仓库现行 id 生成器；Resend 探针留空待 R8。这些不改变 D-002 冻结语义，无需再问询。
4. 审计模式：scaffold none → 关门 self（沿 GOAL-002～006 先例）。出现安全/迁移语义分歧时升级 P-004 问询。

## 未选方案

- 先问一轮「实现细节」再开工：I-010/I-011 已 verified 且 D-002 条款可直接实施，P-004 无未决裁决点，不应空转问询。
- 本目标顺带做设置页/试发：违反 R6/R7 阶段切分。
