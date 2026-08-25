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
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
)

// createR2Fixture builds a database shaped exactly like the pre-migration R2
// store (users + refresh_tokens, no schema_migrations) so OpenSeeded() has to run the
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
	st, err := OpenSeeded(path, "admin", "hash", true)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	authRepository := authsession.NewRepository(st)
	operationRepository := operationlog.NewRepository(st)
	applied, err := st.appliedMigrations()
	if err != nil {
		t.Fatalf("applied: %v", err)
	}
	if len(applied) != 56 || applied[0].version != 1 || applied[1].version != 2 || applied[2].version != 3 || applied[3].version != 4 || applied[4].version != 5 || applied[5].version != 6 || applied[6].version != 7 || applied[7].version != 8 || applied[8].version != 9 || applied[9].version != 10 || applied[10].version != 11 || applied[11].version != 12 || applied[12].version != 13 || applied[13].version != 14 || applied[14].version != 15 || applied[15].version != 16 || applied[16].version != 17 || applied[17].version != 18 || applied[18].version != 19 || applied[19].version != 20 || applied[20].version != 21 || applied[21].version != 22 || applied[22].version != 23 || applied[23].version != 24 || applied[24].version != 25 || applied[25].version != 26 || applied[26].version != 27 || applied[27].version != 28 || applied[28].version != 29 || applied[29].version != 30 || applied[30].version != 31 || applied[31].version != 32 || applied[32].version != 33 || applied[33].version != 34 || applied[34].version != 35 || applied[35].version != 36 || applied[36].version != 37 || applied[36].name != "notifications_message_keys" || applied[37].version != 38 || applied[37].name != "must_change_password" || applied[38].version != 39 || applied[38].name != "dict_entry_badge_style" || applied[39].version != 40 || applied[39].name != "site_footer" || applied[40].version != 41 || applied[40].name != "operation_log_correlation" || applied[41].version != 42 || applied[41].name != "async_jobs" || applied[42].version != 43 || applied[42].name != "operation_log_wallet_jobs" || applied[43].version != 44 || applied[43].name != "service_credentials" || applied[44].version != 45 || applied[44].name != "operation_log_service_credentials" || applied[45].version != 46 || applied[45].name != "site_operation_log_retention" || applied[46].version != 47 || applied[46].name != "operation_log_archive" || applied[47].version != 48 || applied[47].name != "operation_log_session" || applied[48].version != 49 || applied[48].name != "seed_admin_must_change_password" || applied[49].version != 50 || applied[49].name != "wallet_ledger_order_repair" || applied[50].version != 51 || applied[50].name != "mail_outbox" || applied[51].version != 52 || applied[51].name != "mail_config" || applied[52].version != 53 || applied[52].name != "operation_log_mail_events" || applied[53].version != 54 || applied[53].name != "account_email_identity" || applied[54].version != 55 || applied[54].name != "email_verification_challenges" || applied[55].version != 56 || applied[55].name != "password_recovery_challenges" {
		t.Fatalf("applied = %+v, want versions [1..56] ending in password_recovery_challenges", applied)
	}
	for _, tbl := range []string{
		"users", "refresh_tokens", "schema_migrations",
		"roles", "user_roles", "permissions", "role_permissions", "menu_items", "role_menu_items",
		"operation_log", "operation_log_correlation", "operation_log_archive", "operation_log_session", "site_settings", "system_data_reconcile", "system_data_grants", "jobs", "service_credentials",
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
	if err := operationRepository.RecordOperation(operationlog.Operation{
		ID: "op-fresh", Event: operationlog.EventUserCreate, ActorID: "user-admin", ActorName: "Admin",
		CreatedAt: time.Now().UTC(),
	}); err != nil {
		t.Fatalf("users.create on fresh operation_log: %v", err)
	}
	if err := operationRepository.RecordOperation(operationlog.Operation{
		ID: "op-settings", Event: operationlog.EventSettingsUpdate, ActorID: "user-admin", ActorName: "Admin",
		CreatedAt: time.Now().UTC(),
	}); err != nil {
		t.Fatalf("settings.update on fresh operation_log: %v", err)
	}
	// A fresh empty DB has nothing to recover: no snapshot should exist for ANY
	// version. This globs all pre-vN snapshots, not just pre-v0002 — mid-batch
	// system-data seeding used to trip the dbHasRows guard and produce ~39
	// wasteful full-database copies per fresh open (the full handler test
	// timeout); asserting the whole pattern locks that regression out.
	if snaps, _ := filepath.Glob(path + ".pre-v*.sqlite"); len(snaps) != 0 {
		t.Fatalf("fresh DB produced snapshots %v", snaps)
	}
	u, err := authRepository.UserByUsername("admin")
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
	st2, err := OpenSeeded(path, "admin", "hash-v2", true)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer st2.Close()
	authRepository2 := authsession.NewRepository(st2)
	u2, err := authRepository2.UserByUsername("admin")
	if err != nil {
		t.Fatalf("admin after reopen: %v", err)
	}
	if u2.PasswordHash != "hash" {
		t.Fatalf("password_hash = %q after reopen, want hash (seed must be no-op)", u2.PasswordHash)
	}
	applied2, _ := st2.appliedMigrations()
	if len(applied2) != 56 {
		t.Fatalf("migrations re-applied on reopen: %v", applied2)
	}
	// Reopen of a fully-migrated fresh DB: no new snapshots for any version.
	if snaps, _ := filepath.Glob(path + ".pre-v*.sqlite"); len(snaps) != 0 {
		t.Fatalf("reopen produced snapshots %v", snaps)
	}
}

