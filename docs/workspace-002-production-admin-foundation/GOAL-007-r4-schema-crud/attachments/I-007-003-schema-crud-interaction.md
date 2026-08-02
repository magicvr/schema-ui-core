---
title: I-007-003 · Schema CRUD 读写交互契约
status: active
doc_type: info-collection
created: 2026-08-02
updated: 2026-08-02
parent: GOAL-007-r4-schema-crud
version: 0.2.2
related_info: I-007-003
related_decision: D-005
---

# I-007-003 · Schema CRUD 页面/Node/action 绑定与交互/权限契约

> **结论**：本附件与 D-005 关闭 `I-007-003`，并完成 S4/S5 的交互契约冻结。以下页面绑定、字段映射、交互状态与权限矩阵是 S4/S5 实施输入，**不是**已交付的 Schema 写交互代码或页面事实。
> **扫描日期**：2026-08-02。工作区 `shared_materials_catalog: none`；事实来自本仓库 Web Renderer、permission/action 引擎、records client、schema fixtures 与既有测试。
> **修订（v0.2.0 · 响应 A-005）**：闭合 F-002/F-003/F-004（required）并处理 R-001/R-002——§2 冻结**唯一**页面结构（modal `create-form` + modal `edit-form` + 行 `delete`，各自独立 action/submitAction）；§5 冻结 `records.write` → `permissions.edit/delete` 的**唯一**表达式写法与 cascade 挂载点，禁止「仅 `permissionIntent` 无表达式」；新增 §9 最小冻结实现规格（顶层 actions、capabilities、`$row` 绑定、search 归属、预填来源、Renderer 文件白名单）。修订决策见 D-005 补记。
> **修订（v0.2.1 · 响应 A-006）**：闭合 F-005/F-006（required）——§9.1 改写为对齐 `action.schema` 的 **5 个顶层 action**（3×RequestAction + 2×ModalAction），`onSuccess` 用 **`behavior`**、挂载字段用 **`actionRef`**、confirm 移到 **rowAction** 项、delete 的 `requestMapping.path.id` 留在 rowAction；新增 §9.1a 冻结 **form submit 的 `{id}` 槽绑定**（从打开 modal 时捕获的行上下文解析，`formAction` 的有界扩展，落入 §9.5 白名单并补测试）；§9.5 白名单扩展允许一次性新增 modal/confirm 渲染文件（R-001）。修订决策见 D-005 v0.2.1 补记。
> **修订（v0.2.2 · 响应 A-007）**：闭合 F-007（required）——§9.2 delete 的 `confirm` 由 `{ text }` 对象改为 **string**（`confirm: "Delete this record?"`，与 registry `table.props.actions[].confirm: string` 一致）；§9.5 白名单补入 `request-construction.ts` / `row-action.ts`（R-001）；§9.1 注明 `reload` 隐含关闭 modal（R-002）。修订决策见 D-005 v0.2.2 补记。

## 1. 当前可继承基线（现状证据）

