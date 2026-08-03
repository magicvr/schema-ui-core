// Migration runner for the SQLite auth store (GOAL-006 D-002 / I-006-001).
//
// R2 started with an idempotent `CREATE TABLE IF NOT EXISTS` blob in Open; R3
// replaces that with a versioned, checksummed ledger so that upgrades are
// ordered, auditable, and recoverable. The contract that callers must keep:
//
//   - Store.Open still returns after migrating + seeding with the same signature.
//   - Applied migration numbers form a contiguous prefix starting at 0001;
//     unknown versions, gaps, duplicate names and checksum drift all fail closed.
//   - Each migration runs in its own transaction (DDL + data transform + the
//     ledger insert all commit or roll back together).
//   - A non-empty file database is snapshotted (VACUUM INTO) before the first
//     data-mutating migration so an upgrade has a recoverable pre-state; the
//     snapshot is named after the first pending version (pre-v0002 / pre-v0003).
package store

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// schemaMigrationsDDL is the migration ledger. checksum is the SHA-256 (lower
// hex) of the canonical SQL + data-transformer id for the applied migration.
const schemaMigrationsDDL = `CREATE TABLE schema_migrations (
  version    INTEGER PRIMARY KEY,
  name       TEXT NOT NULL UNIQUE,
  checksum   TEXT NOT NULL CHECK (length(checksum) = 64),
  applied_at INTEGER NOT NULL
)`

// migration is one versioned, checksummed schema/data step.
type migration struct {
	version int
	name    string
	// stmts are executed in order inside up() on an empty database and are the
	// canonical checksum input. Editing them (or bumping transformID after a
	// data-transform change) changes the ledger checksum, so a database that
	// already recorded the version fails closed on next startup (drift).
	stmts       []string
	transformID string
	up          func(tx *sql.Tx) error
}

// appliedMigration is one row read back from the ledger.
type appliedMigration struct {
	version  int
	name     string
	checksum string
}

// compiledMigrations is the ordered, immutable migration set for this store.
// Versions must be strictly ascending and unique; names must be unique.
var compiledMigrations = []migration{
	{
		version:     1,
		name:        "r2_baseline",
		transformID: "0001:r2-baseline:v1",
		stmts:       r2BaselineDDL,
		up:          migrate0001,
	},
	{
		version:     2,
		name:        "rbac_expand",
		transformID: "0002:rbac-expand:v1",
		stmts:       rbacExpandDDL,
		up:          migrate0002,
	},
	{
		version:     3,
		name:        "records_persist",
		transformID: "0003:records-persist:v1",
		stmts:       recordsPersistDDL,
		up:          migrate0003,
	},
	{
		version:     4,
		name:        "operation_log",
		transformID: "0004:operation-log:v1",
		stmts:       operationLogDDL,
		up:          migrate0004,
	},
	{
		version:     5,
		name:        "operation_log_expand",
		transformID: "0005:operation-log-expand:v1",
		stmts:       operationLogExpandDDL,
		up:          migrate0005,
	},
	{
		version:     6,
		name:        "records_retire",
		transformID: "0006:records-retire:v1",
		stmts:       recordsRetireDDL,
		up:          migrate0006,
	},
}

