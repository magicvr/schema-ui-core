---
doc_type: goal-audit
id: A-006-independent-closure-review-2
status: closed
parent: GOAL-002-r1-contract-freeze
created: 2026-09-05
updated: 2026-09-05
date: 2026-09-05
source: independent
auditor: codex (gpt-5.6-sol, medium)
audit_type: finding-closure
scope: 仅复审 A-005 对 A-004 F-001～F-003 的关闭证据；核对 D-002 现行文本、workspace.md R1 投影与相关代码事实
verdict: pass
open_required: 0
version: 0.1.0
---

# A-006 · A-004 finding-closure 独立复审（第 2 轮，2026-09-05）

- **source**：independent
- **auditor**：codex (gpt-5.6-sol, medium)
- **类型 / scope**：finding-closure · 仅复审 A-005 是否真实、充分、可核对地闭合 A-004 F-001～F-003；不重审整个 R1，也不把 R2/R3 尚未发生的数字 Offer 实现与验收测试纳入本次合同冻结分母。
- **verdict**：**pass**
- **open required**：0

## 范围与区间

- 核对对象：A-004 原 finding 与关闭方向、A-005 关闭证据表、D-002 现行 §4.2/§4.4/§5.2/§7、`workspace.md` R1 纲领阶段行，以及 wallet、`kernel.Store`、SQLite/PostgreSQL `Store.Run` 与唯一违反分类的现行代码事实。[`03-audit/A-004-independent-closure-review.md:29-66`；`03-audit/A-005-self-response-a004.md:10-38`]
- 工作区绑定、canonical scope 与 VP 挂接有效；共享资料目录为 `none`，本次未使用外部共享材料。[`docs/workspaces/workspace-031-digital-offer-entitlement/workspace.md:20-31`；`docs/workspaces/workspace-031-digital-offer-entitlement/GOAL-002-r1-contract-freeze/00-meta.md:1-15`]
- 本次判断的是“修订后的合同是否已消除 A-004 指出的不可执行/不真实主张，并能按现有平台接口无歧义实施”，不是证明尚未进入 R2/R3 的数字 Offer 代码已经实现。[`03-audit/A-004-independent-closure-review.md:23-27`；`01-decision/D-002-digital-offer-contract.md:162,222`]

## 关闭证据核对表

| Finding | A-005 闭合声明 | 核对结论 + 证据 |
|---------|----------------|-----------------|
| A-004 F-001 · wallet ledger 唯一竞争 sentinel 对域调用方不可见 | `fixed` | **闭合成立**。修订后的 §4.4 不再要求 `digitaloffer` 识别 wallet/store 包内未导出的 `errIdempotencyRace`：每个 attempt 先回读 `(subject_id, request_id)`；仅对本域 `digital_purchases` INSERT 的唯一违反使用公开的 `kernel.IsUniqueViolation`；显式终态为余额不足、币种不匹配、非 on_sale、subject 不存在与 `ErrRequestIdConflict`，其余非终态错误统一进入最多 3 次的整事务有界重试。`kernel.IsUniqueViolation` 为公开函数并覆盖 SQLite/PostgreSQL 唯一违反；wallet 的未导出 sentinel 仍只在包内产生，但已无需跨包分类。[`01-decision/D-002-digital-offer-contract.md:122-160`；`apps/api/kernel/unique_violation.go:8-35`；`apps/api/modules/wallet/store/repository.go:48-58,637-648`] 同 request 并发时，失败 attempt 整体回滚；胜者提交后，失败方下一 attempt 在 wallet mutation 前回读并按相同 offer 幂等重放、按不同 offer 返回 `ErrRequestIdConflict`。SQLite 以 `BEGIN IMMEDIATE`、5 秒 busy timeout 令写者排队；PostgreSQL 的唯一约束/乐观更新竞争在事务完成后可由新事务重读，因此该收敛论证与平台事实相容。[`D-002 §4.4，:124-160`；`apps/api/internal/store/store.go:48-57`；`apps/api/modules/wallet/store/repository.go:623-648`] 两笔 wallet 调用仍冻结为互异键 `<request_id>:freeze` / `<request_id>:deduct`，且 wallet 确实按调用方传入的 `IdempotencyKey` 回读并落账。[`D-002 §4.4，:128-138`；`repository.go:583-600,637-642`] |
| A-004 F-002 · Consume 重试循环与 `Store.Run` 事务边界矛盾 | `fixed` | **闭合成立**。§5.2 已把 `for attempt in 1..3` 放到 `store.Run` 外，每次 attempt 调用一个全新事务；callback 内无 COMMIT/ROLLBACK，而以 `nil`/`committed-OK`、`deterministic-insufficient`、其余竞争失败三类结果驱动提交、回滚或下一 attempt。[`D-002 §5.2，:184-218`] 这与 `kernel.Store.Run`“一次 Run 一个事务、禁止嵌套 Run”、`kernel.Tx` 不暴露 Commit/Rollback，以及 SQLite/PostgreSQL runner 统一在 callback 外 BeginTx/rollback/commit 的代码事实一致。[`apps/api/kernel/store.go:27-53`；`apps/api/internal/store/store.go:145-170`；`apps/api/internal/store/postgres.go:231-259`] 隔离级别声明也已改为平台不固定：PostgreSQL 实现实际调用 `BeginTx(ctx, nil)`，合同明确不依赖 serializable，而依赖 UPDATE 最终谓词重检、全新事务重读及有界重试。[`D-002 §5.2，:219-220`；`postgres.go:231-240`] 不足路径通过 callback 返回错误令本 attempt 回滚，部分扣减不提交；成功 attempt 才原子提交全部扣减，故“不足时无消费”的全有或全无语义保留。[`D-002 §5.2，:207-220`；`store.go:155-170`；`postgres.go:243-259`] |
| A-004 F-003 · workspace 投影滞后 | `fixed` | **闭合成立**。R1 行现已逐项投影 A-003 响应、A-004 independent closure `fail`（2 required：跨包错误契约 / Consume 事务边界）以及 A-005 已修订、待 A-006 closure 复审，与三份审计条目的 source/verdict/open-required/响应关系一致。[`docs/workspaces/workspace-031-digital-offer-entitlement/workspace.md:33-40`；`03-audit/A-003-self-response-a002.md:1-24`；`03-audit/A-004-independent-closure-review.md:1-21`；`03-audit/A-005-self-response-a004.md:1-38`] |

