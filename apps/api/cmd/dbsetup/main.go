// Command dbsetup provisions the PostgreSQL databases this checkout needs before
// the API can start or the test suite can run.
//
// The application never runs CREATE DATABASE: startup pings the configured
// database and then applies the migration catalog, so on a fresh or reset server
// `dev.cmd start` fails with
//
//	FATAL: database "schema_ui_dev" does not exist (SQLSTATE 3D000)
//
// Run this once (or after resetting the server) to create the databases:
//
//	go run ./cmd/dbsetup              # dev (DB_NAME) + test (PG_TEST_DB)
//	go run ./cmd/dbsetup --dev-only   # only the dev database
//	go run ./cmd/dbsetup --test-only  # only the test database
//	go run ./cmd/dbsetup extra_db     # also ensure the named databases (dev role)
//
// It is idempotent: existing databases are reported and left untouched, and it
// never drops anything. Migrations are applied by the server/test run itself, not
// here. Connection details come from the process environment first, then
// apps/api/configs/.env — the dev database uses the DB_* namespace and the test
// database PG_TEST_*, the same namespaces the API server and internal/pgtest use.
package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/pgsetup"
)

func main() {
	pgsetup.LoadEnvFile()

	devOnly, testOnly := false, false
	var extra []string
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-h", "--help", "help":
			usage()
			return
		case "--dev-only":
			devOnly = true
		case "--test-only":
			testOnly = true
		default:
			if strings.HasPrefix(arg, "-") {
				fmt.Fprintln(os.Stderr, "unknown flag:", arg)
				usage()
				os.Exit(2)
			}
			extra = append(extra, arg)
		}
	}
	if devOnly && testOnly {
		fmt.Fprintln(os.Stderr, "--dev-only and --test-only are mutually exclusive")
		os.Exit(2)
	}

	// Each target carries its role: the role picks both the credentials and the
	// database name, so a machine that points dev and test at different servers
	// still works.
	type target struct {
		role string
		name string
	}
	var targets []target
	if !testOnly {
		targets = append(targets, target{role: pgsetup.RoleDev, name: pgsetup.DatabaseName(pgsetup.RoleDev)})
	}
	if !devOnly {
		targets = append(targets, target{role: pgsetup.RoleTest, name: pgsetup.DatabaseName(pgsetup.RoleTest)})
	}
	for _, name := range extra {
		targets = append(targets, target{role: pgsetup.RoleDev, name: name})
	}

	seen := map[string]bool{}
	unique := make([]target, 0, len(targets))
	for _, tgt := range targets {
		key := tgt.role + "/" + tgt.name
		if strings.TrimSpace(tgt.name) == "" || seen[key] {
			continue
		}
		seen[key] = true
		unique = append(unique, tgt)
	}
	if len(unique) == 0 {
		fmt.Fprintln(os.Stderr, "nothing to do: no database name configured (DB_NAME / PG_TEST_DB)")
		os.Exit(2)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	created, existing := 0, 0
	for _, role := range []string{pgsetup.RoleDev, pgsetup.RoleTest} {
		wanted := make([]string, 0, len(unique))
		for _, tgt := range unique {
			if tgt.role == role {
				wanted = append(wanted, tgt.name)
			}
		}
		if len(wanted) == 0 {
			continue
		}
		server := pgsetup.ServerForRole(role)
		conn, maintenanceDB, err := server.ConnectMaintenance(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s role: cannot reach PostgreSQL: %v\n", role, err)
			fmt.Fprintf(os.Stderr, "  check %s in the environment or apps/api/configs/.env\n", envHint(role))
			os.Exit(1)
		}
		_ = conn.Close(ctx)
		fmt.Printf("%s role: %s@%s via maintenance database %s\n", role, server.User, server.Host, maintenanceDB)
		for _, name := range wanted {
			made, err := server.Ensure(ctx, name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ensure %s (%s): %v\n", name, role, err)
				os.Exit(1)
			}
			if made {
				fmt.Printf("created %s (%s)\n", name, role)
				created++
				continue
			}
			fmt.Printf("exists  %s (%s)\n", name, role)
			existing++
		}
	}
	fmt.Printf("done: %d created, %d already present\n", created, existing)
	fmt.Println("migrations are applied by the server/test run itself (start the API or run the tests)")
}

func envHint(role string) string {
	if role == pgsetup.RoleTest {
		return "PG_TEST_HOST/PG_TEST_PORT/PG_TEST_USER/PG_TEST_PASSWORD"
	}
	return "DB_HOST/DB_PORT/DB_USER/DB_PASSWORD"
}

func usage() {
	fmt.Fprint(os.Stderr, `usage: dbsetup [--dev-only|--test-only] [name ...]

Creates the PostgreSQL databases this checkout needs (idempotent, never drops):
  --dev-only    only DB_NAME (the dev server database, DB_* credentials)
  --test-only   only PG_TEST_DB (the test database, PG_TEST_* credentials)
  name ...      also ensure these databases (dev role)

Connection details come from the process environment, then apps/api/configs/.env.
SQLite needs no provisioning: the file is created on first start.
`)
}
