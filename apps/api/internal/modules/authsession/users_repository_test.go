package authsession

import (
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"
)

// D-001 P0 · batch-delete atomicity: DeleteUsersBatch commits the whole
// selection in one transaction — a mid-batch failure (self, last-admin,
// not-found) rolls every earlier delete back.
func TestDeleteUsersBatchAtomicRollback(t *testing.T) {
	repository, st := openRepository(t, "users-batch-atomic.db", true)
	now := time.Now().UTC()

	seed := func(id, username string, roles []string) {
		t.Helper()
		if _, err := repository.CreateUserManagement(User{
			ID: id, Username: username, Name: username, Roles: roles,
			PasswordHash: "h", CreatedAt: now, UpdatedAt: now,
		}); err != nil {
			t.Fatalf("create %s: %v", username, err)
		}
	}
	seed("user-alice", "alice", []string{"viewer"})
	seed("user-bob", "bob", []string{"viewer"})

	// [alice, actor-self, bob]: alice's delete is validated and applied inside
	// the transaction, then the self-guard fails — the whole batch must roll
	// back so alice survives.
	deleted, err := repository.DeleteUsersBatch(
		[]string{"user-alice", "user-admin", "user-bob"}, "user-admin",
	)
	if !errors.Is(err, ErrSelfOperation) {
		t.Fatalf("batch err = %v, want ErrSelfOperation", err)
	}
	if deleted != 0 {
		t.Fatalf("deleted = %d, want 0 on rollback", deleted)
	}
	for _, id := range []string{"user-alice", "user-bob"} {
		if _, err := repository.GetUser(id); err != nil {
			t.Fatalf("%s rolled back but is gone: %v", id, err)
		}
	}

	// Batch containing the only admin (user-admin — admin2 was not yet created)
	// fails with ErrLastAdmin and rolls back.
	if _, err := repository.DeleteUsersBatch(
		[]string{"user-alice", "user-admin"}, "user-external",
	); !errors.Is(err, ErrLastAdmin) {
		t.Fatalf("last-admin batch err = %v, want ErrLastAdmin", err)
	}
	if _, err := repository.GetUser("user-alice"); err != nil {
		t.Fatalf("alice rolled back but is gone: %v", err)
	}

	// With a second admin present, deleting the first admin now succeeds.
	seed("user-admin2", "admin2", []string{"admin"})
	if _, err := repository.DeleteUsersBatch(
		[]string{"user-alice", "user-admin"}, "user-external",
	); err != nil {
		t.Fatalf("batch with two admins: %v", err)
	}
	if _, err := repository.GetUser("user-alice"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("alice should be deleted, got err = %v", err)
	}

	// A not-found id in the selection aborts the batch before any delete.
	if _, err := repository.DeleteUsersBatch(
		[]string{"user-bob", "usr-missing"}, "user-external",
	); !errors.Is(err, ErrNotFound) {
		t.Fatalf("not-found batch err = %v, want ErrNotFound", err)
	}
	if _, err := repository.GetUser("user-bob"); err != nil {
		t.Fatalf("bob rolled back but is gone: %v", err)
	}

	// All-valid batch deletes everything and drops the refresh tokens. A third
	// admin keeps admin2 deletable (admin2 alone would be the last admin).
	seed("user-admin3", "admin3", []string{"admin"})
	if _, err := repository.DeleteUsersBatch(
		[]string{"user-bob", "user-admin2"}, "user-external",
	); err != nil {
		t.Fatalf("valid batch: %v", err)
	}
	if count := repositoryQueryInt(t, st, `SELECT COUNT(*) FROM users`); count != 1 {
		t.Fatalf("users after valid batch = %d, want 1 (only admin3)", count)
	}
	if count := repositoryQueryInt(t, st, `SELECT COUNT(*) FROM user_roles`); count != 1 {
		t.Fatalf("user_roles after valid batch = %d, want 1", count)
	}

	// Dedup: repeated ids delete once.
	seed("user-carl", "carl", []string{"viewer"})
	deleted, err = repository.DeleteUsersBatch(
		[]string{"user-carl", "user-carl"}, "user-external",
	)
	if err != nil {
		t.Fatalf("dedup batch: %v", err)
	}
	if deleted != 1 {
		t.Fatalf("deleted = %d, want 1 (deduped)", deleted)
	}
}

