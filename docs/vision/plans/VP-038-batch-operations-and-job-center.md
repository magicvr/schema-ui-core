---
doc_type: vision-plan
id: VP-038-batch-operations-and-job-center
title: Admin 批量操作与异步结果中心
status: active
vision_ref: schema-ui-core-admin-foundation@0.4.0
lead_workspace: workspace-038-batch-operations-and-job-center
created: 2026-09-19
updated: 2026-09-19
version: 0.2.0
parent: null
---

# VP-038 · Admin 批量操作与异步结果中心

## 状态、激活与关门门禁

| 项 | 值 |
|-----|-----|
| status | **`active`**（2026-09-19 · v0.2.0 · 用户确认激活；lead `workspace-038-batch-operations-and-job-center`） |
| 组合位置 | **Admin 功能分支 · 体验增强**；承接 VP-037 之后的工作流连续性下一拍，也是 `roadmap.md`「体验增强」清单中「批量结果中心」的承接者 |
| Vision Review | 计划阶段 [VRev-098](../reviews/VRev-098-vp038-batch-operations-job-center-planned.md) self `pass`；激活就绪 [VRev-099](../reviews/VRev-099-vp038-batch-operations-job-center-activation.md) self `pass`；两条报告 open required = 0 |
| 激活门禁 | **已满足**：`I-038-004` 用户 P-004 裁决 = 方案 A（新建 `admin.jobs`，进入 admin 默认集；Profile 内容扩展，不暂挂 `go`）；`I-038-005` Admin 类 freshness **PASS**（`0c29c08` → `7e5ce891`）；用户确认 workspace/Root slug |
| 基础设施边界 | 首波不消耗 Redis、MQ、多实例、外部队列或专用搜索引擎 trigger；不重开 VP-012 |

## 用户已裁决（2026-09-19 · P-004）

| 项 | 裁决 |
|----|------|
| 组合层下一拍 | **VP-038**（备选：架构 C1 `timestamptz` 合同 / 版本与维护提示 / 不立项） |
| 架构候选 C1（`RES-T03-tz`） | **保持登记、待本波之后单独裁决**，不并入 VP-038 |
| `I-038-004` 模块边界 | **方案 A**：新建 **`admin.jobs`** 模块，进入 **admin 默认集**；`mvp` / `demo` 不启用。定性 = **Profile 内容扩展**（沿用 S 系列先例：file-library / data-dictionary / system-monitoring / scheduled-tasks / recycle-bin / mfa / wallet 等），**不改装配语义**（`ResolveProfile` 逻辑、Manifest 聚合规则、协议 pin、共同门禁语义零改动）→ **不暂挂 VP-008 `go`** |
| workspace / Root slug | `workspace-038-batch-operations-and-job-center` / `GOAL-001-batch-operations-and-job-center` |

## 意图

`VP-012` 已交付 profile 无关的 **Job 六态运行时**（queued / running / succeeded / failed / cancelled / expired + 进度 / 重试 / 取消 / 结果过期），但它当时**明确把「通用 Job 管理页」写为非目标**；`core.jobs` 至今是**迁移专用模块**（`Provider.Register` 空实现，注释声明其刻意不进入运行期 Profile）。结果是：Job 运行时只有 `wallet.reconcile` 一个消费者，且只能经 wallet 模块自己的 `/api/wallet/jobs/*` 访问，没有任何跨模块的作业可见性。

同时，ADR-0022 的批量请求（`actions.batch.request` / `batchMapping`）已在 users / roles / data-dictionary / scheduled-tasks 落地，但只有 `batch-delete` 一种、且是**同步**的（返回 `{"deleted": n}` 后由前端 reload）；`data-transfer` 的导出与导入同样是同步请求。长操作既没有进度，也没有可回看的结果。

