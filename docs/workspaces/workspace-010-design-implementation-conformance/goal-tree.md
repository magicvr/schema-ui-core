---
title: 目标树 · workspace-010-design-implementation-conformance
status: active
created: 2026-08-11
updated: 2026-08-13
parent: null
version: 0.3.0
workspace_id: workspace-010-design-implementation-conformance
---

# 目标树 · 设计意图与实现符合性（持续对齐程序）

> 工作区：`workspace-010-design-implementation-conformance`
> canonical：`docs/workspaces/workspace-010-design-implementation-conformance/`
> Root：`GOAL-001-design-implementation-conformance`（**长期程序容器 · active**）
> primary_plan：`VP-010-design-implementation-conformance`（**active**）

## 树

```text
GOAL-001-design-implementation-conformance [active]  · 持续符合性程序
├── GOAL-002-w1-examples-optional-module [done]       · W1 范例面可选化
├── GOAL-003-demo-profile [done]                      · W2 demo Profile：mvp + 范例
├── GOAL-004-w3-schema-host-protocol-conformance [done] · W3 协议优先的 Host/App 符合性整改
└── GOAL-005-w4-long-content-presentation [active]    · W4 长内容列截断与详情换行
```

Root **保持 active**。W1/W2/W3 均关门；W3 六检查点全部完成（2026-08-13 关门：S6 cross 审计
A-007/A-008，BLOCKING 清零，用户 P-004 裁决 account-locked 实现生产源）；不推导 Root/VP done。
W4（GOAL-005）为当前进行中波次。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-design-implementation-conformance | 设计意图与实现符合性（持续对齐程序） | null | active | —（程序容器，不用 n/n→done） | 2026-08-13 |
| GOAL-002-w1-examples-optional-module | W1 · 范例/演示产品面可选模块化 | GOAL-001-design-implementation-conformance | done | 6/6 | 2026-08-11 |
| GOAL-003-demo-profile | W2 · `demo` Profile：mvp + 范例页面 | GOAL-001-design-implementation-conformance | done | 6/6 | 2026-08-11 |
| GOAL-004-w3-schema-host-protocol-conformance | W3 · Schema-UI 语义对齐与 Host/App 协议增补 | GOAL-001-design-implementation-conformance | done | 6/6 | 2026-08-13 |
| GOAL-005-w4-long-content-presentation | W4 · 长内容列的列表截断与详情换行（以角色页权限/菜单为代表） | GOAL-001-design-implementation-conformance | active | 5/6 | 2026-08-13 |

## 维护说明

- Root 是长期能力容器；`status: done` 仅在程序废弃或 `primary_plan` 迁移且用户确认时使用。
- 波次 progress 只写在子目标；不得用波次完成数推导 Root done。
- 层级唯一来源是目标 `00-meta.md` 的 `parent`。
