---
title: I-PROTO-FULL-001 · schema-ui-docs@v2.7.0 整份契约覆盖纳入/排除表（冻结版）
status: active
doc_type: info-coverage-freeze
created: 2026-08-08
updated: 2026-08-08
parent: GOAL-001-full-protocol-contract-v2-7-0
version: 1.0.0
related_info: I-PROTO-FULL-001
related_decision: D-002
source_inventory: docs/vision/protocol-inventory-v2.7.0.md
inventory_pin: ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b
freeze_status: frozen
frozen_at: 2026-08-08
freeze_decision: D-002
supersedes: I-PROTO-001 v0.1.3（仅作历史 MVP 基线，只读；不覆盖其文件）
---

# I-PROTO-FULL-001 · 整份 v2.7.0 契约覆盖纳入/排除表（冻结版）

> **性质**：VP-006 整份契约覆盖的**现行权威**（S1 冻结）。对 `schema-ui-docs@v2.7.0` pin（commit `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b`）下的**能力域、component registry 与 conformance fixture 套件逐项**给出 disposition。
> **冻结版本**：`v1.0.0`（2026-08-08，Root D-002）。
> **不是**：实现/验收完成证明；S2–S5 仍须以真实实现与测试闭合本表的 `include` 行。
> **权威输入**：[protocol-inventory-v2.7.0.md](../../vision/protocol-inventory-v2.7.0.md)（全量清单）；S0 差集证据 [I-S0-001-gap-analysis-v0-1-3-to-full.md](I-S0-001-gap-analysis-v0-1-3-to-full.md)（I-001 closed）。
> **历史基线（只读）**：workspace-001 `I-PROTO-001 v0.1.3`（MVP 子集，回归对照；**本表升版不改写其文件**）。

## 0. 冻结原则

1. **默认 disposition = `include`**（VP-006 exit 1）：对 inventory / registry / fixture 承诺面的可验证兼容是默认目标；`include-partial` 仅允许保真边角与可列明的次要语义缺口，且不得表达「整域或主要子面不打算做」。
2. **范围收缩纪律**：任何 `exclude` / 范围收缩必须是**用户书面接受的有界 residual**（范围 + 复审触发）——**本表无 exclude、无收缩**（S0 差集全部可纳入，I-002 = N/A）。
3. **结构契约随纳入域带入**：被纳入域依赖的 `docs/schemas/*.json`（6 个）全部进入验证基线（已 vendor + SHA pin）。
4. **保真度**：达到「契约语义可验证」；不要求 VP-005 级视觉产品化（VP-005 实施冻结至本 VP closed）。
5. **fail-closed**：白名单外/未支持能力必须显式拒绝或报错，禁止静默忽略或静默降级。
6. **范例即验证**：每个纳入域须有可发现范例 + 验证入口（S4 登记于本表「验证入口」列并指向真实路径）。
7. **验证驱动真实出货函数**：in-repo 测试直接驱动 renderer / 后端真实实现，无 mock 被测单元、无重实现。

## 1. 能力域（domain_id）冻结范围 — 12/12 include

