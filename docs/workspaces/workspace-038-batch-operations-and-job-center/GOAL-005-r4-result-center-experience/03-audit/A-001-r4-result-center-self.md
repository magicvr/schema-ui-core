---
id: A-001-r4-result-center-self
doc: audit-entry
parent: GOAL-005-r4-result-center-experience
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# A-001 · R4 结果中心自审（GOAL-005 C1～C3）

## A-001 · R4 C1～C3 自审（2026-09-19）

- **source**：self
- **auditor**：`/govern` 编排器（DeepSeek Harness 会话）
- **类型** / **scope**：`stage` · GOAL-005 C1～C3（写操作面 / 结果中心呈现 / 体验与测试收敛）
- **verdict**：**pass**（0 required；4 recommended）
- **完整意见**：本文件

### 范围与区间

- 被审目标：`docs/workspaces/workspace-038-batch-operations-and-job-center/GOAL-005-r4-result-center-experience/`
- 工作区校验：`root_goal` = `GOAL-001-batch-operations-and-job-center`、`canonical_scope` = 本区路径、`vision_role: delivery`、`plan_refs`/`primary_plan` = `VP-038-…` —— 绑定一致，未跨区。
- 审计区间：`e1893a1a`（R3 关门）→ `ab215ebc`（R4 C1～C3 + F-001 加固）。**不含** C4 本身与 R5。
- 审计材料：`00-meta.md`、`01-decision/D-001-r4-structure-and-write-contract-freeze.md`、`02-execution/E-001-r4-implementation.md`；代码 `internal/jobs/actions.go`、`internal/handler/jobs.go`、`modules/jobs/provider.go`、`modules/jobs/schema/jobs.json`、`kernel/profile.go`、`internal/composition/composition.go`、`apps/web/src/lib/job-result-download.ts`、`apps/web/src/renderer/render.tsx`、`apps/web/src/components/jobs-auto-refresh.tsx`、`apps/web/src/components/jobs-batch-export.tsx`。

### 成果（有证据）

| # | 成果 | 证据 |
|---|------|------|
| 1 | 管理作用域取消/重试：两条写路由 + 仓储/runner 方法，**状态集合与错误码逐字镜像** actor 作用域 | `internal/jobs/actions.go`；`TestRequestCancelAny*`、`TestRetryAny*` |
| 2 | **既有 actor 隔离未被放松**：`repository.go`/`runner.go` 在本区间 **0 行改动**（`git diff --stat e1893a1a..HEAD -- …/repository.go …/runner.go` 为空），且新测试**反向钉住** `RequestCancel`/`Retry` 对非本人仍 `ErrNotFound` | `TestRequestCancelAnyCancelsAnotherActorsQueuedJob`、`TestRetryAnyRequeuesAnotherActorsFailedJob` |
| 3 | 取消语义与既有实现一致：queued → 立即 `cancelled`（含 `finished_at`）；running → `cancel_requested=1` 且经 active 执行注册表**真实停止在飞 handler** | `actions.go`；`TestRunnerCancelAnyStopsAnotherActorsRunningJob` |
| 4 | 重试**不重置 attempt 预算**（VP-038 非目标），并清 progress/lease/result/error/expires | `actions.go`；`TestRetryAnyRequeuesAnotherActorsFailedJob`（断言 attempt=1/3 保留） |
| 5 | 门控 fail-closed：`jobs.write`；匿名 401 / editor 403，且被拒请求**不触及状态机** | `jobs_actions_test.go` 的 `TestJobActionsRequireJobsWrite` |
| 6 | 被拒转换（409）与不存在（404）**不改动数据行** | `TestRequestCancelAnyRefusesTerminalStates`、`TestJobActionsRejectedTransitionIs409`（比较 status/attempt/updatedAt） |
| 7 | 声明=服务：actions 为 nil 时两条写路由既不声明也不挂载；四种绑定组合的描述符与计划描述符逐项一致 | `provider.go`；`TestJobActionRoutesAbsentWithoutActions`、`TestDescriptorDeclaresEveryBoundSurface`、`TestPlanDescriptorMatchesTheFullProvider` |
| 8 | UI 可用状态由**服务端派生**（`cancellable`/`retryable`/`downloadable`），与写合同前置条件同源，浏览器不复制状态机 | `jobToMap`；`TestJobsProjectionDerivesActionAvailability` |
| 9 | 结果中心在**真实 `jobs.json`** 上可用：六态列 + 状态色 + 进度 + 详情抽屉 + 取消/重试/下载 + 自动刷新 | `renderer/jobs-result-center.test.tsx`（9 例，经 D-VAL → loadPageDocument → RenderPage） |
| 10 | 下载为**单一实现**：CSV 用服务端 `fileName`，非导出作业落 JSON；R3 组件与渲染器分支共用（关闭 `I-038-013`） | `lib/job-result-download.ts`；`job-result-download.test.ts`；`jobs-batch-export.test.tsx` 的下载用例 |
| 11 | R3 遗留交互契约落地：空选择禁用、**只接受 202**、**不调用 `reloadList()`**、轮询至终态 | `jobs-batch-export.test.tsx`（4 例） |
| 12 | 中英文双目录键集合对齐（1220/1220）；无固定调色板类（整页 token 断言）；自动刷新控件可访问（`aria-label`） | `i18n/messages/*.json`；`jobs-result-center.test.tsx` 的主题用例 |
| 13 | 未触碰 pinned 协议工件，未改 Job 六态合同，未同步化 `batch-delete`，未新增权限键 | `git diff --stat e1893a1a..HEAD -- docs/schemas apps/web/src/protocol/upstream apps/api/modules/jobs/migration …` 为空 |

