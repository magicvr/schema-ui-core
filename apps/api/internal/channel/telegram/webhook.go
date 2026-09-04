package telegram

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	telegramstore "github.com/magicvr/schema-ui-core/apps/api/modules/channel/telegram/store"
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
	botIDGetter  func() (int64, error)
	inboundStore *telegramstore.Repository
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
	BotIDGetter  func() (int64, error)
	InboundStore *telegramstore.Repository
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
		botIDGetter:  cfg.BotIDGetter,
		inboundStore: cfg.InboundStore,
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
// subject-mapping, persistence, and dispatcher path as webhook delivery.
// Polling has no client IP face, so only the chat/user buckets are applicable.
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
	normalized, update, supported := normalizeInbound(payload, now)
	if !supported {
		// Non-text and malformed updates are deliberately accepted and skipped;
		// they are outside the C2 transcript boundary.
		return nil
	}
	chatID := strconv.FormatInt(normalized.ChatID, 10)
	userID := ""
	if normalized.UserID != 0 {
		userID = strconv.FormatInt(normalized.UserID, 10)
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

	if h.inboundStore != nil {
		if h.botIDGetter == nil {
			return fmt.Errorf("telegram: bot identity is unavailable")
		}
		botID, err := h.botIDGetter()
		if err != nil {
			return fmt.Errorf("telegram: get bot identity: %w", err)
		}
		if botID <= 0 {
			return fmt.Errorf("telegram: bot identity is unavailable")
		}
		normalized.BotID = botID
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

	if h.inboundStore != nil {
		isNew, err := h.inboundStore.RecordInbound(ctx, normalized)
		if err != nil {
			// The receipt is the retry boundary. Do not confirm or advance
			// polling when the transaction did not commit.
			return err
		}
		if !isNew {
			return nil
		}
	}

	// Dispatch Update. Handler failures are logged but do not make Telegram
	// retry a successfully accepted update; persistence failures above do.
	if h.dispatcher != nil {
		update.SubjectID = subjectID
		if err := h.dispatcher.Dispatch(ctx, update, h.sender); err != nil {
			slog.Warn("telegram: dispatch handler error", "err", err, "command", update.Command, "chat_id", chatID)
		}
	}
	return nil
}

func normalizeInbound(payload UpdatePayload, now time.Time) (telegramstore.InboundMessage, kernel.TelegramUpdate, bool) {
	if payload.Message != nil {
		message := payload.Message
		if message.Chat == nil || strings.TrimSpace(message.Text) == "" {
			return telegramstore.InboundMessage{}, kernel.TelegramUpdate{}, false
		}

		var userID int64
		var senderUsername string
		chatTitle := message.Chat.Title
		if message.From != nil {
			userID = message.From.ID
			senderUsername = message.From.Username
			if strings.EqualFold(message.Chat.Type, "private") && strings.TrimSpace(chatTitle) == "" {
				chatTitle = strings.TrimSpace(strings.Join([]string{message.From.FirstName, message.From.LastName}, " "))
			}
		}
		command := ""
		kind := "text"
		if strings.HasPrefix(message.Text, "/") {
			kind = "command"
			parts := strings.Fields(message.Text)
			if len(parts) > 0 {
				command = parts[0]
			}
		}
		return telegramstore.InboundMessage{
				UpdateID:       payload.UpdateID,
				ChatID:         message.Chat.ID,
				ChatType:       message.Chat.Type,
				ChatTitle:      chatTitle,
				ChatUsername:   message.Chat.Username,
				UserID:         userID,
				MessageID:      message.MessageID,
				MessageKind:    kind,
				Text:           message.Text,
				SenderUsername: senderUsername,
				ReceivedAt:     now,
			}, kernel.TelegramUpdate{
				ChatID:  strconv.FormatInt(message.Chat.ID, 10),
				UserID:  optionalTelegramID(userID),
				Command: command,
				Text:    message.Text,
			}, true
	}

	if payload.CallbackQuery != nil && payload.CallbackQuery.Message != nil && payload.CallbackQuery.Message.Chat != nil {
		callback := payload.CallbackQuery
		chat := callback.Message.Chat
		var userID int64
		var senderUsername string
		if callback.From != nil {
			userID = callback.From.ID
			senderUsername = callback.From.Username
		}
		return telegramstore.InboundMessage{
				UpdateID:        payload.UpdateID,
				ChatID:          chat.ID,
				ChatType:        chat.Type,
				ChatTitle:       chat.Title,
				ChatUsername:    chat.Username,
				UserID:          userID,
				MessageID:       callback.Message.MessageID,
				CallbackQueryID: callback.ID,
				MessageKind:     "callback",
				CallbackData:    callback.Data,
				SenderUsername:  senderUsername,
				ReceivedAt:      now,
			}, kernel.TelegramUpdate{
				ChatID:       strconv.FormatInt(chat.ID, 10),
				UserID:       optionalTelegramID(userID),
				CallbackData: callback.Data,
			}, true
	}

	return telegramstore.InboundMessage{}, kernel.TelegramUpdate{}, false
}

func optionalTelegramID(id int64) string {
	if id == 0 {
		return ""
	}
	return strconv.FormatInt(id, 10)
}
