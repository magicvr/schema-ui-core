---
id: A-002-r4-c1-c3-independent
doc: audit-entry
parent: GOAL-005-r4-result-center-experience
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# A-002 · R4 C1～C3 独立交叉审计（GOAL-005）

## A-002 · R4 C1～C3 independent（2026-09-19）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high · `/audit`）
- **类型** / **scope**：`stage` · GOAL-005 C1～C3 实施复审——① 新增管理作用域写面是否真的 `jobs.write` 门控、fail-closed、且不越权；② 既有 actor 作用域 `RequestCancel`/`Retry` 是否逐字未放宽；③ 可取消/可重试状态集合与错误码是否与既有合同和 Job 六态合同一致；④ 服务端派生字段是否与写合同同源；⑤ 前端是否真的打到冻结路由、下载是否单一实现、`jobs-auto-refresh` 是否引入未声明副作用；⑥ 前端交互级测试（含 R3 遗留 4 例）是否真的能失败；⑦ 是否越界改 pinned 工件 / Job 六态 / 同步 `batch-delete` / 既有 VP；⑧ 文档登记是否与代码事实一致
- **verdict**：**pass**（0 required；4 recommended）
- **完整意见**：本文件

### 范围与区间

- 被审目标：`docs/workspaces/workspace-038-batch-operations-and-job-center/GOAL-005-r4-result-center-experience/`
- 工作区校验：`workspace.md` `id` = `workspace-038-batch-operations-and-job-center`，`root_goal` = `GOAL-001-batch-operations-and-job-center`，`canonical_scope` 与本区路径一致，`vision_role: delivery`，`plan_refs`/`primary_plan` = `VP-038-batch-operations-and-job-center`。`shared_materials_catalog: none`，本意见未把共享资料当事实或关闭证据。未读取其他工作区目标状态。
- 审计区间：R3 关门 `e1893a1a` → R4 C1 `ea6e4006` → C2 `f1351534` → C3 `bafb766f` → F-001 加固 `ab215ebc` → self 文档 `b44c5ea6`。**不含** C4 本身、Root R4 投影、R5。
- 约束输入：R1 `GOAL-002/01-decision/D-001`（方案 B / 首波 1 条 / K-1～K-7）；R2 `GOAL-003/01-decision/D-001`（`jobs.read`/`jobs.write` 归属、索引形状）；R3 `GOAL-004/01-decision/D-001`（双重门禁、进度、前端口径）；本目标 `D-001` §0 两项用户 P-004 裁决原文。
- 工作树在审计开始时已有**未提交**改动（`apps/web` 两份 i18n 目录 + `jobs-result-center.test.tsx` 多 1 例）。本条把 **HEAD 已提交交付**与 **工作树未登记补丁**分开陈述；前端抽验跑在工作树（含该补丁）。变异只改 `internal/handler/jobs.go` 一处派生字段，已还原，`git diff` 对该文件为空。

### 范围与区间 · P-005

| ID | 级别 | 00-meta | 最晚阶段 | 本 scope 判定 |
|----|------|---------|----------|----------------|
| I-038-014 | required | **verified**（方案 A） | C2 前 | 用户书面裁决已落 `D-001` §0/§3；实施为既有 `jobs` 页原地增强，无新页面/导航 |
| I-038-015 | required | **verified**（逐字镜像） | C1 前 | 用户书面裁决已落 `D-001` §0/§1；代码与测试与裁决一致 |
| I-038-016 | non-blocking | **verified** | R4 前 | `downloadable = succeeded` 关闭过期下载入口；410 仍为直接访问兜底 |
| I-038-013 | non-blocking | **verified** | C3 | 单一实现 `lib/job-result-download.ts` + 4 例守卫 |
| I-038-010 | non-blocking | 00-meta 称 C3 承接 | C3 | 双目录键集合对齐测试本会话绿；不阻断 |
| I-038-006 | deferred · non-blocking | 保持 | — | 不在本 scope |

无到期且影响本 scope 的 required 信息项处于未关闭状态。`accepted-residual` 不适用。`01-decision.md` 索引表仍把 I-038-013～016 写成 `open`、D-ID 写成「暂无」——与 `00-meta` / 磁盘上的 `D-001` 不一致，见 F-004。

### 成果（有证据）