| 事实 | 证据 | R4 继承 |
|------|------|---------|
| Renderer node whitelist：layout（grid/section/tabs）、data/action（text/table/recordView/actionButton）、form；未知 node 类型 fail closed | `apps/web/src/renderer/render.tsx` `dispatchParsedNode` | S4 沿用；不新增 node 类型 |
| 表格渲染 `SchemaTable`：从 `props.dataSource`（默认 `/api/records`）加载，提供 loading / error / empty 与列排序、分页计数；无行操作、无 toolbar | `schema-table.tsx`、`components/data-table.tsx` | **一次性补齐** table `props.actions`/`props.toolbar` 渲染与 form-to-query 绑定 |
| 表单渲染 `FormControls`：冻结控件白名单（input/select/textarea/switch/checkbox/radio/cascader/checkboxGroup/richText/password），wire kind 与 capability 门禁 | `form-controls.ts`、`form-controls.tsx` | 字段映射复用白名单；S4 增加「submit → API」绑定 |
| `$context` 表达式/反应引擎：`$context.user.<path>` / `$context.features.<path>` 与 `==`/`!=`/`contains`；reactions 翻转字段 visible/disabled | `protocol/app-manifest.ts`（evaluateExpression）、`renderer/reactions.ts` | S4 交互与权限表达式沿用同一冻结语法 |
| permission 执行引擎：`permissionCascade`/`permissionIntent`（edit/delete）+ `executeAction` 序列（visible → permission → disabled/requiresSelection → confirm）；target kinds：formField/formSubmit/rowAction/toolbarTrigger/actionButton/column | `renderer/permissions.ts` | **复用**作为 create/edit/delete 的 UI 门禁与确认入口 |
| 行 action / toolbar 目标已在 permission 引擎建模（`table.props.actions`/`props.toolbar` 校验与求值）但**渲染器尚未渲染** | `permissions.ts` `collectTargets`（rowAction/toolbarTrigger） | S4 渲染层补齐 |
| records client：`fetchRecords`/`parseRecordList`/`updateRecord`（PATCH）/`deleteRecord`（DELETE）；**尚无 createRecord（POST）** | `renderer/records.ts` | S4 新增 createRecord 客户端函数 |
| 现有页面 fixtures：`catalog`（table 3 列）、`data-table`、`search-form-table`（search form + table 5 列）、`form-controls`、`form-with-reactions`（reactions）、`list-edit-lifecycle`（recordView + edit form，**仅展示**，无 submit 绑定/行操作/确认） | `apps/api/internal/handler/fixtures/schema/*.json` | `list-edit-lifecycle` 演进为代表性 CRUD 页；`search-form-table` 纳入 form-to-query 绑定 |
| 菜单/导航投影：`me.features`（`menu_list_edit_lifecycle` 等）经 `visibleWhen` 门控导航 | `auth.FeaturesForUser`、`navigation.ts`、`app-manifest.ts` `isNavigationItemVisible` | S4 菜单门禁不变；写 affordance 走 `$context.user.permissions` |
| 会话上下文：`me` 返回 `{ user: { id, name, roles, permissions }, features }`；`NavigationContext = { user, features }` | `account/session.go`、`handler/account.go`、`protocol/app-manifest.ts` | 权限表达式数据源 |
| 后端 CRUD API 已冻结并实现（I-007-001/002 + S3）：GET/POST/PATCH/DELETE、统一 envelope、`records.read`/`records.write` | `apps/api/internal/handler/records.go` | S4 客户端调用该 API；403/401 为权威负向 |

## 2. 代表页面与 Node/action 绑定（冻结）

**页面**：`list-edit-lifecycle`（fixture `list-edit-lifecycle.json`）为唯一代表性 CRUD 生命周期页；与 R3 种子菜单 `list-edit-lifecycle`（feature `menu_list_edit_lifecycle`，admin 可见）对齐。**搜索不并入本页**：form-to-query 绑定归 `search-form-table` 页（见 §9），避免本页出现第二套查询入口。

### 2.1 唯一页面结构（A-005 F-002 闭合）

本页 body 由**一个 table** + **两个 modal（各含一个 default form）** + **行删除确认**组成。每个写操作有**独立** action / `submitAction`，**禁止**用单个 form 的 `submitAction` 同时表达 POST 与 PATCH：

| 结构 | 载体 | action / submitAction | HTTP |
|------|------|----------------------|------|
| 列表 | `table`（`props.dataSource: "/api/records"`，列 name/status/owner/updatedAt + 操作列） | 无 | `GET /api/records` |
| 新建 | toolbar `create`（toolbarTrigger，`permissionIntent: edit`）→ **modal `create-form`**（default form，空表单） | `submitAction: "createRecord"` | `POST /api/records` |
| 编辑 | rowAction `edit`（`permissionIntent: edit`）→ **modal `edit-form`**（default form，按选中行**预填**） | `submitAction: "updateRecord"` | `PATCH /api/records/{id}` |
| 删除 | rowAction `delete`（`permissionIntent: delete`）→ 确认 → 请求 | `deleteRecord`（`actions.row.request`） | `DELETE /api/records/{id}` |
| 详情 | `recordView`（渲染**选中行拷贝**，只读） | 无 | 不发起 GET |

