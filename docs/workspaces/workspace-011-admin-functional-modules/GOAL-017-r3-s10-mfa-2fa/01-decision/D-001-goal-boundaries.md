---
id: D-001
goal: GOAL-017-r3-s10-mfa-2fa
title: 立项边界：因子选型、模块身份与审计策略
date: 2026-08-15
status: accepted
parent: GOAL-017-r3-s10-mfa-2fa
created: 2026-08-15
updated: 2026-08-15
version: 1.0.0
---

# D-001 · 立项边界（S-10 MFA/2FA）

## 决策

1. **因子选型（I-001）**：**TOTP（RFC 6238）优先**；不引入短信/邮件通道（邮件依赖 B-09 模板基建，S1 复核）。恢复码策略（一次性、可吊销）S1 冻结。
2. **模块身份**：候选 admin.mfa（标准 Admin 模块，管理/自助面）+ core.auth 集成扩展（挑战端点与登录流程）；S1 冻结最终身份与 Descriptor 依赖（预期 core.auth-session / core.operationlog / core.navigation-capability）。
3. **审计策略**：MFA 属 **security 门禁** → S1 方案冻结与 S5 关门必须 **grok build independent**（用户书面偏好：grok-4.6 · reasoning high）；S2/S3 以 self 审计。
4. **无越界**：不改变既有失败锁定/限流/改密吊销语义（只叠加第二因素）；不改变 Profile 默认集语义 / 协议 pin；共享基架问题回流 VP-009/VP-010；go 失效触发时门闩自动关闭。

## 未选方案

- 短信/邮件 OTP 作为首选因子 → 依赖外部通道与模板基建（B-09），安全与成本权衡差，S1 如需再复核。
