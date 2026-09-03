package telegram

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
)

// SettingsHandler exposes Admin diagnostic status and hot-switch management endpoints.
type SettingsHandler struct {
	runtime *RuntimeManager
}

// NewSettingsHandler constructs a new SettingsHandler.
func NewSettingsHandler(runtime *RuntimeManager) *SettingsHandler {
	return &SettingsHandler{runtime: runtime}
}

// updateSettingsRequest defines the JSON payload for PATCH /api/channel/telegram/settings.
type updateSettingsRequest struct {
	BotToken      *string `json:"bot_token"`
	WebhookSecret *string `json:"webhook_secret"`
}

func (h *SettingsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.IdentityFrom(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: authentication required", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet:
		if !slices.Contains(user.Permissions, "settings.read") {
			http.Error(w, "Forbidden: permission required: settings.read", http.StatusForbidden)
			return
		}
		h.handleGet(w, r)
	case http.MethodPatch, http.MethodPut:
		if !slices.Contains(user.Permissions, "settings.write") {
			http.Error(w, "Forbidden: permission required: settings.write", http.StatusForbidden)
			return
		}
		h.handleUpdate(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *SettingsHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if h.runtime == nil {
		_ = json.NewEncoder(w).Encode(RuntimeStatus{})
		return
	}
	_ = json.NewEncoder(w).Encode(h.runtime.Status())
}

func (h *SettingsHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	if h.runtime == nil {
		http.Error(w, "telegram runtime not available", http.StatusInternalServerError)
		return
	}

	var req updateSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	currentToken := h.runtime.GetToken()
	currentSecret := h.runtime.GetSecret()

	newToken := currentToken
	if req.BotToken != nil {
		newToken = strings.TrimSpace(*req.BotToken)
	}

	newSecret := currentSecret
	if req.WebhookSecret != nil {
		newSecret = strings.TrimSpace(*req.WebhookSecret)
	}

	h.runtime.Update(newToken, newSecret)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(h.runtime.Status())
}
