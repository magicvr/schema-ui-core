package authsession

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

func TestCreateUserDoubleWritesRoles(t *testing.T) {
	repository, st := openRepository(t, "doublewrite.db", false)
	now := time.Now().UTC()
	if err := repository.CreateUser(User{
		ID: "u1", Username: "alice", Name: "Alice",
		Roles: []string{"editor", "viewer"}, PasswordHash: "h", CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}

	var rolesJSON string
	if err := st.WithTx(context.Background(), func(tx *sql.Tx) error {
		return tx.QueryRow(`SELECT roles FROM users WHERE id = 'u1'`).Scan(&rolesJSON)
	}); err != nil {
		t.Fatal(err)
	}
	if rolesJSON != `["editor","viewer"]` {
		t.Fatalf("legacy roles = %s", rolesJSON)
	}
	if roles := repositoryQueryInt(t, st, `SELECT COUNT(*) FROM roles`); roles != 2 {
		t.Fatalf("roles = %d, want 2", roles)
	}
	if links := repositoryQueryInt(t, st, `SELECT COUNT(*) FROM user_roles WHERE user_id = 'u1'`); links != 2 {
		t.Fatalf("user_roles = %d, want 2", links)
	}
	user, err := repository.UserByUsername("alice")
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"editor", "viewer"}; !reflect.DeepEqual(user.Roles, want) {
		t.Fatalf("roles = %v, want %v", user.Roles, want)
	}
}

func TestNormalizedReadSortedAndDetectsMismatch(t *testing.T) {
	repository, st := openRepository(t, "normalized.db", false)
	now := time.Now().UTC()
	if err := repository.CreateUser(User{
		ID: "u1", Username: "alice", Name: "Alice", Roles: []string{"viewer", "admin"},
		PasswordHash: "h", CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatal(err)
	}
	user, err := repository.UserByID("u1")
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"admin", "viewer"}; !reflect.DeepEqual(user.Roles, want) {
		t.Fatalf("roles = %v, want %v", user.Roles, want)
	}

	if err := repositoryExec(t, st, `DELETE FROM user_roles WHERE user_id = 'u1'`); err != nil {
		t.Fatal(err)
	}
	if _, err := repository.UserByUsername("alice"); err == nil {
		t.Fatal("expected role mismatch after normalized relation removal")
	}
	if _, err := repository.UserByID("u1"); err == nil {
		t.Fatal("expected role mismatch from id lookup")
	}
}

func TestBootstrapRelationsStayIdempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "seed.db")
	st, err := testsupport.OpenStore(path, "admin", "hash-v1", true)
	if err != nil {
		t.Fatal(err)
	}
	repository := NewRepository(st)
	if links := repositoryQueryInt(t, st, `SELECT COUNT(*) FROM user_roles WHERE user_id = 'user-admin'`); links != 2 {
		t.Fatalf("seed links = %d, want 2", links)
	}
	user, err := repository.UserByUsername("admin")
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"admin", "editor"}; !reflect.DeepEqual(user.Roles, want) {
		t.Fatalf("seed roles = %v, want %v", user.Roles, want)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	st, err = testsupport.OpenStore(path, "admin", "hash-v2", true)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	repository = NewRepository(st)
	user, err = repository.UserByUsername("admin")
	if err != nil {
		t.Fatal(err)
	}
	if user.PasswordHash != "hash-v1" {
		t.Fatalf("password = %q, want original hash-v1", user.PasswordHash)
	}
	if links := repositoryQueryInt(t, st, `SELECT COUNT(*) FROM user_roles WHERE user_id = 'user-admin'`); links != 2 {
		t.Fatalf("reopen links = %d, want 2", links)
	}
}

func TestUserRolesForeignKeysAndCascade(t *testing.T) {
	repository, st := openRepository(t, "fk.db", false)
	now := time.Now().UTC()
	if err := repository.CreateUser(User{
		ID: "u2", Username: "bob", Name: "Bob", Roles: []string{"editor"},
		PasswordHash: "h", CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatal(err)
	}

	err := st.WithTx(context.Background(), func(tx *sql.Tx) error {
		if _, err := tx.Exec(`INSERT INTO user_roles (user_id, role_id) VALUES ('u2', 'role-nope')`); err == nil {
			return errors.New("expected FK violation for unknown role")
		}
		if _, err := tx.Exec(`DELETE FROM roles WHERE id = 'role-editor'`); err == nil {
			return errors.New("expected RESTRICT violation for in-use role")
		}
		_, err := tx.Exec(`DELETE FROM users WHERE id = 'u2'`)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	if links := repositoryQueryInt(t, st, `SELECT COUNT(*) FROM user_roles WHERE user_id = 'u2'`); links != 0 {
		t.Fatalf("links after user delete = %d, want 0", links)
	}
}

var _ TxRunner = (*store.Store)(nil)
