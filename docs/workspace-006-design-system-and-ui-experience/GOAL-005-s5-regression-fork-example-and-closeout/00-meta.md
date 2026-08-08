---
id: GOAL-005-s5-regression-fork-example-and-closeout
title: S5 · 视觉回归 + fork Token 示例 + 过程关门
status: done
parent: GOAL-001-design-system-and-ui-experience
created: 2026-08-09
updated: 2026-08-09
version: 0.2.0
progress: 2/2
---

# GOAL-005 · S5 · 视觉回归 + fork Token 示例 + 过程关门

## 概述

承接 Root [GOAL-001-design-system-and-ui-experience](../GOAL-001-design-system-and-ui-experience/00-meta.md) 的 **S5 阶段**：确认既有回归套件（vitest 全量 + build + e2e）在 S4/S5 改动后保持全绿不回退；新增一个可复核的 fork 品牌定制最小示例（改 Token/主题变量即可换品牌色，非新建独立设计系统包）；随后驱动 Root 走向"证据齐备的可关门状态"（关门本身须用户书面确认，本目标不单方面写 `status: done`）。

**实施依据**：Root D-002/D-003/D-004（Token/视觉方向已 accepted）；GOAL-004（S4）已完成。

## 方案

- 新增 `apps/web/src/theme/brand.example.css`：只覆盖 `index.css` 已声明的语义 Token（`--primary`/`--chart-1..5`/`--radius`），附使用说明注释；不引入第二套 Token 真相源、不新建独立 npm 包。
- 新增结构测试 `apps/web/src/theme/brand-example.test.ts`：解析 `index.css` 与示例文件的 `--custom-property` 声明集合，断言示例是基线声明的严格子集（防止示例文件漂移出发明新 token 名）。
- 新增 `apps/web/src/theme/README.md`：说明 fork 如何复制并使用该示例。
- 全量回归：`npm run test`（vitest）、`npm run build`、`npx playwright test`（e2e，若沙箱可运行浏览器）。

## 成功标准（可验收 · 等权检查点 · 共 2 项）

- [x] **C1**：fork 品牌定制最小示例存在（`brand.example.css` + `brand-example.test.ts` + `README.md`），示例只覆盖基线已声明的 Token 名（结构测试通过），可复核。
- [x] **C2**：既有回归套件在 S4+S5 改动后保持全绿不回退：vitest 全量通过、`npm run build` exit 0、Playwright e2e（`schema-crud.spec.ts` + `shell.spec.ts`）2/2 通过（本次沙箱可运行浏览器，非诚实退化证据）。

## 派生进度展示

`progress: 2/2` 由上方 2 个等权检查点派生；仅为展示，不放行阶段、不关闭 finding、不推导 Root `status: done`。

## 独立交叉审计与关门提案

**独立交叉审计已完成并落盘**（`source: independent`，grok build / grok-4.5 / 高思考强度，覆盖 GOAL-004 S4 + 本目标 S5 合并 scope）：`03-audit/A-002-independent-cross-audit-s4-s5.md`，`verdict: conditional`，2 条 required finding。编排已响应（`03-audit/A-003-response-a002.md`）：2 条 required finding 均 `fixed`（本目标五件套补齐 + `00-meta.md` 改写为如实描述；`useDisplayData` 错误不自清也已在代码中修复）；3 条 non-blocking finding 中 1 条 `fixed`（补测试）、2 条 `accepted-residual`（范围明确、不阻断本阶段）。**开放 required findings（本目标 scope）= 0。**

**关门提案（P-004 · 待用户裁决）**：在完成上述响应之后，Root S1–S5 五个阶段检查点均已完成（见 Root `00-meta.md` 同步更新为 `progress: 5/5`）；Root `03-audit.md` 台账开放 required findings = 0（F-002 已 fixed）。**Root `status` 是否置为 `done` 仍须用户书面确认**；本目标与本次实施不单方面代为裁决，只负责把关门所需的证据（五件套、回归、独立审计响应）备齐并如实记录。

## 信息就绪与未知项

无独立信息项；继承 Root I-001~I-005（除 I-004 外均 closed；I-004 为 non-blocking 且默认否，不影响本阶段）。

## 父目标

[GOAL-001-design-system-and-ui-experience](../GOAL-001-design-system-and-ui-experience/00-meta.md)（Root；本目标为 S5 阶段子目标）

## 台账布局

本目标使用 ledger 目录：`01-decision/`、`02-execution/`、`03-audit/`。索引文件保留 frontmatter 与条目表；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*`。
