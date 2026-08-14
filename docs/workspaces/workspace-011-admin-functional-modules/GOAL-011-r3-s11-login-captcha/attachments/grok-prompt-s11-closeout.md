You are an independent close-out auditor (grok build) for a goal-governance workspace. Write a structured audit opinion ONLY — you must NOT modify any files, statuses, or code. Read the repository files to verify claims. Respond in the structured format below, in English or Chinese as convenient, with concrete file:line evidence.

# Scope

Goal GOAL-011-r3-s11-login-captcha (workspace-011-admin-functional-modules): S-11 登录验证码 (login captcha) — an optional arithmetic-captcha gate on POST /api/auth/login, default OFF, with a public preflight (GET /api/auth/captcha), admin settings (GET/PATCH /api/captcha/settings), page "captcha" + menu_captcha, permissions captcha.read/captcha.write, audit event captcha.settings-update, migrations 0023 (logincaptcha) + 0024 (operationlog CHECK extension).

Frozen design: docs/workspaces/workspace-011-admin-functional-modules/GOAL-011-r3-s11-login-captcha/01-decision/D-002-s1-plan-freeze.md (and D-001-goal-boundaries.md).

# Covered artifacts

- apps/api/internal/modules/logincaptcha/challenge.go, store/repository.go, provider.go, migration/migration.go, migration/provider.go, schema/captcha.json, manifest/fragment.json
- apps/api/internal/handler/auth.go (login gate integration), handler/captcha.go (routes), handler/health.go (variadic verifier)
- apps/api/internal/modules/operationlog/migration/migration.go (0024 CHECK), operationlog/repository.go (EventCaptchaSettingsUpdate)
- apps/api/internal/kernel/profile.go, apps/api/internal/composition/composition.go, apps/api/internal/testsupport/store.go, apps/api/internal/modules/compiled/persistence.go
- apps/api/internal/errorcatalog/errorcatalog.go (INVALID_CAPTCHA), handler/error_contract_test.go
- tests: modules/logincaptcha/*_test.go, handler/captcha_test.go, store/migrate_test.go + restart/operations tests (22→24), composition_test.go (22 perms / 12 nav)
- web: apps/web/src/test-fixtures/app-manifest.admin.json, apps/web/src/protocol/upstream-fixtures.test.ts (STATIC_MANIFEST_SHA256), apps/web/src/i18n/messages/*.json, e2e shell.spec.ts, scripts/smoke.sh

# Verification focus (security gate)

1. Login gate: verify order decode → rate limiter → captcha (only when h.captcha != nil && Required()) → Login; failure → 400 INVALID_CAPTCHA; captcha failures do NOT count toward lockout/rate-limit; consumed one-time on ANY attempt; answer stored hashed (sha256(id+":"+answer)); default OFF (captcha_config enabled=0) so the login contract is byte-identical when the module is disabled or gate off.
2. Preflight GET /api/auth/captcha is public (no auth middleware) and returns {enabled:false} with no challenge when off; no information leak.
3. Settings GET/PATCH /api/captcha/settings are fail-closed (401 anonymous, 403 without captcha.read/captcha.write), PATCH validates body {"enabled": bool} → 400 INVALID_SETTINGS_BODY otherwise, writes audit event captcha.settings-update.
4. Migrations: 0023 (captcha_challenges id/answer_hash/expires_at/created_at; captcha_config single-row id=1 enabled DEFAULT 0), 0024 rebuild CHECK + captcha.settings-update; ownership test frozen checksums; versions 1..24 continuous.
5. Provider registration: descriptor routes/pages/navigation/permissions/fragments declared exactly; composition wires a single service instance to both verifier and provider; nil verifier when module disabled; profile admin-only content extension (mvp/demo unchanged 8 perms / 5 nav).
6. Error code contract: INVALID_CAPTCHA added to frozen set + catalog bilingual + web i18n.
7. Tests cover: challenge lifecycle (generate/verify/one-time/expiry/lazy purge), gate off = original behavior, gate on wrong → 400 INVALID_CAPTCHA, gate on correct → 200, settings permission/body, audit event, provider surface, migration counts.

# Output format

## 范围与区间
## 成果（有证据）
(claims table: 主张 | 证据 file:line)
## 对照成功标准（本 scope）
(标准 | 结论 with 满足/部分满足/不满足)
## Findings
### F-NNN · title
| level | required | recommended |
| status | open |
| evidence |
## 必改项汇总
## 与既有意见的异同
(compare with A-001/A-002 self audits if readable)
## 结论 + 建议
verdict: pass | conditional | fail
