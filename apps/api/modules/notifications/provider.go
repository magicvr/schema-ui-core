// Package notifications provides the admin.notifications module surface as a
// kernel.Provider (F-04 · GOAL-006 D-002). Self-service in-app notifications:
// system-event hooks produce rows; the module owns the endpoints and page.
package notifications

import (
	"context"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	authsession "github.com/magicvr/schema-ui-core/apps/api/modules/authsession"
	authsessiondata "github.com/magicvr/schema-ui-core/apps/api/modules/authsession/systemdata"
	"github.com/magicvr/schema-ui-core/apps/api/modules/notifications/manifest"
	notificationsschema "github.com/magicvr/schema-ui-core/apps/api/modules/notifications/schema"
)

const ModuleID = "admin.notifications"

// Provider implements kernel.Provider for admin.notifications.
type Provider struct {
	a          *auth.Authenticator
	repository *authsession.Repository
}

// New constructs the notifications provider.
func New(a *auth.Authenticator, repository *authsession.Repository) *Provider {
	return &Provider{a: a, repository: repository}
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
				"GET /api/notifications", "POST /api/notifications/{id}/read",
				"POST /api/notifications/read-all", "GET /api/notifications/unread-count",
				"GET /api/notifications/settings", "PATCH /api/notifications/settings",
			},
			Pages:      []string{"notifications"},
			Navigation: []string{"menu_notifications"},
			Fragments:  []string{"notifications"},
		},
	}
}

func (p *Provider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return nil, nil // migrations owned by the migration package provider
}

func (p *Provider) Register(ctx context.Context, reg kernel.Registrar) error {
	for _, route := range handler.NotificationRoutes(p.a, p.repository, ModuleID) {
		if err := reg.HTTP(route); err != nil {
			return err
		}
	}
	if err := reg.Schema(kernel.PageContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "notifications"},
		PageID:               "notifications",
		Resources:            []string{"notifications"},
		Actions:              []string{"list", "update"},
		DataSource:           "/api/notifications",
		Owner:                ModuleID,
		Document:             notificationsschema.SchemaDocuments()["notifications"],
	}); err != nil {
		return err
	}
	if err := reg.Navigation(kernel.NavigationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "menu_notifications"},
		NodeID:               "menu_notifications",
		PageID:               "notifications",
		Order:                2,
		Label:                "Notifications",
		Visibility:           authsessiondata.PolicyAdminEditorViewer,
		Permission:           "",
		SystemDataVersion:    authsessiondata.SystemDataVersion,
	}); err != nil {
		return err
	}
	return reg.Manifest(kernel.FragmentContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "notifications"},
		FragmentID:           "notifications",
		ProtocolVersion:      "2.7",
		RequiredCapabilities: []string{"manifest", "navigation"},
		JSON:                 manifest.FragmentJSON,
	})
}