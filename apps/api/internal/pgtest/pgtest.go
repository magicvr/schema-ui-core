// Package pgtest provides the local PostgreSQL test-connection DSN for
// integration tests. Connection details (host / port / user / password /
// database / sslmode) come only from the environment — CI variables or
// apps/api/configs/.env — never from code or committed config, so the same
// tests run against any developer's local postgres without editing sources.
package pgtest

import (
	"net"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	pgEnvOnce sync.Once
)

// loadPGOptions loads apps/api/configs/.env (gitignored) into the process env
// for the PG_TEST_* namespace without overriding already-set variables, so
// local developers can keep test postgres credentials in configs/.env while CI
// passes real environment variables. The repo root is located from this file's
// path so it works regardless of the `go test` working directory.
func loadPGOptions() {
	pgEnvOnce.Do(func() {
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
			if !strings.HasPrefix(k, "PG_TEST_") {
				continue
			}
			if _, set := os.LookupEnv(k); !set {
				v := strings.TrimSpace(line[eq+1:])
				_ = os.Setenv(k, strings.Trim(v, `"'`))
			}
		}
	})
}

func repoConfigsEnvFile() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	dir := filepath.Dir(file) // apps/api/internal/pgtest
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
	loadPGOptions()
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

// DSN returns the postgres DSN for Go integration tests. Priority:
// SCHEMA_UI_R2_PG_DSN (legacy alias) > PG_TEST_DSN > (PG_TEST_HOST / PORT /
// USER / PASSWORD / DB / SSLMODE). Returns "" when postgres testing is not
// configured so the gated tests skip cleanly (CI alone runs the postgres job).
func DSN() string {
	loadPGOptions()
	if dsn := strings.TrimSpace(os.Getenv("SCHEMA_UI_R2_PG_DSN")); dsn != "" {
		return dsn
	}
	if dsn := strings.TrimSpace(os.Getenv("PG_TEST_DSN")); dsn != "" {
		return dsn
	}
	if strings.TrimSpace(os.Getenv("PG_TEST_PASSWORD")) == "" {
		return ""
	}
	host := envOr("PG_TEST_HOST", "127.0.0.1")
	port := envOr("PG_TEST_PORT", "5432")
	user := envOr("PG_TEST_USER", "postgres")
	db := envOr("PG_TEST_DB", "postgres")
	sslmode := envOr("PG_TEST_SSLMODE", "disable")
	if host == "" || port == "" || user == "" || db == "" {
		return ""
	}
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, os.Getenv("PG_TEST_PASSWORD")),
		Host:   net.JoinHostPort(host, port),
		Path:   "/" + db,
	}
	q := u.Query()
	q.Set("sslmode", sslmode)
	u.RawQuery = q.Encode()
	return u.String()
}
