---
doc_type: goal-audit
id: A-002-independent-contract-audit
parent: GOAL-002-r1-contract-freeze
date: 2026-09-05
source: independent
auditor: codex (gpt-5.6-sol, medium)
audit_type: design-plan
scope: D-002 数字 Offer 业务域合同 v0.1.0（draft）的设计审计——对照 VP-031 判据、D-001 用户裁决框架、V-F119 与代码先例可达性
verdict: conditional
open_required: 3
version: 0.1.0
---

# A-002 · 数字 Offer 业务域合同独立审计（2026-09-05）

- **source**：independent
- **auditor**：codex (gpt-5.6-sol, medium)
- **类型 / scope**：design-plan · D-002 数字 Offer 业务域合同 v0.1.0（draft）；对照 VP-031 判据 1/2/3/5/7、首波冻结表、D-001 用户裁决、V-F119 与当前代码先例，判断合同是否可照做。
- **verdict**：**conditional**
- **open required**：3（F-001～F-003）
- **完整意见**：本文件（无附件）

## 范围与区间

- 工作区绑定核对：`workspace-031-digital-offer-entitlement` 的 `root_goal`、`canonical_scope`、`plan_refs` 与 `primary_plan` 均指向当前 Root / VP-031，目标 `GOAL-002-r1-contract-freeze` 位于 canonical 范围内。[`workspace.md:2-14`；`goal-tree.md:8-18`；`GOAL-002-r1-contract-freeze/00-meta.md:2-12`]
- 审计区间：D-001、D-002、目标 01/02/03 台账、A-001 self、VP-031，以及用户指定的钱包事务、限流、Telegram、provider 与 composition 代码先例。未审 R2/R3 尚未发生的实现，不以计划文本冒充实施证据。
- P-005：I-031-001～003 已由用户书面裁决为 verified，I-031-004～005 已按 non-blocking 默认项冻结；本轮未发现到期而未处理的 required 信息项。[`01-decision/D-001-info-adjudication.md:14-27`；`VP-031-digital-offer-entitlement.md:97-105`]

## 成果（有证据）

- **D-001 一致性成立**：D-002 已落实“时长/次数二者并存且一 Offer 固定一种”“同步一拍 fulfilled、无 pending、失败整体回滚”“Admin 只读 + 作废、无人工发放”，并冻结 `price` / `buy` / `entitlements` 与模块 id `biz.digital-offer`。[`01-decision/D-001-info-adjudication.md:16-27`；`01-decision/D-002-digital-offer-contract.md:31-35`、`:46-55`、`:65-76`、`:92-119`、`:155-189`]
- **A-001 五项修复均真实落入 D-002**：Consume 已增加 `offerID`；购买事务已加入 `GetOrCreateSubjectAccountInTx` 与币种拒绝；freeze/deduct 已要求 `ref_type/ref_id` 反链；Offer 下架/改价不影响既有权益；同 request_id 不同 offer 返回冲突。[`03-audit/A-001-self-contract-self-review.md:43-66`、`:77-87`；`01-decision/D-002-digital-offer-contract.md:74-75`、`:105-112`、`:144-148`]
- **钱包主原语真实存在且总体可复用**：`kernel.Store.Run` / `kernel.Tx`、`GetOrCreateSubjectAccountInTx`、`MutateInTx`、`EntryFreeze`、`EntryDeductFrozen`、`RefType` / `RefID` / `IdempotencyKey`、`ErrInsufficient`、`ErrVersionConflict` 均存在；voucher Redeem 已提供“同一事务内 get-or-create + MutateInTx”的直接先例。[`apps/api/kernel/store.go:27-47`；`apps/api/modules/wallet/store/repository.go:48-66`、`:155-165`、`:421-459`、`:576-648`；`apps/api/modules/wallet/voucher/service.go:174-257`]
- **V-F119 已落实为本模块调用约束**：三个桶的 key、窗口、阈值、拒绝语义已冻结；每次请求使用 `AllowRecord`，拒绝后取 `RetryAfterSeconds`，明确禁止 key-wide `Clear` 与 `Reserve/Cancel`。现有 `Clear` 的确是删除该 key 全部历史，因此该禁令必要且方向正确。[`01-decision/D-002-digital-offer-contract.md:193-205`；`apps/api/kernel/ratelimit.go:35-79`]
- **Telegram 与装配先例可达**：`TelegramDispatcher.RegisterCommand` 与 `TelegramUpdate.SubjectID` 真实存在；`plan.HasModule` 条件装配及 provider 的 routes / permissions / pages / navigation 贡献均有 wallet 先例。[`apps/api/kernel/telegram.go:107-126`；`apps/api/modules/wallet/provider.go:254-349`；`apps/api/internal/composition/composition.go:526-614`]

