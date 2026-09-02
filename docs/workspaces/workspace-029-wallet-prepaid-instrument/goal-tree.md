# goal-tree · workspace-029-wallet-prepaid-instrument

*自动同步工作区扁平目标树（树 + 状态表）。更新任一目标状态/进度后必须同步本文件。更新：2026-09-02*

## 目标树

```text
GOAL-001-wallet-prepaid-instrument (钱包预付资金凭证与外部主体接缝 · active · 3/4)
├── GOAL-002-subject-module-and-wallet-integration (外部主体接缝与钱包预付凭证入金集成 · done · 4/4)
├── GOAL-003-admin-voucher-surface-and-export (管理后台预付凭证批次管理、导出与操作审计 · done · 4/4)
（R1 合同冻结 [done] → R2 主体接缝+账本入金 [GOAL-002 done] → R3 Admin 批次面+导出 [GOAL-003 done] → R4 证据与关门）
```

## 状态表

| id | title | status | progress | parent | notes |
|----|-------|--------|----------|--------|-------|
| GOAL-001-wallet-prepaid-instrument | 钱包预付资金凭证与外部主体接缝 | **active** | 3/4 | null | 2026-09-02 开区：VP-029 v0.2.0 active · VRev-066 independent `pass`；R1 合同冻结已完成（D-002）；R2 主体接缝与预付凭证入金已完成（GOAL-002 done）；R3 Admin 批次面已完成（GOAL-003 done） |
| GOAL-002-subject-module-and-wallet-integration | 外部主体接缝与钱包预付凭证入金集成 | **done** | 4/4 | GOAL-001-wallet-prepaid-instrument | R2 实施子目标：主体表/凭证表迁移 0064、主体接缝、核销原子入金与并发防双花（A-001 pass + A-002 conditional → A-003 fixed 关门） |
| GOAL-003-admin-voucher-surface-and-export | 管理后台预付凭证批次管理、导出与操作审计 | **done** | 4/4 | GOAL-001-wallet-prepaid-instrument | R3 实施子目标：批次生成/导出/作废/查询 HTTP 端点、wallet.voucher.issue 权限键与操作审计（A-001 pass + A-002 conditional → A-003 fixed 关门） |
