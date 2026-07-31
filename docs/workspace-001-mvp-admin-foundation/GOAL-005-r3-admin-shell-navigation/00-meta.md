---
id: GOAL-005-r3-admin-shell-navigation
title: R3 · Admin 外壳与导航
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-07-31
version: 0.1.0
---

# GOAL-005 · R3 · Admin 外壳与导航

## 概述

在 `apps/web` 的 R1 React 工程骨架之上，规划并在后续实施阶段交付由 App manifest 驱动的 Admin 外壳、导航入口与路由边界。本目标当前只完成立项、范围登记和高层路线图；manifest loader、导航壳和 router 尚未实现。

范围依据：Root 已将 R3 定义为“Admin 外壳与导航”；协议资料将 `D-APP` 映射为 React 侧的“装载与导航壳”，并固定 `S-09`、App Manifest schema 及 `app-manifest` / `app-navigation` fixture 的上游路径。固定资料版本为 `schema-ui-docs` artifact `2.7.0`，source commit 为 `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b`。

## 范围边界

### 纳入

- App manifest 的来源、最小可用子集、装载入口和失败边界的决策与实现计划。
- Admin shell 的页面级布局、导航容器和应用入口边界。
- manifest 到导航项/路由入口的映射，以及默认路由、fallback、active-route 语义的决策与实现计划。
- 面向 `app-manifest` / `app-navigation` 的结构或行为验证路径；具体证据须在实施阶段产生。

### 排除

- R4 的账号、权限模型与权限继承实现；父目标的 `I-PROTO-002` 仍为 R4 的开放 required 门禁。
- R5 的协议 Renderer 全量、业务域范例页和并行 `mock-app` 业务演示；父目标的 `I-PROTO-003` 仍为 R5 的开放 required 门禁。
- “完整协议支持”或完整 conformance 主张；上游 excluded reference/runner 不能单独作为兼容证明。
- 完整主题/设计系统产品化，以及与 R3 无关的业务页面实现。

## 高层路线图（P-001）

1. **契约发现与信息就绪**：解析 manifest schema、导航 fixture 和 shell 产品边界，关闭或按 P-004 记录影响方案/实施的 required 信息项。
2. **方案冻结**：记录 manifest 最小子集、路由映射、默认/fallback/active-route 语义和 shell 边界；未完成前不进入实现。
3. **R3 实施**：在 `apps/web` 落地 manifest 装载、Admin shell 和导航路由入口，保持 R4/R5 边界。
4. **验证与关门**：使用可核对的结构/行为/运行时证据复核成功标准，完成阶段自审和 required finding 闭合后再申请关门。

以上检查点均为计划，当前没有 R3 实施完成事实。

## 成功标准

- [ ] R3 必需信息项已由证据验证，或有用户书面接受的有界 residual；未知项没有被默认为已知。
- [ ] manifest 装载使用已冻结的 schema/版本/最小子集，并有可核对的无效输入或装载失败处理边界。
- [ ] Admin shell 能由已冻结的导航数据进入页面，默认路由、fallback 和 active-route 语义可通过测试或运行时证据复核。
- [ ] `app-manifest` / `app-navigation` 的验证路径已执行并记录结果；未用 excluded reference/runner 替代语义验证。
- [ ] R4 权限、R5 Renderer/范例和完整协议支持仍保持边界外；目标关门前无开放 required finding。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-005-001 | required | App manifest 的最小 schema、来源、版本和本地校验入口是什么？ | 方案冻结 / 实施 | 方案冻结前 | 对照 `S-09`、`docs/schemas/app-manifest.schema.json` 和固定 artifact；确认本仓库接入方式 | open | 不延期；方案冻结前复核 | 上游路径和版本已登记；本地接入仍待确认 |
| I-005-002 | required | manifest 条目如何映射为导航项、路由入口和层级关系？ | 方案冻结 / 实施 | 方案冻结前 | 对照 `S-09` 与 `conformance/fixtures/app-navigation/`，形成映射决策并记录反例 | open | 不延期；方案冻结前复核 | 当前未找到本地映射实现；待确认 |
| I-005-003 | required | 默认路由、未知路由和 fallback 页面/行为分别是什么？ | 方案冻结 / 实施 / 验收 | 实施前 | 从导航 fixture 和 shell 边界中收集规则，补充路由测试案例 | open | 不延期；实施前复核 | 当前未发现 fallback 契约；待确认 |
| I-005-004 | required | active-route 的来源、匹配优先级、重定向和 URL 语义是什么？ | 方案冻结 / 实施 / 验收 | 实施前 | 设计并验证 URL/route 状态矩阵，记录用户可见行为 | open | 不延期；实施前复核 | 当前 `main.tsx` 无 router；待确认 |
| I-005-005 | required | Admin shell 的产品边界包含哪些固定区域，哪些留给业务页或后续产品化？ | 目标方案 / 实施 | 方案冻结前 | 对照 Charter、VP、R3 边界并取得用户确认；将未决取舍写入决策 | open | 不延期；方案冻结前复核 | 当前只能确认“Admin 外壳 + 导航壳”；具体边界待确认 |

父目标的 `I-PROTO-001=verified` 仅表示 R2 冻结范围；不解除本目标的未知项，也不代表 R3-R5 已实施或完整协议已支持。

## 父目标

- [GOAL-001-mvp-admin-foundation](../GOAL-001-mvp-admin-foundation/00-meta.md)

## 备注

本目标的 `status: active` 表示 R3 已进入规划阶段，不表示代码实现、验收或关门已完成。
