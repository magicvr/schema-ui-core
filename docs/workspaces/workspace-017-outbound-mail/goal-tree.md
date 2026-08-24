---
title: 目标树 · workspace-017-outbound-mail
status: done
created: 2026-08-22
updated: 2026-08-22
parent: null
version: 0.6.0
workspace_id: workspace-017-outbound-mail
---

# 目标树 · 出站邮件

> 工作区：`workspace-017-outbound-mail`
> canonical：`docs/workspaces/workspace-017-outbound-mail/`
> Root：`GOAL-001-outbound-mail`（**done · 4/4**；A-001 self pass + A-002 independent pass，recommended 全响应）
> primary_plan：`VP-017-outbound-mail`（**`closed`** · 2026-08-24 组合层有界关门；VRev-039 / VR-042）

## 树

```text
GOAL-001-outbound-mail [done 4/4]       · 出站邮件（SMTP 发送端口）
├── GOAL-002-port-contract-freeze [done 3/3]        · R1 发送端口与合同冻结（A-001 self pass；I-001/I-002 verified）
├── GOAL-003-smtp-dial-config [done 3/3]            · R2 SMTP 接入与配置面（A-001 self pass；隐式 TLS 465 唯一路径；I-003/I-004 verified）
├── GOAL-004-default-sink-surface-sweep [done 3/3]  · R3 默认 sink 落地与公共面 sweep（A-001 self pass；sweep 零泄漏）
└── GOAL-005-r4-readyz-evidence [done 3/3]          · R4 显式路径证据与 readyz 扩依赖（A-001 self pass；I-005/I-006 verified）
```

Root 关门：R1～R4 全部完成；关门审计 = A-001 self `pass` + A-002 independent（grok build）`pass`，三条 recommended 均已响应（F-001/F-002 `fixed`、F-003 delegated `/vision`）。开放 required finding = 0。VP-017 关门记录由 `/vision` 收尾。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-outbound-mail | 出站邮件（SMTP 发送端口） | null | done | 4/4 | 2026-08-22 |
| GOAL-002-port-contract-freeze | R1 发送端口与合同冻结 | GOAL-001-outbound-mail | done | 3/3 | 2026-08-22 |
| GOAL-003-smtp-dial-config | R2 SMTP 接入与配置面 | GOAL-001-outbound-mail | done | 3/3 | 2026-08-22 |
| GOAL-004-default-sink-surface-sweep | R3 默认 sink 落地与公共面 sweep | GOAL-001-outbound-mail | done | 3/3 | 2026-08-22 |
| GOAL-005-r4-readyz-evidence | R4 显式路径证据与 readyz 扩依赖 | GOAL-001-outbound-mail | done | 3/3 | 2026-08-22 |