// r2BaselineDDL is the canonical R2 schema (users + refresh_tokens). It is
// executed verbatim on an empty database and is the fingerprint contract for an
// existing R2 database that has no migration ledger yet.
var r2BaselineDDL = []string{
	`CREATE TABLE users (
  id            TEXT PRIMARY KEY,
  username      TEXT NOT NULL UNIQUE,
  name          TEXT NOT NULL,
  roles         TEXT NOT NULL, -- JSON array; R3 normalizes
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

// rbacExpandDDL is the canonical R3 normalized schema created by 0002.
// Delete semantics (I-006-001 §4): users cascade their role relations; a role
// must have no user relations before deletion, then cascades grants; a
// permission / menu item must have no grants before deletion.
var rbacExpandDDL = []string{
	`CREATE TABLE roles (
  id         TEXT PRIMARY KEY,
  key        TEXT NOT NULL UNIQUE CHECK (key <> ''),
  name       TEXT NOT NULL,
  system     INTEGER NOT NULL DEFAULT 0 CHECK (system IN (0, 1)),
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
)`,
	`CREATE TABLE user_roles (
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role_id TEXT NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
  PRIMARY KEY (user_id, role_id)
)`,
	`CREATE INDEX idx_user_roles_role_id ON user_roles(role_id)`,
	`CREATE TABLE permissions (
  id          TEXT PRIMARY KEY,
  key         TEXT NOT NULL UNIQUE CHECK (key <> ''),
  description TEXT NOT NULL DEFAULT '',
  created_at  INTEGER NOT NULL,
  updated_at  INTEGER NOT NULL
)`,
	`CREATE TABLE role_permissions (
  role_id       TEXT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
  permission_id TEXT NOT NULL REFERENCES permissions(id) ON DELETE RESTRICT,
  PRIMARY KEY (role_id, permission_id)
)`,
	`CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id)`,
	`CREATE TABLE menu_items (
  id          TEXT PRIMARY KEY,
  page_ref    TEXT NOT NULL UNIQUE CHECK (page_ref <> ''),
  feature_key TEXT NOT NULL UNIQUE CHECK (feature_key <> ''),
  sort_order  INTEGER NOT NULL DEFAULT 0,
  enabled     INTEGER NOT NULL DEFAULT 1 CHECK (enabled IN (0, 1)),
  created_at  INTEGER NOT NULL,
  updated_at  INTEGER NOT NULL
)`,
	`CREATE TABLE role_menu_items (
  role_id      TEXT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
  menu_item_id TEXT NOT NULL REFERENCES menu_items(id) ON DELETE RESTRICT,
  PRIMARY KEY (role_id, menu_item_id)
)`,
	`CREATE INDEX idx_role_menu_items_menu_item_id ON role_menu_items(menu_item_id)`,
}

// recordsPersistDDL is the canonical R4 records schema created by 0003
// (GOAL-007 D-003 / I-007-002 §2). updated_at is Unix milliseconds (D-004); the
// trim-non-empty CHECKs are a DB backstop on top of handler-level validation.
var recordsPersistDDL = []string{
	`CREATE TABLE records (
  id         TEXT PRIMARY KEY,
  name       TEXT NOT NULL CHECK (length(trim(name)) > 0),
  status     TEXT NOT NULL CHECK (length(trim(status)) > 0),
  owner      TEXT NOT NULL CHECK (length(trim(owner)) > 0),
  updated_at INTEGER NOT NULL
)`,
	`CREATE INDEX idx_records_name ON records(name)`,
	`CREATE INDEX idx_records_updated_at ON records(updated_at)`,
	`CREATE INDEX idx_records_owner ON records(owner)`,
}

// migrate0003 creates the records table and its indexes. Business seed rows are
// intentionally NOT inserted here: seeds live in seedRecords so user deletes and
// re-starts never resurrect or overwrite rows (D-003 §3).
func migrate0003(tx *sql.Tx) error {
	for _, stmt := range recordsPersistDDL {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("create records: %w", err)
		}
	}
	return nil
}

// operationLogDDL is the canonical R5 S6 append-only operation log schema
// (I-008-003 §3). The event CHECK enumerates the frozen event set; created_at is
// Unix milliseconds matching the records updated_at precision (D-004).
var operationLogDDL = []string{
	`CREATE TABLE operation_log (
  id         TEXT PRIMARY KEY,
  event      TEXT NOT NULL CHECK (event IN ('records.create','records.update','records.delete','auth.login','auth.logout','auth.refresh')),
  actor_id   TEXT NOT NULL,
  actor_name TEXT NOT NULL,
  record_id  TEXT,
  detail     TEXT,
  created_at INTEGER NOT NULL
)`,
	`CREATE INDEX idx_operation_log_created_at ON operation_log(created_at DESC)`,
}

// migrate0004 creates the operation log table and its index (R5 S6 optional
// bonus checkpoint, I-008-003).
func migrate0004(tx *sql.Tx) error {
	for _, stmt := range operationLogDDL {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("create operation_log: %w", err)
		}
	}
	return nil
}

// operationLogExpandDDL is the canonical 0005 operation_log definition
// (GOAL-011 S2 · I-011-001 §5): the event CHECK is expanded with the users/roles
// write events while records.* and auth.* stay valid for historical rows. SQLite
// cannot ALTER a CHECK constraint, so migrate0005 rebuilds the table in place.
var operationLogExpandDDL = []string{
	`CREATE TABLE operation_log (
  id         TEXT PRIMARY KEY,
  event      TEXT NOT NULL CHECK (event IN ('records.create','records.update','records.delete','auth.login','auth.logout','auth.refresh','users.create','users.update','users.delete','roles.create','roles.update','roles.delete')),
  actor_id   TEXT NOT NULL,
  actor_name TEXT NOT NULL,
  record_id  TEXT,
  detail     TEXT,
  created_at INTEGER NOT NULL
)`,
	`CREATE INDEX idx_operation_log_created_at ON operation_log(created_at DESC)`,
}

// migrate0005 rebuilds operation_log with the expanded event CHECK, preserving
// existing rows and the created_at index. All steps run in one transaction so a
// failure rolls the rebuild back (fail closed).
func migrate0005(tx *sql.Tx) error {
	if _, err := tx.Exec(`ALTER TABLE operation_log RENAME TO operation_log_old`); err != nil {
		return fmt.Errorf("rename operation_log: %w", err)
	}
	if _, err := tx.Exec(operationLogExpandDDL[0]); err != nil {
		return fmt.Errorf("create operation_log expanded: %w", err)
	}
	if _, err := tx.Exec(
		`INSERT INTO operation_log (id, event, actor_id, actor_name, record_id, detail, created_at)
		 SELECT id, event, actor_id, actor_name, record_id, detail, created_at FROM operation_log_old`,
	); err != nil {
		return fmt.Errorf("migrate operation_log rows: %w", err)
	}
	if _, err := tx.Exec(`DROP TABLE operation_log_old`); err != nil {
		return fmt.Errorf("drop operation_log_old: %w", err)
	}
	if _, err := tx.Exec(operationLogExpandDDL[1]); err != nil {
		return fmt.Errorf("create operation_log index: %w", err)
	}
	return nil
}

// recordsRetireDDL is the canonical 0006 records retirement (GOAL-011 S3 ·
// I-011-002 §2.1): drop the records table and prune the records permission/menu
// rows and their grants. The join rows are deleted before their parents (FK
// RESTRICT). records.* operation-log events stay valid for historical rows.
var recordsRetireDDL = []string{
	`DROP TABLE IF EXISTS records`,
	`DELETE FROM role_permissions WHERE permission_id IN ('perm-records-read','perm-records-write')`,
	`DELETE FROM role_menu_items WHERE menu_item_id = 'menu-list-edit-lifecycle'`,
	`DELETE FROM permissions WHERE id IN ('perm-records-read','perm-records-write')`,
	`DELETE FROM menu_items WHERE id = 'menu-list-edit-lifecycle'`,
}

// migrate0006 retires records from the product surface: the table is dropped and
// the records permission/menu rows cleaned (idempotent deletes). The pre-v0006
// snapshot taken by the runner before this migration is the data-recovery
// backstop (I-011-002 §2.4).
func migrate0006(tx *sql.Tx) error {
	for _, stmt := range recordsRetireDDL {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("records retire: %w", err)
		}
	}
	return nil
}

// migrate applies pending migrations. It is called by Open and must remain the
// only entry point into schema management so the ledger stays authoritative.
func (s *Store) migrate() error {
	if err := validateCompiled(); err != nil {
		return err
	}
	if err := s.assertForeignKeysOn(); err != nil {
		return err
	}

	applied, err := s.appliedMigrations()
	if err != nil {
		return err
	}
	if applied == nil {
		// No ledger yet: 0001 bootstraps the R2 baseline and the ledger in one
		// transaction. An existing R2 database is fingerprint-checked first so
		// a partial or unknown structure is never silently taken over.
		if err := s.applyMigration(compiledMigrations[0]); err != nil {
			return err
		}
		applied, err = s.appliedMigrations()
		if err != nil {
			return err
		}
	}
	if err := validateApplied(applied); err != nil {
		return err
	}

	pending := pendingMigrations(applied)
	// Snapshot before EVERY pending data-mutating migration (version >= 2) so
	// each step has a recoverable pre-state (I-011-002 §2.3 v0.2.0, A-002 F-002).
	// With 0005+0006 both pending, both pre-v0005 and pre-v0006 exist — pre-v0006
	// (post-0005, records table still present) is the records data-recovery
	// backstop. A fresh DB (no rows) produces no snapshot.
	for _, m := range pending {
		if m.version >= 2 {
			if err := s.snapshotBeforePending(m.version); err != nil {
				return err
			}
		}
		if err := s.applyMigration(m); err != nil {
			return err
		}
	}
	return s.verifyIntegrity()
}

// applyMigration runs one migration and its ledger insert in a single
// transaction; any failure rolls the whole step back (fail closed).
func (s *Store) applyMigration(m migration) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("begin migration %d (%s): %w", m.version, m.name, err)
	}
	if err := m.up(tx); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("migration %d (%s): %w", m.version, m.name, err)
	}
	if _, err := tx.Exec(
		`INSERT INTO schema_migrations (version, name, checksum, applied_at) VALUES (?, ?, ?, ?)`,
		m.version, m.name, migrationChecksum(m), time.Now().UTC().Unix()); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("record migration %d (%s): %w", m.version, m.name, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migration %d (%s): %w", m.version, m.name, err)
	}
	return nil
}

// validateCompiled checks the compiled migration list invariants.
func validateCompiled() error {
	if compiledMigrations[0].version != 1 {
		return fmt.Errorf("migrations must start at version 1, have %d", compiledMigrations[0].version)
	}
	names := map[string]int{}
	prev := 0
	for _, m := range compiledMigrations {
		if m.version <= prev {
			return fmt.Errorf("migrations must be strictly ascending, saw %d after %d", m.version, prev)
		}
		if p, ok := names[m.name]; ok {
			return fmt.Errorf("migration name %q reused at versions %d and %d", m.name, p, m.version)
		}
		if len(migrationChecksum(m)) != 64 {
			return fmt.Errorf("migration %d checksum length %d, want 64", m.version, len(migrationChecksum(m)))
		}
		names[m.name] = m.version
		prev = m.version
	}
	return nil
}

// validateApplied verifies the ledger is a contiguous prefix of the compiled
// set with matching names and checksums; any drift fails closed.
func validateApplied(applied []appliedMigration) error {
	known := make(map[int]migration, len(compiledMigrations))
	for _, m := range compiledMigrations {
		known[m.version] = m
	}
	if len(applied) == 0 {
		return errors.New("store: schema_migrations exists but is empty (partial bootstrap)")
	}
	if applied[0].version != 1 {
		return fmt.Errorf("store: migration ledger starts at version %d, want 1", applied[0].version)
	}
	for i, a := range applied {
		if i > 0 && a.version != applied[i-1].version+1 {
			return fmt.Errorf("store: migration ledger missing intermediate version before %d", a.version)
		}
		m, ok := known[a.version]
		if !ok {
			return fmt.Errorf("store: unknown applied migration version %d (%s)", a.version, a.name)
		}
		if m.name != a.name {
			return fmt.Errorf("store: applied migration %d name %q, code has %q", a.version, a.name, m.name)
		}
		if got := migrationChecksum(m); got != a.checksum {
			return fmt.Errorf("store: migration %d checksum drift (ledger %s, code %s)", a.version, a.checksum, got)
		}
	}
	return nil
}

func pendingMigrations(applied []appliedMigration) []migration {
	appliedSet := make(map[int]bool, len(applied))
	for _, a := range applied {
		appliedSet[a.version] = true
	}
	var pending []migration
	for _, m := range compiledMigrations {
		if !appliedSet[m.version] {
			pending = append(pending, m)
		}
	}
	return pending
}

// appliedMigrations returns nil when no ledger table exists, else the applied
// rows ordered by version. An existing but empty ledger is a partial bootstrap
// and is rejected.
func (s *Store) appliedMigrations() ([]appliedMigration, error) {
	var exists bool
	if err := s.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM sqlite_master WHERE type = 'table' AND name = 'schema_migrations')`,
	).Scan(&exists); err != nil {
		return nil, fmt.Errorf("store: check migration ledger: %w", err)
	}
	if !exists {
		return nil, nil
	}
	rows, err := s.db.Query(`SELECT version, name, checksum FROM schema_migrations ORDER BY version`)
	if err != nil {
		return nil, fmt.Errorf("store: read migration ledger: %w", err)
	}
	defer rows.Close()
	var applied []appliedMigration
	for rows.Next() {
		var a appliedMigration
		if err := rows.Scan(&a.version, &a.name, &a.checksum); err != nil {
			return nil, fmt.Errorf("store: scan migration ledger: %w", err)
		}
		applied = append(applied, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: read migration ledger: %w", err)
	}
	if len(applied) == 0 {
		return nil, errors.New("store: schema_migrations exists but is empty (partial bootstrap)")
	}
	return applied, nil
}

