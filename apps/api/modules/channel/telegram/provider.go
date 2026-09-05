// Package telegram provides the channel.telegram module surface as a kernel.Provider
// (VP-030 / GOAL-003 R2). It exposes the public webhook route for Telegram Bot updates.
package telegram

import (
	"context"
	"fmt"
	"net/http"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	authsessiondata "github.com/magicvr/schema-ui-core/apps/api/modules/authsession/systemdata"
	telegrammanifest "github.com/magicvr/schema-ui-core/apps/api/modules/channel/telegram/manifest"
	telegrammigration "github.com/magicvr/schema-ui-core/apps/api/modules/channel/telegram/migration"
	telegramschema "github.com/magicvr/schema-ui-core/apps/api/modules/channel/telegram/schema"
)

// ModuleID is the official module identifier for the Telegram channel.
const ModuleID = "channel.telegram"

// Provider implements kernel.Provider for channel.telegram.
type Provider struct {
	webhookHandler  http.Handler
	settingsHandler http.Handler
	leaseHandler    http.Handler
	operatorHandler http.Handler
}

// New constructs the channel.telegram module provider.
func New(webhookHandler http.Handler, handlers ...http.Handler) *Provider {
	var sh http.Handler
	var lh http.Handler
	if len(handlers) > 0 {
		sh = handlers[0]
	}
	if len(handlers) > 1 {
		lh = handlers[1]
	}
	var oh http.Handler
	if len(handlers) > 2 {
		oh = handlers[2]
	}
	return &Provider{
		webhookHandler:  webhookHandler,
		settingsHandler: sh,
		leaseHandler:    lh,
		operatorHandler: oh,
	}
}

func (p *Provider) Descriptor() kernel.Module {
	return kernel.Module{
		ID:             ModuleID,
		Version:        "2.0.0",
		KernelAPIRange: ">=2.0 <3.0",
		// admin.settings is a hard dependency so its settings.read/settings.write
		// permission contributions are always in the same ContributionSet as
		// menu_telegram (R-001 / A-002): nav.Permission must be declared by some
		// provider in the set, and permission keys are globally unique — we reuse
		// admin.settings's rather than mint a new key (no-new-permission red line).
		DependsOn: []string{"core.server-registration", "core.schema-render", "core.navigation-capability", "admin.settings"},
		Requires:  []kernel.Capability{kernel.CapabilityHTTP, kernel.CapabilitySchema, kernel.CapabilityNavigation},
		Contributions: kernel.ContributionKeys{
			Routes: []string{
				"GET /api/channel/telegram/settings",
				"PATCH /api/channel/telegram/settings",
				"POST /api/channel/telegram/lease/acquire",
				"POST /api/channel/telegram/lease/heartbeat",
				"POST /api/channel/telegram/lease/release",
				"POST /api/channel/telegram/webhook",
				"GET /api/channel/telegram/operator/sessions",
				"GET /api/channel/telegram/operator/sessions/{chat_id}/capability",
				"GET /api/channel/telegram/operator/sessions/{chat_id}/messages",
				"POST /api/channel/telegram/operator/sessions/{chat_id}/messages",
				"POST /api/channel/telegram/operator/sessions/{chat_id}/messages/{request_id}/retry",
			},
			Permissions: []string{"telegram.operator.read", "telegram.operator.write"},
			// The settings page remains the sidebar entry; the operator page is
			// an inner route reached from that page (no new permission keys).
			Pages:      []string{"telegram-settings", "telegram-operator"},
			Navigation: []string{"menu_telegram"},
			Fragments:  []string{"telegram-settings"},
		},
	}
}

func (p *Provider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return telegrammigration.Descriptors()
}

