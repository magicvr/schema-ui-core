// Self-recovery HTTP surface tests (workspace-019 R2 · GOAL-003 D-001 §2):
// the pre-auth envelope contract — enumeration-neutral start, uniform
// invalid-code answers, second-factor gate ordering, password baseline, and
// 204 success without token issuance.
package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	kernel "github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authsession "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
)

type fakeRecoveryRepo struct {
	target      *authsession.RecoveryTarget
	resolveErr  error
	startErr    error
	outcome     authsession.RecoveryOutcome
	voidOnNext  bool
	completed   []string
	userName    string
}

func (f *fakeRecoveryRepo) ResolveRecoveryTarget(identifier string) (*authsession.RecoveryTarget, error) {
	if f.resolveErr != nil {
		return nil, f.resolveErr
	}
	return f.target, nil
}
func (f *fakeRecoveryRepo) StartRecovery(identifier string, sender kernel.MailSender, now time.Time) error {
	return f.startErr
}
func (f *fakeRecoveryRepo) EvaluateRecoveryCode(userID, code string, now time.Time) (authsession.RecoveryOutcome, error) {
	return f.outcome, nil
}
func (f *fakeRecoveryRepo) ConsumeRecoveryAttempt(userID string) bool { return f.voidOnNext }
func (f *fakeRecoveryRepo) DropStaleRecoveryChallenge(userID string, now time.Time) {}
func (f *fakeRecoveryRepo) CompleteRecovery(userID, passwordHash, actorID string, now time.Time) error {
	f.completed = append(f.completed, userID+"|"+passwordHash)
	return nil
}
func (f *fakeRecoveryRepo) UserByID(id string) (*authsession.User, error) {
	return &authsession.User{ID: id, Name: f.userName}, nil
}

type fakeGate struct {
	required bool
	verifyErr error
}

func (g *fakeGate) Required(userID string) bool { return g.required }
func (g *fakeGate) VerifySecondFactor(userID, code, recoveryCode string, now time.Time) error {
	return g.verifyErr
}


func newRecoveryMux(repo RecoveryRepository, gate RecoverySecondFactor) *http.ServeMux {
	mux := http.NewServeMux()
	RegisterRecovery(mux, nil, repo, nil, nil, gate)
	return mux
}

func postJSON(t *testing.T, mux *http.ServeMux, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func TestRecoveryStartEnvelope(t *testing.T) {
	repo := &fakeRecoveryRepo{target: &authsession.RecoveryTarget{UserID: "u1", Enabled: true}}
	mux := newRecoveryMux(repo, nil)

	// invalid body → cataloged 400
	if rec := postJSON(t, mux, "/api/auth/recovery/start", `{}`); rec.Code != http.StatusBadRequest {
		t.Fatalf("empty account status = %d", rec.Code)
	} else if !strings.Contains(rec.Body.String(), "INVALID_RECOVERY_BODY") {
		t.Fatalf("body = %s, want INVALID_RECOVERY_BODY", rec.Body.String())
	}
	// no-path account → SAME dispatched shape (anti-enumeration)
	repo.startErr = authsession.ErrRecoveryNotAvailable
	rec := postJSON(t, mux, "/api/auth/recovery/start", `{"account":"ghost"}`)
	if rec.Code != http.StatusAccepted || !strings.Contains(rec.Body.String(), "dispatched") {
		t.Fatalf("no-path start = (%d, %s), want 202 dispatched", rec.Code, rec.Body.String())
	}
	// cooldown surfaces distinctly (a real user needs this signal)
	repo.startErr = authsession.ErrRecoveryCooldown
	if rec := postJSON(t, mux, "/api/auth/recovery/start", `{"account":"admin"}`); rec.Code != http.StatusTooManyRequests {
		t.Fatalf("cooldown status = %d, want 429", rec.Code)
	}
}

func TestRecoveryCompleteUniformInvalidCode(t *testing.T) {
	repo := &fakeRecoveryRepo{target: &authsession.RecoveryTarget{UserID: "u1"}}
	cases := []struct {
		name     string
		resolve  error
		outcome  authsession.RecoveryOutcome
	}{
		{"unknown identifier", errors.New("nope"), authsession.RecoveryNotPending},
		{"no live challenge", nil, authsession.RecoveryNotPending},
		{"wrong code", nil, authsession.RecoveryMismatch},
	}
	for _, tc := range cases {
		repo.resolveErr = tc.resolve
		repo.outcome = tc.outcome
		mux := newRecoveryMux(repo, nil)
		rec := postJSON(t, mux, "/api/auth/recovery/complete", `{"account":"x","code":"123456","newPassword":"new-password-1"}`)
		if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), "RECOVERY_CODE_INVALID") {
			t.Fatalf("%s = (%d, %s), want 400 RECOVERY_CODE_INVALID", tc.name, rec.Code, rec.Body.String())
		}
	}
	// expired is distinct so a real user knows to re-request
	repo.resolveErr = nil
	repo.outcome = authsession.RecoveryExpired
	mux := newRecoveryMux(repo, nil)
	if rec := postJSON(t, mux, "/api/auth/recovery/complete", `{"account":"x","code":"123456","newPassword":"new-password-1"}`); !strings.Contains(rec.Body.String(), "RECOVERY_CODE_EXPIRED") {
		t.Fatalf("expired body = %s", rec.Body.String())
	}
}

