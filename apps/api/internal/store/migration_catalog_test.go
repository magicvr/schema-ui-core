package store

import (
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	compiledmodules "github.com/magicvr/schema-ui-core/apps/api/internal/modules/compiled"
)

// compiledMigrations is test-only compatibility for focused runner tests. The
// production store has no built-in migration registry.
var compiledMigrations = mustCompiledMigrationCatalog()

func mustCompiledMigrationCatalog() []kernel.MigrationContribution {
	catalog, err := compiledmodules.PersistenceCatalog()
	if err != nil {
		panic(err)
	}
	return catalog
}

func MigrationCatalog() []kernel.MigrationContribution {
	return append([]kernel.MigrationContribution(nil), compiledMigrations...)
}

// Open exists only in store's test build. Production callers must supply the
// compiled catalog explicitly through OpenWithCatalog.
func Open(path, adminUsername, adminPasswordHash string, seedAdmin bool) (*Store, error) {
	return OpenWithCatalog(path, adminUsername, adminPasswordHash, seedAdmin, MigrationCatalog())
}
