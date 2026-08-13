package kernel

import (
	"fmt"
	"sort"
	"strings"
)

type ProfileName string

const (
	ProfileMVP    ProfileName = "mvp"
	ProfileAdmin  ProfileName = "admin"
	ProfileDemo   ProfileName = "demo"
	ProfileCustom ProfileName = "custom"
)

type ProfileResolution struct {
	Name       ProfileName
	Modules    []string
	Source     string
	Precedence []string
}

var profileDefaults = map[ProfileName][]string{
	ProfileMVP: {
		"core.server-registration",
		"core.auth-session",
		"core.manifest-route",
		"core.navigation-capability",
		"core.schema-render",
		"core.operationlog",
		"admin.users",
		"admin.roles",
		// F-03 (GOAL-005 D-002 §6): self-service password/session management is
		// an account-security baseline; Profile content extension, not an
		// assembly-semantics change.
		"admin.account",
		// F-01 (GOAL-003 D-002 §3): production home dashboard — Profile content
		// extension (explicit 必办-3 declaration).
		"admin.dashboard",
		// F-04 (GOAL-006 D-002 §6): in-app security notifications — account
		// safety baseline (same rationale as F-03).
		"admin.notifications",
	},
	ProfileAdmin: {
		"core.server-registration",
		"core.auth-session",
		"core.manifest-route",
		"core.navigation-capability",
		"core.schema-render",
		"core.operationlog",
		"admin.users",
		"admin.roles",
		"admin.settings",
		"admin.activity",
		"admin.account",
		// F-02 (GOAL-004 D-002 §6): admin.data-transfer — admin-only profile;
		// export/import is a management-surface capability (content extension).
		"admin.data-transfer",
		// F-01 (GOAL-003 D-002 §3): production home dashboard.
		"admin.dashboard",
		// F-04 (GOAL-006 D-002 §6): in-app security notifications.
		"admin.notifications",
		// S-02 (GOAL-007 D-001 §2): admin.file-library — admin-only profile;
		// unified file/attachment library over the shared upload store.
		"admin.file-library",
	},
	// ProfileDemo is the non-production demonstration profile (W2, GOAL-003 /
	// workspace-010): the full mvp capability surface plus the optional
	// dev.examples module, so a single APP_PROFILE=demo boots the protocol
	// examples alongside the real mvp pages. It is never a production default.
	ProfileDemo: {
		"core.server-registration",
		"core.auth-session",
		"core.manifest-route",
		"core.navigation-capability",
		"core.schema-render",
		"core.operationlog",
		"admin.users",
		"admin.roles",
		"dev.examples",
		"admin.account",
		"admin.dashboard",
		// F-04 (GOAL-006): in-app security notifications (inherited via demo).
		"admin.notifications",
	},
}

func ParseModuleList(raw string) ([]string, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	parts := strings.Split(raw, ",")
	modules := make([]string, 0, len(parts))
	for _, part := range parts {
		id := strings.TrimSpace(part)
		if id == "" {
			return nil, fmt.Errorf("module list contains an empty id")
		}
		modules = append(modules, id)
	}
	return modules, nil
}

func ResolveProfile(name string, explicitModules []string) (ProfileResolution, error) {
	profileName := ProfileName(strings.ToLower(strings.TrimSpace(name)))
	if profileName == "" {
		profileName = ProfileMVP
	}
	defaults, known := profileDefaults[profileName]
	if !known && profileName != ProfileCustom {
		return ProfileResolution{}, kernelError(CodeProfileUnknown, string(profileName), "profile is not compiled")
	}
	if profileName == ProfileCustom && len(explicitModules) == 0 {
		return ProfileResolution{}, kernelError(CodeProfileModulesRequired, string(profileName), "custom profile requires APP_MODULES_ENABLED")
	}
	modules := append([]string(nil), defaults...)
	source := "profile.default"
	if len(explicitModules) > 0 {
		modules = append([]string(nil), explicitModules...)
		source = "modules.enabled"
	}
	return ProfileResolution{
		Name:       profileName,
		Modules:    modules,
		Source:     source,
		Precedence: []string{"compiled-profile-default", "modules.enabled", "environment"},
	}, nil
}

