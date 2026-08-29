package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	authsession "github.com/magicvr/schema-ui-core/apps/api/modules/authsession"
)

// fakeRBACCatalog is an in-memory RBACCatalog test double.
type fakeRBACCatalog struct {
	permissions []authsession.PermissionCatalogEntry
	menuItems   []authsession.MenuItemCatalogEntry
}

func (f *fakeRBACCatalog) ListPermissionCatalog() ([]authsession.PermissionCatalogEntry, error) {
	return f.permissions, nil
}

func (f *fakeRBACCatalog) ListMenuItemCatalog() ([]authsession.MenuItemCatalogEntry, error) {
	return f.menuItems, nil
}

// W11 · U-02: GET /api/permissions + GET /api/menu-items expose the RBAC
// catalogs to roles.read holders; anonymous/unauthorized are rejected.
func TestRBACCatalogRoutes(t *testing.T) {
	env := newAuthTestEnv(t)
	catalog := &fakeRBACCatalog{
		permissions: []authsession.PermissionCatalogEntry{
			{Key: "users.read", Description: "users read gate"},
			{Key: "tasks.write", Description: "tasks write gate"},
		},
		menuItems: []authsession.MenuItemCatalogEntry{
			{ID: "menu-users", PageRef: "users", Label: "Users"},
			{ID: "menu-scheduled-tasks", PageRef: "scheduled-tasks", Label: "Scheduled tasks"},
		},
	}
	for _, route := range CatalogRoutes(env.a, catalog, "admin.roles") {
		env.mux.Handle(route.Method+" "+route.Pattern, route.Handler)
	}

	// Anonymous → 401.
	req := httptest.NewRequest(http.MethodGet, "/api/permissions", nil)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	expectError(t, rr, http.StatusUnauthorized, "UNAUTHENTICATED")

	// Viewer (roles.read granted by default policy) → catalog rows.
	env.addUser(t, "viewer1", "viewer-password", []string{"viewer"})
	token := env.login(t, "viewer1", "viewer-password")

	req = bearer(t, token, http.MethodGet, "/api/permissions", "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK || !strings.Contains(rr.Body.String(), "tasks.write") {
		t.Fatalf("permissions catalog = %d: %s", rr.Code, rr.Body.String())
	}

	req = bearer(t, token, http.MethodGet, "/api/menu-items", "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK || !strings.Contains(rr.Body.String(), "menu-scheduled-tasks") || !strings.Contains(rr.Body.String(), "Scheduled tasks") {
		t.Fatalf("menu items catalog = %d: %s", rr.Code, rr.Body.String())
	}

	// F-007 (grok A-002): an authenticated user WITHOUT roles.read → 403.
	env.addUser(t, "auditor1", "auditor-password", []string{"custom-auditor"})
	auditorToken := env.login(t, "auditor1", "auditor-password")
	req = bearer(t, auditorToken, http.MethodGet, "/api/permissions", "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	expectError(t, rr, http.StatusForbidden, "FORBIDDEN")
	req = bearer(t, auditorToken, http.MethodGet, "/api/menu-items", "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	expectError(t, rr, http.StatusForbidden, "FORBIDDEN")
}
