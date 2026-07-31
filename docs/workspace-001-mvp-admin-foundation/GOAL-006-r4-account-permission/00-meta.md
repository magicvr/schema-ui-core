---
id: GOAL-006-r4-account-permission
title: R4 · 核心账号与权限
status: done
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-07-31
version: 0.4.0
---

# GOAL-006 · R4 · 核心账号与权限

## 概述

在 R2 冻结的 MVP 覆盖基线与 R3 Admin 外壳/导航之上，交付核心账号与权限能力：账号权限最小 API 的设计、`D-PERM` 映射冻结，以及前后端可核对的鉴权链路。R4 方案冻结前验证本目标信息项 `I-006-001`；R4 **实施**前闭合父目标 `I-PROTO-002`。

范围依据：Root 将 R4 定义为「核心账号与权限」；协议资料将 `D-PERM` 映射为账号/权限核心能力，并配套 `permissions-inheritance` 等 behavioral fixture 作为结构/行为验证来源。固定协议版本沿用 `schema-ui-docs` artifact `2.7.0`、source commit `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b`。

## 范围边界

### 纳入

- 账号权限最小 API 的设计决策与 `D-PERM` 结构/行为映射。
- 前后端账号/权限链路的实现与可核对验证（R4 实施阶段）。
- 与 R3 Admin 外壳导航上下文的衔接：真实身份/权限来源替换 R3 的默认空 context。

### 排除

- R5 的协议 Renderer 全量、纳入域范例页与 `mock-app` 业务演示；父目标 `I-PROTO-003` 为 R5 验收/关门 required 门禁。
- 完整权限继承产品化、完整主题/设计系统与「完整协议支持」主张；R3 的 `I-005-*` 验证边界不自动扩展到权限语义。
- 任何未在 `I-PROTO-001` v0.1.3 覆盖表内纳入的协议能力。

## 高层路线图（P-001）

1. **契约发现与信息就绪**：**完成**；验证 `I-006-001`（账号权限最小 API 与 `D-PERM` 映射），并在方案冻结时闭合父目标 `I-PROTO-002`（D-004，证据见 `attachments/dperm/`）。
2. **方案冻结**：**完成**；D-004 冻结最小 API、权限模型映射与前后端集成边界（账号会话最小闭环；`$context` 只读快照；Go 独立鉴权）。
3. **R4 实施**：**完成**；Go 会话与 `/api/accounts/me`、Go 独立鉴权、Web `$context` 挂载与 D-PERM 求值引擎已按 02-execution 时间线落地；`go test`/`go build`、web 94 项测试、`npm run build`、HTTP 运行时与代理联调证据已入账（2026-07-31，A-001 self pass + A-002 independent pass）。
4. **验证与关门**：**完成**；结构/行为/运行时证据已齐（见步骤 3），A-004 关门自审（self）与 A-005 独立关门复审（independent）均 **pass**、开放 required=0，经用户确认后 `done` 评估完成。

## 成功标准

- [x] R4 必需信息项 `I-006-001` 已由证据验证（D-004 + 固定资料 SHA-256）；未知项没有被默认为已知。
- [x] 账号权限最小 API 与 `D-PERM` 映射方案已冻结（D-004）；父目标 `I-PROTO-002` 在 R4 实施前合法闭合（`verified`）。
- [x] 前后端账号/权限链路具备可核对的实现与验证路径（结构/行为/运行时证据），R3 导航 context 得到真实身份/权限来源衔接。
- [x] R5 Renderer/范例与完整协议支持保持边界外；关门前无开放 required finding。

## 关门结论

- 成功标准四项全部达成；信息门禁 closed（`I-006-001` verified、父目标 `I-PROTO-002` verified）；意见台账开放 required=0。
- A-004 关门自审（self）与 A-005 独立关门复审（independent）均 **pass**，独立意见同意 A-004 的关门主张；A-006 经用户 `/govern` 授权执行关门。
- F-002～F-004（Renderer 接线 / token 会话 / 双端一致性 oracle）为 recommended 跟踪项，随 R5 / 生产化 / `I-PROTO-004` 解决，**不**阻断 R4 关门。
- 本目标 `done` 仅覆盖 R4 冻结子集与边界，不推导完整协议 conformance，也不改变父目标 R5/R6 门禁。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-006-001 | required | 账号权限最小 API 与 `D-PERM` 映射是否完整？ | 方案冻结 / 实施 | 方案冻结前 | 对照 `D-PERM` 与 `permissions-inheritance` 等固定 fixture，设计最小 API 并映射 React/Go 实现路径 | **verified** | 不延期；2026-07-31 方案冻结时验证 | D-004 冻结映射；`attachments/dperm/` 固定资料（cases.json 17 例、node.schema、ADR-0023、node-protocol §3.9/3.9.1、component-registry、renderer-spec §7.1）SHA-256 核验 |

父目标 `I-PROTO-002`（required，R4 **实施**门禁）、`I-PROTO-003`（required，R5 验收/关门门禁）由 Root 维护，本目标不修改其状态。`I-PROTO-004`（vendor vs pin，non-blocking）仍 open；关闭该项时须补 schema-conformance 等价性校验或显式记录等价范围（GOAL-005 A-007 F-002 跟进项）。

## 父目标

- [GOAL-001-mvp-admin-foundation](../GOAL-001-mvp-admin-foundation/00-meta.md)

## 备注

- 立项日期：2026-07-31，承接 R3 关门（GOAL-005 `done`）。
- Root 纲领进度仍为 `3/6`；本目标方案冻结不抬升 progress，也不放行 R4 实施。
- 2026-07-31：D-004 方案冻结，`I-006-001` 与父目标 `I-PROTO-002` 均 `verified`；实施仍待用户指令并须记实施事实。
- 2026-07-31：R4 实施完成（02-execution 时间线）；A-001（self）与 A-002（independent）实施阶段均 pass、开放 required=0；R4 关门自审与 `done` 评估待用户确认后进行。
- 2026-07-31：**R4 关门完成**；A-004 关门自审（self）与 A-005 独立关门复审（independent）均 pass、开放 required=0；经用户 `/govern` 授权，本目标标 `done`，Root 纲领 R4 检查点完成（`progress` → 4/6），goal-tree 同步。
- R4 实施阶段：实施计划、`$context` 会话方案（静态/注入 vs token）、前后端实现、fixture 对照测试与运行时证据均已在 02-execution 记录；不把「方案已冻结」写成「已实现」。