**预填来源（A-005 R-002 闭合）**：create-form 为空；edit-form 与 recordView 均使用**列表中已加载的选中行拷贝**（单一数据源 = 列表），**不**使用 `form.recordSource` / 独立 GET；因此 `form.record.load` / `record.view.load` **不在**所需 capabilities 内。列表刷新后清空选中态。

**一次性渲染补齐（S4 实施范围，非每页）**：渲染 table `props.actions`/`props.toolbar`、modal 内容与确认；绑定 form `submitAction` → 对应顶层 action；成功/错误/确认反馈；行选中传递。此后「新增或调整代表页面」仅改 fixture，**不修改 Renderer 主路径代码**（S4 成功标准）。允许改动的文件白名单见 §9（R-001）。

## 3. 字段映射（冻结）

| 字段 | 控件 | wire kind | create body | PATCH body | 约束 |
|------|------|-----------|-------------|------------|------|
| `id` | —（表格/详情展示） | string | 服务端生成，客户端不得传 | 不适用 | 只读 |
| `updatedAt` | —（表格/详情展示） | string | 服务端写入 | 服务端刷新 | 只读；RFC3339 含毫秒 |
| `name` | input | string | **必填** | present 键 | trim 后非空 |
| `status` | select（options active/pending/archived，**仅 UI 提示**） | string | **必填** | present 键 | trim 后非空；**API 非枚举** |
| `owner` | input | string | **必填** | present 键 | trim 后非空 |

- create body：`{ name, status, owner }`（全必填；缺/非字符串/空 → `INVALID_CREATE_FIELD`）。
- PATCH body：仅 present 键 `{ name?, status?, owner? }`（空 → `INVALID_PATCH_FIELD`）。
- 表单提交前对 name/owner 做 trim；status 以 select 约束合法选项（客户端层），后端仍按 I-007-001 接受任意非空 string。

## 4. 交互状态（成功 / 加载 / 空态 / 错误 / 确认）

| 状态 | 行为 |
|------|------|
| 加载 | 列表 `DataTable` loading 态；create/edit 提交中禁用提交按钮并显示进行中 |
| 空态 | 列表 "No records match."（DataTable emptyMessage）；详情/表单不适用 |
| 成功 | create → 成功提示 + 列表刷新（新行可见）；edit → 成功提示 + 该行刷新（updatedAt 更新）；delete → 成功提示 + 行移除 |
| 错误 | 统一 envelope `{error,message}` → `role=alert` 页级/表单级提示；code 映射见下 |
| 确认 | delete 复用 `executeAction` `confirm: true`：未确认 → `CONFIRM_CANCELLED`，不发起 DELETE |

**错误 code 映射**（稳定意图）：

| code | 呈现 |
|------|------|
| `UNAUTHENTICATED` | 会话失效 → 引导重新登录（LoginPage）；不把 401 当普通表单错误 |
| `FORBIDDEN` | 权限不足提示；页面保持只读 |
| `INVALID_CREATE_FIELD` / `INVALID_PATCH_FIELD` | 表单字段级错误（message 指明字段） |
| `INVALID_CREATE_BODY` / `INVALID_PATCH_BODY` | 表单整体错误（body 非 JSON/超限） |
| `RECORD_NOT_FOUND` | 行已删除/陈旧 → 提示 + 列表刷新 |
| `INVALID_SORT_*` / `INVALID_PAGE*` | 列表查询错误提示 |

## 5. 权限矩阵（冻结）

**后端为权威**：`records.read` 门禁 GET；`records.write` 门禁 POST/PATCH/DELETE；匿名 401、缺权限 403。**前端隐藏/禁用仅为 UX，不是安全边界**（S5：后端 403 不被前端隐藏替代）。

| 身份 | 列表/详情 | 新建/编辑/删除 | 备注 |
|------|-----------|----------------|------|
| admin（records.read + records.write） | 可见 | 可见 + 可操作 | 菜单 feature 可见 |
| editor / viewer（records.read） | 可见（只读） | affordance 隐藏或禁用 | 直接调用 API 仍 403（T-API-09） |
| anonymous | 无会话 → LoginPage | 无会话 | 直接 API 401（T-API-08） |

