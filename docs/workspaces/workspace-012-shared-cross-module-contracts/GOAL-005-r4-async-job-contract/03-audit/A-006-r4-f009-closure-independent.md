---
id: A-006-r4-f009-closure-independent
goal: GOAL-005-r4-async-job-contract
doc: audit-entry
record_id: A-006
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: A-004 F-009 关闭复核；D-002 v0.2.0 recover-cancel + CompleteWithCommit；I-003；S0→S1；无新可达死状态 / wallet run↔Job 分裂
audit_type: finding-closure
verdict: pass
status: recorded
parent: GOAL-005-r4-async-job-contract
created: 2026-08-18
updated: 2026-08-18
version: 0.1.0
responds_to: A-004
reviews: A-005
---

# A-006 · A-004 F-009 关闭复核（2026-08-18）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high；项目级路径见 `docs/architecture/independent-audit-execution.md`）
- **类型**：ad-hoc / finding-closure
- **scope**：GOAL-005 S0；仅复核 A-004 required F-009 的关闭候选。只读 A-004、A-005、E-003、D-002 v0.2.0，以及必要的现有 store / wallet 证据。判断 F-009 是否可 `fixed`、I-003 是否可 `verified`、D-002/S0 是否可放行 S1；核对未引入新的可达死状态或 wallet run / Job 状态分裂。
- **verdict**：pass

## 范围与区间

- **工作区**：`workspace-012-shared-cross-module-contracts`（`workspace.md`：`root_goal` = `GOAL-001-shared-cross-module-contracts`；`canonical_scope` 与本目标路径一致；`shared_materials_catalog: none`；`vision_role: delivery`；`primary_plan` = `VP-012-shared-cross-module-contracts`）
- **covered**：A-004 F-009 原文三项关闭条件；A-005 / E-003 候选响应；D-002 v0.2.0 §3～§5 / §9 的 `recover-cancel` 与 `complete-with-commit`；对照现有 `WithTx` / 单连接、`wallet_reconciliation_runs` PK、`ReconcileRun` 自开事务、条件 `UPDATE` + `RowsAffected`
- **excluded**：S1～S4 实施与关门、未重新执行测试套件、其他工作区上下文、共享资料内容（目录为 `none`）、改写 D-002 / `00-meta` / status / progress / 实现代码、重开已由 A-004 标 `fixed` 的 A-002 F-001～F-008
- **本轮未复验**：未运行 `go test`；「可实施」只核到现网事务/表结构与 D-002 谓词是否闭合，不把未写的 `CompleteWithCommit` / `ReconcileOnceTx` 当成已交付

## 工作区与对齐（只读）

| 检查项 | 结论 | 证据 |
|--------|------|------|
| 工作区绑定 | 通过 | `workspace.md` Root / canonical / `plan_refs`+`primary_plan` 与 GOAL-005 `parent`、`primary_plan` 一致 |
| 共享资料引用 | 无引用，不构成关闭证据 | `shared_materials_catalog: none` |
| 既有 Goal 审计 | A-001 self = pass；A-002 independent = conditional（F-001～F-006 已由 A-004 标 fixed）；A-004 independent = conditional（唯一开放 required = F-009）；A-005 self = conditional（F-009 候选 fixed，不自闭） | `03-audit.md`；`A-005-r4-a004-response-self.md:19-24` |
| P-004 冲突 | 无 | A-005 未宣称已闭合；本条与 A-004 无相反 verdict 需用户裁 |

## 成果（有证据）