本 VP 把「批量操作的异步进度与结果中心」作为有界产品能力交付：让**已注册的作业种类**与**纳入首波的批量操作**在 Admin 里可观察、可操作、可追溯，而不是再建一套队列或搜索基础设施。

本 VP 是面向用户的 Admin 体验增强，不是 VP-010 的符合性整改，不承载业务域，也不重开 VP-012。

## 首波范围与边界

| 范围 | 本 VP 首波 | 不在本 VP |
|------|-----------|-----------|
| 作业可见性 | 已注册 Job 种类的列表、详情、六态、进度、attempt、错误码/消息、结果读取 | 调度器/worker 内部语义改造；Job 运行时六态合同本身的重新定义 |
| 作用域 | 按当前权限与数据范围过滤（管理作用域 vs actor 作用域由 R1 冻结） | 新的权限绕过、组织/多租户权限模型、SSO |
| 批量操作 | **至少一条**真实批量操作以异步 Job 承接（202 + jobId → 进度 → 结果） | 把已交付的同步 `batch-delete` 改成 breaking 异步语义 |
| 结果中心 | 进行中/成功/失败/取消/过期结果可读取；导出类结果可下载；失败可按合同重试 | 通用 BI/报表中心、跨模块数据仓库、批量审批工作流 |
| 体验 | 中英文、浅色/深色、加载/空态/错误态、可访问、与既有列表/反馈约定一致 | 实体全文检索、Saved Views 重做、Toast 全局重做、新业务域 |
| 基础设施 | 消费既有 `internal/jobs` 运行时与既有 Store；保留外部队列接缝 | Redis、外部队列 / broker、多实例、跨进程索引、专用搜索引擎 |

## 与相邻 VP / 路线图的边界

| VP / 方向 | 关系 |
|-----------|------|
| **VP-012** | **消费**已交付的 Job 六态合同与 `wallet.reconcile` 先例，**不重开** VP-012，不改其首波冻结表与已交付语义；「通用 Job 管理页」正是本 VP 承接的、当时被显式排除的项 |
| **VP-011** | 消费其已交付的列表/资源/导出导入模块面；不重开 VP-011 |
| **VP-037** | 复用其已交付的统一反馈语义与列表页基线（成功/失败/重试/维护反馈）；不重开 VP-037，不改 Saved View 存储格式或查询/重置合同 |
| **VP-036** | 若作业页面需要进入 Command Palette 检索分母，按 VP-036 已交付的 provider 接缝注册；不重开 VP-036 |
| **VP-009** | 共享基架安全问题（含作业/批量路径的鉴权与越权）仍归持续生产加固 |
| **VP-010** | 本 VP 是向好演进的产品能力，不是 as-designed / as-built 偏差整改；若实现发现既有协议符合性问题，转 VP-010 |
| **VP-008 `go`** | **关键边界（已由 2026-09-19 P-004 裁决落定）**：新建 `admin.jobs` 进 admin 默认集属 **Profile 内容扩展**，沿用 S 系列先例（file-library / data-dictionary / system-monitoring / scheduled-tasks / recycle-bin / mfa / wallet），**不改装配语义**（`ResolveProfile` 逻辑、Manifest 聚合规则、协议 pin、共同门禁语义零改动）→ **不暂挂 `go`**。若实施期改为改动装配语义或默认集**结构**，仍须按 freshness / `go` 规则暂停与复核 |
| **架构分支** | 外部队列 / broker（`RT-Q02`）、Redis（`RT-Q03`）、多实例（A3）保持 `trigger-gated`；本 VP 只在既有进程内 Job 运行时上做产品面 |
| **业务域分支** | 不新增 Catalog、订单、支付、库存、CMS 等业务域；作业中心是基架能力，不把某个业务域的流程搬进来 |

## 方向级退出判据

在同时满足下列方向时，本 VP **可以**有界或完整关门（证据必须在工作区目标内）：

