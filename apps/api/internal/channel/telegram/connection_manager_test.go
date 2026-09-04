package telegram

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

func TestConnectionManager_WebhookEstablishment(t *testing.T) {
	rm, err := NewRuntimeManagerWithSettings("bot-token", "webhook-secret", TelegramModeWebhook, "https://console.example", nil, testMasterKey())
	if err != nil {
		t.Fatal(err)
	}
	var mu sync.Mutex
	var paths []string
	var webhookPayload setWebhookPayload
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		paths = append(paths, r.URL.Path)
		mu.Unlock()
		if r.URL.Path == "/botbot-token/setWebhook" {
			body, _ := io.ReadAll(r.Body)
			if err := json.Unmarshal(body, &webhookPayload); err != nil {
				t.Errorf("decode setWebhook: %v", err)
			}
			_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
			return
		}
		if r.URL.Path == "/botbot-token/getMe" {
			_, _ = w.Write([]byte(`{"ok":true,"result":{"id":7,"is_bot":true,"username":"console_bot"}}`))
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	manager := NewConnectionManager(rm, NewDispatcher(), NewBotAPIClient(rm, server.Client(), server.URL), NewPollingBotAPIClient(rm, server.Client(), server.URL), nil)
	if err := manager.Start(context.Background()); err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer manager.Stop(context.Background())
	status := manager.Status()
	if status.State != ConnectionStateRunning || status.Receiver != ReceiverWebhook || status.BotID != 7 || status.BotUsername != "console_bot" {
		t.Fatalf("webhook status = %+v", status)
	}
	mu.Lock()
	gotPaths := append([]string(nil), paths...)
	mu.Unlock()
	if !reflect.DeepEqual(gotPaths, []string{"/botbot-token/getMe", "/botbot-token/setWebhook"}) {
		t.Fatalf("webhook call order = %v", gotPaths)
	}
	if webhookPayload.URL != "https://console.example/api/channel/telegram/webhook" || webhookPayload.SecretToken != "webhook-secret" {
		t.Fatalf("setWebhook payload = %+v", webhookPayload)
	}
}

func TestConnectionManager_PollingWithoutDemandEstablishesAndStaysIdle(t *testing.T) {
	rm := newTestRuntimeManager(t, "bot-token", "", nil)
	var getUpdatesCalls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/botbot-token/getMe":
			_, _ = w.Write([]byte(`{"ok":true,"result":{"id":8,"is_bot":true,"username":"idle_bot"}}`))
		case "/botbot-token/deleteWebhook":
			_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
		case "/botbot-token/getUpdates":
			getUpdatesCalls.Add(1)
			t.Errorf("getUpdates must not run without a lease or business handler")
			_, _ = w.Write([]byte(`{"ok":true,"result":[]}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	manager := NewConnectionManager(rm, NewDispatcher(), NewBotAPIClient(rm, server.Client(), server.URL), NewPollingBotAPIClient(rm, server.Client(), server.URL), nil)
	if err := manager.Start(context.Background()); err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer manager.Stop(context.Background())
	status := manager.Status()
	if status.State != ConnectionStateIdle || status.Receiver != ReceiverNone {
		t.Fatalf("idle polling status = %+v", status)
	}
	time.Sleep(20 * time.Millisecond)
	if getUpdatesCalls.Load() != 0 {
		t.Fatalf("unexpected getUpdates calls: %d", getUpdatesCalls.Load())
	}
}

func TestConnectionManager_PollingDispatchesAndDrains(t *testing.T) {
	rm := newTestRuntimeManager(t, "bot-token", "", nil)
	dispatcher := NewDispatcher()
	handled := make(chan struct{})
	if err := dispatcher.RegisterCommand("status", func(ctx context.Context, upd kernel.TelegramUpdate) error {
		if upd.ChatID != "123" || upd.UserID != "456" {
			t.Errorf("unexpected polling update: %+v", upd)
		}
		select {
		case <-handled:
		default:
			close(handled)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	webhook := NewWebhookHandler(HandlerConfig{Dispatcher: dispatcher, Sender: NewCaptureSender()})
	var getUpdatesCalls atomic.Int32
	secondPoll := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/botbot-token/getMe":
			_, _ = w.Write([]byte(`{"ok":true,"result":{"id":9,"is_bot":true,"username":"poll_bot"}}`))
		case "/botbot-token/deleteWebhook":
			_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	pollingClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if getUpdatesCalls.Add(1) == 1 {
			return botAPIJSONResponse(`{"ok":true,"result":[{"update_id":11,"message":{"from":{"id":456},"chat":{"id":123},"text":"/status"}}]}`), nil
		}
		select {
		case <-secondPoll:
		default:
			close(secondPoll)
		}
		<-req.Context().Done()
		return nil, req.Context().Err()
	})}
	manager := NewConnectionManager(rm, dispatcher, NewBotAPIClient(rm, server.Client(), server.URL), NewPollingBotAPIClient(rm, pollingClient, "https://telegram.test"), webhook.HandlePollingUpdate)
	if err := manager.Start(context.Background()); err != nil {
		t.Fatalf("Start: %v", err)
	}
	select {
	case <-handled:
	case <-time.After(time.Second):
		t.Fatal("polling update was not dispatched")
	}
	select {
	case <-secondPoll:
	case <-time.After(time.Second):
		t.Fatal("polling drain request was not established")
	}
	stopCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stopErr := manager.Stop(stopCtx)
	if stopErr != nil {
		t.Fatalf("Stop: %v", stopErr)
	}
	status := manager.Status()
	if status.State != ConnectionStateIdle || status.Receiver != ReceiverNone {
		t.Fatalf("stopped polling status = %+v", status)
	}
	if getUpdatesCalls.Load() < 1 {
		t.Fatal("expected at least one getUpdates call")
	}
}

func TestConnectionManager_PollingEmptyResultContinues(t *testing.T) {
	rm := newTestRuntimeManager(t, "bot-token", "", nil)
	secondPoll := make(chan struct{})
	var getUpdatesCalls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/botbot-token/getMe":
			_, _ = w.Write([]byte(`{"ok":true,"result":{"id":10,"is_bot":true,"username":"lease_bot"}}`))
		case "/botbot-token/deleteWebhook":
			_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
		case "/botbot-token/getUpdates":
			if getUpdatesCalls.Add(1) == 1 {
				_, _ = w.Write([]byte(`{"ok":true,"result":[]}`))
				return
			}
			select {
			case <-secondPoll:
			default:
				close(secondPoll)
			}
			_, _ = w.Write([]byte(`{"ok":true,"result":[]}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	manager := NewConnectionManager(rm, NewDispatcher(), NewBotAPIClient(rm, server.Client(), server.URL), NewPollingBotAPIClient(rm, server.Client(), server.URL), nil)
	if err := manager.AcquireLease(context.Background(), "console-1"); err != nil {
		t.Fatal(err)
	}
	if err := manager.Start(context.Background()); err != nil {
		t.Fatalf("Start: %v", err)
	}
	select {
	case <-secondPoll:
	case <-time.After(time.Second):
		t.Fatal("polling loop did not continue after an empty result")
	}
	stopCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := manager.Stop(stopCtx); err != nil {
		t.Fatalf("Stop: %v", err)
	}
}

func TestConnectionManager_PollingErrorClearsReceiver(t *testing.T) {
	rm := newTestRuntimeManager(t, "bot-token", "", nil)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/botbot-token/getMe":
			_, _ = w.Write([]byte(`{"ok":true,"result":{"id":13,"is_bot":true,"username":"error_bot"}}`))
		case "/botbot-token/deleteWebhook":
			_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	pollingClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return botAPIJSONResponse(`{"ok":false,"error_code":502,"description":"upstream unavailable"}`), nil
	})}
	manager := NewConnectionManager(rm, NewDispatcher(), NewBotAPIClient(rm, server.Client(), server.URL), NewPollingBotAPIClient(rm, pollingClient, "https://telegram.test"), nil)
	if err := manager.AcquireLease(context.Background(), "console-1"); err != nil {
		t.Fatal(err)
	}
	if err := manager.Start(context.Background()); err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer manager.Stop(context.Background())

	deadline := time.NewTimer(time.Second)
	defer deadline.Stop()
	for {
		status := manager.Status()
		if status.State == ConnectionStateError {
			if status.Receiver != ReceiverNone || status.LastError == "" {
				t.Fatalf("polling error status = %+v", status)
			}
			break
		}
		select {
		case <-deadline.C:
			t.Fatalf("polling error was not reported; status=%+v", status)
		case <-time.After(time.Millisecond):
		}
	}
	manager.stateMu.Lock()
	defer manager.stateMu.Unlock()
	if manager.pollCancel != nil || manager.pollDone != nil {
		t.Fatalf("polling handles retained after asynchronous error: cancel=%v done=%v", manager.pollCancel != nil, manager.pollDone != nil)
	}
}

