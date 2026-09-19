// Package jobs provides the admin.jobs module surface as a kernel.Provider
// (GOAL-003 R2 · VP-038): the management-scope read surface over the durable
// async Job runtime — a cross-actor job list, job detail and result download,
// gated by jobs.read.
//
// Scope boundary (GOAL-002 D-001 §2.1): this module answers "every job" for an
// operator. It does NOT replace the actor-scoped admin.wallet job routes
// (/api/wallet/jobs/*), whose id+kind+actor_id isolation stays frozen. The two
// views share the same runtime and table; they differ only in visibility.
//
// The jobs table itself is owned by the migration-only core.jobs provider
// (modules/jobs/migration); this module contributes no persistence.
package jobs

import (
	"context"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	authsessiondata "github.com/magicvr/schema-ui-core/apps/api/modules/authsession/systemdata"
	"github.com/magicvr/schema-ui-core/apps/api/modules/jobs/manifest"
	jobsschema "github.com/magicvr/schema-ui-core/apps/api/modules/jobs/schema"
)

// ModuleID is the stable admin.jobs module identifier.
const ModuleID = "admin.jobs"

// Provider implements kernel.Provider for admin.jobs.
type Provider struct {
	a      *auth.Authenticator
	reader handler.JobReader
}

// New constructs the jobs provider over the shared Job read surface.
func New(a *auth.Authenticator, reader handler.JobReader) *Provider {
	return &Provider{a: a, reader: reader}
}

func (p *Provider) Descriptor() kernel.Module {
	return kernel.Module{
		ID:             ModuleID,
		Version:        "2.0.0",
		KernelAPIRange: ">=2.0 <3.0",
		DependsOn:      []string{"core.auth-session", "core.navigation-capability", "core.schema-render"},
		Requires:       kernel.StandardAdminCapabilities(),
		Contributions: kernel.ContributionKeys{
			Routes: []string{
				"GET /api/jobs",
				"GET /api/jobs/{id}",
				"GET /api/jobs/{id}/result",
			},
			Pages:       []string{"jobs"},
			Navigation:  []string{"menu_jobs"},
			Permissions: []string{"jobs.read"},
			Fragments:   []string{"jobs"},
		},
	}
}

func (p *Provider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return nil, nil // the jobs table is owned by the migration-only core.jobs provider
}

func (p *Provider) Register(ctx context.Context, reg kernel.Registrar) error {
	for _, route := range handler.JobsRoutes(p.a, p.reader, ModuleID) {
		if err := reg.HTTP(route); err != nil {
			return err
		}
	}
	if err := reg.Schema(kernel.PageContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "jobs"},
		PageID:               "jobs",
		Resources:            []string{"jobs"},
		Actions:              []string{"list", "detail"},
		DataSource:           "/api/jobs",
		Owner:                ModuleID,
		Document:             jobsschema.SchemaDocuments()["jobs"],
	}); err != nil {
		return err
	}
	// jobs.read is the management-scope read key. PolicyAdmin mirrors
	// admin.scheduled-tasks' tasks.read: a job row carries another actor's
	// error messages and correlation ids, so visibility is an operator
	// concern (GOAL-003 D-001 §4).
	if err := reg.Authorization(kernel.PermissionContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "jobs.read"},
		Permission:           "jobs.read",
		Resource:             "jobs",
		Action:               "read",
		PolicyID:             authsessiondata.PolicyAdmin,
		SystemDataVersion:    authsessiondata.SystemDataVersion,
	}); err != nil {
		return err
	}
	if err := reg.Navigation(kernel.NavigationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "menu_jobs"},
		NodeID:               "menu_jobs",
		PageID:               "jobs",
		Order:                7,
		Label:                "Jobs",
		Group:                &kernel.NavigationGroup{Key: "operations", Order: 30, Label: "Operations", LabelKey: "manifest.nav.group.operations", Secondary: "OPS"},
		Visibility:           authsessiondata.PolicyAdmin,
		Permission:           "jobs.read",
		SystemDataVersion:    authsessiondata.SystemDataVersion,
	}); err != nil {
		return err
	}
	return reg.Manifest(kernel.FragmentContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "jobs"},
		FragmentID:           "jobs",
		ProtocolVersion:      "2.7",
		RequiredCapabilities: []string{"manifest", "navigation"},
		JSON:                 manifest.FragmentJSON,
	})
}