## 一致性核对

- §4.2 定义单个购买 attempt 内的 `freeze → deduct_frozen → purchase → entitlement` 原子步骤与失败整体回滚；§4.4 只补充 attempt 外层重试、幂等回读和错误分类，两者层次一致，无“已提交 freeze 后再补偿”的分段事务矛盾。[`D-002 §4.2，:99-120`；`§4.4，:122-162`]
- §5.2 的消费事务独立于购买事务，但同样由 `Store.Run` callback 返回值保持全有或全无；§7 的 operationlog fail-closed 约束适用于 Admin 域写操作，购买路径则以不可变购买凭证和 wallet ref 反链追溯。两处没有放宽失败回滚或要求 callback 自行提交，未发现直接矛盾。[`D-002 §5.2，:184-222`；`§7，:245-267`]
- 修订内容仍符合 D-001 的同步 `fulfilled`、无 `pending`、Admin 不人工发放裁决，也未与 VP-031 的同事务购买、失败 fail-closed、并发验收及 required finding 闭合判据冲突。[`01-decision/D-001-info-adjudication.md:10-18`；`docs/vision/plans/VP-031-digital-offer-entitlement.md:63-72`]

## 新 Findings

### F-001 · §5.2 “读时已锁行”措辞与实际 SELECT/隔离声明不一致

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **证据**：D-002 §5.2 的 SELECT 未写 `FOR UPDATE`，且同节明确 PostgreSQL 使用 `BeginTx(ctx, nil)`、平台不固定隔离级别；但确定性不足分支注释写成“读时已锁行语义下成立”。[`01-decision/D-002-digital-offer-contract.md:191-220`；`apps/api/internal/store/postgres.go:231-240`]
- **影响**：这是证明措辞不精确，不推翻本次 F-002 closure。现行算法的有效依据是最终条件 UPDATE 的谓词重检、失败 attempt 回滚，以及新事务重读；并不需要把普通 SELECT 描述为锁行读取。
- **建议方向**：后续由 `/govern` 将该注释改为基于条件 UPDATE/线性化点的表述；若确实意图使用锁行读，则应显式冻结跨方言可执行 SQL 与测试分母。该项不阻断 R1 合同 closure，但应在进入相应实现前消除歧义。

## 结论 + 建议下一步

- A-005 对 A-004 F-001、F-002、F-003 的关闭声明均有现行文本与代码接口事实支撑；A-004 的 2 项 med required 已按 `fixed` 路径充分闭合，**open required = 0**。
- 本次未发现与 VP-031 判据、D-001 裁决或现行平台接口直接冲突的新 required 缺陷；仅新增 1 项 low recommended 的证明措辞清理。
- 建议由 `/govern` 响应 A-006：登记 A-004 F-001～F-003 为 fixed closure，处理或排期 A-006 F-001，再按既有门禁决定是否冻结 D-002/推进 R1；本审计不代替治理状态推进。

## 声明

本意见仅追加 `source: independent` 审计记录，不修改目标 `status` / `progress`、goal-tree 状态列、D-001/D-002 正文或 01/02 台账；响应、finding 状态闭合与阶段推进由 `/govern` 处理。