// V-MIG-02 + V-REC-01/02 · an existing R2 DB is fingerprinted and registered by
// 0001, snapshotted before 0002, upgraded with backfill, and verified.
func TestMigrateExistingR2DB(t *testing.T) {
	path := filepath.Join(t.TempDir(), "existing.db")
	createR2Fixture(t, path)

	st, err := OpenSeeded(path, "admin", "hash", false)
	if err != nil {
		t.Fatalf("open existing R2 DB: %v", err)
	}
	defer st.Close()
	authRepository := authsession.NewRepository(st)

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
	u, err := authRepository.UserByID("user-admin")
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
	rt, err := authRepository.RefreshTokenByHash("abc123")
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

	st, err := OpenSeeded(path, "admin", "hash", false)
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

	st, err := OpenSeeded(path, "admin", "hash", false)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer st.Close()
	authRepository := authsession.NewRepository(st)

	want := []string{"admin", "editor"}
	u, err := authRepository.UserByID("user-admin")
	if err != nil {
		t.Fatalf("UserByID after migration: %v", err)
	}
	if !reflect.DeepEqual(u.Roles, want) {
		t.Fatalf("UserByID roles = %v, want %v (deduped, sorted by key)", u.Roles, want)
	}
	u2, err := authRepository.UserByUsername("admin")
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
	st, err := OpenSeeded(path, "admin", "hash", false)
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

	if _, err := OpenSeeded(path, "admin", "hash", false); err == nil {
		t.Fatal("expected fail closed for unknown applied version")
	}
}

// V-MIG-03 · a missing intermediate version fails closed.
func TestMigrateFailClosedMissingIntermediate(t *testing.T) {
	path := filepath.Join(t.TempDir(), "gap.db")
	st, err := OpenSeeded(path, "admin", "hash", false)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	st.Close()

	db := rawOpen(t, path)
	if _, err := db.Exec(`DELETE FROM schema_migrations WHERE version = 1`); err != nil {
		t.Fatalf("delete version 1: %v", err)
	}
	db.Close()

	if _, err := OpenSeeded(path, "admin", "hash", false); err == nil {
		t.Fatal("expected fail closed for missing intermediate version")
	}
}

// A-001 F-001 (S6 闭合) · a true middle gap — ledger {1,3} with 2 absent — is
// detected by the `a.version != applied[i-1].version+1` branch (the existing
// MissingIntermediate test only exercises the "ledger does not start at 1" path).
func TestMigrateFailClosedMissingMiddle(t *testing.T) {
	path := filepath.Join(t.TempDir(), "middlegap.db")
	st, err := OpenSeeded(path, "admin", "hash", false)
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

	if _, err := OpenSeeded(path, "admin", "hash", false); err == nil {
		t.Fatal("expected fail closed for ledger {1,3} with a missing middle version")
	}
}

