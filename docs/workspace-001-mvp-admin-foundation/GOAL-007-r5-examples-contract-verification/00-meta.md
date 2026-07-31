---
id: GOAL-007-r5-examples-contract-verification
title: R5 · 纳入域范例与契约验证
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-07-31
version: 0.1.0
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

1. **契约发现与登记**：**未开始**；建立本目标 `I-007-001`（每纳入域范例路径 + 自动化/手工验证入口登记表），对齐 [I-PROTO-001 v0.1.3] §3 候选与协议清单 §2.5。
2. **范例页/场景实现**：**未开始**；按登记表为未覆盖域落地可观察范例页（含必要的 React 页面/组件与 Go 数据/动作路径支撑），复用 R3/R4 已有产物。
3. **结构/行为验证**：**未开始**；`node`/`page` schema 校验、`component-format` 五 case 与已纳入 fixtures 对照，登记可执行验证命令/步骤。
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
| I-007-001 | required | 每条纳入能力的范例页路径与自动化/手工验证入口？ | R5 验收 / 关门 | R5 验收前 | 对照 [I-PROTO-001 v0.1.3] §3 与协议清单 §2.5，登记每纳入域的范例路径与验证命令/步骤 | open | — | 待 R5 契约发现登记 |

父目标 `I-PROTO-003`（required，R5 验收/关门门禁）由 Root 维护，本目标通过闭合 `I-007-001` 为其提供证据。`I-PROTO-004`（vendor vs pin，non-blocking）仍 open；不阻断 R5 规划，实施前可结合 schema/fixture 校验策略一并处理。

## 父目标

- [GOAL-001-mvp-admin-foundation](../GOAL-001-mvp-admin-foundation/00-meta.md)

## 备注

- 立项日期：2026-07-31，承接 R4 关门（GOAL-006 `done`）。
- Root 纲领进度仍为 `4/6`；本目标立项不抬升 progress，不放行 Root `done`，不改变 `I-PROTO-003` / `I-PROTO-004` 状态。
- R4 关门时登记的 recommended 跟踪项（F-002 Renderer 接线 / F-003 token 会话 / F-004 双端一致性 oracle）标注「随 R5 / 生产化 / `I-PROTO-004` 解决」——本目标为这些跟踪项提供落地窗口，但 recommended 不阻断 R5 推进。
- 结论与进度只写已发生事实；「未开始」阶段不得写成已完成。
