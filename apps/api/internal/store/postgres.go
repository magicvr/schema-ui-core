package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"

	_ "github.com/jackc/pgx/v5/stdlib" // registers driver name "pgx"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// postgres implements kernel.Store for the postgres dialect. R2 delivers
// connect + Ping + WasFresh (probe open); the compiled catalog is NOT applied
// here — dual-dialect apply is R3.
type postgres struct {
	db              *sql.DB
	fresh           bool
	systemDataReady atomic.Bool
}

// openPostgres opens a probe connection: DSN connect, Ping, WasFresh. A
// non-empty catalog fails closed because the R2 catalog still contains
// SQLite-specific SQL (R1 v1.4 §2).
func openPostgres(ctx context.Context, opts OpenOptions, catalog []kernel.MigrationContribution) (*postgres, error) {
	if opts.DSN == "" {
		return nil, errors.New("store: postgres requires a non-empty DSN")
	}
	if len(catalog) > 0 {
		return nil, errors.New("store: postgres (R2) cannot apply the compiled catalog — it contains SQLite-specific SQL; dual-dialect apply lands in R3, refusing to half-execute")
	}
	db, err := sql.Open("pgx", opts.DSN)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}
	applyPostgresPoolOptions(db, opts)

	probeCtx := ctx
	cancel := func() {}
	if opts.ConnectTimeout > 0 {
		probeCtx, cancel = context.WithTimeout(ctx, opts.ConnectTimeout)
	}
	defer cancel()
	if err := db.PingContext(probeCtx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("postgres ping: %w", err)
	}
	fresh, err := postgresWasFresh(probeCtx, db)
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("postgres WasFresh: %w", err)
	}
	return &postgres{db: db, fresh: fresh}, nil
}

func applyPostgresPoolOptions(db *sql.DB, opts OpenOptions) {
	if opts.PoolMaxOpenConns > 0 {
		db.SetMaxOpenConns(opts.PoolMaxOpenConns)
	}
	if opts.PoolMaxIdleConns > 0 {
		db.SetMaxIdleConns(opts.PoolMaxIdleConns)
	}
	if opts.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(opts.ConnMaxLifetime)
	}
}

func (p *postgres) Dialect() kernel.Dialect { return kernel.DialectPostgres }

// Run executes fn in one transaction; placeholders are rebound '?' -> $n.
// Nested Run is detected per-callback via the ctx passed to fn (R1 v1.4
// A-008 F-004).
func (p *postgres) Run(ctx context.Context, fn func(kernel.Tx) error) error {
	if err := enterRun(); err != nil {
		return err
	}
	defer leaveRun()
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// R1 v1.4 §2: a panicking fn rolls the transaction back and re-panics, so
	// the tx handle is never abandoned and the caller still observes the panic.
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		}
	}()
	if err := fn(pgTx{tx: tx}); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return err
	}
	return nil
}

func (p *postgres) Ping(ctx context.Context) error { return p.db.PingContext(ctx) }
func (p *postgres) Close() error                   { return p.db.Close() }
func (p *postgres) WasFresh() bool                 { return p.fresh }
func (p *postgres) MarkSystemDataReady()           { p.systemDataReady.Store(true) }

func (p *postgres) SystemDataReady() error {
	if !p.systemDataReady.Load() {
		return errors.New("store: system-data reconciliation is not ready")
	}
	return nil
}

// pgTx adapts *sql.Tx to kernel.Tx with postgres placeholder rebinding.
type pgTx struct{ tx *sql.Tx }

func (t pgTx) Exec(ctx context.Context, query string, args ...any) (kernel.Result, error) {
	return t.tx.ExecContext(ctx, rebindPostgres(query), args...)
}

func (t pgTx) Query(ctx context.Context, query string, args ...any) (kernel.Rows, error) {
	return t.tx.QueryContext(ctx, rebindPostgres(query), args...)
}

func (t pgTx) QueryRow(ctx context.Context, query string, args ...any) kernel.Row {
	return t.tx.QueryRowContext(ctx, rebindPostgres(query), args...)
}

// postgresWasFresh reports whether the database has zero user base tables in
// the first existing USER (non-system) schema of the server-resolved
// search_path (R1 v1.4 §2). A search_path with no existing user schema is
// fresh; system schemas (pg_catalog / information_schema / pg_toast / pg_temp_*)
// never count as the probed user schema.
func postgresWasFresh(ctx context.Context, db *sql.DB) (bool, error) {
	var raw string
	if err := db.QueryRowContext(ctx, `SHOW search_path`).Scan(&raw); err != nil {
		return false, fmt.Errorf("read search_path: %w", err)
	}
	schema, err := firstExistingSchema(ctx, db, raw)
	if err != nil {
		return false, err
	}
	if schema == "" {
		return true, nil
	}
	var count int
	err = db.QueryRowContext(ctx, `
		SELECT count(*)
		FROM pg_catalog.pg_class c
		JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
		WHERE n.nspname = $1
		  AND c.relkind = 'r'
		  AND n.nspname NOT IN ('pg_catalog', 'information_schema', 'pg_toast')
	`, schema).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("count base tables in %q: %w", schema, err)
	}
	return count == 0, nil
}

// searchPathCandidates parses the SHOW search_path value into ordered schema
// candidates, honoring standard double-quoted identifiers and lowercasing
// unquoted ones (like PostgreSQL does). "$user" stays literal for callers to
// resolve against current_user.
func searchPathCandidates(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		p := strings.TrimSpace(part)
		if p == "" {
			continue
		}
		if strings.HasPrefix(p, `"`) {
			if len(p) >= 2 && strings.HasSuffix(p, `"`) {
				p = strings.TrimSuffix(strings.TrimPrefix(p, `"`), `"`)
				p = strings.ReplaceAll(p, `""`, `"`)
			}
		} else {
			p = strings.ToLower(p)
		}
		out = append(out, p)
	}
	return out
}

// firstExistingSchema walks the ordered search_path candidates and returns the
// first schema that actually exists (postgres skips non-existent entries) AND
// is a user schema — system schemas (pg_catalog / information_schema /
// pg_toast / pg_temp_*) are skipped like postgres treats them as non-resolvable
// for user objects. "$user" is resolved against current_user when it appears.
// Empty result = no existing user schema = a fresh database.
func firstExistingSchema(ctx context.Context, db *sql.DB, raw string) (string, error) {
	resolvedUser := ""
	for _, cand := range searchPathCandidates(raw) {
		name := cand
		if name == "$user" {
			if resolvedUser == "" {
				if err := db.QueryRowContext(ctx, `SELECT current_user`).Scan(&resolvedUser); err != nil {
					return "", fmt.Errorf("resolve $user: %w", err)
				}
			}
			name = resolvedUser
		}
		if isSystemSchema(name) {
			continue
		}
		var exists bool
		if err := db.QueryRowContext(ctx,
			`SELECT EXISTS(SELECT 1 FROM information_schema.schemata WHERE schema_name = $1)`,
			name,
		).Scan(&exists); err != nil {
			return "", fmt.Errorf("probe schema %q: %w", name, err)
		}
		if exists {
			return name, nil
		}
	}
	return "", nil
}

// isSystemSchema reports whether name is a PostgreSQL system schema that must
// never be treated as the probed user schema for WasFresh.
func isSystemSchema(name string) bool {
	switch name {
	case "pg_catalog", "information_schema", "pg_toast":
		return true
	}
	return strings.HasPrefix(name, "pg_temp_") || strings.HasPrefix(name, "pg_toast_temp_")
}