| # | 主张 | 独立核验 |
|---|------|----------|
| 1 | **既有 actor 隔离零退化**（重点） | `git diff --stat e1893a1a..HEAD -- apps/api/internal/jobs/repository.go apps/api/internal/jobs/runner.go` **为空**。`RequestCancel` 仍 `getForActorTx` + SQL `AND actor_id=?`；`Retry` 仍 `AND actor_id=?` + `requireAffectedForActor`。`Runner.Cancel`/`Retry` 仍走这两条，**不**调用 `*Any`。钱包写面仍 `jobService.Cancel/Retry(..., user.ID)` + `wallet.write`。全仓生产调用 `RequestCancelAny`/`RetryAny`/`CancelAny`/`RetryAny` 仅 handler 管理写面、composition 绑定与测试。反向钉：`TestRequestCancelAnyCancelsAnotherActorsQueuedJob` / `TestRetryAnyRequeuesAnotherActorsFailedJob` 断言非本人 `ErrNotFound` |
| 2 | **写门禁真实 fail-closed**（重点） | `POST /api/jobs/{id}/cancel` 与 `/retry` 均 `requirePermission(w, r, "jobs.write")`，`!ok` 则 **return**（`jobs.go:78-80, 89-91`）。`requirePermission`：无身份 401，权限集不含该键 403。被拒请求不触及状态机（`TestJobActionsRequireJobsWrite`）。actions 为 nil 时路由不挂载（`TestJobActionRoutesAbsentWithoutActions`）。**不能**由 editor 403 区分 `jobs.write` vs `jobs.read`（两键同 `PolicyAdmin`）——见 F-001 |
| 3 | **状态集合与错误码逐字镜像，未改合同**（重点） | `RequestCancelAny`：`queued` → 立即 `cancelled`（置 `finished_at`，`cancel_requested=0`）；`running` → 只置 `cancel_requested=1`，不置 `finished_at`，由既有 `FinalizeCancel` 收尾；其余 `ErrNotCancellable`。`RetryAny`：`status='failed' AND attempt < max_attempts`，**不**重置 attempt；其余 `ErrNotRetryable`。HTTP 映射仅 `JOB_NOT_FOUND` / `JOB_NOT_CANCELLABLE` / `JOB_NOT_RETRYABLE`（`error_contract_test.go` 既有冻结码，**未新增**）。本会话 `go test ./internal/jobs/ -run Any -count=1` **全绿** |
| 4 | **派生字段与写合同同源（规则同一，落字两处）** | `jobToMap`：`cancellable = queued\|running`；`retryable = failed && attempt < maxAttempts`；`downloadable = succeeded`。与 `actions.go` 前置条件同一条规则，但是布尔复述而非共享函数。`TestJobsProjectionDerivesActionAvailability` 含耗尽 attempt 行。本会话**独立变异**：去掉 `attempt < maxAttempts` 后该用例变红（`retryable:true` 且 `attempt=3/3`）；还原后变绿，`jobs.go` diff 为空。`downloadable=succeeded` 与结果路由三段语义一致：queued/running → 409、expired → 410、succeeded → 200 附件；failed/cancelled 的 GET 仍 200 投影但 UI 不提供下载（与 D-001 §2/§4 一致） |
| 5 | **投影为追加，不破坏 R2 冻结面** | `jobToMap` 保留 `id/kind/status/progress/attempt/maxAttempts/resultUrl/error` 等原字段；新增 `cancellable/retryable/downloadable/statusStyle/errorCode/errorMessage`。`resultUrl` 仍仅 succeeded 且走 `jobs.ResultURL`。读路由仍 `jobs.read` |
| 6 | **前端真实打到冻结路由，下载单一实现**（重点） | 测试读取出厂 `apps/api/modules/jobs/schema/jobs.json`（非手写替身）。`cancelJob`/`retryJob` URL 为 `POST /api/jobs/{id}/cancel`/`/retry`；交互测试断言精确 URL `/api/jobs/job-running/cancel`、`/api/jobs/job-failed/retry`。`click()` 对缺失按钮 `not.toBeNull`，故「无 POST」**不是**因为没渲染出行操作。只读主体用例先断言 Cancel/Retry **存在且 disabled**，再派发 click，再断言 POST 为空；download 对 `jobs.read` 仍可用。`jobs.downloadResult` 白名单 URL `GET /api/jobs/{id}/result`，经 `downloadJobResultDocument`；R3 组件同一助手 |
| 7 | **R3 遗留 4 例可判别** | 空选择：按钮 disabled 且 `/api/jobs/batch-export` 请求为空。只接受 202：组件 `response.status !== 202`；用例喂 200，断言无 progress、`jobListCalls()=0`、出现错误文案——**200 被排除在轮询之外**。不 `reloadList()`：`reloadList` 会 `setSelections({})` **且** bump `reloadToken` 强制重拉；用例同时钉 `/api/users` 次数 = 1 与按钮仍含 `(2)`。即使列表请求被合并，选择被清仍会红。终态下载走共享助手与服务端 `fileName` |
| 8 | **自动刷新副作用在声明范围内** | `jobs-auto-refresh.tsx` tick **只**调 `crud?.reloadList()`；`useEffect` 在 `intervalMs<=0` 时不建定时器，卸载/`intervalMs` 变化时 `clearInterval`。无全局监听。默认 Off。`jobs.json` 表节点**无** `props.selection`；交互测试断言无 checkbox。R-1.2 前提成立 |
| 9 | **schema 合法；CustomAction 无 `onSuccess`** | `jobs.json`：表节点 `permissionCascade.keys=[edit]` + `permissions.edit = jobs.write`；download 行动作本地 `permissions.edit = jobs.read`（不严于服务端）；`downloadJobResult` 无 `onSuccess`。`badgeStyleField` 为仓库既有本地扩展。本会话 `all-module-schemas-dval` 39、`capability-declaration.guard` 36、`schema-keys.structural` 4 **全绿** |
| 10 | **未越界** | 四 checkpoint `--stat` 只含 `apps/api/**` 与 `apps/web/**`。`git diff --stat e1893a1a..HEAD -- docs/schemas apps/web/src/protocol/upstream apps/api/modules/jobs/migration` **为空**。未改 Job 六态字面量、未改同步 `batch-delete`、未新增权限键（仍 `jobs.read` + `jobs.write`）、未重开既有 VP。Root 仍 `active · 3/5`，R4 检查点未勾 |
| 11 | **回归本会话独立复跑** | 见下表。C1～C3 写面/呈现/交互主张在抽验范围内成立 |