// 0001 r2_baseline: on an empty database, create the R2 tables and ledger; on an
// existing R2 database (no ledger), fingerprint-check the structure first. The
// ledger table itself is only created at the end of the transaction so a failed
// fingerprint leaves nothing behind.
func migrate0001(tx *sql.Tx) error {
	empty, err := isEmptyDatabase(tx)
	if err != nil {
		return err
	}
	if empty {
		for _, stmt := range r2BaselineDDL {
			if _, err := tx.Exec(stmt); err != nil {
				return fmt.Errorf("create baseline: %w", err)
			}
		}
	} else {
		if err := fingerprintR2(tx); err != nil {
			return err
		}
	}
	if _, err := tx.Exec(schemaMigrationsDDL); err != nil {
		return fmt.Errorf("create migration ledger: %w", err)
	}
	return nil
}

func isEmptyDatabase(tx *sql.Tx) (bool, error) {
	var count int
	err := tx.QueryRow(
		`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%'`,
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("store: count tables: %w", err)
	}
	return count == 0, nil
}

// fingerprintR2 verifies an existing, ledger-less database matches the R2
// minimal structure (users + refresh_tokens, their columns, the FK and the
// user_id index). Any mismatch fails closed and rolls back.
func fingerprintR2(tx *sql.Tx) error {
	got := map[string]bool{}
	rows, err := tx.Query(`SELECT name FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%'`)
	if err != nil {
		return fmt.Errorf("fingerprint: list tables: %w", err)
	}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			rows.Close()
			return err
		}
		got[name] = true
	}
	if err := rows.Close(); err != nil {
		return err
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if len(got) != 2 || !got["users"] || !got["refresh_tokens"] {
		return fmt.Errorf("fingerprint: unexpected table set %v, want {users refresh_tokens}", tableSet(got))
	}

	if err := fingerprintColumns(tx, "users", map[string]string{
		"id": "TEXT", "username": "TEXT", "name": "TEXT", "roles": "TEXT",
		"password_hash": "TEXT", "created_at": "INTEGER", "updated_at": "INTEGER",
	}); err != nil {
		return err
	}
	if err := fingerprintColumns(tx, "refresh_tokens", map[string]string{
		"id": "TEXT", "user_id": "TEXT", "token_hash": "TEXT", "expires_at": "INTEGER",
		"revoked_at": "INTEGER", "created_at": "INTEGER",
	}); err != nil {
		return err
	}

	// refresh_tokens.user_id -> users.id
	var fks int
	fkRows, err := tx.Query(`PRAGMA foreign_key_list(refresh_tokens)`)
	if err != nil {
		return fmt.Errorf("fingerprint: fk list: %w", err)
	}
	for fkRows.Next() {
		var id, seq, onUpdate, onDelete, match string
		var table, from, to *string
		if err := fkRows.Scan(&id, &seq, &table, &from, &to, &onUpdate, &onDelete, &match); err != nil {
			fkRows.Close()
			return fmt.Errorf("fingerprint: scan fk: %w", err)
		}
		if table != nil && *table == "users" && from != nil && *from == "user_id" && to != nil && *to == "id" {
			fks++
		}
	}
	if err := fkRows.Close(); err != nil {
		return err
	}
	if err := fkRows.Err(); err != nil {
		return err
	}
	if fks == 0 {
		return errors.New("fingerprint: refresh_tokens.user_id FK -> users.id missing")
	}

	// idx_refresh_tokens_user_id on refresh_tokens(user_id)
	var idxCount int
	if err := tx.QueryRow(
		`SELECT COUNT(*) FROM sqlite_master WHERE type = 'index' AND name = 'idx_refresh_tokens_user_id' AND tbl_name = 'refresh_tokens'`,
	).Scan(&idxCount); err != nil {
		return fmt.Errorf("fingerprint: check index: %w", err)
	}
	if idxCount != 1 {
		return errors.New("fingerprint: index idx_refresh_tokens_user_id missing")
	}
	return nil
}

