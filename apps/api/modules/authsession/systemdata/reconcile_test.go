package systemdata

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	authsession "github.com/magicvr/schema-ui-core/apps/api/modules/authsession"
	"github.com/magicvr/schema-ui-core/apps/api/modules/compiled"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
)

func openTestStore(t *testing.T) *store.Store {
	t.Helper()
	catalog, err := compiled.PersistenceCatalog()
	if err != nil {
		t.Fatalf("catalog: %v", err)
	}
	st, err := store.OpenWithCatalog(filepath.Join(t.TempDir(), "system-data.db"), catalog)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return st
}

func TestValidateInputsRejectsWellFormedUnknownPolicies(t *testing.T) {
	permission := kernel.PermissionContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.sample", Key: "sample.read"},
		Permission:           "sample.read", Resource: "sample", Action: "read", PolicyID: "system.custom", SystemDataVersion: 1,
	}
	if err := validateInputs([]kernel.PermissionContribution{permission}, nil); err == nil || !strings.Contains(err.Error(), `unknown policy "system.custom"`) {
		t.Fatalf("permission policy err = %v", err)
	}

	navigation := kernel.NavigationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.sample", Key: "menu_sample"},
		NodeID:               "menu_sample", PageID: "sample", Label: "Sample", Visibility: "system.custom", SystemDataVersion: 1,
	}
	if err := validateInputs(nil, []kernel.NavigationContribution{navigation}); err == nil || !strings.Contains(err.Error(), `unknown visibility policy "system.custom"`) {
		t.Fatalf("visibility policy err = %v", err)
	}
}

func sampleSystemData() ([]kernel.PermissionContribution, []kernel.NavigationContribution) {
	permissions := []kernel.PermissionContribution{
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.users", Key: "users.read"},
			Permission:           "users.read", Resource: "users", Action: "read",
			PolicyID: PolicyAdminEditorViewer, SystemDataVersion: 1,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.users", Key: "users.write"},
			Permission:           "users.write", Resource: "users", Action: "write",
			PolicyID: PolicyAdmin, SystemDataVersion: 1,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.settings", Key: "settings.read"},
			Permission:           "settings.read", Resource: "settings", Action: "read",
			PolicyID: PolicyAdmin, SystemDataVersion: 1,
		},
	}
	navigation := []kernel.NavigationContribution{
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.users", Key: "menu_users"},
			NodeID:               "menu_users", PageID: "users", Order: 1, Label: "Users",
			Visibility: PolicyAdmin, Permission: "users.read", SystemDataVersion: 1,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.settings", Key: "menu_settings"},
			NodeID:               "menu_settings", PageID: "settings", Order: 1, Label: "Settings",
			Visibility: PolicyAdmin, Permission: "settings.read", SystemDataVersion: 1,
		},
	}
	return permissions, navigation
}

func queryInt(t *testing.T, st *store.Store, query string, args ...any) int {
	t.Helper()
	var value int
	if err := st.WithTx(context.Background(), func(tx *sql.Tx) error {
		return tx.QueryRow(query, args...).Scan(&value)
	}); err != nil {
		t.Fatalf("query %q: %v", query, err)
	}
	return value
}

