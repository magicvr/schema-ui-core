---
id: GOAL-002-r1-scope-semantics-freeze
title: R1 列表分母与工作流语义冻结
status: done
parent: GOAL-001-admin-workflow-continuity
created: 2026-09-16
updated: 2026-09-17
version: 0.4.0
progress: 3/3
plan_refs:
  - VP-037-admin-workflow-continuity
primary_plan: VP-037-admin-workflow-continuity
vision_ref: schema-ui-core-admin-foundation@0.4.0
---

# GOAL-002 · R1 列表分母与工作流语义冻结

## 概述

为 Root 的 R1 阶段建立可复核的实现入口：从当前 manifest、Schema 页面和 Renderer 事实中冻结列表页分母、Saved View 状态边界、dirty-state 场景及反馈/权限映射。该目标只承载信息收集与方案冻结，不提前实现 R2～R4。

## 范围与边界

- 覆盖当前 mvp、admin、demo、custom profile 的页面注册、可发现页面、隐藏但可由动作抵达的页面，以及全部已注册 table/form 节点。
- 记录真实的查询状态、筛选/排序/列配置能力、表单来源和权限/错误合同。
- 冻结用户隔离、持久化、序列化、失效和 dirty-state 取舍后，才允许 R2～R4 进入方案/实施。
- 不新增业务域、全文搜索、批量结果中心、跨用户共享或协作能力。

## 成功检查点

- [x] C1：profile/page/table/form 分母和权限覆盖形成机器可核对矩阵（24 个列表表面、58 个表单节点）。
- [x] C2：Saved View、dirty-state、Toast/API 错误/恢复语义的取舍形成已确认决策。
- [x] C3：I-037-001～004 关闭，R1 self `A-001` 与 independent `A-002` 均 `pass`；A-003 响应后无开放 required / 必改 finding，Root R1 检查点可投影为完成。

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-037-001 | required | 当前列表页、状态字段、Profile/权限覆盖的精确分母 | R1 / R2 / R5 | R1 | 扫描 manifest、Schema 与 profile matrix | verified | 2026-09-17 已由矩阵核对 | `attachments/r1-denominator-matrix.json`、`attachments/r1-form-matrix.json` |
| I-037-002 | required | Saved View 所有权、持久化/序列化、权限变化后的失效语义 | R1 / R2 | R1 | 用户确认 localStorage 方案 A；矩阵与 D-003 冻结 allowlist、失效和异常边界 | verified | 2026-09-17；R2 仍需留存实现/回归证据 | `01-decision/D-003-saved-view-localstorage-accepted.md`、`attachments/r1-state-feedback-matrix.md` |
| I-037-003 | required | dirty-state 在内部路由、浏览器离开、提交、重置、取消中的统一语义 | R1 / R3 | R1 | 盘点 App 导航、FormInner、modal 生命周期和浏览器事件；D-004 冻结 | verified | 2026-09-17；R3 仍需留存浏览器/自动化证据 | `01-decision/D-004-workflow-semantics-frozen.md`、`attachments/r1-state-feedback-matrix.md` |
| I-037-004 | required | Toast、API 错误、重试、维护/不可用反馈的分类与可访问呈现 | R1 / R4 | R1 | 对照 `readResourceApiError`、FeedbackRegion、DataTable 和 Host failure；D-005 冻结 | verified | 2026-09-17；R4 仍需留存跨页面回归证据 | `01-decision/D-005-feedback-recovery-semantics-frozen.md`、`attachments/r1-state-feedback-matrix.md` |

## 父目标

- `GOAL-001-admin-workflow-continuity`（R1 阶段子目标；本目标已 `done · 3/3`，Root 仍为 `active`，并按 R1→R2/R3/R4→R5 路线推进）。

## 台账布局

本目标使用平铺 ledger：`01-decision/`、`02-execution/`、`03-audit/`，并保留 `attachments/` 作为证据载体。
