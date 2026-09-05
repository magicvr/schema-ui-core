# goal-tree · workspace-031-digital-offer-entitlement

*自动同步工作区扁平目标树（树 + 状态表）。更新任一目标状态/进度后必须同步本文件。更新：2026-09-05（关门撤回：A-006 runtime-integration `fail` 4 required；D-003 裁决 + F-003/F-006 修复完成，待 closure 复审重新关门）*

## 目标树

```text
GOAL-001-digital-offer-entitlement (数字 Offer 与权益 · active · 3/4)
├── GOAL-002-r1-contract-freeze (R1 合同冻结 · done · 3/3)
├── GOAL-003-r2-offer-purchase-wallet (R2 Offer CRUD + 购买 + 钱包扣款 · done · 4/4)
├── GOAL-004-r3-entitlement-validation-telegram (R3 权益核验/消耗 + 可选 Telegram 注册 · done · 3/3)
└── GOAL-005-r4-evidence-closeout (R4 证据矩阵/边界核账/关门审计 · active · 1/2)
（R1～R3 关门；R4 重开：A-006 runtime-integration 整改中）
```

## 状态表

| id | title | status | progress | parent | notes |
|----|-------|--------|----------|--------|-------|
| GOAL-001-digital-offer-entitlement | 数字 Offer 与权益 | **active** | 3/4 | null | VP-031 lead Root。R1 关门（2026-09-05）：D-001 用户裁决 + D-002 合同 accepted（现行为 v1.2.0，R2 期间两次附录；A-001～A-007 审计闭合，A-006 independent `pass` open required 0）。R2 关门（2026-09-05）：GOAL-003 done 4/4（A-001～A-008 审计循环：self ×4 / independent ×4，A-008 independent `pass` 0 required）。R3 关门（2026-09-05）：GOAL-004 done 3/3（A-002 independent `pass` 0 required）。R4 关门（2026-09-05）：GOAL-005 done 2/2（A-004 independent closure `pass` 0 required）。**Root 关门**：成功标准 1～8 全部达成（GOAL-005 E-001 证据矩阵）。 |
| GOAL-002-r1-contract-freeze | R1 合同冻结（Offer 字段/购买状态机/权益形态/命令清单/事务与限流边界） | **done** | 3/3 | GOAL-001-digital-offer-entitlement | 2026-09-05 关门：C1 用户裁决（I-031-001～005 verified）+ C2 D-002 v1.0.0 accepted + C3 审计循环 A-001～A-007（self ×3 / independent codex gpt-5.6-sol ×2，A-004 fail 整改后 A-006 pass）。 |
| GOAL-003-r2-offer-purchase-wallet | R2 Offer CRUD + 购买 + 钱包扣款 | **done** | 4/4 | GOAL-001-digital-offer-entitlement | 2026-09-05 关门：C1～C3 实施 + 双库验收；C4 循环 A-001～A-008（A-002/A-004 fail 整改，A-006 conditional 收尾，A-008 independent `pass` 0 required）。 |
| GOAL-004-r3-entitlement-validation-telegram | R3 权益核验/消耗 + 可选 Telegram 注册 | **done** | 3/3 | GOAL-001-digital-offer-entitlement | 2026-09-05 关门：C1/C2 实施（Check/Consume/Telegram/查询桶）+ C3 审计（A-002 independent `pass` 0 required；A-003 处置 2 项 recommended）。 |
| GOAL-005-r4-evidence-closeout | R4 证据矩阵/边界核账/关门审计 | **active** | 1/2 | GOAL-001-digital-offer-entitlement | 2026-09-05 关门撤回：A-006 independent runtime-integration `fail`（F-003 模块注册表/F-004 CRUD 语义/F-005 迁移策略/F-006 组合根验收）——D-003 裁决 + F-003/F-006 修复完成（BuiltinModules 注册 + 组合根验收测试），待 closure 复审重新关门。 |
