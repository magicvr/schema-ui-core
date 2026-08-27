package store

import (
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// OpenOptions carries the configured dialect connection surface (R1 v1.4 §2 /
// §5). Only connection-facing fields are allowed; driver types never surface.
type OpenOptions struct {
	// Dialect selects the store implementation. Config normalizes an empty
	// value to kernel.DialectSQLite before calling Open.
	Dialect kernel.Dialect
	// Path is the SQLite file path, or for postgres the file-path-shaped value
	// whose filepath.Dir(path) derives the file storage root (not a SQL
	// connection). Defaults are applied by config.
	Path string
	// DSN is the postgres SQL connection string; required for postgres and
	// empty for sqlite (enforced by config validation).
	DSN string

	// Connection-pool / lifetime knobs. sqlite honors PoolMaxOpenConns
	// (default 4 for file DBs; in-memory DBs keep MaxOpenConns=1 because each
	// connection would be a separate database). postgres uses all four. Zero
	// leaves the driver default.
	PoolMaxOpenConns int
	PoolMaxIdleConns int
	ConnMaxLifetime  time.Duration
	// ConnectTimeout caps the initial connect + Ping probe (postgres).
	ConnectTimeout time.Duration
}
