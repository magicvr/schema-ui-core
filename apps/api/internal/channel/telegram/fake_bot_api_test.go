package telegram

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type fakeBotAPIResponse struct {
	status int
	body   string
}

// fakeBotAPI is a deterministic local Bot API surface for C5. It records the
// method sequence, captures setWebhook payloads, and can hold getUpdates until
// the request context is canceled so manager drain is observable.
type fakeBotAPI struct {
	server *httptest.Server

	mu             sync.Mutex
	responses      map[string]fakeBotAPIResponse
	calls          []string
	webhookPayload setWebhookPayload
}

func newFakeBotAPI(t *testing.T) *fakeBotAPI {
	t.Helper()
	fake := &fakeBotAPI{
		responses: map[string]fakeBotAPIResponse{
			"getMe":         {status: http.StatusOK, body: `{"ok":true,"result":{"id":101,"is_bot":true,"username":"fake_bot"}}`},
			"setWebhook":    {status: http.StatusOK, body: `{"ok":true,"result":true}`},
			"deleteWebhook": {status: http.StatusOK, body: `{"ok":true,"result":true}`},
			"getUpdates":    {status: http.StatusOK, body: `{"ok":true,"result":[]}`},
		},
	}
	fake.server = httptest.NewServer(http.HandlerFunc(fake.serveHTTP))
	t.Cleanup(fake.server.Close)
	return fake
}

func (f *fakeBotAPI) setResponse(method string, status int, body string) {
	f.mu.Lock()
	f.responses[method] = fakeBotAPIResponse{status: status, body: body}
	f.mu.Unlock()
}

func (f *fakeBotAPI) callsSnapshot() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.calls...)
}

func (f *fakeBotAPI) webhookPayloadSnapshot() setWebhookPayload {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.webhookPayload
}

func (f *fakeBotAPI) serveHTTP(w http.ResponseWriter, r *http.Request) {
	method := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
	f.mu.Lock()
	f.calls = append(f.calls, method)
	response := f.responses[method]
	f.mu.Unlock()

	if method == "setWebhook" {
		body, _ := io.ReadAll(r.Body)
		var payload setWebhookPayload
		if err := json.Unmarshal(body, &payload); err == nil {
			f.mu.Lock()
			f.webhookPayload = payload
			f.mu.Unlock()
		}
	}
	if response.status == 0 {
		response.status = http.StatusOK
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.status)
	_, _ = io.WriteString(w, response.body)
}

func TestFakeBotAPI_ConnectionLifecycle(t *testing.T) {
	t.Run("polling lease release cancels long poll", func(t *testing.T) {
		fake := newFakeBotAPI(t)
		rm := newTestRuntimeManager(t, "bot-token", "", nil)
		pollStarted := make(chan struct{})
		pollCanceled := make(chan struct{})
		var pollStartOnce, pollCancelOnce sync.Once
		var pollCalls atomic.Int32
		pollingClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			pollCalls.Add(1)
			pollStartOnce.Do(func() { close(pollStarted) })
			<-req.Context().Done()
			pollCancelOnce.Do(func() { close(pollCanceled) })
			return nil, req.Context().Err()
		})}
		manager := NewConnectionManager(rm, NewDispatcher(), NewBotAPIClient(rm, fake.server.Client(), fake.server.URL), NewPollingBotAPIClient(rm, pollingClient, "https://telegram.test"), nil)
		t.Cleanup(func() { _ = manager.Stop(context.Background()) })

		if err := manager.Start(context.Background()); err != nil {
			t.Fatalf("Start: %v", err)
		}
		if got := fake.callsSnapshot(); !reflect.DeepEqual(got, []string{"getMe", "deleteWebhook"}) {
			t.Fatalf("idle startup calls = %v", got)
		}
		if err := manager.AcquireLease(context.Background(), "console-session"); err != nil {
			t.Fatalf("AcquireLease: %v", err)
		}
		select {
		case <-pollStarted:
		case <-time.After(time.Second):
			t.Fatal("fake long poll did not start")
		}
		status := manager.Status()
		if status.State != ConnectionStateRunning || status.Receiver != ReceiverPolling {
			t.Fatalf("polling status = %+v", status)
		}

		if err := manager.ReleaseLease(context.Background(), "console-session"); err != nil {
			t.Fatalf("ReleaseLease: %v", err)
		}
		select {
		case <-pollCanceled:
		case <-time.After(time.Second):
			t.Fatal("ReleaseLease returned before the fake long poll observed cancellation")
		}
		status = manager.Status()
		if status.State != ConnectionStateIdle || status.Receiver != ReceiverNone {
			t.Fatalf("post-release status = %+v", status)
		}
		if err := manager.Stop(context.Background()); err != nil {
			t.Fatalf("Stop: %v", err)
		}
		if got := fake.callsSnapshot(); !reflect.DeepEqual(got, []string{"getMe", "deleteWebhook"}) {
			t.Fatalf("polling lifecycle calls = %v, polling transport calls = %d", got, pollCalls.Load())
		}
		if got := pollCalls.Load(); got != 1 {
			t.Fatalf("polling transport calls = %d, want 1", got)
		}
	})

	t.Run("webhook setup sends the explicit secret", func(t *testing.T) {
		fake := newFakeBotAPI(t)
		rm, err := NewRuntimeManagerWithSettings("bot-token", "webhook-secret", TelegramModeWebhook, "https://console.example", nil, testMasterKey())
		if err != nil {
			t.Fatal(err)
		}
		manager := NewConnectionManager(rm, NewDispatcher(), NewBotAPIClient(rm, fake.server.Client(), fake.server.URL), NewPollingBotAPIClient(rm, fake.server.Client(), fake.server.URL), nil)
		if err := manager.Start(context.Background()); err != nil {
			t.Fatalf("Start: %v", err)
		}
		defer manager.Stop(context.Background())

		status := manager.Status()
		if status.State != ConnectionStateRunning || status.Receiver != ReceiverWebhook {
			t.Fatalf("webhook status = %+v", status)
		}
		if got := fake.callsSnapshot(); !reflect.DeepEqual(got, []string{"getMe", "setWebhook"}) {
			t.Fatalf("webhook calls = %v", got)
		}
		payload := fake.webhookPayloadSnapshot()
		if payload.URL != "https://console.example/api/channel/telegram/webhook" || payload.SecretToken != "webhook-secret" {
			t.Fatalf("webhook payload = %+v", payload)
		}
	})
}

