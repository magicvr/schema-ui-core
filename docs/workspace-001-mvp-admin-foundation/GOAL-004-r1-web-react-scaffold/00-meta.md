---
id: GOAL-004-r1-web-react-scaffold
title: R1 · React Web 工程骨架
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-07-31
version: 0.1.0
---

# GOAL-004 · R1 · React Web 工程骨架

## 概述

在 `apps/web` 建立可本地运行的 React 前端骨架：Vite + TypeScript + npm，并接入 Tailwind CSS 与 shadcn/ui **工具链基线**（空壳主题/示例组件级即可）。  
**不**实现 Admin 业务页与协议 Renderer 全量能力（外壳属 R3；协议域范例属 R5）。

结构可参考平行仓 `allinme.web-client`（`dev`）的 `host` / `protocol` / `renderer` 分层思想；**不**整仓拷贝，**不**以自研 tokens 替代 Charter 要求的 Tailwind/shadcn 方向。

## 成功标准

- [ ] `apps/web/package.json` + `package-lock.json`；脚本含 `dev` / `build`（及既有测试占位可选）
- [ ] `npm install` + `npm run dev`（或文档命令）可启动
- [ ] TypeScript + Vite + React 19 基线可构建（`npm run build`）
- [ ] Tailwind 已接入；shadcn/ui 初始化或等价组件目录约定已落盘（至少 1 个无业务示例）
- [ ] 浅色/深色切换的**最小占位**存在（完整产品化可延后 R3）
- [ ] 未把平行仓 mock 业务域（订单等）作为默认路由树

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-004-001 | non-blocking | shadcn 风格预设（new-york / default）与 CSS 变量策略 | 组件目录生成 | 首次 shadcn init 前 | 选预设并写入决策 | open | — | 默认倾向 new-york；实施时确认 |
| I-004-002 | non-blocking | 是否保留 host/protocol/renderer 包边界于 R1 空壳 | 目录深度 | R1 骨架冻结前 | 最小 `src/app` 即可 vs 预建分层 | open | — | 建议 R1 预建空分层目录防后续大挪 |

## 父目标

- [GOAL-001-mvp-admin-foundation](../GOAL-001-mvp-admin-foundation/00-meta.md)

## 备注

- 纲领阶段：**R1**。路径 `apps/web` 由 Root D-004 锁定。
- 平行仓现状：React 19 + Vite 8 + Vitest + oxlint + npm；**无** Tailwind/shadcn —— 本目标需**新增** UI 基线，而非 pure copy。