**表达式门控（冻结 `$context` 语法 · A-005 F-003 闭合）**：

- **唯一推荐写法**：写权限统一挂在 **table 祖先**（行/toolbar action 的 mount 点）**一次**，行/toolbar 仅标 `permissionIntent`；**禁止「仅 `permissionIntent` 无表达式」**（intent 不是 API 权限键，无表达式时有效权限默认为 true）。

  ```jsonc
  // table 节点（mount 点）：records.write → DSL edit/delete 一次声明，行/toolbar 经 cascade 继承
  "permissionCascade": { "keys": ["edit", "delete"] },
  "permissions": {
    "edit":   "$context.user.permissions contains \"records.write\"",
    "delete": "$context.user.permissions contains \"records.write\""
  }
  // rowAction edit → { "key": "edit", "permissionIntent": "edit" }
  // rowAction delete → { "key": "delete", "permissionIntent": "delete" }
  // toolbar create → { "key": "create", "permissionIntent": "edit" }
  ```

- **两个 modal form 各自声明**：modal content 是**新 permission 根**（`collectTargets` newRoot），不继承页面 body 的 cascade——`create-form` / `edit-form` 均须自身携带写表达式：

  ```jsonc
  // form 节点（modal content 内）
  "permissionCascade": { "keys": ["edit"] },
  "permissions": { "edit": "$context.user.permissions contains \"records.write\"" },
  "props": { "submitAction": "createRecord" /* 或 "updateRecord" */ }
  ```

  formSubmit target（`isDefaultForm && typeof submitAction === "string"`）为隐式 edit 意图，经该表达式求值；`records.write` → DSL `edit`/`delete` 键，API 无独立 delete 权限键（delete 与 edit 同用 `records.write`）。
- 读 affordance：页面访问由菜单 feature（`menu_list_edit_lifecycle`）+ 后端 401/403 把关；如需页级只读门禁，可在 body 顶层 section 加 `permissions.view: "$context.user.permissions contains \"records.read\""`（可选，不改变后端权威）。
- 菜单：`$context.features.menu_list_edit_lifecycle`（已实现，不变）。
- 使用 permission 字段的页面 `meta.requiredCapabilities` 必须含 `permissions.inheritance`；协议版本 ≥ 2.3（`permissions.ts` 冻结门禁）。

## 6. 测试矩阵（S4/S5 验收最低断言，T-UI-01～10）

| ID | 断言 |
|----|------|
| T-UI-01 | 代表页列表从 `/api/records` 加载：loading → ready；空表 → "No records match." |
| T-UI-02 | 列排序头切换 sort/order；分页计数正确显示 |
| T-UI-03 | `search-form-table` 页：搜索 form 提交 → 列表按 `q` 过滤（form-to-query 绑定） |
| T-UI-04 | `list-edit-lifecycle` 新建 form（modal create-form，name/status/owner）→ POST → 201；成功后列表含新行；空白字段 → `INVALID_CREATE_FIELD` 表单错误 |
| T-UI-05 | `list-edit-lifecycle` 编辑 form（modal edit-form，选中行预填）→ PATCH → 200；该行 name/status/owner 更新且 `updatedAt` 刷新 |
| T-UI-06 | `list-edit-lifecycle` 删除行 → 确认框 → 确认 → DELETE 204 → 行移除；取消 → 无请求 |
| T-UI-07 | API 错误 envelope 正确呈现（`INVALID_CREATE_FIELD` / `RECORD_NOT_FOUND` / `FORBIDDEN`） |
| T-UI-08 | admin 可见全部 affordance；viewer/editor 只读（隐藏/禁用写 affordance，`records.write` 表达式求值为 false）；匿名重定向登录 |
| T-UI-09 | 后端权威不被前端隐藏替代：即便隐藏，直接 POST/PATCH/DELETE 由 API 403/401 拦截（API 侧 T-API-08/09 已承担；UI 断言呈现一致） |
| T-UI-10 | 新增/调整代表页面不修改 Renderer 主路径代码——页级变更仅允许触碰 §9.5 白名单外的 fixtures；用文件 diff 断言 |

