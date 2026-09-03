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
		// S-01 (GOAL-008 D-001 §2): admin.data-dictionary — admin-only profile;
		// two-level dictionary types and entries.
		"admin.data-dictionary",
		// S-03 (GOAL-009 D-001 §2): admin.system-monitoring — admin-only profile;
		// read-only monitoring surface.
		"admin.system-monitoring",
		// S-04 (GOAL-010 D-001 §2): admin.scheduled-tasks — admin-only profile;
		// cron task management + in-process scheduler.
		"admin.scheduled-tasks",
		// S-11 (GOAL-011 D-001 §2): admin.login-captcha — admin-only profile;
		// optional login captcha gate + settings page (default off).
		"admin.login-captcha",
		// S-12 (GOAL-012 D-001 §1): admin.recycle-bin — admin-only profile;
		// deleted-row snapshots with browse/restore/purge.
		"admin.recycle-bin",
		// S-09 (GOAL-016 D-001 §2): admin.data-permission — admin-only profile;
		// row-level scope policy + assignment management (v1 wires no resource).
		"admin.data-permission",
		// S-10 (GOAL-017 D-001 §1): admin.mfa — admin-only profile; TOTP
		// second-factor gate + self-service surface (default: nobody enrolled,
		// login behavior unchanged).
		"admin.mfa",
		// S-14 (GOAL-019 D-002 §3): admin.wallet — admin-only profile; wallet
		// accounts + immutable ledger + reconciliation (content extension).
		"admin.wallet",
	},
	// ProfileDemo is the non-production demonstration profile (W2, GOAL-003 /
	// workspace-010): the full mvp capability surface plus the optional
	// dev.examples module, so a single app.profile: demo (or app.modules.preset:
	// demo) boots the protocol
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
		return ProfileResolution{}, kernelError(CodeProfileModulesRequired, string(profileName), "custom profile requires an explicit app.modules list or preset (config.yaml)")
	}
	modules := append([]string(nil), defaults...)
	source := "profile.default"
	if len(explicitModules) > 0 {
		modules = append([]string(nil), explicitModules...)
		source = "modules.list"
	}
	return ProfileResolution{
		Name:       profileName,
		Modules:    modules,
		Source:     source,
		Precedence: []string{"compiled-profile-default", "modules.preset", "modules.list"},
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
		{ID: "admin.users", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.operationlog"}, Requires: StandardAdminCapabilities(), Contributions: ContributionKeys{Routes: []string{"GET /api/users", "GET /api/users/{id}", "POST /api/users", "PATCH /api/users/{id}", "DELETE /api/users/{id}", "POST /api/users/batch-delete", "GET /api/users/invites", "POST /api/users/invites", "DELETE /api/users/invites/{id}", "POST /api/users/invites/{id}/resend"}, Pages: []string{"users", "users-invites"}, Navigation: []string{"menu_users"}, Permissions: []string{"users.read", "users.write", "users.invite"}, Fragments: []string{"users"}}},
		{ID: "admin.roles", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.operationlog"}, Requires: StandardAdminCapabilities(), Contributions: ContributionKeys{Routes: []string{"GET /api/roles", "GET /api/roles/{id}", "POST /api/roles", "PATCH /api/roles/{id}", "DELETE /api/roles/{id}", "POST /api/roles/batch-delete", "GET /api/permissions", "GET /api/menu-items"}, Pages: []string{"roles"}, Navigation: []string{"menu_roles"}, Permissions: []string{"roles.read", "roles.write", "roles.assign"}, Fragments: []string{"roles"}}},
		// S-09 admin.data-permission (GOAL-016): row-level scope policies/assignments.
		{ID: "admin.data-permission", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.operationlog"}, Requires: StandardAdminCapabilities(), Contributions: ContributionKeys{Routes: []string{"GET /api/data-permission/policies", "PATCH /api/data-permission/policies", "GET /api/data-permission/scopes", "PATCH /api/data-permission/scopes"}, Pages: []string{"data-permission"}, Navigation: []string{"menu_data_permission"}, Permissions: []string{"data-permission.read", "data-permission.write"}, Fragments: []string{"data-permission"}}},
		// S-10 admin.mfa (GOAL-017): TOTP second factor — no page/nav/fragment
		// (D-002 §4: personal-center block + users row action).
		{ID: "admin.mfa", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.operationlog"}, Requires: StandardAdminCapabilities(), Contributions: ContributionKeys{Routes: []string{"POST /api/auth/mfa/verify", "GET /api/mfa/status", "POST /api/mfa/enroll", "POST /api/mfa/confirm", "POST /api/mfa/disable", "POST /api/mfa/recovery/rotate", "POST /api/users/{id}/mfa/reset"}, Permissions: []string{"users.mfa-reset"}}},
		// W26 (GOAL-038): + standalone mail console / outbound log pages and
		// their sidebar nodes — same settings.read permission, no new keys.
		{ID: "admin.settings", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.operationlog"}, Requires: StandardAdminCapabilities(), Contributions: ContributionKeys{Routes: []string{"GET /api/branding", "GET /api/settings", "GET /api/settings/{id}", "PATCH /api/settings/{id}", "POST /api/settings/{id}/reset", "POST /api/branding/assets", "GET /api/branding/assets/{id}", "GET /api/settings/password-policy", "PATCH /api/settings/password-policy"}, Pages: []string{"settings", "mail", "mail-outbox"}, Navigation: []string{"menu_settings", "menu_mail", "menu_mail_outbox"}, Permissions: []string{"settings.read", "settings.write"}, ConfigNamespaces: []string{"settings.branding"}, Fragments: []string{"settings"}}},
		{ID: "admin.activity", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.auth-session", "core.navigation-capability", "core.manifest-route", "core.operationlog"}, Requires: StandardAdminCapabilities(), Contributions: ContributionKeys{Routes: []string{"GET /api/operations", "GET /api/operations/{id}", "GET /api/operations/export"}, Pages: []string{"activity"}, Navigation: []string{"menu_activity"}, Permissions: []string{"operations.read"}, Fragments: []string{"activity"}}},
		// F-03 admin.account (GOAL-005): self-service + enable/disable/unlock.
		// workspace-018 R3: + email bind/verify/resend (GOAL-004 D-001 §3).
		{ID: "admin.account", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.operationlog"}, Requires: StandardAdminCapabilities(), Contributions: ContributionKeys{Routes: []string{"GET /api/account/profile", "PATCH /api/account/profile", "POST /api/account/avatar", "GET /api/account/avatars/{id}", "POST /api/account/password", "GET /api/account/sessions", "POST /api/account/sessions/{id}/revoke", "POST /api/account/sessions/revoke-others", "POST /api/account/email/bind", "POST /api/account/email/verify", "POST /api/account/email/resend", "POST /api/users/{id}/enable", "POST /api/users/{id}/disable", "POST /api/users/{id}/unlock"}, Pages: []string{"account"}, Navigation: []string{"menu_account"}, Permissions: []string{"users.enable", "users.disable"}, Fragments: []string{"account"}}},
		// F-02 admin.data-transfer (GOAL-004): CSV export/import shared capability.
		{ID: "admin.data-transfer", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.auth-session", "core.schema-render", "core.operationlog"}, Requires: StandardAdminCapabilities(), Contributions: ContributionKeys{Routes: []string{"GET /api/export/{resource}", "POST /api/import/{resource}", "GET /api/import/{resource}/template"}, Permissions: []string{"data.export", "data.import"}}},
		// F-01 admin.dashboard (GOAL-003): production home dashboard (no routes).
		{ID: "admin.dashboard", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.manifest-route"}, Requires: StandardAdminCapabilities(), Contributions: ContributionKeys{Pages: []string{"dashboard"}, Navigation: []string{"menu_dashboard"}, Fragments: []string{"dashboard"}}},
		// F-04 admin.notifications (GOAL-006): in-app notifications.
		{ID: "admin.notifications", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.operationlog"}, Requires: StandardAdminCapabilities(), Contributions: ContributionKeys{Routes: []string{"GET /api/notifications", "POST /api/notifications/{id}/read", "POST /api/notifications/read-all", "GET /api/notifications/unread-count", "GET /api/notifications/settings", "PATCH /api/notifications/settings"}, Pages: []string{"notifications"}, Navigation: []string{"menu_notifications"}, Fragments: []string{"notifications"}}},
		// S-02 admin.file-library (GOAL-007): unified file/attachment library.
		{ID: "admin.file-library", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.operationlog"}, Requires: StandardAdminCapabilities(), Contributions: ContributionKeys{Routes: []string{"GET /api/library/files", "GET /api/library/files/{id}", "GET /api/library/files/{id}/download", "DELETE /api/library/files/{id}", "POST /api/library/files/upload"}, Pages: []string{"file-library"}, Navigation: []string{"menu_files"}, Permissions: []string{"files.read", "files.delete"}, Fragments: []string{"file-library"}}},
		// S-01 admin.data-dictionary (GOAL-008): two-level dictionary types/entries.
		{ID: "admin.data-dictionary", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.operationlog"}, Requires: StandardAdminCapabilities(), Contributions: ContributionKeys{Routes: []string{"GET /api/data-dictionary/types", "GET /api/data-dictionary/types/{id}", "POST /api/data-dictionary/types", "PATCH /api/data-dictionary/types/{id}", "DELETE /api/data-dictionary/types/{id}", "POST /api/data-dictionary/types/batch-delete", "GET /api/data-dictionary/entries", "GET /api/data-dictionary/entries/{id}", "POST /api/data-dictionary/entries", "PATCH /api/data-dictionary/entries/{id}", "DELETE /api/data-dictionary/entries/{id}", "POST /api/data-dictionary/entries/batch-delete"}, Pages: []string{"data-dictionary", "dictionary-entries"}, Navigation: []string{"menu_dictionary"}, Permissions: []string{"dictionary.read", "dictionary.write"}, Fragments: []string{"data-dictionary"}}},
		// S-03 admin.system-monitoring (GOAL-009): read-only monitoring surface.
		{ID: "admin.system-monitoring", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.operationlog"}, Requires: StandardAdminCapabilities(), Contributions: ContributionKeys{Routes: []string{"GET /api/system-monitoring/status", "GET /api/system-monitoring/errors", "GET /api/system-monitoring/errors/{id}"}, Pages: []string{"system-monitoring"}, Navigation: []string{"menu_monitoring"}, Permissions: []string{"monitoring.read"}, Fragments: []string{"system-monitoring"}}},
		// S-04 admin.scheduled-tasks (GOAL-010): cron task management + scheduler.
		{ID: "admin.scheduled-tasks", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.operationlog"}, Requires: StandardAdminCapabilities(), Contributions: ContributionKeys{Routes: []string{"GET /api/scheduled-tasks", "GET /api/scheduled-tasks/{id}", "GET /api/scheduled-tasks/handlers", "POST /api/scheduled-tasks/cron/preview", "POST /api/scheduled-tasks", "PATCH /api/scheduled-tasks/{id}", "DELETE /api/scheduled-tasks/{id}", "POST /api/scheduled-tasks/batch-delete", "POST /api/scheduled-tasks/{id}/run", "GET /api/scheduled-tasks/{id}/runs", "GET /api/task-runs", "GET /api/task-runs/{id}"}, Pages: []string{"scheduled-tasks", "task-runs"}, Navigation: []string{"menu_scheduled_tasks"}, Permissions: []string{"tasks.read", "tasks.write"}, Fragments: []string{"scheduled-tasks"}}},
		// S-11 admin.login-captcha (GOAL-011): optional login captcha gate.
		// D-003 (user ruling): no page/nav/fragment — the switch lives in the
		// admin.settings security section.
		{ID: "admin.login-captcha", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.operationlog"}, Requires: StandardAdminCapabilities(), Contributions: ContributionKeys{Routes: []string{"GET /api/auth/captcha", "GET /api/captcha/settings", "PATCH /api/captcha/settings"}, Permissions: []string{"captcha.read", "captcha.write"}}},
		// S-12 admin.recycle-bin (GOAL-012): deleted-row snapshots.
		{ID: "admin.recycle-bin", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.operationlog"}, Requires: StandardAdminCapabilities(), Contributions: ContributionKeys{Routes: []string{"GET /api/recycle-bin", "GET /api/recycle-bin/{id}", "POST /api/recycle-bin/{id}/restore", "DELETE /api/recycle-bin/{id}", "POST /api/recycle-bin/purge-all"}, Pages: []string{"recycle-bin"}, Navigation: []string{"menu_recycle_bin"}, Permissions: []string{"recycle.read", "recycle.write"}, Fragments: []string{"recycle-bin"}}},
		// S-14 admin.wallet (GOAL-019): wallet/ledger — accounts + immutable
		// ledger + reconciliation. Money-path mutations use wallet.adjust.
		{ID: "admin.wallet", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.operationlog"}, Requires: StandardAdminCapabilities(), Contributions: ContributionKeys{Routes: []string{"GET /api/wallet/accounts", "POST /api/wallet/accounts", "PATCH /api/wallet/accounts/{id}", "GET /api/wallet/accounts/{id}/entries", "GET /api/wallet/entries", "GET /api/wallet/by-owner/{ownerId}", "POST /api/wallet/by-owner/{ownerId}", "POST /api/wallet/by-owner/{ownerId}/adjust", "POST /api/wallet/accounts/{id}/adjust", "POST /api/wallet/accounts/{id}/freeze", "POST /api/wallet/accounts/{id}/unfreeze", "POST /api/wallet/accounts/{id}/deduct-frozen", "POST /api/wallet/reconcile", "GET /api/wallet/reconcile/runs", "GET /api/wallet/jobs/{id}", "POST /api/wallet/jobs/{id}/cancel", "POST /api/wallet/jobs/{id}/retry", "GET /api/wallet/jobs/{id}/result", "GET /api/wallet/me", "POST /api/wallet/me", "GET /api/wallet/me/entries", "POST /api/wallet/me/redeem", "POST /api/wallet/vouchers/batches", "GET /api/wallet/vouchers", "GET /api/wallet/vouchers/{id}", "POST /api/wallet/vouchers/{id}/void"}, Pages: []string{"wallet", "wallet-entries", "my-wallet", "wallet-vouchers"}, Navigation: []string{"menu_wallet", "menu_wallet_self", "menu_wallet_vouchers"}, Permissions: []string{"wallet.read", "wallet.write", "wallet.adjust", "wallet.voucher.issue"}, Fragments: []string{"wallet"}}},
		// dev.examples is the optional demonstration module (W1, GOAL-002): it owns
		// the 8 example pages + Examples navigation as a horizontal demo surface.
		// It is compiled but never enabled by mvp/admin defaults; enable via
		// app.modules (config.yaml) or a dedicated dogfood profile (D-003 §3).
		{ID: "dev.examples", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.schema-render", "core.navigation-capability"}, Contributions: ContributionKeys{Pages: []string{"overview", "data-table", "search-form-table", "form-controls", "form-with-reactions", "form-with-upload", "data-display", "admin-list-batch"}, Fragments: []string{"examples"}}},
		// VP-030 (GOAL-003/GOAL-004): channel.telegram — Telegram bot channel runtime.
		// Exposes public webhook POST /api/channel/telegram/webhook and settings endpoints.
		// Compiled candidate; not enabled in mvp/admin defaults (channel extension).
		{ID: "channel.telegram", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", DependsOn: []string{"core.server-registration"}, Requires: []Capability{CapabilityHTTP}, Contributions: ContributionKeys{Routes: []string{"GET /api/channel/telegram/settings", "PATCH /api/channel/telegram/settings", "POST /api/channel/telegram/webhook"}}},
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
