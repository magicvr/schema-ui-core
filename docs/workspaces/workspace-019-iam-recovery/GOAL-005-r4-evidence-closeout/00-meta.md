---
id: GOAL-005-r4-evidence-closeout
title: R4 端到端证据与关门
status: done
parent: GOAL-001-iam-recovery
created: 2026-08-25
updated: 2026-08-25
version: 1.0.0
progress: 3/3
plan_refs:
  - VP-019-iam-recovery
primary_plan: VP-019-iam-recovery
serves_summary: 承接 Root R4：三条链（恢复/邀请/策略强制）经真实 HTTP + mock 渠道取信；无越界核对；关门审计后 Root 4/4。
---

# GOAL-005 · R4 端到端证据与关门

## 检查点（progress 来源）

| # | 检查点 | 证据 |
|---|--------|------|
| C1 | 恢复链 HTTP e2e：bind/verify → start → 渠道取码 → complete → 新密码登录 | **完成**：TestR4RecoveryChainOverHTTP（app 代码入账） |
| C2 | 邀请链 HTTP e2e：admin 建邀 → 渠道取链接 → accept → 受邀角色登录 + 一次性回放拒绝；策略强化经 HTTP 拦住弱创建 | **完成**：TestR4InviteChainOverHTTP / TestR4PolicyEnforcementOverHTTP |
| C3 | 关门审计：self + independent（grok build）开放 required = 0，无越界 | **完成**：A-001 independent conditional→F-001/F-002 fixed 归零；A-002 self pass |

`progress` = 1/3。