| domain_id | 冻结 disposition | 纳入边界 | 明确排除 | 验证入口（S4 登记 · 真实路径） |
|-----------|------------------|----------|----------|----------------------------------------------|
| D-NODE | **include** | Node/Page 树解析、递归渲染、结构契约（node/page schema） | 无（基座） | structural：`apps/web/src/protocol/conformance/schema-validate.ts`（stage3 结构测试全绿）；范例页 `apps/api/internal/modules/schemarender/schema/*.json`（8 页） |
| D-EXPR | **include** | 全量表达式语法（`$deps`/`$self`/`$context`/`$row` 白名单 + `== != > >= < <= contains && || ! ( )` + 严格类型/码点比较）与快照多轮 reaction 引擎（fulfill/otherwise/observers/baselines/externalUpdates/loop protection/深等检测/`MULTIPLE_VALUE_WRITES`） | 无 | fixtures：`upstream/reactions.cases.json` **16/16 全绿**（stage3）；引擎：`apps/web/src/renderer/reaction-engine.ts` + `reaction-expression.ts`；单元：`reaction-expression.test.ts`（14）；渲染集成：`render.tsx` FormView（值提交/显隐/阻断）；范例 `form-with-reactions` |
| D-COMP | **include** | registry **24/24 type** 渲染与 props 语义（布局/数据/操作/表单全类别）；`component-format` 数据类型与非强制转换 | 无 | fixtures：`component-format` 5/5；registry 门禁：`render.ts`（10 节点 type）/`form-controls.ts`（13 控件 type + 门禁）；展示渲染：`render.tsx` StatCardView/ChartView + 集成测试（`render.test.tsx`）；范例 `form-controls`、`data-table`、`data-display` |
| D-DATA | **include** | 列表/详情 datasource、response mapping、query 序列化（ADR-0010）、static-data、request 构造（非批量） | 无 | fixtures：`request-construction` 64/64 非 batch、`response-mapping` 23/23、`query-serialization` 16/16、`static-data` 9/9；Go 通用 resource CRUD `apps/api/internal/handler/resources.go`（users/roles 真实端点） |
| D-ACT | **include** | 页面/行级/批量 action：request/navigate/modal/upload/custom 类型、request lifecycle、OutcomeBehavior、批量 Trigger（ADR-0022 D4/D5：requiresSelection/batchMapping/EMPTY_SELECTION/reload 清选） | 无（v0.1.3 的批量排除解除） | fixtures：`actions` 11/11、`request-lifecycle` 4/4、`request-construction` batch **11/11**；构造器：`request-construction.ts`（buildBatchRequest）；执行：`render.tsx` runBatchRequest/invokeBatchAction + 集成测试（批量端到端 + 空选 fail-closed）；Go `POST {path}/batch-delete`（`resources.go` + `users_batch_test.go`）；范例 `admin-list-batch` |
| D-PERM | **include** | 账号会话最小闭环 + 权限继承 / intent（ADR-0023）→ UI 显隐禁用 + Go 鉴权模型；17/17 行为 case | 完整 IAM 产品（SSO 联邦、细粒度审计后台等，维持 v0.1.3 边界） | fixtures：`permissions-inheritance` 17/17（SHA `ac124fa1…` pin，`upstream/permissions-inheritance.cases.json`）；`apps/web/src/renderer/permissions.ts` + `permissions-inheritance.test.ts`；Go `internal/account`、`internal/auth`（真实鉴权端点） |
| D-APP | **include** | App manifest 装载、导航壳、路由、真实端点（`/.well-known/schema-ui/app-manifest.json`、`/api/schema/{pageId}`） | 多租户应用市场、动态远程 manifest 编排（维持 v0.1.3 边界） | fixtures：`app-manifest` 37/37、`app-navigation` 16/16（`upstream-fixtures.test.ts`）；`apps/web/src/protocol/app-manifest.ts`；Go `internal/handler/manifest.go`、`schema.go`；集成：`representative-pages.integration.test.tsx` |
| D-TABLE | **include** | 排序声明、搜索表、基础列表交互 + **多选批量（ADR-0022 D2）**：selection.mode=multiple、选中键规范化（去重保序/count=keys.length）、清选时机（筛选/翻页/排序/reload）、全选本页 | 跨页全选、筛选结果全集、部分成功回填（ADR-0022 非目标） | fixtures：`table-sort` 14/14、`search-table` 11/11（含 selection 状态机）；UI：`schema-table.tsx`（多选列/全选/清选/requiresSelection）+ `schema-table.test.tsx`；集成：批量端到端测试；范例 `data-table`、`admin-list-batch` |
| D-FORM | **include** | 全部表单控件（基座 input/select/inputNumber/datePicker/dateRangePicker/form + 2.6 extended + 2.7 advanced + upload）+ `defaultValue`（ADR-0033）+ wire 规则（ADR-0028–0033）+ 校验展示 + 编辑回填（ADR-0021） | 无（v0.1.3 白名单外 4 type 纳入） | `apps/web/src/renderer/form-controls.ts`（13 控件 + 门禁）+ `form-controls.test.ts`（含 new gates）；`form-controls.tsx`（8 类控件渲染）；提交投影展开（dateRange/upload）；范例 `form-controls`（全控件面） |
| D-UPLOAD | **include** | 上传控件（action/actionRef 双模式）、顶层 `type: upload` action、capability `actions.upload` 门禁、客户端编排（ADR-0012：逐文件 multipart、约束前置、原子失败、url 优先取值、幂等重试）、**后端上传端点**（multipart 接收、独立校验、`{url,id,name,size}` 响应、语义错误码） | 无（v0.1.3 整域排除解除） | fixtures：`upstream/uploads.cases.json` **13/13 全绿**；编排：`apps/web/src/protocol/conformance/upload-orchestration.ts`；控件+传输：`form-controls.tsx` UploadField + `render.tsx` uploadFiles；集成：上传端到端测试；Go `internal/handler/upload.go` + `upload_test.go`（POST /api/upload、GET /api/files/{id}）；范例 `form-with-upload` |
| D-VER | **include** | `supportedCapabilities` / 版本协商（ADR-0009）、runtime defaults（ADR-0017，baseURL 等） | 多协议大版本并行矩阵（维持 v0.1.3 边界） | fixtures：`version-negotiation` 44/44、`runtime-defaults` 9/9；`apps/web/src/protocol/load-page.ts` + `conformance/version-negotiate.ts` |
| D-VAL | **include** | 构建时/加载时结构校验（6 schemas，vendor + SHA pin）、registry 语义门禁 | 自研替代上游 schema 语义 | `docs/schemas/*.json`（6 个，provenance SHA pin）；`runtime-schema-validate.ts`（浏览器侧）；stage3 结构测试；加载路径 `load-page.ts`（PAGE_SCHEMA_INVALID fail-closed） |