### 对照检查点

| 检查点 | 状态 | 证据 |
|--------|------|------|
| C1 写操作面 | **达成** | 成果 1～5、10；`Any` 用例与 handler 写面测试本会话绿 |
| C2 结果中心呈现 | **达成** | 成果 4～6、8、9；出厂 `jobs.json` + 交互测试 |
| C3 体验与测试收敛 | **达成**（有 recommended 缺口） | 成果 6、7、9、11；前端 error 目录缺口见 F-002 |
| C4 审计与投影 | 进行中 | 本条为 independent 腿；不改 status/progress |

### 回归证据可信度

本会话**独立复跑**（均 `-count=1`）：

| 命令 | 结果 |
|------|------|
| `cd apps/api && go test ./internal/jobs/ -run Any -count=1 -v` | **PASS**（8 个顶层用例，含四终态拒绝、attempt 保留、runner 停在飞 handler） |
| `cd apps/api && go test ./internal/handler/ -run TestJobsProjectionDerivesActionAvailability -count=1 -v` | 基线 **PASS**；变异后 **FAIL**（见下）；还原后 **PASS** |
| `cd apps/api && go test ./... -count=1` | **全绿**（含 `internal/handler` 50.9s、`internal/store` 65.8s、`internal/composition` 32.9s、`internal/jobs` 2.8s、`modules/jobs` 0.5s） |
| `cd apps/web && npm test -- src/renderer/jobs-result-center.test.tsx src/components/jobs-batch-export.test.tsx src/lib/job-result-download.test.ts src/protocol/all-module-schemas-dval.test.ts src/protocol/capability-declaration.guard.test.ts src/i18n/schema-keys.structural.test.ts` | **6 files / 97 tests 全绿**（result-center **10** 例含工作树未提交的本地化例；batch-export 4；download 4；dval 39；capability 36；schema-keys 4） |

未在本会话复跑 web 全量 1455/1457 与 `tsc -b` / `vite build`。E-001 / A-001 的全量数字标为「执行/自审自称」；本条独立复证的是 Go 全量 + 任务书点名的 web 抽验。

**独立变异（本会话，非复述 self）**：把 `jobToMap` 的 `retryable` 改为忽略 `max_attempts`（`job.Status == jobs.StatusFailed`）。`TestJobsProjectionDerivesActionAvailability` 变红：`an exhausted failure still advertises an action`，`retryable:true` 且 `attempt=3 maxAttempts=3`。还原后该用例绿，`git diff -- apps/api/internal/handler/jobs.go` 为空。结论：投影测试对「可重试集合被放宽」**可判别**。

抽验前端断言（非恒真）：取消用例 `posted.map(url) === ["/api/jobs/job-running/cancel"]`；只读用例先断言按钮存在再断言 POST 空；202 用例喂 200 后 `jobListCalls()===0`。

### Findings

#### F-001 · HTTP 写门禁测试仍无法区分 `jobs.write` 与 `jobs.read`

