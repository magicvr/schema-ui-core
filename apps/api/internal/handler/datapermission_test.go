package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	datapermissionstore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/datapermission/store"
)

// fakeDataPermissionService is an in-memory DataPermissionService test double
// (the real store-backed service is covered by the module tests).
type fakeDataPermissionService struct {
	policies    map[string]datapermissionstore.Policy
	assignments map[string]map[string]string // userID → resource → scope
	enforceable map[string]bool
}

func newFakeDataPermissionService(enforceable ...string) *fakeDataPermissionService {
	set := map[string]bool{}
	for _, r := range enforceable {
		set[r] = true
	}
	return &fakeDataPermissionService{
		policies:    map[string]datapermissionstore.Policy{},
		assignments: map[string]map[string]string{},
		enforceable: set,
	}
}

func (s *fakeDataPermissionService) ScopeFor(userID, resource string) (*ScopeConstraint, error) {
	if !s.enforceable[resource] {
		return nil, nil
	}
	p, ok := s.policies[resource]
	if !ok || !p.Enabled {
		return nil, nil
	}
	effective := p.DefaultScope
	if userScopes, ok := s.assignments[userID]; ok {
		if v, ok := userScopes[resource]; ok {
			effective = v
		}
	}
	if effective != datapermissionstore.ScopeSelf {
		return nil, nil
	}
	return &ScopeConstraint{Resource: resource, ScopeType: datapermissionstore.ScopeSelf, OwnerColumn: p.OwnerColumn, ActorID: userID}, nil
}

func (s *fakeDataPermissionService) ListPolicies() ([]datapermissionstore.Policy, error) {
	items := make([]datapermissionstore.Policy, 0, len(s.policies))
	for _, p := range s.policies {
		items = append(items, p)
	}
	return items, nil
}

func (s *fakeDataPermissionService) UpsertPolicy(resource, ownerColumn, defaultScope string, enabled bool, now time.Time) error {
	if !s.enforceable[resource] {
		return datapermissionstore.ErrNotEnforceable
	}
	if !datapermissionstore.ValidScope(defaultScope) {
		return datapermissionstore.ErrInvalidScope
	}
	s.policies[resource] = datapermissionstore.Policy{Resource: resource, OwnerColumn: ownerColumn, DefaultScope: defaultScope, Enabled: enabled, UpdatedAt: now}
	return nil
}

func (s *fakeDataPermissionService) ListAssignments(userID string) ([]datapermissionstore.Assignment, error) {
	items := []datapermissionstore.Assignment{}
	for resource, scope := range s.assignments[userID] {
		items = append(items, datapermissionstore.Assignment{UserID: userID, Resource: resource, ScopeType: scope})
	}
	return items, nil
}

func (s *fakeDataPermissionService) UpsertAssignments(userID string, scopes map[string]string, now time.Time) error {
	for resource := range scopes {
		if !s.enforceable[resource] {
			return datapermissionstore.ErrNotEnforceable
		}
		if !datapermissionstore.ValidScope(scopes[resource]) {
			return datapermissionstore.ErrInvalidScope
		}
	}
	if s.assignments[userID] == nil {
		s.assignments[userID] = map[string]string{}
	}
	for resource, scope := range scopes {
		s.assignments[userID][resource] = scope
	}
	return nil
}

func mountDataPermissionRoutes(t *testing.T, env *authTestEnv, service DataPermissionService) {
	t.Helper()
	mountRoutes := func(routes []kernel.RouteContribution) {
		for _, r := range routes {
			env.mux.Handle(r.Method+" "+r.Pattern, r.Handler)
		}
	}
	mountRoutes(DataPermissionRoutes(env.a, service, env.operations, "admin.data-permission"))
}

