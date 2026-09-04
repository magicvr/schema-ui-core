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
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func botAPIJSONResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

func TestBotAPIClient_ManagementMethodsAndPayloads(t *testing.T) {
	rm := newTestRuntimeManager(t, "bot-token", "webhook-secret", nil)
	var paths []string
	var webhookPayload setWebhookPayload
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path)
		if r.URL.Path == "/botbot-token/setWebhook" {
			body, _ := io.ReadAll(r.Body)
			if err := json.Unmarshal(body, &webhookPayload); err != nil {
				t.Errorf("decode setWebhook payload: %v", err)
			}
			_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
			return
		}
		if r.URL.Path == "/botbot-token/getMe" {
			_, _ = w.Write([]byte(`{"ok":true,"result":{"id":42,"is_bot":true,"username":"ops_bot"}}`))
			return
		}
		if r.URL.Path == "/botbot-token/deleteWebhook" {
			_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	client := NewBotAPIClient(rm, server.Client(), server.URL)
	bot, err := client.GetMe(context.Background())
	if err != nil {
		t.Fatalf("GetMe: %v", err)
	}
	if bot.ID != 42 || bot.Username != "ops_bot" {
		t.Fatalf("unexpected bot identity: %+v", bot)
	}
	if err := client.SetWebhook(context.Background(), "https://console.example/api/channel/telegram/webhook", "webhook-secret"); err != nil {
		t.Fatalf("SetWebhook: %v", err)
	}
	if err := client.DeleteWebhook(context.Background()); err != nil {
		t.Fatalf("DeleteWebhook: %v", err)
	}

	wantPaths := []string{
		"/botbot-token/getMe",
		"/botbot-token/setWebhook",
		"/botbot-token/deleteWebhook",
	}
	if !reflect.DeepEqual(paths, wantPaths) {
		t.Fatalf("management call order = %v, want %v", paths, wantPaths)
	}
	if webhookPayload.URL != "https://console.example/api/channel/telegram/webhook" || webhookPayload.SecretToken != "webhook-secret" {
		t.Fatalf("unexpected setWebhook payload: %+v", webhookPayload)
	}
}

func TestPollingBotAPIClient_RequestAndClientTimeouts(t *testing.T) {
	rm := newTestRuntimeManager(t, "bot-token", "", nil)
	var payload getUpdatesPayload
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Errorf("decode getUpdates payload: %v", err)
		}
		_, _ = w.Write([]byte(`{"ok":true,"result":[]}`))
	}))
	defer server.Close()

	client := NewPollingBotAPIClient(rm, server.Client(), server.URL)
	if client.contextTimeout != PollingRequestContextTimeout {
		t.Fatalf("context timeout = %s, want %s", client.contextTimeout, PollingRequestContextTimeout)
	}
	if client.client.Timeout != PollingHTTPClientTimeout {
		t.Fatalf("client timeout = %s, want %s", client.client.Timeout, PollingHTTPClientTimeout)
	}
	if !(GetUpdatesRequestTimeout < client.contextTimeout && client.contextTimeout < client.client.Timeout) {
		t.Fatalf("polling timeout ordering = API %s, context %s, client %s", GetUpdatesRequestTimeout, client.contextTimeout, client.client.Timeout)
	}
	if _, err := client.GetUpdates(context.Background(), 17); err != nil {
		t.Fatalf("GetUpdates: %v", err)
	}
	if payload.Offset != 17 || payload.Timeout != 30 {
		t.Fatalf("getUpdates payload = %+v, want offset 17 timeout 30", payload)
	}
}

func TestPollingBotAPIClient_ContextDeadlineLeavesLongPollGrace(t *testing.T) {
	rm := newTestRuntimeManager(t, "bot-token", "", nil)
	deadlineSeen := make(chan time.Time, 1)
	client := NewPollingBotAPIClient(rm, &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		deadline, ok := req.Context().Deadline()
		if !ok {
			t.Error("polling request has no context deadline")
		} else {
			deadlineSeen <- deadline
		}
		return botAPIJSONResponse(`{"ok":true,"result":[]}`), nil
	})}, "https://telegram.test")
	if _, err := client.GetUpdates(context.Background(), 0); err != nil {
		t.Fatalf("GetUpdates: %v", err)
	}
	select {
	case deadline := <-deadlineSeen:
		remaining := time.Until(deadline)
		if remaining <= GetUpdatesRequestTimeout-time.Second {
			t.Fatalf("context deadline leaves insufficient long-poll grace: %s", remaining)
		}
		if remaining > PollingRequestContextTimeout {
			t.Fatalf("context deadline exceeds configured timeout: %s", remaining)
		}
	case <-time.After(time.Second):
		t.Fatal("polling request deadline was not observed")
	}
}

func TestBotAPIClient_TransportErrorsAreSanitized(t *testing.T) {
	rm := newTestRuntimeManager(t, "bot-token", "webhook-secret", nil)
	client := NewBotAPIClient(rm, &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("dial failed")
	})}, "https://telegram.test")
	_, err := client.GetMe(context.Background())
	if err == nil {
		t.Fatal("GetMe unexpectedly succeeded")
	}
	if strings.Contains(err.Error(), "bot-token") || strings.Contains(err.Error(), "webhook-secret") {
		t.Fatalf("transport error leaked secret: %v", err)
	}
	if !strings.Contains(err.Error(), "execute request failed") {
		t.Fatalf("transport error lost safe diagnostic: %v", err)
	}
}

func TestBotAPIClient_NonOKAndProtocolErrorsFailClosed(t *testing.T) {
	rm := newTestRuntimeManager(t, "bot-token", "", nil)
	responses := []string{
		`{"ok":false,"error_code":401,"description":"Unauthorized"}`,
		`<html>gateway error</html>`,
	}
	for _, response := range responses {
		response := response
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if response[0] == '<' {
				w.WriteHeader(http.StatusOK)
			}
			_, _ = w.Write([]byte(response))
		}))
		client := NewBotAPIClient(rm, server.Client(), server.URL)
		if _, err := client.GetMe(context.Background()); err == nil {
			t.Fatalf("GetMe unexpectedly succeeded for response %q", response)
		}
		server.Close()
	}
}

func TestPollingBotAPIClient_ContextCancellation(t *testing.T) {
	rm := newTestRuntimeManager(t, "bot-token", "", nil)
	client := NewPollingBotAPIClient(rm, &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		<-req.Context().Done()
		return nil, req.Context().Err()
	})}, "https://telegram.test")
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, err := client.GetUpdates(ctx, 0)
		done <- err
	}()
	cancel()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("GetUpdates cancellation returned nil")
		}
	case <-time.After(time.Second):
		t.Fatal("GetUpdates did not observe context cancellation")
	}
}