| 主张 | 证据 | 核验 |
|------|------|------|
| D-002 已升至 v0.2.0，并声明 F-009 须经本复核才转 `accepted` | `01-decision/D-002-r4-precise-contract.md:1-14,106-114` | 通过：`status` 仍为 `proposed`；`version: 0.2.0` |
| A-005 / E-003 未自闭 F-009 | `03-audit/A-005-r4-a004-response-self.md:22-24`；`02-execution/E-003-r4-f009-response.md:13` | 通过：候选 fixed / 等待 independent |
| `recover-cancel` 已写入转换表，且不依赖旧 fencing owner | `D-002-r4-precise-contract.md:54,68` | 通过：`from+guard` = `running, cancel_requested=1, lease_expires_at <= now` → `cancelled` |
| `complete-with-commit` 把 consumer durable write 与 Job `succeeded` 放进同一事务 | `D-002-r4-precise-contract.md:51,58,65,72-73,77` | 通过：任一步失败整体回滚；cancel/recover-cancel 先赢则 callback 不执行 |
| 现网同一 SQLite + `WithTx` 可承载该原子提交 | `apps/api/internal/store/store.go:49,67-80`；`wallet/store/repository.go:21-34,642-703` | 通过：wallet 已走平台 `TxRunner`；`SetMaxOpenConns(1)` 使外层事务持锁期间 scanner 无法交错提交 |
| 现网 `ReconcileRun` **仍自开** `WithTx`，故契约必须用 `ReconcileOnceTx(tx)` 而非直接调 `ReconcileRun` | `wallet/store/repository.go:640-703` | 通过：D-002:72-73 已点名 `ReconcileOnceTx(tx, …)`；嵌套 `WithTx` 会破坏「同事务回滚」 |
| run 表 PK + UNIQUE 冲突检测仍在；无需改 wallet schema | `wallet/migration/migration.go:56-64`；`repository.go:695-700,806-808` | 通过：与 A-004 F-003 关闭条件一致 |
| 现网 HTTP 仍每次 `newID` + 同步 200 | `wallet/provider.go:128-134`；`handler/wallet.go:287-307` | 通过：属 S3 实施缺口，不否定本条对 S0 谓词的闭合 |

## 对照成功标准（F-009 / S0 关闭复审）

| 标准 | 本轮 | 证据 |
|------|------|------|
| 1. 六态与合法/非法转换可测；lease 过期后无活 fencing 时不得留在 `running` | F-009 原死状态已有事件 | `D-002:54,68`；下表「可达 running 组合」 |
| 2. 领取 / 不重复执行；已提交 wallet run 不得被写成 `cancelled` | 用同一事务消除「run 已提交、Job 未 complete」窗口 | `D-002:51,58,65,72-73,77` |
| 3. 202 + 轮询 + 410；已提交 run 的 complete 不再被死 `running` 挡住 | 窗口消除后 GET result 不会对已提交 run 永久返回 `JOB_RESULT_NOT_READY` | `D-002:77,87` |
| 4. Profile / 模块矩阵 / Manifest 不变 | 本轮未改 §1；不重开 F-004 | `D-002:16-26` |

## 信息门禁核对（P-005）

| ID | 级别 | 最晚阶段 | 当前（00-meta，本条不改） | 本轮结论 |
|----|------|----------|---------------------------|----------|
| I-001 | required | S0 结束前 | verified | 保持；本轮不重开关 |
| I-002 | required | S0 结束前（影响 S1） | verified | 保持；A-004 已确认 |
| I-003 | required | S0 结束前（影响 S1/S2） | collecting | **可转 `verified`**。F-009 原文三项关闭条件已成文且可重复核对；精确转换在「无可达无事件 running」意义上已闭合 |
| I-004 | required | S1 实施前 | verified | 保持；不能代替 S4 关门独立审 |

无 `deferred`。无用户书面 `accepted-residual`。I-003 的 meta 状态变更由 `/govern` 写入，本条不改 `00-meta`。

## F-009 原文三项关闭条件（对照 A-004:272-280）

