---
id: A-012-r4-s4-closeout-independent
goal: GOAL-005-r4-async-job-contract
doc: audit-entry
record_id: A-012
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: R4 S4 close-out；D-002 v0.2.0；四条成功标准；A-001～A-011；migration/Profile 边界、状态机/恢复、CompleteWithCommit、wallet actor HTTP/result、终态审计、migration 43 correlation、错误码契约、全量验证
audit_type: close-out
verdict: pass
status: recorded
parent: GOAL-005-r4-async-job-contract
created: 2026-08-18
updated: 2026-08-18
version: 0.1.0
reviews: A-011
---

# A-012 · R4 S4 independent close-out（2026-08-18）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high；项目级路径见 `docs/architecture/independent-audit-execution.md`）
- **类型**：close-out / S4
- **scope**：GOAL-005 S4 关门。核对 D-002 v0.2.0 不变式、四条成功标准、A-001～A-011 开放项与关闭证据，以及当前代码/测试/提交。重点：migration/Profile 边界、状态机与恢复、`CompleteWithCommit` 原子性、wallet actor HTTP/result、终态审计、migration 43 correlation 保留、错误码契约、全量验证。
- **verdict**：pass

## 范围与区间

- **工作区**：`workspace-012-shared-cross-module-contracts`（`workspace.md`：`root_goal` = `GOAL-001-shared-cross-module-contracts`；`canonical_scope` 与本目标路径一致；`shared_materials_catalog: none`；`vision_role: delivery`；`primary_plan` = `VP-012-shared-cross-module-contracts`）
- **covered**：`00-meta` / D-002 v0.2.0 / E-005～E-008 / A-001～A-011；`apps/api` 当前 `internal/jobs`、`modules/jobs/migration`、`modules/wallet`、`handler/wallet.go`、`composition`、`kernel/profile.go`、`errorcatalog`、`operationlog` migration 43、`store` catalog；提交 `2013e7f`、`e670b56`、`c8305bb`、`3ce848b`、`425215a`；本轮独立复跑测试
- **excluded**：不改 `status` / `progress` / 路线图 / D-002 / 实现代码；不读取或比较其他工作区治理上下文（Appendix A 增量仅作为 E-008 引用的 Q2 证据核验）；不审通用 Job UI、外部队列、scheduled-task 改造、R5/R6
- **本轮复验**：本会话独立执行并通过 `go test -timeout 15m ./...`；另复跑 `./internal/jobs`、`./internal/modules/wallet`、`./internal/store`、定向 `handler`/`kernel`/`composition`/`docscheck`，以及 `go test -race` 与 `-count=10`（jobs + wallet）；`git diff --check` 通过

## 工作区与对齐（只读）

| 检查项 | 结论 | 证据 |
|--------|------|------|
| 工作区绑定 | 通过 | `workspace.md` Root / canonical / `plan_refs`+`primary_plan` 与 GOAL-005 `parent`、`primary_plan` 一致 |
| 共享资料引用 | 无引用，不构成关闭证据 | `shared_materials_catalog: none` |
| 对齐链 | 未发现与 Root R4 / VP-012 的明显冲突 | 六态 Job + wallet 有界真实消费；无新默认 Profile 模块 ID；无 Tier D |
| Vision Review required | 本 scope 未见开放 required | 本意见不审 Vision Review 本身 |
| 既有 Goal 审计 | A-001～A-010 无开放 required；A-011 self = pass，声明不单独放行 S4 | `03-audit.md` |
| P-004 冲突 | 无 | A-011 与本条无相反 verdict；本条不重开已 `fixed` 的 F-001～F-010 |

## 成果（有证据）

