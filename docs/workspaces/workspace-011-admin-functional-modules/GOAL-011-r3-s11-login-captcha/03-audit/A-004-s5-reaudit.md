---
id: A-004
goal: GOAL-011-r3-s11-login-captcha
source: independent
date: 2026-08-14
scope: S5 关门 · 复审（grok A-003 required 修复验证 + 新增 recommended）
verdict: conditional
auditor: grok-build
audit_type: close-out-reaudit
status: recorded
parent: GOAL-011-r3-s11-login-captcha
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# A-004 · independent 复审（S-11 · 修复后）

## 结论

**verdict: conditional**。A-003 全部 5 个 required（F-001~F-005）+ F-006 recommended 确认 **fixed**（代码 + 执行测试：go logincaptcha/handler 全绿、vitest LoginPage/auth-client 27/27）。新增 4 个 recommended（F-009~F-012）。原始意见：`attachments/grok-audit-s11-reaudit.txt`。

## 逐项复核（A-003）

| finding | 结论 |
|--------|------|
| F-001 过期未强制 | **fixed**：ConsumeChallenge 单事务 SELECT answer_hash, expires_at → DELETE → expiresAt>now 才 matched（store/repository.go:57-76）；TestConsumeChallengeExpiredFails |
| F-002 Web 未接线 | **fixed**：auth-client captcha 字段 + 400 INVALID_CAPTCHA 映射；LoginPage 预检/输入/提交；单测 3 个 |
| F-003 settings bool/string | **fixed**：GET 字符串 "true"/"false"，PATCH parseBoolValue（notifications 模式） |
| F-004 删除失败 fail-open | **fixed**：Delete 失败 → error → ErrInvalidCaptcha；无 `_ =` 吞错 |
| F-005 真服务 HTTP 测试 | **fixed**：解题→200；editor 403；25 次 captcha 失败不计锁定/限流 |
| F-006 Required fail-open | **fixed**：配置读失败 → gate ON |

## 新增 Findings（recommended）

- **F-009**（recommended）：INVALID_CAPTCHA 失败后 UI 不刷新挑战（挑战已被消费）。→ **fixed**：LoginPage 失败后 refreshCaptcha() + 清空答案；测试断言二次预检。
- **F-010**（recommended）：TestVerifyExpiredChallengeFails 用 "hash-x"（与哈希方案不一致，混入错答拒绝）。→ **fixed**：改用 answerHash(id,"42") 隔离过期拒绝；store 级过期测试保留。
- **F-011**（recommended）：Required() fail-closed 无直接单测。→ **fixed**：TestRequiredFailsClosedOnConfigError（failingRunner）。
- **F-012**（recommended）：公开预检 enabled 时每次 Generate（无限流；受 5 分钟 TTL 约束）。→ **documented residual**：预检为匿名公开端点（D-002 §1 明确），挑战生成为 sha256 哈希 + 随机数，成本低；TTL 5 分钟限制表膨胀；登录本身受限流/锁定/captcha 门禁叠加。复核触发：若未来预检端点出现滥用证据（日志观测 QPS），为预检加每 IP 限流。

## 必改项汇总

- required：无（A-003 全部闭合）。
- recommended：F-009/F-010/F-011 已 fixed；F-012 documented residual（本文件留痕）。

## 结论

required 门禁无开放项 → S5 可关门（P-003 三路径：全部 required 走 fixed）。
