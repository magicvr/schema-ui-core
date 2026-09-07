package kernel

import (
	"context"
	"errors"
	"strings"
	"testing"
)

// Compile-time interface checks for stub implementations.
type stubTelegramSender struct {
	lastMsg TelegramMessage
}

func (s *stubTelegramSender) Send(ctx context.Context, msg TelegramMessage) error {
	if err := msg.Validate(); err != nil {
		return err
	}
	s.lastMsg = msg
	return nil
}

var _ TelegramSender = (*stubTelegramSender)(nil)

type stubTelegramDispatcher struct {
	commands  map[string]TelegramHandler
	callbacks map[string]TelegramHandler
}

func newStubDispatcher() *stubTelegramDispatcher {
	return &stubTelegramDispatcher{
		commands:  make(map[string]TelegramHandler),
		callbacks: make(map[string]TelegramHandler),
	}
}

func (d *stubTelegramDispatcher) RegisterCommand(name string, h TelegramHandler) error {
	norm, err := NormalizeTelegramCommand(name)
	if err != nil {
		return err
	}
	if h == nil {
		return ErrTelegramHandlerNil
	}
	if _, exists := d.commands[norm]; exists {
		return ErrTelegramCommandConflict
	}
	d.commands[norm] = h
	return nil
}

func (d *stubTelegramDispatcher) UnregisterCommand(name string) {
	norm, err := NormalizeTelegramCommand(name)
	if err != nil {
		return
	}
	delete(d.commands, norm)
}

func (d *stubTelegramDispatcher) RegisterCallback(data string, h TelegramHandler) error {
	if err := ValidateTelegramCallback(data); err != nil {
		return err
	}
	if h == nil {
		return ErrTelegramHandlerNil
	}
	if _, exists := d.callbacks[data]; exists {
		return ErrTelegramCallbackConflict
	}
	d.callbacks[data] = h
	return nil
}

func (d *stubTelegramDispatcher) UnregisterCallback(data string) {
	delete(d.callbacks, data)
}

var _ TelegramDispatcher = (*stubTelegramDispatcher)(nil)

func TestTelegramMessage_Validate(t *testing.T) {
	longCallbackData := strings.Repeat("a", 65)
	validMaxCallbackData := strings.Repeat("a", 64)

	tests := []struct {
		name    string
		msg     TelegramMessage
		wantErr bool
		errIs   error
	}{
		{
			name: "valid minimal message",
			msg: TelegramMessage{
				ChatID: "123456789",
				Text:   "Hello from Telegram bot",
			},
			wantErr: false,
		},
		{
			name: "valid group chat with negative id",
			msg: TelegramMessage{
				ChatID: "-1001234567890",
				Text:   "Group announcement",
			},
			wantErr: false,
		},
		{
			name: "valid with inline keyboard rows and max length callback",
			msg: TelegramMessage{
				ChatID: "12345",
				Text:   "Choose an option:",
				Buttons: [][]TelegramInlineButton{
					{
						{Text: "Option A", CallbackData: "opt_a"},
						{Text: "Option B", CallbackData: validMaxCallbackData},
					},
					{
						{Text: "Cancel", CallbackData: "cancel"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing ChatID",
			msg: TelegramMessage{
				ChatID: "",
				Text:   "Hello",
			},
			wantErr: true,
			errIs:   ErrTelegramMessageInvalid,
		},
		{
			name: "invalid ChatID non-numeric",
			msg: TelegramMessage{
				ChatID: "abc123",
				Text:   "Hello",
			},
			wantErr: true,
			errIs:   ErrTelegramMessageInvalid,
		},
		{
			name: "invalid ChatID special characters",
			msg: TelegramMessage{
				ChatID: "@some_channel",
				Text:   "Hello",
			},
			wantErr: true,
			errIs:   ErrTelegramMessageInvalid,
		},
		{
			name: "empty Text",
			msg: TelegramMessage{
				ChatID: "123456",
				Text:   "",
			},
			wantErr: true,
			errIs:   ErrTelegramMessageInvalid,
		},
		{
			name: "whitespace only Text",
			msg: TelegramMessage{
				ChatID: "123456",
				Text:   "   \t\n  ",
			},
			wantErr: true,
			errIs:   ErrTelegramMessageInvalid,
		},
		{
			name: "button with empty Text",
			msg: TelegramMessage{
				ChatID: "123456",
				Text:   "Menu",
				Buttons: [][]TelegramInlineButton{
					{
						{Text: "", CallbackData: "data_1"},
					},
				},
			},
			wantErr: true,
			errIs:   ErrTelegramMessageInvalid,
		},
		{
			name: "button with whitespace only Text",
			msg: TelegramMessage{
				ChatID: "123456",
				Text:   "Menu",
				Buttons: [][]TelegramInlineButton{
					{
						{Text: "  ", CallbackData: "data_1"},
					},
				},
			},
			wantErr: true,
			errIs:   ErrTelegramMessageInvalid,
		},
		{
			name: "button with empty CallbackData",
			msg: TelegramMessage{
				ChatID: "123456",
				Text:   "Menu",
				Buttons: [][]TelegramInlineButton{
					{
						{Text: "Click", CallbackData: ""},
					},
				},
			},
			wantErr: true,
			errIs:   ErrTelegramMessageInvalid,
		},
		{
			name: "button with CallbackData exceeding 64 bytes",
			msg: TelegramMessage{
				ChatID: "123456",
				Text:   "Menu",
				Buttons: [][]TelegramInlineButton{
					{
						{Text: "Click", CallbackData: longCallbackData},
					},
				},
			},
			wantErr: true,
			errIs:   ErrTelegramMessageInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if tt.errIs != nil && !errors.Is(err, tt.errIs) {
				t.Fatalf("Validate() error = %v, expected errors.Is %v", err, tt.errIs)
			}
		})
	}
}

func TestNormalizeTelegramCommand(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    string
		wantErr bool
		errIs   error
	}{
		{name: "standard with slash", raw: "/start", want: "start"},
		{name: "bare name", raw: "help", want: "help"},
		{name: "with bot username suffix", raw: "/settings@MyTestBot", want: "settings"},
		{name: "with whitespace", raw: "  /info   ", want: "info"},
		{name: "empty string", raw: "", wantErr: true, errIs: ErrTelegramCommandEmpty},
		{name: "slash only", raw: "/", wantErr: true, errIs: ErrTelegramCommandEmpty},
		{name: "slash with whitespace only", raw: "/   ", wantErr: true, errIs: ErrTelegramCommandEmpty},
		{name: "contains space inside", raw: "start now", wantErr: true, errIs: ErrTelegramCommandInvalid},
		{name: "contains inner slash", raw: "foo/bar", wantErr: true, errIs: ErrTelegramCommandInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeTelegramCommand(tt.raw)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NormalizeTelegramCommand(%q) error = %v, wantErr = %v", tt.raw, err, tt.wantErr)
			}
			if tt.errIs != nil && !errors.Is(err, tt.errIs) {
				t.Fatalf("NormalizeTelegramCommand(%q) error = %v, expected errors.Is %v", tt.raw, err, tt.errIs)
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf("NormalizeTelegramCommand(%q) = %q, want %q", tt.raw, got, tt.want)
			}
		})
	}
}

