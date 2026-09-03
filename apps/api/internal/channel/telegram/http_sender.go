package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// DefaultTelegramAPIBaseURL is the standard Telegram Bot API base URL.
const DefaultTelegramAPIBaseURL = "https://api.telegram.org"

// OutboundHTTPTimeout is the strict 10s timeout budget for Bot API calls (D-002 §4).
const OutboundHTTPTimeout = 10 * time.Second

// HTTPSender implements kernel.TelegramSender with standard library net/http.
// When unconfigured, it automatically downgrades to the in-memory CaptureSender.
type HTTPSender struct {
	runtime    *RuntimeManager
	client     *http.Client
	apiBaseURL string
}

var _ kernel.TelegramSender = (*HTTPSender)(nil)

// NewHTTPSender constructs a new HTTPSender.
func NewHTTPSender(runtime *RuntimeManager, client *http.Client, apiBaseURL string) *HTTPSender {
	if client == nil {
		client = &http.Client{Timeout: OutboundHTTPTimeout}
	}
	if strings.TrimSpace(apiBaseURL) == "" {
		apiBaseURL = DefaultTelegramAPIBaseURL
	}
	return &HTTPSender{
		runtime:    runtime,
		client:     client,
		apiBaseURL: strings.TrimRight(apiBaseURL, "/"),
	}
}

// Telegram Bot API sendMessage wire request types.
type sendMessagePayload struct {
	ChatID      string                 `json:"chat_id"`
	Text        string                 `json:"text"`
	ReplyMarkup *inlineKeyboardMarkup  `json:"reply_markup,omitempty"`
}

type inlineKeyboardMarkup struct {
	InlineKeyboard [][]inlineKeyboardButton `json:"inline_keyboard"`
}

type inlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

// Send delivers one message to Telegram Bot API or falls back to CaptureSender.
func (s *HTTPSender) Send(ctx context.Context, msg kernel.TelegramMessage) error {
	// 1. Contract-level validation fails closed immediately.
	if err := msg.Validate(); err != nil {
		return err
	}

	// 2. Read active token from runtime.
	var token string
	if s.runtime != nil {
		token = s.runtime.GetToken()
	}

	// 3. If unconfigured, safely downgrade to memory CaptureSender (判据 #3).
	if token == "" {
		if s.runtime != nil && s.runtime.Mock() != nil {
			return s.runtime.Mock().Send(ctx, msg)
		}
		return nil
	}

	// 4. Construct wire payload.
	payload := sendMessagePayload{
		ChatID: msg.ChatID,
		Text:   msg.Text,
	}
	if len(msg.Buttons) > 0 {
		keyboard := make([][]inlineKeyboardButton, len(msg.Buttons))
		for rIdx, row := range msg.Buttons {
			keyboard[rIdx] = make([]inlineKeyboardButton, len(row))
			for cIdx, btn := range row {
				keyboard[rIdx][cIdx] = inlineKeyboardButton{
					Text:         btn.Text,
					CallbackData: btn.CallbackData,
				}
			}
		}
		payload.ReplyMarkup = &inlineKeyboardMarkup{
			InlineKeyboard: keyboard,
		}
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("telegram: marshal sendMessage payload: %w", err)
	}

	// 5. Send HTTP request with 10s timeout budget.
	sendCtx, cancel := context.WithTimeout(ctx, OutboundHTTPTimeout)
	defer cancel()

	url := fmt.Sprintf("%s/bot%s/sendMessage", s.apiBaseURL, token)
	req, err := http.NewRequestWithContext(sendCtx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("telegram: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("telegram: execute sendMessage: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("telegram: sendMessage failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
