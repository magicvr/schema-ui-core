package store

import (
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

// V-SEED-01 · seeding with no prior users creates the stable roles
// (system=1), permissions, representative menu item and grants; editor and
// viewer stay read-only.
func TestSeedRBACEntitiesAndGrants(t *testing.T) {
	st, err := Open(filepath.Join(t.TempDir(), "seed-rbac.db"), "admin", "hash", true)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer st.Close()

	rows, err := st.db.Query(`SELECT key FROM roles WHERE system = 1 ORDER BY key`)
	if err != nil {
		t.Fatalf("system roles: %v", err)
	}
	var systemRoles []string
	for rows.Next() {
		var k string
		if err := rows.Scan(&k); err != nil {
			rows.Close()
			t.Fatal(err)
		}
		systemRoles = append(systemRoles, k)
	}
	if err := rows.Close(); err != nil {
		t.Fatal(err)
	}
	if want := []string{"admin", "editor", "viewer"}; !reflect.DeepEqual(systemRoles, want) {
		t.Fatalf("system roles = %v, want %v", systemRoles, want)
	}

	prows, err := st.db.Query(`SELECT key FROM permissions ORDER BY key`)
	if err != nil {
		t.Fatalf("permissions: %v", err)
	}
	var perms []string
	for prows.Next() {
		var k string
		if err := prows.Scan(&k); err != nil {
			prows.Close()
			t.Fatal(err)
		}
		perms = append(perms, k)
	}
	if err := prows.Close(); err != nil {
		t.Fatal(err)
	}
	if want := []string{
		"operations.read",
		"roles.assign", "roles.read", "roles.write",
		"settings.read", "settings.write",
		"users.read", "users.write",
	}; !reflect.DeepEqual(perms, want) {
		t.Fatalf("permissions = %v, want %v", perms, want)
	}

	var feature string
	if err := st.db.QueryRow(`SELECT feature_key FROM menu_items WHERE id = 'menu-users'`).Scan(&feature); err != nil {
		t.Fatalf("menu item: %v", err)
	}
	if feature != "menu_users" {
		t.Fatalf("menu feature_key = %q, want menu_users", feature)
	}

	count := func(q string) int {
		var n int
		if err := st.db.QueryRow(q).Scan(&n); err != nil {
			t.Fatalf("%s: %v", q, err)
		}
		return n
	}
	// admin 8 + editor 3 + viewer 2 = 13 (GOAL-013 adds settings/operations).
	if n := count(`SELECT COUNT(*) FROM role_permissions`); n != 13 {
		t.Fatalf("role_permissions = %d, want 13 (admin 8, editor 3, viewer 2)", n)
	}
	if n := count(`SELECT COUNT(*) FROM role_menu_items`); n != 5 {
		t.Fatalf("role_menu_items = %d, want 5 (admin 4 + editor activity)", n)
	}
	// admin: users/roles/settings/operations; editor: users/roles read + operations; viewer: users/roles read.
	if n := count(`SELECT COUNT(*) FROM role_permissions WHERE role_id = 'role-admin'`); n != 8 {
		t.Fatalf("admin grants = %d, want 8", n)
	}
	if n := count(`SELECT COUNT(*) FROM role_permissions WHERE role_id = 'role-editor'`); n != 3 {
		t.Fatalf("editor grants = %d, want 3", n)
	}
	if n := count(`SELECT COUNT(*) FROM role_permissions WHERE role_id = 'role-viewer'`); n != 2 {
		t.Fatalf("viewer grants = %d, want 2", n)
	}
	var write int
	if err := st.db.QueryRow(`SELECT COUNT(*) FROM role_permissions WHERE permission_id = 'perm-users-write'`).Scan(&write); err != nil {
		t.Fatal(err)
	}
	if write != 1 {
		t.Fatalf("users.write granted to %d roles, want only admin", write)
	}
}

// V-SEED-01 · an existing user does not cause the relation seed to be skipped;
// repeated startup is idempotent and does not overwrite non-seed user fields.
func TestSeedRBACIncrementalWithExistingUsers(t *testing.T) {
	path := filepath.Join(t.TempDir(), "seed-incremental.db")
	st, err := Open(path, "admin", "hash", false)
	if err != nil {
		t.Fatalf("open (no seed): %v", err)
	}
	now := time.Now().UTC()
	if err := st.CreateUser(User{
		ID: "u1", Username: "alice", Name: "Alice",
		Roles: []string{"viewer"}, PasswordHash: "h", CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	// Reopen with seeding enabled: finalized system-data contributions repair
	// system relations, but fresh bootstrap does not create a second admin.
	st2, err := Open(path, "admin", "hash", true)
	if err != nil {
		t.Fatalf("reopen with seed: %v", err)
	}
	count := func(db *Store, q string) int {
		var n int
		if err := db.db.QueryRow(q).Scan(&n); err != nil {
			t.Fatalf("%s: %v", q, err)
		}
		return n
	}
	if n := count(st2, `SELECT COUNT(*) FROM role_permissions`); n != 13 {
		t.Fatalf("role_permissions = %d, want 13", n)
	}
	if n := count(st2, `SELECT COUNT(*) FROM role_menu_items`); n != 5 {
		t.Fatalf("role_menu_items = %d, want 5", n)
	}
	// Non-seed user untouched.
	u, err := st2.UserByID("u1")
	if err != nil {
		t.Fatalf("non-seed user: %v", err)
	}
	if u.Name != "Alice" || u.PasswordHash != "h" {
		t.Fatalf("non-seed user overwritten: %+v", u)
	}
	if want := []string{"viewer"}; !reflect.DeepEqual(u.Roles, want) {
		t.Fatalf("non-seed user roles = %v, want %v", u.Roles, want)
	}
	// The existing non-fresh database is not bootstrapped a second time.
	if n := count(st2, `SELECT COUNT(*) FROM user_roles WHERE user_id = 'user-admin'`); n != 0 {
		t.Fatalf("unexpected fresh-bootstrap user_roles on existing DB = %d", n)
	}

	// A third open must not duplicate any relation.
	if err := st2.Close(); err != nil {
		t.Fatal(err)
	}
	st3, err := Open(path, "admin", "hash", true)
	if err != nil {
		t.Fatalf("third open: %v", err)
	}
	defer st3.Close()
	if n := count(st3, `SELECT COUNT(*) FROM role_permissions`); n != 13 {
		t.Fatalf("after 3rd open role_permissions = %d, want 13 (idempotent)", n)
	}
	if n := count(st3, `SELECT COUNT(*) FROM role_menu_items`); n != 5 {
		t.Fatalf("after 3rd open role_menu_items = %d, want 5 (idempotent)", n)
	}
	if n := count(st3, `SELECT COUNT(*) FROM user_roles WHERE user_id = 'user-admin'`); n != 0 {
		t.Fatalf("after 3rd open unexpected bootstrap user_roles = %d", n)
	}
}

// S4 gate source · PermissionsForUser resolves a user's permission keys from
// the seeded role-permission relations.
func TestPermissionsForUser(t *testing.T) {
	st, err := Open(filepath.Join(t.TempDir(), "perms.db"), "admin", "hash", true)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer st.Close()

	adminPerms, err := st.PermissionsForUser("user-admin")
	if err != nil {
		t.Fatalf("admin permissions: %v", err)
	}
	if want := []string{
		"operations.read",
		"roles.assign", "roles.read", "roles.write",
		"settings.read", "settings.write",
		"users.read", "users.write",
	}; !reflect.DeepEqual(adminPerms, want) {
		t.Fatalf("admin perms = %v, want %v", adminPerms, want)
	}

	now := time.Now().UTC()
	if err := st.CreateUser(User{
		ID: "u2", Username: "viewer", Name: "Viewer",
		Roles: []string{"viewer"}, PasswordHash: "h", CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatalf("create viewer: %v", err)
	}
	viewerPerms, err := st.PermissionsForUser("u2")
	if err != nil {
		t.Fatalf("viewer permissions: %v", err)
	}
	if want := []string{"roles.read", "users.read"}; !reflect.DeepEqual(viewerPerms, want) {
		t.Fatalf("viewer perms = %v, want %v", viewerPerms, want)
	}

	if err := st.CreateUser(User{
		ID: "u3", Username: "auditor", Name: "Auditor",
		Roles: []string{"custom"}, PasswordHash: "h", CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatalf("create custom: %v", err)
	}
	customPerms, err := st.PermissionsForUser("u3")
	if err != nil {
		t.Fatalf("custom permissions: %v", err)
	}
	if len(customPerms) != 0 {
		t.Fatalf("custom perms = %v, want empty (no grants)", customPerms)
	}
}

// V-MENU-01 · FeaturesForUser projects every registered menu feature key: true
// only when any of the user's roles holds the grant (multi-role OR).
func TestFeaturesForUser(t *testing.T) {
	st, err := Open(filepath.Join(t.TempDir(), "feat.db"), "admin", "hash", true)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer st.Close()

	assert := func(userID string, want bool, label string) {
		t.Helper()
		feat, err := st.FeaturesForUser(userID)
		if err != nil {
			t.Fatalf("%s features: %v", label, err)
		}
		if v, ok := feat["menu_users"]; !ok {
			t.Fatalf("%s: menu_users missing from %v", label, feat)
		} else if v != want {
			t.Fatalf("%s: menu_users = %v, want %v", label, v, want)
		}
	}

	assert("user-admin", true, "admin")
	now := time.Now().UTC()
	mkUser := func(id, username string, roles []string) {
		t.Helper()
		if err := st.CreateUser(User{
			ID: id, Username: username, Name: username,
			Roles: roles, PasswordHash: "h", CreatedAt: now, UpdatedAt: now,
		}); err != nil {
			t.Fatalf("create %s: %v", username, err)
		}
	}
	mkUser("u-v", "viewer", []string{"viewer"})
	mkUser("u-e", "editor", []string{"editor"})
	mkUser("u-av", "adminviewer", []string{"admin", "viewer"})
	mkUser("u-0", "norole", nil)
	assert("u-v", false, "viewer")
	assert("u-e", false, "editor")
	assert("u-av", true, "admin+viewer (OR)")
	assert("u-0", false, "no-role")
}
