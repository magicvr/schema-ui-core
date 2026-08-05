package testsupport

import (
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/compiled"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
)

// OpenStore opens a test database through the same compiled module catalog as
// the production composition root.
func OpenStore(path, adminUsername, adminPasswordHash string, seedAdmin bool) (*store.Store, error) {
	catalog, err := compiled.PersistenceCatalog()
	if err != nil {
		return nil, err
	}
	return store.OpenWithCatalog(path, adminUsername, adminPasswordHash, seedAdmin, catalog)
}
