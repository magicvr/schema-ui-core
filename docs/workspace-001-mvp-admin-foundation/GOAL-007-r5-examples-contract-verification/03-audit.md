---
id: GOAL-007-r5-examples-contract-verification
doc: audit
status: done
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-08-01
version: 0.14.0
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

## A-005 · R5 阶段 2 独立复核：冻结契约一致性与阶段 3 就绪性（2026-08-01）

- **source**：independent
- **auditor**：OpenAI Codex
- **类型**：execution-facts
- **scope**：复核 GOAL-007 阶段 2 批次 2a/2b/2c 的实现事实、v0.1.3 冻结边界与进入阶段 3 的就绪性；不审计尚未执行的 schema/fixture 正式对照、`I-PROTO-003` 闭合、R5 验收或关门。
- **verdict**：conditional

### 范围与区间

当前工作区为 `workspace-001-mvp-admin-foundation`，Root 为 `GOAL-001-mvp-admin-foundation`，canonical 范围为 `docs/workspace-001-mvp-admin-foundation/`。工作区、Root、VP-001 与 active Charter 的 `plan_refs` / `primary_plan` / `vision_ref` 链可解析；`shared_materials_catalog: none`，本意见未把共享资料作为事实或关闭证据。

本意见以 Root 已冻结的 [I-PROTO-001 v0.1.3](../GOAL-001-mvp-admin-foundation/attachments/I-PROTO-001-coverage-draft.md) 为范围权威。阶段 3 尚未执行是既定路线，不单独构成 finding；以下 findings 只针对阶段 2 已实现代码与“阶段 2 全部落地”主张之间的可复核差异。

### 成果（有证据）

- 当前工作树复跑 `cd apps/web && npm test` 通过：14 个测试文件、166 项测试；`npm run build` 通过。
- 当前工作树复跑 `cd apps/api && go test ./...`、`go build ./...` 通过。
- 批次 2a 的 Go records API、Web records client、DataTable 与两条范例路由存在，列表/详情/搜索/排序/PATCH/DELETE 的单元证据可复核。
- 批次 2b/2c 的表单控件、row action、reactions、最小 RenderPage 与范例路由存在；实现并非“未发生”。但当前测试把部分 fail-open / 缺门禁行为写成通过预期，测试全绿不能关闭下列契约差异。

### 对照成功标准

| 标准 | 结论 | 证据 |
|------|------|------|
| 未覆盖域具备可运行实现/范例路径 | 部分满足 | D-DATA、D-TABLE、D-FORM、D-ACT、D-EXPR 与最小 D-COMP 路径存在；但 D-FORM / D-COMP 的冻结门禁与 type 范围不完整，见 F-001～F-004 |
| 结构/行为正式验证 | 未到阶段 | `node`/`page`/`action`/`reaction` schema 与纳入 fixtures 正式执行属于阶段 3；本意见不把“尚未执行”改写为失败 |
| `I-PROTO-003` / R5 验收关门 | 未到阶段且保持阻断 | Root `I-PROTO-003` 仍 open / required；本意见不关闭该信息项，也不放行 R5 验收或关门 |

### Findings

#### F-001 · 无效 action 表达式在 visible 门禁上 fail-open

- **level**：required
- **severity**：medium
- **status**：open
- **影响门禁**：阶段 2 D-ACT / D-COMP 完成认定；阶段 3 行为验证放行。
- **证据**：`apps/web/src/renderer/render.ts` 的 `gateAction()` 对无效字符串返回调用方 `defaultValue`，而 `tableActionGate()` 对 `visibleWhen` 使用默认 `true`；因此显式但非法的 `visibleWhen` 会得到 `visible: true`。`apps/web/src/renderer/render.test.ts` 还以“fail-closed”为测试名断言非法 `$deps.x == true` 返回 `true`。同仓 `evaluateExpression()` 对不匹配语法返回 `false`，权限值未知时也按 deny 处理。
- **必改**：区分“属性缺省”与“显式非法表达式”；显式非法 `visibleWhen` 必须拒绝/隐藏并产生可核对错误，补 `tableActionGate` 的非法 visible/disabled 回归测试，不能用默认可见关闭该 finding。

#### F-002 · RenderPage 未执行其声明的 D-FORM 版本 / capability / type 门禁

- **level**：required
- **severity**：medium
- **status**：open
- **影响门禁**：阶段 2 D-FORM / D-COMP 完成认定；阶段 3 结构验证放行。
- **证据**：`apps/web/src/renderer/render.tsx` 的 `FormView` 只解析 reactions，并把原始 `fields` 强转为 `FormControlField[]` 交给 `FormControls`；没有调用 `checkFormCapabilities()`，也没有用 `document.meta` 执行版本/capability 门禁。`render.test.tsx` 的测试页只声明 `app.manifest`，却包含需要 `form.controls.extended` 的 `textarea` / `switch`，并断言它们正常渲染。未知 field type 也会在 `FormControls` 的 switch 尾部静默返回 `null`，没有 fail-closed 错误。
- **必改**：RenderPage 在渲染前须解析字段并执行 D-FORM whitelist、版本与 capability 门禁；非法 type、低版本或缺 capability 应产生确定错误并拒绝受影响字段/页面，补对应 Renderer 级测试。

#### F-003 · `defaultValue` 只校验 wire 类型，未执行冻结的 2.7 + advanced 门禁

- **level**：required
- **severity**：medium
- **status**：open
- **影响门禁**：阶段 2 D-FORM 完成认定；阶段 3 表单结构/行为验证放行。
- **证据**：[I-PROTO-001 v0.1.3 §5](../GOAL-001-mvp-admin-foundation/attachments/I-PROTO-001-coverage-draft.md) 明确任一纳入字段的 `props.defaultValue` 要求 `protocolVersion >= 2.7` 且包含 `form.controls.advanced`。`form-controls.ts` 的字段循环只对 select-multiple 做 2.6 门禁，然后调用 `validateDefaultValue()` 校验 wire 类型；`form-controls.test.ts` 还断言缺少 `form.controls.advanced` 的 base meta 可接受 input `defaultValue: "ok"`。
- **必改**：为任何出现 `defaultValue` 的字段增加 2.7 与 `form.controls.advanced` 双门禁，并补 2.6、2.7 缺 capability、2.7 完整 capability 三组回归测试。

