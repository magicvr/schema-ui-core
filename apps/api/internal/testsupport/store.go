package testsupport

import (
	"context"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authsessiondata "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession/systemdata"
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
	st, err := store.OpenWithCatalog(path, catalog)
	if err != nil || !seedAdmin {
		return st, err
	}
	needsBootstrap, err := authsessiondata.NeedsBootstrap(context.Background(), st)
	if err != nil {
		_ = st.Close()
		return nil, err
	}
	if needsBootstrap {
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
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.users", Key: "users.read"}, Permission: "users.read", Resource: "users", Action: "read", PolicyID: authsessiondata.PolicyAdminEditorViewer, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.users", Key: "users.write"}, Permission: "users.write", Resource: "users", Action: "write", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.roles", Key: "roles.read"}, Permission: "roles.read", Resource: "roles", Action: "read", PolicyID: authsessiondata.PolicyAdminEditorViewer, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.roles", Key: "roles.write"}, Permission: "roles.write", Resource: "roles", Action: "write", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.roles", Key: "roles.assign"}, Permission: "roles.assign", Resource: "roles", Action: "assign", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.settings", Key: "settings.read"}, Permission: "settings.read", Resource: "settings", Action: "read", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.settings", Key: "settings.write"}, Permission: "settings.write", Resource: "settings", Action: "write", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.activity", Key: "operations.read"}, Permission: "operations.read", Resource: "operations", Action: "read", PolicyID: authsessiondata.PolicyAdminEditor, SystemDataVersion: authsessiondata.SystemDataVersion},
		// W4 P0-2: files.write is a central shared-capability permission
		// (upload endpoint is centrally registered), admin-only by default.
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "core.server-registration", Key: "files.write"}, Permission: "files.write", Resource: "files", Action: "write", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		// F-03 (GOAL-005): admin.account enable/disable keys, admin-only.
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.account", Key: "users.enable"}, Permission: "users.enable", Resource: "users", Action: "enable", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.account", Key: "users.disable"}, Permission: "users.disable", Resource: "users", Action: "disable", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		// F-02 (GOAL-004): admin.data-transfer keys.
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.data-transfer", Key: "data.export"}, Permission: "data.export", Resource: "data", Action: "export", PolicyID: authsessiondata.PolicyAdminEditor, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.data-transfer", Key: "data.import"}, Permission: "data.import", Resource: "data", Action: "import", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
	}
	navigation := []kernel.NavigationContribution{
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.users", Key: "menu_users"}, NodeID: "menu_users", PageID: "users", Order: 1, Label: "Users", Visibility: authsessiondata.PolicyAdmin, Permission: "users.read", SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.roles", Key: "menu_roles"}, NodeID: "menu_roles", PageID: "roles", Order: 2, Label: "Roles", Visibility: authsessiondata.PolicyAdmin, Permission: "roles.read", SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.settings", Key: "menu_settings"}, NodeID: "menu_settings", PageID: "settings", Order: 1, Label: "Settings", Visibility: authsessiondata.PolicyAdmin, Permission: "settings.read", SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.activity", Key: "menu_activity"}, NodeID: "menu_activity", PageID: "activity", Order: 2, Label: "Activity", Visibility: authsessiondata.PolicyAdminEditor, Permission: "operations.read", SystemDataVersion: authsessiondata.SystemDataVersion},
		// F-03 (GOAL-005): self-service account page for every standard role.
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.account", Key: "menu_account"}, NodeID: "menu_account", PageID: "account", Order: 1, Label: "Account", Visibility: authsessiondata.PolicyAdminEditorViewer, Permission: "", SystemDataVersion: authsessiondata.SystemDataVersion},
		// F-01 (GOAL-003): dashboard page for every standard role.
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.dashboard", Key: "menu_dashboard"}, NodeID: "menu_dashboard", PageID: "dashboard", Order: 0, Label: "Dashboard", Visibility: authsessiondata.PolicyAdminEditorViewer, Permission: "", SystemDataVersion: authsessiondata.SystemDataVersion},
		// F-04 (GOAL-006): notifications page for every standard role.
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.notifications", Key: "menu_notifications"}, NodeID: "menu_notifications", PageID: "notifications", Order: 2, Label: "Notifications", Visibility: authsessiondata.PolicyAdminEditorViewer, Permission: "", SystemDataVersion: authsessiondata.SystemDataVersion},
	}
	return permissions, navigation
}