func TestFakeBotAPI_ErrorMatrix(t *testing.T) {
	shapes := []struct {
		name   string
		status int
		body   string
	}{
		{name: "http status", status: http.StatusBadGateway, body: `{"ok":false,"description":"gateway"}`},
		{name: "api error", status: http.StatusOK, body: `{"ok":false,"error_code":401,"description":"Unauthorized"}`},
		{name: "malformed envelope", status: http.StatusOK, body: "{"},
		{name: "missing result", status: http.StatusOK, body: `{"ok":true}`},
		{name: "malformed result", status: http.StatusOK, body: `{"ok":true,"result":"not-the-expected-shape"}`},
		{name: "false result", status: http.StatusOK, body: `{"ok":true,"result":false}`},
	}
	methods := []struct {
		name    string
		method  string
		polling bool
		invoke  func(*BotAPIClient) error
	}{
		{name: "getMe", method: "getMe", invoke: func(client *BotAPIClient) error {
			_, err := client.GetMe(context.Background())
			return err
		}},
		{name: "setWebhook", method: "setWebhook", invoke: func(client *BotAPIClient) error {
			return client.SetWebhook(context.Background(), "https://console.example/api/channel/telegram/webhook", "secret")
		}},
		{name: "deleteWebhook", method: "deleteWebhook", invoke: func(client *BotAPIClient) error {
			return client.DeleteWebhook(context.Background())
		}},
		{name: "getUpdates", method: "getUpdates", polling: true, invoke: func(client *BotAPIClient) error {
			_, err := client.GetUpdates(context.Background(), 0)
			return err
		}},
	}

	for _, method := range methods {
		method := method
		for _, shape := range shapes {
			shape := shape
			t.Run(method.name+"/"+shape.name, func(t *testing.T) {
				fake := newFakeBotAPI(t)
				fake.setResponse(method.method, shape.status, shape.body)
				rm := newTestRuntimeManager(t, "bot-token", "webhook-secret", nil)
				var client *BotAPIClient
				if method.polling {
					client = NewPollingBotAPIClient(rm, fake.server.Client(), fake.server.URL)
				} else {
					client = NewBotAPIClient(rm, fake.server.Client(), fake.server.URL)
				}
				if err := method.invoke(client); err == nil {
					t.Fatalf("%s unexpectedly accepted %s", method.name, shape.name)
				}
				if got := fake.callsSnapshot(); !reflect.DeepEqual(got, []string{method.method}) {
					t.Fatalf("calls = %v, want [%s]", got, method.method)
				}
			})
		}
	}
}

func TestFakeBotAPI_ResponseBodyLimit(t *testing.T) {
	fake := newFakeBotAPI(t)
	fake.setResponse("getMe", http.StatusOK, `{"ok":true,"result":{"id":1,"is_bot":true,"username":"`+strings.Repeat("x", BotAPIResponseBodyLimit)+`"}}`)
	rm := newTestRuntimeManager(t, "bot-token", "", nil)
	client := NewBotAPIClient(rm, fake.server.Client(), fake.server.URL)
	if _, err := client.GetMe(context.Background()); err == nil {
		t.Fatal("oversized response unexpectedly succeeded")
	}
}

func TestBotAPIClient_TransportErrorMatrix(t *testing.T) {
	methods := []struct {
		name    string
		polling bool
		invoke  func(*BotAPIClient) error
	}{
		{name: "getMe", invoke: func(client *BotAPIClient) error { _, err := client.GetMe(context.Background()); return err }},
		{name: "setWebhook", invoke: func(client *BotAPIClient) error {
			return client.SetWebhook(context.Background(), "https://console.example", "secret")
		}},
		{name: "deleteWebhook", invoke: func(client *BotAPIClient) error { return client.DeleteWebhook(context.Background()) }},
		{name: "getUpdates", polling: true, invoke: func(client *BotAPIClient) error { _, err := client.GetUpdates(context.Background(), 0); return err }},
	}
	for _, method := range methods {
		method := method
		t.Run(method.name, func(t *testing.T) {
			rm := newTestRuntimeManager(t, "bot-token", "webhook-secret", nil)
			transport := roundTripFunc(func(*http.Request) (*http.Response, error) {
				return nil, errors.New("dial failed")
			})
			client := NewBotAPIClient(rm, &http.Client{Transport: transport}, "https://telegram.test")
			if method.polling {
				client = NewPollingBotAPIClient(rm, &http.Client{Transport: transport}, "https://telegram.test")
			}
			err := method.invoke(client)
			if err == nil || !strings.Contains(err.Error(), "execute request failed") {
				t.Fatalf("transport error = %v", err)
			}
			if strings.Contains(err.Error(), "bot-token") || strings.Contains(err.Error(), "webhook-secret") {
				t.Fatalf("transport error leaked a secret: %v", err)
			}
		})
	}
}
