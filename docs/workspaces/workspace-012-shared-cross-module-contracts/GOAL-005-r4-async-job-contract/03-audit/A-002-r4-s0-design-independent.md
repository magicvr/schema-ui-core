---
id: A-002-r4-s0-design-independent
goal: GOAL-005-r4-async-job-contract
doc: audit-entry
record_id: A-002
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: R4 S0 · D-001/E-001/I-001～I-004；migration-only core.jobs owner 与 Profile/Manifest 边界；queued/running/succeeded/failed/cancelled/expired；claim/lease/recovery；取消与完成竞态；retry attempt；结果过期；actor 隔离；wallet reconcile 真实消费
audit_type: design-plan+execution-facts
verdict: conditional
status: recorded
parent: GOAL-005-r4-async-job-contract
created: 2026-08-18
updated: 2026-08-18
version: 0.1.0
---

# A-002 · R4 S0 设计/计划 + 扫描事实独立审计（2026-08-18）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high；项目级路径见 `docs/architecture/independent-audit-execution.md`）
- **类型**：stage（S0）/ design-plan + execution-facts
- **scope**：GOAL-005 S0 仅限 `00-meta.md`、`D-001`、`E-001`、`A-001`，以及必要的现有 API kernel / store / wallet 代码。核验 migration-only `core.jobs` owner 是否保持 Profile/Manifest 边界、六态转换、claim/lease/recovery、取消与完成竞态、retry attempt、结果过期、actor 隔离、wallet reconcile 真实消费。
- **verdict**：conditional

## 范围与区间

- **工作区**：`workspace-012-shared-cross-module-contracts`（`workspace.md`：`root_goal` = `GOAL-001-shared-cross-module-contracts`；`canonical_scope` 与本目标路径一致；`shared_materials_catalog: none`；`vision_role: delivery`；`primary_plan` = `VP-012-shared-cross-module-contracts`）
- **covered**：GOAL-005 定义与非目标、S0 路线图、D-001 提议、E-001 扫描事实、I-001～I-004、A-001 self、对照 `kernel.CollectPersistence` / `RegisterContributions` / `profileDefaults` / `BuiltinModules`、store 事务与连接池、`admin.wallet` provider/handler/reconcile、`task_runs`、import/export 同步路径
- **excluded**：S1～S4 实施与关门、未重新执行测试套件、其他工作区上下文、共享资料内容（目录为 `none`）、通用 Job 管理 UI、外部队列、把 scheduled task 改造成 Job
- **本轮未复验**：未运行 `go test` / Web 测试；「测试覆盖」仅核到测试代码与断言存在，运行结果标为证据不足

## 工作区与对齐（只读）

| 检查项 | 结论 | 证据 |
|--------|------|------|
| 工作区绑定 | 通过 | `workspace.md` Root / canonical / `plan_refs`+`primary_plan` 与 GOAL-005 `parent`、`primary_plan` 一致 |
| 共享资料引用 | 无引用，不构成关闭证据 | `shared_materials_catalog: none`；GOAL-005 未引用 `material_id`/`sha256` |
| 对齐链 | 未发现与 Root R4 / VP-012 方向的明显冲突 | VP-012 方向级范围含六态、进度、重试、取消、结果下载/过期；Root R4 依赖 R1/R3；GOAL-005 以 wallet reconcile 为有界真实消费。非目标排除新默认 Profile 模块 ID、Manifest 默认集、Tier D、外部队列 |
| Vision Review required | 本 scope 未见开放 required | 本意见不审 Vision Review 本身；未把 reviews 当作 Goal 放行证据 |
| 既有 Goal 审计 | A-001 self = pass，开放 required = 0；声明不放行 S1 | `03-audit.md`；`A-001-r4-s0-design-self.md` |

## 成果（有证据）

