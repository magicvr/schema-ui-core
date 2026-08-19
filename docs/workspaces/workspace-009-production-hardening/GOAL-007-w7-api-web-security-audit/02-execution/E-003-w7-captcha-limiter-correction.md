---
id: E-003
goal: GOAL-007-w7-api-web-security-audit
title: F-006 captcha limiter correction after independent A-003 finding
date: 2026-08-19
status: recorded
parent: GOAL-001-production-hardening
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
---

# E-003 · F-006 captcha limiter correction

## 触发

- 独立审计 A-003（grok-4.6 high）verdict **conditional**：A-001 F-006 仍 open，因为 `captchaGenerateLimiter` 只调用了 `allow()`，而 `loginRateLimiter.allow()` 只检查已有失败记录、不会创建计数；匿名客户端仍可无限生成验证码。

## 已发生事实

- 修正 `apps/api/internal/handler/captcha.go`：在 `allow()` 通过后调用 `captchaGenerateLimiter.record()` 记录本次生成请求，使滑动窗口真正计数（10 次/分钟）。
- 新增回归测试 `TestCaptchaPreflightRateLimited`：连续 10 次 200 后第 11 次 429 `RATE_LIMITED`。

## 回归证据

- `cd apps/api && go test ./internal/handler -run 'TestCaptcha' -count=1 -timeout 90s`：通过。
- 该修正与已跑通的 `go test ./...` 全绿（此前全绿；本修正为增量小改动且已测）。
- 待独立复审 A-004 确认后闭合 F-006。