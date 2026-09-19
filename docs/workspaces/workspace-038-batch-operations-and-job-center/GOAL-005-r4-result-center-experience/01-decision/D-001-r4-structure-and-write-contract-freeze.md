---
title: D-001 · R4 结构选型与写面合同冻结
status: active
created: 2026-09-19
updated: 2026-09-19
parent: GOAL-005-r4-result-center-experience
version: 1.0.0
---

# D-001 · R4 结构选型与写面合同冻结

本决策冻结 GOAL-005（R4 结果中心与体验收敛）的两项**用户已裁决**的关键决策，以及由它们派生的实现口径。用户裁决原文见 §0；§1 起为派生口径，可经交叉审计质询，但不得与 §0 冲突。

## 0. 用户裁决（P-004 询问已发生，此处留痕）

| 信息项 | 用户裁决 | 影响的检查点 |
|--------|----------|--------------|
| `I-038-014` 结构选型 | **方案 A：在既有 `jobs` 列表页上收敛**（行操作 + `recordView` 详情抽屉）；**不新增**结果中心独立页、不新增导航节点 | C2 |
| `I-038-015` 取消/重试合同口径 | **逐字镜像既有 actor 作用域合同**：可取消/可重试状态集合、转换结果、错误码、重试次数上限一律与 `internal/jobs` 现有实现相同；**不重置 attempt 预算**；门控沿用 `jobs.write` | C1 |

两项均为用户书面裁决（会话内提问 → 用户选择），非编排器推断。

## 1. 写操作面冻结合同（C1 · 派生自 §0 的 `I-038-015`）

新增两条管理作用域路由，与 actor 作用域 wallet 写路径**并行且不改动**：

| 方法 / 路径 | 权限 | 语义 |
|-------------|------|------|
| `POST /api/jobs/{id}/cancel` | `jobs.write` | 取消**他人或本人**的作业 |
| `POST /api/jobs/{id}/retry` | `jobs.write` | 重试**他人或本人**的失败作业 |

### 1.1 状态集合与错误码（逐字镜像，不扩展）

| 操作 | 允许的前置状态 | 结果 | 其它状态 | 不存在 |
|------|----------------|------|----------|--------|
| cancel | `queued` | 立即 `cancelled`（`cancel_requested=0`、置 `finished_at`） | `409 JOB_NOT_CANCELLABLE` | `404 JOB_NOT_FOUND` |
| cancel | `running` | `cancel_requested=1`（worker 心跳察觉后走既有 `FinalizeCancel`） | 同上 | 同上 |
| retry | `failed` 且 `attempt < max_attempts` | 置 `queued`，清 `progress/cancel_requested/lease/result/error/finished_at/expires_at` | `409 JOB_NOT_RETRYABLE` | `404 JOB_NOT_FOUND` |

**明确不做**：不改动 Job 六态合同；不重置 `attempt`；不为 `succeeded` 作业提供重跑；不新增错误码（沿用已冻结的 `JOB_NOT_CANCELLABLE` / `JOB_NOT_RETRYABLE` / `JOB_NOT_FOUND`，避免 `error_contract_test` 需要新增注册）。

### 1.2 实现落点与"不放松既有隔离"

- 新方法写在**新文件** `apps/api/internal/jobs/actions.go`（`RequestCancelAny` / `RetryAny` / `(*Runner).CancelAny` / `(*Runner).RetryAny`），`repository.go` 与 `runner.go` 的既有 actor 作用域方法**不动**。
- `Any` 变体 = 既有实现**去掉 `actor_id` 谓词**，其余逐字相同（同一状态分支、同一 SQL 集合、同一错误码）。
- 既有 actor 作用域隔离由新增测试**反向钉住**：`RequestCancel` / `Retry` 对非本人仍返回 `ErrNotFound`（防止"新增通用路径"被误改为"放宽原路径"）。
- 组合根直接把 `*jobs.Runner` 绑定为 `handler.JobActions`（其 `CancelAny`/`RetryAny` 结构性满足该接口），避免第二份实现进入审计面。

### 1.3 描述符与计划同步

- provider `Descriptor().Contributions.Routes` 与 `kernel.BuiltinModules()` 的 `admin.jobs` 计划描述符**必须同时**含这两条路由；`Permissions` 保持 `jobs.read` + `jobs.write`（**不新增权限键**）。
- 路由仅在 `actions != nil` 时挂载；只读组合下这两条路由既不声明也不服务（声明=服务的既有不变量）。

## 2. 读投影的可操作性扩展（C2 前置）

`jobToMap`（`internal/handler/jobs.go`，R2 读投影）**新增 4 个派生字段**，作为 UI 行操作可用性的**唯一真相来源**（服务端派生，避免前端复制状态机）：

| 字段 | 派生规则 | 消费方 |
|------|----------|--------|
| `cancellable` | `status ∈ {queued, running}` | cancel 行操作的 `disabledWhen` |
| `retryable` | `status = failed && attempt < maxAttempts` | retry 行操作的 `disabledWhen` |
| `downloadable` | `status = succeeded` | download 行操作的 `disabledWhen`（过期态因此天然不可下载） |
| `statusStyle` | 六态 → `success`/`warning`/`destructive`/`info`/`neutral` | 状态列 `badgeStyleField` |

**理由**：行操作的可用状态必须与 §1.1 的服务端合同逐字一致；让服务端派生可保证二者不可能漂移。字段为**追加**（R2 已冻结的字段与语义不变，`resultUrl`/`progress` 等逐字保持），R2 测试无"投影键集合"断言，追加不破坏冻结语义。