#### F-004 · D-COMP Renderer whitelist 静默窄于 v0.1.3 冻结初表

- **level**：required
- **severity**：medium
- **status**：open
- **影响门禁**：阶段 2 D-COMP 完成认定；`I-PROTO-003` 的逐纳入能力证据；阶段 3 schema / registry 对照。
- **证据**：[I-PROTO-001 v0.1.3 §5](../GOAL-001-mvp-admin-foundation/attachments/I-PROTO-001-coverage-draft.md) 的初始 type 白名单包含布局 `grid` / `section` / `tabs`，数据与操作 `text` / `table` / `recordView` / `actionButton`，以及表单 type；冻结后缩放范围要求新决策、新版本与重估。当前 `render.ts` 只把 `form` / `section` / `table` 作为 Renderer whitelist，`render.test.ts` 明确断言 `grid` 被拒绝；`I-007-001-registry.md` 的 D-COMP 行仍将该三 type Renderer 记为“已实现”，没有逐项解释其余冻结 type 的实现证据或合法范围变更。
- **必改**：由 `/govern` 逐项映射冻结 whitelist 的现有实现与验证路径并补齐缺口；若确需缩小 v0.1.3 范围，必须按冻结规则取得用户书面裁决、追加 Root 决策与新版覆盖表，再同步登记表，不能以“最小 Renderer”静默改写冻结范围。

#### F-005 · Root 的 R5 路线投影落后于目标与 goal-tree

- **level**：recommended
- **severity**：low
- **status**：open
- **影响**：阶段事实可追踪性，不单独阻断阶段 3 修正工作。
- **证据**：Root `00-meta.md` 的 R5 行仍写“批次 2c 与阶段 3/4 未开始”，而 GOAL-007 `00-meta` / `02-execution` 与工作区 `goal-tree.md` 已写批次 2c 完成、阶段 2 全部落地。
- **建议**：在 `/govern` 响应本意见时同步 Root 路线投影；不要借此改 Root `progress`（仍为 4/6）。

### 必改项汇总

1. `F-001`：非法 action 表达式改为 fail-closed，并补行为测试。
2. `F-002`：RenderPage 接入 D-FORM type / version / capability 门禁与错误呈现。
3. `F-003`：补 `defaultValue` 的 2.7 + `form.controls.advanced` 双门禁。
4. `F-004`：补齐冻结 D-COMP type 映射/实现，或按治理规则书面修订冻结范围。

开放 required = 4。可进入“修正 + 阶段 3 验证实现”工作，但不得把当前阶段 2 证据解释为无条件完成、不得关闭 `I-PROTO-003`、不得放行 R5 验收或关门。

### 与既有意见的异同

- 与 A-002～A-004 一致：阶段 2 的实现、166 项 Web 测试、Web build 与 Go test/build 均有可复核事实；阶段 3 的 schema/fixture 正式执行尚未发生，且 `I-PROTO-003` 仍 open。
- 与 A-003/A-004 的“无 required finding / pass”不同：本意见依据同一冻结范围发现 F-001～F-004。对同一批次完成主张，一方认为无必改、一方认为 required，构成 P-004 需要由 `/govern` 展示并请用户裁决的同 scope 分歧；在 fixed / accepted-residual / user-overruled 留痕前，不取乐观侧放行。

### 结论 + 建议给编排器/用户的下一步

阶段 2 不是“未实现”，但当前证据不足以支持无条件完成认定，故 verdict 为 **conditional**。建议使用 `/govern workspace-001-mvp-admin-foundation GOAL-007-r5-examples-contract-verification` 响应 A-005：优先按 `fixed` 路径修正 F-001～F-004，将修正测试与阶段 3 schema/fixture 执行统一纳入可重复验证；同时处理 P-004 同 scope 分歧和 F-005 投影同步。

### 声明

本意见不修改 status/progress、路线图检查点、`I-PROTO-003` / `I-PROTO-004` 状态或 goal-tree；finding 响应与阶段推进由 `/govern` 处理。

## 响应（对 A-005）

| date | actor | summary |
|------|-------|---------|
| 2026-08-01 | `/govern` | 响应 A-005（independent, conditional）：用户裁决「不需要自审，直接推进」（P-004 §3.1），并按 **`fixed`** 路径闭合 F-001～F-004；F-005（recommended）同步。实施与证据：F-001 → `fixed`（`render.ts` 以 `resolveActionGate` 区分「属性缺省」与「显式非法表达式」，显式非法 `visibleWhen`/`disabledWhen` fail-closed 并出 `ACTION_GATE_EXPRESSION_INVALID` 可核对错误；`tableActionGate` 返回 `{visible, disabled, errors}`，补非法 visible/disabled 回归测试）；F-002 → `fixed`（`render.tsx` `FormView` 经 `gateRenderFormFields` 在渲染前执行 D-FORM §5 whitelist、版本与 capability 门禁，非法 type / 低版本 / 缺 capability 拒绝受影响字段并出 `role="alert"` 错误；补 Renderer 级测试）；F-003 → `fixed`（`checkFormCapabilities` 为任一 `defaultValue` 字段增加 `protocolVersion >= 2.7` + `form.controls.advanced` 双门禁，补 2.6 / 2.7 缺 capability / 2.7 完整 capability 三组回归测试）；F-004 → `fixed`（Renderer whitelist 从 form/section/table 扩展至冻结 §5 全部 node type：grid/section/tabs/text/table/recordView/actionButton + form，逐 type 落 dispatch 与测试，不再静默窄于 v0.1.3 初表；无范围缩小故无需 Root 覆盖表修订）。F-005 → `fixed`（Root 00-meta R5 行同步「批次 2c 完成 / 阶段 2 全部落地」）。测试 / 构建 / 运行时：`npm test` 14 文件 **173 项** 全绿（+7：render.test.ts 重写 gate 断言 + render.test.tsx 新增 D-FORM 门禁拒绝与 §5 node dispatch 测试）、`npm run build`（tsc + vite）通过、`go test ./...` / `go build ./...` 通过、Edge headless 实测 `/form-with-reactions`（4 字段 + 无 alert）与 `/form-controls`、`/list-edit-lifecycle`、`/data-table`、`/search-form-table`（0 门禁错误）。D-007 已留痕。`I-PROTO-003` 仍 open（验收/关门门禁，阶段 3/4 未开始）。 |

