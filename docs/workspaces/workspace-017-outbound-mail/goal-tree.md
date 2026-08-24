---
title: 目标树 · workspace-017-outbound-mail
status: active
created: 2026-08-22
updated: 2026-08-24
parent: null
version: 1.0.0
workspace_id: workspace-017-outbound-mail
---

# 目标树 · 出站邮件

> 工作区：`workspace-017-outbound-mail`
> canonical：`docs/workspaces/workspace-017-outbound-mail/`
> Root：`GOAL-001-outbound-mail`（**done · 8/8**；2026-08-24 用户否决旧关门并升级分母后，R5～R8 重开交付完毕，现行分母再关门审计 A-003/A-004 pass 放行）
> primary_plan：`VP-017-outbound-mail`（v0.4.0 → v0.5.0 `closed`；VRev-041 / VRev-042）

## 树

```text
GOAL-001-outbound-mail [done 8/8]         · 出站邮件（渠道供应商模型）
├── GOAL-002-port-contract-freeze [done 3/3]        · R1 发送端口与合同冻结（历史；不回退）
├── GOAL-003-smtp-dial-config [done 3/3]            · R2 SMTP 接入与配置面（历史；不回退）
├── GOAL-004-default-sink-surface-sweep [done 3/3]  · R3 默认 sink 落地与公共面 sweep（历史；不回退）
├── GOAL-005-r4-readyz-evidence [done 3/3]          · R4 显式路径证据与 readyz 扩依赖（历史；不回退）
├── GOAL-006-channel-provider-contract [done 3/3]   · R5 渠道供应商合同冻结（D-002；I-010/I-011 verified）
├── GOAL-007-mock-resend-delivery [done 4/4]        · R6 mock 站内出站记录与 Resend 渠道落地（A-001 self pass）
├── GOAL-008-mail-admin-surface [done 4/4]          · R7 设置「邮件」tab：配置 / 热切换 / 试发（A-001 self pass）
└── GOAL-009-r8-evidence-readyz [done 4/4]          · R8 生产渠道探针与现行分母关门证据（A-001 self pass；live 投递 PASS）
```

历史 Root 关门（A-001 self pass + A-002 independent pass）原文保留，**已由用户否决其作为当时分母的效力**。本次再关门对照**现行分母**新开审计：A-003（self 关门向）+ A-004（independent · 子代理交叉核对）均 pass；independent F-001/F-002（台账现势性）已随本提交 fixed。开放 required Goal finding = 0。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-outbound-mail | 出站邮件（渠道供应商模型） | null | done | 8/8 | 2026-08-24 |
| GOAL-002-port-contract-freeze | R1 发送端口与合同冻结 | GOAL-001-outbound-mail | done | 3/3 | 2026-08-22 |
| GOAL-003-smtp-dial-config | R2 SMTP 接入与配置面 | GOAL-001-outbound-mail | done | 3/3 | 2026-08-22 |
| GOAL-004-default-sink-surface-sweep | R3 默认 sink 落地与公共面 sweep | GOAL-001-outbound-mail | done | 3/3 | 2026-08-22 |
| GOAL-005-r4-readyz-evidence | R4 显式路径证据与 readyz 扩依赖 | GOAL-001-outbound-mail | done | 3/3 | 2026-08-22 |
| GOAL-006-channel-provider-contract | R5 渠道供应商合同冻结 | GOAL-001-outbound-mail | done | 3/3 | 2026-08-24 |
| GOAL-007-mock-resend-delivery | R6 mock 站内出站记录与 Resend 渠道落地 | GOAL-001-outbound-mail | done | 4/4 | 2026-08-24 |
| GOAL-008-mail-admin-surface | R7 设置「邮件」tab：配置 / 热切换 / 试发 | GOAL-001-outbound-mail | done | 4/4 | 2026-08-24 |
| GOAL-009-r8-evidence-readyz | R8 生产渠道探针与现行分母关门证据 | GOAL-001-outbound-mail | done | 4/4 | 2026-08-24 |
