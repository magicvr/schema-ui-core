---
id: GOAL-002-r1-denominator-and-contract-freeze
title: R1 分母与契约冻结
status: done
parent: GOAL-001-batch-operations-and-job-center
created: 2026-09-19
updated: 2026-09-19
version: 0.3.0
progress: 4/4
plan_refs:
  - VP-038-batch-operations-and-job-center
primary_plan: VP-038-batch-operations-and-job-center
vision_ref: schema-ui-core-admin-foundation@0.4.0
---

# GOAL-002 · R1 分母与契约冻结

## 概述

承接 Root `GOAL-001` 的纲领阶段 **R1**：把「纳入首波的 Job 种类与可见作用域」「批量异步契约形态与协议面影响」「首波批量操作分母」「权限/Profile 边界与排除项」冻结成**可机器核对**的矩阵与决策，使 R2（通用作业读面）与 R3（批量操作异步承接）可以在无未知项的前提下开工。

本目标是**信息收集 + 方案冻结**阶段，**不实现** R2/R3 的任何产品代码（`admin.jobs` 模块建立、Job 列表 API、批量异步端点均不在本目标范围）。

## 范围与非目标

### 本目标范围

- `I-038-001`：Job 种类×作用域矩阵（现存及首波纳入的种类、可见作用域口径、与 `wallet.reconcile` 既有 actor 作用域的兼容关系、读面所需 repository 查询的缺口清单）。
- `I-038-002`：批量异步契约形态判定（本地扩展 vs 上游协议变更）、同步 `batch-delete` 兼容口径、pinned `schema-ui-docs@v2.9.0` 协议面影响结论。
- `I-038-003`：首波批量/长操作分母（逐项给出「保持同步 / 改异步 / 不进首波」），并承接 `V-F126`。
- 权限/Profile 边界与排除项冻结（含 `admin.jobs` 进 admin 默认集的边界确认，承接 `I-038-004` 已裁决结论）。
- R1 阶段审计（self + 按风险判定的 independent）与 R1 检查点投影到 Root。

### 明确非目标

- 实现 `admin.jobs` 模块、Job 列表/详情/结果 API、批量异步端点、结果中心 UI（属 R2～R4）。
- 改变 Job 六态合同；把同步 `batch-delete` 改成 breaking 异步语义。
- 新增 capability 字符串或修改任何 pinned 协议工件（`docs/schemas/**`、`apps/web/src/protocol/upstream/*.cases.json`）。
- 实体全文检索、组织/数据权限、新业务域；Redis / 外部队列 / 多实例 / 专用搜索引擎（保持 gated）。

## 成功检查点

以下 4 个检查点构成 `progress: 4/4` 的派生来源。

- [x] **C1 分母与作用域矩阵**：`I-038-001` 关闭——Job 种类×作用域矩阵（种类、注册点、可见作用域、读面缺口、索引、后台周期任务口径、运行时门控）落盘且可机器核对（`attachments/r1-job-kind-scope-matrix.md`）。
- [x] **C2 契约形态冻结**：`I-038-002` 关闭——批量异步契约形态经用户 P-004 裁决（**方案 B**）并落盘（含未选方案、同步 `batch-delete` 兼容口径、协议 pin 零影响结论；`01-decision/D-001-…` §1）。
- [x] **C3 首波分母冻结**：`I-038-003` 关闭——首波 = **仅「新建批量导出所选」**（1 条），含保持同步 S-1～S-5、排除 X-1～X-11、Breaking 标记与回归面（`attachments/r1-first-wave-denominator-matrix.md`）；承接 `V-F126`。
- [x] **C4 R1 审计与投影**：self 审计 + 按风险判定的 independent 审计落盘（`03-audit/A-NNN`），开放 required = 0，Root R1 检查点可投影为完成。

## 审计模式（P-002 实施前确定）

**`cross`（self + independent）**。理由：R1 冻结的批量异步契约直接影响权限与作用域语义（管理作用域 Job 读面 = 跨 actor 可见性），且契约形态决定 R2/R3 是否触碰 pinned 协议面；属「协议/跨边界」高影响门禁。independent provider 按项目级决策 [independent-audit-execution.md](../../../../architecture/independent-audit-execution.md) = 本地 grok build（模型 grok 4.6 · 思考强度 high · `/audit`）。

## 信息就绪与未知项（P-005）

本目标承接 Root 的三项 required 门禁；侦察证据位于 Root 目标 `../GOAL-001-batch-operations-and-job-center/attachments/R1-recon-*.md` 与 Root `../GOAL-001-batch-operations-and-job-center/02-execution/E-002-r1-recon.md`；本目标自有的两份冻结矩阵在 `attachments/r1-*-matrix.md`。

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-038-001 | required | Job 种类、可见作用域（管理 vs actor）与 `wallet.reconcile` 既有 actor 作用域的兼容关系；读面所需 repository 查询缺口 | R1 C1、R2 读面 | R1 | 侦察报告 + 本目标 C1 矩阵 | **verified**（2026-09-19） | — | Root `attachments/R1-recon-I-038-001-job-kinds-and-scopes.md`；`attachments/r1-job-kind-scope-matrix.md`；`01-decision/D-001-…` §2 |
| I-038-002 | required | 批量异步契约是本地扩展还是上游协议变更；同步 `batch-delete` 是否保持；是否触碰 pinned `v2.9.0` 协议面 | R1 C2、R3 实施 | R1 | 侦察报告 + 用户 P-004 裁决 + C2 决策落盘 | **verified**（2026-09-19 用户裁决方案 B） | — | Root `attachments/R1-recon-I-038-002-batch-async-contract.md`；`01-decision/D-001-…` §1 |
| I-038-003 | required | 哪些现有批量/长操作进入异步首波，哪些保持同步 | R1 C3、R3 实施 | R1 | 侦察报告 + 用户 P-004 裁决 + C3 矩阵落盘 | **verified**（2026-09-19 用户裁决首波 = 新建批量导出所选） | — | Root `attachments/R1-recon-I-038-003-batch-operation-inventory.md`；`attachments/r1-first-wave-denominator-matrix.md`；`01-decision/D-001-…` §3；承接 `V-F126` |

## 父目标

- `GOAL-001-batch-operations-and-job-center`（Root 的 R1 纲领阶段子目标；Root 保持 `active`，按 R1→R2→R3→R4→R5 推进）。

## 台账布局

本目标使用平铺 ledger：`01-decision/`、`02-execution/`、`03-audit/`，并保留 `attachments/` 作为证据载体。

## 备注

- `progress: 4/4` 只由上方 4 个显式检查点派生；不放行阶段、不关闭 finding、不覆盖 status。C4 已完成（开放 required = 0），Root R1 检查点可投影。
- 侦察报告为 `status: draft` 的只读证据，**不是**决策；C1～C3 的冻结结论以本目标 `01-decision/` 为准。
- 本目标关门后 Root R1 检查点方可投影为完成；R2/R3 子目标在 R1 冻结后按 P-001 立项。
