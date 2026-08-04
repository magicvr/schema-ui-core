---
doc_type: protocol-inventory
title: schema-ui-docs v2.7.0 协议实施清单（本地提取）
status: active
source_repo: https://github.com/magicvr/schema-ui-docs
source_tag: v2.7.0
source_commit: ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b
artifact_version: 2.7.0
protocol_version: "2.7"
created: 2026-07-31
updated: 2026-08-04
parent: null
version: 0.1.1
vision_ref: schema-ui-core-admin-foundation@0.2.0
source_vision_ref: schema-ui-core-admin-foundation@0.1.0
serves: VP-001-mvp-admin-foundation
---

# 协议实施清单 · schema-ui-docs@v2.7.0

本清单在 Charter `@0.1.0` 下提取，2026-08-04 仅将现行机读引用 re-align 到 `@0.2.0`。外部 pin、清单内容和既有 `I-PROTO-001 v0.1.3` 覆盖证据均未改变。

> **用途**：闭合 `F-V001` 所需的**本地可复核**协议能力与结构/行为契约清单，并给出 React / Go / 范例 / 验证路径映射。  
> **不是**：MVP 覆盖子集的冻结声明，也不是“已实现协议兼容”的证据。  
> **门禁**：覆盖范围冻结与实施核验在开区后由 **`/govern`** 信息项/决策推进；本文件只提供清单与映射权威。

## 1. 固定源

| 项 | 值 |
|----|-----|
| canonical source | https://github.com/magicvr/schema-ui-docs |
| release ref | `v2.7.0` |
| pinned commit | `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b` |
| manifest | https://raw.githubusercontent.com/magicvr/schema-ui-docs/ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b/protocol-manifest.json |
| 提取日期 | 2026-07-31 |
| manifest 核验 | `artifactVersion`=`2.7.0`，`protocolVersion`=`2.7`，含 semantic / structural / behavioral authority |

上游分工原则（总纲）：后端定义页面语义、数据与结构；前端实现 Renderer、组件库与皮肤。本仓 MVP 方向为 **React 前端 + Go 后端**。

## 2. 权威分层（来自 protocol-manifest）

### 2.1 语义规范（semanticSpecs）

| id | 上游路径 | 主题 |
|----|----------|------|
| S-00 | `PROJECT_CHARTER.md` | 上游项目章程 / 权威边界 |
| S-01 | `docs/00-overview.md` | 总纲、术语、版本与能力声明 |
| S-02 | `docs/01-node-protocol.md` | Node 核心结构 |
| S-03 | `docs/02-reaction-expression.md` | 联动表达式 |
| S-04 | `docs/03-component-registry.md` | 组件 type 注册表语义 |
| S-05 | `docs/04-datasource-contract.md` | 数据源 / API 契约 |
| S-06 | `docs/06-validation.md` | 校验规则与工具链 |
| S-07 | `docs/07-actions-contract.md` | Action 行为契约 |
| S-08 | `docs/08-renderer-spec.md` | Renderer 实现规范 |
| S-09 | `docs/09-app-manifest.md` | App Manifest / 导航 |
| S-ADR-0001…0033 | `docs/decisions/0001-…`～`0033-…` | 架构决策（见 §2.4） |

### 2.2 结构契约（structuralContracts）

固定 commit 下 `docs/schemas/*.json` 共 **6** 个文件：

| 文件 | 角色 |
|------|------|
| `docs/schemas/node.schema.json` | Node 结构 |
| `docs/schemas/page.schema.json` | Page 结构 |
| `docs/schemas/action.schema.json` | Action 结构 |
| `docs/schemas/reaction.schema.json` | Reaction 结构 |
| `docs/schemas/app-manifest.schema.json` | App Manifest 结构 |
| `docs/schemas/component-registry.json` | 组件注册 DSL（非纯 JSON Schema 文档，属结构权威） |

### 2.3 行为契约（behavioralContracts）

| 类别 | 路径 | 说明 |
|------|------|------|
| 套件 schema | `conformance/schemas/fixture-suite.schema.json` | fixture 套件形状 |
| fixture 套件目录 | 见下表（各含 `cases.json`） | 行为用例权威 |

