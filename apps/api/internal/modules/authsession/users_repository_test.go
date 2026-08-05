package authsession

import (
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestUsersLastAdminProtection(t *testing.T) {
	repository, _ := openRepository(t, "last-admin.db", true)
	now := time.Now().UTC()
	if _, err := repository.UpdateUser(
		"user-admin", UserPatch{Roles: &[]string{"editor"}}, "other-user", now,
	); !errors.Is(err, ErrLastAdmin) {
		t.Fatalf("demote only admin = %v, want ErrLastAdmin", err)
	}
	if err := repository.DeleteUser("user-admin", "other-user"); !errors.Is(err, ErrLastAdmin) {
		t.Fatalf("delete only admin = %v, want ErrLastAdmin", err)
	}
	if _, err := repository.CreateUserManagement(User{
		ID: "user-admin2", Username: "admin2", Name: "Admin 2", Roles: []string{"admin"},
		PasswordHash: "h", CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatalf("create second admin: %v", err)
	}
	if _, err := repository.UpdateUser(
		"user-admin", UserPatch{Roles: &[]string{"editor"}}, "user-admin2", now,
	); err != nil {
		t.Fatalf("demote with second admin: %v", err)
	}
}

func TestUsersCreateValidationAndRoundTrip(t *testing.T) {
	repository, st := openRepository(t, "users.db", true)
	now := time.Now().UTC()
	unknown := User{
		ID: "user-new", Username: "newuser", Name: "New", Roles: []string{"ghost"},
		PasswordHash: "h", CreatedAt: now, UpdatedAt: now,
	}
	if _, err := repository.CreateUserManagement(unknown); !errors.Is(err, ErrInvalidRole) {
		t.Fatalf("unknown role = %v, want ErrInvalidRole", err)
	}
	if count := repositoryQueryInt(t, st, `SELECT COUNT(*) FROM roles WHERE key = 'ghost'`); count != 0 {
		t.Fatalf("ghost role count = %d, want 0", count)
	}

	user := unknown
	user.Roles = []string{"viewer", "editor"}
	got, err := repository.CreateUserManagement(user)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if want := []string{"editor", "viewer"}; !reflect.DeepEqual(got.Roles, want) {
		t.Fatalf("roles = %v, want %v", got.Roles, want)
	}
	duplicate := user
	duplicate.ID = "user-duplicate"
	if _, err := repository.CreateUserManagement(duplicate); !errors.Is(err, ErrUsernameTaken) {
		t.Fatalf("duplicate username = %v, want ErrUsernameTaken", err)
	}
	if _, err := repository.UpdateUser(
		user.ID, UserPatch{Roles: &[]string{"viewer"}}, "other", now,
	); err != nil {
		t.Fatalf("update roles: %v", err)
	}
	got, err = repository.GetUser(user.ID)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"viewer"}; !reflect.DeepEqual(got.Roles, want) {
		t.Fatalf("updated roles = %v, want %v", got.Roles, want)
	}
}

func TestDeleteUserSerializesLastAdminCheck(t *testing.T) {
	repository, st := openRepository(t, "concurrent-admin.db", true)
	now := time.Now().UTC()
	if _, err := repository.CreateUserManagement(User{
		ID: "user-admin2", Username: "admin2", Name: "Admin 2", Roles: []string{"admin"},
		PasswordHash: "h", CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatal(err)
	}

	start := make(chan struct{})
	errs := make(chan error, 2)
	var group sync.WaitGroup
	for _, id := range []string{"user-admin", "user-admin2"} {
		group.Add(1)
		go func(userID string) {
			defer group.Done()
			<-start
			errs <- repository.DeleteUser(userID, "external-actor")
		}(id)
	}
	close(start)
	group.Wait()
	close(errs)

	var succeeded, blocked int
	for err := range errs {
		switch {
		case err == nil:
			succeeded++
		case errors.Is(err, ErrLastAdmin):
			blocked++
		default:
			t.Fatalf("concurrent delete: %v", err)
		}
	}
	if succeeded != 1 || blocked != 1 {
		t.Fatalf("succeeded=%d blocked=%d, want 1/1", succeeded, blocked)
	}
	if admins := repositoryQueryInt(t, st, `SELECT COUNT(*) FROM user_roles ur JOIN roles r ON r.id = ur.role_id WHERE r.key = 'admin'`); admins != 1 {
		t.Fatalf("admins = %d, want 1", admins)
	}
}
