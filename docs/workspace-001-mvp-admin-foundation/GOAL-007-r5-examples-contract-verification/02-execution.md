---
id: GOAL-007-r5-examples-contract-verification
doc: execution
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-08-01
version: 0.8.0
---

# 执行记录 · GOAL-007

## 时间线

> 涉及 P-005 信息项时，记录本次收集/验证的实际动作、I-00N、级别、证据路径，以及新发现的未知。计划中的收集动作必须明确标为计划，不能把 `open`、`deferred` 或 `accepted-residual` 写成已验证事实。

### 2026-07-31 · 目标立项（R5 规划）

- 承接 R4 关门（GOAL-006 `done`，Root `progress` 4/6），经用户 `/govern` 确认 slug 立项 `GOAL-007-r5-examples-contract-verification`。
- 写入五件套；`00-meta` 定义 R5 范围（11 纳入域范例 + 结构/行为验证）与 `I-007-001`（required，R5 验收前）；`01-decision` 记录 D-001 立项与 D-002 路线。
- 同步 `goal-tree.md` 登记本目标（`active`），Root 路线图 R5 标记为「规划中」。
- **未**修改 `apps/*` 实现；范例页/验证路径的收集与实现尚未开始（`I-007-001` 仍 open）。

### 2026-07-31 · R5 阶段 1：契约发现与登记完成

- 响应 A-001（independent，`verdict: pass`）：用户裁决「不需要自审，直接推进」；`01-decision` 记录 D-003。
- 落盘 [attachments/I-007-001-registry.md](attachments/I-007-001-registry.md) v0.1.0：逐纳入域登记范例路径 + 结构/行为验证入口，对齐 [I-PROTO-001 v0.1.3 §3] 与协议清单 §2.5；明确排除 D-UPLOAD、多选批量、完整 registry、scenarios 自动化门禁。
- 核验并登记已有可执行验证入口：`cd apps/web && npm test`（6 测试文件 / 94 项，覆盖 app-manifest / app-navigation / permissions-inheritance / account context / navigation / App integration）、`npm run build`、`cd apps/api && go test ./...`、`go build ./...`；D-APP（App.tsx / navigation.ts / app-manifest.ts / 静态 manifest）与 D-PERM（permissions.ts / account context / Go session/permission/handler）复用产物入账。
- `I-007-001` → `verified`（**登记层面**）：登记表已建立并核验复用命令；阶段 3 的逐域结构/行为验证执行仍未开始，`I-PROTO-003` 保持 open。**未**修改 `apps/*` 实现。

### 2026-07-31 · R5 阶段 2：实施子方案定稿

- 用户裁决「不需要自审，直接推进」（P-004 §3.1 应答），并选择「先落实施子方案再实现」。
- 记录决策 D-004：阶段 2 分三批实施（2a D-DATA/D-TABLE + Go 数据支撑 → 2b D-FORM/D-ACT → 2c D-EXPR + Renderer 接线）；`I-PROTO-004`（vendor vs pin）决策时点定在阶段 3 结构校验实现前，不阻断阶段 2。
- 批次 2a 完成前不修改 `apps/*`；阶段 2 落地以 03 逐批记录证据，全部完成后再进入阶段 3。

### 2026-07-31 · R5 阶段 2：批次 2a 实施完成（D-DATA / D-TABLE）

**Go 列表/详情数据 API**（`apps/api/internal/handler/records.go`）：
- 新增 `GET /api/records`（查询参数 `q` / `sort` / `order` / `page` / `pageSize`，`sort` 白名单 name/status/owner/updatedAt）与 `GET /api/records/{id}`（404 `RECORD_NOT_FOUND`）。静态 dev 数据集 8 条（`staticRecords()`，去业务化）；错误 fail-closed（非法参数 400）。
- 注册于 `handler.Register`（health.go）；新增 `records_test.go` 6 项测试（默认列表、搜索、排序 desc、分页、非法参数、详情/404）。`go test ./...` 全绿（**21 项**）、`go build ./...` 通过。

