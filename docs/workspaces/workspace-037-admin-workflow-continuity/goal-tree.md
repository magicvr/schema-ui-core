---
title: 目标树 · workspace-037-admin-workflow-continuity
status: active
created: 2026-09-16
updated: 2026-09-18
parent: null
version: 1.7.0
workspace_id: workspace-037-admin-workflow-continuity
---

# 目标树 · Admin 工作流连续性与安全反馈

> 工作区：`workspace-037-admin-workflow-continuity`
> canonical：`docs/workspaces/workspace-037-admin-workflow-continuity/`
> Root：`GOAL-001-admin-workflow-continuity`（**active · 4/6**）
> primary_plan：`VP-037-admin-workflow-continuity`（active · v1.2.0）

## 目标树

```text
GOAL-001-admin-workflow-continuity [active · 4/6]
├── GOAL-002-r1-scope-semantics-freeze [done · 3/3]
├── GOAL-003-r2-saved-views [done · 4/4]
├── GOAL-004-r3-unsaved-change-protection [done · 4/4]
├── GOAL-005-r4-unified-feedback-recovery [done · 4/4]
├── GOAL-006-r5-composition-acceptance [active · 3/4]
├── GOAL-007-list-page-visual-alignment [active · 7/8]
└── GOAL-008-typecheck-evidence-convention [active · 2/4]
```

## 纲领路线图

```text
R1 分母与语义冻结 [done · GOAL-002 · 3/3]
  → { R2 Saved Views [done · GOAL-003 · 4/4]
      R3 未保存变更保护 [done · GOAL-004 · 4/4]
      R4 统一反馈与恢复 [done · GOAL-005 · 4/4] }
  → R5 组合验收与关门 [active · GOAL-006 · 3/4]
  → R6 列表页视觉与筛选体验收敛 [active · GOAL-007 · 7/8]
```

> R2/R3/R4 在 R1 完成后可按独立证据与并行价值并行；上图表示门禁顺序，不表示任何实现已完成。

## 状态表

| id | title | status | progress | parent | notes |
|----|-------|--------|----------|--------|-------|
| GOAL-001-admin-workflow-continuity | Admin 工作流连续性与安全反馈交付 | **active** | 4/6 | null | R1/R2/R3/R4 done；R6 `GOAL-007` 已完成原 C1～C4 后按用户选择 A 回开，现为 active 7/8（C5/C7/C8 已完成；C6 审计 A-002 conditional，F-005 跨区部分移交 GOAL-008）；GOAL-008 `active · 2/4` 为非纲领整改子目标；I-037-001～004/007 verified；I-037-006 verified；Vision open required = 0；R5 `GOAL-006` active 3/4；R5-I-004 用户书面关门确认仍开放 |
| GOAL-002-r1-scope-semantics-freeze | R1 列表分母与工作流语义冻结 | **done** | 3/3 | GOAL-001-admin-workflow-continuity | C1/C2/C3 已完成；A-001 self、A-002 independent 均 pass；F-002～F-004 为非阻断 recommended，已转入 R2 |
| GOAL-003-r2-saved-views | R2 用户级 Saved Views 闭环 | **done** | 4/4 | GOAL-001-admin-workflow-continuity | D-001 已冻结 24 个 `type: table` 分母、custom 排除与活 Schema allowlist；C1～C4 完成，A-001/A-002/A-003 pass，checkpoint `39c744ef` |
| GOAL-004-r3-unsaved-change-protection | R3 未保存变更保护与离开确认 | **done** | 4/4 | GOAL-001-admin-workflow-continuity | D-001 已承接 R1 D-004 dirty-state 合同；C1～C4 完成，A-003 independent recheck、A-004 self pass，checkpoint `d2b39189` |
| GOAL-005-r4-unified-feedback-recovery | R4 统一反馈与恢复 | **done** | 4/4 | GOAL-001-admin-workflow-continuity | C1～C4 完成；A-002 F-001 经 A-003 independent recheck fixed/pass，A-004 self pass，checkpoint `89666e5c`；F-002 Host/resource 对照为不阻断 recommended |
| GOAL-006-r5-composition-acceptance | R5 组合验收与关门准备 | **active** | 3/4 | GOAL-001-admin-workflow-continuity | 已按 D-010 开设并完成 C1～C3 组合核对与最终验证；继续审计与用户确认门禁，Root/VP 尚未关门 |
| GOAL-007-list-page-visual-alignment | R6 列表页视觉与筛选体验收敛 | **active** | 7/8 | GOAL-001-admin-workflow-continuity | 原 C1～C4 与 A-001 self 已完成；按 D-012 回开，C5 实现/回归（E-004）、C7 控制位纠偏（D-003/E-005）与 C8 控件语义（D-004/E-006）已完成；C6 审计 A-002 已记录（conditional：F-001/F-002/F-004 已闭合，F-005 跨区部分移交 GOAL-008）；不改变顶部功能栏、左侧导航与查询/重置合同 |
| GOAL-008-typecheck-evidence-convention | 类型检查证据约定纠偏与防复发 | **active** | 2/4 | GOAL-001-admin-workflow-continuity | 承接 R6 A-002 F-005（high required）的跨工作区部分；用户 P-004 裁决为方案 A（立独立目标）。C1 空转证明与影响面、C2 `npm run typecheck` 入口已固化；C3 守卫形态待决策，C4 待审计。非纲领阶段，不改变 Root 六阶段分母 |

## 维护说明

- Root `progress: 4/6` 由 `00-meta.md` 的 6 个显式纲领检查点派生；R6 回开与追加 C7/C8 不改变六阶段分母（C7/C8 只改变 R6 子目标自身的 8 项检查点），`GOAL-008` 为非纲领整改子目标、不计入 Root 分母。它只表示阶段检查点展示，不放行方案、实施、验收或关门。R5 用户确认门禁仍开放，因此 Root/VP 尚未完成。
- 新建阶段子目标前，先在 Root 决策/路线图中冻结阶段边界；目标文件夹在本工作区根平铺。
- status/progress/parent 或新增子目标发生变化时，必须同步本文件树与状态表。
