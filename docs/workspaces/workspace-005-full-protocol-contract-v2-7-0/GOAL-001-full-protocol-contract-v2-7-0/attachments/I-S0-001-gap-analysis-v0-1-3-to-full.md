---
title: S0 差距盘点 · I-PROTO-001 v0.1.3 (MVP) vs schema-ui-docs@v2.7.0 整份契约
status: active
doc_type: evidence-gap-analysis
created: 2026-08-08
updated: 2026-08-08
parent: GOAL-001-full-protocol-contract-v2-7-0
version: 0.1.0
related_info: I-001
related_goal_checkpoint: S0
---

# S0 差距盘点（I-001 闭合证据）

> 本文件是 S0 检查点的事实证据：以 **可复核的本地资产**（`docs/schemas/` 6 个结构契约、`apps/web/src/protocol/upstream/` 15 套 vendor fixture + provenance、`docs/vision/protocol-inventory-v2.7.0.md` 全量清单、`I-PROTO-001 v0.1.3` 冻结表、`apps/web` / `apps/api` 代码现状）逐一比对，产出**可审计差集**：前端保真债 vs 未纳入 type 分列。上游 pin 仓库（`ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b`）本次可达，`uploads` / `permissions-inheritance` fixture 与关键规范文档已拉取复核（S2/S3 实现权威）。

## 1. 基线对照

| 项 | I-PROTO-001 v0.1.3（MVP 冻结，只读） | 整份 v2.7.0 契约（VP-006 目标） |
|----|--------------------------------------|---------------------------------|
| 能力域 | 7 include + 4 include-partial + 1 exclude | 12 域**逐项** disposition，默认 include |
| Fixture 套件 | 15 vendor + permissions-inheritance（GOAL-006 附件）；reactions 全排除；request-construction batch 排除；uploads 未 vendor | 16 行为套件全部可执行（含 uploads）；scenarios 保持 support-only |
| registry type | 18/24 白名单 | 24/24 |
| 表达引擎 | `$context` 显隐/禁用（`==`/`!=`/`contains` 子集） | `$deps`/`$self`/`$context`/`$row` 全语法 + 快照多轮求值 |

## 2. 能力域差集（domain × disposition）

| domain_id | v0.1.3 | 整份契约目标 | 差集性质 | S0 证据 |
|-----------|--------|--------------|----------|---------|
| D-NODE | include | include | **无缺口** | node/page schema 结构校验全绿（stage3-fixtures 20 项 + representative pages） |
| D-EXPR | include（$context 子集语义） | include（整引擎） | **保真债**：`reactions` 16/16 排除；缺多轮 `$deps` 快照引擎（fulfill.value/otherwise/observers/baselines/externalUpdates/loop protection/深等检测/`MULTIPLE_VALUE_WRITES` 警告/码点字符串比较/严格类型） | `upstream/reactions.cases.json`（16 case）；`stage3-fixtures.test.ts` L351-368 全排除声明；上游 `docs/02-reaction-expression.md` |
| D-COMP | include-partial（18/24） | include（24/24） | **未纳入 type**：`statCard`、`chart`（数据展示，supportsData）；`inputNumber`、`datePicker`、`dateRangePicker`（表单基座） | `docs/schemas/component-registry.json` 24 key；`render.ts` WHITELISTED_NODE_TYPES 8 项；`form-controls.ts` 10 项 |
| D-DATA | include | include | **无缺口** | request-construction 64/75 非 batch case 全绿；response-mapping 23、static-data 9、query-serialization 16 全绿 |
| D-ACT | include-partial（非批量） | include（含批量） | **保真债**：`batchRequest` 11 case 排除（Q1=否）；批量 Trigger 执行序（requiresSelection/EMPTY_SELECTION/selectionAfterSuccessReload）未实现 | `upstream/request-construction.cases.json`（11 batch case 清单见 §4）；ADR-0022；`07-actions-contract.md` §3.5 |
| D-PERM | include | include | **无缺口** | permissions-inheritance 17/17 case 全绿（GOAL-006 附件 SHA 钉死 `ac124fa1…`）；含执行序（visible→permission→disabled→confirm） |
| D-APP | include | include | **无缺口** | app-manifest 37 + app-navigation 16 全绿；Manifest 真实端点 `/api/schema/*` 与 `/.well-known/schema-ui/app-manifest.json` 挂载 |
| D-TABLE | include-partial（无批量多选） | include（含 selection） | **保真债**：selection 状态机在 conformance 层已有（search-table 11 case 全绿含 selection 事件），但 UI 无多选列、toolbar 无 requiresSelection/批量触发、batchMapping 未接入 | `conformance/search-table.ts`；`schema-table.tsx`（无 selection props）；ADR-0022 D2 |
| D-FORM | include-partial（14/18 控件 + defaultValue） | include（全部） | **未纳入 type**：`inputNumber`、`datePicker`、`dateRangePicker`、`upload` | `form-controls.ts` FORM_CONTROLS 10 项 vs registry 20 表单 type |
| D-UPLOAD | **exclude** | **include** | **未纳入 type + 整域新增**：上传控件、顶层 `type: upload` action、capability `actions.upload` 门禁、客户端编排（13 case）、后端上传端点 | `upstream/` 无 uploads；上游 `07-actions-contract.md` §7 + ADR-0012 + `uploads/cases.json`（13 case，本次拉取核对） |
| D-VER | include | include | **无缺口** | version-negotiation 44 + runtime-defaults 9 全绿 |
| D-VAL | include | include | **无缺口** | 6 schema 全 vendor + SHA pin + 结构校验执行 |

