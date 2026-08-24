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

// R7 evidence (workspace-017 GOAL-008 D-001 over Root D-007 / GOAL-006
// D-002): newMailRuntime builds THE switching kernel.MailSender — mock outbox
// publisher by default with nil probe (readyz semantics unchanged), boot-time
// SMTP contributing its ESMTP Ping probe exactly as frozen in R4, and Resend
// without a probe until R8.

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

func TestMailRuntimeDefaultsToSwitcherOverMockWithoutProbe(t *testing.T) {
	cfg := lifecycleAppConfig(t, "mvp", "127.0.0.1:0")
	st := outboxTestStore(t)
	sender, probe, err := newMailRuntime(cfg, st, testLogger())
	if err != nil {
		t.Fatalf("newMailRuntime: %v", err)
	}
	if probe != nil {
		t.Fatal("mock default must NOT extend readyz (仅显式生产渠道后 readyz 扩依赖)")
	}
	msg := kernel.MailMessage{To: "user@example.com", Subject: "wiring", TextBody: "hello"}
	if err := sender.Send(context.Background(), msg); err != nil {
		t.Fatalf("Send through the port: %v", err)
	}
	view, err := sender.PublicView()
	if err != nil || view.Channel != config.MailChannelMock {
		t.Fatalf("runtime channel = %+v, %v; want mock", view, err)
	}
	reader := mail.NewOutboxSink(st, 0)
	if n, err := reader.Count(context.Background()); err != nil || n != 1 {
		t.Fatalf("outbox count = %d, %v; want 1 persisted record", n, err)
	}
}

func TestMailRuntimeBootSMTPContributesProbe(t *testing.T) {
	cfg := lifecycleAppConfig(t, "mvp", "127.0.0.1:0")
	cfg.MailChannel = config.MailChannelSMTP
	cfg.MailSMTPHost = "smtp.example.com"
	cfg.MailSMTPPort = 465
	cfg.MailSMTPUsername = "api@example.com"
	cfg.MailSMTPPassword = "secret"
	cfg.MailSMTPFrom = "no-reply@example.com"
	st := outboxTestStore(t)
	_, probe, err := newMailRuntime(cfg, st, testLogger())
	if err != nil {
		t.Fatalf("newMailRuntime: %v", err)
	}
	if probe == nil {
		t.Fatal("boot SMTP channel must contribute a readyz probe")
	}
}

func TestMailRuntimeResendHasNoProbeUntilR8(t *testing.T) {
	cfg := lifecycleAppConfig(t, "mvp", "127.0.0.1:0")
	cfg.MailResendAPIKey = "re-key"
	cfg.MailResendFrom = "no-reply@example.com"
	st := outboxTestStore(t)
	sender, probe, err := newMailRuntime(cfg, st, testLogger())
	if err != nil {
		t.Fatalf("newMailRuntime: %v", err)
	}
	if probe != nil {
		t.Fatal("resend must not extend readyz before the R8 production probes")
	}
	view, err := sender.PublicView()
	if err != nil || view.Channel != config.MailChannelResend || !view.Secrets.ResendAPIKeySet {
		t.Fatalf("seeded runtime view = %+v, %v; want resend with api key set", view, err)
	}
	// The PublicView struct carries presence booleans only — no secret field
	// exists to leak (verified by construction).
}

func TestNewMailRuntimeRejectsIncompleteBlockDefensively(t *testing.T) {
	// ValidateProd already gates partial blocks; this proves the composition
	// constructor itself also fails closed if a hand-built Config bypasses it.
	cfg := lifecycleAppConfig(t, "mvp", "127.0.0.1:0")
	cfg.MailSMTPHost = "smtp.example.com" // everything else missing
	st := outboxTestStore(t)
	if _, _, err := newMailRuntime(cfg, st, testLogger()); err == nil {
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
	if _, _, err := newMailRuntime(ambiguous, st, testLogger()); err == nil {
		t.Fatal("two fully configured production channels without an explicit selector must fail closed in composition")
	}
}