func tableSet(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func fingerprintColumns(tx *sql.Tx, table string, want map[string]string) error {
	rows, err := tx.Query(fmt.Sprintf(`PRAGMA table_info(%s)`, table))
	if err != nil {
		return fmt.Errorf("fingerprint %s: %w", table, err)
	}
	defer rows.Close()
	got := map[string]string{}
	for rows.Next() {
		var cid, notnull, pk int
		var name, typ string
		var dflt sql.NullString
		if err := rows.Scan(&cid, &name, &typ, &notnull, &dflt, &pk); err != nil {
			return fmt.Errorf("fingerprint %s: scan: %w", table, err)
		}
		got[name] = strings.ToUpper(typ)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("fingerprint %s: %w", table, err)
	}
	if len(got) != len(want) {
		return fmt.Errorf("fingerprint: %s has %d columns, want %d", table, len(got), len(want))
	}
	for name, typ := range want {
		if t, ok := got[name]; !ok || t != typ {
			return fmt.Errorf("fingerprint: %s.%s = %q, want %q", table, name, got[name], typ)
		}
	}
	return nil
}

// 0002 rbac_expand: create the normalized RBAC tables and backfill roles from
// each user's legacy roles JSON, validating keys and deduping within a user.
// Any malformed roles array, empty/invalid key or constraint conflict rolls the
// whole migration back (no partial backfill).
func migrate0002(tx *sql.Tx) error {
	for _, stmt := range rbacExpandDDL {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("create rbac: %w", err)
		}
	}
	return backfillRoles(tx)
}

