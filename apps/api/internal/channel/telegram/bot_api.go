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
)

const (
	// GetUpdatesRequestTimeout is the long-poll request budget fixed by D-001.
	GetUpdatesRequestTimeout = 30 * time.Second
	// PollingHTTPClientTimeout is deliberately larger than GetUpdatesRequestTimeout
	// so a normal long-poll response is not reported as a client timeout (D-001).
	PollingHTTPClientTimeout = 40 * time.Second
	// BotAPIResponseBodyLimit bounds error and protocol bodies from the remote API.
	BotAPIResponseBodyLimit = 2 << 20
)

// BotAPIClient is an internal standard-library Telegram Bot API adapter. It
// deliberately stays below the kernel boundary so public contracts carry no
// SDK or HTTP-client types.
type BotAPIClient struct {
	runtime        *RuntimeManager
	client         *http.Client
	apiBaseURL     string
	requestTimeout time.Duration
}

// BotUser is the small getMe result needed by the connection status surface.
type BotUser struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
}

type botAPIEnvelope struct {
	OK          bool            `json:"ok"`
	Result      json.RawMessage `json:"result"`
	Description string          `json:"description,omitempty"`
	ErrorCode   int             `json:"error_code,omitempty"`
}

type setWebhookPayload struct {
	URL         string `json:"url"`
	SecretToken string `json:"secret_token"`
}

type getUpdatesPayload struct {
	Offset  int64 `json:"offset,omitempty"`
	Timeout int   `json:"timeout"`
}

// NewBotAPIClient constructs the short-budget management client used for
// getMe/setWebhook/deleteWebhook. A supplied client is cloned so the adapter
// never mutates a caller-owned http.Client.
func NewBotAPIClient(runtime *RuntimeManager, client *http.Client, apiBaseURL string) *BotAPIClient {
	return newBotAPIClient(runtime, client, apiBaseURL, OutboundHTTPTimeout, OutboundHTTPTimeout)
}

// NewPollingBotAPIClient constructs the dedicated long-polling client. It is
// intentionally separate from the 10-second sendMessage client (D-001).
func NewPollingBotAPIClient(runtime *RuntimeManager, client *http.Client, apiBaseURL string) *BotAPIClient {
	return newBotAPIClient(runtime, client, apiBaseURL, GetUpdatesRequestTimeout, PollingHTTPClientTimeout)
}

func newBotAPIClient(runtime *RuntimeManager, client *http.Client, apiBaseURL string, requestTimeout, clientTimeout time.Duration) *BotAPIClient {
	if client == nil {
		client = &http.Client{}
	} else {
		clone := *client
		client = &clone
	}
	client.Timeout = clientTimeout
	if strings.TrimSpace(apiBaseURL) == "" {
		apiBaseURL = DefaultTelegramAPIBaseURL
	}
	return &BotAPIClient{
		runtime:        runtime,
		client:         client,
		apiBaseURL:     strings.TrimRight(strings.TrimSpace(apiBaseURL), "/"),
		requestTimeout: requestTimeout,
	}
}

// GetMe verifies the token and returns the bot identity without exposing the
// token in any returned error.
func (c *BotAPIClient) GetMe(ctx context.Context) (BotUser, error) {
	var user BotUser
	if err := c.call(ctx, "getMe", nil, &user); err != nil {
		return BotUser{}, err
	}
	return user, nil
}

// SetWebhook configures Telegram's webhook target and its secret token.
func (c *BotAPIClient) SetWebhook(ctx context.Context, webhookURL, secret string) error {
	return c.call(ctx, "setWebhook", setWebhookPayload{
		URL:         webhookURL,
		SecretToken: secret,
	}, nil)
}

// DeleteWebhook removes Telegram's remote webhook configuration.
func (c *BotAPIClient) DeleteWebhook(ctx context.Context) error {
	return c.call(ctx, "deleteWebhook", nil, nil)
}

// GetUpdates performs one long-poll request. An empty result is a normal
// response and is returned as an empty slice with nil error.
func (c *BotAPIClient) GetUpdates(ctx context.Context, offset int64) ([]UpdatePayload, error) {
	var updates []UpdatePayload
	if err := c.call(ctx, "getUpdates", getUpdatesPayload{
		Offset:  offset,
		Timeout: int(GetUpdatesRequestTimeout / time.Second),
	}, &updates); err != nil {
		return nil, err
	}
	if updates == nil {
		updates = []UpdatePayload{}
	}
	return updates, nil
}

func (c *BotAPIClient) call(ctx context.Context, method string, payload any, result any) error {
	if c == nil || c.runtime == nil {
		return fmt.Errorf("telegram: %s: runtime is unavailable", method)
	}
	token := strings.TrimSpace(c.runtime.GetToken())
	if token == "" {
		return fmt.Errorf("telegram: %s: bot token is not configured", method)
	}

	var body io.Reader
	if payload != nil {
		bodyBytes, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("telegram: %s: marshal request: %w", method, err)
		}
		body = bytes.NewReader(bodyBytes)
	}

	requestCtx := ctx
	var cancel context.CancelFunc
	if c.requestTimeout > 0 {
		requestCtx, cancel = context.WithTimeout(ctx, c.requestTimeout)
		defer cancel()
	}
	requestURL := fmt.Sprintf("%s/bot%s/%s", c.apiBaseURL, token, method)
	req, err := http.NewRequestWithContext(requestCtx, http.MethodPost, requestURL, body)
	if err != nil {
		return fmt.Errorf("telegram: %s: create request: %w", method, err)
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("telegram: %s: execute request: %w", method, err)
	}
	defer resp.Body.Close()
	responseBody, readErr := io.ReadAll(io.LimitReader(resp.Body, BotAPIResponseBodyLimit))
	if readErr != nil {
		return fmt.Errorf("telegram: %s: read response: %w", method, readErr)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram: %s failed with HTTP status %d", method, resp.StatusCode)
	}

	var envelope botAPIEnvelope
	if err := json.Unmarshal(responseBody, &envelope); err != nil {
		return fmt.Errorf("telegram: %s: non-JSON response: %w", method, err)
	}
	if !envelope.OK {
		return fmt.Errorf("telegram: %s returned ok=false: %s (code %d)", method, envelope.Description, envelope.ErrorCode)
	}
	if result == nil {
		return nil
	}
	if len(envelope.Result) == 0 || string(envelope.Result) == "null" {
		return fmt.Errorf("telegram: %s: successful response has no result", method)
	}
	if err := json.Unmarshal(envelope.Result, result); err != nil {
		return fmt.Errorf("telegram: %s: decode result: %w", method, err)
	}
	return nil
}
