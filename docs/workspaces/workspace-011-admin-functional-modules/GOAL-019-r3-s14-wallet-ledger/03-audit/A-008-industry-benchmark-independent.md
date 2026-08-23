---
id: A-008
goal: GOAL-019-r3-s14-wallet-ledger
title: 业界对标独立审计（冻结扣款/复式记账/转账/风控/对账性能）
date: 2026-08-16
source: independent
auditor: 用户提供的独立审视（业界钱包/账务系统对标）
scope: admin.wallet 模块整体（GOAL-019 方案与实现 + GOAL-020 自动开户产物）
audit_type: ad-hoc
verdict: conditional
status: recorded
parent: GOAL-019-r3-s14-wallet-ledger
created: 2026-08-16
updated: 2026-08-16
version: 1.0.0
---

# A-008 · 业界对标独立审计（independent）

- **source**：independent（用户代贴的独立审视，未由本编排器起草）
- **auditor**：用户提供的独立审视（业界主流钱包/账务系统对标）
- **scope**：admin.wallet 模块整体——GOAL-019（S-14 钱包/账务）方案与实现、GOAL-020（自动开户）产物
- **verdict**：**conditional**（存在未闭合 P0 缺陷：冻结扣款原语缺失 + 幂等载荷比对遗漏 refType/refId）
- **完整意见**：[attachments/audit-A-008-independent.md](../attachments/audit-A-008-independent.md)

## 范围与区间

- 审计对象：apply 表语义、流水不可变设计、乐观锁/幂等、权限分级、惰性开户、对账实现、前端展示、风控与合规面。
- 与既有意见关系：A-007（grok close-out pass）结论不因此撤回；本条为**行业演进视角**的补充意见，其 findings 中 P0 项按 required 处理（未闭合前相关演进门禁不放行），P1/P2 为 recommended/non-blocking 演进方向。

## 成果（做得好的方面，有证据）

| 项 | 证据 |
|----|------|
| 三态余额模型 + DDL 硬约束（恒等式 CHECK、balance_* >= 0） | migration.go 0031；Apply() 应用层 + DDL 双校验 |
| 不可变账本（append-only + balance_after 快照） | wallet_ledger_entries 仅 INSERT/SELECT；快照链对账 |
| 并发与幂等（version 乐观锁 + (account_id, idempotency_key) 复合唯一） | store.Mutate；TestMutateIdempotency |
| 资金权限分级（read/write/adjust） | handler 端点绑定；A-007 核对 |
| 用户钱包惰性开户（get-or-create + auto 审计） | GOAL-020 D-001；GetOrCreateUserAccount |

## Findings

### F-001 ·【P0 · required】缺失冻结扣款/结算（deduct_frozen）原语

业界预授权/押金闭环：available → freeze → **deduct_frozen（冻结消费）** 或 unfreeze。当前 Apply 无 deduct_frozen：冻结资金只能先 unfreeze 再 adjust(-d)，解冻瞬间资金暴露为可用（可被并发消费），且两步间失败会滞留资金——破坏预授权原子性与资金安全。

### F-002 ·【P0 · required】幂等载荷比对遗漏 refType/refId

Mutate 幂等判定仅比较 EntryType/AmountDelta/Memo；同 idempotencyKey 换单据号重试会被误判为同载荷返回旧流水——新单据未入账却返回成功（静默丢单）。

### F-003 ·【P1 · recommended】单边账 vs 复式记账（会计模型局限）

无借贷分录/资金头寸概念，无法验证平台总余额 = 总准备金；调账资金凭空产生/消失，无出资方追踪。

### F-004 ·【P1 · recommended】缺原子转账能力

跨账户资金转移需业务方拆两次 adjust，存在单边账风险（A 扣成功 B 加失败）。

### F-005 ·【P1 · recommended】交易业务类型/状态机缺失

全部变动归 adjust + memo 字符串；业界有严格枚举（充值/消费/退款/提现/分润/罚扣等）便于核算与报表。

### F-006 ·【P1 · recommended】对账 O(N) 全量内存重放（热点账户风险）

checkAccountChain 从创世块逐笔捞取全量历史重放；热点账户（大商户/系统账户）数十万流水时长事务 + OOM 风险。业界：期初快照 + 期间增量重放，或离线批处理。

### F-007 ·【P1 · recommended】纯乐观锁高并发重试风暴

热点账户并发调账大量 409；业界：悲观行锁排队 / 子账户分片 / 异步串行入账。

### F-008 ·【P2 · recommended】大额调账风控缺失（Maker-Checker/限额）

wallet.adjust 可单人任意额度调账，无单笔/日累计限额与审批流。

### F-009 ·【P2 · recommended】operationlog 与记账非同事务

账本事务提交后才写操作日志，崩溃窗口审计脱节（底层流水 actor_id 保底）；GOAL-019 D-002 §2 已文档化为残余。

### F-010 ·【P2 · non-blocking】单币种与固定精度

硬编码 CNY + 分精度；积分/JPY/代币（4-18 位小数）无法适配。

### F-011 ·【P2 · non-blocking】前端金额格式化与筛选维度

表格直接显示整数分（如 1000 而非 ￥10.00）；筛选仅 ownerId LIKE，缺时间范围/金额区间/流水类型组合筛选。

## 必改项汇总

| ID | 级别 | 一句话 |
|----|------|--------|
| F-001 | required · P0 | 补 deduct_frozen 冻结扣款原语（apply 表 + 端点 + 审计） |
| F-002 | required · P0 | 幂等载荷比对纳入 RefType/RefID |

F-003~F-011 为 recommended / non-blocking 演进方向。

## 结论 + 建议

基础数据完整性/不可变审计/并发控制扎实；真实业务闭环（冻结扣款）、会计模型（复式）、性能（对账/热点）与风控合规存在演进盲区。建议 /govern：F-002 立即 fixed（纯缺陷）；F-001 经用户裁决后纳入实现；F-003~F-011 登记演进方向（按需立项）。

## 声明

本意见不修改 status / progress / 方案正文。响应由 /govern 处理。保证等级 L0。
