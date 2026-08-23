// A2 backfill · seed admin must_change_password=1
// Verifies migration 0049 for two scenarios:
//   a) upgrade path: old-schema admin row (must_change_password=0) is set to 1.
//   b) idempotency: running the full catalog twice leaves must_change_password=1
//      and does not fail.
package store

import (
	"path/filepath"
	"testing"
	"time"
)

// TestMigrate0049BackfillsSeedAdminMustChangePassword simulates an upgrade
// install: the database was created before migration 0038 added
// must_change_password (DEFAULT 0), so the seed admin row has value 0.
// After applying the full catalog (including 0049), the row must have value 1.
func TestMigrate0049BackfillsSeedAdminMustChangePassword(t *testing.T) {
	path := filepath.Join(t.TempDir(), "seed-admin-0048.db")
	catalog := MigrationCatalog()

	// Open with first 48 migrations (all versions up to and including 48).
	st, err := OpenWithCatalog(path, catalog[:48])
	if err != nil {
		t.Fatal(err)
	}

	// Simulate an old-install seed admin: id='user-admin', must_change_password=0.
	now := time.Now().UTC().Unix()
	if _, err := st.db.Exec(
		`INSERT INTO users (id, username, name, roles, password_hash, must_change_password, created_at, updated_at)
		 VALUES ('user-admin', 'admin', 'Admin', '["admin","editor"]', 'legacy-hash', 0, ?, ?)`,
		now, now,
	); err != nil {
		t.Fatalf("insert legacy seed admin: %v", err)
	}

	// Also insert a non-seed user that must NOT be touched.
	if _, err := st.db.Exec(
		`INSERT INTO users (id, username, name, roles, password_hash, must_change_password, created_at, updated_at)
		 VALUES ('user-operator', 'operator', 'Operator', '["editor"]', 'op-hash', 0, ?, ?)`,
		now, now,
	); err != nil {
		t.Fatalf("insert non-seed user: %v", err)
	}

	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	// Apply the full catalog (0049 is the last one).
	st, err = OpenWithCatalog(path, catalog)
	if err != nil {
		t.Fatalf("open with full catalog: %v", err)
	}
	defer st.Close()

	// Seed admin must now have must_change_password = 1.
	var adminFlag int
	if err := st.db.QueryRow(
		`SELECT must_change_password FROM users WHERE id = 'user-admin'`,
	).Scan(&adminFlag); err != nil {
		t.Fatalf("scan seed admin flag: %v", err)
	}
	if adminFlag != 1 {
		t.Errorf("seed admin must_change_password = %d after 0049, want 1", adminFlag)
	}

	// Non-seed user must be untouched (0 → 0).
	var opFlag int
	if err := st.db.QueryRow(
		`SELECT must_change_password FROM users WHERE id = 'user-operator'`,
	).Scan(&opFlag); err != nil {
		t.Fatalf("scan non-seed user flag: %v", err)
	}
	if opFlag != 0 {
		t.Errorf("non-seed user must_change_password = %d after 0049, want 0 (untouched)", opFlag)
	}
}

// TestMigrate0049IdempotentOnFreshDatabase verifies that 0049 is a no-op when
// the seed admin was created by the current bootstrap (must_change_password=1
// already) and that re-opening the same database does not change anything.
func TestMigrate0049IdempotentOnFreshDatabase(t *testing.T) {
	path := filepath.Join(t.TempDir(), "seed-admin-fresh.db")
	catalog := MigrationCatalog()

	// Fresh DB: all migrations including 0049.
	st, err := OpenWithCatalog(path, catalog)
	if err != nil {
		t.Fatalf("first open: %v", err)
	}

	// Bootstrap admin with must_change_password=1 (current behavior).
	now := time.Now().UTC().Unix()
	if _, err := st.db.Exec(
		`INSERT INTO users (id, username, name, roles, password_hash, must_change_password, created_at, updated_at)
		 VALUES ('user-admin', 'admin', 'Admin', '["admin","editor"]', 'new-hash', 1, ?, ?)`,
		now, now,
	); err != nil {
		t.Fatalf("insert fresh seed admin: %v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	// Re-open (0049 is already applied; apply body is a no-op UPDATE).
	st, err = OpenWithCatalog(path, catalog)
	if err != nil {
		t.Fatalf("second open: %v", err)
	}
	defer st.Close()

	var flag int
	if err := st.db.QueryRow(
		`SELECT must_change_password FROM users WHERE id = 'user-admin'`,
	).Scan(&flag); err != nil {
		t.Fatalf("scan flag after reopen: %v", err)
	}
	if flag != 1 {
		t.Errorf("must_change_password = %d after idempotent reopen, want 1", flag)
	}
}
