// Package jobs provides the admin.jobs module surface as a kernel.Provider
// (GOAL-003 R2 · VP-038): the management-scope surface over the durable async
// Job runtime — a cross-actor job list, job detail and result download gated by
// jobs.read, the batch-export submit route (R3) and the management-scope
// cancel/retry actions (R4) gated by jobs.write.
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
	a         *auth.Authenticator
	reader    handler.JobReader
	submitter handler.JobBatchSubmitter
	actions   handler.JobActions
}

// New constructs the jobs provider over the shared Job read surface.
//
// submitter (R3 batch export) and actions (R4 cancel/retry) may each be nil for
// a composition that only reads; the corresponding routes are then omitted from
// the contribution set entirely — rather than mounted and failing at request
// time — so the routes the descriptor declares always match the routes served.
// The parameters are typed on purpose: an adapter that does not satisfy the
// seam is a compile error, not a silently missing route.
func New(a *auth.Authenticator, reader handler.JobReader, submitter handler.JobBatchSubmitter, actions handler.JobActions) *Provider {
	return &Provider{a: a, reader: reader, submitter: submitter, actions: actions}
}

func (p *Provider) Descriptor() kernel.Module {
	routes := []string{
		"GET /api/jobs",
		"GET /api/jobs/{id}",
		"GET /api/jobs/{id}/result",
	}
	permissions := []string{"jobs.read"}
	if p.submitter != nil {
		// R3 (GOAL-004): the batch-export submit route.
		routes = append(routes, "POST /api/jobs/batch-export")
	}
	if p.actions != nil {
		// R4 (GOAL-005): management-scope cancel/retry (result center).
		routes = append(routes, "POST /api/jobs/{id}/cancel", "POST /api/jobs/{id}/retry")
	}
	if p.submitter != nil || p.actions != nil {
		// jobs.write authorizes mutating the async runtime: submitting a job,
		// cancelling one, or retrying one. Declared once, shared by both.
		permissions = append(permissions, "jobs.write")
	}
	return kernel.Module{
		ID:             ModuleID,
		Version:        "2.0.0",
		KernelAPIRange: ">=2.0 <3.0",
		DependsOn:      []string{"core.auth-session", "core.navigation-capability", "core.schema-render"},
		Requires:       kernel.StandardAdminCapabilities(),
		Contributions: kernel.ContributionKeys{
			Routes:      routes,
			Pages:       []string{"jobs"},
			Navigation:  []string{"menu_jobs"},
			Permissions: permissions,
			Fragments:   []string{"jobs"},
		},
	}
}

func (p *Provider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return nil, nil // the jobs table is owned by the migration-only core.jobs provider
}

func (p *Provider) Register(ctx context.Context, reg kernel.Registrar) error {
	for _, route := range handler.JobsRoutes(p.a, p.reader, p.submitter, p.actions, ModuleID) {
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
	contributions := []kernel.PermissionContribution{{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "jobs.read"},
		Permission:           "jobs.read",
		Resource:             "jobs",
		Action:               "read",
		PolicyID:             authsessiondata.PolicyAdmin,
		SystemDataVersion:    authsessiondata.SystemDataVersion,
	}}
	if p.submitter != nil || p.actions != nil {
		// jobs.write authorizes mutating the async runtime: submitting an async
		// job (R3), cancelling one or retrying one (R4). It is deliberately
		// paired with data.export at the submit route, so this key alone cannot
		// move data out (GOAL-004 D-001 §1); cancel/retry gate on this key alone
		// because they neither read nor move row data (GOAL-005 D-001 §2).
		contributions = append(contributions, kernel.PermissionContribution{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "jobs.write"},
			Permission:           "jobs.write",
			Resource:             "jobs",
			Action:               "write",
			PolicyID:             authsessiondata.PolicyAdmin,
			SystemDataVersion:    authsessiondata.SystemDataVersion,
		})
	}
	for _, permission := range contributions {
		if err := reg.Authorization(permission); err != nil {
			return err
		}
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
