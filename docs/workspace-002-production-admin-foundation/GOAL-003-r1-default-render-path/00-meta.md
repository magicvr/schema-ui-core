---
id: GOAL-003-r1-default-render-path
title: R1 · 默认 Renderer 主路径与示例降级
status: active
created: 2026-08-01
updated: 2026-08-01
parent: GOAL-001-production-admin-foundation
version: 0.3.0
progress: 3/4
---

# GOAL-003 · R1 · 默认 Renderer 主路径与示例降级

## 概述

将 Admin 匹配路由后的**默认页面能力**切换为 Schema 驱动的 `RenderPage`（或等价宿主），并把 5 个手写 `EXAMPLE_PAGES` 迁移为 Schema 文档（`I-003-001` 已定策略：迁移，见 D-003）。依赖 GOAL-002 的加载与校验管线；代表性 Node 页内容与回归证据由 GOAL-004 补齐。

依据 Root **D-004** / **D-005**。

## 成功标准

- [x] 匹配 manifest 路由后，**默认**走 Schema 加载 → 校验 → `RenderPage`（不再以 `EXAMPLE_PAGES` 为默认分支）
- [ ] 既有 5 个手写示例**迁移为 Schema 文档**，经默认 Schema 主路径渲染；应用内不再存在手写示例作为独立页面路径（D-003）
- [x] 非示例页不再展示「renderer remains a later protocol boundary」类占位作为主交付面
- [x] 有自动化测试：默认路径渲染 Schema 页（含迁移后示例页）；缺失/非法 Schema 时统一错误面可预期（fail-closed）

## 派生进度

`progress: 3/4` 由上方 4 条成功标准等权派生；标准 2 待 GOAL-004 落地 5 份迁移 Schema 文档并经默认主路径渲染后闭合（本目标已移除手写默认分支）。

## 信息需求

| ID | 问题 / 所需信息 | 级别 | 影响门禁 | 最晚阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据或结论 |
|----|-----------------|------|----------|----------|-----------------|------|-------------|------------|
| `I-003-001` | 既有 5 个 EXAMPLE_PAGES 在默认切换后如何保留（双分支 / 仅测试 / 迁移为 Schema）？ | required | 降级策略 | 实施切换前 | 对照 registry 与产品演示需求；记决策 | closed | 不延期 | **迁移为 Schema**（D-003）：语义改写为 Schema 文档经默认主路径渲染；显式入口 = schemaUrl 链（route → schemaUrl → 加载+校验 → RenderPage）；5 份文档由 GOAL-004 落地 |
| `I-003-002` | table 数据注入（现由示例拥有）在默认路径如何挂接？ | required | 列表页可渲染 | 默认路径验收前 | 与 GOAL-002/004 对齐；可先静态/既有 records API | open | 与 004 联调复核 | 见 I-001 矩阵 D-DATA 部分 |

## 依赖与边界

| 项 | 说明 |
|----|------|
| 依赖 | **GOAL-002** 加载器与校验可用（硬依赖） |
| 协作 | GOAL-004 提供可演示 Node 文档与回归 |
| In | App 路由默认分支、`RenderPage` 宿主、示例降级声明与测试 |
| Out | 加载器本体、覆盖扩域、真实认证、完整业务 CRUD 验收 |

## 父目标

- [GOAL-001-production-admin-foundation](../GOAL-001-production-admin-foundation/00-meta.md)
