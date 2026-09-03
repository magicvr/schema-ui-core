// Package telegram provides the channel.telegram module surface as a kernel.Provider
// (VP-030 / GOAL-003 R2). It exposes the public webhook route for Telegram Bot updates.
package telegram

import (
	"context"
	"net/http"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// ModuleID is the official module identifier for the Telegram channel.
const ModuleID = "channel.telegram"

// Provider implements kernel.Provider for channel.telegram.
type Provider struct {
	webhookHandler  http.Handler
	settingsHandler http.Handler
}

// New constructs the channel.telegram module provider.
func New(webhookHandler http.Handler, settingsHandlers ...http.Handler) *Provider {
	var sh http.Handler
	if len(settingsHandlers) > 0 {
		sh = settingsHandlers[0]
	}
	return &Provider{
		webhookHandler:  webhookHandler,
		settingsHandler: sh,
	}
}

func (p *Provider) Descriptor() kernel.Module {
	return kernel.Module{
		ID:             ModuleID,
		Version:        "2.0.0",
		KernelAPIRange: ">=2.0 <3.0",
		DependsOn:      []string{"core.server-registration"},
		Requires:       []kernel.Capability{kernel.CapabilityHTTP},
		Contributions: kernel.ContributionKeys{
			Routes: []string{
				"GET /api/channel/telegram/settings",
				"PATCH /api/channel/telegram/settings",
				"POST /api/channel/telegram/webhook",
			},
		},
	}
}

func (p *Provider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return nil, nil
}

func (p *Provider) Register(ctx context.Context, reg kernel.Registrar) error {
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
	return nil
}
