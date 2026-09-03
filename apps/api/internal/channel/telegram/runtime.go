package telegram

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/mail"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// defaultMasterKey is the deterministic fallback key when no explicit master key is supplied.
var defaultMasterKey = sha256.Sum256([]byte("schema-ui-core:channel:telegram:master-key:v1"))

// RuntimeStatus contains diagnostic and status information about the Telegram channel.
// Sensitive secrets are never exposed in plaintext or partial masks (F-002 / R-005 / R-008).
type RuntimeStatus struct {
	Configured            bool `json:"configured"`
	TokenSet              bool `json:"token_set"`
	SecretSet             bool `json:"secret_set"`
	CapturedMessagesCount int  `json:"captured_messages_count"`
	CapturedCount         int  `json:"captured_count"` // backwards compatibility alias
}

// TxRunner is the persistence boundary for storing runtime channel configurations.
type TxRunner interface {
	Run(context.Context, func(kernel.Tx) error) error
}

// RuntimeManager manages dynamic channel configuration (Bot Token and Webhook Secret)
// with thread-safe hot switching (I-030-005) and persistent encrypted database storage (F-002).
type RuntimeManager struct {
	mu        sync.RWMutex
	token     string
	secret    string
	mock      *CaptureSender
	runner    TxRunner
	masterKey []byte
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
		token:     strings.TrimSpace(seedToken),
		secret:    strings.TrimSpace(seedSecret),
		mock:      mock,
		runner:    runner,
		masterKey: defaultMasterKey[:],
	}

	if runner != nil {
		rm.initPersistence(seedToken, seedSecret)
	}

	return rm
}

func (r *RuntimeManager) initPersistence(seedToken, seedSecret string) {
	ctx := context.Background()
	_ = r.runner.Run(ctx, func(tx kernel.Tx) error {
		var count int
		row := tx.QueryRow(ctx, `SELECT COUNT(*) FROM telegram_config WHERE id = 1`)
		if err := row.Scan(&count); err != nil || count == 0 {
			// Seed row if table exists but empty
			tokenEnc, err1 := mail.EncryptSecret(r.masterKey, strings.TrimSpace(seedToken))
			secretEnc, err2 := mail.EncryptSecret(r.masterKey, strings.TrimSpace(seedSecret))
			if err1 == nil && err2 == nil {
				now := time.Now().Unix()
				_, _ = tx.Exec(ctx, `INSERT INTO telegram_config (id, bot_token_enc, webhook_secret_enc, updated_at) VALUES (1, ?, ?, ?)`,
					tokenEnc, secretEnc, now)
			}
			return nil
		}

		var dbTokenEnc, dbSecretEnc string
		var updatedAt int64
		row2 := tx.QueryRow(ctx, `SELECT bot_token_enc, webhook_secret_enc, updated_at FROM telegram_config WHERE id = 1`)
		if err := row2.Scan(&dbTokenEnc, &dbSecretEnc, &updatedAt); err == nil {
			decToken, err1 := mail.DecryptSecret(r.masterKey, dbTokenEnc)
			decSecret, err2 := mail.DecryptSecret(r.masterKey, dbSecretEnc)
			if err1 == nil && err2 == nil {
				r.mu.Lock()
				if strings.TrimSpace(decToken) != "" {
					r.token = strings.TrimSpace(decToken)
				}
				if strings.TrimSpace(decSecret) != "" {
					r.secret = strings.TrimSpace(decSecret)
				}
				r.mu.Unlock()
			}
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

// Update hot-switches the active Bot Token and Webhook Secret.
// Persists encrypted secrets to the database before modifying in-memory state (fail-closed).
func (r *RuntimeManager) Update(ctx context.Context, token, secret string) error {
	trimmedToken := strings.TrimSpace(token)
	trimmedSecret := strings.TrimSpace(secret)

	if r.runner != nil {
		tokenEnc, err := mail.EncryptSecret(r.masterKey, trimmedToken)
		if err != nil {
			return fmt.Errorf("telegram: encrypt token: %w", err)
		}
		secretEnc, err := mail.EncryptSecret(r.masterKey, trimmedSecret)
		if err != nil {
			return fmt.Errorf("telegram: encrypt secret: %w", err)
		}

		now := time.Now().Unix()
		err = r.runner.Run(ctx, func(tx kernel.Tx) error {
			_, err := tx.Exec(ctx, `INSERT INTO telegram_config (id, bot_token_enc, webhook_secret_enc, updated_at)
				VALUES (1, ?, ?, ?)
				ON CONFLICT(id) DO UPDATE SET
					bot_token_enc = excluded.bot_token_enc,
					webhook_secret_enc = excluded.webhook_secret_enc,
					updated_at = excluded.updated_at`,
				tokenEnc, secretEnc, now)
			return err
		})
		if err != nil {
			return fmt.Errorf("telegram: persist config update: %w", err)
		}
	}

	// Memory state only updates after persistence succeeds (fail-closed).
	r.mu.Lock()
	r.token = trimmedToken
	r.secret = trimmedSecret
	r.mu.Unlock()

	return nil
}

// Status returns a masked snapshot of the runtime channel configuration.
// No secret fragments or partial masks are returned (R-005).
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
		CapturedMessagesCount: captured,
		CapturedCount:         captured,
	}
}

// Mock returns the underlying capture sender.
func (r *RuntimeManager) Mock() *CaptureSender {
	return r.mock
}
