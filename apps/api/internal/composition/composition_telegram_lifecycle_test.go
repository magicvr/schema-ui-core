package composition

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"go.uber.org/fx"

	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

func TestTelegramFxShutdownDrainsPollingReceiver(t *testing.T) {
	pollStarted := make(chan struct{})
	pollCanceled := make(chan struct{})
	var startOnce, cancelOnce sync.Once
	client := &http.Client{Transport: telegramLifecycleRoundTripFunc(func(r *http.Request) (*http.Response, error) {
		switch r.URL.Path {
		case "/botlive-bot-token/getMe":
			return telegramLifecycleJSONResponse(`{"ok":true,"result":{"id":201,"is_bot":true,"username":"shutdown_bot"}}`), nil
		case "/botlive-bot-token/deleteWebhook":
			return telegramLifecycleJSONResponse(`{"ok":true,"result":true}`), nil
		case "/botlive-bot-token/getUpdates":
			startOnce.Do(func() { close(pollStarted) })
			<-r.Context().Done()
			cancelOnce.Do(func() { close(pollCanceled) })
			return nil, r.Context().Err()
		default:
			return telegramLifecycleJSONResponse(`{"ok":false,"description":"unexpected path"}`), nil
		}
	})}

	cfg := &config.Config{
		ProfileName:           string(kernel.ProfileCustom),
		ModulesEnabled:        []string{"core.server-registration", "core.auth-session", "core.schema-render", "core.manifest-route", "core.navigation-capability", "core.operationlog", "admin.settings", "channel.telegram"},
		TelegramBotToken:      "live-bot-token",
		TelegramWebhookSecret: "correct-secret",
		TelegramMasterKey:     "test-master-key",
		HTTPAddr:              "127.0.0.1:0",
		DBPath:                filepath.Join(t.TempDir(), "telegram_fx_shutdown.db"),
	}

	var injected *TelegramRuntime
	app, err := newAppWithOptions(
		cfg,
		"test-secret",
		"hash",
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		fx.Populate(&injected),
		fx.Supply(&telegramRuntimeOptions{APIBaseURL: "https://telegram.test", HTTPClient: client}),
	)
	if err != nil {
		t.Fatalf("newAppWithOptions: %v", err)
	}
	t.Cleanup(func() { _ = app.Stop(context.Background()) })

	startCtx, startCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer startCancel()
	if err := app.Start(startCtx); err != nil {
		t.Fatalf("app.Start: %v", err)
	}
	if injected == nil || injected.Connection == nil {
		t.Fatal("expected Fx-injected Telegram connection manager")
	}
	if err := injected.Connection.AcquireLease(context.Background(), "fx-shutdown-session"); err != nil {
		t.Fatalf("AcquireLease: %v", err)
	}
	select {
	case <-pollStarted:
	case <-time.After(time.Second):
		t.Fatal("polling receiver did not start after acquiring the lease")
	}

	stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer stopCancel()
	if err := app.Stop(stopCtx); err != nil {
		t.Fatalf("app.Stop: %v", err)
	}
	select {
	case <-pollCanceled:
	case <-time.After(time.Second):
		t.Fatal("Fx shutdown returned without canceling the Telegram long poll")
	}
	status := injected.Connection.Status()
	if status.State != telegramConnectionStateIdle || status.Receiver != telegramReceiverNone {
		t.Fatalf("Telegram status after Fx shutdown = %+v", status)
	}
}

const (
	telegramConnectionStateIdle = "idle"
	telegramReceiverNone        = "none"
)

type telegramLifecycleRoundTripFunc func(*http.Request) (*http.Response, error)

func (f telegramLifecycleRoundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func telegramLifecycleJSONResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}