## A-006 · A-005 finding 闭合复审 + 阶段 2 完成后就绪性（2026-08-01）

- **source**：independent
- **auditor**：Grok 4.5 · xAI · `/audit`
- **类型**：finding-closure + execution-facts
- **scope**：复核 GOAL-007 对 A-005 F-001～F-004（required）与 F-005（recommended）的 `fixed` 闭合证据；核对「阶段 2 全部落地」主张在修复后是否可成立；评估进入阶段 3（结构/行为验证）的就绪性。**不**审计尚未执行的 schema/fixture 正式对照、`I-PROTO-003` 闭合、R5 验收或关门。
- **verdict**：pass

### 范围与区间

当前工作区 `workspace-001-mvp-admin-foundation`，Root `GOAL-001-mvp-admin-foundation`，canonical `docs/workspace-001-mvp-admin-foundation/`。工作区 `plan_refs` / `primary_plan` = `VP-001-mvp-admin-foundation`；`shared_materials_catalog: none`，本意见未把共享资料引用当作事实或关闭证据。范围权威仍为 Root 冻结 [I-PROTO-001 v0.1.3](../GOAL-001-mvp-admin-foundation/attachments/I-PROTO-001-coverage-draft.md)。

本轮只读核验：目标五件套 + `attachments/I-007-001-registry.md` v0.5.0 + 直接读代码（`render.ts` / `render.tsx` / `form-controls.ts` 及对应测试）+ 当前工作树复跑命令。

### 成果（有证据）

**命令复跑（本审计会话，2026-08-01）**

| 命令 | 结果 |
|------|------|
| `cd apps/web && npm test` | 14 文件 / **173** 项全绿 |
| `cd apps/web && npm run build` | tsc + vite 通过 |
| `cd apps/api && go test ./...` | account + handler 包全绿（含 records list/detail/PATCH/DELETE 等） |
| `cd apps/api && go build ./...` | 通过 |

**A-005 F-001～F-005 闭合核对**

| Finding | 声明闭合 | 本轮核对 | 结论 |
|---------|----------|----------|------|
| F-001 非法 action 表达式 fail-open | fixed | `resolveActionGate` 区分 absent → default vs 显式非法 → `ACTION_GATE_EXPRESSION_INVALID`；`tableActionGate` 返回 `{visible, disabled, errors}`（非法 visible → false，非法 disabled → true）；`render.test.ts` 含缺省 vs 非法、非法 visible/disabled 回归 | **闭合成立** |
| F-002 RenderPage 未执行 D-FORM 门禁 | fixed | `FormView` 调用 `gateRenderFormFields`；whitelist / version / capability 拒绝字段并以 `role="alert"` 列出；`render.test.tsx` 覆盖缺 capability、未知 type | **闭合成立** |
| F-003 defaultValue 缺 2.7+advanced | fixed | `checkFormCapabilities` 对任一 `defaultValue` 要求 `protocolVersion >= 2.7` + `form.controls.advanced`；`form-controls.test.ts` 含 2.6 / 缺 advanced / 完整 capability 三组 | **闭合成立** |
| F-004 Renderer whitelist 窄于 §5 | fixed | `WHITELISTED_NODE_TYPES` + `isWhitelistedNodeType` 覆盖 grid/section/tabs/text/table/recordView/actionButton/form；`render.tsx` 逐 type dispatch；测试断言 §5 全 true、chart/modal/upload 仍 false；非缩小范围故无需新版覆盖表 | **闭合成立** |
| F-005 Root R5 投影落后 | fixed | Root `00-meta.md` R5 行已写「阶段 1 + 批次 2a/2b/2c 完成（阶段 2 全部落地），阶段 3/4 未开始」；`progress` 仍 4/6 | **闭合成立** |

**阶段 2 实现事实（与 A-002～A-004 一致，修复后仍成立）**

- 批次 2a：Go records API + Web records/DataTable + `/data-table` `/search-form-table`。
- 批次 2b：D-FORM 控件表面 + D-ACT `runRowAction` + PATCH/DELETE + `/form-controls` `/list-edit-lifecycle`。
- 批次 2c：D-EXPR `reactions.ts` + D-COMP `RenderPage` + `/form-with-reactions`。
- 登记表 [I-007-001-registry.md](attachments/I-007-001-registry.md) v0.5.0；`I-007-001` = verified（**登记层面**，非逐域验证已执行）。
- 阶段 3/4 与父目标 `I-PROTO-003`（required，open）均未声称完成——文档与代码边界一致。

### 对照成功标准（本 scope）

| 标准 | 结论 | 证据 |
|------|------|------|
| A-005 required findings 合法闭合 | 满足 | 上表 F-001～F-004；路径 `fixed` + 代码/测试可重复核对 |
| 阶段 2「全部落地」在修复后可成立 | 满足（条件见下） | 实现 + 173 项 web + build + go test/build；A-005 所指 fail-open / 窄 whitelist 已修 |
| 结构/行为正式验证 | 未到阶段 | schema vendor/pin（`I-PROTO-004`）与纳入 fixtures 对照属阶段 3；本意见不把「未开始」改写为失败 |
| `I-PROTO-003` / R5 验收关门 | 未到阶段且保持阻断 | Root 信息项仍 open/required；本意见不关闭、不放行 |

### Findings

无 **required** finding。开放 required = **0**。

以下为 recommended / 如实记录，**不阻断**阶段 3 启动，亦**不**重新打开 A-005 F-001～F-004：

#### F-001 · ActionButtonView 在 `visible: false` 时不渲染 `gate.errors`（recommended）

- **level**：recommended
- **severity**：low
- **status**：open
- **说明**：`tableActionGate` API 与单元测试已对非法 `visibleWhen` 产出 `ACTION_GATE_EXPRESSION_INVALID`（F-001 核心 fail-closed 已满足）。`ActionButtonView` 在 `!gate.visible` 时直接 `return null`，UI 不展示 errors；非法 `disabledWhen` 路径则会在禁用按钮旁列出 errors。若阶段 3 希望「页面可观察错误」与 API 一致，可在隐藏路径保留可核对 alert（或文档明确仅 API/测试为核对面）。
- **证据**：`apps/web/src/renderer/render.tsx` `ActionButtonView`；对比 `tableActionGate` + `render.test.ts`。

