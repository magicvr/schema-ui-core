// Package mfa provides the admin.mfa module surface as a kernel.Provider
// (S-10 · GOAL-017 D-002): the MFAVerifier login gate, the public verify
// endpoint, the self-service /api/mfa/* endpoints, the users.mfa-reset
// permission and the mfa.* audit events. No page/navigation/fragment — the
// management surface lives in the account page (personal center) and the
// users page row action (D-002 §4). The MFAVerifier is fed to the auth
// handler via RegisterWithMFA so the login contract stays byte-identical
// when the module is disabled.
package mfa

import (
	"context"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	authsessiondata "github.com/magicvr/schema-ui-core/apps/api/modules/authsession/systemdata"
	"github.com/magicvr/schema-ui-core/apps/api/modules/mfa/store"
	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
)

// ModuleID is the stable admin.mfa module identifier.
const ModuleID = "admin.mfa"

// Provider implements kernel.Provider for admin.mfa.
type Provider struct {
	a          *auth.Authenticator
	service    *Service
	operations operationlog.Recorder
	revoker    handler.SessionRevoker
	limiters   kernel.RateLimiterProvider
}

// New constructs the MFA provider. revoker (the auth-session repository)
// powers the disable/admin-reset session invalidation (A-004 F-002).
func New(a *auth.Authenticator, service *Service, operations operationlog.Recorder, revoker handler.SessionRevoker, limiters kernel.RateLimiterProvider) *Provider {
	return &Provider{a: a, service: service, operations: operations, revoker: revoker, limiters: limiters}
}

func (p *Provider) Descriptor() kernel.Module {
	return kernel.Module{
		ID:             ModuleID,
		Version:        "2.0.0",
		KernelAPIRange: ">=2.0 <3.0",
		DependsOn:      []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.operationlog"},
		Requires:       kernel.StandardAdminCapabilities(),
		Contributions: kernel.ContributionKeys{
			Routes: []string{
				"POST /api/auth/mfa/verify",
				"GET /api/mfa/status", "POST /api/mfa/enroll", "POST /api/mfa/confirm",
				"POST /api/mfa/disable", "POST /api/mfa/recovery/rotate",
				"POST /api/users/{id}/mfa/reset",
			},
			Permissions: []string{"users.mfa-reset"},
		},
	}
}

func (p *Provider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return nil, nil // tables are owned by the mfa/migration provider (0029)
}

func (p *Provider) Register(ctx context.Context, reg kernel.Registrar) error {
	for _, route := range handler.MFARoutes(p.a, p.service, p.operations, p.revoker, ModuleID, p.limiters) {
		if err := reg.HTTP(route); err != nil {
			return err
		}
	}
	for _, permission := range []kernel.PermissionContribution{
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "users.mfa-reset"}, Permission: "users.mfa-reset", Resource: "users", Action: "mfa-reset", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
	} {
		if err := reg.Authorization(permission); err != nil {
			return err
		}
	}
	return nil
}

var (
	_ handler.MFAVerifier   = (*Service)(nil)
	_ handler.MFASelfService = (*Service)(nil)
	_ auth.MFAEnforcer      = (*Service)(nil)
)

// keep store import referenced in provider tests via NewService call sites.
var _ = store.ErrNotFound
