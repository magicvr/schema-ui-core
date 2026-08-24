---
title: 目标树 · workspace-017-outbound-mail
status: active
created: 2026-08-22
updated: 2026-08-24
parent: null
version: 0.8.0
workspace_id: workspace-017-outbound-mail
---

# 目标树 · 出站邮件

> 工作区：`workspace-017-outbound-mail`
> canonical：`docs/workspaces/workspace-017-outbound-mail/`
> Root：`GOAL-001-outbound-mail`（**active · 5/8**；用户 2026-08-24 否决组合层关门；R1～R4 实施史不回退）
> primary_plan：`VP-017-outbound-mail`（**`active`** · v0.4.0；VRev-041）

## 树

```text
GOAL-001-outbound-mail [active 6/8]       · 出站邮件（渠道供应商模型）
├── GOAL-002-port-contract-freeze [done 3/3]        · R1 发送端口与合同冻结（历史；不回退）
├── GOAL-003-smtp-dial-config [done 3/3]            · R2 SMTP 接入与配置面（历史；不回退）
├── GOAL-004-default-sink-surface-sweep [done 3/3]  · R3 默认 sink 落地与公共面 sweep（历史；不回退）
├── GOAL-005-r4-readyz-evidence [done 3/3]          · R4 显式路径证据与 readyz 扩依赖（历史；不回退）
├── GOAL-006-channel-provider-contract [done 3/3]   · R5 渠道供应商合同冻结（D-002；I-010/I-011 verified）
├── GOAL-007-mock-resend-delivery [done 4/4]        · R6 mock 站内出站记录与 Resend 渠道落地（A-001 self pass）
└── GOAL-008-mail-admin-surface [active 0/4]        · R7 设置「邮件」tab：配置 / 热切换 / 试发（消费 D-007）
```

历史 Root 关门（A-001 self pass + A-002 independent pass）原文保留，**已由用户否决其作为现行 `done` 的效力**。开放 required Goal finding = 0。信息项全部闭合（I-009 由 Root D-007 关闭；GOAL-007 F-001 residual 经用户裁决 accepted-residual）。下一阶段 = R8（证据 + readyz），R7 关门后按 P-001 开设。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-outbound-mail | 出站邮件（渠道供应商模型） | null | active | 6/8 | 2026-08-24 |
| GOAL-002-port-contract-freeze | R1 发送端口与合同冻结 | GOAL-001-outbound-mail | done | 3/3 | 2026-08-22 |
| GOAL-003-smtp-dial-config | R2 SMTP 接入与配置面 | GOAL-001-outbound-mail | done | 3/3 | 2026-08-22 |
| GOAL-004-default-sink-surface-sweep | R3 默认 sink 落地与公共面 sweep | GOAL-001-outbound-mail | done | 3/3 | 2026-08-22 |
| GOAL-005-r4-readyz-evidence | R4 显式路径证据与 readyz 扩依赖 | GOAL-001-outbound-mail | done | 3/3 | 2026-08-22 |
| GOAL-006-channel-provider-contract | R5 渠道供应商合同冻结 | GOAL-001-outbound-mail | done | 3/3 | 2026-08-24 |
| GOAL-007-mock-resend-delivery | R6 mock 站内出站记录与 Resend 渠道落地 | GOAL-001-outbound-mail | done | 4/4 | 2026-08-24 |
| GOAL-008-mail-admin-surface | R7 设置「邮件」tab：配置 / 热切换 / 试发 | GOAL-001-outbound-mail | active | 0/4 | 2026-08-24 |
