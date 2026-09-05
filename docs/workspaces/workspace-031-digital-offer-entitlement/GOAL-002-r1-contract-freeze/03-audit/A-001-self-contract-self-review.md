---
doc_type: goal-audit
id: A-001-self-contract-self-review
parent: GOAL-002-r1-contract-freeze
date: 2026-09-05
status: closed
version: 1.0.0
---

# A-001 · R1 合同自审（self · design-plan）

## A-001 · R1 合同自审（2026-09-05）

- **source**：self
- **auditor**：编排器（/govern 会话内自审）
- **类型 / scope**：design-plan · D-002 数字 Offer 业务域合同 v0.1.0 草案（对照 VP-031 判据 1/2/3/5/7、V-F119、D-001 裁决框架、代码先例可达性）
- **verdict**：**conditional**（2 med required + 3 recommended；逐条修复后复核通过，见文末闭合记录）

### 范围与区间

仅审 D-002 合同正文与 D-001 裁决一致性；不审 R2+ 实现（尚未发生）。工作区绑定：`workspace-031-digital-offer-entitlement` / Root `GOAL-001-digital-offer-entitlement`，与 canonical 一致。

### 成果（有证据）

- D-002 覆盖 VP-031「首波冻结」表全部行（Offer/购买/权益/资金/通道/Profile/事件），判据映射表（§0）逐条可追溯。
- 单事务购买边界（§4.2）与代码事实吻合：`walletstore.Repository.MutateInTx(tx kernel.Tx, ...)` 存在（repository.go:579），`kernel.Tx` 为 dialect 中性事务（kernel/store.go:43）。
- 权益惰性有效判定（§3）满足判据 3 三态可测，无后台 job 依赖。
- 限流桶（§8）覆盖 V-F119：请求计数语义、禁 key-wide `Clear`、阈值冻结。
- D-001 三项 required 裁决与合同条款无冲突。

### 对照成功标准（GOAL-002 方向级）

| 标准 | 状态 | 证据 |
|------|------|------|
| 1 · I-031-001～003 verified | 达成 | D-001；Root/VP-031 台账已回写 |
| 2 · D-002 覆盖首波冻结表 | 达成 | D-002 §0–§10 |
| 3 · 事务边界可执行 | **部分** | F-002：subject 账户获取路径未落合同 |
| 4 · 限流语义冻结 | 达成 | D-002 §8 |
| 5 · 审计闭合 | 未开始 | 本条即 self 审；independent 待跑 |

### Findings

- **F-001 · Consume 缺 offer 维度**
  - 严重度：med；建议：required；状态：open → **closed（fixed）**
  - 描述：D-002 §5.2 `Consume(subjectID, n)` 无 offer 限定，跨 offer 消耗会错扣（购买 A 的次数被服务 B 消耗）。§5.1 Check 以 (subjectID, offerID) 为界，消耗应同界。
  - 修复：合同改为 `Consume(subjectID, offerID, n)`，仅消耗该 offer 的 count 型权益行。

- **F-002 · 购买链路未指定 subject 钱包账户获取路径**
  - 严重度：med；建议：required；状态：open → **closed（fixed）**
  - 描述：§4.2 步骤 4 直接说「subject 钱包账户 currency」，但未落明账户的 get-or-create 语义。代码先例：`GetOrCreateSubjectAccountInTx(tx, subjectID, now)`（repository.go:424，voucher.Redeem 同路径；每 subject 单账户单币种）。
  - 修复：§4.2 增步骤：事务内 `GetOrCreateSubjectAccountInTx` 取得账户；账户 `currency ≠ offer.currency` → `BIZOFFER_CURRENCY_MISMATCH` 拒绝（subject 单账户单币种，跨币种 offer 首波不可购）。

- **F-003 · 钱包流水缺反向引用**
  - 严重度：low；建议：recommended；状态：open → **closed（fixed）**
  - 描述：purchase 存 freeze/deduct entry id，但 wallet_ledger_entries 有 `ref_type/ref_id`（repository.go:121-122）可反链购买，双向对账更强。
  - 修复：§4.2 freeze/deduct 两笔流水的 `ref_type='biz_offer_purchase'`、`ref_id=purchase.id`。

- **F-004 · 下架对既有权益的影响未声明**
  - 严重度：low；建议：recommended；状态：open → **closed（fixed）**
  - 描述：§2 允许 off_sale，但 §3 未声明既有权益是否继续有效。
  - 修复：§3 增注：offer 下架/改价不影响已发放权益的有效性（有效性只由 §3 谓词决定）。

- **F-005 · 幂等冲突检测细节未落**
  - 严重度：low；建议：recommended；状态：open → **closed（fixed）**
  - 描述：§4.2 步骤 3 说「命中 → 原样返回」，但同 `(subject_id, request_id)` 不同 offer 的冲突未指明检测方式。
  - 修复：§4.2 增注：幂等负载 = offerID；既有凭证 `offer_id ≠ 请求 offerID` → `BIZOFFER_REQUEST_CONFLICT`（对齐 wallet ErrIdempotencyConflict 语义）。

### 必改项汇总

- required：F-001、F-002（均已 fixed）。

### 结论 + 建议下一步

- 修复后 D-002 可进入 independent 审计（codex · gpt 5.6 sol · medium）；verdict 由 conditional 转 pass 的前提 = F-001/F-002 修复可核对（见下）。
- 建议下一步：应用修复 → E 记录 → codex independent A-002。

### 闭合记录（2026-09-05 · self 响应）

| Finding | 闭合路径 | 证据 |
|---------|----------|------|
| F-001 | fixed | D-002 §5.2 改为 `Consume(subjectID, offerID, n)` |
| F-002 | fixed | D-002 §4.2 增 `GetOrCreateSubjectAccountInTx` 步骤与币种不匹配拒绝 |
| F-003 | fixed | D-002 §4.2 增钱包流水 ref_type/ref_id 反链 |
| F-004 | fixed | D-002 §3 增下架不影响既有权益声明 |
| F-005 | fixed | D-002 §4.2 增幂等负载冲突检测 |

修复已直接应用到 `01-decision/D-002-digital-offer-contract.md`（同日，版本仍 0.1.0 draft，待 A-002 后转 accepted 1.0.0）；事实记录见 GOAL-002 `02-execution/E-002-contract-audit-cycle.md`。
