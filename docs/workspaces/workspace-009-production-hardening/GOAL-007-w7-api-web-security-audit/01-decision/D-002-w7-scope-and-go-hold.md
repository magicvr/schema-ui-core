---
id: D-002
goal: GOAL-007-w7-api-web-security-audit
title: 用户确认 S2 修复范围与 I-002 go 暂挂裁决
date: 2026-08-19
status: accepted
parent: GOAL-001-production-hardening
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
---

# D-002 · 用户确认 S2 修复范围与 I-002 go 暂挂裁决

### 触发

2026-08-19 目标轮次通过 P-004 用户裁决点向用户确认：

- S2：GOAL-007 修复范围是否整单采纳 A-001 的 required findings（F-001～F-012）。
- I-002：2 条 high required（F-001 MFA 存储错误 fail-open、F-002 mfa-reset 越权）期间是否暂挂 VP-008 `go` 消费有效性宣称。

### 用户书面裁决

1. **S2 修复范围**：整单采纳 A-001 F-001～F-012 为本波 required 修复范围。
2. **I-002**：在 F-001/F-002 两条 high required 闭合前，**不对外宣称 VP-008 `go` 消费有效性**；闭合后恢复宣称前应复核。

### 决定

1. 本波实施范围 = A-001 必改项汇总第 1～12 条（F-001..F-012）。不逐条做 residual/overruled。
2. `00-meta.md` 成功标准 S2 标记为完成；S3 开始实施；S4 待 required=0 + cross 审计后关门。
3. I-002 在 01-decision 信息表与 meta 中状态改为 verified/closed：暂挂对外宣称，证据 = 本 D-002 + 实施后 F-001/F-002 闭合证据。
4. 审计模式按 D-001 §4：实施后 self；关门前 independent（workspace.md 默认 grok build · grok-4.6 · high · `/audit`）→ cross 门禁。

### 为什么

- P-004 要求修复范围与 go 暂挂不能由编排器静默裁决；用户已书面确认，可放行 S3。
- fail-closed 语义：存储异常时应拒绝登录，不能让已知密码攻击者绕过 2FA；mfa-reset 必须镜像 users 资源的 admin 目标边界。
- 对外 go 宣称属于共享基架信任声明，2 条 high 未闭合期间暂挂是保守且可追溯的处置。

### 未选方案

- 只修 2 条 high、其余 required 延后：用户未选，范围不成立。
- 不暂挂 go 宣称：用户未选，不据此执行。
- 静默按 recommended 忽略 F-013～F-016：A-001 标注 recommended，本波不纳入 S3/S4 required 闭合范围，但实施中若顺手修复仍可在 E 记录。