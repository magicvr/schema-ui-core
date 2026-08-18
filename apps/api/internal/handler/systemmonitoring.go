// System monitoring surface (S-03 · GOAL-009 D-002): a read-only admin
// monitoring module — a process-level status summary (in-process liveness /
// readiness probes, version, uptime, module set, DB size) and a recent-events
// view over the operation_log (best-effort "error/event" surface, v1 boundary
// D-002 `1). No write paths, no audit events, no migrations.
package handler

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
	"github.com/magicvr/schema-ui-core/apps/api/pkg/version"
)

// MonitoringStatusRow is one status summary row. The endpoint serves it as a
// single-row list envelope ({items,total,page,pageSize}) because the Host
// loads statCard dataSource values through fetchResourceList (A-003 F-001);
// the flat object alone would fail closed on the page.
type MonitoringStatusRow struct {
	Status           string   `json:"status"`
	AvailabilityMode string   `json:"availabilityMode"`
	Ready            bool     `json:"ready"`
	Version          string   `json:"version"`
	Commit           string   `json:"commit"`
	UptimeSeconds    int64    `json:"uptimeSeconds"`
	ModuleCount      int      `json:"moduleCount"`
	Modules          []string `json:"modules"`
	DBSizeBytes      int64    `json:"dbSizeBytes"`
}

// monitoringEntity adapts the operation-log reader to the read-only errors
// resource (same shape as the activity operations resource, monitoring.read).
type monitoringEntity struct {
	repository operationlog.Reader
}

func (e *monitoringEntity) List(f resourceFilter) ([]map[string]any, int, error) {
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

func (e *monitoringEntity) Get(id string) (map[string]any, error) {
	op, err := e.repository.GetOperation(id)
	if err != nil {
		// GetOperation wraps the sentinel (fmt.Errorf %w) — use errors.Is
		// (A-003 F-002), matching the activity operations entity.
		if errors.Is(err, operationlog.ErrNotFound) {
			return nil, errResourceNotFound
		}
		return nil, err
	}
	return operationToMap(*op), nil
}

func (e *monitoringEntity) Create(map[string]any, string, time.Time, account.User) (map[string]any, error) {
	return nil, errReadOnlyResource
}

func (e *monitoringEntity) Update(string, map[string]any, time.Time, account.User) (map[string]any, error) {
	return nil, errReadOnlyResource
}

func (e *monitoringEntity) Delete(string, account.User) error {
	return errReadOnlyResource
}

// MonitoringRoutes returns the admin.system-monitoring HTTP surface.
func MonitoringRoutes(a *auth.Authenticator, st *store.Store, plan kernel.Plan, ready func() bool, dbPath string, startTime time.Time, repository operationlog.Reader, availabilityMode, moduleID string) []kernel.RouteContribution {
	routes := ResourceRoutes(a, Resource{
		ID:              "monitoring-errors",
		Path:            "/api/system-monitoring/errors",
		Listable:        true,
		ReadOnly:        true,
		SortFields:      []string{"createdAt", "event", "actorName"},
		QSearch:         true,
		Entity:          &monitoringEntity{repository: repository},
		PermissionRead:  "monitoring.read",
		PermissionWrite: "monitoring.write", // unused when ReadOnly
		NotFoundCode:    "OPERATION_NOT_FOUND",
	}, moduleID)

	routes = append(routes, kernel.RouteContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey("GET", "/api/system-monitoring/status")},
		Method:               "GET",
		Pattern:              "/api/system-monitoring/status",
		Handler: a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := requirePermission(w, r, "monitoring.read"); !ok {
				return
			}
			// In-process probes (D-002 `1): store ping + module-graph readiness
			// gate — equivalent to /healthz + /readyz without the HTTP round trip.
			status := "ok"
			readyOK := true
			ctx, cancel := context.WithTimeout(r.Context(), time.Second)
			defer cancel()
			if err := st.Ping(ctx); err != nil {
				status = "unavailable"
				readyOK = false
			}
			if ready != nil && !ready() {
				status = "not-ready"
				readyOK = false
			}
			dbSize := int64(0)
			if info, err := os.Stat(dbPath); err == nil {
				dbSize = info.Size()
			}
			row := MonitoringStatusRow{
				Status:           status,
				AvailabilityMode: availabilityMode,
				Ready:            readyOK,
				Version:          version.Version,
				Commit:           version.Commit,
				UptimeSeconds:    int64(time.Since(startTime).Seconds()),
				ModuleCount:      len(plan.IDs()),
				Modules:          plan.IDs(),
				DBSizeBytes:      dbSize,
			}
			writeJSON(w, http.StatusOK, resourceList{Items: []map[string]any{statusRowToMap(row)}, Total: 1, Page: 1, PageSize: 1})
		})),
	})
	return routes
}

func statusRowToMap(row MonitoringStatusRow) map[string]any {
	return map[string]any{
		"status":           row.Status,
		"availabilityMode": row.AvailabilityMode,
		"ready":            row.Ready,
		"version":          row.Version,
		"commit":           row.Commit,
		"uptimeSeconds":    row.UptimeSeconds,
		"moduleCount":      row.ModuleCount,
		"modules":          row.Modules,
		"dbSizeBytes":      row.DBSizeBytes,
	}
}
