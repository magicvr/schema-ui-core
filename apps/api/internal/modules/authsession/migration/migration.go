package migration

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

const ModuleID = "core.auth-session"

const schemaMigrationsDDL = `CREATE TABLE schema_migrations (
  version    INTEGER PRIMARY KEY,
  name       TEXT NOT NULL UNIQUE,
  checksum   TEXT NOT NULL CHECK (length(checksum) = 64),
  applied_at INTEGER NOT NULL
)`

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

var systemDataReconcileDDL = []string{
	`CREATE TABLE system_data_reconcile (
  module_id        TEXT NOT NULL,
  kind             TEXT NOT NULL CHECK (kind IN ('base','authorization','navigation')),
  contribution_key TEXT NOT NULL,
  version          INTEGER NOT NULL CHECK (version > 0),
  checksum         TEXT NOT NULL CHECK (length(checksum) = 64),
  applied_at       INTEGER NOT NULL,
  PRIMARY KEY (module_id, kind, contribution_key)
)`,
	`CREATE TABLE system_data_grants (
  module_id        TEXT NOT NULL,
  kind             TEXT NOT NULL CHECK (kind IN ('authorization','navigation')),
  contribution_key TEXT NOT NULL,
  role_key         TEXT NOT NULL,
  target_id        TEXT NOT NULL,
  PRIMARY KEY (module_id, kind, contribution_key, role_key, target_id)
)`,
}

// Descriptors returns the immutable 0001-0002 auth/session migration history.
func Descriptors() []kernel.MigrationContribution {
	return []kernel.MigrationContribution{
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "r2_baseline"},
			Version:              1,
			Name:                 "r2_baseline",
			Checksum:             kernel.MigrationChecksum(r2BaselineDDL, "0001:r2-baseline:v1"),
			Apply:                migrateBaseline,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "rbac_expand"},
			Version:              2,
			Name:                 "rbac_expand",
			Checksum:             kernel.MigrationChecksum(rbacExpandDDL, "0002:rbac-expand:v1"),
			Apply:                migrateRBAC,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "system_data_reconcile"},
			Version:              9,
			Name:                 "system_data_reconcile",
			Checksum:             kernel.MigrationChecksum(systemDataReconcileDDL, "0009:system-data-reconcile:v1"),
			Apply:                migrateSystemDataReconcile,
		},
	}
}

func migrateBaseline(tx *sql.Tx) error {
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
	} else if err := fingerprintR2(tx); err != nil {
		return err
	}
	if _, err := tx.Exec(schemaMigrationsDDL); err != nil {
		return fmt.Errorf("create migration ledger: %w", err)
	}
	return nil
}

func migrateRBAC(tx *sql.Tx) error {
	for _, stmt := range rbacExpandDDL {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("create rbac: %w", err)
		}
	}
	return backfillRoles(tx)
}

func migrateSystemDataReconcile(tx *sql.Tx) error {
	for _, stmt := range systemDataReconcileDDL {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("create system-data reconcile tables: %w", err)
		}
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

	var foreignKeys int
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
			foreignKeys++
		}
	}
	if err := fkRows.Close(); err != nil {
		return err
	}
	if err := fkRows.Err(); err != nil {
		return err
	}
	if foreignKeys == 0 {
		return errors.New("fingerprint: refresh_tokens.user_id FK -> users.id missing")
	}

	var indexCount int
	if err := tx.QueryRow(
		`SELECT COUNT(*) FROM sqlite_master WHERE type = 'index' AND name = 'idx_refresh_tokens_user_id' AND tbl_name = 'refresh_tokens'`,
	).Scan(&indexCount); err != nil {
		return fmt.Errorf("fingerprint: check index: %w", err)
	}
	if indexCount != 1 {
		return errors.New("fingerprint: index idx_refresh_tokens_user_id missing")
	}
	return nil
}

func tableSet(items map[string]bool) []string {
	out := make([]string, 0, len(items))
	for item := range items {
		out = append(out, item)
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
		var cid, notNull, primaryKey int
		var name, typ string
		var defaultValue sql.NullString
		if err := rows.Scan(&cid, &name, &typ, &notNull, &defaultValue, &primaryKey); err != nil {
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
		if actual, ok := got[name]; !ok || actual != typ {
			return fmt.Errorf("fingerprint: %s.%s = %q, want %q", table, name, got[name], typ)
		}
	}
	return nil
}

var roleKeyPattern = regexp.MustCompile(`^[a-z][a-z0-9_-]*$`)

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

func linkUserRole(tx *sql.Tx, userID, key string, now int64) error {
	if !roleKeyPattern.MatchString(key) {
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
			if !roleKeyPattern.MatchString(key) {
				return fmt.Errorf("backfill user %s: invalid role key %q", id, key)
			}
			if seen[key] {
				continue
			}
			seen[key] = true
			if err := linkUserRole(tx, id, key, now); err != nil {
				return fmt.Errorf("backfill user %s: %w", id, err)
			}
		}
	}
	return rows.Err()
}
