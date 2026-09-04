package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/errorcatalog"
	"github.com/magicvr/schema-ui-core/apps/api/internal/requestid"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	telegramstore "github.com/magicvr/schema-ui-core/apps/api/modules/channel/telegram/store"
)

const (
	telegramOperatorReadPermission  = "telegram.operator.read"
	telegramOperatorWritePermission = "telegram.operator.write"
	telegramOperatorMaxBodyBytes    = 8 << 10
)

var (
	telegramOperatorChatIDPattern = regexp.MustCompile(`^-?[0-9]+$`)
	telegramOperatorRequestID     = regexp.MustCompile(`^[A-Za-z0-9._-]{1,128}$`)
)

type TelegramOperatorRepository interface {
	ListSessions(context.Context, int64, int, int) ([]telegramstore.Session, int, error)
	ListTimeline(context.Context, int64, int64, int, int) ([]telegramstore.TimelineEntry, int, error)
	CreatePending(context.Context, int64, int64, string, string) (telegramstore.OutboundMessage, bool, error)
	CreateRetry(context.Context, int64, int64, string, string) (telegramstore.OutboundMessage, bool, error)
	MarkSent(context.Context, int64, string) error
	MarkFailed(context.Context, int64, string, string) error
}

// TelegramRuntimeProbe returns the current non-secret runtime state and token
// readiness. Keeping this as a function seam avoids coupling the shared HTTP
// handler package to the channel adapter package (whose webhook adapter already
// depends on handler).
type TelegramRuntimeProbe func() (botID int64, state, receiver, token string)

// TelegramBusinessHandlerProbe reports whether business Telegram handlers have
// occupied the dispatcher and therefore make the operator surface unavailable.
type TelegramBusinessHandlerProbe func() bool

// TelegramOperatorHandler serves the authenticated C3 operator surface. The
// composition root wraps it in auth.Middleware before the Provider receives it;
// the handler still performs a permission check so direct mounts fail closed.
type TelegramOperatorHandler struct {
	runtimeProbe  TelegramRuntimeProbe
	businessProbe TelegramBusinessHandlerProbe
	repository    TelegramOperatorRepository
	sender        kernel.TelegramSender
}

// NewTelegramOperatorHandler constructs the C3 session/transcript/send handler.
func NewTelegramOperatorHandler(runtimeProbe TelegramRuntimeProbe, businessProbe TelegramBusinessHandlerProbe, repository TelegramOperatorRepository, sender kernel.TelegramSender) *TelegramOperatorHandler {
	return &TelegramOperatorHandler{
		runtimeProbe:  runtimeProbe,
		businessProbe: businessProbe,
		repository:    repository,
		sender:        sender,
	}
}

type telegramOperatorRoute struct {
	method     string
	permission string
	chatID     string
	requestID  string
	kind       telegramOperatorRouteKind
}

type telegramOperatorRouteKind uint8

const (
	telegramOperatorSessions telegramOperatorRouteKind = iota + 1
	telegramOperatorMessages
	telegramOperatorSend
	telegramOperatorRetry
)

