package auth

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

// GOAL-014 D-002 (W13 F-007 targeted-DoS fix) regression locks for the
// layered login-lockout model:
//
//  1. Per-(account|source) locks deny only the abusive source — the
//     legitimate user keeps logging in from elsewhere.
//  2. The GLOBAL ceiling (100 consecutive failures, 24h sliding restart)
//     still brakes distributed guessing and fires OnLockOpened once.
//  3. Login failures NEVER revoke refresh tokens anymore: after a
//     failure-driven global lock EXPIRES, the pre-lock refresh token still
//     rotates (forced logout was the weaponizable edge of the old model).

func failuresFrom(t *testing.T, a *Authenticator, ip string, n int, start time.Time) {
	t.Helper()
	for i := 0; i < n; i++ {
		_, _, _, err := a.Login("admin", "wrong", start.Add(time.Duration(i)*time.Second), ip)
		if !errors.Is(err, ErrInvalidCredentials) {
			t.Fatalf("failure %d from %s = %v, want ErrInvalidCredentials", i+1, ip, err)
		}
	}
}

func TestLoginSourceScopedLockout(t *testing.T) {
	a := newTestAuth(t, false)
	start := now()

	failuresFrom(t, a, "10.0.0.1", IPSourceLockThreshold, start)

	// The abusive source is locked.
	if _, _, _, err := a.Login("admin", "pw", start.Add(time.Minute), "10.0.0.1"); !errors.Is(err, ErrAccountLocked) {
		t.Fatalf("locked-source correct password = %v, want ErrAccountLocked", err)
	}
	// Even a WRONG password from another source stays a plain credential
	// failure (its own bucket is empty).
	if _, _, _, err := a.Login("admin", "wrong", start.Add(time.Minute), "10.0.0.2"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("fresh source wrong password = %v, want ErrInvalidCredentials", err)
	}
	// THE FIX: the legitimate user logs in fine from their own source while
	// the attacker's source sits locked.
	access, _, _, err := a.Login("admin", "pw", start.Add(2*time.Minute), "10.0.0.2")
	if err != nil {
		t.Fatalf("legitimate login from unaffected source = %v, want success", err)
	}
	if access == "" {
		t.Fatal("expected an access token")
	}
}

func TestLoginGlobalCeilingBrakesDistributedGuessing(t *testing.T) {
	a := newTestAuth(t, false)
	opened := 0
	a.OnLockOpened = func(string) { opened++ }
	start := now()

	// Spread 100 failures over many sources: no single pair reaches its own
	// threshold, but the account-wide ceiling must open the global lock and
	// fire the admin-visibility hook exactly once.
	for i := range LockThresholdFailures {
		ip := fmt.Sprintf("10.%d.%d.%d", 1+i/250%250, (i/8)%250, i%8+2)
		_, _, _, err := a.Login("admin", "wrong", start.Add(time.Duration(i)*time.Second), ip)
		if !errors.Is(err, ErrInvalidCredentials) {
			t.Fatalf("distributed failure %d (%s) = %v, want ErrInvalidCredentials", i+1, ip, err)
		}
	}
	// A brand-new source is denied by the GLOBAL lock...
	if _, _, _, err := a.Login("admin", "pw", start.Add(2*time.Minute), "10.9.9.9"); !errors.Is(err, ErrAccountLocked) {
		t.Fatalf("post-ceiling fresh-source login = %v, want ErrAccountLocked", err)
	}
	if opened != 1 {
		t.Fatalf("OnLockOpened fired %d times, want exactly 1 (global lock open)", opened)
	}

	// ...and it self-heals: past the lock window the same source succeeds,
	// with both counters reset.
	recovered := start.Add(2*time.Minute + LockWindow + time.Minute)
	if _, _, _, err := a.Login("admin", "pw", recovered, "10.9.9.9"); err != nil {
		t.Fatalf("login after lock window = %v, want success", err)
	}
}

func TestLoginFailureNoLongerRevokesRefreshTokens(t *testing.T) {
	a := newTestAuth(t, false)
	start := now()

	// Two sessions = two live refresh tokens ("two devices").
	_, refreshA, _, err := a.Login("admin", "pw", start)
	if err != nil {
		t.Fatalf("Login A: %v", err)
	}
	_, refreshB, _, err := a.Login("admin", "pw", start.Add(time.Second))
	if err != nil {
		t.Fatalf("Login B: %v", err)
	}

	// Drive the GLOBAL lock open through real failed logins.
	for i := range LockThresholdFailures {
		ip := fmt.Sprintf("10.7.%d.%d", i/125, i%125+3)
		if _, _, _, err := a.Login("admin", "wrong", start.Add(time.Duration(i)*time.Second), ip); !errors.Is(err, ErrInvalidCredentials) {
			t.Fatalf("setup failure %d = %v", i+1, err)
		}
	}

	// Presenting token A while the global lock window is open: denied
	// (locked-account contract) — and consumed by the pre-existing
	// rotate-before-checks semantics of Refresh itself.
	if _, _, _, err := a.Refresh(refreshA, start.Add(2*time.Minute)); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("refresh during global lock = %v, want ErrInvalidToken", err)
	}

	// THE REMOVAL UNDER TEST: token B was never presented during the lock and
	// must STILL be live after the window — the old failure-driven
	// RevokeAllRefreshTokensForUser would have blanket-revoked it at lock
	// open (the forced-logout weapon).
	if _, _, _, err := a.Refresh(refreshB, start.Add(2*time.Minute+LockWindow+time.Minute)); err != nil {
		t.Fatalf("untouched refresh token after lock expiry = %v, want success (failures must not revoke)", err)
	}
}
