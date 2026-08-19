// Recycle bin HTTP surface (S-12 · GOAL-012 D-002 §3): browse, restore and
// purge deleted-row snapshots. The delete hooks live in the resource factory
// (Resource.Trash → TrashRecorder); this file only owns the module routes.
package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
	recyclestore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/recyclebin/store"
)

// RecycleItem is the wire projection of one recycle snapshot.
type RecycleItem struct {
	ID         string
	Resource   string
	ResourceID string
	Payload    map[string]any
	ActorID    string
	ActorName  string
	DeletedAt  time.Time
	RestoredAt time.Time // zero = unrestored
}

// RecycleBinService is the surface the recycle routes consume (satisfied
// structurally by the admin.recycle-bin module Service).
type RecycleBinService interface {
	ListItems(resource, q, sortField, order string, page, pageSize int) ([]RecycleItem, int, error)
	GetItem(id string) (*RecycleItem, error)
	Restore(id string, now time.Time) (map[string]any, error)
	Purge(id string) error
	PurgeAll() (int, error)
}

// RecycleBinRoutes returns the admin.recycle-bin HTTP surface.
func RecycleBinRoutes(a *auth.Authenticator, service RecycleBinService, operations operationlog.Recorder, moduleID string) []kernel.RouteContribution {
	var routes []kernel.RouteContribution

	routes = append(routes, kernel.RouteContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey("GET", "/api/recycle-bin")},
		Method:               "GET",
		Pattern:              "/api/recycle-bin",
		Handler: a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := requirePermission(w, r, "recycle.read"); !ok {
				return
			}
			page, ok := intParam(r.URL.Query().Get("page"), 1)
			if !ok {
				writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PAGE", "page must be a positive integer")
				return
			}
			pageSize, ok := intParam(r.URL.Query().Get("pageSize"), 20)
			if !ok || pageSize > maxPageSize {
				writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PAGE_SIZE", "pageSize must be a positive integer not exceeding 100")
				return
			}
			sortField := r.URL.Query().Get("sort")
			if sortField == "" {
				sortField = "deletedAt"
			}
			if sortField != "deletedAt" && sortField != "resource" && sortField != "actorName" {
				writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SORT_FIELD", "unsupported sort field")
				return
			}
			order := r.URL.Query().Get("order")
			if order == "" {
				order = "desc"
			}
			if order != "asc" && order != "desc" {
				writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SORT_ORDER", "order must be asc or desc")
				return
			}
			items, total, err := service.ListItems(r.URL.Query().Get("resource"), r.URL.Query().Get("q"), sortField, order, page, pageSize)
			if err != nil {
				writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not list recycle items")
				return
			}
			rows := make([]map[string]any, 0, len(items))
			for _, item := range items {
				rows = append(rows, recycleItemToMap(item))
			}
			writeJSON(w, http.StatusOK, resourceList{Items: rows, Total: total, Page: page, PageSize: pageSize})
		})),
	})

	routes = append(routes, kernel.RouteContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey("GET", "/api/recycle-bin/{id}")},
		Method:               "GET",
		Pattern:              "/api/recycle-bin/{id}",
		Handler: a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := requirePermission(w, r, "recycle.read"); !ok {
				return
			}
			item, err := service.GetItem(r.PathValue("id"))
			if err != nil {
				writeRecycleError(w, r, err)
				return
			}
			writeJSON(w, http.StatusOK, recycleItemToMap(*item))
		})),
	})

	routes = append(routes, kernel.RouteContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey("POST", "/api/recycle-bin/{id}/restore")},
		Method:               "POST",
		Pattern:              "/api/recycle-bin/{id}/restore",
		Handler: a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := requirePermission(w, r, "recycle.write")
			if !ok {
				return
			}
			now := time.Now().UTC()
			row, err := service.Restore(r.PathValue("id"), now)
			if err != nil {
				writeRecycleError(w, r, err)
				return
			}
			recordRecycleEvent(operations, operationlog.EventRecycleRestore, user, r.PathValue("id"), now)
			writeJSON(w, http.StatusOK, row)
		})),
	})

	routes = append(routes, kernel.RouteContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey("POST", "/api/recycle-bin/purge-all")},
		Method:               "POST",
		Pattern:              "/api/recycle-bin/purge-all",
		Handler: a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := requirePermission(w, r, "recycle.write")
			if !ok {
				return
			}
			now := time.Now().UTC()
			purged, err := service.PurgeAll()
			if err != nil {
				writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not purge recycle items")
				return
			}
			recordRecycleEvent(operations, operationlog.EventRecyclePurge, user, "purge-all", now)
			writeJSON(w, http.StatusOK, map[string]any{"purged": purged})
		})),
	})

	routes = append(routes, kernel.RouteContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey("DELETE", "/api/recycle-bin/{id}")},
		Method:               "DELETE",
		Pattern:              "/api/recycle-bin/{id}",
		Handler: a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := requirePermission(w, r, "recycle.write")
			if !ok {
				return
			}
			now := time.Now().UTC()
			if err := service.Purge(r.PathValue("id")); err != nil {
				writeRecycleError(w, r, err)
				return
			}
			recordRecycleEvent(operations, operationlog.EventRecyclePurge, user, r.PathValue("id"), now)
			w.WriteHeader(http.StatusNoContent)
		})),
	})
	return routes
}

func recycleItemToMap(item RecycleItem) map[string]any {
	// T-05 (GOAL-013 D-006): deletedAt/restoredAt serialize as UTC ISO-8601
	// (the frozen 3-digit-millisecond shape), matching every other list —
	// never raw Unix seconds. The frontend renders them with formatDisplayTime.
	row := map[string]any{
		"id":         item.ID,
		"resource":   item.Resource,
		"resourceId": item.ResourceID,
		"actorId":    item.ActorID,
		"actorName":  item.ActorName,
		"deletedAt":  item.DeletedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
		"restored":   !item.RestoredAt.IsZero(),
	}
	if !item.RestoredAt.IsZero() {
		row["restoredAt"] = item.RestoredAt.UTC().Format("2006-01-02T15:04:05.000Z07:00")
	}
	if item.Payload != nil {
		row["payload"] = item.Payload
	}
	return row
}

func writeRecycleError(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, recyclestore.ErrItemNotFound) {
		writeLocalizedError(w, r, http.StatusNotFound, "RECYCLE_ITEM_NOT_FOUND", "no recycle item with that id")
		return
	}
	if errors.Is(err, recyclestore.ErrItemAlreadyRestored) {
		writeLocalizedError(w, r, http.StatusConflict, "RECYCLE_ITEM_ALREADY_RESTORED", "recycle item is already restored")
		return
	}
	var de *DomainError
	if errors.As(err, &de) {
		// Localized envelope so the conflict carries the catalog messageKey
		// (grok A-003 F-003): RECYCLE_RESTORE_CONFLICT is cataloged.
		writeLocalizedError(w, r, de.Status, de.Code, de.Message)
		return
	}
	writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not process recycle item")
}

func recordRecycleEvent(operations operationlog.Recorder, event string, user account.User, id string, now time.Time) {
	recordAudit(operations, user, event, id, nil, now.UTC(), nil)
}
