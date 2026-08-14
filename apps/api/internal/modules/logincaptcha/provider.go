// Package logincaptcha provides the admin.login-captcha module surface as a
// kernel.Provider (S-11 · GOAL-011 D-002): an optional arithmetic challenge
// gate on POST /api/auth/login, the public challenge preflight
// (GET /api/auth/captcha), admin settings (GET/PATCH /api/captcha/settings),
// captcha.read / captcha.write permission keys and the captcha.settings-update
// audit event. The gate is default-off (D-001 §5).
package logincaptcha

import (
	"context"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authsessiondata "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession/systemdata"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/logincaptcha/manifest"
	captchaschema "github.com/magicvr/schema-ui-core/apps/api/internal/modules/logincaptcha/schema"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
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
			Pages:       []string{"captcha"},
			Navigation:  []string{"menu_captcha"},
			Permissions: []string{"captcha.read", "captcha.write"},
			Fragments:   []string{"captcha"},
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
	for _, pageID := range []string{"captcha"} {
		if err := reg.Schema(kernel.PageContribution{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: pageID},
			PageID:               pageID,
			Resources:            []string{"captcha"},
			Actions:              []string{"read", "update"},
			DataSource:           "/api/captcha/settings",
			Owner:                ModuleID,
			Document:             captchaschema.SchemaDocuments()[pageID],
		}); err != nil {
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
	if err := reg.Navigation(kernel.NavigationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "menu_captcha"},
		NodeID:               "menu_captcha",
		PageID:               "captcha",
		Order:                7,
		Label:                "Login captcha",
		Visibility:           authsessiondata.PolicyAdmin,
		Permission:           "captcha.read",
		SystemDataVersion:    authsessiondata.SystemDataVersion,
	}); err != nil {
		return err
	}
	return reg.Manifest(kernel.FragmentContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "captcha"},
		FragmentID:           "captcha",
		ProtocolVersion:      "2.7",
		RequiredCapabilities: []string{"manifest", "navigation"},
		JSON:                 manifest.FragmentJSON,
	})
}
