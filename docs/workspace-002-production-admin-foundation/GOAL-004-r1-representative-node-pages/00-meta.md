---
id: GOAL-004-r1-representative-node-pages
title: R1 · 代表性 Node 页面与回归证据
status: active
created: 2026-08-01
updated: 2026-08-01
parent: GOAL-001-production-admin-foundation
version: 0.1.0
progress: 0/5
---

# GOAL-004 · R1 · 代表性 Node 页面与回归证据

## 概述

在 `v0.1.3` 白名单内，交付至少一组**可加载的 Node 树页面**（列表、表单、组合/详情类），并留下自动化或可重复回归证据，证明「改 Schema 即可出现页面」在主路径上成立。页面消费依赖 GOAL-002/003；本目标侧重**页面资产 + 回归**，不重做加载器或路由策略。

依据 Root **D-004** / **D-005** 与 VP-002 阶段 1 验收建议。

## 成功标准

- [ ] 至少 **1** 个列表向 Node 页面（table 或 search+table 结构）可经 schemaUrl 主路径渲染
- [ ] 至少 **1** 个表单向 Node 页面（含白名单控件；可选 `$context` reaction）可经主路径渲染
- [ ] 至少 **1** 个组合或详情向页面（section/grid/tabs + recordView/text 等）可经主路径渲染
- [ ] 未知节点类型或非法页面在代表性路径上 fail-closed 且错误可观察
- [ ] 有自动化测试或等价可重复脚本覆盖上述成功/失败路径中的关键断言

## 派生进度

`progress: 0/5` 由上方 5 条成功标准等权派生。

## 信息需求

| ID | 问题 / 所需信息 | 级别 | 影响门禁 | 最晚阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据或结论 |
|----|-----------------|------|----------|----------|-----------------|------|-------------|------------|
| `I-004-001` | 代表性页面是否改写现有 5 个示例语义，还是新增独立 schema 资源？ | required | 页面资产布局 | 实施开始前 | 对照 EXAMPLE_PAGES 与演示需要；记决策 | open | 不延期 | 倾向新增 schema 资源并保留示例兼容 |
| `I-004-002` | 列表数据是否继续用 `/api/records` 演示 API？ | non-blocking | 列表可观察 | 列表页验收前 | 复用既有 records 路径即可 | open | 可与实施并行 | I-001 允许 R1 reuse |

## 依赖与边界

| 项 | 说明 |
|----|------|
| 依赖 | GOAL-002 加载/校验；GOAL-003 默认主路径（完整「新增只改 Schema」证明需要 003） |
| 可并行 | 可在 002/003 完成前起草 Node JSON 与单测用 `RenderPage` 直渲 |
| In | 代表性 Node 文档、失败路径、回归证据 |
| Out | 真实认证、持久化 IAM、R4 CRUD 实体定稿、覆盖表扩域、D-UPLOAD |

## 父目标

- [GOAL-001-production-admin-foundation](../GOAL-001-production-admin-foundation/00-meta.md)
