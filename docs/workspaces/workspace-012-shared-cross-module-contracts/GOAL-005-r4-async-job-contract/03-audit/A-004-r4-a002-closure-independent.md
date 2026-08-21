---
id: A-004-r4-a002-closure-independent
goal: GOAL-005-r4-async-job-contract
doc: audit-entry
record_id: A-004
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: A-002 F-001～F-006 关闭复核；D-002 精确契约；I-002/I-003；S0→S1 门禁
audit_type: finding-closure
verdict: conditional
status: recorded
parent: GOAL-005-r4-async-job-contract
created: 2026-08-18
updated: 2026-08-18
version: 0.1.0
responds_to: A-002
reviews: A-003
---

# A-004 · A-002 F-001～F-006 关闭复核（2026-08-18）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high；项目级路径见 `docs/architecture/independent-audit-execution.md`）
- **类型**：ad-hoc / finding-closure
- **scope**：GOAL-005 S0；复核 A-002 required F-001～F-006 的关闭候选。只读 A-002、A-003、D-002 及必要现有 kernel / store / wallet 证据。判断每条是否可按 `fixed` 关闭、I-002/I-003 是否可转 `verified`、S0 是否可放行 S1。
- **verdict**：conditional

## 范围与区间

- **工作区**：`workspace-012-shared-cross-module-contracts`（`workspace.md`：`root_goal` = `GOAL-001-shared-cross-module-contracts`；`canonical_scope` 与本目标路径一致；`shared_materials_catalog: none`；`vision_role: delivery`；`primary_plan` = `VP-012-shared-cross-module-contracts`）
- **covered**：A-002 F-001～F-008 原文；A-003 候选响应；D-002 精确契约；I-002/I-003 门禁；对照现有 `CollectPersistence` / `RegisterContributions` / `profileDefaults` / `BuiltinModules` / `descriptorsMatch`、store 条件更新、`admin.wallet` provider/handler/reconcile/PK、R1 `writeLocalizedError` / errorcatalog
- **excluded**：S1～S4 实施与关门、未重新执行测试套件、其他工作区上下文、共享资料内容（目录为 `none`）、改写 D-002 / `00-meta` / status / progress / 实现代码
- **本轮未复验**：未运行 `go test`；「可实施」只核到现网模式与 D-002 谓词是否闭合，不把未写代码当成已交付

## 工作区与对齐（只读）

| 检查项 | 结论 | 证据 |
|--------|------|------|
| 工作区绑定 | 通过 | `workspace.md` Root / canonical / `plan_refs`+`primary_plan` 与 GOAL-005 `parent`、`primary_plan` 一致 |
| 共享资料引用 | 无引用，不构成关闭证据 | `shared_materials_catalog: none` |
| 既有 Goal 审计 | A-001 self = pass（0 required）；A-002 independent = conditional（6 required）；A-003 self = conditional（候选 fixed，未自闭） | `03-audit.md` |
| P-004 冲突 | 无 | A-003 未宣称已闭合；本条与 A-002 无相反 verdict 需用户裁 |

## 成果（有证据）

