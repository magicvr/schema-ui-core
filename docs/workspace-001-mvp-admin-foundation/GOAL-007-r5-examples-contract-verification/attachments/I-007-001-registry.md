---
title: I-007-001 · R5 纳入域范例路径与验证入口登记表
status: active
doc_type: info-registry
created: 2026-07-31
updated: 2026-08-01
parent: GOAL-007-r5-examples-contract-verification
version: 0.4.0
related_info: I-007-001
coverage_freeze: GOAL-001-mvp-admin-foundation/attachments/I-PROTO-001-coverage-draft.md v0.1.3
inventory: docs/vision/protocol-inventory-v2.7.0.md
inventory_pin: ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b
---

# I-007-001 · R5 纳入域范例路径与验证入口登记表

> **性质**：R5 阶段 1（契约发现与登记）的登记产物，回答 `I-007-001`「每条纳入能力的范例页路径与自动化/手工验证入口」。
> **本表是信息登记**：标注「已有」的行指向现存产物与可执行验证入口；标注「规划」的行登记拟建路径与验证入口，**执行结果**属 R5 阶段 2/3，本表不把它们写成已验证事实。
> **范围基线**：[I-PROTO-001 v0.1.3 §3] 范例候选 + [protocol-inventory §2.5] 场景 + [§3] 映射。排除项（D-UPLOAD、多选批量、完整 registry、scenarios 自动化门禁）不登记。

## 1. 逐域登记