func (h *TelegramOperatorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, ok := classifyTelegramOperatorRoute(r.URL.Path, r.Method)
	if !ok {
		http.NotFound(w, r)
		return
	}
	if _, ok := requirePermission(w, r, route.permission); !ok {
		return
	}
	if r.Method != route.method {
		writeLocalizedError(w, r, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}
	if !h.runtimeAvailable() {
		writeLocalizedError(w, r, http.StatusConflict, "TELEGRAM_OPERATOR_UNAVAILABLE", "telegram operator is unavailable")
		return
	}

	switch route.kind {
	case telegramOperatorSessions:
		h.listSessions(w, r)
	case telegramOperatorMessages:
		h.listMessages(w, r, route.chatID)
	case telegramOperatorSend:
		h.sendMessage(w, r, route.chatID)
	case telegramOperatorRetry:
		h.retryMessage(w, r, route.chatID, route.requestID)
	default:
		http.NotFound(w, r)
	}
}

func classifyTelegramOperatorRoute(path, method string) (telegramOperatorRoute, bool) {
	path = strings.TrimSuffix(path, "/")
	const prefix = "/api/channel/telegram/operator/sessions"
	if path == prefix {
		return telegramOperatorRoute{
			method:     http.MethodGet,
			permission: telegramOperatorReadPermission,
			kind:       telegramOperatorSessions,
		}, true
	}
	parts := strings.Split(path, "/")
	if len(parts) < 8 || strings.Join(parts[:6], "/") != prefix {
		return telegramOperatorRoute{}, false
	}
	if len(parts) == 8 && parts[7] == "messages" {
		if method == http.MethodPost {
			return telegramOperatorRoute{
				method:     http.MethodPost,
				permission: telegramOperatorWritePermission,
				chatID:     parts[6],
				kind:       telegramOperatorSend,
			}, true
		}
		return telegramOperatorRoute{
			method:     http.MethodGet,
			permission: telegramOperatorReadPermission,
			chatID:     parts[6],
			kind:       telegramOperatorMessages,
		}, true
	}
	if len(parts) == 10 && parts[7] == "messages" && parts[9] == "retry" {
		return telegramOperatorRoute{
			method:     http.MethodPost,
			permission: telegramOperatorWritePermission,
			chatID:     parts[6],
			requestID:  parts[8],
			kind:       telegramOperatorRetry,
		}, true
	}
	return telegramOperatorRoute{}, false
}

func (h *TelegramOperatorHandler) runtimeAvailable() bool {
	if h == nil || h.runtimeProbe == nil || h.businessProbe == nil || h.repository == nil {
		return false
	}
	botID, state, receiver, _ := h.runtimeProbe()
	if state != "running" || botID <= 0 {
		return false
	}
	if receiver != "webhook" && receiver != "polling" {
		return false
	}
	return !h.businessProbe()
}

func (h *TelegramOperatorHandler) senderReady() bool {
	if !h.runtimeAvailable() || h.sender == nil || h.runtimeProbe == nil {
		return false
	}
	_, _, _, token := h.runtimeProbe()
	return strings.TrimSpace(token) != ""
}

type telegramSessionListResponse struct {
	Items    []telegramSessionResponse `json:"items"`
	Total    int                       `json:"total"`
	Page     int                       `json:"page"`
	PageSize int                       `json:"pageSize"`
}

type telegramSessionResponse struct {
	ChatID        string `json:"chatId"`
	ChatType      string `json:"chatType"`
	Title         string `json:"title"`
	Username      string `json:"username"`
	LastMessageAt string `json:"lastMessageAt"`
}

func (h *TelegramOperatorHandler) listSessions(w http.ResponseWriter, r *http.Request) {
	page, pageSize, ok := parseTelegramOperatorPagination(w, r)
	if !ok {
		return
	}
	botID, _, _, _ := h.runtimeProbe()
	sessions, total, err := h.repository.ListSessions(r.Context(), botID, page, pageSize)
	if err != nil {
		writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not list telegram sessions")
		return
	}
	items := make([]telegramSessionResponse, 0, len(sessions))
	for _, session := range sessions {
		items = append(items, telegramSessionResponse{
			ChatID:        strconv.FormatInt(session.ChatID, 10),
			ChatType:      session.ChatType,
			Title:         session.Title,
			Username:      session.Username,
			LastMessageAt: session.LastMessageAt.UTC().Format(time.RFC3339),
		})
	}
	writeJSON(w, http.StatusOK, telegramSessionListResponse{Items: items, Total: total, Page: page, PageSize: pageSize})
}

type telegramTimelineResponse struct {
	Items    []telegramTimelineItem `json:"items"`
	Total    int                    `json:"total"`
	Page     int                    `json:"page"`
	PageSize int                    `json:"pageSize"`
}

type telegramTimelineItem struct {
	ChatID         string  `json:"chatId"`
	Direction      string  `json:"direction"`
	Status         string  `json:"status"`
	OccurredAt     string  `json:"occurredAt"`
	UpdateID       string  `json:"updateId,omitempty"`
	MessageID      string  `json:"messageId,omitempty"`
	UserID         string  `json:"userId,omitempty"`
	SenderUsername string  `json:"senderUsername,omitempty"`
	RequestID      string  `json:"requestId,omitempty"`
	RetryOf        *string `json:"retryOf"`
	Text           string  `json:"text"`
}

func (h *TelegramOperatorHandler) listMessages(w http.ResponseWriter, r *http.Request, rawChatID string) {
	chatID, ok := parseTelegramChatID(w, r, rawChatID)
	if !ok {
		return
	}
	page, pageSize, ok := parseTelegramOperatorPagination(w, r)
	if !ok {
		return
	}
	botID, _, _, _ := h.runtimeProbe()
	entries, total, err := h.repository.ListTimeline(r.Context(), botID, chatID, page, pageSize)
	if err != nil {
		writeTelegramStoreError(w, r, err)
		return
	}
	items := make([]telegramTimelineItem, 0, len(entries))
	for _, entry := range entries {
		item := telegramTimelineItem{
			ChatID:     strconv.FormatInt(entry.ChatID, 10),
			Direction:  entry.Direction,
			Status:     entry.Status,
			OccurredAt: entry.OccurredAt.UTC().Format(time.RFC3339),
			Text:       entry.Text,
		}
		if entry.UpdateID != 0 {
			item.UpdateID = strconv.FormatInt(entry.UpdateID, 10)
		}
		if entry.MessageID != nil {
			item.MessageID = strconv.FormatInt(*entry.MessageID, 10)
		}
		if entry.UserID != nil {
			item.UserID = strconv.FormatInt(*entry.UserID, 10)
		}
		item.SenderUsername = entry.SenderUsername
		item.RequestID = entry.RequestID
		if entry.Direction == "outbound" {
			if entry.RetryOf != "" {
				retryOf := entry.RetryOf
				item.RetryOf = &retryOf
			}
		}
		items = append(items, item)
	}
	writeJSON(w, http.StatusOK, telegramTimelineResponse{Items: items, Total: total, Page: page, PageSize: pageSize})
}

type telegramSendBody struct {
	RequestID string `json:"requestId"`
	Text      string `json:"text"`
}

type telegramRetryBody struct {
	RequestID string `json:"requestId"`
}

type telegramOutboundResponse struct {
	ChatID       string  `json:"chatId"`
	RequestID    string  `json:"requestId"`
	RetryOf      *string `json:"retryOf"`
	Text         string  `json:"text"`
	Status       string  `json:"status"`
	ErrorMessage string  `json:"errorMessage,omitempty"`
	CreatedAt    string  `json:"createdAt"`
	UpdatedAt    string  `json:"updatedAt"`
}

func (h *TelegramOperatorHandler) sendMessage(w http.ResponseWriter, r *http.Request, rawChatID string) {
	chatID, ok := parseTelegramChatID(w, r, rawChatID)
	if !ok {
		return
	}
	var body telegramSendBody
	if !decodeTelegramOperatorBody(w, r, &body) || !validateTelegramSendBody(w, r, body) {
		return
	}
	botID, _, _, _ := h.runtimeProbe()
	message, created, err := h.repository.CreatePending(r.Context(), botID, chatID, body.RequestID, body.Text)
	if err != nil {
		if writeTelegramStoreError(w, r, err) {
			return
		}
		writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not persist telegram message")
		return
	}
	if !created {
		writeJSON(w, http.StatusOK, telegramOutboundResponseFrom(message))
		return
	}
	if !h.senderReady() {
		if err := h.repository.MarkFailed(r.Context(), botID, body.RequestID, "telegram operator became unavailable"); err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not finalize telegram message state")
			return
		}
		writeLocalizedError(w, r, http.StatusConflict, "TELEGRAM_OPERATOR_UNAVAILABLE", "telegram operator is unavailable")
		return
	}
	if err := h.sender.Send(r.Context(), kernel.TelegramMessage{ChatID: strconv.FormatInt(chatID, 10), Text: body.Text}); err != nil {
		if markErr := h.repository.MarkFailed(r.Context(), botID, body.RequestID, "telegram send failed"); markErr != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not finalize telegram message state")
			return
		}
		message.Status = "failed"
		message.ErrorMessage = "telegram send failed"
		message.UpdatedAt = time.Now().UTC()
		writeTelegramSendFailure(w, r, message)
		return
	}
	if err := h.repository.MarkSent(r.Context(), botID, body.RequestID); err != nil {
		// The external send already happened. Keep the durable row pending and
		// fail closed; never call sender again for this request id.
		writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "telegram message state could not be finalized")
		return
	}
	message.Status = "sent"
	message.UpdatedAt = time.Now().UTC()
	writeJSON(w, http.StatusCreated, telegramOutboundResponseFrom(message))
}

