---
id: D-002-r4-precise-contract
goal: GOAL-005-r4-async-job-contract
status: proposed
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-005-r4-async-job-contract
version: 0.1.0
responds_to: A-002
---

# D-002 · R4 S0 精确契约冻结稿

本稿响应 A-002 F-001～F-006；经 independent 复核后才转 `accepted`。

## 1. 数据与所有权不变式

- migration 42 由 migration-only `core.jobs` provider 拥有，只加入 `compiled.PersistenceProviders()`；其 `Register` 为空。
- `core.jobs` **不得**进入 `profileDefaults`、`BuiltinModules()` 或 composition 的 Plan-gated runtime provider 列表。
- runtime 位于 `internal/jobs`，不是 module/provider；不新增 page、navigation、fragment、permission 或默认模块 ID。
- `admin.wallet` 只新增四个 route keys，provider Descriptor 与 `kernel.BuiltinModules()` 必须同步：
  - `GET /api/wallet/jobs/{id}`
  - `POST /api/wallet/jobs/{id}/cancel`
  - `POST /api/wallet/jobs/{id}/retry`
  - `GET /api/wallet/jobs/{id}/result`
- `POST /api/wallet/reconcile` 保留 path、改为 202；权限仍全部使用 `wallet.read`。允许的 catalog 变化仅是上述 wallet route 内容扩展；模块 ID 集、Profile 默认集和 Manifest 装配算法/fragment 集保持不变。

## 2. Job 行与默认值

`id, kind, status, payload, progress, cancel_requested, attempt, max_attempts, lease_owner, lease_version, lease_expires_at, result, error_code, error_message, actor_id, correlation_id, created_at, updated_at, finished_at, expires_at`。

- 新 Job：`queued`、progress=0、attempt=0、max_attempts=3、lease_version=0；payload/actor/correlation 不可变。
- progress 为整数 0～100；running 更新只允许 `new >= current && new <= 99`，完成原子设 100。
- lease=30s、heartbeat=10s、恢复扫描=10s，均由 runner option 注入以便确定性测试；生产默认固定如上。
- 成功结果保留 24h；`expires_at = finished_at + 24h`。v1 只有 `succeeded` 转 `expired`，failed/cancelled 无下载结果且不转 expired。

## 3. 精确转换表

所有 repository 事件均为单条条件 UPDATE 或单事务内的条件 UPDATE；`RowsAffected != 1` 返回稳定 transition error。

| event | from + guard | to | 原子字段变化 |
|-------|--------------|----|--------------|
| submit | new id | queued | 写不可变字段与默认值 |
| claim | queued, attempt < max_attempts | running | attempt+1；new lease_owner；lease_version+1；lease_expires_at=now+lease；清 finished/expires/error |
| reclaim | running, lease_expires_at <= now, cancel_requested=0, attempt < max_attempts | running | attempt+1；new lease_owner；lease_version+1；续写 lease；保留 progress |
| exhaust | running, lease expired, attempt >= max_attempts | failed | error_code=`JOB_ATTEMPTS_EXHAUSTED`；finished_at=now；清 lease |
| heartbeat | running, exact lease_owner+lease_version | running | lease_expires_at=now+lease |
| progress | running, exact owner+version, new>=current, new<=99 | running | progress=new |
| request-cancel queued | queued, actor match | cancelled | finished_at=now；清 lease；cancel_requested=0 |
| request-cancel running | running, actor match | running | cancel_requested=1；runner cancel signal best-effort |
| complete | running, exact owner+version | succeeded | progress=100；result；finished_at=now；expires_at=now+24h；清 lease/error/cancel flag |
| fail | running, exact owner+version, cancel_requested=0 | failed | error fields；finished_at=now；清 lease/result |
| finalize-cancel | running, exact owner+version, cancel_requested=1 | cancelled | finished_at=now；清 lease/result/error |
| retry | failed, actor match, attempt < max_attempts | queued | progress=0；清 cancel/lease/result/error/finished/expires；保留 id/payload/actor/correlation/attempt/max |
| expire | succeeded, expires_at <= now | expired | 同一 UPDATE 清 result 并写 updated_at；保留 finished/expires |

`complete` 表示 handler 已提交业务结果，因此即使 cancel flag 已置位也允许它用 owner+fencing token 赢得终态；取消是 best-effort，不承诺回滚已提交副作用。非成功返回时 runner 重新读取 cancel flag：为 1 则 `finalize-cancel`，否则 `fail`；二者均由条件 UPDATE 唯一决胜。终态不再接受取消；只有 failed 可 retry；expired 不可 retry。