#### F-002 · 范例页直连 `FormControls` 未强制 `gateRenderFormFields` 硬阻断（recommended）

- **level**：recommended
- **severity**：low
- **status**：open
- **说明**：RenderPage 路径已闭合 F-002。`form-controls-page` 调用 `checkFormCapabilitiesRaw` 并展示错误，但 `gates.length > 0` 时仍渲染控件；`list-edit-lifecycle` 编辑表单直接使用 `FormControls` 未走 Renderer 门禁。属范例页一致性/教学路径，不否定 Renderer 级修复；阶段 3 可统一「展示路径必过门禁」纪律。
- **证据**：`apps/web/src/app/examples/form-controls-page.tsx`、`list-edit-lifecycle-page.tsx`。

#### F-003 · 登记表 §4 复用产物计数与 v0.5.0 §2 不一致（recommended）

- **level**：recommended
- **severity**：low
- **status**：open
- **说明**：§2 验证入口已更新为 web 173 项 / form-controls 14 / render 14 / render.tsx 7；§4 复用表仍残留 v0.3/v0.4 计数（如 form-controls 13、render 11、render.tsx 4）。不改变实现事实，建议 `/govern` 同步文案，避免阶段 3 对照时误读。
- **证据**：`attachments/I-007-001-registry.md` §2 vs §4。

### 必改项汇总

无。开放 required = 0。

### 信息门禁（P-005）

| 项 | 状态 | 本 scope |
|----|------|----------|
| `I-007-001` | verified（登记层面） | 不阻断阶段 3 启动 |
| 父目标 `I-PROTO-003` | open / required；最晚 R5 验收前 | **阻断验收/关门**；不阻断阶段 3 实施 |
| 父目标 `I-PROTO-004` | open / non-blocking | 阶段 3 **结构校验实现前**须决策 vendor vs pin（既定时点） |
| 到期 required 信息项 | 无 | 本轮未到 R5 验收 |

### 与既有意见的异同

- **与 A-005**：同 scope 下 A-005 为 conditional（4 required open）；本意见在 D-007 / `fixed` 落地后复核闭合证据，**pass**。不否定 A-005 当时的发现，只确认响应后证据链可重复。
- **与 A-002～A-004（self, pass）**：一致认为阶段 2 有可复核实现与测试；A-005 曾指出 self 漏检 fail-open / 窄 whitelist——本轮确认该缺口已修，与 self 的「批次可推进」结论在修复后可对齐。
- **与响应表（A-005 响应）**：响应声明的测试计数（173）、闭合路径与本轮复跑一致；未发现「声明 fixed 但代码仍 fail-open」的关闭不实。

### 结论 + 建议给编排器/用户的下一步

A-005 四项 required finding 的 `fixed` 闭合可复核，阶段 2 在修复后可作为进入阶段 3 的基线；**verdict = pass**。不得据此关闭 `I-PROTO-003`、放行 R5 验收/关门，或抬升 Root `progress`（仍 4/6）。

建议使用：

```text
/govern workspace-001-mvp-admin-foundation GOAL-007-r5-examples-contract-verification
```

响应 A-006（无 required）：可选同步登记表 §4 文案与 recommended F-001/F-002 UI 一致性；**主路径**进入阶段 3——先决策 `I-PROTO-004`（vendor vs pin），再落地 `node`/`page`/`action`/`reaction` schema 校验与已纳入 fixtures（含 `component-format` 五 case、`reactions`、`request-*`、`table-sort`/`search-table`、`runtime-defaults` 等）的可执行对照。

### 声明

本意见不修改 status/progress、路线图检查点、`I-PROTO-003` / `I-PROTO-004` 状态或 goal-tree；finding 响应与阶段推进由 `/govern` 处理。

## 响应（对 A-006）

| date | actor | summary |
|------|-------|---------|
| 2026-08-01 | `/govern` | 响应 A-006（independent, pass）：用户裁决「不需要自审，直接推进」（P-004 §3.1）。采纳 pass；开放 required=0。A-006 recommended：**F-003**（登记表 §4 计数）→ 随登记表 v0.6.0 **fixed**；**F-001/F-002**（ActionButtonView errors UI / 范例页门禁一致性）→ 跟踪，不阻断阶段 3。主路径：`I-PROTO-004`=**vendor**（D-008）+ 阶段 3 落地——schemas/fixtures vendor+SHA pin；Ajv 结构校验；纳入 fixture suites 执行/覆盖记账（`npm test` **326** 项 / build / go test 全绿）。`I-PROTO-003` 仍 open（验收/关门）；Root `progress` 仍 4/6。 |

## A-007 · R5 阶段 3 独立复核：结构/行为验证完成主张与阶段 4 就绪性（2026-08-01）

- **source**：independent
- **auditor**：Grok 4.5 · xAI · `/audit`
- **类型**：execution-facts
- **scope**：复核 GOAL-007 阶段 3（`I-PROTO-004`=vendor + schemas/fixtures pin + Ajv 结构校验 + 纳入 fixture 对照/覆盖记账）的完成主张是否可成立；评估进入阶段 4（验收 / `I-PROTO-003` 闭合 / 关门）的就绪性与证据缺口。**不**放行 `I-PROTO-003`、R5 验收或关门；**不**抬升 Root `progress`。
- **verdict**：conditional

### 范围与区间

当前工作区 `workspace-001-mvp-admin-foundation`，Root `GOAL-001-mvp-admin-foundation`，canonical `docs/workspace-001-mvp-admin-foundation/`。`plan_refs` / `primary_plan` = `VP-001-mvp-admin-foundation`；`shared_materials_catalog: none`，本意见未把共享资料引用当作事实或关闭证据。

范围权威：

- 冻结覆盖 [I-PROTO-001 v0.1.3](../GOAL-001-mvp-admin-foundation/attachments/I-PROTO-001-coverage-draft.md) §2 fixture 映射与 §5.1 结构验证要求
- GOAL-007 路线图阶段 3 定义与 D-008（vendor + 阶段 3）
- 登记表 [I-007-001-registry.md](attachments/I-007-001-registry.md) v0.6.0