| # | A-004 要求 | D-002 v0.2.0 响应 | 本条 |
|---|------------|-------------------|------|
| 1 | lease 过期且**没有**活 fencing owner 时，Job 不能留在 `running` | `recover-cancel`：`running ∧ cancel_requested=1 ∧ lease_expires_at <= now` → `cancelled`，不依赖旧 owner；§4 写明 scanner 执行且不受 attempt 是否耗尽影响 | **满足**。`D-002:54,68` |
| 2 | 若该 Job ID 已有 `wallet_reconciliation_runs` 行 → 不得写成 `cancelled`，应能 `complete` 或等价 adopt-success | 不用「先查 run 再分支」；用 `CompleteWithCommit` 使「已提交 run ∧ Job 仍为 running」成为不可达。先赢 cancel/recover-cancel ⇒ callback 不执行、无 run；先赢 complete-with-commit ⇒ Job 已是 `succeeded`，`recover-cancel` 的 `from=running` 不再匹配 | **满足**（等价构造，不是 A-004 示例的 reclaim-then-adopt）。`D-002:51,58,72-73,77` |
| 3 | 若尚无 run → 可原子转 `cancelled` | `recover-cancel` 单条条件 UPDATE，不依赖旧 owner | **满足**。`D-002:39,54` |

A-004 给出的示例（`reclaim` 在 `cancel_requested=1` 时仍换 fencing，或 scanner 按 run 是否存在分支）**不是**唯一合法写法。v0.2.0 选择「禁止在 cancel 旗标下 reclaim + 原子 consumer commit + 无 owner 的 recover-cancel」，与三项条件等价，且不再要求恢复路径去读 wallet 表。

## 可达 `running` 组合（查新死状态）

| cancel | lease | attempt | 可发射事件 | 是否死状态 |
|--------|-------|---------|------------|------------|
| 0 | 未过期 | < max | heartbeat / progress / complete-with-commit / fail / request-cancel | 否 |
| 0 | 未过期 | ≥ max | 同上（claim 已加过 attempt；仍可 complete） | 否 |
| 1 | 未过期 | 任意 | heartbeat / progress / complete-with-commit / finalize-cancel | 否 |
| 0 | 过期 | < max | reclaim | 否 |
| 0 | 过期 | ≥ max | exhaust | 否 |
| 1 | 过期 | < max | **recover-cancel**（A-004 原文洞） | **否（已补）** |
| 1 | 过期 | ≥ max | recover-cancel；表上的 exhaust 谓词也匹配 | **否**（有事件；见 F-010 谓词重叠，不是死状态） |

崩溃且 lease 尚未过期时仍会短暂停在 `running`，上界为 lease=30s，与 A-004 已接受的 F-002 设计相同，不是新死状态。

## wallet run / Job 配对（查分裂）

在「成功路径必须走 `CompleteWithCommit`，且 `ReconcileOnceTx` 使用**同一** `*sql.Tx`」前提下：

| Job 状态 | 已提交 `wallet_reconciliation_runs` | 如何到达 |
|----------|-------------------------------------|----------|
| queued / running | 无 | submit/claim 不写 run；未提交事务崩溃则一起回滚 |
| succeeded | 有（同一事务） | complete-with-commit 提交 |
| failed | 无 | `fail` 不写 run；callback 失败则整事务回滚 |
| cancelled | 无 | request-cancel queued / finalize-cancel / recover-cancel 先赢 ⇒ callback 不执行 |
| expired | 有（业务史保留；Job `result` 被清空） | 仅 `succeeded→expired`；不是 running/cancelled 分裂 |

「run 已提交且 Job 仍为 `running`」以及「run 已提交且 Job 为 `cancelled`」在契约下不可达。这正是 A-004:268 指出的 F-003 违反窗口。

现网 `ReconcileRun` 自己 `BeginTx`（`repository.go:642`）**不能**直接塞进 callback；若 S3 误调它，内层事务会先提交 run，外层 Job 更新失败时就会重新打开 F-009。该约束已写在 D-002:72-73，属 S3 必测，不回退本条关闭。

## Disposition（保留 A-004 F-009 原文）

| ID | 原文标题（A-004） | level | 候选（A-005） | 本条 disposition | 关闭路径 |
|----|-------------------|-------|---------------|------------------|----------|
| F-009 | `running + cancel_requested=1 + lease expired + attempt < max_attempts` 无恢复事件 | required / high | candidate fixed | **fixed** | D-002 v0.2.0 §3 `recover-cancel` + §3/§4/§5 `CompleteWithCommit` |

`fixed` 仅表示 A-004 F-009 原文三项关闭条件已成文且可重复核对，**不是** S1～S3 已实施。

