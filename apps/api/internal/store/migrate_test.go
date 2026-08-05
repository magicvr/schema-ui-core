package store

import (
	"database/sql"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// createR2Fixture builds a database shaped exactly like the pre-migration R2
// store (users + refresh_tokens, no schema_migrations) so Open() has to run the
// 0001 fingerprint/registration path.
func createR2Fixture(t *testing.T, path string) {
	t.Helper()
	createR2FixtureRoles(t, path, `["admin","editor"]`)
}

func createR2FixtureRoles(t *testing.T, path, rolesJSON string) {
	t.Helper()
	_ = os.Remove(path)
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	ddl := []string{
		`CREATE TABLE users (
  id            TEXT PRIMARY KEY,
  username      TEXT NOT NULL UNIQUE,
  name          TEXT NOT NULL,
  roles         TEXT NOT NULL,
  password_hash TEXT NOT NULL,
  created_at    INTEGER NOT NULL,
  updated_at    INTEGER NOT NULL
)`,
		`CREATE TABLE refresh_tokens (
  id         TEXT PRIMARY KEY,
  user_id    TEXT NOT NULL REFERENCES users(id),
  token_hash TEXT NOT NULL UNIQUE,
  expires_at INTEGER NOT NULL,
  revoked_at INTEGER,
  created_at INTEGER NOT NULL
)`,
		`CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id)`,
	}
	for _, s := range ddl {
		if _, err := db.Exec(s); err != nil {
			t.Fatalf("fixture ddl: %v", err)
		}
	}
	now := time.Now().UTC().Unix()
	if _, err := db.Exec(
		`INSERT INTO users (id, username, name, roles, password_hash, created_at, updated_at) VALUES (?,?,?,?,?,?,?)`,
		"user-admin", "admin", "Admin", rolesJSON, "hash-v1", now, now,
	); err != nil {
		t.Fatalf("fixture user: %v", err)
	}
	if _, err := db.Exec(
		`INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at) VALUES (?,?,?,?,?,?)`,
		"rt1", "user-admin", "abc123", now+3600, nil, now,
	); err != nil {
		t.Fatalf("fixture refresh token: %v", err)
	}
}

// upgradeR2ToV2 applies migrations 1 and 2 only to an R2 fixture, producing the
// exact {1,2} ledger state that existed before R4 — the upgrade-test input.
func upgradeR2ToV2(t *testing.T, path string) {
	t.Helper()
	db := rawOpen(t, path)
	s := &Store{db: db, path: path}
	for _, v := range []int{1, 2} {
		if err := s.applyMigration(compiledMigrations[v-1]); err != nil {
			db.Close()
			t.Fatalf("apply 000%d: %v", v, err)
		}
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
}

func rawOpen(t *testing.T, path string) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("raw open: %v", err)
	}
	db.SetMaxOpenConns(1)
	return db
}

