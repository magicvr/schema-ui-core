package store

import (
	"context"
	"database/sql"
	"errors"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authsession "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
	authmigration "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession/migration"
	compiledmodules "github.com/magicvr/schema-ui-core/apps/api/internal/modules/compiled"
	logincaptchastore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/logincaptcha/store"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"

	"github.com/magicvr/schema-ui-core/apps/api/internal/pgtest"
)

func TestRebindPostgres(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{in: "SELECT 1", want: "SELECT 1"},
		{in: "SELECT * FROM t WHERE a = ? AND b = ?", want: "SELECT * FROM t WHERE a = $1 AND b = $2"},
		{in: "INSERT INTO t (a, b) VALUES (?, ?)", want: "INSERT INTO t (a, b) VALUES ($1, $2)"},
		{in: "UPDATE t SET a = ? WHERE id = ?", want: "UPDATE t SET a = $1 WHERE id = $2"},
		{in: "DELETE FROM t WHERE id = ?", want: "DELETE FROM t WHERE id = $1"},
	}
	for _, c := range cases {
		if got := rebindPostgres(c.in); got != c.want {
			t.Errorf("rebindPostgres(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestSearchPathCandidates(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{
		{in: `"$user", public`, want: []string{"$user", "public"}},
		{in: `"$user", app, public`, want: []string{"$user", "app", "public"}},
		{in: `"MySchema", public`, want: []string{"MySchema", "public"}},
		{in: `pg_catalog`, want: []string{"pg_catalog"}},
		{in: `"weird "" name", public`, want: []string{`weird " name`, "public"}},
		{in: ``, want: []string{}},
	}
	for _, c := range cases {
		if got := searchPathCandidates(c.in); !stringsEqual(got, c.want) {
			t.Errorf("searchPathCandidates(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func stringsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestIsSystemSchema(t *testing.T) {
	for _, s := range []string{"pg_catalog", "information_schema", "pg_toast", "pg_temp_1", "pg_toast_temp_1"} {
		if !isSystemSchema(s) {
			t.Errorf("isSystemSchema(%q) = false, want true", s)
		}
	}
	for _, s := range []string{"public", "app", "probe", "$user", "MySchema"} {
		if isSystemSchema(s) {
			t.Errorf("isSystemSchema(%q) = true, want false", s)
		}
	}
}

// TestAuthsessionPostgresApplyIntegration proves the R3 T3 ported authsession
// migration bodies apply on a live postgres fresh bootstrap: the schema lands,
// Unix time columns are BIGINT, and service_credentials.name is CITEXT (the
// COLLATE NOCASE equivalent). Runs in a dedicated scratch database so the
// shared probe DB stays clean.
func TestAuthsessionPostgresApplyIntegration(t *testing.T) {
	dsn := pgtest.DSN()
	if dsn == "" {
		t.Skip("postgres test env not set (PG_TEST_*); skipping authsession postgres apply integration")
	}
	ctx := context.Background()

	u, err := url.Parse(dsn)
	if err != nil {
		t.Fatal(err)
	}
	const dbName = "r3auth"
	adminDSN := u.String()
	u.Path = "/" + dbName
	authDSN := u.String()

	admin, err := sql.Open("pgx", adminDSN)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = admin.ExecContext(context.Background(), `DROP DATABASE IF EXISTS `+dbName+` WITH (FORCE)`)
		_ = admin.Close()
	})
	if _, err := admin.ExecContext(ctx, `DROP DATABASE IF EXISTS `+dbName+` WITH (FORCE)`); err != nil {
		t.Fatalf("drop prior scratch db: %v", err)
	}
	if _, err := admin.ExecContext(ctx, `CREATE DATABASE `+dbName); err != nil {
		t.Fatalf("create scratch db: %v", err)
	}

	st, err := Open(ctx, OpenOptions{
		Dialect:        kernel.DialectPostgres,
		DSN:            authDSN,
		ConnectTimeout: 10 * time.Second,
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	pg, ok := st.(*postgres)
	if !ok {
		t.Fatal("expected *postgres store")
	}
	t.Cleanup(func() { _ = st.Close() })
	if !st.WasFresh() {
		t.Errorf("fresh scratch db: WasFresh() = false, want true")
	}

	descs := authmigration.Descriptors()
	byVersion := map[int]kernel.MigrationContribution{}
	for _, m := range descs {
		byVersion[m.Version] = m
	}
	// Dependency order of the authsession-owned chain (v1 bootstraps users).
	for _, v := range []int{1, 2, 9, 11, 12, 38, 44} {
		m, ok := byVersion[v]
		if !ok {
			t.Fatalf("missing authsession migration %d", v)
		}
		apply := m.Apply
		if m.ApplyPostgres != nil {
			apply = m.ApplyPostgres
		}
		if err := st.Run(ctx, func(tx kernel.Tx) error { return apply(tx) }); err != nil {
			t.Fatalf("apply %d (%s) on postgres: %v", v, m.Name, err)
		}
	}

	var reg *string
	if err := pg.db.QueryRowContext(ctx, `SELECT to_regclass('schema_migrations')`).Scan(&reg); err != nil {
		t.Fatal(err)
	}
	if reg == nil {
		t.Fatal("schema_migrations ledger not created by the postgres baseline")
	}
	// R1 v1.3: Unix time columns are BIGINT on postgres.
	for _, tc := range []struct{ table, col string }{
		{"users", "created_at"}, {"refresh_tokens", "expires_at"},
		{"roles", "updated_at"}, {"service_credentials", "created_at"},
	} {
		var typ string
		if err := pg.db.QueryRowContext(ctx,
			`SELECT data_type FROM information_schema.columns WHERE table_name = $1 AND column_name = $2`,
			tc.table, tc.col,
		).Scan(&typ); err != nil {
			t.Fatalf("read type %s.%s: %v", tc.table, tc.col, err)
		}
		if typ != "bigint" {
			t.Errorf("%s.%s data_type = %q, want bigint", tc.table, tc.col, typ)
		}
	}
	// R1 v1.4 F-002: COLLATE NOCASE becomes CITEXT on postgres (USER-DEFINED).
	var nameTyp string
	if err := pg.db.QueryRowContext(ctx,
		`SELECT data_type FROM information_schema.columns WHERE table_name = 'service_credentials' AND column_name = 'name'`,
	).Scan(&nameTyp); err != nil {
		t.Fatal(err)
	}
	if nameTyp != "USER-DEFINED" {
		t.Errorf("service_credentials.name data_type = %q, want USER-DEFINED (citext)", nameTyp)
	}
	_ = st.Close()
}

// TestFullCatalogPostgresBootstrapIntegration drives T3 to completion: the
// COMPLETE compiled catalog must fresh-bootstrap on a live postgres via the
// R3 runner (PostgresApply ?? Apply per migration), creating the ledger with
// sqlite-bound checksums. It is the close-out evidence test for GOAL-004.
func TestFullCatalogPostgresBootstrapIntegration(t *testing.T) {
	dsn := pgtest.DSN()
	if dsn == "" {
		t.Skip("postgres test env not set (PG_TEST_*); skipping full catalog postgres bootstrap")
	}
	ctx := context.Background()
	u, err := url.Parse(dsn)
	if err != nil {
		t.Fatal(err)
	}
	const dbName = "r3full"
	adminDSN := u.String()
	u.Path = "/" + dbName
	fullDSN := u.String()

	admin, err := sql.Open("pgx", adminDSN)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = admin.ExecContext(context.Background(), `DROP DATABASE IF EXISTS `+dbName+` WITH (FORCE)`)
		_ = admin.Close()
	})
	if _, err := admin.ExecContext(ctx, `DROP DATABASE IF EXISTS `+dbName+` WITH (FORCE)`); err != nil {
		t.Fatalf("drop prior scratch db: %v", err)
	}
	if _, err := admin.ExecContext(ctx, `CREATE DATABASE `+dbName); err != nil {
		t.Fatalf("create scratch db: %v", err)
	}

	catalog, err := compiledmodules.PersistenceCatalog()
	if err != nil {
		t.Fatal(err)
	}
	if len(catalog) == 0 {
		t.Fatal("empty compiled catalog")
	}

	st, err := Open(ctx, OpenOptions{
		Dialect:        kernel.DialectPostgres,
		DSN:            fullDSN,
		ConnectTimeout: 15 * time.Second,
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	pg := st.(*postgres)
	t.Cleanup(func() { _ = st.Close() })
	if !st.WasFresh() {
		t.Errorf("fresh scratch db: WasFresh() = false, want true")
	}

	if err := pg.migrate(catalog); err != nil {
		t.Fatalf("full catalog postgres bootstrap failed at: %v", err)
	}
	_ = st.Close()

	// Reopen and migrate again = idempotent; ledger matches catalog.
	st2, err := Open(ctx, OpenOptions{Dialect: kernel.DialectPostgres, DSN: fullDSN, ConnectTimeout: 15 * time.Second}, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st2.Close() })
	if st2.WasFresh() {
		t.Errorf("reopened bootstrapped db: WasFresh() = true, want false")
	}
	if err := st2.(*postgres).migrate(catalog); err != nil {
		t.Fatalf("re-migrate full catalog should be idempotent: %v", err)
	}
	var migCount int
	if err := st2.(*postgres).db.QueryRowContext(ctx, `SELECT count(*) FROM schema_migrations`).Scan(&migCount); err != nil {
		t.Fatal(err)
	}
	if migCount != len(catalog) {
		t.Fatalf("ledger rows = %d, want %d (full catalog)", migCount, len(catalog))
	}

	// R1 v1.3 / §3 compliance: ported Unix time + money columns are BIGINT.
	// (operation_log is the known outstanding module — see GOAL-004 T3 notes.)
	assertPGType := func(table, col, want string) {
		t.Helper()
		var typ string
		if err := st2.(*postgres).db.QueryRowContext(ctx,
			`SELECT data_type FROM information_schema.columns WHERE table_name = $1 AND column_name = $2`,
			table, col,
		).Scan(&typ); err != nil {
			t.Fatalf("read type %s.%s: %v", table, col, err)
		}
		if typ != want {
			t.Errorf("%s.%s data_type = %q, want %q", table, col, typ, want)
		}
	}
	for _, tc := range []struct{ table, col string }{
		{"users", "created_at"}, {"refresh_tokens", "expires_at"},
		{"service_credentials", "created_at"},
		{"jobs", "created_at"}, {"jobs", "lease_expires_at"},
		{"notifications", "created_at"}, {"dict_types", "updated_at"},
		{"data_scope_policies", "updated_at"}, {"user_mfa", "created_at"},
		{"mfa_proofs", "expires_at"}, {"captcha_challenges", "expires_at"},
		{"scheduled_tasks", "created_at"}, {"task_runs", "started_at"},
		{"site_settings", "updated_at"},
		{"recycle_items", "deleted_at"},
		{"wallet_accounts", "balance_total"}, {"wallet_accounts", "created_at"},
		{"wallet_ledger_entries", "amount_delta"},
		{"operation_log", "created_at"},
		{"operation_log_archive", "created_at"}, {"operation_log_archive", "archived_at"},
	} {
		assertPGType(tc.table, tc.col, "bigint")
	}

	// R1 v1.3 hard rule — no Unix time column may remain integer/int4 on
	// postgres: any column whose name looks like a time stamp must be bigint.
	var leftover int
	timeNames := []string{
		`created_at`, `updated_at`, `expires_at`, `applied_at`, `archived_at`,
		`lease_expires_at`, `finished_at`, `started_at`, `read_at`, `revoked_at`,
		`last_used_at`, `restored_at`, `deleted_at`, `locked_until`,
	}
	q := `SELECT count(*) FROM information_schema.columns
WHERE table_schema = 'public' AND data_type = 'integer' AND column_name = ANY($1)`
	if err := st2.(*postgres).db.QueryRowContext(ctx, q, timeNames).Scan(&leftover); err != nil {
		t.Fatal(err)
	}
	if leftover != 0 {
		t.Fatalf("%d Unix time column(s) are still integer/int4 on postgres (violates R1 v1.3)", leftover)
	}

	// R4 evidence: a repository migrated to the kernel.Tx port runs on the
	// bootstrapped postgres store (working end-to-end, not just probing).
	captchaRepo := logincaptchastore.NewRepository(st2.(*postgres))
	now := time.Now()
	if err := captchaRepo.SetEnabled(true, now); err != nil {
		t.Fatalf("captcha SetEnabled on postgres: %v", err)
	}
	on, err := captchaRepo.Enabled()
	if err != nil {
		t.Fatalf("captcha Enabled on postgres: %v", err)
	}
	if !on {
		t.Errorf("captcha Enabled() after SetEnabled(true) = false, want true")
	}
	id := "r4repo" + strconv.FormatInt(now.UnixNano(), 10)
	if err := captchaRepo.CreateChallenge(id, "hash", now.Add(time.Minute), now); err != nil {
		t.Fatalf("captcha CreateChallenge on postgres: %v", err)
	}
	matched, err := captchaRepo.ConsumeChallenge(id, "hash", now)
	if err != nil {
		t.Fatalf("captcha ConsumeChallenge on postgres: %v", err)
	}
	if !matched {
		t.Errorf("ConsumeChallenge with correct answer = false, want true")
	}

	// R4 A-004 F-001: runtime substring search must be portable — no sqlite
	// instr(); LOWER(col) LIKE '%' || CAST($1 AS TEXT) || '%' executes on
	// postgres (the -count=1 run above migrated all 48 via the same search
	// shapes without "function instr does not exist").
	var portable int
	if err := st2.(*postgres).db.QueryRowContext(ctx,
		`SELECT count(*) FROM roles WHERE lower(name) LIKE '%' || CAST($1 AS TEXT) || '%'`,
		"admin",
	).Scan(&portable); err != nil {
		t.Fatalf("portable LIKE-concat search on postgres: %v", err)
	}
}

func TestOpenPostgresRequiresDSN(t *testing.T) {
	_, err := Open(context.Background(), OpenOptions{Dialect: kernel.DialectPostgres}, nil)
	if err == nil || !strings.Contains(err.Error(), "DSN") {
		t.Fatalf("postgres without DSN must fail closed, got %v", err)
	}
}

// TestOpenPostgresAppliesNonEmptyCatalogIntegration proves the R3 unblock:
// Open with a non-empty (dual-dialect) catalog runs the postgres migrate
// runner during open (fresh bootstrap), instead of the R2-era fail-closed.
func TestOpenPostgresAppliesNonEmptyCatalogIntegration(t *testing.T) {
	dsn := pgtest.DSN()
	if dsn == "" {
		t.Skip("postgres test env not set (PG_TEST_*); skipping postgres non-empty-catalog open integration")
	}
	ctx := context.Background()
	u, err := url.Parse(dsn)
	if err != nil {
		t.Fatal(err)
	}
	const dbName = "r3open"
	adminDSN := u.String()
	u.Path = "/" + dbName
	openDSN := u.String()

	admin, err := sql.Open("pgx", adminDSN)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = admin.ExecContext(context.Background(), `DROP DATABASE IF EXISTS `+dbName+` WITH (FORCE)`)
		_ = admin.Close()
	})
	if _, err := admin.ExecContext(ctx, `DROP DATABASE IF EXISTS `+dbName+` WITH (FORCE)`); err != nil {
		t.Fatalf("drop prior scratch db: %v", err)
	}
	if _, err := admin.ExecContext(ctx, `CREATE DATABASE `+dbName); err != nil {
		t.Fatalf("create scratch db: %v", err)
	}

	bootstrap := []string{
		`CREATE TABLE schema_migrations (version INTEGER PRIMARY KEY, name TEXT NOT NULL, checksum TEXT NOT NULL, applied_at BIGINT NOT NULL)`,
		`CREATE TABLE r3_open_probe (id SERIAL PRIMARY KEY, note TEXT NOT NULL)`,
	}
	catalog := []kernel.MigrationContribution{
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: "r3.test", Key: "bootstrap"},
			Version:              1,
			Name:                 "bootstrap",
			Checksum:             kernel.MigrationChecksum(bootstrap, "r3:open-bootstrap:v1"),
			Apply: func(tx kernel.Tx) error {
				for _, s := range bootstrap {
					if _, err := tx.Exec(ctx, s); err != nil {
						return err
					}
				}
				return nil
			},
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: "r3.test", Key: "seed"},
			Version:              2,
			Name:                 "seed",
			Checksum:             kernel.MigrationChecksum([]string{`INSERT INTO r3_open_probe (note) VALUES (?)`}, "r3:open-seed:v1"),
			Apply: func(tx kernel.Tx) error {
				_, err := tx.Exec(ctx, `INSERT INTO r3_open_probe (note) VALUES (?)`, "hello")
				return err
			},
		},
	}

	st, err := Open(ctx, OpenOptions{
		Dialect:        kernel.DialectPostgres,
		DSN:            openDSN,
		ConnectTimeout: 10 * time.Second,
	}, catalog)
	if err != nil {
		t.Fatalf("open with non-empty catalog: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	pg := st.(*postgres)

	var migCount, rows int
	if err := pg.db.QueryRowContext(ctx, `SELECT count(*) FROM schema_migrations`).Scan(&migCount); err != nil {
		t.Fatal(err)
	}
	if err := pg.db.QueryRowContext(ctx, `SELECT count(*) FROM r3_open_probe`).Scan(&rows); err != nil {
		t.Fatal(err)
	}
	if migCount != 2 || rows != 1 {
		t.Fatalf("ledger=%d rows=%d, want 2/1 (open-time bootstrap)", migCount, rows)
	}
	if !st.WasFresh() {
		t.Errorf("WasFresh recorded at open (pre-migrate) should be true")
	}
}

// scratchDSN creates an empty database named dbName on the PG_TEST server and
// returns a DSN pointing at it. The database is dropped in t.Cleanup.
//
// WasFresh assertions cannot run against PG_TEST_DB itself: postgresWasFresh
// counts any user table in the first existing search_path schema, and the
// shared probe database typically already holds a leftover catalog (or the
// developer's app schema). Later tests in this file already isolate via
// dedicated scratch databases; R2-era probe/migrate tests must do the same.
func scratchDSN(t *testing.T, adminDSN, dbName string) string {
	t.Helper()
	ctx := context.Background()
	u, err := url.Parse(adminDSN)
	if err != nil {
		t.Fatal(err)
	}
	admin, err := sql.Open("pgx", u.String())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = admin.ExecContext(context.Background(), `DROP DATABASE IF EXISTS `+dbName+` WITH (FORCE)`)
		_ = admin.Close()
	})
	if _, err := admin.ExecContext(ctx, `DROP DATABASE IF EXISTS `+dbName+` WITH (FORCE)`); err != nil {
		t.Fatalf("drop prior scratch db: %v", err)
	}
	if _, err := admin.ExecContext(ctx, `CREATE DATABASE `+dbName); err != nil {
		t.Fatalf("create scratch db: %v", err)
	}
	u.Path = "/" + dbName
	return u.String()
}