| domain | disposition | 范例页 / 场景路径 | 结构验证入口 | 行为验证入口 | 现状 | 依赖 / 备注 |
|--------|-------------|--------------------|--------------|--------------|------|-------------|
| D-NODE | include | 任意合法 page（全站基座，无专属业务页） | `node`/`page` schema 校验（须先 vendor 或 pin，随 I-PROTO-004） | —（基座） | 规划 | 全站 page 结构先决；阶段 3 落地校验命令 |
| D-EXPR | include | `form-with-reactions`（改写，去业务化） | `reaction.schema.json`（随 I-PROTO-004） | fixtures `reactions`（阶段 3 接入 suite） | **已实现（批次 2c）**：`apps/web/src/renderer/reactions.ts`（`ReactionRule`/`parseReactionRule` fail-closed + `evaluateReactions` 复用 `evaluateExpression`，frozen $context 表达式子集；`parseAndEvaluateReactions`）+ `reactions.test.ts` 12 项；范例页 `form-with-reactions`（`app/examples/form-with-reactions-page.tsx`，Admin/Viewer 角色 + audit feature 切换演示字段显隐/禁用，经最小 Renderer `RenderPage` 渲染） | 复用 `evaluateExpression`（导航 visibleWhen 同源） |
| D-COMP | include-partial | `data-table` / `form-controls-*`（§5 白名单 type） | `node`/`page` schema + registry membership（§5） | fixtures `component-format`（5 case：currency/percent/datetime） | **已实现（批次 2c）**：`apps/web/src/renderer/render.ts`（`parseRenderNode`/`collectFieldIds`/`resolveFormReactions`/`gateAction`/`tableActionGate`，whitelist form/section/table fail-closed）+ `render.tsx` `RenderPage`（form 经 `FormControls` + reactions 应用；未知 type 出 alert）+ `render.test.ts` 11 项 + `render.test.tsx` 4 项 | 白名单 §5 已冻结；`reaction.schema.json` 结构校验随 I-PROTO-004 |
| D-DATA | include | `search-form-table` / `data-table`（Go 列表/详情 + 前端） | `page` schema | fixtures `request-construction`、`response-mapping`、`query-serialization`、`static-data` | **已实现（批次 2a + 2b 扩展）**：Go `GET /api/records`（q/sort/order/page/pageSize）+ `GET /api/records/{id}` + `PATCH`/`DELETE /api/records/{id}`（`apps/api/internal/handler/records.go`）；Web `fetchRecords`/`buildRecordsQuery`/`parseRecordList`/`updateRecord`/`deleteRecord`（`apps/web/src/renderer/records.ts`）+ `useRecords` hook（`use-records.ts`）；范例页 `search-form-table` / `data-table` / `list-edit-lifecycle`（`apps/web/src/app/examples/`）。Go 测试（records 12 项）、Web 测试（records 15 项 + app-examples 5 项）入账；浏览器实测（Edge headless）2a 两页渲染真实数据，2b 页面由 jsdom 表面测试覆盖 | 
| D-ACT | include-partial | `row-backend-actions` / `admin-list-edit-lifecycle`（非批量） | `action.schema.json`（随 I-PROTO-004） | fixtures `actions`、`request-lifecycle`（非批量子集，阶段 3 接入） | **已实现（批次 2b）**：`apps/web/src/renderer/row-action.ts` `runRowAction` 包装 R4 `executeAction` 时序引擎（visible/confirm/confirmed/disabled/requiresSelection 透传，fail-closed）；`row-action.test.ts` 5 项；范例页 `list-edit-lifecycle`（Edit/Delete 经 D-ACT 门禁 + PATCH/DELETE 回写）；Go 新增 `PATCH`/`DELETE /api/records/{id}`（`records.go`，mutex 保护，`validatePatch` fail-closed，`records_test.go` +5 项） | Q1=否排除批量 |
| D-PERM | include | `permission-inheritance`（**已有**：`apps/web/src/renderer/permissions.ts` + `permissions-inheritance.test.ts`） | `validatePermissions`（L2 fail-closed 校验，permissions.ts） | fixtures `permissions-inheritance`（GOAL-006 attachments `dperm/cases.json` 17 例，SHA-256 pin 于测试） | **已有（R4）** | 仅登记复核，不重复实现 |
| D-APP | include | Admin 外壳 + 导航（**已有**：`apps/web/src/app/App.tsx`、`app/navigation.ts`、`public/.well-known/schema-ui/app-manifest.json`） | `docs/schemas/app-manifest.schema.json`（已 vendor，SHA pin） | fixtures `app-manifest`、`app-navigation`（`apps/web/src/protocol/upstream/*.cases.json` + `upstream-fixtures.test.ts`） | **已有（R3）** | 仅登记复核，不重复实现 |
| D-TABLE | include-partial | `data-table` / `search-table`（排序声明 + 基础列表交互） | `page` schema | fixtures `table-sort`、`search-table` | **已实现（批次 2a）**：`apps/web/src/components/data-table.tsx`（可排序列表，aria-sort、asc/desc 切换，渲染/加载/空/错误态）+ `data-table.test.tsx` 6 项；sort 声明经 `useRecords`/`fetchRecords` 序列化为 `sort`/`order` 查询参数（对照 table-sort / query-serialization 语义） | 排除多选批量（Q1=否）；依赖 D-DATA Go API + 前端表格组件 |
| D-FORM | include-partial | `form-controls-advanced`（§5 全部 2.6/2.7 控件） | `node`/`page` schema + registry + 版本/capability 下限 | `component-format`（格式子集）；scenarios（手工路径，Q5=否） | **已实现（批次 2b）**：`apps/web/src/renderer/form-controls.ts`（`FormControlType` §5 白名单、`wireKindOf`、`coerceFieldValue`、`validateDefaultValue` fail-closed、`checkFormCapabilities` 版本/capability 门禁）+ `form-controls.tsx` `FormControls` 组件（10 种控件）+ `form-controls.test.ts` 13 项；范例页 `form-controls` | `defaultValue` 属性规则随 2.7 |
| D-VER | include | 全站（`supportedCapabilities` / 版本协商）——已有子集：`validateAppManifest` + `upstream-fixtures.test.ts` `negotiateFixture` | — | fixtures `version-negotiation`（已在 app-manifest fixtures 覆盖）、`runtime-defaults`（待接入） | 部分已有 | 版本协商已随 D-APP fixtures 落地；`runtime-defaults` 阶段 3 接入 |
| D-VAL | include | —（构建时/加载时结构校验） | `docs/06-validation.md` + 6 schemas（已 vendor `app-manifest.schema.json`） | — | 部分已有 | 其余 5 schemas 随 I-PROTO-004 决策 vendor/pin |