## Findings

### F-009 · `running + cancel_requested=1 + lease expired + attempt < max_attempts` 无恢复事件

| 字段 | 值 |
|------|-----|
| level | required |
| severity | high |
| status | **closed** |
| disposition | **fixed** |
| 影响门禁 | I-003；成功标准 1–3；S0 冻结；S1 实施 |
| 关联 | I-003；约束 F-003 第 2 项在恢复路径上可执行 |
| evidence | `01-decision/D-002-r4-precise-contract.md:51,54,58,65-68,72-73,77`；`apps/api/internal/store/store.go:49,67-80`；`apps/api/internal/modules/wallet/store/repository.go:21-34,640-703,806-808`；`apps/api/internal/modules/wallet/migration/migration.go:56-64` |
| closure | 本条 A-006；响应记录见 A-005 / E-003 / D-002 §9 |

**原文（A-004:240-280，保留）**：

> D-002 表对 `running` 的 scanner / 恢复事件只有：
>
> | event | guard | to |
> |-------|-------|----|
> | reclaim | lease expired **且 `cancel_requested=0`** 且 attempt < max | running（新 owner，attempt+1） |
> | exhaust | lease expired 且 attempt >= max | failed |
> | finalize-cancel | **exact owner+version** 且 cancel=1 | cancelled |
> | complete | **exact owner+version** | succeeded |
>
> 下列状态**可达**且无任何事件：
>
> 1. `submit` → `claim`（attempt=1）→ `request-cancel running`（`cancel_requested=1`，status 仍为 running）
> 2. handler / 进程在 `complete` 或 `finalize-cancel` 之前崩溃；进程重启后 active set 为空
> 3. 租约过期后：`reclaim` 被 `cancel_requested=1` 挡住；`exhaust` 因 `attempt < max` 不成立；`finalize-cancel` / `complete` 需要已消失的 fencing token
>
> 结果：Job **永久停在 `running`**。跨重启也不会前进，因为 reclaim 被禁止，attempt 不再增加，exhaust 永远不会触发。
>
> 若崩溃发生在 `ReconcileOnce` 已提交之后，还违反 D-002 自己写的 F-003 规则：业务 run 已在 `wallet_reconciliation_runs`（`migration.go:56-64`；`repository.go:695-700`），Job 却不能 `complete`，GET result 会按 `JOB_RESULT_NOT_READY` 对待仍为 running 的行（`D-002:87`）。
>
> 这不是 A-002 原文三项的重复，而是 D-002 在补表时引入的**表闭合缺口**。S1 若按现表开工，实现者必须自行发明恢复事件；不同选择会改写成功标准 1–3 与 F-003「不伪称 cancelled」。
>
> 关闭本条至少要在 D-002 增加一条可测恢复事件，且必须同时满足：
>
> - lease 过期后，**没有**活 fencing owner 时，Job 不能留在 `running`
> - 若该 Job ID 已有 `wallet_reconciliation_runs` 行 → 不得写成 `cancelled`（服从 F-003 / `D-002:73`），应能 `complete` 或等价 adopt-success
> - 若尚无 run → 可原子转 `cancelled`

**复核（v0.2.0）**：

1. **恢复事件已存在且覆盖原文轨迹。** 原文三步到达 `running ∧ cancel_requested=1 ∧ lease expired ∧ attempt < max` 后，scanner 可发 `recover-cancel`（`D-002:54`），不再需要已消失的 owner+version。`reclaim` 仍要求 `cancel_requested=0`（`D-002:46`），与「取消后不换新 handler」一致。
2. **「已提交 run 不得伪称 cancelled」改为构造上不可达。** `complete-with-commit`（`D-002:51,58`）先锁定/复核 fencing，再跑 callback，再写 `succeeded`；任一步失败整体回滚。因此：
   - 进程在 callback 提交前崩溃 → 无 run，随后 `recover-cancel` → `cancelled`（条件 3）；
   - callback 与 Job success 已提交 → 行不再是 `running`，`recover-cancel` 不匹配，Job 保持 `succeeded`（条件 2）；
   - `recover-cancel` 先提交 → fencing 失败，callback 不执行，无 run（条件 2/3）。
