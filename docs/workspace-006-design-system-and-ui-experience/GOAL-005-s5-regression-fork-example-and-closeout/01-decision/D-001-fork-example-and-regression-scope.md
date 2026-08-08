---
id: D-001-fork-example-and-regression-scope
title: fork Token 示例形态与回归核对范围
date: 2026-08-09
status: accepted
parent: GOAL-005-s5-regression-fork-example-and-closeout
---

# D-001 · fork Token 示例形态与回归核对范围

## 决策

1. **fork 示例形态**：新增 `apps/web/src/theme/brand.example.css` + `apps/web/src/theme/README.md`。
   示例只覆盖 `src/index.css` 已声明的语义 Token（`--primary`/`--chart-1..5`/`--radius` 等），
   不新建平行主题系统、不新增 Token 命名。fork 复制该文件、改值、在 `main.tsx` 里于
   `index.css` **之后** 多 import 一行即可换品牌色/图表色/圆角，不改任何组件代码
   （因为 S1–S4 已完成 Token 消费闭环）。
2. **可复核性**：新增结构测试 `brand-example.test.ts`，用正则解析 `index.css` 与
   `brand.example.css` 两份文件声明的 `--custom-property` 名称集合，断言示例集合是
   基线集合的**严格子集**（不会发明基线之外的 Token 名），且至少覆盖 primary/chart/radius
   三个品牌相关族。这把"示例是否会漂移出既有 Token 系统"变成一个可自动核对的事实，
   不依赖人工复查。
3. **回归核对范围**：`apps/web` vitest 全量 + `npm run build`（gating，plan 步骤 2/3）；
   额外探测并实际运行 Playwright e2e（`npx playwright --version` 确认已装 → 直接
   `npm run test:e2e` 而非仅记录降级）。本沙箱环境可运行 Chromium headless，
   两条 e2e spec（`schema-crud.spec.ts`/`shell.spec.ts`）均通过，作为比"诚实退化"更强的证据。

## 未选方案

- 新建独立的 `ThemeProvider`/JS 主题配置对象来表达 fork 品牌：与既有纯 CSS 变量机制重复，
  引入运行时开销和第二套真相源，放弃。
- 把示例做成可执行的 CLI 脚手架命令（如 `npm run fork:brand`）：超出"最小示例"范围，
  且需要额外的脚本维护面，放弃；一份可复制的 CSS 文件 + README 已满足验收要求。