func TestOpenPostgresProbeIntegration(t *testing.T) {
	dsn := pgtest.DSN()
	if dsn == "" {
		t.Skip("postgres test env not set (PG_TEST_*); skipping postgres probe integration (no PG = dev/fast-test keeps working)")
	}
	ctx := context.Background()
	probeDSN := scratchDSN(t, dsn, "r2probe")
	probe := func() (kernel.Store, error) {
		return Open(ctx, OpenOptions{
			Dialect:        kernel.DialectPostgres,
			DSN:            probeDSN,
			ConnectTimeout: 10 * time.Second,
		}, nil)
	}

	// An empty scratch database must report WasFresh()=true on the default
	// $user/public search_path (also exercises $user resolution).
	st, err := probe()
	if err != nil {
		t.Fatal(err)
	}
	if st.Dialect() != kernel.DialectPostgres {
		t.Fatalf("dialect = %q, want postgres", st.Dialect())
	}
	if err := st.Ping(ctx); err != nil {
		t.Fatalf("ping: %v", err)
	}
	if !st.WasFresh() {
		t.Errorf("empty scratch db: WasFresh() = false, want true (search_path resolution)")
	}

	// One Run = one tx; placeholders rebound '?' -> $n; commit persists.
	if err := st.Run(ctx, func(tx kernel.Tx) error {
		if _, err := tx.Exec(ctx, "CREATE TABLE IF NOT EXISTS _r2_wasfresh_probe (id serial PRIMARY KEY, name text)"); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, "INSERT INTO _r2_wasfresh_probe (name) VALUES (?)", "hello"); err != nil {
			return err
		}
		var n int
		return tx.QueryRow(ctx, "SELECT count(*) FROM _r2_wasfresh_probe WHERE name = ?", "hello").Scan(&n)
	}); err != nil {
		t.Fatalf("run on postgres: %v", err)
	}
	_ = st.Close()

	// A user base table now exists in the first existing schema (public):
	// WasFresh must flip to false on a fresh open.
	st2, err := probe()
	if err != nil {
		t.Fatal(err)
	}
	defer st2.Close()
	if st2.WasFresh() {
		t.Errorf("db with a user table: WasFresh() = true, want false")
	}
}

