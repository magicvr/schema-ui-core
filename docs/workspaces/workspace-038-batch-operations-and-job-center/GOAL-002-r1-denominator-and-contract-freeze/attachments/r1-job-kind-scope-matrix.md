---
title: R1 冻结矩阵 · Job 种类×作用域（I-038-001 / GOAL-002 C1）
status: frozen
created: 2026-09-19
updated: 2026-09-19
parent: GOAL-002-r1-denominator-and-contract-freeze
version: 1.0.0
frozen_by: D-001-r1-contract-and-denominator-freeze
---

# R1 冻结矩阵 · Job 种类×作用域

> **性质**：GOAL-002 C1 的可机器核对交付物，冻结依据 `01-decision/D-001-r1-contract-and-denominator-freeze.md` §2。
> **证据基线**：`attachments/R1-recon-I-038-001-job-kinds-and-scopes.md` + 编排器独立复核（Root `E-002` F-1～F-4、GOAL-002 `E-001` G-1～G-3）。
> **机器可核对**：每条含 `symbol` 与 `evidence` 列，可用 `grep` 逐条复验。

## 1. Job 种类分母（冻结）

| # | kind | 状态 | 声明点 | handler | 注册方式 | 提交入口 | 首波归属 |
|---|------|------|--------|---------|---------|---------|---------|
| 1 | `wallet.reconcile` | **已交付** | `modules/wallet/jobs.go:19` `const ReconcileJobKind` | `(*JobService).runReconcile` `jobs.go:116` | `RegisterWithTerminalHook` `jobs.go:39` | `POST /api/wallet/reconcile` `handler/wallet.go:379` | **不进首波**（已是先例/模板，`out-of-scope`） |
| 2 | 异步批量导出（kind 字符串待 R3 冻结，建议 `jobs.batch-export`） | **首波新建** | R3 实施 | R3 实施 | `RegisterWithTerminalHook`（runner 启动前） | R3 新增端点 | **首波唯一纳入** |
| — | `panic.kind` | 仅测试 | `internal/jobs/runner_panic_test.go:30` | 测试内联 | 测试内 | — | 不计入生产分母 |

**冻结结论**：生产 Job 种类**当前 = 1**（`wallet.reconcile`）；首波后 = **2**。kind 注册表是内存 `map[string]registration`（`runner.go:62`），**无元数据/枚举**，重复 kind 注册会启动失败（`runner.go:101-106`）——分母只能由 `Register*` 调用点扫描得出，无中心清单。

## 2. 作用域矩阵（冻结）

| 维度 | 现行（`admin.wallet`） | **首波冻结（`admin.jobs`）** |
|------|----------------------|---------------------------|
| 读列表 | **不存在** | **管理作用域**（跨 actor），`jobs.read` |
| 读详情 | actor 作用域：`GetForActor(id, kind, actorID)`（`repository.go:66-74`，SQL `WHERE id=? AND kind=? AND actor_id=?`） | 新增管理作用域详情（新方法 + 新路由），`jobs.read` |
| 取消 | actor 作用域：`RequestCancel(id, actorID)`（`repository.go:124-149`，`getForActorTx` = `WHERE id=? AND actor_id=?`；**repository 不校验 kind**，kind 校验在 wallet service `jobs.go:96`） | 新增管理作用域写路径（写权限），**既有 actor 路径不动** |
| 重试 | actor 作用域：`Retry(id, actorID)`（`repository.go:220-238`，`WHERE id=? AND actor_id=? AND status='failed' AND attempt<max_attempts`） | 同上 |
| 结果 | actor 作用域 + wallet 专属硬编码 URL | 新增管理作用域结果读取；结果 URL 须**通用化**（现行 `walletJobToMap:1003` 硬编码 `/api/wallet/jobs/{id}/result`，是 wallet 权限门下的路径，不可直接复用） |
| 未授权/越界 | 跨 actor 与错 kind 均塌缩为 `ErrNotFound` → `404 JOB_NOT_FOUND`（`repository.go:109-111`） | **保持**同一塌缩语义（不泄露存在性） |
| 权限 | 无 job 专属权限；复用 `wallet.read` / `wallet.write`（`wallet.go:358,388,401,414,427`） | **新增** `jobs.read`（+ 写权限键，R2 冻结），策略 `PolicyAdmin` |
| 模块维度 | 表中**无 module 列**（`migration.go:15-44`） | 首波不新增 module 列；跨模块可见性由 kind 白名单/权限表达（R2 冻结） |

## 3. 读面缺口矩阵（冻结）

