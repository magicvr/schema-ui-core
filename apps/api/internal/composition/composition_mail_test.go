package composition

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/internal/mail"
)

// R4 evidence (workspace-017 GOAL-005 D-001): newMailSender resolves ONE
// kernel.MailSender plus an optional readyz probe — capture sink with nil
// probe by default (readyz semantics unchanged), explicit SMTP adapter whose
// ESMTP Ping joins readyz only when explicitly configured. Wiring mirrors the
// objectStore precedent: direct construction inside NewApp, probe passed to
// RegisterWithMFAProbes.

func testLogger() *slog.Logger { return slog.New(slog.NewTextHandler(io.Discard, nil)) }

func TestMailSenderDefaultsToCaptureSinkWithoutProbe(t *testing.T) {
	cfg := lifecycleAppConfig(t, "mvp", "127.0.0.1:0")
	sender, probe, err := newMailSender(cfg, testLogger())
	if err != nil {
		t.Fatalf("newMailSender: %v", err)
	}
	sink, ok := sender.(*mail.CaptureSink)
	if !ok {
		t.Fatalf("unconfigured mail sender = %T, want *mail.CaptureSink", sender)
	}
	if probe != nil {
		t.Fatal("capture default must NOT extend readyz (仅显式配置后 readyz 扩依赖)")
	}
	msg := kernel.MailMessage{To: "user@example.com", Subject: "wiring", TextBody: "hello"}
	if err := sender.Send(context.Background(), msg); err != nil {
		t.Fatalf("Send through the port: %v", err)
	}
	if got, has := sink.Last(); !has || got != msg {
		t.Fatalf("Last = %+v has=%v, want %+v", got, has, msg)
	}
}

func TestMailSenderSelectsExplicitSMTPWithProbe(t *testing.T) {
	cfg := lifecycleAppConfig(t, "mvp", "127.0.0.1:0")
	cfg.MailSMTPHost = "smtp.example.com"
	cfg.MailSMTPPort = 465
	cfg.MailSMTPUsername = "api@example.com"
	cfg.MailSMTPPassword = "secret"
	cfg.MailSMTPFrom = "no-reply@example.com"
	sender, probe, err := newMailSender(cfg, testLogger())
	if err != nil {
		t.Fatalf("newMailSender: %v", err)
	}
	if _, ok := sender.(*mail.SMTP); !ok {
		t.Fatalf("explicitly configured mail sender = %T, want *mail.SMTP", sender)
	}
	if probe == nil {
		t.Fatal("explicitly configured endpoint must contribute a readyz probe")
	}
}

func TestNewMailSenderRejectsIncompleteBlockDefensively(t *testing.T) {
	// ValidateProd already gates partial blocks; this proves the composition
	// constructor itself also fails closed if a hand-built Config bypasses it.
	cfg := lifecycleAppConfig(t, "mvp", "127.0.0.1:0")
	cfg.MailSMTPHost = "smtp.example.com" // everything else missing
	if _, _, err := newMailSender(cfg, testLogger()); err == nil {
		t.Fatal("partial mail.smtp block must fail closed in composition")
	}
}
