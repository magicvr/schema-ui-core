---
workspace_id: workspace-006-design-system-and-ui-experience
root_goal: GOAL-001-design-system-and-ui-experience
canonical_scope: docs/workspace-006-design-system-and-ui-experience/
status: active
created: 2026-08-09
updated: 2026-08-09
version: 0.3.0
parent: null
---

# 目标树 · 设计系统与 Schema 驱动 UI/UX 工作区

| 字段 | 值 |
|------|----|
| 工作区 | `workspace-006-design-system-and-ui-experience` |
| canonical 范围 | `docs/workspace-006-design-system-and-ui-experience/` |
| Root Goal | `GOAL-001-design-system-and-ui-experience` |
| primary plan | `VP-005-design-system-and-ui-experience` |

## ASCII 树

```text
GOAL-001-design-system-and-ui-experience [active] (4/5)
├── GOAL-002-s1-design-tokens-and-primitives [done] (6/6)
├── GOAL-003-s2-s3-renderer-and-shell [done] (2/2)
└── GOAL-004-s4-state-and-feedback [done] (3/3)
```

## 状态表

| ID | 标题 | Parent | Status | Progress | Updated |
|----|------|--------|--------|----------|---------|
| GOAL-001-design-system-and-ui-experience | 现代设计系统与 Schema 驱动 UI/UX 体验产品化 | `null` | **active** | `4/5` | 2026-08-09 |
| GOAL-002-s1-design-tokens-and-primitives | S1 · Design Token / 主题 / shadcn primitives | `GOAL-001-design-system-and-ui-experience` | **done** | `6/6` | 2026-08-09 |
| GOAL-003-s2-s3-renderer-and-shell | S2+S3 · Renderer 视觉接入 + Shell 断点升级 | `GOAL-001-design-system-and-ui-experience` | **done** | `2/2` | 2026-08-09 |
| GOAL-004-s4-state-and-feedback | S4 · 状态与反馈一致性 | `GOAL-001-design-system-and-ui-experience` | **done** | `3/3` | 2026-08-09 |

## 维护说明

- `4/5` 由 Root `00-meta.md` 中 S1～S5 五个等权检查点派生（S1/S2/S3/S4 已完成；S5 未开始）。
- `6/6` 由 GOAL-002 六个等权检查点派生（C1–C6，均已完成；GOAL-002 status: done）。
- `2/2` 由 GOAL-003 两个等权检查点派生（C1 S2/C2 S3，均已完成；GOAL-003 status: done）。
- `3/3` 由 GOAL-004 三个等权检查点派生（C1 纯函数+单测/C2 Skeleton 消费点统一/C3 回归不回退，均已完成；GOAL-004 status: done）。
- 本区为 VP-005 唯一 lead / delivery（2026-08-09 `/govern` scaffold；slug 用户书面确认）。
- 建区 **不**勾选检查点；**不**宣称视觉产品化已交付。
- 视觉范围权威：VP-005 type 表 + `I-PROTO-FULL-001`（workspace-005 只读引用）；不得扩张协议 disposition。
- **D-004（2026-08-09）**：Stitch 视觉方向已冻结（I-005 closed）；证据 E-004 + `GOAL-001/.../attachments/visual-direction-stitch-summary.md`；本地截图在 `raw/`（gitignore）。开放 required finding 仍 **F-002**（S1 完成门禁）。
- **GOAL-002（2026-08-09）**：S1 子目标五件套补齐（承接 S1 治理上下文）；实施输入 = Root D-002/D-003/D-004；未勾选任何检查点。
- **GOAL-004（2026-08-09）**：S4 子目标五件套补齐并完成；三条主渲染路径（DataTable/StatCardView/ChartView）loading 态统一改用 `Skeleton` primitive；新增纯函数 `resolveAsyncDisplayState` 承载判定逻辑；vitest 607→613 全绿；build exit 0。
