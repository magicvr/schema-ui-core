package store

import (
	"path/filepath"
	"testing"
	"time"
)

func TestSeedAdminIdempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")

	st, err := Open(path, "admin", "hash-v1", true)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	u, err := st.UserByUsername("admin")
	if err != nil {
		t.Fatalf("UserByUsername: %v", err)
	}
	if u.PasswordHash != "hash-v1" {
		t.Fatalf("password_hash = %q, want hash-v1", u.PasswordHash)
	}
	if len(u.Roles) != 2 || u.Roles[0] != "admin" {
		t.Fatalf("roles = %v, want [admin editor]", u.Roles)
	}
	if err := st.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	// Reopening with a different seed must not duplicate or overwrite.
	st2, err := Open(path, "admin", "hash-v2", true)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer st2.Close()
	u2, err := st2.UserByUsername("admin")
	if err != nil {
		t.Fatalf("UserByUsername after reopen: %v", err)
	}
	if u2.PasswordHash != "hash-v1" {
		t.Fatalf("password_hash = %q after reopen, want hash-v1 (seed must be no-op)", u2.PasswordHash)
	}
}

func TestCreateUserAndLookup(t *testing.T) {
	st, err := Open(filepath.Join(t.TempDir(), "test.db"), "admin", "hash", false)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer st.Close()

	now := time.Now().UTC()
	if err := st.CreateUser(User{ID: "u1", Username: "alice", Name: "Alice", Roles: []string{"editor"}, PasswordHash: "h", CreatedAt: now, UpdatedAt: now}); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	got, err := st.UserByUsername("alice")
	if err != nil {
		t.Fatalf("UserByUsername: %v", err)
	}
	if got.ID != "u1" || len(got.Roles) != 1 || got.Roles[0] != "editor" {
		t.Fatalf("got = %+v, want u1/editor", got)
	}
	if _, err := st.UserByUsername("nobody"); err != ErrNotFound {
		t.Fatalf("UserByUsername(nobody) = %v, want ErrNotFound", err)
	}
}

func TestRefreshTokenLifecycle(t *testing.T) {
	st, err := Open(filepath.Join(t.TempDir(), "test.db"), "admin", "hash", true)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer st.Close()

	now := time.Now().UTC()
	rt := RefreshToken{ID: "rt1", UserID: "user-admin", TokenHash: "abc", ExpiresAt: now.Add(time.Hour), CreatedAt: now}
	if err := st.CreateRefreshToken(rt); err != nil {
		t.Fatalf("CreateRefreshToken: %v", err)
	}
	got, err := st.RefreshTokenByHash("abc")
	if err != nil {
		t.Fatalf("RefreshTokenByHash: %v", err)
	}
	if got.RevokedAt != nil {
		t.Fatalf("revoked_at = %v, want nil", got.RevokedAt)
	}
	if _, err := st.RefreshTokenByHash("missing"); err != ErrNotFound {
		t.Fatalf("RefreshTokenByHash(missing) = %v, want ErrNotFound", err)
	}

	// Revoke once succeeds; second revoke is a no-op result.
	if err := st.RevokeRefreshToken("rt1", now.Add(time.Minute)); err != nil {
		t.Fatalf("RevokeRefreshToken: %v", err)
	}
	got2, _ := st.RefreshTokenByHash("abc")
	if got2.RevokedAt == nil {
		t.Fatalf("revoked_at = nil, want set")
	}
	if err := st.RevokeRefreshToken("rt1", now.Add(2*time.Minute)); err != ErrAlreadyRevoked {
		t.Fatalf("second revoke = %v, want ErrAlreadyRevoked", err)
	}
	if err := st.RevokeRefreshToken("rt-missing", now); err != ErrNotFound {
		t.Fatalf("revoke missing = %v, want ErrNotFound", err)
	}
}
