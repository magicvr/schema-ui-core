package mail

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// Runtime channel model (VP-017 R7 / workspace-017 GOAL-008; semantics frozen
// by Root D-007 over the GOAL-006 D-002 contract): the admin-selected channel
// and its parameters live in the mail_config row (migration 0052) and apply to
// every subsequent Send — hot switch, single-process. The file/env layer
// (config.Config) only SEEDS the row on first boot; admin intent then wins and
// survives restarts.
//
// Switch failure keeps the old sender: Update validates the candidate adapter
// BEFORE persisting, so a failed switch leaves both the stored state and the
// cached adapter untouched.

const (
	RuntimeChannelMock   = "mock"
	RuntimeChannelResend = "resend"
	RuntimeChannelSMTP   = "smtp"
)

// SeedConfig carries the file/env-layer resolution used the first time the
// runtime row is created (fresh installs keep the documented config-file
// behavior; existing deployments see zero behavior change until an admin
// saves from the settings tab).
type SeedConfig struct {
	Channel       string
	MockRetention int
	ResendFrom    string
	ResendAPIKey  string // plaintext at the seam; encrypted before storage
	SMTPHost      string
	SMTPPort      int
	SMTPUsername  string
	SMTPPassword  string // plaintext at the seam; encrypted before storage
	SMTPFrom      string
}

// RuntimeConfig is one decoded snapshot of the mail_config row (secrets
// already decrypted for adapter construction; never serialized to HTTP).
type RuntimeConfig struct {
	Channel       string
	MockRetention int
	ResendFrom    string
	ResendAPIKey  string
	SMTPHost      string
	SMTPPort      int
	SMTPUsername  string
	SMTPPassword  string
	SMTPFrom      string
	UpdatedAt     int64
}

// PublicView is the read face of the runtime configuration: it carries NO
// secret values — only whether each secret is set (D-007 write-only rule).
type PublicView struct {
	Channel       string          `json:"channel"`
	MockRetention int             `json:"mockRetention"`
	Resend        PublicResend    `json:"resend"`
	SMTP          PublicSMTP      `json:"smtp"`
	Secrets       PublicSecrets   `json:"secrets"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type PublicResend struct {
	From string `json:"from"`
}

type PublicSMTP struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	From     string `json:"from"`
}

type PublicSecrets struct {
	ResendAPIKeySet bool `json:"resendApiKeySet"`
	SMTPPasswordSet bool `json:"smtpPasswordSet"`
}

// UpdateRequest is the write face: empty secret fields mean "keep the stored
// value"; non-empty secrets replace the stored value after validation.
type UpdateRequest struct {
	Channel       string
	MockRetention *int
	ResendFrom    *string
	ResendAPIKey  *string
	SMTPHost      *string
	SMTPPort      *int
	SMTPUsername  *string
	SMTPPassword  *string
	SMTPFrom      *string
}

// ErrUnknownChannel guards the closed first-wave channel set.
var ErrUnknownChannel = errors.New("mail: unknown channel")

// Switcher is THE composed kernel.MailSender once R7 lands: every Send
// resolves the current runtime state (cached by updated_at; the cache is
// refreshed whenever the row changes) and delegates to the matching adapter.
// It is safe for concurrent use.
type Switcher struct {
	store     kernel.Store
	masterKey []byte
	logger    *slog.Logger

	// resendBaseURL overrides the Resend endpoint for tests/harness only;
	// empty selects the production URL. It swaps WHICH host receives the
	// request, never the request shape or auth scheme.
	resendBaseURL string

	mu    sync.Mutex
	cached *cachedAdapter
}

type cachedAdapter struct {
	updatedAt int64
	sender    kernel.MailSender
}

// NewSwitcher seeds the runtime row on first boot (seed-once; DB wins after)
// and returns the switching sender.
func NewSwitcher(store kernel.Store, masterKey []byte, seed SeedConfig, logger *slog.Logger) (*Switcher, error) {
	if logger == nil {
		logger = slog.Default()
	}
	s := &Switcher{store: store, masterKey: masterKey, logger: logger}
	if err := s.ensureSeeded(seed); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Switcher) ensureSeeded(seed SeedConfig) error {
	err := s.store.Run(context.Background(), func(tx kernel.Tx) error {
		var n int
		if err := tx.QueryRow(context.Background(), `SELECT COUNT(*) FROM mail_config WHERE id = 1`).Scan(&n); err != nil {
			return err
		}
		if n > 0 {
			return nil
		}
		apiKeyEnc, err := EncryptSecret(s.masterKey, seed.ResendAPIKey)
		if err != nil {
			return err
		}
		passwordEnc, err := EncryptSecret(s.masterKey, seed.SMTPPassword)
		if err != nil {
			return err
		}
		retention := seed.MockRetention
		if retention <= 0 {
			retention = DefaultOutboxCap
		}
		_, err = tx.Exec(context.Background(),
			`INSERT INTO mail_config (id, channel, mock_retention, resend_from, resend_api_key_enc,
			  smtp_host, smtp_port, smtp_username, smtp_password_enc, smtp_from, updated_at)
			 VALUES (1, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			seed.Channel, retention, seed.ResendFrom, apiKeyEnc,
			seed.SMTPHost, seed.SMTPPort, seed.SMTPUsername, passwordEnc, seed.SMTPFrom,
			time.Now().UnixMilli(),
		)
		return err
	})
	if err != nil {
		return fmt.Errorf("mail: seed runtime config: %w", err)
	}
	return nil
}