## 对照成功标准

| 标准 | 审计结论 | 证据 / 限制 |
|------|----------|-------------|
| VP-031 判据 1 · Offer CRUD / Admin / 权限 / 审计 / C 端列表 | **部分满足** | Offer、路由、权限、协议页与公开列表已冻结；审计事务与失败语义未冻结，见 F-003。D-002 §2、§5.3、§7。 |
| 判据 2 · 余额不足、freeze→deduct_frozen、同事务 fail-closed、幂等、并发 | **部分满足** | 事务方向及钱包原语可达；精确调用参数、ledger 幂等键、并发重放/冲突恢复未冻结，见 F-001。D-002 §4.2；wallet repository `:576-648`。 |
| 判据 3 · valid / expired / exhausted 与统一核验、原子 Consume | **部分满足** | 惰性判定无需后台 job，方向成立；Consume 的 PostgreSQL 并发算法不足以直接照做，见 F-002。D-002 §3、§5。 |
| 判据 5 · Telegram 可选 | **满足（设计层）** | RegisterCommand、SubjectID 与 disabled dispatcher 路径有现存接缝；空 SubjectID 已要求 fail closed。D-002 §6；`kernel/telegram.go:107-126`。 |
| 判据 7 · 边界保持 | **满足（设计层）** | 红线、Profile 排除、subject-only 与不解禁通用 Entitlement 接缝均明确。D-002 §10。 |
| V-F119 | **满足（设计层）** | 请求计数桶使用 AllowRecord，明确禁止 key-wide Clear。D-002 §8；`kernel/ratelimit.go:35-79`。 |

## Findings

### F-001 · 单事务购买尚未冻结可执行的 wallet mutation 与并发幂等协议

- **严重度**：med
- **建议**：required
- **状态**：open
- **证据**：`01-decision/D-002-digital-offer-contract.md:99-120`；`apps/api/modules/wallet/store/repository.go:155-165`、`:576-600`、`:623-648`、`:722-729`；`apps/api/modules/wallet/voucher/service.go:241-257`
- **问题**：D-002 使用 `walletRepo.MutateInTx(tx, freeze)` / `(..., deduct_frozen)` 的缩写，但真实签名还要求 `accountID`、完整 `LedgerEntryInput`、独立 `entryID` 与 `now`。合同未冻结两笔流水各自的 `IdempotencyKey`、`entryID` 生成与 payload，也未规定 `ErrVersionConflict`、ledger 唯一竞争或并发 `(subject_id, request_id)` 插入冲突后的重试/回读策略。真实 `MutateInTx` 不自动重试；同一 idempotency key 搭配不同 entry type / payload 会返回冲突。
- **风险**：不同实现可选择相同 key、空 key 或不稳定 key；并发重复购买可能返回 wallet version conflict 而非既有凭证，或无法证明“一凭证、一扣款、一份权益”的重放合同。合同目前方向可行，但不能仅按现文无歧义实现和验收。
- **修复 / 验证方向**：在 D-002 冻结 purchase id、freeze/deduct entry id 的预生成顺序；两笔 `LedgerEntryInput` 的 EntryType、AmountDelta、RefType、RefID、稳定且互异的 IdempotencyKey；明确传入 `account.ID`；冻结 `ErrVersionConflict` / ledger idempotency race / purchase UNIQUE race 的有界重试与“回读既有凭证后比较负载”规则。并发测试至少覆盖相同 request 双发、不同 offer 复用 request、不同 request 同账户竞争、重试后无重复 ledger/凭证/权益。