**Web 数据客户端与表格组件**（`apps/web/src/`）：
- `renderer/records.ts`：`RecordItem`/`RecordList` 类型、`buildRecordsQuery`（query-serialization）、`parseRecordList`（response-mapping，fail-closed）、`fetchRecords`（request-construction）。新增 `records.test.ts` 11 项。
- `components/data-table.tsx`：泛型可排序表格（`DataTable<T>`，aria-sort、asc/desc 切换、加载/错误/空态、自定义渲染）。新增 `data-table.test.tsx` 6 项。
- `renderer/use-records.ts`：`useRecords` hook（query 状态、请求、loading/error、取消）。

**范例页与接线**（`apps/web/src/app/`）：
- `examples/data-table-page.tsx` 与 `examples/search-form-table-page.tsx`（搜索表单 + 表格）；`examples/registry.tsx` 按 pageId 分发到 `DataTable`/`SearchFormTable`；`App.tsx` `PageSurface` 命中范例页时渲染真实页面，否则维持 manifest fallback。
- `public/.well-known/schema-ui/app-manifest.json` 登记 `data-table`（`/data-table`）与 `search-form-table`（`/search-form-table`）页面，sidebar 新增「Examples」分组；`App.tsx` icon registry 补 `table`/`search`。
- 静态 manifest SHA 变更 → 同步 `upstream-fixtures.test.ts` 的 `STATIC_MANIFEST_SHA256`；`app-manifest.test.ts` 页面清单断言更新。

**测试 / 构建 / 运行时证据**：
- `npm test` 全绿（**9 测试文件 / 114 项**）；`npm run build`（tsc + vite）通过。
- 运行时实测（起 `apps/api` + `apps/web` dev，Edge headless）：`/data-table` 渲染真实数据（Acme Console / Northwind Sales / Hooli Connect / Umbrella Ops…，含 active/pending 状态、alice 所有者），`/search-form-table` 渲染「8 records」+ `rec-8` Globex Admin；两页均无 Loading/Error。API 冒烟：list 默认（8 条）、`q=alice`（3 条）、`sort=updatedAt&order=desc`、detail `rec-3`、404 均正确；Vite 代理 `/api` 正常。
- 登记表 [attachments/I-007-001-registry.md](attachments/I-007-001-registry.md) 升 **v0.2.0**：D-DATA/D-TABLE 行「现状」→ 已实现，验证入口命令与证据入账，阶段 3 待接入清单移出已落地映射路径。

**未改动**：`I-PROTO-003`（父目标，open）与 `I-PROTO-004`（open）状态；Root `progress` 4/6；批次 2b（D-FORM/D-ACT）、2c（D-EXPR/Renderer）未开始。

### 2026-07-31 · R5 阶段 2：批次 2b 实施完成（D-FORM / D-ACT）

**D-FORM 控件表面**（`apps/web/src/renderer/form-controls.ts`）：
- `FormControlType` 白名单（§5 冻结：base input/select + 2.6 extended textarea/switch/checkbox/radio + 2.7 advanced cascader/checkboxGroup/richText/password）+ `isWhitelistedFormControl`；`wireKindOf`（string/boolean/string-array 映射）、`coerceFieldValue`（defaultValue 应用 + 类型强转）、`validateDefaultValue`（wire 类型失配 fail-closed）。
- `checkFormCapabilities` / `checkFormCapabilitiesRaw`：版本/capability 门禁（extended ≥2.6 + `form.controls.extended`；advanced ≥2.7 + `form.controls.advanced`；select multiple ≥2.6；非白名单 `FORM_TYPE_NOT_WHITELISTED` fail-closed）。新增 `form-controls.test.ts` 13 项。
- `form-controls.tsx`：`FormControls` 组件渲染全部 10 种白名单控件（input/password/select 单多选/radio/checkboxGroup/cascader/switch/checkbox/textarea/richText）。

**D-ACT 非批量动作**（`apps/web/src/renderer/row-action.ts`）：
- `runRowAction` 包装 R4 `executeAction` 时序引擎：`visible`/`confirm`/`confirmed`/`disabled`/`requiresSelection` 参数透传，返回 outcome/reason/permissionDenied/confirmed。新增 `row-action.test.ts` 5 项（执行/权限拒绝/隐藏 NOT_VISIBLE/确认取消+确认/disabled）。

