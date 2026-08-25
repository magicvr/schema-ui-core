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
> Root：`GOAL-001-iam-recovery`（**active** · 0/4；2026-08-25 开区）
> primary_plan：`VP-019-iam-recovery`（**active** · VRev-043 independent `pass`）

## 树

```text
GOAL-001-iam-recovery [active 0/4]                · IAM：密码策略 / 邀请入职 / 自助恢复状态机
└─ GOAL-002-iam-contract-freeze [active 1/3]      · R1 IAM 合同冻结（恢复 / 策略 / 邀请）
```

R1 已立项进行中（2026-08-25，E-002/E-001）：I-001/I-002/I-009 经用户裁决 **verified**（Root D-002），C1 满；余下 I-003/I-004/I-005（required）与 I-007/I-008 投影确认在 GOAL-002 内继续，随后落合同条款（C2）与 self 审计（C3）。R2 自助恢复全链依赖 R1 冻结。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-iam-recovery | IAM：密码策略 / 邀请入职 / 自助恢复状态机 | null | active | 0/4 | 2026-08-25 |
| GOAL-002-iam-contract-freeze | R1 IAM 合同冻结（恢复 / 策略 / 邀请） | GOAL-001-iam-recovery | active | 1/3 | 2026-08-25 |