| 主张 | 证据 | 核验 |
|------|------|------|
| D-002 已写成 S0 精确契约冻结稿，并声明须经本复核才转 `accepted` | `01-decision/D-002-r4-precise-contract.md:1-14,102-104` | 通过：status 仍为 `proposed` |
| A-003 未自闭 A-002 required | `03-audit/A-003-r4-a002-response-self.md:25-29` | 通过：候选 fixed / pending independent |
| migration-only 模式仍可行 | `apps/api/internal/modules/compiled/persistence.go:23-41`；`kernel/persistence.go:37-69`；`kernel/provider.go:67-94`；`composition.go:134-142,349-353`；`corepersistence/migration/provider.go:12-20` | 通过：Persistence 与 Plan-gated runtime 仍分离；`core.persistence` 先例 `Register` 为空 |
| `core.jobs` 仍不在默认 Profile / BuiltinModules | `kernel/profile.go:25-114,158-206` 无 `core.jobs` / `core.persistence` | 通过 |
| `admin.wallet` 已在 ProfileAdmin 与 BuiltinModules；Descriptor 双写仍 fail-closed | `profile.go:90-92,199`；`wallet/provider.go:154-186`；`kernel/provider.go:107-135` `descriptorsMatch` | 通过：D-002 四条新 route key 必须两侧同改 |
| 条件 `UPDATE` + `RowsAffected` 仍可用 | `store/store.go:49,67-80`；`wallet/store/repository.go:320-346` | 通过 |
| `wallet_reconciliation_runs.id` 已是 `TEXT PRIMARY KEY`；UNIQUE 冲突可检测 | `wallet/migration/migration.go:56-64`；`repository.go:695-700,806-808` | 通过：Job ID 作 run ID / `ReconcileOnce` 主键决胜可落在现表，无需改 wallet schema |
| 现网 `Reconcile` / HTTP 仍非幂等、仍是同步 200 | `wallet/provider.go:128-134` 每次 `newID`；`repository.go:636-708` 每次 `INSERT`；`handler/wallet.go:287-307` `StatusOK` | 通过：属 S3 实施缺口，不否定 D-002 冻结 |
| 现网 list runs 仍全局、只查 `wallet.read` | `handler/wallet.go:310-333`；`repository.go:768-784` | 通过：与 D-002 §6「业务历史保持全局」一致 |
| R1 包络仍在；catalog 尚无 `JOB_*` / 410 | `handler/localize.go:17-32`；`errorcatalog/errorcatalog.go` 无 `JOB_` 键 | 通过：码已在 D-002 §7 命名；入 catalog 是 S3 |

## 对照成功标准（S0 关闭复审）

| 标准 | S0 状态 | 证据 |
|------|---------|------|
| 1. 六态与合法/非法转换可测 | 大部分已冻；**一处可达死状态** | D-002:37-57；F-009 |
| 2. 领取 / 不重复执行 | claim/reclaim/fencing/幂等已冻 | D-002:45-46,59-65,67-73 |
| 3. 202 + 轮询 + 410 | HTTP/TTL/410 已冻；死状态会挡住已提交 run 的 complete | D-002:35,55,73-81,87；F-009 |
| 4. Profile / 模块矩阵 / Manifest 不变 | 不变式与 route 清单已冻 | D-002:16-26 |

## 信息门禁核对（P-005）

| ID | 级别 | 最晚阶段 | 当前（00-meta） | 本轮结论 |
|----|------|----------|-----------------|----------|
| I-001 | required | S0 结束前 | verified | 保持；本轮不重开关 |
| I-002 | required | S0 结束前（影响 S1） | collecting | **可转 `verified`**。F-004 原文关闭条件已写入 D-002 §1；现网装配模式仍支持 |
| I-003 | required | S0 结束前（影响 S1/S2） | collecting | **不可转 `verified`**。F-001～F-003/F-005/F-006 原文项已写入，但 D-002 表存在可达且无事件的 `running` 死状态（F-009），精确转换仍未闭合 |
| I-004 | required | S1 实施前 | verified | 保持；不能代替 S4 关门独立审 |

无 `deferred`。无用户书面 `accepted-residual`。

## Disposition（保留 A-002 原文编号）

| ID | 原文标题（A-002） | level | 候选（A-003） | 本条 disposition | 关闭路径 |
|----|-------------------|-------|---------------|------------------|----------|
| F-001 | I-003 缺少可冻结的状态转换表与唯一决胜谓词 | required / high | candidate fixed | **fixed** | D-002 §3/§7 |
| F-002 | claim / lease / recovery 未定义 attempt 语义与 in-flight 互斥 | required / high | candidate fixed | **fixed** | D-002 §2/§4 |
| F-003 | wallet reconcile 未绑定「一次 Job ↔ 至多一次业务 run」 | required / high | candidate fixed | **fixed** | D-002 §5 + 现表 PK |
| F-004 | I-002 模式成立但未冻成可检查不变式 | required / med | candidate fixed | **fixed** | D-002 §1 + 现网装配 |
| F-005 | actor 隔离与 `wallet.read` 全局交叉未定义 | required / med | candidate fixed | **fixed** | D-002 §6 |
| F-006 | 结果过期缺少 TTL、触发器与原子清空 | required / med | candidate fixed | **fixed** | D-002 §2/§3/§7 |
| F-007 | retry 字段复位与 `max_attempts` 默认值未写 | recommended / med | 已写入 | **fixed** | D-002:32-35,54 |
| F-008 | 202 与审计时点未登记为 S3 必测 | recommended / low | 已写入 | **fixed** | D-002:73-81 |
| F-009 | （新）`cancel_requested=1` 且租约过期且 `attempt < max` 无恢复事件 | required / high | — | **open** | 阻断 I-003 / S1 |

