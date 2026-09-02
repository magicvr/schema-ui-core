package testsupport

import (
	"context"
	"database/sql"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	authsessiondata "github.com/magicvr/schema-ui-core/apps/api/modules/authsession/systemdata"
	"github.com/magicvr/schema-ui-core/apps/api/modules/compiled"
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
	// W16-F01: the production bootstrap seeds must_change_password=1. Test
	// environments use the seeded admin as a normal pre-change account, so clear
	// the flag after bootstrap to keep existing test contracts stable.
	if err := st.WithTx(context.Background(), func(tx *sql.Tx) error {
		_, err := tx.Exec(`UPDATE users SET must_change_password = 0`)
		return err
	}); err != nil {
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
		// workspace-019 R3 (GOAL-004): invitation management key, admin-only.
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.users", Key: "users.invite"}, Permission: "users.invite", Resource: "users", Action: "invite", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.roles", Key: "roles.read"}, Permission: "roles.read", Resource: "roles", Action: "read", PolicyID: authsessiondata.PolicyAdminEditorViewer, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.roles", Key: "roles.write"}, Permission: "roles.write", Resource: "roles", Action: "write", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.roles", Key: "roles.assign"}, Permission: "roles.assign", Resource: "roles", Action: "assign", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.settings", Key: "settings.read"}, Permission: "settings.read", Resource: "settings", Action: "read", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.settings", Key: "settings.write"}, Permission: "settings.write", Resource: "settings", Action: "write", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.activity", Key: "operations.read"}, Permission: "operations.read", Resource: "operations", Action: "read", PolicyID: authsessiondata.PolicyAdminEditor, SystemDataVersion: authsessiondata.SystemDataVersion},
		// W4 P0-2: files.write is a central shared-capability permission
		// (upload endpoint is centrally registered), admin-only by default.
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "core.server-registration", Key: "files.write"}, Permission: "files.write", Resource: "files", Action: "write", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "core.auth-session", Key: "service-credentials.read"}, Permission: "service-credentials.read", Resource: "service-credentials", Action: "read", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "core.auth-session", Key: "service-credentials.write"}, Permission: "service-credentials.write", Resource: "service-credentials", Action: "write", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		// F-03 (GOAL-005): admin.account enable/disable keys, admin-only.
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.account", Key: "users.enable"}, Permission: "users.enable", Resource: "users", Action: "enable", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.account", Key: "users.disable"}, Permission: "users.disable", Resource: "users", Action: "disable", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		// S-10 (GOAL-017): admin.mfa admin-reset key.
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.mfa", Key: "users.mfa-reset"}, Permission: "users.mfa-reset", Resource: "users", Action: "mfa-reset", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		// S-09 (GOAL-016): admin.data-permission keys.
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.data-permission", Key: "data-permission.read"}, Permission: "data-permission.read", Resource: "data-permission", Action: "read", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.data-permission", Key: "data-permission.write"}, Permission: "data-permission.write", Resource: "data-permission", Action: "write", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		// S-14 (GOAL-019): admin.wallet keys.
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.wallet", Key: "wallet.read"}, Permission: "wallet.read", Resource: "wallet", Action: "read", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.wallet", Key: "wallet.write"}, Permission: "wallet.write", Resource: "wallet", Action: "write", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.wallet", Key: "wallet.adjust"}, Permission: "wallet.adjust", Resource: "wallet", Action: "adjust", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.wallet", Key: "wallet.voucher.issue"}, Permission: "wallet.voucher.issue", Resource: "wallet", Action: "voucher.issue", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		// F-02 (GOAL-004): admin.data-transfer keys.
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.data-transfer", Key: "data.export"}, Permission: "data.export", Resource: "data", Action: "export", PolicyID: authsessiondata.PolicyAdminEditor, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.data-transfer", Key: "data.import"}, Permission: "data.import", Resource: "data", Action: "import", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		// S-02 (GOAL-007): admin.file-library keys, admin-only.
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.file-library", Key: "files.read"}, Permission: "files.read", Resource: "files", Action: "read", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.file-library", Key: "files.delete"}, Permission: "files.delete", Resource: "files", Action: "delete", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		// S-01 (GOAL-008): admin.data-dictionary keys, admin-only.
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.data-dictionary", Key: "dictionary.read"}, Permission: "dictionary.read", Resource: "dictionary", Action: "read", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.data-dictionary", Key: "dictionary.write"}, Permission: "dictionary.write", Resource: "dictionary", Action: "write", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		// S-03 (GOAL-009): admin.system-monitoring key, admin-only.
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.system-monitoring", Key: "monitoring.read"}, Permission: "monitoring.read", Resource: "system-monitoring", Action: "read", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		// S-04 (GOAL-010): admin.scheduled-tasks keys, admin-only.
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.scheduled-tasks", Key: "tasks.read"}, Permission: "tasks.read", Resource: "scheduled-tasks", Action: "read", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.scheduled-tasks", Key: "tasks.write"}, Permission: "tasks.write", Resource: "scheduled-tasks", Action: "write", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		// S-11 (GOAL-011): admin.login-captcha keys, admin-only.
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.login-captcha", Key: "captcha.read"}, Permission: "captcha.read", Resource: "captcha", Action: "read", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.login-captcha", Key: "captcha.write"}, Permission: "captcha.write", Resource: "captcha", Action: "write", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		// S-12 (GOAL-012): admin.recycle-bin keys, admin-only.
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.recycle-bin", Key: "recycle.read"}, Permission: "recycle.read", Resource: "recycle-bin", Action: "read", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.recycle-bin", Key: "recycle.write"}, Permission: "recycle.write", Resource: "recycle-bin", Action: "write", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
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
		// S-02 (GOAL-007): file library page (admin-only management surface).
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.file-library", Key: "menu_files"}, NodeID: "menu_files", PageID: "file-library", Order: 3, Label: "File library", Visibility: authsessiondata.PolicyAdmin, Permission: "files.read", SystemDataVersion: authsessiondata.SystemDataVersion},
		// S-01 (GOAL-008): data dictionary page (admin-only management surface).
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.data-dictionary", Key: "menu_dictionary"}, NodeID: "menu_dictionary", PageID: "data-dictionary", Order: 4, Label: "Data dictionary", Visibility: authsessiondata.PolicyAdmin, Permission: "dictionary.read", SystemDataVersion: authsessiondata.SystemDataVersion},
		// S-03 (GOAL-009): system monitoring page (admin-only read surface).
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.system-monitoring", Key: "menu_monitoring"}, NodeID: "menu_monitoring", PageID: "system-monitoring", Order: 5, Label: "System monitoring", Visibility: authsessiondata.PolicyAdmin, Permission: "monitoring.read", SystemDataVersion: authsessiondata.SystemDataVersion},
		// S-04 (GOAL-010): scheduled tasks page (admin-only management surface).
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.scheduled-tasks", Key: "menu_scheduled_tasks"}, NodeID: "menu_scheduled_tasks", PageID: "scheduled-tasks", Order: 6, Label: "Scheduled tasks", Visibility: authsessiondata.PolicyAdmin, Permission: "tasks.read", SystemDataVersion: authsessiondata.SystemDataVersion},
		// S-12 (GOAL-012): recycle bin page (admin-only surface).
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.recycle-bin", Key: "menu_recycle_bin"}, NodeID: "menu_recycle_bin", PageID: "recycle-bin", Order: 8, Label: "Recycle bin", Visibility: authsessiondata.PolicyAdmin, Permission: "recycle.read", SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.data-permission", Key: "menu_data_permission"}, NodeID: "menu_data_permission", PageID: "data-permission", Order: 9, Label: "Data permission", Visibility: authsessiondata.PolicyAdmin, Permission: "data-permission.read", SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.wallet", Key: "menu_wallet"}, NodeID: "menu_wallet", PageID: "wallet", Order: 3, Label: "Wallet", Visibility: authsessiondata.PolicyAdmin, Permission: "wallet.read", SystemDataVersion: authsessiondata.SystemDataVersion},
	}
	return permissions, navigation
}
