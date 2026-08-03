// Roles management domain tests (GOAL-011 S2 · I-011-001 §3): key format,
// duplicate detection, and system/in-use delete protection.
package store

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestRolesStoreCreateValidation(t *testing.T) {
	st := openSeedStore(t)
	now := time.Now().UTC()

	if _, err := st.CreateRole("Bad Key!", "Bad", now); !errors.Is(err, ErrInvalidKey) {
		t.Fatalf("invalid key err = %v, want ErrInvalidKey", err)
	}
	r, err := st.CreateRole("ops", "Operator", now)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if r.ID != "role-ops" || r.Key != "ops" || r.System {
		t.Fatalf("created = %+v, want role-ops/ops/system=false", r)
	}
	if _, err := st.CreateRole("ops", "Dup", now); !errors.Is(err, ErrRoleTaken) {
		t.Fatalf("duplicate key err = %v, want ErrRoleTaken", err)
	}
}

func TestRolesStoreSystemAndInUseProtection(t *testing.T) {
	st := openSeedStore(t)
	now := time.Now().UTC()

	// system roles cannot be renamed or deleted
	if _, err := st.UpdateRole("role-admin", "Root", now); !errors.Is(err, ErrRoleSystem) {
		t.Fatalf("system update err = %v, want ErrRoleSystem", err)
	}
	if err := st.DeleteRole("role-admin"); !errors.Is(err, ErrRoleSystem) {
		t.Fatalf("system delete err = %v, want ErrRoleSystem", err)
	}

	// a role assigned to a user cannot be deleted
	if _, err := st.CreateRole("ops", "Operator", now); err != nil {
		t.Fatal(err)
	}
	u := User{
		ID: "user-ops", Username: "opsuser", Name: "Ops",
		Roles: []string{"ops"}, PasswordHash: "h", CreatedAt: now, UpdatedAt: now,
	}
	if _, err := st.CreateUserManagement(u); err != nil {
		t.Fatalf("create user: %v", err)
	}
	if err := st.DeleteRole("role-ops"); !errors.Is(err, ErrRoleInUse) {
		t.Fatalf("in-use delete err = %v, want ErrRoleInUse", err)
	}

	// a free role can be renamed and deleted
	if _, err := st.UpdateRole("role-ops", "Ops 2", now); err != nil {
		t.Fatalf("update free role: %v", err)
	}
	// free the role (remove the user relation) then delete succeeds
	if _, err := st.UpdateUser("user-ops", UserPatch{Roles: &[]string{}}, "other", now); err != nil {
		t.Fatalf("clear user roles: %v", err)
	}
	if err := st.DeleteRole("role-ops"); err != nil {
		t.Fatalf("delete free role: %v", err)
	}
	if _, err := st.GetRole("role-ops"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("get deleted role err = %v, want ErrNotFound", err)
	}
}

func TestRoleGrantLifecycle(t *testing.T) {
	st := openSeedStore(t)
	now := time.Now().UTC()
	r, err := st.CreateRoleWithGrants(
		"support", "Support", []string{"users.read"}, []string{"menu-users"}, now,
	)
	if err != nil {
		t.Fatalf("create with grants: %v", err)
	}
	if !reflect.DeepEqual(r.Permissions, []string{"users.read"}) || !reflect.DeepEqual(r.MenuItems, []string{"menu-users"}) {
		t.Fatalf("created grants = permissions %v menus %v", r.Permissions, r.MenuItems)
	}
	permissions := []string{"roles.read"}
	menus := []string{}
	r, err = st.UpdateRoleWithGrants("role-support", RolePatch{Permissions: &permissions, MenuItems: &menus}, now)
	if err != nil {
		t.Fatalf("update grants: %v", err)
	}
	if !reflect.DeepEqual(r.Permissions, []string{"roles.read"}) || len(r.MenuItems) != 0 {
		t.Fatalf("updated grants = permissions %v menus %v", r.Permissions, r.MenuItems)
	}
	invalid := []string{"ghost.permission"}
	if _, err := st.UpdateRoleWithGrants("role-support", RolePatch{Permissions: &invalid}, now); !errors.Is(err, ErrInvalidPermission) {
		t.Fatalf("invalid permission err = %v, want ErrInvalidPermission", err)
	}
	invalidMenu := []string{"menu-ghost"}
	if _, err := st.UpdateRoleWithGrants("role-support", RolePatch{MenuItems: &invalidMenu}, now); !errors.Is(err, ErrInvalidMenuItem) {
		t.Fatalf("invalid menu err = %v, want ErrInvalidMenuItem", err)
	}
	if _, err := st.UpdateRoleWithGrants("role-admin", RolePatch{Permissions: &permissions}, now); !errors.Is(err, ErrRoleSystem) {
		t.Fatalf("system grant update err = %v, want ErrRoleSystem", err)
	}
}