## 2. 验证入口命令登记（已存在的可执行入口）

| 入口 | 命令 | 覆盖域 | 证据 / 说明 |
|------|------|--------|-------------|
| Web 单元 + upstream fixtures | `cd apps/web && npm test` | D-APP、D-PERM、D-EXPR、D-DATA、D-TABLE、D-FORM、D-ACT、D-COMP | 14 测试文件 / **166 项**（含 `upstream-fixtures.test.ts`、`permissions-inheritance.test.ts`、`app-manifest.test.ts`、`navigation.test.ts`、`App.integration.test.tsx`、`account/context.test.ts`、`data-table.test.tsx`、`records.test.ts`、`app-examples.test.tsx`、`form-controls.test.ts`、`row-action.test.ts`、`reactions.test.ts`、`render.test.ts`、`render.test.tsx`）；上游 fixture SHA 与 provenance 已 pin |
| Web 构建 | `cd apps/web && npm run build` | D-APP、D-DATA、D-TABLE、D-FORM、D-ACT | tsc + vite build |
| Go 测试 | `cd apps/api && go test ./...` | D-PERM、D-VER(子集)、D-DATA、D-ACT | `internal/account/permission_test.go`、`internal/handler/account_test.go`、`internal/handler/records_test.go`（records 12 项：list/detail/PATCH/DELETE）、`health_test.go`（**18 顶层 / 26 含 Evaluate 子测试**） |
| Go 构建 | `cd apps/api && go build ./...` | D-PERM、D-VER(子集)、D-DATA、D-ACT | server 可运行 |
| Go 数据 API 冒烟 | `curl http://127.0.0.1:8080/api/records`（q/sort/order/page/pageSize）、`/api/records/{id}`、PATCH/DELETE `/api/records/{id}`、404 | D-DATA、D-TABLE、D-ACT | 批次 2a 运行时冒烟（list/detail/排序/分页/搜索/404）；批次 2b PATCH/DELETE 由 `records_test.go` 覆盖（update/delete/404/invalid）；详见 02-execution |
| 浏览器渲染 | 起 `apps/api` + `apps/web` 后访问 `/data-table`、`/search-form-table`、`/list-edit-lifecycle`、`/form-controls`、`/form-with-reactions` | D-DATA、D-TABLE、D-APP、D-FORM、D-ACT、D-EXPR、D-COMP | Edge headless 实测 2a 两页渲染真实数据（Acme Console…8 records），无 Loading/Error；2b 页面渲染由 `app-examples.test.tsx` jsdom 表面测试覆盖；批次 2c `/form-with-reactions` Edge headless 实测渲染 Admin/Viewer 切换 + 全部字段 + 无 alert |

**阶段 3 待接入的验证入口**（登记为计划，执行结果未发生）：`reactions`、`component-format`（5 case）、`actions`、`request-lifecycle`（非批量）、`runtime-defaults`；以及 `node`/`page`/`action`/`reaction` schema 校验命令（依赖 I-PROTO-004 vendor/pin 决策）。`request-construction` / `response-mapping` / `query-serialization` / `static-data`（D-DATA）与 `table-sort` / `search-table`（D-TABLE）的**实现路径已落地**（批次 2a：`renderer/records.ts`、`components/data-table.tsx`、Go `handler/records.go`），对照上游 cases.json 的正式接入执行仍在阶段 3。

## 3. 与 I-PROTO-001 v0.1.3 §3 候选的对齐

