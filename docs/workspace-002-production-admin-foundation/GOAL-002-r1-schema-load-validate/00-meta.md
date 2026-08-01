---
id: GOAL-002-r1-schema-load-validate
title: R1 · Schema 加载、校验与统一错误面
status: active
created: 2026-08-01
updated: 2026-08-01
parent: GOAL-001-production-admin-foundation
version: 0.1.0
progress: 0/4
---

# GOAL-002 · R1 · Schema 加载、校验与统一错误面

## 概述

在不扩大 `I-PROTO-001 v0.1.3` 的前提下，交付运行时 **Schema 页面文档加载**、**加载时结构校验（D-VAL）** 与 **统一可观察错误面**，使 manifest 中的 `page.schemaUrl` 从「仅建模」变为可消费管线。本目标不切换 App 默认渲染分支（见 GOAL-003），也不交付业务代表性页面全集（见 GOAL-004）。

依据 Root **D-004**（I-001 矩阵 §4 In）与 **D-005**（R1 子目标拆分）。

## 成功标准

- [ ] 存在可调用的页面 Schema 加载入口：按已解析的 `schemaUrl`（含路由参数模板展开）获取 page 文档
- [ ] 加载路径强制结构校验（至少 page/node 或等价 Ajv 入口）；非法文档 fail-closed，不进入渲染
- [ ] 无效 Schema、网络/解析失败、校验失败暴露统一错误结构（可测、可在 UI 或宿主层观察）
- [ ] 有自动化测试覆盖：成功加载样例、校验失败、未知/非法 body 的拒绝路径

## 派生进度

`progress: 0/4` 由上方 4 条成功标准等权派生；不放行 R1 Root 检查点。

## 信息需求

| ID | 问题 / 所需信息 | 级别 | 影响门禁 | 最晚阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据或结论 |
|----|-----------------|------|----------|----------|-----------------|------|-------------|------------|
| `I-002-001` | Schema 静态资源放置策略（Vite public / 构建内嵌 / mock fetch）最小方案？ | required | 实施路径选择 | 实施开始前 | 对照现有 manifest `schemaUrl` 与 dev server；记决策 | open | 不延期 | 待实施时定 |
| `I-002-002` | 错误码/错误包络是否与现有 ManifestError / RenderError 对齐？ | non-blocking | 错误面一致性 | 验收前 | 对照 `app-manifest` 与 `render` 错误类型 | open | 验收前复核 | 倾向复用既有 code 风格 |

## 依赖与边界

| 项 | 说明 |
|----|------|
| 父目标 | `GOAL-001-production-admin-foundation` |
| 同阶段 | 可与 GOAL-003 / GOAL-004 规划并行；**实施上** GOAL-003 默认依赖本目标加载器可用 |
| In | 加载器、D-VAL 串联、统一错误、测试 |
| Out | App 默认分支切换、代表性业务页全集、真实认证、覆盖表扩域 |

## 父目标

- [GOAL-001-production-admin-foundation](../GOAL-001-production-admin-foundation/00-meta.md)