本轮只读核验：目标五件套 + 登记表 + Root `I-PROTO-003`/`I-PROTO-004` 状态 + 直接读代码（`conformance/*`、`upstream/provenance.json`、`docs/schemas/*`）+ **当前工作树复跑命令** + provenance SHA 全量核对。

### 成果（有证据）

**命令复跑（本审计会话，2026-08-01）**

| 命令 | 结果 |
|------|------|
| `cd apps/web && npm test` | **15 文件 / 326 项**全绿（含 `stage3-fixtures.test.ts` **153** 项） |
| `cd apps/web && npm run build` | tsc + vite 通过 |
| `cd apps/api && go test ./...` | account + handler 包全绿 |
| `cd apps/api && go build ./...` | 通过 |
| provenance SHA 全量 | **20/20** artifacts 与本地文件一致（`sourceCommit=ca9e5fe…`，`artifactVersion=2.7.0`） |

**阶段 3 可复核落地（与 02-execution / D-008 一致）**

| 面 | 证据 |
|----|------|
| I-PROTO-004=vendor | Root D-012 + GOAL-007 D-008；`docs/schemas/` 含 node/page/action/reaction/component-registry（+既有 app-manifest）；fixtures 于 `apps/web/src/protocol/upstream/*.cases.json` |
| SHA pin | `apps/web/src/protocol/upstream/provenance.json` 锁定 6 schema + 14 fixture suites；本轮 20/20 哈希一致 |
| 结构校验 | `conformance/schema-validate.ts`（Ajv draft-07）；page/node 有合法样本 `ok=true` 与缺 type / 缺 protocolVersion 拒绝断言 |
| 行为对照（已执行） | component-format 5、query-serialization 16、static-data 9、request-lifecycle 4、runtime-defaults 9、response-mapping 23、search-table 11、table-sort 14、version-negotiation 44、actions 11（均 `assertCoverage` 全量执行） |
| 覆盖记账（未执行 host） | reactions **16/16 排除**（MVP `$context` vs 上游 multi-round `$deps`，D-008 已决策）；request-construction **75/75**（batch=Q1 排除 + non-batch=deferred 统一引擎） |
| 登记表 | v0.6.0 逐域现状与验证入口已更新；§4 计数与 §2 对齐（A-006 F-003 响应闭合） |
| 门禁诚实性 | 父目标 `I-PROTO-003` 仍 **open/required**；阶段 4 **未**声称完成；Root `progress` 仍 **4/6** |

**Ajv 探针（本轮）**

- `reaction` 缺 `dependencies` → `ok: false`；含 `dependencies: []` 合法样本 → `ok: true`
- `action` `{type,method,url}` request 形状 → `ok: true`

### 对照成功标准（阶段 3 / 进入阶段 4）

| 标准 | 结论 | 证据 |
|------|------|------|
| 结构验证可执行（node/page schema） | **部分满足** | page/node 样本 + 拒绝路径有断言；action/reaction 导出可调但测试断言薄弱（见 F-002） |
| component-format 五 case 对照 | **满足** | 5/5 执行通过 |
| 行为验证与 R2 基线一致、不越界 | **条件满足** | 已执行 suite 有对照；排除/deferred 有记账理由；**未**静默伪通过 reactions/request-construction |
| 冻结 include suite 均有可执行对照或合法边界 | **缺口** | `reactions`、`request-construction` 为冻结 §2 **include**，执行数为 **0**（见 F-001） |
| 每纳入域范例路径 + 验证入口（I-PROTO-003） | **未闭合**（正确） | 登记表 + 阶段 3 证据可支撑验收准备，但 Root `I-PROTO-003` 仍 open；不得据此关门 |
| 排除项明确 | **满足** | D-UPLOAD、batch Q1、scenarios 非自动化门禁与文档一致 |

### Findings

#### F-001 · 冻结 include suite `request-construction` / `reactions` 执行数为 0，阶段 3「完成」易被误读为全 suite 已对照

- **level**：required
- **severity**：medium
- **status**：open
- **影响门禁**：阶段 3 完成认定的无条件表述；阶段 4 / `I-PROTO-003` 闭合叙述（若把「阶段 3 完成」直接当作全部冻结行为入口已绿）
- **证据**：
  - [I-PROTO-001 v0.1.3 §2](../GOAL-001-mvp-admin-foundation/attachments/I-PROTO-001-coverage-draft.md)：`request-construction` → **include**（D-DATA）；`reactions` → **include**（D-EXPR）
  - `stage3-fixtures.test.ts`：两 suite 均 `assertCoverage(..., [], exclusions)`，**executedIds 为空**；reactions 16 全排除、request-construction 75 全 deferred/Q1
  - 登记表 v0.6.0 已诚实写「记账排除 / deferred」，但 `00-meta` / `02-execution` / goal-tree 摘要统一写「阶段 3 **完成**」，未附 include-suite 执行矩阵
  - D-008 对 reactions 的「vendor+account 排除、拒绝伪通过」决策成立且可核对；request-construction non-batch 的「deferred 统一引擎」是工程延期，**不是**冻结 Q1 式边界排除
- **必改**（进入阶段 4 / 闭合 `I-PROTO-003` 前至少满足其一组合）：
  1. 在登记表（或 `02-execution`）落盘 **include suite 执行矩阵**：suite → executed N / excluded(reason) / deferred(reason)，避免「阶段 3 完成 = 全部 include 已绿」；
  2. 对 `request-construction` **non-batch** 子集：落地可执行 host/adapter 对照（`fixed`），或用户书面 **`accepted-residual`**（范围 + 复审触发），或新决策将冻结 suite 改为 include-partial 并写清边界；
  3. 对 `reactions`：将 MVP `$context` 验证入口（`reactions.test.ts` + `form-with-reactions`）在登记表升格为 **正式 D-EXPR 行为验证入口**，并显式声明上游 multi-round suite **不在** MVP 语义子集（与 D-008 对齐）；若需与冻结 §2 `include` 字面完全一致，则 residual 或修订冻结。
- **禁止**：在未处理上述项时，用「阶段 3 完成 + 326 全绿」直接关闭 `I-PROTO-003` 或主张「全部纳入 fixture 已对照通过」。

#### F-002 · action / reaction 结构校验缺少「合法样本必须 ok」的可核对断言

