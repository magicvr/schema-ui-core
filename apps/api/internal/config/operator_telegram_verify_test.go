package config

import (
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// TestCustomModulesResolveWithTelegram guards the VP-030 enablement mechanism:
// a custom app.modules list that includes channel.telegram must resolve to a
// plan containing it AND its hard dependency admin.settings (whose
// settings.read permission menu_telegram references — GOAL-006 R-001).
// It uses an inline module list rather than the operator configs/config.yaml,
// because that file is an operator/deployment choice (its default profile is
// mvp and may legitimately differ across environments).
func TestCustomModulesResolveWithTelegram(t *testing.T) {
	modules := []string{
		"core.server-registration", "core.auth-session", "core.schema-render",
		"core.manifest-route", "core.navigation-capability", "core.operationlog",
		"admin.settings", "channel.telegram",
	}
	reg, err := kernel.NewRegistry(kernel.BuiltinModules())
	if err != nil {
		t.Fatalf("NewRegistry: %v", err)
	}
	plan, err := reg.Resolve(modules)
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}
	if !plan.HasModule("channel.telegram") {
		t.Fatalf("resolved plan missing channel.telegram")
	}
	if !plan.HasModule("admin.settings") {
		t.Fatalf("resolved plan missing admin.settings (channel.telegram dependency)")
	}
}
