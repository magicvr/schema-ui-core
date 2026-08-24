// workspace-018 R2 · GOAL-003: migration 0054 account_email_identity.
// Verifies the frozen R1 contract (GOAL-002 D-001 §1/§2/§3/§6) at the schema
// layer:
//   a) upgrade path: pre-0054 rows land as (email=NULL, email_status=NULL);
//   b) fresh catalog: both columns and the lower(email) unique index exist;
//   c) bind-reserves-slot semantics: case-folded duplicates rejected by the
//      index, multiple NULL emails coexist, CHECK rejects unknown statuses.
package store

import (
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func mustInsertUser(t *testing.T, st *Store, id, username, email, emailStatus string) {
	t.Helper()
	now := time.Now().UTC().Unix()
	if _, err := st.db.Exec(
		`INSERT INTO users (id, username, name, roles, password_hash, email, email_status, created_at, updated_at)
		 VALUES (?, ?, 'User', '["viewer"]', 'hash', ?, ?, ?, ?)`,
		id, username, nullableString(email), nullableString(emailStatus), now, now,
	); err != nil {
		t.Fatalf("insert %s: %v", id, err)
	}
}

func nullableString(s string) any {
	if s == "" {
		return nil
	}
	return s
}

// TestMigrate0054UpgradePathLandsUnboundRows simulates a pre-0054 install:
// rows created before the migration must read back as (NULL, NULL) — i.e.
// unbound — after the full catalog applies.
func TestMigrate0054UpgradePathLandsUnboundRows(t *testing.T) {
	path := filepath.Join(t.TempDir(), "pre-0054.db")
	catalog := MigrationCatalog()

	st, err := OpenWithCatalog(path, catalog[:53])
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC().Unix()
	if _, err := st.db.Exec(
		`INSERT INTO users (id, username, name, roles, password_hash, created_at, updated_at)
		 VALUES ('legacy-user', 'legacy', 'Legacy', '["viewer"]', 'hash', ?, ?)`,
		now, now,
	); err != nil {
		t.Fatalf("insert legacy user: %v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	st, err = OpenWithCatalog(path, catalog)
	if err != nil {
		t.Fatalf("open with full catalog: %v", err)
	}
	defer st.Close()

	var email, status *string
	if err := st.db.QueryRow(
		`SELECT email, email_status FROM users WHERE id = 'legacy-user'`,
	).Scan(&email, &status); err != nil {
		t.Fatalf("scan legacy identity columns: %v", err)
	}
	if email != nil || status != nil {
		t.Fatalf("legacy row = (%v, %v), want (NULL, NULL)", email, status)
	}
}

func assertEmailIdentityObjects(t *testing.T, st *Store) {
	t.Helper()
	got := map[string]string{}
	rows, err := st.db.Query(`PRAGMA table_info(users)`)
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		var cid, notNull, pk int
		var name, typ string
		var dflt any
		if err := rows.Scan(&cid, &name, &typ, &notNull, &dflt, &pk); err != nil {
			rows.Close()
			t.Fatal(err)
		}
		got[name] = strings.ToUpper(typ)
	}
	rows.Close()
	if got["email"] != "TEXT" || got["email_status"] != "TEXT" {
		t.Fatalf("users identity columns = %v, want email/email_status TEXT", got)
	}

	var n int
	if err := st.db.QueryRow(
		`SELECT COUNT(*) FROM sqlite_master WHERE type = 'index' AND name = 'idx_users_email_lower' AND tbl_name = 'users'`,
	).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("idx_users_email_lower rows = %d, want 1", n)
	}
}

// TestMigrate0054FreshCatalogSemantics verifies the contract semantics on a
// fresh fully-migrated database: object shape plus bind-reserves-slot
// uniqueness (case-folded), NULL coexistence, and the CHECK guard.
func TestMigrate0054FreshCatalogSemantics(t *testing.T) {
	path := filepath.Join(t.TempDir(), "fresh-0054.db")
	st, err := OpenWithCatalog(path, MigrationCatalog())
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	assertEmailIdentityObjects(t, st)

	// Bind reserves the slot.
	mustInsertUser(t, st, "u-alice", "alice", "Alice@Example.COM", "pending")

	// Same address differing only by case must be rejected by the
	// lower(email) unique index.
	if _, err := st.db.Exec(
		`INSERT INTO users (id, username, name, roles, password_hash, email, email_status, created_at, updated_at)
		 VALUES ('u-alice-lower', 'alice-lower', 'User', '["viewer"]', 'hash', 'alice@example.com', 'verified', ?, ?)`,
		time.Now().UTC().Unix(), time.Now().UTC().Unix(),
	); err == nil {
		t.Fatal("case-insensitive duplicate email was accepted; want unique-index rejection")
	} else if !strings.Contains(strings.ToUpper(err.Error()), "UNIQUE") {
		t.Fatalf("duplicate-email error %v, want a UNIQUE violation", err)
	}

	// Multiple accounts without an email coexist (NULLs mutually distinct).
	mustInsertUser(t, st, "u-null-1", "nullone", "", "")
	mustInsertUser(t, st, "u-null-2", "nulltwo", "", "")

	// Unknown verification status is rejected by the CHECK guard.
	if _, err := st.db.Exec(
		`INSERT INTO users (id, username, name, roles, password_hash, email, email_status, created_at, updated_at)
		 VALUES ('u-bogus', 'bogus', 'User', '["viewer"]', 'hash', 'bogus@example.com', 'bogus', ?, ?)`,
		time.Now().UTC().Unix(), time.Now().UTC().Unix(),
	); err == nil {
		t.Fatal("email_status 'bogus' was accepted; want CHECK violation")
	}

	// Verified state is reachable and distinct addresses coexist.
	mustInsertUser(t, st, "u-carol", "carol", "Carol@Example.com", "verified")
}
