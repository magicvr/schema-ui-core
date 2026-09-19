---
id: GOAL-005-r4-result-center-experience
title: R4 结果中心与体验收敛
status: active
parent: GOAL-001-batch-operations-and-job-center
created: 2026-09-19
updated: 2026-09-19
version: 0.2.0
progress: 3/4
plan_refs:
  - VP-038-batch-operations-and-job-center
primary_plan: VP-038-batch-operations-and-job-center
vision_ref: schema-ui-core-admin-foundation@0.4.0
---

# GOAL-005 · R4 结果中心与体验收敛

## 概述

承接 Root `GOAL-001` 的纲领阶段 **R4**：把 R2 的作业读面与 R3 的异步批量导出收敛成**可用的结果中心体验**——进行中、成功、失败、取消、过期与结果已过期六类呈现可用；导出类结果可下载；失败在合同允许时可重试；并可取消进行中的作业。同时收敛中英文、浅色/深色、加载/空态/错误态与既有产品约定。

R1～R3 已冻结并交付：契约归本地（方案 B，ADR-0022 冻结）、`jobs.read` 管理读面、`jobs.write` 提交写面、`jobs.batch-export` 作业与真实进度。R4 补上**写操作面**（取消/重试）与**呈现收敛**。

## 范围与非目标

### 本目标范围（承接 R3 移交项与 Root 路线图 R4）

- **写操作面**：管理作用域的 `POST /api/jobs/{id}/cancel` 与 `POST /api/jobs/{id}/retry`（`jobs.write` 门控）；与既有 actor 作用域 wallet 写路径**并行且不改动**。
- **结果中心呈现**：六类状态（进行中/成功/失败/取消/过期/结果已过期）可读；进度可观察；导出类结果**可下载**；失败可按合同重试；进行中可取消。
- **体验收敛**：中英文、浅色/深色、加载/空态/错误态、可访问性与既有列表/反馈约定一致。
- **前端交互级测试**：覆盖空选择禁用、只接受 202、**不调用 `reloadList()`**、轮询至终态、终态下载（R3 审计遗留的 recommended）。
- **`I-038-013`**：导出文件名/格式与两个导出入口的 UI 文案一致性。

### 明确非目标

- 改变 Job 六态合同；把同步 `batch-delete` 改成异步。
- 新增 Redis / 外部队列 / 多实例 / 专用搜索引擎（保持 gated）。
- 重开 VP-012/011/037/036；触碰任何 pinned 协议工件。
- 通用 BI/报表中心、跨模块数据仓库。

## 成功检查点

以下 4 个检查点构成 `progress: 3/4` 的派生来源。

- [x] **C1 写操作面**：管理作用域取消/重试路由 + 仓储方法（`jobs.write` 门控、fail-closed）；既有 actor 隔离写路径逐字不变。证据：`02-execution/E-001-r4-implementation.md` §1/§3；`internal/jobs/actions_test.go` 反向钉住 `RequestCancel`/`Retry` 对非本人仍 `ErrNotFound`。
- [x] **C2 结果中心呈现**：六类状态 + 进度 + 下载 + 重试 + 取消在 UI 可用；空态/加载态/错误态与既有约定一致。证据：`modules/jobs/schema/jobs.json`；`renderer/jobs-result-center.test.tsx`（9 例，渲染真实 `jobs.json`）。
- [x] **C3 体验与测试收敛**：中英文、浅色/深色、可访问性；前端交互级测试落地；`I-038-013` 关闭；全量回归绿。证据：i18n 双目录 1220 键对齐；主题 token 整页断言；`jobs-batch-export.test.tsx`（R3 遗留 4 例）；`job-result-download.test.ts`（`I-038-013` 单源守卫）；vitest 117 files / 1455 tests 全绿 + `go test ./...` 全绿。
- [ ] **C4 R4 审计与投影**：self + independent 审计落盘，开放 required = 0，Root R4 检查点可投影。

## 审计模式（P-002 实施前确定）

**`cross`（self + independent）**。理由：本阶段新增**管理作用域的写操作**（可取消/重试他人作业）并收敛权限可见面——属权限/数据高影响门禁。independent provider 按项目级决策 [independent-audit-execution.md](../../../../architecture/independent-audit-execution.md) = 本地 grok build（模型 grok 4.6 · 思考强度 high · `/audit`）。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-038-014 | required | **结果中心结构选型**：在既有 `jobs` 列表页上收敛（行操作 + 详情抽屉/内联）还是新增独立结果中心页？两者对能力声明、导航、i18n 与回归面的影响 | C2 | C2 前 | 读既有页面结构、行操作能力与导航约定；给出两方案影响面 | **verified** | — | **用户 P-004 裁决 = 方案 A**（既有 `jobs` 页原地收敛：行操作 + `recordView`；不新增页面/导航）。见 `01-decision/D-001` §0/§3 |
| I-038-015 | required | 管理作用域取消/重试的**合同口径**：可取消/可重试的状态集合、取消进行中作业的语义（`cancel_requested` → worker 察觉）、重试的次数上限与既有 `ErrNotCancellable`/`ErrNotRetryable` 复用 | C1 | C1 前 | 读 `internal/jobs/repository.go` 的 actor 作用域实现与 runner 取消路径 | **verified** | — | **用户 P-004 裁决 = 逐字镜像**（同状态集合/同错误码/不重置 attempt/`jobs.write`）。见 `01-decision/D-001` §0/§1；证据 `internal/jobs/actions_test.go` |
| I-038-016 | non-blocking | 结果过期后的呈现与「下载」入口的失效语义（410 `JOB_RESULT_EXPIRED` 在 UI 上的呈现） | C2 | R4 前 | 对照既有错误呈现约定 | **verified** | — | 由 `jobToMap` 派生 `downloadable`（仅 `succeeded`）关闭：过期态按钮禁用，410 仅在直接访问 URL 时兜底。见 `D-001` §4 |
| I-038-013 | non-blocking | 导出文件名/格式与两个导出入口的 UI 文案一致性 | C3 | R4（C3） | 对照同步导出与异步导出的命名/格式 | **verified** | — | 唯一实现 `apps/web/src/lib/job-result-download.ts`（CSV 用服务端 `fileName`，否则 JSON），R3 组件与渲染器分支共用；守卫 `job-result-download.test.ts`（4 例） |

R3 已关闭的 `I-038-011`/`012` 不再重复登记；`I-038-013` 已由本目标 C3 关闭为 `verified`；`I-038-006` 保持 `deferred · non-blocking`；`I-038-010`（导航与 i18n 键位）已由本目标 C3 承接（双目录键集合对齐 1220/1220，`schema-keys.structural` + 能力声明 guard 全绿）。

## 审计意见状态

尚无审计条目（C4 待跑：self + grok build independent，模式 `cross`）。

## 父目标

- `GOAL-001-batch-operations-and-job-center`（Root 的 R4 纲领阶段子目标；Root 现为 `active · 3/5`）。

## 台账布局

本目标使用平铺 ledger：`01-decision/`、`02-execution/`、`03-audit/`，并保留 `attachments/` 作为证据载体。

## 备注

- `progress: 0/4` 只由上方 4 个显式检查点派生；不放行阶段、不关闭 finding、不覆盖 status。
- 本目标是**最后一个实现阶段**；R5 只做证据、审计与关门。
- `I-038-014` 是结构选型，按用户指令属**关键决策**，须询问用户（不得静默）。
