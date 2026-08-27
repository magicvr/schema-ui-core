package config

import (
	"strings"
	"testing"
)

// R6 surface (workspace-017 GOAL-007; contract frozen by GOAL-006 D-002 §2/§4):
// the mail.channel selector, the mail.resend block, and the frozen resolution
// algorithm — explicit selection wins; empty derives (one complete production
// block wins / both = ambiguous fail-closed / none = mock); a touched block
// must be complete in every environment.
func TestMailChannelResolution(t *testing.T) {
	fullSMTP := func() *Config {
		return &Config{AppEnv: "development", MailSMTPHost: "h", MailSMTPUsername: "u",
			MailSMTPPassword: "p", MailSMTPFrom: "f@example.com"}
	}
	fullResend := func() *Config {
		return &Config{AppEnv: "development", MailResendAPIKey: "re-key", MailResendFrom: "f@example.com"}
	}

	t.Run("empty selector with no production blocks resolves mock", func(t *testing.T) {
		c := &Config{AppEnv: "development"}
		ch, err := c.ResolveMailChannel()
		if err != nil || ch != MailChannelMock {
			t.Fatalf("default resolution = %q, %v; want mock", ch, err)
		}
	})

	t.Run("empty selector keeps legacy smtp derivation", func(t *testing.T) {
		c := fullSMTP()
		ch, err := c.ResolveMailChannel()
		if err != nil || ch != MailChannelSMTP {
			t.Fatalf("smtp-only derivation = %q, %v; want smtp (backward compat)", ch, err)
		}
		if err := c.ValidateProd(); err != nil {
			t.Fatalf("existing mail.smtp deployments must keep passing validation: %v", err)
		}
	})

	t.Run("empty selector picks resend when only resend is configured", func(t *testing.T) {
		c := fullResend()
		ch, err := c.ResolveMailChannel()
		if err != nil || ch != MailChannelResend {
			t.Fatalf("resend-only derivation = %q, %v; want resend", ch, err)
		}
	})

	t.Run("both production blocks without selector is ambiguous and fails closed", func(t *testing.T) {
		c := fullSMTP()
		c.MailResendAPIKey = "re-key"
		c.MailResendFrom = "f@example.com"
		if _, err := c.ResolveMailChannel(); err == nil || !strings.Contains(err.Error(), "mail.channel") {
			t.Fatalf("ambiguous config must demand an explicit channel, got %v", err)
		}
		if err := c.ValidateProd(); err == nil {
			t.Fatal("ambiguous config must fail ValidateProd")
		}
	})

	t.Run("explicit channel wins over derivation", func(t *testing.T) {
		c := fullSMTP()
		c.MailResendAPIKey = "re-key"
		c.MailResendFrom = "f@example.com"
		c.MailChannel = MailChannelResend
		ch, err := c.ResolveMailChannel()
		if err != nil || ch != MailChannelResend {
			t.Fatalf("explicit resend over two blocks = %q, %v; want resend", ch, err)
		}
	})

	t.Run("explicit mock ignores configured production blocks", func(t *testing.T) {
		c := fullSMTP()
		c.MailChannel = MailChannelMock
		ch, err := c.ResolveMailChannel()
		if err != nil || ch != MailChannelMock {
			t.Fatalf("explicit mock = %q, %v; want mock", ch, err)
		}
	})

	t.Run("unknown channel value fails closed", func(t *testing.T) {
		c := &Config{AppEnv: "development", MailChannel: "sendgrid"}
		if err := c.ValidateProd(); err == nil || !strings.Contains(err.Error(), "mail.channel") {
			t.Fatalf("unknown mail.channel must fail closed naming the key, got %v", err)
		}
	})

	t.Run("explicit channel with missing block fails closed", func(t *testing.T) {
		t.Run("smtp without block", func(t *testing.T) {
			c := &Config{AppEnv: "development", MailChannel: MailChannelSMTP}
			err := c.ValidateProd()
			if err == nil || !strings.Contains(err.Error(), "mail.smtp") {
				t.Fatalf("channel=smtp without block must fail closed, got %v", err)
			}
		})
		t.Run("resend partial block names first missing key", func(t *testing.T) {
			c := &Config{AppEnv: "development", MailChannel: MailChannelResend, MailResendAPIKey: "re-key"}
			err := c.ValidateProd()
			if err == nil || !strings.Contains(err.Error(), "mail.resend.from") {
				t.Fatalf("partial resend block must name mail.resend.from, got %v", err)
			}
		})
	})

	t.Run("touched resend block requires completeness even under mock", func(t *testing.T) {
		c := &Config{AppEnv: "development", MailChannel: MailChannelMock, MailResendFrom: "f@example.com"}
		err := c.ValidateProd()
		if err == nil || !strings.Contains(err.Error(), "mail.resend.api-key") {
			t.Fatalf("touched-but-incomplete resend must fail closed naming api-key, got %v", err)
		}
	})
}

func TestMailResendKeysLoadAndOverride(t *testing.T) {
	t.Run("yaml resend block parses", func(t *testing.T) {
		y := "app:\n  env: development\nmail:\n  channel: resend\n  resend:\n    api-key: re-yaml\n    from: yaml@example.com\n"
		writeConfig(t, y)
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v", cfg.LoadError)
		}
		if cfg.MailChannel != "resend" || cfg.MailResendAPIKey != "re-yaml" || cfg.MailResendFrom != "yaml@example.com" {
			t.Fatalf("mail.resend block = %+v", cfg)
		}
	})

	t.Run("env overrides win over yaml and normalize case", func(t *testing.T) {
		t.Setenv("MAIL_CHANNEL", " RESEND ")
		t.Setenv("MAIL_RESEND_API_KEY", "env-key")
		t.Setenv("MAIL_RESEND_FROM", "env-from@example.com")
		y := "app:\n  env: development\nmail:\n  channel: mock\n  resend:\n    api-key: yaml-key\n    from: yaml@example.com\n"
		writeConfig(t, y)
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v", cfg.LoadError)
		}
		if cfg.MailChannel != "resend" || cfg.MailResendAPIKey != "env-key" || cfg.MailResendFrom != "env-from@example.com" {
			t.Fatalf("env override = %q/%q/%q", cfg.MailChannel, cfg.MailResendAPIKey, cfg.MailResendFrom)
		}
	})
}

// The resend pairing error names the KEY, never the secret value.
func TestMailResendMisconfigDoesNotLeakSecret(t *testing.T) {
	const secret = "supersecret-resend-value-42"
	c := &Config{AppEnv: "development", MailChannel: MailChannelResend, MailResendAPIKey: secret}
	err := c.ValidateProd()
	if err == nil {
		t.Fatal("missing-from misconfig must fail closed")
	}
	if strings.Contains(err.Error(), secret) {
		t.Fatal("error must never carry the secret value")
	}
	if !strings.Contains(err.Error(), "mail.resend.from") {
		t.Fatalf("error must name the missing key, got %v", err)
	}
}
