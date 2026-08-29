package kernel

import (
	"context"
	"database/sql"
)

// Dialect identifies the active store dialect. Modules and jobs must never
// branch on Dialect; only the store implementation, composition, readyz and
// non-SQL diagnostics may read it (R1 v1.4 §1).
type Dialect string

const (
	// DialectSQLite is the embedded dev / mvp / fast-test default.
	DialectSQLite Dialect = "sqlite"
	// DialectPostgres is the production-authoritative dialect (VP-013).
	DialectPostgres Dialect = "postgres"
)

// ErrNoRows is the kernel "no rows" sentinel exposed to modules. Until the
// public surface stops importing database/sql (R4), it is an alias of
// sql.ErrNoRows so both errors.Is(err, kernel.ErrNoRows) and
// errors.Is(err, sql.ErrNoRows) hold (R1 v1.4 §3). Modules use
// errors.Is(err, kernel.ErrNoRows).
var ErrNoRows = sql.ErrNoRows

// Store is the kernel persistence port (R1 v1.4 §2). Domain code consumes the
// transaction boundary and never receives a driver handle. One Run call is one
// transaction; nested Run is forbidden (fail closed).
type Store interface {
	Dialect() Dialect
	Run(ctx context.Context, fn func(Tx) error) error
	Ping(ctx context.Context) error
	Close() error
	WasFresh() bool
	MarkSystemDataReady()
	SystemDataReady() error
}

// Tx is a dialect-neutral transaction handed to Run callbacks. It is only
// valid inside the current Run callback (R1 v1.4 §2). Placeholders are always
// written as '?' in module/migration SQL and rebound by the implementation.
type Tx interface {
	Exec(ctx context.Context, query string, args ...any) (Result, error)
	Query(ctx context.Context, query string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) Row
}

// Result mirrors database/sql.Result (RowsAffected only; LastInsertId is
// SQLite-specific and forbidden by R1 v1.4 §3).
type Result interface {
	RowsAffected() (int64, error)
}

// Rows mirrors database/sql.Rows.
type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Close() error
	Err() error
}

// Row mirrors database/sql.Row.
type Row interface {
	Scan(dest ...any) error
}
