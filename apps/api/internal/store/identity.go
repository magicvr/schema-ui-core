package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/temporal"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// Startup identity + plan (GOAL-032 / W21).
//
// schema_migrations is the history table — same job as EF Core
// __EFMigrationsHistory: it is the authority for "which catalog versions
// already ran". Identify reads it first. Plan then chooses refuse / noop /
// pending-only / fresh / adopt / restore. A missing history table does NOT
// mean "run every migration" (that is the EF footgun that produces 42P07 on
// a dirty database). Missing history plus foreign tables → refuse.

type dbIdentityKind string

const (
	identityEmpty                dbIdentityKind = "empty"
	identityOursLedger           dbIdentityKind = "ours-ledger"
	identityOursR2               dbIdentityKind = "ours-r2"
	identityOursCompleteNoLedger dbIdentityKind = "ours-complete-no-ledger"
	identityOursPartialNoLedger  dbIdentityKind = "ours-partial-no-ledger"
	identityLostLedgerUnsafe     dbIdentityKind = "lost-ledger-unsafe"
	identityForeign              dbIdentityKind = "foreign"
)

type dbIdentity struct {
	Kind      dbIdentityKind
	Tables    []string
	Applied   []appliedMigration
	OursUsers bool
}

type startupAction string

const (
	actionRefuse           startupAction = "refuse"
	actionNoop             startupAction = "noop"
	actionFresh            startupAction = "fresh"
	actionAdoptR2          startupAction = "adopt-r2"
	actionAdoptThenPending startupAction = "adopt-then-pending"
	actionRestoreLedger    startupAction = "restore-ledger"
	actionApplyPending     startupAction = "apply-pending"
)

type startupPlan struct {
	Action startupAction
	Reason string
}

// sqliteLedgerDDL is the restore-path ledger literal. applied_at is TEXT holding
// the canonical fixed-6 UTC form (workspace-040 v73 / Root D-008); v73 rebuilds
// every pre-existing ledger onto this shape.
const sqliteLedgerDDL = `CREATE TABLE IF NOT EXISTS schema_migrations (
  version    INTEGER PRIMARY KEY,
  name       TEXT NOT NULL UNIQUE,
  checksum   TEXT NOT NULL CHECK (length(checksum) = 64),
  applied_at TEXT NOT NULL
)`

// postgresLedgerDDL is the postgres restore-path ledger literal:
// applied_at is timestamptz(6) (workspace-040 v73).
const postgresLedgerDDL = `CREATE TABLE IF NOT EXISTS schema_migrations (
  version    INTEGER PRIMARY KEY,
  name       TEXT NOT NULL UNIQUE,
  checksum   TEXT NOT NULL CHECK (length(checksum) = 64),
  applied_at timestamptz(6) NOT NULL
)`

func contentTables(tables []string) []string {
	out := make([]string, 0, len(tables))
	for _, name := range tables {
		if name != "schema_migrations" {
			out = append(out, name)
		}
	}
	return out
}

func tableNameSet(tables []string) map[string]bool {
	have := make(map[string]bool, len(tables))
	for _, name := range tables {
		have[name] = true
	}
	return have
}

// completeFingerprintCatalogHead is the compiled catalog max version the
// restore-ledger object set was reviewed against. TestCompleteFingerprintTracksCatalogHead
// fails when the catalog grows past this so the table list is updated.
//
// 87 = the workspace-040 R2 v73–v87 timestamp conversions, which rebuild
// existing tables and create no new objects (lockedHeadExtraTables[87] = {}).
const completeFingerprintCatalogHead = 87

// completeLostLedgerTables must include a table created at/after the catalog
// head (v44 service_credentials, v48 operation_log_session, v51 mail_outbox, v52 mail_config, v64 subjects/vouchers, v65 voucher_batches, v66 telegram_config, v67 telegram_config_connection, v68 telegram ingress tables, v69 telegram outbound table, v70 digital-offer tables)
// so restore-ledger cannot stamp current catalog while later objects are
// missing (A-001 F-001).
var completeLostLedgerTables = []string{
	"users", "refresh_tokens", "operation_log", "jobs",
	"service_credentials", "operation_log_session",
	"mail_outbox",
	"mail_config",
	"email_verification_challenges",
	"password_recovery_challenges",
	"password_policy",
	"user_password_history",
	"user_invites",
	"login_failures",
	"subjects",
	"vouchers",
	"voucher_batches",
	"telegram_config",
	"telegram_sessions",
	"telegram_inbound_messages",
	"telegram_outbound_messages",
	"digital_offers",
	"digital_purchases",
	"digital_entitlements",
}

