package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/ratelimit"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// fakeMFAService is an in-memory MFAVerifier + MFASelfService test double
// (the real store-backed service is covered by the module tests).
type fakeMFAService struct {
	required map[string]bool
	proofs   map[string]string // proof → userID
	enrolled map[string]bool
}

func newFakeMFAService() *fakeMFAService {
	return &fakeMFAService{
		required: map[string]bool{},
		proofs:   map[string]string{},
		enrolled: map[string]bool{},
	}
}

func (s *fakeMFAService) Required(userID string) bool { return s.required[userID] }

func (s *fakeMFAService) BeginChallenge(userID string, now time.Time) (string, error) {
	proof := "proof-" + userID
	s.proofs[proof] = userID
	return proof, nil
}

func (s *fakeMFAService) Verify(proof, code, recoveryCode string, now time.Time) (string, error) {
	userID, ok := s.proofs[proof]
	if !ok {
		return "", ErrMFAProofExpired
	}
	if code != "123456" && recoveryCode != "RECOVERY" {
		return "", ErrMFAInvalid
	}
	delete(s.proofs, proof)
	return userID, nil
}

func (s *fakeMFAService) Status(userID string) (bool, time.Time, error) {
	return s.enrolled[userID], time.Time{}, nil
}

func (s *fakeMFAService) Enroll(userID, name string, now time.Time) (string, string, []string, error) {
	return "SECRET", "otpauth://totp/", []string{"C1", "C2"}, nil
}

func (s *fakeMFAService) Confirm(userID, code string, now time.Time) error {
	if code != "123456" {
		return ErrMFAInvalid
	}
	s.enrolled[userID] = true
	s.required[userID] = true
	return nil
}

func (s *fakeMFAService) Disable(userID, code, recoveryCode string, now time.Time) error {
	if code != "123456" && recoveryCode != "RECOVERY" {
		return ErrMFAInvalid
	}
	delete(s.enrolled, userID)
	delete(s.required, userID)
	return nil
}

func (s *fakeMFAService) RotateRecovery(userID, code, recoveryCode string, now time.Time) ([]string, error) {
	if code != "123456" && recoveryCode != "RECOVERY" {
		return nil, ErrMFAInvalid
	}
	return []string{"N1", "N2"}, nil
}

func (s *fakeMFAService) AdminReset(userID string) (bool, error) {
	wasActive := s.required[userID]
	delete(s.enrolled, userID)
	delete(s.required, userID)
	return wasActive, nil
}

// fakeSessionRevoker records revocation calls.
type fakeSessionRevoker struct {
	revoked []string
}

func (r *fakeSessionRevoker) BumpTokenVersionAndRevokeAll(userID string, now time.Time) error {
	r.revoked = append(r.revoked, userID)
	return nil
}

func mountMFASurface(t *testing.T, env *authTestEnv, service MFASelfService, revoker SessionRevoker) {
	t.Helper()
	mountRoutes := func(routes []kernel.RouteContribution) {
		for _, r := range routes {
			env.mux.Handle(r.Method+" "+r.Pattern, r.Handler)
		}
	}
	mountRoutes(MFARoutes(env.a, service, env.operations, revoker, "admin.mfa", ratelimit.NewProvider()))
}

