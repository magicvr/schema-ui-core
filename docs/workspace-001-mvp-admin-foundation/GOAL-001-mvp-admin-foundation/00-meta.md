---
id: GOAL-001-mvp-admin-foundation
title: MVP Admin 基架
status: active
parent: null
created: 2026-07-31
updated: 2026-07-31
version: 0.11.0
progress: 4/6
plan_refs: VP-001-mvp-admin-foundation
primary_plan: VP-001-mvp-admin-foundation
serves_summary: 交付 VP-001 可 fork 的 React+Go Admin MVP，以 schema-ui-docs@v2.7.0 为协议边界
---

# GOAL-001 · MVP Admin 基架

## 概述

在本工作区交付可 fork 的 React 前端与 Go 后端 Admin MVP：以固定的 `schema-ui-docs@v2.7.0`（commit `ca9e5fe…`）为兼容边界，完成核心账号与权限链路，并使每一**纳入范围**的协议能力具备可观察范例与验证路径。

本目标是工作区 Root（`parent: null`），服务意图 [VP-001-mvp-admin-foundation](../../vision/plans/VP-001-mvp-admin-foundation.md)。

## 愿景对齐

| 字段 | 值 |
|------|-----|
| Charter | `schema-ui-core-admin-foundation@0.1.0` |
| `plan_refs` | `VP-001-mvp-admin-foundation` |
| `primary_plan` | `VP-001-mvp-admin-foundation` |
| `serves_summary` | 交付 VP-001 可 fork 的 React+Go Admin MVP，以 schema-ui-docs@v2.7.0 为协议边界 |
| 工作区 | `workspace-001-mvp-admin-foundation`（`vision_role: primary`） |

不在此扩写第二套愿景边界；协议清单权威见 [protocol-inventory-v2.7.0.md](../../vision/protocol-inventory-v2.7.0.md)。

## 成功标准（方向级 · 可验证）

- [x] 存在可运行、可 fork 的 React 前端与 Go 后端工程骨架，并以固定协议版本为兼容边界（文档/配置可指回 pin）。
- [x] MVP 协议覆盖子集已书面冻结（`I-PROTO-001` verified）。
- [ ] 每一纳入项具备可核对的前后端实现路径（R3-R5）。
- [x] 核心账号与权限链路具备可验证的前后端集成（对照 `D-PERM` / `I-PROTO-002`）。
- [ ] 每一纳入能力有可观察范例页面（或场景）与可执行验证入口（`I-PROTO-003`）。
- [x] 未主张“支持全部协议功能”；未纳入项有明确边界说明。

## 纲领路线图（P-001）

| 阶段 | 名称 | 状态 | 说明 |
|------|------|------|------|
| R1 | 工程骨架与仓库约定 | **完成** | 子目标 GOAL-002/003/004 均 `done`（各 A-003 self pass）；Root A-001 independent pass → A-002/D-006 已响应并维持；`apps/*` + monorepo 约定已交付；I-STACK-001/002 verified |
| R2 | MVP 协议覆盖子集冻结 | **完成** | 用户书面确认按 `I-PROTO-001-coverage-draft.md` v0.1.3 冻结；D-009 记录决定，A-005 已以 `fixed` 闭合 F-001/F-002，`I-PROTO-001` = `verified` |
| R3 | Admin 外壳与导航 | **完成** | GOAL-005 已记录 D-005 冻结、R3 实现、73 项测试、构建、fixture/HTTP 入口复核；A-006 实施 self-audit 通过，F-003 已 fixed，浅色/深色基线可后置产品化 VP |
| R4 | 核心账号与权限 | **完成** | 依赖 R2；GOAL-006 D-004 冻结最小 API 与 D-PERM 映射，`I-006-001` verified；`I-PROTO-002`（R4 **实施**门禁）已于方案冻结时闭合；R4 实施（Go 会话/鉴权、Web `$context`、D-PERM 引擎、17 例 fixture）与关门（A-004 self + A-005 independent 均 pass）完成，GOAL-006 `done` |
| R5 | 纳入域范例与契约验证 | 阶段 1 + 批次 2a/2b/2c 完成（阶段 2 全部落地） | 每纳入域范例页 + 结构/行为验证路径（`I-PROTO-003`）；GOAL-007 阶段 1（`I-007-001` 登记表）与阶段 2 批次 2a（D-DATA/D-TABLE）、批次 2b（D-FORM/D-ACT）、批次 2c（D-EXPR/D-COMP）已完成（2026-07-31/08-01），阶段 3/4 未开始 |
| R6 | 集成验收与 VP 证据 | 未开始 | 对照 VP 退出判据收集工作区证据；不自动改 VP status |

