// Package datapermission provider tests (S-09 · GOAL-016 D-002): the module
// registers the policy/scope routes, the data-permission page schema,
// data-permission.read / data-permission.write keys, menu_data_permission
// navigation and the fragment; ScopeFor resolves all/self per policy +
// assignment, and the enforceability gate rejects unwired resources.
package datapermission

import (
	"context"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authsession "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
	datapermissionstore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/datapermission/store"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

func newDataPermissionTestEnv(t *testing.T) (*auth.Authenticator, *Service, *operationlog.Repository) {
	t.Helper()
	hash, err := auth.HashPassword("test-password", 4)
	if err != nil {
		t.Fatalf("hash seed password: %v", err)
	}
	st, err := testsupport.OpenStore(filepath.Join(t.TempDir(), "test.db"), "admin", hash, true)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	repository := authsession.NewRepository(st)
	a := auth.NewWithRepository([]byte("test-secret"), 15*time.Minute, 30*24*time.Hour, repository, false)
	return a, NewService(datapermissionstore.NewRepository(st), []string{"orders"}), operationlog.NewRepository(st)
}

func planWithDataPermission(t *testing.T) kernel.Plan {
	t.Helper()
	registry, err := kernel.NewRegistry(kernel.BuiltinModules())
	if err != nil {
		t.Fatalf("registry: %v", err)
	}
	plan, err := registry.Resolve([]string{
		"core.server-registration", "core.auth-session", "core.manifest-route",
		"core.navigation-capability", "core.schema-render", "core.operationlog",
		"admin.data-permission",
	})
	if err != nil {
		t.Fatalf("resolve plan: %v", err)
	}
	return plan
}

func TestDataPermissionProviderRegistersSurfaces(t *testing.T) {
	a, service, operations := newDataPermissionTestEnv(t)
	provider := New(a, service, operations)
	set, err := kernel.RegisterContributions(context.Background(), planWithDataPermission(t), []kernel.Provider{provider})
	if err != nil {
		t.Fatalf("RegisterContributions: %v", err)
	}
	wantRoutes := []string{
		"GET /api/data-permission/policies", "PATCH /api/data-permission/policies/{resource}",
		"GET /api/data-permission/scopes", "PATCH /api/data-permission/scopes",
	}
	for _, key := range wantRoutes {
		if !slices.ContainsFunc(set.Routes, func(r kernel.RouteContribution) bool { return r.Key == key }) {
			t.Fatalf("route %q missing from provider set", key)
		}
	}
	if !slices.ContainsFunc(set.Pages, func(p kernel.PageContribution) bool { return p.PageID == "data-permission" }) {
		t.Fatalf("page data-permission missing")
	}
	if !slices.ContainsFunc(set.Navigation, func(n kernel.NavigationContribution) bool { return n.NodeID == "menu_data_permission" }) {
		t.Fatalf("navigation menu_data_permission missing")
	}
	for _, perm := range []string{"data-permission.read", "data-permission.write"} {
		if !slices.ContainsFunc(set.Permissions, func(p kernel.PermissionContribution) bool { return p.Permission == perm }) {
			t.Fatalf("permission %q missing", perm)
		}
	}
	if !slices.ContainsFunc(set.Fragments, func(f kernel.FragmentContribution) bool { return f.FragmentID == "data-permission" }) {
		t.Fatalf("fragment data-permission missing")
	}
}

// ScopeFor: unwired resource → nil; wired with policy self → constraint;
// assignment overrides default; disabled policy → nil.
func TestServiceScopeFor(t *testing.T) {
	a, service, _ := newDataPermissionTestEnv(t)
	_ = a

	// Unwired resource: no constraint even with a policy row (PATCH rejects
	// unwired writes, so this can only arise from direct store writes).
	now := time.Now().UTC()
	if err := service.repo.UpsertPolicy("users", "owner_id", datapermissionstore.ScopeSelf, true, now); err != nil {
		t.Fatalf("seed unwired policy: %v", err)
	}
	c, err := service.ScopeFor("u1", "users")
	if err != nil || c != nil {
		t.Fatalf("unwired scope = %v, %v; want nil", c, err)
	}

	// Wired resource: default self → constraint.
	if err := service.UpsertPolicy("orders", "owner_id", datapermissionstore.ScopeSelf, true, now); err != nil {
		t.Fatalf("seed policy: %v", err)
	}
	c, err = service.ScopeFor("u1", "orders")
	if err != nil || c == nil || c.ScopeType != "self" || c.OwnerColumn != "owner_id" || c.ActorID != "u1" {
		t.Fatalf("self scope = %+v, %v", c, err)
	}

	// Assignment "all" overrides the self default → nil.
	if err := service.UpsertAssignments("u1", map[string]string{"orders": "all"}, now); err != nil {
		t.Fatalf("assign all: %v", err)
	}
	c, err = service.ScopeFor("u1", "orders")
	if err != nil || c != nil {
		t.Fatalf("assigned-all scope = %v, %v; want nil", c, err)
	}

	// No policy → nil.
	c, err = service.ScopeFor("u1", "catalog")
	if err != nil || c != nil {
		t.Fatalf("no-policy scope = %v, %v; want nil", c, err)
	}
}

// The enforceability gate lives on the policy write (A-005): unwired
// resources are rejected before any store write.
func TestServiceUpsertPolicyRejectsUnwired(t *testing.T) {
	a, service, _ := newDataPermissionTestEnv(t)
	_ = a
	err := service.UpsertPolicy("users", "owner_id", datapermissionstore.ScopeSelf, true, time.Now().UTC())
	if err != datapermissionstore.ErrNotEnforceable {
		t.Fatalf("unwired upsert err = %v, want ErrNotEnforceable", err)
	}
	if err := service.UpsertPolicy("orders", "owner_id", datapermissionstore.ScopeSelf, true, time.Now().UTC()); err != nil {
		t.Fatalf("wired upsert err = %v", err)
	}
	if err := service.UpsertPolicy("orders", "owner_id", "team", true, time.Now().UTC()); err != datapermissionstore.ErrInvalidScope {
		t.Fatalf("invalid scope err = %v, want ErrInvalidScope", err)
	}
}

var _ handler.RowScopeProvider = (*Service)(nil)
