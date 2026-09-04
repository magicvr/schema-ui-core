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
	BotToken             *string `json:"bot_token"`
	WebhookSecret        *string `json:"webhook_secret"`
	Mode                 *string `json:"mode"`
	WebhookPublicBaseURL *string `json:"webhook_public_base_url"`
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

	if req.Mode != nil {
		mode := strings.ToLower(strings.TrimSpace(*req.Mode))
		if !ValidTelegramMode(mode) {
			http.Error(w, "invalid telegram mode", http.StatusBadRequest)
			return
		}
		req.Mode = &mode
	}

	if req.WebhookPublicBaseURL != nil {
		webhookPublicBaseURL := strings.TrimSpace(*req.WebhookPublicBaseURL)
		if err := validateWebhookPublicBaseURL(webhookPublicBaseURL); err != nil {
			http.Error(w, "invalid webhook public base URL", http.StatusBadRequest)
			return
		}
		req.WebhookPublicBaseURL = &webhookPublicBaseURL
	}

	if err := h.runtime.UpdateSettingsPatch(r.Context(), req.BotToken, req.WebhookSecret, req.Mode, req.WebhookPublicBaseURL); err != nil {
		http.Error(w, "failed to persist telegram settings: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(h.runtime.Status())
}
