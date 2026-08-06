---
doc_type: vision-charter
vision_id: schema-ui-core-admin-foundation
title: Schema UI Core 中型项目 Admin 基架
status: active
version: 0.2.0
effective_date: 2026-08-04
primary_workspace: workspace-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-08-06
parent: null
---

# Charter · Schema UI Core 中型项目 Admin 基架

## 目的陈述

以 `magicvr/schema-ui-docs` 所定义的前后端协议为兼容边界，构建一个面向中型项目、可被后续项目 fork 的基础 Admin 框架：前端采用 React，后端采用 Go，并让协议驱动的页面、数据与交互能力有可运行、可验证的实现路径。基架最终以单一代码主线、薄内核和可组合模块承载不同 fork 起点，避免用长期平行代码线交换短期裁剪便利。

## 协议来源

| 项 | 值 |
|----|----|
| canonical source | https://github.com/magicvr/schema-ui-docs |
| release ref | `v2.7.0` |
| pinned commit | `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b` |
| manifest | https://raw.githubusercontent.com/magicvr/schema-ui-docs/ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b/protocol-manifest.json |

该外部协议是语义、结构与行为契约的来源。本仓库当前未 vendor 该协议全文；**本地实施清单与前后端映射**已提取于 [protocol-inventory-v2.7.0.md](protocol-inventory-v2.7.0.md)（`F-V001` → `fixed`）。MVP 覆盖子集已在 workspace-001 冻结；任何覆盖扩张与实现核验仍由对应工作区内 **`/govern`** 推进。

## 方向级成功边界

在本愿景仍为 `active` 的前提下，下列方向成立即视为仍在愿景内；它们不是可关门的执行 checklist：

1. 提供可 fork 的 React 前端与 Go 后端 Admin 基架，并对 `schema-ui-docs` `v2.7.0` 的协议能力形成可验证的兼容实现与示例路径。
2. MVP 覆盖最核心的账号与权限能力；每一纳入范围的协议功能均有范例页面和对应的验证路径。
3. 前端经产品化后可被 fork 项目直接使用，采用 Tailwind CSS 与 shadcn/ui 风格组件，支持浅色和深色模式，并以 Linear 与 Vercel Dashboard 的克制、工作导向体验为参考。
4. 以单一代码主线、薄内核、框架无关模块契约和启动时 Profile 提供不同 fork 起点；MVP 与完整 Admin 是同一架构的配置形态，不维护长期平行演进代码线。
5. 后端聚合已启用模块的 Manifest、Schema、导航、权限与数据生命周期贡献；同一前端 build 能随 Profile 组合标准模块，增减模块不要求修改 Renderer 或 Shell 的中央注册路径。

## 非目标

本 Charter 不要求、也不把下列事项写成愿景成功条件：

- 不建设特定业务领域的终端产品；钱包、订单、类目、通知等属于后续 VP 的候选能力。
- 不在本项目内重新定义或替代 `schema-ui-docs` 的协议语义；协议变更应回到上游契约或形成明确的兼容决策。
- 不建设运行时插件市场、远程模块下载、`.so` 加载或运行中热插拔；Profile 只在已编译候选集中选择模块。
- 不承诺 Profile 从二进制中物理移除未启用模块；需要物理裁剪时由 fork 或独立构建目标负责。

## 原则摘要

- 协议优先：固定上游版本并以可复核的协议清单约束前后端实现，避免用本地假设替代契约。
- 可 fork 优先：把通用、可扩展的基架与特定业务模块分开，降低后继项目的二次开发成本。
- 范例即验证：协议支持不能只停留在声明，每个纳入能力都应有可观察的示例与测试路径。
- 单主线模块化：薄内核不依赖业务模块，组合根静态汇集候选模块，Profile 只决定启动时启用集合；模块边界、依赖、数据与失败语义必须可验证。
- 操作原则以 [docs/architecture/principles.md](../architecture/principles.md) P-001 至 P-006 为准。

## 战略假设与未知

| id | 假设 / 未知 | 影响 | 状态 |
|----|-------------|------|------|
| H-001 | 必须从固定的 `schema-ui-docs` `v2.7.0` 提取完整协议能力清单、结构 schema 与 conformance 范围，**并据此**冻结 MVP 的协议覆盖边界。 | VP-001 的协议实现计划与“完整协议覆盖”声明 | **分列**：① 清单提取 = `verified`（已落盘 [protocol-inventory-v2.7.0.md](protocol-inventory-v2.7.0.md)）；② 覆盖子集冻结 = `verified`（Root D-009 冻结 [v0.1.3 覆盖表](../workspace-001-mvp-admin-foundation/GOAL-001-mvp-admin-foundation/attachments/I-PROTO-001-coverage-draft.md)）。它只确认 MVP 子集，**不**主张完整协议支持或实现/验收完成。`F-V006` → `fixed`。 |

## 与工作区 / VP 的关系

- 本 Charter 是对齐链源头；不使用 Goal 的 `done` 状态，也不维护 progress%。
- 已关闭的 [VP-001](plans/VP-001-mvp-admin-foundation.md)、[VP-002](plans/VP-002-production-admin-foundation.md)、[VP-003](plans/VP-003-modular-admin-architecture.md) 与 [VP-004](plans/VP-004-module-contribution-readiness.md) 保留各自交付历史，并已精确 re-align 到本版本而不重开。VP-003 终态架构由 [module-architecture.md](../architecture/module-architecture.md) 固化；一方模块贡献操作契约由 [module-contribution-playbook.md](../architecture/module-contribution-playbook.md)（VP-004）固化。
- **当前无 active 交付 VP**。订单/钱包/类目/通知等业务能力仍属后续独立 VP 候选，建 VP 前须 `/vision` 复核是否需要 strategic 修订。
- 工作区与 Root 必须挂接 `plan_refs` / `primary_plan`。现行 primary 工作区：`workspace-001-mvp-admin-foundation`（Root `GOAL-001-mvp-admin-foundation`，`primary_plan` = VP-001）；已关闭 VP 的 delivery 区历史绑定保留，不改变 primary。

## 现行版本

| 项 | 值 |
|----|----|
| `vision_id` | `schema-ui-core-admin-foundation` |
| 版本 | `0.2.0` |
| 状态 | `active` |
| 引用格式 | `schema-ui-core-admin-foundation@0.2.0` |

修订史见 [revisions.md](revisions.md)，愿景审视台账见 [reviews.md](reviews.md)。
