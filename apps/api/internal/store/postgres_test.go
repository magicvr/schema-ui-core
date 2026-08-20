package store

import (
	"context"
	"database/sql"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authmigration "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession/migration"
	compiledmodules "github.com/magicvr/schema-ui-core/apps/api/internal/modules/compiled"
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
	dsn := os.Getenv("SCHEMA_UI_R2_PG_DSN")
	if dsn == "" {
		t.Skip("SCHEMA_UI_R2_PG_DSN not set; skipping authsession postgres apply integration")
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
	dsn := os.Getenv("SCHEMA_UI_R2_PG_DSN")
	if dsn == "" {
		t.Skip("SCHEMA_UI_R2_PG_DSN not set; skipping full catalog postgres bootstrap")
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
	} {
		assertPGType(tc.table, tc.col, "bigint")
	}
}

func TestOpenPostgresRequiresDSN(t *testing.T) {
	_, err := Open(context.Background(), OpenOptions{Dialect: kernel.DialectPostgres}, nil)
	if err == nil || !strings.Contains(err.Error(), "DSN") {
		t.Fatalf("postgres without DSN must fail closed, got %v", err)
	}
}

func TestOpenPostgresFailsClosedOnNonEmptyCatalog(t *testing.T) {
	// This must fail BEFORE any network connection: R2 cannot apply the
	// SQLite-specific compiled catalog to postgres.
	_, err := Open(context.Background(), OpenOptions{
		Dialect: kernel.DialectPostgres,
		DSN:     "postgres://user:pass@127.0.0.1:5432/nonexistent",
	}, []kernel.MigrationContribution{{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: "m", Key: "m1"},
		Name:                 "m1",
	}})
	if err == nil || !strings.Contains(err.Error(), "R3") {
		t.Fatalf("postgres with non-empty catalog must fail closed, got %v", err)
	}
}

func TestOpenPostgresProbeIntegration(t *testing.T) {
	dsn := os.Getenv("SCHEMA_UI_R2_PG_DSN")
	if dsn == "" {
		t.Skip("SCHEMA_UI_R2_PG_DSN not set; skipping postgres probe integration (no PG = dev/fast-test keeps working)")
	}
	ctx := context.Background()
	probe := func() (kernel.Store, error) {
		return Open(ctx, OpenOptions{
			Dialect:        kernel.DialectPostgres,
			DSN:            dsn,
			ConnectTimeout: 10 * time.Second,
		}, nil)
	}

	// A freshly created, empty probe database must report WasFresh()=true on
	// the default $user/public search_path (also exercises $user resolution).
	// The probe database is shared across tests/processes, so first clean any
	// leftover scratch table to make the assertion deterministic.
	seed, err := probe()
	if err != nil {
		t.Fatal(err)
	}
	if err := seed.Run(ctx, func(tx kernel.Tx) error {
		_, err := tx.Exec(ctx, "DROP TABLE IF EXISTS _r2_wasfresh_probe, r3_users, schema_migrations")
		return err
	}); err != nil {
		t.Fatalf("clean slate: %v", err)
	}
	_ = seed.Close()

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
		t.Errorf("empty probe db: WasFresh() = false, want true (search_path resolution)")
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
	// WasFresh must flip to false on a fresh open; the table also drops cleanly
	// so the test stays idempotent.
	st2, err := probe()
	if err != nil {
		t.Fatal(err)
	}
	defer st2.Close()
	if st2.WasFresh() {
		t.Errorf("db with a user table: WasFresh() = true, want false")
	}
	if err := st2.Run(ctx, func(tx kernel.Tx) error {
		_, err := tx.Exec(ctx, "DROP TABLE IF EXISTS _r2_wasfresh_probe")
		return err
	}); err != nil {
		t.Fatalf("cleanup: %v", err)
	}
}

// TestPostgresMigrateRunnerIntegration proves the R3 T2 postgres migration
// runner: a portable scratch catalog (rebindable '?') applies on a live PG,
// the ledger records sqlite-bound checksums, re-open is idempotent, and
// checksum drift fails closed. Gated by SCHEMA_UI_R2_PG_DSN (no PG = skip).
func TestPostgresMigrateRunnerIntegration(t *testing.T) {
	dsn := os.Getenv("SCHEMA_UI_R2_PG_DSN")
	if dsn == "" {
		t.Skip("SCHEMA_UI_R2_PG_DSN not set; skipping postgres migrate runner integration")
	}
	ctx := context.Background()
	openPG := func() *postgres {
		st, err := Open(ctx, OpenOptions{
			Dialect:        kernel.DialectPostgres,
			DSN:            dsn,
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

	// Deterministic clean slate BEFORE WasFresh assertions: drop any leftover
	// scratch tables from a prior run, then open a fresh probe.
	seed := openPG()
	if _, err := seed.db.ExecContext(ctx, `DROP TABLE IF EXISTS r3_users, schema_migrations, _r2_wasfresh_probe`); err != nil {
		t.Fatalf("clean slate: %v", err)
	}
	_ = seed.Close()

	st := openPG()
	t.Cleanup(func() {
		// Use a dedicated connection (st may already be closed in the body) so
		// the scratch tables are always removed for the next test / process.
		if db, err := sql.Open("pgx", dsn); err == nil {
			_, _ = db.ExecContext(context.Background(), `DROP TABLE IF EXISTS r3_users, schema_migrations, _r2_wasfresh_probe`)
			_ = db.Close()
		}
	})

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
