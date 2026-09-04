package telegram

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/modules/wallet/subject"
)

// HeaderTelegramSecretToken is the official Telegram secret header.
const HeaderTelegramSecretToken = "X-Telegram-Bot-Api-Secret-Token"

// Rate limit parameters (frozen in GOAL-002 D-002 §6).
const (
	RateLimitWindowIP   = time.Minute
	RateLimitMaxIP      = 60
	RateLimitWindowChat = time.Minute
	RateLimitMaxChat    = 30
	RateLimitWindowUser = time.Minute
	RateLimitMaxUser    = 20
)

// WebhookHandler handles inbound Telegram Webhook POST requests.
type WebhookHandler struct {
	tokenGetter  func() string
	secretGetter func() string
	ipLimiter    kernel.RateLimiter
	chatLimiter  kernel.RateLimiter
	userLimiter  kernel.RateLimiter
	subjectStore *subject.Store
	dispatcher   *Dispatcher
	sender       kernel.TelegramSender
	now          func() time.Time
}

// HandlerConfig carries dependencies for WebhookHandler.
type HandlerConfig struct {
	TokenGetter  func() string
	SecretGetter func() string
	RateLimiters kernel.RateLimiterProvider
	SubjectStore *subject.Store
	Dispatcher   *Dispatcher
	Sender       kernel.TelegramSender
	Now          func() time.Time
}

// NewWebhookHandler constructs a new WebhookHandler.
func NewWebhookHandler(cfg HandlerConfig) *WebhookHandler {
	nowFn := cfg.Now
	if nowFn == nil {
		nowFn = func() time.Time { return time.Now().UTC() }
	}

	var ipLim, chatLim, userLim kernel.RateLimiter
	if cfg.RateLimiters != nil {
		ipLim = cfg.RateLimiters.NewRateLimiter(RateLimitWindowIP, RateLimitMaxIP, 0)
		chatLim = cfg.RateLimiters.NewRateLimiter(RateLimitWindowChat, RateLimitMaxChat, 0)
		userLim = cfg.RateLimiters.NewRateLimiter(RateLimitWindowUser, RateLimitMaxUser, 0)
	}

	return &WebhookHandler{
		tokenGetter:  cfg.TokenGetter,
		secretGetter: cfg.SecretGetter,
		ipLimiter:    ipLim,
		chatLimiter:  chatLim,
		userLimiter:  userLim,
		subjectStore: cfg.SubjectStore,
		dispatcher:   cfg.Dispatcher,
		sender:       cfg.Sender,
		now:          nowFn,
	}
}

