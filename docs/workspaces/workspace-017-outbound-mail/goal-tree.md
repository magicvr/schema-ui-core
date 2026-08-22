---
title: 目标树 · workspace-017-outbound-mail
status: active
created: 2026-08-22
updated: 2026-08-22
parent: null
version: 0.2.0
workspace_id: workspace-017-outbound-mail
---

# 目标树 · 出站邮件

> 工作区：`workspace-017-outbound-mail`
> canonical：`docs/workspaces/workspace-017-outbound-mail/`
> Root：`GOAL-001-outbound-mail`（**active** · 1/4）
> primary_plan：`VP-017-outbound-mail`（`active`）

## 树

```text
GOAL-001-outbound-mail [active 1/4]     · 出站邮件（SMTP 发送端口）
└── GOAL-002-port-contract-freeze [done 3/3]   · R1 发送端口与合同冻结（A-001 self pass；I-001/I-002 verified）
```

R1 已冻结（Root D-002 + GOAL-002 D-001，kernel 端口代码落地）。下一阶段 R2 前须关闭 I-003/I-004。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-outbound-mail | 出站邮件（SMTP 发送端口） | null | active | 1/4 | 2026-08-22 |
| GOAL-002-port-contract-freeze | R1 发送端口与合同冻结 | GOAL-001-outbound-mail | done | 3/3 | 2026-08-22 |