### 2.1 汇总计数（S1 覆盖表输入）

- **无缺口域**：7（D-NODE、D-DATA、D-PERM、D-APP、D-VER、D-VAL + D-COMP 已纳入 18 type 的主路径）
- **保真债**：4 项——D-EXPR 整引擎、D-ACT 批量执行、D-TABLE 多选 UI/批量、D-COMP statCard/chart 渲染
- **未纳入 type**：6 项——statCard、chart、inputNumber、datePicker、dateRangePicker、upload
- **整域新增**：1 项——D-UPLOAD（前端 + 后端 + fixture + 范例）
- **v0.1.3 排除面转 include**：D-UPLOAD（exclude→include）；reactions 16 case、batchRequest 11 case、6 个 registry type

## 3. 前端保真债 vs 未纳入 type 分列（S0 要求的分类）

| 类别 | 清单 | 性质 |
|------|------|------|
| **A. 前端保真债**（能力已纳入、语义/交互深度不足） | ① `$deps` 多轮 reaction 引擎（16 case）；② 表格多选 UI + 批量 Trigger 执行序（11 batchRequest case + ADR-0022 D2/D4/D5）；③ statCard/chart 数据展示渲染；④ inputNumber/datePicker/dateRangePicker 控件 wire 与门禁 | 提升到契约语义可验证，不要求视觉产品化（VP-005 冻结） |
| **B. 未纳入 type**（v0.1.3 白名单外） | ① `upload` 控件（D-UPLOAD 整域）；② 后端上传端点；③ uploads fixture 13 case | v0.1.3 明确 exclude；整份契约默认 include |

## 4. Fixture 套件逐项现状（17 套）

| suite | 上游 case 数 | 本地执行 | v0.1.3 | 整份契约处置（S1 提案） |
|-------|-------------|----------|--------|--------------------------|
| actions | 11 | ✅ 11/11 | include | include（无变化） |
| app-manifest | 37 | ✅ | include | include |
| app-navigation | 16 | ✅ | include | include |
| component-format | 5 | ✅ | include | include |
| permissions-inheritance | 17 | ✅（GOAL-006 附件 SHA 钉死） | include | include（vendor 于本区 upstream 并继续 SHA pin） |
| query-serialization | 16 | ✅ | include | include |
| reactions | 16 | ❌ 0/16（全排除） | 排除 | **include（整引擎实现）** |
| request-construction | 75（含 11 batch） | ✅ 64/75 | include-partial | **include（batch 11 case 纳入）** |
| request-lifecycle | 4 | ✅ | include | include |
| response-mapping | 23 | ✅ | include | include |
| runtime-defaults | 9 | ✅ | include | include |
| search-table | 11 | ✅（含 selection 状态机） | include | include |
| static-data | 9 | ✅ | include | include |
| table-sort | 14 | ✅ | include | include |
| version-negotiation | 44 | ✅ | include | include |
| uploads | 13 | ❌ 未 vendor | exclude | **include（vendor + 实现）** |
| scenarios | — | support-only | support-only | support-only（范例源） |