- 严重度：med
- 建议：**recommended**
- 描述：与 self A-001 F-001 同因：`TestJobActionsRequireJobsWrite` 只证明「无 `jobs.write` 者被拒」。两键均为 `PolicyAdmin`，editor 两者皆无；把路由门禁改成 `jobs.read` 后该用例仍绿。self 已加 `TestJobActionGateIsNotDiscriminableToday`（策略集一旦分离即 fail closed 要求补反例）+ web 只读主体对 download 的可判别对照。
- 独立评估：① **静态第三条路径成立**：handler 源码逐字 `requirePermission(..., "jobs.write")` 且 `!ok` 则 `return`，本条已核对；描述符只声明模块级权限集合，**不**按路由钉死门禁键。② **测试内可构造真反例**：身份来自 `role_permissions`，测试不必受 seeded 角色限制，可向身份库写入「只持 `jobs.read`」的授予，或在请求上下文注入该权限集——self 未走这条。③ 嵌套守卫是**触发器**不是判别器；web 对照钉的是 **UI download 门控**（`jobs.read`），不是 HTTP cancel/retry 门禁键。④ 去掉门禁（而非换键）会被现有 editor 403 抓住。
- 处置建议：维持嵌套守卫即可；策略集分离前不阻断。若要一次钉死键名，用自定义授予构造 `jobs.read`-only 主体打 POST cancel，期望 403 且行不变。
- 证据：`apps/api/internal/handler/jobs.go:78-91`；`jobs_actions_test.go` `TestJobActionsRequireJobsWrite` / `TestJobActionGateIsNotDiscriminableToday`；`apps/web/src/renderer/jobs-result-center.test.tsx` 只读主体例。
- 状态：`open`（recommended；与 self F-001 残余同向，不升级为 required）

#### F-002 · 已提交的前端目录缺少 R4 写面将下发的 `error.job*` messageKey

- 严重度：med
- 建议：**recommended**
- 描述：服务端 `errorcatalog` 已为 `JOB_NOT_FOUND` / `JOB_NOT_CANCELLABLE` / `JOB_NOT_RETRYABLE` / `JOB_RESULT_NOT_READY` / `JOB_RESULT_EXPIRED` / `JOB_ATTEMPTS_EXHAUSTED` / `JOB_HANDLER_FAILED` 配置 `error.job*` messageKey。R4 把这些码第一次送到 **jobs 结果中心**写操作反馈。`git show HEAD:apps/web/src/i18n/messages/en-US.json` **不含**这些键。客户端优先目录、未知 key 回退服务端字符串并记 missing-translation。服务端目录本身有中英文，故不是「完全无文案」，但 C3「中英文」在客户端目录上对这条新写面不闭合。
- 工作树（**未提交、未进 E-001 checkpoint**）已补 7 键 + `it("localizes a refused transition through the API message key")`。该例喂 409 + `messageKey: error.jobNotCancellable`，断言 UI 出现目录文案「this job can no longer be cancelled」且**不**出现服务端英文「job cannot be cancelled」——**可判别**。本会话抽验含此例故 10/10 绿；HEAD 仅 9 例。
- 处置建议：由 `/govern` 把该补丁纳入登记并提交，或书面接受「写面错误走服务端目录回退」。
- 证据：`apps/api/internal/errorcatalog/errorcatalog.go:163-169`；HEAD vs 工作树 `apps/web/src/i18n/messages/{en-US,zh-CN}.json`；工作树 `jobs-result-center.test.tsx` 新增例。
- 状态：`open`（recommended）

#### F-003 · 前端夹具的 `retryable` 未编码 attempt 预算

- 严重度：low
- 建议：**recommended**
- 描述：`jobs-result-center.test.tsx` 的 `row()` 令 `retryable: overrides.status === "failed"`，忽略 `attempt < maxAttempts`；种子行也无「失败且预算耗尽」行。服务端投影测试覆盖了耗尽态（本会话变异已证明可红）。若 schema 的 `disabledWhen` 被改成「status==failed」而忽略服务端 `retryable`，前端交互测试仍会绿。
- 处置建议：加一行 `status=failed, attempt=maxAttempts, retryable=false`，断言 Retry disabled。
- 证据：`apps/web/src/renderer/jobs-result-center.test.tsx` `row()` / `ROWS`；对照 `TestJobsProjectionDerivesActionAvailability`。
- 状态：`open`（recommended）

#### F-004 · 文档索引与执行登记未完全跟上代码事实

