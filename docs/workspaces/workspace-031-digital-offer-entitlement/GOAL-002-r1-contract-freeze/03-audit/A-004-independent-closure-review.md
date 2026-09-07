---
doc_type: goal-audit
id: A-004-independent-closure-review
parent: GOAL-002-r1-contract-freeze
date: 2026-09-05
source: independent
auditor: codex (gpt-5.6-sol, medium)
audit_type: finding-closure
scope: 仅复审 A-003 对 A-002 F-001～F-005 的关闭证据；核对 D-002 现行合同与 workspace.md R1 阶段投影
verdict: fail
open_required: 2
version: 0.1.0
---

# A-004 · A-002 finding-closure 独立复审（2026-09-05）

- **source**：independent
- **auditor**：codex (gpt-5.6-sol, medium)
- **类型 / scope**：finding-closure · 仅复审 A-003 的关闭证据是否真实、充分、可核对地闭合 A-002 F-001～F-005；不重审整个目标，不把 R2/R3 尚未发生的实现或测试当作本次合同闭合分母。
- **verdict**：**fail**
- **open required**：2（A-002 F-001、F-002 仍未充分闭合）

## 范围与区间

- 工作区绑定有效：`workspace-031-digital-offer-entitlement` 的 `root_goal`、`canonical_scope`、`plan_refs` 与 `primary_plan` 均覆盖当前目标；共享资料目录为 `none`，本次未使用外部共享材料。[`workspace.md:2-14`；`GOAL-002-r1-contract-freeze/00-meta.md:2-12`]
- 核对区间：A-002 原 finding 与修复方向、A-003 关闭证据表、D-002 于提交 `a81c35f5` 后的现行文本、`workspace.md` R1 行，以及 wallet / store / operationlog 的现行接口事实。[`03-audit/A-002-independent-contract-audit.md:49-106`；`03-audit/A-003-self-response-a002.md:22-34`；`git show a81c35f5 --stat`]
- 本复审判断的是“合同修订能否按现有接口无歧义照做”，不是要求数字 Offer 实现已经存在；R2/R3 实施与验收测试仍是后续阶段证据。[`01-decision/D-002-digital-offer-contract.md:289-293`]

## 关闭证据核对表

| Finding | A-003 闭合声明 | 核对结论 + 证据 |
|---------|----------------|-----------------|
| A-002 F-001 · 购买 mutation 与并发幂等协议 | `fixed` | **未充分闭合**。§4.4 已真实补齐 `MutateInTx(tx, accountID, in, entryID, now)`、正数 `AmountDelta` 的 apply 语义、两笔均传 `account.ID`、互异预生成 id、完整 `LedgerEntryInput`、互异幂等键、purchase UNIQUE 回读、三次整事务重试及 Telegram `request_id` 派生；这些与 wallet 签名/apply 表一致。[D-002 §4.4，`:122-152`；`apps/api/modules/wallet/store/repository.go:155-165`、`:529-573`、`:576-648`] 但合同要求调用方识别并重试 wallet `errIdempotencyRace`，该 sentinel 是 wallet/store 包内未导出符号，`MutateInTx` 会直接返回它，而只有同包的 `Mutate` 包装器能识别并回读；拟建的 `digitaloffer` 包无法按合同分支可靠识别该错误。因此 ledger 唯一竞争的“可执行冻结”仍不成立。[D-002 §4.4，`:141-149`；`repository.go:48-58`、`:637-648`、`:664-681`] |
| A-002 F-002 · Consume 跨方言并发与 void 线性化 | `fixed` | **未充分闭合**。§5.2 已补齐候选谓词与稳定排序、最终 UPDATE 的主体/offer/形态/active/余额谓词、`RowsAffected` 竞争、部分扣减回滚以及 void 谓词重检，算法方向可证明且测试分母明确。[D-002 §5.2，`:174-204`] 但伪代码把 `for attempt in 1..3` 放在单个 `Run(ctx, tx)` 内，并在循环中直接 `COMMIT` / `ROLLBACK`；现行 `kernel.Tx` 不暴露 Commit/Rollback，且一个 `Store.Run` 固定只代表一个事务、嵌套 Run 被禁止。回滚后也不能在同一 tx 上进入下一 attempt。故“下一 attempt 重读候选”的事务边界不可照做；重试循环必须位于 `Store.Run` 外，每次 attempt 新开事务，并以 callback 返回值驱动 commit/rollback。[D-002 §5.2，`:179-201`；`apps/api/kernel/store.go:27-47`；`apps/api/internal/store/store.go:145-170`；`apps/api/internal/store/postgres.go:231-259`] 此外，`BeginTx(ctx, nil)` 只表明未显式固定 TxOptions，代码本身未强制 PostgreSQL `READ COMMITTED`；合同不应把可配置的数据库默认值写成由平台代码保证的隔离事实。[D-002 §5.2，`:179`；`postgres.go:234-240`] |
| A-002 F-003 · Admin 同事务 fail-closed 审计 | `fixed` | **闭合成立**。§1 已将依赖改为 `operationlog.TransactionalRecorder`；§7 冻结 caller-owned 同事务、审计失败整体回滚、四个事件名、域行 record id、会话 actor、before/after detail、重复 void 不追加事件、购买路径豁免理由与失败注入测试。真实接口为 `RecordOperationTx(kernel.Tx, Operation) error`，并使用传入 tx 写入 operation log，错误向调用方返回。[D-002 §1，`:29-35`；§7，`:235-247`；`apps/api/modules/operationlog/repository.go:112-121`、`:149-196`] |
| A-002 F-004 · 聚合优先级与错误码映射 | `fixed` | **闭合成立**。§5.1 已冻结 `valid → no_entitlement → expired → exhausted → voided` 的确定性聚合顺序；§9 新增 `BIZOFFER_ENTITLEMENT_INSUFFICIENT` / `ErrEntitlementInsufficient` / HTTP 409，并冻结 Check reason、Consume sentinel、wallet 余额不足及 HTTP/Telegram 一致映射与表驱动测试分母。[D-002 §5.1，`:166-172`；§9，`:265-280`] |
| A-002 F-005 · workspace 阶段投影 | `fixed` | **原 finding 已闭合**。R1 行已从“待开始；I-031-001～005 尚未处理”改为“进行中”，并记录 C1 于 2026-09-05 关门、I-031-001～005 verified 与合同审计循环；该变更包含在 `a81c35f5`。[`workspace.md:31-38`；`git diff a81c35f5^ a81c35f5 -- docs/workspaces/workspace-031-digital-offer-entitlement/workspace.md`] 当前措辞仍停在 “A-001/A-002 审计循环中”，未投影 A-003 已响应、等待 A-004 的最新事实，见新 recommended finding。 |

