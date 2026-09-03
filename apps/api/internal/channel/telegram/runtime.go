package telegram

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// RuntimeStatus contains public diagnostic and status information about the Telegram channel.
// Sensitive secrets are never exposed in plaintext.
type RuntimeStatus struct {
	Configured            bool   `json:"configured"`
	TokenSet              bool   `json:"token_set"`
	SecretSet             bool   `json:"secret_set"`
	TokenMasked           string `json:"token_masked,omitempty"`
	SecretMasked          string `json:"secret_masked,omitempty"`
	CapturedMessagesCount int    `json:"captured_messages_count"`
	CapturedCount         int    `json:"captured_count"`
}

// TxRunner is the persistence boundary for storing runtime channel configurations.
type TxRunner interface {
	Run(context.Context, func(kernel.Tx) error) error
}

// RuntimeManager manages dynamic channel configuration (Bot Token and Webhook Secret)
// with thread-safe hot switching (I-030-005) and persistent database storage (F-002).
type RuntimeManager struct {
	mu     sync.RWMutex
	token  string
	secret string
	mock   *CaptureSender
	runner TxRunner
}

// NewRuntimeManager constructs a RuntimeManager initialized with the given token and secret,
// loading any persisted state from the database if available.
func NewRuntimeManager(seedToken, seedSecret string, mock *CaptureSender, runners ...TxRunner) *RuntimeManager {
	if mock == nil {
		mock = NewCaptureSender()
	}

	var runner TxRunner
	if len(runners) > 0 {
		runner = runners[0]
	}

	rm := &RuntimeManager{
		token:  strings.TrimSpace(seedToken),
		secret: strings.TrimSpace(seedSecret),
		mock:   mock,
		runner: runner,
	}

	if runner != nil {
		rm.initPersistence(seedToken, seedSecret)
	}

	return rm
}

func (r *RuntimeManager) initPersistence(seedToken, seedSecret string) {
	ctx := context.Background()
	_ = r.runner.Run(ctx, func(tx kernel.Tx) error {
		_, err := tx.Exec(ctx, `CREATE TABLE IF NOT EXISTS telegram_config (
			id INTEGER PRIMARY KEY,
			bot_token TEXT NOT NULL,
			webhook_secret TEXT NOT NULL,
			updated_at BIGINT NOT NULL
		)`)
		if err != nil {
			return err
		}

		var count int
		row := tx.QueryRow(ctx, `SELECT COUNT(*) FROM telegram_config WHERE id = 1`)
		if err := row.Scan(&count); err != nil || count == 0 {
			now := time.Now().Unix()
			_, err = tx.Exec(ctx, `INSERT INTO telegram_config (id, bot_token, webhook_secret, updated_at) VALUES (1, ?, ?, ?)`,
				strings.TrimSpace(seedToken), strings.TrimSpace(seedSecret), now)
			return err
		}

		var dbToken, dbSecret string
		var updatedAt int64
		row2 := tx.QueryRow(ctx, `SELECT bot_token, webhook_secret, updated_at FROM telegram_config WHERE id = 1`)
		if err := row2.Scan(&dbToken, &dbSecret, &updatedAt); err == nil {
			r.mu.Lock()
			// If DB has configuration, it overrides initial seed
			if strings.TrimSpace(dbToken) != "" {
				r.token = strings.TrimSpace(dbToken)
			}
			if strings.TrimSpace(dbSecret) != "" {
				r.secret = strings.TrimSpace(dbSecret)
			}
			r.mu.Unlock()
		}
		return nil
	})
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

// Update hot-switches the active Bot Token and Webhook Secret in memory and persists to database.
func (r *RuntimeManager) Update(ctx context.Context, token, secret string) error {
	trimmedToken := strings.TrimSpace(token)
	trimmedSecret := strings.TrimSpace(secret)

	r.mu.Lock()
	r.token = trimmedToken
	r.secret = trimmedSecret
	r.mu.Unlock()

	if r.runner != nil {
		now := time.Now().Unix()
		err := r.runner.Run(ctx, func(tx kernel.Tx) error {
			_, err := tx.Exec(ctx, `INSERT INTO telegram_config (id, bot_token, webhook_secret, updated_at)
				VALUES (1, ?, ?, ?)
				ON CONFLICT(id) DO UPDATE SET
					bot_token = excluded.bot_token,
					webhook_secret = excluded.webhook_secret,
					updated_at = excluded.updated_at`,
				trimmedToken, trimmedSecret, now)
			return err
		})
		if err != nil {
			return fmt.Errorf("telegram: persist config update: %w", err)
		}
	}

	return nil
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
		Configured:            r.token != "",
		TokenSet:              r.token != "",
		SecretSet:             r.secret != "",
		TokenMasked:           maskSecret(r.token),
		SecretMasked:          maskSecret(r.secret),
		CapturedMessagesCount: captured,
		CapturedCount:         captured,
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
