// Package digitaloffer provides the biz.digital-offer module surface as a
// kernel.Provider (VP-031 · GOAL-002 D-002 v1.0.0 §1/§7): digital offers with
// a frozen entitlement form, thin purchases over the wallet money primitives,
// per-subject entitlements, digitaloffer.read / digitaloffer.offer.manage /
// digitaloffer.entitlement.void permission keys, the offers/entitlements pages
// and bizoffer.* audit events. NOT in the mvp/admin default profile sets; the
// module is assembled only when the plan enables biz.digital-offer.
package digitaloffer

import (
	"context"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	authsessiondata "github.com/magicvr/schema-ui-core/apps/api/modules/authsession/systemdata"
	"github.com/magicvr/schema-ui-core/apps/api/modules/digitaloffer/manifest"
	"github.com/magicvr/schema-ui-core/apps/api/modules/digitaloffer/schema"
	digitalofferservice "github.com/magicvr/schema-ui-core/apps/api/modules/digitaloffer/service"
)

// ModuleID is the stable biz.digital-offer module identifier (D-001 I-031-005).
const ModuleID = "biz.digital-offer"

// Provider implements kernel.Provider for biz.digital-offer.
type Provider struct {
	a          *auth.Authenticator
	service    *digitalofferservice.Service
	limiters   kernel.RateLimiterProvider
	dispatcher kernel.TelegramDispatcher // nil = channel surface not assembled
	sender     kernel.TelegramSender
}

// New constructs the digital-offer provider. dispatcher may be nil when the
// channel surface is not assembled; otherwise the §6 commands register on it.
func New(a *auth.Authenticator, service *digitalofferservice.Service, limiters kernel.RateLimiterProvider, dispatcher kernel.TelegramDispatcher, sender kernel.TelegramSender) *Provider {
	return &Provider{a: a, service: service, limiters: limiters, dispatcher: dispatcher, sender: sender}
}

func (p *Provider) Descriptor() kernel.Module {
	return kernel.Module{
		ID:             ModuleID,
		Version:        "1.0.0",
		KernelAPIRange: ">=2.0 <3.0",
		DependsOn:      []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.operationlog"},
		Requires:       kernel.StandardAdminCapabilities(),
		Contributions: kernel.ContributionKeys{
			Routes: []string{
				"GET /api/digitaloffer/offers", "POST /api/digitaloffer/offers",
				"PATCH /api/digitaloffer/offers/{id}",
				"GET /api/digitaloffer/purchases", "GET /api/digitaloffer/entitlements",
				"POST /api/digitaloffer/entitlements/{id}/void",
				// D-002 §5.3: public C-end catalog (unauthenticated read).
				"GET /api/biz/offers",
			},
			Pages:       []string{"digitaloffer-offers", "digitaloffer-entitlements", "digitaloffer-purchases"},
			Navigation:  []string{"menu_digitaloffer_offers", "menu_digitaloffer_entitlements", "menu_digitaloffer_purchases"},
			Permissions: []string{"digitaloffer.read", "digitaloffer.offer.manage", "digitaloffer.entitlement.void"},
			Fragments:   []string{"digitaloffer"},
		},
	}
}

func (p *Provider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return nil, nil // tables are owned by the digitaloffer/migration provider (0070)
}

