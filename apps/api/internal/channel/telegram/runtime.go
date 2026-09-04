package telegram

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/mail"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

const (
	TelegramModePolling = "polling"
	TelegramModeWebhook = "webhook"
)

const (
	ConnectionStateUnconfigured = "unconfigured"
	ConnectionStateStarting     = "starting"
	ConnectionStateRunning      = "running"
	ConnectionStateStopping     = "stopping"
	ConnectionStateError        = "error"
	ConnectionStateIdle         = "idle"

	ReceiverNone    = "none"
	ReceiverWebhook = "webhook"
	ReceiverPolling = "polling"
)

// ConnectionStatus is the non-secret operational state of the Telegram
// receiver. It is intentionally separate from the kernel Telegram contract.
type ConnectionStatus struct {
	State       string
	Receiver    string
	BotID       int64
	BotUsername string
	LastError   string
}

// RuntimeStatus contains diagnostic and status information about the Telegram channel.
// Sensitive secrets are never exposed in plaintext or partial masks (F-002 / R-005 / R-008).
type RuntimeStatus struct {
	Configured            bool   `json:"configured"`
	TokenSet              bool   `json:"token_set"`
	SecretSet             bool   `json:"secret_set"`
	Mode                  string `json:"mode"`
	WebhookPublicBaseURL  string `json:"webhook_public_base_url"`
	ConnectionState       string `json:"connection_state"`
	Receiver              string `json:"receiver"`
	BotID                 int64  `json:"bot_id,omitempty"`
	BotUsername           string `json:"bot_username,omitempty"`
	LastError             string `json:"last_error,omitempty"`
	CapturedMessagesCount int    `json:"captured_messages_count"`
	CapturedCount         int    `json:"captured_count"` // backwards compatibility alias
}

// TxRunner is the persistence boundary for storing runtime channel configurations.
type TxRunner interface {
	Run(context.Context, func(kernel.Tx) error) error
}

// RuntimeManager manages dynamic channel configuration (Bot Token and Webhook Secret)
// with thread-safe hot switching (I-030-005) and persistent encrypted database storage (F-002).
type RuntimeManager struct {
	mu                   sync.RWMutex
	updateMu             sync.Mutex
	token                string
	secret               string
	mode                 string
	webhookPublicBaseURL string
	connectionState      string
	receiver             string
	botID                int64
	botUsername          string
	lastError            string
	mock                 *CaptureSender
	runner               TxRunner
	masterKey            []byte
	settingsChangedMu    sync.RWMutex
	settingsChanged      func(context.Context) error
}

// NewRuntimeManager constructs a RuntimeManager initialized with the given token and secret,
// loading any persisted state from the database if available. masterKey is the at-rest
// encryption key — required, and NEVER derived from a source constant (F-002 / A-006:
// "主密钥离开源码"). A nil/empty key is a construction error (fail-closed).
func NewRuntimeManager(seedToken, seedSecret string, mock *CaptureSender, masterKey []byte, runners ...TxRunner) (*RuntimeManager, error) {
	return NewRuntimeManagerWithSettings(seedToken, seedSecret, TelegramModePolling, "", mock, masterKey, runners...)
}

