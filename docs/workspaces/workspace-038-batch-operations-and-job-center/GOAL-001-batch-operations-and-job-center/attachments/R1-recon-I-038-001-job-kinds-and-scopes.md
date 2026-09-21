---
title: R1 侦察 · I-038-001 Job 种类与可见作用域分母
status: draft
created: 2026-09-19
updated: 2026-09-19
parent: null
version: 0.1.0
---

# R1 侦察 · Job 种类与可见作用域分母

> **性质**：只读代码侦察（READ-ONLY reconnaissance），服务 `I-038-001`（required，R1 前关闭）。
> 本文件只记录**已读代码的事实**与 file:line 证据；不含方案决策。未读到/无法判定的内容集中在末尾 `## 待确认 / 未知`。
> 侦察时间基线：工作树 `apps/**` 无未提交变更（同 `D-001` §2 freshness 记录）。

---

## 0. 结论速览（TL;DR）

| 问题 | 事实 |
|------|------|
| 全仓已注册 Job 种类数 | **1**（`wallet.reconcile`）。测试专用 `panic.kind` 不属生产面 |
| 生产 Submit 调用点数 | **1**（`JobService.SubmitReconcile`） |
| 暴露 Job 的 HTTP 模块 | **1**（`admin.wallet`），4 条路由 + 1 条提交路由，**无列表路由** |
| 作用域模型 | `actor_id` 等值 + `kind` 等值（"本 actor + 本 kind"） |
| 跨 actor / 管理面读路径 | **不存在**（全仓零命中） |
| 通用列表查询能力 | **不存在**。Repository 唯一列表方法是调度器用的 `ListRunnable`（无过滤/无 offset/无 total） |
| Job 专用权限字符串 | **无**。复用 `wallet.read` / `wallet.write` |
| 可支撑 admin 全量列表的索引 | **无**（`idx_jobs_actor` 以 `actor_id` 打头，跨 actor `ORDER BY updated_at DESC` 不被覆盖） |
| 同类 admin 跨 actor 列表先例 | **有**：`admin.activity`（operations）与 `admin.scheduled-tasks`（task-runs 全局历史） |

---

## 1. JOB KINDS

### 1.1 注册机制（先看契约）

`Runner` 的 kind 注册表是**内存 map，键为裸字符串**，没有任何 kind 常量枚举、元数据表或清单：

- `apps/api/internal/jobs/runner.go:62` — `handlers map[string]registration`
- `apps/api/internal/jobs/runner.go:89` — `func (r *Runner) Register(kind string, handler Handler) error`
- `apps/api/internal/jobs/runner.go:95` — `func (r *Runner) RegisterWithTerminalHook(kind string, handler Handler, terminal TerminalHook) error`
- `apps/api/internal/jobs/runner.go:96-98` — 只校验 `kind == "" || handler == nil` → `ErrInvalid`
- `apps/api/internal/jobs/runner.go:104-106` — 重复 kind 报错：`return fmt.Errorf("jobs: handler %q already registered", kind)`
- `apps/api/internal/jobs/runner.go:101-103` — `if r.started { return errors.New("jobs: handlers must be registered before start") }`
- `apps/api/internal/jobs/runner.go:269` — 派发时按 `r.handlers[job.Kind]` 查表；未注册 → `Fail(..., "JOB_HANDLER_FAILED", "job handler not registered", ...)`（runner.go:271-273）

**含义（事实性）**：kind 的"分母"只能靠扫描 `Register*` 调用点得出；运行期没有 `SELECT DISTINCT kind` 之外的枚举手段，也没有任何 kind → 展示名 / payload schema / 描述 / 归属模块的元数据注册点。

`ReconcileJobKind` 是本仓唯一以常量形式声明的 kind；没有名为 `*JobKind` 的其他常量（全仓 grep `JobKind` 仅命中 `ReconcileJobKind` 与测试）。

### 1.2 生产 kind 表（完整）

| 项 | 值 | 证据 |
|----|----|------|
| **kind 字符串** | `wallet.reconcile` | `apps/api/modules/wallet/jobs.go:19` `const ReconcileJobKind = "wallet.reconcile"` |
| **声明包/文件** | `package wallet` · `apps/api/modules/wallet/jobs.go` | 同上 |
| **注册调用** | `runner.RegisterWithTerminalHook(ReconcileJobKind, s.runReconcile, s.recordTerminal)` | `apps/api/modules/wallet/jobs.go:39` |
| **handler 函数** | `func (s *JobService) runReconcile(ctx, job jobs.Job, reporter jobs.Reporter) (jobs.CommitFunc, error)` | `apps/api/modules/wallet/jobs.go:116` |
| **terminal hook** | `func (s *JobService) recordTerminal(job jobs.Job)` | `apps/api/modules/wallet/jobs.go:162` |
| **谁调用 Submit** | `func (s *JobService) SubmitReconcile(ctx, accountID string, actor account.User, correlationID string) (*jobs.Job, error)` | `apps/api/modules/wallet/jobs.go:45` |
| ↳ Submit 实际调用 | `s.runner.Submit(ctx, jobs.CreateInput{ID: id, Kind: ReconcileJobKind, Payload: payload, ActorID: actor.ID, CorrelationID: correlationID, MaxAttempts: jobs.DefaultMaxAttempts, Now: now})` | `apps/api/modules/wallet/jobs.go:69-72` |
| ↳ 上游 HTTP 调用点 | `job, err := jobService.SubmitReconcile(r.Context(), strings.TrimSpace(body.AccountID), user, correlationID)` | `apps/api/internal/handler/wallet.go:379`（`POST /api/wallet/reconcile`） |
| **装配点** | `walletJobs, err := walletmodule.NewJobService(walletService, jobRuntime.repository, jobRuntime.runner, operations)` | `apps/api/internal/composition/composition.go:584` |
| **payload 形状** | `type reconcileJobPayload struct { AccountID string \`json:"accountId,omitempty"\` }` | `apps/api/modules/wallet/jobs.go:21-23`；序列化于 `:65` |
| **result 形状** | `reconciliationResult(run walletstore.ReconciliationRun) map[string]any` → `{id, accountId, result, mismatchCount, details, actorId, createdAt}` | `apps/api/modules/wallet/jobs.go:214-221`；落盘于 `:158` `return json.Marshal(reconciliationResult(*run))` |
| **MaxAttempts** | `jobs.DefaultMaxAttempts` = **3** | 传入 `jobs.go:71`；常量 `apps/api/internal/jobs/model.go:27` `const DefaultMaxAttempts = 3` |
| **result TTL** | **24h**（Runner 全局选项，非按 kind） | `apps/api/internal/jobs/runner.go:51` `ResultTTL: 24 * time.Hour`（`DefaultRunnerOptions()`） |
| **Job ID 生成** | `jobs.NewID(now)` = 16 位 hex 毫秒 + 12 字节随机 | `apps/api/modules/wallet/jobs.go:47`；`apps/api/internal/jobs/model.go:78-84` |