func tableExistsDB(t *testing.T, db *sql.DB, name string) bool {
	t.Helper()
	var n int
	if err := db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?`, name).Scan(&n); err != nil {
		t.Fatalf("check table %s: %v", name, err)
	}
	return n == 1
}

// V-MIG-01 · fresh empty DB applies the compiled 0001-0009 history once;
// reopening is a no-op.
func TestMigrateFreshDB(t *testing.T) {
	path := filepath.Join(t.TempDir(), "fresh.db")
	st, err := Open(path, "admin", "hash", true)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	applied, err := st.appliedMigrations()
	if err != nil {
		t.Fatalf("applied: %v", err)
	}
	if len(applied) != 9 || applied[0].version != 1 || applied[1].version != 2 || applied[2].version != 3 || applied[3].version != 4 || applied[4].version != 5 || applied[5].version != 6 || applied[6].version != 7 || applied[7].version != 8 || applied[8].version != 9 {
		t.Fatalf("applied = %+v, want versions [1 2 3 4 5 6 7 8 9]", applied)
	}
	for _, tbl := range []string{
		"users", "refresh_tokens", "schema_migrations",
		"roles", "user_roles", "permissions", "role_permissions", "menu_items", "role_menu_items",
		"operation_log", "site_settings", "system_data_reconcile", "system_data_grants",
	} {
		if !tableExistsDB(t, st.db, tbl) {
			t.Fatalf("table %s missing after fresh migration", tbl)
		}
	}
	// 0006 records_retire drops the records table on a fresh install too
	// (GOAL-011 S3 · I-011-002 §5).
	if tableExistsDB(t, st.db, "records") {
		t.Fatal("records table must not exist after fresh migration (0006)")
	}
	// The expanded operation_log CHECK accepts users.* and settings.update (0008).
	if err := st.RecordOperation(Operation{
		ID: "op-fresh", Event: EventUserCreate, ActorID: "user-admin", ActorName: "Admin",
		CreatedAt: time.Now().UTC(),
	}); err != nil {
		t.Fatalf("users.create on fresh operation_log: %v", err)
	}
	if err := st.RecordOperation(Operation{
		ID: "op-settings", Event: EventSettingsUpdate, ActorID: "user-admin", ActorName: "Admin",
		CreatedAt: time.Now().UTC(),
	}); err != nil {
		t.Fatalf("settings.update on fresh operation_log: %v", err)
	}
	// A fresh empty DB has nothing to recover: no snapshot should exist.
	if snaps, _ := filepath.Glob(path + ".pre-v0002-*.sqlite"); len(snaps) != 0 {
		t.Fatalf("fresh DB produced snapshots %v", snaps)
	}
	u, err := st.UserByUsername("admin")
	if err != nil {
		t.Fatalf("seeded admin: %v", err)
	}
	if u.PasswordHash != "hash" {
		t.Fatalf("password_hash = %q, want hash", u.PasswordHash)
	}
	if err := st.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	// Reopen: migrations not re-applied, seed not overwritten, no new snapshot.
	st2, err := Open(path, "admin", "hash-v2", true)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer st2.Close()
	u2, err := st2.UserByUsername("admin")
	if err != nil {
		t.Fatalf("admin after reopen: %v", err)
	}
	if u2.PasswordHash != "hash" {
		t.Fatalf("password_hash = %q after reopen, want hash (seed must be no-op)", u2.PasswordHash)
	}
	applied2, _ := st2.appliedMigrations()
	if len(applied2) != 9 {
		t.Fatalf("migrations re-applied on reopen: %v", applied2)
	}
	if snaps, _ := filepath.Glob(path + ".pre-v0002-*.sqlite"); len(snaps) != 0 {
		t.Fatalf("reopen produced snapshots %v", snaps)
	}
}

// V-MIG-02 + V-REC-01/02 · an existing R2 DB is fingerprinted and registered by
// 0001, snapshotted before 0002, upgraded with backfill, and verified.
func TestMigrateExistingR2DB(t *testing.T) {
	path := filepath.Join(t.TempDir(), "existing.db")
	createR2Fixture(t, path)

	st, err := Open(path, "admin", "hash", false)
	if err != nil {
		t.Fatalf("open existing R2 DB: %v", err)
	}
	defer st.Close()

	// V-REC-01 · the pre-v0002 snapshot is a single recoverable copy.
	snaps, err := filepath.Glob(path + ".pre-v0002-*.sqlite")
	if err != nil || len(snaps) != 1 {
		t.Fatalf("snapshots = %v, want exactly 1 (err %v)", snaps, err)
	}
	snap := snaps[0]
	snapDB := rawOpen(t, snap)
	defer snapDB.Close()
	var integ string
	if err := snapDB.QueryRow(`PRAGMA integrity_check`).Scan(&integ); err != nil || integ != "ok" {
		t.Fatalf("snapshot integrity_check = %q, err %v", integ, err)
	}
	var name, roles string
	if err := snapDB.QueryRow(`SELECT name, roles FROM users WHERE id = 'user-admin'`).Scan(&name, &roles); err != nil {
		t.Fatalf("snapshot user query: %v", err)
	}
	if name != "Admin" || roles != `["admin","editor"]` {
		t.Fatalf("snapshot user = %s / %s, want Admin / [admin editor]", name, roles)
	}
	var th string
	if err := snapDB.QueryRow(`SELECT token_hash FROM refresh_tokens WHERE id = 'rt1'`).Scan(&th); err != nil || th != "abc123" {
		t.Fatalf("snapshot refresh token = %q, err %v", th, err)
	}
	if tableExistsDB(t, snapDB, "roles") {
		t.Fatalf("snapshot unexpectedly already contains RBAC tables")
	}

	// V-REC-02 · main DB is intact and fully verified after migration.
	if err := st.verifyIntegrity(); err != nil {
		t.Fatalf("post-migration integrity: %v", err)
	}
	u, err := st.UserByID("user-admin")
	if err != nil {
		t.Fatalf("user after upgrade: %v", err)
	}
	if len(u.Roles) != 2 || u.Roles[0] != "admin" {
		t.Fatalf("roles = %v, want [admin editor]", u.Roles)
	}
	var roleCount, urCount int
	if err := st.db.QueryRow(`SELECT COUNT(*) FROM roles`).Scan(&roleCount); err != nil {
		t.Fatal(err)
	}
	if err := st.db.QueryRow(`SELECT COUNT(*) FROM user_roles`).Scan(&urCount); err != nil {
		t.Fatal(err)
	}
	if roleCount != 2 || urCount != 2 {
		t.Fatalf("backfill roles=%d user_roles=%d, want 2/2", roleCount, urCount)
	}
	var key string
	if err := st.db.QueryRow(`SELECT key FROM roles WHERE id = 'role-admin'`).Scan(&key); err != nil || key != "admin" {
		t.Fatalf("role-admin key = %q, err %v", key, err)
	}
	rt, err := st.RefreshTokenByHash("abc123")
	if err != nil {
		t.Fatalf("refresh token after upgrade: %v", err)
	}
	if rt.UserID != "user-admin" {
		t.Fatalf("refresh token user = %q, want user-admin", rt.UserID)
	}
}

// backfill dedupes duplicate role keys within one user's roles list.
func TestMigrateExistingR2DedupeRoles(t *testing.T) {
	path := filepath.Join(t.TempDir(), "dedupe.db")
	createR2FixtureRoles(t, path, `["admin","admin","editor"]`)

	st, err := Open(path, "admin", "hash", false)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer st.Close()
	var roleCount, urCount int
	st.db.QueryRow(`SELECT COUNT(*) FROM roles`).Scan(&roleCount)
	st.db.QueryRow(`SELECT COUNT(*) FROM user_roles`).Scan(&urCount)
	if roleCount != 2 || urCount != 2 {
		t.Fatalf("after dedupe roles=%d user_roles=%d, want 2/2", roleCount, urCount)
	}
}

// A-002 F-004 · an existing R2 user whose legacy roles JSON contains duplicates
// is readable after migration: 0002 backfill dedupes the relations and the read
// comparator follows set semantics (I-006-001 §5), so UserByID / UserByUsername
// succeed and return the deduped roles sorted by key.
func TestMigrateExistingR2DuplicateRolesReadable(t *testing.T) {
	path := filepath.Join(t.TempDir(), "legacy-dup-read.db")
	createR2FixtureRoles(t, path, `["admin","admin","editor"]`)

	st, err := Open(path, "admin", "hash", false)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer st.Close()

	want := []string{"admin", "editor"}
	u, err := st.UserByID("user-admin")
	if err != nil {
		t.Fatalf("UserByID after migration: %v", err)
	}
	if !reflect.DeepEqual(u.Roles, want) {
		t.Fatalf("UserByID roles = %v, want %v (deduped, sorted by key)", u.Roles, want)
	}
	u2, err := st.UserByUsername("admin")
	if err != nil {
		t.Fatalf("UserByUsername after migration: %v", err)
	}
	if !reflect.DeepEqual(u2.Roles, want) {
		t.Fatalf("UserByUsername roles = %v, want %v (deduped, sorted by key)", u2.Roles, want)
	}
}

// V-MIG-03 · an unknown applied version fails closed.
func TestMigrateFailClosedUnknownVersion(t *testing.T) {
	path := filepath.Join(t.TempDir(), "unknown.db")
	st, err := Open(path, "admin", "hash", false)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	st.Close()

	db := rawOpen(t, path)
	if _, err := db.Exec(
		`INSERT INTO schema_migrations (version, name, checksum, applied_at) VALUES (99, 'future', ?, 1)`,
		strings.Repeat("0", 64),
	); err != nil {
		t.Fatalf("inject unknown version: %v", err)
	}
	db.Close()

	if _, err := Open(path, "admin", "hash", false); err == nil {
		t.Fatal("expected fail closed for unknown applied version")
	}
}

// V-MIG-03 · a missing intermediate version fails closed.
func TestMigrateFailClosedMissingIntermediate(t *testing.T) {
	path := filepath.Join(t.TempDir(), "gap.db")
	st, err := Open(path, "admin", "hash", false)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	st.Close()

	db := rawOpen(t, path)
	if _, err := db.Exec(`DELETE FROM schema_migrations WHERE version = 1`); err != nil {
		t.Fatalf("delete version 1: %v", err)
	}
	db.Close()

	if _, err := Open(path, "admin", "hash", false); err == nil {
		t.Fatal("expected fail closed for missing intermediate version")
	}
}

// A-001 F-001 (S6 闭合) · a true middle gap — ledger {1,3} with 2 absent — is
// detected by the `a.version != applied[i-1].version+1` branch (the existing
// MissingIntermediate test only exercises the "ledger does not start at 1" path).
func TestMigrateFailClosedMissingMiddle(t *testing.T) {
	path := filepath.Join(t.TempDir(), "middlegap.db")
	st, err := Open(path, "admin", "hash", false)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	st.Close()

	// Drop the middle row from the now-{1,2,3} ledger to build {1,3}.
	db := rawOpen(t, path)
	if _, err := db.Exec(`DELETE FROM schema_migrations WHERE version = 2`); err != nil {
		t.Fatalf("delete version 2: %v", err)
	}
	db.Close()

	if _, err := Open(path, "admin", "hash", false); err == nil {
		t.Fatal("expected fail closed for ledger {1,3} with a missing middle version")
	}
}

// V-MIG-03 · a ledger checksum that no longer matches the code fails closed.
func TestMigrateFailClosedChecksumDrift(t *testing.T) {
	path := filepath.Join(t.TempDir(), "drift.db")
	st, err := Open(path, "admin", "hash", false)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	st.Close()

	db := rawOpen(t, path)
	if _, err := db.Exec(
		`UPDATE schema_migrations SET checksum = ? WHERE version = 1`,
		strings.Repeat("a", 64),
	); err != nil {
		t.Fatalf("tamper checksum: %v", err)
	}
	db.Close()

	if _, err := Open(path, "admin", "hash", false); err == nil {
		t.Fatal("expected fail closed for checksum drift")
	}
}

// V-MIG-03 · a partial baseline (users only) is rejected and leaves no ledger.
func TestMigrateFailClosedPartialBaseline(t *testing.T) {
	path := filepath.Join(t.TempDir(), "partial.db")
	db := rawOpen(t, path)
	if _, err := db.Exec(`CREATE TABLE users (
  id            TEXT PRIMARY KEY,
  username      TEXT NOT NULL UNIQUE,
  name          TEXT NOT NULL,
  roles         TEXT NOT NULL,
  password_hash TEXT NOT NULL,
  created_at    INTEGER NOT NULL,
  updated_at    INTEGER NOT NULL
)`); err != nil {
		t.Fatalf("fixture: %v", err)
	}
	db.Close()

	if _, err := Open(path, "admin", "hash", false); err == nil {
		t.Fatal("expected fail closed for partial baseline (missing refresh_tokens)")
	}
	check := rawOpen(t, path)
	defer check.Close()
	if tableExistsDB(t, check, "schema_migrations") {
		t.Fatal("partial baseline must not leave an empty migration ledger")
	}
}

// V-MIG-03 · invalid or non-array roles roll 0002 back completely.
func TestMigrateFailClosedInvalidRoles(t *testing.T) {
	for name, rolesJSON := range map[string]string{
		"invalid-key": `["Bad Key!"]`,
		"non-array":   `"not-an-array"`,
	} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "badroles-"+name+".db")
			createR2FixtureRoles(t, path, rolesJSON)

			if _, err := Open(path, "admin", "hash", false); err == nil {
				t.Fatalf("expected fail closed for roles %s", rolesJSON)
			}
			check := rawOpen(t, path)
			defer check.Close()
			// 0001 committed (ledger exists, version 1), 0002 fully rolled back.
			if !tableExistsDB(t, check, "schema_migrations") {
				t.Fatal("0001 should have committed before 0002 failed")
			}
			if tableExistsDB(t, check, "roles") {
				t.Fatal("0002 must roll back completely: roles table present")
			}
			var v int
			if err := check.QueryRow(`SELECT MAX(version) FROM schema_migrations`).Scan(&v); err != nil || v != 1 {
				t.Fatalf("ledger max version = %d, err %v, want 1", v, err)
			}
		})
	}
}

// V-MIG-04 · foreign_keys is asserted ON for the store connection.
func TestForeignKeyEnabled(t *testing.T) {
	st, err := Open(filepath.Join(t.TempDir(), "fk.db"), "admin", "hash", false)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer st.Close()
	now := time.Now().UTC().Unix()
	_, err = st.db.Exec(
		`INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at) VALUES (?,?,?,?,?,?)`,
		"rt-x", "missing-user", "xyz", now+3600, nil, now,
	)
	if err == nil {
		t.Fatal("expected FK violation for refresh token referencing missing user")
	}
}

// A-001 F-002 (S6 闭合) · the full V-MIG-04 matrix: reverse indexes, unique and
// CHECK constraints, and CASCADE|RESTRICT delete semantics across the RBAC tables
// are asserted on the store connection (not just declared in DDL).
func TestRBACConstraintsAndIndexes(t *testing.T) {
	st, err := Open(filepath.Join(t.TempDir(), "rbac-mig.db"), "admin", "hash", true)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer st.Close()

	// Reverse indexes exist.
	for _, idx := range []string{
		"idx_user_roles_role_id", "idx_role_permissions_permission_id", "idx_role_menu_items_menu_item_id",
	} {
		var n int
		if err := st.db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type = 'index' AND name = ?`, idx).Scan(&n); err != nil || n != 1 {
			t.Fatalf("index %s = %d, err %v, want 1", idx, n, err)
		}
	}

	// Unique: roles.key, permissions.key, menu_items.page_ref / feature_key.
	if _, err := st.db.Exec(
		`INSERT INTO roles (id, key, name, system, created_at, updated_at) VALUES ('role-x', 'admin', 'X', 0, 1, 1)`,
	); err == nil {
		t.Fatal("expected unique violation for duplicate roles.key")
	}
	if _, err := st.db.Exec(
		`INSERT INTO permissions (id, key, description, created_at, updated_at) VALUES ('p-x', 'users.read', '', 1, 1)`,
	); err == nil {
		t.Fatal("expected unique violation for duplicate permissions.key")
	}
	if _, err := st.db.Exec(
		`INSERT INTO menu_items (id, page_ref, feature_key, sort_order, enabled, created_at, updated_at) VALUES ('m-x', 'users', 'fx_x', 0, 1, 1, 1)`,
	); err == nil {
		t.Fatal("expected unique violation for duplicate menu_items.page_ref")
	}
	if _, err := st.db.Exec(
		`INSERT INTO menu_items (id, page_ref, feature_key, sort_order, enabled, created_at, updated_at) VALUES ('m-y', 'other-page', 'menu_users', 0, 1, 1, 1)`,
	); err == nil {
		t.Fatal("expected unique violation for duplicate menu_items.feature_key")
	}
	// Primary keys: duplicate user_roles link is rejected.
	if _, err := st.db.Exec(
		`INSERT INTO user_roles (user_id, role_id) VALUES ('user-admin', 'role-admin')`,
	); err == nil {
		t.Fatal("expected PK violation for duplicate user_roles")
	}

	// CHECK: roles.system and menu_items.enabled are 0/1.
	if _, err := st.db.Exec(
		`INSERT INTO roles (id, key, name, system, created_at, updated_at) VALUES ('role-y', 'y', 'Y', 2, 1, 1)`,
	); err == nil {
		t.Fatal("expected CHECK violation for roles.system = 2")
	}
	if _, err := st.db.Exec(
		`INSERT INTO menu_items (id, page_ref, feature_key, sort_order, enabled, created_at, updated_at) VALUES ('m-z', 'z-page', 'fz_z', 0, 2, 1, 1)`,
	); err == nil {
		t.Fatal("expected CHECK violation for menu_items.enabled = 2")
	}

	// CASCADE: deleting a role (no user relations) cascades its permission grants.
	if _, err := st.db.Exec(`DELETE FROM roles WHERE id = 'role-viewer'`); err != nil {
		t.Fatalf("delete viewer role: %v", err)
	}
	var rp int
	if err := st.db.QueryRow(`SELECT COUNT(*) FROM role_permissions WHERE role_id = 'role-viewer'`).Scan(&rp); err != nil || rp != 0 {
		t.Fatalf("role_viewer grants after cascade = %d, err %v, want 0", rp, err)
	}

	// RESTRICT: deleting a permission / menu item still granted fails.
	if _, err := st.db.Exec(`DELETE FROM permissions WHERE id = 'perm-users-read'`); err == nil {
		t.Fatal("expected RESTRICT deleting an in-use permission")
	}
	if _, err := st.db.Exec(`DELETE FROM menu_items WHERE id = 'menu-users'`); err == nil {
		t.Fatal("expected RESTRICT deleting an in-use menu item")
	}
}