### F-002 · Consume 的跨行并发与 void 竞争规则不足，`事务串行化` 不是现有跨方言保证

- **严重度**：med
- **建议**：required
- **状态**：open
- **证据**：`01-decision/D-002-digital-offer-contract.md:142-149`；`apps/api/kernel/store.go:27-47`；`apps/api/internal/store/postgres.go:231-259`
- **问题**：合同只给出逐行 `UPDATE ... WHERE id=? AND remaining_count>=?`，没有冻结候选行读取/锁定方式、`RowsAffected=0` 后继续选行还是重算、跨多行部分扣减如何在竞争后重新分配，以及与 Admin `void` 的顺序/谓词。PostgreSQL 先例使用 `BeginTx(ctx, nil)`，并不承诺 serializable 事务；因此“RowsAffected 守卫 + 事务串行化”不能作为现有平台的通用事实。示例 UPDATE 也未把 `subject_id`、`offer_id`、`form='count'`、`status='active'` 纳入最终写入谓词。
- **风险**：实现者可产生并发下的假性不足、错误消费已作废/非目标 Offer 行，或在跨行消费时得到方言相关行为；这直接影响判据 3 与 D-002 自称的“并发消费不得超扣”。
- **修复 / 验证方向**：冻结一个 SQLite/PostgreSQL 均可证明的算法：候选条件与稳定排序、最终 UPDATE 的主体/Offer/形态/active/余额谓词、冲突后重算或有界重试、任一失败回滚全部已扣行、与 void 的线性化规则。测试覆盖单行与多行、总量恰够、两个消费者都应成功、仅一方应成功、RowsAffected 竞争、并发 void，且双数据库运行。

### F-003 · Admin 写操作的审计只声明 `Recorder`，未冻结同事务 fail-closed 边界

- **严重度**：med
- **建议**：required
- **状态**：open
- **证据**：`01-decision/D-001-info-adjudication.md:18-20`；`01-decision/D-002-digital-offer-contract.md:34-35`、`:54-55`、`:168-190`；`apps/api/modules/operationlog/repository.go:112-120`、`:149-158`
- **问题**：D-001 要求权益作废有独立权限键与审计，VP-031 判据 1 要求 Offer 变更有审计；D-002 却只把 `operationlog.Recorder` 作为依赖，并未规定审计失败时是否回滚域写入。仓库已提供 `TransactionalRecorder.RecordOperationTx(kernel.Tx, Operation)`，专用于必须与业务变更 fail-closed 的 mutation audit；普通 `Recorder.RecordOperation` 自启事务，调用方可选择 best-effort。
- **风险**：Offer 上下架/改价或权益作废成功后，审计写入失败仍可能留下不可追溯变更，无法稳定满足判据 1 与 D-001 的作废审计裁决。
- **修复 / 验证方向**：冻结 Offer create/edit/on-sale/off-sale 与 entitlement void 的事件名、actor、record id、before/after 或必要 detail、correlation id；要求域写入与 `RecordOperationTx` 共用 caller-owned transaction，审计失败整体回滚。补 forced audit failure 测试，证明业务行与审计行同成同败；同时定义“重复 void 返回成功”是否追加审计及其幂等规则。

