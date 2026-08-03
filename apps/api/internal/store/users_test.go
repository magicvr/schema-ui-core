// Users management domain tests (GOAL-011 S2 · I-011-001 §2.4/§2.3): last-admin
// protection (reachable only via a non-admin actor) and create-time role
// validation / no-implicit-role.
package store

import (
	"errors"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func openSeedStore(t *testing.T) *Store {
	t.Helper()
	st, err := Open(filepath.Join(t.TempDir(), "users.db"), "admin", "hash", true)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return st
}

// I-011-001 §2.4 · a non-admin actor cannot delete or demote the only admin.
func TestUsersLastAdminProtection(t *testing.T) {
	st := openSeedStore(t)
	now := time.Now().UTC()

	// demote the only admin by a non-admin actor → ErrLastAdmin
	_, err := st.UpdateUser("user-admin", UserPatch{Roles: &[]string{"editor"}}, "other-user", now)
	if !errors.Is(err, ErrLastAdmin) {
		t.Fatalf("demote only admin err = %v, want ErrLastAdmin", err)
	}
	// delete the only admin by a non-admin actor → ErrLastAdmin
	err = st.DeleteUser("user-admin", "other-user")
	if !errors.Is(err, ErrLastAdmin) {
		t.Fatalf("delete only admin err = %v, want ErrLastAdmin", err)
	}
	// a second admin makes demotion legal (actor is still an admin afterward)
	u := User{
		ID: "user-admin2", Username: "admin2", Name: "Admin 2",
		Roles: []string{"admin"}, PasswordHash: "h", CreatedAt: now, UpdatedAt: now,
	}
	if _, err := st.CreateUserManagement(u); err != nil {
		t.Fatalf("create second admin: %v", err)
	}
	if _, err := st.UpdateUser("user-admin", UserPatch{Roles: &[]string{"editor"}}, "user-admin2", now); err != nil {
		t.Fatalf("demote with second admin: %v", err)
	}
}

// I-011-001 §2.3 · CreateUserManagement rejects unknown role keys and does not
// implicitly create role rows.
func TestUsersCreateManagementRoleValidation(t *testing.T) {
	st := openSeedStore(t)
	now := time.Now().UTC()
	u := User{
		ID: "user-new", Username: "newuser", Name: "New",
		Roles: []string{"ghost"}, PasswordHash: "h", CreatedAt: now, UpdatedAt: now,
	}
	_, err := st.CreateUserManagement(u)
	if !errors.Is(err, ErrInvalidRole) {
		t.Fatalf("unknown role err = %v, want ErrInvalidRole", err)
	}
	var n int
	if err := st.db.QueryRow(`SELECT COUNT(*) FROM roles WHERE key = 'ghost'`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("ghost role was implicitly created: count = %d", n)
	}
	// duplicate username → ErrUsernameTaken
	u2 := u
	u2.ID = "user-new2"
	u2.Roles = []string{"viewer"}
	if _, err := st.CreateUserManagement(u2); err != nil {
		t.Fatalf("create user: %v", err)
	}
	u3 := u2
	u3.ID = "user-new3"
	if _, err := st.CreateUserManagement(u3); !errors.Is(err, ErrUsernameTaken) {
		t.Fatalf("duplicate username err = %v, want ErrUsernameTaken", err)
	}
}

// I-011-001 §2.3 · a created user's roles round-trip through GetUser with the
// normalized relation, and the legacy JSON + user_roles stay set-equal.
func TestUsersCreateManagementRolesRoundTrip(t *testing.T) {
	st := openSeedStore(t)
	now := time.Now().UTC()
	u := User{
		ID: "user-rt", Username: "roundtrip", Name: "RT",
		Roles: []string{"viewer", "editor"}, PasswordHash: "h", CreatedAt: now, UpdatedAt: now,
	}
	got, err := st.CreateUserManagement(u)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if want := []string{"editor", "viewer"}; !reflect.DeepEqual(got.Roles, want) {
		t.Fatalf("roles = %v, want %v (ascending)", got.Roles, want)
	}
	// UpdateUser role change is reflected on read.
	if _, err := st.UpdateUser("user-rt", UserPatch{Roles: &[]string{"viewer"}}, "other", now); err != nil {
		t.Fatalf("update roles: %v", err)
	}
	got, err = st.GetUser("user-rt")
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"viewer"}; !reflect.DeepEqual(got.Roles, want) {
		t.Fatalf("roles after update = %v, want %v", got.Roles, want)
	}
}
