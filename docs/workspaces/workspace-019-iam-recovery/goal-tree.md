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
└─ GOAL-002-iam-contract-freeze [active 2/3]      · R1 IAM 合同冻结（恢复 / 策略 / 邀请）
```

R1 进行中（2026-08-25，E-002/E-003）：I-001～I-009 全部经用户裁决 **verified**（Root D-002 + GOAL-002 D-001），合同条款 §1～§5 已落盘；剩 C3 self 关门审计后 GOAL-002 关门、Root R1 记完成。R2 自助恢复全链依赖 R1 冻结。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-iam-recovery | IAM：密码策略 / 邀请入职 / 自助恢复状态机 | null | active | 0/4 | 2026-08-25 |
| GOAL-002-iam-contract-freeze | R1 IAM 合同冻结（恢复 / 策略 / 邀请） | GOAL-001-iam-recovery | active | 2/3 | 2026-08-25 |