### F-004 · 权益 reason 与 Consume sentinel 到冻结错误码的映射不完整

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **证据**：`01-decision/D-002-digital-offer-contract.md:134-149`、`:207-220`
- **问题**：`Check` 定义了 `voided` / `expired` / `exhausted` / `no_entitlement`，但混合多行状态的优先级未完全定义；`Consume` 返回 `ErrEntitlementInsufficient`，错误码表仅给出总括的 `BIZOFFER_ENTITLEMENT_INVALID`，未明确 sentinel/reason 到 HTTP/Telegram 错误的映射。
- **风险**：不同入口可能对同一权益集合返回不同 reason、状态码或文案，削弱判据 3 的可测性。
- **修复 / 验证方向**：冻结混合状态优先级和 `ErrEntitlementInsufficient` 映射；用表驱动测试覆盖多份权益的 active/voided/expired/exhausted 组合及 HTTP/Telegram 一致性。

### F-005 · 工作区阶段说明已落后于目标树与 D-001 事实

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **证据**：`docs/workspaces/workspace-031-digital-offer-entitlement/workspace.md:33-40`；`docs/workspaces/workspace-031-digital-offer-entitlement/goal-tree.md:3-18`；`01-decision/D-001-info-adjudication.md:36-39`
- **问题**：`workspace.md` 仍写 R1“待开始、I-031-001～005 尚未处理”，而 goal-tree 与 D-001 已记录 R1 进行中且信息项 verified。
- **风险**：不改变当前目标状态真相，但会误导后续只读扫描与审计范围判断。
- **修复 / 验证方向**：由 `/govern` 在响应本意见时同步 workspace 阶段投影；本 independent 审计不修改该文件。

## 必改项汇总

- **F-001（required）**：冻结可直接调用的 wallet mutation 参数、两笔流水幂等键/entry id，以及并发购买冲突的重试与回读协议。
- **F-002（required）**：补齐 Consume 在 PostgreSQL/SQLite 下的跨行并发、RowsAffected 竞争和 void 线性化算法。
- **F-003（required）**：Admin 关键写操作改为同事务 `TransactionalRecorder` fail-closed 审计合同，并冻结事件与失败测试。
- open required = **3**。在三项按 `fixed`、`accepted-residual` 或 `user-overruled` 合法闭合前，不应将 D-002 从 draft 冻结为 accepted，也不应放行依赖这些条款的 R2/R3 实施门禁。

## 与既有意见的异同（对照 A-001）

- **一致**：认可 A-001 的正面判断——D-002 与 D-001 框架一致，钱包事务原语、惰性有效判定、V-F119 和 Telegram/provider 接缝方向可达。
- **复核确认**：A-001 F-001～F-005 的五项文本修复均已真实落入 D-002；本意见不重新打开这些原 finding。
- **新增差异**：A-001 主要检查“是否有对应条款”；本 independent 审计进一步按“能否照做”检查真实函数签名、并发隔离与审计失败路径，新增 F-001～F-003 required 及 F-004～F-005 recommended。
- **verdict 差异**：A-001 为 conditional 且其五项已 fixed；A-002 仍为 conditional，因为新发现 3 项 med required 尚未闭合。不存在相反 verdict 冲突，但 A-001 的闭合不能替代本条 findings 的响应。

## 结论 + 建议给编排器的下一步

D-002 的业务边界与主路径总体正确，D-001、VP-031 判据 5/7 与 V-F119 对齐良好，钱包/Telegram/provider 主接缝也真实存在；但购买幂等协议、Consume 跨方言并发算法、Admin 审计事务边界尚未冻结到可无歧义实施的程度。因此本次 **verdict = conditional**，不是无条件放行依据。

建议由 `/govern`：

1. 逐条响应 F-001～F-003，优先选择 `fixed` 并修订 D-002；F-004～F-005 可一并处理。
2. 修订后追加执行事实与 self 响应，不覆盖 A-001/A-002 原文。
3. 在进入 R2/R3 前做 focused independent finding-closure 复审；未获用户书面 residual/overrule 时，不得绕过 open required。

## 声明

本意见仅写入 independent audit ledger；不修改 `status`、`progress`、goal-tree 状态列、D-001/D-002 正文或任何 01/02 台账。finding 响应与状态推进由 `/govern` 处理。