### 1.1 汇总计数

| disposition | 数量 | domain_id |
|-------------|------|-----------|
| include | **12** | D-NODE, D-EXPR, D-COMP, D-DATA, D-ACT, D-PERM, D-APP, D-TABLE, D-FORM, D-UPLOAD, D-VER, D-VAL |
| include-partial | **0** | — |
| exclude | **0** | — |

## 2. Component registry 冻结范围 — 24/24 include

> 固定源：`docs/schemas/component-registry.json`（SHA pin，vendor 于本仓）。`since` / capability 语义以 registry 与迁移文档为准。

| 子面 | type（全部 include） | 范围与版本规则 |
|------|----------------------|----------------|
| 布局 | `grid`, `section`, `tabs` | 既有（v0.1.3 已纳入） |
| 数据与操作 | `text`, `table`, `recordView`, `actionButton`, **`statCard`**, **`chart`** | statCard/chart 为 supportsData 展示型（format plain/currency/percent、unit、valueField、chartType/xField/yField）；S3 实现渲染与格式应用；table 新增 `selection`/`toolbar`/batch 语义（ADR-0022） |
| 表单基座 | `form`, `input`, `select`（单值/多值）, **`inputNumber`**, **`datePicker`**, **`dateRangePicker`** | inputNumber wire=number（min/max/step/precision）；datePicker wire=ISO string（format/min/max）；dateRangePicker 绑 startField/endField 两字段（禁 fulfill.value/otherwise.value，ADR-0028/0033 语义）；S2 实现 |
| 表单 2.6 扩展 | `textarea`, `switch`, `checkbox`, `radio`（+ `select.mode: multiple`） | 既有；capability `form.controls.extended`；wire 已冻结 |
| 表单 2.7 进阶 | `cascader`, `checkboxGroup`, `richText`, `password`（+ 任一字段 `props.defaultValue`） | 既有；capability `form.controls.advanced`；wire 已冻结 |
| 上传 | **`upload`** | action/actionRef 双模式（oneOf）；actionRef 模式要求页面声明 `actions.upload`；accept/maxSize/multiple 唯一约束源；S2/S3 实现 |

## 3. Fixture 套件冻结映射 — 16/16 行为套件 include（scenarios support-only）

| suite id | 冻结 disposition | 绑定 domain | 说明（S4 执行证据） |
|----------|------------------|-------------|----------------------|
| actions | include | D-ACT | 11/11 全绿（stage3） |
| app-manifest | include | D-APP | 37/37 全绿（upstream-fixtures） |
| app-navigation | include | D-APP | 16/16 全绿（upstream-fixtures） |
| component-format | include | D-COMP | 5/5 全绿（stage3） |
| permissions-inheritance | include | D-PERM | 17/17 全绿（`upstream/permissions-inheritance.cases.json`，SHA `ac124fa1…`；renderer 套件） |
| query-serialization | include | D-DATA | 16/16 全绿（stage3） |
| reactions | **include（v0.1.3 全排除解除）** | D-EXPR | **16/16 全绿**（stage3 · reaction-engine.ts） |
| request-construction | **include（batch 11 case 纳入）** | D-DATA / D-ACT | **75/75 全绿**（stage3，含 11 batchRequest） |
| request-lifecycle | include | D-ACT | 4/4 全绿（stage3） |
| response-mapping | include | D-DATA | 23/23 全绿（stage3） |
| runtime-defaults | include | D-VER | 9/9 全绿（stage3） |
| search-table | include | D-TABLE | 11/11 全绿（stage3，含 selection 事件） |
| static-data | include | D-DATA | 9/9 全绿（stage3） |
| table-sort | include | D-TABLE | 14/14 全绿（stage3） |
| version-negotiation | include | D-VER | 44/44 全绿（stage3） |
| uploads | **include**（S1 已 vendor + SHA pin `aaeb9683…`） | D-UPLOAD | **13/13 全绿**（stage3 · upload-orchestration.ts） |
| scenarios | support-only | 多域 | 信息性场景；范例候选源；不是独立 conformance 门禁 |

> 上游 `conformance/reference-*` / `runner` 仍为 manifest **excluded** 参考实现，不得单独证明兼容。

## 4. 相对 v0.1.3 的可审计差集摘要

### 4.1 转为 include 计数