// postV1CatalogTables: any of these without a complete fingerprint means a
// mid-catalog lost ledger — refuse rather than CREATE TABLE on existing objects
// (A-001 F-002).
var postV1CatalogTables = []string{
	"roles", "user_roles", "permissions", "role_permissions",
	"menu_items", "role_menu_items",
	"operation_log", "operation_log_correlation", "operation_log_archive", "operation_log_session",
	"jobs", "service_credentials", "notifications",
	"system_data_reconcile", "system_data_grants", "site_settings",
	"dict_types", "dict_entries", "wallet_accounts", "recycle_items",
	"captcha_challenges", "scheduled_tasks", "user_mfa",
	"mail_outbox",
	"mail_config",
	"email_verification_challenges",
	"password_recovery_challenges",
	"password_policy",
	"user_password_history",
	"user_invites",
	"login_failures",
	"telegram_outbound_messages",
}

func lostLedgerLooksComplete(tables []string) bool {
	have := tableNameSet(contentTables(tables))
	for _, name := range completeLostLedgerTables {
		if !have[name] {
			return false
		}
	}
	return true
}

func hasPostV1CatalogTables(tables []string) bool {
	have := tableNameSet(contentTables(tables))
	for _, name := range postV1CatalogTables {
		if have[name] {
			return true
		}
	}
	return false
}

func r2Exact(tables []string) bool {
	have := tableNameSet(contentTables(tables))
	return len(have) == 2 && have["users"] && have["refresh_tokens"]
}

func classifyIdentity(tables []string, applied []appliedMigration, oursUsers bool) dbIdentity {
	id := dbIdentity{
		Tables:    append([]string(nil), tables...),
		Applied:   applied,
		OursUsers: oursUsers,
	}
	if applied != nil {
		id.Kind = identityOursLedger
		return id
	}
	if len(contentTables(tables)) == 0 {
		id.Kind = identityEmpty
		return id
	}
	if oursUsers && lostLedgerLooksComplete(tables) {
		id.Kind = identityOursCompleteNoLedger
		return id
	}
	if oursUsers && r2Exact(tables) {
		id.Kind = identityOursR2
		return id
	}
	if oursUsers && hasPostV1CatalogTables(tables) {
		id.Kind = identityLostLedgerUnsafe
		return id
	}
	if oursUsers {
		id.Kind = identityOursPartialNoLedger
		return id
	}
	id.Kind = identityForeign
	return id
}

func planStartup(id dbIdentity, catalog []kernel.MigrationContribution) (startupPlan, error) {
	switch id.Kind {
	case identityEmpty:
		return startupPlan{Action: actionFresh, Reason: "empty database: create schema_migrations and apply catalog"}, nil
	case identityOursLedger:
		pending := pendingMigrations(id.Applied, catalog)
		if len(pending) == 0 {
			return startupPlan{Action: actionNoop, Reason: "schema_migrations matches catalog: no migrate"}, nil
		}
		return startupPlan{Action: actionApplyPending, Reason: fmt.Sprintf("schema_migrations prefix ok: apply %d pending", len(pending))}, nil
	case identityOursR2:
		return startupPlan{Action: actionAdoptR2, Reason: "ledger-less R2 {users, refresh_tokens}: v1 adopt then pending"}, nil
	case identityOursCompleteNoLedger:
		return startupPlan{Action: actionRestoreLedger, Reason: "schema-ui catalog tables present without schema_migrations: restore ledger, do not re-apply"}, nil
	case identityOursPartialNoLedger:
		return startupPlan{Action: actionAdoptThenPending, Reason: "schema-ui users present without complete catalog: v1 adopt then pending"}, nil
	case identityLostLedgerUnsafe:
		return startupPlan{
			Action: actionRefuse,
			Reason: fmt.Sprintf("identity=lost-ledger-unsafe tables=%v; refusing to stamp current catalog or CREATE over existing post-v1 objects", contentTables(id.Tables)),
		}, nil
	case identityForeign:
		return startupPlan{
			Action: actionRefuse,
			Reason: fmt.Sprintf("identity=foreign tables=%v; refusing to migrate a database that is not a schema-ui store", contentTables(id.Tables)),
		}, nil
	default:
		return startupPlan{}, fmt.Errorf("store: unknown identity %q", id.Kind)
	}
}

