// Package systemmonitoring provides the admin.system-monitoring module
// surface as a kernel.Provider (S-03 · GOAL-009 D-002): a read-only monitoring
// page — status summary (in-process probes) and a recent-events view. No
// persistence, no audit events.
package systemmonitoring

import (
	"context"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authsessiondata "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession/systemdata"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/systemmonitoring/manifest"
	monitoringschema "github.com/magicvr/schema-ui-core/apps/api/internal/modules/systemmonitoring/schema"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
)

// ModuleID is the stable admin.system-monitoring module identifier.
const ModuleID = "admin.system-monitoring"

// Provider implements kernel.Provider for admin.system-monitoring.
type Provider struct {
	a                *auth.Authenticator
	st               *store.Store
	plan             kernel.Plan
	ready            func() bool
	dbPath           string
	startTime        time.Time
	operations       operationlog.Reader
	availabilityMode string
}

// New constructs the monitoring provider with framework-agnostic dependencies.
func New(a *auth.Authenticator, st *store.Store, plan kernel.Plan, ready func() bool, dbPath string, startTime time.Time, operations operationlog.Reader, mode ...string) *Provider {
	availabilityMode := "normal"
	if len(mode) > 0 && mode[0] != "" {
		availabilityMode = mode[0]
	}
	return &Provider{a: a, st: st, plan: plan, ready: ready, dbPath: dbPath, startTime: startTime, operations: operations, availabilityMode: availabilityMode}
}

func (p *Provider) Descriptor() kernel.Module {
	return kernel.Module{
		ID:             ModuleID,
		Version:        "2.0.0",
		KernelAPIRange: ">=2.0 <3.0",
		DependsOn:      []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.operationlog"},
		Requires:       kernel.StandardAdminCapabilities(),
		Contributions: kernel.ContributionKeys{
			Routes: []string{
				"GET /api/system-monitoring/status",
				"GET /api/system-monitoring/errors", "GET /api/system-monitoring/errors/{id}",
			},
			Pages:       []string{"system-monitoring"},
			Navigation:  []string{"menu_monitoring"},
			Permissions: []string{"monitoring.read"},
			Fragments:   []string{"system-monitoring"},
		},
	}
}

func (p *Provider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return nil, nil // read-only module; no tables
}

func (p *Provider) Register(ctx context.Context, reg kernel.Registrar) error {
	for _, route := range handler.MonitoringRoutes(p.a, p.st, p.plan, p.ready, p.dbPath, p.startTime, p.operations, p.availabilityMode, ModuleID) {
		if err := reg.HTTP(route); err != nil {
			return err
		}
	}
	if err := reg.Schema(kernel.PageContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "system-monitoring"},
		PageID:               "system-monitoring",
		Resources:            []string{"monitoring-errors"},
		Actions:              []string{"list"},
		DataSource:           "/api/system-monitoring/errors",
		Owner:                ModuleID,
		Document:             monitoringschema.SchemaDocuments()["system-monitoring"],
	}); err != nil {
		return err
	}
	if err := reg.Authorization(kernel.PermissionContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "monitoring.read"},
		Permission:           "monitoring.read",
		Resource:             "system-monitoring",
		Action:               "read",
		PolicyID:             authsessiondata.PolicyAdmin,
		SystemDataVersion:    authsessiondata.SystemDataVersion,
	}); err != nil {
		return err
	}
	if err := reg.Navigation(kernel.NavigationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "menu_monitoring"},
		NodeID:               "menu_monitoring",
		PageID:               "system-monitoring",
		Order:                5,
		Label:                "System monitoring",
		Visibility:           authsessiondata.PolicyAdmin,
		Permission:           "monitoring.read",
		SystemDataVersion:    authsessiondata.SystemDataVersion,
	}); err != nil {
		return err
	}
	return reg.Manifest(kernel.FragmentContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "system-monitoring"},
		FragmentID:           "system-monitoring",
		ProtocolVersion:      "2.7",
		RequiredCapabilities: []string{"manifest", "navigation"},
		JSON:                 manifest.FragmentJSON,
	})
}