| 主张 | 证据 | 核验 |
|------|------|------|
| 当前无通用 jobs 表 / repository / runner | 全仓 `CREATE TABLE job` / `jobs` 无匹配；`compiled/persistence.go` PersistenceProviders 仅 account/dictionary/captcha/datapermission/mfa/recycle/wallet/scheduledtasks/notifications/auth/corepersistence/operationlog/settings | 通过 |
| 迁移最大版本为 41，下一全局版本可用 42 | `operationlog/migration/migration.go:354` `Version: 41`（`operation_log_correlation`）；无 Version ≥ 42 | 通过 |
| `task_runs` 只有同步终态 `ran/failed` | `scheduledtasks/migration/migration.go:29-32` `CHECK (status IN ('ran','failed'))`；无 queued/running/lease/attempt/progress/cancel | 通过 |
| scheduled task 手动触发是同步请求 | `handler/scheduledtasks.go:424-444` `runner.Execute` 后直接 `204` | 通过 |
| import/export 与文件下载为同步请求 | `handler/import.go:42-52`、`handler/export.go:42-48` 在同一请求内处理；Web 无 reconcile/job 轮询（`apps/web/src` 无 `reconcile` 匹配） | 通过 |
| wallet reconcile 已是独立、可持久化的业务 run | `wallet/migration/migration.go:56-64` `wallet_reconciliation_runs`；`Service.Reconcile` `provider.go:128-134`；`Repository.ReconcileRun` `repository.go:636-708` 每次 `INSERT` 新 run | 通过：适合作首个消费方；**不**具备 Job 幂等 |
| `admin.wallet` 已有 profile contribution 与 `wallet.read` | `kernel/profile.go:90-92` ProfileAdmin 含 `admin.wallet`；`wallet/provider.go:162-179` 已声明 `POST /api/wallet/reconcile` 与 `wallet.read`；`handler/wallet.go:288-289` 权限为 `wallet.read` | 通过 |
| Persistence 与 Profile 装配分离；migration-only owner 可不进默认模块集 | `compiled/persistence.go:1-3,23-41`；`kernel/persistence.go:37-68` 收集全部 compiled persistence；`composition.go:134-138` store 走 PersistenceCatalog，不看 Plan；`kernel/provider.go:67-90` 未启用 provider 不注册 surface；`provider_test.go:363-374` | 通过：I-002 **模式可行** |
| 已有 migration-only 先例，且不在默认 Profile | `corepersistence/migration/provider.go:12-20` `ModuleID = core.persistence`、`Register` 空实现；`profile.go:26-114` `profileDefaults` 与 `BuiltinModules()` 均无 `core.persistence` / `core.jobs` | 通过 |
| `admin.wallet` runtime 本身不贡献 migration | `wallet/provider.go:184-186` `CompiledPersistence() (nil, nil)`；表由 `wallet/migration` owner 管理 | 通过：与「runtime ≠ migration owner」一致 |
| store 可用条件 `UPDATE` + `RowsAffected` 做原子决胜 | `store/store.go:49` `SetMaxOpenConns(1)`；`67-80` `WithTx`；`wallet/store/repository.go:320-346` `WHERE id=? AND version=?`，`affected==0` → 冲突 | 通过：A-001 此点成立。单连接**不能**单独保证「handler 仍在跑时不会被二次 claim」 |
| 现网 reconcile HTTP 是同步 200，不是 202 Job | `handler/wallet.go:287-307` `writeJSON(..., StatusOK, reconcileRunToMap(*run))`；`wallet_test.go:222-234` 断言 `200` 且 body `result==consistent` | 通过：D-001 决策 8 是明确破坏性变更，尚未实施 |
| I-001 扫描结论可独立复核 | 上表 + E-001 | 通过；保持 `verified` |
| I-004 模式/provider 可唯一判定 | `00-meta` I-004；`docs/architecture/independent-audit-execution.md:23-27` grok-build / grok-4.6 / `/audit` | 通过；本条不能代替 S4 关门独立审 |

## 对照成功标准（S0 适用部分）

GOAL-005 四条成功标准均属 S1–S4 交付物。S0 只评估「是否已具备冻结契约 / 进入实施的信息」。

| 标准 | S0 状态 | 证据 |
|------|---------|------|
| 1. 六态与合法/非法转换、单调进度、attempt/lease、取消、重试、过期均有确定性测试 | 未开始；转换谓词未冻到可测 | D-001:17-21 只有散文，无 from/event/guard/to 表（F-001/F-002/F-006/F-007） |
| 2. Job 持久化；runner 可领取 `queued` 或租约过期的 `running`；重复领取不重复执行 | 未开始；「不重复执行」在 reclaim + 仍在跑的 handler 上无互斥定义 | D-001:18；`SetMaxOpenConns(1)` 只串行化连接，不串行化 goroutine 业务副作用（F-002/F-003） |
| 3. wallet reconcile 提交 202 Job；可轮询并读最终结果；过期 410 | 未开始；现网仍是 200 + run body | `wallet.go:306`；过期 TTL/扫描/错误码未冻（F-006）；202 破坏现有测试（F-008） |
| 4. API/迁移测试通过；Profile 默认集、模块矩阵和 Manifest 装配语义不变 | 所有权模式可保持默认模块集不变；wallet 新路由键会改 catalog 声明，D-001 未列出 | `profile.go:26-114` 无 `core.jobs`；`descriptorsMatch` `provider.go:107-135`（F-004） |

