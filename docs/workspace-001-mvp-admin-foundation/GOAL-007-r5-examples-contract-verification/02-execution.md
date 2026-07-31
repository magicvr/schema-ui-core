---
id: GOAL-007-r5-examples-contract-verification
doc: execution
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-07-31
version: 0.5.0
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

## 待办

1. **批次 2c**：D-EXPR `form-with-reactions` + Renderer 接线（D-COMP，resolve F-002）。
2. 落地 `node`/`page` schema 校验与已纳入 fixtures 对照（`reactions`、`component-format`、`request-construction`、`response-mapping`、`query-serialization`、`static-data`、`actions`、`request-lifecycle`、`table-sort`、`search-table`、`runtime-defaults` 等），登记可执行验证入口；`I-PROTO-004` 决策先于阶段 3 结构校验实现。
3. 闭合父目标 `I-PROTO-003` 并完成 R5 自审/关门。

## 进度评估

**阶段 1/4 完成（契约发现与登记）；阶段 2 批次 2a（D-DATA/D-TABLE）与批次 2b（D-FORM/D-ACT）完成**。批次 2c（阶段 2 剩余）、阶段 3（结构/行为验证）、阶段 4（验收/关门）未开始（进度仅为展示，不放行阶段、不推导 `done`）。
