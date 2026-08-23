package config

import (
	"strings"
	"testing"
)

// R2 surface (workspace-017 GOAL-003 D-001): mail.smtp keys, env injection
// and the fail-closed pairing rules. Untouched = valid embedded default;
// any single key makes host/username/password/from required.
func TestMailSMTPConfig(t *testing.T) {
	t.Run("untouched config passes ValidateProd", func(t *testing.T) {
		c := &Config{AppEnv: "development"}
		if err := c.ValidateProd(); err != nil {
			t.Fatalf("empty mail block must pass, got %v", err)
		}
	})

	t.Run("yaml smtp block parses", func(t *testing.T) {
		y := "app:\n  env: development\nmail:\n  smtp:\n    host: smtp.example.com\n    port: 2465\n    username: api@example.com\n    password: shhh\n    from: no-reply@example.com\n"
		writeConfig(t, y)
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v", cfg.LoadError)
		}
		if cfg.MailSMTPHost != "smtp.example.com" || cfg.MailSMTPPort != 2465 ||
			cfg.MailSMTPUsername != "api@example.com" || cfg.MailSMTPPassword != "shhh" ||
			cfg.MailSMTPFrom != "no-reply@example.com" {
			t.Fatalf("mail.smtp block = %+v", cfg)
		}
	})

	t.Run("env overrides win over yaml", func(t *testing.T) {
		t.Setenv("MAIL_SMTP_HOST", "env-smtp.example.com")
		t.Setenv("MAIL_SMTP_USERNAME", "env-user@example.com")
		t.Setenv("MAIL_SMTP_PASSWORD", "env-secret")
		t.Setenv("MAIL_SMTP_FROM", "env-from@example.com")
		y := "app:\n  env: development\nmail:\n  smtp:\n    host: yaml-smtp.example.com\n    username: u\n    password: p\n    from: f@example.com\n"
		writeConfig(t, y)
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v", cfg.LoadError)
		}
		if cfg.MailSMTPHost != "env-smtp.example.com" || cfg.MailSMTPUsername != "env-user@example.com" ||
			cfg.MailSMTPPassword != "env-secret" || cfg.MailSMTPFrom != "env-from@example.com" {
			t.Fatalf("env override = %s/%s", cfg.MailSMTPHost, cfg.MailSMTPFrom)
		}
	})

	t.Run("invalid explicit port fails closed at load", func(t *testing.T) {
		t.Setenv("MAIL_SMTP_PORT", "not-a-number")
		writeConfig(t, "app:\n  env: development\n")
		cfg := Load()
		if cfg.LoadError == nil || !strings.Contains(cfg.LoadError.Error(), "MAIL_SMTP_PORT") {
			t.Fatalf("non-numeric MAIL_SMTP_PORT must fail closed, got %v", cfg.LoadError)
		}
	})

	t.Run("partial block fails closed naming first missing key", func(t *testing.T) {
		for _, tc := range []struct{ name string; cfg *Config; want string }{
			{"host only", &Config{AppEnv: "development", MailSMTPHost: "h"}, "mail.smtp.username"},
			{"missing from", &Config{AppEnv: "development", MailSMTPHost: "h", MailSMTPUsername: "u", MailSMTPPassword: "p"}, "mail.smtp.from"},
			{"port only", &Config{AppEnv: "development", MailSMTPPort: 465}, "mail.smtp.host"},
		} {
			t.Run(tc.name, func(t *testing.T) {
				err := tc.cfg.ValidateProd()
				if err == nil || !strings.Contains(err.Error(), tc.want) {
					t.Fatalf("ValidateProd must name %s, got %v", tc.want, err)
				}
			})
		}
	})

	t.Run("out-of-range port fails closed", func(t *testing.T) {
		c := &Config{AppEnv: "development",
			MailSMTPHost: "h", MailSMTPPort: 70000, MailSMTPUsername: "u", MailSMTPPassword: "p", MailSMTPFrom: "f@example.com"}
		if err := c.ValidateProd(); err == nil || !strings.Contains(err.Error(), "port") {
			t.Fatalf("out-of-range port must fail closed, got %v", err)
		}
	})

	t.Run("display-name or invalid from fails closed", func(t *testing.T) {
		c := &Config{AppEnv: "development",
			MailSMTPHost: "h", MailSMTPUsername: "u", MailSMTPPassword: "p", MailSMTPFrom: "Ops <ops@example.com>"}
		if err := c.ValidateProd(); err == nil || !strings.Contains(err.Error(), "mail.smtp.from") {
			t.Fatalf("display-name From must fail closed, got %v", err)
		}
	})
}

// The pairing error names the KEY, never the secret value.
func TestMailMisconfigErrorDoesNotLeakSecret(t *testing.T) {
	const secret = "supersecret-mail-value-42"
	c := &Config{AppEnv: "development", MailSMTPHost: "h", MailSMTPUsername: "u", MailSMTPPassword: secret}
	err := c.ValidateProd()
	if err == nil {
		t.Fatal("missing-from misconfig must fail closed")
	}
	if strings.Contains(err.Error(), secret) {
		t.Fatal("error must never carry the secret value")
	}
	if !strings.Contains(err.Error(), "mail.smtp.from") {
		t.Fatalf("error must name the missing key, got %v", err)
	}
}
