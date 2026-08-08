---
id: E-001-s2-s3-implementation
title: S2+S3 实施 — Renderer chart Token 化 + Shell 移动抽屉
date: 2026-08-09
status: done
parent: GOAL-003-s2-s3-renderer-and-shell
---

# E-001 · S2+S3 实施完成记录

## 事实摘要

**日期**：2026-08-09  
**scope**：S2 Renderer 视觉接入 + S3 Shell 移动端升级

### S2 — Renderer Token 化

| 消费点 | 原始值 | 迁移后 |
|--------|--------|--------|
| `render.tsx` ChartView pie 切片颜色 | `hsl(${(index * 137.5) % 360} 65% 55%)` | `var(--color-chart-${(index % 5) + 1})` |

confirm/modal 的 `bg-overlay` / `shadow-md|lg` 迁移已在 S1（GOAL-002）完成。  
FeedbackRegion success 颜色迁移已在 S1 完成。

### S3 — Shell 移动端抽屉

- `App.tsx`：新增 `mobileDrawerOpen` state（默认 false）。
- 汉堡按钮：`aria-label="Open navigation menu"`，仅 `lg:hidden` 显示。
- 导航抽屉（固定定位，`z-40`）：包含关闭按钮 `aria-label="Close navigation menu"`。
- 背景遮罩：`className="fixed inset-0 z-30 bg-overlay lg:hidden"`（语义 Token）。
- 抽屉面板：`shadow-lg`（语义 elevation Token）。
- 导航后关闭：`onNavigate` 中调用 `setMobileDrawerOpen(false)`。
- 移除旧的 `overflow-x-auto` 水平滚动导航条。
- 移除已无调用者的 `flattenNavigation` 函数（TypeScript unused-var）。

### 新增测试

`apps/web/src/app/shell.test.ts`：

- 6 个移动抽屉状态机纯逻辑测试（open/close/navigate/idempotent）。
- 6 个 App.tsx 结构断言（aria-label / bg-overlay / shadow-lg / mobileDrawerOpen）。

### 验证结果

- vitest run：**607 tests / 27 test files — 全绿**（含 12 个新 shell.test.ts 测试）。
- `npm run build`：**exit 0**，dist 产物正常（index-*.css 27.22 kB / index-*.js 491.34 kB）。
