---
id: E-001-r4-implementation
doc: execution-entry
parent: GOAL-005-r4-result-center-experience
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# E-001 · R4 C1～C3 实施（结果中心与体验收敛）

## 事实（2026-09-19）

### 0. 前置：两项用户 P-004 裁决

| 信息项 | 用户裁决 | 落点 |
|--------|----------|------|
| `I-038-014` 结构选型 | **方案 A**：既有 `jobs` 页原地收敛（行操作 + `recordView`），不新增页面/导航 | `D-001` §0/§3 |
| `I-038-015` 取消/重试口径 | **逐字镜像**既有 actor 作用域合同（状态集合/错误码/不重置 attempt/`jobs.write`） | `D-001` §0/§1 |

### 1. 落地产物

**C1 · 写操作面**

| 产物 | 位置 |
|------|------|
| `RequestCancelAny` / `RetryAny`（= 既有实现去掉 `actor_id` 谓词）+ `(*Runner).CancelAny` / `RetryAny`（含 active 执行取消 + scan 唤醒） | `internal/jobs/actions.go`（**新文件**） |
| `POST /api/jobs/{id}/cancel`、`POST /api/jobs/{id}/retry`（`jobs.write` 门控）+ `JobActions` 接口 + `writeJobActionError`（复用冻结码 `JOB_NOT_FOUND`/`JOB_NOT_CANCELLABLE`/`JOB_NOT_RETRYABLE`） | `internal/handler/jobs.go` |
| 投影派生字段 `cancellable`/`retryable`/`downloadable`/`statusStyle` + 扁平 `errorCode`/`errorMessage` | `internal/handler/jobs.go` `jobToMap` |
| provider：`New(a, reader, submitter, actions)`（typed 参数）+ 按绑定面声明路由与 `jobs.write` | `modules/jobs/provider.go` |
| 计划描述符补两条路由（权限集合不变） | `kernel/profile.go` |
| 组合根把 `*jobs.Runner` 直接绑定为 `handler.JobActions` | `internal/composition/composition.go` |

**C2 · 结果中心呈现**

| 产物 | 位置 |
|------|------|
| 能力 +`actions.row.request`/`actions.page.trigger`/`permissions.inheritance`/`record.view.load`；页级 `cancelJob`/`retryJob`/`downloadJobResult`；表级 `permissionCascade{edit}` + `permissions.edit = jobs.write`；行操作 cancel/retry/download（`requestMapping.path.id=$row.id` + `disabledWhen`）；列 +`errorCode`/`resultExpiresAt`/`correlationId`；`status` 加 `badgeStyleField`；`recordView` `job-detail` 15 字段 | `modules/jobs/schema/jobs.json` |
| 结果信封 → 文件决策唯一实现（CSV 用服务端 `fileName`，否则 JSON） | `apps/web/src/lib/job-result-download.ts`（新） |
| 白名单 handler `jobs.downloadResult`；`triggerBlobDownload` 移入共享模块 | `apps/web/src/renderer/render.tsx` |
| `jobs-auto-refresh`（off/5/10/30s，tick 调 `reloadList`） | `apps/web/src/components/jobs-auto-refresh.tsx`（新） |
| R3 组件改调共享下载助手（消除两份 CSV 抽取实现） | `apps/web/src/components/jobs-batch-export.tsx` |
| i18n 双目录各 +17 键 | `apps/web/src/i18n/messages/*.json` |

**C3 · 体验与测试**

| 产物 | 位置 |
|------|------|
| R4 交互契约 9 例（行操作启用态、冻结路由、下载、只读主体、自动刷新、recordView、主题 token） | `apps/web/src/renderer/jobs-result-center.test.tsx`（新） |
| R3 遗留交互测试 4 例（空选择、只接受 202、不 reloadList、轮询至终态下载） | `apps/web/src/components/jobs-batch-export.test.tsx`（新） |
| 共享下载助手单测 4 例（`I-038-013` 单源守卫） | `apps/web/src/lib/job-result-download.test.ts`（新） |
| 后端：写合同/隔离反向钉住/门禁/投影派生/描述符一致 | `internal/jobs/actions_test.go`、`internal/handler/jobs_actions_test.go`、`internal/handler/jobs_test.go`、`modules/jobs/provider_test.go` |

### 2. 实施中修正的真实契约问题

| # | 问题 | 处置 |
|---|------|------|
| 1 | `jobs.json` 最初把 `permissions` 写在 `props` 内 → 行操作门禁**未生效**：ADR-0023 中表节点对其 actions 的门禁只在 `permissionCascade.keys` 声明后才参与 | 改为表节点级 `permissionCascade{keys:[edit]}` + `permissions`（与 `scheduled-tasks.json` 同形）。**变异验证**：移除后「只读主体」用例变红 |
| 2 | 新增自定义节点后两个渲染型测试报 `unknown custom component: jobs-auto-refresh` | 按 R3 先例在 `main.tsx` + 10 个渲染型测试补 side-effect import |
| 3 | `download` 行操作无法用 `permissionIntent` 表达「服务端门控是 `jobs.read`」（冻结 intent 集合仅 `{edit,delete}`） | 改用**行动作本地 `permissions.edit`**（`permissions.ts` `INTENT_KEYS` 循环支持的既有路径），使 UI 不严于服务端；注释与测试固定该口径 |
| 4 | `CustomAction` schema 不允许 `onSuccess`（`additionalProperties:false`） | 移除 `downloadJobResult` 的 `onSuccess`，与 `exportUsers` 同形 |