| 差集项 | v0.1.3 状态 | 本表状态 | 计数 |
|--------|-------------|----------|------|
| 能力域 | 7 include + 4 include-partial + 1 exclude | **12 include** | 4 个 partial 域（D-COMP/D-ACT/D-TABLE/D-FORM）→ include；1 个 exclude 域（D-UPLOAD）→ include |
| registry type | 18/24 白名单 | **24/24 include** | +6 type（statCard、chart、inputNumber、datePicker、dateRangePicker、upload） |
| reactions fixture | 16/16 排除 | **16/16 include** | +16 case |
| request-construction | 64/75（batch 排除） | **75/75 include** | +11 batchRequest case |
| uploads fixture | 未 vendor | **13/13 include** | +13 case（S1 已 vendor + SHA pin） |
| 后端服务端契约 | 无批量/上传端点 | 批量端点 + 上传端点 | +2 端点族 |

**累计**：v0.1.3 纳入面（7 include 域 + 4 partial 域已纳入子面 + **280 case** + 18 type）**保持不回退**；整份契约新增纳入 = 5 域升格/新增（D-COMP/D-ACT/D-TABLE/D-FORM/D-UPLOAD 全量）+ 6 registry type + **40 fixture case**（reactions 16 + batchRequest 11 + uploads 13）+ 2 后端端点族。16 个行为套件合计 320 case，本表目标 = **320/320 全绿**。

### 4.2 仍 residual 清单

**空**（本表无 exclude、无 include-partial、无范围收缩；I-002 = N/A，无用户 residual 需求）。

### 4.3 与 v0.1.3 保持的边界（非收缩，属上游/产品既有边界）

| 边界 | 依据 |
|------|------|
| 业务领域模块（订单/钱包/类目/通知等终端产品） | Charter 非目标；inventory §3 业务样例仅作结构模板并改写 |
| 完整 IAM 产品（SSO 联邦、细粒度审计后台） | D-PERM v0.1.3 既有边界（`include` 主面不受影响） |
| 跨页全选 / 筛选结果全集 / 批量部分成功回填 | ADR-0022 非目标（上游协议边界，非本仓收缩） |
| 多租户应用市场 / 动态远程 manifest 编排 | D-APP v0.1.3 既有边界 |
| 多协议大版本并行矩阵 | D-VER v0.1.3 既有边界 |
| `scenarios` 套件自动化门禁 | Q5 延续：support-only（范例源） |
| 上游 reference-js/python/runner | manifest excluded 参考实现 |

## 5. 实施批次（I-004 输入，Root D-002 采纳）

| 批次 | 范围 | 阶段 |
|------|------|------|
| B1 | D-FORM/D-COMP：inputNumber、datePicker、dateRangePicker、statCard、chart | S2 |
| B2 | D-EXPR：`$deps` 全语法 + 快照多轮 reaction 引擎（16 case） | S2 |
| B3 | D-TABLE/D-ACT：多选 UI + batchMapping/requiresSelection/EMPTY_SELECTION + 11 batchRequest case + Go 批量端点 | S2 |
| B4 | D-UPLOAD：upload 控件 + `type: upload` action + 13 case 编排 + Go 上传端点 + 范例 | S2/S3 |
| B5 | S3 保真：fail-closed 边界测试、表达式/权限边角、升级 example 页面 | S3 |
| B6 | S4 范例 + conformance 登记（本表验证入口列回填） | S4 |

## 6. 用户裁决与冻结证据

| # | 问题 | 用户方向 | 冻结证据 |
|---|------|----------|----------|
| — | S0 差集是否全部纳入（无 exclude/收缩）？ | 无需裁决：差集全部可纳入，默认 include 生效（VP-006 exit 1 纪律；I-002 = N/A） | E-002 §4；本表 §4.2 空 residual 清单 |
| — | S1 冻结独立审计 | 用户指定 provider：**grok build**、模型 **grok 4.5**、思考强度 **high** | 本区 `03-audit/A-001-*.md`（source: independent）及其响应 |

冻结证据：用户此前书面目标「**必须**支持 `schema-ui-docs@v2.7.0` **整份契约**」（2026-08-08，VP-006 用户裁决节）→ 本表 v1.0.0 → Root [D-002](../01-decision/D-002-full-coverage-freeze.md) → `I-PROTO-FULL-001=closed`。

## 7. 变更规则（冻结后）

- 本文件是 v1.0.0 冻结基线；变更覆盖 disposition = 新 Root 决策 + 新版本表；**不得静默改写**。
- S2–S4 允许回填「验证入口」列的具体路径（属执行证据登记，不改变 disposition）；disposition 变更仍须新决策。
- 历史 `I-PROTO-001 v0.1.3` 文件保持只读未改；任何「已完整支持 v2.7.0」声明必须以本表 + 关门证据背书。
