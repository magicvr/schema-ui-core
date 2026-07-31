---
id: GOAL-007-r5-examples-contract-verification
doc: audit
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-08-01
version: 0.5.0
---

# 审计 · GOAL-007

> 本文件是目标的唯一正式意见台账（P-003）。正式意见必须为可扫描的 `A-00N` 编号节。  
> A-001（independent）为立项基线审计；响应见文末。

## 信息就绪核对（开区基线）

| 核对项 | 状态 | 备注 |
|--------|------|------|
| 影响本 scope 的 I-00N | 已登记 | 见 [00-meta.md](00-meta.md)：`I-007-001`（required，R5 验收前） |
| 到期 required 是否已 verified / residual | 无到期 | `I-007-001` 最晚阶段 R5 验收前，本轮未到；父目标 `I-PROTO-003` 由 Root 维护 |
| 资料引用是否固定且用户确认 | 无 | `shared_materials_catalog: none`；范例候选源为协议清单与冻结覆盖表，非共享资料目录 |

## A-001 · R5 目标定义与信息就绪基线（2026-07-31）

- **source**：independent
- **auditor**：GitHub Copilot
- **类型**：goal-definition
- **scope**：GOAL-007 的 R5 立项范围、与 Root / VP 的对齐、冻结覆盖边界、信息项与当前规划阶段的可推进性。
- **verdict**：pass

### 范围与区间

本意见仅审计 R5 的目标定义与规划基线，不审计尚未实施的范例页、逐域 schema 校验或 R5 关门。当前工作区为 `workspace-001-mvp-admin-foundation`，Root 为 `GOAL-001-mvp-admin-foundation`，canonical 范围与 `VP-001-mvp-admin-foundation` 绑定一致；`shared_materials_catalog: none`，本 scope 未声明或依赖共享资料引用。

### 成果（有证据）

- GOAL-007 的父级、范围、四阶段路线图与 R5 成功标准已在 [00-meta.md](00-meta.md) 明确；目标将实现范围限制为 R2 v0.1.3 冻结的 11 个 include / include-partial 域，并明确排除 D-UPLOAD、未列白名单 type 和多选批量 action。
- 冻结覆盖表 [I-PROTO-001 v0.1.3](../GOAL-001-mvp-admin-foundation/attachments/I-PROTO-001-coverage-draft.md) 已逐域给出 disposition、范例候选和 fixture 映射；该基线与目标定义一致，且没有把冻结误述为 R5 已实施。
- `I-007-001` 与父目标 `I-PROTO-003` 均被明确为 R5 验收 / 关门前的 required 门禁，当前均未被伪称为已关闭。
- 可复用的 R3/R4 基线存在且可执行：`apps/web` 的 `npm test` 于 2026-07-31 通过，6 个测试文件共 94 项测试通过，覆盖 app manifest / navigation 与 permissions-inheritance 等既有路径。该结果只证明复用前提未退化，不证明 R5 的逐域范例或验证已完成。

### 对照成功标准

- 「每个纳入域的范例路径与验证入口」尚未达成，但当前路线图的首阶段正是为此建立 `I-007-001` 登记表；其最晚需要阶段为 R5 验收前，未出现到期 required 信息项。
- R5 的未开始实施状态在 [02-execution.md](02-execution.md) 中与当前代码证据相符：现有 Web 产物可指向 D-APP / D-PERM 既有路径，但未发现 R5 专属逐域范例页或验证登记，因而没有将这些既有产物计入 R5 完成事实。

### Findings

无 required finding。

### 必改项汇总

无。后续实施前应先由 `/govern` 执行路线图第 1 阶段，建立 `I-007-001` 的逐域范例路径与验证入口登记；在该信息项和父目标 `I-PROTO-003` 合法闭合前，不得放行 R5 验收或关门。

### 与既有意见的异同

本目标此前没有正式 `A-00N` 意见，因此不存在冲突或待响应的既有 finding。

