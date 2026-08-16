---
id: D-001
goal: GOAL-021-wallet-deduct-frozen
title: 方案冻结：deduct_frozen + 幂等修复 + 演进登记
date: 2026-08-16
status: accepted
parent: GOAL-021-wallet-deduct-frozen
created: 2026-08-16
updated: 2026-08-16
version: 1.0.0
---

# D-001 · 方案冻结（冻结扣款 + 幂等修复）

> 依据：A-008 F-001/F-002；用户裁决（2026-08-16）；GOAL-019 D-002 apply 表先例。

## 1. deduct_frozen 语义（I-001 闭合）

| 字段 | 值 |
|------|-----|
| 动作 | available 不变；total -= d；frozen -= d |
| 符号 | d > 0（正值消费额） |
| 拒绝 | d <= 0（ErrInvalidEntry）；frozen < d（ErrInsufficient）；账户 disabled（ErrDisabled） |
| 快照 | balance_after_* 三列照常记录（恒等式 total = available + frozen 仍成立） |
| 权限 | wallet.adjust（金额变动专用键） |
| 幂等 | 沿用 idempotencyKey（复合 UNIQUE + 载荷比对） |
| 审计 | wallet.deduct-frozen 事件（detail: accountId/entryId/amountDelta） |

## 2. 幂等载荷比对修正（I-003 闭合）

- Mutate 幂等判定由 EntryType/AmountDelta/Memo 扩展为 **+ RefType + RefID**；同 key 异单据 → LEDGER_IDEMPOTENCY_CONFLICT。
- 兼容：既有同 key 重放均为同载荷（含单据），语义不变；测试补异单据用例。

## 3. 迁移（I-002 闭合）

- **0033**：wallet_ledger_entries 重建——entry_type CHECK 扩为 ('adjust','freeze','unfreeze','deduct_frozen')；rename→create→copy→drop + 索引/UNIQUE 重建（operationlog rebuild 先例；数据原样保留）。
- **0034**：operationlog 事件 CHECK 重建加 wallet.deduct-frozen（0032 同款模式）。
- 编号衔接：当前 max 32 → 33/34；迁移冻结测试（migrate_test want 列表）同步。

## 4. 端点与前端

- `POST /api/wallet/accounts/{id}/deduct-frozen`（wallet.adjust）：body {amount, memo, idempotencyKey?, refType?, refId?} → {account, entry}。
- schema：钱包页行操作加「扣减冻结」modal（amount/memo/idempotencyKey 字段）；i18n en/zh 键。
- Descriptor/BuiltinModules 路由声明同步。

## 5. 演进登记（F-003~F-011，不实现）

| ID | 演进方向 | 触发条件 |
|----|----------|----------|
| F-003 | 复式记账（借贷分录 + 资金头寸） | 出现外部对账/审计合规要求时 |
| F-004 | 原子转账（Transfer 事务） | 出现账户间资金划转业务时 |
| F-005 | 交易业务类型枚举（充值/消费/退款/提现等） | 财务分类核算需求明确时 |
| F-006 | 对账快照/增量重放（替代全量） | 热点账户流水量达阈值（如 >10 万）或对账超时实测时 |
| F-007 | 热点账户并发（悲观锁/分片/异步串行） | 实测高并发 409 风暴时 |
| F-008 | 调账风控（单笔/日累计限额 + Maker-Checker） | 资金内控审计要求时 |
| F-009 | operationlog 同事务（维持现有残余） | 审计完整性要求升级时（流水 actor 保底已覆盖） |
| F-010 | 多币种与可变精度 | 多币种/积分/代币业务立项时 |
| F-011 | 前端金额格式化 + 组合筛选 | 财务报表体验需求时 |

## 6. 测试与验证

- Apply：deduct 正常/超额/零/负/disabled；快照恒等式。
- Mutate 集成：deduct 后三余额、幂等重放、异单据冲突。
- handler：端点门禁/审计事件/错误码。
- 迁移：0033/0034 冻结（want 列表 32→34）、重建保数据（有流水后升级）。
- 全量回归：go + web。

## 7. 未选方案（留痕）

- F-001 走 accepted-residual：冻结资金扣减闭环是资金安全缺口，不采纳残余。
- 演进项本批实现：用户裁决全部登记（避免一次摊大）。
- 新错误码：复用 INSUFFICIENT_BALANCE/INVALID_LEDGER_ENTRY，不新增。