// Two-step login: password factor OK + MFA required → proof, no tokens; then
// /api/auth/mfa/verify completes the login with a real token pair.
func TestMFALoginTwoStep(t *testing.T) {
	env := newAuthTestEnv(t)
	fake := newFakeMFAService()
	revoker := &fakeSessionRevoker{}
	// A dedicated mux so the MFA gate can be injected into the login handler.
	mux := http.NewServeMux()
	RegisterWithMFA(mux, env.a, env.st, env.operations, testAdminPlan(t), nil, ratelimit.NewProvider(), nil, fake)
	mountRoutes := func(routes []kernel.RouteContribution) {
		for _, r := range routes {
			mux.Handle(r.Method+" "+r.Pattern, r.Handler)
		}
	}
	mountRoutes(MFARoutes(env.a, fake, env.operations, revoker, "admin.mfa", ratelimit.NewProvider()))

	// Without an enrollment: normal login issues tokens.
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"username":"admin","password":"test-password"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("plain login = %d: %s", rr.Code, rr.Body.String())
	}

	// With MFA required: no tokens; a proof is issued.
	fake.required["user-admin"] = true
	req = httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"username":"admin","password":"test-password"}`))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("mfa login = %d: %s", rr.Code, rr.Body.String())
	}
	var pending map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&pending)
	if pending["mfaRequired"] != true || pending["mfaProof"] != "proof-user-admin" {
		t.Fatalf("pending = %v", pending)
	}

	// Wrong code → MFA_INVALID; correct code → real tokens.
	req = httptest.NewRequest(http.MethodPost, "/api/auth/mfa/verify", strings.NewReader(`{"proof":"proof-user-admin","code":"000000"}`))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	expectError(t, rr, http.StatusUnauthorized, "MFA_INVALID")

	req = httptest.NewRequest(http.MethodPost, "/api/auth/mfa/verify", strings.NewReader(`{"proof":"proof-user-admin","code":"123456"}`))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("verify = %d: %s", rr.Code, rr.Body.String())
	}
	var done map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&done)
	if done["accessToken"] == nil || done["refreshToken"] == nil {
		t.Fatalf("verify response missing tokens: %v", done)
	}

	// Replayed proof → MFA_PROOF_EXPIRED (one-shot).
	req = httptest.NewRequest(http.MethodPost, "/api/auth/mfa/verify", strings.NewReader(`{"proof":"proof-user-admin","code":"123456"}`))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	expectError(t, rr, http.StatusUnauthorized, "MFA_PROOF_EXPIRED")
}

// Self-service surface: gates, enroll/confirm/disable + session invalidation,
// recovery rotate, admin reset permission.
func TestMFASelfService(t *testing.T) {
	env := newAuthTestEnv(t)
	fake := newFakeMFAService()
	revoker := &fakeSessionRevoker{}
	mountMFASurface(t, env, fake, revoker)
	token := adminToken(t, env)

	// Anonymous → 401.
	req := httptest.NewRequest(http.MethodGet, "/api/mfa/status", nil)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	expectError(t, rr, http.StatusUnauthorized, "UNAUTHENTICATED")

	// Enroll → one-time payload (W7 F-007: requires current password).
	req = bearer(t, token, http.MethodPost, "/api/mfa/enroll", `{"currentPassword":"test-password"}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK || !strings.Contains(rr.Body.String(), "recoveryCodes") {
		t.Fatalf("enroll = %d: %s", rr.Code, rr.Body.String())
	}

	// Confirm with a wrong code → 400 MFA_INVALID (client validation, NOT a
	// lost session — W11 M-03: the web auth wrapper must not sign the user out).
	req = bearer(t, token, http.MethodPost, "/api/mfa/confirm", `{"code":"000000"}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	expectError(t, rr, http.StatusBadRequest, "MFA_INVALID")

	// Confirm activates.
	req = bearer(t, token, http.MethodPost, "/api/mfa/confirm", `{"code":"123456"}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("confirm = %d: %s", rr.Code, rr.Body.String())
	}
	req = bearer(t, token, http.MethodGet, "/api/mfa/status", "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK || !strings.Contains(rr.Body.String(), "\"enabled\":true") {
		t.Fatalf("status = %d: %s", rr.Code, rr.Body.String())
	}

	// Disable requires a valid code and revokes sessions. A wrong code is a
	// client validation failure (400, W11 M-02) — never a 401 session loss.
	req = bearer(t, token, http.MethodPost, "/api/mfa/disable", `{"code":"wrong"}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	expectError(t, rr, http.StatusBadRequest, "MFA_INVALID")
	req = bearer(t, token, http.MethodPost, "/api/mfa/disable", `{"code":"123456"}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("disable = %d: %s", rr.Code, rr.Body.String())
	}
	if len(revoker.revoked) != 1 || revoker.revoked[0] != "user-admin" {
		t.Fatalf("revoker = %v, want [user-admin]", revoker.revoked)
	}

	// Recovery rotate: a wrong code is a 400 validation failure (same
	// self-service semantics as confirm/disable — W11 M-02/M-03); a valid
	// recovery code returns a fresh set.
	req = bearer(t, token, http.MethodPost, "/api/mfa/recovery/rotate", `{"code":"wrong"}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	expectError(t, rr, http.StatusBadRequest, "MFA_INVALID")
	req = bearer(t, token, http.MethodPost, "/api/mfa/recovery/rotate", `{"recoveryCode":"RECOVERY"}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK || !strings.Contains(rr.Body.String(), "N1") {
		t.Fatalf("rotate = %d: %s", rr.Code, rr.Body.String())
	}

	// Admin reset: editor (no users.mfa-reset) → 403; admin → 204 + revoke.
	env.addUser(t, "editor1", "editor-password", []string{"editor"})
	editorToken := env.login(t, "editor1", "editor-password")
	req = bearer(t, editorToken, http.MethodPost, "/api/users/user-editor1/mfa/reset", "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	expectError(t, rr, http.StatusForbidden, "FORBIDDEN")
	// Mark editor1 as having ACTIVE MFA so the admin reset reports an active
	// removal and revokes sessions (W7 F-002: no-enrollment resets must NOT
	// force-logout; this positive case proves the active path still revokes).
	fake.required["user-editor1"] = true
	fake.enrolled["user-editor1"] = true
	req = bearer(t, token, http.MethodPost, "/api/users/user-editor1/mfa/reset", "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("admin reset = %d: %s", rr.Code, rr.Body.String())
	}
	if len(revoker.revoked) != 2 || revoker.revoked[1] != "user-editor1" {
		t.Fatalf("revoker after reset = %v", revoker.revoked)
	}
}

var _ = auth.ErrInvalidCredentials

// TestMFAResetAdminTargetBoundary verifies W7 F-002 / A-003 F-004: a delegated
// (non-admin) actor who DOES hold users.mfa-reset must still be blocked from
// resetting an ADMIN account's MFA (ADMIN_ACCOUNT_FORBIDDEN), while remaining
// able to reset a non-admin account.
func TestMFAResetAdminTargetBoundary(t *testing.T) {
	env := newAuthTestEnv(t)
	fake := newFakeMFAService()
	revoker := &fakeSessionRevoker{}
	mountMFASurface(t, env, fake, revoker)

	now := time.Now().UTC()
	if _, err := env.authRepository.CreateRoleWithGrants(
		"mfa-resetter", "MFA resetter",
		[]string{"users.read", "users.mfa-reset"}, nil, now,
	); err != nil {
		t.Fatalf("create mfa-resetter role: %v", err)
	}
	env.addUser(t, "resetter", "resetter-password", []string{"mfa-resetter"})
	resetterToken := env.login(t, "resetter", "resetter-password")

	// The seeded admin user is a real admin (has the admin role).
	adminID := "user-" + testSeedUsername

	// Resetting an admin target from a non-admin delegated actor → 403
	// ADMIN_ACCOUNT_FORBIDDEN (the boundary, not the permission gate — the
	// actor legitimately holds users.mfa-reset).
	req := bearer(t, resetterToken, http.MethodPost, "/api/users/"+adminID+"/mfa/reset", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	var out map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&out)
	if rr.Code != http.StatusForbidden || out["error"] != "ADMIN_ACCOUNT_FORBIDDEN" {
		t.Fatalf("delegated reset of admin = %d %v, want 403 ADMIN_ACCOUNT_FORBIDDEN", rr.Code, out)
	}
	if len(revoker.revoked) != 0 {
		t.Fatalf("revoker = %v, want none (boundary rejected before reset)", revoker.revoked)
	}

	// The same delegated actor may reset a NON-admin user's MFA.
	env.addUser(t, "staff", "staff-password", []string{"editor"})
	fake.required["user-staff"] = true
	fake.enrolled["user-staff"] = true
	req = bearer(t, resetterToken, http.MethodPost, "/api/users/user-staff/mfa/reset", "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("delegated reset of non-admin = %d, want 204: %s", rr.Code, rr.Body.String())
	}
	if len(revoker.revoked) != 1 || revoker.revoked[0] != "user-staff" {
		t.Fatalf("revoker after non-admin reset = %v, want [user-staff]", revoker.revoked)
	}
}

// A6 · MFA verify independent rate limit: POST /api/auth/mfa/verify has its
// own counting bucket (independent of the login limiter); repeated failed
// attempts from one IP trigger 429 RATE_LIMITED with a Retry-After header.
// A single successful-path attempt is never blocked.
func TestMFAVerifyRateLimit(t *testing.T) {
	env := newAuthTestEnv(t)
	fake := newFakeMFAService()
	revoker := &fakeSessionRevoker{}
	mux := http.NewServeMux()
	RegisterWithMFA(mux, env.a, env.st, env.operations, testAdminPlan(t), nil, ratelimit.NewProvider(), nil, fake)
	mountRoutes := func(routes []kernel.RouteContribution) {
		for _, r := range routes {
			mux.Handle(r.Method+" "+r.Pattern, r.Handler)
		}
	}
	mountRoutes(MFARoutes(env.a, fake, env.operations, revoker, "admin.mfa", ratelimit.NewProvider()))

	// Enroll the admin user so BeginChallenge works.
	fake.required["user-admin"] = true

	// Helper: POST /api/auth/mfa/verify with an invalid code and an invalid
	// proof (so it hits the rate-limiter record path on each failure).
	verifyInvalid := func() *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/mfa/verify",
			strings.NewReader(`{"proof":"no-such-proof","code":"000000"}`))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		return rr
	}

	// The first mfaVerifyRateLimiterMax attempts must return a non-429 code
	// (the service rejects the bogus proof with 401 MFA_PROOF_EXPIRED, not a
	// rate limit error).
	for i := range mfaVerifyRateLimiterMax {
		rr := verifyInvalid()
		if rr.Code == http.StatusTooManyRequests {
			t.Fatalf("attempt %d (within budget) got 429, want <429", i+1)
		}
	}

	// The very next attempt (one over the limit) must be 429 RATE_LIMITED.
	rr := verifyInvalid()
	if rr.Code != http.StatusTooManyRequests {
		t.Fatalf("over-limit attempt = %d, want 429", rr.Code)
	}
	var out map[string]string
	_ = json.NewDecoder(rr.Body).Decode(&out)
	if out["error"] != "RATE_LIMITED" {
		t.Fatalf("over-limit code = %q, want RATE_LIMITED", out["error"])
	}
	if rr.Header().Get("Retry-After") == "" {
		t.Fatal("RATE_LIMITED response missing Retry-After header")
	}
}