- **level**：required
- **severity**：medium
- **status**：open
- **影响门禁**：阶段 3 结构验证完成认定（D-VAL / D-ACT / D-EXPR schema 面）；§5.1「结构契约随纳入域带入」
- **证据**：
  - `stage3-fixtures.test.ts`「validates a minimal reaction document」使用**缺少** `dependencies` 的文档，却只断言 `typeof ok.ok === "boolean"`，**不**要求 `ok === true` 或 `ok === false`
  - 本轮探针：该样本实际 `ok: false`（`must have required property 'dependencies'`）；合法 `{dependencies:[], when, fulfill}` 才 `ok: true`——测试未锁住任一路径
  - **action** schema 在 `validateAgainstSchema` 中可编译且本轮探针 request/navigate 合法样本 `ok: true`，但 stage3 **无任何** action 正/负向断言
  - page/node 路径有完整正/负向断言，形成对比
- **必改**：为 `reaction` 与 `action` 各补至少：1 条合法样本 `expect(ok).toBe(true)` + 1 条非法样本 `expect(ok).toBe(false)`（reaction 合法样本须含 `dependencies`）；不得用「typeof boolean」充当结构校验完成证据。

#### F-003 · Conformance 适配器与生产 host 路径分离（recommended）

- **level**：recommended
- **severity**：medium
- **status**：open
- **说明**：query-serialization / response-mapping / table-sort 等 stage3 绿测跑在 `conformance/*` 纯适配器上；生产 `records.ts` 的 `buildRecordsQuery` / `parseRecordList` 是独立实现，**不**调用 `serializeQuery` / `mapResponse`。阶段 3 证明「冻结 fixture 语义可复现」，**不**自动证明生产路径与 ADR 全量一致。验收叙述应区分「conformance oracle」与「host 单元/范例证据」。
- **证据**：`conformance/query-serialize.ts` vs `renderer/records.ts`；`conformance/response-mapping.ts` vs `parseRecordList`。

#### F-004 · `component-registry.json` 已 vendor+pin 但未参与 membership 校验（recommended）

- **level**：recommended
- **severity**：low
- **status**：open
- **说明**：冻结 §5.1 要求 type 解析到白名单 **与 registry**；当前 Renderer 用硬编码 `WHITELISTED_NODE_TYPES`，registry 文件仅 pin 存在。硬编码表与 §5 一致则风险可控；阶段 4 可补「registry 键 ⊆/⊇ §5」对照或加载 registry 做 membership。
- **证据**：`docs/schemas/component-registry.json` 在 provenance；`render.ts` `WHITELISTED_NODE_TYPES`；无读取 registry 的校验代码。

#### F-005 · A-006 recommended F-001/F-002 仍开放（recommended · 跟踪）

- **level**：recommended
- **severity**：low
- **status**：open
- **说明**：ActionButtonView 在 `visible:false` 时不展示 `gate.errors`；部分范例页直连 `FormControls` 未硬阻断。A-006 响应已标跟踪、不阻断阶段 3；阶段 4 可一并处理或继续 residual 跟踪。
- **证据**：A-006 响应表；`render.tsx` / `form-controls-page.tsx`。

### 必改项汇总

1. **F-001**：include suite 执行矩阵 + `request-construction` non-batch 与 `reactions` MVP 入口的 fixed / residual / 冻结修订之一；禁止用「阶段 3 完成」静默覆盖 0 执行 suite。
2. **F-002**：为 action/reaction 结构校验补合法/非法可核对断言。

开放 required = **2**。可继续准备阶段 4 材料，但 **不得**无条件认定阶段 3 行为/结构证据已完整覆盖冻结 include 面，**不得**关闭 `I-PROTO-003` / 放行 R5 验收或关门。

### 信息门禁（P-005）

| 项 | 状态 | 本 scope |
|----|------|----------|
| `I-007-001` | verified（登记 + 阶段 3 执行证据） | 不单独等于 `I-PROTO-003` |
| 父目标 `I-PROTO-003` | **open / required**；最晚 R5 验收前 | **阻断验收/关门** |
| 父目标 `I-PROTO-004` | **verified**（vendor） | 阶段 3 前置已满足 |
| 到期 required 信息项 | 无（未进入正式验收） | 进入阶段 4 时 `I-PROTO-003` 即到期门禁 |

### 与既有意见的异同

- **与 A-006（independent, pass）**：一致认为 A-005 修复可重复、可进入阶段 3；本意见在阶段 3 **落地后**复核，发现结构断言空洞与 include suite 0 执行缺口 → **conditional**。
- **与 D-008 / 02-execution**：采纳 vendor、reactions 不伪通过、request-construction 记账的方向正确；本意见要求把「记账/deferred」提升为验收前可裁决的显式矩阵与闭合路径，而不是否定阶段 3 实现事实。
- **与 A-002～A-004（self）**：阶段 2 实现证据仍成立；本意见不重开 A-005 已 fixed 项。
- **无同 scope self 审计**覆盖阶段 3：用户此前对 A-006 裁决跳过自审；本独立意见作为阶段 3 交叉台账。P-004：若用户还要阶段 3 自审，由 `/govern` 询问，本意见不代替。

### 结论 + 建议给编排器/用户的下一步

阶段 3 的 vendor/pin、Ajv 入口、多数 include suite 对照与 326 测试/构建证据**真实可复核**；reactions 排除与 request-construction 记账**未**伪装为 pass。但冻结 include 面存在 **0 执行** suite，且 action/reaction 结构断言不完整，故 **verdict = conditional**。不得关闭 `I-PROTO-003`，不得放行 R5 验收/关门，不得抬升 Root `progress`（仍 4/6）。

建议使用：

```text
/govern workspace-001-mvp-admin-foundation GOAL-007-r5-examples-contract-verification
```

响应 A-007：优先 `fixed` 路径闭合 F-002（补 action/reaction 结构断言）；对 F-001 落盘 suite 执行矩阵并选择 request-construction non-batch 与 reactions MVP 入口的 fixed / accepted-residual / 冻结修订；然后进入阶段 4 验收材料与 `I-PROTO-003` 闭合评审（仍须用户确认）。P-004：可询问是否需要阶段 3 自审。

### 声明

