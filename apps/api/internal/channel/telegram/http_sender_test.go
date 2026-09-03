package telegram

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

func TestHTTPSender_UnconfiguredToken_DowngradesToMock(t *testing.T) {
	mock := NewCaptureSender()
	rm := newTestRuntimeManager(t, "", "", mock)
	sender := NewHTTPSender(rm, nil, "")

	msg := kernel.TelegramMessage{
		ChatID: "12345678",
		Text:   "Hello mock fallback",
	}

	if err := sender.Send(context.Background(), msg); err != nil {
		t.Fatalf("unexpected error on mock fallback: %v", err)
	}

	last := mock.Last()
	if last == nil || last.Text != "Hello mock fallback" || last.ChatID != "12345678" {
		t.Fatalf("expected message captured in mock: %+v", last)
	}
}

func TestHTTPSender_ValidationFailClosed(t *testing.T) {
	rm := newTestRuntimeManager(t, "test-token", "", nil)
	sender := NewHTTPSender(rm, nil, "http://invalid-host")

	// Invalid empty text
	msg := kernel.TelegramMessage{
		ChatID: "12345678",
		Text:   "",
	}

	err := sender.Send(context.Background(), msg)
	if !errors.Is(err, kernel.ErrTelegramMessageInvalid) {
		t.Fatalf("expected ErrTelegramMessageInvalid, got %v", err)
	}
}

func TestHTTPSender_ValidMessageDelivery(t *testing.T) {
	var capturedURL string
	var capturedPayload sendMessagePayload

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedURL = r.URL.String()
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &capturedPayload)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	rm := newTestRuntimeManager(t, "my-bot-token", "my-secret", nil)
	sender := NewHTTPSender(rm, server.Client(), server.URL)

	msg := kernel.TelegramMessage{
		ChatID: "-1001234567890",
		Text:   "Production message with keyboard",
		Buttons: [][]kernel.TelegramInlineButton{
			{
				{Text: "Approve", CallbackData: "act_approve"},
				{Text: "Reject", CallbackData: "act_reject"},
			},
		},
	}

	if err := sender.Send(context.Background(), msg); err != nil {
		t.Fatalf("Send failed: %v", err)
	}

	expectedURL := "/botmy-bot-token/sendMessage"
	if capturedURL != expectedURL {
		t.Fatalf("expected URL %q, got %q", expectedURL, capturedURL)
	}
	if capturedPayload.ChatID != "-1001234567890" || capturedPayload.Text != "Production message with keyboard" {
		t.Fatalf("unexpected payload: %+v", capturedPayload)
	}
	if capturedPayload.ReplyMarkup == nil || len(capturedPayload.ReplyMarkup.InlineKeyboard) != 1 {
		t.Fatalf("expected inline keyboard with 1 row")
	}
	buttons := capturedPayload.ReplyMarkup.InlineKeyboard[0]
	if len(buttons) != 2 || buttons[0].Text != "Approve" || buttons[1].CallbackData != "act_reject" {
		t.Fatalf("unexpected buttons: %+v", buttons)
	}
}

func TestHTTPSender_APIErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"ok": false, "description": "Chat not found"}`))
	}))
	defer server.Close()

	rm := newTestRuntimeManager(t, "my-bot-token", "", nil)
	sender := NewHTTPSender(rm, server.Client(), server.URL)

	msg := kernel.TelegramMessage{
		ChatID: "999999",
		Text:   "Test message",
	}

	err := sender.Send(context.Background(), msg)
	if err == nil {
		t.Fatalf("expected error on 400 Bad Request, got nil")
	}
}

func TestHTTPSender_Status200_ButOKFalse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": false, "error_code": 403, "description": "Forbidden: bot was blocked by the user"}`))
	}))
	defer server.Close()

	rm := newTestRuntimeManager(t, "my-bot-token", "", nil)
	sender := NewHTTPSender(rm, server.Client(), server.URL)

	msg := kernel.TelegramMessage{
		ChatID: "123456",
		Text:   "Test message",
	}

	err := sender.Send(context.Background(), msg)
	if err == nil {
		t.Fatalf("expected error when ok=false, got nil")
	}
	if !strings.Contains(err.Error(), "bot was blocked by the user") {
		t.Fatalf("expected error message to contain description, got %v", err)
	}
}

func TestDisabledSenderAndDispatcher(t *testing.T) {
	sender := NewDisabledSender()
	msg := kernel.TelegramMessage{
		ChatID: "123",
		Text:   "Hello",
	}
	if err := sender.Send(context.Background(), msg); !errors.Is(err, kernel.ErrTelegramDisabled) {
		t.Fatalf("expected ErrTelegramDisabled from DisabledSender, got %v", err)
	}

	disp := NewDisabledDispatcher()
	if err := disp.RegisterCommand("start", func(ctx context.Context, upd kernel.TelegramUpdate) error { return nil }); err != nil {
		t.Fatalf("expected nil from DisabledDispatcher RegisterCommand, got %v", err)
	}
	disp.UnregisterCommand("start")
	if err := disp.RegisterCallback("cb", func(ctx context.Context, upd kernel.TelegramUpdate) error { return nil }); err != nil {
		t.Fatalf("expected nil from DisabledDispatcher RegisterCallback, got %v", err)
	}
	disp.UnregisterCallback("cb")
}

func TestHTTPSender_TimeoutBudget(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	rm := newTestRuntimeManager(t, "my-bot-token", "", nil)
	sender := NewHTTPSender(rm, server.Client(), server.URL)

	// Context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	msg := kernel.TelegramMessage{
		ChatID: "123456",
		Text:   "Will timeout",
	}

	err := sender.Send(ctx, msg)
	if err == nil {
		t.Fatalf("expected timeout error, got nil")
	}
}
