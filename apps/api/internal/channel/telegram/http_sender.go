package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

// ErrTelegramTokenMissing indicates that an HTTPSender has no usable Bot API
// token and no explicit in-memory CaptureSender fallback.
var ErrTelegramTokenMissing = errors.New("telegram: bot token is not configured")

// HTTPSender implements kernel.TelegramSender with standard library net/http.
// When unconfigured, it delegates only to an explicit in-memory CaptureSender;
// without that test/development fallback it fails closed.
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
	ChatID      string                `json:"chat_id"`
	Text        string                `json:"text"`
	ReplyMarkup *inlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

type inlineKeyboardMarkup struct {
	InlineKeyboard [][]inlineKeyboardButton `json:"inline_keyboard"`
}

type inlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

type botAPIResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description,omitempty"`
	ErrorCode   int    `json:"error_code,omitempty"`
}

// Send delivers one message to Telegram Bot API or an explicit CaptureSender.
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

	// 3. Preserve the explicit in-memory test/development fallback, but never
	// report success when there is no runtime or fallback sender.
	if token == "" {
		if s.runtime != nil && s.runtime.Mock() != nil {
			return s.runtime.Mock().Send(ctx, msg)
		}
		return ErrTelegramTokenMissing
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

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 8192))

	if resp.StatusCode != http.StatusOK {
		return newTelegramAPIError("sendMessage", resp.StatusCode, respBody)
	}

	// 6. Check Telegram Bot API response payload "ok": true (R-004 / A-008).
	// Fail closed on a non-JSON 200 body: the Bot API always answers JSON, so a
	// body that cannot be unmarshalled is treated as a failure, not success.
	var apiResp botAPIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return fmt.Errorf("telegram: sendMessage returned non-JSON 200 body: %w: %s", err, strings.TrimSpace(string(respBody)))
	}
	if !apiResp.OK {
		return &TelegramAPIError{
			Method:      "sendMessage",
			HTTPStatus:  resp.StatusCode,
			ErrorCode:   apiResp.ErrorCode,
			Description: apiResp.Description,
		}
	}

	return nil
}