func usersLooksLikeSQLite(db *sql.DB) (bool, error) {
	var n int
	if err := db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = 'users'`).Scan(&n); err != nil {
		return false, fmt.Errorf("store: probe sqlite users: %w", err)
	}
	if n == 0 {
		return false, nil
	}
	rows, err := db.Query(`PRAGMA table_info(users)`)
	if err != nil {
		return false, fmt.Errorf("store: pragma users: %w", err)
	}
	defer rows.Close()
	got := map[string]string{}
	for rows.Next() {
		var cid, notNull, pk int
		var name, typ string
		var dflt sql.NullString
		if err := rows.Scan(&cid, &name, &typ, &notNull, &dflt, &pk); err != nil {
			return false, fmt.Errorf("store: scan pragma users: %w", err)
		}
		got[name] = strings.ToUpper(typ)
	}
	if err := rows.Err(); err != nil {
		return false, err
	}
	for _, col := range []string{"id", "username", "name", "roles", "password_hash"} {
		if got[col] != "TEXT" {
			return false, nil
		}
	}
	return true, nil
}

func usersLooksLikePostgres(db *sql.DB) (bool, error) {
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM pg_catalog.pg_class c
			JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
			WHERE n.nspname = current_schema()
			  AND c.relkind = 'r'
			  AND c.relname = 'users'
		)`).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("store: probe postgres users: %w", err)
	}
	if !exists {
		return false, nil
	}
	rows, err := db.Query(`
		SELECT column_name, udt_name
		FROM information_schema.columns
		WHERE table_schema = current_schema() AND table_name = 'users'`)
	if err != nil {
		return false, fmt.Errorf("store: postgres users columns: %w", err)
	}
	defer rows.Close()
	got := map[string]string{}
	for rows.Next() {
		var name, typ string
		if err := rows.Scan(&name, &typ); err != nil {
			return false, fmt.Errorf("store: scan postgres users columns: %w", err)
		}
		got[name] = strings.ToLower(typ)
	}
	if err := rows.Err(); err != nil {
		return false, err
	}
	want := map[string]string{
		"id": "text", "username": "text", "name": "text",
		"roles": "text", "password_hash": "text",
	}
	for col, typ := range want {
		if got[col] != typ {
			return false, nil
		}
	}
	return true, nil
}

// ledgerWrite is one prepared ledger row write: the statement to run and the
// value bound to applied_at.
type ledgerWrite struct {
	statement string
	value     any
}

// sqliteLedgerWrite binds a ledger timestamp in the shape of the live
// schema_migrations.applied_at column.
//
// The column is INTEGER unix seconds until v73 converts it (workspace-040 R2)
// and canonical fixed-6 UTC TEXT afterwards, and that conversion happens inside
// the same migration batch — so the shape is probed on the transaction that
// performs the insert (store is the ledger's sole writer, Root D-011) instead of
// assumed. Binding the wrong shape either violates the column contract or
// silently stores a non-canonical value (the exact hazard the v73 design names).
//
// SQLite is dynamically typed, so one statement text serves both shapes.
func sqliteLedgerWrite(ctx context.Context, tx kernel.Tx, now time.Time) (ledgerWrite, error) {
	legacy, err := ledgerAppliedAtIsInteger(ctx, tx,
		`SELECT type FROM pragma_table_info('schema_migrations') WHERE name = 'applied_at'`)
	if err != nil {
		return ledgerWrite{}, err
	}
	write := ledgerWrite{
		statement: `INSERT INTO schema_migrations (version, name, checksum, applied_at) VALUES (?, ?, ?, ?)`,
		value:     temporal.NewValue(now).String(),
	}
	if legacy {
		write.value = now.Unix()
	}
	return write, nil
}

