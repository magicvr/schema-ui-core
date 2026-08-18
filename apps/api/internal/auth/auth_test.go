package auth

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	authsession "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

func newTestAuth(t *testing.T, devSession bool) *Authenticator {
	t.Helper()
	hash, err := HashPassword("pw", 4)
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	st, err := testsupport.OpenStore(filepath.Join(t.TempDir(), "test.db"), "admin", hash, true)
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
	if sub.UserID != "user-admin" {
		t.Fatalf("subject = %v, want user-admin", sub.UserID)
	}
}

func TestServiceCredentialMiddlewarePrecedesDevFallback(t *testing.T) {
	hash, err := HashPassword("pw", 4)
	if err != nil {
		t.Fatal(err)
	}
	st, err := testsupport.OpenStore(filepath.Join(t.TempDir(), "service-auth.db"), "admin", hash, true)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	repository := authsession.NewRepository(st)
	a := NewWithRepository([]byte("secret"), 15*time.Minute, 30*24*time.Hour, repository, true)
	raw, tokenHash, prefix, err := NewServiceCredentialToken()
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	credential := authsession.ServiceCredential{
		ID: "0123456789abcdef0123456789abcdef", Name: "Build Agent",
		TokenPrefix: prefix, TokenHash: tokenHash, Scopes: []string{"records.read"},
		ExpiresAt: now.Add(time.Hour), CreatedBy: "user-admin", CreatedAt: now, UpdatedAt: now,
	}
	if err := repository.CreateServiceCredential(credential, nil); err != nil {
		t.Fatal(err)
	}

	called := false
	protected := a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		identity, ok := IdentityFrom(r.Context())
		if !ok || identity.ID != "service-credential:"+credential.ID || identity.CredentialID != credential.ID || !identity.IsServiceCredential() {
			t.Fatalf("identity = %+v, ok=%v", identity, ok)
		}
		if _, ok := UserIdentityFrom(r.Context()); ok {
			t.Fatal("service principal passed UserIdentityFrom")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	request := httptest.NewRequest(http.MethodGet, "/api/resources/widgets", nil)
	request.Header.Set("Authorization", "Bearer "+raw)
	response := httptest.NewRecorder()
	protected.ServeHTTP(response, request)
	if response.Code != http.StatusNoContent || !called {
		t.Fatalf("valid credential response=%d called=%v body=%s", response.Code, called, response.Body.String())
	}
	got, err := repository.ServiceCredentialByID(credential.ID)
	if err != nil || got.LastUsedAt == nil {
		t.Fatalf("last used credential = %+v, err=%v", got, err)
	}

	for name, token := range map[string]string{
		"unknown": serviceCredentialPrefix + "unknown",
		"revoked": raw,
	} {
		if name == "revoked" {
			if _, _, err := repository.RevokeServiceCredential(credential.ID, now, nil); err != nil {
				t.Fatal(err)
			}
		}
		called = false
		request = httptest.NewRequest(http.MethodGet, "/api/resources/widgets", nil)
		request.Header.Set("Authorization", "Bearer "+token)
		response = httptest.NewRecorder()
		protected.ServeHTTP(response, request)
		if response.Code != http.StatusUnauthorized || called {
			t.Fatalf("%s service credential response=%d called=%v body=%s", name, response.Code, called, response.Body.String())
		}
	}
	expiredRaw, expiredHash, expiredPrefix, err := NewServiceCredentialToken()
	if err != nil {
		t.Fatal(err)
	}
	if err := repository.CreateServiceCredential(authsession.ServiceCredential{
		ID: "fedcba9876543210fedcba9876543210", Name: "Expired Agent",
		TokenPrefix: expiredPrefix, TokenHash: expiredHash, Scopes: []string{"records.read"},
		ExpiresAt: now.Add(-time.Minute), CreatedBy: "user-admin", CreatedAt: now, UpdatedAt: now,
	}, nil); err != nil {
		t.Fatal(err)
	}
	called = false
	request = httptest.NewRequest(http.MethodGet, "/api/resources/widgets", nil)
	request.Header.Set("Authorization", "Bearer "+expiredRaw)
	response = httptest.NewRecorder()
	protected.ServeHTTP(response, request)
	if response.Code != http.StatusUnauthorized || called {
		t.Fatalf("expired service credential response=%d called=%v body=%s", response.Code, called, response.Body.String())
	}
}

func TestNewServiceCredentialTokenContract(t *testing.T) {
	raw, hash, prefix, err := NewServiceCredentialToken()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(raw, serviceCredentialPrefix) || len(raw) != len(serviceCredentialPrefix)+43 {
		t.Fatalf("raw shape = %q", raw)
	}
	if prefix != raw[:15] || hash != HashToken(raw) || len(hash) != 64 {
		t.Fatalf("hash/prefix contract failed: hash=%q prefix=%q", hash, prefix)
	}
	if id := NewServiceCredentialID(); len(id) != 32 {
		t.Fatalf("credential id length = %d, want 32", len(id))
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

// Refresh must not mint a new pair for a locked account even if a refresh
// token is still live (revoke-on-lock is best-effort).
func TestRefreshRejectsLockedAccount(t *testing.T) {
	hash, err := HashPassword("pw", 4)
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	st, err := testsupport.OpenStore(filepath.Join(t.TempDir(), "locked-refresh.db"), "admin", hash, true)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	a := New([]byte("secret"), 15*time.Minute, 30*24*time.Hour, st, false)

	_, refresh, _, err := a.Login("admin", "pw", now())
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	lockUntil := now().Add(time.Hour).Unix()
	if err := st.WithTx(context.Background(), func(tx *sql.Tx) error {
		_, err := tx.Exec(`UPDATE users SET locked_until = ? WHERE username = 'admin'`, lockUntil)
		return err
	}); err != nil {
		t.Fatalf("lock user: %v", err)
	}

	if _, _, _, err := a.Refresh(refresh, now()); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("refresh locked account = %v, want ErrInvalidToken", err)
	}
}

// C2 hardening: two concurrent rotations of the same refresh token must yield
// exactly one new session pair — the second caller loses the atomic guarded
// revoke and fails with ErrTokenRevoked instead of issuing a second live pair.
func TestRefreshConcurrentRotationSingleWinner(t *testing.T) {
	a := newTestAuth(t, false)
	_, refresh, _, err := a.Login("admin", "pw", now())
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	start := make(chan struct{})
	results := make(chan error, 2)
	for range 2 {
		go func() {
			<-start
			_, _, _, err := a.Refresh(refresh, now().Add(time.Minute))
			results <- err
		}()
	}
	close(start)
	var winners, revoked int
	for range 2 {
		err := <-results
		switch {
		case err == nil:
			winners++
		case errors.Is(err, ErrTokenRevoked):
			revoked++
		default:
			t.Fatalf("unexpected refresh error: %v", err)
		}
	}
	if winners != 1 || revoked != 1 {
		t.Fatalf("concurrent rotation: winners=%d revoked=%d, want 1/1", winners, revoked)
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
	if uid, err := a.Logout(refresh, now().Add(time.Minute)); err != nil {
		t.Fatalf("Logout: %v", err)
	} else if uid != "user-admin" {
		t.Fatalf("Logout user id = %q, want user-admin", uid)
	}
	// Idempotent: logging out the same token again is a no-op success.
	if _, err := a.Logout(refresh, now().Add(2*time.Minute)); err != nil {
		t.Fatalf("second Logout = %v, want nil", err)
	}
	if _, _, _, err := a.Refresh(refresh, now().Add(3*time.Minute)); !errors.Is(err, ErrTokenRevoked) {
		t.Fatalf("refresh after logout = %v, want ErrTokenRevoked", err)
	}
}

func TestParseAccessTokenExpiredAndWrongSecret(t *testing.T) {
	// A token minted with a negative TTL is already expired at signing.
	expired, err := SignAccessToken([]byte("secret"), "user-admin", 0, -time.Minute, now())
	if err != nil {
		t.Fatalf("SignAccessToken: %v", err)
	}
	if _, err := ParseAccessToken([]byte("secret"), expired); err == nil {
		t.Fatalf("ParseAccessToken(expired) = nil, want error")
	}
	// A token signed with a different secret must be rejected.
	other, err := SignAccessToken([]byte("other"), "user-admin", 0, time.Minute, now())
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

// buildLegacyR2DB creates a pre-migration R2 database (users + refresh_tokens,
// no schema_migrations) so the compiled-catalog store open has to run the 0001 fingerprint + 0002
// backfill path. users holds [id, username, name, rolesJSON, passwordHash].
func buildLegacyR2DB(t *testing.T, path string, users [][]string) {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("open legacy db: %v", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	ddl := []string{
		`CREATE TABLE users (
  id            TEXT PRIMARY KEY,
  username      TEXT NOT NULL UNIQUE,
  name          TEXT NOT NULL,
  roles         TEXT NOT NULL,
  password_hash TEXT NOT NULL,
  created_at    INTEGER NOT NULL,
  updated_at    INTEGER NOT NULL
)`,
		`CREATE TABLE refresh_tokens (
  id         TEXT PRIMARY KEY,
  user_id    TEXT NOT NULL REFERENCES users(id),
  token_hash TEXT NOT NULL UNIQUE,
  expires_at INTEGER NOT NULL,
  revoked_at INTEGER,
  created_at INTEGER NOT NULL
)`,
		`CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id)`,
	}
	for _, s := range ddl {
		if _, err := db.Exec(s); err != nil {
			t.Fatalf("legacy ddl: %v", err)
		}
	}
	now := time.Now().UTC().Unix()
	for _, u := range users {
		if _, err := db.Exec(
			`INSERT INTO users (id, username, name, roles, password_hash, created_at, updated_at)
			 VALUES (?,?,?,?,?,?,?)`,
			u[0], u[1], u[2], u[3], u[4], now, now,
		); err != nil {
			t.Fatalf("legacy insert %s: %v", u[0], err)
		}
	}
}

// A-002 F-004 · a legacy R2 user whose roles JSON contains duplicates survives
// migration and can authenticate: 0002 dedupes the relations and the read
// comparator follows set semantics. Login exercises UserByUsername and Refresh
// exercises UserByID on the migrated user.
func TestLoginAndRefreshAfterMigrateDuplicateRoles(t *testing.T) {
	path := filepath.Join(t.TempDir(), "legacy-dup.db")
	hash, err := HashPassword("pw", 4)
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	buildLegacyR2DB(t, path, [][]string{{"u-alice", "alice", "Alice", `["admin","admin","editor"]`, hash}})

	st, err := testsupport.OpenStore(path, "admin", "hash", false)
	if err != nil {
		t.Fatalf("open migrated store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	a := New([]byte("secret"), 15*time.Minute, testRefreshTTL, st, false)

	access, refresh, user, err := a.Login("alice", "pw", now())
	if err != nil {
		t.Fatalf("Login after migration with duplicate legacy roles: %v", err)
	}
	if want := []string{"admin", "editor"}; !reflect.DeepEqual(user.Roles, want) {
		t.Fatalf("roles = %v, want %v (deduped, sorted by key)", user.Roles, want)
	}
	sub, err := ParseAccessToken([]byte("secret"), access)
	if err != nil {
		t.Fatalf("ParseAccessToken: %v", err)
	}
	if sub.UserID != "u-alice" {
		t.Fatalf("subject = %q, want u-alice", sub.UserID)
	}
	if _, _, _, err := a.Refresh(refresh, now().Add(time.Minute)); err != nil {
		t.Fatalf("Refresh after migration: %v", err)
	}
}
