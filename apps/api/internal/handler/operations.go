// Operations (activity) read-only resource (GOAL-013): list + detail over the
// existing operation_log table. No create/update/delete — append is internal.
package handler

import (
	"net/http"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
)

func operationsResource(st *store.Store) Resource {
	return Resource{
		ID:              "operations",
		Path:            "/api/operations",
		Listable:        true,
		ReadOnly:        true,
		SortFields:      []string{"createdAt", "event", "actorName"},
		QSearch:         true,
		Entity:          &operationsEntity{st: st},
		PermissionRead:  "operations.read",
		PermissionWrite: "operations.write", // unused when ReadOnly
		NotFoundCode:    "OPERATION_NOT_FOUND",
	}
}

type operationsEntity struct {
	st *store.Store
}

func operationToMap(op store.Operation) map[string]any {
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
	items, total, err := e.st.ListOperationsFiltered(store.OperationFilter{
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
	op, err := e.st.GetOperation(id)
	if err != nil {
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

// registerOperations mounts the read-only operations resource. Kept as a thin
// wrapper so the Activity module can own its registration boundary.
func registerOperations(mux *http.ServeMux, a *auth.Authenticator, st *store.Store) {
	registerResource(mux, a, operationsResource(st))
}

// RegisterActivity exposes the Activity module registration adapter to the
// composition root. Operation-log writes remain available without this route.
func RegisterActivity(mux *http.ServeMux, a *auth.Authenticator, st *store.Store) {
	registerOperations(mux, a, st)
}

// OperationsResource exposes the read-only operations resource descriptor to
// module providers (R4 C4.2).
func OperationsResource(st *store.Store) Resource {
	return operationsResource(st)
}