纲领阶段串行；同一阶段内可并行子目标。阶段完成须更新本表标记，并不得假装未满足退出条件。

## 派生进度展示

纲领检查点 **4/6** 完成（R1、R2、R3、R4）→ frontmatter `progress: 4/6`（goal-tree 同步）。progress **不**放行后续阶段、不关闭 finding、不推导 Root `done`。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-PROTO-001 | required | 哪些 `domain_id` / fixture suite 纳入 VP-001 MVP？ | R2 冻结基线 / 后续实施范围 | R2 结束前 | 对照 [protocol-inventory](../../vision/protocol-inventory-v2.7.0.md) §3 决策并落盘纳入/排除表 | **verified** | 2026-07-31 用户书面确认 v0.1.3；覆盖变更须新决策与版本 | 冻结基线：[attachments/I-PROTO-001-coverage-draft.md](attachments/I-PROTO-001-coverage-draft.md) v0.1.3；D-009；A-005 以 `fixed` 留痕 F-001/F-002 |
| I-PROTO-002 | required | 账号权限最小 API 与 `D-PERM` 映射是否完整？ | R4 实施 | R4 实施前 | 设计最小 API + 对照 permissions-inheritance fixtures | **verified** | 2026-07-31 方案冻结时闭合（GOAL-006 D-004）；闭合仅覆盖设计/映射，不放行实施本身 | GOAL-006 D-004 + [attachments/dperm/](../../workspace-001-mvp-admin-foundation/GOAL-006-r4-account-permission/attachments/dperm/) 固定资料（cases.json 17 例、node.schema、ADR-0023 等，SHA-256 核验） |
| I-PROTO-003 | required | 每条纳入能力的范例页路径与自动化/手工验证入口？ | R5 验收 / 关门 | R5 验收前 | 为纳入域登记范例路径与验证命令/步骤 | open | — | 待 R2 纳入表 |
| I-PROTO-004 | non-blocking | 是否 vendor 上游 schemas/fixtures，或 pin 远程校验？ | R1–R5 工程策略 | R1/实施前为宜 | 决策 vendor vs pin；记录维护成本 | open | — | 待确认 |
| I-STACK-001 | required | 前端/后端具体脚手架与包管理（目录布局、模块边界）？ | R1 实施 | R1 实施前 | 用户确认或有界 spike 后写入决策 | verified | — | 2026-07-31 D-004：`apps/web`+`apps/api`；Web=npm+Vite+React+TS+Tailwind/shadcn；API=Go modules；结构参考平行仓择优移植 |
| I-STACK-002 | non-blocking | monorepo vs 前后端分仓、默认端口与 env 约定 | R1 约定 | R1 内 | 决策落盘即可 | verified | — | 2026-07-31 D-004：本仓 monorepo `apps/*`；默认端口/env 在 GOAL-002/003 约定中细化 |

## 父目标

- `null`（Root）

## 备注