### 结论 + 建议给编排器/用户的下一步

目标定义、范围边界、对齐链与当前阶段门禁可支持进入 R5 的契约发现登记；不能据此放行范例实现之外的范围、验收或关门。建议使用 `/govern workspace-001 GOAL-007` 响应本意见并推进 `I-007-001` 的登记工作。

### 声明

本意见不修改 status/progress；响应、finding 关闭和阶段推进由 `/govern` 处理。

## 响应（对 A-001）

| date | actor | summary |
|------|-------|---------|
| 2026-07-31 | `/govern` | 采纳 A-001 `pass`（无 required finding）。用户裁决「不需要自审，直接推进」；执行 R5 阶段 1：落盘 `I-007-001` 登记表（attachments/I-007-001-registry.md v0.1.0，逐纳入域范例路径 + 结构/行为验证入口），核验 D-APP/D-PERM 复用产物与可执行验证命令（web 94 项测试 / build、Go test / build），`I-007-001` → `verified`（登记层面）。阶段 2-4（范例实现、验证执行、验收关门）未开始；父目标 `I-PROTO-003` 仍 open，R5 验收/关门前须以阶段 3 可执行证据闭合。D-003 已留痕。 |

## A-002 · R5 阶段 2 批次 2a 自审：D-DATA / D-TABLE 范例与 Go 支撑（2026-07-31）

- **source**：self
- **auditor**：Claude Code · `/govern`
- **类型**：execution-facts
- **scope**：GOAL-007 批次 2a 实施——D-DATA（列表/详情数据 API + 前端数据客户端）与 D-TABLE（可排序表格 + 搜索表范例页）的实现、登记表同步与阶段 2 内可推进性。
- **verdict**：pass

### 范围与区间

本意见只审计批次 2a 的实施事实与阶段门禁，**不**审计批次 2b/2c、阶段 3 的 fixtures 正式执行、`I-PROTO-003` 闭合或 R5 关门。当前工作区 `workspace-001-mvp-admin-foundation`，Root `GOAL-001-mvp-admin-foundation`，canonical 范围一致；`shared_materials_catalog: none`，本 scope 不依赖共享资料引用（范例候选源为冻结覆盖表与协议清单，非资料目录）。

### 成果（有证据）

- **Go 列表/详情 API**（`apps/api/internal/handler/records.go`）：`GET /api/records` 支持 `q`/`sort`/`order`/`page`/`pageSize`（sort 白名单 name/status/owner/updatedAt，非法参数 400 fail-closed）；`GET /api/records/{id}`（404 `RECORD_NOT_FOUND`）。注册于 `handler.Register`（health.go）。`records_test.go` 6 项：默认列表、搜索、desc 排序、分页、非法参数、详情/404。`go test ./...` 全绿（21 项）、`go build ./...` 通过。
- **Web 数据客户端**（`apps/web/src/renderer/records.ts`）：`RecordItem`/`RecordList` 类型 + `buildRecordsQuery`（query-serialization）、`parseRecordList`（response-mapping，fail-closed 形状校验）、`fetchRecords`（request-construction）。`records.test.ts` 11 项通过。
- **表格组件**（`apps/web/src/components/data-table.tsx`）：泛型 `DataTable<T>`，可排序表头（aria-sort、asc/desc 切换）、加载/错误/空态、自定义单元格渲染。`data-table.test.tsx` 6 项通过。
- **范例页与接线**（`apps/web/src/app/examples/`）：`data-table-page.tsx` 与 `search-form-table-page.tsx`（搜索表单 + 表格）；`registry.tsx` 按 pageId 分发；`App.tsx` `PageSurface` 命中范例页时渲染真实页面，其余维持 manifest fallback。`app-manifest.json` 登记 `data-table` 与 `search-form-table` 页面 + sidebar「Examples」分组；icon registry 补 `table`/`search`；静态 manifest SHA 与 `app-manifest.test.ts` 页面清单断言已同步。
- **运行时证据**：起 `apps/api` + `apps/web` dev 后，Edge headless 实测 `/data-table`（Acme Console / Northwind Sales / Hooli Connect / Umbrella Ops，含 active/pending/alice）与 `/search-form-table`（8 records、rec-8 Globex Admin）均渲染真实 API 数据，无 Loading/Error；API 冒烟（默认 8 条 / `q=alice` 3 条 / `sort=updatedAt&order=desc` / detail / 404）与 Vite `/api` 代理均正确。
- **登记表同步**：`attachments/I-007-001-registry.md` 升 **v0.2.0**，D-DATA/D-TABLE 行「现状」→ 已实现，验证入口命令（web 114 项 / go 21 项 / API 冒烟 / 浏览器渲染）入账，阶段 3 待接入清单移出已落地映射路径。