// A6 · MFA verify rate limit is per-IP and independent of the login limiter:
// two different client IPs each get their own bucket; exhausting one does not
// affect the other.
func TestMFAVerifyRateLimitPerIP(t *testing.T) {
	env := newAuthTestEnv(t)
	fake := newFakeMFAService()
	revoker := &fakeSessionRevoker{}
	mux := http.NewServeMux()
	RegisterWithMFA(mux, env.a, env.st, env.operations, testAdminPlan(t), nil, ratelimit.NewProvider(), nil, fake)
	mountRoutes := func(routes []kernel.RouteContribution) {
		for _, r := range routes {
			mux.Handle(r.Method+" "+r.Pattern, r.Handler)
		}
	}
	mountRoutes(MFARoutes(env.a, fake, env.operations, revoker, "admin.mfa", ratelimit.NewProvider()))

	verifyFromIP := func(ip string) int {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/mfa/verify",
			strings.NewReader(`{"proof":"no-such-proof","code":"000000"}`))
		req.Header.Set("Content-Type", "application/json")
		req.RemoteAddr = fmt.Sprintf("%s:50000", ip)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		return rr.Code
	}

	// Exhaust IP-A's bucket.
	for range mfaVerifyRateLimiterMax {
		verifyFromIP("10.0.0.1")
	}
	if code := verifyFromIP("10.0.0.1"); code != http.StatusTooManyRequests {
		t.Fatalf("IP-A over limit = %d, want 429", code)
	}

	// IP-B's bucket is untouched — must not be rate-limited.
	if code := verifyFromIP("10.0.0.2"); code == http.StatusTooManyRequests {
		t.Fatalf("IP-B erroneously rate-limited after exhausting IP-A's bucket")
	}
}