3. 现网模式支持该构造：wallet 与未来 `core.jobs` 共用平台 `WithTx`（`store.go:67-80`；`repository.go:21-34`）；`wallet_reconciliation_runs.id` 已是 `TEXT PRIMARY KEY`（`migration.go:56-57`）。现网 `ReconcileRun` 仍自开事务（`repository.go:642-703`）是 S3 必须改成 `ReconcileOnceTx(tx)` 的事实，不是 S0 未冻。

原缺口已闭合。表上 `exhaust` 与 `recover-cancel` 在 `cancel=1 ∧ attempt≥max ∧ lease expired` 重叠，不是死状态，也不打开 run/Job 分裂；记 recommended F-010，不回退本条。

### F-010 · `exhaust` 与 `recover-cancel` 在表上的 guard 重叠

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | low |
| status | open |
| 影响门禁 | 不阻断 I-003 / S1 |
| evidence | `D-002-r4-precise-contract.md:47,54,67-68` |

当 `running ∧ cancel_requested=1 ∧ lease_expires_at <= now ∧ attempt >= max_attempts` 时，§3 表同时匹配：

- `exhaust`（`D-002:47`：只写 `lease expired, attempt >= max`）→ `failed` / `JOB_ATTEMPTS_EXHAUSTED`
- `recover-cancel`（`D-002:54`）→ `cancelled`

§4 散文已裁定「`recover-cancel` 不受 attempt 是否耗尽影响」（`D-002:68`），意图唯一。S1 按 §3+§4 实现时应把 `exhaust` 写成 `cancel_requested=0`，不必再发明事件。这不是可达死状态，也不是 wallet run / Job 分裂；不阻止 F-009 `fixed`、I-003 `verified` 或放行 S1。编排器若要消掉表/散文差，可在进入 S1 前给 `exhaust` 补 `cancel_requested=0`（本意见不改 D-002）。

## 必改项汇总

无开放 required。F-009 可按 `fixed` 关闭。

F-010 为 recommended，不阻断 S1。

## 与既有意见的异同

| 点 | A-004 | A-005 | 本条 |
|----|-------|-------|------|
| F-009 | **open** required；阻断 I-003 / S1 | 候选 fixed，不自闭 | **同意关闭为 fixed** |
| I-003 | 不可 verified | 保持 collecting | **可 verified** |
| D-002 / S1 | 不可放行 | 保持 proposed / 不进 S1 | **可 `accepted` 并放行 S1** |
| 新死状态 / run↔Job 分裂 | 发现 F-009 死状态 + 已提交 run 无法 complete | 以原子 commit 消除窗口 | 未发现新的可达无事件 `running`；契约下无「run 已提交 ∧ Job 为 running/cancelled」 |
| 冲突 | — | 未宣称已闭合 | 无 P-004 冲突 |

A-001 两条 recommended（S2 扫描路径、S3/S4 descriptor conformance）与 A-004 已关闭的 F-001～F-008 不被本条改写。

## 结论 + 建议给编排器/用户的下一步

**verdict = pass。** A-004 F-009 可按 `fixed` 关闭。I-003 **可**转 `verified`。D-002 **可**转 `accepted`。S0 **可**放行 S1。

未发现新的可达死状态。在 `CompleteWithCommit` + `ReconcileOnceTx(tx)` 前提下，未发现 wallet run 与 Job 状态分裂。

建议 `/govern`：响应本条 A-006；将 A-004 F-009 标 `fixed`；把 I-003 标 `verified`；将 D-002 转 `accepted`；更新 S0 检查点并进入 S1。F-010 可选：给 `exhaust` 补 `cancel_requested=0`。S3 必测：callback 不得调用自开事务的 `ReconcileRun`；崩溃回滚后不得留下孤立 run。

## 声明

本意见不修改 `status` / `progress` / 检查点 / `00-meta` / `D-002` / 实现代码，也不自行关闭目标状态。响应由 `/govern` 处理。