### 3. 验证证据

| 验证 | 结果 |
|------|------|
| `go build ./...`（apps/api） | exit 0 |
| `go vet ./...`（apps/api） | 无输出 |
| `go test ./...`（apps/api） | **全绿**（含 `internal/jobs`、`internal/handler`、`internal/composition`、`kernel`、`modules/jobs`） |
| `npm run typecheck`（tsc -b + e2e tsconfig） | exit 0 |
| `npm run build`（vite） | exit 0 |
| `npm test`（vitest） | **117 files / 1455 tests 全绿**（R3 基线 114/1440 → +3 files / +15 tests） |

**新增测试覆盖**

| 测试 | 断言 |
|------|------|
| `TestRequestCancelAnyCancelsAnotherActorsQueuedJob` | Any 变体接受他人作业且立即 `cancelled`；**同例断言 `RequestCancel` 对非本人仍 `ErrNotFound`** |
| `TestRequestCancelAnyMarksAnotherActorsRunningJob` | running → `cancel_requested=1` 且不置 `finished_at` |
| `TestRequestCancelAnyRefusesTerminalStates` | 四终态 → `ErrNotCancellable` 且**行未被改动** |
| `TestRetryAnyRequeuesAnotherActorsFailedJob` | 清 progress/lease/result/error/finished/expires；**attempt 不重置**；`Retry` 对非本人仍 `ErrNotFound` |
| `TestRetryAnyRefusesEverythingButFailedWithBudget` | 耗尽的 failed → `ErrNotRetryable` 且**不重新武装**；queued/running/succeeded/cancelled 与不存在 |
| `TestRunnerCancelAnyStopsAnotherActorsRunningJob` | 经 active 执行注册表真实停止在飞 handler |
| `TestJobCancelAnySucceeds` / `TestJobRetryAnySucceeds` | 200 + 共享投影；retry 清错误、保留已消耗 attempt |
| `TestJobActionsMissingJobIs404` / `TestJobActionsRejectedTransitionIs409` | 冻结码 + **被拒不改行** |
| `TestJobActionsRequireJobsWrite` | 匿名 401 / editor 403 且**不触及状态机** |
| `TestJobActionRoutesAbsentWithoutActions` | actions 为 nil 时路由不挂载、不产生副作用 |
| `TestJobsProjectionDerivesActionAvailability` | 派生字段与写合同前置条件逐态一致（含 attempt 耗尽态） |
| `TestDescriptorDeclaresEveryBoundSurface` / `TestPlanDescriptorMatchesTheFullProvider` | 四种绑定组合的路由/权限集合；与 `kernel` 计划描述符一致 |
| `jobs-result-center.test.tsx`（9） | 六态列/状态色/无选择列/主题 token/逐行启用态/冻结路由/CSV 下载/只读主体/自动刷新/recordView |
| `jobs-batch-export.test.tsx`（4） | 空选择禁用且不触网/提交体与**不 reloadList**/只接受 202/轮询至终态并下载 |
| `job-result-download.test.ts`（4） | 信封识别/文件名优先级/CSV 与 JSON 两条落盘路径 |

### 4. 边界

- **未**触碰任何 pinned 协议工件（`docs/schemas/**`、`apps/web/src/protocol/upstream/**`）。
- **未**改 Job 六态合同；**未**同步化 `batch-delete`；**未**重开 VP-012/011/037/036。
- **未**改 `repository.go`/`runner.go` 的既有 actor 作用域方法（隔离语义与测试逐字保持）。
- `jobs.json` **未**声明 `table.selection` 或 `actions.batch.request`；未新增权限键（仍 `jobs.read` + `jobs.write`）。
- `npm run build` 会重写 `apps/web/public/protocol/conformance-claim*.json`（buildId 随 HEAD 变化）；本次已还原该生成物，未纳入提交（不作为 R4 事实）。

### 5. Git checkpoints

| hash | 内容 |
|------|------|
| `ea6e4006` | R4 C1 后端写面（actions/handler/provider/profile/composition + 后端测试） |
| `f1351534` | R4 C2 前端结果中心（jobs.json/下载助手/渲染器/自动刷新/i18n/夹具） |
| `bafb766f` | R4 C3 交互级与单源测试（3 个新测试文件） |

### 6. 未做（移交 C4）

- **未**执行审计：C4 待跑（self + grok build independent，模式 `cross`）。
- **未**投影 Root R4 检查点（Root 仍 `active · 3/5`）。
- **未**在真实浏览器会话验证（R5 退出矩阵/自动化回归承接）。
