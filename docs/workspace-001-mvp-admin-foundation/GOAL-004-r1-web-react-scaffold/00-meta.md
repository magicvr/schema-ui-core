---
id: GOAL-004-r1-web-react-scaffold
title: R1 · React Web 工程骨架
status: done
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-07-31
version: 0.3.0
---

# GOAL-004 · R1 · React Web 工程骨架

## 概述

在 `apps/web` 建立可本地运行的 React 前端骨架：Vite + TypeScript + npm，并接入 Tailwind CSS 与 shadcn/ui **工具链基线**（空壳主题/示例组件级即可）。  
**不**实现 Admin 业务页与协议 Renderer 全量能力（外壳属 R3；协议域范例属 R5）。

结构可参考平行仓 `allinme.web-client`（`dev`）的 `host` / `protocol` / `renderer` 分层思想；**不**整仓拷贝，**不**以自研 tokens 替代 Charter 要求的 Tailwind/shadcn 方向。

## 成功标准

- [x] `apps/web/package.json` + `package-lock.json`；脚本含 `dev` / `build`
- [x] `npm install` + `npm run dev`（或文档命令）可启动
- [x] TypeScript + Vite + React 19 基线可构建（`npm run build`）
- [x] Tailwind 已接入；**可指回 shadcn 初始化痕迹**（`components.json` 与/或文档记载的 init + `components/ui`）；至少 1 个无业务示例组件
- [x] 浅色/深色切换的**最小占位**存在（完整产品化可延后 R3）
- [x] R1 预建空 `host` / `protocol` / `renderer` 分层目录 + 边界 README（D-002）
- [x] 未把平行仓 mock 业务域（订单等）作为默认路由树；**不含** App manifest 导航壳 / 多业务路由（属 R3）

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-004-001 | non-blocking | shadcn 风格预设（new-york / default）与 CSS 变量策略 | 组件目录生成 | 首次 shadcn init 前 | 选预设并写入决策 | **verified** | — | D-002：new-york + CSS variables（shadcn 默认） |
| I-004-002 | **required** | 是否保留 host/protocol/renderer 包边界于 R1 | 骨架目录冻结 | R1 骨架目录冻结前 | 扁平 vs 预建分层二选一 | **verified** | — | D-002 方案 **(B)**：R1 预建空分层 + README 边界 |

## 父目标

- [GOAL-001-mvp-admin-foundation](../GOAL-001-mvp-admin-foundation/00-meta.md)

## 备注

- 纲领阶段：**R1**。路径 `apps/web` 由 Root D-004 锁定；目录所有权服从 GOAL-002 D-002。
- 平行仓现状：React 19 + Vite 8 + Vitest + oxlint + npm；**无** Tailwind/shadcn —— 本目标需**新增** UI 基线，而非 pure copy。
