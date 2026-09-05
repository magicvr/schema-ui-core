# goal-tree · workspace-031-digital-offer-entitlement

*自动同步工作区扁平目标树（树 + 状态表）。更新任一目标状态/进度后必须同步本文件。更新：2026-09-05（第 3 轮关门撤回：A-012 independent runtime closure re-audit `conditional` 2 required（F-007/F-008 证据缺口）→ D-004 fixed ×2 + A-013 closed ×2；Root/GOAL-005 active、VP-031 active，待 focused independent closure 复审）*

## 目标树

```text
GOAL-001-digital-offer-entitlement (数字 Offer 与权益 · active · 3/4)
├── GOAL-002-r1-contract-freeze (R1 合同冻结 · done · 3/3)
├── GOAL-003-r2-offer-purchase-wallet (R2 Offer CRUD + 购买 + 钱包扣款 · done · 4/4)
├── GOAL-004-r3-entitlement-validation-telegram (R3 权益核验/消耗 + 可选 Telegram 注册 · done · 3/3)
└── GOAL-005-r4-evidence-closeout (R4 证据矩阵/边界核账/关门审计 · active · 1/2)
（第 3 轮关门撤回：A-012 conditional 2 required → D-004 fixed ×2 + A-013 closed ×2 · 待 focused independent closure 复审）
```

## 状态表

| id | title | status | progress | parent | notes |
|----|-------|--------|----------|--------|-------|
| GOAL-001-digital-offer-entitlement | 数字 Offer 与权益 | **active** | 3/4 | null | VP-031 lead Root。R1 关门（2026-09-05）：D-001 用户裁决 + D-002 合同 accepted（现行为 v1.2.0，R2 期间两次附录；A-001～A-007 审计闭合，A-006 independent `pass` open required 0）。R2 关门（2026-09-05）：GOAL-003 done 4/4（A-001～A-008 审计循环：self ×4 / independent ×4，A-008 independent `pass` 0 required）。R3 关门（2026-09-05）：GOAL-004 done 3/3（A-002 independent `pass` 0 required）。R4 曾关门两次（2026-09-05）：GOAL-005 done 2/2（A-004、A-010 各一轮 independent closure `pass`）。**Root 第 3 轮关门撤回（2026-09-05）**：A-012 independent runtime closure re-audit `conditional` 2 required（F-007 迁移 Apply 中途失败/reopen 双方言证据、F-008 Telegram-enabled 真实组合根 + 结构化 Manifest）→ D-004 fixed ×2 + A-013 closed ×2（迁移失败/reopen 测试 SQLite+PG；Telegram-enabled 组合根测试 + seam 透传 + DELETE 负向；apps/api `go test ./...` 全绿）；待 focused independent closure 复审 `pass` 后重新关门。 |
| GOAL-002-r1-contract-freeze | R1 合同冻结（Offer 字段/购买状态机/权益形态/命令清单/事务与限流边界） | **done** | 3/3 | GOAL-001-digital-offer-entitlement | 2026-09-05 关门：C1 用户裁决（I-031-001～005 verified）+ C2 D-002 v1.0.0 accepted + C3 审计循环 A-001～A-007（self ×3 / independent codex gpt-5.6-sol ×2，A-004 fail 整改后 A-006 pass）。 |
| GOAL-003-r2-offer-purchase-wallet | R2 Offer CRUD + 购买 + 钱包扣款 | **done** | 4/4 | GOAL-001-digital-offer-entitlement | 2026-09-05 关门：C1～C3 实施 + 双库验收；C4 循环 A-001～A-008（A-002/A-004 fail 整改，A-006 conditional 收尾，A-008 independent `pass` 0 required）。 |
| GOAL-004-r3-entitlement-validation-telegram | R3 权益核验/消耗 + 可选 Telegram 注册 | **done** | 3/3 | GOAL-001-digital-offer-entitlement | 2026-09-05 关门：C1/C2 实施（Check/Consume/Telegram/查询桶）+ C3 审计（A-002 independent `pass` 0 required；A-003 处置 2 项 recommended）。 |
| GOAL-005-r4-evidence-closeout | R4 证据矩阵/边界核账/关门审计 | **active** | 1/2 | GOAL-001-digital-offer-entitlement | 曾两次关门（2026-09-05）。**第 3 轮撤回**：A-012 independent `conditional` 2 required（F-007/F-008 证据缺口）→ D-004 fixed ×2 + A-013 closed ×2（E-008：迁移失败/reopen 测试 SQLite+PG `Test(Migrate\|PostgresMigrate)MidApplyFailureNoResidueThenReopen`；Telegram-enabled 组合根 `TestDigitalOfferTelegramCompositionRoot` + seam 透传 + DELETE 负向）。C2 重开中，待 focused independent closure 复审 `pass` 0 required 后重新关门（2/2）。 |
