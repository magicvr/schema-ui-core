package mail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// DefaultResendBaseURL is the production Resend HTTP API root
// (workspace-017 GOAL-007; contract frozen by GOAL-006 D-002 §4). Tests and
// the harness-equivalent live check override it via ResendOptions.BaseURL —
// the override swaps WHICH endpoint receives the request, never the request
// shape or the auth scheme.
const DefaultResendBaseURL = "https://api.resend.com"

// ResendOptions carries the explicitly configured production channel.
// APIKey/From are required (fail-closed at NewResend, mirroring NewSMTP).
// BaseURL and Client are test seams: empty/nil select the production values.
type ResendOptions struct {
	APIKey string // SECRET: arrives via MAIL_RESEND_API_KEY env / configs/.env
	From   string // default sender stamped on every message (bare addr-spec)

	BaseURL string
	Client  *http.Client
}

// Resend is the production kernel.MailSender adapter over the Resend HTTP
// API: one POST per message (synchronous, no queue, no retry), matching the
// R1 port contract. It is safe for concurrent use: every Send issues its own
// HTTP request through a shared http.Client (which is concurrency-safe).
type Resend struct {
	apiKey  string
	from    string
	baseURL string
	client  *http.Client
}

// NewResend validates options fail-closed and returns the adapter.
func NewResend(opts ResendOptions) (*Resend, error) {
	for _, pair := range []struct{ name, value string }{
		{"api key", opts.APIKey},
		{"from", opts.From},
	} {
		if strings.TrimSpace(pair.value) == "" {
			return nil, fmt.Errorf("mail: explicit Resend requires %s", pair.name)
		}
	}
	from := strings.TrimSpace(opts.From)
	parsed, err := mail.ParseAddress(from)
	if err != nil || parsed.Address != from {
		return nil, fmt.Errorf("mail: Resend from must be a bare address")
	}
	baseURL := strings.TrimRight(strings.TrimSpace(opts.BaseURL), "/")
	if baseURL == "" {
		baseURL = DefaultResendBaseURL
	}
	client := opts.Client
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}
	return &Resend{apiKey: opts.APIKey, from: from, baseURL: baseURL, client: client}, nil
}

// resendEmailRequest is the /v1-ish JSON body of POST {base}/emails. The From
// header never travels on MailMessage: the configured sender is stamped here,
// at the adapter boundary (R1 contract).
type resendEmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Text    string   `json:"text"`
}

// Send delivers one message synchronously:
//
//	POST {base}/emails  Authorization: Bearer <api key>
//	  {"from": <configured from>, "to": [<msg.To>], "subject", "text"}
//
// A 2xx status accepts delivery (the port contract ends there); any non-2xx
// or transport failure returns an error naming the STATUS, never the key.
func (r *Resend) Send(ctx context.Context, msg kernel.MailMessage) error {
	if err := msg.Validate(); err != nil {
		return fmt.Errorf("mail: %v", err)
	}
	body, err := json.Marshal(resendEmailRequest{
		From:    r.from,
		To:      []string{msg.To},
		Subject: msg.Subject,
		Text:    msg.TextBody,
	})
	if err != nil {
		return fmt.Errorf("mail: encode resend request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.baseURL+"/emails", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("mail: build resend request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+r.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("mail: resend transport: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		snippet, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("mail: resend rejected with status %d: %s", resp.StatusCode, strings.TrimSpace(string(snippet)))
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	return nil
}