## 5. Registry type 差集（24 key）

| 类别 | v0.1.3 已纳入（18） | 整份契约新增（6） |
|------|--------------------|-------------------|
| 布局 | grid、section、tabs | — |
| 数据/操作 | text、table、recordView、actionButton | **statCard、chart** |
| 表单基座 | form、input、select、textarea、switch、checkbox、radio、select.multiple | **inputNumber、datePicker、dateRangePicker** |
| 表单 2.7 | cascader、checkboxGroup、richText、password | — |
| 上传 | — | **upload** |

## 6. 后端（Go）差集

| 面 | 现状 | 整份契约缺口 |
|----|------|--------------|
| 列表/详情/动作 | 通用 resource CRUD 工厂（list/detail/create/update/delete）+ 权限门 + operationlog | **无**（users/roles 实例可验证） |
| Manifest/Schema | `/.well-known/schema-ui/app-manifest.json` + `/api/schema/{pageId}` 真实端点 | **无** |
| 批量服务端契约 | 无 | **新增**：批量端点（如 `POST /api/{resource}/batch-delete`，按 `$selection.keys` 语义，权限门 + 成功后清选） |
| 上传服务端契约 | 无（settings logo 仅 URL 文本） | **新增**：`multipart/form-data` 上传端点（fieldName 可配、size/type 独立校验、响应 `{url,id,name,size}`、语义错误码 FILE_TOO_LARGE / UNSUPPORTED_FILE_TYPE / STORAGE_UNAVAILABLE） |
| 权限模型 | account/permission + auth 中间件 | **无**（D-PERM 渲染层语义 + 后端鉴权已分列，见协议总纲） |

## 7. 结论

1. 差集**全部可纳入**：无任何域需要 `exclude` 或范围收缩；S1 覆盖表可默认 `include` 全量（12/12 域、24/24 registry type、16/16 行为 fixture 套件）。
2. 实现批次建议（I-004 输入）：
   - **批次 1（D-FORM/D-COMP 控件与展示）**：inputNumber、datePicker、dateRangePicker（wire + 门禁 + 渲染 + 范例）、statCard、chart（数据展示 + 范例）
   - **批次 2（D-EXPR）**：`$deps` 全语法表达式解析/校验 + 快照多轮 reaction 引擎 + 16 case 全绿
   - **批次 3（D-TABLE/D-ACT 批量）**：多选 UI + selection 状态机接线 + batchMapping/requiresSelection/EMPTY_SELECTION + 11 batchRequest case + 后端批量端点
   - **批次 4（D-UPLOAD）**：upload 控件 + `type: upload` action + `actions.upload` 门禁 + 13 case 编排 + 后端上传端点 + 范例
3. `scenarios` 保持 support-only（范例源），不设独立自动化门禁（Q5 延续）。
4. 回归基线：v0.1.3 主路径（64/75 request-construction + 全量其余套件 + Go 全量测试）保持全绿不回退。

## 8. 勘误（A-001 F-003 响应，2026-08-08）

- §2 D-UPLOAD 行「`upstream/` 无 uploads」为盘点时点事实；**S1 冻结期间**已 vendor `uploads.cases.json`（13 case，SHA `aaeb968369e145d422d163d2f086b0ed3754b4486ddb889ca30375b64c8e22e2`）与 `permissions-inheritance.cases.json`（17 case，SHA `ac124fa1…`）至 `apps/web/src/protocol/upstream/` 并更新 `provenance.json`。本附件其余盘点结论不受影响（差集判断未变）。
