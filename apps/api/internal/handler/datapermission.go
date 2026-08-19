// Data-permission management surface (S-09 · GOAL-016 D-002 §3): scope
// policy registration (owner column + required default scope) and user ×
// resource scope assignments. The enforcement side is the RowScopeProvider
// interface consumed by the resource factory (resources.go); this file only
// owns the module's HTTP contributions.
package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	datapermissionstore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/datapermission/store"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
)

// DataPermissionService is the surface the data-permission routes consume. It
// is satisfied structurally by the admin.data-permission module Service (no
// handler import of the module package — the direction is module → handler).
type DataPermissionService interface {
	// ScopeFor is the RowScopeProvider contract consumed by the resource
	// factory; nil means no constraint.
	ScopeFor(userID, resource string) (*ScopeConstraint, error)
	ListPolicies() ([]datapermissionstore.Policy, error)
	// UpsertPolicy registers a resource policy. default_scope is required;
	// the resource must be wired for scoping (SCOPE_NOT_ENFORCEABLE).
	UpsertPolicy(resource, ownerColumn, defaultScope string, enabled bool, now time.Time) error
	ListAssignments(userID string) ([]datapermissionstore.Assignment, error)
	// UpsertAssignments upserts one user's assignment map (resource → scope).
	UpsertAssignments(userID string, scopes map[string]string, now time.Time) error
}

// DataPermissionRoutes returns the admin.data-permission HTTP surface.
func DataPermissionRoutes(a *auth.Authenticator, service DataPermissionService, operations operationlog.Recorder, moduleID string) []kernel.RouteContribution {
	var routes []kernel.RouteContribution
	add := func(method, pattern string, h http.Handler) {
		routes = append(routes, kernel.RouteContribution{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey(method, pattern)},
			Method:               method,
			Pattern:              pattern,
			Handler:              h,
		})
	}

	// Policies: list registered scope policies (unified list envelope
	// I-010-001 §3 + camelCase row projection the schema-driven table
	// surface renders: resource / ownerColumn / defaultScope / enabled).
	add("GET", "/api/data-permission/policies", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r, "data-permission.read"); !ok {
			return
		}
		policies, err := service.ListPolicies()
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not list scope policies")
			return
		}
		if policies == nil {
			policies = []datapermissionstore.Policy{}
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
		total := len(policies)
		start := (page - 1) * pageSize
		if start > total {
			start = total
		}
		end := start + pageSize
		if end > total {
			end = total
		}
		pagePolicies := policies[start:end]
		items := make([]map[string]any, 0, len(pagePolicies))
		for _, policy := range pagePolicies {
			items = append(items, map[string]any{
				"resource":     policy.Resource,
				"ownerColumn":  policy.OwnerColumn,
				"defaultScope": policy.DefaultScope,
				"enabled":      policy.Enabled,
				"updatedAt":    policy.UpdatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
			})
		}
		writeJSON(w, http.StatusOK, resourceList{Items: items, Total: total, Page: page, PageSize: pageSize})
	})))

	// Policies: register/update one resource policy (audited). The resource
	// rides in the body (not a path slot): the schema-driven form action only
	// binds `{id}` slots, and the modal form owns the resource field
	// (W10 fix, GOAL-011).
	add("PATCH", "/api/data-permission/policies", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requirePermission(w, r, "data-permission.write")
		if !ok {
			return
		}
		var body struct {
			Resource     string `json:"resource"`
			OwnerColumn  string `json:"ownerColumn"`
			DefaultScope string `json:"defaultScope"`
			Enabled      *bool  `json:"enabled"`
		}
		r.Body = http.MaxBytesReader(w, r.Body, 4<<10)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SCOPE_BODY", "body must be JSON with resource, ownerColumn and defaultScope")
			return
		}
		resource := strings.TrimSpace(body.Resource)
		if resource == "" {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SCOPE_BODY", "resource is required")
			return
		}
		// default_scope is REQUIRED (A-004 F-001: no implicit default).
		if strings.TrimSpace(body.DefaultScope) == "" {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SCOPE_BODY", "defaultScope is required")
			return
		}
		enabled := true
		if body.Enabled != nil {
			enabled = *body.Enabled
		}
		now := time.Now().UTC()
		if err := service.UpsertPolicy(resource, body.OwnerColumn, body.DefaultScope, enabled, now); err != nil {
			writeScopeError(w, r, err)
			return
		}
		recordAudit(operations, user, "data-permission.policy-update", "", auditDetail("policy-update", map[string]any{"resource": resource}), now, r.Context())
		writeJSON(w, http.StatusOK, map[string]any{"resource": resource, "ownerColumn": body.OwnerColumn, "defaultScope": body.DefaultScope, "enabled": enabled})
	})))

	// Scopes: list one user's assignments.
	add("GET", "/api/data-permission/scopes", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r, "data-permission.read"); !ok {
			return
		}
		userID := strings.TrimSpace(r.URL.Query().Get("userId"))
		if userID == "" {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SCOPE_BODY", "userId is required")
			return
		}
		assignments, err := service.ListAssignments(userID)
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not list scope assignments")
			return
		}
		if assignments == nil {
			assignments = []datapermissionstore.Assignment{}
		}
		writeJSON(w, http.StatusOK, map[string]any{"userId": userID, "items": assignments})
	})))

	// Scopes: upsert one user's assignments (audited).
	add("PATCH", "/api/data-permission/scopes", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requirePermission(w, r, "data-permission.write")
		if !ok {
			return
		}
		var body struct {
			UserID string            `json:"userId"`
			Scopes map[string]string `json:"scopes"`
		}
		r.Body = http.MaxBytesReader(w, r.Body, 16<<10)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SCOPE_BODY", "body must be JSON with userId and scopes")
			return
		}
		if strings.TrimSpace(body.UserID) == "" {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SCOPE_BODY", "userId is required")
			return
		}
		if len(body.Scopes) == 0 {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SCOPE_BODY", "scopes must not be empty")
			return
		}
		now := time.Now().UTC()
		if err := service.UpsertAssignments(body.UserID, body.Scopes, now); err != nil {
			writeScopeError(w, r, err)
			return
		}
		recordAudit(operations, user, "data-permission.scope-update", "", auditDetail("scope-update", map[string]any{"userId": body.UserID}), now, r.Context())
		writeJSON(w, http.StatusOK, map[string]any{"userId": body.UserID, "updated": len(body.Scopes)})
	})))

	return routes
}

// writeScopeError maps data-permission domain errors to the frozen wire codes.
func writeScopeError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, datapermissionstore.ErrInvalidScope):
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SCOPE", "scope must be all or self")
	case errors.Is(err, datapermissionstore.ErrNotEnforceable):
		writeLocalizedError(w, r, http.StatusBadRequest, "SCOPE_NOT_ENFORCEABLE", "resource is not wired for row-level scoping")
	default:
		writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not update data permission")
	}
}


