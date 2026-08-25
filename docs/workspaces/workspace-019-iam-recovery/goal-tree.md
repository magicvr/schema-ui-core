---
title: 目标树 · workspace-019-iam-recovery
status: active
created: 2026-08-25
updated: 2026-08-25
parent: null
version: 0.1.0
workspace_id: workspace-019-iam-recovery
---

# 目标树 · IAM（密码策略 / 邀请入职 / 自助恢复状态机）

> 工作区：`workspace-019-iam-recovery`
> canonical：`docs/workspaces/workspace-019-iam-recovery/`
> Root：`GOAL-001-iam-recovery`（**active** · 1/4；R1 合同冻结关门 2026-08-25）
> primary_plan：`VP-019-iam-recovery`（**active** · VRev-043 independent `pass`）

## 树

```text
GOAL-001-iam-recovery [active 1/4]                · IAM：密码策略 / 邀请入职 / 自助恢复状态机
└─ GOAL-002-iam-contract-freeze [done 3/3]        · R1 IAM 合同冻结（恢复 / 策略 / 邀请）
```

R1 关门（2026-08-25，E-003）：I-001～I-009 全部经用户裁决 **verified**（Root D-002 + GOAL-002 D-001），合同条款 §1～§5 落盘，A-001 self `pass`（0 required）。下一步 R2 自助恢复全链立项（GOAL-003），迁移与实施切片按既定模式走 independent（grok build · grok-4.6 · high）。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-iam-recovery | IAM：密码策略 / 邀请入职 / 自助恢复状态机 | null | active | 1/4 | 2026-08-25 |
| GOAL-002-iam-contract-freeze | R1 IAM 合同冻结（恢复 / 策略 / 邀请） | GOAL-001-iam-recovery | done | 3/3 | 2026-08-25 |