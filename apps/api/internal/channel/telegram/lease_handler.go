package telegram

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
)

// LeaseHandler exposes the authenticated console-session lease used to keep a
// polling receiver alive while the Telegram admin page is open. The session id
// is always taken from the verified request identity; callers cannot select or
// impersonate another session.
type LeaseHandler struct {
	manager *ConnectionManager
}

// NewLeaseHandler constructs the process-local Telegram console lease handler.
func NewLeaseHandler(manager *ConnectionManager) *LeaseHandler {
	return &LeaseHandler{manager: manager}
}

type leaseResponse struct {
	Active           bool   `json:"active"`
	TTLSeconds       int    `json:"ttl_seconds"`
	ActiveLeaseCount int    `json:"active_lease_count"`
	ConnectionState  string `json:"connection_state"`
	Receiver         string `json:"receiver"`
}

func (h *LeaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.IdentityFrom(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: authentication required", http.StatusUnauthorized)
		return
	}
	if !slices.Contains(user.Permissions, "settings.read") {
		http.Error(w, "Forbidden: permission required: settings.read", http.StatusForbidden)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if h == nil || h.manager == nil {
		http.Error(w, "telegram connection manager not available", http.StatusInternalServerError)
		return
	}

	sessionID := strings.TrimSpace(user.SessionID)
	if sessionID == "" {
		// A user id is deliberately not a fallback: multiple browser sessions
		// for one account must be independently leased and released.
		http.Error(w, "Unauthorized: session identity required", http.StatusUnauthorized)
		return
	}

	var err error
	active := true
	switch r.URL.Path {
	case "/api/channel/telegram/lease/acquire":
		err = h.manager.AcquireLease(r.Context(), sessionID)
	case "/api/channel/telegram/lease/heartbeat":
		err = h.manager.HeartbeatLease(r.Context(), sessionID)
	case "/api/channel/telegram/lease/release":
		err = h.manager.ReleaseLease(r.Context(), sessionID)
		active = false
	default:
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "telegram lease operation failed", http.StatusInternalServerError)
		return
	}

	status := h.manager.Status()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(leaseResponse{
		Active:           active,
		TTLSeconds:       int(PollingLeaseTTL.Seconds()),
		ActiveLeaseCount: h.manager.ActiveLeaseCount(),
		ConnectionState:  status.State,
		Receiver:         status.Receiver,
	})
}
