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
> Root：`GOAL-001-account-email-identity`（**active · 2/4**；D-003 已解冻；R1/R2 已关门）
> primary_plan：`VP-018-account-email-identity`（`active` · 已解冻）

## 树

```text
GOAL-001-account-email-identity [active 2/4]        · 账号邮箱身份（绑定与校验）
├── GOAL-002-identity-contract-freeze [done 3/3]    · R1 身份合同冻结
└── GOAL-003-dual-dialect-email-schema [done 4/4]   · R2 双方言 schema + 唯一性
```

开区 scaffold 保留。2026-08-24 解冻（D-003）后：R1 由 GOAL-002 承接并关门（身份合同七条款冻结 · A-001 self pass）；R2 由 GOAL-003 承接并关门（迁移 0054 双方言落地 · 六处黄金断言同步 · A-001 independent pass / grok build grok-4.6 high）。下一步 R3 绑定/校验流消费 `MailSender`：I-005 数值冻结为接入前门禁；承接清单 = PG 语义 harness（GOAL-003 A-001 F-001，可选）、email/email_status 配对不变量（F-002）、SQLite lower() ASCII 补偿的仓储归一（N-1）。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-account-email-identity | 账号邮箱身份（绑定与校验） | null | active | 2/4 | 2026-08-24 |
| GOAL-002-identity-contract-freeze | R1 身份合同冻结 | GOAL-001-account-email-identity | done | 3/3 | 2026-08-24 |
| GOAL-003-dual-dialect-email-schema | R2 双方言 schema + 唯一性 | GOAL-001-account-email-identity | done | 4/4 | 2026-08-24 |