// TestPostgresMigrateRunnerIntegration proves the R3 T2 postgres migration
// runner: a portable scratch catalog (rebindable '?') applies on a live PG,
// the ledger records sqlite-bound checksums, re-open is idempotent, and
// checksum drift fails closed. Gated by SCHEMA_UI_R2_PG_DSN (no PG = skip).
func TestPostgresMigrateRunnerIntegration(t *testing.T) {
	dsn := pgtest.DSN()
	if dsn == "" {
		t.Skip("postgres test env not set (PG_TEST_*); skipping postgres migrate runner integration")
	}
	ctx := context.Background()
	scratch := scratchDSN(t, dsn, "r3runner")
	openPG := func() *postgres {
		st, err := Open(ctx, OpenOptions{
			Dialect:        kernel.DialectPostgres,
			DSN:            scratch,
			ConnectTimeout: 10 * time.Second,
		}, nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = st.Close() })
		return st.(*postgres)
	}

	ledgerDDL := `CREATE TABLE schema_migrations (version INTEGER PRIMARY KEY, name TEXT NOT NULL, checksum TEXT NOT NULL, applied_at BIGINT NOT NULL)`
	usersDDL := `CREATE TABLE r3_users (id SERIAL PRIMARY KEY, username TEXT NOT NULL UNIQUE, enabled SMALLINT NOT NULL DEFAULT 1, created_at BIGINT NOT NULL)`
	seedSQL := `INSERT INTO r3_users (username, enabled, created_at) VALUES (?, ?, ?)`
	bootstrapSQL := []string{ledgerDDL, usersDDL}

	catalog := []kernel.MigrationContribution{
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: "r3.test", Key: "bootstrap"},
			Version:              1,
			Name:                 "bootstrap",
			Checksum:             kernel.MigrationChecksum(bootstrapSQL, "r3:bootstrap:v1"),
			Apply: func(tx kernel.Tx) error {
				for _, s := range bootstrapSQL {
					if _, err := tx.Exec(ctx, s); err != nil {
						return err
					}
				}
				return nil
			},
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: "r3.test", Key: "seed"},
			Version:              2,
			Name:                 "seed",
			Checksum:             kernel.MigrationChecksum([]string{seedSQL}, "r3:seed:v1"),
			Apply: func(tx kernel.Tx) error {
				rows := [][3]any{{"alice", 1, int64(1)}, {"bob", 1, int64(2)}}
				for _, r := range rows {
					if _, err := tx.Exec(ctx, seedSQL, r[0], r[1], r[2]); err != nil {
						return err
					}
				}
				return nil
			},
		},
	}

	st := openPG()
	if err := st.migrate(catalog); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	var migCount, userRows int
	if err := st.db.QueryRowContext(ctx, `SELECT count(*) FROM schema_migrations`).Scan(&migCount); err != nil {
		t.Fatal(err)
	}
	if err := st.db.QueryRowContext(ctx, `SELECT count(*) FROM r3_users`).Scan(&userRows); err != nil {
		t.Fatal(err)
	}
	if migCount != 2 || userRows != 2 {
		t.Fatalf("ledger=%d rows=%d, want 2/2", migCount, userRows)
	}
	if !st.WasFresh() {
		t.Errorf("WasFresh recorded at open (pre-migrate) should be true")
	}
	_ = st.Close()

	// Re-open + migrate = idempotent (validateApplied prefix ok, nothing new).
	st2 := openPG()
	if err := st2.migrate(catalog); err != nil {
		t.Fatalf("re-migrate should be idempotent: %v", err)
	}
	if st2.WasFresh() {
		t.Errorf("reopened db with user tables: WasFresh() = true, want false")
	}

	// Checksum drift fails closed (same ledger, drifted catalog copy).
	st3 := openPG()
	drifted := append([]kernel.MigrationContribution(nil), catalog...)
	drifted[1].Checksum = strings.Repeat("a", 64)
	if err := st3.migrate(drifted); err == nil || !strings.Contains(err.Error(), "drift") {
		t.Fatalf("checksum drift must fail closed, got %v", err)
	}
}

