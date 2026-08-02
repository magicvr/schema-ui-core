package store

import (
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

// V-MIG-05 · CreateUser double-writes the legacy JSON and the normalized
// user_roles relation in one transaction; reads return the normalized source.
func TestCreateUserDoubleWritesRoles(t *testing.T) {
	st, err := Open(filepath.Join(t.TempDir(), "doublewrite.db"), "admin", "hash", false)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer st.Close()

	now := time.Now().UTC()
	if err := st.CreateUser(User{
		ID: "u1", Username: "alice", Name: "Alice",
		Roles: []string{"editor", "viewer"}, PasswordHash: "h", CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}

	var rolesJSON string
	if err := st.db.QueryRow(`SELECT roles FROM users WHERE id = 'u1'`).Scan(&rolesJSON); err != nil {
		t.Fatal(err)
	}
	if rolesJSON != `["editor","viewer"]` {
		t.Fatalf("legacy roles = %s", rolesJSON)
	}
	var roleCount, urCount int
	st.db.QueryRow(`SELECT COUNT(*) FROM roles`).Scan(&roleCount)
	st.db.QueryRow(`SELECT COUNT(*) FROM user_roles WHERE user_id = 'u1'`).Scan(&urCount)
	if roleCount != 2 || urCount != 2 {
		t.Fatalf("double-write roles=%d user_roles=%d, want 2/2", roleCount, urCount)
	}
	var key string
	if err := st.db.QueryRow(`SELECT key FROM roles WHERE id = 'role-viewer'`).Scan(&key); err != nil || key != "viewer" {
		t.Fatalf("role-viewer key = %q, err %v", key, err)
	}

	u, err := st.UserByUsername("alice")
	if err != nil {
		t.Fatalf("UserByUsername: %v", err)
	}
	if want := []string{"editor", "viewer"}; !reflect.DeepEqual(u.Roles, want) {
		t.Fatalf("roles = %v, want %v", u.Roles, want)
	}
}

// V-MIG-05 · normalized output order is deterministic: read roles are sorted by
// key ascending regardless of the legacy JSON order.
func TestNormalizedReadSortedByKey(t *testing.T) {
	st, err := Open(filepath.Join(t.TempDir(), "sorted.db"), "admin", "hash", false)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer st.Close()

	now := time.Now().UTC()
	if err := st.CreateUser(User{
		ID: "u1", Username: "alice", Name: "Alice",
		Roles:        []string{"viewer", "admin"}, // intentionally unsorted
		PasswordHash: "h", CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}

	u, err := st.UserByID("u1")
	if err != nil {
		t.Fatalf("UserByID: %v", err)
	}
	if want := []string{"admin", "viewer"}; !reflect.DeepEqual(u.Roles, want) {
		t.Fatalf("roles = %v, want %v (ascending by key)", u.Roles, want)
	}
}

// V-MIG-05 · a user whose legacy JSON and normalized relation diverge errors
// loudly instead of returning a silently inconsistent identity.
func TestReadDetectsRoleMismatch(t *testing.T) {
	st, err := Open(filepath.Join(t.TempDir(), "mismatch.db"), "admin", "hash", false)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer st.Close()

	now := time.Now().UTC().Unix()
	// Simulate a pre-S2 gap: legacy JSON present, no normalized user_roles rows.
	if _, err := st.db.Exec(
		`INSERT INTO users (id, username, name, roles, password_hash, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		"u1", "alice", "Alice", `["admin"]`, "h", now, now,
	); err != nil {
		t.Fatalf("insert legacy-only user: %v", err)
	}
	if _, err := st.UserByUsername("alice"); err == nil {
		t.Fatal("expected role-mismatch error for legacy-only user")
	}

	// The same gap on the primary-key lookup also fails closed.
	if _, err := st.UserByID("u1"); err == nil {
		t.Fatal("expected role-mismatch error from UserByID")
	}
}

// seedAdmin double-writes the normalized relations, closing the S1 fresh-DB
// intermediate state; reopening stays idempotent.
func TestSeedAdminDoubleWrites(t *testing.T) {
	path := filepath.Join(t.TempDir(), "seed.db")
	st, err := Open(path, "admin", "hash", true)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	var ur int
	if err := st.db.QueryRow(`SELECT COUNT(*) FROM user_roles WHERE user_id = 'user-admin'`).Scan(&ur); err != nil || ur != 2 {
		t.Fatalf("seed user_roles = %d, err %v, want 2", ur, err)
	}
	u, err := st.UserByUsername("admin")
	if err != nil {
		t.Fatalf("seed read: %v", err)
	}
	if want := []string{"admin", "editor"}; !reflect.DeepEqual(u.Roles, want) {
		t.Fatalf("seed roles = %v, want %v", u.Roles, want)
	}
	if err := st.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	// Reopen with a different hash: password untouched, relations not duplicated.
	st2, err := Open(path, "admin", "hash-v2", true)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer st2.Close()
	if err := st2.db.QueryRow(`SELECT COUNT(*) FROM user_roles WHERE user_id = 'user-admin'`).Scan(&ur); err != nil || ur != 2 {
		t.Fatalf("reopen seed user_roles = %d, err %v, want 2", ur, err)
	}
	u2, err := st2.UserByUsername("admin")
	if err != nil {
		t.Fatalf("reopen read: %v", err)
	}
	if u2.PasswordHash != "hash" {
		t.Fatalf("password_hash = %q, want hash (seed must not overwrite)", u2.PasswordHash)
	}
	if want := []string{"admin", "editor"}; !reflect.DeepEqual(u2.Roles, want) {
		t.Fatalf("reopen roles = %v, want %v", u2.Roles, want)
	}
}

// F-002 (S2 part) · user_roles FK and CASCADE/RESTRICT semantics are enforced on
// the store connection (not just declared in DDL).
func TestUserRolesFKAndCascade(t *testing.T) {
	st, err := Open(filepath.Join(t.TempDir(), "fk.db"), "admin", "hash", false)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer st.Close()

	now := time.Now().UTC()
	if err := st.CreateUser(User{
		ID: "u2", Username: "bob", Name: "Bob",
		Roles: []string{"editor"}, PasswordHash: "h", CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}

	// FK: linking to a role id that does not exist is rejected.
	if _, err := st.db.Exec(
		`INSERT INTO user_roles (user_id, role_id) VALUES ('u2', 'role-nope')`,
	); err == nil {
		t.Fatal("expected FK violation for unknown role id")
	}

	// RESTRICT: deleting a role that a user still holds is rejected.
	if _, err := st.db.Exec(`DELETE FROM roles WHERE id = 'role-editor'`); err == nil {
		t.Fatal("expected RESTRICT violation deleting an in-use role")
	}

	// CASCADE: deleting the user removes their role relations.
	if _, err := st.db.Exec(`DELETE FROM users WHERE id = 'u2'`); err != nil {
		t.Fatalf("delete user: %v", err)
	}
	var cnt int
	if err := st.db.QueryRow(`SELECT COUNT(*) FROM user_roles WHERE user_id = 'u2'`).Scan(&cnt); err != nil || cnt != 0 {
		t.Fatalf("user_roles after user delete = %d, err %v, want 0 (CASCADE)", cnt, err)
	}
}