func (h *TelegramOperatorHandler) retryMessage(w http.ResponseWriter, r *http.Request, rawChatID, sourceRequestID string) {
	chatID, ok := parseTelegramChatID(w, r, rawChatID)
	if !ok {
		return
	}
	if !telegramOperatorRequestID.MatchString(sourceRequestID) {
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_BODY", "request id is invalid")
		return
	}
	var body telegramRetryBody
	if !decodeTelegramOperatorBody(w, r, &body) || !validateTelegramRequestID(w, r, body.RequestID) {
		return
	}
	botID, _, _, _ := h.runtimeProbe()
	message, created, err := h.repository.CreateRetry(r.Context(), botID, chatID, sourceRequestID, body.RequestID)
	if err != nil {
		if writeTelegramStoreError(w, r, err) {
			return
		}
		writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not persist telegram retry")
		return
	}
	if !created {
		writeJSON(w, http.StatusOK, telegramOutboundResponseFrom(message))
		return
	}
	if !h.senderReady() {
		if err := h.repository.MarkFailed(r.Context(), botID, body.RequestID, "telegram operator became unavailable"); err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not finalize telegram message state")
			return
		}
		writeLocalizedError(w, r, http.StatusConflict, "TELEGRAM_OPERATOR_UNAVAILABLE", "telegram operator is unavailable")
		return
	}
	if err := h.sender.Send(r.Context(), kernel.TelegramMessage{ChatID: strconv.FormatInt(chatID, 10), Text: message.Text}); err != nil {
		if markErr := h.repository.MarkFailed(r.Context(), botID, body.RequestID, "telegram send failed"); markErr != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not finalize telegram message state")
			return
		}
		message.Status = "failed"
		message.ErrorMessage = "telegram send failed"
		message.UpdatedAt = time.Now().UTC()
		writeTelegramSendFailure(w, r, message)
		return
	}
	if err := h.repository.MarkSent(r.Context(), botID, body.RequestID); err != nil {
		writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "telegram message state could not be finalized")
		return
	}
	message.Status = "sent"
	message.UpdatedAt = time.Now().UTC()
	writeJSON(w, http.StatusCreated, telegramOutboundResponseFrom(message))
}

