---
title: I-007-003 · Schema CRUD 读写交互契约
status: active
doc_type: info-collection
created: 2026-08-02
updated: 2026-08-02
parent: GOAL-007-r4-schema-crud
version: 0.1.0
related_info: I-007-003
related_decision: D-005
---

# I-007-003 · Schema CRUD 页面/Node/action 绑定与交互/权限契约

> **结论**：本附件与 D-005 关闭 `I-007-003`，并完成 S4/S5 的交互契约冻结。以下页面绑定、字段映射、交互状态与权限矩阵是 S4/S5 实施输入，**不是**已交付的 Schema 写交互代码或页面事实。
> **扫描日期**：2026-08-02。工作区 `shared_materials_catalog: none`；事实来自本仓库 Web Renderer、permission/action 引擎、records client、schema fixtures 与既有测试。

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

**页面**：`list-edit-lifecycle`（fixture `list-edit-lifecycle.json`）为唯一代表性 CRUD 生命周期页；与 R3 种子菜单 `list-edit-lifecycle`（feature `menu_list_edit_lifecycle`，admin 可见）对齐。

| 节点 / 属性 | 绑定 | API |
|------------|------|-----|
| `table`（`props.dataSource`） | 记录列表 | `GET /api/records`（q/sort/order/page/pageSize） |
| `table.props.actions` rowAction `edit`（`permissionIntent: edit`） | 打开编辑 form 并预填该行字段 | 预填来自行数据；提交见下 |
| `table.props.actions` rowAction `delete`（`permissionIntent: delete`） | 确认框（`confirm`）→ 删除 → 行移除 | `DELETE /api/records/{id}` |
| `table.props.toolbar` toolbarTrigger `create`（`permissionIntent: edit`） | 打开新建 form（空表单） | 提交见下 |
| `form` `submitAction`（formSubmit） | create 模式 → 新建；edit 模式 → 更新；成功/失败反馈 + 列表刷新 | `POST /api/records` / `PATCH /api/records/{id}` |
| 搜索 `form`（`search-form-table` 模式） | 提交 → table `q` 绑定 | `GET /api/records?q=…` |
| `recordView` | 详情展示（只读） | `GET /api/records/{id}` |

**一次性渲染补齐（S4 实施范围，非每页）**：渲染 table `props.actions`/`props.toolbar`；绑定 form `submitAction` → 对应 POST/PATCH；成功/错误/确认反馈。此后「新增或调整代表页面」仅改 fixture，**不修改 Renderer 主路径代码**（S4 成功标准）。

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

**表达式门控（冻结 `$context` 语法）**：
- 写 affordance（create/edit/delete）：`visibleWhen: "$context.user.permissions contains \"records.write\""`，或经 `permissionIntent: edit`/`delete` + `permissions`（复用 `permissions.inheritance` 引擎）。
- 读 affordance（页面/列表/详情）：`"$context.user.permissions contains \"records.read\""`。
- 菜单：`$context.features.menu_list_edit_lifecycle`（已实现，不变）。
- 使用 permission 字段的页面 `meta.requiredCapabilities` 必须含 `permissions.inheritance`；协议版本 ≥ 2.3（`permissions.ts` 冻结门禁）。

## 6. 测试矩阵（S4/S5 验收最低断言，T-UI-01～10）

| ID | 断言 |
|----|------|
| T-UI-01 | 代表页列表从 `/api/records` 加载：loading → ready；空表 → "No records match." |
| T-UI-02 | 列排序头切换 sort/order；分页计数正确显示 |
| T-UI-03 | 搜索 form 提交 → 列表按 `q` 过滤 |
| T-UI-04 | 新建 form（name/status/owner）→ POST → 201；成功后列表含新行；空白字段 → `INVALID_CREATE_FIELD` 表单错误 |
| T-UI-05 | 编辑 form 预填行数据 → PATCH → 200；该行 name/status/owner 更新且 `updatedAt` 刷新 |
| T-UI-06 | 删除行 → 确认框 → 确认 → DELETE 204 → 行移除；取消 → 无请求 |
| T-UI-07 | API 错误 envelope 正确呈现（`INVALID_CREATE_FIELD` / `RECORD_NOT_FOUND` / `FORBIDDEN`） |
| T-UI-08 | admin 可见全部 affordance；viewer/editor 只读（隐藏/禁用写 affordance）；匿名重定向登录 |
| T-UI-09 | 后端权威不被前端隐藏替代：即便隐藏，直接 POST/PATCH/DELETE 由 API 403/401 拦截（API 侧 T-API-08/09 已承担；UI 断言呈现一致） |
| T-UI-10 | 新增/调整代表页面不修改 Renderer 主路径代码（页面变更仅触 fixture；用文件 diff 断言） |

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
- I-007-001（v0.2.0）、I-007-002（v0.2.0）、D-005
