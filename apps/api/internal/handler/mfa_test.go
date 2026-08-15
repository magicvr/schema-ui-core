package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
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
	return []string{"N1", "N2"}, nil
}

func (s *fakeMFAService) AdminReset(userID string) error {
	delete(s.enrolled, userID)
	delete(s.required, userID)
	return nil
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
	mountRoutes(MFARoutes(env.a, service, env.operations, revoker, "admin.mfa"))
}

// Two-step login: password factor OK + MFA required → proof, no tokens; then
// /api/auth/mfa/verify completes the login with a real token pair.
func TestMFALoginTwoStep(t *testing.T) {
	env := newAuthTestEnv(t)
	fake := newFakeMFAService()
	revoker := &fakeSessionRevoker{}
	// A dedicated mux so the MFA gate can be injected into the login handler.
	mux := http.NewServeMux()
	RegisterWithMFA(mux, env.a, env.st, env.operations, testAdminPlan(t), nil, nil, fake)
	mountRoutes := func(routes []kernel.RouteContribution) {
		for _, r := range routes {
			mux.Handle(r.Method+" "+r.Pattern, r.Handler)
		}
	}
	mountRoutes(MFARoutes(env.a, fake, env.operations, revoker, "admin.mfa"))

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

	// Enroll → one-time payload.
	req = bearer(t, token, http.MethodPost, "/api/mfa/enroll", "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK || !strings.Contains(rr.Body.String(), "recoveryCodes") {
		t.Fatalf("enroll = %d: %s", rr.Code, rr.Body.String())
	}

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

	// Disable requires a valid code and revokes sessions.
	req = bearer(t, token, http.MethodPost, "/api/mfa/disable", `{"code":"wrong"}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	expectError(t, rr, http.StatusUnauthorized, "MFA_INVALID")
	req = bearer(t, token, http.MethodPost, "/api/mfa/disable", `{"code":"123456"}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("disable = %d: %s", rr.Code, rr.Body.String())
	}
	if len(revoker.revoked) != 1 || revoker.revoked[0] != "user-admin" {
		t.Fatalf("revoker = %v, want [user-admin]", revoker.revoked)
	}

	// Recovery rotate returns a fresh set.
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