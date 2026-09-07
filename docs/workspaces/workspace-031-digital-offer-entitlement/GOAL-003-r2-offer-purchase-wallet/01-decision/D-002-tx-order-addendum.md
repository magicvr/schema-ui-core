---
doc_type: goal-decision
id: D-002-tx-order-addendum
parent: GOAL-003-r2-offer-purchase-wallet
date: 2026-09-05
status: accepted
version: 1.0.0
---

# D-002 · 合同附录：购买事务内步骤重排（D-002 v1.1.0 → v1.2.0）

## 决定

A-002（independent）审计期间在真 PostgreSQL 上运行并发双发验收时发现：并发同 `(subject_id, request_id)` 的各 Purchase 调用**各自预生成不同 purchaseID**，但共享同一 wallet 幂等键 `<request_id>:freeze` / `:deduct`；在 PG READ COMMITTED 下败者事务可携带「同键不同 RefID 负载」触达 ledger，触发 wallet `ErrIdempotencyConflict`（终态）——破坏 §4.4 收敛论证（SQLite 因写者串行未暴露该窗口）。

**修订**：购买事务内步骤重排为「校验 → **INSERT 凭证（请求守卫）** → freeze → deduct_frozen → INSERT 权益」。唯一约束 `UNIQUE(subject_id, request_id)` 在任何钱包流 水之前裁决；败者在守卫处失败并按 §4.4 回读重放，**永不触碰 wallet ledger**。happy-path 的原子性语义不变（同事务全有或全无，凭证与流 水同回滚）。

已同步 D-002 §4.2 / §4.4 条款并升 **v1.2.0**；实施代码与测试（含 PG 并发矩阵）按 v1.2.0 执行并通过。

## 未选方案

- purchaseID 改为确定性哈希：破坏 id 的时间有序设计（wallet v50 order-repair 教训）且同样属于合同偏离。
- 将 wallet `ErrIdempotencyConflict` 在本域降级为可重试：会掩盖真实的幂等负载漂移，弱化资金路径 fail-closed。
- 维持原顺序、限制并发：调用方可控性不足，收敛性依赖调用方纪律，不可验收。

## 影响

- §4.2 步骤编号与 §4.4 attempt 顺序同步更新；验收矩阵（SQLite + 真 PG）覆盖重排后的并发收敛。
- 判据、红线、D-001 裁决框架不变。
