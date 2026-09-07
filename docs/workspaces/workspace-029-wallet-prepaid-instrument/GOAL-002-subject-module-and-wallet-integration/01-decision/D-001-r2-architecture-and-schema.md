---
doc_type: goal-decision
id: D-001-r2-architecture-and-schema
parent: GOAL-002-subject-module-and-wallet-integration
date: 2026-09-02
status: accepted
version: 0.1.0
---

# D-001 · R2 实施架构设计与 DDL 规格

## 触发

承接 `GOAL-001` D-002 裁决，推进外部主体接缝与凭证核销入金的核心代码落地。

## 决定

1. **迁移 0064 规划**：
   - 建立 `subjects` 表：`id`, `issuer`, `external_id`, `created_at`，唯一约束 `UNIQUE(issuer, external_id)`。
   - 建立 `vouchers` 表：`id`, `batch_id`, `code_hash`, `code_prefix`, `amount`, `currency`, `status`, `expires_at`, `redeemed_by`, `redeemed_at`, `created_at`, `updated_at`，唯一约束 `UNIQUE(code_hash)`。
   - 保证双方言（SQLite / PG）的一致性。
2. **外部主体模型**：
   - 通道无关独立包 `modules/wallet/subject` 或直接扩展至模块内独立存储，支持 `GetOrCreateSubject(ctx, issuer, externalID)`。
   - 钱包 `OwnerExistsFunc` 适配：区分主体与系统用户，未登记主体禁止开户，禁止回退查 `admin.users`。
3. **预付凭证核销**：
   - 单事务内 CAS 校验并标记 `status='redeemed'`，影响行数=0 则报错。
   - 同事务调用钱包调账原语（`adjust`），`ref_type='voucher'`, `ref_id=voucher.id`，幂等键填 `voucher.id`。
   - 外部调用仅暴露 `Redeem(ctx, subjectID, code)` 模块方法。