## 信息门禁核对（P-005）

| ID | 级别 | 最晚阶段 | 状态 | 是否到期 | 本轮结论 |
|----|------|----------|------|----------|----------|
| I-001 | required | S0 结束前 | verified | 本轮复核扫描事实 | **扫描判断成立**，E-001 主干可重复核对。`verified` 不得被读成 I-002/I-003 已关闭 |
| I-002 | required | S0 结束前（影响 S1 实施） | collecting | **阻断 S0 冻结 / S1** | 所有权**模式**已独立证实可行；关闭前必须把「不进 Profile/BuiltinModules + 仅 PersistenceProviders + wallet 路由键双写清单」写入 D-001（F-004）。本条不把 I-002 标 `verified` |
| I-003 | required | S0 结束前（影响 S1/S2 方案冻结） | collecting | **阻断 S0 冻结 / S1** | D-001 初版散文**不是**可冻结 transition table。取消/完成竞态谓词、claim vs reclaim 的 attempt/lease、handler 重入、过期操作语义、actor 隔离范围均未精确到可实现（F-001/F-002/F-003/F-005/F-006）。**不得**在本意见后关闭 I-003 |
| I-004 | required | S1 实施前 | verified | **未阻断 S0 设计审计** | data/migration/compatibility 按 P-003 可唯一判定 `independent`。provider = 项目级 grok-build（grok-4.6 reasoning high）。本条 **不能**关闭 S4 关门独立审 |

无 `deferred` 项。无用户书面 `accepted-residual`。

## Findings

### F-001 · I-003 缺少可冻结的状态转换表与唯一决胜谓词

| 字段 | 值 |
|------|-----|
| level | required |
| severity | high |
| status | open |
| 影响门禁 | I-003；S0 方案冻结；S1 实施 |
| 关联 | I-003 |
| evidence | `01-decision/D-001-r4-job-contract.md:17-21`；`00-meta.md:48,60`；`03-audit/A-001-r4-s0-design-self.md:30` |

D-001 决策 3–7 列出了六态、终态、`queued→running` claim、`queued→cancelled`、`running` 先写 `cancel_requested`、失败可重试、到期转 `expired`。这覆盖 VP-012 方向级名词，但 I-003 要求的是**精确状态转换**，而当前文件没有一张可实施、可测的表。

S1 若现在开工，实现者必须自行发明至少下列谓词，而不同选择会改变成功标准 1–2 的含义：

1. **完成/失败 vs 取消**：`running` 的取消只把 `cancel_requested` 置位、不改 `status`。若 `Complete`/`Fail` 的 `UPDATE` 只要求 `status=running`，则业务已完成后取消请求会被覆盖为 `succeeded`/`failed`。若还要求 `cancel_requested=0`，则取消请求在协作取消完成前就能挡住终态。D-001:19 只写「数据库条件更新唯一决胜」，没有写出两条 `UPDATE` 的 `WHERE`。
2. **非法转换与错误码**：`succeeded/failed/cancelled/expired` 上的取消、对 `expired` 重试、`queued→expired`、`running→expired`（跳过执行终态）均未列。R1 包络已存在（`handler/localize.go:17-32` 写入 catalog 码 + correlation），但 Job 非法转换码未进 D-001。
3. **进度**：只写「仅 `running` 单调增加」，未写类型/上界、reclaim 后是否保留、非 running 写进度的拒绝码。

在补齐 from / event / guard / to / 字段突变 / 非法码之前，**不得关闭 I-003，不得开始 S1**。

### F-002 · claim / lease / recovery 未定义 attempt 语义与 in-flight 互斥

| 字段 | 值 |
|------|-----|
| level | required |
| severity | high |
| status | open |
| 影响门禁 | I-003；成功标准 2；S1/S2 |
| 关联 | I-003 |
| evidence | `D-001-r4-job-contract.md:18`；`store/store.go:49,67-80`；`A-001-r4-s0-design-self.md:27,36` |