**Go 编辑生命周期支撑**（`apps/api/internal/handler/records.go`）：
- 新增 `PATCH /api/records/{id}`（指针字段部分更新 + `validatePatch` 空值 fail-closed，`INVALID_PATCH_BODY`/`INVALID_PATCH_FIELD`）与 `DELETE /api/records/{id}`（204 `NoContent`）。`recordHandler` 增加 `sync.RWMutex` 保护共享数据集（MVP 无 DB）。`records_test.go` 新增 5 项（update、update invalid、update 404、delete、delete 404）。`go test ./...` 全绿（**18 顶层 / 26 含 Evaluate 子测试**）、`go build ./...` 通过。

**Web 客户端编辑/删除**（`apps/web/src/renderer/records.ts`）：
- `updateRecord`（PATCH，部分字段 + fail-closed 形状校验）与 `deleteRecord`（DELETE，204 成功）。`records.test.ts` 新增 4 项（update 成功/400、delete 204/404）。

**范例页与接线**（`apps/web/src/app/`）：
- `examples/form-controls-page.tsx`：全部白名单控件 + 版本/capability 门禁展示 + Serialize values。
- `examples/list-edit-lifecycle-page.tsx`：列表 + 行编辑（Edit/Delete 经 D-ACT `runRowAction` 门禁，PATCH 保存、DELETE 确认）+ `DataTable` 复用。
- `examples/registry.tsx` 注册 `form-controls`/`list-edit-lifecycle`；`App.tsx` icon registry 补 `form`/`pen`；`app-manifest.json` 登记两页面 + sidebar「Examples」组；`app-manifest.test.ts` 页面清单断言同步；`upstream-fixtures.test.ts` `STATIC_MANIFEST_SHA256` 更新；`app-examples.test.tsx` 新增 2 项表面测试（list-edit-lifecycle Edit/Delete 门禁、form-controls capability 门禁通过）。

**测试 / 构建证据**：
- `npm test` 全绿（**11 测试文件 / 138 项**）；`npm run build`（tsc + vite）通过（修复一处 TS6133 未用变量后）。
- 登记表 [attachments/I-007-001-registry.md](attachments/I-007-001-registry.md) 升 **v0.3.0**：D-FORM/D-ACT 行「现状」→ 已实现，验证命令与证据入账。

**未改动**：`I-PROTO-003`（父目标，open）与 `I-PROTO-004`（open）状态；Root `progress` 4/6；批次 2c（D-EXPR/Renderer）与阶段 3/4 未开始。

### 2026-08-01 · R5 阶段 2：批次 2c 实施完成（D-EXPR / D-COMP）

**D-EXPR 反应引擎**（`apps/web/src/renderer/reactions.ts`）：
- `ReactionRule`/`ReactionApply`（fieldId + 显式 visible/disabled 布尔）、`parseReactionRule`（fail-closed：非对象、缺 id/when、非法表达式 `REACTION_EXPRESSION_INVALID`、apply 缺 fieldId/非布尔 `REACTION_APPLY_INVALID`）、`evaluateReactions`（复用 `evaluateExpression`，frozen $context 子集；未知 apply 字段 `REACTION_APPLY_FIELD_UNKNOWN` fail-closed）、`parseAndEvaluateReactions`。新增 `reactions.test.ts` **12 项**。
- `evaluateExpression` 表达式语法校验从 `app-manifest.ts` 以 `isValidExpression` 导出复用（同一 frozen 语法，不重复实现）。