## 4. claim、fencing 与恢复

- 每次 claim/reclaim 都代表一次 handler 执行并增加 attempt；manual retry 仅 requeue，下一次 claim 才增加 attempt。
- runner 启动时立即扫描，之后每 10s 扫描 queued 与 lease-expired running；同进程 active Job set 会跳过仍在运行的 ID。
- handler 运行期间 heartbeat 续租。每次 claim 生成随机 lease_owner 并增加 lease_version；progress/heartbeat/terminal 更新均要求 owner+version，旧执行者不能写 Job 状态。
- 多进程或暂停超过 lease 的 handler 仍可能重入，因此业务 handler **必须按 Job ID 幂等**；fencing 保护 Job 状态，消费方幂等保护副作用。
- lease-expired 且 attempt 已耗尽的 Job 由 scanner 原子转 failed；不会永久停在 running。

## 5. wallet 消费幂等与结果

- `wallet.reconcile` payload 只含可选 accountId；Job ID 同时作为 reconcile run ID。
- 新增 `ReconcileOnce(accountID, runID, actorID, now)`：相同 runID 已存在时返回既有 run；并发插入由主键唯一约束决胜，败者 reload。一个 Job 至多一行 `wallet_reconciliation_runs`，无需改 wallet 表。
- `consistent` 与 `inconsistent` 都表示对账成功，Job 均为 `succeeded`；只有执行/存储错误令 Job failed。
- result 是稳定 reconcile run JSON：`id/accountId/result/mismatchCount/details/actorId/createdAt`。result endpoint 以 JSON attachment 返回。
- handler 在 run 持久化后返回成功；此后取消到达时 complete 可胜出。若 cancel 在业务调用前被观察到则不执行；若 run 已提交，Job 不伪称 cancelled。
- 提交时记录 `wallet.reconcile.queued`（jobId/correlationId）；成功终态记录既有 `wallet.reconcile`（jobId/runId/result），failed/cancelled 分别记录终态事件。旧同步 200 测试改为 202+轮询+结果断言。

## 6. actor 与 HTTP

- 所有五个 wallet Job endpoint 先要求 `wallet.read`；再要求 `job.kind == wallet.reconcile && job.actor_id == 当前认证用户`。
- 跨 actor、非 wallet kind与不存在统一返回 404 `JOB_NOT_FOUND`，不泄露存在性。HTTP 不提供 system actor 豁免。
- `GET /api/wallet/reconcile/runs` 保持现有 wallet.read 全局业务历史；它不授予 Job 查询权，也不改变 actor predicate。
- `POST reconcile` 返回 202 Job representation；GET Job 返回状态/进度/attempt/max/cancelRequested/error/resultUrl，不内嵌 payload/result。

## 7. 错误与过期

所有错误走 R1 localized envelope：`JOB_NOT_FOUND` 404、`JOB_NOT_CANCELLABLE` 409、`JOB_NOT_RETRYABLE` 409、`JOB_RESULT_NOT_READY` 409、`JOB_RESULT_EXPIRED` 410、`JOB_ATTEMPTS_EXHAUSTED`（持久化内部终态码）、`JOB_HANDLER_FAILED`（持久化内部终态码）。

GET Job/result 在读取前调用同一 repository `ExpireIfDue`；周期 scanner 也调用它。expire 与 result 清空为同一条件 UPDATE。result 对 queued/running 返回 `JOB_RESULT_NOT_READY`，failed/cancelled 返回对应终态但无下载，expired 返回 410。

## 8. A-002 disposition

| finding | 响应 |
|---------|------|
| F-001 required | 第 3/7 节补齐 from/event/guard/to、字段变化与错误码 |
| F-002 required | 第 2/4 节冻结 attempt、heartbeat、fencing、active set 与恢复 |
| F-003 required | 第 5 节用 Job ID 作为 run ID，冻结一次 Job 至多一次业务 run |
| F-004 required | 第 1 节冻结 migration-only/Profile/Descriptor 不变式与 route 清单 |
| F-005 required | 第 6 节冻结 actor predicate、404 与全局业务历史边界 |
| F-006 required | 第 2/3/7 节冻结 TTL、触发器、原子清空、410 与仅 succeeded 过期 |
| F-007 recommended | 第 2/3 节冻结默认 attempt/max 与 retry 字段复位 |
| F-008 recommended | 第 5/6 节登记 202、轮询/result 与 audit 时点 |

## 门禁

A-002 required 当前仅为“候选 fixed”，须经 independent 复核。复核通过后：D-002 转 `accepted`，I-002/I-003 转 `verified`，S0 才完成并放行 S1。