`fixed` 仅表示 A-002 原文三项/四项关闭条件已成文且可重复核对，**不是**实施完成，也不是 I-003 已关闭。

## Findings

### F-001 · I-003 缺少可冻结的状态转换表与唯一决胜谓词

| 字段 | 值 |
|------|-----|
| level | required |
| severity | high |
| status | **closed** |
| disposition | **fixed** |
| 关联 | I-003 |
| evidence | `D-002-r4-precise-contract.md:37-57,83-87` |

**原文（A-002，保留）**：D-001 决策 3–7 只有散文，没有 from / event / guard / to / 字段突变 / 非法码。S1 会被迫发明 (1) complete/fail 与 cancel 的 `WHERE`；(2) 终态取消、expired 重试、`queued→expired`、`running→expired` 的非法码；(3) 进度类型/上界、reclaim 后是否保留、非 running 写进度的拒绝。

**复核**：三项均已写入。

1. `complete` = `running` + exact owner+version（即使 cancel 已置位也可胜出）；`fail` 另要求 `cancel_requested=0`；`finalize-cancel` 要求 `cancel_requested=1`；`RowsAffected != 1` 为稳定 transition error。`D-002:51-57`
2. 终态不再接受取消（`JOB_NOT_CANCELLABLE` 409）；仅 `failed` 且 `attempt < max_attempts` 可 retry（`JOB_NOT_RETRYABLE` 409）；v1 只有 `succeeded→expired`，无 `queued→expired` / `running→expired`。`D-002:35,54-55,85`
3. progress 为整数 0～100；running 更新 `new >= current && new <= 99`；complete 原子设 100；reclaim 保留 progress；不在表内的进度写入走 `RowsAffected != 1`。`D-002:33,46,49,51`

原缺口已闭合。D-002 表另有**新的**可达死状态，不回退本条，记 F-009。

### F-002 · claim / lease / recovery 未定义 attempt 语义与 in-flight 互斥

| 字段 | 值 |
|------|-----|
| level | required |
| severity | high |
| status | **closed** |
| disposition | **fixed** |
| 关联 | I-003 |
| evidence | `D-002-r4-precise-contract.md:32-34,45-47,59-65`；`store/store.go:49,67-80` |

**原文（A-002，保留）**：未冻 reclaim 是否 `attempt+1`、lease 续约、以及「重复领取不重复执行」的 in-flight 互斥。`SetMaxOpenConns(1)` 不能阻止双执行。

**复核**：

1. claim 与 reclaim **都** `attempt+1`；manual retry 只 requeue，下一次 claim 才加 attempt。`D-002:45-46,61`
2. lease=30s、heartbeat=10s；heartbeat / progress / terminal 均要求 exact `lease_owner+lease_version`。`D-002:34,48-51,63`
3. 同进程 active set 跳过仍在跑的 ID；跨进程重入明确交给业务 Job-ID 幂等。`D-002:62-64`

原缺口已闭合。reclaim 递增 attempt 会使连续崩溃在 `attempt >= max` 后无法 manual retry——这是已冻结选择，不是未定义项。

### F-003 · wallet reconcile 真实消费未绑定「一次 Job ↔ 至多一次业务 run」

