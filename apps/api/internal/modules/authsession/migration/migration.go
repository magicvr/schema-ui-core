package migration

import (
	"context"
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

// accessTokenRevocationDDL (W4 P0-3): a per-user monotonic token_version column
// lets the server revoke already-issued access tokens immediately. Every
// access-token JWT carries the user's token_version at issue time; the auth
// middleware compares it to the persisted value and rejects a stale token after
// a password change (which bumps the version). ALTER TABLE ADD COLUMN is
// additive and idempotent within a single fresh-run migration ledger.
var accessTokenRevocationDDL = []string{
	`ALTER TABLE users ADD COLUMN token_version INTEGER NOT NULL DEFAULT 0`,
}

// accountLockDDL adds the account-lock production source (GOAL-004 S4-6):
// failed_login_count feeds the lock threshold; locked_until (unix seconds,
// 0 = not locked) is the lock window with automatic expiry.
var accountLockDDL = []string{
	`ALTER TABLE users ADD COLUMN failed_login_count INTEGER NOT NULL DEFAULT 0`,
	`ALTER TABLE users ADD COLUMN locked_until INTEGER NOT NULL DEFAULT 0`,
}

// mustChangePasswordDDL adds the forced-password-change flag (W16-F01):
// 1 = the user must replace the initial/reset password before business APIs
// are usable; 0 = normal. Additive and backward compatible.
var mustChangePasswordDDL = []string{
	`ALTER TABLE users ADD COLUMN must_change_password INTEGER NOT NULL DEFAULT 0`,
}

// seedAdminMustChangePasswordSQL (0049 · A2 backfill): sets
// must_change_password = 1 on the seed admin row (id = 'user-admin') when it
// was created before migration 0038 introduced the column.  The WHERE clause
// is intentionally narrow: it targets only the stable seed-admin primary key
// ('user-admin') AND requires must_change_password = 0, so the UPDATE is a
// no-op on fresh databases (bootstrap already inserts with value 1) and on any
// re-run (idempotent).  Human-created accounts with id != 'user-admin' are
// never touched.
var seedAdminMustChangePasswordSQL = []string{
	`UPDATE users SET must_change_password = 1 WHERE id = 'user-admin' AND must_change_password = 0`,
}

// accountEmailIdentityDDL (workspace-018 R2 · GOAL-003 D-001): account email
// identity per the frozen R1 contract (GOAL-002 D-001 §1/§2/§3/§6).
//   - email TEXT NULL: NULL = unbound; NULLs are mutually distinct in unique
//     indexes on both dialects, so accounts without email never collide.
//   - email_status TEXT NULL CHECK ('pending'|'verified'): meaningful only
//     when email IS NOT NULL — a NULL email means unbound regardless.
//   - idx_users_email_lower on lower(email): the physical carrier of
//     bind-reserves-slot case-insensitive uniqueness. All three statements
//     are byte-identical on sqlite and postgres, so ApplyPostgres stays nil;
//     sqlite lower() folds ASCII only, compensated by application-layer
//     normalization in R3 repositories (GOAL-002 A-001 F-2).
var accountEmailIdentityDDL = []string{
	`ALTER TABLE users ADD COLUMN email TEXT`,
	`ALTER TABLE users ADD COLUMN email_status TEXT CHECK (email_status IN ('pending','verified'))`,
	`CREATE UNIQUE INDEX idx_users_email_lower ON users(lower(email))`,
}

// emailVerificationDDL (workspace-018 R3 · GOAL-004 D-001 §1): one active
// verification challenge per user — PK on user_id makes the upsert an
// idempotent replace, ON DELETE CASCADE cleans up with the account. Times are
// unix seconds; the INTEGER/BIGINT split follows the accountLock precedent,
// so this migration ships paired dialect bodies (not portable).
var emailVerificationDDL = []string{
	`CREATE TABLE email_verification_challenges (
  user_id       TEXT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  code_hash     TEXT NOT NULL,
  expires_at    INTEGER NOT NULL,
  sent_at       INTEGER NOT NULL,
  attempt_count INTEGER NOT NULL DEFAULT 0
)`,
}

var postgresEmailVerificationDDL = []string{
	`CREATE TABLE email_verification_challenges (
  user_id       TEXT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  code_hash     TEXT NOT NULL,
  expires_at    BIGINT NOT NULL,
  sent_at       BIGINT NOT NULL,
  attempt_count INTEGER NOT NULL DEFAULT 0
)`,
}

// passwordRecoveryDDL (workspace-019 R2 · GOAL-003 D-001 §1): one active
// self-recovery challenge per user — same shape as the email verification
// table (0055) so the frozen numbers (TTL 10 min / cooldown 60 s / 5 failed
// attempts void) reuse the identical bookkeeping. PK on user_id makes the
// upsert an idempotent replace, ON DELETE CASCADE cleans up with the account.
// Times are unix seconds; INTEGER/BIGINT split follows the 0055 precedent,
// so this migration ships paired dialect bodies (not portable).
var passwordRecoveryDDL = []string{
	`CREATE TABLE password_recovery_challenges (
  user_id       TEXT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  code_hash     TEXT NOT NULL,
  expires_at    INTEGER NOT NULL,
  sent_at       INTEGER NOT NULL,
  attempt_count INTEGER NOT NULL DEFAULT 0
)`,
}

var passwordRecoveryPGDDL = []string{
	`CREATE TABLE password_recovery_challenges (
  user_id       TEXT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  code_hash     TEXT NOT NULL,
  expires_at    BIGINT NOT NULL,
  sent_at       BIGINT NOT NULL,
  attempt_count INTEGER NOT NULL DEFAULT 0
)`,
}

// ---- postgres-flavored apply bodies (R3 dual-dialect ledger; R1 v1.4 §3/§4).
// The sqlite/canonical SQL above is untouched so its checksums stay stable;
// these bodies run only on the postgres runner. Unix time columns are BIGINT,
// COLLATE NOCASE becomes CITEXT. migrateBaselinePG mirrors sqlite's empty-vs-
// fingerprint split: a ledger-less database that already has the R2 tables is
// adopted (create the ledger only), not CREATE-TABLE'd again.

const postgresSchemaMigrationsDDL = `CREATE TABLE schema_migrations (
  version    INTEGER PRIMARY KEY,
  name       TEXT NOT NULL UNIQUE,
  checksum   TEXT NOT NULL CHECK (length(checksum) = 64),
  applied_at BIGINT NOT NULL
)`

// postgresBaselineDDL mirrors r2BaselineDDL + schemaMigrationsDDL with BIGINT
// time columns; the ledger is created here (the sqlite Apply creates it inside
// migrateBaseline, this postgres body inlines it).
var postgresBaselineDDL = []string{
	postgresSchemaMigrationsDDL,
	`CREATE TABLE users (
  id            TEXT PRIMARY KEY,
  username      TEXT NOT NULL UNIQUE,
  name          TEXT NOT NULL,
  roles         TEXT NOT NULL,
  password_hash TEXT NOT NULL,
  created_at    BIGINT NOT NULL,
  updated_at    BIGINT NOT NULL
)`,
	`CREATE TABLE refresh_tokens (
  id         TEXT PRIMARY KEY,
  user_id    TEXT NOT NULL REFERENCES users(id),
  token_hash TEXT NOT NULL UNIQUE,
  expires_at BIGINT NOT NULL,
  revoked_at BIGINT,
  created_at BIGINT NOT NULL
)`,
	`CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id)`,
}

// postgresRBACDDL mirrors rbacExpandDDL with BIGINT time columns.
var postgresRBACDDL = []string{
	`CREATE TABLE roles (
  id         TEXT PRIMARY KEY,
  key        TEXT NOT NULL UNIQUE CHECK (key <> ''),
  name       TEXT NOT NULL,
  system     INTEGER NOT NULL DEFAULT 0 CHECK (system IN (0, 1)),
  created_at BIGINT NOT NULL,
  updated_at BIGINT NOT NULL
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
  created_at  BIGINT NOT NULL,
  updated_at  BIGINT NOT NULL
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
  created_at  BIGINT NOT NULL,
  updated_at  BIGINT NOT NULL
)`,
	`CREATE TABLE role_menu_items (
  role_id      TEXT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
  menu_item_id TEXT NOT NULL REFERENCES menu_items(id) ON DELETE RESTRICT,
  PRIMARY KEY (role_id, menu_item_id)
)`,
	`CREATE INDEX idx_role_menu_items_menu_item_id ON role_menu_items(menu_item_id)`,
}

// postgresSystemDataDDL mirrors systemDataReconcileDDL with BIGINT applied_at.
var postgresSystemDataDDL = []string{
	`CREATE TABLE system_data_reconcile (
  module_id        TEXT NOT NULL,
  kind             TEXT NOT NULL CHECK (kind IN ('base','authorization','navigation')),
  contribution_key TEXT NOT NULL,
  version          INTEGER NOT NULL CHECK (version > 0),
  checksum         TEXT NOT NULL CHECK (length(checksum) = 64),
  applied_at       BIGINT NOT NULL,
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

// postgresAccountLockDDL mirrors accountLockDDL; locked_until is a unix-seconds
// time column and must be BIGINT on postgres.
var postgresAccountLockDDL = []string{
	`ALTER TABLE users ADD COLUMN failed_login_count INTEGER NOT NULL DEFAULT 0`,
	`ALTER TABLE users ADD COLUMN locked_until BIGINT NOT NULL DEFAULT 0`,
}

// postgresServiceCredentialsDDL mirrors serviceCredentialsDDL with BIGINT time
// columns and the name unique key on CITEXT (case-insensitive like sqlite
// COLLATE NOCASE; R1 v1.4 F-002 — citing the extension is the DDL-level
// equivalent).
var postgresServiceCredentialsDDL = []string{
	`CREATE EXTENSION IF NOT EXISTS citext`,
	`CREATE TABLE service_credentials (
  id           TEXT PRIMARY KEY CHECK (length(id) = 32),
  name         CITEXT NOT NULL UNIQUE CHECK (length(trim(name)) BETWEEN 1 AND 100),
  token_prefix TEXT NOT NULL CHECK (length(token_prefix) = 15),
  token_hash   TEXT NOT NULL UNIQUE CHECK (length(token_hash) = 64),
  scopes       TEXT NOT NULL,
  expires_at   BIGINT NOT NULL,
  revoked_at   BIGINT,
  last_used_at BIGINT,
  created_by   TEXT NOT NULL,
  created_at   BIGINT NOT NULL,
  updated_at   BIGINT NOT NULL
)`,
	`CREATE INDEX idx_service_credentials_created_at ON service_credentials(created_at DESC, id DESC)`,
	`CREATE INDEX idx_service_credentials_expires_at ON service_credentials(expires_at)`,
}

var serviceCredentialsDDL = []string{
	`CREATE TABLE service_credentials (
  id           TEXT PRIMARY KEY CHECK (length(id) = 32),
  name         TEXT NOT NULL COLLATE NOCASE UNIQUE CHECK (length(trim(name)) BETWEEN 1 AND 100),
  token_prefix TEXT NOT NULL CHECK (length(token_prefix) = 15),
  token_hash   TEXT NOT NULL UNIQUE CHECK (length(token_hash) = 64),
  scopes       TEXT NOT NULL,
  expires_at   INTEGER NOT NULL,
  revoked_at   INTEGER,
  last_used_at INTEGER,
  created_by   TEXT NOT NULL,
  created_at   INTEGER NOT NULL,
  updated_at   INTEGER NOT NULL
)`,
	`CREATE INDEX idx_service_credentials_created_at ON service_credentials(created_at DESC, id DESC)`,
	`CREATE INDEX idx_service_credentials_expires_at ON service_credentials(expires_at)`,
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
			ApplyPostgres:        migrateBaselinePG,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "rbac_expand"},
			Version:              2,
			Name:                 "rbac_expand",
			Checksum:             kernel.MigrationChecksum(rbacExpandDDL, "0002:rbac-expand:v1"),
			Apply:                migrateRBAC,
			ApplyPostgres:        migrateRBACPG,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "system_data_reconcile"},
			Version:              9,
			Name:                 "system_data_reconcile",
			Checksum:             kernel.MigrationChecksum(systemDataReconcileDDL, "0009:system-data-reconcile:v1"),
			Apply:                migrateSystemDataReconcile,
			ApplyPostgres:        migrateSystemDataPG,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "access_token_revocation"},
			Version:              11,
			Name:                 "access_token_revocation",
			Checksum:             kernel.MigrationChecksum(accessTokenRevocationDDL, "0011:access-token-revocation:v1"),
			Apply:                migrateAccessTokenRevocation,
			// ApplyPostgres nil: portable additive ALTER (INTEGER token_version).
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "account_lock"},
			Version:              12,
			Name:                 "account_lock",
			Checksum:             kernel.MigrationChecksum(accountLockDDL, "0012:account-lock:v1"),
			Apply:                migrateAccountLock,
			ApplyPostgres:        migrateAccountLockPG,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "must_change_password"},
			Version:              38,
			Name:                 "must_change_password",
			Checksum:             kernel.MigrationChecksum(mustChangePasswordDDL, "0038:must-change-password:v1"),
			Apply:                migrateMustChangePassword,
			// ApplyPostgres nil: portable additive ALTER (INTEGER flag).
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "service_credentials"},
			Version:              44,
			Name:                 "service_credentials",
			Checksum:             kernel.MigrationChecksum(serviceCredentialsDDL, "0044:service-credentials:v1"),
			Apply:                migrateServiceCredentials,
			ApplyPostgres:        migrateServiceCredentialsPG,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "seed_admin_must_change_password"},
			Version:              49,
			Name:                 "seed_admin_must_change_password",
			Checksum:             kernel.MigrationChecksum(seedAdminMustChangePasswordSQL, "0049:seed-admin-must-change-password:v1"),
			Apply:                migrateSeedAdminMustChangePassword,
			// ApplyPostgres nil: portable UPDATE (no dialect difference).
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "account_email_identity"},
			Version:              54,
			Name:                 "account_email_identity",
			Checksum:             kernel.MigrationChecksum(accountEmailIdentityDDL, "0054:account-email-identity:v1"),
			Apply:                migrateAccountEmailIdentity,
			// ApplyPostgres nil: portable DDL (TEXT column + literal CHECK +
			// lower() expression unique index are identical on both dialects).
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "email_verification_challenges"},
			Version:              55,
			Name:                 "email_verification_challenges",
			Checksum:             kernel.MigrationChecksum(emailVerificationDDL, "0055:email-verification-challenges:v1"),
			Apply:                migrateEmailVerification,
			ApplyPostgres:        migrateEmailVerificationPG,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "password_recovery_challenges"},
			Version:              56,
			Name:                 "password_recovery_challenges",
			Checksum:             kernel.MigrationChecksum(passwordRecoveryDDL, "0056:password-recovery-challenges:v1"),
			Apply:                migratePasswordRecovery,
			ApplyPostgres:        migratePasswordRecoveryPG,
		},
	}
}

func migrateBaseline(tx kernel.Tx) error {
	empty, err := isEmptyDatabase(tx)
	if err != nil {
		return err
	}
	if empty {
		for _, stmt := range r2BaselineDDL {
			if _, err := tx.Exec(context.Background(), stmt); err != nil {
				return fmt.Errorf("create baseline: %w", err)
			}
		}
	} else if err := fingerprintR2(tx); err != nil {
		return err
	}
	if _, err := tx.Exec(context.Background(), schemaMigrationsDDL); err != nil {
		return fmt.Errorf("create migration ledger: %w", err)
	}
	return nil
}

func migrateRBAC(tx kernel.Tx) error {
	for _, stmt := range rbacExpandDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create rbac: %w", err)
		}
	}
	return backfillRoles(tx)
}

func migrateSystemDataReconcile(tx kernel.Tx) error {
	for _, stmt := range systemDataReconcileDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create system-data reconcile tables: %w", err)
		}
	}
	return nil
}

func migrateAccessTokenRevocation(tx kernel.Tx) error {
	for _, stmt := range accessTokenRevocationDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("add users.token_version: %w", err)
		}
	}
	return nil
}

func migrateAccountLock(tx kernel.Tx) error {
	for _, stmt := range accountLockDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("add account-lock columns: %w", err)
		}
	}
	return nil
}

func migrateMustChangePassword(tx kernel.Tx) error {
	for _, stmt := range mustChangePasswordDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("add users.must_change_password: %w", err)
		}
	}
	return nil
}

func migrateSeedAdminMustChangePassword(tx kernel.Tx) error {
	for _, stmt := range seedAdminMustChangePasswordSQL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("backfill seed admin must_change_password: %w", err)
		}
	}
	return nil
}

func migrateAccountEmailIdentity(tx kernel.Tx) error {
	for _, stmt := range accountEmailIdentityDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("add account email identity columns/index: %w", err)
		}
	}
	return nil
}

func migrateEmailVerification(tx kernel.Tx) error {
	for _, stmt := range emailVerificationDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create email verification challenges: %w", err)
		}
	}
	return nil
}

func migrateEmailVerificationPG(tx kernel.Tx) error {
	for _, stmt := range postgresEmailVerificationDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create email verification challenges (postgres): %w", err)
		}
	}
	return nil
}

func migratePasswordRecovery(tx kernel.Tx) error {
	for _, stmt := range passwordRecoveryDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create password recovery challenges: %w", err)
		}
	}
	return nil
}

func migratePasswordRecoveryPG(tx kernel.Tx) error {
	for _, stmt := range passwordRecoveryPGDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create password recovery challenges (postgres): %w", err)
		}
	}
	return nil
}

func migrateServiceCredentials(tx kernel.Tx) error {
	for _, stmt := range serviceCredentialsDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create service credentials: %w", err)
		}
	}
	return nil
}

// migrateBaselinePG is the postgres variant of migrateBaseline. The runner
// only calls this when schema_migrations is absent.
//
//   - empty schema: create the full baseline (ledger + users + refresh_tokens)
//   - existing schema-ui `users` (later columns allowed): create the ledger and
//     skip objects that already exist (42P07), so a second start does not fail
//   - a `users` table that is not ours: fail closed (do not CREATE around it)
func migrateBaselinePG(tx kernel.Tx) error {
	empty, err := isEmptyDatabasePG(tx)
	if err != nil {
		return err
	}
	if empty {
		for _, stmt := range postgresBaselineDDL {
			if _, err := tx.Exec(context.Background(), stmt); err != nil {
				return fmt.Errorf("create baseline (postgres): %w", err)
			}
		}
		return nil
	}
	ours, err := usersLooksLikeSchemaUIPG(tx)
	if err != nil {
		return err
	}
	if !ours {
		hasUsers, err := tableExistsPG(tx, "users")
		if err != nil {
			return err
		}
		if hasUsers {
			return errors.New("fingerprint (postgres): relation users already exists but is not a schema-ui users table; use a dedicated empty database (not the cluster 'postgres' database) or restore schema_migrations")
		}
		return fmt.Errorf("fingerprint (postgres): schema is not empty and has no schema-ui users table so it cannot be adopted — use an empty dedicated database or restore the ledger")
	}
	// Postgres aborts the whole tx on 42P07, so adoption must probe then
	// CREATE only missing objects — never "create and ignore duplicate".
	hasLedger, err := tableExistsPG(tx, "schema_migrations")
	if err != nil {
		return err
	}
	if !hasLedger {
		if _, err := tx.Exec(context.Background(), postgresSchemaMigrationsDDL); err != nil {
			return fmt.Errorf("create migration ledger (postgres): %w", err)
		}
	}
	hasRefresh, err := tableExistsPG(tx, "refresh_tokens")
	if err != nil {
		return err
	}
	if !hasRefresh {
		for _, stmt := range postgresBaselineDDL[2:] {
			if _, err := tx.Exec(context.Background(), stmt); err != nil {
				return fmt.Errorf("create baseline (postgres): %w", err)
			}
		}
	}
	return nil
}

func migrateRBACPG(tx kernel.Tx) error {
	for _, stmt := range postgresRBACDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create rbac (postgres): %w", err)
		}
	}
	// On a postgres fresh bootstrap there are no pre-existing users, so this is
	// a no-op here and stays for parity with the sqlite apply.
	return backfillRoles(tx)
}

func migrateSystemDataPG(tx kernel.Tx) error {
	for _, stmt := range postgresSystemDataDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create system-data reconcile tables (postgres): %w", err)
		}
	}
	return nil
}

func migrateAccountLockPG(tx kernel.Tx) error {
	for _, stmt := range postgresAccountLockDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("add account-lock columns (postgres): %w", err)
		}
	}
	return nil
}

func migrateServiceCredentialsPG(tx kernel.Tx) error {
	for _, stmt := range postgresServiceCredentialsDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create service credentials (postgres): %w", err)
		}
	}
	return nil
}

func isEmptyDatabase(tx kernel.Tx) (bool, error) {
	var count int
	err := tx.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%'`,
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("store: count tables: %w", err)
	}
	return count == 0, nil
}