func BuiltinModules() []Module {
	return []Module{
		{ID: "core.server-registration", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", Provides: []Capability{CapabilityHTTP}},
		{ID: "core.auth-session", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.server-registration"}, Provides: []Capability{CapabilityAuthorization, CapabilityPersistence}, Contributions: ContributionKeys{Routes: []string{"/api/auth"}}},
		{ID: "core.schema-render", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.server-registration"}, Provides: []Capability{CapabilitySchema, CapabilityValidation}},
		{ID: "core.manifest-route", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.server-registration", "core.schema-render"}, Provides: []Capability{CapabilityManifest}, Contributions: ContributionKeys{Routes: []string{"/.well-known/schema-ui/app-manifest.json", "/.well-known/schema-ui/host-bootstrap.json"}}},
		{ID: "core.navigation-capability", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.auth-session", "core.manifest-route"}, Provides: []Capability{CapabilityNavigation, CapabilityExpressions}},
		{ID: "core.operationlog", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.server-registration"}, Provides: []Capability{CapabilityOperationLog}},
		{ID: "admin.users", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.operationlog"}, Requires: StandardAdminCapabilities(), Contributions: ContributionKeys{Routes: []string{"GET /api/users", "GET /api/users/{id}", "POST /api/users", "PATCH /api/users/{id}", "DELETE /api/users/{id}", "POST /api/users/batch-delete"}, Pages: []string{"users"}, Navigation: []string{"menu_users"}, Permissions: []string{"users.read", "users.write"}, Fragments: []string{"users"}}},
		{ID: "admin.roles", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.operationlog"}, Requires: StandardAdminCapabilities(), Contributions: ContributionKeys{Routes: []string{"GET /api/roles", "GET /api/roles/{id}", "POST /api/roles", "PATCH /api/roles/{id}", "DELETE /api/roles/{id}", "POST /api/roles/batch-delete"}, Pages: []string{"roles"}, Navigation: []string{"menu_roles"}, Permissions: []string{"roles.read", "roles.write", "roles.assign"}, Fragments: []string{"roles"}}},
		{ID: "admin.settings", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.operationlog"}, Requires: StandardAdminCapabilities(), Contributions: ContributionKeys{Routes: []string{"GET /api/branding", "GET /api/settings", "GET /api/settings/{id}", "PATCH /api/settings/{id}", "POST /api/settings/{id}/reset"}, Pages: []string{"settings"}, Navigation: []string{"menu_settings"}, Permissions: []string{"settings.read", "settings.write"}, ConfigNamespaces: []string{"settings.branding"}, Fragments: []string{"settings"}}},
		{ID: "admin.activity", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.auth-session", "core.navigation-capability", "core.manifest-route", "core.operationlog"}, Requires: StandardAdminCapabilities(), Contributions: ContributionKeys{Routes: []string{"GET /api/operations", "GET /api/operations/{id}"}, Pages: []string{"activity"}, Navigation: []string{"menu_activity"}, Permissions: []string{"operations.read"}, Fragments: []string{"activity"}}},
		// F-03 admin.account (GOAL-005): self-service + enable/disable/unlock.
		{ID: "admin.account", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.operationlog"}, Requires: StandardAdminCapabilities(), Contributions: ContributionKeys{Routes: []string{"GET /api/account/profile", "PATCH /api/account/profile", "POST /api/account/password", "GET /api/account/sessions", "POST /api/account/sessions/{id}/revoke", "POST /api/users/{id}/enable", "POST /api/users/{id}/disable", "POST /api/users/{id}/unlock"}, Pages: []string{"account"}, Navigation: []string{"menu_account"}, Permissions: []string{"users.enable", "users.disable"}, Fragments: []string{"account"}}},
		// F-02 admin.data-transfer (GOAL-004): CSV export/import shared capability.
		{ID: "admin.data-transfer", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.auth-session", "core.schema-render", "core.operationlog"}, Requires: StandardAdminCapabilities(), Contributions: ContributionKeys{Routes: []string{"GET /api/export/{resource}", "POST /api/import/{resource}"}, Permissions: []string{"data.export", "data.import"}}},
		// F-01 admin.dashboard (GOAL-003): production home dashboard (no routes).
		{ID: "admin.dashboard", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.manifest-route"}, Requires: StandardAdminCapabilities(), Contributions: ContributionKeys{Pages: []string{"dashboard"}, Navigation: []string{"menu_dashboard"}, Fragments: []string{"dashboard"}}},
		// F-04 admin.notifications (GOAL-006): in-app notifications.
		{ID: "admin.notifications", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.operationlog"}, Requires: StandardAdminCapabilities(), Contributions: ContributionKeys{Routes: []string{"GET /api/notifications", "POST /api/notifications/{id}/read", "POST /api/notifications/read-all", "GET /api/notifications/unread-count", "GET /api/notifications/settings", "PATCH /api/notifications/settings"}, Pages: []string{"notifications"}, Navigation: []string{"menu_notifications"}, Fragments: []string{"notifications"}}},
		// S-02 admin.file-library (GOAL-007): unified file/attachment library.
		{ID: "admin.file-library", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.operationlog"}, Requires: StandardAdminCapabilities(), Contributions: ContributionKeys{Routes: []string{"GET /api/library/files", "GET /api/library/files/{id}", "GET /api/library/files/{id}/download", "DELETE /api/library/files/{id}", "POST /api/library/files/upload"}, Pages: []string{"file-library"}, Navigation: []string{"menu_files"}, Permissions: []string{"files.read", "files.delete"}, Fragments: []string{"file-library"}}},
		// dev.examples is the optional demonstration module (W1, GOAL-002): it owns
		// the 8 example pages + Examples navigation as a horizontal demo surface.
		// It is compiled but never enabled by mvp/admin defaults; enable via
		// APP_MODULES_ENABLED or a dedicated dogfood profile (D-003 §3).
		{ID: "dev.examples", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.schema-render", "core.navigation-capability"}, Contributions: ContributionKeys{Pages: []string{"overview", "data-table", "search-form-table", "form-controls", "form-with-reactions", "form-with-upload", "data-display", "admin-list-batch"}, Fragments: []string{"examples"}}},
	}
}

// StandardAdminCapabilities is the requirement set of every standard Admin
// functional module (freeze package: core six; "按需" never overrides them).
// Exported so module providers declare Requires without duplicating the set.
func StandardAdminCapabilities() []Capability {
	return []Capability{CapabilityHTTP, CapabilitySchema, CapabilityAuthorization, CapabilityNavigation, CapabilityManifest, CapabilityPersistence}
}

func SortedModuleIDs(modules []Module) []string {
	ids := make([]string, 0, len(modules))
	for _, module := range modules {
		ids = append(ids, module.ID)
	}
	sort.Strings(ids)
	return ids
}