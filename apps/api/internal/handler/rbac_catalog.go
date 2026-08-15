// RBAC catalog surface (W11 · U-02): read-only catalogs consumed by the
// schema-driven role forms' dynamic option sources (optionsSource). They expose
// the same reconciled tables the roles resource validates grants against, so
// the UI can only offer keys/ids the backend would accept. Local extension —
// the upstream protocol pin is unchanged.
package handler

import (
	"net/http"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authsession "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
)

// RBACCatalog is the persistence surface the catalog routes read from
// (satisfied structurally by the auth-session repository).
type RBACCatalog interface {
	ListPermissionCatalog() ([]authsession.PermissionCatalogEntry, error)
	ListMenuItemCatalog() ([]authsession.MenuItemCatalogEntry, error)
}

// CatalogRoutes returns the RBAC catalog HTTP surface (auth + roles.read).
func CatalogRoutes(a *auth.Authenticator, catalog RBACCatalog, moduleID string) []kernel.RouteContribution {
	return []kernel.RouteContribution{
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: "GET /api/permissions"},
			Method:               http.MethodGet,
			Pattern:              "/api/permissions",
			Handler: a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if _, ok := requirePermission(w, r, "roles.read"); !ok {
					return
				}
				entries, err := catalog.ListPermissionCatalog()
				if err != nil {
					writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not read permission catalog")
					return
				}
				items := make([]map[string]any, 0, len(entries))
				for _, entry := range entries {
					items = append(items, map[string]any{"key": entry.Key, "description": entry.Description})
				}
				writeJSON(w, http.StatusOK, map[string]any{"items": items})
			})),
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: "GET /api/menu-items"},
			Method:               http.MethodGet,
			Pattern:              "/api/menu-items",
			Handler: a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if _, ok := requirePermission(w, r, "roles.read"); !ok {
					return
				}
				entries, err := catalog.ListMenuItemCatalog()
				if err != nil {
					writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not read menu item catalog")
					return
				}
				items := make([]map[string]any, 0, len(entries))
				for _, entry := range entries {
					items = append(items, map[string]any{"id": entry.ID, "pageRef": entry.PageRef, "label": entry.Label})
				}
				writeJSON(w, http.StatusOK, map[string]any{"items": items})
			})),
		},
	}
}