func telegramOutboundResponseFrom(message telegramstore.OutboundMessage) telegramOutboundResponse {
	response := telegramOutboundResponse{
		ChatID:       strconv.FormatInt(message.ChatID, 10),
		RequestID:    message.RequestID,
		Text:         message.Text,
		Status:       message.Status,
		ErrorMessage: message.ErrorMessage,
		CreatedAt:    message.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:    message.UpdatedAt.UTC().Format(time.RFC3339),
	}
	if message.RetryOf != "" {
		retryOf := message.RetryOf
		response.RetryOf = &retryOf
	}
	return response
}

func parseTelegramOperatorPagination(w http.ResponseWriter, r *http.Request) (int, int, bool) {
	query := r.URL.Query()
	page, ok := parsePositiveOperatorParam(query.Get("page"), 1)
	if !ok {
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PAGE", "page must be a positive integer")
		return 0, 0, false
	}
	pageSize, ok := parsePositiveOperatorParam(query.Get("pageSize"), DefaultPageSize)
	if !ok || pageSize > maxPageSize {
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PAGE_SIZE", "pageSize must be a positive integer not exceeding 100")
		return 0, 0, false
	}
	return page, pageSize, true
}

func parsePositiveOperatorParam(raw string, fallback int) (int, bool) {
	if raw == "" {
		return fallback, true
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value < 1 {
		return 0, false
	}
	return value, true
}

func parseTelegramChatID(w http.ResponseWriter, r *http.Request, raw string) (int64, bool) {
	if !telegramOperatorChatIDPattern.MatchString(raw) {
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_BODY", "chat id must be a decimal Telegram id")
		return 0, false
	}
	chatID, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || chatID == 0 {
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_BODY", "chat id must be a valid non-zero Telegram id")
		return 0, false
	}
	return chatID, true
}

func decodeTelegramOperatorBody(w http.ResponseWriter, r *http.Request, destination any) bool {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, telegramOperatorMaxBodyBytes))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_BODY", "body must be JSON")
		return false
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_BODY", "body must contain one JSON object")
		return false
	}
	return true
}