| 字段 | 值 |
|------|-----|
| level | required |
| severity | high |
| status | **closed** |
| disposition | **fixed** |
| 关联 | I-003 |
| evidence | `D-002-r4-precise-contract.md:67-74`；`wallet/migration/migration.go:56-64`；`wallet/store/repository.go:636-708,806-808`；`wallet/provider.go:128-134` |

**原文（A-002，保留）**：`Reconcile` / `ReconcileRun` 每次 `INSERT` 新行；须冻 (1) 同一 Job 至多一条 run；(2) 取消胜出但 run 已落库的语义；(3) result 投影与 `inconsistent` ≠ Job `failed`。

**复核**：

1. Job ID = reconcile run ID；`ReconcileOnce` 遇已存在 runID 回读；并发插入靠主键决胜。现表 `id TEXT PRIMARY KEY` + `isUniqueViolation` 使该规则可实施且「无需改 wallet 表」成立。`D-002:69-70`；`migration.go:56-57`；`repository.go:806-808`
2. 活 owner：run 已提交后 complete 可胜出；「不伪称 cancelled」。`D-002:51,73`
3. `consistent`/`inconsistent` 都是 Job `succeeded`；result JSON 字段已列。`D-002:71-72`

现网 `Reconcile` 仍每次 `newID`+`INSERT`（`provider.go:128-134`；`repository.go:695-700`）是 S3 事实，不是 S0 未冻。活 owner 规则成立；**死 owner + cancel 旗标**下 complete 无法发出，由 F-009 约束恢复事件必须服从本条第 2 项。

### F-004 · I-002 所有权模式成立，但未冻成可检查的 Profile/Manifest 不变式

| 字段 | 值 |
|------|-----|
| level | required |
| severity | med |
| status | **closed** |
| disposition | **fixed** |
| 关联 | I-002 |
| evidence | `D-002-r4-precise-contract.md:16-26`；`compiled/persistence.go:23-41`；`kernel/provider.go:67-94,107-135`；`kernel/profile.go:25-114,158-206,199`；`wallet/provider.go:162-186`；`composition.go:134-142,349-353`；`corepersistence/migration/provider.go:12-20` |

**原文（A-002，保留）**：须写入 (1) `core.jobs` 只进 PersistenceProviders、空 `Register`、不进 Profile/BuiltinModules/runtime 表；(2) wallet 路由键清单与 Descriptor 双写；(3) 无新 page/nav/fragment/permission，权限仍 `wallet.read`；(4) 「模块矩阵不变」= 模块 ID 集与装配算法不变，wallet 路由内容扩展是显式例外。

**复核**：D-002 §1 四条均已成文。现网仍支持该模式：`PersistenceProviders()` 与 `plan.HasModule` runtime 分离；`profileDefaults` / `BuiltinModules()` 无 `core.jobs`；`descriptorsMatch` 仍 fail-closed；`admin.wallet` 已在 ProfileAdmin（`profile.go:90-92,199`）。四条新 key：

- `GET /api/wallet/jobs/{id}`
- `POST /api/wallet/jobs/{id}/cancel`
- `POST /api/wallet/jobs/{id}/retry`
- `GET /api/wallet/jobs/{id}/result`

`POST /api/wallet/reconcile` 保留 path、改 202，不是新模块 ID。测试锁默认模块集仍是 S3/S4（A-001 F-002），不阻断本条关闭。

**I-002 可转 `verified`。**

### F-005 · actor 隔离与 `wallet.read` 全局面的交叉未定义

| 字段 | 值 |
|------|-----|
| level | required |
| severity | med |
| status | **closed** |
| disposition | **fixed** |
| 关联 | I-003 |
| evidence | `D-002-r4-precise-contract.md:76-81`；`handler/wallet.go:288-334`；`wallet/store/repository.go:768-784` |

**原文（A-002，保留）**：未冻 GET/cancel/retry 是仅提交 actor 还是任何 `wallet.read`；跨 actor 404 还是 403；业务 runs 全局与 Job 按 actor 如何对应；系统 actor 是否豁免。

