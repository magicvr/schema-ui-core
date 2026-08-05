// Operations (activity) read-only resource (GOAL-013): list + detail over the
// existing operation_log table. No create/update/delete — append is internal.
package handler

import (
	"errors"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
)

func operationsResource(repository operationlog.Reader) Resource {
	return Resource{
		ID:              "operations",
		Path:            "/api/operations",
		Listable:        true,
		ReadOnly:        true,
		SortFields:      []string{"createdAt", "event", "actorName"},
		QSearch:         true,
		Entity:          &operationsEntity{repository: repository},
		PermissionRead:  "operations.read",
		PermissionWrite: "operations.write", // unused when ReadOnly
		NotFoundCode:    "OPERATION_NOT_FOUND",
	}
}

type operationsEntity struct {
	repository operationlog.Reader
}

func operationToMap(op operationlog.Operation) map[string]any {
	row := map[string]any{
		"id":        op.ID,
		"event":     op.Event,
		"actorId":   op.ActorID,
		"actorName": op.ActorName,
		"createdAt": op.CreatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
	}
	if op.RecordID != nil {
		row["recordId"] = *op.RecordID
	} else {
		row["recordId"] = ""
	}
	if op.Detail != nil {
		row["detail"] = *op.Detail
	} else {
		row["detail"] = ""
	}
	return row
}

func (e *operationsEntity) List(f resourceFilter) ([]map[string]any, int, error) {
	items, total, err := e.repository.ListOperationsFiltered(operationlog.OperationFilter{
		Q: f.Q, Sort: f.Sort, Order: f.Order, Page: f.Page, PageSize: f.PageSize,
	})
	if err != nil {
		return nil, 0, err
	}
	out := make([]map[string]any, 0, len(items))
	for _, it := range items {
		out = append(out, operationToMap(it))
	}
	return out, total, nil
}

func (e *operationsEntity) Get(id string) (map[string]any, error) {
	op, err := e.repository.GetOperation(id)
	if err != nil {
		if errors.Is(err, operationlog.ErrNotFound) {
			return nil, errResourceNotFound
		}
		return nil, err
	}
	return operationToMap(*op), nil
}

func (e *operationsEntity) Create(map[string]any, string, time.Time, account.User) (map[string]any, error) {
	return nil, errReadOnlyResource
}

func (e *operationsEntity) Update(string, map[string]any, time.Time, account.User) (map[string]any, error) {
	return nil, errReadOnlyResource
}

func (e *operationsEntity) Delete(string, account.User) error {
	return errReadOnlyResource
}

// RegisterActivity is removed in R6 C6.1: the Activity module mounts its
// read-only operations surface via the module provider (ResourceRoutes +
// kernel.RegisterContributions); the handler test environment mounts
// ResourceRoutes directly. Operation-log writes remain owned by core.operationlog.

// OperationsResource exposes the read-only operations resource descriptor to
// module providers (R4 C4.2).
func OperationsResource(repository operationlog.Reader) Resource {
	return operationsResource(repository)
}
