---
title: I-001 · 协议实施差量矩阵（R1 方案输入）
status: active
doc_type: info-gap-matrix
created: 2026-08-01
updated: 2026-08-01
parent: GOAL-001-production-admin-foundation
version: 0.1.1
related_info: I-001
related_decision: D-004
coverage_baseline: I-PROTO-001 v0.1.3
coverage_path: docs/workspaces/workspace-001-mvp-admin-foundation/GOAL-001-mvp-admin-foundation/attachments/I-PROTO-001-coverage-draft.md
inventory_pin: ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b
---

# I-001 · 协议实施差量矩阵

> **性质**：回答「冻结协议到当前代码 / fixture / Renderer 运行路径的实施差量是什么」。  
> **权威范围输入**：`I-PROTO-001 v0.1.3`（Q2 路径见 frontmatter；**不**改写冻结表）。  
> **不是**：R1 实施完成、验收通过或覆盖面扩大；不把本矩阵勾选 Root 路线图 R1。  
> **扫描日期**：2026-08-01（仓库只读对照；无产品代码变更）。

## 0. 总览结论

| 维度 | 结论 |
|------|------|
| 协议子集是否冻结 | 是：`v0.1.3` 由 workspace-001 Root D-009 冻结；本区继承 |
| 库级 / 测试级能力 | **大量已有**：manifest、导航、权限 L2、form 白名单、`$context` reactions、table/query/action adapters、Ajv 结构校验、上游 fixture 执行 |
| 产品默认页面路径 | **未产品化**：`App` 优先走手写 `EXAMPLE_PAGES`；非示例页仅 shell 占位；`schemaUrl` **未**驱动加载与 `RenderPage` |
| R1 核心差量 | **默认主路径**：Schema 加载 → 校验 → Node 解释 → 渲染；失败统一；示例退居次要 |
| R1 不必做 | 真实登录（R2）、DB 持久化权限（R3）、业务 CRUD 验收实体选型（R4）、Docker/fork 计时（R5）、扩大 `v0.1.3` |

### R1 必须差量（方案边界）

1. **Schema 运行时加载**：按 manifest `page.schemaUrl`（及路由参数解析）获取页面文档；无加载器 = 未产品化。  
2. **默认渲染主路径**：匹配路由后默认 `RenderPage`（或等价），而非 `EXAMPLE_PAGES[pageId]`。  
3. **加载前/中结构校验**：至少 page/node（及既有 capability）fail-closed；错误可观察、可统一。  
4. **未知节点 / 无效 Schema / 运行时错误**：确定性拒绝与统一错误面（库内已有部分 fail-closed；缺产品级串联）。  
5. **代表性页面**：至少一组列表、详情/record、表单、组合页以 **Node 树** 交付；新增页优先改 Schema。  
6. **示例降级**：手写 React 示例保留为兼容/演示，不得再作为「新增业务页」默认方式。

> **2026-08-01 细化（GOAL-003 D-003）**：经用户裁决，降级机制定为**迁移为 Schema**——5 个手写示例语义改写为 Schema 文档，经默认主路径（route → schemaUrl → 加载+校验 → RenderPage）渲染，应用内不再保留手写页面路径；本文档「保留为兼容/演示」的表述被该决策细化取代。I-001 差量结论仍成立：示例不得作为「Schema 驱动完成」的默认主路径证明。

### 明确非 R1（本矩阵标注 defer）

| 项 | 归属 |
|----|------|
| 真实登录 / Token / 登出 | R2（I-002） |
| 用户/角色/菜单持久化与种子 | R3（I-003） |
| 代表性 CRUD 业务实体与 API 语义定稿 | R4（I-004） |
| 15 分钟 fork、Docker、部署基线 | R5（I-005） |
| 操作日志 | I-006 non-blocking |
| D-UPLOAD 全域 | `v0.1.3` exclude；完整实现线候选 |
| 上游 multi-round `$deps` reactions | 已记为 MVP 边界外；保持 exclude 除非新决策 |

## 1. 能力域差量（对照 v0.1.3 disposition）