### 对照检查点

| 检查点 | 状态 | 证据 |
|--------|------|------|
| C1 写操作面 | **达成** | 成果 1～8 |
| C2 结果中心呈现 | **达成** | 成果 8～10、12 |
| C3 体验与测试收敛 | **达成** | 成果 10～12；Go 全绿 + web 117 files / 1457 tests 全绿 + `tsc -b` exit 0 |
| C4 审计与投影 | 进行中 | 本条即 self 腿；independent 腿待跑 |

### 独立复核抽查（编排器对自身实现的复验）

| 主张 | 复核结果 |
|------|---------|
| "逐字镜像既有合同" | **成立**：状态分支、SQL 字段集合、错误码与 actor 版一一对应；`requireAffectedNotRetryable` 与 `requireAffectedForActor` 同形（仅去掉 actor 谓词） |
| "既有隔离未被放松" | **成立**：`repository.go`/`runner.go` 区间零改动；两条新测试反向断言非本人仍 `ErrNotFound` |
| "UI 不会与合同漂移" | **成立**：可用性字段由服务端 `jobToMap` 派生，前端仅消费；Go 侧逐态断言 |
| "download 门控不比服务端更严" | **成立且可判别**：服务端 `GET /api/jobs/{id}/result` 门控为 `jobs.read`，schema 的本地 gate 亦为 `jobs.read`；**变异验证**：改成 `jobs.write` 后只读主体用例变红 |
| "表级权限真的生效" | **成立**：首版把 `permissions` 写进 `props` 导致行操作门禁不生效；改为表节点级 `permissionCascade`+`permissions` 后生效。**变异验证**：移除级联后只读主体用例变红 |
| "写合同门禁是 jobs.write" | **不能由现有测试区分**（见 F-001）：改为 `jobs.read` 后 editor 用例仍绿 |
| "取消状态集合正确" | **成立**：**变异验证**：把 `case StatusRunning` 扩为 `StatusRunning, StatusSucceeded` 后终态用例变红（`cancel succeeded = <nil>`） |
| "全量回归" | **成立**：`go build ./...` exit 0、`go vet ./...` 无输出、`go test ./...` 全绿；web `npm test` 117/1457 全绿、`npm run typecheck` 与 `npm run build` exit 0 |
| "无 unknown custom 警告" | **成立**：web 回归日志中 `unknown custom component` 出现 0 次（10 个渲染型测试已补 side-effect import） |

未发现与证据矛盾的陈述。

### Findings

#### F-001 · 写面门禁测试无法区分 `jobs.write` 与 `jobs.read`

- 严重度：med
- 建议：**recommended**
- 描述：`TestJobActionsRequireJobsWrite` 只证明「无 `jobs.write` 者被拒」；把 cancel 路由的 `requirePermission` 改成 `jobs.read` 后该用例**仍然绿**（本自审实测）。原因是两个权限键当前同为 `PolicyAdmin`，editor 两者皆无，不存在「持 read 而无 write」的 seeded 主体，故**不可构造**判别性反例。若只留现有断言，未来误把写门禁换成读门禁不会被发现。
- 处置：**本区间已加固**（`ab215ebc`）——① handler 新增 `TestJobActionGateIsNotDiscriminableToday` 钉住嵌套不变式（一旦 `jobs.read` 与 `jobs.write` 的持有者集合分离即变红并要求补反例测试）；② web 只读主体用例补**可判别对照**：download 对 `jobs.read` 主体仍可用，改门禁即变红（已变异验证）。
- 证据：`internal/handler/jobs_actions_test.go`；`apps/web/src/renderer/jobs-result-center.test.tsx`；变异记录见 `02-execution/E-001` §2 与本次自审。
- 状态：`fixed`（嵌套守卫 + 可判别对照已落地；真正的 write-vs-read 反例在策略集分离前不可构造，已由守卫显式承接）