D-001:18：`queued→running` 原子 claim **并增加 attempt、写 lease**；租约过期的 `running` 可被重新 claim。未冻结：

1. **reclaim 是否再 `attempt+1`**。若是，长操作因进程重启/租约过短即可耗尽 `max_attempts`，与「仅 `failed` 且 `attempt < max_attempts` 才重试」（D-001:20）叠加后语义混乱。若否，须写明「仅从 `queued` 的 claim 增加 attempt；同一次执行的 lease 续约/reclaim 不增加」。
2. **lease 续约**。无 heartbeat 时，handler 仍在执行、lease 已过期，周期扫描（A-001 F-001 已提醒必须存在）会第二次 claim。`SetMaxOpenConns(1)` 只让两个 `WithTx` 排队，**不能**阻止第一个 `ReconcileRun` 提交后第二个 handler 再插入另一条 run。
3. **进程内互斥**。成功标准 2 要求「重复领取不重复执行」。至少要冻其一：lease 必须在运行中续约且 reclaim 不得发给仍持有 lease_owner 的活 handler；或 Job id 有 in-process 执行锁，扫描跳过仍在跑的 owner。

A-001 F-001 只把 startup/周期扫描列为 recommended 实现门禁。独立意见：**扫描路径是必要的，但没有互斥谓词则扫描本身会制造双执行**。这是 S0 契约缺口，不是 S2 才第一次出现的实现细节。

### F-003 · wallet reconcile 真实消费未绑定「一次 Job ↔ 至多一次业务 run」

| 字段 | 值 |
|------|-----|
| level | required |
| severity | high |
| status | open |
| 影响门禁 | I-003；成功标准 2–3；S3 消费 |
| 关联 | I-003 |
| evidence | `wallet/store/repository.go:636-708`（每次 `INSERT` 新 `wallet_reconciliation_runs`）；`wallet/provider.go:128-134`（每次预生成新 `runID`）；`handler/wallet.go:287-307`（同步调用后立即 200 + 审计）；`D-001-r4-job-contract.md:16,22` |

D-001 正确拒绝「把 reconcile run 扩成通用 Job」，并要求 handler 调用同一 service、把 run 投影为 Job result。独立核验后，这条消费契约仍不能实施：

1. **`Reconcile` / `ReconcileRun` 非幂等**。一次成功调用就落一行新历史。Job 在 `Reconcile` 已提交、`status` 尚未写成 `succeeded` 时被 reclaim 或重试，会得到**第二条**业务 run，而 Job 仍是同一个 id。这直接违反成功标准 2 在消费层的含义。
2. **取消胜出后的副作用**。若条件更新让 Job 成为 `cancelled`，但 handler 已插入 run，则出现「Job 取消、钱包历史已有一次 reconcile」。D-001:22 把 run 当作业务真相源，必须写明这是接受的残余还是必须用「先占 run id / 按 Job id 幂等插入」避免。
3. **投影字段未冻**。result JSON 含哪些 run 字段、失败时 Job `error_code` 与 run `result=inconsistent` 是否区分（链不一致是业务成功还是 Job failed）未写。现网 `result` 只有 `consistent|inconsistent`（`migration.go:59`），二者都是一次成功的对账。

S0 必须冻结至少一条可测规则，例如：handler 以 Job id（或 payload 内一次性 `runID`）作为 reconcile 幂等键；同一 Job 的 reclaim/retry **不得**再插入第二行；测试覆盖「handler 成功但 complete 失败后的第二次 claim」。没有这条，wallet 不能被称为已设计好的真实消费路径。

### F-004 · I-002 所有权模式成立，但未冻成可检查的 Profile/Manifest 不变式

| 字段 | 值 |
|------|-----|
| level | required |
| severity | med |
| status | open |
| 影响门禁 | I-002；成功标准 4；Root「模块矩阵不变」 |
| 关联 | I-002 |
| evidence | `D-001-r4-job-contract.md:15-16,32`；`kernel/profile.go:26-114,159+`；`kernel/provider.go:67-90,107-135`；`wallet/provider.go:162-175,184-186`；`compiled/persistence.go:1-3,23-41`；`composition.go:134-138,349-363`；`A-001-r4-s0-design-self.md:26,37` |

独立核验同意 A-001：把 Job 表放进 PersistenceProviders、runtime 放 internal、HTTP 挂在已启用的 `admin.wallet`，**可以**不把 `core.jobs` 加进 `profileDefaults`，因此默认模块 ID 集与 Manifest fragment 装配路径可以不变。

