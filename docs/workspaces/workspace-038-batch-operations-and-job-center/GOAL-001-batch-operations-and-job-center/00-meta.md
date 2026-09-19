---
id: GOAL-001-batch-operations-and-job-center
title: Admin 批量操作与异步结果中心交付
status: active
parent: null
created: 2026-09-19
updated: 2026-09-19
version: 0.2.0
progress: 1/5
plan_refs:
  - VP-038-batch-operations-and-job-center
primary_plan: VP-038-batch-operations-and-job-center
vision_ref: schema-ui-core-admin-foundation@0.4.0
---

# GOAL-001 · Admin 批量操作与异步结果中心交付

## 概述

在 VP-012 已交付的 profile 无关 Job 六态运行时、VP-011 已交付的批量/导出导入面、VP-037 已交付的统一反馈与列表基线之上，交付 VP-038 的批量操作与异步结果中心：让已注册 Job 种类与纳入首波的批量操作在 Admin 中可观察、可操作、可追溯。

Root 只承接 VP-038 的实现层路线图（R1→R5），不把实体全文检索、组织/数据权限、新业务域或架构 gated 项（Redis / 外部队列 / 多实例 / 专用搜索引擎）写入本目标，也不重开 VP-012/011/037/036。

工作区与 Root 已建立（2026-09-19）；纲领阶段 **R1 已由子目标 `GOAL-002-r1-denominator-and-contract-freeze` 交付并关门**（`done · 4/4`），`progress: 1/5` 是显式检查点的派生展示。

## 子目标

| id | 纲领阶段 | status | progress |
|----|---------|--------|----------|
| GOAL-002-r1-denominator-and-contract-freeze | R1 分母与契约冻结 | **done** | 4/4 |

R2～R5 子目标在 R1 冻结后按 P-001 逐阶段立项。R1 已于 2026-09-19 关门（审计模式 `cross`：self `A-001` `pass` + grok build 4.6 high independent `A-002` `conditional` → `A-003` 响应 required 全 `fixed`，开放 required = 0）。

## 愿景对齐

- 工作区：`workspace-038-batch-operations-and-job-center`
- Charter：`schema-ui-core-admin-foundation@0.4.0`
- VP：`VP-038-batch-operations-and-job-center`（`active` v0.2.0）
- `plan_refs` / `primary_plan`：均为 `VP-038-batch-operations-and-job-center`
- `serves_summary`：把已交付的 Job 六态运行时与批量请求能力收敛为可观察、可操作的 Admin 产品面（作业可见性 + 至少一条异步批量承接 + 结果中心），保持既有权限/Profile 语义与同步 `batch-delete` 合同不回退；不新增业务域或解除 Redis/MQ/多实例/搜索 gated 条件。

## 范围与非目标

### 本目标范围

- 已注册 Job 种类的列表、详情、六态、进度、attempt、错误码/消息、结果读取与过期语义的对外契约。
- 作业可见作用域（管理作用域 vs actor 作用域）与按权限/数据范围的 fail-closed 过滤。
- 至少一条真实批量操作以异步 Job 承接（202 + jobId → 进度 → 结果）。
- 结果中心体验：进行中/成功/失败/取消/过期/结果已过期呈现、导出类结果下载、失败重试与取消。
- 中英文、浅色/深色、加载/空态/错误态与既有产品约定一致性；Profile/权限边界回归。
- `admin.jobs` 模块建立与 admin 默认集内容扩展（`I-038-004` 用户 P-004 裁决）。

### 明确非目标

- 实体全文索引、专用搜索引擎 `RT-X01` / `RT-X02`、跨进程索引。
- Redis、外部队列 / broker（`RT-Q02`）、多实例（A3）、第二持久化栈。
- 组织/部门/岗位、数据权限 `org`、SSO/多租户、新业务域。
- 重开 VP-012/011/037/036；改变 Job 六态合同本身；把同步 `batch-delete` 改成 breaking 异步语义。
- 通用 BI/报表中心、跨模块数据仓库、批量审批工作流。

## 成功标准与纲领路线图

以下 5 个检查点构成 Root 的派生 progress 来源；纲领阶段按 R1 → R2 → R3 → R4 → R5 串行推进。