1. **分母与契约冻结**：纳入首波的 Job 种类、可见作用域（管理 vs actor）、六态与进度/结果/过期的对外语义、批量异步契约（请求形状、202/jobId、轮询或读取路径）形成可机器核对的矩阵；排除项显式点名。
2. **通用作业可见性**：已注册作业可按权限列出并查看详情（状态、进度、attempt/max_attempts、错误码与消息、时间戳、correlation）；无权限、越作用域或不存在时 fail-closed，不泄露他人作业。
3. **批量操作异步承接**：至少一条真实批量操作以异步 Job 承接，可观察到 queued→running→终态与进度，并可在终态读取结果；**同时**既有 ADR-0022 同步 `batch-delete` 的已交付语义与回归不退化。
4. **结果中心体验**：进行中、成功、失败、取消、过期与结果已过期六类呈现可用；导出类结果可下载；失败在合同允许时可重试；中英文、浅色/深色、加载/空态/错误态与既有产品约定一致。
5. **权限与 Profile 安全**：作业与结果按当前 Profile 与权限过滤，不绕过既有路由守卫与数据范围（含 self scope）；mvp / admin（及适用的 demo/custom）覆盖有矩阵证据；直接 URL 行为与现有守卫一致。
6. **基础设施与范围保持**：未实现或解除 Redis、外部队列/broker、多实例、跨进程索引或专用搜索引擎；未重开 VP-012/011/037；未改 Charter 目的/边界/非目标；未把新业务域混入本 VP。
7. **证据与审计**：退出矩阵、浏览器/自动化回归、必要的独立意见已落盘；开放 required finding = 0；组合投影同步且经用户书面确认关门。

## P-005 信息需求

| 编号 | 所需信息 | 级别 | 影响 | 最晚需要 | 收集/验证动作 | 状态 | 延期/复核与证据 |
|------|----------|------|------|----------|----------------|------|------------------|
| I-038-001 | Job 分母与可见作用域：现存及首波纳入的 Job 种类、是否展示全量/本 actor/按数据范围、与 `wallet.reconcile` 既有 actor 作用域的兼容关系 | required | R1 范围冻结、R2 读面、判据 1/2 | R1 | 扫描 `internal/jobs` 与各模块 `Register*JobKind` 消费点，形成种类×作用域矩阵；确认是否需要新增 repository 列表查询 | open | — |
| I-038-002 | 批量异步契约：是否扩展 ADR-0022 以支持异步变体，还是为长操作另立独立契约；同步 `batch-delete` 是否保持；是否触碰 pinned `schema-ui-docs@v2.9.0` 协议面 | required | R1 方案冻结、R3 实施、判据 1/3/6 | R1 | 对照 `apps/web/src/protocol`（`actions.batch.request` / `batchMapping`）与上游 v2.9.0 契约；判定「本地扩展」vs「上游协议变更」，并给出兼容与回归口径 | open | — |
| I-038-003 | 首波批量操作分母：哪些现有批量/长操作（批量删除、批量启停、导出、导入、对账等）进入异步首波，哪些保持同步 | required | R1 范围冻结、R3 实施、判据 3 | R1 | 盘点 `resources.go` 批量面、`data-transfer` 导出导入、wallet reconcile 现状，逐项给出「保持同步 / 改异步 / 不进首波」 | open | — |
| I-038-004 | **Profile / 模块矩阵边界**：作业中心以新模块（如 `admin.jobs`）承载还是挂在既有模块；是否进入默认 Profile 集 | required | 激活、R1、VP-008 `go` 消费有效性、判据 5/6 | **激活前**（须用户 P-004 裁决） | 读 `kernel/profile.go` 现行模块矩阵与 `mvp`/`admin`/`demo` 集合；给出「新增模块进默认集」与「挂既有模块」两方案的影响面与红线结论 | **verified**（2026-09-19 用户 P-004 裁决方案 A：新建 `admin.jobs` 进 admin 默认集；Profile 内容扩展，不改装配语义，不暂挂 `go`） | VRev-099；`kernel/profile.go:46`–`93` ProfileAdmin 先例 |
| I-038-005 | 激活前 Admin 类 freshness 与 VP-008 `go` 消费有效性 | required | 激活与开区 | 激活前 | 执行 Admin 类 freshness review，核对协议 pin、依赖锁、迁移台账、Profile 默认集与装配、provenance 及区间变更 | **verified**（2026-09-19：`0c29c08` → `7e5ce891` 五域 PASS；不暂挂 `go`） | VRev-099 |
| I-038-006 | 历史作业保留与清理策略（`expires_at` 已存在；是否需要归档/清理/容量上限） | non-blocking | 后续运维波次边界 | 关门后或出现容量触发 | 不纳入首波；出现真实容量或合规需求时由 `/vision` 复核 | deferred | 延期理由：首波聚焦可见性与结果读取，不新建数据生命周期程序；责任人：`/vision`；复核触发：作业表容量/保留期出现真实需求 |