// A6 · normal single verify is not affected: a valid proof + code succeeds on
// the first attempt without hitting the rate limiter.
func TestMFAVerifyRateLimitDoesNotBlockNormalFlow(t *testing.T) {
	env := newAuthTestEnv(t)
	fake := newFakeMFAService()
	revoker := &fakeSessionRevoker{}
	mux := http.NewServeMux()
	RegisterWithMFA(mux, env.a, env.st, env.operations, testAdminPlan(t), nil, ratelimit.NewProvider(), nil, fake)
	mountRoutes := func(routes []kernel.RouteContribution) {
		for _, r := range routes {
			mux.Handle(r.Method+" "+r.Pattern, r.Handler)
		}
	}
	mountRoutes(MFARoutes(env.a, fake, env.operations, revoker, "admin.mfa", ratelimit.NewProvider()))

	// Set up a valid proof as the login flow would.
	fake.required["user-admin"] = true
	proof, err := fake.BeginChallenge("user-admin", time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/auth/mfa/verify",
		strings.NewReader(`{"proof":"`+proof+`","code":"123456"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("normal verify = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	var resp map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&resp)
	if resp["accessToken"] == nil {
		t.Fatal("normal verify response missing accessToken")
	}
}

// W13 F-002 (GOAL-013 A-001): failed second-factor proofs on the identity-
// scoped disable / recovery-rotate surfaces are bounded by the shared step-up
// bucket — after mfaStepUpLimiterMax failures per operation|IP|user, even a
// VALID code answers 429 RATE_LIMITED instead of leaving an unbounded TOTP
// guessing oracle behind a session token.
func TestMFAStepUpDisableAndRotateRateLimited(t *testing.T) {
	env := newAuthTestEnv(t)
	fake := newFakeMFAService()
	revoker := &fakeSessionRevoker{}
	mountMFASurface(t, env, fake, revoker)
	token := adminToken(t, env)

	// Activate a real enrollment so the surface matches production state.
	req := bearer(t, token, http.MethodPost, "/api/mfa/enroll", `{"currentPassword":"test-password"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("enroll = %d: %s", rr.Code, rr.Body.String())
	}
	req = bearer(t, token, http.MethodPost, "/api/mfa/confirm", `{"code":"123456"}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("confirm = %d: %s", rr.Code, rr.Body.String())
	}

	// Exhaust the DISABLE bucket with wrong codes (each records one failure).
	for i := range mfaStepUpLimiterMax {
		req = bearer(t, token, http.MethodPost, "/api/mfa/disable", `{"code":"wrong"}`)
		rr = httptest.NewRecorder()
		env.mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("disable failure %d status = %d (%s), want 400 MFA_INVALID", i+1, rr.Code, rr.Body.String())
		}
		expectError(t, rr, http.StatusBadRequest, "MFA_INVALID")
	}
	// Over budget → 429 even for a valid code.
	req = bearer(t, token, http.MethodPost, "/api/mfa/disable", `{"code":"123456"}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusTooManyRequests || !strings.Contains(rr.Body.String(), "RATE_LIMITED") {
		t.Fatalf("disable over budget = %d body %s, want 429 RATE_LIMITED", rr.Code, rr.Body.String())
	}
	if len(revoker.revoked) != 0 {
		t.Fatalf("rate-limited disable must not revoke sessions: %v", revoker.revoked)
	}

	// The RECOVERY-ROTATE bucket is keyed per operation — still usable.
	for i := range mfaStepUpLimiterMax {
		req = bearer(t, token, http.MethodPost, "/api/mfa/recovery/rotate", `{"code":"wrong"}`)
		rr = httptest.NewRecorder()
		env.mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("rotate failure %d status = %d (%s), want 400 MFA_INVALID", i+1, rr.Code, rr.Body.String())
		}
	}
	req = bearer(t, token, http.MethodPost, "/api/mfa/recovery/rotate", `{"recoveryCode":"RECOVERY"}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Fatalf("rotate over budget = %d body %s, want 429 RATE_LIMITED", rr.Code, rr.Body.String())
	}
}

// W13 F-003 (GOAL-013 A-001): wrong currentPassword guesses on enroll share
// the login-limiter model — bounded failures then 429 RATE_LIMITED.
func TestMFAEnrollWrongPasswordRateLimited(t *testing.T) {
	env := newAuthTestEnv(t)
	fake := newFakeMFAService()
	revoker := &fakeSessionRevoker{}
	mountMFASurface(t, env, fake, revoker)
	token := adminToken(t, env)

	for range mfaStepUpLimiterMax {
		req := bearer(t, token, http.MethodPost, "/api/mfa/enroll", `{"currentPassword":"wrong-password"}`)
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, req)
		expectError(t, rr, http.StatusBadRequest, "INVALID_PASSWORD")
	}
	req := bearer(t, token, http.MethodPost, "/api/mfa/enroll", `{"currentPassword":"test-password"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusTooManyRequests || !strings.Contains(rr.Body.String(), "RATE_LIMITED") {
		t.Fatalf("enroll over budget = %d body %s, want 429 RATE_LIMITED", rr.Code, rr.Body.String())
	}
}