---
title: I-007-001 · R5 纳入域范例路径与验证入口登记表
status: active
doc_type: info-registry
created: 2026-07-31
updated: 2026-08-01
parent: GOAL-007-r5-examples-contract-verification
version: 0.6.0
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
| D-NODE | include | 任意合法 page（全站基座，无专属业务页） | `node`/`page` schema（Ajv，`docs/schemas/*.json` vendor） | —（基座） | **阶段 3 已接入**：`validateAgainstSchema` + sample 白名单 page；`npm test` stage3 structural | I-PROTO-004=vendor（2026-08-01） |
| D-EXPR | include | `form-with-reactions`（改写，去业务化） | `reaction.schema.json`（vendor） | fixtures `reactions`：**vendor+SHA pin**；**全量记账排除**（上游 multi-round `$deps` 字段写引擎 ∉ MVP `$context` visible/disabled 子集） | **已实现（批次 2c）** + 阶段 3 覆盖记账 | MVP 复用 `evaluateExpression`；上游 suite 不伪映射为 pass |
| D-COMP | include-partial | `data-table` / `form-controls-*`（§5 白名单 type） | `node`/`page` schema + registry membership（§5） | fixtures `component-format`：**5/5 执行通过**（`conformance/component-format.ts`） | **已实现（批次 2c + A-005 响应）** + 阶段 3 component-format 对照 | 白名单 §5；schema vendor |
| D-DATA | include | `search-form-table` / `data-table`（Go 列表/详情 + 前端） | `page` schema | fixtures：`query-serialization` **16/16**、`static-data` **9/9**、`response-mapping` **23/23**；`request-construction` vendor+记账（batch Q1 排除，其余 deferred 统一引擎） | **已实现（批次 2a + 2b 扩展）** + 阶段 3 对照 |  |
| D-ACT | include-partial | `row-backend-actions` / `admin-list-edit-lifecycle`（非批量） | `action.schema.json`（vendor） | fixtures `actions` **11/11**、`request-lifecycle` **4/4**（非批量 transport/lifecycle） | **已实现（批次 2b）** + 阶段 3 对照 | Q1=否排除批量 request-construction |
| D-PERM | include | `permission-inheritance`（**已有**：`apps/web/src/renderer/permissions.ts` + `permissions-inheritance.test.ts`） | `validatePermissions`（L2 fail-closed 校验，permissions.ts） | fixtures `permissions-inheritance`（GOAL-006 attachments `dperm/cases.json` 17 例，SHA-256 pin 于测试） | **已有（R4）** | 仅登记复核，不重复实现 |
| D-APP | include | Admin 外壳 + 导航（**已有**：`apps/web/src/app/App.tsx`、`app/navigation.ts`、`public/.well-known/schema-ui/app-manifest.json`） | `docs/schemas/app-manifest.schema.json`（已 vendor，SHA pin） | fixtures `app-manifest`、`app-navigation`（`apps/web/src/protocol/upstream/*.cases.json` + `upstream-fixtures.test.ts`） | **已有（R3）** | 仅登记复核，不重复实现 |
| D-TABLE | include-partial | `data-table` / `search-table`（排序声明 + 基础列表交互） | `page` schema | fixtures `table-sort` **14/14**、`search-table` **11/11** | **已实现（批次 2a）** + 阶段 3 对照 | 排除多选批量（Q1=否） |
| D-FORM | include-partial | `form-controls-advanced`（§5 全部 2.6/2.7 控件） | `node`/`page` schema + registry + 版本/capability 下限 | `component-format`（格式子集）；scenarios（手工路径，Q5=否） | **已实现（批次 2b + A-005 响应）**：`apps/web/src/renderer/form-controls.ts`（`FormControlType` §5 白名单、`wireKindOf`、`coerceFieldValue`、`validateDefaultValue` fail-closed、`checkFormCapabilities` 版本/capability 门禁，含 **defaultValue 2.7+advanced 双门禁**）+ `form-controls.tsx` `FormControls` 组件（10 种控件）+ `form-controls.test.ts` 14 项；范例页 `form-controls` | `defaultValue` 属性规则随 2.7 |
| D-VER | include | 全站（`supportedCapabilities` / 版本协商）——已有子集：`validateAppManifest` + `upstream-fixtures.test.ts` `negotiateFixture` | — | fixtures `version-negotiation` **44/44**、`runtime-defaults` **9/9** | **阶段 3 已接入** + R3 app-manifest negotiate 子集 | |
| D-VAL | include | —（构建时/加载时结构校验） | 6 schemas **全部 vendor**（`docs/schemas/`） | Ajv `validateAgainstSchema` | **阶段 3 已接入** | I-PROTO-004=vendor |

## 2. 验证入口命令登记（已存在的可执行入口）

