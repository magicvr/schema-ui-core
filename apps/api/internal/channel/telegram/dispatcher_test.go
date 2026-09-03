package telegram

import (
	"context"
	"errors"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

func TestDispatcher_RegisterAndDispatchCommand(t *testing.T) {
	ctx := context.Background()
	disp := NewDispatcher()
	sender := NewCaptureSender()

	var received kernel.TelegramUpdate
	called := false
	handler := func(ctx context.Context, upd kernel.TelegramUpdate) error {
		called = true
		received = upd
		return nil
	}

	// Register with leading slash and bot username
	if err := disp.RegisterCommand("/start@MyTestBot", handler); err != nil {
		t.Fatalf("unexpected error registering command: %v", err)
	}

	// Conflict error on re-registration
	if err := disp.RegisterCommand("start", handler); !errors.Is(err, kernel.ErrTelegramCommandConflict) {
		t.Fatalf("expected ErrTelegramCommandConflict, got %v", err)
	}

	// Dispatch registered command
	upd := kernel.TelegramUpdate{
		ChatID:    "123456",
		UserID:    "7890",
		SubjectID: "sub-1",
		Command:   "/start",
		Text:      "/start hello",
	}
	if err := disp.Dispatch(ctx, upd, sender); err != nil {
		t.Fatalf("Dispatch failed: %v", err)
	}
	if !called {
		t.Fatalf("expected handler to be called")
	}
	if received.UserID != "7890" || received.ChatID != "123456" || received.SubjectID != "sub-1" {
		t.Fatalf("unexpected received update: %+v", received)
	}
	if sender.Last() != nil {
		t.Fatalf("expected no fallback message for registered command")
	}

	// Unregister
	disp.UnregisterCommand("start")
	called = false

	// Dispatch again: now it's an unknown command -> triggers fallback sender message
	if err := disp.Dispatch(ctx, upd, sender); err != nil {
		t.Fatalf("Dispatch failed: %v", err)
	}
	if called {
		t.Fatalf("expected unregistered handler not to be called")
	}
	lastMsg := sender.Last()
	if lastMsg == nil {
		t.Fatalf("expected fallback message sent for unknown command")
	}
	if lastMsg.ChatID != "123456" || lastMsg.Text != kernel.DefaultTelegramUnknownCommandText {
		t.Fatalf("unexpected fallback message: %+v", lastMsg)
	}
}

func TestDispatcher_RegisterAndDispatchCallback(t *testing.T) {
	ctx := context.Background()
	disp := NewDispatcher()
	sender := NewCaptureSender()

	called := false
	var received kernel.TelegramUpdate
	handler := func(ctx context.Context, upd kernel.TelegramUpdate) error {
		called = true
		received = upd
		return nil
	}

	// Register callback
	if err := disp.RegisterCallback("btn_click_confirm", handler); err != nil {
		t.Fatalf("unexpected error registering callback: %v", err)
	}

	// Conflict
	if err := disp.RegisterCallback("btn_click_confirm", handler); !errors.Is(err, kernel.ErrTelegramCallbackConflict) {
		t.Fatalf("expected ErrTelegramCallbackConflict, got %v", err)
	}

	// Dispatch callback
	upd := kernel.TelegramUpdate{
		ChatID:       "999",
		UserID:       "888",
		CallbackData: "btn_click_confirm",
	}
	if err := disp.Dispatch(ctx, upd, sender); err != nil {
		t.Fatalf("Dispatch failed: %v", err)
	}
	if !called || received.CallbackData != "btn_click_confirm" {
		t.Fatalf("expected callback handler called with correct data")
	}

	// Unknown callback is a silent no-op (no error, no sender message)
	sender.Reset()
	unknownUpd := kernel.TelegramUpdate{
		ChatID:       "999",
		UserID:       "888",
		CallbackData: "unknown_action",
	}
	if err := disp.Dispatch(ctx, unknownUpd, sender); err != nil {
		t.Fatalf("Dispatch unknown callback failed: %v", err)
	}
	if sender.Last() != nil {
		t.Fatalf("expected no message sent for unknown callback")
	}

	// Unregister
	disp.UnregisterCallback("btn_click_confirm")
	called = false
	if err := disp.Dispatch(ctx, upd, sender); err != nil {
		t.Fatalf("Dispatch unregistered callback failed: %v", err)
	}
	if called {
		t.Fatalf("expected unregistered callback handler not to be called")
	}
}

func TestDispatcher_InvalidRegistrations(t *testing.T) {
	disp := NewDispatcher()
	dummy := func(ctx context.Context, upd kernel.TelegramUpdate) error { return nil }

	// Nil handler
	if err := disp.RegisterCommand("start", nil); !errors.Is(err, kernel.ErrTelegramHandlerNil) {
		t.Fatalf("expected ErrTelegramHandlerNil, got %v", err)
	}
	if err := disp.RegisterCallback("btn", nil); !errors.Is(err, kernel.ErrTelegramHandlerNil) {
		t.Fatalf("expected ErrTelegramHandlerNil, got %v", err)
	}

	// Empty command
	if err := disp.RegisterCommand("", dummy); !errors.Is(err, kernel.ErrTelegramCommandEmpty) {
		t.Fatalf("expected ErrTelegramCommandEmpty, got %v", err)
	}
	if err := disp.RegisterCommand("/", dummy); !errors.Is(err, kernel.ErrTelegramCommandEmpty) {
		t.Fatalf("expected ErrTelegramCommandEmpty, got %v", err)
	}

	// Invalid command with spaces or slashes
	if err := disp.RegisterCommand("invalid cmd", dummy); !errors.Is(err, kernel.ErrTelegramCommandInvalid) {
		t.Fatalf("expected ErrTelegramCommandInvalid, got %v", err)
	}
	if err := disp.RegisterCommand("invalid/cmd", dummy); !errors.Is(err, kernel.ErrTelegramCommandInvalid) {
		t.Fatalf("expected ErrTelegramCommandInvalid, got %v", err)
	}

	// Empty callback
	if err := disp.RegisterCallback("", dummy); !errors.Is(err, kernel.ErrTelegramCallbackEmpty) {
		t.Fatalf("expected ErrTelegramCallbackEmpty, got %v", err)
	}

	// Callback exceeding 64 bytes
	tooLong := "12345678901234567890123456789012345678901234567890123456789012345" // 65 bytes
	if err := disp.RegisterCallback(tooLong, dummy); !errors.Is(err, kernel.ErrTelegramCallbackTooLong) {
		t.Fatalf("expected ErrTelegramCallbackTooLong, got %v", err)
	}
}
