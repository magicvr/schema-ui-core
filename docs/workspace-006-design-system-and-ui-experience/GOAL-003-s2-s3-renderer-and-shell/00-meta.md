---
id: GOAL-003-s2-s3-renderer-and-shell
title: S2+S3 · Renderer 视觉接入 + Shell 断点升级
status: active
parent: GOAL-001-design-system-and-ui-experience
created: 2026-08-09
updated: 2026-08-09
version: 0.2.0
progress: 0/2
---

# GOAL-003 · S2+S3 · Renderer 视觉接入 + Shell 断点升级

## 概述

承接 Root S2 和 S3 阶段，**成功标准必须以 Root D-004 / VP-005 exit 2–3 为分母**，不得缩成「换 CSS 变量 / 只加移动抽屉」。

### 已发生的局部实现（事实，非本目标完成证明）

| 子项 | 状态 | 证据 |
|------|------|------|
| chart pie 用 `--chart-*` | 已实现 | `render.tsx`；commit `c2d7b60` |
| confirm/modal overlay/shadow Token | 已在 S1 | GOAL-002 |
| 移动汉堡 + 导航抽屉 | 已实现 | `App.tsx`；`shell.test.ts` |

### 未达标（A-006 / A-002）

| 子项 | 状态 |
|------|------|
| 桌面密表视觉 + **移动卡片列表** | **未做** |
| `recordView` 右栏 Drawer / 移动 Sheet | **未做** |
| 登录 / Overview 壳气质对齐 Stitch | **未做** |
| 表单控件 / 表格主路径消费 design primitives | **未做** |

**实施依据**：Root D-002/D-003/D-004（已 accepted）；S1 Token 体系（GOAL-002 done）；I-PROTO-FULL-001 type 白名单不扩张。

## 成功标准（重写 · 对齐 D-004 · 等权 2 项）

- [ ] **C1（S2 · Renderer 视觉）**：钉死 type 面至少完成 D-004 优先级中的 **列表双端**（桌面密表 + 移动卡片列表）与 **`recordView` Drawer/Sheet 呈现**；表单/登录/展示面有可观察升级或明确分批证据；**禁止**仅以 chart Token 勾选。vitest/build 不回退。
- [ ] **C2（S3 · Shell）**：壳层与登录等对照 Stitch Overview / Sign-in 可复核升级；移动汉堡抽屉可保留为已交付子能力，但 **不足**单独完成 C2；Dialog/Toast 语言与壳气质一致。

> 历史过窄检查点（「仅 pie 颜色」「仅移动抽屉」）已作废，见 `03-audit/A-002-under-delivery-vs-d004.md`。

## 派生进度

`progress: 0/2` 由上方两个等权检查点派生；仅为展示，不放行阶段或关门。

## 父目标

[GOAL-001-design-system-and-ui-experience](../GOAL-001-design-system-and-ui-experience/00-meta.md)（Root；本目标为 S2+S3 阶段子目标；Root 开放 **F-VUI-001 / F-VUI-002**）

## 台账布局

使用 ledger 目录：`01-decision/`、`02-execution/`、`03-audit/`。