- [x] **R1 分母与契约冻结**：Job 种类×作用域矩阵、批量异步契约（含同步 `batch-delete` 兼容口径与协议面影响判定）、首波批量操作分母（保持同步 / 改异步 / 不进首波）、权限/Profile 边界与排除项冻结；`I-038-001`～`003` 关闭。→ 由 `GOAL-002-r1-denominator-and-contract-freeze` 交付（**`done · 4/4`**，2026-09-19）。
- [ ] **R2 通用作业读面**：Job 列表/详情/结果读取 API + 权限与作用域过滤 + fail-closed；不改变 Job 六态合同。
- [ ] **R3 批量操作异步承接**：至少一条真实批量操作走 Job（202 + jobId + 进度 + 结果）；同步 `batch-delete` 既有语义与回归不退化。
- [ ] **R4 结果中心与体验收敛**：列表/详情/进度/终态/过期/重试/取消/下载；中英文、浅色深色、加载空态错误态、可访问。
- [ ] **R5 证据与关门**：退出矩阵、浏览器/自动化回归、独立意见、残余登记、组合投影同步与用户确认；开放 required = 0。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-038-001 | required | 现存及首波纳入的 Job 种类、可见作用域（全量/本 actor/按数据范围）与 `wallet.reconcile` 既有 actor 作用域的兼容关系是什么？ | R1 范围冻结、R2 读面 | R1 | 扫描 `internal/jobs` 与各模块 `Register*JobKind` 消费点，形成种类×作用域矩阵；确认是否需要新增 repository 列表查询 | open | — | 待确认 |
| I-038-002 | required | 批量异步契约是扩展 ADR-0022 还是另立独立契约？同步 `batch-delete` 是否保持？是否触碰 pinned `schema-ui-docs@v2.9.0` 协议面？ | R1 方案冻结、R3 实施 | R1 | 对照 `apps/web/src/protocol`（`actions.batch.request` / `batchMapping`）与上游 v2.9.0 契约；判定「本地扩展」vs「上游协议变更」，给出兼容与回归口径 | open | — | 待确认 |
| I-038-003 | required | 哪些现有批量/长操作进入异步首波，哪些保持同步？ | R1 范围冻结、R3 实施 | R1 | 盘点 `resources.go` 批量面、`data-transfer` 导出导入、wallet reconcile 现状，逐项给出「保持同步 / 改异步 / 不进首波」 | open | — | 待确认（同时承接 `V-F126`） |
| I-038-004 | required | 作业中心以新模块承载还是挂既有模块；是否进入默认 Profile 集？ | 激活、R1、VP-008 `go` 消费有效性 | 激活前 | 读 `kernel/profile.go` 现行模块矩阵与 `mvp`/`admin`/`demo` 集合，给出方案影响面与红线结论 | **verified** | 2026-09-19 用户 P-004 裁决方案 A | 新建 `admin.jobs` 进 admin 默认集（Profile 内容扩展，不改装配语义，不暂挂 `go`）；VRev-099 |
| I-038-005 | required | 激活前 Admin 类 freshness 与 VP-008 `go` 消费有效性是否仍成立？ | 激活与开区 | 激活前 | 执行 Admin 类 freshness review，核对协议 pin、依赖锁、迁移台账、Profile 默认集与装配、provenance 及区间变更 | **verified** | 2026-09-19 已完成 | `0c29c08` → `7e5ce891` 五域 PASS，不暂挂 `go`；VRev-099 |
| I-038-006 | non-blocking | 历史作业保留与清理策略（`expires_at` 已存在；是否需要归档/清理/容量上限）？ | 后续运维波次边界 | 关门后或出现容量触发 | 不纳入首波；出现真实容量或合规需求时由 `/vision` 复核 | deferred | 理由：首波聚焦可见性与结果读取，不新建数据生命周期程序；责任人：`/vision`；复核触发：作业表容量/保留期出现真实需求 | 待确认 |

`I-038-004` 与 `I-038-005` 已于 2026-09-19 关闭（用户 P-004 裁决 + Admin 类 freshness PASS），激活门禁解除。`I-038-001`～`003` 为 R1 冻结前的 required 门禁，仍为 open（侦察证据已收集，见 `02-execution/E-002-r1-recon.md` 与 `attachments/R1-recon-*.md`；冻结口径待 `GOAL-002` C1～C3 决策落盘，C2/C3 须经用户 P-004 裁决）。`I-038-003` 同时承接 `V-F126`。`I-038-006` 是有界延期，不代表已验证或承诺后续实现。

## 父目标

- Root 目标，`parent: null`。

## 台账布局

本目标从第一条记录起使用平铺 ledger：`01-decision/`、`02-execution/`、`03-audit/`。索引文件保留 frontmatter、摘要与条目链接；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*` 文件。

## 备注

- workspace/Root scaffold 是已发生事实；纲领 R1 已由 `GOAL-002` 交付并关门（`done · 4/4`），`progress: 1/5` 只由上方 5 个显式检查点派生，不放行阶段、不关闭 finding、不覆盖 status。
- 建区不代表任何实现阶段完成；VP-038 关门须链接本区证据并经用户确认。
- Vision Review `VRev-098`/`VRev-099` 属愿景层；Goal 审计须写入本目标 `03-audit/`，不能用 Vision Review 代替。
- `admin.jobs` 模块的建立属 VP-038 R2/R3 实现范围，本轮 scaffold **未**改动 `apps/**`。
