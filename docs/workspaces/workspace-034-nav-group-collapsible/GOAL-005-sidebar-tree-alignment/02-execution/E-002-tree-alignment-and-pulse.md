---
id: E-002-tree-alignment-and-pulse
goal_id: GOAL-005-sidebar-tree-alignment
doc: execution-entry
status: recorded
date: 2026-09-08
parent: GOAL-001-nav-group-collapsible
version: 0.1.0
---

# E-002 · 组内导航左翼与 active pulse 实施

## 已发生事实

- `apps/web/src/app/App.tsx` 移除展开组内容容器的 `border-l`、`ml-2`、`pl-3`，使组内叶子导航整体左移并回到组标题基础列。
- active 页面继续使用左侧静态 `bg-primary` 高亮竖线。
- active 页面有 secondary 时显示 secondary 文本且不显示 pulse；无 secondary 时显示 `data-navigation-active-dot="active"` 及 `animate-ping` 光点。
- inactive 页面无 secondary 时不显示 pulse；带 secondary 的页面继续显示注册副文本。
- `nav-groups.test.tsx` 增加 Roles 无 secondary 的 generic active fixture，覆盖无 group border/indent、active marker、secondary/pulse 互斥和 compact sidebar spacing。

## 证据路径

- `apps/web/src/app/App.tsx`
- `apps/web/src/app/nav-groups.test.tsx`
- `docs/workspaces/workspace-034-nav-group-collapsible/GOAL-005-sidebar-tree-alignment/01-decision/D-001-tree-alignment-and-pulse.md`

