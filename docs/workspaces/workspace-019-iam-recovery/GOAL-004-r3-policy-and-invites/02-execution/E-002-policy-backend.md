---
id: E-002
doc: execution-entry
goal: GOAL-004-r3-policy-and-invites
status: recorded
created: 2026-08-25
updated: 2026-08-25
version: 1.0.0
---

# E-002 · C2 迁移落地 + C3 策略半场（四口强制）

2026-08-25 完成：

- **C2 满**：0057 password_policy（单例行，可移植）+ user_password_history；0058 user_invites；checksum `bfc7c4f2…` / `04f77fdc…` / `a35bbb21…` 入冻结台账；identity.go head 56→59、三表清单、四处黄金断言同步；store 包全绿。
- **C3 策略半场**：authsession/password_policy.go（GetPasswordPolicy/UpdatePasswordPolicy/ValidateNewPassword/历史捕获）；UpdateUser 同事务推旧 hash 入历史按深度裁剪；四口接线——users Create、users Update(password)、account_self changePassword、recovery complete（均设密前 ValidateNewPassword，错误统一 INVALID_PASSWORD）。
- 邀请域与配置 API 未开始（C3 另一半）；Web（C4）未开始。
