// Command e2e-pgset provisions the dedicated scratch PostgreSQL database the
// browser E2E suite uses when the harness dialect contract is postgres
// (W24 / GOAL-035). The assistant is the same pattern as internal/pgtest and
// the CI api-postgres job: a fresh dedicated database per run, created before
// and dropped after — never the developer's shared database.
//
// Connection details come from the environment first, then apps/api/configs/.env
// (gitignored), exactly like the API server's own loading so local developers
// get the same behavior with zero extra configuration. CI passes process env.
//
// Usage:
//
//	go run ./cmd/e2e-pgset create <name>   # CREATE DATABASE name
//	go run ./cmd/e2e-pgset drop <name>     # DROP DATABASE name (no active conns)
//	go run ./cmd/e2e-pgset verify <name>   # exit 0 once schema_migrations exists
//	go run ./cmd/e2e-pgset list            # list schema_ui_e2e_* databases (leftovers)
package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

func main() {
	loadDBEnvFile()
	if len(os.Args) == 2 && os.Args[1] == "list" {
		listExisting()
		return
	}
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: e2e-pgset create|drop|verify <name> | list")
		os.Exit(2)
	}
	action, name := os.Args[1], os.Args[2]
	if !identRe.MatchString(name) {
		fmt.Fprintln(os.Stderr, "invalid database name:", name)
		os.Exit(2)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	switch action {
	case "create":
		conn, err := pgx.Connect(ctx, maintenanceDSN())
		if err != nil {
			fmt.Fprintln(os.Stderr, "connect:", err)
			os.Exit(1)
		}
		defer conn.Close(ctx)
		if _, err := conn.Exec(ctx, "CREATE DATABASE "+name); err != nil {
			fmt.Fprintln(os.Stderr, "create:", err)
			os.Exit(1)
		}
		fmt.Println("created", name)
	case "drop":
		conn, err := pgx.Connect(ctx, maintenanceDSN())
		if err != nil {
			fmt.Fprintln(os.Stderr, "connect:", err)
			os.Exit(1)
		}
		defer conn.Close(ctx)
		_, err = conn.Exec(ctx, "DROP DATABASE "+name+" WITH (FORCE)")
		if err != nil {
			// PG < 13 fallback (and any server that rejects the FORCE clause).
			if _, err2 := conn.Exec(ctx, "DROP DATABASE "+name); err2 != nil {
				fmt.Fprintln(os.Stderr, "drop:", err)
				fmt.Fprintln(os.Stderr, "drop (no force):", err2)
				os.Exit(1)
			}
		}
		fmt.Println("dropped", name)
	case "verify":
		conn, err := pgx.Connect(ctx, dbDSN(name))
		if err != nil {
			fmt.Fprintln(os.Stderr, "verify: connect:", err)
			os.Exit(1)
		}
		defer conn.Close(ctx)
		var table any
		if err := conn.QueryRow(ctx, "SELECT to_regclass('public.schema_migrations')").Scan(&table); err != nil {
			fmt.Fprintln(os.Stderr, "verify:", err)
			os.Exit(1)
		}
		if table == nil {
			fmt.Fprintln(os.Stderr, "verify: schema_migrations not present yet in", name)
			os.Exit(1)
		}
		fmt.Println("verified", name)
	default:
		fmt.Fprintln(os.Stderr, "unknown action:", action)
		os.Exit(2)
	}
}

func listExisting() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctx, maintenanceDSN())
	if err != nil {
		fmt.Fprintln(os.Stderr, "connect:", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)
	rows, err := conn.Query(ctx,
		"SELECT datname FROM pg_database WHERE datname LIKE 'schema_ui_e2e_%' ORDER BY datname")
	if err != nil {
		fmt.Fprintln(os.Stderr, "list:", err)
		os.Exit(1)
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			fmt.Fprintln(os.Stderr, "list:", err)
			os.Exit(1)
		}
		fmt.Println(name)
		count++
	}
	fmt.Fprintf(os.Stderr, "%d schema_ui_e2e_* database(s)\n", count)
}

var identRe = regexp.MustCompile(`^[a-z_][a-z0-9_]*$`)

// loadDBEnvFile loads DB_* keys from apps/api/configs/.env without overriding
// already-set process env (mirrors config.Load and internal/pgtest).
func loadDBEnvFile() {
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
		if !strings.HasPrefix(k, "DB_") {
			continue
		}
		if _, set := os.LookupEnv(k); !set {
			v := strings.TrimSpace(line[eq+1:])
			_ = os.Setenv(k, strings.Trim(v, `"'`))
		}
	}
}

func repoConfigsEnvFile() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	dir := filepath.Dir(file) // apps/api/cmd/e2e-pgset
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

func dbDSN(name string) string {
	host := envOr("DB_HOST", "127.0.0.1")
	port := envOr("DB_PORT", "5432")
	user := envOr("DB_USER", "")
	pass := os.Getenv("DB_PASSWORD")
	sslmode := envOr("DB_SSLMODE", "disable")
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, pass),
		Host:   host + ":" + port,
		Path:   "/" + name,
	}
	q := u.Query()
	q.Set("sslmode", sslmode)
	u.RawQuery = q.Encode()
	return u.String()
}

func maintenanceDSN() string {
	return dbDSN(envOr("DB_NAME", "postgres"))
}

func envOr(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}