本意见不修改 status/progress、路线图检查点、`I-PROTO-003` / `I-PROTO-004` 状态或 goal-tree；finding 响应与阶段推进由 `/govern` 处理。

## 响应（对 A-007）

| date | actor | summary |
|------|-------|---------|
| 2026-08-01 | `/govern` | 响应 A-007（independent, conditional）：用户裁决「不需要自审，直接推进」（P-004 §3.1）。**F-002 → `fixed`**：action/reaction 结构正负向断言。**F-001 → 组合闭合**：① §2b 执行矩阵；② `reactions` MVP 正式入口升格；③ `request-construction` 初记 residual（D-009）后由用户澄清 **更正为 fixed**（见下行 D-010）。 |
| 2026-08-01 | `/govern` | **D-010 更正**：`request-construction` non-batch → **`fixed`**。落地 `conformance/request-construction.ts`（`constructRequest`）；stage3 执行 non-batch **64/64** 对照；batch **11** 仍 Q1 排除。登记表 **v0.8.0**。证据：`npm test` **395** 项 / build 全绿。开放 required 仍 **0**。**不**关闭 `I-PROTO-003`；阶段 4 未开始；Root `progress` 仍 4/6。 |

## A-008 · R5 阶段 4 关门自审：验收证据与 I-PROTO-003 闭合（2026-08-01）

- **source**：self
- **auditor**：Grok 4.5 · xAI · `/govern`
- **类型**：close-out
- **scope**：GOAL-007 成功标准全量对照；阶段 1–4 证据链；父目标 `I-PROTO-003` 闭合合法性；是否具备用户授权 `done` 的条件。**不**修改 status 为 `done`；**不**抬升 Root `progress`；**不**放行 VP/R6。
- **verdict**：pass

### 范围与区间

当前工作区 `workspace-001-mvp-admin-foundation`，Root `GOAL-001-mvp-admin-foundation`，canonical `docs/workspace-001-mvp-admin-foundation/`。`plan_refs` / `primary_plan` = `VP-001-mvp-admin-foundation`；`shared_materials_catalog: none`。范围权威为冻结 [I-PROTO-001 v0.1.3] 与登记表 [I-007-001-registry.md](attachments/I-007-001-registry.md) **v0.8.0**。

本意见在阶段 4 复跑命令与 D-011 / Root D-013 落盘后书写。

### 成果（有证据）

**命令复跑（本会话，2026-08-01）**

| 命令 | 结果 |
|------|------|
| `cd apps/web && npm test` | **15 文件 / 395 项**全绿（含 stage3 **222**） |
| `cd apps/web && npm run build` | tsc + vite 通过 |
| `cd apps/api && go test ./...` | 全绿 |
| `cd apps/api && go build ./...` | 通过 |

**阶段与决策**

| 阶段 | 状态 | 关键证据 |
|------|------|----------|
| 1 登记 | 完成 | `I-007-001` verified；登记表 v0.8.0 |
| 2 范例实现 | 完成 | 批次 2a/2b/2c + A-005 fixed |
| 3 结构/行为 | 完成 | vendor+Ajv+conformance；A-007/D-010 required 闭 |
| 4 验收 | 材料完成 | D-011；逐域核对表（02-execution）；Root D-013 |

**意见台账**：A-001～A-007 相关 required 均已合法闭合；本 scope 开放 required finding = **0**。

**信息门禁**

| 项 | 状态 |
|----|------|
| `I-007-001` | verified |
| `I-PROTO-003` | **verified**（Root D-013） |
| `I-PROTO-004` | verified |
| 到期未处理 required 信息项 | **无** |

### 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| 每纳入域范例路径 + 验证入口 | **达成** | 登记表 §1 十一域 + 02 验收表 |
| 未覆盖域前后端范例 | **达成** | D-DATA/TABLE/FORM/ACT/EXPR + Go records API |
| 结构验证可执行 | **达成** | Ajv page/node/action/reaction 正负向；component-format 5/5 |
| 行为验证不越 v0.1.3 | **达成** | §2b 矩阵；排除项诚实记账 |
| `I-PROTO-003` verified；无开放 required finding | **达成** | D-013；开放 required=0 |
| 目标 status=`done` | **未执行** | 须用户另行授权；本自审不改 status |

### Findings

无 **required** finding。开放 required = **0**。

以下 recommended / 跟踪（**不阻断**用户授权 `done`）：

#### F-001 · Conformance oracle 与生产 host 路径分离（recommended · 沿用 A-007 F-003）

- **level**：recommended
- **severity**：medium
- **status**：open
- **说明**：阶段 3 绿测证明冻结 fixture 语义可复现；生产 `records.ts` 等未必调用同一 adapter。R6/产品化应区分叙述或逐步统一。

#### F-002 · component-registry membership 未加载校验（recommended · 沿用 A-007 F-004）

- **level**：recommended
- **severity**：low
- **status**：open
- **说明**：硬编码 `WHITELISTED_NODE_TYPES` 与 §5 一致则风险可控；可补 registry 键对照。

#### F-003 · 范例页 / ActionButton errors UI 一致性（recommended · 沿用 A-006/A-007）

- **level**：recommended
- **severity**：low
- **status**：open
- **说明**：RenderPage 门禁已闭；部分范例直连 `FormControls`、隐藏 action 不展示 gate.errors——教学/UI 一致性，非 I-PROTO-003 缺口。

### 必改项汇总

无。

### 与既有意见的异同

- **与 A-007（conditional）**：required 已由 D-009/D-010 fixed；本意见确认阶段 4 验收后证据仍成立 → **pass**。
- **与 A-002～A-004（self）**：阶段 2 实现证据仍成立；本意见覆盖全目标成功标准与信息门禁。
- **与 Root D-013**：采纳 `I-PROTO-003` verified；本意见不反向打开。

### 结论 + 建议下一步

GOAL-007 在 R5 范围内**具备关门条件**（成功标准可勾选、`I-PROTO-003` verified、无开放 required finding、关门向自审 pass）。**verdict = pass**。

建议：

1. 可选：`/audit` 独立关门复审（与 GOAL-006 对称）。
2. 用户授权后：GOAL-007 → `done`；Root 纲领 R5 → 完成；`progress` → **5/6**；再规划 R6。
3. recommended 跟踪至 R6/产品化，不阻断本目标 `done`。

