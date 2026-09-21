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
// Since W35 (GOAL-047) the maintenance connection falls back to the standard
// maintenance databases when the configured DB_NAME does not exist yet, so this
// tool can bootstrap the very first database on a fresh or reset server.
//
// Usage:
//
//	go run ./cmd/e2e-pgset create <name>   # CREATE DATABASE name (idempotent)
//	go run ./cmd/e2e-pgset drop <name>     # DROP DATABASE name (no active conns)
//	go run ./cmd/e2e-pgset verify <name>   # exit 0 once schema_migrations exists
//	go run ./cmd/e2e-pgset list            # list schema_ui_e2e_* databases (leftovers)
//
// For the dev/test databases use `go run ./cmd/dbsetup` (or `dev.cmd init-db`).
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/pgsetup"
)

func main() {
	pgsetup.LoadEnvFile()
	if len(os.Args) == 2 && os.Args[1] == "list" {
		listExisting()
		return
	}
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: e2e-pgset create|drop|verify <name> | list")
		os.Exit(2)
	}
	action, name := os.Args[1], os.Args[2]
	if !pgsetup.ValidName(name) {
		fmt.Fprintln(os.Stderr, "invalid database name:", name)
		os.Exit(2)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server := pgsetup.ServerForRole(pgsetup.RoleDev)

	switch action {
	case "create":
		made, err := server.Ensure(ctx, name)
		if err != nil {
			fmt.Fprintln(os.Stderr, "create:", err)
			os.Exit(1)
		}
		if made {
			fmt.Println("created", name)
			return
		}
		fmt.Println("exists", name)
	case "drop":
		dropped, err := server.Drop(ctx, name)
		if err != nil {
			fmt.Fprintln(os.Stderr, "drop:", err)
			os.Exit(1)
		}
		if dropped {
			fmt.Println("dropped", name)
			return
		}
		fmt.Println("absent", name)
	case "verify":
		if err := server.VerifyMigrated(ctx, name); err != nil {
			fmt.Fprintln(os.Stderr, "verify:", err)
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
	names, err := pgsetup.ServerForRole(pgsetup.RoleDev).ListE2EDatabases(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "list:", err)
		os.Exit(1)
	}
	for _, name := range names {
		fmt.Println(name)
	}
	fmt.Fprintf(os.Stderr, "%d schema_ui_e2e_* database(s)\n", len(names))
}
