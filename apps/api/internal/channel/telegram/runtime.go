package telegram

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/mail"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

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
// loading any persisted state from the database if available. masterKey is the at-rest
// encryption key — required, and NEVER derived from a source constant (F-002 / A-006:
// "主密钥离开源码"). A nil/empty key is a construction error (fail-closed).
func NewRuntimeManager(seedToken, seedSecret string, mock *CaptureSender, masterKey []byte, runners ...TxRunner) (*RuntimeManager, error) {
	if len(masterKey) == 0 {
		return nil, fmt.Errorf("telegram: master key is required")
	}
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
		masterKey: masterKey,
	}

	if runner != nil {
		if err := rm.initPersistence(seedToken, seedSecret); err != nil {
			return nil, fmt.Errorf("telegram: init persistence: %w", err)
		}
	}

	return rm, nil
}

// initPersistence loads (or seeds) the encrypted telegram_config row. Any DB or
// decryption failure is returned so the composition root fails closed instead of
// silently staying on the seed values (F-002 / A-006).
//
// Once a row exists it is authoritative: the decrypted values are applied
// verbatim, including empty ones, so an admin clearing a token/secret survives
// restart instead of reverting to the env seed (A-008 informational).
func (r *RuntimeManager) initPersistence(seedToken, seedSecret string) error {
	ctx := context.Background()
	return r.runner.Run(ctx, func(tx kernel.Tx) error {
		var count int
		row := tx.QueryRow(ctx, `SELECT COUNT(*) FROM telegram_config WHERE id = 1`)
		if err := row.Scan(&count); err != nil {
			return fmt.Errorf("telegram: read config presence: %w", err)
		}
		if count == 0 {
			// Seed row if the table is empty.
			tokenEnc, err1 := mail.EncryptSecret(r.masterKey, strings.TrimSpace(seedToken))
			secretEnc, err2 := mail.EncryptSecret(r.masterKey, strings.TrimSpace(seedSecret))
			if err1 != nil || err2 != nil {
				return fmt.Errorf("telegram: encrypt seed secrets: %v / %v", err1, err2)
			}
			now := time.Now().Unix()
			if _, err := tx.Exec(ctx, `INSERT INTO telegram_config (id, bot_token_enc, webhook_secret_enc, updated_at) VALUES (1, ?, ?, ?)`,
				tokenEnc, secretEnc, now); err != nil {
				return fmt.Errorf("telegram: seed config: %w", err)
			}
			return nil
		}

		var dbTokenEnc, dbSecretEnc string
		var updatedAt int64
		row2 := tx.QueryRow(ctx, `SELECT bot_token_enc, webhook_secret_enc, updated_at FROM telegram_config WHERE id = 1`)
		if err := row2.Scan(&dbTokenEnc, &dbSecretEnc, &updatedAt); err != nil {
			return fmt.Errorf("telegram: read persisted config: %w", err)
		}
		decToken, err1 := mail.DecryptSecret(r.masterKey, dbTokenEnc)
		decSecret, err2 := mail.DecryptSecret(r.masterKey, dbSecretEnc)
		if err1 != nil || err2 != nil {
			return fmt.Errorf("telegram: decrypt persisted config: %v / %v", err1, err2)
		}
		r.mu.Lock()
		r.token = strings.TrimSpace(decToken)
		r.secret = strings.TrimSpace(decSecret)
		r.mu.Unlock()
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