// roleKeyRe matches stable role keys: lowercase start, then lowercase letters,
// digits, underscore or hyphen.
var roleKeyRe = regexp.MustCompile(`^[a-z][a-z0-9_-]*$`)

// ensureRole idempotently creates the derived role row for a stable role key
// (system=0 until the S3 seed upgrades known roles). It is shared by 0002
// backfill, CreateUser and seedAdmin so the user_roles FK never dangles.
func ensureRole(tx *sql.Tx, key string, now int64) error {
	if _, err := tx.Exec(
		`INSERT INTO roles (id, key, name, system, created_at, updated_at)
		 VALUES (?, ?, ?, 0, ?, ?)
		 ON CONFLICT(id) DO NOTHING`,
		"role-"+key, key, key, now, now,
	); err != nil {
		return fmt.Errorf("ensure role %s: %w", key, err)
	}
	return nil
}

// linkUserRole idempotently links a user to a role, creating the role row first.
// It is used by 0002 backfill, CreateUser and seedAdmin (阶段 A/B double-write).
func linkUserRole(tx *sql.Tx, userID, key string, now int64) error {
	if !roleKeyRe.MatchString(key) {
		return fmt.Errorf("invalid role key %q", key)
	}
	if err := ensureRole(tx, key, now); err != nil {
		return err
	}
	if _, err := tx.Exec(
		`INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)
		 ON CONFLICT(user_id, role_id) DO NOTHING`,
		userID, "role-"+key,
	); err != nil {
		return fmt.Errorf("link user %s role %s: %w", userID, key, err)
	}
	return nil
}

