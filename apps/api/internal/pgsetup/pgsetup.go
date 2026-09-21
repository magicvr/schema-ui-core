// Package pgsetup provisions the local PostgreSQL databases the dev launcher and
// the test suite need before the API can start.
//
// Why this exists: the application never runs CREATE DATABASE. Startup pings the
// configured database and then applies the migration catalog, so on a fresh or
// reset server `dev.cmd start` fails with
//
//	FATAL: database "schema_ui_dev" does not exist (SQLSTATE 3D000)
//
// The helper in this package makes that step explicit and scriptable, and — the
// reason it is not just a thin wrapper — it can bootstrap the very first database
// on an empty server: the maintenance connection must not depend on the
// configured database already existing.
//
// Credentials are per role, mirroring how the rest of the repo already reads
// them: the dev server uses the DB_* namespace (exactly like the API server's own
// loading) and the test suite uses PG_TEST_* (exactly like internal/pgtest). One
// machine may point the two roles at the same server, but the tool must not
// assume it. Both namespaces are read from apps/api/configs/.env when the process
// environment does not provide them. Secrets are never printed.
package pgsetup

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/jackc/pgx/v5"
)

// Roles name the database a caller wants to provision. They select which
// environment namespace supplies the connection details and the database name.
const (
	RoleDev  = "dev"
	RoleTest = "test"
)

// identRe is the unquoted-identifier shape accepted for a database name: the only
// guard between a configured name and the CREATE/DROP DATABASE statements.
var identRe = regexp.MustCompile(`^[a-z_][a-z0-9_]*$`)

// ValidName reports whether a database name is safe to interpolate.
func ValidName(name string) bool { return identRe.MatchString(name) }

// envFilePrefixes are the namespaces read from configs/.env: DB_* is the API
// server's own namespace (dev connection + DB_NAME), PG_TEST_* is the
// integration-test namespace (internal/pgtest) that supplies PG_TEST_DB.
var envFilePrefixes = []string{"DB_", "PG_TEST_"}

// LoadEnvFile loads the namespaces above from apps/api/configs/.env without
// overriding already-set process env. Missing file is fine (CI passes real env
// vars).
func LoadEnvFile() {
	envFile := repoConfigsEnvFile()
	if envFile == "" {
		return
	}
	raw, err := os.ReadFile(envFile)
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(raw), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		eq := strings.IndexByte(line, '=')
		if eq <= 0 {
			continue
		}
		k := strings.TrimSpace(line[:eq])
		if !hasAnyPrefix(k, envFilePrefixes) {
			continue
		}
		if _, set := os.LookupEnv(k); !set {
			v := strings.TrimSpace(line[eq+1:])
			_ = os.Setenv(k, strings.Trim(v, `"'`))
		}
	}
}