**D-COMP 最小 Renderer 接线**（`apps/web/src/renderer/render.ts` / `render.tsx`，resolve R4 推荐跟踪项 F-002）：
- `render.ts`：`parseRenderNode`（whitelist form/section/table，未知 type `RENDER_UNKNOWN_NODE_TYPE` fail-closed，form 缺 fields `RENDER_FORM_FIELD_INVALID`）、`collectFieldIds`、`resolveFormReactions`、`gateAction`（布尔或 $context 表达式）、`tableActionGate`（visible/disabled 独立求值）。新增 `render.test.ts` **11 项**。
- `render.tsx`：`RenderPage` 分发层——form 经 `FormControls` 渲染并应用 reaction 状态（隐藏/禁用字段），section 容器递归，未知 type 出 `role="alert"`；`FormControls` 增加 `fieldDisabled` 按字段禁用（向后兼容）。新增 `render.test.tsx` **4 项**（默认全显、reaction 隐藏、reaction 禁用、未知 type fail-closed）。
- 这是「Renderer 接线」落地：`render.tsx` 消费 D-EXPR（reactions）与 D-FORM（FormControls），resolve GOAL-006 F-003/F-002 跟踪项（renderer 集成层消费引擎）。

**范例页与接线**（`apps/web/src/app/`）：
- `examples/form-with-reactions-page.tsx`：D-EXPR 范例——Admin/Viewer 角色 + audit feature 切换改变 `$context` 快照，`RenderPage` 重新应用字段显隐/禁用（approval 仅 admin 可见、auditNote 无 audit 时禁用）；页面展示 `JSON.stringify(context)`。
- `examples/registry.tsx` 注册 `form-with-reactions`；`App.tsx` icon registry 补 `reaction`（Zap）；`app-manifest.json` 登记页面 + sidebar「Examples」组；`app-manifest.test.ts` 页面清单断言同步；`upstream-fixtures.test.ts` `STATIC_MANIFEST_SHA256` 更新（新值 `0475f7bb…5e120`）；`app-examples.test.tsx` 新增 1 项表面测试（context 快照 + 表单字段渲染）。

**测试 / 构建 / 运行时证据**：
- `npm test` 全绿（**14 测试文件 / 166 项**，+28：reactions 12 + render 11 + render.tsx 4 + app-examples 1）；`npm run build`（tsc + vite）通过（修复一处 TS2322 section 节点类型 + 一处 TS2352 字段转换）。
- Edge headless 实测 `/form-with-reactions`：渲染 Admin/Viewer/audit 切换、全部 4 个字段（Name/Kind/Approval/Audit note）、无 `role="alert"` 错误。Go 侧无改动（`go test ./...`/`go build ./...` 保持通过）。
- 登记表 [attachments/I-007-001-registry.md](attachments/I-007-001-registry.md) 升 **v0.4.0**：D-EXPR/D-COMP 行「现状」→ 已实现，验证命令与证据入账，P2 `form-with-reactions` 标记已实现。

**未改动**：`I-PROTO-003`（父目标，open）与 `I-PROTO-004`（open，non-blocking，阶段 3 结构校验前决策）状态；Root `progress` 4/6；阶段 3（结构/行为验证）与阶段 4（验收/关门）未开始。

### 2026-08-01 · 响应 A-005：F-001～F-004 以 fixed 闭合 + F-005 同步

用户裁决「不需要自审，直接推进」（P-004 §3.1），按 **`fixed`** 路径响应 A-005（independent, conditional）4 项 required findings（D-007 已留痕）。直接读代码核验了 A-005 的主张（`render.ts` `gateAction` fail-open、`render.tsx` `FormView` 未执行 D-FORM 门禁、`checkFormCapabilities` 无 defaultValue 2.7+advanced 门禁、renderer whitelist 窄于 §5 初表）后实施修正：

**F-001 · action 表达式 fail-closed**（`apps/web/src/renderer/render.ts`）：
- 以 `resolveActionGate(expression, context, absentDefault, path)` 替代 `gateAction`：`undefined`/`null` → absent default；boolean → 直通；合法 `$context` 表达式 → 求值；**显式非法**字符串或非表达式值 → `{ kind: "error", code: "ACTION_GATE_EXPRESSION_INVALID" }` fail-closed，不再静默返回 `defaultValue`。
- `tableActionGate` 改为返回 `{ visible, disabled, errors }`：显式非法 `visibleWhen` → `visible: false` + error；显式非法 `disabledWhen` → `disabled: true` + error；缺省门禁不产生 error。
- 测试（`render.test.ts`）：重写 gate 断言为 `resolveActionGate`/`tableActionGate` 新 API；补非法 visible/disabled 回归测试与「缺省 vs 显式非法」区分断言。

