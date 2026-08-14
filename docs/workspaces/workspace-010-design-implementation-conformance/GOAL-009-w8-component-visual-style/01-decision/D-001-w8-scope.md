---
id: D-001
goal: GOAL-009-w8-component-visual-style
date: 2026-08-14
status: accepted
parent: GOAL-009-w8-component-visual-style
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# D-001 · W8 范围与方案（三部分）

## 决策

- **P1 · 语种切换器下拉重构**（已交付）：原生 select → lucide Languages 图标触发（size-9 rounded-md ghost，与通知铃一致）+ 自制无 Portal 下拉；暗色 token 化（bg-popover≈neutral-900 / border-border≈white-10 / hover:bg-accent≈neutral-800 / 选中 ✓）；Escape/外部点击关闭。事实登记于 E-001（提交 46292e5）。
- **P2 · 明暗切换按钮统一**：ThemeToggle 弃用 outline Button，改与语种/铃一致图标幽灵样式（size-9 rounded-md text-muted-foreground hover:bg-accent hover:text-foreground）；aria-label/title/sr-only 改为 i18n 键（shell.theme.toggle，en/zh），消除固定英文「Toggle color theme」。
- **P3 · 下拉暗色审计**：全仓扫描 select/menu/listbox/popover；确认语种切换 bug（原生 select 弹层暗色下默认亮色）已随控件替换消除；表单 SelectField 原生 select 显式声明 `dark:scheme-dark`（控件级 color-scheme 兜底，root 已级联）；自制面板（铃/语种）复核 bg-card/bg-popover token 在暗色即正确。

## 未选方案

- 不引入 Radix/Headless 下拉库（项目自制轻量下拉已够用，避免新依赖与 portal 主题作用域问题）。
- 不整体替换表单 select（协议控件形态；仅修暗色声明）。