**复核**：五个 wallet Job endpoint 先 `wallet.read`，再 `job.kind == wallet.reconcile && job.actor_id == 当前认证用户`；跨 actor / 非 wallet kind / 不存在统一 404 `JOB_NOT_FOUND`；HTTP 无 system actor 豁免；`GET /api/wallet/reconcile/runs` 保持现网全局业务历史，不授予 Job 查询权。`D-002:76-80`。现网 list 确为全表分页（`repository.go:782-784`），与「二者并存」一致。客户端对应关系：Job result 含 run `id`，runs 列表仍是独立业务史。

### F-006 · 结果过期缺少 TTL、触发器与原子清空语义

| 字段 | 值 |
|------|-----|
| level | required |
| severity | med |
| status | **closed** |
| disposition | **fixed** |
| 关联 | I-003 |
| evidence | `D-002-r4-precise-contract.md:35,55,83-87`；`handler/localize.go:17-32`；`errorcatalog/errorcatalog.go`（无 `JOB_`） |

**原文（A-002，保留）**：未冻 `expires_at` 写入时机与 TTL、谁扫描、expire+clear 原子性、410 catalog 码、`cancelled`/`failed` 是否过期、`expired` 可否 retry。

**复核**：成功结果保留 24h，`expires_at = finished_at + 24h`；v1 仅 `succeeded` 转 `expired`；`expire` 与清 `result` 同一条件 UPDATE；GET Job/result 读前 `ExpireIfDue`，周期 scanner 同样调用；410 = `JOB_RESULT_EXPIRED` + R1 包络；failed/cancelled 不转 expired、无下载、不可把 expired 当 retry 源。`D-002:35,55,85-87`。码入 `errorcatalog.Catalog` 是 S3，S0 只要求命名。

GET result 在 failed/cancelled 上「返回对应终态但无下载」的 HTTP 状态未写成单一数字；不回退本条，S3 实现时用 200 表征 + 无 attachment 即可，不必再开 required。

### F-007 · retry 的字段复位表与 `max_attempts` 默认值未写（recommended）

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | med |
| status | **closed** |
| disposition | **fixed** |
| evidence | `D-002-r4-precise-contract.md:32-35,54` |

attempt 初值 0、`max_attempts=3`；retry 清 cancel/lease/result/error/finished/expires，progress=0，保留 id/payload/actor/correlation/attempt/max。可关。

### F-008 · `POST /api/wallet/reconcile` 202 与现网 200/审计时点未登记为 S3 必测（recommended）

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | low |
| status | **closed** |
| disposition | **fixed** |
| evidence | `D-002-r4-precise-contract.md:73-81`；`handler/wallet.go:304-306`；`handler/wallet_test.go:222-234,254` |

提交写 `wallet.reconcile.queued`；成功终态写既有 `wallet.reconcile`；failed/cancelled 分事件；旧 200 测试改为 202+轮询+结果。现网仍断言 200 / `wallet.reconcile`（`wallet.go:306`；`wallet_test.go:224-225,254`）是 S3 必改事实，登记已够。

### F-009 · `running + cancel_requested=1 + lease expired + attempt < max_attempts` 无恢复事件

| 字段 | 值 |
|------|-----|
| level | required |
| severity | high |
| status | open |
| 影响门禁 | I-003；成功标准 1–3；S0 冻结；S1 实施 |
| 关联 | I-003；约束 F-003 第 2 项在恢复路径上可执行 |
| evidence | `D-002-r4-precise-contract.md:45-47,51-53,62-65,73` |

D-002 表对 `running` 的 scanner / 恢复事件只有：

| event | guard | to |
|-------|-------|----|
| reclaim | lease expired **且 `cancel_requested=0`** 且 attempt < max | running（新 owner，attempt+1） |
| exhaust | lease expired 且 attempt >= max | failed |
| finalize-cancel | **exact owner+version** 且 cancel=1 | cancelled |
| complete | **exact owner+version** | succeeded |

下列状态**可达**且无任何事件：

