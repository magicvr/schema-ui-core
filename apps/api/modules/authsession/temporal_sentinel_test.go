package authsession

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
)

// This file pins the workspace-040 R2 sentinel and monotonicity rules that
// D-001 §2 makes normative for core.auth-session:
//
//	#5  users.locked_until              0 -> NULL   (P-5-SENT / P-5-NULL)
//	#6  users.last_login_failure_at     0 -> NULL   (P-6-SENT / P-6-NULL)
//	#20 login_failures.locked_until     0 -> NULL   (P-20-SENT / P-20-NULL)
//	#4  users.updated_at                monotonic   (P-4-MONO)
//	#11 roles.updated_at                monotonic   (P-11-MONO)
//
// The pre-R2 suite covered none of these, and the "0 meant NULL" sentinel is
// exactly the change a regression here would silently undo, so the assertions
// below read the RAW stored value as well as the scanned domain value.

// rawColumn reads one column as stored text, which is the only way to tell SQL
// NULL apart from a stored epoch value ("0", "1970-01-01T00:00:00.000000Z").
func rawColumn(t *testing.T, st *store.Store, query string, args ...any) sql.NullString {
	t.Helper()
	var value sql.NullString
	if err := st.WithTx(context.Background(), func(tx *sql.Tx) error {
		return tx.QueryRow(query, args...).Scan(&value)
	}); err != nil {
		t.Fatalf("query %q: %v", query, err)
	}
	return value
}

func userIDs(users []User) []string {
	ids := make([]string, 0, len(users))
	for i := range users {
		ids = append(ids, users[i].ID)
	}
	return ids
}

