---
id: GOAL-007-r5-examples-contract-verification
title: R5 · 纳入域范例与契约验证
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-08-01
version: 0.10.0
---

# GOAL-007 · R5 · 纳入域范例与契约验证

## 概述

承接 R2 冻结的 MVP 覆盖基线（[I-PROTO-001 v0.1.3](../GOAL-001-mvp-admin-foundation/attachments/I-PROTO-001-coverage-draft.md)）与 R3/R4 已交付的 Admin 外壳、账号权限链路，为每一**纳入域**交付可观察的范例页（或场景）与可执行的结构/行为验证路径，并在 R5 验收前闭合父目标 `I-PROTO-003`（required）。

范围依据：Root 将 R5 定义为「纳入域范例与契约验证」；协议清单 [protocol-inventory-v2.7.0.md](../../vision/protocol-inventory-v2.7.0.md) §3 为每个 domain_id 登记了 React / Go / 范例页面候选 / 验证路径，§2.5 提供可复用的信息性场景。固定协议版本沿用 `schema-ui-docs` artifact `2.7.0`、source commit `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b`。

## 范围边界

### 纳入

- 为 **11 个纳入域**（D-NODE、D-EXPR、D-COMP、D-DATA、D-ACT、D-PERM、D-APP、D-TABLE、D-FORM、D-VER、D-VAL）登记范例页/场景路径 + 结构/行为验证入口。
- 范例页候选源：协议清单 §2.5 信息性场景与 [I-PROTO-001 v0.1.3] §3（P0/P1/P2 候选），优先复用 `05-scenarios` 与上游 `_samples` 结构模板并去业务化改写。
- 结构验证：对纳入页面运行 `node.schema.json` / `page.schema.json`；行为验证：对应已纳入 fixture suite（`component-format` 五 case、`request-construction`、`response-mapping`、`query-serialization`、`static-data`、`reactions`、`table-sort`、`search-table`、`version-negotiation`、`runtime-defaults`、`app-manifest`、`app-navigation`、`permissions-inheritance` 等，见 v0.1.3 §2）。
- 复用并登记 R3/R4 已有产物：D-APP 外壳/导航、D-PERM 权限链路已具备运行路径，作为对应域范例的既有证据，不重复实现。

### 排除

- **D-UPLOAD**（exclude）及其 `uploads` fixture 整域。
- 完整 `component-registry`、未列入 v0.1.3 §5 白名单的 type；多选批量 action/request 语义（Q1=否）。
- `scenarios` 不作为独立自动化门禁（Q5=否）；仅作范例/手工路径。
- 「完整协议支持」主张、VP 关门证据与 R6 集成验收，均不在本目标范围。

## 高层路线图（P-001）

1. **契约发现与登记**：**完成**；`I-007-001` 登记表已落盘 [attachments/I-007-001-registry.md](attachments/I-007-001-registry.md)（2026-07-31），逐纳入域登记范例路径 + 结构/行为验证入口，对齐 [I-PROTO-001 v0.1.3] §3 候选与协议清单 §2.5；D-APP/D-PERM 复用产物与可执行验证命令已核验（`npm test` 94 项 / `npm run build` / `go test ./...` / `go build ./...`）。
2. **范例页/场景实现**：**批次 2a（D-DATA/D-TABLE）、批次 2b（D-FORM/D-ACT）与批次 2c（D-EXPR/D-COMP）完成**；按登记表为未覆盖域落地可观察范例页（含必要的 React 页面/组件与 Go 数据/动作路径支撑），复用 R3/R4 已有产物。批次 2a：Go `GET /api/records`（list/detail）+ Web `records.ts`/`use-records.ts` + `data-table.tsx` + `search-form-table`/`data-table` 范例页；`npm test` 114 项 / `go test` 21 项 / build / Edge 实测全绿。批次 2b：D-FORM §5 白名单控件表面（`form-controls.ts`/`.tsx`，2.6/2.7 版本/capability 门禁）+ D-ACT 非批量动作（`row-action.ts` 复用 R4 `executeAction`）+ Go `PATCH`/`DELETE /api/records/{id}` + `form-controls`/`list-edit-lifecycle` 范例页；`npm test` 138 项 / `go test` 18 顶层 / build 全绿。批次 2c：D-EXPR 反应引擎（`reactions.ts` 复用 `evaluateExpression`）+ D-COMP 最小 Renderer 接线（`render.ts`/`.tsx`，resolve R4 F-002）+ `form-with-reactions` 范例页；`npm test` 166 项 / `npm run build` / `go test` / Edge 实测全绿。**阶段 2 全部落地**；`I-PROTO-004` 在阶段 3 结构校验实现前决策。
3. **结构/行为验证**：**完成（2026-08-01）**；`I-PROTO-004`=vendor；schemas/fixtures vendor+SHA pin；Ajv + conformance 对照；A-007/D-010 后 `npm test` **395** 项（stage3 **222**）。include suite 以登记表 **§2b** 为准：多数 suite 已执行；`reactions` 上游 0 执行（MVP 正式入口 = `reactions.test.ts` + `/form-with-reactions`）；`request-construction` non-batch **64/64 fixed**（batch Q1 排除）。**禁止**将「全绿」读成 batch 或上游 multi-round reactions 已对照。
4. **验收与关门**：**未开始**；闭合父目标 `I-PROTO-003`（每纳入域范例路径 + 验证入口均有可核对证据）→ 自审/关门审计。

## 成功标准（暂定 · 可随契约发现细化）

