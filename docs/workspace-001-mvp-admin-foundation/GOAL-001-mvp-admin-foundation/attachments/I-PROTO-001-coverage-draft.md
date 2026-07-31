---
title: I-PROTO-001 · MVP 协议覆盖纳入/排除表（冻结版）
status: active
doc_type: info-coverage-freeze
created: 2026-07-31
updated: 2026-07-31
parent: GOAL-001-mvp-admin-foundation
version: 0.1.3
related_info: I-PROTO-001
related_decision: D-005, D-007, D-008, D-009
source_inventory: docs/vision/protocol-inventory-v2.7.0.md
inventory_pin: ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b
freeze_status: frozen
frozen_at: 2026-07-31
freeze_decision: D-009
---

# I-PROTO-001 · MVP 协议覆盖纳入/排除表（冻结版）

> **性质**：R2 已正式冻结的 MVP 覆盖基线（`I-PROTO-001` → `verified`；D-009）。
> **冻结版本**：`v0.1.3`，固定于 2026-07-31。
> **不是**：完整协议支持、R3-R5 实现/验收完成或 VP 关门证据；只可作为后续规划与实现范围基线。
> **权威输入**：[protocol-inventory-v2.7.0.md](../../../vision/protocol-inventory-v2.7.0.md) §3（`mvp_candidate` 仅提示）。  
> **冻结证据**：用户书面确认本表 v0.1.3 → Root [D-009](../01-decision.md) → `I-PROTO-001=verified`。

## 0. 冻结原则

1. **Charter / VP 边界优先**：核心账号与权限必纳；每一纳入域须可落前后端路径 + 范例 + 验证（R5 再登记 `I-PROTO-003`）。
2. **清单全量 ≠ 覆盖全量**：inventory 是能力全集；本表只选 VP-001 MVP 子集。
3. **`mvp_candidate` 可推翻**：yes/partial/optional 只作初判；本表为决策输入。
4. **业务非目标**：订单/钱包/类目/通知等内容域不纳入；上游样例仅可作结构模板并改写。
5. **partial 必须写清边界**：写入「纳入子面 / 明确排除子面」，禁止裸 `partial` 充当冻结。
6. **结构契约随纳入域带入**：被纳入域所依赖的 `docs/schemas/*.json` 自动进入验证基线，不另开“只纳 schema、不纳语义”的空洞行。

## 1. 能力域（domain_id）冻结范围

| domain_id | 冻结 disposition | 纳入边界 | 明确排除 / 延后 | 主要验证入口 | 备注 |
|-----------|------------------|------------------|-------------------------|--------------|------|
| D-NODE | **include** | Node/Page 树解析与合法 page 产出/消费 | 无（基座） | structural：`node`/`page` schema | 全站基座 |
| D-EXPR | **include** | 表单/列表常用 reaction 与 visible-when 级表达式 | 冷门比较语义边角、与未纳控件绑定的复杂表达式 | fixtures：`reactions` | 支撑权限显隐与表单联动 |
| D-COMP | **include-partial** | MVP 初始 type 白名单见 §5；涵盖最小布局、数据/操作与 D-FORM；Q4 已确认冻结时固化原则 + 初表 | 完整 registry；§5 未列 type；`upload`、`chart`、`statCard` 等 | registry + node/page structural；`component-format` 仅格式子集 | A-003 F-001 闭合证据见 §5/§5.1；不主张全 registry |
| D-DATA | **include** | 列表/详情 datasource、response mapping、query 序列化、static-data | 与未纳域强绑定的专用 API 形状 | fixtures：`request-construction`、`response-mapping`、`query-serialization`、`static-data` | Go 列表/详情 API 主责 |
| D-ACT | **include-partial** | 仅非批量的 page/row action 触发、request lifecycle、编辑/删除等通用后端 request | **MVP 排除**所有依赖 D-TABLE 多选批量的 action/request 语义、payload 或 mapping；业务订单动作 | fixtures：`actions`、`request-lifecycle` 的非批量子集 | Q1=否；与 D-TABLE 一致；A-005 已闭合 A-003 F-002 |
| D-PERM | **include** | 账号会话最小闭环 + 权限继承 / intent → UI 显隐禁用 + Go 鉴权模型 | 完整 IAM 产品（SSO 联邦、细粒度审计后台等）除非 R4 决策扩写 | fixtures：`permissions-inheritance`；场景 `permission-inheritance` | **核心必纳**（Charter） |
| D-APP | **include** | App manifest 装载、导航壳、路由入口 | 多租户应用市场、动态远程 manifest 编排 | fixtures：`app-manifest`、`app-navigation` | 主要在 R3 落地外壳 |
| D-TABLE | **include-partial** | 排序声明、搜索表、基础列表交互 | **MVP 排除**多选批量（ADR 0022）及其 D-ACT action/request 语义；可作为 R5+ 或完整实现线 | fixtures：`table-sort`、`search-table` | Q1=否；列表 Admin 刚需 |
| D-FORM | **include-partial** | 基础 `form` / `input` / `select`，以及全部 2.6/2.7 控件与 `defaultValue` 规则（见 §5）；校验展示 + 编辑回填 | §5 未列 type；`upload`；完整 registry 的其他表单 type | registry + node/page structural + 版本/capability；`component-format` 仅格式子集；scenarios support-only | D-008；A-003 F-001 闭合证据见 §5/§5.1 |
| D-UPLOAD | **exclude** | — | 上传 UI/端点/fixtures `uploads` 整域 | — | inventory `optional`；完整实现线候选 |
| D-VER | **include** | `supportedCapabilities` / 版本协商 / runtime defaults（baseURL 等） | 多协议大版本并行矩阵 | fixtures：`version-negotiation`、`runtime-defaults` | 兼容边界声明必需 |
| D-VAL | **include** | 构建时/加载时结构校验（schemas）；可选服务端校验不强制 R2 | 自研替代上游 schema 语义 | `docs/06-validation.md` + 6 schemas | 支撑 I-PROTO-004 工程策略 |

