package telegram

import (
	"strings"
	"sync"
)

// RuntimeStatus contains public diagnostic and status information about the Telegram channel.
// Sensitive secrets are never exposed in plaintext.
type RuntimeStatus struct {
	Configured     bool   `json:"configured"`
	TokenMasked    string `json:"token_masked"`
	SecretMasked   string `json:"secret_masked"`
	CapturedCount  int    `json:"captured_count"`
}

// RuntimeManager manages dynamic channel configuration (Bot Token and Webhook Secret)
// with thread-safe hot switching (I-030-005).
type RuntimeManager struct {
	mu     sync.RWMutex
	token  string
	secret string
	mock   *CaptureSender
}

// NewRuntimeManager constructs a RuntimeManager initialized with the given token and secret.
func NewRuntimeManager(token, secret string, mock *CaptureSender) *RuntimeManager {
	if mock == nil {
		mock = NewCaptureSender()
	}
	return &RuntimeManager{
		token:  strings.TrimSpace(token),
		secret: strings.TrimSpace(secret),
		mock:   mock,
	}
}

// GetToken returns the currently active Bot Token.
func (r *RuntimeManager) GetToken() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.token
}

// GetSecret returns the currently active Webhook Secret.
func (r *RuntimeManager) GetSecret() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.secret
}

// Update hot-switches the active Bot Token and Webhook Secret.
func (r *RuntimeManager) Update(token, secret string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.token = strings.TrimSpace(token)
	r.secret = strings.TrimSpace(secret)
}

// Status returns a masked snapshot of the runtime channel configuration.
func (r *RuntimeManager) Status() RuntimeStatus {
	r.mu.RLock()
	defer r.mu.RUnlock()

	captured := 0
	if r.mock != nil {
		captured = len(r.mock.Messages())
	}

	return RuntimeStatus{
		Configured:    r.token != "",
		TokenMasked:   maskSecret(r.token),
		SecretMasked:  maskSecret(r.secret),
		CapturedCount: captured,
	}
}

// Mock returns the underlying capture sender.
func (r *RuntimeManager) Mock() *CaptureSender {
	return r.mock
}

// maskSecret masks a secret string, keeping at most the last 4 characters visible.
func maskSecret(s string) string {
	if s == "" {
		return ""
	}
	if len(s) <= 6 {
		return "******"
	}
	return strings.Repeat("*", len(s)-4) + s[len(s)-4:]
}
