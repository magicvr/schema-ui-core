---
id: GOAL-002-s1-design-tokens-and-primitives
title: S1 · Design Token / 主题 / shadcn primitives
status: done
parent: GOAL-001-design-system-and-ui-experience
created: 2026-08-09
updated: 2026-08-09
version: 0.2.0
progress: 6/6
---

# GOAL-002 · S1 · Design Token / 主题 / shadcn primitives

## 概述

承接 Root [GOAL-001-design-system-and-ui-experience](../GOAL-001-design-system-and-ui-experience/00-meta.md) 的 **S1 阶段**：在现有 shadcn + Tailwind v4 CSS 变量体系上扩展语义化 Token（Color / Typography / Radius / Shadow / Spacing）、深/浅色主题切换（关键壳层无持续 FOUC）、shadcn 风格 primitives 可发现，并完成 **F-002**（`--elevation-*` → `@theme --shadow-*` 无自引用映射）的实施证据闭合。

**方案依据（已 accepted 的 Root 决策）**：D-002（Token 语义分层与命名）、D-003（合并响应三审：elevation→shadow 映射、双层纪律、消费闭环、字阶路径）、D-004（视觉方向冻结，Stitch 定稿）。本目标不重新决策，只实施与验证。

**范围纪律**：不引入第二套 Token 真相源（无 tokens.json / Style Dictionary / `ds.*`）；不扩张 `I-PROTO-FULL-001` disposition；不新增 registry type。

## 成功标准（可验收 · 等权检查点 · 共 6 项）

- [x] **C1**：`index.css` 按区块组织（Color / Typography / Elevation·Shadow / Radius / 引导注释），增量含 `--destructive`、`--success`、`--chart-1…5`、`--overlay`、`--font-sans`、`--font-mono`、`--elevation-sm|md|lg`（`:root`/`.dark` 双态，深色降不透明度）——无同名自引用（F-001/F-005）。
- [x] **C2**：**F-002 映射**：`@theme inline` 中 `--shadow-sm|md|lg: var(--elevation-*)`，`shadow-sm|md|lg` utility 可生成；深/浅色下取值符合预期；confirm/modal 至少一处从 `shadow-xl` 迁移到语义 shadow（D-003 §2 闭合条件 a–d）。
- [x] **C3**：**FOUC 治理**：`index.html` 同步内联引导脚本（读 `localStorage.theme` / `prefers-color-scheme`），首帧即带 `dark` class；`color-scheme` 属性随主题设置；主题应用逻辑抽成可单测纯单元（`src/theme/`）。
- [x] **C4**：**primitives**：shadcn 风格 `ui/*` 原语补齐（Input / Textarea / Label / Card / Badge / Skeleton 等，CVA + `cn`），全部消费语义 Token，可发现（`components.json` 保持指向 `src/index.css`）。
- [x] **C5**：**消费闭环（F-003）**：success 反馈迁到 `--success` Token；chart 消费 `--chart-1…5`；overlay 用语义 `--overlay`；shadow 消费语义 utility；不进 Token 的类别写边界理由。
- [x] **C6**：**验证**：vitest 全绿（含主题纯逻辑单测、Token 结构断言）；`npm run build` 通过；headless 启动浅/深两屏截图无持续 FOUC、壳层渲染；暗色主文案可读性不低于定稿 Overview dark 下限。

## 派生进度展示

`progress: X/6` 由上方 6 个等权检查点派生；仅为展示，不放行阶段、不关闭 finding（**F-002 闭合以 Root 台账为准**）。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | S1 实施输入（Token 约定、F-002 映射、视觉方向）是否齐备 | C1–C6 | 实施前 | 读 Root D-002/D-003/D-004 + I-S1-001 | **closed** | — | Root 决策链已 accepted；I-001/I-002/I-005 closed（Root 00-meta） |
| I-002 | required | 深浅色 headless 验证环境可用性 | C6 | C6 前 | Playwright chromium 安装检查 | **closed** | — | `%LOCALAPPDATA%\ms-playwright\chromium-1234` 已装（2026-08-09） |

## 父目标

- [GOAL-001-design-system-and-ui-experience](../GOAL-001-design-system-and-ui-experience/00-meta.md)（Root；本目标为 S1 阶段子目标）

## 台账布局

本目标使用 ledger 目录：`01-decision/`、`02-execution/`、`03-audit/`。索引文件保留 frontmatter 与条目表；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*`。