func (p *Provider) Register(ctx context.Context, reg kernel.Registrar) error {
	if p.dispatcher != nil {
		if err := p.service.RegisterTelegram(p.dispatcher, p.sender); err != nil {
			return err
		}
	}
	for _, route := range handler.DigitalOfferRoutes(p.a, p.service, ModuleID, p.limiters) {
		if err := reg.HTTP(route); err != nil {
			return err
		}
	}
	for _, route := range handler.DigitalOfferPublicRoutes(p.a, p.service, ModuleID, p.limiters) {
		if err := reg.HTTP(route); err != nil {
			return err
		}
	}
	if err := reg.Schema(kernel.PageContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "digitaloffer-offers"},
		PageID:               "digitaloffer-offers",
		Resources:            []string{"digitaloffer"},
		Actions:              []string{"list", "create", "update"},
		DataSource:           "/api/digitaloffer/offers",
		Owner:                ModuleID,
		Document:             schema.SchemaDocuments()["digitaloffer-offers"],
	}); err != nil {
		return err
	}
	if err := reg.Schema(kernel.PageContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "digitaloffer-entitlements"},
		PageID:               "digitaloffer-entitlements",
		Resources:            []string{"digitaloffer"},
		Actions:              []string{"list", "update"},
		DataSource:           "/api/digitaloffer/entitlements",
		Owner:                ModuleID,
		Document:             schema.SchemaDocuments()["digitaloffer-entitlements"],
	}); err != nil {
		return err
	}
	if err := reg.Schema(kernel.PageContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "digitaloffer-purchases"},
		PageID:               "digitaloffer-purchases",
		Resources:            []string{"digitaloffer"},
		Actions:              []string{"list"},
		DataSource:           "/api/digitaloffer/purchases",
		Owner:                ModuleID,
		Document:             schema.SchemaDocuments()["digitaloffer-purchases"],
	}); err != nil {
		return err
	}
	for _, permission := range []kernel.PermissionContribution{
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "digitaloffer.read"}, Permission: "digitaloffer.read", Resource: "digitaloffer", Action: "read", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "digitaloffer.offer.manage"}, Permission: "digitaloffer.offer.manage", Resource: "digitaloffer", Action: "offer.manage", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "digitaloffer.entitlement.void"}, Permission: "digitaloffer.entitlement.void", Resource: "digitaloffer", Action: "entitlement.void", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
	} {
		if err := reg.Authorization(permission); err != nil {
			return err
		}
	}
	if err := reg.Navigation(kernel.NavigationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "menu_digitaloffer_offers"},
		NodeID:               "menu_digitaloffer_offers",
		PageID:               "digitaloffer-offers",
		Order:                20,
		Label:                "Digital products",
		Group:                &kernel.NavigationGroup{Key: "commerce", Order: 50, Label: "Commerce", LabelKey: "manifest.nav.group.commerce", Secondary: "COMMERCE"},
		Visibility:           authsessiondata.PolicyAdmin,
		Permission:           "digitaloffer.read",
		SystemDataVersion:    authsessiondata.SystemDataVersion,
	}); err != nil {
		return err
	}
	if err := reg.Navigation(kernel.NavigationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "menu_digitaloffer_entitlements"},
		NodeID:               "menu_digitaloffer_entitlements",
		PageID:               "digitaloffer-entitlements",
		Order:                21,
		Label:                "Digital entitlements",
		Group:                &kernel.NavigationGroup{Key: "commerce", Order: 50, Label: "Commerce", LabelKey: "manifest.nav.group.commerce", Secondary: "COMMERCE"},
		Visibility:           authsessiondata.PolicyAdmin,
		Permission:           "digitaloffer.read",
		SystemDataVersion:    authsessiondata.SystemDataVersion,
	}); err != nil {
		return err
	}
	if err := reg.Navigation(kernel.NavigationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "menu_digitaloffer_purchases"},
		NodeID:               "menu_digitaloffer_purchases",
		PageID:               "digitaloffer-purchases",
		Order:                22,
		Label:                "Digital orders",
		Group:                &kernel.NavigationGroup{Key: "commerce", Order: 50, Label: "Commerce", LabelKey: "manifest.nav.group.commerce", Secondary: "COMMERCE"},
		Visibility:           authsessiondata.PolicyAdmin,
		Permission:           "digitaloffer.read",
		SystemDataVersion:    authsessiondata.SystemDataVersion,
	}); err != nil {
		return err
	}
	return reg.Manifest(kernel.FragmentContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "digitaloffer"},
		FragmentID:           "digitaloffer",
		ProtocolVersion:      "2.7",
		RequiredCapabilities: []string{"manifest", "navigation"},
		JSON:                 manifest.FragmentJSON,
	})
}