| domain_id | disposition | 现状（可核对路径） | 差量 | R1 必须 | 证据 / 备注 |
|-----------|-------------|-------------------|------|---------|-------------|
| **D-NODE** | include | `render.ts` `parseRenderNode` / 白名单；`docs/schemas/node.schema.json` + `page.schema.json`；Ajv `schema-validate.ts`（测试） | 无生产「从 schemaUrl 加载 page 文档 → 解析 body」管线；示例内嵌 TS 常量 | **是**（加载+默认路径+校验串联） | `apps/web/src/renderer/render.ts`；`protocol/conformance/schema-validate.ts` |
| **D-EXPR** | include | `reactions.ts` `$context` 引擎；`form-with-reactions` 示例走 `RenderPage`；fixture `upstream/reactions.cases.json` 已 pin | 产品页未默认消费 Schema 内 reactions；`$deps` 多轮仍排除 | **部分**（主路径上可运行 `$context`；不扩 `$deps`） | `renderer/reactions.ts`；`examples/form-with-reactions-page.tsx`；protocol README |
| **D-COMP** | include-partial | `RenderPage` 白名单：grid/section/tabs/text/table/recordView/actionButton/form；未知 type fail-closed | 未成为默认页面能力；table 数据仍多由示例注入 | **是**（主路径默认渲染白名单节点） | `renderer/render.tsx` 注释 R5 D-COMP；`render.test.tsx` |
| **D-DATA** | include | `records.ts` / `use-records`；Go `handler/records.go` 进程内静态列表；query/response/static-data fixtures | Schema 声明的 datasource 未统一由主路径解释；持久化非 R1 | **部分**（主路径能挂既有 list/fetch；不要求 DB） | `renderer/records.ts`；`apps/api/internal/handler/records.go` |
| **D-ACT** | include-partial | `row-action.ts`、`permissions.executeAction`、request-lifecycle / actions fixtures；非批量边界 | Schema 驱动 page/row action 未默认串联；批量仍 exclude | **部分**（非批量 action 可在主路径触发） | `renderer/row-action.ts`；`permissions.ts` |
| **D-PERM** | include | L2 `permissions.ts`；导航 permissions；`/api/accounts/me` + **静态** `StaticDevSession` | 身份为开发会话，非真实认证（R2）；主路径权限失败 UX 需产品化串联 | **部分**（沿用静态会话 + UI/API 门控可演示；不替换为真实登录） | `account/session.go`；`handler/account.go` |
| **D-APP** | include | `app-manifest.ts` 装载/校验；shell 导航；route match/fallback；app-manifest / app-navigation fixtures | **缺口**：`schemaUrl` 已建模但 **PageSurface 不消费**；非示例页占位文案写明 renderer 为 later boundary | **是**（schemaUrl → 页面主路径） | `App.tsx` `PageSurface` L215–273；`protocol/README.md` |
| **D-TABLE** | include-partial | `data-table.tsx`；search/sort fixtures；示例 data-table / search-form-table | Schema 列表页默认渲染与列/查询绑定未产品化；多选批量 exclude | **部分**（代表性 Schema 列表页） | `components/data-table.tsx`；`examples/*` |
| **D-FORM** | include-partial | §5 控件白名单 + 2.6/2.7 capability gate；`FormControls`；form-controls 示例 | 主路径表单依赖 Schema 文档加载；capability 失败需统一错误面 | **部分**（代表性 Schema 表单页） | `form-controls.ts` / `.tsx` |
| **D-UPLOAD** | exclude | 无实现（符合冻结） | 无 R1 差量 | **否** | 保持 exclude |
| **D-VER** | include | version-negotiation / runtime-defaults fixtures；host 接受 protocol `2.7` | 页面级 version/capability 与加载器错误码需统一到主路径 | **部分** | `conformance/version-negotiate.ts`；app-manifest host 2.7 |
| **D-VAL** | include | 构建/测试态 Ajv；vendor schemas 有 SHA pin | **运行时**页面加载校验未接入 App 主路径 | **是**（加载时结构校验） | `docs/schemas/*`；`schema-validate.ts`；`upstream/provenance.json` |

### 1.1 产品路径对照（关键事实）

| 路径 / 行为 | 现状 | R1 目标态 |
|-------------|------|-----------|
| `App.tsx` → `EXAMPLE_PAGES[pageId]` | 有 pageId 映射则手写 React 页 | 次要/兼容；默认不走此分支 |
| `App.tsx` → 非示例匹配页 | shell +「page renderer remains a later protocol boundary」 | 加载 `schemaUrl` → 校验 → `RenderPage` |
| `manifest.pages[].schemaUrl` | 校验与模板解析存在；**运行时不 fetch** | 解析参数后加载页面文档 |
| `RenderPage` | 库 + 单测 + 至少 `form-with-reactions` 内嵌文档 | 成为新增/默认业务页主路径 |
| 无效 Schema / 未知 node | 库/测试 fail-closed 片段存在 | 产品错误面统一、可观察 |
| `renderer/README.md` | 仍写「R1 预建空分层」 | 与实现同步（实施时改；本回合不改代码） |