func validateTelegramSendBody(w http.ResponseWriter, r *http.Request, body telegramSendBody) bool {
	if !validateTelegramRequestID(w, r, body.RequestID) {
		return false
	}
	if strings.TrimSpace(body.Text) == "" || len([]byte(body.Text)) > 4096 {
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_BODY", "text must be non-empty and no more than 4096 UTF-8 bytes")
		return false
	}
	return true
}

func validateTelegramRequestID(w http.ResponseWriter, r *http.Request, requestIDValue string) bool {
	if !telegramOperatorRequestID.MatchString(requestIDValue) {
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_BODY", "request id must match [A-Za-z0-9._-]{1,128}")
		return false
	}
	return true
}

func writeTelegramStoreError(w http.ResponseWriter, r *http.Request, err error) bool {
	switch {
	case errors.Is(err, telegramstore.ErrChatNotFound):
		writeLocalizedError(w, r, http.StatusNotFound, "TELEGRAM_CHAT_NOT_FOUND", "no Telegram session for that chat")
		return true
	case errors.Is(err, telegramstore.ErrRequestNotFound):
		writeLocalizedError(w, r, http.StatusNotFound, "TELEGRAM_REQUEST_NOT_FOUND", "no Telegram request with that id")
		return true
	case errors.Is(err, telegramstore.ErrRequestInProgress):
		writeLocalizedError(w, r, http.StatusConflict, "TELEGRAM_REQUEST_IN_PROGRESS", "another Telegram attempt is still pending")
		return true
	case errors.Is(err, telegramstore.ErrRequestConflict):
		writeLocalizedError(w, r, http.StatusConflict, "TELEGRAM_REQUEST_CONFLICT", "request id is already bound to a different payload")
		return true
	case errors.Is(err, telegramstore.ErrRetryNotAllowed):
		writeLocalizedError(w, r, http.StatusConflict, "TELEGRAM_RETRY_NOT_ALLOWED", "this Telegram attempt cannot be retried")
		return true
	default:
		return false
	}
}

func writeTelegramSendFailure(w http.ResponseWriter, r *http.Request, message telegramstore.OutboundMessage) {
	if message.RequestID == "" {
		writeLocalizedError(w, r, http.StatusBadGateway, "TELEGRAM_SEND_FAILED", "telegram send failed")
		return
	}
	locale := errorcatalog.Negotiate(r)
	body, contentLanguage, cataloged := errorcatalog.Body("TELEGRAM_SEND_FAILED", "", locale)
	if !cataloged {
		body = map[string]any{"error": "TELEGRAM_SEND_FAILED", "message": "telegram send failed"}
	} else {
		w.Header().Set("Content-Language", contentLanguage)
	}
	body["item"] = telegramOutboundResponseFrom(message)
	if id := requestid.FromContext(r.Context()); id != "" {
		body[requestid.BodyName] = id
	}
	writeJSON(w, http.StatusBadGateway, body)
}
