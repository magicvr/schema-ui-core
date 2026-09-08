---
id: D-001-spacing-and-active-marker
goal_id: GOAL-004-sidebar-active-indicator
doc: decision-entry
source: user-request
status: accepted
date: 2026-09-08
parent: GOAL-001-nav-group-collapsible
version: 0.1.0
---

# D-001 · 紧凑组间距与 active 指示语义

## 决策

1. 非 horizontal 导航组间距从 `space-y-6` 调整为 `space-y-2`，对齐 Sidebar Engine 范例的 `space-y-sm` 密度。
2. active 页面左侧显示静态高亮竖线；页面右侧显示注册者提供的 secondary 副文本（如有），不再显示闪烁光点。
3. 没有 secondary 的 active 页面不渲染光点、空 span 或其它右侧占位；inactive 页面保持现有文本/hover 语义。
4. 该调整复用 desktop sidebar 与 mobile drawer 共用的 `NavigationLink`，不改变 active 计算、权限、折叠持久化或路由。

## 理由

- 参考页采用更紧凑的组间节奏；当前 24px 组间距造成侧栏垂直密度过低。
- secondary 是注册层的可选语义，应该替换状态装饰而不是与闪烁光点并列，避免同一右侧位置表达两种含义。
- 左侧竖线提供稳定的选中定位，不依赖动画或状态点。

## 验收指向

- `nav-groups.test.tsx`：active marker、secondary、有/无 secondary 与 compact spacing。
- `navigation.test.ts`：secondary projection 与 active 状态保持。
- `tsc -b` / `npm run build`：类型与生产构建通过。

