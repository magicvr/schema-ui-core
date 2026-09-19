package handler

import (
	"context"
	"fmt"

	"github.com/magicvr/schema-ui-core/apps/api/internal/jobs"
)

// BatchExportSubmitting is the module-side batch-export service the submitter
// drives. Declared here so the handler package does not import the module.
type BatchExportSubmitting interface {
	Supported(resource string) bool
	Submit(ctx context.Context, runner *jobs.Runner, resource string, ids []string, actorID, correlationID string) (*jobs.Job, error)
}

// BatchExportSubmitter adapts the module-side batch-export service to the
// handler's JobBatchSubmitter seam (GOAL-004 R3).
type BatchExportSubmitter struct {
	service BatchExportSubmitting
	runner  *jobs.Runner
}

// NewBatchExportSubmitter binds the service to the runner that executes it.
func NewBatchExportSubmitter(service BatchExportSubmitting, runner *jobs.Runner) *BatchExportSubmitter {
	return &BatchExportSubmitter{service: service, runner: runner}
}

// Supported reports whether the resource can be batch-exported.
func (s *BatchExportSubmitter) Supported(resource string) bool { return s.service.Supported(resource) }

// SubmitBatchExport enqueues the export job.
func (s *BatchExportSubmitter) SubmitBatchExport(ctx context.Context, resource string, ids []string, actorID, correlationID string) (*jobs.Job, error) {
	return s.service.Submit(ctx, s.runner, resource, ids, actorID, correlationID)
}


// BatchExportRowSource adapts the users/roles resource entities to the async
// batch-export row source (GOAL-004 R3). It exists so the Job handler can build
// an export without importing the modules that own those resources, and so the
// column sets stay in one place: ExportSelectedRows delegates to
// SelectedExportRows, which reuses exportHeaders/exportRow — the exact same
// column order, JSON array serialization and formula neutralization the
// synchronous export uses.
type BatchExportRowSource struct {
	users ResourceEntity
	roles ResourceEntity
}

// NewBatchExportRowSource binds the two exportable entities. A nil entity makes
// its resource unsupported rather than panicking.
func NewBatchExportRowSource(users, roles ResourceEntity) *BatchExportRowSource {
	return &BatchExportRowSource{users: users, roles: roles}
}

// SupportedExportResources returns the resources this source can export. It is
// the intersection of the frozen export denominator and the entities actually
// wired, so an unwired resource is never advertised.
func (s *BatchExportRowSource) SupportedExportResources() []string {
	supported := make([]string, 0, len(ExportableResources))
	for _, resource := range ExportableResources {
		if s.entityFor(resource) != nil {
			supported = append(supported, resource)
		}
	}
	return supported
}

// ExportSelectedRows renders the CSV rows for an explicit selection.
func (s *BatchExportRowSource) ExportSelectedRows(resource string, ids []string, onRow func(done int) error) ([]string, [][]string, error) {
	entity := s.entityFor(resource)
	if entity == nil {
		return nil, nil, fmt.Errorf("no export for resource %q", resource)
	}
	return SelectedExportRows(entity, resource, ids, onRow)
}

func (s *BatchExportRowSource) entityFor(resource string) ResourceEntity {
	switch resource {
	case "users":
		return s.users
	case "roles":
		return s.roles
	default:
		return nil
	}
}