**Job ID 复用为业务 run ID**：`s.service.ReconcileOnceTx(..., payload.AccountID, job.ID, job.ActorID, ...)`（`apps/api/modules/wallet/jobs.go:128`）——同一 ID 既是 job id 也是 `wallet_reconciliation_runs.id`。这是 `GET /api/wallet/reconcile/runs` 与 Job result 能对应的原因。

### 1.3 非生产 kind（仅测试，供完整性）

| kind | 位置 | 说明 |
|------|------|------|
| `wallet.reconcile`（重复注册） | `apps/api/internal/handler/wallet_test.go:164`；`apps/api/internal/jobs/runner_test.go:29,109,242,274`；`shutdown_reclaim_test.go:27` | 测试替身 runner |
| `panic.kind` | `apps/api/internal/jobs/runner_panic_test.go:30,42` | 仅验证 panic → 持久化失败 |

### 1.4 六态与状态常量（分母的另一半）

`apps/api/internal/jobs/model.go:18-25`：`StatusQueued/Running/Succeeded/Failed/Cancelled/Expired` = `queued/running/succeeded/failed/cancelled/expired`。
DB 层 CHECK 同集合：`apps/api/modules/jobs/migration/migration.go:18`。

### 1.5 派发/结果写入路径（影响"列表可见性"的事实）

- 成功：`CompleteWithCommit(ctx, lease, now, r.options.ResultTTL, outcome.commit)`（`runner.go:363`）→ `status='succeeded', progress=100, ..., expires_at=?`（`repository.go:202-207`）。
- 过期：`ExpireDue`（`repository.go:254-257`）`status='succeeded' AND expires_at <= ?` → `status='expired', result=NULL`。
- 惰性过期读：`JobService.Get` 在 `StatusSucceeded` 时先 `ExpireIfDue` 再重读（`apps/api/modules/wallet/jobs.go:86-91`）。
- **含义**：`expired` 只在被读或被扫描时才产生，列表实现若不做等价惰性/批量过期，会看到"已过 TTL 但仍显示 succeeded 且 result 可下载"的行。

---

## 2. SUBMIT/READ PATHS TODAY

### 2.1 路由总表（模块 `admin.wallet`）

模块 id 声明：`apps/api/modules/wallet/provider.go:29` `const ModuleID = "admin.wallet"`。
路由贡献清单：`provider.go:262-281`（其中 Job 相关见 `:273-274`）。
挂载：`provider.go:294-299` `for _, route := range handler.WalletRoutes(p.a, p.service, p.jobs, p.operations, ModuleID, p.ownerExists)`。
kernel profile 描述符：`apps/api/kernel/profile.go:202`（`{ID: "admin.wallet", ...}`，Routes 含四条 Job 路由）。

| # | 方法 + 路径 | handler 位置 | 权限 | 成功响应 | 错误 |
|---|-------------|--------------|------|----------|------|
| 1 | `POST /api/wallet/reconcile` | `apps/api/internal/handler/wallet.go:357-385` | `wallet.write`（`:358`） | **202** + `walletJobToMap(*job)`（`:384`） | 400 `INVALID_BODY`（`:372`） |
| 2 | `GET /api/wallet/jobs/{id}` | `wallet.go:387-398` | `wallet.read`（`:388`） | **200** + `walletJobToMap(*job)`（`:397`） | 404 `JOB_NOT_FOUND` |
| 3 | `POST /api/wallet/jobs/{id}/cancel` | `wallet.go:400-411` | `wallet.write`（`:401`） | **200** + `walletJobToMap(*job)`（`:410`） | 409 `JOB_NOT_CANCELLABLE` |
| 4 | `POST /api/wallet/jobs/{id}/retry` | `wallet.go:413-424` | `wallet.write`（`:414`） | **200** + `walletJobToMap(*job)`（`:423`） | 409 `JOB_NOT_RETRYABLE` |
| 5 | `GET /api/wallet/jobs/{id}/result` | `wallet.go:426-451` | `wallet.read`（`:427`） | succeeded → **200** 原始 `job.Result` 字节 + `Content-Disposition: attachment; filename="wallet-reconcile-<id>.json"`（`:444-447`）；failed/cancelled → 200 + `walletJobToMap`（`:442`） | 409 `JOB_RESULT_NOT_READY`（`:438`）；410 `JOB_RESULT_EXPIRED`（`:440`） |

**不存在任何 Job 列表路由**（全仓无 `GET /api/jobs`、无 `/api/wallet/jobs` 集合路由）。

### 2.2 响应映射函数（原文引用）

`apps/api/internal/handler/wallet.go:989-1006`：

```go
func walletJobToMap(job jobs.Job) map[string]any {
	row := map[string]any{
		"id": job.ID, "kind": job.Kind, "status": job.Status,
		"progress": job.Progress, "attempt": job.Attempt, "maxAttempts": job.MaxAttempts,
		"cancelRequested": job.CancelRequested, "createdAt": job.CreatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
		"updatedAt": job.UpdatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
	}
	if job.ErrorCode != "" {
		row["error"] = map[string]any{"code": job.ErrorCode, "message": job.ErrorMessage}
	}
	if job.FinishedAt != nil {
		row["finishedAt"] = job.FinishedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00")
	}
	if job.Status == jobs.StatusSucceeded {
		row["resultUrl"] = "/api/wallet/jobs/" + job.ID + "/result"
	}
	return row
}
```