### 1.1 冻结汇总计数

| disposition | domain_id |
|-------------|-----------|
| include | D-NODE, D-EXPR, D-DATA, D-PERM, D-APP, D-VER, D-VAL（7） |
| include-partial | D-COMP, D-ACT, D-TABLE, D-FORM（4） |
| exclude | D-UPLOAD（1） |

## 2. Fixture suite 冻结映射

| suite id | 冻结 disposition | 绑定 domain | 说明 |
|----------|------------------|-------------|------|
| actions | include-partial | D-ACT | 仅非批量 page/row action；依赖 D-TABLE 多选批量的 case exclude（Q1=否） |
| app-manifest | include | D-APP | |
| app-navigation | include | D-APP | |
| component-format | include | D-COMP | 纳入当前五个 v2.7 format case；仅验证数值/string 非强制转换，不替代 type 白名单或页面结构检查 |
| permissions-inheritance | include | D-PERM | 核心 |
| query-serialization | include | D-DATA | |
| reactions | include | D-EXPR | |
| request-construction | include | D-DATA | |
| request-lifecycle | include-partial | D-ACT | 仅非批量 lifecycle；依赖 D-TABLE 多选批量的 case exclude（Q1=否） |
| response-mapping | include | D-DATA | |
| runtime-defaults | include | D-VER | |
| search-table | include | D-TABLE | |
| static-data | include | D-DATA | |
| table-sort | include | D-TABLE | |
| version-negotiation | include | D-VER | |
| uploads | **exclude** | D-UPLOAD | 与域排除一致 |
| scenarios | **support-only** | 多域 | 信息性场景；**不是**单独 conformance 门禁；可作范例候选源 |

> 上游 `conformance/reference-*` / `runner` 仍为 manifest **excluded** 参考实现，不得单独证明兼容。

## 3. 范例场景候选（冻结范围内 · 非 I-PROTO-003 闭合）

R5 前仅作规划提示；路径以 inventory §2.5 为准：

| 优先级 | 场景 / 样例 | 服务域 |
|--------|-------------|--------|
| P0 | permission-inheritance | D-PERM |
| P0 | Admin 外壳 + 导航（自建，对齐 app-manifest） | D-APP |
| P1 | search-form-table / data-table | D-DATA, D-TABLE |
| P1 | admin-list-edit-lifecycle 或 list-detail（去业务化改写） | D-ACT, D-DATA, D-FORM |
| P2 | form-with-reactions（基础） | D-EXPR, D-FORM |
| 不做（MVP） | form-with-upload；form-controls-advanced/extended 全量；order-* 业务样例原文 | D-UPLOAD / 非目标 |

## 4. 与后续信息项的接口

| 信息项 | 本草案之后 |
|--------|------------|
| I-PROTO-001 | D-009 已冻结本表 v0.1.3 → `verified`；后续范围变更须新决策与新版表 |
| I-PROTO-002 | 仅对 **include 的 D-PERM** 设计最小 API 与映射（R4 前） |
| I-PROTO-003 | 仅为 **include / include-partial** 域登记范例路径与验证入口（R5 前） |
| I-PROTO-004 | vendor vs pin 不阻断本表；建议 R3 前定，以支持 schema/fixture 校验 |

## 5. D-COMP / D-FORM 初始 type 白名单（Q2 全量确认）

> **白名单原则**：只纳入下表 type 与属性规则；任何未列 registry key 均不在 MVP 覆盖内。它是当前 MVP 的初表，不是完整 registry 支持声明。固定源均为 `schema-ui-docs@v2.7.0` commit `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b`。