| suite id | 路径 |
|----------|------|
| actions | `conformance/fixtures/actions/` |
| app-manifest | `conformance/fixtures/app-manifest/` |
| app-navigation | `conformance/fixtures/app-navigation/` |
| component-format | `conformance/fixtures/component-format/` |
| permissions-inheritance | `conformance/fixtures/permissions-inheritance/` |
| query-serialization | `conformance/fixtures/query-serialization/` |
| reactions | `conformance/fixtures/reactions/` |
| request-construction | `conformance/fixtures/request-construction/` |
| request-lifecycle | `conformance/fixtures/request-lifecycle/` |
| response-mapping | `conformance/fixtures/response-mapping/` |
| runtime-defaults | `conformance/fixtures/runtime-defaults/` |
| scenarios | `conformance/fixtures/scenarios/` |
| search-table | `conformance/fixtures/search-table/` |
| static-data | `conformance/fixtures/static-data/` |
| table-sort | `conformance/fixtures/table-sort/` |
| uploads | `conformance/fixtures/uploads/` |
| version-negotiation | `conformance/fixtures/version-negotiation/` |

> 上游 `conformance/reference-js`、`reference-python`、`runner` 在 manifest **excluded** 中：可作参考实现，**不是**本清单的协议语义权威，也不得单独当作“已兼容”证明。

### 2.4 关键 ADR 索引（semantic）

| ADR | 文件 | 与 MVP 相关度（初判，未冻结） |
|-----|------|------------------------------|
| 0001 | single-node-tree | 核心 |
| 0002 | not-two-schema-uischema | 核心 |
| 0003 | context-namespace-and-visible-when | 核心 |
| 0004 | row-level-scope | 表格/行操作 |
| 0005 | response-mapping | 数据契约 |
| 0006 | expression-evaluation-order | 表达式 |
| 0008 | row-action-backend-request | 行操作 → Go |
| 0009 | strict-version-negotiation | 版本协商 |
| 0010 | query-serialization | 查询 → Go |
| 0011 | reserved-query-params | 查询 |
| 0012 | upload-execution | 上传 |
| 0013 | dataref-read-only | 数据引用 |
| 0014 | action-replay-semantics | Action |
| 0015 | request-lifecycle-latest-wins | 请求生命周期 |
| 0016 | expression-comparison-semantics | 表达式 |
| 0017 | runtime-defaults-and-baseurl | 运行时默认 |
| 0018 | component-boundary-semantics | 组件边界 |
| 0019 | v2-admin-scope | Admin 范围 |
| 0020 | page-action-trigger | 页面工具栏 |
| 0021 | record-navigation-and-form-load | 编辑回填 |
| 0022 | table-selection-and-batch-request | 多选批量 |
| 0023 | container-permission-inheritance | **账号/权限相关** |
| 0024 | record-view | 只读详情 |
| 0025 | app-manifest | 应用清单 |
| 0026 | app-navigation | 导航 |
| 0027 | table-sort-declaration | 表格排序 |
| 0028 | form-control-surface | 扩展表单控件（2.6） |
| 0029–0033 | cascader / checkbox-group / rich-text / password / form-default-value | 进阶表单控件（2.7） |

### 2.5 信息性场景与样例（非语义权威，范例候选）

| 场景文档 | 路径 |
|----------|------|
| README | `docs/05-scenarios/README.md` |
| admin-list-batch | `docs/05-scenarios/admin-list-batch.md` |
| admin-list-detail-lifecycle | `docs/05-scenarios/admin-list-detail-lifecycle.md` |
| admin-list-edit-lifecycle | `docs/05-scenarios/admin-list-edit-lifecycle.md` |
| data-table | `docs/05-scenarios/data-table.md` |
| form-controls-advanced | `docs/05-scenarios/form-controls-advanced.md` |
| form-controls-extended | `docs/05-scenarios/form-controls-extended.md` |
| form-with-reactions | `docs/05-scenarios/form-with-reactions.md` |
| form-with-upload | `docs/05-scenarios/form-with-upload.md` |
| grid-dashboard | `docs/05-scenarios/grid-dashboard.md` |
| permission-inheritance | `docs/05-scenarios/permission-inheritance.md` |
| row-backend-actions | `docs/05-scenarios/row-backend-actions.md` |
| search-form-table | `docs/05-scenarios/search-form-table.md` |

Manifest 点名的样例 YAML：`order-list-batch`、`order-detail-lifecycle`、`user-profile-extended`、`user-profile-advanced`（另有 edit/list lifecycle 与 audit-*-bad 负例样例，见上游 `_samples/`）。

## 3. 能力面映射（React / Go / 范例 / 验证）

图例：`primary` = 主责实现面；`support` = 配合面；`n/a` = 通常不在该面落地。  
**`mvp_candidate`** 仅作规划提示（账号权限 + 基架 + 协议范例），**不是**已冻结的 MVP 覆盖集合。