**关键事实**：
- 该映射**故意不内嵌 `payload` / `result`**（`A-012-r4-s4-closeout-independent.md:98` 记录了此冻结口径："GET Job 不内嵌 payload/result；succeeded 只给 `resultUrl`"）。
- `resultUrl` 是**硬编码 wallet 前缀**，不是通用路径——任何通用作业中心要么复用该字段（会把用户导向 wallet 路由），要么改契约。
- 映射函数是 `handler` 包私有函数，**不可被其他模块复用**。

### 2.3 错误码映射

`apps/api/internal/handler/wallet.go:903-914`：

```go
func writeWalletJobError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, jobs.ErrNotFound):
		writeLocalizedError(w, r, http.StatusNotFound, "JOB_NOT_FOUND", "job not found")
	case errors.Is(err, jobs.ErrNotCancellable):
		writeLocalizedError(w, r, http.StatusConflict, "JOB_NOT_CANCELLABLE", "job cannot be cancelled")
	case errors.Is(err, jobs.ErrNotRetryable):
		writeLocalizedError(w, r, http.StatusConflict, "JOB_NOT_RETRYABLE", "job cannot be retried")
	default:
		writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "job operation failed")
	}
}
```

### 2.4 Job 服务接口（actor 边界）

`apps/api/internal/handler/wallet.go:58-65`：

```go
// WalletJobService is the actor-scoped async boundary consumed by wallet
// reconciliation routes.
type WalletJobService interface {
	SubmitReconcile(ctx context.Context, accountID string, actor account.User, correlationID string) (*jobs.Job, error)
	Get(ctx context.Context, id, actorID string) (*jobs.Job, error)
	Cancel(ctx context.Context, id, actorID string) (*jobs.Job, error)
	Retry(ctx context.Context, id, actorID string) (*jobs.Job, error)
}
```

注意接口**没有 List**，且每个方法都以 `actorID` 为参数——作用域是接口层面的显式契约，不只是 SQL 细节。

---

## 3. SCOPING

### 3.1 Repository 三个 actor 相关方法的精确谓词

| 方法 | 位置 | 谓词（原文） | 作用域语义 |
|------|------|--------------|-----------|
| `GetForActor(ctx, id, kind, actorID)` | `apps/api/internal/jobs/repository.go:66-74` | ``SELECT ... FROM jobs WHERE id=? AND kind=? AND actor_id=?``（`:70`） | **id + kind + actor 三重等值** |
| `RequestCancel(ctx, id, actorID, now)` | `repository.go:124-149` | 先 `getForActorTx`（`:127`）；`getForActorTx` = ``WHERE id=? AND actor_id=?``（`:468-470`）；两条 UPDATE 再带 `AND actor_id=?`（`:134`、`:138`） | **id + actor**（**不含 kind**） |
| `Retry(ctx, id, actorID, now)` | `repository.go:220-238` | ``UPDATE jobs SET status='queued', ... WHERE id=? AND actor_id=? AND status='failed' AND attempt < max_attempts``（`:223-227`）；失败复核走 `requireAffectedForActor`（`:231`）→ `getForActorTx`（`:450-462`） | **id + actor + status='failed' + attempt<max** |

`getForActorTx`（`repository.go:468-470`）：

```go
func getForActorTx(ctx context.Context, tx kernel.Tx, id, actorID string) (*Job, error) {
	return scanJob(tx.QueryRow(ctx, `SELECT `+jobColumns+` FROM jobs WHERE id=? AND actor_id=?`, id, actorID))
}
```

**kind 约束的落点**：`GetForActor` 在 SQL 里；`RequestCancel` / `Retry` 在**服务层前置校验**——`apps/api/modules/wallet/jobs.go:96` 与 `:110` 都先做 `s.repository.GetForActor(ctx, id, ReconcileJobKind, actorID)`，再调 `s.runner.Cancel/Retry`。

### 3.2 作用域回答

- **是"同 actor"（`actor_id` 等值）**：三个方法都以 `actor_id=?` 收口。
- **是"同 kind"**：`GetForActor` 在 SQL 内；cancel/retry 由 wallet 服务层补上（repository 层本身不校验 kind）。
- **不是"同 module"**：`jobs` 表**没有 module 列**（见 `migration.go:15-44` 列清单）。模块归属只能由 kind 前缀（`wallet.`）或代码约定推断，DB 层不可查询。
- 跨 actor 或 kind 不匹配统一收敛为 `jobs.ErrNotFound` → HTTP 404 `JOB_NOT_FOUND`（`repository.go:109-111` 把 `kernel.ErrNoRows` 映射为 `ErrNotFound`；`writeWalletJobError` `wallet.go:905-906`）。

### 3.3 是否存在任何 admin/管理面（跨 actor）读路径

**不存在。** 依据：

- 全仓 `FROM jobs` 命中仅 4 处，全部在 jobs 包内部或测试：
  - `apps/api/internal/jobs/repository.go:70, 266, 282, 298, 325, 465, 469`（`GetForActor` / 维护扫描 / `ListRunnable` / `IsCancelRequested` / `getTx` / `getForActorTx`）
  - `apps/api/internal/store/postgres_test.go:844` — `SELECT count(*) FROM jobs`（测试断言）
- `Repository.Get(ctx, id)`（`repository.go:56-64`）**不带 actor 过滤**，但它是 runner 内部读（`Claim`/`CompleteWithCommit`/`transitionJobs` 复用 `getTx`），**没有任何 HTTP 路由暴露它**。
- `ListRunnable`（`repository.go:292-316`）是跨 actor 的，但它是**调度器队列扫描**（只取 queued/可抢占 running），不是管理读面，且不暴露 HTTP。

### 3.4 权限门（route guard / permission string / capability）

- 门函数：`apps/api/internal/handler/resources.go:322-333`

```go
func requirePermission(w http.ResponseWriter, r *http.Request, permission string) (account.User, bool) {
	user, ok := auth.IdentityFrom(r.Context())
	if !ok {
		writeLocalizedError(w, r, http.StatusUnauthorized, "UNAUTHENTICATED", "no active session")
		return account.User{}, false
	}
	if !slices.Contains(user.Permissions, permission) {
		writeLocalizedError(w, r, http.StatusForbidden, "FORBIDDEN", "permission required: "+permission)
		return account.User{}, false
	}
	return user, true
}
```