| 入口 | 命令 | 覆盖域 | 证据 / 说明 |
|------|------|--------|-------------|
| Web 单元 + upstream fixtures + stage3 | `cd apps/web && npm test` | D-APP、D-PERM、D-EXPR、D-DATA、D-TABLE、D-FORM、D-ACT、D-COMP、D-VER、D-VAL、D-NODE | **15 测试文件 / 326 项**（既有 173 + stage3 **153**；含 `conformance/stage3-fixtures.test.ts`、`upstream-fixtures.test.ts`、permissions / app-manifest / render / form-controls / reactions 等）；provenance SHA pin 扩展 |
| Web 构建 | `cd apps/web && npm run build` | D-APP、D-DATA、D-TABLE、D-FORM、D-ACT、D-VAL | tsc + vite build |
| Go 测试 | `cd apps/api && go test ./...` | D-PERM、D-VER(子集)、D-DATA、D-ACT | records 12 项 + account/permission 等 |
| Go 构建 | `cd apps/api && go build ./...` | D-PERM、D-VER(子集)、D-DATA、D-ACT | server 可运行 |
| Go 数据 API 冒烟 | `curl http://127.0.0.1:8080/api/records` 等 | D-DATA、D-TABLE、D-ACT | 批次 2a/2b 运行时证据，见 02-execution |
| 浏览器渲染 | `/data-table`、`/search-form-table`、`/list-edit-lifecycle`、`/form-controls`、`/form-with-reactions` | D-DATA、D-TABLE、D-APP、D-FORM、D-ACT、D-EXPR、D-COMP | Edge headless / jsdom 表面测试，见 02-execution |
| 结构 schema 校验 | `npm test` → stage3 structural + `schema-validate.ts`（Ajv） | D-NODE、D-VAL、D-COMP | vendor schemas in `docs/schemas/` |

**阶段 3 已接入**（2026-08-01，`I-PROTO-004`=vendor）：见上表与 `apps/web/src/protocol/conformance/`。仍排除/记账：`reactions`（MVP 子集外）、`request-construction` batch（Q1）与非 batch 统一引擎 deferred。

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
| D-FORM / D-ACT（批次 2b 新增） | `apps/web/src/renderer/{form-controls.ts,form-controls.tsx,row-action.ts}`、`app/examples/{form-controls-page,list-edit-lifecycle-page}.tsx`、`apps/api/internal/handler/records.go`（PATCH/DELETE） | `npm test`（form-controls **14** 项 + row-action 5 项 + records PATCH/DELETE 4 项 + app-examples）+ `go test` + jsdom |
| D-EXPR / D-COMP（批次 2c 新增） | `apps/web/src/renderer/{reactions.ts,render.ts,render.tsx}`、`app/examples/form-with-reactions-page.tsx` | `npm test`（reactions 12 + render **14** + render.tsx **7** + app-examples）+ build + Edge headless |
| 阶段 3 conformance | `apps/web/src/protocol/conformance/*`、`docs/schemas/*`、`upstream/*.cases.json`、`provenance.json` | `npm test` stage3 **153** 项；结构 Ajv + 纳入 suite 执行/覆盖记账 |

## 5. 变更规则

- 本表是 `I-007-001` 的信息登记，**不是** R5 关门证据；`I-PROTO-003` 的闭合仍须阶段 3 的可执行验证结果与阶段 4 验收。
- 覆盖范围变更须改 [I-PROTO-001 v0.1.3]（新决策 + 新版本），并同步重估本表。
- 阶段 2/3 实现落盘后，更新本表对应行的「现状」与验证命令为已发生事实。

> 变更：v0.4.0（2026-08-01）——批次 2c 实现落盘：D-EXPR 行「现状」→ 已实现（`reactions.ts` + `form-with-reactions` 范例页），D-COMP 行「现状」→ 已实现（`render.ts`/`render.tsx` 最小 Renderer，resolve F-002），P2 `form-with-reactions` 标记已实现；验证命令计数更新（web 166 项 / 14 文件），浏览器渲染新增 `/form-with-reactions` Edge headless 实测。

> 变更：v0.6.0（2026-08-01）——阶段 3 + I-PROTO-004=vendor：schemas/fixtures vendor+SHA pin；结构 Ajv 与纳入 suite 对照入账；web **326** 项 / 15 文件；§4 计数与 §2 对齐（A-006 F-003）。

> 变更：v0.5.0（2026-08-01）——A-005 响应落盘：D-COMP 行更新为「Renderer whitelist = 冻结 §5 全部 node type（grid/section/tabs/text/table/recordView/actionButton + form）」，D-FORM 行更新为含 defaultValue 2.7+advanced 双门禁；验证命令计数更新（web 173 项 / 14 文件），`form-controls.test.ts` 14 项、`render.test.ts` 14 项、`render.test.tsx` 7 项。

> 变更：v0.3.0（2026-07-31）——批次 2b 实现落盘：D-FORM/D-ACT 行「现状」→ 已实现（控件表面 + `runRowAction` + Go PATCH/DELETE + 范例页 `form-controls`/`list-edit-lifecycle`），D-DATA 行扩展 PATCH/DELETE，验证命令计数更新（web 138 项 / go 18 顶层）；P1 `admin-list-edit-lifecycle` 标记已实现。

> 变更：v0.2.0（2026-07-31）——批次 2a 实现落盘：D-DATA/D-TABLE 行「现状」→ 已实现，新增验证入口（web 114 项 / go 21 项 / API 冒烟 / 浏览器渲染）；阶段 3 待接入清单移出已落地映射路径。