### 对照成功标准（批次 2a 相关）

| 标准 | 状态 | 证据 |
|------|------|------|
| 未覆盖域（D-DATA / D-TABLE）具备可运行前后端范例路径 | 已达成 | 上述 Go API + Web 范例页 + 浏览器实测 |
| 结构验证可执行（`node`/`page` schema） | 阶段 3（未开始） | 依赖 `I-PROTO-004` 决策；不在批次 2a 范围 |
| 行为验证与 R2 基线一致（fixtures 正式执行） | 阶段 3（未开始） | 实现路径已落地，upstream cases 对照执行留待阶段 3 |
| 父目标 `I-PROTO-003` 闭合 | 未开始 | 验收/关门门禁；批次 2a 不触碰 |

### Findings

无 required finding。

批次 2a 范围内未发现阻断项。以下为如实记录，不构成必改：

- D-DATA / D-TABLE 的 upstream fixtures（`request-construction`、`response-mapping`、`query-serialization`、`static-data`、`table-sort`、`search-table`）**实现路径已落地**，但对照 `cases.json` 的正式执行属阶段 3，`I-007-001` 只确认登记（非逐域验证已执行）——与 D-002/D-003 既定路线一致。
- `I-PROTO-004`（vendor vs pin，non-blocking）仍 open，决策时点已定在阶段 3 结构校验实现前，不阻断批次 2a/2b/2c。

### 必改项汇总

无。

### 结论 + 建议下一步

批次 2a 实施与证据满足阶段 2 内推进条件：`npm test` 114 项 / `go test` 21 项 / build / 浏览器实测均通过，登记表与执行记录已同步，无开放 required finding，无到期 required 信息项影响本 scope。`I-PROTO-003` 仍 open（验收/关门门禁），本自审不放行验收或关门。建议按 D-004 进入批次 2b（D-FORM §5 白名单控件 + D-ACT 非批量动作）。

## A-003 · R5 阶段 2 批次 2b 自审：D-FORM 控件与 D-ACT 动作（2026-07-31）

- **source**：self
- **auditor**：Claude Code · `/govern`
- **类型**：execution-facts
- **scope**：GOAL-007 批次 2b 实施——D-FORM §5 白名单控件表面（含 2.6/2.7 版本/capability 门禁）、D-ACT 非批量行动作（复用 R4 `executeAction` 时序引擎）、Go PATCH/DELETE 编辑生命周期支撑与范例页的实现、登记表同步与阶段 2 内可推进性。
- **verdict**：pass

### 范围与区间

本意见只审计批次 2b 的实施事实与阶段门禁，**不**审计批次 2c（D-EXPR + Renderer 接线）、阶段 3 的 fixtures 正式执行、`I-PROTO-003` 闭合或 R5 关门。当前工作区 `workspace-001-mvp-admin-foundation`，Root `GOAL-001-mvp-admin-foundation`，canonical 范围一致；`shared_materials_catalog: none`，本 scope 不依赖共享资料引用。

### 成果（有证据）