- P0 候选：`permission-inheritance`（D-PERM，**已有**）、Admin 外壳 + 导航（D-APP，**已有**）。
- P1 候选：`search-form-table` / `data-table`（D-DATA、D-TABLE）→ **已实现（批次 2a）**；`admin-list-edit-lifecycle` / `list-detail`（D-ACT、D-DATA、D-FORM）→ **已实现（批次 2b）**。
- P2 候选：`form-with-reactions`（D-EXPR、D-FORM）→ **已实现（批次 2c）**。
- 不做（MVP）：`form-with-upload`、`form-controls-advanced/extended` 全量、`order-*` 业务样例原文 → 与 §5 白名单、非目标一致，不登记。

## 4. 复用产物（不重复实现）

| 域 | 复用路径 | 复核点 |
|----|----------|--------|
| D-APP | `apps/web/src/app/App.tsx`、`app/navigation.ts`、`protocol/app-manifest.ts`、`public/.well-known/schema-ui/app-manifest.json` | `npm test` / `npm run build` 通过即为复用前提未退化 |
| D-PERM | `apps/web/src/renderer/permissions.ts`、`account/context.ts`、`apps/api/internal/account/{session,permission}.go`、`handler/account.go` | `npm test`（permissions-inheritance 17 例）+ `go test ./...` 通过 |
| D-DATA / D-TABLE（批次 2a 新增） | `apps/web/src/renderer/{records.ts,use-records.ts}`、`components/data-table.tsx`、`app/examples/{data-table-page,search-form-table-page}.tsx`、`apps/api/internal/handler/records.go` | `npm test`（records 15 项 + data-table 6 项 + app-examples 5 项）+ `go test ./...`（records 12 项）+ 浏览器渲染实测（2a） |
| D-FORM / D-ACT（批次 2b 新增） | `apps/web/src/renderer/{form-controls.ts,form-controls.tsx,row-action.ts}`、`app/examples/{form-controls-page,list-edit-lifecycle-page}.tsx`、`apps/api/internal/handler/records.go`（PATCH/DELETE） | `npm test`（form-controls 13 项 + row-action 5 项 + records PATCH/DELETE 4 项 + app-examples 2 项）+ `go test ./...`（records +5 项）+ jsdom 表面测试覆盖 |
| D-EXPR / D-COMP（批次 2c 新增） | `apps/web/src/renderer/{reactions.ts,render.ts,render.tsx}`、`app/examples/form-with-reactions-page.tsx` | `npm test`（reactions 12 项 + render 11 项 + render.tsx 4 项 + app-examples 1 项）+ `npm run build` + Edge headless 实测 `/form-with-reactions`（无 alert） |

## 5. 变更规则

- 本表是 `I-007-001` 的信息登记，**不是** R5 关门证据；`I-PROTO-003` 的闭合仍须阶段 3 的可执行验证结果与阶段 4 验收。
- 覆盖范围变更须改 [I-PROTO-001 v0.1.3]（新决策 + 新版本），并同步重估本表。
- 阶段 2/3 实现落盘后，更新本表对应行的「现状」与验证命令为已发生事实。

> 变更：v0.4.0（2026-08-01）——批次 2c 实现落盘：D-EXPR 行「现状」→ 已实现（`reactions.ts` + `form-with-reactions` 范例页），D-COMP 行「现状」→ 已实现（`render.ts`/`render.tsx` 最小 Renderer，resolve F-002），P2 `form-with-reactions` 标记已实现；验证命令计数更新（web 166 项 / 14 文件），浏览器渲染新增 `/form-with-reactions` Edge headless 实测。

> 变更：v0.3.0（2026-07-31）——批次 2b 实现落盘：D-FORM/D-ACT 行「现状」→ 已实现（控件表面 + `runRowAction` + Go PATCH/DELETE + 范例页 `form-controls`/`list-edit-lifecycle`），D-DATA 行扩展 PATCH/DELETE，验证命令计数更新（web 138 项 / go 18 顶层）；P1 `admin-list-edit-lifecycle` 标记已实现。

> 变更：v0.2.0（2026-07-31）——批次 2a 实现落盘：D-DATA/D-TABLE 行「现状」→ 已实现，新增验证入口（web 114 项 / go 21 项 / API 冒烟 / 浏览器渲染）；阶段 3 待接入清单移出已落地映射路径。
