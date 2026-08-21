---
id: GOAL-003-s2-s3-renderer-and-shell
title: S2+S3 · Renderer 视觉接入 + Shell 断点升级
status: done
parent: GOAL-001-design-system-and-ui-experience
created: 2026-08-09
updated: 2026-08-09
version: 0.3.0
progress: 2/2
---

# GOAL-003 · S2+S3 · Renderer 视觉接入 + Shell 断点升级

## 概述

承接 Root S2 和 S3 阶段，**成功标准以 Root D-004 / VP-005 exit 2–3 为分母**。

### 交付事实（E-002）

| 子项 | 状态 | 证据 |
|------|------|------|
| 桌面密表 + 移动卡片列表 | **done** | `data-table.tsx` dual-end；commit `f16dc9f` |
| `recordView` Drawer / Sheet | **done** | `render.tsx` RecordView；selection 路径测试 |
| 表单 / 展示 primitives | **done** | form-controls Input/Label/Textarea；StatCard Card |
| 壳 topbar + ~256 sidenav + 登录 | **done** | `App.tsx` / `LoginPage.tsx`；commit `5716df9` |
| chart Token / 移动汉堡抽屉 | 保留子项 | E-001 历史 |

## 成功标准（对齐 D-004 · 等权 2 项）

- [x] **C1（S2 · Renderer 视觉）**：列表双端 + recordView Drawer/Sheet；表单/登录/展示可观察升级；非 chart-only。vitest/build 不回退。
- [x] **C2（S3 · Shell）**：壳层与登录对照 D-004 可复核升级；移动汉堡仅为子能力。

## 派生进度

`progress: 2/2`。开放 required = 0（F-003-001 fixed；A-003/A-004/A-005）。

## 父目标

[GOAL-001-design-system-and-ui-experience](../GOAL-001-design-system-and-ui-experience/00-meta.md)

## 台账布局

使用 ledger 目录：`01-decision/`、`02-execution/`、`03-audit/`。