- **D-FORM 控件表面**（`apps/web/src/renderer/form-controls.ts`）：`FormControlType` §5 白名单（base input/select；2.6 extended textarea/switch/checkbox/radio；2.7 advanced cascader/checkboxGroup/richText/password）+ `isWhitelistedFormControl`；`wireKindOf`（string/boolean/string-array）；`coerceFieldValue`（defaultValue 应用 + 强转）；`validateDefaultValue`（wire 类型失配 fail-closed `DEFAULT_VALUE_TYPE_MISMATCH`）；`checkFormCapabilities`/`checkFormCapabilitiesRaw`（extended ≥2.6 + `form.controls.extended`；advanced ≥2.7 + `form.controls.advanced`；select multiple ≥2.6；非白名单 `FORM_TYPE_NOT_WHITELISTED`）。`form-controls.tsx` `FormControls` 渲染全部 10 种控件。`form-controls.test.ts` **13 项**通过。
- **D-ACT 非批量动作**（`apps/web/src/renderer/row-action.ts`）：`runRowAction` 包装 R4 `executeAction`，`visible`/`confirm`/`confirmed`/`disabled`/`requiresSelection` 透传，返回 outcome/reason/permissionDenied/confirmed。`row-action.test.ts` **5 项**（执行、权限拒绝、隐藏 NOT_VISIBLE、确认取消+确认、disabled）通过。
- **Go 编辑生命周期支撑**（`apps/api/internal/handler/records.go`）：新增 `PATCH /api/records/{id}`（指针字段部分更新，`validatePatch` 空值 fail-closed）与 `DELETE /api/records/{id}`（204）；`recordHandler` 以 `sync.RWMutex` 保护共享数据集（MVP 无 DB）。`records_test.go` 新增 **5 项**（update、update invalid、update 404、delete、delete 404）。`go test ./...` 全绿（**18 顶层 / 26 含 Evaluate 子测试**）、`go build ./...` 通过。
- **Web 客户端**（`apps/web/src/renderer/records.ts`）：`updateRecord`（PATCH，部分字段 + fail-closed 形状校验）与 `deleteRecord`（DELETE，204）。`records.test.ts` 新增 **4 项**。
- **范例页与接线**（`apps/web/src/app/`）：`form-controls-page.tsx`（全控件 + 门禁展示 + Serialize）、`list-edit-lifecycle-page.tsx`（列表 + 行编辑/删除经 D-ACT 门禁，PATCH/DELETE 回写，复用 `DataTable`）；`registry.tsx` 注册两页；`App.tsx` icon 补 `form`/`pen`；`app-manifest.json` 登记 `form-controls`/`list-edit-lifecycle` + sidebar「Examples」组；`app-manifest.test.ts` 页面断言、`upstream-fixtures.test.ts` `STATIC_MANIFEST_SHA256` 同步；`app-examples.test.tsx` 新增 **2 项**表面测试（list-edit-lifecycle Edit/Delete 门禁、form-controls capability 门禁）。
- **测试 / 构建证据**：`npm test` 全绿（**11 测试文件 / 138 项**）；`npm run build`（tsc + vite）通过（修复一处 TS6133 未用变量后）。登记表 [I-007-001-registry.md](attachments/I-007-001-registry.md) 升 **v0.3.0**。

### 对照成功标准（批次 2b 相关）

| 标准 | 状态 | 证据 |
|------|------|------|
| 未覆盖域（D-FORM / D-ACT）具备可运行前后端范例路径 | 已达成 | 上述控件表面 + `runRowAction` + Go PATCH/DELETE + 范例页 |
| 结构验证可执行（`node`/`page` schema） | 阶段 3（未开始） | 依赖 `I-PROTO-004` 决策；不在批次 2b 范围 |
| 行为验证与 R2 基线一致（fixtures 正式执行） | 阶段 3（未开始） | 实现路径已落地，upstream cases 对照执行留待阶段 3 |
| 父目标 `I-PROTO-003` 闭合 | 未开始 | 验收/关门门禁；批次 2b 不触碰 |

