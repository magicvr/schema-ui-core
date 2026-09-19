---
id: GOAL-004-r3-async-batch-operation
title: R3 批量操作异步承接
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

# GOAL-004 · R3 批量操作异步承接

## 概述

承接 Root `GOAL-001` 的纲领阶段 **R3**：让**至少一条真实批量操作**以异步 Job 承接——用户在列表页选中若干行、提交后拿到 `202 + jobId`，可观察 `queued → running → 终态` 与**进度**，并在终态读取结果。同时既有 ADR-0022 同步 `batch-delete` 的已交付语义与回归**不得退化**。

R1 已冻结首波分母 = **仅「新建批量导出所选」1 条**（`GOAL-002/01-decision/D-001-…` §3），契约形态 = **方案 B**（本地模块自有异步契约；ADR-0022 逐字冻结；协议 pin 零改动，§1）。R2 已交付 `admin.jobs` 管理读面与共享 `jobs.ResultURL`（`GOAL-003`）。本目标在此之上补上**写面**。

## 范围与非目标

### 本目标范围（承接 R1 `D-001` §5 的 T-5 / T-6）

- **T-5** 异步批量导出端点（`POST` → **202 + job 投影**）+ 新 Job kind + handler；**细粒度 `reporter.Progress`**（R1 `D-001` §1.2 K-6：既有先例只报硬编码 10→100，进度须由新 handler 自己实现）。
- **T-5** 结果经既有 `GET /api/jobs/{id}/result` 读取（R2 已交付，须复用而非另建）。
- **T-6** 前端触发机制（O-1 冻结 = **自定义组件**，形态照 `components/monitoring-auto-refresh.tsx`）+ capability 声明口径（O-2 冻结 = **不声明 `actions.batch.request`**）。
- **T-2 写权限** `jobs.write`（`PolicyAdmin`）声明与接线（R2 未声明，按 `D-001` §4 归本阶段）。
- 同步 `batch-delete` 全量回归不退化；既有 Job 六态合同不变。

### 明确非目标

- 结果中心体验收敛（列表/详情/进度/终态/过期/重试/取消/下载的完整 UI）——属 **R4**。
- 把既有同步 `batch-delete` 改成异步（VP-038 显式非目标）；改任何 pinned 协议工件。
- 导入/purge-all/操作日志导出改异步（R1 `D-001` §3.4 已排除）。
- 新增 Redis / 外部队列 / 多实例 / 专用搜索引擎（保持 gated）。

## 成功检查点

以下 4 个检查点构成 `progress: 4/4` 的派生来源。

- [x] **C1 异步写面与进度**：`jobs.write` 权限 + 批量导出端点（202 + jobId）+ Job kind + handler 落地；**进度可观察**（非硬编码终值）；结果经既有 `/api/jobs/{id}/result` 可读。
- [x] **C2 前端触发**：列表页接入自定义组件触发提交（O-1）并轮询至终态；页面**未**声明 `actions.batch.request`（O-2）；中英文与既有反馈约定一致。
- [x] **C3 同步路径回归**：既有同步 `batch-delete`（含原子性与协议 fixture）与 Job 六态合同**逐字不退化**；全量 Go/Web 回归绿。
- [x] **C4 R3 审计与投影**：self + independent 审计落盘，开放 required = 0，Root R3 检查点可投影。

## 审计模式（P-002 实施前确定）

**`cross`（self + independent）**。理由：本阶段新增**写面**（可触发他人数据导出）并新增 Job kind，涉及权限边界、数据外带面与既有同步批量路径的不回退保证——属权限/数据高影响门禁。independent provider 按项目级决策 [independent-audit-execution.md](../../../../architecture/independent-audit-execution.md) = 本地 grok build（模型 grok 4.6 · 思考强度 high · `/audit`）。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-038-011 | required | 导出**数据面**分母与权限口径：批量导出所选支持哪些资源；是否沿用 `data.export` 还是 `jobs.write`；导出内容是否含敏感字段（对齐既有 `GET /api/export/{resource}` 的列集与转义） | C1 | C1 前 | 读既有导出实现与列集、权限门；给出导出资源与列的白名单 | **verified**（2026-09-19：users/roles 同分母；列集与转义复用同步导出；**双重门禁 jobs.write + data.export**） | — | `attachments/R3-recon-frontend-batch-trigger.md`；`01-decision/D-001-…` §1 |
| I-038-012 | required | 前端能否从自定义组件读到**当前表格选择集**；若不能，触发机制如何取得选中行 | C2 | C2 前 | 读 `render.tsx` 的 selection 状态与 `custom-components.ts` 的 context 契约 | **verified**（2026-09-19：可，经 `useSchemaCrud().selection(tableId)`；位置与 capability 口径见 D-001 §2） | — | `attachments/R3-recon-frontend-batch-trigger.md`；`01-decision/D-001-…` §2 |
| I-038-013 | non-blocking | 导出结果的文件名/格式与既有 CSV 约定的一致性细节 | C2/R4 | R4 前 | 对照既有导出文件名与 BOM/转义 | open（文件名 `<resource>-selection.csv`；BOM/RFC4180 已复用，UI 文案差异待 R4 收敛） | — | `01-decision/D-001-…` §1.1 |

R1 已关闭的 `I-038-001`～`003` 与 R2 已关闭的 `I-038-007`～`009` 不再重复登记；`I-038-006` 保持 `deferred · non-blocking`；`I-038-010`（导航与 i18n 键位）最晚 R4。

## 父目标

- `GOAL-001-batch-operations-and-job-center`（Root 的 R3 纲领阶段子目标；Root 现为 `active · 2/5`）。

## 台账布局

本目标使用平铺 ledger：`01-decision/`、`02-execution/`、`03-audit/`，并保留 `attachments/` 作为证据载体。

## 备注

- `progress: 0/4` 只由上方 4 个显式检查点派生；不放行阶段、不关闭 finding、不覆盖 status。
- R1/R2 的冻结结论是本目标的**约束输入**：契约归属本地（不碰 `batchMapping`）、首波只做 1 条、O-1/O-2 口径已定、`jobs.read` 已交付。