| domain_id | 能力域 | React（Renderer / UI） | Go（数据 / 动作 / 权限） | 范例页面候选 | 验证路径 | mvp_candidate |
|-----------|--------|------------------------|--------------------------|--------------|----------|---------------|
| D-NODE | Node / Page 树 | primary：解析、递归渲染 | support：产出 YAML/JSON 页面配置 | 任意合法 page | structural：`node`/`page` schema | yes |
| D-EXPR | Reaction / 表达式 | primary：求值与 fulfill | n/a（语义在配置） | form-with-reactions | fixtures：`reactions` | yes |
| D-COMP | 组件注册表 | primary：type→组件 | support：只写语义 props | form-controls-*、data-table | fixtures：`component-format`；registry JSON | partial（MVP 子集待 /govern） |
| D-DATA | Datasource / 响应映射 | primary：请求构造与 mapping | primary：列表/详情 API 形状 | search-form-table、data-table | fixtures：`request-construction`、`response-mapping`、`static-data`、`query-serialization` | yes |
| D-ACT | Actions（含 row/page/batch） | primary：触发与生命周期 | primary：后端 request 语义 | row-backend-actions、admin-list-* | fixtures：`actions`、`request-lifecycle` | yes |
| D-PERM | 权限继承 / intent | primary：UI 显隐与禁用 | primary：鉴权与权限模型 | permission-inheritance | fixtures：`permissions-inheritance` | **yes（核心账号权限）** |
| D-APP | App manifest / 导航 | primary：装载与导航壳 | support：manifest 托管/生成 | Admin 外壳 + 导航 | fixtures：`app-manifest`、`app-navigation` | yes（外壳） |
| D-TABLE | 表格排序 / 多选批量 | primary | support：批量 API | admin-list-batch、data-table | fixtures：`table-sort`、`search-table` | partial |
| D-FORM | 扩展/进阶表单控件 | primary | support：字段校验 API | form-controls-extended/advanced | structural + scenarios | partial |
| D-UPLOAD | 上传 | primary：上传 UI/协议 | primary：上传端点 | form-with-upload | fixtures：`uploads` | optional |
| D-VER | 版本协商 / capabilities | primary：`supportedCapabilities` | support：协商头/错误 | 全站 | fixtures：`version-negotiation`、`runtime-defaults` | yes |
| D-VAL | 配置校验 | support：构建时/加载时 | support：可选服务端校验 | — | `docs/06-validation.md` + schemas | yes |

### 3.1 前后端职责摘要

| 面 | 职责边界（对齐上游总纲） |
|----|--------------------------|
| **React Renderer** | Node 树解析；`type` 分发；表达式与 reactions；Action 触发与 request lifecycle；capabilities 匹配；Admin 壳 + 主题（Tailwind/shadcn，浅/深色） |
| **Go 后端** | 账号与权限；列表/详情/动作 API 符合 datasource 与 actions 契约；可选 page/manifest 托管；上传与批量等服务端语义 |
| **范例** | 每个**纳入 MVP 的**能力域至少一条可观察页面或场景（优先复用 `05-scenarios` 与 `_samples`） |
| **验证** | 结构：JSON Schema；行为：对应 conformance fixtures；集成：账号权限链路 e2e；不得把 excluded 的 reference runner 当作唯一门禁 |

## 4. 与 VP-001 / Charter 的关系

| 项 | 状态 |
|----|------|
| 外部固定源 | 已 pin（Charter / VP 已记录） |
| 完整实施清单 | **本文件**（F-V001 闭合证据） |
| MVP 协议覆盖子集冻结 | **已完成**；workspace-001 Root D-009 冻结 `I-PROTO-001 v0.1.3`，不等于全量协议支持 |
| “支持全部协议功能”主张 | **禁止**，直至覆盖子集冻结且实现证据闭合 |
| H-001 | 清单已提取；覆盖子集已由 workspace-001 Root D-009 冻结为 `I-PROTO-001 v0.1.3` |

## 5. 历史开区信息项建议（已交 `/govern`）

下列不是愿景 finding，仅作实现层登记模板：

| 建议 id | 问题 | 级别 | 最晚阶段 |
|---------|------|------|----------|
| I-PROTO-001 | 哪些 domain_id / fixture suite 纳入 VP-001 MVP？ | required | 方案冻结前 |
| I-PROTO-002 | 账号权限最小 API 与 `D-PERM` 映射是否完整？ | required | 实施前 |
| I-PROTO-003 | 每条纳入能力的范例页路径与自动化/手工验证入口？ | required | 验收前 |
| I-PROTO-004 | 是否 vendor 上游 schemas/fixtures 到本仓，或 pin 远程校验？ | non-blocking | 实施前 |

## 6. 变更规则

- 上游 pin 变更 → 重新提取本清单并升 `version`；同步 Charter/VP 协议引用。  
- 仅修正映射笔误 → editorial，不改 pin。  
- 本文件**不得**写入 goal progress% 或关闭 Goal finding。
