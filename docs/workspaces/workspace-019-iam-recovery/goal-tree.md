---
title: 目标树 · workspace-019-iam-recovery
status: active
created: 2026-08-25
updated: 2026-08-25
parent: null
version: 0.2.0
workspace_id: workspace-019-iam-recovery
---

# 目标树 · IAM（密码策略 / 邀请入职 / 自助恢复状态机）

> 工作区：`workspace-019-iam-recovery`
> canonical：`docs/workspaces/workspace-019-iam-recovery/`
> Root：`GOAL-001-iam-recovery`（**active** · 2/4；R1/R2 关门 2026-08-25）
> primary_plan：`VP-019-iam-recovery`（**active** · VRev-043 independent `pass`）

## 树

```text
GOAL-001-iam-recovery [active 2/4]                · IAM：密码策略 / 邀请入职 / 自助恢复状态机
├─ GOAL-002-iam-contract-freeze [done 3/3]        · R1 IAM 合同冻结（恢复 / 策略 / 邀请）
├─ GOAL-003-r2-self-recovery-flow [done 5/5]      · R2 自助恢复全链（后端 + Web）
└─ GOAL-004-r3-policy-and-invites [active 2/5]    · R3 密码策略 + 邀请入职
```

R1 关门（GOAL-002 · A-001 self pass）：I-001～I-009 全部经用户裁决 **verified**，合同条款 §1～§5 落盘。R2 关门（GOAL-003 · A-001 independent conditional→F-001～F-004 fixed 归零 + A-002 self pass）：迁移 0056 + start/complete 公开面 + MFA 第二因子门 + Web 两步恢复流落地。下一步 R3 密码策略 + 邀请入职立项（GOAL-004），实施切片继续走 independent（grok build · grok-4.6 · high）。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-iam-recovery | IAM：密码策略 / 邀请入职 / 自助恢复状态机 | null | active | 2/4 | 2026-08-25 |
| GOAL-002-iam-contract-freeze | R1 IAM 合同冻结（恢复 / 策略 / 邀请） | GOAL-001-iam-recovery | done | 3/3 | 2026-08-25 |
| GOAL-003-r2-self-recovery-flow | R2 自助恢复全链（后端 + Web） | GOAL-001-iam-recovery | done | 5/5 | 2026-08-25 |
| GOAL-004-r3-policy-and-invites | R3 密码策略 + 邀请入职 | GOAL-001-iam-recovery | active | 2/5 | 2026-08-25 |