// NewRuntimeManagerWithSettings is the R2 constructor. YAML/env values are
// seeds only; once the singleton DB row exists, its values are authoritative.
func NewRuntimeManagerWithSettings(seedToken, seedSecret, seedMode, seedWebhookPublicBaseURL string, mock *CaptureSender, masterKey []byte, runners ...TxRunner) (*RuntimeManager, error) {
	if len(masterKey) == 0 {
		return nil, fmt.Errorf("telegram: master key is required")
	}
	if mock == nil {
		mock = NewCaptureSender()
	}
	seedMode = normalizeTelegramMode(seedMode)
	if !ValidTelegramMode(seedMode) {
		return nil, fmt.Errorf("telegram: mode must be polling or webhook (got %q)", seedMode)
	}
	seedWebhookPublicBaseURL = strings.TrimSpace(seedWebhookPublicBaseURL)
	if err := validateWebhookPublicBaseURL(seedWebhookPublicBaseURL); err != nil {
		return nil, err
	}

	var runner TxRunner
	if len(runners) > 0 {
		runner = runners[0]
	}

	trimmedSeedToken := strings.TrimSpace(seedToken)
	initialConnectionState := ConnectionStateUnconfigured
	if trimmedSeedToken != "" {
		initialConnectionState = ConnectionStateIdle
	}
	rm := &RuntimeManager{
		token:                trimmedSeedToken,
		secret:               strings.TrimSpace(seedSecret),
		mode:                 seedMode,
		webhookPublicBaseURL: seedWebhookPublicBaseURL,
		connectionState:      initialConnectionState,
		receiver:             ReceiverNone,
		mock:                 mock,
		runner:               runner,
		masterKey:            masterKey,
	}

	if runner != nil {
		if err := rm.initPersistence(seedToken, seedSecret, seedMode, seedWebhookPublicBaseURL); err != nil {
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
func (r *RuntimeManager) initPersistence(seedToken, seedSecret, seedMode, seedWebhookPublicBaseURL string) error {
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
			if _, err := tx.Exec(ctx, `INSERT INTO telegram_config (id, bot_token_enc, webhook_secret_enc, mode, webhook_public_base_url, updated_at) VALUES (1, ?, ?, ?, ?, ?)`,
				tokenEnc, secretEnc, seedMode, seedWebhookPublicBaseURL, now); err != nil {
				return fmt.Errorf("telegram: seed config: %w", err)
			}
			return nil
		}

		var dbTokenEnc, dbSecretEnc, dbMode, dbWebhookPublicBaseURL string
		var updatedAt int64
		row2 := tx.QueryRow(ctx, `SELECT bot_token_enc, webhook_secret_enc, mode, webhook_public_base_url, updated_at FROM telegram_config WHERE id = 1`)
		if err := row2.Scan(&dbTokenEnc, &dbSecretEnc, &dbMode, &dbWebhookPublicBaseURL, &updatedAt); err != nil {
			return fmt.Errorf("telegram: read persisted config: %w", err)
		}
		decToken, err1 := mail.DecryptSecret(r.masterKey, dbTokenEnc)
		decSecret, err2 := mail.DecryptSecret(r.masterKey, dbSecretEnc)
		if err1 != nil || err2 != nil {
			return fmt.Errorf("telegram: decrypt persisted config: %v / %v", err1, err2)
		}
		dbMode = normalizeTelegramMode(dbMode)
		if !ValidTelegramMode(dbMode) {
			return fmt.Errorf("telegram: persisted mode is invalid: %q", dbMode)
		}
		dbWebhookPublicBaseURL = strings.TrimSpace(dbWebhookPublicBaseURL)
		if err := validateWebhookPublicBaseURL(dbWebhookPublicBaseURL); err != nil {
			return fmt.Errorf("telegram: persisted webhook URL: %w", err)
		}
		r.mu.Lock()
		r.token = strings.TrimSpace(decToken)
		r.secret = strings.TrimSpace(decSecret)
		r.mode = dbMode
		r.webhookPublicBaseURL = dbWebhookPublicBaseURL
		if r.token == "" {
			r.connectionState = ConnectionStateUnconfigured
		} else {
			r.connectionState = ConnectionStateIdle
		}
		r.receiver = ReceiverNone
		r.botID = 0
		r.botUsername = ""
		r.lastError = ""
		r.mu.Unlock()
		return nil
	})
}

func normalizeTelegramMode(mode string) string {
	mode = strings.ToLower(strings.TrimSpace(mode))
	if mode == "" {
		return TelegramModePolling
	}
	return mode
}

func ValidTelegramMode(mode string) bool {
	switch normalizeTelegramMode(mode) {
	case TelegramModePolling, TelegramModeWebhook:
		return true
	default:
		return false
	}
}

func validateWebhookPublicBaseURL(raw string) error {
	value := strings.TrimSpace(raw)
	if value == "" {
		return nil
	}
	parsed, err := url.Parse(value)
	if err != nil || value != raw || (parsed.Scheme != "http" && parsed.Scheme != "https") ||
		parsed.Host == "" || parsed.User != nil || parsed.Path != "" ||
		parsed.RawQuery != "" || parsed.Fragment != "" || strings.ContainsAny(value, " \t\r\n") {
		return fmt.Errorf("telegram: webhook_public_base_url must be an absolute http(s) origin with no path, query, fragment, credentials, or whitespace (got %q)", raw)
	}
	return nil
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

// GetMode returns the currently persisted Telegram receiver mode.
func (r *RuntimeManager) GetMode() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.mode
}

// GetWebhookPublicBaseURL returns the explicit webhook origin.
func (r *RuntimeManager) GetWebhookPublicBaseURL() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.webhookPublicBaseURL
}

// ConnectionStatus returns the current non-secret receiver state.
func (r *RuntimeManager) ConnectionStatus() ConnectionStatus {
	if r == nil {
		return ConnectionStatus{State: ConnectionStateUnconfigured, Receiver: ReceiverNone}
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	return ConnectionStatus{
		State:       r.connectionState,
		Receiver:    r.receiver,
		BotID:       r.botID,
		BotUsername: r.botUsername,
		LastError:   r.lastError,
	}
}

func (r *RuntimeManager) setConnectionStatus(status ConnectionStatus) {
	if r == nil {
		return
	}
	r.mu.Lock()
	r.connectionState = status.State
	r.receiver = status.Receiver
	r.botID = status.BotID
	r.botUsername = status.BotUsername
	r.lastError = status.LastError
	r.mu.Unlock()
}

// SetSettingsChangedHandler installs the process-local callback used by the
// connection manager to reconcile a running receiver after an Admin PATCH.
// It is an internal composition seam and does not alter the kernel contract.
func (r *RuntimeManager) SetSettingsChangedHandler(handler func(context.Context) error) {
	if r == nil {
		return
	}
	r.settingsChangedMu.Lock()
	r.settingsChanged = handler
	r.settingsChangedMu.Unlock()
}

func (r *RuntimeManager) settingsChangedHandler() func(context.Context) error {
	if r == nil {
		return nil
	}
	r.settingsChangedMu.RLock()
	defer r.settingsChangedMu.RUnlock()
	return r.settingsChanged
}

func (r *RuntimeManager) currentConnectionSettings() (string, string) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.mode, r.webhookPublicBaseURL
}

