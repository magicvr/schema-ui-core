package auth

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
)

func newTestAuth(t *testing.T, devSession bool) *Authenticator {
	t.Helper()
	hash, err := HashPassword("pw", 4)
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	st, err := store.Open(filepath.Join(t.TempDir(), "test.db"), "admin", hash, true)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return New([]byte("secret"), 15*time.Minute, 30*24*time.Hour, st, devSession)
}

const testRefreshTTL = 30 * 24 * time.Hour

func now() time.Time {
	return time.Now().UTC()
}

func TestLoginSuccess(t *testing.T) {
	a := newTestAuth(t, false)
	access, refresh, user, err := a.Login("admin", "pw", now())
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if access == "" || refresh == "" {
		t.Fatalf("empty tokens: access=%q refresh=%q", access, refresh)
	}
	if user.ID != "user-admin" {
		t.Fatalf("user.id = %v, want user-admin", user.ID)
	}

	sub, err := ParseAccessToken([]byte("secret"), access)
	if err != nil {
		t.Fatalf("ParseAccessToken: %v", err)
	}
	if sub != "user-admin" {
		t.Fatalf("subject = %v, want user-admin", sub)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	a := newTestAuth(t, false)
	_, _, _, err := a.Login("admin", "wrong", now())
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("err = %v, want ErrInvalidCredentials", err)
	}
}

func TestLoginUnknownUser(t *testing.T) {
	a := newTestAuth(t, false)
	_, _, _, err := a.Login("nobody", "pw", now())
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("err = %v, want ErrInvalidCredentials (no user enumeration)", err)
	}
}

func TestRefreshRotatesAndRevokesOld(t *testing.T) {
	a := newTestAuth(t, false)
	_, oldRefresh, _, err := a.Login("admin", "pw", now())
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	_, newRefresh, _, err := a.Refresh(oldRefresh, now().Add(time.Minute))
	if err != nil {
		t.Fatalf("Refresh: %v", err)
	}
	if newRefresh == "" || newRefresh == oldRefresh {
		t.Fatalf("new refresh = %q, want non-empty and different", newRefresh)
	}

	// The old token is now revoked and must be rejected.
	if _, _, _, err := a.Refresh(oldRefresh, now().Add(2*time.Minute)); !errors.Is(err, ErrTokenRevoked) {
		t.Fatalf("reuse of old refresh = %v, want ErrTokenRevoked", err)
	}
}

func TestRefreshUnknownToken(t *testing.T) {
	a := newTestAuth(t, false)
	if _, _, _, err := a.Refresh("bogus", now()); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("err = %v, want ErrInvalidToken", err)
	}
}

func TestRefreshExpired(t *testing.T) {
	a := newTestAuth(t, false)
	_, refresh, _, err := a.Login("admin", "pw", now())
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	later := now().Add(30*24*time.Hour + time.Second)
	if _, _, _, err := a.Refresh(refresh, later); !errors.Is(err, ErrExpiredToken) {
		t.Fatalf("err = %v, want ErrExpiredToken", err)
	}
}

func TestLogoutRevokes(t *testing.T) {
	a := newTestAuth(t, false)
	_, refresh, _, err := a.Login("admin", "pw", now())
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if err := a.Logout(refresh, now().Add(time.Minute)); err != nil {
		t.Fatalf("Logout: %v", err)
	}
	// Idempotent: logging out the same token again is a no-op success.
	if err := a.Logout(refresh, now().Add(2*time.Minute)); err != nil {
		t.Fatalf("second Logout = %v, want nil", err)
	}
	if _, _, _, err := a.Refresh(refresh, now().Add(3*time.Minute)); !errors.Is(err, ErrTokenRevoked) {
		t.Fatalf("refresh after logout = %v, want ErrTokenRevoked", err)
	}
}

func TestParseAccessTokenExpiredAndWrongSecret(t *testing.T) {
	// A token minted with a negative TTL is already expired at signing.
	expired, err := SignAccessToken([]byte("secret"), "user-admin", -time.Minute, now())
	if err != nil {
		t.Fatalf("SignAccessToken: %v", err)
	}
	if _, err := ParseAccessToken([]byte("secret"), expired); err == nil {
		t.Fatalf("ParseAccessToken(expired) = nil, want error")
	}
	// A token signed with a different secret must be rejected.
	other, err := SignAccessToken([]byte("other"), "user-admin", time.Minute, now())
	if err != nil {
		t.Fatalf("SignAccessToken: %v", err)
	}
	if _, err := ParseAccessToken([]byte("secret"), other); err == nil {
		t.Fatalf("ParseAccessToken(other secret) = nil, want error")
	}
}

func TestOpaqueTokenHashStable(t *testing.T) {
	raw, hash, err := NewOpaqueToken()
	if err != nil {
		t.Fatalf("NewOpaqueToken: %v", err)
	}
	if raw == "" || hash == "" {
		t.Fatalf("empty token: raw=%q hash=%q", raw, hash)
	}
	if got := HashToken(raw); got != hash {
		t.Fatalf("HashToken(raw) = %q, want %q", got, hash)
	}
}