// W25/I-001 (2026-08-23): deleting a user must purge user_roles links and
// user_mfa rows. Before this fix the links survived as orphans: roles.assigned
// Users stayed > 0 forever, roles.deletable stayed false and the browser e2e
// schema-crud flow (assign role → delete user → delete role) failed.
func TestDeleteUserCleansRoleAndMfaLinks(t *testing.T) {
	repository, st := openRepository(t, "delete-user-links.db", true)
	now := time.Now().UTC()
	if _, err := repository.CreateUserManagement(User{
		ID: "user-carol", Username: "carol", Name: "Carol", Roles: []string{"viewer"},
		PasswordHash: "h", CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatal(err)
	}
	if err := repositoryExec(
		t, st,
		`INSERT INTO user_mfa (user_id, status, totp_secret_ciphertext, recovery_codes_hash, created_at, updated_at)
		 VALUES ('user-carol', 'active', 'x', 'y', 1, 1)`,
	); err != nil {
		t.Fatalf("seed user_mfa: %v", err)
	}
	if n := repositoryQueryInt(t, st, `SELECT COUNT(*) FROM user_roles WHERE user_id = 'user-carol'`); n != 1 {
		t.Fatalf("user_roles before delete = %d, want 1", n)
	}
	if err := repository.DeleteUser("user-carol", "user-admin"); err != nil {
		t.Fatalf("delete user: %v", err)
	}
	if n := repositoryQueryInt(t, st, `SELECT COUNT(*) FROM user_roles WHERE user_id = 'user-carol'`); n != 0 {
		t.Fatalf("orphan user_roles after delete = %d, want 0", n)
	}
	if n := repositoryQueryInt(t, st, `SELECT COUNT(*) FROM user_mfa WHERE user_id = 'user-carol'`); n != 0 {
		t.Fatalf("orphan user_mfa after delete = %d, want 0", n)
	}
}

// W25/I-001: the batch delete path applies the same per-user link purge.
func TestDeleteUsersBatchCleansRoleAndMfaLinks(t *testing.T) {
	repository, st := openRepository(t, "delete-batch-links.db", true)
	now := time.Now().UTC()
	for _, id := range []string{"user-dave", "user-erin"} {
		if _, err := repository.CreateUserManagement(User{
			ID: id, Username: id, Name: id, Roles: []string{"viewer"},
			PasswordHash: "h", CreatedAt: now, UpdatedAt: now,
		}); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := repository.DeleteUsersBatch([]string{"user-dave", "user-erin"}, "user-external"); err != nil {
		t.Fatalf("batch delete: %v", err)
	}
	if n := repositoryQueryInt(t, st, `SELECT COUNT(*) FROM user_roles WHERE user_id IN ('user-dave','user-erin')`); n != 0 {
		t.Fatalf("orphan user_roles after batch = %d, want 0", n)
	}
}

// A-002 F-001 regression: a batch containing EVERY admin must fail with
// ErrLastAdmin — per-id last-admin counting is unsound for a batch (each admin
// still sees the others present), so the batch-level guard rejects and rolls
// back instead of leaving zero admins.
func TestDeleteUsersBatchRejectsRemovingAllAdmins(t *testing.T) {
	repository, st := openRepository(t, "users-batch-all-admins.db", true)
	now := time.Now().UTC()
	if _, err := repository.CreateUserManagement(User{
		ID: "user-admin2", Username: "admin2", Name: "Admin 2", Roles: []string{"admin"},
		PasswordHash: "h", CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatal(err)
	}

	// Both admins in the selection → ErrLastAdmin, nothing deleted.
	if _, err := repository.DeleteUsersBatch(
		[]string{"user-admin", "user-admin2"}, "user-external",
	); !errors.Is(err, ErrLastAdmin) {
		t.Fatalf("all-admins batch err = %v, want ErrLastAdmin", err)
	}
	if count := repositoryQueryInt(t, st, `SELECT COUNT(*) FROM users`); count != 2 {
		t.Fatalf("users after rejected all-admins batch = %d, want 2", count)
	}

	// One admin plus one non-admin admin2 (the actor is admin2): still the only
	// admin remaining outside the batch → rejected.
	seedViewer := func(id, username string) {
		t.Helper()
		if _, err := repository.CreateUserManagement(User{
			ID: id, Username: username, Name: username, Roles: []string{"viewer"},
			PasswordHash: "h", CreatedAt: now, UpdatedAt: now,
		}); err != nil {
			t.Fatalf("create %s: %v", username, err)
		}
	}
	seedViewer("user-viewer", "viewer1")
	if _, err := repository.DeleteUsersBatch(
		[]string{"user-viewer", "user-admin", "user-admin2"}, "user-external",
	); !errors.Is(err, ErrLastAdmin) {
		t.Fatalf("all-admins+viewer batch err = %v, want ErrLastAdmin", err)
	}
	if _, err := repository.GetUser("user-viewer"); err != nil {
		t.Fatalf("viewer rolled back but is gone: %v", err)
	}

	// A third admin outside the batch makes the same batch legal.
	if _, err := repository.CreateUserManagement(User{
		ID: "user-admin3", Username: "admin3", Name: "Admin 3", Roles: []string{"admin"},
		PasswordHash: "h", CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := repository.DeleteUsersBatch(
		[]string{"user-viewer", "user-admin", "user-admin2"}, "user-external",
	); err != nil {
		t.Fatalf("batch with surviving admin3: %v", err)
	}
	if count := repositoryQueryInt(t, st, `SELECT COUNT(*) FROM users`); count != 1 {
		t.Fatalf("users after legal batch = %d, want 1 (admin3)", count)
	}
}

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
