# goal-tree · workspace-031-digital-offer-entitlement

*自动同步工作区扁平目标树（树 + 状态表）。更新任一目标状态/进度后必须同步本文件。更新：2026-09-05（创建 GOAL-002 R1 子目标；C1 信息裁决关门；D-002 合同 draft）*

## 目标树

```text
GOAL-001-digital-offer-entitlement (数字 Offer 与权益 · active · 0/4)
└── GOAL-002-r1-contract-freeze (R1 合同冻结 · active · 1/3)
（R1 进行中：C1 关门；C2 合同 D-002 draft 待审计；R2～R4 未立项，按阶段渐进开设）
```

## 状态表

| id | title | status | progress | parent | notes |
|----|-------|--------|----------|--------|-------|
| GOAL-001-digital-offer-entitlement | 数字 Offer 与权益 | **active** | 0/4 | null | VP-031 lead Root。R1 进行中（GOAL-002 承载）：I-031-001～005 已 verified（用户裁决 + 默认冻结，GOAL-002 D-001）；合同 D-002 v0.1.0 draft 待 self+codex 审计；V-F119 已纳入合同 §8。 |
| GOAL-002-r1-contract-freeze | R1 合同冻结（Offer 字段/购买状态机/权益形态/命令清单/事务与限流边界） | **active** | 1/3 | GOAL-001-digital-offer-entitlement | C1 信息裁决关门（2026-09-05 用户书面）；C2 D-002 draft 落盘待审；C3 待 self A-001 + codex independent A-002（gpt 5.6 sol · medium）。 |