// V-MIG-03 · a ledger checksum that no longer matches the code fails closed.
func TestMigrateFailClosedChecksumDrift(t *testing.T) {
	path := filepath.Join(t.TempDir(), "drift.db")
	st, err := OpenSeeded(path, "admin", "hash", false)
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

	if _, err := OpenSeeded(path, "admin", "hash", false); err == nil {
		t.Fatal("expected fail closed for checksum drift")
	}
}

func TestMigrateRestoresLostLedgerSQLite(t *testing.T) {
	path := filepath.Join(t.TempDir(), "lostled.db")
	st, err := OpenSeeded(path, "admin", "hash", true)
	if err != nil {
		t.Fatalf("bootstrap: %v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	db := rawOpen(t, path)
	if _, err := db.Exec(`DROP TABLE schema_migrations`); err != nil {
		db.Close()
		t.Fatalf("drop ledger: %v", err)
	}
	db.Close()

	st2, err := OpenSeeded(path, "admin", "hash", false)
	if err != nil {
		t.Fatalf("reopen after lost ledger: %v", err)
	}
	defer st2.Close()
	applied, err := st2.appliedMigrations()
	if err != nil {
		t.Fatal(err)
	}
	if len(applied) != len(compiledMigrations) {
		t.Fatalf("restored ledger rows = %d, want %d", len(applied), len(compiledMigrations))
	}
	var username string
	if err := st2.db.QueryRow(`SELECT username FROM users WHERE username = 'admin'`).Scan(&username); err != nil {
		t.Fatalf("preserved admin: %v", err)
	}
}

// V-MIG-03 · a foreign sqlite file (no schema-ui users) is refused and leaves no ledger.
func TestMigrateFailClosedForeignSQLite(t *testing.T) {
	path := filepath.Join(t.TempDir(), "foreign.db")
	db := rawOpen(t, path)
	if _, err := db.Exec(`CREATE TABLE orders (id TEXT PRIMARY KEY)`); err != nil {
		t.Fatalf("fixture: %v", err)
	}
	db.Close()

	if _, err := OpenSeeded(path, "admin", "hash", false); err == nil {
		t.Fatal("expected fail closed for foreign sqlite")
	} else if !strings.Contains(err.Error(), "identity=foreign") {
		t.Fatalf("want identity=foreign, got %v", err)
	}
	check := rawOpen(t, path)
	defer check.Close()
	if tableExistsDB(t, check, "schema_migrations") {
		t.Fatal("foreign sqlite must not create a migration ledger")
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

	if _, err := OpenSeeded(path, "admin", "hash", false); err == nil {
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

			if _, err := OpenSeeded(path, "admin", "hash", false); err == nil {
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
	st, err := OpenSeeded(filepath.Join(t.TempDir(), "fk.db"), "admin", "hash", false)
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
	st, err := OpenSeeded(filepath.Join(t.TempDir(), "rbac-mig.db"), "admin", "hash", true)
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
		{"admin.settings", "site_settings_v2", "b593aa2dd003e1339710b35478c87b105e6bb1762be0b4b08f3c986a5063a047"},
		{"core.auth-session", "access_token_revocation", "c3ea720aa0d0f10c67ee0ea734fe439db928de70a82b87c265391214dbce4688"},
		{"core.auth-session", "account_lock", "b9039118ebf3444bc2309ea481daac7ffdb1c0d627b4621642b5803c4fa3deb4"},
		{"admin.account", "account_enable_state", "ca2b3f38793d54c4d440f8e5af034dbed6cf32e50240304e9115d538ea05c539"},
		{"core.operationlog", "operation_log_account_events", "643d0f44bdf0f62a7689a0c4bdae2b3b511b077c5f7291db1a08763af389567d"},
		{"core.operationlog", "operation_log_data_transfer", "d8aa6a97fe57978126e7d06d03dbfcb1bb529c8ff680199e61f76791978f24b9"},
		{"admin.notifications", "notifications", "70d1357dc638fec1b100e0e9287dba46375a387914b884a899294f273f0fcbaa"},
		{"admin.notifications", "notifications_enabled", "16191741bded09598f2628da66f7b6251ef3886e616d839b5fd8b325a10d7264"},
		{"core.operationlog", "operation_log_file_events", "3351b6e6993dea21abd85f96049483eb8d9cfea4ad45bef34d3f5a824ac49249"},
		{"admin.data-dictionary", "dictionary", "8f2c2a18037a4c3eb67704354bbdd01acbbe319b477731ced054070aec6f6587"},
		{"core.operationlog", "operation_log_dictionary", "ac6d2c28eb213f90030b13f6e22bda470e4d7fa2c072d55939cda00e77ffd059"},
		{"admin.scheduled-tasks", "scheduled_tasks", "076c8fa37fb28deecc3737ac1886e2eb98261081b1ee0188da11c4bfd2e44286"},
		{"core.operationlog", "operation_log_tasks", "cb64c05e37e247127a0b4466a0812e113a96c0fce8f127bdd5ca6e5890326b04"},
		{"admin.login-captcha", "login_captcha", "6bc0b55675f4d00a231de3a0d1f89f52688e3efb292ab3cc5975c225c308b853"},
		{"core.operationlog", "operation_log_captcha", "51d6a3c11031d4fd2b3ff08573030e32747fff51b88d614af86e66a13548e264"},
		{"admin.recycle-bin", "recycle_items", "39087fe4f00cf1832cc5f1021cfd8aa0f2d153c9866a28732da17a6eae2f8408"},
		{"core.operationlog", "operation_log_recycle", "681f3bdce9f7a1fa7956849823d7e599e963638474b9d89b42151410a9dcb361"},
		{"admin.data-permission", "data_permission", "f3ce4c717b2f7c43090183c53f05cc1ecbd2299622d7c4ad06296f0c1f1a4318"},
		{"core.operationlog", "operation_log_data_permission", "a18c42e714030251c8f8eb7a1dec1b098e56f42138664f688991891bffbfc2e8"},
		{"admin.mfa", "user_mfa", "daa2592f09da53ccc77062d27c2f1d5bdf590ffa5886af62fe323318331764d3"},
		{"core.operationlog", "operation_log_mfa", "70464abf1bc3e9eb4dac3c836efaabc4f046610fba62c51486555cba91c40b9d"},
		// S-14 (GOAL-019): admin.wallet tables (0031) + operationlog wallet events (0032).
		{"admin.wallet", "wallet", "bc92082f6fadfc1812f16037275685f7901c1bb026334dc6f1256fac6db7358f"},
		{"core.operationlog", "operation_log_wallet", "1c27e86cbf362cdd08e7721fbfe416f21e2752bc84fa688253ef7000232f74dc"},
		// GOAL-021 (D-001 §3): ledger CHECK rebuild (deduct_frozen) + operationlog event.
		{"admin.wallet", "wallet_ledger_deduct", "b3135b2888dec0aa6da032121026e1ebfa07bed8bc75396bdc60166a09b3077d"},
		{"core.operationlog", "operation_log_wallet_deduct", "b6b54bee8b1baff9b5c8222a6619074ef3be54f211c55bceabe96f1c3a291467"},
		// W13 T-05 (GOAL-014): account self-service avatar column + avatar event.
		{"admin.account", "account_avatar_url", "a2872e8d7142a955851092b294f548394f2aea88567edbe9f486bd78b32c9f2a"},
		{"core.operationlog", "operation_log_avatar_events", "3f4a67b35244036f081728891c6b942abdab25d55b74da9a1b0095a877afb35c"},
		// W14 F-04 (GOAL-016): notification i18n message keys.
		{"admin.notifications", "notifications_message_keys", "4f0a99c0e14940e3df488cc031af240161c9b5d7843920f078df08a1c22159a0"},
		// W16 F-01 (GOAL-025): forced initial-password-change flag.
		{"core.auth-session", "must_change_password", "922df3e58653dcc491b96b11cc86217080d92958193aa95eb5cde55117a9e47d"},
		// W16 F-09 (GOAL-027): dictionary entry badge style.
		{"admin.data-dictionary", "dict_entry_badge_style", "b1a53caae83bfb7bac824185059830014d1beab916bf5d69d6e134201acdf987"},
		// W16 F-10 (GOAL-027): site footer copyright/ICP columns.
		{"admin.settings", "site_footer", "5277f8b095001e658958c47f597f8f31e9869bbeced572b94c65fdc8829b2aba"},
		// VP-012 R1: correlation IDs are persisted separately from operation detail.
		{"core.operationlog", "operation_log_correlation", "17922088b0b6911c47ad338faccf40bcd0654d500a4277fbe392ece62293842a"},
		// VP-012 R4: migration-only core.jobs durable state machine.
		{"core.jobs", "async_jobs", "55e1d3f88de080bd0b6015841e76f1ce32604444619d180a3b228123f99dec68"},
		// VP-012 R4: wallet Job enqueue and terminal audit events.
		{"core.operationlog", "operation_log_wallet_jobs", "210a2ebe080081478242c3b1fd4a348f9c0f3922dbd1d98f4b31c8c86345c24e"},
		// VP-012 R6: service credential storage and lifecycle audit events.
		{"core.auth-session", "service_credentials", "5c51ade23d7f97e8142b3d1d53b2d084438b8e1a7dc08c3739b000010e4b8f5d"},
		{"core.operationlog", "operation_log_service_credentials", "4cacb1a9262aedaf68df5b38b0294cb58d0735d28454b2a7640805c6b91766b6"},
		{"admin.settings", "site_operation_log_retention", "5038cd0dc684f801409db24bc61eba671d3f12327bc3615d086f58892b33e14b"},
		{"core.operationlog", "operation_log_archive", "6228b4e840bf28ce9afc534435dd201822b80f79474fd448522880a355774d40"},
		{"core.operationlog", "operation_log_session", "1427328e3942b8bddf0d0970ac173d72a4bbefeee427f32a819e16cb3935edf5"},
		// W22 A2: seed admin must_change_password backfill for upgraded databases.
		{"core.auth-session", "seed_admin_must_change_password", "b2c6bbf113733ce7fb89933ac33735eb4dc448112ef99ef41d519fb88badbc32"},
		// GOAL-037 / F-008: 0050 data-only ledger-order repair for legacy
		// same-millisecond disordered entries.
		{"admin.wallet", "wallet_ledger_order_repair", "835902cb80352790f36c56bb57ce186071dd26dac50c13816b2f401e9c340720"},
		// VP-017 R6 (workspace-017 GOAL-007; GOAL-006 D-002 §3): mock-channel
		// outbound record table.
		{"core.persistence", "mail_outbox", "4b42ffe91c15ed143237f3cefa5c6be7227efcb66ffd718ce05753843cf0187c"},
		// VP-017 R7 (workspace-017 GOAL-008; Root D-007): single-row runtime
		// channel state (secrets stored encrypted).
		{"core.persistence", "mail_config", "51633d08fffa540eefcebd6b14e551c9c35506af587fad2cd9b0be3687906887"},
		// VP-017 R7 (GOAL-008): operation_log event enum gains the mail admin events.
		{"core.operationlog", "operation_log_mail_events", "b447a4b6249893390b111fdd359178766bc48f03771b0ef30d3cddd1797fd075"},
		// workspace-018 R2 (GOAL-003 D-001): account email identity columns +
		// lower(email) unique expression index.
		{"core.auth-session", "account_email_identity", "f9a0bc654dffece5610e30097c04730654a7e9b40f4bdbe253ab04ec87032b0b"},
		// workspace-018 R3 (GOAL-004 D-001 §1): per-user active verification
		// challenge table (bind-reserves-slot delivery state).
		{"core.auth-session", "email_verification_challenges", "1556bda28a7fb995807eea2b376a35ea79cf497fa76c73171b2973304ce5b754"},
		// workspace-019 R2 (GOAL-003 D-001 §1): per-user active self-recovery
		// challenge table (forgot-password delivery state).
		{"core.auth-session", "password_recovery_challenges", "e19db1a293a013e801abbf60c47be55174dec7f8722fd2a5cd05eb971b4520c3"},
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