缺口是 D-001 没有把「如何不改变」写成可检查约束，I-002 因此不能从 `collecting` 改为 `verified`：

1. **`core.jobs` 禁止进入** `profileDefaults`（mvp/admin/demo）、`BuiltinModules()` 以及 `composition` 的 Plan-gated runtime provider 列表。只允许加入 `compiled.PersistenceProviders()`，`Register` 必须空实现（对标 `core.persistence`）。
2. **允许的唯一 catalog 变更**是 `admin.wallet` 的 `Contributions.Routes` 增列 Job 查询/取消/重试/结果路径，且必须**同时**改 `wallet/provider.go` Descriptor 与 `kernel/profile.go` BuiltinModules，否则 `descriptorsMatch` 会 fail-closed。D-001 未列出这些 path key。
3. **禁止**为本波次新增 pages / navigation / fragments / permission keys。权限继续只用 `wallet.read`。
4. 成功标准 4 与 Root 成功标准 3 的「模块矩阵不变」必须解释为：模块 ID 集与 Manifest 装配算法不变；wallet **路由声明**的内容扩展是显式允许的例外，并用 conformance 测试锁默认模块集。

A-001 F-002 把 descriptor 双写列为 recommended 实现门禁。独立意见：双写清单本身属于 **I-002 关闭条件**，应在 S0 写入 D-001；测试证明则仍是 S3/S4。

### F-005 · actor 隔离与 `wallet.read` 全局面的交叉未定义

| 字段 | 值 |
|------|-----|
| level | required |
| severity | med |
| status | open |
| 影响门禁 | I-003；S3 HTTP 契约 |
| 关联 | I-003 |
| evidence | `D-001-r4-job-contract.md:16`；`handler/wallet.go:288-334`；`wallet/store/repository.go:768-784` |

D-001:16：「Job 查询按提交 actor 隔离」。现网 `POST /api/wallet/reconcile` 与 `GET /api/wallet/reconcile/runs` 都只检查 `wallet.read`；list runs **不过滤** `actor_id`（`ORDER BY created_at DESC` 全表分页）。任意持 `wallet.read` 的管理员都可以对空 `accountId` 触发全账户对账（`ReconcileRun` `accountID==""` 分支，`repository.go:649-664`）。

未冻结：

- GET Job / result / cancel / retry 是「仅 `actor_id == 当前用户`」还是「任何 `wallet.read`」？
- 跨 actor 读是 404 还是 403？稳定错误码是什么？
- 业务 runs 列表保持全局、Job 列表按 actor——二者并存时客户端如何对应？
- 系统/内部 actor 是否豁免隔离？

没有这些谓词，S3 HTTP 测试无法写，也无法证明「真实消费」不会把别人的 Job result 泄漏给任一 wallet 管理员。

### F-006 · 结果过期缺少 TTL、触发器与原子清空语义

| 字段 | 值 |
|------|-----|
| level | required |
| severity | med |
| status | open |
| 影响门禁 | I-003；成功标准 1、3 |
| 关联 | I-003 |
| evidence | `D-001-r4-job-contract.md:17,21,26`；`handler/localize.go:17-32`；`errorcatalog` 无 Job/410 码（本轮检索无 `JOB_` / `410` catalog 项） |

D-001 正确把 `expired` 限于「已结束且到达保留期限」，避免把「结果不可取」写成执行失败。但操作语义仍不足：

1. `expires_at` 何时写入（`finished_at + TTL`）？TTL 默认值？谁扫描（周期 runner vs 读时惰性）？
2. 转 `expired` 与清空 `result` 必须同一条件更新，否则读路径可能看到 `expired` 仍带 result，或 `succeeded` 已无 result。
3. 410 必须走 R1 包络（`writeLocalizedError` + catalog 码）。D-001 只写「稳定 410」，没有 code。
4. `cancelled`/`failed` 无结果时是否仍转 `expired`？`expired` 可否 retry？

这些不补，成功标准 1/3 的过期测试无法确定期望。

### F-007 · retry 的字段复位表与 `max_attempts` 默认值未写（recommended）

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | med |
| status | open |
| 影响门禁 | S1 repository；不单独阻断 I-003（被 F-001 覆盖） |
| 关联 | I-003 |
| evidence | `D-001-r4-job-contract.md:20,26` |

