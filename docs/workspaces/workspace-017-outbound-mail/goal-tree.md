---
title: 目标树 · workspace-017-outbound-mail
status: active
created: 2026-08-22
updated: 2026-08-24
parent: null
version: 0.7.0
workspace_id: workspace-017-outbound-mail
---

# 目标树 · 出站邮件

> 工作区：`workspace-017-outbound-mail`
> canonical：`docs/workspaces/workspace-017-outbound-mail/`
> Root：`GOAL-001-outbound-mail`（**active · 4/8**；用户 2026-08-24 否决组合层关门；R1～R4 实施史不回退）
> primary_plan：`VP-017-outbound-mail`（**`active`** · v0.4.0；VRev-041）

## 树

```text
GOAL-001-outbound-mail [active 4/8]       · 出站邮件（渠道供应商模型）
├── GOAL-002-port-contract-freeze [done 3/3]        · R1 发送端口与合同冻结（历史；不回退）
├── GOAL-003-smtp-dial-config [done 3/3]            · R2 SMTP 接入与配置面（历史；不回退）
├── GOAL-004-default-sink-surface-sweep [done 3/3]  · R3 默认 sink 落地与公共面 sweep（历史；不回退）
├── GOAL-005-r4-readyz-evidence [done 3/3]          · R4 显式路径证据与 readyz 扩依赖（历史；不回退）
└── GOAL-006-channel-provider-contract [active 0/3] · R5 渠道供应商合同冻结
```

历史 Root 关门（A-001 self pass + A-002 independent pass）原文保留，**已由用户否决其作为现行 `done` 的效力**。开放 required Goal finding = 0。R6～R8 子目标按 P-001 尚未创建。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-outbound-mail | 出站邮件（渠道供应商模型） | null | active | 4/8 | 2026-08-24 |
| GOAL-002-port-contract-freeze | R1 发送端口与合同冻结 | GOAL-001-outbound-mail | done | 3/3 | 2026-08-22 |
| GOAL-003-smtp-dial-config | R2 SMTP 接入与配置面 | GOAL-001-outbound-mail | done | 3/3 | 2026-08-22 |
| GOAL-004-default-sink-surface-sweep | R3 默认 sink 落地与公共面 sweep | GOAL-001-outbound-mail | done | 3/3 | 2026-08-22 |
| GOAL-005-r4-readyz-evidence | R4 显式路径证据与 readyz 扩依赖 | GOAL-001-outbound-mail | done | 3/3 | 2026-08-22 |
| GOAL-006-channel-provider-contract | R5 渠道供应商合同冻结 | GOAL-001-outbound-mail | active | 0/3 | 2026-08-24 |
