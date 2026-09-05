# goal-tree · workspace-031-digital-offer-entitlement

*自动同步工作区扁平目标树（树 + 状态表）。更新任一目标状态/进度后必须同步本文件。更新：2026-09-05（R1 关门：GOAL-002 done 3/3，D-002 v1.0.0 accepted；创建 GOAL-003 启动 R2）*

## 目标树

```text
GOAL-001-digital-offer-entitlement (数字 Offer 与权益 · active · 1/4)
├── GOAL-002-r1-contract-freeze (R1 合同冻结 · done · 3/3)
└── GOAL-003-r2-offer-purchase-wallet (R2 Offer CRUD + 购买 + 钱包扣款 · active · 0/4)
（R1 关门；R2 进行中；R3～R4 未立项，按阶段渐进开设）
```

## 状态表

| id | title | status | progress | parent | notes |
|----|-------|--------|----------|--------|-------|
| GOAL-001-digital-offer-entitlement | 数字 Offer 与权益 | **active** | 1/4 | null | VP-031 lead Root。R1 关门（2026-09-05）：D-001 用户裁决 + D-002 v1.0.0 accepted（A-001～A-007 审计闭合，A-006 independent `pass` open required 0）。R2 进行中（GOAL-003），实施分母 = D-002 v1.0.0。 |
| GOAL-002-r1-contract-freeze | R1 合同冻结（Offer 字段/购买状态机/权益形态/命令清单/事务与限流边界） | **done** | 3/3 | GOAL-001-digital-offer-entitlement | 2026-09-05 关门：C1 用户裁决（I-031-001～005 verified）+ C2 D-002 v1.0.0 accepted + C3 审计循环 A-001～A-007（self ×3 / independent codex gpt-5.6-sol ×2，A-004 fail 整改后 A-006 pass）。 |
| GOAL-003-r2-offer-purchase-wallet | R2 Offer CRUD + 购买 + 钱包扣款 | **active** | 0/4 | GOAL-001-digital-offer-entitlement | 2026-09-05 设立。检查点：C1 切片方案 → C2 模块/迁移/装配落地 → C3 购买链路 + 并发验收测试 → C4 审视与关门（self + independent）。 |
