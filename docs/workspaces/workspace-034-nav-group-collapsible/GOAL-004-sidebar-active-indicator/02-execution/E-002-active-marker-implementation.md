---
id: E-002-active-marker-implementation
goal_id: GOAL-004-sidebar-active-indicator
doc: execution-entry
status: recorded
date: 2026-09-08
parent: GOAL-001-nav-group-collapsible
version: 0.1.0
---

# E-002 · 紧凑间距与 active 指示实施

## 已发生事实

- `apps/web/src/app/App.tsx` 将非 horizontal `NavigationItems` 的组间距从 `space-y-6` 调整为 `space-y-2`，desktop sidebar 与 mobile drawer 共用该间距。
- `NavigationLink` 在非 horizontal surface 为 active 页面增加静态左侧高亮竖线 `data-navigation-active-marker="active"`；inactive 页面保留透明对齐槽位。
- active 页面继续使用注册 projection 提供的 secondary 副文本；没有 secondary 时不渲染 secondary 元素，也不渲染闪烁光点或占位元素。
- 更新 `nav-groups.test.tsx`，覆盖 compact spacing、active marker、secondary 显示与 `.animate-ping` 缺省断言；既有 navigation/deep-link matrix 保持通过。

## 证据路径

- `apps/web/src/app/App.tsx`
- `apps/web/src/app/nav-groups.test.tsx`
- `apps/web/src/app/navigation.ts`