## 7. 与 I-007-001/002 的接口

- 客户端调用已冻结并实现的 records API（I-007-001 + S3）：POST/PATCH body 形状、201/200/204、稳定 error code、`updatedAt` 毫秒格式均不变。
- 权限键与 `me.features` 投影来自 I-007-002/GOAL-006 已冻结模型；本契约只绑定「哪个页面节点调用哪个 API」与「UI 如何呈现状态/权限」。
- 重启/端到端验收与 `I-007-004` 无耦合；S4/S5 页面交互验收在单元/集成测试完成，进程重启 E2E 归 S6。

## 8. 证据索引

- `apps/web/src/renderer/{render.ts,render.tsx,schema-table.tsx,form-controls.ts,reactions.ts,permissions.ts,row-action.ts,records.ts,use-records.ts}`
- `apps/web/src/components/data-table.tsx`
- `apps/web/src/app/{App.tsx,navigation.ts}`
- `apps/web/src/protocol/app-manifest.ts`、`load-page.ts`
- `apps/api/internal/handler/fixtures/schema/{list-edit-lifecycle,search-form-table,catalog,form-with-reactions}.json`
- `apps/api/internal/handler/records.go`、`account.go`；`apps/api/internal/account/session.go`
- `docs/schemas/{action.schema.json,component-registry.json}`；`apps/web/src/protocol/conformance/request-construction.ts`（`$row.*` / `requestMapping`）
- I-007-001（v0.2.0）、I-007-002（v0.2.0）、D-005

## 9. 最小冻结实现规格（v0.2.0 补 · A-005 F-004 闭合 + R-001/R-002）

### 9.1 顶层 `actions`（`page.actions`，对齐 `action.schema.json`）

**5 个顶层 action**：3 个 RequestAction + 2 个 ModalAction（v0.2.1 · A-006 F-006 闭合）。

| id | 类型 | 形状 |
|----|------|------|
| `createRecord` | RequestAction | `method: POST`，`url: "/api/records"`，`bodyMapping: { name: "name", status: "status", owner: "owner" }`，`onSuccess: { behavior: "reload" }` |
| `updateRecord` | RequestAction | `method: PATCH`，`url: "/api/records/{id}"`，`bodyMapping` 同 create；`onSuccess: { behavior: "reload" }`；**`{id}` 槽由 §9.1a 规则绑定**，**不得**写 `requestMapping`/`$row`（form submit 不解析它们） |
| `deleteRecord` | RequestAction | `method: DELETE`，`url: "/api/records/{id}"`，`onSuccess: { behavior: "reload" }`；**无** `confirm`、**无** `requestMapping`（两者都在 rowAction 项上，见 §9.2） |
| `openCreate` | ModalAction | `type: "modal"`，`content`: default form（create-form，空表单，`submitAction: "createRecord"`） |
| `openEdit` | ModalAction | `type: "modal"`，`content`: default form（edit-form，选中行预填，`submitAction: "updateRecord"`） |

- `onSuccess` 用 **`behavior`**（enum: `toast` \| `navigate` \| `reload` \| `closeModal`；`OutcomeBehavior` `additionalProperties: false`），**无 `type` 字段**。
- **create/edit/delete 均 `behavior: "reload"`**：S4 实现注记——**`reload` 隐含关闭当前 modal 并清空选中态**（A-007 R-002），避免 SPA 列表局部刷新后 modal 残留；`closeModal`/`toast` 不单独使用。
- 两 modal 的 form 字段均为 name/status/owner（见 §3）；status 为 select（active/pending/archived 仅 UI 提示）。

### 9.1a form submit 的 `{id}` 槽绑定（A-006 F-005 闭合 · 有界扩展）

`buildFormAction`（`request-construction.ts`）对 form 提交**不做** path 槽 / `$row` 绑定：url 按字面 + baseURL 组合，body 仅来自 formValues/bodyMapping。为支持 edit 提交 `PATCH /api/records/{id}`，冻结规则：

