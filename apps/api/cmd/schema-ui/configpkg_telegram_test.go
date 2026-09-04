package main

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestExportTelegramSecretsAreExcluded(t *testing.T) {
	source := writeTemp(t, "telegram-secrets.yaml", `app:
  env: development
telegram:
  bot_token: ${TELEGRAM_BOT_TOKEN}
  webhook_secret: ${TELEGRAM_WEBHOOK_SECRET}
  mode: webhook
  webhook_public_base_url: https://console.example
`)
	tree, exclude, err := buildExportTree(source)
	if err != nil {
		t.Fatalf("buildExportTree: %v", err)
	}
	if tree.Telegram.Mode != "webhook" || tree.Telegram.WebhookPublicBaseURL != "https://console.example" {
		t.Fatalf("exported Telegram non-secret settings = %+v", tree.Telegram)
	}
	want := map[string]string{
		"telegram.bot_token":      "TELEGRAM_BOT_TOKEN",
		"telegram.webhook_secret": "TELEGRAM_WEBHOOK_SECRET",
	}
	got := make(map[string]string)
	for _, entry := range exclude {
		got[entry.Key] = entry.Env
	}
	for key, env := range want {
		if got[key] != env {
			t.Errorf("sensitive export entry %q = %q, want env %q", key, got[key], env)
		}
	}
	raw, err := yaml.Marshal(&tree)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), "TELEGRAM_BOT_TOKEN") || strings.Contains(string(raw), "TELEGRAM_WEBHOOK_SECRET") ||
		strings.Contains(string(raw), "bot_token") || strings.Contains(string(raw), "webhook_secret") {
		t.Fatalf("Telegram secret value or key leaked into exported config tree:\n%s", raw)
	}
}
