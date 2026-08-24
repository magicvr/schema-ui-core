---
title: 目标树 · workspace-018-account-email-identity
status: active
created: 2026-08-24
updated: 2026-08-24
parent: null
version: 0.2.0
workspace_id: workspace-018-account-email-identity
---

# 目标树 · 账号邮箱身份

> 工作区：`workspace-018-account-email-identity`
> canonical：`docs/workspaces/workspace-018-account-email-identity/`
> Root：`GOAL-001-account-email-identity`（**active · 1/4**；D-003 已解冻——VP-017 v0.5.0 按现行分母再关门；R1 已关门）
> primary_plan：`VP-018-account-email-identity`（`active` · 已解冻）

## 树

```text
GOAL-001-account-email-identity [active 1/4]        · 账号邮箱身份（绑定与校验）
└── GOAL-002-identity-contract-freeze [done 3/3]    · R1 身份合同冻结
```

开区 scaffold 保留。2026-08-24 解冻（D-003）后 R1 由 GOAL-002 承接并关门：身份合同七条款冻结（验证码 · 绑定即占槽 · lower(email) 唯一），A-001 self pass。下一步 R2 双方言 schema（independent 审计前置）；I-005 数值冻结归 R3 接入前。`users` DDL 变更随 R2 实施。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-account-email-identity | 账号邮箱身份（绑定与校验） | null | active | 1/4 | 2026-08-24 |
| GOAL-002-identity-contract-freeze | R1 身份合同冻结 | GOAL-001-account-email-identity | done | 3/3 | 2026-08-24 |
