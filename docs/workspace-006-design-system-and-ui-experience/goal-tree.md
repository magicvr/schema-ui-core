---
workspace_id: workspace-006-design-system-and-ui-experience
root_goal: GOAL-001-design-system-and-ui-experience
canonical_scope: docs/workspace-006-design-system-and-ui-experience/
status: active
created: 2026-08-09
updated: 2026-08-09
version: 0.8.0
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
GOAL-001-design-system-and-ui-experience [active] (2/5)
├── GOAL-002-s1-design-tokens-and-primitives [done] (6/6)
├── GOAL-003-s2-s3-renderer-and-shell [active] (0/2)
├── GOAL-004-s4-state-and-feedback [done] (3/3)
└── GOAL-005-s5-regression-fork-example-and-closeout [done] (2/2)
```

## 状态表

| ID | 标题 | Parent | Status | Progress | Updated |
|----|------|--------|--------|----------|---------|
| GOAL-001-design-system-and-ui-experience | 现代设计系统与 Schema 驱动 UI/UX 体验产品化 | `null` | **active** | `2/5` | 2026-08-09 |
| GOAL-002-s1-design-tokens-and-primitives | S1 · Design Token / 主题 / shadcn primitives | `GOAL-001-design-system-and-ui-experience` | **done** | `6/6` | 2026-08-09 |
| GOAL-003-s2-s3-renderer-and-shell | S2+S3 · Renderer 视觉接入 + Shell 断点升级 | `GOAL-001-design-system-and-ui-experience` | **active** | `0/2` | 2026-08-09 |
| GOAL-004-s4-state-and-feedback | S4 · 状态与反馈一致性 | `GOAL-001-design-system-and-ui-experience` | **done** | `3/3` | 2026-08-09 |
| GOAL-005-s5-regression-fork-example-and-closeout | S5 · 视觉回归 + fork Token 示例 + 过程关门 | `GOAL-001-design-system-and-ui-experience` | **done** | `2/2` | 2026-08-09 |

## 维护说明

- `2/5`：Root 仅 **S1、S4** 勾选；**S2/S3/S5 取消勾选**（D-006 / A-006）。
- GOAL-002 / GOAL-004 / GOAL-005 保持 `done`（局部交付真实）；**不得**单独推导 Root `done`。
- GOAL-003 回 `active`（`0/2`）：过窄 C1/C2 不得代表 Root S2/S3；开放 F-003-001 与 Root **F-VUI-001 / F-VUI-002**。
- **D-005 已 superseded**；**D-006 accepted**（用户要求回退完成状态）。
- 开放 required（Root）：**F-VUI-001、F-VUI-002**。F-VUI-003 fixed（状态回退）。
- 视觉范围权威：VP-005 type 表 + `I-PROTO-FULL-001` + D-004 Stitch 约束。