## 2. Fixture suite 差量（对照 v0.1.3 §2）

| suite | disposition | 现状 | R1 必须 | 备注 |
|-------|-------------|------|---------|------|
| app-manifest | include | 执行 + pin | 否（已支撑 shell） | `upstream-fixtures.test.ts` |
| app-navigation | include | 执行 + pin | 否 | 同上 |
| actions | include-partial | adapter + stage3 | 否* | *产品串联非 fixture 扩面 |
| component-format | include | adapter | 否 | 五 case format 子集 |
| permissions-inheritance | include | renderer 单测 + 场景 | 否* | R2 真认证另计 |
| query-serialization | include | adapter | 否 | |
| reactions | include | pin；引擎 `$context`；多轮 `$deps` 排除 | 否* | 主路径消费见域表 |
| request-construction | include | adapter | 否 | |
| request-lifecycle | include-partial | adapter | 否* | 非批量 |
| response-mapping | include | adapter | 否 | |
| runtime-defaults | include | adapter | 否 | |
| search-table | include | adapter | 否* | |
| static-data | include | adapter | 否 | |
| table-sort | include | adapter | 否* | |
| version-negotiation | include | adapter | 否 | |
| uploads | exclude | 未纳入 | 否 | 符合冻结 |
| scenarios | support-only | 手写示例对齐部分场景 | 否（非自动化门禁） | 范例候选 |

\*「否」指 **不必新开 fixture 覆盖工程**；R1 验收仍需产品路径可观察证据（页面 + 测试）。

## 3. 验证入口对照

| 验证面 | 现状 | R1 验收建议（方案层） |
|--------|------|----------------------|
| 结构 schema（node/page/…） | 测试内 Ajv | 加载管线强制；失败可见 |
| 行为 fixture | stage3 + upstream tests | 保持 pin；不宣称全协议 |
| 手写示例页 | 5 个 EXAMPLE_PAGES | 迁移为 Schema（2026-08-01 · GOAL-003 D-003）；**不得**作为默认主路径证明「Schema 驱动完成」 |
| E2E / 集成 | 现有 shell / example 测试 | 增加「schemaUrl → 渲染」路径测试 |
| 后端 records / me | 静态会话 + 进程内 records | R1 可继续用；不伪装生产身份 |

## 4. R1 方案边界建议（供 D-004 / 后续子目标）

在 **不扩大** `v0.1.3` 的前提下，R1 实施范围建议冻结为：

1. **In**：Schema 加载器；默认 `RenderPage` 主路径；加载时 D-VAL；白名单 D-COMP/D-FORM；`$context` D-EXPR；代表性列表/表单/组合 Node 页；统一失败面；示例降级。  
2. **Reuse**：现有 manifest shell、fixture pin、records 演示 API、静态 dev session（明确 non-production）。  
3. **Out（本阶段）**：真实认证、持久化 IAM、覆盖表扩域、D-UPLOAD、`$deps` reactions、业务模块、fork/Docker 关门证据。

建议后续子目标拆分（**本回合不创建**，仅规划提示）：

| 候选 | 范围 |
|------|------|
| Schema 加载 + 校验 + 错误面 | D-APP schemaUrl、D-VAL、失败统一 |
| 默认 Renderer 主路径 + 示例降级 | App 路由分支、`RenderPage` 默认 |
| 代表性 Node 页面与回归 | 列表/表单/组合 + 测试证据 |

## 5. 信息项闭合声明

| 字段 | 值 |
|------|-----|
| I-001 问题 | 冻结协议 → 代码/fixture/Renderer 运行路径差量是什么？ |
| 答案载体 | 本文件 §0–§4 |
| 状态建议 | `verified`（差量已可核对；**非** R1 实施完成） |
| 残余风险 | 扫描基于 2026-08-01 树；实施中若发现新域差量须回流信息表（P-005） |
| 禁止表述 | 不得因本文件把 Root `progress` 勾选 R1，或宣称生产级 Renderer 已交付 |
