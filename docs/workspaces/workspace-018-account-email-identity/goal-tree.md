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
> Root：`GOAL-001-account-email-identity`（**active · 4/4**；四阶段全关；Root 关门审计 A-002 待 grok 恢复——E-009）
> primary_plan：`VP-018-account-email-identity`（`active` · 已解冻）

## 树

```text
GOAL-001-account-email-identity [active 4/4]        · 账号邮箱身份（绑定与校验）
├── GOAL-002-identity-contract-freeze [done 3/3]    · R1 身份合同冻结
├── GOAL-003-dual-dialect-email-schema [done 4/4]   · R2 双方言 schema + 唯一性
├── GOAL-004-r3-binding-flow [done 4/4]             · R3 绑定/校验消费 MailSender
└── GOAL-005-r4-evidence [done 3/3]                 · R4 证据与关门就绪
```

开区 scaffold 保留。2026-08-24 解冻（D-003）后连续关门：R1 = GOAL-002（身份合同七条款 · self pass）；R2 = GOAL-003（迁移 0054 双方言落地 · independent pass）；R3 = GOAL-004（迁移 0055 挑战表 + bind/verify/resend 消费 MailSender + I-006 代填 HTTP 链路 + 最小绑定卡 · A-001 conditional → F-001 fixed 后 required 归零）。剩 **R4 证据**：从 VP-017 当时默认渠道取信、唯一性 fail-closed 可核对、无 IAM 越界，关门审计后 Root 结项。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-account-email-identity | 账号邮箱身份（绑定与校验） | null | active | 4/4 | 2026-08-24 |
| GOAL-002-identity-contract-freeze | R1 身份合同冻结 | GOAL-001-account-email-identity | done | 3/3 | 2026-08-24 |
| GOAL-003-dual-dialect-email-schema | R2 双方言 schema + 唯一性 | GOAL-001-account-email-identity | done | 4/4 | 2026-08-24 |
| GOAL-004-r3-binding-flow | R3 绑定/校验消费 MailSender | GOAL-001-account-email-identity | done | 4/4 | 2026-08-24 |
| GOAL-005-r4-evidence | R4 证据与关门就绪 | GOAL-001-account-email-identity | done | 3/3 | 2026-08-24 |