## 3. 前端结构（C2 · 派生自 §0 的 `I-038-014` = 方案 A）

`apps/api/modules/jobs/schema/jobs.json` **原地增强**，不新增页面/导航/路由：

1. **能力声明追加**：`actions.row.request`、`actions.page.trigger`、`permissions.inheritance`、`record.view.load`（均由本页实际使用触发，`capability-declaration.guard.test.ts` 的 marker 全部命中）。
2. **页级 actions**：`cancelJob`（POST `/api/jobs/{id}/cancel`）、`retryJob`（POST `/api/jobs/{id}/retry`）、`downloadJobResult`（`type: custom`，`handler: jobs.downloadResult`）。
3. **表级 `permissions`**：`edit` = `$context.user.permissions contains "jobs.write"`（cancel/retry 的 `permissionIntent: "edit"`）。
4. **行操作**（`props.actions`）：cancel（带 `confirm`）、retry、download；三者均有 `requestMapping.path.id = "$row.id"`（download 由 custom handler 的 `{id}` 槽解析）与 §2 的 `disabledWhen`。行操作数 3 > `MAX_INLINE_ROW_ACTIONS`(2)，第 3 项落入既有 overflow 菜单，不新增渲染能力。
5. **列追加**：`error.code`（错误码）、`resultExpiresAt`、`correlationId`；`status` 列加 `badgeStyleField: "statusStyle"`。
6. **`recordView` 详情节点** `job-detail`：展示六态共 12 个字段（含 `error.message`、`resultUrl`、`correlationId`），闭包"六类状态可读"。

**download 的门控口径**：download 行操作**不使用** `permissionIntent`，改用**行动作本地 `permissions.edit` = `$context.user.permissions contains "jobs.read"`**。理由：服务端 `GET /api/jobs/{id}/result` 的门控是 `jobs.read`（R2 冻结），UI 不得比服务端更严；而冻结的 intent 集合仅 `{edit, delete}`，行动作本地 `permissions` 是被支持但此前未被使用的路径（`permissions.ts` `INTENT_KEYS` 循环）。该选择以测试钉住（本地 gate 生效、`jobs.write`-only 主体不被误放行/误拦截）。

**状态文案口径**：状态列显示**原始状态码 + 颜色徽标**，列头与筛选器选项走 i18n。理由：与既有产品约定一致（`wallet.status`、`scheduledTasks.enabled`、`users.mfaEnabled` 均显示原始值），且逐值本地化需要协议级扩展 → 触碰 pinned 工件，属本 VP 明确非目标。此条为**决定**（非未知），如需逐值本地化应作为后续 VP/协议的独立议题。

## 4. 结果下载与过期语义（C2 · 关闭 `I-038-016`）

- 渲染器新增白名单自定义 handler `jobs.downloadResult` → `GET /api/jobs/{id}/result`（`{id}` 槽沿用既有行上下文绑定）。
- 结果文档为批量导出时（含 `csv` + `fileName`）落盘为 CSV；否则落盘为 `job-<id>.json`（保持服务端 attachment 语义）。
- 该解析逻辑抽到共享模块 `apps/web/src/lib/job-result-download.ts`，被 R3 自定义组件与渲染器分支**共用**，避免两份 CSV 抽取实现。
- `I-038-016`（过期态呈现）：**`downloadable` 派生字段已把过期态排除在可下载之外**（按钮禁用 + `title` 提示），无需新增 410 专项 UI；服务端 410 `JOB_RESULT_EXPIRED` 仍在直接访问 URL 时兜底。结论落为"已由 §2 的派生字段关闭"，非"未处理"。

## 5. 进度可观察性（C2）

新增自定义组件 `jobs-auto-refresh`（`apps/web/src/components/jobs-auto-refresh.tsx`），形态与既有 `monitoring-auto-refresh` 一致：`off / 5s / 10s / 30s` 选择器，tick 调用 `crud.reloadList()`（页面级已声明的公开 seam）。

**已知副作用与边界**：`reloadList()` 会清空**本页**所有表选择（ADR-0022 D2）。`jobs` 页的表**未**声明 `table.selection`，无选择态可清；组件注释与测试钉住该前提，若将来 jobs 表启用选择，必须改用定向刷新或显式排除。

## 6. 非目标复述（本次决策未触碰）

- 不改 Job 六态合同；不同步化 `batch-delete`；不新增 Redis/外部队列/多实例/搜索引擎。
- 不触碰 pinned 工件（`docs/schemas/**`、`apps/web/src/protocol/upstream/**`）；`jobs.json` 只使用 pinned 允许的属性 + 仓库既有的本地扩展（`badgeStyleField`）。
- 不重开 VP-012/011/037/036；`jobs.json`/`users.json` 不声明 `table.selection` 或 `actions.batch.request`。

## 7. 残余与风险（R-1）

| 编号 | 内容 | 处置 |
|------|------|------|
| R-1.1 | `jobs.read`-only 主体在 UI 上可下载（与服务的 read 门控一致），但 `jobs.write` 主体集合目前等价（两者均 `PolicyAdmin`），差异不可现场构造 | 以嵌套性测试钉住；策略集一旦分离则需判别性测试 |
| R-1.2 | `jobs-auto-refresh` 依赖 `reloadList()` 的"清空全部选择"语义 | 组件内注释 + 交互级测试断言 jobs 表无选择态；jobs 表启用选择时必须重审 |
| R-1.3 | 状态列不做逐值本地化 | 见 §3 状态文案口径；如需变更走协议/后续 VP |
