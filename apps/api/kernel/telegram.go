package kernel

// Kernel Telegram channel port (VP-030 / workspace-030 GOAL-002 D-002, R1).
//
// The port defines the in-process Telegram channel contract for the kernel and
// modules: outbound message delivery with optional callback buttons (TelegramSender),
// inbound command/callback registration and dispatch (TelegramDispatcher), and
// thin Update representations. Public types carry no third-party Telegram SDK
// or HTTP client types. Domain and extension code consumes TelegramSender /
// TelegramDispatcher, while webhook HTTP routing and Bot API client details
// stay inside internal/channel/telegram adapters.
//
// Contract frozen by workspace-030 GOAL-002 D-002:
//
//   - Send is synchronous: no queue, no background retry. Failures return as
//     errors for the caller to handle (mirroring MailSender).
//   - TelegramMessage supports plain text and optional inline keyboard rows
//     containing callback_data buttons only (URL/Login buttons are prohibited).
//   - TelegramUpdate is a thin inbound view (ChatID, UserID, SubjectID, Command,
//     Text, CallbackData) without raw SDK payloads.
//   - TelegramDispatcher provides static command and callback handler registration.
//     Duplicate registrations fail closed; unknown commands fall back cleanly.
//   - When channel.telegram is disabled, the composition root provides a no-op
//     or fail-closed stub (Send returns ErrTelegramDisabled).

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// DefaultTelegramUnknownCommandText is the product-frozen fallback message sent
// when an incoming command is not registered by any handler (GOAL-002 D-002 §3).
const DefaultTelegramUnknownCommandText = "Sorry, unrecognized command."

// TelegramCallbackDataMaxBytes is the Telegram Bot API upper bound on callback_data payload size.
const TelegramCallbackDataMaxBytes = 64

// chatIDPattern validates decimal chat ID strings (e.g. "123456789" or "-100123456789").
var chatIDPattern = regexp.MustCompile(`^-?[0-9]+$`)

// Sentinel errors (GOAL-002 D-002 §1/§3/§4). Callers use errors.Is.
var (
	ErrTelegramDisabled         = errors.New("kernel: telegram channel disabled")
	ErrTelegramHandlerNil       = errors.New("kernel: telegram handler is nil")
	ErrTelegramCommandEmpty     = errors.New("kernel: telegram command name is empty")
	ErrTelegramCommandInvalid   = errors.New("kernel: telegram command name is invalid")
	ErrTelegramCommandConflict  = errors.New("kernel: telegram command already registered")
	ErrTelegramCallbackEmpty    = errors.New("kernel: telegram callback data is empty")
	ErrTelegramCallbackConflict = errors.New("kernel: telegram callback already registered")
	ErrTelegramCallbackTooLong  = errors.New("kernel: telegram callback data exceeds maximum allowed bytes")
	ErrTelegramMessageInvalid   = errors.New("kernel: telegram message is invalid")
)

// TelegramInlineButton is an inline keyboard button with a callback_data payload.
// Prohibited by D-002: URL, Login, Switch, WebApp buttons.
type TelegramInlineButton struct {
	Text         string
	CallbackData string
}

// TelegramMessage is one outbound Telegram message.
type TelegramMessage struct {
	ChatID  string
	Text    string
	Buttons [][]TelegramInlineButton // Optional; nil or empty for no inline keyboard
}

// Validate enforces contract-level outbound rules (GOAL-002 D-002 §4):
// - ChatID must be non-empty and well-formed decimal digits (optional leading minus).
// - Text must be non-empty (Telegram rejects empty text messages).
// - Each button must have non-empty Text and CallbackData.
// - CallbackData must not exceed 64 bytes.
func (m TelegramMessage) Validate() error {
	chatID := strings.TrimSpace(m.ChatID)
	if chatID == "" || !chatIDPattern.MatchString(chatID) {
		return fmt.Errorf("%w: invalid or missing ChatID %q", ErrTelegramMessageInvalid, m.ChatID)
	}
	if strings.TrimSpace(m.Text) == "" {
		return fmt.Errorf("%w: message Text must not be empty", ErrTelegramMessageInvalid)
	}
	for rowIdx, row := range m.Buttons {
		for colIdx, btn := range row {
			if strings.TrimSpace(btn.Text) == "" {
				return fmt.Errorf("%w: button at [%d][%d] has empty Text", ErrTelegramMessageInvalid, rowIdx, colIdx)
			}
			if btn.CallbackData == "" {
				return fmt.Errorf("%w: button at [%d][%d] has empty CallbackData", ErrTelegramMessageInvalid, rowIdx, colIdx)
			}
			if len([]byte(btn.CallbackData)) > TelegramCallbackDataMaxBytes {
				return fmt.Errorf("%w: button at [%d][%d] CallbackData exceeds %d bytes (got %d)",
					ErrTelegramMessageInvalid, rowIdx, colIdx, TelegramCallbackDataMaxBytes, len([]byte(btn.CallbackData)))
			}
		}
	}
	return nil
}

// TelegramSender is the kernel outbound Telegram messaging port (R1).
// Implementations validate the message before sending, failing closed on violations.
type TelegramSender interface {
	Send(ctx context.Context, msg TelegramMessage) error
}

// TelegramUpdate is a thin inbound view delivered to registered handlers.
// Prohibited from carrying raw SDK types.
type TelegramUpdate struct {
	ChatID       string // Telegram chat_id, decimal string (may be negative for groups/channels)
	UserID       string // Telegram user_id, decimal string
	SubjectID    string // Populated after identity mapping in R2; allowed empty in R1
	Command      string // Command name without leading slash or @bot suffix; empty for non-commands
	Text         string // Message text (verbatim); may be empty for callback queries
	CallbackData string // callback_query.data; empty for standard messages
}

// TelegramHandler is the callback function for an inbound Telegram update.
type TelegramHandler func(ctx context.Context, upd TelegramUpdate) error

// TelegramDispatcher provides registration for Telegram command and callback handlers.
type TelegramDispatcher interface {
	RegisterCommand(name string, h TelegramHandler) error
	UnregisterCommand(name string)
	RegisterCallback(data string, h TelegramHandler) error
	UnregisterCallback(data string)
}

// NormalizeTelegramCommand normalizes a raw command string by stripping a leading
// slash and any optional @BotName suffix (e.g. "/start@MyBot" -> "start").
// Returns ErrTelegramCommandEmpty or ErrTelegramCommandInvalid if malformed.
func NormalizeTelegramCommand(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	s = strings.TrimPrefix(s, "/")
	if atIdx := strings.Index(s, "@"); atIdx >= 0 {
		s = s[:atIdx]
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return "", ErrTelegramCommandEmpty
	}
	if strings.Contains(s, "/") || strings.Contains(s, " ") {
		return "", ErrTelegramCommandInvalid
	}
	return s, nil
}

// ValidateTelegramCallback validates callback data (non-empty and <= 64 bytes).
func ValidateTelegramCallback(data string) error {
	if data == "" {
		return ErrTelegramCallbackEmpty
	}
	if len([]byte(data)) > TelegramCallbackDataMaxBytes {
		return ErrTelegramCallbackTooLong
	}
	return nil
}
