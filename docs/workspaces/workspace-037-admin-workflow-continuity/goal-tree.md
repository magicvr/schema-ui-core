---
title: 目标树 · workspace-037-admin-workflow-continuity
status: active
created: 2026-09-16
updated: 2026-09-18
parent: null
version: 1.10.0
workspace_id: workspace-037-admin-workflow-continuity
---

# 目标树 · Admin 工作流连续性与安全反馈

> 工作区：`workspace-037-admin-workflow-continuity`
> canonical：`docs/workspaces/workspace-037-admin-workflow-continuity/`
> Root：`GOAL-001-admin-workflow-continuity`（**active · 5/6**）
> primary_plan：`VP-037-admin-workflow-continuity`（active · v1.4.0）

## 目标树

```text
GOAL-001-admin-workflow-continuity [active · 5/6]
├── GOAL-002-r1-scope-semantics-freeze [done · 3/3]
├── GOAL-003-r2-saved-views [done · 4/4]
├── GOAL-004-r3-unsaved-change-protection [done · 4/4]
├── GOAL-005-r4-unified-feedback-recovery [done · 4/4]
├── GOAL-006-r5-composition-acceptance [active · 3/4]
├── GOAL-007-list-page-visual-alignment [done · 8/8]
├── GOAL-008-typecheck-evidence-convention [done · 4/4]
└── GOAL-009-list-visual-e2e-guard [done · 4/4]
```

## 纲领路线图

```text
R1 分母与语义冻结 [done · GOAL-002 · 3/3]
  → { R2 Saved Views [done · GOAL-003 · 4/4]
      R3 未保存变更保护 [done · GOAL-004 · 4/4]
      R4 统一反馈与恢复 [done · GOAL-005 · 4/4] }
  → R5 组合验收与关门 [active · GOAL-006 · 3/4]
  → R6 列表页视觉与筛选体验收敛 [done · GOAL-007 · 8/8]
  · 整改（非纲领）类型检查证据约定 [done · GOAL-008 · 4/4]
  · 整改（非纲领）列表视觉 e2e 守卫 [done · GOAL-009 · 4/4]
```

> R2/R3/R4 在 R1 完成后可按独立证据与并行价值并行；上图表示门禁顺序，不表示任何实现已完成。

## 状态表

| id | title | status | progress | parent | notes |
|----|-------|--------|----------|--------|-------|
| GOAL-001-admin-workflow-continuity | Admin 工作流连续性与安全反馈交付 | **active** | 5/6 | null | R1/R2/R3/R4/R6 done；R6 `GOAL-007` 以 `done · 8/8` 关门（C6 审计 A-002 conditional，F-005 经 GOAL-008 闭环，F-003 经 GOAL-009 闭环）；整改子目标 `GOAL-008` `done · 4/4`、`GOAL-009` `done · 4/4`（非纲领，不计入分母）；I-037-001～004/007 verified；I-037-006 verified；Vision open required = 0；R5 `GOAL-006` active 3/4；R5-I-004 用户书面关门确认仍开放 |
| GOAL-002-r1-scope-semantics-freeze | R1 列表分母与工作流语义冻结 | **done** | 3/3 | GOAL-001-admin-workflow-continuity | C1/C2/C3 已完成；A-001 self、A-002 independent 均 pass；F-002～F-004 为非阻断 recommended，已转入 R2 |
| GOAL-003-r2-saved-views | R2 用户级 Saved Views 闭环 | **done** | 4/4 | GOAL-001-admin-workflow-continuity | D-001 已冻结 24 个 `type: table` 分母、custom 排除与活 Schema allowlist；C1～C4 完成，A-001/A-002/A-003 pass，checkpoint `39c744ef` |
| GOAL-004-r3-unsaved-change-protection | R3 未保存变更保护与离开确认 | **done** | 4/4 | GOAL-001-admin-workflow-continuity | D-001 已承接 R1 D-004 dirty-state 合同；C1～C4 完成，A-003 independent recheck、A-004 self pass，checkpoint `d2b39189` |
| GOAL-005-r4-unified-feedback-recovery | R4 统一反馈与恢复 | **done** | 4/4 | GOAL-001-admin-workflow-continuity | C1～C4 完成；A-002 F-001 经 A-003 independent recheck fixed/pass，A-004 self pass，checkpoint `89666e5c`；F-002 Host/resource 对照为不阻断 recommended |
| GOAL-006-r5-composition-acceptance | R5 组合验收与关门准备 | **active** | 3/4 | GOAL-001-admin-workflow-continuity | 已按 D-010 开设并完成 C1～C3 组合核对与最终验证；2026-09-18 按用户要求产出关门证据一页（`attachments/r5-closeout-evidence-digest.md`/E-006），R5-I-004 仍 collecting，Root/VP 尚未关门 |
| GOAL-007-list-page-visual-alignment | R6 列表页视觉与筛选体验收敛 | **done** | 8/8 | GOAL-001-admin-workflow-continuity | C1～C8 全部完成：C5 布局修订（E-004）、C7 控制位纠偏（D-003/E-005）、C8 控件语义（D-004/E-006）；C6 审计 A-002 `conditional`（F-001/F-002/F-004 fixed；F-005 跨区部分经 GOAL-008 闭环）；不改变顶部功能栏、左侧导航与查询/重置合同 |
| GOAL-008-typecheck-evidence-convention | 类型检查证据约定纠偏与防复发 | **done** | 4/4 | GOAL-001-admin-workflow-continuity | 承接 R6 A-002 F-005（high required）跨工作区部分；用户 P-004 裁决方案 A。C1 空转证明与影响面、C2 `npm run typecheck` 入口、C3 守卫（6 断言 + CI 门禁，5/5 变异捕获，含 e2e 第二层缺口修复）、C4 self 审计 `pass`（开放 required = 0）。2026-09-18 按用户 D-015 授权执行跨区勘误（E-006：11 处注记 + 2 处复核有效），I-008-004 转 verified；同期实测发现 `-p tsconfig.json` 同样空转，记 `A-002 F-001`（守卫缺口，recommended open）。非纲领阶段，不改变 Root 六阶段分母 |
| GOAL-009-list-visual-e2e-guard | 列表视觉浏览器级回归守卫 | **done** | 4/4 | GOAL-001-admin-workflow-continuity | 承接 R6 A-002 F-003（recommended）。新增 `e2e/list-visual-surface.spec.ts`（roles 页；C7 配对/图标/视图表单、C8 高度/token/开关存在性、C5 布局顺序与列表内 footer）；mvp/admin 两 profile 均通过，6/6 变异捕获（含当年由用户发现的 C5 配对与 C7 表单位置两类）；self 审计 `pass`。非纲领阶段 |

## 维护说明

- Root `progress: 5/6` 由 `00-meta.md` 的 6 个显式纲领检查点派生；`GOAL-008`/`GOAL-009` 为非纲领整改子目标、不计入 Root 分母。它只表示阶段检查点展示，不放行方案、实施、验收或关门。R5 用户确认门禁仍开放，因此 Root/VP 尚未完成。
- `GOAL-008` 的跨区勘误（2026-09-18）只改 `docs/workspaces/workspace-002/009/010/011` 中受影响条目的**注记**，不改这些工作区的 `status`/`progress`/goal-tree；其台账不在本文件登记。
- 新建阶段子目标前，先在 Root 决策/路线图中冻结阶段边界；目标文件夹在本工作区根平铺。
- status/progress/parent 或新增子目标发生变化时，必须同步本文件树与状态表。
