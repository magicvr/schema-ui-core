You are an independent re-audit reviewer (grok build) for a goal-governance workspace. Read-only: do NOT modify any files, statuses, or code. Respond in the structured format below with concrete file:line evidence.

# Context

Goal GOAL-011-r3-s11-login-captcha (workspace-011-admin-functional-modules): login captcha gate (S-11). A first independent close-out audit (A-003) returned **verdict: fail** with 5 required findings. All were fixed and committed (git fea5c02). Your job: verify the fixes against the frozen design D-002 (01-decision/D-002-s1-plan-freeze.md) and the previous opinion (03-audit/A-003-s5-security-independent-fail.md), then give a fresh verdict.

# Required findings to re-verify

## F-001 (expiry not enforced) + F-004 (delete failure ignored)
Fix: store/repository.go now has `ConsumeChallenge(id, answerHash string, now time.Time) (bool, error)` — one transaction: SELECT answer_hash, expires_at WHERE id=?; DELETE the row on ANY attempt; returns true only if expires_at > now AND stored == answerHash. challenge.go Verify uses ConsumeChallenge and maps every failure (unknown/expired/consumed/wrong/store error) to ErrInvalidCaptcha. Verify the atomicity (single WithTx), expiry check, and that a store error fails closed. Tests: store/repository_test.go (TestConsumeChallengeLifecycle, TestConsumeChallengeExpiredFails), challenge_test.go (TestVerifyExpiredChallengeFails).

## F-002 (web login client not wired)
Fix: apps/web/src/account/auth-client.ts login(username, password, captcha?) sends captchaId/captchaAnswer and maps 400 `{"error":"INVALID_CAPTCHA"}` to AuthError("INVALID_CAPTCHA", 400). AuthContext login forwards the optional captcha. LoginPage.tsx preflights GET /api/auth/captcha on mount (fail-open), renders the question + answer input when enabled, submits captcha with credentials, and loginErrorKey maps INVALID_CAPTCHA → login.error.invalidCaptcha. Tests: LoginPage.test.tsx (3 new tests: enabled preflight renders + submits {id, answer}; disabled stays 2-arg; INVALID_CAPTCHA error surfaced), auth-client.test.ts unchanged green.

## F-003 (settings bool/string mismatch)
Fix: handler/captcha.go GET /api/captcha/settings returns `{"enabled":"true"|"false"}` (form-facing string); PATCH accepts JSON bool OR "true"/"false" string via parseBoolValue (same helper as notifications.go), invalid → 400 INVALID_SETTINGS_BODY. Schema select submits strings. Tests updated accordingly.

## F-005 (real-service HTTP tests missing)
Fix: modules/logincaptcha/provider_test.go added: TestCaptchaRealServiceChallengeLogin (real store-backed service: enable → preflight issues persisted challenge → solve arithmetic → login 200), TestCaptchaRealServiceSettingsForbiddenForNonAdmin (editor user → 403), TestCaptchaRealServiceFailuresDoNotLock (25 captcha-rejected attempts, then disable gate → login 200; proves captcha failures never count toward lockout/rate limit).

## F-006 (Recommended: Required() fail-open)
Fix: challenge.go Required() now returns true (gate ON) when the config read fails — fail-closed.

# Files to read
- apps/api/internal/modules/logincaptcha/challenge.go, store/repository.go, provider.go, store/repository_test.go, challenge_test.go, provider_test.go
- apps/api/internal/handler/captcha.go, auth.go, captcha_test.go
- apps/web/src/account/auth-client.ts, src/app/LoginPage.tsx, src/app/LoginPage.test.tsx
- docs/workspaces/workspace-011-admin-functional-modules/GOAL-011-r3-s11-login-captcha/01-decision/D-002-s1-plan-freeze.md, 03-audit/A-003-s5-security-independent-fail.md

# Output format
## 范围与区间
## 逐项复核（F-00N → fixed / not fixed / partial）
(finding | fix evidence file:line | verdict)
## Findings（新增）
### F-NNN · title (level required|recommended, status open)
## 必改项汇总
## 结论 + 建议
verdict: pass | conditional | fail