func isEmptyDatabasePG(tx kernel.Tx) (bool, error) {
	var count int
	err := tx.QueryRow(context.Background(), `
		SELECT COUNT(*)
		FROM pg_catalog.pg_class c
		JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
		WHERE n.nspname = current_schema()
		  AND c.relkind = 'r'`).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("store: count tables (postgres): %w", err)
	}
	return count == 0, nil
}

func tableExistsPG(tx kernel.Tx, name string) (bool, error) {
	var exists bool
	err := tx.QueryRow(context.Background(), `
		SELECT EXISTS (
			SELECT 1
			FROM pg_catalog.pg_class c
			JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
			WHERE n.nspname = current_schema()
			  AND c.relkind = 'r'
			  AND c.relname = ?
		)`, name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("store: probe table %s (postgres): %w", name, err)
	}
	return exists, nil
}

// usersLooksLikeSchemaUIPG reports whether current_schema.users has the R2
// identity columns (text id/username/name/roles/password_hash). Extra columns
// from later migrations are allowed so a fully-migrated ledger-less database
// can still be adopted.
func usersLooksLikeSchemaUIPG(tx kernel.Tx) (bool, error) {
	exists, err := tableExistsPG(tx, "users")
	if err != nil || !exists {
		return false, err
	}
	rows, err := tx.Query(context.Background(), `
		SELECT column_name, udt_name
		FROM information_schema.columns
		WHERE table_schema = current_schema() AND table_name = 'users'`)
	if err != nil {
		return false, fmt.Errorf("fingerprint (postgres) users: %w", err)
	}
	defer rows.Close()
	got := map[string]string{}
	for rows.Next() {
		var name, typ string
		if err := rows.Scan(&name, &typ); err != nil {
			return false, fmt.Errorf("fingerprint (postgres) users: scan: %w", err)
		}
		got[name] = strings.ToLower(typ)
	}
	if err := rows.Err(); err != nil {
		return false, fmt.Errorf("fingerprint (postgres) users: %w", err)
	}
	want := map[string]string{
		"id": "text", "username": "text", "name": "text",
		"roles": "text", "password_hash": "text",
	}
	for col, typ := range want {
		if actual, ok := got[col]; !ok || actual != typ {
			return false, nil
		}
	}
	return true, nil
}

