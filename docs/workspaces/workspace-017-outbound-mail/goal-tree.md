---
title: 目标树 · workspace-017-outbound-mail
status: active
created: 2026-08-22
updated: 2026-08-22
parent: null
version: 0.3.0
workspace_id: workspace-017-outbound-mail
---

# 目标树 · 出站邮件

> 工作区：`workspace-017-outbound-mail`
> canonical：`docs/workspaces/workspace-017-outbound-mail/`
> Root：`GOAL-001-outbound-mail`（**active** · 2/4）
> primary_plan：`VP-017-outbound-mail`（`active`）

## 树

```text
GOAL-001-outbound-mail [active 2/4]     · 出站邮件（SMTP 发送端口）
├── GOAL-002-port-contract-freeze [done 3/3]   · R1 发送端口与合同冻结（A-001 self pass；I-001/I-002 verified）
└── GOAL-003-smtp-dial-config [done 3/3]       · R2 SMTP 接入与配置面（A-001 self pass；隐式 TLS 465 唯一路径；I-003/I-004 verified）
```

R1/R2 已完成（合同冻结 + 适配器/配置面，测试绿）。下一阶段 R3：capture sink 落地 + composition 接线 + 公共面去客户端类型 sweep。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-outbound-mail | 出站邮件（SMTP 发送端口） | null | active | 2/4 | 2026-08-22 |
| GOAL-002-port-contract-freeze | R1 发送端口与合同冻结 | GOAL-001-outbound-mail | done | 3/3 | 2026-08-22 |
| GOAL-003-smtp-dial-config | R2 SMTP 接入与配置面 | GOAL-001-outbound-mail | done | 3/3 | 2026-08-22 |
