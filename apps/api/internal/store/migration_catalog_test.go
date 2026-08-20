package store

import (
	"context"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authsessiondata "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession/systemdata"
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

// OpenSeeded exists only in store's test build. Production callers must supply
// the compiled catalog explicitly through OpenWithCatalog (the kernel port
// entry point is store.Open; R1 v1.4 sec.2).
func OpenSeeded(path, adminUsername, adminPasswordHash string, seedAdmin bool) (*Store, error) {
	st, err := OpenWithCatalog(path, MigrationCatalog())
	if err != nil || !seedAdmin {
		return st, err
	}
	if st.WasFresh() {
		if err := authsessiondata.Bootstrap(context.Background(), st, adminUsername, adminPasswordHash); err != nil {
			_ = st.Close()
			return nil, err
		}
	}
	permissions, navigation := testSystemDataContributions()
	if err := authsessiondata.Reconcile(context.Background(), st, permissions, navigation); err != nil {
		_ = st.Close()
		return nil, err
	}
	st.MarkSystemDataReady()
	return st, nil
}

func testSystemDataContributions() ([]kernel.PermissionContribution, []kernel.NavigationContribution) {
	permissions := []kernel.PermissionContribution{
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.users", Key: "users.read"}, Permission: "users.read", Resource: "users", Action: "read", PolicyID: authsessiondata.PolicyAdminEditorViewer, SystemDataVersion: 1},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.users", Key: "users.write"}, Permission: "users.write", Resource: "users", Action: "write", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: 1},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.roles", Key: "roles.read"}, Permission: "roles.read", Resource: "roles", Action: "read", PolicyID: authsessiondata.PolicyAdminEditorViewer, SystemDataVersion: 1},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.roles", Key: "roles.write"}, Permission: "roles.write", Resource: "roles", Action: "write", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: 1},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.roles", Key: "roles.assign"}, Permission: "roles.assign", Resource: "roles", Action: "assign", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: 1},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.settings", Key: "settings.read"}, Permission: "settings.read", Resource: "settings", Action: "read", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: 1},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.settings", Key: "settings.write"}, Permission: "settings.write", Resource: "settings", Action: "write", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: 1},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.activity", Key: "operations.read"}, Permission: "operations.read", Resource: "operations", Action: "read", PolicyID: authsessiondata.PolicyAdminEditor, SystemDataVersion: 1},
	}
	navigation := []kernel.NavigationContribution{
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.users", Key: "menu_users"}, NodeID: "menu_users", PageID: "users", Order: 1, Label: "Users", Visibility: authsessiondata.PolicyAdmin, Permission: "users.read", SystemDataVersion: 1},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.roles", Key: "menu_roles"}, NodeID: "menu_roles", PageID: "roles", Order: 2, Label: "Roles", Visibility: authsessiondata.PolicyAdmin, Permission: "roles.read", SystemDataVersion: 1},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.settings", Key: "menu_settings"}, NodeID: "menu_settings", PageID: "settings", Order: 1, Label: "Settings", Visibility: authsessiondata.PolicyAdmin, Permission: "settings.read", SystemDataVersion: 1},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.activity", Key: "menu_activity"}, NodeID: "menu_activity", PageID: "activity", Order: 2, Label: "Activity", Visibility: authsessiondata.PolicyAdminEditor, Permission: "operations.read", SystemDataVersion: 1},
	}
	return permissions, navigation
}
