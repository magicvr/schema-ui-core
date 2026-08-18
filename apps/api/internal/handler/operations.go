// Operations (activity) read-only resource (GOAL-013): list + detail over the
// existing operation_log table. No create/update/delete — append is internal.
package handler

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
)

func operationsResource(repository operationlog.Reader) Resource {
	return Resource{
		ID:         "operations",
		Path:       "/api/operations",
		Listable:   true,
		ReadOnly:   true,
		SortFields: []string{"createdAt", "event", "actorName"},
		QSearch:    true,
		// W14 F-03: structured audit filters (event / actor / time range).
		ExtraQuery:      []string{"event", "actorName", "from", "to"},
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
	row["correlationId"] = op.CorrelationID
	return row
}

func (e *operationsEntity) List(f resourceFilter) ([]map[string]any, int, error) {
	filter, err := operationFilterFromResource(f)
	if err != nil {
		return nil, 0, err
	}
	items, total, err := e.repository.ListOperationsFiltered(filter)
	if err != nil {
		return nil, 0, err
	}
	out := make([]map[string]any, 0, len(items))
	for _, it := range items {
		out = append(out, operationToMap(it))
	}
	return out, total, nil
}

// operationFilterFromResource converts a validated resource list filter into
// the operationlog filter, parsing the F-03 date-range query params.
func operationFilterFromResource(f resourceFilter) (operationlog.OperationFilter, error) {
	filter := operationlog.OperationFilter{
		Q: f.Q, Event: f.Extra["event"], ActorName: f.Extra["actorName"],
		Sort: f.Sort, Order: f.Order, Page: f.Page, PageSize: f.PageSize,
	}
	if raw := f.Extra["from"]; raw != "" {
		from, parseErr := parseOperationTime(raw, false)
		if parseErr != nil {
			return operationlog.OperationFilter{}, parseErr
		}
		filter.From = &from
	}
	if raw := f.Extra["to"]; raw != "" {
		to, parseErr := parseOperationTime(raw, true)
		if parseErr != nil {
			return operationlog.OperationFilter{}, parseErr
		}
		filter.To = &to
	}
	return filter, nil
}

// parseOperationTime accepts YYYY-MM-DD (inclusive day boundaries) or RFC3339.
func parseOperationTime(raw string, endOfDay bool) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}, nil
	}
	if parsed, err := time.Parse("2006-01-02", raw); err == nil {
		if endOfDay {
			return parsed.Add(24*time.Hour - time.Nanosecond).UTC(), nil
		}
		return parsed.UTC(), nil
	}
	if parsed, err := time.Parse(time.RFC3339, raw); err == nil {
		return parsed.UTC(), nil
	}
	return time.Time{}, &DomainError{Status: http.StatusBadRequest, Code: "INVALID_DATE_FILTER", Message: "from/to must be YYYY-MM-DD or RFC3339"}
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
