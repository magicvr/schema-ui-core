package authsession

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestRolesRepositoryCreateValidation(t *testing.T) {
	repository, _ := openRepository(t, "roles-create.db", true)
	now := time.Now().UTC()

	if _, err := repository.CreateRole("Bad Key!", "Bad", now); !errors.Is(err, ErrInvalidKey) {
		t.Fatalf("invalid key err = %v, want ErrInvalidKey", err)
	}
	role, err := repository.CreateRole("ops", "Operator", now)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if role.ID != "role-ops" || role.Key != "ops" || role.System {
		t.Fatalf("created = %+v, want role-ops/ops/system=false", role)
	}
	if _, err := repository.CreateRole("ops", "Dup", now); !errors.Is(err, ErrRoleTaken) {
		t.Fatalf("duplicate key err = %v, want ErrRoleTaken", err)
	}
}

func TestRolesRepositorySystemAndInUseProtection(t *testing.T) {
	repository, _ := openRepository(t, "roles-protection.db", true)
	now := time.Now().UTC()

	if _, err := repository.UpdateRole("role-admin", "Root", now); !errors.Is(err, ErrRoleSystem) {
		t.Fatalf("system update err = %v, want ErrRoleSystem", err)
	}
	if err := repository.DeleteRole("role-admin"); !errors.Is(err, ErrRoleSystem) {
		t.Fatalf("system delete err = %v, want ErrRoleSystem", err)
	}
	if _, err := repository.CreateRole("ops", "Operator", now); err != nil {
		t.Fatal(err)
	}
	if _, err := repository.CreateUserManagement(User{
		ID: "user-ops", Username: "opsuser", Name: "Ops", Roles: []string{"ops"},
		PasswordHash: "h", CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatalf("create user: %v", err)
	}
	if err := repository.DeleteRole("role-ops"); !errors.Is(err, ErrRoleInUse) {
		t.Fatalf("in-use delete err = %v, want ErrRoleInUse", err)
	}
	if _, err := repository.UpdateRole("role-ops", "Ops 2", now); err != nil {
		t.Fatalf("update role: %v", err)
	}
	if _, err := repository.UpdateUser("user-ops", UserPatch{Roles: &[]string{}}, "other", now); err != nil {
		t.Fatalf("clear user roles: %v", err)
	}
	if err := repository.DeleteRole("role-ops"); err != nil {
		t.Fatalf("delete free role: %v", err)
	}
	if _, err := repository.GetRole("role-ops"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("get deleted role err = %v, want ErrNotFound", err)
	}
}

func TestRoleGrantLifecycle(t *testing.T) {
	repository, _ := openRepository(t, "role-grants.db", true)
	now := time.Now().UTC()
	role, err := repository.CreateRoleWithGrants(
		"support", "Support", []string{"users.read"}, []string{"menu-users"}, now,
	)
	if err != nil {
		t.Fatalf("create with grants: %v", err)
	}
	if !reflect.DeepEqual(role.Permissions, []string{"users.read"}) || !reflect.DeepEqual(role.MenuItems, []string{"menu-users"}) {
		t.Fatalf("created grants = permissions %v menus %v", role.Permissions, role.MenuItems)
	}
	permissions := []string{"roles.read"}
	menus := []string{}
	role, err = repository.UpdateRoleWithGrants(
		"role-support", RolePatch{Permissions: &permissions, MenuItems: &menus}, now,
	)
	if err != nil {
		t.Fatalf("update grants: %v", err)
	}
	if !reflect.DeepEqual(role.Permissions, []string{"roles.read"}) || len(role.MenuItems) != 0 {
		t.Fatalf("updated grants = permissions %v menus %v", role.Permissions, role.MenuItems)
	}
	invalidPermissions := []string{"ghost.permission"}
	if _, err := repository.UpdateRoleWithGrants(
		"role-support", RolePatch{Permissions: &invalidPermissions}, now,
	); !errors.Is(err, ErrInvalidPermission) {
		t.Fatalf("invalid permission err = %v, want ErrInvalidPermission", err)
	}
	invalidMenus := []string{"menu-ghost"}
	if _, err := repository.UpdateRoleWithGrants(
		"role-support", RolePatch{MenuItems: &invalidMenus}, now,
	); !errors.Is(err, ErrInvalidMenuItem) {
		t.Fatalf("invalid menu err = %v, want ErrInvalidMenuItem", err)
	}
	if _, err := repository.UpdateRoleWithGrants(
		"role-admin", RolePatch{Permissions: &permissions}, now,
	); !errors.Is(err, ErrRoleSystem) {
		t.Fatalf("system grant update err = %v, want ErrRoleSystem", err)
	}
}