## 新 Findings

### F-001 · wallet ledger 唯一竞争 sentinel 对域调用方不可见

- **严重度**：med
- **建议**：required
- **状态**：open
- **关联**：A-002 F-001 closure
- **证据**：D-002 §4.4 `:141-149`；`apps/api/modules/wallet/store/repository.go:48-58`、`:637-648`、`:664-681`
- **问题**：合同要求 `digitaloffer` 对 wallet `errIdempotencyRace` 做整事务重试，但该错误未导出，跨包调用方无法可靠分类；真实可识别/回读逻辑只存在于不能用于 caller-owned 购买事务的 `Mutate` 包装器中。
- **关闭方向**：在进入 R2 前冻结一个跨包可执行的错误契约，例如由 wallet 暴露稳定 retryable sentinel/分类接口，或让 `MutateInTx` 在不破坏 caller-owned 事务的前提下返回可识别结果；随后同步 D-002 并以针对该竞争路径的测试作为后续实施证据。不得依赖错误字符串。

### F-002 · Consume 重试循环与 `Store.Run` 事务边界矛盾

- **严重度**：med
- **建议**：required
- **状态**：open
- **关联**：A-002 F-002 closure
- **证据**：D-002 §5.2 `:179-201`；`apps/api/kernel/store.go:27-47`；`apps/api/internal/store/store.go:145-170`；`apps/api/internal/store/postgres.go:231-259`
- **问题**：合同伪代码在一个 `Run(ctx, tx)` 内执行多次 attempt 并直接 COMMIT/ROLLBACK，但 `kernel.Tx` 没有这些方法，一个 Run 也只能提交或回滚一次。因此 RowsAffected 竞争后的“下一 attempt”没有合法的新事务边界，算法不能按现文实现。
- **关闭方向**：把有界 attempt 循环冻结在 `store.Run` 外；每次 attempt 调用一次 `store.Run`，callback 内只返回成功、确定性不足或 retryable-race，外层据此停止或重试。同步说明 PostgreSQL 隔离是部署前提还是显式 TxOptions 保证，不把 `BeginTx(ctx, nil)` 写成强制隔离保证。

### F-003 · workspace R1 审计投影未包含 A-003/A-004 状态

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **证据**：`workspace.md:35`；`03-audit/A-003-self-response-a002.md:14-16`、`:32-34`
- **问题**：原 F-005 的“R1 仍待开始”已修复，但当前 R1 行仍写“A-001/A-002 审计循环中”，未反映 A-003 已完成 fixed 响应并等待 A-004 closure 复审。该偏差不改变本次 required verdict，但会继续形成阶段投影滞后。
- **关闭方向**：由 `/govern` 在响应本意见时同步工作区 R1 行为当前真实审计状态。

## 结论 + 建议下一步

A-003 对 F-003、F-004 与原 F-005 的关闭证据成立；F-001、F-002 的合同修订也覆盖了大部分原要求，但分别与现行 wallet 错误可见性、`Store.Run` 单事务边界直接矛盾，尚不能按现文无歧义实现。A-003 “全部 5 项均 fixed” 的关闭声明因此不实，本次 **verdict = fail**，仍有 **2 项 med required**。

建议由 `/govern`：

1. 响应 A-004 F-001/F-002，优先走 `fixed`，修订 D-002 的跨包 retryable error 契约与 Consume 外层重试事务结构；不得以未来实现测试替代当前合同缺口。
2. 同步处理 recommended F-003 的 workspace 阶段投影。
3. 修订后追加 self 响应，并再做 focused independent finding-closure；required 未合法闭合前不得把 D-002 转 accepted、关门 C3 或推进受其阻断的 R2/R3。

## 声明

本意见仅追加 independent audit ledger；未修改 `00-meta` 的 `status`/`progress`、goal-tree 状态列、D-001/D-002 正文或 01/02 台账。finding 响应与状态推进由 `/govern` 处理。