func backfillRoles(tx *sql.Tx) error {
	rows, err := tx.Query(`SELECT id, roles FROM users`)
	if err != nil {
		return fmt.Errorf("backfill: list users: %w", err)
	}
	defer rows.Close()
	now := time.Now().UTC().Unix()
	for rows.Next() {
		var id, rolesJSON string
		if err := rows.Scan(&id, &rolesJSON); err != nil {
			return err
		}
		var keys []string
		if err := json.Unmarshal([]byte(rolesJSON), &keys); err != nil {
			return fmt.Errorf("backfill user %s: roles %q is not a JSON array: %w", id, rolesJSON, err)
		}
		seen := map[string]bool{}
		for _, key := range keys {
			if !roleKeyRe.MatchString(key) {
				return fmt.Errorf("backfill user %s: invalid role key %q", id, key)
			}
			if seen[key] {
				continue // dedupe within a user's roles list
			}
			seen[key] = true
			if err := linkUserRole(tx, id, key, now); err != nil {
				return fmt.Errorf("backfill user %s: %w", id, err)
			}
		}
	}
	return rows.Err()
}

// assertForeignKeysOn enables and verifies PRAGMA foreign_keys on the single
// connection the store uses (SetMaxOpenConns(1)). This is a hard invariant for
// every migration and query.
func (s *Store) assertForeignKeysOn() error {
	if _, err := s.db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		return fmt.Errorf("store: enable foreign_keys: %w", err)
	}
	var on int
	if err := s.db.QueryRow(`PRAGMA foreign_keys`).Scan(&on); err != nil {
		return fmt.Errorf("store: read foreign_keys: %w", err)
	}
	if on != 1 {
		return errors.New("store: PRAGMA foreign_keys could not be enabled")
	}
	return nil
}