- Job 路由的权限字符串：**无 Job 专用键**，复用 `wallet.read`（GET）与 `wallet.write`（提交/取消/重试）——`wallet.go:358, 388, 401, 414, 427`。
- 权限声明：`apps/api/modules/wallet/provider.go:341-350`（四条 `kernel.PermissionContribution`，`PolicyID: authsessiondata.PolicyAdmin`）；profile 描述符 `apps/api/kernel/profile.go:202` `Permissions: []string{"wallet.read", "wallet.write", "wallet.adjust", "wallet.voucher.issue"}`。
- **capability 声明**：Job 相关路由**没有任何 feature/navigation capability 键**；wallet 的 navigation 用 `menu_wallet` / `menu_wallet_self` / `menu_wallet_vouchers`（`provider.go:351-391`），与 Job 无关。
- 中间件：所有 Job 路由都包在 `a.Middleware(...)`（认证会话中间件）内，再叠加 `requirePermission`。

---

## 4. LIST/QUERY CAPABILITY GAP

### 4.1 Repository 现有全部方法（`apps/api/internal/jobs/repository.go`）

| # | 方法签名 | 行 | 过滤 / 排序 / 限制 |
|---|----------|----|--------------------|
| 1 | `Create(ctx, input CreateInput) (*Job, error)` | 25 | 无（写入；校验 ID/Kind/ActorID/CorrelationID 非空、payload 合法 JSON） |
| 2 | `Get(ctx, id string) (*Job, error)` | 56 | 仅 `id`。**无 actor 过滤** |
| 3 | `GetForActor(ctx, id, kind, actorID string) (*Job, error)` | 66 | `id AND kind AND actor_id` |
| 4 | `Claim(ctx, id, owner string, now, leaseDuration) (*Job, Lease, error)` | 76 | `id AND ((status='queued' AND attempt<max_attempts) OR (status='running' AND cancel_requested=0 AND lease_expires_at<=? AND attempt<max_attempts))` |
| 5 | `Heartbeat(ctx, lease, now, leaseDuration) error` | 106 | lease 三元组 + `status='running'` |
| 6 | `UpdateProgress(ctx, lease, progress, now) error` | 115 | lease + `progress <= ?` 单调守卫 |
| 7 | `RequestCancel(ctx, id, actorID, now) (*Job, error)` | 124 | `id AND actor_id`（+ 状态分支） |
| 8 | `FinalizeCancel(ctx, lease, now) error` | 151 | lease + `cancel_requested=1` |
| 9 | `Fail(ctx, lease, code, message, now) error` | 162 | lease + `cancel_requested=0` |
| 10 | `CompleteWithCommit(ctx, lease, now, resultTTL, commit) (*Job, error)` | 176 | lease + `status='running'`，同事务提交业务结果 |
| 11 | `Retry(ctx, id, actorID, now) (*Job, error)` | 220 | `id AND actor_id AND status='failed' AND attempt<max_attempts` |
| 12 | `ExpireIfDue(ctx, id, now) (*Job, error)` | 240 | `id AND status='succeeded' AND expires_at<=?` |
| 13 | `ExpireDue(ctx, now) (int64, error)` | 254 | **批量**：`status='succeeded' AND expires_at<=?` → `expired` |
| 14 | `RecoverCancelledDue(ctx, now) (int64, error)` | 259 | 计数包装 |
| 15 | `RecoverCancelledDueJobs(ctx, now) ([]Job, error)` | 264 | `status='running' AND cancel_requested=1 AND lease_expires_at<=?` |
| 16 | `ExhaustExpired(ctx, now) (int64, error)` | 275 | 计数包装 |
| 17 | `ExhaustExpiredJobs(ctx, now) ([]Job, error)` | 280 | `status='running' AND cancel_requested=0 AND lease_expires_at<=? AND attempt>=max_attempts` |
| 18 | **`ListRunnable(ctx, now, limit int) ([]Job, error)`** | **292** | **唯一列表方法**：`(status='queued' AND attempt<max_attempts) OR (status='running' AND cancel_requested=0 AND lease_expires_at<=? AND attempt<max_attempts)`；`ORDER BY created_at, id LIMIT ?`。**无 offset、无 total、无 kind/actor/时间过滤** |
| 19 | `IsCancelRequested(ctx, lease) (bool, error)` | 318 | lease 三元组 |
| — | 私有：`updateLease` (337) / `updateGuardedLease` (347) / `bulkTransition` (371) / `transitionJobs` (384) / `getTx` (464) / `getForActorTx` (468) | — | 内部 |

`ListRunnable` 原文（`repository.go:298-301`）：

```go
rows, err := tx.Query(ctx, `SELECT `+jobColumns+` FROM jobs
WHERE (status='queued' AND attempt < max_attempts)
   OR (status='running' AND cancel_requested=0 AND lease_expires_at <= ? AND attempt < max_attempts)
ORDER BY created_at, id LIMIT ?`, toMillis(now), limit)
```

**它不是用户读面**：无 actor 维度、无状态选择、无分页，且只返回"可运行"子集（永远看不到 succeeded/failed/cancelled/expired）。

### 4.2 通用 admin 作业列表的缺口清单

| 能力 | 现状 | 缺口 |
|------|------|------|
| filter by **kind** | 无。仅 `GetForActor` 的等值（需已知 id） | ❌ 完全缺失 |
| filter by **status** | 无（`ListRunnable` 硬编码 queued/running 谓词） | ❌ 完全缺失 |
| filter by **actor** | 无（`GetForActor` 需已知 id；无"按 actor 列全部"） | ❌ 完全缺失 |
| filter by **time range** | 无（`expires_at<=?` 是维护语义，非用户时间范围） | ❌ 完全缺失 |
| **pagination (limit/offset)** | 仅 `ListRunnable` 有 `limit`，**无 offset** | ❌ 缺失 offset；无 page/pageSize |
| **total count** | 无（全仓唯一 `COUNT(*)` over jobs 在 `apps/api/internal/store/postgres_test.go:844`，属测试） | ❌ 完全缺失 |
| **sort options** | 固定 `ORDER BY created_at, id`（`repository.go:301`） | ❌ 无 `updated_at` / `status` / `kind` 排序 |
| 单条 **admin** 读 | `Get(ctx, id)` 存在但无 HTTP 暴露 | ⚠️ 能力在、路由不在 |
| 列表**响应信封** | 无 | ❌ 需新建（先例信封见 §5） |