### Findings

无 required finding。

批次 2b 范围内未发现阻断项。以下为如实记录，不构成必改：

- D-FORM / D-ACT 的 upstream fixtures（`component-format`、`actions`、`request-lifecycle` 非批量子集）**实现路径已落地**，但对照 `cases.json` 的正式执行属阶段 3，`I-007-001` 只确认登记（非逐域验证已执行）。
- 批次 2b 范例页（`form-controls` / `list-edit-lifecycle`）的渲染由 `app-examples.test.tsx` **jsdom 表面测试**覆盖；与批次 2a 的 Edge headless 浏览器实测不同，本轮未重跑真实浏览器渲染。该差异不阻断阶段 2 内推进，阶段 3/4 可补浏览器复核。
- `I-PROTO-004`（vendor vs pin，non-blocking）仍 open，决策时点定在阶段 3 结构校验实现前，不阻断批次 2b/2c。

### 必改项汇总

无。

### 结论 + 建议下一步

批次 2b 实施与证据满足阶段 2 内推进条件：`npm test` 138 项 / `go test` 18 顶层 / build 通过，登记表与执行记录已同步，无开放 required finding，无到期 required 信息项影响本 scope。`I-PROTO-003` 仍 open（验收/关门门禁），本自审不放行验收或关门。建议按 D-004 进入批次 2c（D-EXPR `form-with-reactions` + Renderer 接线）。

## A-004 · R5 阶段 2 批次 2c 自审：D-EXPR 反应引擎与 D-COMP Renderer 接线（2026-08-01）

- **source**：self
- **auditor**：Claude Code · `/govern`
- **类型**：execution-facts
- **scope**：GOAL-007 批次 2c 实施——D-EXPR 反应引擎（复用 `evaluateExpression`）、D-COMP 最小 Renderer 接线（resolve R4 F-002）、`form-with-reactions` 范例页的实现、登记表同步与阶段 2 内可推进性。
- **verdict**：pass

### 范围与区间

本意见只审计批次 2c 的实施事实与阶段门禁，**不**审计阶段 3 的 fixtures/schema 正式执行、`I-PROTO-003` 闭合或 R5 关门。当前工作区 `workspace-001-mvp-admin-foundation`，Root `GOAL-001-mvp-admin-foundation`，canonical 范围一致；`shared_materials_catalog: none`，本 scope 不依赖共享资料引用。

### 成果（有证据）

