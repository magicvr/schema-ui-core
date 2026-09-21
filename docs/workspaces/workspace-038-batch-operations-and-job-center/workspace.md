---
id: workspace-038-batch-operations-and-job-center
title: Admin 批量操作与异步结果中心工作区
status: done
root_goal: GOAL-001-batch-operations-and-job-center
canonical_scope: docs/workspaces/workspace-038-batch-operations-and-job-center/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-038-batch-operations-and-job-center
primary_plan: VP-038-batch-operations-and-job-center
created: 2026-09-19
updated: 2026-09-19
version: 1.0.0
parent: null
---

# 工作区上下文 · Admin 批量操作与异步结果中心

本工作区是 [VP-038-batch-operations-and-job-center](../../vision/plans/VP-038-batch-operations-and-job-center.md) 的唯一 `delivery` workspace，承接 VP-012 首波显式排除的「通用 Job 管理页」与 roadmap「体验增强」清单中的「批量结果中心」。它不重开 VP-012/011/037/036，不属于 VP-010 符合性整改，不承载实体全文检索、组织/数据权限、新业务域或 Redis/MQ/多实例/外部队列。

- VP-038 于 2026-09-19 经用户确认从 `planned` 激活为 **`active` v0.2.0**（计划 self = `VRev-098` `pass`；激活 self = `VRev-099` `pass`，0 required）。
- **P-004 裁决（`I-038-004`，2026-09-19）**：新建 **`admin.jobs`** 模块并进入 **admin 默认集**；`mvp` / `demo` 不启用。定性 = **Profile 内容扩展**（沿用 S 系列先例），**不改装配语义** → **不暂挂 VP-008 `go`**。
- **Admin 类 freshness PASS**：`0c29c08`（VP-037 激活基线）→ `7e5ce891`（HEAD）；协议 pin `v2.9.0` / `81aa1d8`、依赖锁、迁移台账、Profile 默认集与装配、provenance 五域零变更；区间 `apps/**` 变更全部为 VP-037 已审结目。
- Root `[workspace-038-batch-operations-and-job-center] GOAL-001-batch-operations-and-job-center`：**`done · 5/5`**；纲领 R1→R5 全部完成；VP-038 经**用户书面确认**于 2026-09-19 `closed` v1.0.0。
- **R1 已完成并关门**（2026-09-19 · 子目标 `GOAL-002-r1-denominator-and-contract-freeze` **`done · 4/4`**）：C1 分母与作用域矩阵 / C2 契约形态冻结 / C3 首波分母冻结 / C4 R1 审计与投影。审计模式 `cross`——self `A-001` `pass` → grok build（grok-4.6 · high · `/audit`）independent `A-002` `conditional`（3 required）→ `A-003` 响应 required 全 `fixed`，**开放 required = 0**。
- **R2 已完成并关门**（2026-09-19 · 子目标 `GOAL-003-r2-generic-job-read-surface` **`done · 4/4`**）：C1 模块与权限接线 / C2 查询与索引 / C3 读面 API 与作用域 / C4 R2 审计与投影。审计模式 `cross`——self `A-001` `pass` → grok build independent `A-002` **`pass`** → `A-003` 响应 6 条 recommended 全 `fixed`，**开放 required = 0**。
- **R3 已完成并关门**（2026-09-19 · 子目标 `GOAL-004-r3-async-batch-operation` **`done · 4/4`**）：C1 异步写面与进度 / C2 前端触发 / C3 同步路径回归 / C4 R3 审计与投影。审计模式 `cross`——self `A-001` `pass` → grok build independent `A-002` **`pass`** → `A-003` 响应 8 条 recommended 全 `fixed`，**开放 required = 0**。
- **R4 已完成并关门**（2026-09-19 · 子目标 `GOAL-005-r4-result-center-experience` **`done · 4/4`**）：C1 写操作面 / C2 结果中心呈现 / C3 体验与测试收敛 / C4 审计与投影。审计模式 `cross`——self `A-001` `pass` + grok build independent `A-002` **`pass`** → `A-003` 响应 5 条 recommended `fixed` + 3 条经**用户书面裁决** `accepted-residual`，**开放 required = 0**。**R4 两项关键决策由用户 P-004 裁决**：`I-038-014` = 方案 A（既有 `jobs` 页原地收敛：行操作 + `recordView`；不新增页面/导航）；`I-038-015` = **逐字镜像**既有 actor 作用域取消/重试合同（同状态集合、同冻结错误码、不重置 attempt、`jobs.write` 门控）。
- **R4 残余移交**（用户 2026-09-19 指令）：`GOAL-005 A-001` 三条 low 级项（通用列值本地化 / 表格定向刷新 seam / 空闲不轮询）经本区**有界接受**后移交 `[workspace-010-design-implementation-conformance] GOAL-044-w32-r4-residual-seams`，该目标**当日 `done · 4/4`** 交付并已回填本区 `A-003` 为 `fixed`。
- **R5 已完成并关门**（2026-09-19 · 子目标 `GOAL-006-r5-evidence-and-closeout` **`done · 4/4`**）：C1 退出矩阵 / C2 浏览器回归 / C3 独立意见与残余登记 / C4 组合投影与关门提请。审计模式 `cross`——self `A-001` `pass` + grok build independent `A-002` **`pass`** → `A-003`/`A-004` 响应后 **开放 required = 0**。e2e 双 profile 全绿（mvp 16 passed / 5 skipped / 0 failed；admin **17 passed / 4 skipped / 0 failed**，含新增 `e2e/jobs-result-center.spec.ts` 端到端用例：选择行 → 导出所选 → 真实进度 → 终态下载 → 结果中心读同一作业 → 行操作再次下载）。**VP-038 经用户书面确认 `closed` v1.0.0；Root `done · 5/5`**。关门后残余 1 条 bounded residual（e2e fresh-seed 顺序契约，已登记 roadmap）；关门 Vision Review 未执行（如实登记）。
- **R1 冻结要点（用户 P-004 裁决 2026-09-19）**：C2 = **方案 B**（另立本地模块自有异步契约，ADR-0022 同步 `batch-delete` 语义逐字冻结，协议 pin 零改动）；C3 首波 = **仅「新建批量导出所选」1 条**；C1 = **管理作用域 + 新增 `jobs.read` 权限**（`PolicyAdmin`，既有 actor 隔离语义与测试不动）。
- **R2 交付要点**：新建 `admin.jobs` 模块（进 admin 默认集）+ `jobs.read`（`PolicyAdmin`）+ `GET /api/jobs`、`/{id}`、`/{id}/result`（管理作用域 + fail-closed）+ 迁移 v72 管理列表索引 `(created_at DESC, id DESC)`；并修 R-1（`jobRuntime.enabled` 改为 `admin.jobs` 或 `admin.wallet` 任一存在即启用）。跨 VP 触碰：wallet 结果地址改经共享 `jobs.ResultURL`，**输出字符串逐字不变**（独立审计复核通过）。
- **R3 交付要点**：`jobs.batch-export` 作业 + `POST /api/jobs/batch-export`（**202 + jobId**，**双重门禁** `jobs.write` 且 `data.export`，使异步路径不扩大数据外带面）+ **真实细粒度进度**（按选中行上报，非硬编码）+ `jobs-batch-export` 自定义组件（users 页多选 + 退避轮询 + 下载，**不经** ADR-0022 `batchMapping`、**不调用** `reloadList()`）。同步 `batch-delete` 与 Job 六态合同逐字未退化（独立审计以空 diff 核实）。
- 用户已确认 workspace slug = `workspace-038-batch-operations-and-job-center`；Root slug = `GOAL-001-batch-operations-and-job-center`。
- Vision open required：0；`V-F125` → `fixed`；`V-F126` 保持 `open · recommended`，承接动作已由 `I-038-003` 完成（闭合登记属 `/vision`，见首波矩阵 §7 交接项）。
- `I-038-001`～`003` 与 `I-038-007`～`009` 已关闭为 `verified`；`I-038-013`～`016` 已由 `GOAL-005` 关闭为 `verified`；R3 移交项见 `D-001` §5 的 T-5/T-6（O-1/O-2 方案已在 R2 冻结，实现属 R3/R4）。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-038-batch-operations-and-job-center` | 与本区目标及资料引用的 `workspace_id` 一致；当前无固定共享资料 |
| Root Goal | `GOAL-001-batch-operations-and-job-center` | `parent: null`；**done · 5/5** |
| canonical 范围 | `docs/workspaces/workspace-038-batch-operations-and-job-center/` | 本区唯一目标状态范围 |
| 共享资料目录 | `none` | 本区暂无固定共享资料；不得声明共享资料引用 |
| 愿景角色 | `delivery` | VP-038 唯一 delivery workspace；不改变 Charter primary workspace |
| 规划对齐 | `primary_plan` = `VP-038-batch-operations-and-job-center` | `plan_refs` 必填且已精确绑定 |

## 愿景对齐

- Charter：`schema-ui-core-admin-foundation@0.4.0`
- VP：[VP-038-batch-operations-and-job-center](../../vision/plans/VP-038-batch-operations-and-job-center.md)（**`closed` · v1.0.0** · 2026-09-19 用户书面确认关门）
- 计划审视：[VRev-098](../../vision/reviews/VRev-098-vp038-batch-operations-job-center-planned.md) self `pass`
- 激活审视：[VRev-099](../../vision/reviews/VRev-099-vp038-batch-operations-and-job-center-activation.md) self `pass`
- Vision open required：0；`V-F126` `open · recommended`（承接动作已由 `I-038-003` 完成；闭合登记属 `/vision`）；`I-038-004`/`I-038-005` `verified`；`I-038-001`～`003` **`verified`**；`I-038-013`～`016` **`verified`**；`I-038-006` `deferred · non-blocking`

## 纲领阶段

| 阶段 | 目的 | 状态 |
|------|------|------|
| R1 | 分母与契约冻结：Job 种类×作用域矩阵、批量异步契约（含同步 `batch-delete` 兼容口径）、首波批量操作分母、权限/Profile 边界与排除项 | **done**（`GOAL-002` `done · 4/4`；`I-038-001`～`003` verified；审计 `A-001` self pass + `A-002` grok independent conditional → `A-003` required 全 fixed，开放 required = 0） |
| R2 | 通用作业读面：Job 列表/详情/结果读取 API + 权限与作用域过滤 + fail-closed；不改变 Job 六态合同 | **done**（`GOAL-003` `done · 4/4`；`admin.jobs` 模块 + `jobs.read` + 迁移 v72；审计 `A-001` self pass + `A-002` grok independent **pass** → `A-003` 全 fixed，开放 required = 0） |
| R3 | 批量操作异步承接：至少一条真实批量操作走 Job（202 + jobId + 进度 + 结果）；同步 `batch-delete` 回归不退化 | **done**（`GOAL-004` `done · 4/4`；`jobs.batch-export` + 双重门禁 + 真实进度 + 前端组件；审计 `A-001` self pass + `A-002` grok independent **pass** → `A-003` 全 fixed，开放 required = 0） |
| R4 | 结果中心与体验收敛：列表/详情/进度/终态/过期/重试/取消/下载；中英文、浅色深色、加载空态错误态、可访问 | **done**（`GOAL-005` `done · 4/4`；管理作用域取消/重试 + `jobs` 页行操作/`recordView`/下载/自动刷新 + 交互级测试（含 R3 遗留 4 例）与 `I-038-013` 单源守卫；审计 `A-001` self pass + `A-002` grok independent **pass** → `A-003` 5 fixed + 3 用户裁决 accepted-residual（移交 `[workspace-010]` `GOAL-044`），开放 required = 0） |
| R5 | 证据与关门：退出矩阵、浏览器/自动化回归、独立意见、残余登记、组合投影同步与用户确认 | **done**（`GOAL-006` `done · 4/4`；判据 1～7 达成；cross 审计 self `A-001` + grok independent `A-002` 均 `pass` → `A-003`/`A-004` 开放 required = 0；e2e 双 profile 全绿含 jobs 端到端用例；**用户书面确认 VP-038 关门**） |

纲领阶段按 R1 → R2 → R3 → R4 → R5 串行推进；工作区建立本身不代表任何实现阶段完成。

## 固定共享资料引用

> `shared-materials/index.json` 只能提供候选路径与摘要。本区 `shared_materials_catalog: none`，当前无可作为事实依据的共享资料引用。

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | — | — | none |

## 备注

本工作区只保存 VP-038 的实现层目标、决策、执行事实与 Goal 审计；不得将 Vision Review 或 VP 状态复制成第二套目标状态源。愿景组合编排仍以 `docs/vision/` 为准；本区 `goal-tree.md` 与目标五件套是实现层状态真相源。