### 声明

本意见不修改 status/progress；`done` 与 Root 检查点变更须 `/govern` + 用户确认。

## A-009 · R5 独立关门复审：验收证据、信息门禁与当前运行态（2026-08-01）

- **source**：independent
- **auditor**：GitHub Copilot
- **类型**：close-out
- **scope**：GOAL-007 的关门条件，包括成功标准、阶段 3/4 可执行证据、`I-007-001` 与父目标 `I-PROTO-003` / `I-PROTO-004`、既有 required finding 的闭合，以及当前 Web / API 运行态。**不**修改 GOAL-007 `status`、Root `progress`、VP 或 R6 状态。
- **verdict**：pass

### 范围与区间

当前工作区为 `workspace-001-mvp-admin-foundation`，Root 为 `GOAL-001-mvp-admin-foundation`，canonical 范围为 `docs/workspace-001-mvp-admin-foundation/`。工作区的 `plan_refs` / `primary_plan` 指向 `VP-001-mvp-admin-foundation`，该 VP 的 `vision_ref` 与 active Charter `schema-ui-core-admin-foundation@0.1.0` 一致。`shared_materials_catalog: none`，本次没有把共享资料引用当作范围事实或关门证据。

本轮独立核对了 GOAL-007 五件套、登记表 v0.8.1、Root D-013、A-007/A-008 的 finding 响应、阶段 3 conformance 实现和当前工作树的命令结果。

### 成果（有证据）

- **信息门禁**：`I-007-001`、父目标 `I-PROTO-003`（Root D-013）及 `I-PROTO-004` 均为 `verified`；其证据分别可回指登记表、D-013 和 vendor + SHA pin。没有到期且未闭合的 required 信息项。
- **逐域验收**：登记表 §1 覆盖 11 个纳入域，且 §2b 明确区分已执行、排除与正式补充入口。`request-construction` non-batch 为 **64/64** 可执行对照；batch 的 11 项仍按 Q1 排除。上游 `reactions` 的 16 项被诚实记为 MVP `$context` 子集外，正式验证入口为 `reactions.test.ts` 与 `/form-with-reactions`，没有把该排除表述为 host 对照通过。
- **结构与固定协议**：`stage3-fixtures.test.ts` 对 page/node/action/reaction 均含合法通过与非法拒绝断言。复核 `provenance.json` 的 **20/20** 个 vendored schema/fixture，SHA-256 mismatch = **0**，来源仍为 `schema-ui-docs` commit `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b`、artifact version `2.7.0`。
- **当前可执行证据**：本审计会话复跑 `cd apps/web && npm test`，结果为 **15 文件 / 395 项**通过（stage3 **222**）；`npm run build` 成功。`cd apps/api && go test ./...` 与 `go build ./...` 均成功。
- **运行时入口**：预览服务的 `/.well-known/schema-ui/app-manifest.json`、`/` 与 `/form-with-reactions` 均返回 HTTP 200；manifest 为 `protocolVersion=2.7`、9 个页面、`homePageRef=overview`。

### 对照成功标准

| 标准 | 结论 | 证据 |
|------|------|------|
| 11 个纳入域均有范例路径与可执行验证入口 | 满足 | 登记表 v0.8.1 §1、§2、§2b |
| 未覆盖域具备可运行的前后端范例 | 满足 | records API、五条 Examples 路由、Web/API 测试与构建 |
| 结构与行为验证不超出冻结边界 | 满足 | Ajv 正负向断言、fixture 执行矩阵、Q1/reactions 边界 |
| `I-PROTO-003` 在验收前合法闭合 | 满足 | Root D-013、D-011、登记表 v0.8.1 与本轮复跑 |
| 无开放 required finding 或 required 信息项阻断关门 | 满足 | A-007 required 已由 D-009/D-010 fixed；A-008 与本轮复核均为 pass |
| 完整协议、VP 或 Root 已关门 | 未主张 | 本审计不将 R5 证据扩大为这些结论 |

### Findings

无 required finding。开放 required = **0**。

A-008 的 recommended 跟踪项（conformance oracle 与生产 host 分离、registry membership、范例页/ActionButton UI 一致性）仍保持 open，但其范围不构成 `I-PROTO-003` 或 GOAL-007 关门门禁，且本轮未发现将它们伪装为已修正或“完整协议支持”的表述。

### 必改项汇总

无。GOAL-007 的关门范围内不存在需要以 `fixed`、`accepted-residual` 或 `user-overruled` 路径处理的开放 required finding。

### 与既有意见的异同

- 与 A-007（independent，conditional）一致承认其当时指出的阶段 3 缺口；D-009/D-010 的 `fixed` 响应已补足 action/reaction 结构断言、执行矩阵和 non-batch request-construction 对照，本轮直接复跑后未发现 reopening 证据。
- 与 A-008（self，pass）同向，但本意见重新执行 Web/Go 验证、核对 vendored artifact SHA 和预览运行态，而非仅复述其结论。
- 没有与本关门 scope 相冲突的 relevant verdict，也没有开放 required finding 需要 P-004 裁决。

### 结论 + 建议给编排器/用户的下一步

GOAL-007 的 R5 关门条件可复核成立，**verdict = pass**。本意见不授权或执行 `done`；由用户通过 `/govern` 决定是否将 GOAL-007 标为 `done`，并在授权后把 Root 的 R5 纲领检查点更新为完成、派生 `progress` 更新为 `5/6`。R6、VP 关门和 recommended 跟踪项仍应作为后续范围处理。

### 声明

本意见不修改 status/progress、路线图检查点、Root/VP 状态或 goal-tree；finding 响应与正式关门由 `/govern` 处理。

## 响应（对 A-009）

| date | actor | summary |
|------|-------|---------|
| 2026-08-01 | `/govern` | 响应 A-009（independent, pass）：采纳 **pass**；开放 required=0。recommended 跟踪项（A-008 F-001～F-003 / A-007 F-003～F-005）**不阻断**关门，继续跟踪至 R6/产品化。用户书面授权 GOAL-007 → **`done`**；Root 纲领 R5 → **完成**；派生 `progress` → **5/6**（goal-tree 同步）。D-012 / Root D-014 已留痕。**不**放行 R6、VP 关门或 Root `done`；不主张完整协议支持。 |