#### F-002 · 状态列显示原始状态码，未做逐值本地化

- 严重度：low
- 建议：**recommended**
- 描述：结果中心的状态列显示 `queued`/`running`/… 原始码 + 颜色徽标，列头与筛选器选项才走 i18n。这与既有产品约定一致（`wallet.status`、`scheduledTasks.enabled`、`users.mfaEnabled` 均显示原始值），但中文用户读状态列时仍是英文码。
- 处置：**已定性为决策而非缺陷**（`D-001` §3）：逐值本地化需要协议级 `valueLabels` 类扩展 → 触碰 pinned 工件，属 VP-038 明确非目标；六态词汇已通过筛选器本地化选项对用户可见。若用户要求逐值本地化，应作为后续 VP/协议议题。
- 证据：`modules/jobs/schema/jobs.json`（`badgeStyleField`）；`D-001` §3「状态文案口径」；`i18n` 的 `schema.jobs.status.*` 六键在筛选器中使用。
- 状态：`open`（recommended；已记录为有界决策，等 independent 腿与用户是否接受）

#### F-003 · `jobs-auto-refresh` 依赖 `reloadList()` 清空选择的语义

- 严重度：low
- 建议：**recommended**
- 描述：组件 tick 调用页面级 `reloadList()`，该 seam 会清空**本页所有表选择**（ADR-0022 D2）。当前 `jobs` 表未声明 `props.selection`，无选择可清（测试已断言无选择列），故为惰性；但若将来 jobs 表启用行选择，自动刷新会静默丢弃用户选择。
- 处置：已在 `D-001` §5 与组件注释中登记为 R-1.2（含复核触发条件），并由交互测试钉住「jobs 表无选择」这一前提。
- 证据：`components/jobs-auto-refresh.tsx` 注释；`renderer/render.tsx` `reloadList`（清空 `selections`）；`jobs-result-center.test.tsx` 断言 `input[type=checkbox]` 为 null。
- 状态：`open`（recommended；有界风险 + 已登记复核触发）

#### F-004 · 自动刷新在无进行中作业时仍按固定档位请求

- 严重度：low
- 建议：**recommended**
- 描述：tick 不区分「列表内是否仍有非终态作业」，因此操作员开启 5s 后即使所有作业已终结也会持续请求。这是有意的简化（与 `monitoring-auto-refresh` 同形），且默认为 Off、由操作员显式开启；但存在无谓的后台流量。
- 处置：登记为已知取舍；若 independent 腿或用户要求，可在 R5 或后续波次改为「有非终态行才轮询」。
- 证据：`components/jobs-auto-refresh.tsx`（`OPTIONS` 默认 0 = Off）；`jobs-result-center.test.tsx` 的档位用例（Off 后 30s 零新增请求）。
- 状态：`open`（recommended）

### 必改项汇总（required）

**无。** 未发现 high 级未关闭 required；`I-038-013`～`016` 均已 `verified`。

### 结论 + 建议下一步

C1～C3 交付物**如实、可核对、边界干净**：管理作用域的取消/重试确实按用户裁决逐字镜像既有合同（且既有 actor 隔离以零 diff + 反向断言双重钉住），结果中心在出厂 schema 上真实可用（六态、详情、下载、重试、取消、自动刷新），R3 遗留的交互级测试与 `I-038-013` 已闭合，全量回归绿。

四条 recommended 均为测试可判别性、呈现口径与轮询取舍，不阻断。其中 F-001 已在同一区间以**嵌套守卫 + 可判别对照**加固，并保留「策略集一旦分离即须补反例测试」的显式触发。

**verdict = pass**，C4 的 self 腿通过。**建议下一步**：按项目级决策 [independent-audit-execution.md](../../../../architecture/independent-audit-execution.md) 调用本地 grok build（grok 4.6 · high · `/audit`）执行 independent 腿——**须专门核验**：① 既有 actor 隔离是否真的零退化；② 写门禁的真实性（gate 变异可判别性）；③ 服务端派生字段与写合同是否真同源；④ 前端下载/取消/重试是否真的打到冻结路由；⑤ 是否越界改 pinned 工件或 Job 六态合同。

**本条不修改** `status`、检查点或派生 `progress`。
