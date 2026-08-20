package store

import (
	"context"
	"fmt"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// Open opens a store for the configured dialect (R1 v1.4 §2 terminal shape).
//
//   - sqlite: applies the supplied compiled catalog (default dev path
//     unchanged; this is also the fast-test path).
//   - postgres (R3 dual-dialect ledger): connect + Ping + WasFresh, then,
//     when a non-empty catalog is supplied, applies it through the postgres
//     migrate runner (fresh bootstrap / incremental / ledger + checksum). A
//     nil/empty catalog is the probe-open path (driver, pool, readyz).
func Open(ctx context.Context, opts OpenOptions, catalog []kernel.MigrationContribution) (kernel.Store, error) {
	switch opts.Dialect {
	case kernel.DialectPostgres:
		return openPostgres(ctx, opts, catalog)
	case kernel.DialectSQLite:
		normalized, err := normalizeCatalog(catalog)
		if err != nil {
			return nil, err
		}
		return open(opts.Path, normalized)
	default:
		return nil, fmt.Errorf("store: unknown dialect %q", opts.Dialect)
	}
}
