---
id: E-002
goal: GOAL-021-wallet-deduct-frozen
title: S2 实现 + S3 验证 + S4 判定 + S5 关门
date: 2026-08-16
status: recorded
parent: GOAL-021-wallet-deduct-frozen
created: 2026-08-16
updated: 2026-08-16
version: 1.0.0
---

# E-002 · 实现与关门（2026-08-16）

## 事实

- **S2 实现**：store Apply 加 EntryDeductFrozen（total-=d、frozen-=d、available 不变；d<=0 → ErrInvalidEntry；frozen<d → ErrInsufficient）；幂等比对纳入 RefType/RefID（F-002）；迁移 0033（ledger 表 CHECK 超集重建，rename/create/copy/drop 保数据）+ 0034（operationlog 事件 CHECK）；handler POST /api/wallet/accounts/{id}/deduct-frozen（wallet.adjust + wallet.deduct-frozen 审计）；schema 行操作/modal + i18n en/zh；Descriptor/BuiltinModules 同步。
- **S3 验证**：go 全量全绿 + web 全量全绿（D-VAL 全模块 21 文档 + schema-keys）；新增测试：TestMutateDeductFrozen（750/600/150 + 超额拒绝 + 对账一致）、TestMutateIdempotencyRefCompare（同 key 异单据 conflict + 余额不变）、TestApplyDeductFrozenWithFrozenBalance、TestDeductFrozenPreciseSentinels（ErrInvalidEntry/ErrInsufficient/ErrDisabled 精确）、TestWalletDeductFrozenEndpoint（含 409 码体 INSUFFICIENT_BALANCE + 审计事件）、TestMigrate0033PreservesWalletLedgerRows（有流水升级保数据）、迁移冻结 32→34（checksum 0033 b3135b28…、0034 b6b54bee…）。
- **S4 go 判定（D-002）**：加法路由、无权限/导航/Profile/协议变更 → 不暂挂。
- **S5 关门**：A-002（grok independent · data 门禁）**pass（0 required）**；recommended F-001（有流水升级用例）→ fixed（TestMigrate0033PreservesWalletLedgerRows）；F-002（精确哨兵/码体）→ fixed；F-003（台账）→ 本条目落盘。
- **A-008 响应闭环（GOAL-019）**：F-001（冻结扣款原语）→ fixed；F-002（幂等比对）→ fixed——已在 GOAL-019 03-audit 响应记录留痕；F-003~F-011 演进登记于 D-001 §5。