// TestPostgresMigrateAdoptsLedgerlessR2 is the postgres twin of
// TestMigrateExistingR2DB: a schema with users+refresh_tokens and no
// schema_migrations must be fingerprinted and adopted, not fail with
// 42P07 relation "users" already exists.
func TestPostgresMigrateAdoptsLedgerlessR2(t *testing.T) {
	dsn := pgtest.DSN()
	if dsn == "" {
		t.Skip("postgres test env not set (PG_TEST_*); skipping ledger-less R2 adopt")
	}
	ctx := context.Background()
	scratch := scratchDSN(t, dsn, "r3r2adopt")

	setup, err := sql.Open("pgx", scratch)
	if err != nil {
		t.Fatal(err)
	}
	for _, stmt := range []string{
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
		`INSERT INTO users (id, username, name, roles, password_hash, created_at, updated_at)
		 VALUES ('user-admin', 'admin', 'Admin', '["admin"]', 'hash-v1', 1, 1)`,
	} {
		if _, err := setup.ExecContext(ctx, stmt); err != nil {
			_ = setup.Close()
			t.Fatalf("r2 fixture: %v", err)
		}
	}
	if err := setup.Close(); err != nil {
		t.Fatal(err)
	}

	catalog, err := compiledmodules.PersistenceCatalog()
	if err != nil {
		t.Fatal(err)
	}
	st, err := Open(ctx, OpenOptions{
		Dialect:        kernel.DialectPostgres,
		DSN:            scratch,
		ConnectTimeout: 10 * time.Second,
	}, catalog)
	if err != nil {
		t.Fatalf("open ledger-less R2 postgres: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })

	pg := st.(*postgres)
	var migCount int
	if err := pg.db.QueryRowContext(ctx, `SELECT count(*) FROM schema_migrations`).Scan(&migCount); err != nil {
		t.Fatal(err)
	}
	if migCount != len(catalog) {
		t.Fatalf("ledger rows = %d, want %d (adopted R2 then remaining catalog)", migCount, len(catalog))
	}
	var username string
	if err := pg.db.QueryRowContext(ctx, `SELECT username FROM users WHERE id = 'user-admin'`).Scan(&username); err != nil {
		t.Fatal(err)
	}
	if username != "admin" {
		t.Fatalf("adopted user = %q, want admin", username)
	}
}

// TestPostgresMigrateRejectsConflictingUsersTable: a leftover users table that
// is not the R2 shape must fail closed with a fingerprint error, not 42P07.
func TestPostgresMigrateRejectsConflictingUsersTable(t *testing.T) {
	dsn := pgtest.DSN()
	if dsn == "" {
		t.Skip("postgres test env not set (PG_TEST_*); skipping conflicting-users fingerprint")
	}
	ctx := context.Background()
	scratch := scratchDSN(t, dsn, "r3conflict")

	setup, err := sql.Open("pgx", scratch)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := setup.ExecContext(ctx, `CREATE TABLE users (id SERIAL PRIMARY KEY, email TEXT NOT NULL)`); err != nil {
		_ = setup.Close()
		t.Fatalf("conflict fixture: %v", err)
	}
	if err := setup.Close(); err != nil {
		t.Fatal(err)
	}

	catalog, err := compiledmodules.PersistenceCatalog()
	if err != nil {
		t.Fatal(err)
	}
	st, err := Open(ctx, OpenOptions{
		Dialect:        kernel.DialectPostgres,
		DSN:            scratch,
		ConnectTimeout: 10 * time.Second,
	}, catalog)
	if err == nil {
		_ = st.Close()
		t.Fatal("open against a conflicting users table must fail closed")
	}
	msg := err.Error()
	if !strings.Contains(msg, "identity=foreign") && !strings.Contains(msg, "refusing to migrate") {
		t.Fatalf("want foreign-identity refuse, got %v", err)
	}
	if strings.Contains(msg, "42P07") || strings.Contains(msg, "SQLSTATE 42P07") {
		t.Fatalf("must not surface CREATE TABLE 42P07, got %v", err)
	}
}

// TestPostgresMigrateRestoresLostLedger: a fully-migrated postgres that lost
// schema_migrations (the live-dev failure mode) must restore the ledger on
// reopen instead of failing with 42P07 or replaying operation_log rebuilds.
func TestPostgresMigrateRestoresLostLedger(t *testing.T) {
	dsn := pgtest.DSN()
	if dsn == "" {
		t.Skip("postgres test env not set (PG_TEST_*); skipping lost-ledger restore")
	}
	ctx := context.Background()
	scratch := scratchDSN(t, dsn, "r3lostled")
	catalog, err := compiledmodules.PersistenceCatalog()
	if err != nil {
		t.Fatal(err)
	}

	st, err := Open(ctx, OpenOptions{
		Dialect:        kernel.DialectPostgres,
		DSN:            scratch,
		ConnectTimeout: 15 * time.Second,
	}, catalog)
	if err != nil {
		t.Fatalf("initial bootstrap: %v", err)
	}
	pg := st.(*postgres)
	if _, err := pg.db.ExecContext(ctx, `INSERT INTO users (id, username, name, roles, password_hash, created_at, updated_at)
		VALUES ('user-keep', 'keeper', 'Keeper', '["admin"]', 'hash', 1, 1)`); err != nil {
		_ = st.Close()
		t.Fatalf("seed user: %v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	admin, err := sql.Open("pgx", scratch)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := admin.ExecContext(ctx, `DROP TABLE schema_migrations`); err != nil {
		_ = admin.Close()
		t.Fatalf("drop ledger: %v", err)
	}
	if err := admin.Close(); err != nil {
		t.Fatal(err)
	}

	st2, err := Open(ctx, OpenOptions{
		Dialect:        kernel.DialectPostgres,
		DSN:            scratch,
		ConnectTimeout: 15 * time.Second,
	}, catalog)
	if err != nil {
		t.Fatalf("reopen after lost ledger: %v", err)
	}
	t.Cleanup(func() { _ = st2.Close() })
	pg2 := st2.(*postgres)

	var migCount int
	if err := pg2.db.QueryRowContext(ctx, `SELECT count(*) FROM schema_migrations`).Scan(&migCount); err != nil {
		t.Fatal(err)
	}
	if migCount != len(catalog) {
		t.Fatalf("restored ledger rows = %d, want %d", migCount, len(catalog))
	}
	var username string
	if err := pg2.db.QueryRowContext(ctx, `SELECT username FROM users WHERE id = 'user-keep'`).Scan(&username); err != nil {
		t.Fatal(err)
	}
	if username != "keeper" {
		t.Fatalf("preserved user = %q, want keeper", username)
	}
	var jobs int
	if err := pg2.db.QueryRowContext(ctx, `SELECT count(*) FROM jobs`).Scan(&jobs); err != nil {
		t.Fatalf("jobs table missing after ledger restore: %v", err)
	}
}

// TestPostgresMigrateRefusesIncompleteLostLedger: jobs-era schema without
// catalog-head objects must not stamp v44+ as applied (A-001 F-001).
func TestPostgresMigrateRefusesIncompleteLostLedger(t *testing.T) {
	dsn := pgtest.DSN()
	if dsn == "" {
		t.Skip("postgres test env not set (PG_TEST_*); skipping incomplete lost-ledger refuse")
	}
	ctx := context.Background()
	scratch := scratchDSN(t, dsn, "r3lostinc")
	catalog, err := compiledmodules.PersistenceCatalog()
	if err != nil {
		t.Fatal(err)
	}

	st, err := Open(ctx, OpenOptions{
		Dialect:        kernel.DialectPostgres,
		DSN:            scratch,
		ConnectTimeout: 15 * time.Second,
	}, catalog)
	if err != nil {
		t.Fatalf("initial bootstrap: %v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	admin, err := sql.Open("pgx", scratch)
	if err != nil {
		t.Fatal(err)
	}
	for _, stmt := range []string{
		`DROP TABLE schema_migrations`,
		`DROP TABLE IF EXISTS service_credentials CASCADE`,
		`DROP TABLE IF EXISTS operation_log_session CASCADE`,
	} {
		if _, err := admin.ExecContext(ctx, stmt); err != nil {
			_ = admin.Close()
			t.Fatalf("%s: %v", stmt, err)
		}
	}
	if err := admin.Close(); err != nil {
		t.Fatal(err)
	}

	st2, err := Open(ctx, OpenOptions{
		Dialect:        kernel.DialectPostgres,
		DSN:            scratch,
		ConnectTimeout: 15 * time.Second,
	}, catalog)
	if err == nil {
		_ = st2.Close()
		t.Fatal("incomplete lost ledger must fail closed, not stamp current catalog")
	}
	if !strings.Contains(err.Error(), "lost-ledger-unsafe") {
		t.Fatalf("want lost-ledger-unsafe refuse, got %v", err)
	}
	check, err := sql.Open("pgx", scratch)
	if err != nil {
		t.Fatal(err)
	}
	defer check.Close()
	var exists bool
	if err := check.QueryRowContext(ctx, `SELECT to_regclass('schema_migrations') IS NOT NULL`).Scan(&exists); err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Fatal("refuse must not create schema_migrations")
	}
}

// TestPostgresCrossModuleSharedTx proves R5 U3: a caller-owned transaction can
// span multiple repositories on postgres — authsession credential creation
// plus an operation-log audit row commit together, and both roll back when the
// audit fails. Uses the logincaptcha-repo trick of running directly over a
// bootstrapped postgres store.
func TestPostgresCrossModuleSharedTx(t *testing.T) {
	dsn := pgtest.DSN()
	if dsn == "" {
		t.Skip("postgres test env not set (PG_TEST_*); skipping postgres cross-module tx integration")
	}
	ctx := context.Background()
	u, err := url.Parse(dsn)
	if err != nil {
		t.Fatal(err)
	}
	const dbName = "r5tx"
	adminDSN := u.String()
	u.Path = "/" + dbName
	txDSN := u.String()

	admin, err := sql.Open("pgx", adminDSN)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = admin.ExecContext(context.Background(), `DROP DATABASE IF EXISTS `+dbName+` WITH (FORCE)`)
		_ = admin.Close()
	})
	if _, err := admin.ExecContext(ctx, `DROP DATABASE IF EXISTS `+dbName+` WITH (FORCE)`); err != nil {
		t.Fatalf("drop prior scratch db: %v", err)
	}
	if _, err := admin.ExecContext(ctx, `CREATE DATABASE `+dbName); err != nil {
		t.Fatalf("create scratch db: %v", err)
	}

	catalog, err := compiledmodules.PersistenceCatalog()
	if err != nil {
		t.Fatal(err)
	}
	st, err := Open(ctx, OpenOptions{Dialect: kernel.DialectPostgres, DSN: txDSN, ConnectTimeout: 15 * time.Second}, catalog)
	if err != nil {
		t.Fatalf("open+bootstrap: %v", err)
	}
	pg := st.(*postgres)
	t.Cleanup(func() { _ = st.Close() })

	authRepo := authsession.NewRepository(st)
	opRepo := operationlog.NewRepository(st)
	now := time.Date(2026, 8, 20, 12, 0, 0, 0, time.UTC)

	credential := authsession.ServiceCredential{
		ID: "0123456789abcdef0123456789abcdef", Name: "cross-module",
		TokenPrefix: "sui_sc_crossmod", TokenHash: strings.Repeat("a", 64),
		Scopes: []string{"records.read"}, ExpiresAt: now.Add(24 * time.Hour),
		CreatedBy: "user-admin", CreatedAt: now, UpdatedAt: now,
	}
	op := operationlog.Operation{
		ID: "op-cross-" + strconv.FormatInt(now.UnixNano(), 10), Event: operationlog.EventServiceCredentialUse,
		ActorID: "service-credential:" + credential.ID, ActorName: credential.Name,
		CorrelationID: "corr-cross", CreatedAt: now,
	}

	// Commit path: audit writes the op-log row in the same tx; both persist.
	opID := op.ID
	if err := authRepo.CreateServiceCredential(credential, func(tx kernel.Tx) error {
		return opRepo.RecordOperationTx(tx, op)
	}); err != nil {
		t.Fatalf("create credential with audit (commit): %v", err)
	}
	var credCount, opCount int
	if err := pg.db.QueryRowContext(ctx, `SELECT count(*) FROM service_credentials WHERE id = $1`, credential.ID).Scan(&credCount); err != nil {
		t.Fatal(err)
	}
	if err := pg.db.QueryRowContext(ctx, `SELECT count(*) FROM operation_log WHERE id = $1`, opID).Scan(&opCount); err != nil {
		t.Fatal(err)
	}
	if credCount != 1 || opCount != 1 {
		t.Fatalf("cross-module commit: cred=%d op=%d, want 1/1", credCount, opCount)
	}

	// Rollback path: audit failure rolls back BOTH the credential and audit row.
	cred2 := credential
	cred2.ID = "1123456789abcdef0123456789abcdef"
	cred2.Name = "cross-module-rollback"
	op2 := op
	op2.ID = "op-rollback"
	if err := authRepo.CreateServiceCredential(cred2, func(tx kernel.Tx) error {
		_ = opRepo.RecordOperationTx(tx, op2)
		return errors.New("audit forced failure")
	}); err == nil {
		t.Fatal("expected audit failure to fail the creation")
	}
	var cred2Count, op2Count int
	if err := pg.db.QueryRowContext(ctx, `SELECT count(*) FROM service_credentials WHERE id = $1`, cred2.ID).Scan(&cred2Count); err != nil {
		t.Fatal(err)
	}
	if err := pg.db.QueryRowContext(ctx, `SELECT count(*) FROM operation_log WHERE id = $1`, op2.ID).Scan(&op2Count); err != nil {
		t.Fatal(err)
	}
	if cred2Count != 0 || op2Count != 0 {
		t.Fatalf("cross-module rollback: cred=%d op=%d, want 0/0", cred2Count, op2Count)
	}
}

// TestPostgresDataMigrationPrototype closes R5 A-001 F-002: a minimal
// SQLite→Postgres logical data-migration prototype through the domain
// repository — read a user from a seeded SQLite store, write it to a
// fresh-bootstrapped postgres store, verify it round-trips.
func TestPostgresDataMigrationPrototype(t *testing.T) {
	dsn := pgtest.DSN()
	if dsn == "" {
		t.Skip("postgres test env not set (PG_TEST_*); skipping postgres data migration prototype")
	}
	ctx := context.Background()

	// Source: a seeded sqlite store.
	src, err := OpenSeeded(filepath.Join(t.TempDir(), "src.db"), "admin", "hash", true)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = src.Close() })
	srcAuth := authsession.NewRepository(src)
	admin, err := srcAuth.UserByUsername("admin")
	if err != nil {
		t.Fatalf("read source user: %v", err)
	}

	// Target: a fresh postgres store (48 migrations), then copy the user.
	u, err := url.Parse(dsn)
	if err != nil {
		t.Fatal(err)
	}
	const dbName = "r5mig"
	adminDSN := u.String()
	u.Path = "/" + dbName
	migDSN := u.String()
	adm, err := sql.Open("pgx", adminDSN)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = adm.ExecContext(context.Background(), `DROP DATABASE IF EXISTS `+dbName+` WITH (FORCE)`)
		_ = adm.Close()
	})
	if _, err := adm.ExecContext(ctx, `DROP DATABASE IF EXISTS `+dbName+` WITH (FORCE)`); err != nil {
		t.Fatalf("drop prior scratch db: %v", err)
	}
	if _, err := adm.ExecContext(ctx, `CREATE DATABASE `+dbName); err != nil {
		t.Fatalf("create scratch db: %v", err)
	}
	catalog, err := compiledmodules.PersistenceCatalog()
	if err != nil {
		t.Fatal(err)
	}
	target, err := Open(ctx, OpenOptions{Dialect: kernel.DialectPostgres, DSN: migDSN, ConnectTimeout: 15 * time.Second}, catalog)
	if err != nil {
		t.Fatalf("open+bootstrap target: %v", err)
	}
	t.Cleanup(func() { _ = target.Close() })
	targetAuth := authsession.NewRepository(target)

	// Copy the source user's persisted fields to the target (logical migration).
	migrated := authsession.User{
		ID: admin.ID, Username: admin.Username, Name: admin.Name,
		Roles: admin.Roles, PasswordHash: admin.PasswordHash,
		MustChangePassword: admin.MustChangePassword,
		CreatedAt:          admin.CreatedAt, UpdatedAt: admin.UpdatedAt,
	}
	if err := targetAuth.CreateUser(migrated); err != nil {
		t.Fatalf("write migrated user to postgres: %v", err)
	}
	got, err := targetAuth.UserByUsername(admin.Username)
	if err != nil {
		t.Fatalf("read migrated user from postgres: %v", err)
	}
	if got.ID != admin.ID || got.Name != admin.Name {
		t.Fatalf("migrated user mismatch: got %q/%q want %q/%q", got.ID, got.Name, admin.ID, admin.Name)
	}
	if len(got.Roles) != len(admin.Roles) {
		t.Fatalf("roles mismatch: %v vs %v", got.Roles, admin.Roles)
	}
}