func TestValidateTelegramCallback(t *testing.T) {
	if err := ValidateTelegramCallback(""); !errors.Is(err, ErrTelegramCallbackEmpty) {
		t.Fatalf("expected ErrTelegramCallbackEmpty for empty string, got %v", err)
	}
	if err := ValidateTelegramCallback(strings.Repeat("x", 64)); err != nil {
		t.Fatalf("expected nil for 64-byte callback, got %v", err)
	}
	if err := ValidateTelegramCallback(strings.Repeat("x", 65)); !errors.Is(err, ErrTelegramCallbackTooLong) {
		t.Fatalf("expected ErrTelegramCallbackTooLong for 65-byte callback, got %v", err)
	}
}

func TestStubDispatcher(t *testing.T) {
	ctx := context.Background()
	disp := newStubDispatcher()

	called := false
	handler := func(ctx context.Context, upd TelegramUpdate) error {
		called = true
		return nil
	}

	// Nil handler check
	if err := disp.RegisterCommand("test", nil); !errors.Is(err, ErrTelegramHandlerNil) {
		t.Fatalf("expected ErrTelegramHandlerNil, got %v", err)
	}

	// Register command
	if err := disp.RegisterCommand("/test@bot", handler); err != nil {
		t.Fatalf("unexpected error registering command: %v", err)
	}

	// Conflict check
	if err := disp.RegisterCommand("test", handler); !errors.Is(err, ErrTelegramCommandConflict) {
		t.Fatalf("expected ErrTelegramCommandConflict, got %v", err)
	}

	// Execute registered command
	h, ok := disp.commands["test"]
	if !ok {
		t.Fatalf("expected command 'test' to be found in dispatcher")
	}
	if err := h(ctx, TelegramUpdate{Command: "test"}); err != nil || !called {
		t.Fatalf("handler execution failed or not called")
	}

	// Unregister
	disp.UnregisterCommand("test")
	if _, ok := disp.commands["test"]; ok {
		t.Fatalf("expected command 'test' to be unregistered")
	}

	// Register callback
	if err := disp.RegisterCallback("action_click", handler); err != nil {
		t.Fatalf("unexpected error registering callback: %v", err)
	}
	if err := disp.RegisterCallback("action_click", handler); !errors.Is(err, ErrTelegramCallbackConflict) {
		t.Fatalf("expected ErrTelegramCallbackConflict, got %v", err)
	}
	disp.UnregisterCallback("action_click")
	if _, ok := disp.callbacks["action_click"]; ok {
		t.Fatalf("expected callback 'action_click' to be unregistered")
	}
}