func TestUserLockAndFailureSentinelsAreStoredAsNull(t *testing.T) {
	repository, st := openRepository(t, "user-lock-sentinels.db", true)
	admin, err := repository.UserByUsername("admin")
	if err != nil {
		t.Fatal(err)
	}

	// A fresh account is not locked, and "not locked" is SQL NULL — not 0 and
	// not the epoch instant.
	if admin.LockedUntil.Valid {
		t.Fatalf("fresh account locked_until = %v, want not locked", admin.LockedUntil.Time)
	}
	if got := rawColumn(t, st, `SELECT locked_until FROM users WHERE id = ?`, admin.ID); got.Valid {
		t.Fatalf("users.locked_until stored as %q, want SQL NULL (the legacy 0 sentinel is gone)", got.String)
	}
	if got := rawColumn(t, st, `SELECT last_login_failure_at FROM users WHERE id = ?`, admin.ID); got.Valid {
		t.Fatalf("users.last_login_failure_at stored as %q, want SQL NULL", got.String)
	}

	now := time.Now().UTC()
	// P-6-NULL: a NULL last_login_failure_at must count as "no recent failure",
	// so the first failure restarts the consecutive run at 1.
	if _, err := repository.RecordLoginFailure(admin.ID, 5, now.Add(time.Hour), now); err != nil {
		t.Fatalf("RecordLoginFailure: %v", err)
	}
	afterFirst, err := repository.UserByID(admin.ID)
	if err != nil {
		t.Fatal(err)
	}
	if afterFirst.FailedLoginCount != 1 {
		t.Fatalf("failed_login_count after the first failure = %d, want 1", afterFirst.FailedLoginCount)
	}
	if got := rawColumn(t, st, `SELECT last_login_failure_at FROM users WHERE id = ?`, admin.ID); !got.Valid {
		t.Fatal("last_login_failure_at was not written as an instant")
	}

	// Threshold reached: the lock window is persisted as a real instant.
	lockUntil := now.Add(time.Hour)
	if _, err := repository.RecordLoginFailure(admin.ID, 2, lockUntil, now.Add(time.Second)); err != nil {
		t.Fatalf("open lock window: %v", err)
	}
	locked, err := repository.UserByID(admin.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !locked.LockedUntil.Valid || !locked.LockedUntil.Time.Equal(lockUntil.Truncate(time.Microsecond)) {
		t.Fatalf("locked_until = %+v, want %s", locked.LockedUntil, lockUntil.Truncate(time.Microsecond))
	}

	// P-5-NULL: clearing the lock writes NULL, never 0.
	if _, err := repository.UnlockUser(admin.ID, now.Add(2*time.Hour)); err != nil {
		t.Fatalf("UnlockUser: %v", err)
	}
	if got := rawColumn(t, st, `SELECT locked_until FROM users WHERE id = ?`, admin.ID); got.Valid {
		t.Fatalf("UnlockUser left users.locked_until = %q, want SQL NULL", got.String)
	}
	if got := rawColumn(t, st, `SELECT last_login_failure_at FROM users WHERE id = ?`, admin.ID); got.Valid {
		t.Fatalf("UnlockUser left users.last_login_failure_at = %q, want SQL NULL", got.String)
	}
	cleared, err := repository.UserByID(admin.ID)
	if err != nil {
		t.Fatal(err)
	}
	if cleared.LockedUntil.Valid || cleared.FailedLoginCount != 0 {
		t.Fatalf("after unlock: locked_until = %+v count = %d, want NULL/0", cleared.LockedUntil, cleared.FailedLoginCount)
	}

	// The successful-login reset path clears both sentinel columns too.
	if _, err := repository.RecordLoginFailure(admin.ID, 1, lockUntil, now.Add(3*time.Hour)); err != nil {
		t.Fatalf("re-lock: %v", err)
	}
	if err := repository.ResetLoginFailures(admin.ID, now.Add(4*time.Hour)); err != nil {
		t.Fatalf("ResetLoginFailures: %v", err)
	}
	if got := rawColumn(t, st, `SELECT locked_until FROM users WHERE id = ?`, admin.ID); got.Valid {
		t.Fatalf("ResetLoginFailures left users.locked_until = %q, want SQL NULL", got.String)
	}
	if got := rawColumn(t, st, `SELECT last_login_failure_at FROM users WHERE id = ?`, admin.ID); got.Valid {
		t.Fatalf("ResetLoginFailures left users.last_login_failure_at = %q, want SQL NULL", got.String)
	}
}

func TestUsersLockedFilterTreatsNullAndClosedWindowsAsUnlocked(t *testing.T) {
	repository, _ := openRepository(t, "users-locked-filter.db", true)
	now := time.Now().UTC()
	admin, err := repository.UserByUsername("admin")
	if err != nil {
		t.Fatal(err)
	}
	if err := repository.CreateUser(User{
		ID: "user-viewer", Username: "viewer", Name: "Viewer", Roles: []string{"viewer"},
		PasswordHash: "h", CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatal(err)
	}

	// admin holds an OPEN window; user-viewer holds a window that already
	// closed. Neither is allowed to be classified by the epoch sentinel any
	// more, so the pair covers both the "IS NULL" and the "<= now" arms.
	if _, err := repository.RecordLoginFailure(admin.ID, 1, now.Add(time.Hour), now); err != nil {
		t.Fatalf("lock admin: %v", err)
	}
	if _, err := repository.RecordLoginFailure("user-viewer", 1, now.Add(-time.Minute), now.Add(-2*time.Minute)); err != nil {
		t.Fatalf("lock viewer: %v", err)
	}

	lockedTrue, lockedFalse := true, false
	locked, total, err := repository.ListUsers(UserFilter{Locked: &lockedTrue, Page: 1, PageSize: 50})
	if err != nil {
		t.Fatalf("ListUsers(locked=true): %v", err)
	}
	if total != 1 || len(locked) != 1 || locked[0].ID != admin.ID {
		t.Fatalf("locked=true = %v (total %d), want only %s", userIDs(locked), total, admin.ID)
	}
	unlocked, total, err := repository.ListUsers(UserFilter{Locked: &lockedFalse, Page: 1, PageSize: 50})
	if err != nil {
		t.Fatalf("ListUsers(locked=false): %v", err)
	}
	if total != 1 || len(unlocked) != 1 || unlocked[0].ID != "user-viewer" {
		t.Fatalf("locked=false = %v (total %d), want only user-viewer", userIDs(unlocked), total)
	}

	// The unlocked arm is an OR and the WHERE clauses are AND-joined, so it must
	// stay parenthesized: with the bug, the closed-window row would escape the
	// keyword filter entirely.
	searched, total, err := repository.ListUsers(UserFilter{Q: "admin", Locked: &lockedFalse, Page: 1, PageSize: 50})
	if err != nil {
		t.Fatalf("ListUsers(q=admin, locked=false): %v", err)
	}
	if total != 0 || len(searched) != 0 {
		t.Fatalf("q=admin + locked=false = %v (total %d), want no match", userIDs(searched), total)
	}
}

func TestSourceLockWindowIsNullUntilOpened(t *testing.T) {
	repository, st := openRepository(t, "source-lock-sentinels.db", true)
	admin, err := repository.UserByUsername("admin")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	const source = "198.51.100.7"

	locked, err := repository.RecordLoginFailureFor(admin.ID, source, 5, now.Add(time.Hour), now)
	if err != nil {
		t.Fatalf("RecordLoginFailureFor insert: %v", err)
	}
	if locked {
		t.Fatal("the first source failure opened a lock")
	}
	// P-20-NULL: the insert carries no lock window (NULL), not the legacy 0.
	if got := rawColumn(t, st, `SELECT locked_until FROM login_failures WHERE user_id = ? AND ip = ?`, admin.ID, source); got.Valid {
		t.Fatalf("login_failures.locked_until on insert = %q, want SQL NULL", got.String)
	}
	if isLocked, err := repository.LoginLockedFor(admin.ID, source, now); err != nil || isLocked {
		t.Fatalf("LoginLockedFor with a NULL window = %v (err %v), want false", isLocked, err)
	}

	// Threshold reached: the window is written as an instant and is then read
	// back with Valid && After(now).
	window := now.Add(time.Hour)
	locked, err = repository.RecordLoginFailureFor(admin.ID, source, 1, window, now.Add(time.Second))
	if err != nil {
		t.Fatalf("RecordLoginFailureFor lock: %v", err)
	}
	if !locked {
		t.Fatal("the threshold did not open the source lock")
	}
	if isLocked, err := repository.LoginLockedFor(admin.ID, source, now); err != nil || !isLocked {
		t.Fatalf("LoginLockedFor inside the window = %v (err %v), want true", isLocked, err)
	}
	if isLocked, err := repository.LoginLockedFor(admin.ID, source, window.Add(time.Second)); err != nil || isLocked {
		t.Fatalf("LoginLockedFor after the window closed = %v (err %v), want false", isLocked, err)
	}
}

func TestUpdatedAtIsMonotonicAtMicrosecondGranularity(t *testing.T) {
	repository, _ := openRepository(t, "monotonic-updated-at.db", true)
	base := time.Date(2026, 9, 20, 12, 0, 0, 123456000, time.UTC)
	if _, err := repository.CreateUserManagement(User{
		ID: "u-mono", Username: "mono", Name: "Mono", Roles: []string{"viewer"},
		PasswordHash: "h", CreatedAt: base, UpdatedAt: base,
	}); err != nil {
		t.Fatalf("CreateUserManagement: %v", err)
	}

	// P-4-MONO: writing at the SAME instant as the stored value must still
	// advance the column (by one microsecond) instead of colliding with it.
	name1 := "Mono 1"
	first, err := repository.UpdateUser("u-mono", UserPatch{Name: &name1}, "actor", base)
	if err != nil {
		t.Fatalf("first UpdateUser: %v", err)
	}
	if !first.UpdatedAt.After(base) {
		t.Fatalf("updated_at %s did not advance past %s", first.UpdatedAt, base)
	}
	name2 := "Mono 2"
	second, err := repository.UpdateUser("u-mono", UserPatch{Name: &name2}, "actor", base)
	if err != nil {
		t.Fatalf("second UpdateUser: %v", err)
	}
	if got := second.UpdatedAt.Sub(first.UpdatedAt); got != time.Microsecond {
		t.Fatalf("same-instant advance = %s, want exactly 1µs (%s -> %s)", got, first.UpdatedAt, second.UpdatedAt)
	}

	// Wall-clock rollback must not move the column backwards.
	name3 := "Mono 3"
	third, err := repository.UpdateUser("u-mono", UserPatch{Name: &name3}, "actor", base.Add(-time.Hour))
	if err != nil {
		t.Fatalf("rollback UpdateUser: %v", err)
	}
	if !third.UpdatedAt.After(second.UpdatedAt) {
		t.Fatalf("rollback update moved updated_at backwards: %s then %s", second.UpdatedAt, third.UpdatedAt)
	}

	// P-11-MONO: roles.updated_at carries the same invariant.
	role, err := repository.CreateRoleWithGrants("mono_role", "Mono Role", nil, nil, base)
	if err != nil {
		t.Fatalf("CreateRoleWithGrants: %v", err)
	}
	if !role.UpdatedAt.Equal(base) {
		t.Fatalf("created role updated_at = %s, want %s", role.UpdatedAt, base)
	}
	renamed1 := "Mono Role 2"
	firstRole, err := repository.UpdateRoleWithGrants(role.ID, RolePatch{Name: &renamed1}, base)
	if err != nil {
		t.Fatalf("first UpdateRoleWithGrants: %v", err)
	}
	renamed2 := "Mono Role 3"
	secondRole, err := repository.UpdateRoleWithGrants(role.ID, RolePatch{Name: &renamed2}, base)
	if err != nil {
		t.Fatalf("second UpdateRoleWithGrants: %v", err)
	}
	if got := secondRole.UpdatedAt.Sub(firstRole.UpdatedAt); got != time.Microsecond {
		t.Fatalf("same-instant role advance = %s, want exactly 1µs (%s -> %s)", got, firstRole.UpdatedAt, secondRole.UpdatedAt)
	}
}