// R6 C6.2 slice 3: the compiled registry is the sole migration authority and
// carries the frozen 0001-0008 owner mapping into the platform runner.
func TestCompiledMigrationCatalogOwnership(t *testing.T) {
	catalog := MigrationCatalog()
	want := []struct {
		moduleID string
		name     string
		checksum string
	}{
		{"core.auth-session", "r2_baseline", "bde1a83172d99932e5a90b6653808c8c6c510bbc8be32fdb52da9686428f6ff4"},
		{"core.auth-session", "rbac_expand", "1a7f630d1916e69aa2901f30d30d7d0b58ad9e96879740e7836af4ca07acc4ee"},
		{"core.persistence", "records_persist", "b195f6c68fac904d7958e4b17add1f64e6c1a696cd9d9382681fa0f06d30af7d"},
		{"core.operationlog", "operation_log", "53fd75b5b827a9a02648b7fee3c6d4c3e0c8972bca8ecb7adac8519cd8d4cf72"},
		{"core.operationlog", "operation_log_expand", "73369d73a15d915e3338bfaa0ba1d7d8f5bf0e838493ac42fe295c1db8c02f02"},
		{"core.persistence", "records_retire", "175ac09f0c67658161a6852d2779781f59985488aa75308f3fa419a06c5f926b"},
		{"admin.settings", "site_settings", "6ffb1d0d978d7475ebd807f4dc1aab609d255186ddefa08e28a5398d265b7dfa"},
		{"core.operationlog", "operation_log_settings", "ec3635f99db24907eb4a371ebd8c8f328c80a69e07715b866e1bca319f518d6c"},
		{"core.auth-session", "system_data_reconcile", "3e1c1e6d95c1f94c38a17ead999ee2cda685ec1e78d2148b4d12111d1eca74b6"},
	}
	if len(catalog) != len(want) {
		t.Fatalf("catalog len = %d, want %d", len(catalog), len(want))
	}
	for index, migration := range catalog {
		if migration.Version != index+1 {
			t.Fatalf("catalog[%d].Version = %d, want %d", index, migration.Version, index+1)
		}
		expected := want[index]
		if migration.ModuleID != expected.moduleID || migration.Name != expected.name || migration.Key != expected.name {
			t.Fatalf("catalog[%d] identity = {%q %q %q}, want {%q %q %q}", index,
				migration.ModuleID, migration.Key, migration.Name, expected.moduleID, expected.name, expected.name)
		}
		if migration.Checksum != expected.checksum {
			t.Fatalf("catalog[%d] %s checksum = %s, want frozen %s", index, migration.Name, migration.Checksum, expected.checksum)
		}
		if migration.Tombstone || migration.Apply == nil {
			t.Fatalf("catalog[%d] %s must preserve its frozen executable Apply", index, migration.Name)
		}
	}
}