func TestBootstrapAndReconcileAreSeparateAndIdempotent(t *testing.T) {
	st := openTestStore(t)
	repository := authsession.NewRepository(st)
	if !st.WasFresh() {
		t.Fatal("new database must be marked fresh")
	}
	if _, err := repository.UserByUsername("admin"); !errors.Is(err, authsession.ErrNotFound) {
		t.Fatalf("migration-only open user = %v, want no bootstrap admin", err)
	}
	if err := Bootstrap(context.Background(), st, "admin", "hash-v1"); err != nil {
		t.Fatalf("bootstrap: %v", err)
	}
	permissions, navigation := sampleSystemData()
	if err := Reconcile(context.Background(), st, permissions, navigation); err != nil {
		t.Fatalf("reconcile: %v", err)
	}
	if got := queryInt(t, st, `SELECT COUNT(*) FROM system_data_reconcile`); got != 1+len(permissions)+len(navigation) {
		t.Fatalf("ledger entries = %d, want %d", got, 1+len(permissions)+len(navigation))
	}
	if got := queryInt(t, st, `SELECT COUNT(*) FROM roles WHERE system = 1`); got != 3 {
		t.Fatalf("system roles = %d, want 3", got)
	}
	if err := st.WithTx(context.Background(), func(tx *sql.Tx) error {
		if _, err := tx.Exec(`UPDATE system_data_reconcile SET applied_at = 123 WHERE module_id = 'admin.users' AND kind = 'authorization' AND contribution_key = 'users.read'`); err != nil {
			return err
		}
		_, err := tx.Exec(`UPDATE roles SET updated_at = 123 WHERE id = 'role-admin'`)
		return err
	}); err != nil {
		t.Fatal(err)
	}
	if err := Reconcile(context.Background(), st, permissions, navigation); err != nil {
		t.Fatalf("second reconcile: %v", err)
	}
	if got := queryInt(t, st, `SELECT applied_at FROM system_data_reconcile WHERE module_id = 'admin.users' AND kind = 'authorization' AND contribution_key = 'users.read'`); got != 123 {
		t.Fatalf("unchanged reconcile rewrote ledger applied_at = %d, want 123", got)
	}
	if got := queryInt(t, st, `SELECT updated_at FROM roles WHERE id = 'role-admin'`); got != 123 {
		t.Fatalf("unchanged reconcile rewrote system role updated_at = %d, want 123", got)
	}
	if got := queryInt(t, st, `SELECT COUNT(*) FROM role_permissions`); got != 5 {
		t.Fatalf("role permission grants = %d, want 5", got)
	}
	if got := queryInt(t, st, `SELECT COUNT(*) FROM role_menu_items`); got != 2 {
		t.Fatalf("role menu grants = %d, want 2", got)
	}
	if err := Bootstrap(context.Background(), st, "admin", "hash-v2"); err != nil {
		t.Fatalf("repeat bootstrap: %v", err)
	}
	admin, err := repository.UserByUsername("admin")
	if err != nil {
		t.Fatal(err)
	}
	if admin.PasswordHash != "hash-v1" {
		t.Fatalf("bootstrap overwrote password: %q", admin.PasswordHash)
	}
}