**F-002 · RenderPage 执行 D-FORM 门禁**（`render.ts` / `render.tsx`）：
- 新增 `gateRenderFormFields(metaValue, rawFields, path)`：解析原始字段为 `FormControlField`，校验 §5 whitelist（未知 type → `FORM_TYPE_NOT_WHITELISTED`），并按字段运行 `checkFormCapabilitiesRaw`（version/capability/defaultValue），只返回通过门禁的字段，拒绝字段以确定错误返回。
- `FormView`（`render.tsx`）渲染前调用 `gateRenderFormFields`：门禁错误以 `role="alert"` 列出，受影响字段不渲染（不再直接强转后交给 `FormControls` 静默渲染）。
- 测试（`render.test.tsx`）：fixture meta 补 `form.controls.extended`（此前缺 capability 却断言 textarea/switch 正常渲染——即 A-005 所指）；新增「缺 capability 字段被拒并出 `FORM_CAPABILITY_REQUIRED`」「未知 type 被拒并出 `FORM_TYPE_NOT_WHITELISTED`」。

**F-003 · defaultValue 2.7 + advanced 双门禁**（`apps/web/src/renderer/form-controls.ts`）：
- `checkFormCapabilities` 字段循环为任一 `field.defaultValue !== undefined` 增加：`protocolVersion >= 2.7`（否则 `FORM_VERSION_TOO_LOW`）+ `requiredCapabilities` 含 `form.controls.advanced`（否则 `FORM_CAPABILITY_REQUIRED`），与既有 wire 类型校验叠加。
- 测试（`form-controls.test.ts`）：补 2.6（version too low）、2.7 缺 advanced capability、2.7 完整 capability（pass）三组回归；修正原「base meta 接受 input defaultValue」的 fail-open 断言。

**F-004 · Renderer whitelist 对齐冻结 §5**（`render.ts` / `render.tsx`）：
- `RenderNodeType` 从 form/section/table **扩展至冻结 §5 全部 node type**：layout `grid`/`section`/`tabs`；data/action `text`/`table`/`recordView`/`actionButton`；`form`。`parseRenderNode` 与 `RenderPage`/dispatch 层为每个 type 落渲染（GridView/TabsView/TextView/RecordView/ActionButtonView），未知 type 仍 fail-closed 出 alert。
- **未做范围缩小** → 不需 Root 覆盖表修订或新版 v0.1.3 冻结；A-005 F-004 所指「静默窄于初表」以补齐缺口闭合。
- 测试：`render.test.ts` `isWhitelistedNodeType` 断言 8 个 §5 type 全 `true`（非 §5 如 chart/modal/upload 仍 false）；`render.test.tsx` 新增 §5 全 node dispatch 渲染测试。

**F-005 · Root 投影同步**：Root `00-meta.md` R5 行更新为「阶段 1 + 批次 2a/2b/2c 完成（阶段 2 全部落地），阶段 3/4 未开始」；`progress` 仍 4/6。

**测试 / 构建 / 运行时证据**：
- `npm test` 全绿（**14 测试文件 / 173 项**，+7：render.test.ts 重写 + render.test.tsx 新增 3 项）；`npm run build`（tsc + vite）通过（修复 option 类型收窄、`node.props` optional 链、tabs label 类型断言）。
- `go test ./...` / `go build ./...` 通过（Go 侧无改动）。
- Edge headless 实测 `/form-with-reactions`：4 字段（Name/Kind/Approval/Audit note）+ Context snapshot + 无 `role="alert"`；`/form-controls`、`/list-edit-lifecycle`、`/data-table`、`/search-form-table` 均 0 门禁错误。
- 登记表 [attachments/I-007-001-registry.md](attachments/I-007-001-registry.md) D-COMP 行同步为「Renderer whitelist = 冻结 §5 全部 node type」（v0.5.0）。