1. `submit` → `claim`（attempt=1）→ `request-cancel running`（`cancel_requested=1`，status 仍为 running）
2. handler / 进程在 `complete` 或 `finalize-cancel` 之前崩溃；进程重启后 active set 为空
3. 租约过期后：`reclaim` 被 `cancel_requested=1` 挡住；`exhaust` 因 `attempt < max` 不成立；`finalize-cancel` / `complete` 需要已消失的 fencing token

结果：Job **永久停在 `running`**。跨重启也不会前进，因为 reclaim 被禁止，attempt 不再增加，exhaust 永远不会触发。

若崩溃发生在 `ReconcileOnce` 已提交之后，还违反 D-002 自己写的 F-003 规则：业务 run 已在 `wallet_reconciliation_runs`（`migration.go:56-64`；`repository.go:695-700`），Job 却不能 `complete`，GET result 会按 `JOB_RESULT_NOT_READY` 对待仍为 running 的行（`D-002:87`）。

这不是 A-002 原文三项的重复，而是 D-002 在补表时引入的**表闭合缺口**。S1 若按现表开工，实现者必须自行发明恢复事件；不同选择会改写成功标准 1–3 与 F-003「不伪称 cancelled」。

关闭本条至少要在 D-002 增加一条可测恢复事件，且必须同时满足：

- lease 过期后，**没有**活 fencing owner 时，Job 不能留在 `running`
- 若该 Job ID 已有 `wallet_reconciliation_runs` 行 → 不得写成 `cancelled`（服从 F-003 / `D-002:73`），应能 `complete` 或等价 adopt-success
- 若尚无 run → 可原子转 `cancelled`

示例（供编排器写入，本意见不改 D-002）：允许 `reclaim` 在 `cancel_requested=1` 时仍换 fencing，新 handler 见已有 run 则 `complete`、见 cancel 且无 run 则 `finalize-cancel`；或 scanner 在 lease 过期后按「run 是否存在」做一条不依赖旧 owner 的条件 UPDATE。

在本条按 `fixed` / `accepted-residual` / `user-overruled` 合法闭合前，**不得**把 I-003 标 `verified`，**不得**放行 S1。

## 必改项汇总

1. **F-009（唯一开放 required）**：为 `running ∧ cancel_requested=1 ∧ lease_expires_at <= now ∧ attempt < max_attempts` 补一条恢复事件；并写明已提交 wallet run 时不得伪称 `cancelled`。

F-001～F-006 原文项本条同意 `fixed`。I-002 可转 `verified`。I-003 与 S1 仍被 F-009 阻断。

## 与既有意见的异同

| 点 | A-002 | A-003 | 本条 |
|----|-------|-------|------|
| F-001～F-006 原文缺口 | 开放 required | 候选 fixed，不自闭 | **同意关闭为 fixed** |
| I-002 | 保持 collecting | 保持 collecting | **可 verified** |
| I-003 | 保持 collecting；不得 S1 | 保持 collecting；不得 S1 | **仍 collecting**；不得 S1 |
| D-002 表是否完整 | （当时无表） | 声称精确转换表已成文 | 表覆盖 A-002 原文；**新发现可达死状态 F-009** |
| 冲突 | — | 未宣称已闭合 | 无 P-004 冲突 |

A-001 两条 recommended（S2 扫描路径、S3/S4 descriptor conformance）仍是实施门禁，不被本条改写。

## 结论 + 建议给编排器/用户的下一步

**verdict = conditional。** A-002 F-001～F-006（及 F-007/F-008）可按 `fixed` 关闭。I-002 可转 `verified`。I-003 **不可**转 `verified`。S0 **不可**放行 S1。

建议 `/govern`：响应本条 A-004；将 A-002 F-001～F-006 标 `fixed`；把 F-009 补进 D-002（恢复事件 + 已提交 run 不得伪 cancel）；再请独立复核 F-009 后，才把 I-003 标 `verified`、D-002 转 `accepted` 并进入 S1。

## 声明

本意见不修改 `status` / `progress` / 检查点 / `00-meta` / `D-002` / 实现代码，也不自行关闭目标状态。响应由 `/govern` 处理。