func TestReconcilePreservesUserFieldsAndDisabledProfileData(t *testing.T) {
	st := openTestStore(t)
	if err := Bootstrap(context.Background(), st, "admin", "hash"); err != nil {
		t.Fatal(err)
	}
	permissions, navigation := sampleSystemData()
	if err := Reconcile(context.Background(), st, permissions, navigation); err != nil {
		t.Fatal(err)
	}
	if err := st.WithTx(context.Background(), func(tx *sql.Tx) error {
		statements := []string{
			`UPDATE roles SET name = 'Operations Admin', updated_at = 123 WHERE id = 'role-admin'`,
			`UPDATE permissions SET description = 'operator description' WHERE id = 'perm-users-read'`,
			`UPDATE menu_items SET sort_order = 99, enabled = 0 WHERE id = 'menu-users'`,
			`INSERT INTO permissions (id, key, description, created_at, updated_at) VALUES ('perm-custom','custom.read','custom',1,1)`,
		}
		for _, statement := range statements {
			if _, err := tx.Exec(statement); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	// Reconcile an MVP-like subset. Settings rows remain because profile
	// disablement never deletes system or user data.
	if err := Reconcile(context.Background(), st, permissions[:2], navigation[:1]); err != nil {
		t.Fatal(err)
	}
	if err := st.WithTx(context.Background(), func(tx *sql.Tx) error {
		var roleName, description string
		var system, roleUpdatedAt, order, enabled int
		if err := tx.QueryRow(`SELECT name, system, updated_at FROM roles WHERE id = 'role-admin'`).Scan(&roleName, &system, &roleUpdatedAt); err != nil {
			return err
		}
		if err := tx.QueryRow(`SELECT description FROM permissions WHERE id = 'perm-users-read'`).Scan(&description); err != nil {
			return err
		}
		if err := tx.QueryRow(`SELECT sort_order, enabled FROM menu_items WHERE id = 'menu-users'`).Scan(&order, &enabled); err != nil {
			return err
		}
		if roleName != "Operations Admin" || system != 1 || roleUpdatedAt != 123 || description != "operator description" || order != 99 || enabled != 0 {
			return errors.New("reconcile overwrote operator-owned fields")
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"settings.read", "custom.read"} {
		if got := queryInt(t, st, `SELECT COUNT(*) FROM permissions WHERE key = ?`, key); got != 1 {
			t.Fatalf("permission %s retained count = %d, want 1", key, got)
		}
	}
	if got := queryInt(t, st, `SELECT COUNT(*) FROM menu_items WHERE feature_key = 'menu_settings'`); got != 1 {
		t.Fatalf("disabled-profile settings menu retained count = %d, want 1", got)
	}
}

func TestReconcileDetectsDriftAndAppliesVersionedGrantPolicy(t *testing.T) {
	st := openTestStore(t)
	repository := authsession.NewRepository(st)
	if err := Bootstrap(context.Background(), st, "admin", "hash"); err != nil {
		t.Fatal(err)
	}
	permissions, _ := sampleSystemData()
	permission := permissions[0]
	if err := Reconcile(context.Background(), st, []kernel.PermissionContribution{permission}, nil); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	if err := repository.CreateUser(authsession.User{
		ID: "user-viewer", Username: "viewer", Name: "Viewer", Roles: []string{"viewer"},
		PasswordHash: "hash", CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatal(err)
	}
	viewerPermissions, err := repository.PermissionsForUser("user-viewer")
	if err != nil || !reflect.DeepEqual(viewerPermissions, []string{"users.read"}) {
		t.Fatalf("viewer permissions = %v, err %v", viewerPermissions, err)
	}

	drifted := permission
	drifted.Resource = "accounts"
	if err := Reconcile(context.Background(), st, []kernel.PermissionContribution{drifted}, nil); err == nil || !strings.Contains(err.Error(), "checksum drift") {
		t.Fatalf("drift error = %v, want checksum drift", err)
	}

	upgraded := permission
	upgraded.SystemDataVersion = 2
	upgraded.PolicyID = PolicyAdmin
	if err := Reconcile(context.Background(), st, []kernel.PermissionContribution{upgraded}, nil); err != nil {
		t.Fatalf("versioned policy upgrade: %v", err)
	}
	viewerPermissions, err = repository.PermissionsForUser("user-viewer")
	if err != nil || len(viewerPermissions) != 0 {
		t.Fatalf("viewer permissions after policy upgrade = %v, err %v", viewerPermissions, err)
	}
	if got := queryInt(t, st, `SELECT version FROM system_data_reconcile WHERE module_id = 'admin.users' AND kind = 'authorization' AND contribution_key = 'users.read'`); got != 2 {
		t.Fatalf("ledger version = %d, want 2", got)
	}
	if err := Reconcile(context.Background(), st, []kernel.PermissionContribution{permission}, nil); err == nil || !strings.Contains(err.Error(), "newer than code") {
		t.Fatalf("version downgrade error = %v, want newer-than-code failure", err)
	}
}

func TestReconcileIdentityConflictRollsBackLedger(t *testing.T) {
	st := openTestStore(t)
	if err := Bootstrap(context.Background(), st, "admin", "hash"); err != nil {
		t.Fatal(err)
	}
	if err := st.WithTx(context.Background(), func(tx *sql.Tx) error {
		_, err := tx.Exec(`INSERT INTO permissions (id, key, description, created_at, updated_at) VALUES ('custom-id','users.read','custom',1,1)`)
		return err
	}); err != nil {
		t.Fatal(err)
	}
	permissions, _ := sampleSystemData()
	if err := Reconcile(context.Background(), st, permissions[:1], nil); err == nil {
		t.Fatal("identity conflict must fail closed")
	}
	if got := queryInt(t, st, `SELECT COUNT(*) FROM system_data_reconcile`); got != 0 {
		t.Fatalf("failed reconcile wrote %d ledger rows, want rollback", got)
	}
}

// A-002 F-004: NeedsBootstrap reports true on a user-less store and false once
// a bootstrap admin exists — the C4 retry gate the composition root relies on.
func TestNeedsBootstrapTracksUserPresence(t *testing.T) {
	st := openTestStore(t)
	needed, err := NeedsBootstrap(context.Background(), st)
	if err != nil {
		t.Fatalf("NeedsBootstrap on empty store: %v", err)
	}
	if !needed {
		t.Fatal("NeedsBootstrap on a user-less store must be true")
	}
	if err := Bootstrap(context.Background(), st, "admin", "hash"); err != nil {
		t.Fatalf("bootstrap: %v", err)
	}
	needed, err = NeedsBootstrap(context.Background(), st)
	if err != nil {
		t.Fatalf("NeedsBootstrap after bootstrap: %v", err)
	}
	if needed {
		t.Fatal("NeedsBootstrap after bootstrap must be false")
	}
}