- 严重度：low
- 建议：**recommended**
- 描述：
  1. `01-decision.md` 信息表仍写 I-038-013～016 `open`、决策索引「暂无」；磁盘已有 `D-001`，`00-meta` 已 `verified`。
  2. `00-meta` 「审计意见状态」仍写「尚无审计条目」，但 `A-001` 已落盘；备注仍写 `progress: 0/4`，frontmatter 为 `3/4`。
  3. `E-001` §5 checkpoint 只列 `ea6e4006` / `f1351534` / `bafb766f`，未列 F-001 加固 `ab215ebc`（A-001 已引用）。
  4. 全量 web 测试计数：`00-meta` C3 / `E-001` 写 1455，A-001 写 1457。本条未复跑全量，**不采信任一侧为独立事实**。
  5. `D-001` §3.4 写三行操作「均有 `requestMapping.path.id`」；`jobs.json` 的 download **没有**该字段（由 custom handler `{id}` 槽解析，同节后半句已说明）。小口径漂移。
- 这些不推翻 C1～C3 代码事实；Root / goal-tree / workspace 的 `GOAL-005 active · 3/4`、Root `active · 3/5`、R4 检查点未勾 **与事实一致**。
- 证据：上述文件相对 `D-001` / `A-001` / `ab215ebc`。
- 状态：`open`（recommended）

### 必改项汇总（required）

**无。** 未发现 high 级未关闭 required；影响本 scope 的 required 信息项 I-038-014 / I-038-015 均为 `verified`。

### 与既有意见的异同（self A-001）

| 项 | self A-001 | 本条 independent |
|----|------------|------------------|
| verdict | pass（0 required + 4 recommended） | **pass**（0 required + 4 recommended） |
| 隔离零退化 | 空 diff + 反向断言 | **同意并复证**；并核钱包路径未改走 `*Any` |
| 写门禁 F-001 | recommended，已嵌套守卫，标 `fixed` | **同意残余**：换键仍不可由 editor 用例抓住。独立补充「身份库可构造真反例」与「静态源码是第三条核对路径」。不升级 required，本条 F-001 保持 open recommended |
| 状态集合 / 错误码 | 变异扩宽 cancel 集合变红 | **同意**；本条另变异 `retryable` 忽略预算，投影测试变红后已还原 |
| 派生字段同源 | 成立 | **同意**（规则同一、落字两处，测试可判别） |
| 前端路由 / 下载单源 | 成立 | **同意**；并核「无 POST」断言非因未渲染按钮 |
| R3 遗留可判别性 | 成立 | **同意**；并核 `reloadList` 必清选择 + 必重拉，双断言不会因缓存漏判；202 用例真排除 200 |
| self F-002 状态列原始码 | recommended，已定性为决策 | **同意**，不另开号 |
| self F-003 刷新清选择 | recommended，jobs 表无选择 | **同意**（前提本会话复证） |
| self F-004 无进行中仍轮询 | recommended 取舍 | **同意**，不另开号 |
| 前端 error.* 目录 | 未报 | **本条新增 F-002**（HEAD 缺键；工作树有未提交补丁） |
| 前端夹具 attempt 预算 | 未报 | **本条新增 F-003** |
| 文档索引漂移 | 未报 | **本条新增 F-004** |
| 越界 / pinned | 空 diff | **同意并复证** |

无与 self 在 required / 必改项上的冲突，不触发 P-004 冲突裁决。

### 结论 + 建议给编排器/用户的下一步

C1～C3 的关键安全与合同主张**独立成立**：管理作用域 cancel/retry 确为 `jobs.write` 且 fail-closed；既有 actor 隔离以空 diff + 非间接调用 + 反向测试钉住，未被 `*Any` 稀释；状态集合与冻结错误码未扩展；派生字段与写合同同一规则且变异可判别；前端打到冻结路由、下载单源、R3 遗留交互测试可红；未触 pinned / 六态 / `batch-delete`。

四条 recommended 均为门禁可判别性残余、客户端错误目录未提交、夹具覆盖与文档索引，**不阻断** C4 在开放 required = 0 条件下响应本条后投影。

**建议下一步**（`/govern`）：

1. 响应本条：F-001 维持残余或补 `jobs.read`-only HTTP 反例；F-002 提交工作树 i18n+测试或书面接受服务端回退；F-003/F-004 可随手修或 residual。
2. 闭合后若开放 required 仍为 0，再勾 C4 并投影 Root R4（当前 Root 仍 `3/5`，**不应**在 C4 前勾选）。

**本条不修改** `status`、检查点或派生 `progress`，不改 goal-tree，不改方案正文，不改 `apps/**` 实现（变异已还原）。

### 声明

本意见不修改 status/progress；响应由 /govern 处理。
