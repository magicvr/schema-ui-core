---
doc_type: goal-execution
id: E-002-handler-migration
parent: GOAL-003-r2-handler-migration
date: 2026-09-03
checkpoint: b08798d4
status: completed
---

# E-002 · 14 处生产调用点迁移完成与回归验证

## 1. 事实简述

依据 D-002 冻结合同及 GOAL-003 D-001 决策，将 HEAD（`98edb03e`）扫描确定的 **14 处生产 Allow→Record 调用点** 100% 迁移至原子 `AllowRecord`，生产环境完全消除两段式 Allow→Record 调用与 TOCTOU 窗口。

## 2. 改造清单核账（14/14）

| # | 位置 | 模式 | 改造事实 |
|---|------|------|----------|
| 1 | `apps/api/internal/handler/auth.go` (登录失败桶) | 失败预算 | 入口 `rateLimiter.AllowRecord`；移除 ErrAccountLocked / ErrAccountDisabled / ErrInvalidCredentials / MFARequired 二次 Record；成功保持 `rateLimiter.Clear` |
| 2 | `apps/api/internal/handler/captcha.go` (验证码) | 立即消费 | `limiter.Allow` + `Record` → `!limiter.AllowRecord` (429) |
| 3 | `apps/api/internal/handler/account_self.go` (密码修改) | 失败预算 | 入口 `passwordLimiter.AllowRecord`；移除密码错误时的二次 `Record`；成功保持 `passwordLimiter.Clear` |
| 4 | `apps/api/internal/handler/recovery.go` (自助恢复 start) | 失败预算 | 入口 `rateLimiter.AllowRecord`；移除 ErrRecoveryNotAvailable 二次 `Record`；成功 `rateLimiter.Clear` |
| 5 | `apps/api/internal/handler/recovery.go` (自助恢复 complete) | 失败预算 | 入口 `rateLimiter.AllowRecord`；`recordFailure` 改为 no-op（不再二次 Record）；RecoveryMatch 后非猜测试错（密码长度/复杂度/二次认证要求）与最终成功均调用 `Clear` |
| 6 | `apps/api/internal/handler/mfa.go` (MFA verify) | 失败预算 | 入口 `mfaVerifyLimiter.AllowRecord`；移除校验失败的二次 `Record`；校验成功新增 `mfaVerifyLimiter.Clear` |
| 7 | `apps/api/internal/handler/mfa.go` (step-up enroll) | 失败预算 | `guardMFAStepUp` 内部改为 `limiter.AllowRecord`；移除密码校验失败二次 `Record`；成功保持 `Clear` |
| 8 | `apps/api/internal/handler/mfa.go` (step-up disable) | 失败预算 | `guardMFAStepUp` 改为 `AllowRecord`；移除 Disable 失败二次 `Record`；成功保持 `Clear` |
| 9 | `apps/api/internal/handler/mfa.go` (step-up recovery-rotate) | 失败预算 | `guardMFAStepUp` 改为 `AllowRecord`；移除 Rotate 失败二次 `Record`；成功保持 `Clear` |
| 10 | `apps/api/internal/handler/invites.go` (邀请接受) | 失败预算 | 入口 `limiter.AllowRecord`；移除 PeekInviteToken 与 AcceptInvite 失败二次 `Record`；成功新增 `limiter.Clear` |
| 11 | `apps/api/internal/handler/wallet_self.go` (钱包核销) | 失败预算 | 入口 `redeemLimiter.AllowRecord`；移除无效 body 及 RedeemForUser 失败二次 `Record`；成功保持 `Clear` |
| 12 | `apps/api/internal/channel/telegram/webhook.go` (IP 桶) | 立即消费 | `ipLimiter.Allow` + `Record` → `!ipLimiter.AllowRecord` (429) |
| 13 | `apps/api/internal/channel/telegram/webhook.go` (Chat 桶) | 立即消费 | `chatLimiter.Allow` + `Record` → `!chatLimiter.AllowRecord` (429) |
| 14 | `apps/api/internal/channel/telegram/webhook.go` (User 桶) | 立即消费 | `userLimiter.Allow` + `Record` → `!userLimiter.AllowRecord` (429) |

静态核验结果：除 `resources.go` 的垃圾回收领域记录方法 `Trash.Record` 外，`apps/api/internal/handler` 与 `apps/api/internal/channel/telegram` 中生产代码不再存在任何 `.Allow(` 或 `.Record(` 调用。

## 3. 测试与回归证据

1. **Telegram Webhook 回归与并发无穿透验证**：
   - 既有 IP/Chat/User 限流单测全部通过（`TestWebhook_RateLimiting_*`）。
   - 新增 `TestWebhook_RateLimiting_ConcurrentNoTOCTOU`：100 个并发请求打入单个 IP 桶（阈值 60），精确断言通过数为 60，429 拦截数为 40，无穿透。
   - `go test -v ./internal/channel/telegram/...` 全部通过（PASS，1.72s）。
2. **Handler 既有限流单测全绿**：
   - `TestLoginRateLimit`、`TestLoginClientIPTrustsXRealIPOnlyFromTrustedPeer` 全部通过。
   - `TestPasswordChangeRateLimited` 全部通过。
   - `TestRecoveryCompleteRateLimitedAfterTwentyFailures`、`TestRecoveryCompleteNonGuessFailuresDoNotRecord` 全部通过。
   - `TestMFAVerifyRateLimit`、`TestMFAVerifyRateLimitPerIP`、`TestMFAVerifyRateLimitDoesNotBlockNormalFlow`、`TestMFAStepUpDisableAndRotateRateLimited`、`TestMFAEnrollWrongPasswordRateLimited` 全部通过。
   - `TestInviteAcceptPublicSurface`、`TestInviteAcceptRateLimitBoundsUnauthenticatedSpray` 全部通过。
   - `TestWalletSelfRedeemRateLimited` 全部通过。
3. **Handler 新增并发与净状态等价测试**：
   - 新增 `TestLoginRateLimit_ConcurrentNoTOCTOUPenetration`：50 个并发登录请求打入同一 IP\|user 桶（阈值 20），精确断言通过数为 20，429 为 30，并发无穿透。
   - 新增 `TestLoginRateLimit_SuccessfulLoginClearsFailureBucket`：循环 10 轮「3 次失败 + 1 次成功」，累计 30 次失败请求，因成功登录执行 `Clear`，全程未触发 20 次拦截，验证净状态等价与 Clear 有效性。
4. **Race 检测器全绿**：
   - `go test -race -run "Test(LoginRateLimit|PasswordChange|MFA|Recovery|InviteAccept|WalletSelf|Captcha)" ./internal/handler` 通过（PASS，94.381s）。
   - `go test -race ./internal/channel/telegram/...` 通过（PASS，6.382s）。

## 4. Git Checkpoint

- Commit: `b08798d4`
- Commit Message: `feat(ratelimit): 迁移 14 处生产调用点至原子 AllowRecord 并补齐并发回归`
- 仅暂存 10 个显式 owned paths，无侵入其它模块或文件。