**结论**：`admin.jobs` 若要做通用列表，必须**新增 Repository 查询方法**（`I-038-001` 的验证动作已预期此事："确认是否需要新增 repository 列表查询"）。现有方法无一可复用。

### 4.3 现有索引（`apps/api/modules/jobs/migration/migration.go`）

SQLite 与 Postgres 变体**索引完全一致**（`:45-47` 与 `:84-86`）：

```sql
CREATE INDEX idx_jobs_runnable ON jobs(status, cancel_requested, lease_expires_at, created_at)
CREATE INDEX idx_jobs_actor    ON jobs(actor_id, kind, updated_at DESC)
CREATE INDEX idx_jobs_expiry   ON jobs(status, expires_at)
```

表列清单（`migration.go:15-44`）确认：**无 module 列、无 tenant 列**；`kind`/`status`/`actor_id` 均 `NOT NULL`（`:17, :18, :30`），`actor_id` 有 `length(trim(actor_id)) > 0` CHECK。

### 4.4 索引对缺失查询形状的支持判定

| 目标查询形状 | 是否被现有索引覆盖 | 理由 |
|--------------|-------------------|------|
| **`ORDER BY updated_at DESC` over 全表（跨 actor）** | **❌ 否** | `idx_jobs_actor` 以 `actor_id` 打头；跨 actor 时无法用其有序性，退化为全表扫描 + 排序。`idx_jobs_runnable` 以 `status` 打头且第 4 列才是 `created_at`，与 `updated_at` 无关 |
| `WHERE actor_id=? ORDER BY updated_at DESC` | ✅ 是 | `idx_jobs_actor` 前缀 + 该索引第 3 列即 `updated_at DESC`，**完全匹配** |
| `WHERE actor_id=? AND kind=? ORDER BY updated_at DESC` | ✅ 是 | 同上，前两列精确匹配（这正是 `GetForActor` 的 actor+kind 形状，但按 `updated_at` 排序） |
| `WHERE kind=?`（全 actor） | ❌ 否 | `kind` 只是 `idx_jobs_actor` 第 2 列，**无前导 `actor_id` 约束时不可用** |
| `WHERE status=?`（全 actor） | ⚠️ 部分 | `idx_jobs_runnable` 前导列 `status` 可做等值定位，但第 2/3 列（`cancel_requested`, `lease_expires_at`）插入后，`created_at` 排序/范围无法直接受益；`idx_jobs_expiry` 仅覆盖 `expires_at` 配对 |
| `WHERE status=? AND created_at >= ? AND created_at < ?` | ❌ 否 | `created_at` 在 `idx_jobs_runnable` 是第 4 列，前面隔着两列，范围扫描失效；`idx_jobs_expiry` 的次列是 `expires_at` 不是 `created_at` |
| `WHERE status=? AND expires_at<=?`（过期扫描） | ✅ 是 | `idx_jobs_expiry` 精确匹配（`ExpireDue` / `ExpireIfDue` 受益） |
| `WHERE status='succeeded' AND expires_at<=?` 的 runnable 抢占谓词 | ✅ 是 | `idx_jobs_runnable`（`status, cancel_requested, lease_expires_at`）+ `idx_jobs_expiry` |
| `ORDER BY created_at` 全局 | ⚠️ 部分 | 仅 `idx_jobs_runnable` 末列，且需先满足 `status` 等值；全表 `ORDER BY created_at` 无专用索引 |
| `ORDER BY finished_at` / 按 `finished_at` 过滤 | ❌ 否 | **`finished_at` 上完全没有索引** |

**给 R1 的判定结论（事实层）**：现有三个索引是**运行时状态机索引**，不是**管理查询索引**。通用 admin 列表的默认形状（跨 actor、按时间倒序、按 kind/status 过滤）**不被任何现有索引覆盖**；若采用该形状，需要新迁移（新索引或组合索引，如 `(kind, updated_at DESC)` / `(status, updated_at DESC)` / `(updated_at DESC)`）。

**迁移台账约束**：`core.jobs` 迁移的 checksum 已被测试钉死——`apps/api/internal/store/migrate_test.go:693`：

```go
{"core.jobs", "async_jobs", "55e1d3f88de080bd0b6015841e76f1ce32604444619d180a3b228123f99dec68"},
```

且 `apps/api/modules/jobs/migration/migration_test.go:15` 钉死 `d.Version != 42` / `Name != "async_jobs"`。**任何 DDL/索引改动都会同时打破这两处断言**，必须走新迁移版本 + 同步更新期望值。

---

## 5. ADMIN VS ACTOR SCOPE PRECEDENT

仓内**已有两个成熟的"admin 跨 actor 列表"先例**，都是"管理面读别人创建的行"。

### 5.1 先例 A · `admin.activity`（operation_log 全局历史）——最贴近