| 能力 | 现状 | 冻结要求 |
|------|------|---------|
| 按 kind 过滤 | ❌ 无 | R2 新增 |
| 按 status 过滤 | ❌ 无 | R2 新增（须含六态全量，不得只覆盖 queued/running） |
| 按 actor 过滤 | ❌ 无 | R2 新增 |
| 时间范围过滤 | ❌ 无 | R2 新增 |
| 分页（page/pageSize + offset） | ❌ 仅 `ListRunnable` 的裸 `LIMIT` | R2 新增，沿用 `internal/pagination`（`pagination.Bounds`/`Offset`）与 `resourceList{items,total,page,pageSize}` 信封（`resources.go:253-258`） |
| 总数 COUNT | ❌ 无 | R2 新增（与列表同 WHERE，先例 `operationlog/repository.go:223`） |
| 排序选项 | ❌ 固定 `ORDER BY created_at, id`（仅 `ListRunnable`） | R2 新增，须与索引决策（O-3）一致 |

## 4. 索引矩阵（冻结事实 + 待决）

| 索引 | 定义 | 覆盖的查询形状 | 管理列表所需形状是否覆盖 |
|------|------|---------------|------------------------|
| `idx_jobs_runnable` | `(status, cancel_requested, lease_expires_at, created_at)` | worker 队列扫描 | ❌（`created_at` 是第 4 列，范围/排序不受益） |
| `idx_jobs_actor` | `(actor_id, kind, updated_at DESC)` | 单 actor（+kind）按 `updated_at` 排序 | ⚠️ 仅 actor 作用域；**跨 actor 的 `ORDER BY updated_at DESC` 不被覆盖** |
| `idx_jobs_expiry` | `(status, expires_at)` | 结果过期扫描 | ❌ |

定义位置：`modules/jobs/migration/migration.go:45-47`（sqlite）与 `:84-86`（postgres），两侧一致。

**待决（O-3，属 R2 方案）**：是否为跨 actor 默认排序新增索引。若新增，须以**新迁移版本号**进行，并同步更新两处冻结断言：
- `internal/store/migrate_test.go:692-693`（`core.jobs/async_jobs` checksum `55e1d3f8…`）
- `modules/jobs/migration/migration_test.go:15`（`Version == 42` / `Name == "async_jobs"`）

## 5. 后台周期任务口径（冻结）

仓内共 **4** 个后台周期循环，**均不走 `jobs` 表**，首波一律 `out-of-scope`：

| 周期任务 | 位置 | 间隔 | 首波口径 |
|---------|------|------|---------|
| scheduled-tasks cron 调度器 | `modules/scheduledtasks/scheduler.go:66` | 30s | out-of-scope（非用户触发） |
| 审计日志保留期清理 | `modules/operationlog/retention.go:103` | 1h | out-of-scope（同上） |
| Job 轮询循环 | `internal/jobs/runner.go:216` | 10s | 基础设施本身 |
| Telegram 租约调和 | `channel/telegram/connection_manager.go:390` | 1s | out-of-scope（同上） |

**冻结结论**：首波**不**把上述非 Job 后台任务迁移进 `jobs`；作业中心的可见性分母 = `jobs` 表内容，不含这些循环。（不排除后续波次迁移，但本轮不承诺。）

## 6. 运行时门控事实（冻结，供 R2/R3 遵循）

| # | 事实 | 证据 | 对首波的含义 |
|---|------|------|-------------|
| R-1 | `jobRuntime` 在组合根**无条件构造**，但 `enabled` 仅在 `plan.HasModule("admin.wallet")` 时置 true；`Start()`/`Stop()` 在 `enabled=false` 时为 no-op | `composition.go:170-197`、`:582-588`、`:1101` | **`admin.jobs` 必须自行置 `enabled`**（或改为按 `admin.jobs`/`admin.wallet` 任一存在即启用），否则「含 `admin.jobs` 但不含 `admin.wallet`」的 Profile 下 runner 不启动。admin 默认集含 wallet 故当前无矛盾，但 R2 须显式处理 |
| R-2 | 迁移 42 已存在 `jobs` 表（含六态状态机 CHECK 与三索引）；`payload`/`result` 为 TEXT JSON | `modules/jobs/migration/migration.go:14-48` | 新增异步操作**无需新建表**；仅索引决策（O-3）可能需新迁移 |
| R-3 | Runner 默认参数：Lease 30s / Heartbeat 10s / Scan 10s / **ResultTTL 24h** / BatchSize 32 | `internal/jobs/runner.go:48-54` | 结果保留 24h 后转 `expired`（结果被清空）——结果中心须呈现该终态 |
| R-4 | `expired` 由 10s 扫描器（`ScanOnce`）或 `JobService.Get` 惰性触发 | `runner.go:201`；`modules/wallet/jobs.go:86-91` | 列表侧的过期语义须在 R2 显式定义（惰性 vs 依赖扫描器） |
| R-5 | 错误码已冻结：`JOB_NOT_FOUND` / `JOB_NOT_CANCELLABLE` / `JOB_NOT_RETRYABLE` / `JOB_RESULT_NOT_READY` / `JOB_RESULT_EXPIRED`（wire）；`JOB_ATTEMPTS_EXHAUSTED` / `JOB_HANDLER_FAILED`（stored） | `internal/handler/error_contract_test.go:75,92-95` | 复用同一码族，不新造同义码 |