- **D-EXPR 反应引擎**（`apps/web/src/renderer/reactions.ts`）：`ReactionRule`/`ReactionApply`（显式 visible/disabled 布尔）；`parseReactionRule` fail-closed（非对象 `REACTION_APPLY_INVALID`、缺 id/when、非法表达式 `REACTION_EXPRESSION_INVALID`、apply 缺 fieldId/非布尔）；`evaluateReactions` 复用 `evaluateExpression`（frozen $context 子集，含 `contains`/`==`/`!=`），未知 apply 字段 `REACTION_APPLY_FIELD_UNKNOWN` fail-closed；`parseAndEvaluateReactions`。`reactions.test.ts` **12 项**通过。`app-manifest.ts` 新增 `isValidExpression` 导出（同一 frozen 语法，`app-manifest.test.ts` 不回归）。
- **D-COMP 最小 Renderer 接线**（`apps/web/src/renderer/render.ts`/`.tsx`）：`parseRenderNode` whitelist form/section/table，未知 type `RENDER_UNKNOWN_NODE_TYPE` fail-closed；form 缺 fields `RENDER_FORM_FIELD_INVALID`；`collectFieldIds`/`resolveFormReactions`/`gateAction`/`tableActionGate`。`render.test.ts` **11 项** + `render.test.tsx` **4 项**（默认全显、reaction 隐藏、reaction 禁用、未知 type alert）。`FormControls` 增加 `fieldDisabled` 按字段禁用（向后兼容，`form-controls.test.ts` 不回归）。
- **resolve F-002**：Renderer 集成层（`RenderPage`）消费 D-EXPR（reactions）与 D-FORM（FormControls）——即 GOAL-006 F-003/F-002 跟踪项「Renderer 接线」落地；本 scope 不主张「页面运行时已全面应用 D-PERM」（权限引擎经 R4 既有 `row-action.ts` 路径消费，非本批次新增）。
- **范例页与接线**（`apps/web/src/app/`）：`form-with-reactions-page.tsx`（Admin/Viewer + audit feature 切换 → `$context` 快照 → 字段显隐/禁用）；`registry.tsx`/`App.tsx` icon（`reaction`=Zap）/`app-manifest.json`（页面 + sidebar）/`app-manifest.test.ts` 断言/`upstream-fixtures.test.ts` `STATIC_MANIFEST_SHA256`（`0475f7bb…5e120`）/`app-examples.test.tsx`（+1 项）同步。
- **测试 / 构建 / 运行时证据**：`npm test` 全绿（**14 文件 / 166 项**，+28）；`npm run build`（tsc + vite）通过；Edge headless 实测 `/form-with-reactions` 渲染 Admin/Viewer/audit 切换 + 全部 4 字段 + 无 `role="alert"`。Go 侧无改动（`go test ./...`/`go build ./...` 通过）。登记表升 **v0.4.0**。

### 对照成功标准（批次 2c 相关）

| 标准 | 状态 | 证据 |
|------|------|------|
| D-EXPR 具备可运行范例路径 | 已达成 | `reactions.ts` + `form-with-reactions` 范例页 + Edge 实测 |
| D-COMP Renderer 接线消费引擎 | 已达成 | `RenderPage` 消费 reactions/FormControls；whitelist fail-closed |
| 结构验证可执行（`node`/`page`/`reaction` schema） | 阶段 3（未开始） | 依赖 `I-PROTO-004` 决策；不在批次 2c 范围 |
| 行为验证与 R2 基线一致（`reactions`/`component-format` fixtures 正式执行） | 阶段 3（未开始） | 实现路径已落地，upstream cases 对照执行留待阶段 3 |
| 父目标 `I-PROTO-003` 闭合 | 未开始 | 验收/关门门禁；批次 2c 不触碰 |

### Findings

无 required finding。

批次 2c 范围内未发现阻断项。以下为如实记录，不构成必改：

- D-EXPR 的 upstream `reactions` fixture 与 D-COMP 的 `component-format`（5 case）**实现路径已落地**，但对照 `cases.json` 的正式执行属阶段 3，`I-007-001` 只确认登记（非逐域验证已执行）。
- `reaction.schema.json` / `node`/`page` schema 结构校验依赖 `I-PROTO-004`（vendor vs pin）决策，批次 2c 未做 schema 级验证。
- 批次 2c 范例页（`form-with-reactions`）有 Edge headless 实测；reactions/render 逻辑另有单元测试覆盖。jsdom 表面测试覆盖 App 级接线。
- `I-PROTO-004`（vendor vs pin，non-blocking）仍 open，决策时点定在阶段 3 结构校验实现前，不阻断批次 2c。

### 必改项汇总

无。

### 结论 + 建议下一步

批次 2c 实施与证据满足阶段 2 内推进条件：`npm test` 166 项 / `npm run build` / `go test` / Edge 实测均通过，登记表与执行记录已同步，无开放 required finding，无到期 required 信息项影响本 scope。`I-PROTO-003` 仍 open（验收/关门门禁），本自审不放行验收或关门。阶段 2 全部落地，建议进入阶段 3（结构/行为验证）：先按 `I-PROTO-004` 决策 vendor vs pin，再接入 `node`/`page`/`reaction` schema 校验与已纳入 fixtures 对照。