| 维度 | 事实 | 证据 |
|------|------|------|
| 模块 id | `admin.activity` | `apps/api/modules/activity/provider.go:19` |
| 路由 | `GET /api/operations`, `GET /api/operations/{id}`, `GET /api/operations/export` | `provider.go:40`；挂载 `provider.go:53-61`（`handler.ResourceRoutes(p.a, handler.OperationsResource(p.operations), ModuleID)`） |
| 资源描述符 | `Resource{ID: "operations", Path: "/api/operations", Listable: true, ReadOnly: true, SortFields: []string{"createdAt","event","actorName"}, QSearch: true, ExtraQuery: []string{"event","actorName","from","to"}, PermissionRead: "operations.read", NotFoundCode: "OPERATION_NOT_FOUND"}` | `apps/api/internal/handler/operations.go:15-30` |
| **权限门** | `operations.read`，策略 `PolicyAdminEditor`（admin + editor） | 声明 `activity/provider.go:73-78`；路由级由 `resources.go:322` `requirePermission` 执行；navigation `provider.go:79-91` `Permission: "operations.read"` |
| 查询方法 | `ListOperationsFiltered(filter OperationFilter) ([]Operation, int, error)` | `apps/api/modules/operationlog/repository.go:218-254` |
| **过滤形状** | Q（LIKE event/actor_name/detail/record_id）+ `event = ?` + `lower(actor_name) = lower(?)` + `created_at >= ?` + `created_at <= ?` | `operationsWhere` `repository.go:316-343` |
| **排序** | `ORDER BY <col> <ASC\|DESC>, id DESC`；列 ∈ {created_at, event, actor_name} | `operationsSortSQL` `repository.go:345-357` |
| **分页** | `LIMIT ? OFFSET ?` + `pagination.Offset(page, pageSize, total)`，**同一 WHERE 先 `COUNT(*)`** | `repository.go:223`（count）、`:233-234`（limit/offset） |
| 分页助手 | `pagination.Bounds` / `pagination.Offset`（防溢出） | `apps/api/internal/pagination/pagination.go:16, 39` |
| **响应信封** | `resourceList{Items, Total, Page, PageSize}` → `{"items":[...],"total":N,"page":P,"pageSize":S}` | 类型 `apps/api/internal/handler/resources.go:253-258`；写出 `:479` |
| 行映射 | `operationToMap(op)` | `handler/operations.go:36-57` |
| 索引先例 | `CREATE INDEX idx_operation_log_created_at ON operation_log(created_at DESC)` | `apps/api/modules/operationlog/migration/migration.go:23`（多版本重复出现） |
| 关键性质 | **查询完全不按 actor 过滤**——是真正的全局管理读面 | `operationsWhere` 无 actor 谓词 |

> 对照：`admin.activity` 用"全局列表 + 单一 read 权限 + 结构化过滤（event/actor/time）+ 统一信封"解决"读别人的行"；Job 侧**完全没有对应物**。

### 5.2 先例 B · `admin.scheduled-tasks` 的全局 run 历史（`task-runs`）

| 维度 | 事实 | 证据 |
|------|------|------|
| 模块 id | `admin.scheduled-tasks` | `apps/api/modules/scheduledtasks/provider.go:22` |
| 路由 | `GET /api/task-runs`, `GET /api/task-runs/{id}`（全局历史）+ `GET /api/scheduled-tasks/{id}/runs`（按 task） | `provider.go:60-61` |
| 权限门 | `tasks.read` | 声明 `provider.go:94-101`；路由级 `apps/api/internal/handler/scheduledtasks.go:473` `requirePermission(w, r, "tasks.read")` |
| 全局查询 | `ListAllRuns(filter ListFilter) ([]TaskRun, int, error)` | `apps/api/modules/scheduledtasks/store/repository.go:317-...` |
| 过滤形状 | `Q`（detail LIKE + task key 子查询）+ `status = ?`（**括号优先级修复记录在 `:324-329`**） | `:321-338` |
| 分页 | `COUNT(*)` `:339`；`ORDER BY started_at DESC LIMIT ? OFFSET ?` + `pagination.Offset` `:342-346` | 同上 |
| 信封 | `resourceList{Items, Total, Page, PageSize}` | `handler/scheduledtasks.go:496` |
| 行映射 | `taskRunToMap(rn)` → `{id, taskId, status, startedAt, detail, finishedAt, ...}` | `handler/scheduledtasks.go:510-519` |
| 关键性质 | **无 actor 维度**（run 属于 task，不属于 actor） | — |

### 5.3 无先例的部分

- **没有任何模块读取 `jobs` 表**（§3.3）。因此"跨 actor 读 jobs"在仓内**没有直接先例**，只能沿用上表的**模式**（权限 + 全局 WHERE + COUNT + LIMIT/OFFSET + `resourceList` 信封 + `pagination.Offset`），不能沿用现成代码。
- `admin.data-transfer` 未发现 job/batch 列表（grep `job|Job|ORDER BY|func (r *Repository) List` 在该模块零命中）——**待确认**是否其导出/导入面以其他命名存在（见 §8）。

---

## 6. PERMISSION / CAPABILITY WIRING（新 admin 模块接线清单）

以 `admin.scheduled-tasks` 为主参考（`admin.activity` 为次参考）。**加粗**为必做项。

