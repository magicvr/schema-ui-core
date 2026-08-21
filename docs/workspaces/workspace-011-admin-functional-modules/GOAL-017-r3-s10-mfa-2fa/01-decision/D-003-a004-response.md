---
id: D-003
goal: GOAL-017-r3-s10-mfa-2fa
title: A-004 响应：S1 独立审计 required 全 fixed
date: 2026-08-15
status: accepted
parent: GOAL-017-r3-s10-mfa-2fa
created: 2026-08-15
updated: 2026-08-15
version: 1.0.0
---

# D-003 · A-004 响应（2026-08-15）

## 审计结论

A-004（grok build · grok-4.6 · high · independent）verdict **conditional**，2 required：

- **F-001（required med）**：数据模型与控制不一致——§2「行存在即启用」vs §4 pending；mfa_proofs 无 fail_count（却写 5 次失效）；无 last_used_step（却写同窗拒重放）；enroll 后会话过期会自锁。
- **F-002（required med）**：admin reset 未规定 token_version+1 + 吊销 refresh，弱于自助 disable。

## 响应（全 fixed）

1. **状态机一致化**：user_mfa 增 status（pending/active）——仅 active 触发登录 MFA（enroll 后 pending，confirm 后 active，防自锁）；mfa_proofs 增 fail_count（达 5 失效落库）；user_mfa 增 last_used_step（同窗重放拒绝落库）。
2. **admin reset 强化**：重置 = 解绑 + token_version+1 + 吊销全部 refresh（与自助 disable 同级）。

闭合路径：**fixed**（D-002 方案修正，写 0029 迁移前可核对）。复审：A-005 grok reaudit。
