---
title: I-007-001 · R5 纳入域范例路径与验证入口登记表
status: active
doc_type: info-registry
created: 2026-07-31
updated: 2026-07-31
parent: GOAL-007-r5-examples-contract-verification
version: 0.1.0
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
| D-EXPR | include | `form-with-reactions`（改写，去业务化） | `reaction.schema.json`（随 I-PROTO-004） | fixtures `reactions`（阶段 3 接入 suite） | 规划 | 已有子集：导航 `visibleWhen` 表达式（`apps/web/src/protocol/app-manifest.ts` `evaluateExpression`）；Go `apps/api/internal/account/permission.go` `Evaluate` |
| D-COMP | include-partial | `data-table` / `form-controls-*`（§5 白名单 type） | `node`/`page` schema + registry membership（§5） | fixtures `component-format`（5 case：currency/percent/datetime） | 规划 | 白名单 §5 已冻结；`apps/web/src/renderer/` 为空分层（R5 实现 Renderer 接线） |
| D-DATA | include | `search-form-table` / `data-table`（Go 列表/详情 + 前端） | `page` schema | fixtures `request-construction`、`response-mapping`、`query-serialization`、`static-data` | 规划 | Go 现仅 `accounts/me` + `health`；需补列表/详情 API 支撑路径 |
| D-ACT | include-partial | `row-backend-actions` / `admin-list-edit-lifecycle`（非批量） | `action.schema.json`（随 I-PROTO-004） | fixtures `actions`、`request-lifecycle`（非批量子集，阶段 3 接入） | 规划 | D-PERM `executeAction` 时序引擎已有（R4，`apps/web/src/renderer/permissions.ts`）；Q1=否排除批量 |
| D-PERM | include | `permission-inheritance`（**已有**：`apps/web/src/renderer/permissions.ts` + `permissions-inheritance.test.ts`） | `validatePermissions`（L2 fail-closed 校验，permissions.ts） | fixtures `permissions-inheritance`（GOAL-006 attachments `dperm/cases.json` 17 例，SHA-256 pin 于测试） | **已有（R4）** | 仅登记复核，不重复实现 |
| D-APP | include | Admin 外壳 + 导航（**已有**：`apps/web/src/app/App.tsx`、`app/navigation.ts`、`public/.well-known/schema-ui/app-manifest.json`） | `docs/schemas/app-manifest.schema.json`（已 vendor，SHA pin） | fixtures `app-manifest`、`app-navigation`（`apps/web/src/protocol/upstream/*.cases.json` + `upstream-fixtures.test.ts`） | **已有（R3）** | 仅登记复核，不重复实现 |
| D-TABLE | include-partial | `data-table` / `search-table`（排序声明 + 基础列表交互） | `page` schema | fixtures `table-sort`、`search-table` | 规划 | 排除多选批量（Q1=否）；依赖 D-DATA Go API + 前端表格组件 |
| D-FORM | include-partial | `form-controls-extended` / `form-controls-advanced`（§5 全部 2.6/2.7 控件） | `node`/`page` schema + registry + 版本/capability 下限 | `component-format`（格式子集）；scenarios（手工路径，Q5=否） | 规划 | 白名单 §5 已冻结；`defaultValue` 属性规则随 2.7 |
| D-VER | include | 全站（`supportedCapabilities` / 版本协商）——已有子集：`validateAppManifest` + `upstream-fixtures.test.ts` `negotiateFixture` | — | fixtures `version-negotiation`（已在 app-manifest fixtures 覆盖）、`runtime-defaults`（待接入） | 部分已有 | 版本协商已随 D-APP fixtures 落地；`runtime-defaults` 阶段 3 接入 |
| D-VAL | include | —（构建时/加载时结构校验） | `docs/06-validation.md` + 6 schemas（已 vendor `app-manifest.schema.json`） | — | 部分已有 | 其余 5 schemas 随 I-PROTO-004 决策 vendor/pin |

## 2. 验证入口命令登记（已存在的可执行入口）

| 入口 | 命令 | 覆盖域 | 证据 / 说明 |
|------|------|--------|-------------|
| Web 单元 + upstream fixtures | `cd apps/web && npm test` | D-APP、D-PERM、D-EXPR(子集) | 6 测试文件 / 94 项（含 `upstream-fixtures.test.ts`、`permissions-inheritance.test.ts`、`app-manifest.test.ts`、`navigation.test.ts`、`App.integration.test.tsx`、`account/context.test.ts`）；上游 fixture SHA 与 provenance 已 pin |
| Web 构建 | `cd apps/web && npm run build` | D-APP | tsc + vite build |
| Go 测试 | `cd apps/api && go test ./...` | D-PERM、D-VER(子集) | `internal/account/permission_test.go`、`internal/handler/account_test.go`、`health_test.go` |
| Go 构建 | `cd apps/api && go build ./...` | D-PERM、D-VER(子集) | server 可运行 |

**阶段 3 待接入的验证入口**（登记为计划，执行结果未发生）：`reactions`、`component-format`（5 case）、`request-construction`、`response-mapping`、`query-serialization`、`static-data`、`actions`、`request-lifecycle`（非批量）、`table-sort`、`search-table`、`runtime-defaults`；以及 `node`/`page`/`action`/`reaction` schema 校验命令（依赖 I-PROTO-004 vendor/pin 决策）。

## 3. 与 I-PROTO-001 v0.1.3 §3 候选的对齐

- P0 候选：`permission-inheritance`（D-PERM，**已有**）、Admin 外壳 + 导航（D-APP，**已有**）。
- P1 候选：`search-form-table` / `data-table`（D-DATA、D-TABLE）、`admin-list-edit-lifecycle` / `list-detail`（D-ACT、D-DATA、D-FORM）→ 登记为规划，阶段 2 实现。
- P2 候选：`form-with-reactions`（D-EXPR、D-FORM）→ 登记为规划。
- 不做（MVP）：`form-with-upload`、`form-controls-advanced/extended` 全量、`order-*` 业务样例原文 → 与 §5 白名单、非目标一致，不登记。

## 4. 复用产物（不重复实现）

| 域 | 复用路径 | 复核点 |
|----|----------|--------|
| D-APP | `apps/web/src/app/App.tsx`、`app/navigation.ts`、`protocol/app-manifest.ts`、`public/.well-known/schema-ui/app-manifest.json` | `npm test` / `npm run build` 通过即为复用前提未退化 |
| D-PERM | `apps/web/src/renderer/permissions.ts`、`account/context.ts`、`apps/api/internal/account/{session,permission}.go`、`handler/account.go` | `npm test`（permissions-inheritance 17 例）+ `go test ./...` 通过 |

## 5. 变更规则

- 本表是 `I-007-001` 的信息登记，**不是** R5 关门证据；`I-PROTO-003` 的闭合仍须阶段 3 的可执行验证结果与阶段 4 验收。
- 覆盖范围变更须改 [I-PROTO-001 v0.1.3]（新决策 + 新版本），并同步重估本表。
- 阶段 2/3 实现落盘后，更新本表对应行的「现状」与验证命令为已发生事实。