| # | 动作 | 文件 | 符号 / 行 |
|---|------|------|-----------|
| 1 | **声明 ModuleID** | `apps/api/modules/<mod>/provider.go` | `const ModuleID = "admin.scheduled-tasks"` — `scheduledtasks/provider.go:22`（同值另见 `schema/schema.go:9`、`migration/migration.go:13`） |
| 2 | **实现 Descriptor** | 同上 | `func (p *Provider) Descriptor() kernel.Module` — `:46-69`：`ID`, `Version: "2.0.0"`, `KernelAPIRange: ">=2.0 <3.0"`, `DependsOn: []string{"core.auth-session","core.navigation-capability","core.schema-render","core.operationlog"}`, **`Requires: kernel.StandardAdminCapabilities()`**, `Contributions: kernel.ContributionKeys{Routes, Pages, Navigation, Permissions, Fragments}` |
| 3 | 能力集定义 | `apps/api/kernel/profile.go` | `func StandardAdminCapabilities() []Capability` — `:234-236` → `{CapabilityHTTP, CapabilitySchema, CapabilityAuthorization, CapabilityNavigation, CapabilityManifest, CapabilityPersistence}` |
| 4 | **实现 Register（路由挂载）** | `apps/api/modules/<mod>/provider.go` | `func (p *Provider) Register(ctx, reg kernel.Registrar) error` — `:75-122`；`for _, route := range handler.XxxRoutes(...) { reg.HTTP(route) }` `:76-80` |
| 5 | **路由贡献 + 认证中间件** | `apps/api/internal/handler/<mod>.go` | `kernel.RouteContribution{ContributionIdentity{ModuleID, Key: kernel.RouteKey(method, pattern)}, Method, Pattern, Handler: a.Middleware(...)}` — 范例 `handler/scheduledtasks.go:468-498`；wallet 的 `add(...)` 工厂 `handler/wallet.go:103-110` |
| 6 | **路由级权限** | handler 内 | `if _, ok := requirePermission(w, r, "tasks.read"); !ok { return }` — `handler/scheduledtasks.go:473`；门实现 `handler/resources.go:322-333` |
| 7 | **权限贡献** | provider | `reg.Authorization(kernel.PermissionContribution{ContributionIdentity{ModuleID, Key: "tasks.read"}, Permission: "tasks.read", Resource: "scheduled-tasks", Action: "read", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion})` — `scheduledtasks/provider.go:94-101`（循环 `:94-101`）。`PolicyAdmin` = 仅管理员；`PolicyAdminEditor` = 管理员+编辑（activity 用后者，`activity/provider.go:75`） |
| 8 | **页面 schema 贡献** | provider + `apps/api/modules/<mod>/schema/schema.go` | `reg.Schema(kernel.PageContribution{... PageID, Resources, Actions, DataSource, Owner, Document: xxxschema.SchemaDocuments()[pageID]})` — `scheduledtasks/provider.go:81-93`；schema 包 `schema/schema.go:9-23`（`//go:embed` + `PageIDs()` + `SchemaDocuments()`） |
| 9 | **导航贡献** | provider | `reg.Navigation(kernel.NavigationContribution{NodeID, PageID, Order, Label, Group, Visibility: authsessiondata.PolicyAdmin, Permission: "tasks.read", SystemDataVersion})` — `scheduledtasks/provider.go:102-114` |
| 10 | **Manifest fragment** | provider + `apps/api/modules/<mod>/manifest/` | `reg.Manifest(kernel.FragmentContribution{FragmentID, ProtocolVersion: "2.7", RequiredCapabilities: []string{"manifest","navigation"}, JSON: manifest.FragmentJSON})` — `scheduledtasks/provider.go:115-121`；`manifest/manifest.go` + `manifest/fragment.json` |
| 11 | **fragment 页面/路由** | `manifest/fragment.json` | `pages[]` 每项 `{pageId, title, schemaUrl: "/api/schema/<page>", route: "/<page>", titleKey}`；`navigation.sidebar[]` 用 `visibleWhen: {"when": "$context.features.menu_scheduled_tasks == true"}` — `scheduledtasks/manifest/fragment.json:12-39` |
| 12 | **进 admin 默认集** | `apps/api/kernel/profile.go` | ① `profileDefaults[ProfileAdmin]` 追加模块 id — `:46-93`（例 `:74-76` `"admin.scheduled-tasks"`）；② `BuiltinModules` 追加完整描述符 — `:193`（admin.wallet 在 `:202`） |
| 13 | **组合根按 plan 装配** | `apps/api/internal/composition/composition.go` | `if plan.HasModule("admin.scheduled-tasks") { providers = append(providers, ...) }` — `:555-557`；wallet 先例 `:582-600` |
| 14 | 迁移（若需新表/索引） | `apps/api/modules/<mod>/migration/migration.go` + `apps/api/modules/compiled/persistence.go` | 注册 `kernel.MigrationContribution`；compiled 目录登记 — `compiled/persistence.go:29-51` |
| 15 | **Web 侧 schema 路径登记** | `apps/web/src/i18n/schema-keys.structural.test.ts` | `"scheduledtasks/schema/task-runs.json"` — `:45`（新页面必须加入该清单） |
| 16 | Web 侧渲染分母登记 | `apps/web/src/i18n/s5-denominator-render.test.tsx` | `"task-runs": resolve(MODULES, "scheduledtasks/schema/task-runs.json")` — `:72` |
| 17 | Web 侧导航/分组 | `apps/web/src/app/nav-groups-r4.test.ts` | `{ path: "/task-runs", groupKey: GROUP_KEYS.operations, pageRef: "scheduled-tasks" }` — `:33` |
| 18 | Web 侧搜索/权限矩阵 | `apps/web/src/app/searchable-profile-matrix.test.ts` | page id `:24`、page→menu 映射 `:56`、`"action:scheduled-tasks:create"` `:254` |
| 19 | Web 侧父子层级（子页高亮） | `apps/web/src/app/navigation.ts` | `NAVIGATION_PAGE_PARENTS` — `:21-27`（例 `"task-runs": "scheduled-tasks"` `:23`） |
| 20 | Web 侧 manifest 期望 | `apps/web/src/protocol/app-manifest.test.ts` | 页面 id 清单 `:303-305` |
| 21 | i18n 文案 | `apps/web/src/i18n/messages/{zh-CN,en-US}.json` | 新 `manifest.title.*` / `manifest.nav.*` / `schema.*` 键 |

**必须同步更新的计数型断言（否则 CI 红）**：

| 断言 | 文件:行 | 现值 |
|------|---------|------|
| admin profile 权限/导航计数 | `apps/api/internal/composition/composition_test.go:529` | `{profile: "admin", wantPermissions: 34, wantNavigation: 18}` |
| 迁移 catalog 与 checksum | `apps/api/internal/store/migrate_test.go:693` | `{"core.jobs", "async_jobs", "55e1d3f8..."}` |
| core.jobs 迁移描述符 | `apps/api/modules/jobs/migration/migration_test.go:15` | `Version != 42` / `Name != "async_jobs"` |
| provider 路由/页面集合一致性 | `apps/api/kernel/provider.go:178`（`stringSetEqual`） | Descriptor `Contributions` 必须与实际 `reg.*` 调用**逐键一致** |
| 错误码冻结表 | `apps/api/internal/handler/error_contract_test.go:75` | Job 码清单 |

---

## 7. 会约束通用 Job 读面的现有测试