func TestRecoveryCompleteSecondFactorGate(t *testing.T) {
	repo := &fakeRecoveryRepo{target: &authsession.RecoveryTarget{UserID: "u1"}, outcome: authsession.RecoveryMatch}
	gate := &fakeGate{required: true}
	mux := newRecoveryMux(repo, gate)

	// missing factor → explicit demand, NO attempt burned for an empty field
	rec := postJSON(t, mux, "/api/auth/recovery/complete", `{"account":"x","code":"123456","newPassword":"new-password-1"}`)
	if !strings.Contains(rec.Body.String(), "RECOVERY_SECOND_FACTOR_REQUIRED") {
		t.Fatalf("missing factor body = %s", rec.Body.String())
	}
	// wrong factor → MFA_INVALID
	gate.verifyErr = ErrMFAInvalid
	if rec := postJSON(t, mux, "/api/auth/recovery/complete", `{"account":"x","code":"123456","secondFactorCode":"000000","newPassword":"new-password-1"}`); !strings.Contains(rec.Body.String(), "MFA_INVALID") {
		t.Fatalf("wrong factor body = %s", rec.Body.String())
	}
	// right factor → completion lands
	gate.verifyErr = nil
	rec = postJSON(t, mux, "/api/auth/recovery/complete", `{"account":"x","code":"123456","recoveryCode":"RECOV-1","newPassword":"new-password-1"}`)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("complete status = %d body %s, want 204", rec.Code, rec.Body.String())
	}
	if len(repo.completed) != 1 || !strings.HasPrefix(repo.completed[0], "u1|") {
		t.Fatalf("completed = %v", repo.completed)
	}
}

func TestRecoveryCompletePasswordBaseline(t *testing.T) {
	repo := &fakeRecoveryRepo{target: &authsession.RecoveryTarget{UserID: "u1"}, outcome: authsession.RecoveryMatch}
	mux := newRecoveryMux(repo, nil)
	for _, pw := range []string{"short", strings.Repeat("x", 73)} {
		body := `{"account":"x","code":"123456","newPassword":"` + pw + `"}`
		if rec := postJSON(t, mux, "/api/auth/recovery/complete", body); !strings.Contains(rec.Body.String(), "INVALID_PASSWORD") {
			t.Fatalf("password %q len %d body = %s, want INVALID_PASSWORD", pw, len(pw), rec.Body.String())
		}
	}
	if len(repo.completed) != 0 {
		t.Fatalf("baseline violations must not complete: %v", repo.completed)
	}
}

// A-001 F-001: complete failures feed the IP|identifier limiter bucket, so a
// known-account brute force hits 429 RATE_LIMITED even though the per-challenge
// attempt cap (≤5) would otherwise be the only brake. The bucket fills the same
// way for existing and non-existing accounts — no existence oracle leaks.
func TestRecoveryCompleteRateLimitedAfterTwentyFailures(t *testing.T) {
	repo := &fakeRecoveryRepo{target: &authsession.RecoveryTarget{UserID: "u1"}, outcome: authsession.RecoveryMismatch}
	mux := newRecoveryMux(repo, nil)
	const budget = 20 // this handler's loginRateLimiter max
	for i := 0; i < budget; i++ {
		rec := postJSON(t, mux, "/api/auth/recovery/complete", `{"account":"x","code":"000000","newPassword":"new-password-1"}`)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("failure %d status = %d (%s), want 400", i+1, rec.Code, rec.Body.String())
		}
	}
	rec := postJSON(t, mux, "/api/auth/recovery/complete", `{"account":"x","code":"000000","newPassword":"new-password-1"}`)
	if rec.Code != http.StatusTooManyRequests || !strings.Contains(rec.Body.String(), "RATE_LIMITED") {
		t.Fatalf("post-budget status = %d body %s, want 429 RATE_LIMITED", rec.Code, rec.Body.String())
	}
}

// A-001 F-001 corollary: the second-factor DEMAND (missing field — not a
// guess) and INVALID_PASSWORD after a matched code consume neither bucket nor
// challenge budget; a legitimate user mid-flow must not lock themselves out.
func TestRecoveryCompleteNonGuessFailuresDoNotRecord(t *testing.T) {
	repo := &fakeRecoveryRepo{target: &authsession.RecoveryTarget{UserID: "u1"}, outcome: authsession.RecoveryMatch}
	gate := &fakeGate{required: true}
	mux := newRecoveryMux(repo, gate)
	for i := 0; i < 25; i++ {
		rec := postJSON(t, mux, "/api/auth/recovery/complete", `{"account":"x","code":"123456","secondFactorCode":"654321","newPassword":"short"}`)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("iteration %d status = %d, want 400 (limiter must stay idle)", i+1, rec.Code)
		}
		if !strings.Contains(rec.Body.String(), "INVALID_PASSWORD") {
			t.Fatalf("iteration %d body = %s, want INVALID_PASSWORD", i+1, rec.Body.String())
		}
	}
	if len(repo.completed) != 0 {
		t.Fatalf("weak-password completes must never land: %v", repo.completed)
	}
}