// Send implements kernel.MailSender over the CURRENT runtime channel.
// W26 (GOAL-038 D-001 §2.1): every channel's outbound mail is recorded in
// mail_outbox. The mock transport records itself (delivered on accept); for
// resend/smtp the switcher appends one record AFTER the adapter returns —
// sent / failed per the frozen vocabulary, inside the shared bounded
// retention. A record-write failure never changes the send result.
func (s *Switcher) Send(ctx context.Context, msg kernel.MailMessage) error {
	cfg, sender, err := s.currentSender()
	if err != nil {
		return fmt.Errorf("mail: resolve active channel: %w", err)
	}
	serr := sender.Send(ctx, msg)
	if cfg.Channel == RuntimeChannelMock {
		return serr // OutboxSink already published the delivered record
	}
	status := DeliverySent
	if serr != nil {
		status = DeliveryFailed
	}
	if rerr := publishOutboundRecord(s.store, cfg.MockRetention, time.Now().UTC(), msg, cfg.Channel, status); rerr != nil {
		s.logger.Error("outbound record write failed", "channel", cfg.Channel, "err", rerr)
	}
	return serr
}

// currentSender returns the runtime row plus the adapter for the stored
// state, rebuilding it when the row changed since the cache was filled.
func (s *Switcher) currentSender() (*RuntimeConfig, kernel.MailSender, error) {
	cfg, err := s.LoadRuntime()
	if err != nil {
		return nil, nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cached != nil && s.cached.updatedAt == cfg.UpdatedAt {
		return cfg, s.cached.sender, nil
	}
	sender, err := s.buildAdapter(*cfg)
	if err != nil {
		return nil, nil, err
	}
	s.cached = &cachedAdapter{updatedAt: cfg.UpdatedAt, sender: sender}
	s.logger.Info("outbound mail channel active", "channel", cfg.Channel)
	return cfg, sender, nil
}

// buildAdapter constructs the adapter for a candidate config. It is shared by
// Send-time caching and Update-time validation so a switch can never persist
// a configuration that would not build.
func (s *Switcher) buildAdapter(cfg RuntimeConfig) (kernel.MailSender, error) {
	switch strings.TrimSpace(cfg.Channel) {
	case RuntimeChannelMock:
		return NewOutboxSink(s.store, cfg.MockRetention), nil
	case RuntimeChannelResend:
		sender, err := NewResend(ResendOptions{APIKey: cfg.ResendAPIKey, From: cfg.ResendFrom, BaseURL: s.resendBaseURL})
		if err != nil {
			return nil, fmt.Errorf("resend: %w", err)
		}
		return sender, nil
	case RuntimeChannelSMTP:
		sender, err := NewSMTP(SMTPOptions{
			Host:     cfg.SMTPHost,
			Port:     cfg.SMTPPort,
			Username: cfg.SMTPUsername,
			Password: cfg.SMTPPassword,
			From:     cfg.SMTPFrom,
		})
		if err != nil {
			return nil, fmt.Errorf("smtp: %w", err)
		}
		return sender, nil
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnknownChannel, cfg.Channel)
	}
}

