---
id: E-008-fluid-shell-width-fix
title: 修复主内容区不随浏览器宽度自适应（用户驳回关门缺口）
date: 2026-08-09
status: recorded
parent: GOAL-001-design-system-and-ui-experience
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
---

# E-008 · 流体壳宽修复

## 触发

用户书面驳回 closeout（非 D-008 确认）：「所有页面用浏览器打开的时候，除了顶部导航条，其他都还没有自适应浏览器的宽度。」

## 根因

`App.tsx` 将 **sidenav + main** 包在 `mx-auto max-w-[1440px]` 中，顶栏全宽、内容区成固定宽度岛，与浏览器宽度不同步。

## 修复事实

| 改动 | 路径 |
|------|------|
| 去掉 body 级 max-w；`data-shell-width=fluid` + `w-full` | `apps/web/src/app/App.tsx` |
| main / page `w-full min-w-0`；table 横向滚动 | `App.tsx`、`data-table.tsx`、`schema-table.tsx` |
| Schema grid 小屏 1 列、md+ 多列 | `render.tsx` + `index.css` |
| 结构测试 | `shell.test.ts` |

## 验证

- vitest + build 见本轮 regression 日志  