`I-038-004` 与 `I-038-005` 已于 2026-09-19 关闭（用户 P-004 裁决 + Admin 类 freshness PASS），激活门禁解除。`I-038-001`～`I-038-003` 为 R1 冻结前的 required 门禁，仍为 open。`I-038-006` 是有界延期，不代表已验证或承诺后续实现。

## 纲领路线图

```text
R1 分母与契约冻结
   Job 种类×作用域矩阵、批量异步契约（含同步 batch-delete 兼容口径）、首波批量操作分母、
   权限/Profile 边界与排除项；I-038-001～004 关闭
  → R2 通用作业读面
     Job 列表/详情/结果读取 API + 权限与作用域过滤 + fail-closed；不改变 Job 六态合同
    → R3 批量操作异步承接
       至少一条真实批量操作走 Job（202 + jobId + 进度 + 结果）；同步 batch-delete 回归不退化
      → R4 结果中心与体验收敛
         列表/详情/进度/终态/过期/重试/取消/下载；中英文、浅色深色、加载空态错误态、可访问
        → R5 证据与关门
           退出矩阵、浏览器/自动化回归、独立意见、残余登记、组合投影同步与用户确认
```

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-038-batch-operations-and-job-center | GOAL-001-batch-operations-and-job-center | lead | 2026-09-19 | 唯一 lead delivery 工作区；用户确认 slug；由 `/govern` scaffold（Root 纲领 R1→R5） |

## 关门记录

（仅 `closed` / `abandoned` 时填写。）

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| — | — | — | — | — |

## 规划修订短史

| date | change |
|------|--------|
| 2026-09-19 | 初创（`planned` · v0.1.0 · 0 区）。用户 2026-09-19 确认组合层下一拍 = 本方向（备选：架构 C1 `timestamptz` 合同 / 版本与维护提示 / 不立项）。计划阶段 self Review = [VRev-098](../reviews/VRev-098-vp038-batch-operations-job-center-planned.md) `pass`（0 required）。C1 按用户同日裁决保持登记、待本波之后单独裁决，不并入本 VP。 |
| 2026-09-19 | 激活（`planned → active` · v0.2.0）。用户确认「激活，然后开设工作区」。**P-004 裁决 `I-038-004` = 方案 A**：新建 `admin.jobs` 进 admin 默认集（Profile 内容扩展，不改装配语义，不暂挂 VP-008 `go`）。Admin 类 freshness **PASS**（`0c29c08` → `7e5ce891`：协议 pin / 依赖锁 / 迁移台账 / Profile 默认集与装配 / provenance 五域零变更）。用户确认 slug `workspace-038-batch-operations-and-job-center` / `GOAL-001-batch-operations-and-job-center`。激活就绪 self = [VRev-099](../reviews/VRev-099-vp038-batch-operations-and-job-center-activation.md) `pass`（0 required）。lead 交 `/govern` scaffold。 |
