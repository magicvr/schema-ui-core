package authsession

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

// W9 A-005 R-F-003: regression lock for the F-004 atomic failure accounting —
// the counter increments by one per failure, reaching the threshold opens the
// lock window AND resets the counter, and the next cycle counts from zero.
// (SQLite serializes writers, so lost updates cannot hide here; the atomic
// UPDATE shape keeps postgres READ COMMITTED equally correct.)
func TestRecordLoginFailureThresholdAndReset(t *testing.T) {
	hash := "test-hash"
	st, err := testsupport.OpenStore(filepath.Join(t.TempDir(), "test.db"), "admin", hash, true)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	repo := NewRepository(st)

	seed, err := repo.UserByUsername("admin")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	lockUntil := now.Add(15 * time.Minute)

	for i := 1; i <= 4; i++ {
		locked, err := repo.RecordLoginFailure(seed.ID, 5, lockUntil, now.Add(time.Duration(i)*time.Second))
		if err != nil {
			t.Fatalf("failure %d: %v", i, err)
		}
		if locked {
			t.Fatalf("failure %d opened the lock early", i)
		}
		u, err := repo.UserByID(seed.ID)
		if err != nil {
			t.Fatal(err)
		}
		if u.FailedLoginCount != i {
			t.Fatalf("count after failure %d = %d, want %d (lost update?)", i, u.FailedLoginCount, i)
		}
		if u.LockedUntil != 0 {
			t.Fatalf("locked_until set before threshold: %d", u.LockedUntil)
		}
	}

	locked, err := repo.RecordLoginFailure(seed.ID, 5, lockUntil, now.Add(5*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	if !locked {
		t.Fatal("5th failure did not open the lock window")
	}
	u, err := repo.UserByID(seed.ID)
	if err != nil {
		t.Fatal(err)
	}
	if u.LockedUntil == 0 || u.FailedLoginCount != 0 {
		t.Fatalf("after lock: locked_until=%d count=%d, want set/0", u.LockedUntil, u.FailedLoginCount)
	}

	// The counter restarts from zero after a lock opens.
	locked, err = repo.RecordLoginFailure(seed.ID, 5, lockUntil, now.Add(6*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	if locked {
		t.Fatal("first failure after lock re-opened immediately")
	}
	u, _ = repo.UserByID(seed.ID)
	if u.FailedLoginCount != 1 {
		t.Fatalf("count after post-lock failure = %d, want 1", u.FailedLoginCount)
	}

	// Unknown user fails closed with ErrNotFound (no silent zero-row success).
	if _, err := repo.RecordLoginFailure("user-missing", 5, lockUntil, now); !errors.Is(err, ErrNotFound) {
		t.Fatalf("unknown user err = %v, want ErrNotFound", err)
	}
}