| 测试 | 文件:行 | 约束内容 | 对通用读面的含义 |
|------|---------|----------|------------------|
| `TestWalletJobServiceCompletesAtomicallyAndAudits` | `apps/api/modules/wallet/jobs_test.go:153-155` | `service.Get(ctx, job.ID, "other-user")` 必须 `errors.Is(err, jobs.ErrNotFound)` | **最强约束**：actor 作用域不可被放宽；通用 admin 读面必须走**另一条**路径（新方法/新路由），不能改 `GetForActor` 语义 |
| `TestCancellationPathsAndActorIsolation` | `apps/api/internal/jobs/repository_test.go:97-99` | 跨 actor `RequestCancel` → `ErrNotFound` | 同上，取消/重试的 actor 隔离冻结 |
| `TestCreateValidationAndGetForActor` | `apps/api/internal/jobs/repository_test.go:218-224` | 跨 actor `GetForActor` → `ErrNotFound`；同 actor 命中且 `MaxAttempts == DefaultMaxAttempts` | 同上 + `DefaultMaxAttempts` 期望 |
| `TestWalletRoutesGates` | `apps/api/internal/handler/wallet_test.go:257-315` | 匿名 401 + editor 403，**逐条列出 4 条 Job 路由**（`:274-277`, `:303-306`） | 新路由必须自带权限门；job 路由清单是硬编码列表 |
| `TestWalletReconcileBadBodyAndWriteGate` | `wallet_test.go:680-708` | submit/cancel/retry 需 `wallet.write`；garbage body → 400 | 写操作权限口径冻结 |
| 结果/终态码断言 | `wallet_test.go:435-452` | `JOB_NOT_CANCELLABLE` / `JOB_NOT_RETRYABLE` 409；`JOB_NOT_FOUND` 404；`:461` `JOB_RESULT_EXPIRED` 410 | 通用读面必须沿用同一错误码语义 |
| `TestWalletErrorCodesCataloged` | `wallet_test.go:662-670` | 18 个 wallet 码（含全部 `JOB_*`）必须在 `errorcatalog.Catalog` 有 En/Zh/MessageKey | 新错误码必须入 catalog |
| 冻结 wire 码 | `apps/api/internal/handler/error_contract_test.go:75` | `"JOB_NOT_FOUND", "JOB_NOT_CANCELLABLE", "JOB_NOT_RETRYABLE", "JOB_RESULT_NOT_READY", "JOB_RESULT_EXPIRED"` | 冻结契约，不可删改 |
| 冻结持久化终态码 | `error_contract_test.go:92-95` | `frozenStoredCodes = []string{"JOB_ATTEMPTS_EXHAUSTED", "JOB_HANDLER_FAILED"}` | 客户端读 failed job 时会看到这两个 code |
| kind 重复注册守卫 | `apps/api/internal/jobs/runner.go:104-106`；测试注册点 `handler/wallet_test.go:164` | 同 kind 二次 `Register` 报错 | 新模块若注册 `wallet.reconcile` 会启动失败 |
| 迁移描述符 | `apps/api/modules/jobs/migration/migration_test.go:15` | `len(descriptors) != 1`、`Version != 42`、`Name != "async_jobs"` | 加索引/改 DDL 必须新版本 + 改断言 |
| 迁移 checksum | `apps/api/internal/store/migrate_test.go:693` | `core.jobs/async_jobs` checksum 冻结 | 同上 |
| 组合计数 | `apps/api/internal/composition/composition_test.go:529` | `wantPermissions: 34, wantNavigation: 18` | 新模块的权限/导航会改这两个数 |

---

## 8. 待确认 / 未知

1. **`admin.data-transfer` 的批量/作业面未确认**：对该模块 grep `job|Job|ORDER BY|func (r *Repository) List` **零命中**。无法排除其导出/导入面以其他命名（如 `export`/`import`/`transfer`）实现，或根本是同步实现。→ 影响 `I-038-003`（首波批量操作分母），**本次未读完该模块全部文件**。
2. **`biz.digital-offer` / `channel.telegram` 是否间接产生 Job**：本次以 `runner.Submit(` / `jobs.CreateInput` 全仓扫描为准，命中仅 wallet 与测试；但**未逐一阅读**这两个模块的全部文件以排除其自建队列/重试表（非 jobs 表）。
3. **`jobRuntime.enabled` 的门控语义**：`composition.go:588` 仅在 `plan.HasModule("admin.wallet")` 时 `jobRuntime.enabled.Store(true)`；`:1101` 的 `jobs.Start()` 在 `enabled=false` 时是 no-op（`:185-190`）。因此**在没有 admin.wallet 的 Profile 下 Job runner 不会启动**——一个不含 wallet 但含 `admin.jobs` 的 Profile 会怎样，代码未表达。→ 对 `I-038-004` 的"进 admin 默认集"结论有潜在影响（admin 集含 wallet，故当前无矛盾；其他 Profile 待确认）。
4. **`expired` 状态在列表中的产生时机**：`ExpireDue` 由 `ScanOnce` 周期调用（`runner.go:201`），`ScanInterval` 默认 10s（`runner.go:51`）。列表读面是否需要显式惰性过期（对齐 `JobService.Get` `jobs.go:86-91`）**未在任何代码或文档中定义**。
5. **`payload`/`result` 是否可在通用列表中暴露**：`A-012-r4-s4-closeout-independent.md:98` 记录"GET Job 不内嵌 payload/result"是 VP-012 的冻结口径，但**该冻结的适用范围（仅 wallet 路由 vs 全仓 Job 表示）在代码中无强制机制**（`walletJobToMap` 只是私有函数）。→ 需 R1 决策，非本次可判定。
6. **是否有跨 actor 读取 `jobs` 的既有设计文档**：本次只读了 VP-012/VP-038 相关文档片段，未穷举 `docs/`。不排除存在已裁定的"通用 Job 管理页为非目标"以外更细的口径。
7. **`idx_jobs_actor` 的实际查询计划**：以上索引判定基于列序的静态推理，**未运行 `EXPLAIN QUERY PLAN`**（只读任务，且未执行数据库命令）。若 R1 需要硬证据，应在实施阶段补 EXPLAIN 证据。
8. **`resultUrl` 硬编码前缀的兼容口径**：`walletJobToMap` 写死 `/api/wallet/jobs/<id>/result`（`wallet.go:1003`）。通用作业中心复用该字段是否会把用户导向 wallet 权限边界（`wallet.read`）**未确认**。
9. **本报告未覆盖**：Web 端是否存在 wallet 作业的既有 UI 消费者。grep `apps/web` 的 `reconcile|job` 仅命中 i18n 文案 `schema.wallet.toolbar.reconcile`（`apps/web/src/i18n/messages/zh-CN.json:858`）与一个 profile matrix 期望 `"action:wallet:reconcile"`（`apps/web/src/app/searchable-profile-matrix.test.ts:260`）；**未见 job 列表/结果中心页面**，但未穷举 `apps/web/src/components/`。
