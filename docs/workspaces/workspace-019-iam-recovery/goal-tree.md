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
> Root：`GOAL-001-iam-recovery`（**done** · 4/4；2026-08-25 全链关门）
> primary_plan：`VP-019-iam-recovery`（**active** · VRev-043 independent `pass`）

## 树

```text
GOAL-001-iam-recovery [active 2/4]                · IAM：密码策略 / 邀请入职 / 自助恢复状态机
├─ GOAL-002-iam-contract-freeze [done 3/3]        · R1 IAM 合同冻结（恢复 / 策略 / 邀请）
├─ GOAL-003-r2-self-recovery-flow [done 5/5]      · R2 自助恢复全链（后端 + Web）
└─ GOAL-004-r3-policy-and-invites [active 3/5]    · R3 密码策略 + 邀请入职
```

R1 关门（GOAL-002）：I-001～I-009 全闭 + 合同条款落盘。R2 关门（GOAL-003）：迁移 0056 + 恢复全链 + MFA 第二因子门。R3 关门（GOAL-004 · A-001 independent conditional→F-001～F-004 fixed 归零 + A-002 self pass）：迁移 0057–0059 + 策略四口强制 + 邀请全链 + Web 面（含新建用户表单角色选择）。**根目标已收官（done 4/4 · 2026-08-25）**：R1～R4 全链关门，开放 required = 0，VP-019 三件交付完成。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-iam-recovery | IAM：密码策略 / 邀请入职 / 自助恢复状态机 | null | done | 4/4 | 2026-08-25 |
| GOAL-002-iam-contract-freeze | R1 IAM 合同冻结（恢复 / 策略 / 邀请） | GOAL-001-iam-recovery | done | 3/3 | 2026-08-25 |
| GOAL-003-r2-self-recovery-flow | R2 自助恢复全链（后端 + Web） | GOAL-001-iam-recovery | done | 5/5 | 2026-08-25 |
| GOAL-004-r3-policy-and-invites | R3 密码策略 + 邀请入职 | GOAL-001-iam-recovery | done | 5/5 | 2026-08-25 |
| GOAL-005-r4-evidence-closeout | R4 端到端证据与关门 | GOAL-001-iam-recovery | done | 3/3 | 2026-08-25 |
