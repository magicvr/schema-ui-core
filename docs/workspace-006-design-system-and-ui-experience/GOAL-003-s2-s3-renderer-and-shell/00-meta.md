---
id: GOAL-003-s2-s3-renderer-and-shell
title: S2+S3 · Renderer 视觉接入 + Shell 断点升级
status: done
parent: GOAL-001-design-system-and-ui-experience
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
progress: 2/2
---

# GOAL-003 · S2+S3 · Renderer 视觉接入 + Shell 断点升级

## 概述

承接 Root S2 和 S3 阶段：

- **S2**：Schema Renderer 中残余硬编码色值迁移至语义 Token；chart 改用 `--chart-*` CSS 变量；confirm/modal 已在 S1 完成；checklist 见下。
- **S3**：Shell 移动端汉堡菜单 + 导航抽屉（`bg-overlay` + `shadow-lg`）；关闭事件（navigate/X/backdrop）；结构断言测试。

**实施依据**：Root D-002/D-003/D-004（已 accepted）；S1 Token 体系（GOAL-002 done）；I-PROTO-FULL-001 type 白名单不扩张。

## 成功标准

- [x] **C1（S2）**：chart `pie` 切片改用 `var(--color-chart-N)` CSS 变量代替 HSL 计算字符串；vitest 通过。
- [x] **C2（S3）**：App.tsx 含移动抽屉状态（`mobileDrawerOpen`）；汉堡按钮 aria-label；X 关闭；backdrop `bg-overlay`；navigate 时关闭；shell.test.ts 结构断言通过。

## 派生进度

`progress: 2/2` 由上方两个等权检查点派生；仅为展示，不放行阶段或关门。

## 父目标

[GOAL-001-design-system-and-ui-experience](../GOAL-001-design-system-and-ui-experience/00-meta.md)（Root；本目标为 S2+S3 阶段子目标）

## 台账布局

使用 ledger 目录：`01-decision/`、`02-execution/`、`03-audit/`。