func fingerprintR2(tx kernel.Tx) error {
	got := map[string]bool{}
	rows, err := tx.Query(context.Background(), `SELECT name FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%'`)
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
	fkRows, err := tx.Query(context.Background(), `PRAGMA foreign_key_list(refresh_tokens)`)
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
	if err := tx.QueryRow(context.Background(),
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

func fingerprintColumns(tx kernel.Tx, table string, want map[string]string) error {
	rows, err := tx.Query(context.Background(), fmt.Sprintf(`PRAGMA table_info(%s)`, table))
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

func ensureRole(tx kernel.Tx, key string, now int64) error {
	if _, err := tx.Exec(context.Background(),
		`INSERT INTO roles (id, key, name, system, created_at, updated_at)
		 VALUES (?, ?, ?, 0, ?, ?)
		 ON CONFLICT(id) DO NOTHING`,
		"role-"+key, key, key, now, now,
	); err != nil {
		return fmt.Errorf("ensure role %s: %w", key, err)
	}
	return nil
}

func linkUserRole(tx kernel.Tx, userID, key string, now int64) error {
	if !roleKeyPattern.MatchString(key) {
		return fmt.Errorf("invalid role key %q", key)
	}
	if err := ensureRole(tx, key, now); err != nil {
		return err
	}
	if _, err := tx.Exec(context.Background(),
		`INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)
		 ON CONFLICT(user_id, role_id) DO NOTHING`,
		userID, "role-"+key,
	); err != nil {
		return fmt.Errorf("link user %s role %s: %w", userID, key, err)
	}
	return nil
}

func backfillRoles(tx kernel.Tx) error {
	rows, err := tx.Query(context.Background(), `SELECT id, roles FROM users`)
	if err != nil {
		return fmt.Errorf("backfill: list users: %w", err)
	}
	type userRoles struct {
		id, rolesJSON string
	}
	var users []userRoles
	for rows.Next() {
		var u userRoles
		if err := rows.Scan(&u.id, &u.rolesJSON); err != nil {
			rows.Close()
			return err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("backfill: iterate users: %w", err)
	}
	if err := rows.Close(); err != nil {
		return err
	}
	// Collect-then-write: postgres/pgx forbids Exec on the same connection
	// while another query's Rows are still open (sqlite allowed it).
	now := time.Now().UTC().Unix()
	for _, u := range users {
		var keys []string
		if err := json.Unmarshal([]byte(u.rolesJSON), &keys); err != nil {
			return fmt.Errorf("backfill user %s: roles %q is not a JSON array: %w", u.id, u.rolesJSON, err)
		}
		seen := map[string]bool{}
		for _, key := range keys {
			if !roleKeyPattern.MatchString(key) {
				return fmt.Errorf("backfill user %s: invalid role key %q", u.id, key)
			}
			if seen[key] {
				continue
			}
			seen[key] = true
			if err := linkUserRole(tx, u.id, key, now); err != nil {
				return fmt.Errorf("backfill user %s: %w", u.id, err)
			}
		}
	}
	return nil
}
