# goal-tree · workspace-029-wallet-prepaid-instrument

*自动同步工作区扁平目标树（树 + 状态表）。更新任一目标状态/进度后必须同步本文件。更新：2026-09-02*

## 目标树

```text
GOAL-001-wallet-prepaid-instrument (钱包预付资金凭证与外部主体接缝 · done · 4/4)
├── GOAL-002-subject-module-and-wallet-integration (外部主体接缝与钱包预付凭证入金集成 · done · 4/4)
├── GOAL-003-admin-voucher-surface-and-export (管理后台预付凭证批次管理、导出与操作审计 · done · 4/4)
├── GOAL-004-evidence-closure-and-closeout (工作区29 证据矩阵闭环、越界核账与根目标关门 · done · 4/4)
（R1 合同冻结 [done] → R2 主体接缝+账本入金 [GOAL-002 done] → R3 Admin 批次面+导出 [GOAL-003 done] → R4 证据与关门 [GOAL-004 done]）
```

## 状态表

| id | title | status | progress | parent | notes |
|----|-------|--------|----------|--------|-------|
| GOAL-001-wallet-prepaid-instrument | 钱包预付资金凭证与外部主体接缝 | **done** | 4/4 | null | 2026-09-02 开区并关门：VP-029 判据 #1~#7 全部实证闭环；R1 合同冻结（D-002）、R2 主体接缝与预付入金（GOAL-002）、R3 Admin 批次面与协议页面（GOAL-003）、R4 证据矩阵与红线核账（GOAL-004）全量达成；审计链 A-001 pass → A-002 conditional（3 required）→ A-003/A-004 fixed 复审闭合 → A-005 conditional（F-001 主体门禁 OR 回归）→ A-006/A-007 fixed 闭合 → A-008 recommended 全处置（F-002 声明化导出 / F-003 PG e2e / F-004 0065 批次注册表 / F-005 expiresAt 秒范围）——open required = 0，维持 done；E-007 用户反馈修复：F-002 声明载体改表单 props（pin 兼容）、voucher 导航图标与排序（Wallet 下一位） |
| GOAL-002-subject-module-and-wallet-integration | 外部主体接缝与钱包预付凭证入金集成 | **done** | 4/4 | GOAL-001-wallet-prepaid-instrument | R2 实施子目标：主体表/凭证表迁移 0064、主体接缝、核销原子入金与并发防双花（A-001 pass + A-002 conditional → A-003 fixed 关门） |
| GOAL-003-admin-voucher-surface-and-export | 管理后台预付凭证批次管理、导出与操作审计 | **done** | 4/4 | GOAL-001-wallet-prepaid-instrument | R3 实施子目标：批次生成/导出/作废/查询 HTTP 端点、wallet.voucher.issue 权限键、协议驱动页面与操作审计（A-001 pass + A-002 conditional → A-003 fixed 关门） |
| GOAL-004-evidence-closure-and-closeout | 工作区29 证据矩阵闭环、越界核账与根目标关门 | **done** | 4/4 | GOAL-001-wallet-prepaid-instrument | R4 实施子目标：VP-029 判据 #1~#7 证据矩阵核账、边界核对与双腿关门审计（A-001 self + A-002 grok independent conditional → A-003 fixed 闭合关门） |
