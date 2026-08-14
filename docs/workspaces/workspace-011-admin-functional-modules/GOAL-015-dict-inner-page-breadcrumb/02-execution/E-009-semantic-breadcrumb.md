---
id: E-009
goal: GOAL-015-dict-inner-page-breadcrumb
date: 2026-08-14
status: recorded
parent: GOAL-015-dict-inner-page-breadcrumb
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# E-009 · 面包屑语义层级 + 视觉重构（D-004）

## 事实

### 1. 逻辑重构（访问历史 → 语义层级）

- `resolveBreadcrumbTrail` 重写：trail = [首页(homePageRef)] + [导航组标签...] + [声明父级链...] + [当前页]；首页恒为根（当前页即首页时单级无 UI）；未知父级 fail-safe；父链防环；同一页面无论如何到达 trail 恒定（与访问历史无关）。
- App.tsx：`visitStack` 状态、首屏种子、popstate 截断、onNavigate 压栈全部移除；`BREADCRUMB_PAGE_PARENTS` 常量（dictionary-entries → data-dictionary；task-runs → scheduled-tasks）；trail 解析传入 navigation + parents + homePageId。
- 返回按钮：最近的非首页路由祖先（语义父级），如条目页 → /data-dictionary；无语义父级（如组标签祖先）时不显示。

### 2. 视觉重构（深色现代后台）

- 移除纯文本「← 返回上一页 |」拼接与「ADMIN WORKSPACE」眉标；面包屑与标题间距 mt-2.5（10px）。
- Breadcrumbs 组件：text-xs 常规字重；祖先 text-muted-foreground + hover:text-foreground + hover:underline；分隔符细斜杠 `/`（text-muted-foreground/40）；当前项 text-foreground/90 不可点击；返回按钮 = 圆形幽灵图标按钮（lucide ArrowLeft，size-6 rounded-full hover:bg-muted），aria-label 用 shell.back。

### 3. 测试与验证

- breadcrumbs.test.ts 重写 9 例（首页根/组标签/声明父级/未知父级 fail-safe/防环/与历史无关/首页自身单级）。
- App.integration：深链 /catalog/42 直接显示「Home / Operations / Catalog detail」（无访问历史）；行导航到条目页显示「Home / Data dictionary / Dictionary entries」+ 返回按钮点击回 /data-dictionary（10/10）。
- 全量回归：web 53 文件 / 958 测试绿；tsc 仅剩 4 处历史既有错误（非本阶段引入）。
- 边界：group 标签为纯文本段（不可点击）；homePageRef 未在 pages 中解析时首页段 fail-safe 跳过。
