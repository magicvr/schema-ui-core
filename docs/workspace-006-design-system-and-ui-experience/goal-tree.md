---
workspace_id: workspace-006-design-system-and-ui-experience
root_goal: GOAL-001-design-system-and-ui-experience
canonical_scope: docs/workspace-006-design-system-and-ui-experience/
status: active
created: 2026-08-09
updated: 2026-08-09
version: 0.9.1
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
GOAL-001-design-system-and-ui-experience [active] (4/5) · closeout-ready
├── GOAL-002-s1-design-tokens-and-primitives [done] (6/6)
├── GOAL-003-s2-s3-renderer-and-shell [done] (2/2)
├── GOAL-004-s4-state-and-feedback [done] (3/3)
└── GOAL-005-s5-regression-fork-example-and-closeout [done] (2/2)
```

## 状态表

| ID | 标题 | Parent | Status | Progress | Updated |
|----|------|--------|--------|----------|---------|
| GOAL-001-design-system-and-ui-experience | 现代设计系统与 Schema 驱动 UI/UX 体验产品化 | `null` | **active** | `4/5` | 2026-08-09 |
| GOAL-002-s1-design-tokens-and-primitives | S1 · Design Token / 主题 / shadcn primitives | `GOAL-001-design-system-and-ui-experience` | **done** | `6/6` | 2026-08-09 |
| GOAL-003-s2-s3-renderer-and-shell | S2+S3 · Renderer 视觉接入 + Shell 断点升级 | `GOAL-001-design-system-and-ui-experience` | **done** | `2/2` | 2026-08-09 |
| GOAL-004-s4-state-and-feedback | S4 · 状态与反馈一致性 | `GOAL-001-design-system-and-ui-experience` | **done** | `3/3` | 2026-08-09 |
| GOAL-005-s5-regression-fork-example-and-closeout | S5 · 视觉回归 + fork Token 示例 + 过程关门 | `GOAL-001-design-system-and-ui-experience` | **done** | `2/2` | 2026-08-09 |

## 维护说明

- Root **closeout-ready**：`active`；`4/5` = S1–S4 勾选；**S5 过程关门**待用户显式确认（E-007；D-007 superseded）。
- 开放 required = **0**（F-VUI-001/002 fixed via A-008/A-009）。F-VUI-007 accepted-residual。
- 子目标 GOAL-002～005 均为 `done`（局部交付真实）；**不得**单独推导 Root `done`。
- D-005 superseded；D-006 accepted；D-007 **superseded**。下次关门落盘 **D-008**。
- 视觉范围权威：VP-005 type 表 + `I-PROTO-FULL-001` + D-004。
