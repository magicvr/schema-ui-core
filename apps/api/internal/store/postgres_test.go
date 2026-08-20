package store

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
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
	st, err := Open(ctx, OpenOptions{
		Dialect:        kernel.DialectPostgres,
		DSN:            dsn,
		ConnectTimeout: 10 * time.Second,
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })

	if st.Dialect() != kernel.DialectPostgres {
		t.Fatalf("dialect = %q, want postgres", st.Dialect())
	}
	if err := st.Ping(ctx); err != nil {
		t.Fatalf("ping: %v", err)
	}
	// WasFresh must evaluate without error (empty default search_path may be
	// fresh or not depending on the database, but the query must not fail).
	_ = st.WasFresh()

	// One Run = one tx; placeholders rebound '?' -> $n; temp table lives on the
	// tx connection and is dropped at session close.
	if err := st.Run(ctx, func(tx kernel.Tx) error {
		if _, err := tx.Exec(ctx, "CREATE TEMP TABLE _r2_probe (id serial PRIMARY KEY, name text)"); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, "INSERT INTO _r2_probe (name) VALUES (?)", "hello"); err != nil {
			return err
		}
		var n int
		return tx.QueryRow(ctx, "SELECT count(*) FROM _r2_probe WHERE name = ?", "hello").Scan(&n)
	}); err != nil {
		t.Fatalf("run on postgres: %v", err)
	}
}
