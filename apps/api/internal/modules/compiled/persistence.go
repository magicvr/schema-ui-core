// Package compiled is the static candidate registry for module-owned assets.
// Runtime profile enablement does not filter persistence: every compiled
// provider contributes its immutable global migration history.
package compiled

import (
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	accountmigration "github.com/magicvr/schema-ui-core/apps/api/internal/modules/account/migration"
	authmigration "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession/migration"
	datadictionarymigration "github.com/magicvr/schema-ui-core/apps/api/internal/modules/datadictionary/migration"
	notificationsmigration "github.com/magicvr/schema-ui-core/apps/api/internal/modules/notifications/migration"
	historymigration "github.com/magicvr/schema-ui-core/apps/api/internal/modules/corepersistence/migration"
	operationlogmigration "github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog/migration"
	settingsmigration "github.com/magicvr/schema-ui-core/apps/api/internal/modules/settings/migration"
)

// PersistenceProviders returns all compiled migration owners. The order is not
// semantically significant; kernel.CollectPersistence validates and sorts the
// merged global catalog by version.
func PersistenceProviders() []kernel.Provider {
	return []kernel.Provider{
		accountmigration.Provider{},
		datadictionarymigration.Provider{},
		notificationsmigration.Provider{},
		authmigration.Provider{},
		historymigration.Provider{},
		operationlogmigration.Provider{},
		settingsmigration.Provider{},
	}
}

func PersistenceCatalog() ([]kernel.MigrationContribution, error) {
	return kernel.CollectPersistence(PersistenceProviders())
}