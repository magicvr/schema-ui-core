// Package logincaptcha provides the admin.login-captcha module surface as a
// kernel.Provider (S-11 · GOAL-011 D-002): an optional arithmetic challenge
// gate on POST /api/auth/login, the public challenge preflight
// (GET /api/auth/captcha), the settings endpoints (GET/PATCH
// /api/captcha/settings) consumed by the admin.settings security section,
// captcha.read / captcha.write permission keys and the
// captcha.settings-update audit event. The gate is default-off (D-001 §5).
// No page/navigation/fragment: the switch lives in the settings page
// (D-003, user ruling 2026-08-14).
package logincaptcha

import (
	"context"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	authsessiondata "github.com/magicvr/schema-ui-core/apps/api/modules/authsession/systemdata"
	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
)

// ModuleID is the stable admin.login-captcha module identifier.
const ModuleID = "admin.login-captcha"

// Provider implements kernel.Provider for admin.login-captcha.
type Provider struct {
	a          *auth.Authenticator
	service    *Service
	operations operationlog.Recorder
}

// New constructs the captcha provider.
func New(a *auth.Authenticator, service *Service, operations operationlog.Recorder) *Provider {
	return &Provider{a: a, service: service, operations: operations}
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
				"GET /api/auth/captcha",
				"GET /api/captcha/settings", "PATCH /api/captcha/settings",
			},
			Permissions: []string{"captcha.read", "captcha.write"},
		},
	}
}

func (p *Provider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return nil, nil // tables are owned by the logincaptcha/migration provider (0023)
}

func (p *Provider) Register(ctx context.Context, reg kernel.Registrar) error {
	for _, route := range handler.CaptchaRoutes(p.a, p.service, p.operations, ModuleID) {
		if err := reg.HTTP(route); err != nil {
			return err
		}
	}
	for _, permission := range []kernel.PermissionContribution{
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "captcha.read"}, Permission: "captcha.read", Resource: "captcha", Action: "read", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "captcha.write"}, Permission: "captcha.write", Resource: "captcha", Action: "write", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
	} {
		if err := reg.Authorization(permission); err != nil {
			return err
		}
	}
	return nil
}
