package composition

import (
	"context"
	"io"
	"log/slog"
	"path/filepath"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/internal/mail"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/compiled"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
)

// R6 evidence (workspace-017 GOAL-007 D-001 over the GOAL-006 D-002 frozen
// contract): newMailSender resolves ONE kernel.MailSender by mail.channel —
// the mock outbox publisher with nil probe by default (readyz semantics
// unchanged), explicit Resend (nil probe until R8) and explicit SMTP (ESMTP
// Ping probe) adapters. Wiring mirrors the objectStore precedent: direct
// construction inside NewApp, probe passed to RegisterWithMFAProbes.

func testLogger() *slog.Logger { return slog.New(slog.NewTextHandler(io.Discard, nil)) }

func outboxTestStore(t *testing.T) *store.Store {
	t.Helper()
	catalog, err := compiled.PersistenceCatalog()
	if err != nil {
		t.Fatal(err)
	}
	st, err := store.OpenWithCatalog(filepath.Join(t.TempDir(), "wiring.db"), catalog)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return st
}

func TestMailSenderDefaultsToMockOutboxWithoutProbe(t *testing.T) {
	cfg := lifecycleAppConfig(t, "mvp", "127.0.0.1:0")
	st := outboxTestStore(t)
	sender, probe, err := newMailSender(cfg, st, testLogger())
	if err != nil {
		t.Fatalf("newMailSender: %v", err)
	}
	sink, ok := sender.(*mail.OutboxSink)
	if !ok {
		t.Fatalf("unconfigured mail sender = %T, want *mail.OutboxSink", sender)
	}
	if probe != nil {
		t.Fatal("mock default must NOT extend readyz (仅显式生产渠道后 readyz 扩依赖)")
	}
	msg := kernel.MailMessage{To: "user@example.com", Subject: "wiring", TextBody: "hello"}
	if err := sender.Send(context.Background(), msg); err != nil {
		t.Fatalf("Send through the port: %v", err)
	}
	if n, err := sink.Count(context.Background()); err != nil || n != 1 {
		t.Fatalf("outbox count = %d, %v; want 1 persisted record", n, err)
	}
}

func TestMailSenderSelectsExplicitSMTPWithProbe(t *testing.T) {
	cfg := lifecycleAppConfig(t, "mvp", "127.0.0.1:0")
	cfg.MailChannel = config.MailChannelSMTP
	cfg.MailSMTPHost = "smtp.example.com"
	cfg.MailSMTPPort = 465
	cfg.MailSMTPUsername = "api@example.com"
	cfg.MailSMTPPassword = "secret"
	cfg.MailSMTPFrom = "no-reply@example.com"
	st := outboxTestStore(t)
	sender, probe, err := newMailSender(cfg, st, testLogger())
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

func TestMailSenderSelectsResendWithoutProbeUntilR8(t *testing.T) {
	cfg := lifecycleAppConfig(t, "mvp", "127.0.0.1:0")
	cfg.MailResendAPIKey = "re-key"
	cfg.MailResendFrom = "no-reply@example.com"
	st := outboxTestStore(t)
	sender, probe, err := newMailSender(cfg, st, testLogger())
	if err != nil {
		t.Fatalf("newMailSender: %v", err)
	}
	if _, ok := sender.(*mail.Resend); !ok {
		t.Fatalf("resend-configured mail sender = %T, want *mail.Resend", sender)
	}
	if probe != nil {
		t.Fatal("resend must not extend readyz before the R8 production probes")
	}
}

func TestNewMailSenderRejectsIncompleteBlockDefensively(t *testing.T) {
	// ValidateProd already gates partial blocks; this proves the composition
	// constructor itself also fails closed if a hand-built Config bypasses it.
	cfg := lifecycleAppConfig(t, "mvp", "127.0.0.1:0")
	cfg.MailSMTPHost = "smtp.example.com" // everything else missing
	st := outboxTestStore(t)
	if _, _, err := newMailSender(cfg, st, testLogger()); err == nil {
		t.Fatal("partial mail.smtp block must fail closed in composition")
	}

	ambiguous := lifecycleAppConfig(t, "mvp", "127.0.0.1:0")
	ambiguous.MailSMTPHost = "smtp.example.com"
	ambiguous.MailSMTPPort = 465
	ambiguous.MailSMTPUsername = "api@example.com"
	ambiguous.MailSMTPPassword = "secret"
	ambiguous.MailSMTPFrom = "no-reply@example.com"
	ambiguous.MailResendAPIKey = "re-key"
	ambiguous.MailResendFrom = "no-reply@example.com"
	if _, _, err := newMailSender(ambiguous, st, testLogger()); err == nil {
		t.Fatal("two fully configured production channels without an explicit selector must fail closed in composition")
	}
}