func TestOpenWithCatalogRejectsInvalidAndAppliedDrift(t *testing.T) {
	catalog := MigrationCatalog()
	if len(catalog) == 0 {
		t.Fatal("empty catalog")
	}
	invalid := append([]kernel.MigrationContribution(nil), catalog...)
	invalid[0].Checksum = "bad"
	if _, err := OpenWithCatalog(filepath.Join(t.TempDir(), "bad.db"), invalid); err == nil {
		t.Fatal("invalid checksum must fail closed")
	}
	missing := append([]kernel.MigrationContribution(nil), catalog[:3]...)
	missing = append(missing, catalog[4:]...)
	if _, err := OpenWithCatalog(filepath.Join(t.TempDir(), "gap.db"), missing); err == nil {
		t.Fatal("catalog version gap must fail closed")
	}

	path := filepath.Join(t.TempDir(), "ok.db")
	st, err := OpenWithCatalog(path, catalog)
	if err != nil {
		t.Fatalf("correct catalog open: %v", err)
	}
	_ = st.Close()

	drifted := append([]kernel.MigrationContribution(nil), catalog...)
	drifted[0].Checksum = "0" + drifted[0].Checksum[1:]
	if _, err := OpenWithCatalog(path, drifted); err == nil {
		t.Fatal("applied checksum drift must fail closed")
	}

	renamed := append([]kernel.MigrationContribution(nil), catalog...)
	renamed[0].Key = "r2_baseline_renamed"
	renamed[0].Name = renamed[0].Key
	if _, err := OpenWithCatalog(path, renamed); err == nil || !strings.Contains(err.Error(), "name") {
		t.Fatalf("applied name drift error = %v, want fail-closed name mismatch", err)
	}
}
