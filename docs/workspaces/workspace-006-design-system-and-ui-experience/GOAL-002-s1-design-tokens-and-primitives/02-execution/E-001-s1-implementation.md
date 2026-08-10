---
id: E-001-s1-implementation
title: S1 实施 — Token 扩写、主题模块、primitives、消费闭环
date: 2026-08-09
status: done
parent: GOAL-002-s1-design-tokens-and-primitives
---

# E-001 · S1 实施完成记录

## 事实摘要

**日期**：2026-08-09  
**scope**：C1–C6（GOAL-002 全部检查点）

### C1 — index.css Token 扩写

`apps/web/src/index.css` 按五个注释区块扩写（§1 Color / §2 Typography / §3 Elevation·Shadow / §4 @theme inline / §5 Base resets）。新增或补全：

- `--destructive` / `--destructive-foreground`（`:root` + `.dark`）
- `--success` / `--success-foreground`（`:root` + `.dark`）
- `--chart-1…5`（`:root` + `.dark`，深色轻微去饱和）
- `--overlay`（`:root` 40% 不透明度 / `.dark` 55%）
- `--font-sans` / `--font-mono`（`:root` Typography 区）
- `--elevation-sm|md|lg`（`:root` Elevation 区，原始 oklch shadow 值）
- `--popover` / `--popover-foreground`（`:root` + `.dark`）

### C2 — F-002 映射（@theme inline → shadow alias）

`@theme inline` 中新增三行：

```css
--shadow-sm: var(--elevation-sm);
--shadow-md: var(--elevation-md);
--shadow-lg: var(--elevation-lg);
```

无自引用（`--shadow-*` 均指向 `--elevation-*`，未出现 `var(--shadow-sm)` 循环）。
Tailwind `shadow-sm|md|lg` utility 可生成并消费语义值。

### C3 — FOUC 治理

- `apps/web/index.html` 新增同步内联 `<script>`：读 `localStorage.theme` + `prefers-color-scheme`，首帧即写 `.dark` class + `color-scheme` 样式。
- `apps/web/src/theme/theme.ts` 新建：导出 `resolveTheme`（纯函数，无 DOM 依赖）、`applyThemeToElement`、`initTheme`、`setTheme`。
- `apps/web/src/main.tsx` 更新：`applyStoredTheme` 改为调用 `initTheme()`（冗余保护，index.html inline 已是主引导）。

### C4 — primitives 补齐

新建 `apps/web/src/components/ui/`：

| 文件 | 组件 |
|------|------|
| `input.tsx` | `Input`（CVA 风格 forwardRef） |
| `textarea.tsx` | `Textarea`（CVA 风格 forwardRef） |
| `label.tsx` | `Label`（forwardRef） |
| `card.tsx` | `Card` / `CardHeader` / `CardTitle` / `CardDescription` / `CardContent` / `CardFooter` |
| `badge.tsx` | `Badge`（CVA variants：default / secondary / destructive / success / outline） |
| `skeleton.tsx` | `Skeleton`（animate-pulse + bg-primary/10） |

全部消费语义 Token（border-border、bg-card、text-card-foreground、bg-success 等）。

### C5 — 消费闭环

| 消费点 | 原始类 | 迁移后 |
|--------|--------|--------|
| `renderer/confirm.tsx` backdrop | `bg-black/40` | `bg-overlay` |
| `renderer/confirm.tsx` dialog | `shadow-xl` | `shadow-md` |
| `renderer/modal.tsx` backdrop | `bg-black/40` | `bg-overlay` |
| `renderer/modal.tsx` modal | `shadow-xl` | `shadow-lg` |
| `renderer/render.tsx` FeedbackRegion success | `border-emerald-500/50 bg-emerald-500/10 text-emerald-700` | `border-success/50 bg-success/10 text-success` |

### C6 — 验证

- `npm run test`（vitest run）：**595 tests, 26 test files — 全绿**（含 26 个新主题测试：`resolveTheme` 6 + `applyThemeToElement` 3 + Token 结构断言 17）。
- `npm run build`（tsc -b + vite build）：**exit 0**，dist 产物正常生成（index-*.css 27.10 kB / index-*.js 489.79 kB）。
- 主题逻辑单测覆盖：纯函数路径（dark/light/fallback/OS/unknown stored）+ DOM stub 路径。
- Token 结构断言：CSS 文件包含所有 13 个必需变量；shadow alias 无自引用。

## 产物路径

| 类型 | 路径 |
|------|------|
| Token CSS | `apps/web/src/index.css` |
| 主题模块 | `apps/web/src/theme/theme.ts` |
| 主题单测 | `apps/web/src/theme/theme.test.ts` |
| FOUC 引导 | `apps/web/index.html`（inline script） |
| primitives | `apps/web/src/components/ui/input.tsx` / textarea / label / card / badge / skeleton |
| build log | scratch `build-s1.log` |