- 当 Renderer 执行 **default form submit** 且 `action.url` 含 `{id}`（或通用 `{slot}`）槽时，从**打开该 modal 时捕获的行上下文**（选中行拷贝）解析该槽，生成最终 url（`PATCH /api/records/rec-1`）。
- 这是对 `formAction` 执行器的**有界扩展**：必须落在 §9.5 白名单文件（`render.tsx` / `request-construction.ts` 等），并补「行上下文槽绑定」单测（T-UI-05 覆盖 `{id}` 解析）。
- **禁止**把 row 专属 `requestMapping`/`$row.*` 写在仅 form 提交的 action 上而不说明执行器（本页 `$row.id` 仅用于 rowAction delete，见 §9.2）。

### 9.2 行级 / 工具栏 / modal 绑定（registry 字段 `actionRef`）

- `table.props.toolbar`：`create` 项 → `actionRef: "openCreate"`，`permissionIntent: "edit"`，label `New record`。
- `table.props.actions`：
  - `edit` 项 → `actionRef: "openEdit"`，`permissionIntent: "edit"`，label `Edit`。
  - `delete` 项 → `actionRef: "deleteRecord"`，`permissionIntent: "delete"`，label `Delete`，**`requestMapping: { path: { id: "$row.id" } }`**（`$row.*` 属 rowAction，由 `buildRowAction` 处理），**`confirm: "Delete this record?"`**（registry `confirm` 为 **string**（v0.2.2 · A-007 F-007）；确认文案在 rowAction 项，与 `executeAction` confirm 序列一致）。
- **字段名**：registry 用 **`actionRef`**（toolbar/actions 的挂载字段），**无** `action` 键、**无** `modal:` 前缀；modal 由顶层 ModalAction（`openCreate`/`openEdit`）承载，toolbar/row 经 `actionRef` 引用。
- 行级 `$row.*` 仅经 rowAction 的 `requestMapping.path`（本页只用 `$row.id`）；列表行数据即 `$row` 源。

### 9.3 `meta.requiredCapabilities` 最小集

`["app.manifest", "app.navigation", "permissions.inheritance", "actions.row.request", "actions.page.trigger", "table.sort"]`
（排序列需要 `table.sort`；**不需要** `form.record.load` / `record.view.load` / `form.controls.extended`——input/select 为基础控件，预填为行拷贝）。

### 9.4 search 归属（F-004 闭合）

form-to-query 绑定（`mode: "search"` + `targetTable`）**只**做进 `search-form-table` 页（T-UI-03 挂该页）；`list-edit-lifecycle` 为生命周期页，**不**含搜索 form。

### 9.5 Renderer 文件白名单（T-UI-10 边界 · A-005 R-001 / A-006 R-001 / A-007 R-001）

「一次性渲染补齐」允许改动的路径：
- `apps/web/src/renderer/{render.tsx, schema-table.tsx, render.ts, records.ts, use-records.ts, form-controls.tsx, row-action.ts}`
- `apps/web/src/components/data-table.tsx`
- `apps/web/src/protocol/conformance/request-construction.ts`（**仅当** §9.1a 槽绑定 / rowAction 构造需在构造层实现时；若改须补 conformance/单测）
- **允许一次性新增** modal 宿主 / 确认对话框等渲染文件：`apps/web/src/renderer/modal*.tsx`、`apps/web/src/renderer/confirm*.tsx`（或等价命名）；新增文件须在 T-UI-10 说明用途，作为「一次性补齐」的一部分（A-006 R-001）。

**此后**「新增/调整代表页面」只允许改 `apps/api/internal/handler/fixtures/schema/*.json`（及配套 web 测试 fixtures）；T-UI-10 断言页级变更不触碰白名单外的 src 文件。

### 9.6 详情数据源（A-005 R-002 闭合）

`recordView` 与 edit-form 预填均渲染**选中行拷贝**（列表单一数据源）；不做独立 GET、不用 `recordSource`。列表 reload 后清空选中态。