| 主张 | 证据 | 核验 |
|------|------|------|
| 提交链存在且与 E-005～E-008 一致 | `2013e7f` migration 42；`e670b56` repository；`c8305bb` runner；`3ce848b` wallet 消费；`425215a` 错误码钉死 | 通过（`git show -s`） |
| `core.jobs` 仅 PersistenceProviders，`Register` 为空 | `modules/jobs/migration/provider.go:20`；`compiled/persistence.go:42`；`profile.go` `profileDefaults` / `BuiltinModules()` 无 `core.jobs` | 通过 |
| wallet 四条 route key 双写，无新 page/nav/fragment/permission | `wallet/provider.go:179-188` 与 `kernel/profile.go:199` 同步四 key；Pages/Navigation/Permissions/Fragments 未增 Job 面 | 通过 |
| runner 仅在 `admin.wallet` 启用时启动 | `composition.go:385-392` `enabled.Store(true)`；`jobRuntime.Start` 未启用则 no-op | 通过 |
| 转换表已落地为条件 UPDATE | `jobs/repository.go` claim/reclaim、heartbeat、progress、cancel、fail、finalize-cancel、recover-cancel、exhaust（`cancel_requested=0`）、retry、expire | 通过；F-010 表/散文差已消 |
| `CompleteWithCommit` 先复核 fencing，callback 与 succeeded 同事务 | `repository.go:176-217`；失效 token 不执行 callback；callback 错误回滚后新 token 可提交同一 run | 通过；`repository_test.go:147-185` |
| wallet callback 使用 `ReconcileOnceTx(tx)`，不调自开事务的 `ReconcileRun` | `wallet/jobs.go:110-116`；`store/repository.go:650-722` | 通过；失败不留 run（`jobs_test.go:163-176`） |
| HTTP：202 + actor 谓词 + 统一 404 + result attachment | `handler/wallet.go:300-387,524-534,617-634`；`JobService.Get/Cancel/Retry` 先 `GetForActor(kind, actor)` | 通过 |
| 终态审计 queued/success/failed/cancelled；scanner hook 覆盖 recover-cancel / exhaust | `jobs.go:59-60,119-136`；`runner.go:185-200`；`jobs_test.go` + `runner_test.go:75-128` | 通过 |
| migration 43 重建前备份/恢复 correlation，并允许三类新事件 | `operationlog/migration/migration.go:423-446,490-503`；`store/migrate_0043_test.go` | 通过 |
| 七个 Job 码进入 catalog 与 pinned set | `errorcatalog.go:152-158`；`error_contract_test.go:65-72`；Q2 Appendix A 增量见 workspace-007 Root D-002 | 通过 |
| 本轮测试 | 见下「本轮复验」 | 通过 |

## 对照成功标准

| 标准 | 本轮 | 证据 |
|------|------|------|
| 1. 六态与合法/非法转换、单调进度、attempt/lease、取消、重试、过期均有确定性测试 | **达成** | repository：claim/reclaim、递减进度拒绝、fail-after-cancel 拒绝、queued/running/recover-cancel、exhaust 后不可 retry、expire 清 result；runner：startup reclaim、运行中取消、失败+retry+自动过期、双 runner heartbeat 只执行一次、Stop 保持 running、无 cancel 的 `context.Canceled` → failed |
| 2. Job 持久化；runner 领取 queued / lease-expired running；重复领取不重复执行 | **达成** | migration 42 CHECK + `jobs` 表；`Claim` 合并 claim/reclaim；active set 跳过同进程 in-flight；fencing 使旧 owner 不能写；wallet 一次 Job 至多一行 run |
| 3. wallet reconcile 提交 202；可轮询并读最终结果；过期返回稳定错误 | **达成**（HTTP 410 为映射+repository 证据，见 F-011） | `wallet_test.go:317-368` 202 + poll + attachment `consistent` + 终态 409；`ExpireIfDue` 原子清 result；handler `expired` → 410 `JOB_RESULT_EXPIRED` |
| 4. API 全量测试与迁移测试通过；Profile 默认集、模块矩阵、Manifest 装配语义不变 | **达成** | 本轮 `go test -timeout 15m ./...` PASS；catalog 42/43 checksum 钉死；mvp 模块列表测试无 `core.jobs`；composition admin permissions=30 / navigation=15 未因 Job 路由增加 |

## 信息门禁核对（P-005）

