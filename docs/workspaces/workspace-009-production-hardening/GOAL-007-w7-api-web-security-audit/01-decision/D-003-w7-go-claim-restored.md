---
id: D-003
goal: GOAL-007-w7-api-web-security-audit
title: VP-008 go 宣称恢复：F-001/F-002 闭合经双 independent 复核
date: 2026-08-19
status: accepted
parent: GOAL-001-production-hardening
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
---

# D-003 · VP-008 go 宣称恢复

### 触发

D-002 规定：在 F-001/F-002 两条 high required 闭合前，不对外宣称 VP-008 `go` 消费有效性；闭合后恢复宣称前应复核。

经独立复核（A-004 grok-4.6 + A-005 claude-sonnet-4），确认 F-001（MFA fail-closed）和 F-002（MFA admin reset boundary）在现行代码中已 genuine fixed。A-001 全部 12/12 required 可核对闭合。

### 决定

1. **VP-008 `go` 消费有效性从暂挂恢复为有效。**
2. I-002 从「暂挂」更新为「已恢复」；证据 = A-004 + A-005 + D-002（原暂挂裁决）。
3. VP-008 自身 `status: closed` 不变；go 消费有效性规则不变；后续业务 VP 激活前仍需完成消费前 freshness review（VP-008 §go 消费有效性）。
4. F-006（captcha generate limiter）已不再构成披露缺口（A-004/A-005 确认），不单独作为 go 暂挂依据。

### 为什么

- F-001（MFA fail-closed）和 F-002（MFA admin reset boundary）是 D-002 指定的 go 暂挂条件。两条均已 genuine fixed 并经双 independent 复核（A-004 + A-005），满足 D-002 的「闭合后恢复宣称前应复核」条款。
- A-001 全部 12/12 required 均为 `fixed`（A-004 pass + A-005 pass），开放 required = 0。
- 恢复 go 宣称不改变 VP-008 的 `closed` 状态，也不改变其消费有效性规则（freshness review、失效触发等）。

### 未选方案

- 维持暂挂直至所有 recommended 闭合：用户未选；recommended 不阻断 required 闭合门禁。
- 仅依赖 A-004 单 independent 复核：用户指令要求独立复核（A-005），形成双 independent cross 门禁。