func (p *Provider) Register(ctx context.Context, reg kernel.Registrar) error {
	if p.operatorHandler == nil {
		return fmt.Errorf("%s: operator handler is required", ModuleID)
	}
	if p.settingsHandler != nil {
		if err := reg.HTTP(kernel.RouteContribution{
			ContributionIdentity: kernel.ContributionIdentity{
				ModuleID: ModuleID,
				Key:      "GET /api/channel/telegram/settings",
			},
			Method:  "GET",
			Pattern: "/api/channel/telegram/settings",
			Handler: p.settingsHandler,
			Public:  false,
		}); err != nil {
			return err
		}
		if err := reg.HTTP(kernel.RouteContribution{
			ContributionIdentity: kernel.ContributionIdentity{
				ModuleID: ModuleID,
				Key:      "PATCH /api/channel/telegram/settings",
			},
			Method:  "PATCH",
			Pattern: "/api/channel/telegram/settings",
			Handler: p.settingsHandler,
			Public:  false,
		}); err != nil {
			return err
		}
	}
	if p.leaseHandler != nil {
		for _, pattern := range []string{
			"/api/channel/telegram/lease/acquire",
			"/api/channel/telegram/lease/heartbeat",
			"/api/channel/telegram/lease/release",
		} {
			if err := reg.HTTP(kernel.RouteContribution{
				ContributionIdentity: kernel.ContributionIdentity{
					ModuleID: ModuleID,
					Key:      kernel.RouteKey(http.MethodPost, pattern),
				},
				Method:  http.MethodPost,
				Pattern: pattern,
				Handler: p.leaseHandler,
				Public:  false,
			}); err != nil {
				return err
			}
		}
	}
	for _, route := range []struct {
		method  string
		pattern string
	}{
		{http.MethodGet, "/api/channel/telegram/operator/sessions"},
		{http.MethodGet, "/api/channel/telegram/operator/sessions/{chat_id}/capability"},
		{http.MethodGet, "/api/channel/telegram/operator/sessions/{chat_id}/messages"},
		{http.MethodPost, "/api/channel/telegram/operator/sessions/{chat_id}/messages"},
		{http.MethodPost, "/api/channel/telegram/operator/sessions/{chat_id}/messages/{request_id}/retry"},
	} {
		if err := reg.HTTP(kernel.RouteContribution{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: kernel.RouteKey(route.method, route.pattern)},
			Method:               route.method,
			Pattern:              route.pattern,
			Handler:              p.operatorHandler,
			Public:               false,
		}); err != nil {
			return err
		}
	}
	for _, permission := range []kernel.PermissionContribution{
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "telegram.operator.read"},
			Permission:           "telegram.operator.read",
			Resource:             "telegram.operator",
			Action:               "read",
			PolicyID:             authsessiondata.PolicyAdminEditorViewer,
			SystemDataVersion:    authsessiondata.SystemDataVersion,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "telegram.operator.write"},
			Permission:           "telegram.operator.write",
			Resource:             "telegram.operator",
			Action:               "write",
			PolicyID:             authsessiondata.PolicyAdmin,
			SystemDataVersion:    authsessiondata.SystemDataVersion,
		},
	} {
		if err := reg.Authorization(permission); err != nil {
			return err
		}
	}
	if p.webhookHandler != nil {
		if err := reg.HTTP(kernel.RouteContribution{
			ContributionIdentity: kernel.ContributionIdentity{
				ModuleID: ModuleID,
				Key:      "POST /api/channel/telegram/webhook",
			},
			Method:  "POST",
			Pattern: "/api/channel/telegram/webhook",
			Handler: p.webhookHandler,
			Public:  true,
		}); err != nil {
			return err
		}
	}
	// The settings and operator pages are always contributed (schema documents
	// are static); the underlying settings/operator APIs keep their own
	// per-route auth gates.
	if err := reg.Schema(kernel.PageContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "telegram-settings"},
		PageID:               "telegram-settings",
		Resources:            []string{"telegram-settings"},
		Actions:              []string{"list", "update"},
		DataSource:           "/api/channel/telegram/settings",
		Owner:                ModuleID,
		Document:             telegramschema.SchemaDocuments()["telegram-settings"],
	}); err != nil {
		return err
	}
	if err := reg.Schema(kernel.PageContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "telegram-operator"},
		PageID:               "telegram-operator",
		Resources:            []string{"telegram.operator"},
		Actions:              []string{"list", "update"},
		DataSource:           "/api/channel/telegram/operator/sessions",
		Owner:                ModuleID,
		Document:             telegramschema.SchemaDocuments()["telegram-operator"],
	}); err != nil {
		return err
	}
	// R-001 / A-002: menu_telegram rides settings.read — declared by the
	// admin.settings provider (DependsOn above) in the same ContributionSet,
	// so the nav permission reference resolves without minting a new key.
	if err := reg.Navigation(kernel.NavigationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "menu_telegram"},
		NodeID:               "menu_telegram",
		PageID:               "telegram-settings",
		Order:                2,
		Label:                "Telegram channel",
		Visibility:           authsessiondata.PolicyAdmin,
		Permission:           "settings.read",
		SystemDataVersion:    authsessiondata.SystemDataVersion,
	}); err != nil {
		return err
	}
	return reg.Manifest(kernel.FragmentContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "telegram-settings"},
		FragmentID:           "telegram-settings",
		ProtocolVersion:      "2.7",
		RequiredCapabilities: []string{"manifest", "navigation"},
		JSON:                 telegrammanifest.FragmentJSON,
	})
}