// LoadRuntime reads and decodes the runtime row (secrets decrypted).
func (s *Switcher) LoadRuntime() (*RuntimeConfig, error) {
	var cfg RuntimeConfig
	var apiKeyEnc, passwordEnc string
	err := s.store.Run(context.Background(), func(tx kernel.Tx) error {
		row := tx.QueryRow(context.Background(),
			`SELECT channel, mock_retention, resend_from, resend_api_key_enc,
			        smtp_host, smtp_port, smtp_username, smtp_password_enc, smtp_from, updated_at
			 FROM mail_config WHERE id = 1`)
		return row.Scan(&cfg.Channel, &cfg.MockRetention, &cfg.ResendFrom, &apiKeyEnc,
			&cfg.SMTPHost, &cfg.SMTPPort, &cfg.SMTPUsername, &passwordEnc, &cfg.SMTPFrom, &cfg.UpdatedAt)
	})
	if err != nil {
		if errors.Is(err, kernel.ErrNoRows) {
			return nil, fmt.Errorf("mail: runtime config row missing (migration 0052 must run first)")
		}
		return nil, fmt.Errorf("mail: load runtime config: %w", err)
	}
	if cfg.ResendAPIKey, err = DecryptSecret(s.masterKey, apiKeyEnc); err != nil {
		return nil, err
	}
	if cfg.SMTPPassword, err = DecryptSecret(s.masterKey, passwordEnc); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// PublicView projects the runtime state for the admin read face — no secret
// values ever leave the server (D-007).
func (s *Switcher) PublicView() (*PublicView, error) {
	cfg, err := s.LoadRuntime()
	if err != nil {
		return nil, err
	}
	return s.publicViewOf(*cfg), nil
}

func (s *Switcher) publicViewOf(cfg RuntimeConfig) *PublicView {
	return &PublicView{
		Channel:       cfg.Channel,
		MockRetention: cfg.MockRetention,
		Resend:        PublicResend{From: cfg.ResendFrom},
		SMTP:          PublicSMTP{Host: cfg.SMTPHost, Port: cfg.SMTPPort, Username: cfg.SMTPUsername, From: cfg.SMTPFrom},
		Secrets: PublicSecrets{
			ResendAPIKeySet: cfg.ResendAPIKey != "",
			SMTPPasswordSet: cfg.SMTPPassword != "",
		},
		UpdatedAt: time.UnixMilli(cfg.UpdatedAt).UTC(),
	}
}

// Update merges the request into the stored row, validates the CANDIDATE
// configuration (build +, for SMTP, an ESMTP Ping probe), and only then
// persists with a fresh updated_at. Any failure returns before mutation:
// the previous channel keeps serving (D-007 切失败保留旧 sender).
func (s *Switcher) Update(ctx context.Context, req UpdateRequest) (*PublicView, error) {
	cfg, err := s.LoadRuntime()
	if err != nil {
		return nil, err
	}
	channel := strings.TrimSpace(req.Channel)
	switch channel {
	case RuntimeChannelMock, RuntimeChannelResend, RuntimeChannelSMTP:
		cfg.Channel = channel
	case "":
		// no channel change requested
	default:
		return nil, fmt.Errorf("mail: %w: %q", ErrUnknownChannel, req.Channel)
	}
	if req.MockRetention != nil {
		if *req.MockRetention < 1 || *req.MockRetention > 100000 {
			return nil, fmt.Errorf("mail: mock retention must be between 1 and 100000")
		}
		cfg.MockRetention = *req.MockRetention
	}
	if req.ResendFrom != nil {
		cfg.ResendFrom = strings.TrimSpace(*req.ResendFrom)
	}
	if req.ResendAPIKey != nil && *req.ResendAPIKey != "" {
		cfg.ResendAPIKey = strings.TrimSpace(*req.ResendAPIKey)
	}
	if req.SMTPHost != nil {
		cfg.SMTPHost = strings.TrimSpace(*req.SMTPHost)
	}
	if req.SMTPPort != nil {
		cfg.SMTPPort = *req.SMTPPort
	}
	if req.SMTPUsername != nil {
		cfg.SMTPUsername = strings.TrimSpace(*req.SMTPUsername)
	}
	if req.SMTPPassword != nil && *req.SMTPPassword != "" {
		cfg.SMTPPassword = *req.SMTPPassword
	}
	if req.SMTPFrom != nil {
		cfg.SMTPFrom = strings.TrimSpace(*req.SMTPFrom)
	}

	// Validate the candidate BEFORE touching the stored row. For SMTP the
	// availability check is the ESMTP Ping over the frozen implicit-TLS dial;
	// mock/resend validate by construction (no network side effect belongs in
	// a settings save).
	candidate, err := s.buildAdapter(*cfg)
	if err != nil {
		return nil, err
	}
	if pinger, ok := candidate.(interface{ Ping(context.Context) error }); ok && cfg.Channel == RuntimeChannelSMTP {
		pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := pinger.Ping(pingCtx); err != nil {
			return nil, fmt.Errorf("mail: new SMTP endpoint is unreachable, keeping the previous channel: %w", err)
		}
	}

	apiKeyEnc, err := EncryptSecret(s.masterKey, cfg.ResendAPIKey)
	if err != nil {
		return nil, err
	}
	passwordEnc, err := EncryptSecret(s.masterKey, cfg.SMTPPassword)
	if err != nil {
		return nil, err
	}
	now := time.Now().UnixMilli()
	err = s.store.Run(context.Background(), func(tx kernel.Tx) error {
		_, err := tx.Exec(context.Background(),
			`UPDATE mail_config SET channel = ?, mock_retention = ?, resend_from = ?, resend_api_key_enc = ?,
			    smtp_host = ?, smtp_port = ?, smtp_username = ?, smtp_password_enc = ?, smtp_from = ?, updated_at = ?
			 WHERE id = 1`,
			cfg.Channel, cfg.MockRetention, cfg.ResendFrom, apiKeyEnc,
			cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUsername, passwordEnc, cfg.SMTPFrom, now,
		)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("mail: save runtime config: %w", err)
	}
	cfg.UpdatedAt = now
	s.mu.Lock()
	s.cached = &cachedAdapter{updatedAt: now, sender: candidate}
	s.mu.Unlock()
	s.logger.Info("outbound mail channel switched", "channel", cfg.Channel)
	return s.publicViewOf(*cfg), nil
}