// W14 F-05 (GOAL-019): policies list honours real pagination instead of a fake
// pageSize.
func TestDataPermissionPoliciesPagination(t *testing.T) {
	env := newAuthTestEnv(t)
	service := newFakeDataPermissionService("orders", "invoices")
	_ = service.UpsertPolicy("orders", "owner_id", "self", true, time.Now())
	_ = service.UpsertPolicy("invoices", "owner_id", "all", true, time.Now())
	mountDataPermissionRoutes(t, env, service)
	token := adminToken(t, env)

	req := bearer(t, token, http.MethodGet, "/api/data-permission/policies?pageSize=1", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("list policies = %d", rr.Code)
	}
	var list struct {
		Items []map[string]any `json:"items"`
		Total int              `json:"total"`
		Page  int              `json:"page"`
		PageSize int           `json:"pageSize"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&list); err != nil {
		t.Fatal(err)
	}
	if list.Total != 2 || len(list.Items) != 1 || list.PageSize != 1 || list.Page != 1 {
		t.Fatalf("policies pagination = %+v, want total 2, one item, pageSize 1", list)
	}
}

// S-09 (GOAL-016 D-002 §3): gates — anonymous 401, authenticated without the
// key 403.
func TestDataPermissionRoutesGates(t *testing.T) {
	env := newAuthTestEnv(t)
	mountDataPermissionRoutes(t, env, newFakeDataPermissionService())

	req := httptest.NewRequest(http.MethodGet, "/api/data-permission/policies", nil)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	expectError(t, rr, http.StatusUnauthorized, "UNAUTHENTICATED")

	// An editor without data-permission.read → 403.
	env.addUser(t, "editor1", "editor-password", []string{"editor"})
	editorToken := env.login(t, "editor1", "editor-password")
	req = bearer(t, editorToken, http.MethodGet, "/api/data-permission/policies", "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	expectError(t, rr, http.StatusForbidden, "FORBIDDEN")
}

// S-09 (GOAL-016 D-002 §3): admin can register a policy on a wired resource
// (audited), while unwired resources are rejected with SCOPE_NOT_ENFORCEABLE.
func TestDataPermissionRoutesPolicyLifecycle(t *testing.T) {
	env := newAuthTestEnv(t)
	service := newFakeDataPermissionService("orders")
	mountDataPermissionRoutes(t, env, service)
	token := adminToken(t, env)

	// Unwired resource → 400 SCOPE_NOT_ENFORCEABLE.
	req := bearer(t, token, http.MethodPatch, "/api/data-permission/policies", `{"resource":"users","ownerColumn":"owner_id","defaultScope":"self"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	expectError(t, rr, http.StatusBadRequest, "SCOPE_NOT_ENFORCEABLE")

	// Missing defaultScope → 400 INVALID_SCOPE_BODY.
	req = bearer(t, token, http.MethodPatch, "/api/data-permission/policies", `{"resource":"orders","ownerColumn":"owner_id"}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	expectError(t, rr, http.StatusBadRequest, "INVALID_SCOPE_BODY")

	// Wired resource → 200 + list shows it.
	req = bearer(t, token, http.MethodPatch, "/api/data-permission/policies", `{"resource":"orders","ownerColumn":"owner_id","defaultScope":"self"}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("register policy = %d: %s", rr.Code, rr.Body.String())
	}
	req = bearer(t, token, http.MethodGet, "/api/data-permission/policies", "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("list policies = %d", rr.Code)
	}
	var list struct {
		Items []datapermissionstore.Policy `json:"items"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&list); err != nil || len(list.Items) != 1 || list.Items[0].Resource != "orders" {
		t.Fatalf("policies = %+v, err = %v", list, err)
	}

	// Scopes: invalid scope rejected; valid assignment accepted.
	req = bearer(t, token, http.MethodPatch, "/api/data-permission/scopes", `{"userId":"u1","scopes":{"orders":"team"}}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	expectError(t, rr, http.StatusBadRequest, "INVALID_SCOPE")
	req = bearer(t, token, http.MethodPatch, "/api/data-permission/scopes", `{"userId":"u1","scopes":{"orders":"all"}}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("assign scopes = %d: %s", rr.Code, rr.Body.String())
	}
	req = bearer(t, token, http.MethodGet, "/api/data-permission/scopes?userId=u1", "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK || !strings.Contains(rr.Body.String(), "orders") {
		t.Fatalf("list scopes = %d: %s", rr.Code, rr.Body.String())
	}
}
