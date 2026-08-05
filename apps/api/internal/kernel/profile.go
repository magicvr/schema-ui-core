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
		return ProfileResolution{}, kernelError(CodeProfileModulesRequired, string(profileName), "custom profile requires MODULES_ENABLED")
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
		{ID: "core.manifest-route", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.server-registration", "core.schema-render"}, Provides: []Capability{CapabilityManifest}, Contributions: ContributionKeys{Routes: []string{"/.well-known/schema-ui/app-manifest.json"}}},
		{ID: "core.navigation-capability", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.auth-session", "core.manifest-route"}, Provides: []Capability{CapabilityNavigation, CapabilityExpressions}},
		{ID: "core.operationlog", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.server-registration"}, Provides: []Capability{CapabilityOperationLog}},
		{ID: "admin.users", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.operationlog"}, Requires: StandardAdminCapabilities(), Contributions: ContributionKeys{Routes: []string{"GET /api/users", "GET /api/users/{id}", "POST /api/users", "PATCH /api/users/{id}", "DELETE /api/users/{id}"}, Pages: []string{"users"}, Navigation: []string{"menu_users"}, Permissions: []string{"users.read", "users.write"}, Fragments: []string{"users"}}},
		{ID: "admin.roles", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.operationlog"}, Requires: StandardAdminCapabilities(), Contributions: ContributionKeys{Routes: []string{"GET /api/roles", "GET /api/roles/{id}", "POST /api/roles", "PATCH /api/roles/{id}", "DELETE /api/roles/{id}"}, Pages: []string{"roles"}, Navigation: []string{"menu_roles"}, Permissions: []string{"roles.read", "roles.write", "roles.assign"}, Fragments: []string{"roles"}}},
		{ID: "admin.settings", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.operationlog"}, Requires: StandardAdminCapabilities(), Contributions: ContributionKeys{Pages: []string{"settings"}, Navigation: []string{"menu_settings"}, Permissions: []string{"settings.read", "settings.write"}}},
		{ID: "admin.activity", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.auth-session", "core.navigation-capability", "core.manifest-route", "core.operationlog"}, Requires: StandardAdminCapabilities(), Contributions: ContributionKeys{Pages: []string{"activity"}, Navigation: []string{"menu_activity"}, Permissions: []string{"operations.read"}}},
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