| 子面 | 纳入 type / 属性 | 范围与版本规则 | 排除 / 备注 |
|------|-----------------|----------------|-------------|
| 布局 | `grid`, `section`, `tabs` | 当前 D-APP 最小导航与页面布局 | 其余布局 type 需另决策 |
| 数据与操作 | `text`, `table`, `recordView`, `actionButton` | 支撑 D-DATA、D-TABLE、D-ACT 的列表、详情和触发入口 | `chart`, `statCard` 不在 MVP 初表；反馈 toast 属 D-ACT 行为，不是 registry type |
| 基础表单 | `form`, `input`, `select`（单值） | ADR-0021 最小编辑/回填面；`select` 的多值模式按下行版本能力处理 | `upload` 由 Q3 保持 exclude |
| 2.6 扩展表单 | `textarea`, `switch`, `checkbox`, `radio`, `select.props.mode: multiple` | 页面 `protocolVersion >= "2.6"` 且 `requiredCapabilities` 含 `form.controls.extended`；提交分别为 string / boolean / boolean / 单值 / 数组 | `select.mode=multiple` 是表单值数组，不是 D-TABLE 批量 action |
| 2.7 进阶表单 | `cascader`, `checkboxGroup`, `richText`, `password` | 页面 `protocolVersion >= "2.7"` 且 `requiredCapabilities` 含 `form.controls.advanced`；提交分别为 path 数组 / value 数组 / Markdown string / string | Host 负责 password 遮罩与 richText 展示消毒 |
| 2.7 属性能力 | 任一上列字段的 `props.defaultValue` | `protocolVersion >= "2.7"` + `form.controls.advanced`；wire 类型匹配，`recordSource` 覆盖 | `defaultValue` 不是独立 component type |

**固定源证据**：

- [2.5-to-2.6 migration](https://github.com/magicvr/schema-ui-docs/blob/ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b/docs/migrations/2.5-to-2.6.md) §1-2：2.6 控件、capability、wire 规则。
- [2.6-to-2.7 migration](https://github.com/magicvr/schema-ui-docs/blob/ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b/docs/migrations/2.6-to-2.7.md) §1-2：2.7 控件、`defaultValue`、capability、wire 规则。
- [component registry](https://github.com/magicvr/schema-ui-docs/blob/ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b/docs/schemas/component-registry.json)：上述 type 及 2.6/2.7 `since` 元数据。

### 5.1 F-001 验证子集与边界

| 验证面 | R5 要求 | 证据 / 限制 |
|--------|---------|-------------|
| Node / Page 结构 | 对纳入页面运行 `node.schema.json` 与 `page.schema.json`；每个 node `type` 必须解析到 §5 白名单和 registry；页面 version/capability 必须满足 2.6 / 2.7 下限 | node schema 负责 node 结构，page schema 负责 meta/body；registry membership 与 version/capability 由 Renderer L2 检查 |
| `component-format` | 纳入现有五个 v2.7 case：`currency-accepts-finite-number`、`currency-rejects-string-coercion`、`percent-rejects-boolean-number-coercion`、`datetime-accepts-string`、`datetime-rejects-number-coercion` | 它只覆盖 format 数据类型与非强制转换，**不**证明表单 type 或 props 全覆盖 |
| form-controls scenarios | `form-controls-extended` / `form-controls-advanced` 作为范例和手工可观察路径 | Q5=否：不作为独立自动化门禁；R5 仍须登记范例页和验证入口（I-PROTO-003） |

固定源：[node schema](https://github.com/magicvr/schema-ui-docs/blob/ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b/docs/schemas/node.schema.json)、[page schema](https://github.com/magicvr/schema-ui-docs/blob/ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b/docs/schemas/page.schema.json)、[component-format cases](https://github.com/magicvr/schema-ui-docs/blob/ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b/conformance/fixtures/component-format/cases.json)、[2.6 scenario](https://github.com/magicvr/schema-ui-docs/blob/ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b/docs/05-scenarios/form-controls-extended.md)、[2.7 scenario](https://github.com/magicvr/schema-ui-docs/blob/ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b/docs/05-scenarios/form-controls-advanced.md)。

## 6. 用户裁决与冻结证据

| # | 问题 | 用户方向（2026-07-31） | 冻结证据 / 后续界限 |
|---|------|------------------------|------------|
| Q1 | D-TABLE 多选批量是否进入 MVP？ | **否**（exclude-from-MVP / 完整实现线） | 已写入 D-ACT/D-TABLE 与 fixture 映射；A-005 已闭合 F-002 |
| Q2 | D-FORM 是否纳入任一 2.6/2.7 扩展控件？ | **是：全部** | §5/§5.1 已列 type、属性、结构与验证子集；D-009 正式冻结 |
| Q3 | D-UPLOAD 是否升格 include？ | **否**（保持 exclude） | 无 |
| Q4 | D-COMP 白名单粒度：冻结时写 type 列表，还是 R3 再附？ | 冻结时写**原则 + 初表**，R3 不得静默扩域 | §5 已固化初始表；新 type 仍须另决策 |
| Q5 | `scenarios` suite 是否任何一条进自动化门禁？ | **否**；仅范例/手工路径 | 无 |

用户另确认：草案修订后需要同 scope self audit；A-005 已完成并通过。详见 [D-007](../01-decision.md)。

## 7. 变更规则（冻结后）

- 本文件是 `v0.1.3` 的冻结基线；`status: active` 与 `freeze_status: frozen` 只代表覆盖范围已决，不代表实现或验收完成。
- 变更覆盖子集 = 新决策（D-00N）+ 新版本表 + 重估 I-PROTO-002/003；不得静默改写本版本。
- 仅有会话确认而未写入 Root 冻结决策时，不得把收集表当作放行证据；本版本已由 D-009 留痕。