func TestConnectionManager_ExpiredLeaseDrainsPolling(t *testing.T) {
	rm := newTestRuntimeManager(t, "bot-token", "", nil)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/botbot-token/getMe":
			_, _ = w.Write([]byte(`{"ok":true,"result":{"id":14,"is_bot":true,"username":"lease_bot"}}`))
		case "/botbot-token/deleteWebhook":
			_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	base := time.Now().UTC()
	var clock atomic.Int64
	clock.Store(base.UnixNano())
	pollStarted := make(chan struct{})
	pollStopped := make(chan struct{})
	var startOnce, stopOnce sync.Once
	pollingClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		startOnce.Do(func() { close(pollStarted) })
		<-req.Context().Done()
		stopOnce.Do(func() { close(pollStopped) })
		return nil, req.Context().Err()
	})}
	manager := NewConnectionManager(rm, NewDispatcher(), NewBotAPIClient(rm, server.Client(), server.URL), NewPollingBotAPIClient(rm, pollingClient, "https://telegram.test"), nil)
	manager.now = func() time.Time { return time.Unix(0, clock.Load()).UTC() }
	if err := manager.AcquireLease(context.Background(), "console-1"); err != nil {
		t.Fatal(err)
	}
	if err := manager.Start(context.Background()); err != nil {
		t.Fatalf("Start: %v", err)
	}
	select {
	case <-pollStarted:
	case <-time.After(time.Second):
		t.Fatal("polling did not start for a live lease")
	}

	clock.Store(base.Add(PollingLeaseTTL + time.Second).UnixNano())
	deadline := time.NewTimer(3 * time.Second)
	defer deadline.Stop()
	for {
		status := manager.Status()
		if status.State == ConnectionStateIdle && status.Receiver == ReceiverNone {
			break
		}
		select {
		case <-deadline.C:
			t.Fatalf("expired lease did not drain polling; status=%+v", status)
		case <-time.After(10 * time.Millisecond):
		}
	}
	select {
	case <-pollStopped:
	case <-time.After(time.Second):
		t.Fatal("expired lease did not cancel polling request")
	}
	if got := manager.ActiveLeaseCount(); got != 0 {
		t.Fatalf("active leases after expiry = %d, want 0", got)
	}
	if err := manager.Stop(context.Background()); err != nil {
		t.Fatalf("Stop: %v", err)
	}
}