**未改动**：`I-PROTO-003`（父目标，open）与 `I-PROTO-004`（open，non-blocking，阶段 3 结构校验前决策）状态；Root `progress` 4/6；阶段 3（结构/行为验证）与阶段 4（验收/关门）未开始。

### 2026-08-01 · 响应 A-006 + I-PROTO-004=vendor + 阶段 3 落地

用户裁决「不需要自审，直接推进」；`I-PROTO-004` 选 **vendor**；进入阶段 3（D-008）。

**Vendor 产物**（pin commit `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b`，`artifactVersion` 2.7.0）：

- Schemas → `docs/schemas/`：`node` / `page` / `action` / `reaction` / `component-registry`（叠加既有 `app-manifest`）；SHA 记入 `apps/web/src/protocol/upstream/provenance.json`。
- Fixtures → `apps/web/src/protocol/upstream/*.cases.json`：actions、component-format、query-serialization、reactions、request-construction、request-lifecycle、response-mapping、runtime-defaults、search-table、static-data、table-sort、version-negotiation（叠加既有 app-manifest / app-navigation）。

**结构校验**（`apps/web/src/protocol/conformance/schema-validate.ts`，Ajv draft-07）：

- `validateAgainstSchema('page'|'node'|…)`；§5 白名单 sample page 可过 `page`/`node` schema；缺 type / 缺 protocolVersion 拒绝。

**行为 fixture 适配器**（`apps/web/src/protocol/conformance/`）：

| Suite | 处理 | 说明 |
|-------|------|------|
| component-format | 全量执行 5 | 格式 wire 类型无强制转换 |
| query-serialization | 全量执行 16 | ADR-0010 序列化 |
| static-data | 全量执行 9 | static/ref 形状 |
| request-lifecycle | 全量执行 4 | latest-wins / hide-drop |
| runtime-defaults | 全量执行 9 | baseURL / defaults / form init |
| response-mapping | 全量执行 23 | table/chart/formRecord/recordView |
| search-table | 全量执行 11 | 四层 merge + selection |
| table-sort | 全量执行 14 | 三态 sort + L2 validate |
| version-negotiation | 全量执行 44 | 严格版本 + capability |
| actions | 全量执行 11 | 非批量 transport→events |
| reactions | 全量记账排除 16 | 上游 `$deps` 字段写引擎 ∉ MVP `$context` 子集 |
| request-construction | 全量记账 75 | batch Q1 排除；其余 deferred 统一引擎（partial 由 records/row-action 覆盖） |

**测试 / 构建证据**：

- `cd apps/web && npm test` → **15 文件 / 326 项** 全绿（+153 stage3 + 既有 173）。
- `cd apps/web && npm run build`（tsc + vite）通过。
- `cd apps/api && go test ./...` / `go build ./...` 通过（Go 侧无改动）。
- 依赖：`ajv@8`（devDependency，schema 校验）。

**登记表**：`attachments/I-007-001-registry.md` 升 **v0.6.0**——阶段 3 验证入口命令与逐 suite 执行状态入账；§4 复用计数与 §2 对齐（A-006 F-003）。

**未改动**：父目标 `I-PROTO-003`（required，open，R5 验收/关门门禁）；Root `progress` 4/6；阶段 4（验收/关门）未开始。`I-PROTO-004` → verified（Root 同步）。

## 待办

1. ~~**批次 2c**：D-EXPR `form-with-reactions` + Renderer 接线（D-COMP，resolve F-002）。~~ **完成（2026-08-01）**。
2. ~~落地 `node`/`page` schema 校验与已纳入 fixtures 对照；`I-PROTO-004` vendor。~~ **完成（2026-08-01，阶段 3）**。
3. 闭合父目标 `I-PROTO-003` 并完成 R5 自审/关门（阶段 4）。
4. 可选：补齐 `request-construction` 统一 host 引擎；MVP `$context` reactions 与上游 suite 的差异文档化给验收。

## 进度评估

**阶段 1–3 完成**（契约发现与登记；范例实现；结构/行为验证）。**阶段 4（验收与关门）未开始**。进度仅为展示，不放行验收、不推导 `done`，不抬升 Root `progress`（仍 4/6）。
