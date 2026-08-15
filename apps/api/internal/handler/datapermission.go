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

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
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

	// Policies: list registered scope policies.
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
		writeJSON(w, http.StatusOK, map[string]any{"items": policies})
	})))

	// Policies: register/update one resource policy (audited).
	add("PATCH", "/api/data-permission/policies/{resource}", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requirePermission(w, r, "data-permission.write")
		if !ok {
			return
		}
		resource := strings.TrimSpace(r.PathValue("resource"))
		if resource == "" {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SCOPE", "resource is required")
			return
		}
		var body struct {
			OwnerColumn  string `json:"ownerColumn"`
			DefaultScope string `json:"defaultScope"`
			Enabled      *bool  `json:"enabled"`
		}
		r.Body = http.MaxBytesReader(w, r.Body, 4<<10)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SCOPE", "body must be JSON with ownerColumn and defaultScope")
			return
		}
		// default_scope is REQUIRED (A-004 F-001: no implicit default).
		if strings.TrimSpace(body.DefaultScope) == "" {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SCOPE", "defaultScope is required")
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
		recordDataPermissionEvent(operations, user, "data-permission.policy-update",
			`{"resource":`+jsonQuote(resource)+`}`, now)
		writeJSON(w, http.StatusOK, map[string]any{"resource": resource, "ownerColumn": body.OwnerColumn, "defaultScope": body.DefaultScope, "enabled": enabled})
	})))

	// Scopes: list one user's assignments.
	add("GET", "/api/data-permission/scopes", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r, "data-permission.read"); !ok {
			return
		}
		userID := strings.TrimSpace(r.URL.Query().Get("userId"))
		if userID == "" {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SCOPE", "userId is required")
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
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SCOPE", "body must be JSON with userId and scopes")
			return
		}
		if strings.TrimSpace(body.UserID) == "" {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SCOPE", "userId is required")
			return
		}
		if len(body.Scopes) == 0 {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SCOPE", "scopes must not be empty")
			return
		}
		now := time.Now().UTC()
		if err := service.UpsertAssignments(body.UserID, body.Scopes, now); err != nil {
			writeScopeError(w, r, err)
			return
		}
		recordDataPermissionEvent(operations, user, "data-permission.scope-update",
			`{"userId":`+jsonQuote(body.UserID)+`}`, now)
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

// recordDataPermissionEvent writes a data-permission audit row.
func recordDataPermissionEvent(operations operationlog.Recorder, user account.User, event, detail string, now time.Time) {
	if operations == nil {
		return
	}
	_ = operations.RecordOperation(operationlog.Operation{
		ID: newOperationID(), Event: event,
		ActorID: user.ID, ActorName: user.Name, Detail: &detail, CreatedAt: now,
	})
}
