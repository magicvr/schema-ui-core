---
id: GOAL-003-r2-generic-job-read-surface
title: R2 通用作业读面
status: done
parent: GOAL-001-batch-operations-and-job-center
created: 2026-09-19
updated: 2026-09-19
version: 0.2.0
progress: 4/4
plan_refs:
  - VP-038-batch-operations-and-job-center
primary_plan: VP-038-batch-operations-and-job-center
vision_ref: schema-ui-core-admin-foundation@0.4.0
---

# GOAL-003 · R2 通用作业读面

## 概述

承接 Root `GOAL-001` 的纲领阶段 **R2**：在 R1 已冻结的分母与契约之上，交付 `admin.jobs` 模块与**管理作用域**的作业读面——已注册 Job 种类可按权限列出、查看详情、读取结果，且越权/不存在时 fail-closed。**不改变 Job 六态合同**，**不放宽**既有 actor 隔离语义。

R2 是本 VP 的第一段实现代码；`admin.jobs` 模块此前不存在（现有 `apps/api/modules/jobs/` 只有迁移用的 `core.jobs`）。

## 范围与非目标

### 本目标范围（承接 R1 `D-001` §5 移交清单）

- **T-1** `admin.jobs` 模块建立：`kernel.Provider`（Descriptor / Register / CompiledPersistence）+ 进 admin 默认集（`I-038-004` 用户裁决 = Profile 内容扩展）。
- **T-2** 新权限 `jobs.read`（+ 写权限键，本目标冻结）声明与接线；`composition_test.go:529` 的 `users`/`navigation` 计数同步。
- **T-3** 新增 `internal/jobs` repository 查询方法：按 kind / status / actor / 时间范围过滤 + 分页（page/pageSize + offset）+ total COUNT + 排序。
- **T-4 / O-3** 管理列表索引决策（是否需要新迁移版本 + 同步两处冻结断言）。
- **T-5（读面部分）** 作业列表 / 详情 / 结果读取 API + 权限与作用域过滤 + fail-closed。
- **T-6 / O-1 / O-2** 前端触发机制与 capability 声明口径 —— **本目标只冻结方案，不实现**（实现属 R3/R4）。
- **T-7** 结果中心页面 —— **不在本目标**（属 R4）。
- 模块级验证与回归（provider/契约测试、依赖冲突 fail-closed、Profile 启动路径、页面/权限/导航可观察）。

### 明确非目标

- 批量操作异步承接（R3）；结果中心 UI 与体验收敛（R4）。
- 改变 Job 六态合同、`payload`/`result` 的既有表示口径（VP-012 冻结「GET Job 不内嵌 payload/result」）。
- 放宽 `GetForActor` 的 actor 隔离语义或其冻结测试。
- 把同步 `batch-delete` 改成异步；触碰任何 pinned 协议工件。
- Redis / 外部队列 / 多实例 / 专用搜索引擎（保持 gated）。

## 成功检查点

以下 4 个检查点构成 `progress: 4/4` 的派生来源。

- [x] **C1 模块与权限接线**：`admin.jobs` 模块建立并进 admin 默认集；`jobs.read` 声明、接线、入 catalog；模块/组合测试与计数断言全绿。
- [x] **C2 查询与索引**：`internal/jobs` 新增查询方法（过滤/分页/total/排序）落地并有 repository 测试；O-3 索引决策冻结并实施（含两处冻结断言的同步）。
- [x] **C3 读面 API 与作用域**：列表/详情/结果读取 API 落地；管理作用域由 `jobs.read` 门控；越权/不存在 fail-closed；**既有 actor 隔离测试语义不变**。
- [x] **C4 R2 审计与投影**：self + independent 审计落盘，开放 required = 0，Root R2 检查点可投影。

## 审计模式（P-002 实施前确定）

**`cross`（self + independent）**。理由：R2 引入**跨 actor 的管理读面**（新权限 + 新可见性边界），并可能新增迁移索引（触碰两处冻结断言）——属「权限/数据边界 + 迁移」高影响门禁。independent provider 按项目级决策 [independent-audit-execution.md](../../../../architecture/independent-audit-execution.md) = 本地 grok build（模型 grok 4.6 · 思考强度 high · `/audit`）。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-038-007 | required | 管理列表索引决策（O-3）：既有三索引是否足够，是否新增迁移版本；新增则两处冻结断言如何同步 | C2 | C2 前 | 迁移机制 + 索引矩阵静态判定 + EXPLAIN 实测 | **verified**（2026-09-19 用户裁决：新增 v72，并细化为 `(created_at DESC, id DESC)`） | — | `attachments/R2-recon-index-and-query-shape.md`；`01-decision/D-001-…` §1/§1.2a |
| I-038-008 | required | 结果 URL 泛化口径：现行 `walletJobToMap` 硬编码 `/api/wallet/jobs/{id}/result`（wallet 权限门下），通用读面如何表达结果地址 | C3 | C3 前 | 读 `walletJobToMap` 与 wallet 结果路由；给出通用结果地址方案 | **verified**（2026-09-19 用户裁决：泛化 + 共享 helper，登记为字节等价重构；C4 审计已复核等价性） | — | `01-decision/D-001-…` §2；`03-audit/A-002-…` |
| I-038-009 | required | 前端触发机制与 capability 声明口径（O-1 / O-2）：自定义组件 vs `props` 本地扩展键；新页面是否声明 `actions.batch.request`（guard marker = `/batch-delete/`） | C4（方案冻结）、R3/R4 实施 | R2 方案冻结前 | 对照 `monitoring-auto-refresh.tsx` 先例与 `capability-declaration.guard.test.ts:38-58,110-114` | **verified**（2026-09-19 冻结：O-1 自定义组件；O-2 不声明 `actions.batch.request`——`jobs.json` 已遵守且 guard 通过） | — | `01-decision/D-001-…` §3 |
| I-038-010 | non-blocking | `admin.jobs` 的导航分组与 i18n 键位（VP-034 分组约定） | R4 体验收敛 | R4 前 | 读 `nav-groups-r4.test.ts` 与既有模块 fragment | open | — | R1 侦察 §6 接线清单 |

R1 已关闭的 `I-038-001`～`003` 与 `I-038-004`/`005` 不再重复登记；`I-038-006` 保持 `deferred · non-blocking`（历史作业保留与清理）。

## 父目标

- `GOAL-001-batch-operations-and-job-center`（Root 的 R2 纲领阶段子目标；Root 现为 `active · 1/5`）。

## 台账布局

本目标使用平铺 ledger：`01-decision/`、`02-execution/`、`03-audit/`，并保留 `attachments/` 作为证据载体。

## 备注

- `progress` 只由上方 4 个显式检查点派生（现为 `4/4`）；不放行阶段、不关闭 finding、不覆盖 status。
- **2026-09-19 投影修正（用户授权 · `VRev-100` `V-F130`）**：frontmatter `progress` 由 `3/4` 更正为 `4/4`，与 `workspace-038/goal-tree.md` 的 `done · 4/4` 及上方四个已勾选检查点一致（`A-003` 响应第 1 条的本意即为 `3/4 → 4/4`）；同时同步本文件 L47 派生说明与上一条旧样板文本。`status` / 结论 / 三个 ledger 目录 / `goal-tree` 均未改动。
- R1 冻结结论（`GOAL-002/01-decision/D-001-…`）是本目标的**约束输入**：C2 = 方案 B、C3 首波 = 新建批量导出所选、C1 = 管理作用域 + `jobs.read`。
- 未定项 O-1～O-3 必须在 R2 方案中冻结；其中 O-3 已登记为 `I-038-007`（required）。