func TestConnectionManager_FailedModeSwitchDrainsPolling(t *testing.T) {
	rm := newTestRuntimeManager(t, "bot-token", "old-secret", nil)
	var mu sync.Mutex
	var paths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		paths = append(paths, r.URL.Path)
		mu.Unlock()
		switch r.URL.Path {
		case "/botbot-token/getMe":
			_, _ = w.Write([]byte(`{"ok":true,"result":{"id":15,"is_bot":true,"username":"switch_error_bot"}}`))
		case "/botbot-token/deleteWebhook":
			_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	pollStarted := make(chan struct{})
	pollStopped := make(chan struct{})
	var startOnce, stopOnce sync.Once
	pollingClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		startOnce.Do(func() { close(pollStarted) })
		<-req.Context().Done()
		stopOnce.Do(func() { close(pollStopped) })
		return nil, req.Context().Err()
	})}
	manager := NewConnectionManager(rm, NewDispatcher(), NewBotAPIClient(rm, server.Client(), server.URL), NewPollingBotAPIClient(rm, pollingClient, "https://telegram.test"), nil)
	rm.SetSettingsChangedHandler(manager.Reconcile)
	if err := manager.AcquireLease(context.Background(), "console-1"); err != nil {
		t.Fatal(err)
	}
	if err := manager.Start(context.Background()); err != nil {
		t.Fatalf("Start: %v", err)
	}
	select {
	case <-pollStarted:
	case <-time.After(time.Second):
		t.Fatal("polling did not start before mode switch")
	}
	if err := rm.UpdateSettings(context.Background(), "bot-token", "", TelegramModeWebhook, "https://console.example"); err == nil {
		t.Fatal("mode switch unexpectedly succeeded without webhook secret")
	}
	select {
	case <-pollStopped:
	case <-time.After(time.Second):
		t.Fatal("failed mode switch did not drain polling")
	}
	status := manager.Status()
	if status.State != ConnectionStateError || status.Receiver != ReceiverNone || status.LastError == "" {
		t.Fatalf("failed switch status = %+v", status)
	}
	mu.Lock()
	gotPaths := append([]string(nil), paths...)
	mu.Unlock()
	wantPaths := []string{
		"/botbot-token/getMe", "/botbot-token/deleteWebhook",
		"/botbot-token/getMe",
	}
	if !reflect.DeepEqual(gotPaths, wantPaths) {
		t.Fatalf("failed switch call order = %v, want %v", gotPaths, wantPaths)
	}
	if err := manager.Stop(context.Background()); err != nil {
		t.Fatalf("Stop: %v", err)
	}
}