// ServeHTTP handles the inbound webhook request following GOAL-002 D-002 §6 order.
func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Check Bot Token configured. If empty, fail-safe 503 without reading body or rate limiting.
	var token string
	if h.tokenGetter != nil {
		token = strings.TrimSpace(h.tokenGetter())
	}
	if token == "" {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	now := h.now()

	// 2. IP Rate Limiting (60/min): atomic check+record (even if secret/parse fails).
	clientIP := handler.LoginClientIP(r)
	ipKey := "tg:webhook:" + clientIP
	if h.ipLimiter != nil {
		if !h.ipLimiter.AllowRecord(ipKey, now) {
			w.Header().Set("Retry-After", strconv.Itoa(h.ipLimiter.RetryAfterSeconds(ipKey, now)))
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
	}

	// 3. Secret Token Verification (fail-closed, constant-time compare).
	var expectedSecret string
	if h.secretGetter != nil {
		expectedSecret = strings.TrimSpace(h.secretGetter())
	}
	if expectedSecret == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	gotSecret := r.Header.Get(HeaderTelegramSecretToken)
	gotHash := sha256.Sum256([]byte(gotSecret))
	expectedHash := sha256.Sum256([]byte(expectedSecret))
	if subtle.ConstantTimeCompare(gotHash[:], expectedHash[:]) != 1 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// 4. Read & Parse JSON Body (bounded to 1MB).
	bodyReader := http.MaxBytesReader(w, r.Body, 1<<20)
	defer bodyReader.Close()

	bodyBytes, err := io.ReadAll(bodyReader)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var payload UpdatePayload
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 5. Apply the same chat/user limits, subject mapping, and dispatcher path
	// used by polling. Identity persistence failures remain retryable 500s.
	if err := h.dispatchPayload(r.Context(), payload, now); err != nil {
		var limitErr *rateLimitExceededError
		if errors.As(err, &limitErr) {
			w.Header().Set("Retry-After", strconv.Itoa(limitErr.retryAfter))
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 6. Return 200 OK with empty body.
	w.WriteHeader(http.StatusOK)
}

// HandlePollingUpdate sends a Bot API update through the same rate-limit,
// subject-mapping, and dispatcher path as webhook delivery. Polling has no
// client IP face, so only the chat/user buckets are applicable.
func (h *WebhookHandler) HandlePollingUpdate(ctx context.Context, payload UpdatePayload) error {
	if h == nil {
		return nil
	}
	return h.dispatchPayload(ctx, payload, h.now())
}

type rateLimitExceededError struct {
	retryAfter int
}

func (e *rateLimitExceededError) Error() string {
	return "telegram: inbound rate limit exceeded"
}

func (h *WebhookHandler) dispatchPayload(ctx context.Context, payload UpdatePayload, now time.Time) error {
	// Extract identifiers and message details.
	var chatID, userID, command, text, callbackData string

	if payload.Message != nil {
		if payload.Message.Chat != nil {
			chatID = strconv.FormatInt(payload.Message.Chat.ID, 10)
		}
		if payload.Message.From != nil {
			userID = strconv.FormatInt(payload.Message.From.ID, 10)
		}
		text = payload.Message.Text
		if strings.HasPrefix(text, "/") {
			parts := strings.Fields(text)
			if len(parts) > 0 {
				command = parts[0]
			}
		}
	} else if payload.CallbackQuery != nil {
		if payload.CallbackQuery.From != nil {
			userID = strconv.FormatInt(payload.CallbackQuery.From.ID, 10)
		}
		callbackData = payload.CallbackQuery.Data
		if payload.CallbackQuery.Message != nil && payload.CallbackQuery.Message.Chat != nil {
			chatID = strconv.FormatInt(payload.CallbackQuery.Message.Chat.ID, 10)
		}
	}

	// Chat Rate Limiting (30/min).
	if h.chatLimiter != nil && chatID != "" {
		chatKey := "tg:chat:" + chatID
		if !h.chatLimiter.AllowRecord(chatKey, now) {
			return &rateLimitExceededError{retryAfter: h.chatLimiter.RetryAfterSeconds(chatKey, now)}
		}
	}

	// User Rate Limiting (20/min).
	if h.userLimiter != nil && userID != "" {
		userKey := "tg:user:" + userID
		if !h.userLimiter.AllowRecord(userKey, now) {
			return &rateLimitExceededError{retryAfter: h.userLimiter.RetryAfterSeconds(userKey, now)}
		}
	}

	// Subject Identity Mapping: GetOrCreateSubject("telegram", userID).
	var subjectID string
	if h.subjectStore != nil && userID != "" {
		sub, _, err := h.subjectStore.GetOrCreateSubject(ctx, "telegram", userID, now)
		if err != nil {
			// Fail-closed on persistence failure so Telegram retries delivery.
			return err
		}
		if sub != nil {
			subjectID = sub.ID
		}
	}

	// Dispatch Update. Handler failures are logged but do not make Telegram
	// retry a successfully accepted update; persistence failures above do.
	if h.dispatcher != nil {
		upd := kernel.TelegramUpdate{
			ChatID:       chatID,
			UserID:       userID,
			SubjectID:    subjectID,
			Command:      command,
			Text:         text,
			CallbackData: callbackData,
		}
		if err := h.dispatcher.Dispatch(ctx, upd, h.sender); err != nil {
			slog.Warn("telegram: dispatch handler error", "err", err, "command", command, "chat_id", chatID)
		}
	}
	return nil
}
