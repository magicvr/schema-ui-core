package composition

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"go.uber.org/fx"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/internal/mail"
)

// R3 evidence (workspace-017 GOAL-004 D-001): the fx container resolves ONE
// kernel.MailSender — capture sink by default, explicit SMTP adapter when
// configured. Unconfigured startup stays green (lifecycle dual-profile tests
// keep passing) while tests can read the last captured message back.

func TestMailSenderWiringDefaultsToCaptureSink(t *testing.T) {
	cfg := lifecycleAppConfig(t, "mvp", "127.0.0.1:0")
	var sender kernel.MailSender
	app := fx.New(
		fx.Supply(cfg),
		fx.Supply(slog.New(slog.NewTextHandler(io.Discard, nil))),
		fx.Provide(newMailSender),
		fx.Populate(&sender),
	)
	startCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Start(startCtx); err != nil {
		t.Fatalf("resolve mail sender from container: %v", err)
	}
	sink, ok := sender.(*mail.CaptureSink)
	if !ok {
		t.Fatalf("unconfigured mail sender = %T, want *mail.CaptureSink", sender)
	}
	msg := kernel.MailMessage{To: "user@example.com", Subject: "wiring", TextBody: "hello"}
	if err := sender.Send(context.Background(), msg); err != nil {
		t.Fatalf("Send through the port: %v", err)
	}
	got, has := sink.Last()
	if !has || got != msg {
		t.Fatalf("Last = %+v has=%v, want %+v", got, has, msg)
	}
}

func TestMailSenderWiringSelectsExplicitSMTP(t *testing.T) {
	cfg := lifecycleAppConfig(t, "mvp", "127.0.0.1:0")
	cfg.MailSMTPHost = "smtp.example.com"
	cfg.MailSMTPPort = 465
	cfg.MailSMTPUsername = "api@example.com"
	cfg.MailSMTPPassword = "secret"
	cfg.MailSMTPFrom = "no-reply@example.com"
	var sender kernel.MailSender
	app := fx.New(
		fx.Supply(cfg),
		fx.Supply(slog.New(slog.NewTextHandler(io.Discard, nil))),
		fx.Provide(newMailSender),
		fx.Populate(&sender),
	)
	startCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Start(startCtx); err != nil {
		t.Fatalf("resolve explicit smtp sender: %v", err)
	}
	if _, ok := sender.(*mail.SMTP); !ok {
		t.Fatalf("explicitly configured mail sender = %T, want *mail.SMTP", sender)
	}
}

func TestNewMailSenderRejectsIncompleteBlockDefensively(t *testing.T) {
	// ValidateProd already gates partial blocks; this proves the composition
	// constructor itself also fails closed if a hand-built Config bypasses it.
	cfg := lifecycleAppConfig(t, "mvp", "127.0.0.1:0")
	cfg.MailSMTPHost = "smtp.example.com" // everything else missing
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	if _, err := newMailSender(cfg, logger); err == nil {
		t.Fatal("partial mail.smtp block must fail closed in composition")
	}
}