| ID | 级别 | 最晚阶段 | 当前（00-meta） | 本轮结论 |
|----|------|----------|-----------------|----------|
| I-001 | required | S0 结束前 | verified | 保持；扫描事实仍成立 |
| I-002 | required | S0 结束前（影响 S1） | verified | 保持；实施后 Profile/BuiltinModules/Plan-gated runtime 仍无 `core.jobs` |
| I-003 | required | S0 结束前（影响 S1/S2） | verified | 保持；转换表已实现且无新的可达无事件 `running` |
| I-004 | required | S4 关门 | verified | **本条即为所要求的 independent 关门审**。模式/provider 与项目级路径一致 |

无 `deferred`。无用户书面 `accepted-residual`。无到期未关闭 required 信息项。

## 重点核查

### migration / Profile 边界

`core.jobs` 只出现在 `compiled.PersistenceProviders()`。Provider `Register` 为空。`profileDefaults`（mvp/admin/demo）与 `BuiltinModules()` 均无该 ID。`composition` 不把 Job runtime 当模块 provider；仅当 Plan 含 `admin.wallet` 时注册 `JobService` 并 `enabled=true`。wallet Descriptor 与 BuiltinModules 同步四条 route key；权限仍为 `wallet.read` / `write` / `adjust`。A-001 F-002 / A-002 F-004 的实施条件已满足。

### 状态机 / 恢复

`ScanOnce` 顺序：`recover-cancel` → `exhaust` → `ExpireDue` → `ListRunnable`。`exhaust` WHERE 含 `cancel_requested=0`（A-006 F-010 已消）。`recover-cancel` 不依赖旧 owner，且不受 attempt 是否耗尽影响。lease-expired + attempt 耗尽 → `JOB_ATTEMPTS_EXHAUSTED`。Stop 取消 handler 但不伪写终态。可达 `running` 组合均有事件，无新死状态。

### CompleteWithCommit 原子性

事务内：读行复核 `running`+owner+version → callback → 条件 UPDATE succeeded。任一步失败整体回滚。cancel/recover-cancel 先赢则 fencing 失败、callback 不执行（或 UPDATE 失败回滚）。有效 fencing 即使 `cancel_requested=1` 仍可先赢。wallet callback 只调 `ReconcileOnceTx`；缺失账户失败后 run 数为 0，Job `JOB_HANDLER_FAILED`。A-006 要求的 S3 必测项成立。

### wallet actor HTTP / result

五个 Job endpoint 先 `wallet.read`，再 `kind=wallet.reconcile && actor_id==当前用户`。跨 actor / 不存在 → `JOB_NOT_FOUND`。GET Job 不内嵌 payload/result；succeeded 只给 `resultUrl`。result：queued/running → 409 `JOB_RESULT_NOT_READY`；expired → 410；failed/cancelled → 200 终态无下载；succeeded → JSON attachment。`GET /api/wallet/reconcile/runs` 仍为全局业务史。

### 终态审计

提交写 `wallet.reconcile.queued`（含 jobId/correlationId）。成功写既有 `wallet.reconcile`（jobId/runId/result）。failed/cancelled 分事件。queued 取消由 `JobService.Cancel` 直接记终态；running 取消由 `finalize-cancel` / scanner `recover-cancel` 走 terminal hook。HTTP 测试替身只在成功时写 `wallet.reconcile`（见 F-011），生产路径由 `wallet/jobs_test.go` 覆盖。

### migration 43 correlation

Apply 先把 `operation_log_correlation` 拷入 TEMP、DROP、重建 `operation_log`（扩大 CHECK）、再重建 correlation 表并 INSERT 备份。`TestMigrate0043PreservesOperationCorrelations`：42 头库写入带 correlation 的 `wallet.reconcile`，升 43 后 correlation 仍在，且三类新事件可插入。

### 错误码契约