func hasAnyPrefix(s string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

func repoConfigsEnvFile() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	dir := filepath.Dir(file) // apps/api/internal/pgsetup
	for {
		if _, err := os.Stat(filepath.Join(dir, "AGENTS.md")); err == nil {
			return filepath.Join(dir, "apps", "api", "configs", ".env")
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

func envOr(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

// Server is one PostgreSQL endpoint plus its credentials, resolved from one
// environment namespace.
type Server struct {
	Role     string
	Host     string
	Port     string
	User     string
	Password string
	SSLMode  string
	// Database is the configured database name for this role (DB_NAME or
	// PG_TEST_DB), defaulting to "postgres" so an unconfigured checkout still
	// resolves to a database that exists.
	Database string
}

// ServerForRole resolves the endpoint and credentials for a role.
func ServerForRole(role string) Server {
	if role == RoleTest {
		return Server{
			Role:     RoleTest,
			Host:     envOr("PG_TEST_HOST", "127.0.0.1"),
			Port:     envOr("PG_TEST_PORT", "5432"),
			User:     envOr("PG_TEST_USER", "postgres"),
			Password: os.Getenv("PG_TEST_PASSWORD"),
			SSLMode:  envOr("PG_TEST_SSLMODE", "disable"),
			Database: envOr("PG_TEST_DB", "postgres"),
		}
	}
	return Server{
		Role:     RoleDev,
		Host:     envOr("DB_HOST", "127.0.0.1"),
		Port:     envOr("DB_PORT", "5432"),
		User:     envOr("DB_USER", ""),
		Password: os.Getenv("DB_PASSWORD"),
		SSLMode:  envOr("DB_SSLMODE", "disable"),
		Database: envOr("DB_NAME", "postgres"),
	}
}

// DatabaseName returns the configured database name for a role.
func DatabaseName(role string) string { return ServerForRole(role).Database }

// DSN returns a DSN for one database on this server. Callers must not log it:
// it carries the password.
func (s Server) DSN(name string) string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(s.User, s.Password),
		Host:   net.JoinHostPort(s.Host, s.Port),
		Path:   "/" + name,
	}
	q := u.Query()
	q.Set("sslmode", s.SSLMode)
	u.RawQuery = q.Encode()
	return u.String()
}

// MaintenanceCandidates lists the databases to try, in order, when opening the
// connection used for CREATE/DROP DATABASE. The configured name comes first so
// the common case keeps parity with the API's own target, then the maintenance
// databases every PostgreSQL server ships: this is what lets the tool create the
// configured database on an empty server instead of failing with 3D000 itself.
func (s Server) MaintenanceCandidates() []string {
	candidates := []string{s.Database, "postgres", "template1"}
	seen := map[string]bool{}
	out := make([]string, 0, len(candidates))
	for _, c := range candidates {
		c = strings.TrimSpace(c)
		if c == "" || seen[c] {
			continue
		}
		seen[c] = true
		out = append(out, c)
	}
	return out
}

// ConnectMaintenance connects to the first candidate that accepts a connection
// and returns it together with the database it landed on.
func (s Server) ConnectMaintenance(ctx context.Context) (*pgx.Conn, string, error) {
	var lastErr error
	for _, candidate := range s.MaintenanceCandidates() {
		conn, err := pgx.Connect(ctx, s.DSN(candidate))
		if err == nil {
			return conn, candidate, nil
		}
		lastErr = err
	}
	return nil, "", fmt.Errorf("no maintenance database reachable on %s (tried %s): %w",
		net.JoinHostPort(s.Host, s.Port), strings.Join(s.MaintenanceCandidates(), ", "), lastErr)
}

// Exists reports whether the database is present on this server.
func (s Server) Exists(ctx context.Context, name string) (bool, error) {
	conn, _, err := s.ConnectMaintenance(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Close(ctx)
	var found any
	err = conn.QueryRow(ctx, "SELECT 1 FROM pg_database WHERE datname = $1", name).Scan(&found)
	if err == pgx.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// Ensure creates the database when it is missing and reports whether it created
// it. It is idempotent and never drops or overwrites an existing database.
func (s Server) Ensure(ctx context.Context, name string) (bool, error) {
	if !ValidName(name) {
		return false, fmt.Errorf("invalid database name %q (want ^[a-z_][a-z0-9_]*$)", name)
	}
	present, err := s.Exists(ctx, name)
	if err != nil {
		return false, err
	}
	if present {
		return false, nil
	}
	conn, _, err := s.ConnectMaintenance(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Close(ctx)
	if _, err := conn.Exec(ctx, "CREATE DATABASE "+name); err != nil {
		return false, err
	}
	return true, nil
}

// Drop removes the database and reports whether one was actually dropped. It is a
// separate, explicit action — nothing in the provisioning path calls it — and it
// is idempotent so cleanup can be re-run safely: a missing database is reported
// as (false, nil), never as an error.
func (s Server) Drop(ctx context.Context, name string) (bool, error) {
	if !ValidName(name) {
		return false, fmt.Errorf("invalid database name %q (want ^[a-z_][a-z0-9_]*$)", name)
	}
	present, err := s.Exists(ctx, name)
	if err != nil {
		return false, err
	}
	if !present {
		return false, nil
	}
	conn, _, err := s.ConnectMaintenance(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Close(ctx)
	if _, err := conn.Exec(ctx, "DROP DATABASE "+name+" WITH (FORCE)"); err != nil {
		// PG < 13 (and servers rejecting FORCE) fall back to the plain form.
		if _, err2 := conn.Exec(ctx, "DROP DATABASE "+name); err2 != nil {
			return false, fmt.Errorf("drop %s: %w (without force: %v)", name, err, err2)
		}
	}
	return true, nil
}

// VerifyMigrated reports whether the database already carries the migration
// ledger (nil means "migrated").
func (s Server) VerifyMigrated(ctx context.Context, name string) error {
	if !ValidName(name) {
		return fmt.Errorf("invalid database name %q", name)
	}
	conn, err := pgx.Connect(ctx, s.DSN(name))
	if err != nil {
		return err
	}
	defer conn.Close(ctx)
	var table any
	if err := conn.QueryRow(ctx, "SELECT to_regclass('public.schema_migrations')").Scan(&table); err != nil {
		return err
	}
	if table == nil {
		return fmt.Errorf("schema_migrations not present yet in %s", name)
	}
	return nil
}

// ListE2EDatabases returns the leftover schema_ui_e2e_* databases, in name order
// (the E2E harness provisions a fresh one per run and drops it after).
func (s Server) ListE2EDatabases(ctx context.Context) ([]string, error) {
	conn, _, err := s.ConnectMaintenance(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)
	rows, err := conn.Query(ctx,
		"SELECT datname FROM pg_database WHERE datname LIKE 'schema_ui_e2e_%' ORDER BY datname")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		out = append(out, name)
	}
	return out, rows.Err()
}