决策 6 已排除「回退 attempt / 换 Job id」，方向正确。建议在转换表中显式列出 retry 时清除/保留的字段：`error_*`、`result`、`lease_*`、`cancel_requested`、`progress`、`finished_at`、`expires_at`；并冻结 `attempt` 初值（建议 0）与 `max_attempts` 默认值。否则 S1 会在「retry 后下一次 claim 把 attempt 变成 2 还是保持 1」上分叉。

### F-008 · `POST /api/wallet/reconcile` 202 与现网 200/审计时点未登记为 S3 必测（recommended）

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | low |
| status | open |
| 影响门禁 | S3 HTTP/审计；不阻断 S0 冻结 |
| evidence | `handler/wallet.go:304-306`；`wallet_test.go:222-234,254` 断言 200 与 `wallet.reconcile` 事件 |

D-001:22 已明确改为 202 + Job representation。现网测试与 `wallet.reconcile` 审计写在同步成功之后。S3 须同时改断言，并冻审计时点：提交时写、完成时写、还是两者都写但 event 不同。本条不否定 202 决策。

## 必改项汇总

1. **F-001**：把 I-003 补成可冻结 transition table（from / event / guard / to / 字段突变 / 非法转换与 R1 错误码），并写明 complete/fail 与 cancel 的唯一 `WHERE`。
2. **F-002**：冻 claim vs reclaim 的 attempt 规则、lease 续约、以及「重复领取不重复执行」的 in-flight 互斥。
3. **F-003**：冻 wallet handler 幂等：同一 Job 至多一条 `wallet_reconciliation_runs`；写明取消胜出但 run 已落库时的语义；冻 result 投影与 `inconsistent` ≠ Job `failed`。
4. **F-004**：冻 I-002 不变式：`core.jobs` 仅 PersistenceProviders；不进 Profile/BuiltinModules/runtime provider 表；列出将写入 `admin.wallet` 的路由键并要求 Descriptor 双写；无新 page/nav/fragment/permission。
5. **F-005**：冻 Job GET/result/cancel/retry 的 actor 谓词、跨 actor 错误码，以及与全局 reconcile runs 列表的关系。
6. **F-006**：冻 `expires_at` TTL、触发器、expire+clear 原子性、410 catalog 码、哪些终态会过期。

## 与既有意见的异同

| 点 | A-001 (self) | 本条 (independent) |
|----|----------------|--------------------|
| E-001 扫描主干 | 同意 | 同意；独立复核通过 |
| CollectPersistence ≠ RegisterContributions | 作为 pass 主依据 | 同意模式可行；**不足**以关闭 I-002（F-004） |
| 条件 UPDATE 可做原子 claim | 同意 | 同意能力存在；**不足**以保证不双执行（F-002） |
| 不复用 `task_runs` / 不把 run 当 Job | 同意 | 同意 |
| D-001 状态集合覆盖 VP-012 | 同意名词覆盖 | 同意名词覆盖；**不同意**已构成 I-003 精确表（F-001） |
| 开放 required | 0；pass；仍等 independent 才关 I-002/I-003 | **6** 条 required；verdict = **conditional** |
| A-001 F-001 startup/周期扫描 | recommended / S2 | 同意必须有扫描；升级为 S0 必须先定义互斥（并入 F-002） |
| A-001 F-002 descriptor + conformance | recommended / S3/S4 | 测试仍在 S3/S4；**清单本身**是 I-002 关闭条件（F-004） |

无意见冲突需要 P-004 用户裁：self 并未关闭 I-002/I-003，也未放行 S1。本条收紧的是「现在能否冻结」，不是推翻 E-001 事实。

## 结论 + 建议给编排器/用户的下一步

**verdict = conditional。** E-001 现状事实成立；migration-only `core.jobs` **可以**不进入默认 Profile/Manifest 模块集。D-001 作为 S0 冻结稿仍不完整：I-002/I-003 必须保持 `collecting`，**不得进入 S1**。

建议 `/govern`：响应本条 A-002 F-001～F-006；用修订 D-001（或 D-002）写入转换表、lease/claim 互斥、wallet 幂等、Profile 不变式、actor 谓词与过期语义；再请独立复审关闭这些 required 之后，才把 I-002/I-003 标 `verified` 并将 D-001 改为 `accepted`。

## 声明

本意见不修改 `status` / `progress` / 检查点 / `00-meta` / `D-001` / 实现代码，也不关闭任何 finding。响应由 `/govern` 处理。