// Update hot-switches the active Bot Token and Webhook Secret.
// Persists encrypted secrets to the database before modifying in-memory state (fail-closed).
func (r *RuntimeManager) Update(ctx context.Context, token, secret string) error {
	mode, webhookPublicBaseURL := r.currentConnectionSettings()
	return r.UpdateSettings(ctx, token, secret, mode, webhookPublicBaseURL)
}

// UpdateSettings persists the complete Telegram settings row before changing
// memory. A serialized update path keeps concurrent PATCH requests from
// publishing a mixed token/secret/mode/URL snapshot.
func (r *RuntimeManager) UpdateSettings(ctx context.Context, token, secret, mode, webhookPublicBaseURL string) error {
	r.updateMu.Lock()
	defer r.updateMu.Unlock()
	return r.updateSettingsLocked(ctx, token, secret, mode, webhookPublicBaseURL)
}

// UpdateSettingsPatch merges a partial Admin settings update while holding the
// same serialization lock as the persistence transaction. This prevents two
// concurrent PATCH requests from reading a stale complementary field and
// overwriting each other (A-006 F-005).
func (r *RuntimeManager) UpdateSettingsPatch(ctx context.Context, token, secret, mode, webhookPublicBaseURL *string) error {
	r.updateMu.Lock()
	defer r.updateMu.Unlock()

	r.mu.RLock()
	currentToken := r.token
	currentSecret := r.secret
	currentMode := r.mode
	currentWebhookPublicBaseURL := r.webhookPublicBaseURL
	r.mu.RUnlock()
	if token != nil {
		currentToken = *token
	}
	if secret != nil {
		currentSecret = *secret
	}
	if mode != nil {
		currentMode = *mode
	}
	if webhookPublicBaseURL != nil {
		currentWebhookPublicBaseURL = *webhookPublicBaseURL
	}
	return r.updateSettingsLocked(ctx, currentToken, currentSecret, currentMode, currentWebhookPublicBaseURL)
}

func (r *RuntimeManager) updateSettingsLocked(ctx context.Context, token, secret, mode, webhookPublicBaseURL string) error {

	trimmedToken := strings.TrimSpace(token)
	trimmedSecret := strings.TrimSpace(secret)
	trimmedMode := normalizeTelegramMode(mode)
	if !ValidTelegramMode(trimmedMode) {
		return fmt.Errorf("telegram: mode must be polling or webhook (got %q)", mode)
	}
	trimmedWebhookPublicBaseURL := strings.TrimSpace(webhookPublicBaseURL)
	if err := validateWebhookPublicBaseURL(trimmedWebhookPublicBaseURL); err != nil {
		return err
	}

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
			_, err := tx.Exec(ctx, `INSERT INTO telegram_config (id, bot_token_enc, webhook_secret_enc, mode, webhook_public_base_url, updated_at)
				VALUES (1, ?, ?, ?, ?, ?)
				ON CONFLICT(id) DO UPDATE SET
					bot_token_enc = excluded.bot_token_enc,
					webhook_secret_enc = excluded.webhook_secret_enc,
					mode = excluded.mode,
					webhook_public_base_url = excluded.webhook_public_base_url,
					updated_at = excluded.updated_at`,
				tokenEnc, secretEnc, trimmedMode, trimmedWebhookPublicBaseURL, now)
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
	r.mode = trimmedMode
	r.webhookPublicBaseURL = trimmedWebhookPublicBaseURL
	r.mu.Unlock()

	if handler := r.settingsChangedHandler(); handler != nil {
		return handler(ctx)
	}
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
		Mode:                  r.mode,
		WebhookPublicBaseURL:  r.webhookPublicBaseURL,
		ConnectionState:       r.connectionState,
		Receiver:              r.receiver,
		BotID:                 r.botID,
		BotUsername:           r.botUsername,
		LastError:             r.lastError,
		CapturedMessagesCount: captured,
		CapturedCount:         captured,
	}
}

// Mock returns the underlying capture sender.
func (r *RuntimeManager) Mock() *CaptureSender {
	return r.mock
}
