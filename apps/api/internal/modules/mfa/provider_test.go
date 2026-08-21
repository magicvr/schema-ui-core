// Package mfa provider tests (S-10 · GOAL-017 D-002): the module registers
// the verify/self-service/admin-reset routes and users.mfa-reset; no
// page/navigation/fragment (personal-center block, D-002 §4).
package mfa

import (
	"context"
	"errors"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authsession "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/mfa/store"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

func newMfaTestEnv(t *testing.T) (*auth.Authenticator, *Service, *authsession.Repository, *operationlog.Repository) {
	t.Helper()
	hash, err := auth.HashPassword("test-password", 4)
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	st, err := testsupport.OpenStore(filepath.Join(t.TempDir(), "test.db"), "admin", hash, true)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	repository := authsession.NewRepository(st)
	a := auth.NewWithRepository([]byte("test-secret"), 15*time.Minute, 30*24*time.Hour, repository, false)
	return a, NewService(store.NewRepository(st), []byte("test-secret")), repository, operationlog.NewRepository(st)
}

func planWithMFA(t *testing.T) kernel.Plan {
	t.Helper()
	registry, err := kernel.NewRegistry(kernel.BuiltinModules())
	if err != nil {
		t.Fatalf("registry: %v", err)
	}
	plan, err := registry.Resolve([]string{
		"core.server-registration", "core.auth-session", "core.manifest-route",
		"core.navigation-capability", "core.schema-render", "core.operationlog",
		"admin.mfa",
	})
	if err != nil {
		t.Fatalf("resolve plan: %v", err)
	}
	return plan
}

func TestMFAProviderRegistersSurfaces(t *testing.T) {
	a, service, revoker, operations := newMfaTestEnv(t)
	provider := New(a, service, operations, revoker)
	set, err := kernel.RegisterContributions(context.Background(), planWithMFA(t), []kernel.Provider{provider})
	if err != nil {
		t.Fatalf("RegisterContributions: %v", err)
	}
	wantRoutes := []string{
		"POST /api/auth/mfa/verify",
		"GET /api/mfa/status", "POST /api/mfa/enroll", "POST /api/mfa/confirm",
		"POST /api/mfa/disable", "POST /api/mfa/recovery/rotate",
		"POST /api/users/{id}/mfa/reset",
	}
	for _, key := range wantRoutes {
		if !slices.ContainsFunc(set.Routes, func(r kernel.RouteContribution) bool { return r.Key == key }) {
			t.Fatalf("route %q missing from provider set", key)
		}
	}
	if len(set.Pages) != 0 || len(set.Navigation) != 0 || len(set.Fragments) != 0 {
		t.Fatalf("mfa must not contribute pages/navigation/fragments (D-002 §4)")
	}
	if !slices.ContainsFunc(set.Permissions, func(p kernel.PermissionContribution) bool { return p.Permission == "users.mfa-reset" }) {
		t.Fatalf("permission users.mfa-reset missing")
	}
}

// The auth core gate stays byte-identical with a nil enforcer and issues no
// tokens when Required returns true.
func TestAuthenticatorMFAgate(t *testing.T) {
	a, service, _, _ := newMfaTestEnv(t)
	_ = service

	// nil enforcer → Login is byte-identical (existing behavior).
	a.SetMFAEnforcer(nil)
	access, refresh, user, err := a.Login("admin", "test-password", time.Now().UTC())
	if err != nil || access == "" || refresh == "" || user.ID == "" {
		t.Fatalf("nil-enforcer login = %q %q %+v %v", access, refresh, user, err)
	}

	// Enforcer requiring the admin → MFARequiredError, no tokens.
	now := time.Now().UTC()
	secret, _, _, err := service.Enroll("user-admin", "Admin", now)
	if err != nil {
		t.Fatal(err)
	}
	if err := service.Confirm("user-admin", mustTotp(t, secret, now), now); err != nil {
		t.Fatal(err)
	}
	a.SetMFAEnforcer(service)
	access, refresh, user, err = a.Login("admin", "test-password", now)
	if err == nil {
		t.Fatalf("mfa-gated login must fail without second factor")
	}
	var mfaReq *auth.MFARequiredError
	if !errors.As(err, &mfaReq) || mfaReq.UserID != "user-admin" {
		t.Fatalf("err = %v, want MFARequiredError{user-admin}", err)
	}
	if access != "" || refresh != "" {
		t.Fatalf("mfa-gated login issued tokens")
	}

	// IssueTokensFor completes the login for the verified user.
	access, refresh, user, err = a.IssueTokensFor("user-admin", now)
	if err != nil || access == "" || refresh == "" {
		t.Fatalf("IssueTokensFor = %q %q %+v %v", access, refresh, user, err)
	}
}

func mustTotp(t *testing.T, secret string, now time.Time) string {
	t.Helper()
	code, err := totpCode(secret, now.Unix()/totpPeriodSeconds)
	if err != nil {
		t.Fatal(err)
	}
	return code
}