- 开区日期：2026-07-31。
- Charter H-001：清单提取 verified；覆盖子集冻结仍 open（本目标 `I-PROTO-001`）。
- recommended 愿景项 `F-V003`（双线分支契约）不在本 Root 门禁内；后续双线 VP 前由 `/vision` 处理。
- R1 子目标（2026-07-31）：`GOAL-002-r1-repo-layout-conventions`、`GOAL-003-r1-api-go-scaffold`、`GOAL-004-r1-web-react-scaffold`。
- R1 独立复核：Root A-001 pass；编排响应 A-002 + D-006（2026-07-31）。R2 覆盖基线随后由 D-009 冻结；本事实不放行 R3-R5 实施。
- R4 规划（2026-07-31）：GOAL-005 A-007 independent 关门复审 `pass`、A-008 响应其 recommended 后，立项 `GOAL-006-r4-account-permission`（`active`）。`I-006-001` 在 R4 方案冻结前验证；`I-PROTO-002` 保持 open 作为 R4 **实施**门禁。A-007 F-002（schema 等价性校验）已登记为关闭 `I-PROTO-004` 时的跟进项。
- R4 方案冻结（2026-07-31）：GOAL-006 D-004 冻结账号权限最小 API 与 `D-PERM` 映射（对照 `permissions-inheritance` fixtures 17 例等固定资料，SHA-256 落盘 `GOAL-006/attachments/dperm/`），`I-006-001` → `verified`，`I-PROTO-002` → `verified`（实施门禁闭合，仅覆盖设计/映射结论）。R4 **实施未开始**；`progress` 保持 `3/6`，不放行实施，不推导 `done`。
- R4 实施与关门（2026-07-31）：GOAL-006 实施完成（Go 会话/鉴权、Web `$context`、D-PERM 引擎、17 例 fixture、94 项 web 测试、构建与 HTTP 运行时证据）；A-004 关门自审（self）与 A-005 独立关门复审（independent）均 pass、开放 required=0；经用户 `/govern` 授权 GOAL-006 → `done`，纲领 R4 → 完成，`progress` → `4/6`（goal-tree 同步）。recommended 跟踪项 F-002～F-004 随 R5 / 生产化 / `I-PROTO-004` 解决。
- R5 立项（2026-07-31）：经用户 `/govern` 确认 slug，立项 `GOAL-007-r5-examples-contract-verification`（`active`）进入 R5 规划（D-011）；为 11 个纳入域交付范例 + 结构/行为验证，R5 验收前闭合 `I-PROTO-003`。goal-tree 与路线图已同步；`progress` 维持 `4/6`，`I-PROTO-003` / `I-PROTO-004` 状态不变。
- R5 阶段 1（2026-07-31）：GOAL-007 响应 A-001（independent pass，用户裁决「不需要自审，直接推进」）并完成契约发现登记——`I-007-001` 登记表落盘（attachments/I-007-001-registry.md v0.1.0），D-APP/D-PERM 复用产物与可执行验证命令核验入账，`I-007-001` → `verified`（登记层面）。阶段 2-4 未开始；`I-PROTO-003` 仍 open，验收前须以阶段 3 可执行证据闭合。`progress` 维持 `4/6`。
- R5 阶段 2 批次 2a（2026-07-31）：GOAL-007 按 D-004 实施批次 2a（D-DATA/D-TABLE 范例 + Go 列表/详情 API），`npm test` 114 项 / `go test` 21 项 / build / Edge 实测全绿，A-002 批次自审（self）pass，登记表升 v0.2.0。`I-PROTO-003` 仍 open；`progress` 维持 `4/6`。
- R5 阶段 2 批次 2b（2026-07-31）：GOAL-007 按 D-004 实施批次 2b（D-FORM 控件表面 + D-ACT 非批量动作 + Go PATCH/DELETE + `form-controls`/`list-edit-lifecycle` 范例页），`npm test` 138 项 / `go test` 18 顶层 / build 全绿，A-003 批次自审（self）pass，登记表升 v0.3.0。`I-PROTO-003` 仍 open（验收/关门门禁）；`I-PROTO-004`（non-blocking）在阶段 3 结构校验实现前决策；批次 2c（D-EXPR + Renderer 接线）与阶段 3/4 未开始。`progress` 维持 `4/6`。
- R5 阶段 2 批次 2c（2026-08-01）：GOAL-007 按 D-004 实施批次 2c（D-EXPR `reactions.ts` + D-COMP `render.ts`/`render.tsx` 最小 Renderer 接线，resolve R4 F-002 + `form-with-reactions` 范例页），`npm test` 166 项 / `npm run build` / `go test` / Edge 实测全绿，A-004 批次自审（self）pass，登记表升 v0.4.0。阶段 2 全部落地；阶段 3/4 未开始。
- R5 A-005 响应（2026-08-01）：GOAL-007 响应 A-005（independent, conditional）——用户裁决「不需要自审，直接推进」，F-001～F-004 按 `fixed` 合法闭合（action 表达式 fail-closed、RenderPage 接入 D-FORM 门禁、defaultValue 2.7+advanced 双门禁、Renderer whitelist 扩展至冻结 §5 全部 node type），F-005 本行投影同步；`npm test` 173 项 / build / `go test` / Edge 实测全绿；登记表升 v0.5.0（D-007）。`I-PROTO-003` 仍 open；`progress` 维持 `4/6`。
