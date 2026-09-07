# goal-tree · workspace-029-wallet-prepaid-instrument

*自动同步工作区扁平目标树（树 + 状态表）。更新任一目标状态/进度后必须同步本文件。更新：2026-09-02*

## 目标树

```text
GOAL-001-wallet-prepaid-instrument (钱包预付资金凭证与外部主体接缝 · done · 5/5)
├── GOAL-002-subject-module-and-wallet-integration (外部主体接缝与钱包预付凭证入金集成 · done · 4/4)
├── GOAL-003-admin-voucher-surface-and-export (管理后台预付凭证批次管理、导出与操作审计 · done · 4/4)
├── GOAL-004-evidence-closure-and-closeout (工作区29 证据矩阵闭环、越界核账与根目标关门 · done · 4/4)
└── GOAL-005-my-wallet-voucher-redeem (我的钱包预付凭证自助核销入口 · done · 4/4)
（R1 合同冻结 [done] → R2 主体接缝+账本入金 [GOAL-002 done] → R3 Admin 批次面+导出 [GOAL-003 done] → R4 证据与关门 [GOAL-004 done] → R5 Admin 自助核销 HTTP + 我的钱包入口 [GOAL-005 done]；Root A-009 关门自审 pass · 工作区 029 结项）
```

## 状态表

| id | title | status | progress | parent | notes |
|----|-------|--------|----------|--------|-------|
| GOAL-001-wallet-prepaid-instrument | 钱包预付资金凭证与外部主体接缝 | **done** | 5/5 | null | 2026-09-02 开区，R1～R5 全量交付（判据 #1～#10 全量 verified · 9 项信息需求 closed）。A-009 关门自审 pass，开放 required = 0，根目标与工作区 029 结项关门 |
| GOAL-002-subject-module-and-wallet-integration | 外部主体接缝与钱包预付凭证入金集成 | **done** | 4/4 | GOAL-001-wallet-prepaid-instrument | R2 实施子目标：主体表/凭证表迁移 0064、主体接缝、核销原子入金与并发防双花（A-001 pass + A-002 conditional → A-003 fixed 关门） |
| GOAL-003-admin-voucher-surface-and-export | 管理后台预付凭证批次管理、导出与操作审计 | **done** | 4/4 | GOAL-001-wallet-prepaid-instrument | R3 实施子目标：批次生成/导出/作废/查询 HTTP 端点、wallet.voucher.issue 权限键、协议驱动页面与操作审计（A-001 pass + A-002 conditional → A-003 fixed 关门） |
| GOAL-004-evidence-closure-and-closeout | 工作区29 证据矩阵闭环、越界核账与根目标关门 | **done** | 4/4 | GOAL-001-wallet-prepaid-instrument | R4 实施子目标：VP-029 当时分母判据 #1～#7 证据矩阵核账、边界核对与双腿关门审计（A-001 self + A-002 grok independent conditional → A-003 fixed 闭合关门） |
| GOAL-005-my-wallet-voucher-redeem | 我的钱包预付凭证自助核销入口 | **done** | 4/4 | GOAL-001-wallet-prepaid-instrument | R5：S1 D-002；S2 HTTP+页面；S3 回归；S4 A-001 independent pass + A-003 self pass；A-001 F-001～F-005 fixed；用户确认 D-004 关门 |