// postgresLedgerWrite is the postgres counterpart: BIGINT unix seconds before
// v73, timestamptz(6) afterwards.
//
// The postgres statements differ per shape on purpose. pgx caches a prepared
// statement per SQL text, so reusing one text across the v73 boundary keeps the
// parameter typed BIGINT (measured: 22P02 "invalid input syntax for type bigint:
// 2026-09-20 12:57:15.914522 +0000 UTC" once the column had become timestamptz).
// Naming the parameter type explicitly in two distinct texts removes the
// dependence on the cached description.
func postgresLedgerWrite(ctx context.Context, tx kernel.Tx, now time.Time) (ledgerWrite, error) {
	legacy, err := ledgerAppliedAtIsInteger(ctx, tx, `
		SELECT data_type FROM information_schema.columns
		WHERE table_schema = current_schema()
		  AND table_name = 'schema_migrations'
		  AND column_name = 'applied_at'`)
	if err != nil {
		return ledgerWrite{}, err
	}
	if legacy {
		return ledgerWrite{
			statement: `INSERT INTO schema_migrations (version, name, checksum, applied_at) VALUES (?, ?, ?, CAST(? AS bigint))`,
			value:     now.Unix(),
		}, nil
	}
	return ledgerWrite{
		statement: `INSERT INTO schema_migrations (version, name, checksum, applied_at) VALUES (?, ?, ?, CAST(? AS timestamptz))`,
		value:     now.Truncate(time.Microsecond),
	}, nil
}

// ledgerAppliedAtIsInteger reports whether the probed declared type is an
// integer family type (the pre-v73 ledger shape).
func ledgerAppliedAtIsInteger(ctx context.Context, tx kernel.Tx, query string) (bool, error) {
	var declared sql.NullString
	if err := tx.QueryRow(ctx, query).Scan(&declared); err != nil {
		return false, fmt.Errorf("probe schema_migrations.applied_at: %w", err)
	}
	return strings.Contains(strings.ToUpper(declared.String), "INT"), nil
}

func stampCatalog(tx kernel.Tx, catalog []kernel.MigrationContribution, write ledgerWrite) error {
	for _, migration := range catalog {
		if _, err := tx.Exec(context.Background(), write.statement,
			migration.Version, migration.Name, migration.Checksum, write.value,
		); err != nil {
			return fmt.Errorf("restore migration ledger %d (%s): %w", migration.Version, migration.Name, err)
		}
	}
	return nil
}

func (s *Store) probeIdentity() (dbIdentity, error) {
	tables, err := s.listUserTables()
	if err != nil {
		return dbIdentity{}, err
	}
	applied, err := s.appliedMigrations()
	if err != nil {
		return dbIdentity{}, err
	}
	ours, err := usersLooksLikeSQLite(s.db)
	if err != nil {
		return dbIdentity{}, err
	}
	return classifyIdentity(tables, applied, ours), nil
}

func (s *Store) listUserTables() ([]string, error) {
	rows, err := s.db.Query(`SELECT name FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%'`)
	if err != nil {
		return nil, fmt.Errorf("store: list sqlite tables: %w", err)
	}
	defer rows.Close()
	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("store: scan sqlite table: %w", err)
		}
		tables = append(tables, name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tables, nil
}

func (s *Store) restoreLedger(catalog []kernel.MigrationContribution) error {
	return s.Run(context.Background(), func(tx kernel.Tx) error {
		if _, err := tx.Exec(context.Background(), sqliteLedgerDDL); err != nil {
			return fmt.Errorf("create schema_migrations: %w", err)
		}
		write, err := sqliteLedgerWrite(context.Background(), tx, time.Now().UTC())
		if err != nil {
			return err
		}
		return stampCatalog(tx, catalog, write)
	})
}

func (p *postgres) probeIdentity(ctx context.Context) (dbIdentity, error) {
	tables, err := p.currentSchemaTables(ctx)
	if err != nil {
		return dbIdentity{}, err
	}
	applied, err := p.appliedMigrationsPG(ctx)
	if err != nil {
		return dbIdentity{}, err
	}
	ours, err := usersLooksLikePostgres(p.db)
	if err != nil {
		return dbIdentity{}, err
	}
	return classifyIdentity(tables, applied, ours), nil
}

func (p *postgres) restoreLedger(ctx context.Context, catalog []kernel.MigrationContribution) error {
	return p.Run(ctx, func(tx kernel.Tx) error {
		if _, err := tx.Exec(ctx, postgresLedgerDDL); err != nil {
			return fmt.Errorf("create schema_migrations: %w", err)
		}
		write, err := postgresLedgerWrite(ctx, tx, time.Now().UTC())
		if err != nil {
			return err
		}
		return stampCatalog(tx, catalog, write)
	})
}
