---
id: E-002
goal: GOAL-009-w8-component-visual-style
date: 2026-08-14
status: recorded
parent: GOAL-009-w8-component-visual-style
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# E-002 · S2 实施：明暗切换按钮统一 + 下拉暗色修正

## 事实

### P2 · ThemeToggle 统一（theme-toggle.tsx）

- 弃用 outline Button，改图标幽灵样式：size-9 rounded-md text-muted-foreground hover:bg-accent hover:text-foreground（与语种/铃铛一致）；Sun/Moon size-4。
NaN
NaN

### P3 · 下拉暗色审计与修正

- 全仓扫描（I-001 closed）：原生 select 仅 form-controls SelectField 1 处；自制菜单面板 2 处（notification-bell bg-card、locale-switcher bg-popover）；无 Radix/Headless/AntD（无 portal 主题作用域问题）。
NaN
NaN