// snapshotBeforePending produces a recoverable copy of a non-empty file database
// right before the first pending migration (version >= 2) is applied, using
// SQLite's VACUUM INTO consistency snapshot. The path is passed as a bound SQL
// string literal (quotes escaped), never through a shell.
func (s *Store) snapshotBeforePending(firstPendingVersion int) error {
	if s.path == "" || s.path == ":memory:" {
		return nil
	}
	hasData, err := s.dbHasRows()
	if err != nil {
		return err
	}
	if !hasData {
		return nil // nothing to recover yet
	}
	target := fmt.Sprintf("%s.pre-v%04d-%s.sqlite", s.path, firstPendingVersion, time.Now().UTC().Format("20060102T150405Z"))
	if _, err := s.db.Exec("VACUUM INTO '" + strings.ReplaceAll(target, "'", "''") + "'"); err != nil {
		return fmt.Errorf("pre-v%04d snapshot to %s: %w", firstPendingVersion, target, err)
	}
	if err := checkIntegrityFile(target); err != nil {
		return fmt.Errorf("pre-v%04d snapshot %s invalid: %w", firstPendingVersion, target, err)
	}
	return nil
}

// dbHasRows reports whether any application table has at least one row. It is
// used to decide whether a pre-upgrade snapshot is worth taking: a database with
// no data has nothing to recover, and a fresh DB must not leave stray snapshots.
// Table names are collected (and the first query's rows closed) before counting,
// because the store holds a single connection and a nested query would deadlock.
func (s *Store) dbHasRows() (bool, error) {
	rows, err := s.db.Query(`SELECT name FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%' AND name != 'schema_migrations'`)
	if err != nil {
		return false, fmt.Errorf("snapshot: list tables: %w", err)
	}
	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			rows.Close()
			return false, fmt.Errorf("snapshot: scan table: %w", err)
		}
		tables = append(tables, name)
	}
	if err := rows.Close(); err != nil {
		return false, fmt.Errorf("snapshot: list tables: %w", err)
	}
	if err := rows.Err(); err != nil {
		return false, fmt.Errorf("snapshot: list tables: %w", err)
	}
	for _, name := range tables {
		var n int
		if err := s.db.QueryRow(`SELECT COUNT(*) FROM "` + strings.ReplaceAll(name, `"`, `""`) + `"`).Scan(&n); err != nil {
			return false, fmt.Errorf("snapshot: count %s: %w", name, err)
		}
		if n > 0 {
			return true, nil
		}
	}
	return false, nil
}

// checkIntegrityFile opens a database file and verifies its integrity.
func checkIntegrityFile(path string) error {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return fmt.Errorf("open snapshot: %w", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	return checkIntegrity(db)
}

// verifyIntegrity runs integrity_check (single "ok" row) and foreign_key_check
// (must return no violation rows) on the main database after migration.
func (s *Store) verifyIntegrity() error {
	if err := checkIntegrity(s.db); err != nil {
		return err
	}
	rows, err := s.db.Query(`PRAGMA foreign_key_check`)
	if err != nil {
		return fmt.Errorf("store: foreign_key_check: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		var table, parent string
		var rowid int64
		if err := rows.Scan(&table, &rowid, &parent); err != nil {
			return fmt.Errorf("store: read foreign_key_check: %w", err)
		}
		return fmt.Errorf("store: foreign_key_check violation: %s rowid %d references %s", table, rowid, parent)
	}
	return rows.Err()
}

func checkIntegrity(db *sql.DB) error {
	var result string
	if err := db.QueryRow(`PRAGMA integrity_check`).Scan(&result); err != nil {
		return fmt.Errorf("integrity_check: %w", err)
	}
	if result != "ok" {
		return fmt.Errorf("integrity_check = %q, want ok", result)
	}
	return nil
}

// migrationChecksum is the SHA-256 (lower hex) of the normalized canonical SQL
// joined with the explicit data-transformer id.
func migrationChecksum(m migration) string {
	input := normalizeSQL(strings.Join(m.stmts, "\n")) + "\n" + m.transformID
	sum := sha256.Sum256([]byte(input))
	return hex.EncodeToString(sum[:])
}

func normalizeSQL(s string) string {
	lines := strings.Split(s, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, line)
		}
	}
	return strings.Join(out, "\n")
}