func TestConnectionManager_FailsClosedOnBotAPIError(t *testing.T) {
	rm := newTestRuntimeManager(t, "bot-token", "", nil)
	var paths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path)
		_, _ = w.Write([]byte(`{"ok":false,"error_code":401,"description":"Unauthorized"}`))
	}))
	defer server.Close()
	manager := NewConnectionManager(rm, NewDispatcher(), NewBotAPIClient(rm, server.Client(), server.URL), NewPollingBotAPIClient(rm, server.Client(), server.URL), nil)
	if err := manager.Start(context.Background()); err == nil {
		t.Fatal("Start unexpectedly succeeded")
	}
	status := manager.Status()
	if status.State != ConnectionStateError || status.Receiver != ReceiverNone || status.LastError == "" {
		t.Fatalf("error status = %+v", status)
	}
	if !reflect.DeepEqual(paths, []string{"/botbot-token/getMe"}) {
		t.Fatalf("calls after getMe failure = %v", paths)
	}
}

func TestConnectionManager_WebhookMissingSecretDoesNotSetWebhook(t *testing.T) {
	rm, err := NewRuntimeManagerWithSettings("bot-token", "", TelegramModeWebhook, "https://console.example", nil, testMasterKey())
	if err != nil {
		t.Fatal(err)
	}
	var paths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path)
		_, _ = w.Write([]byte(`{"ok":true,"result":{"id":11,"is_bot":true,"username":"missing_secret_bot"}}`))
	}))
	defer server.Close()
	manager := NewConnectionManager(rm, NewDispatcher(), NewBotAPIClient(rm, server.Client(), server.URL), NewPollingBotAPIClient(rm, server.Client(), server.URL), nil)
	if err := manager.Start(context.Background()); err == nil {
		t.Fatal("Start unexpectedly succeeded")
	}
	if manager.Status().State != ConnectionStateError {
		t.Fatalf("status = %+v", manager.Status())
	}
	if !reflect.DeepEqual(paths, []string{"/botbot-token/getMe"}) {
		t.Fatalf("calls = %v", paths)
	}
}

func TestConnectionManager_SettingsUpdateHotSwitchesMode(t *testing.T) {
	rm := newTestRuntimeManager(t, "bot-token", "webhook-secret", nil)
	dispatcher := NewDispatcher()
	var mu sync.Mutex
	var paths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		paths = append(paths, r.URL.Path)
		mu.Unlock()
		switch r.URL.Path {
		case "/botbot-token/getMe":
			_, _ = w.Write([]byte(`{"ok":true,"result":{"id":12,"is_bot":true,"username":"switch_bot"}}`))
		case "/botbot-token/deleteWebhook", "/botbot-token/setWebhook":
			_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	manager := NewConnectionManager(rm, dispatcher, NewBotAPIClient(rm, server.Client(), server.URL), NewPollingBotAPIClient(rm, server.Client(), server.URL), nil)
	rm.SetSettingsChangedHandler(manager.Reconcile)
	if err := manager.Start(context.Background()); err != nil {
		t.Fatalf("Start: %v", err)
	}
	if err := rm.UpdateSettings(context.Background(), "bot-token", "webhook-secret", TelegramModeWebhook, "https://console.example"); err != nil {
		t.Fatalf("switch to webhook: %v", err)
	}
	if status := manager.Status(); status.State != ConnectionStateRunning || status.Receiver != ReceiverWebhook {
		t.Fatalf("webhook switch status = %+v", status)
	}
	if err := rm.UpdateSettings(context.Background(), "bot-token", "", TelegramModePolling, ""); err != nil {
		t.Fatalf("switch to polling: %v", err)
	}
	if status := manager.Status(); status.State != ConnectionStateIdle || status.Receiver != ReceiverNone {
		t.Fatalf("polling switch status = %+v", status)
	}
	if err := manager.Stop(context.Background()); err != nil {
		t.Fatalf("Stop: %v", err)
	}
	mu.Lock()
	gotPaths := append([]string(nil), paths...)
	mu.Unlock()
	wantPaths := []string{
		"/botbot-token/getMe", "/botbot-token/deleteWebhook",
		"/botbot-token/getMe", "/botbot-token/setWebhook",
		"/botbot-token/getMe", "/botbot-token/deleteWebhook",
	}
	if !reflect.DeepEqual(gotPaths, wantPaths) {
		t.Fatalf("hot-switch call order = %v, want %v", gotPaths, wantPaths)
	}
}
