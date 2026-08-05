package authsession

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

func openRepository(t *testing.T, name string, seed bool) (*Repository, *store.Store) {
	t.Helper()
	st, err := testsupport.OpenStore(filepath.Join(t.TempDir(), name), "admin", "hash-v1", seed)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return NewRepository(st), st
}

func repositoryQueryInt(t *testing.T, st *store.Store, query string, args ...any) int {
	t.Helper()
	var value int
	if err := st.WithTx(context.Background(), func(tx *sql.Tx) error {
		return tx.QueryRow(query, args...).Scan(&value)
	}); err != nil {
		t.Fatalf("query %q: %v", query, err)
	}
	return value
}

func repositoryExec(t *testing.T, st *store.Store, query string, args ...any) error {
	t.Helper()
	return st.WithTx(context.Background(), func(tx *sql.Tx) error {
		_, err := tx.Exec(query, args...)
		return err
	})
}

func TestCreateUserAndLookup(t *testing.T) {
	repository, _ := openRepository(t, "accounts.db", false)
	now := time.Now().UTC()
	if err := repository.CreateUser(User{
		ID: "u1", Username: "alice", Name: "Alice", Roles: []string{"editor"},
		PasswordHash: "h", CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	got, err := repository.UserByUsername("alice")
	if err != nil {
		t.Fatalf("UserByUsername: %v", err)
	}
	if got.ID != "u1" || len(got.Roles) != 1 || got.Roles[0] != "editor" {
		t.Fatalf("got = %+v, want u1/editor", got)
	}
	if _, err := repository.UserByUsername("nobody"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("UserByUsername(nobody) = %v, want ErrNotFound", err)
	}
}

func TestRefreshTokenLifecycle(t *testing.T) {
	repository, _ := openRepository(t, "refresh.db", true)
	now := time.Now().UTC()
	token := RefreshToken{
		ID: "rt1", UserID: "user-admin", TokenHash: "abc",
		ExpiresAt: now.Add(time.Hour), CreatedAt: now,
	}
	if err := repository.CreateRefreshToken(token); err != nil {
		t.Fatalf("CreateRefreshToken: %v", err)
	}
	got, err := repository.RefreshTokenByHash("abc")
	if err != nil || got.RevokedAt != nil {
		t.Fatalf("RefreshTokenByHash = %+v, err %v", got, err)
	}
	if _, err := repository.RefreshTokenByHash("missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing token = %v, want ErrNotFound", err)
	}
	if err := repository.RevokeRefreshToken("rt1", now.Add(time.Minute)); err != nil {
		t.Fatalf("RevokeRefreshToken: %v", err)
	}
	got, err = repository.RefreshTokenByHash("abc")
	if err != nil || got.RevokedAt == nil {
		t.Fatalf("revoked token = %+v, err %v", got, err)
	}
	if err := repository.RevokeRefreshToken("rt1", now.Add(2*time.Minute)); !errors.Is(err, ErrAlreadyRevoked) {
		t.Fatalf("second revoke = %v, want ErrAlreadyRevoked", err)
	}
	if err := repository.RevokeRefreshToken("missing", now); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing revoke = %v, want ErrNotFound", err)
	}
}
