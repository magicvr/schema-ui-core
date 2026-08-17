// Package scheduledtasks provides the admin.scheduled-tasks module surface as
// a kernel.Provider (S-04 · GOAL-010 D-002): task definition CRUD, an
// in-process best-effort scheduler (system.noop handler), run history, manual
// trigger, tasks.read / tasks.write permission keys and scheduled-tasks.*
// audit events.
package scheduledtasks

import (
	"context"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authsessiondata "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession/systemdata"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/scheduledtasks/manifest"
	tasksschema "github.com/magicvr/schema-ui-core/apps/api/internal/modules/scheduledtasks/schema"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/scheduledtasks/store"
)

// ModuleID is the stable admin.scheduled-tasks module identifier.
const ModuleID = "admin.scheduled-tasks"

// Provider implements kernel.Provider for admin.scheduled-tasks.
type Provider struct {
	a          *auth.Authenticator
	repository *store.Repository
	runner     *Scheduler
	operations operationlog.Recorder
	trash      handler.TrashRecorder
}

// New constructs the tasks provider and starts the best-effort scheduler loop.
// trash (S-12 · GOAL-012 D-002 §2), when non-nil, opts the tasks resource into
// recycle-bin snapshot recording on delete.
func New(a *auth.Authenticator, repository *store.Repository, operations operationlog.Recorder, trash ...handler.TrashRecorder) *Provider {
	scheduler := NewScheduler(repository)
	scheduler.Start()
	var recorder handler.TrashRecorder
	if len(trash) > 0 {
		recorder = trash[0]
	}
	return &Provider{a: a, repository: repository, runner: scheduler, operations: operations, trash: recorder}
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
				"GET /api/scheduled-tasks", "GET /api/scheduled-tasks/{id}",
				"GET /api/scheduled-tasks/handlers",
				"POST /api/scheduled-tasks", "PATCH /api/scheduled-tasks/{id}",
				"DELETE /api/scheduled-tasks/{id}", "POST /api/scheduled-tasks/batch-delete",
				"POST /api/scheduled-tasks/{id}/run", "GET /api/scheduled-tasks/{id}/runs",
				"GET /api/task-runs", "GET /api/task-runs/{id}",
			},
			Pages:       []string{"scheduled-tasks", "task-runs"},
			Navigation:  []string{"menu_scheduled_tasks"},
			Permissions: []string{"tasks.read", "tasks.write"},
			Fragments:   []string{"scheduled-tasks"},
		},
	}
}

func (p *Provider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return nil, nil // tables are owned by the scheduledtasks/migration provider (0021)
}

func (p *Provider) Register(ctx context.Context, reg kernel.Registrar) error {
	for _, route := range handler.ScheduledTaskRoutes(p.a, p.repository, p.runner, p.operations, ModuleID, p.trash) {
		if err := reg.HTTP(route); err != nil {
			return err
		}
	}
	for _, pageID := range []string{"scheduled-tasks", "task-runs"} {
		if err := reg.Schema(kernel.PageContribution{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: pageID},
			PageID:               pageID,
			Resources:            []string{"scheduled-tasks", "task-runs"},
			Actions:              []string{"list", "create", "detail", "update", "delete"},
			DataSource:           "/api/scheduled-tasks",
			Owner:                ModuleID,
			Document:             tasksschema.SchemaDocuments()[pageID],
		}); err != nil {
			return err
		}
	}
	for _, permission := range []kernel.PermissionContribution{
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "tasks.read"}, Permission: "tasks.read", Resource: "scheduled-tasks", Action: "read", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "tasks.write"}, Permission: "tasks.write", Resource: "scheduled-tasks", Action: "write", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
	} {
		if err := reg.Authorization(permission); err != nil {
			return err
		}
	}
	if err := reg.Navigation(kernel.NavigationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "menu_scheduled_tasks"},
		NodeID:               "menu_scheduled_tasks",
		PageID:               "scheduled-tasks",
		Order:                6,
		Label:                "Scheduled tasks",
		Visibility:           authsessiondata.PolicyAdmin,
		Permission:           "tasks.read",
		SystemDataVersion:    authsessiondata.SystemDataVersion,
	}); err != nil {
		return err
	}
	return reg.Manifest(kernel.FragmentContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "scheduled-tasks"},
		FragmentID:           "scheduled-tasks",
		ProtocolVersion:      "2.7",
		RequiredCapabilities: []string{"manifest", "navigation"},
		JSON:                 manifest.FragmentJSON,
	})
}
