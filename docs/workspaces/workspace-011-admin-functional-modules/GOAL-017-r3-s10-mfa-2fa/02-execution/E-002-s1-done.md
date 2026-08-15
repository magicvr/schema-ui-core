---
id: E-002
goal: GOAL-017-r3-s10-mfa-2fa
title: S1 方案冻结完成（MFA/2FA TOTP 设计）
date: 2026-08-15
status: recorded
parent: GOAL-017-r3-s10-mfa-2fa
created: 2026-08-15
updated: 2026-08-15
version: 1.0.0
---

# E-002 · S1 方案冻结完成（2026-08-15）

## 事实

- D-002 落盘：TOTP 因子模型（RFC 6238 自实现 + 恢复码 bcrypt + AES-GCM secret 加密）、两段登录（MFAVerifier 接口 + mfaProof + IssueTokensFor）、端点与权限（users.mfa-reset + 自服务）、迁移 0029/0030、S-11 合成裁定（串行并存 + 分轨）、协议对照（AUTH-004/AUTH-006 → 本地身份层扩展）、S2 清单。
- I-001~I-004 全部闭合（证据 = D-002 各节）；progress 0/5 → 1/5（S1 检查点）。
- 关键证据：handler/auth.go（登录链）、internal/auth/auth.go（Login 核心）、account_self.go（自服务先例）、logincaptcha（一次性挑战先例）、I-HOST-APP-001 AUTH-006。
- 审计：A-003 self（S1 方案）。