- [ ] 每个 include / include-partial 域至少有一条范例页（或场景）路径与一条可执行验证入口（自动化命令或手工步骤）。
- [ ] 未覆盖域（如 D-DATA / D-TABLE / D-FORM / D-EXPR / D-ACT）具备可运行的前后端范例路径；R3/R4 已有产物以登记与复核方式纳入，不重复实现。
- [ ] 结构验证可执行：对纳入范例页面运行 `node`/`page` schema；`component-format` 五个已纳入 case 对照通过。
- [ ] 行为验证与 R2 冻结基线一致：纳入 fixture suite 运行结果不超出 v0.1.3 边界；排除项（D-UPLOAD、多选批量等）明确不在验证范围。
- [ ] 父目标 `I-PROTO-003` 在 R5 验收前合法闭合（`verified`）；关门前无开放 required finding。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-007-001 | required | 每条纳入能力的范例页路径与自动化/手工验证入口？ | R5 验收 / 关门 | R5 验收前 | 对照 [I-PROTO-001 v0.1.3] §3 与协议清单 §2.5，登记每纳入域的范例路径与验证命令/步骤 | **verified** | 2026-07-31 登记；2026-08-01 阶段 3 + A-007/D-010：登记表 **v0.8.0**（§2b 矩阵、reactions 正式入口、request-construction non-batch fixed）；**不**单独等于 `I-PROTO-003` 闭合 | [attachments/I-007-001-registry.md](attachments/I-007-001-registry.md) v0.8.0 |

父目标 `I-PROTO-003`（required，R5 验收/关门门禁）由 Root 维护，本目标通过范例路径 + 阶段 3 可执行验证为其提供证据，**验收前**须正式闭合。`I-PROTO-004`（vendor vs pin）→ **verified**（2026-08-01 D-008 vendor）。

## 父目标

- [GOAL-001-mvp-admin-foundation](../GOAL-001-mvp-admin-foundation/00-meta.md)

## 备注

- 立项日期：2026-07-31，承接 R4 关门（GOAL-006 `done`）。
- 批次 2a 落地（2026-07-31）：D-DATA/D-TABLE 范例 + Go 列表/详情支撑实现完成（详见 02-execution）；登记表升 v0.2.0。`I-PROTO-003` 仍 open，验收前须以阶段 3 可执行证据闭合。
- 批次 2b 落地（2026-07-31）：D-FORM 控件表面 + D-ACT 非批量动作 + Go PATCH/DELETE + `form-controls`/`list-edit-lifecycle` 范例页实现完成（详见 02-execution）；A-003 批次自审（self）pass；登记表升 v0.3.0。`I-PROTO-003` 仍 open，验收前须以阶段 3 可执行证据闭合。`I-PROTO-004` 仍 open（non-blocking，阶段 3 结构校验前决策）。
- 批次 2c 落地（2026-08-01）：D-EXPR 反应引擎（`reactions.ts`，复用 `evaluateExpression`）+ D-COMP 最小 Renderer 接线（`render.ts`/`render.tsx`，resolve R4 F-002）+ `form-with-reactions` 范例页实现完成（详见 02-execution）；A-004 批次自审（self）pass；登记表升 v0.4.0。**阶段 2 全部落地**；`I-PROTO-003` 仍 open，验收前须以阶段 3 可执行证据闭合。`I-PROTO-004` 仍 open（non-blocking，阶段 3 结构校验前决策）。
- A-005 响应（2026-08-01）：A-005（independent, conditional）对阶段 2 完成主张提出 F-001～F-004（required）与 A-002～A-004（self, pass）同 scope 分歧；用户裁决「不需要自审，直接推进」，F-001～F-004 按 **`fixed`** 合法闭合（action 表达式 fail-closed、RenderPage 接入 D-FORM 门禁、defaultValue 2.7+advanced 双门禁、Renderer whitelist 扩展至冻结 §5 全部 node type），F-005 同步；`npm test` 173 项 / build / `go test` / Edge 实测全绿；登记表升 v0.5.0；D-007 已留痕。
- A-006 响应 + 阶段 3（2026-08-01）：A-006（independent, pass）闭合复核成立；用户裁决「不需要自审，直接推进」；`I-PROTO-004`=vendor（D-008）；阶段 3 结构/行为验证落地（schemas+fixtures vendor、Ajv、conformance 适配器）；`npm test` **326** 项 / build / go test 全绿；登记表 v0.6.0。阶段 4 与 `I-PROTO-003` 仍未闭合。
- A-007 响应（2026-08-01）：A-007（independent, conditional）F-002 **fixed**、F-001 矩阵 + reactions 正式入口；request-construction 初记 residual 后由用户澄清 → **D-010 fixed**（non-batch 64 执行）；`npm test` **395** 项；登记表 v0.8.0。开放 required=0。阶段 4 与 `I-PROTO-003` 仍未闭合。
- Root 纲领进度仍为 `4/6`；本目标推进不抬升 progress，不放行 Root `done`；`I-PROTO-004` 已 verified，`I-PROTO-003` 仍 open。
- R4 关门时登记的 recommended 跟踪项（F-002 Renderer 接线 / F-003 token 会话 / F-004 双端一致性 oracle）标注「随 R5 / 生产化 / `I-PROTO-004` 解决」——本目标为这些跟踪项提供落地窗口，但 recommended 不阻断 R5 推进。
- 结论与进度只写已发生事实；「未开始」阶段不得写成已完成。