D-002 §7 七码均在 catalog（双语 + messageKey）。handler 发出的五码进入 `frozenLiteralCodes`；两码作为 `frozenStoredCodes`（只存在于 Job representation）。`TestErrorCodeContractPinnedSet` / `TestErrorCatalogCoversFrozenCodesExceptInternal` / `docscheck` 本轮通过。Q2 权威增量：`docs/workspaces/workspace-007-localization-and-system-settings/GOAL-001-localization-and-system-settings/01-decision/D-002-s0-contract-freeze-info-gates.md`（可新增、不改既有语义）。

## 历史 findings disposition（A-001～A-011）

| ID | 原文 | level | 关闭声明 | 本条 |
|----|------|-------|----------|------|
| A-002 F-001～F-006 | 转换表 / lease / 幂等 / Profile / actor / 过期 | required | A-004 `fixed` | **保持 fixed**；现码与 D-002 一致 |
| A-002 F-007～F-008 | attempt 默认；202/审计时点 | recommended | A-004 `fixed` | **保持 fixed** |
| A-004 F-009 | recover-cancel 死状态 | required | A-006 `fixed` | **保持 fixed**；scanner + 测试覆盖 |
| A-006 F-010 | exhaust/recover-cancel guard 重叠 | recommended | A-007：exhaust 补 `cancel_requested=0` | **保持 fixed**；`repository.go:282-288` |
| A-001 F-001 | startup/周期扫描 | recommended | A-009 | **保持 fixed**；`runner.go:214-218` |
| A-001 F-002 | descriptor 双写 + 模块集 conformance | recommended | A-010/A-011 | **保持 fixed**；kernel mvp 列表 + composition 计数 + descriptorsMatch |

无开放 required。无 P-004 冲突。

## Findings

### F-011 · HTTP / JobService 读路径未直接测 ExpireIfDue → 410

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | low |
| status | open |
| 影响门禁 | 不阻断 S4 关门 |
| evidence | `modules/wallet/jobs.go:64-75`；`handler/wallet.go:362-376`；`handler/wallet_test.go:131-133,317-368`；`jobs/repository_test.go:84-90`；`jobs/runner_test.go:204` |

生产 `JobService.Get` 在 `succeeded` 上调用 `ExpireIfDue` 再重读；handler 将 `expired` 映射为 410 `JOB_RESULT_EXPIRED`。repository 与 runner 已证明原子 expire+清 result。缺口是：HTTP 测试替身 `walletJobTestService.Get` **省略** `ExpireIfDue`，也没有「TTL 后再 GET result」用例。这不否定契约已实现，也不否定成功标准 3；编排器若补一条 JobService/HTTP 410 用例可消除替身保真缺口。不升级为 required。

## 必改项汇总

无开放 required。

Recommended：F-011（可选补 HTTP/JobService 410 用例）。

## 与既有意见的异同

| 点 | A-011 (self) | 本条 (independent) |
|----|----------------|--------------------|
| D-002 实施与四条成功标准 | pass | **同意 pass** |
| 历史 F-001～F-010 | 均 fixed | **同意保持 fixed** |
| 全量测试 / Profile 边界 / 错误码 | pass | **本轮独立复跑后同意** |
| HTTP 410 | 记入「错误契约 / 过期」pass | 同意实现成立；另记 recommended F-011（替身未测读路径） |
| 开放 required | 0 | **0** |
| 放行 | 不单独放行，等本条 | **可放行 S4 / 候选关门**；状态变更仍由 `/govern` |

无意见冲突。

## 结论 + 建议给编排器/用户的下一步

**verdict = pass。** D-002 v0.2.0 已在当前代码落地。四条成功标准有可重复核对的实现与测试证据。A-001～A-011 的 required/recommended 关闭声明仍成立。I-001～I-004 无到期未关闭项。本条满足 I-004 要求的 S4 independent 关门审。

建议 `/govern`：响应本条 A-012；将 S4 标完成并把 GOAL-005 标 `done`（派生 progress 按 5/5 重算）；F-011 可在关门后作为非阻断补测，或忽略。不要用本意见直接改 `00-meta` / `goal-tree`——那是编排器的职责。

## 声明

本意见不修改 `status` / `progress` / 检查点 / `00-meta` / D-002 / 实现代码 / goal-tree。响应由 `/govern` 处理。
