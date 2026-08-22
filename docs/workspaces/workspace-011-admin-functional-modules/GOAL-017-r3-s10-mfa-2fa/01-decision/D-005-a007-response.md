---
id: D-005
goal: GOAL-017-r3-s10-mfa-2fa
title: A-007/A-008 响应：required 全 fixed + MFA UI 用户裁决（阻断 → GOAL-018）
date: 2026-08-15
status: accepted
parent: GOAL-017-r3-s10-mfa-2fa
created: 2026-08-15
updated: 2026-08-15
version: 1.0.0
---

# D-005 · A-007/A-008 响应（2026-08-15）

## 审计结论

- A-007（grok build · grok-4.6 · high · independent）verdict **fail**：F-001（high，auth-client 两段登录不可达）、F-002（med，Enroll 覆盖 active 行）。
- A-008（grok reaudit）verdict **pass**：两条 required 已合法闭合（fixed 可核对）。

## 响应

1. **F-001（fixed）**：auth-client login() 的 mfaRequired 分支移到 token 形状校验之前；契约测试（auth-client.test + LoginPage.test 两段流）钉住无 token 第一段响应。
2. **F-002（fixed）**：Service.Enroll 预检（active → ErrActive）+ store UpsertPending WHERE status='pending' 守卫 + ErrActiveConflict；回归测试覆盖 active 拒绝 / pending 覆盖 / disable 后重登记；handler enroll 错误映射 400 MFA_ALREADY_ACTIVE（A-008 recommended F-001）。
3. **MFA UI 残余（用户 2026-08-15 书面裁决，mfa-ui-residual）**：**阻断 GOAL-017 关门**——新建下级子目标 [GOAL-018-mfa-manager-ui](../../GOAL-018-mfa-manager-ui/00-meta.md) 交付自定义 MFA 管理组件；GOAL-018 关门后回归关闭本目标。
4. A-008 recommended F-002（mfaVerify/AuthContext 整链客户端测试）：登记后续补强；F-003（pending 断言 + 注释）：已随回归测试与注释修正。

## 关门状态

**GOAL-017 暂不关门**（用户裁决：GOAL-018 关门后回归关闭）。A-007 required 已闭合（A-008 pass）。
