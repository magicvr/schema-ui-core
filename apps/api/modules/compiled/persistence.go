// Package compiled is the static candidate registry for module-owned assets.
// Runtime profile enablement does not filter persistence: every compiled
// provider contributes its immutable global migration history.
package compiled

import (
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	accountmigration "github.com/magicvr/schema-ui-core/apps/api/modules/account/migration"
	authmigration "github.com/magicvr/schema-ui-core/apps/api/modules/authsession/migration"
	historymigration "github.com/magicvr/schema-ui-core/apps/api/modules/corepersistence/migration"
	datadictionarymigration "github.com/magicvr/schema-ui-core/apps/api/modules/datadictionary/migration"
	datapermissionmigration "github.com/magicvr/schema-ui-core/apps/api/modules/datapermission/migration"
	jobsmigration "github.com/magicvr/schema-ui-core/apps/api/modules/jobs/migration"
	logincaptchamigration "github.com/magicvr/schema-ui-core/apps/api/modules/logincaptcha/migration"
	mfamigration "github.com/magicvr/schema-ui-core/apps/api/modules/mfa/migration"
	notificationsmigration "github.com/magicvr/schema-ui-core/apps/api/modules/notifications/migration"
	operationlogmigration "github.com/magicvr/schema-ui-core/apps/api/modules/operationlog/migration"
	recyclebinmigration "github.com/magicvr/schema-ui-core/apps/api/modules/recyclebin/migration"
	scheduledtasksmigration "github.com/magicvr/schema-ui-core/apps/api/modules/scheduledtasks/migration"
	settingsmigration "github.com/magicvr/schema-ui-core/apps/api/modules/settings/migration"
	walletmigration "github.com/magicvr/schema-ui-core/apps/api/modules/wallet/migration"
)

// PersistenceProviders returns all compiled migration owners. The order is not
// semantically significant; kernel.CollectPersistence validates and sorts the
// merged global catalog by version.
func PersistenceProviders() []kernel.Provider {
	return []kernel.Provider{
		accountmigration.Provider{},
		datadictionarymigration.Provider{},
		logincaptchamigration.Provider{},
		datapermissionmigration.Provider{},
		mfamigration.Provider{},
		recyclebinmigration.Provider{},
		walletmigration.Provider{},
		scheduledtasksmigration.Provider{},
		notificationsmigration.Provider{},
		authmigration.Provider{},
		historymigration.Provider{},
		operationlogmigration.Provider{},
		settingsmigration.Provider{},
		jobsmigration.Provider{},
	}
}

func PersistenceCatalog() ([]kernel.MigrationContribution, error) {
	return kernel.CollectPersistence(PersistenceProviders())
}
