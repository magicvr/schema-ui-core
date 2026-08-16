// F-04 notification surface (GOAL-006 D-002 `4/`5): self-service endpoints +
// best-effort system-event hooks. Notification content is server-side text
// (localized template rendering is R3/B-09); the master switch gates production.
package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authsession "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
)

// NotificationRepository is the persistence surface for notifications.
type NotificationRepository interface {
	NotificationsEnabledFor(string) (bool, error)
	SetNotificationsEnabled(string, bool, time.Time) error
	ListNotifications(string, authsession.NotificationFilter) ([]authsession.Notification, int, error)
	MarkNotificationRead(string, string, time.Time) error
	MarkAllNotificationsRead(string, time.Time) (int, error)
	UnreadNotificationCount(string) (int, error)
}

// NotifyRepository is the minimal surface for best-effort event hooks.
type NotifyRepository interface {
	NotificationsEnabledFor(string) (bool, error)
	CreateNotification(authsession.Notification, time.Time) error
}

// NotificationRoutes returns the notification route contributions.
func NotificationRoutes(a *auth.Authenticator, repository NotificationRepository, moduleID string) []kernel.RouteContribution {
	h := &notificationHandler{repository: repository, now: time.Now}
	var routes []kernel.RouteContribution
	add := func(method, pattern string, handler http.Handler) {
		routes = append(routes, kernel.RouteContribution{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey(method, pattern)},
			Method:               method,
			Pattern:              pattern,
			Handler:              handler,
		})
	}
	add("GET", "/api/notifications", a.Middleware(h.list()))
	add("POST", "/api/notifications/{id}/read", a.Middleware(h.read()))
	add("POST", "/api/notifications/read-all", a.Middleware(h.readAll()))
	add("GET", "/api/notifications/unread-count", a.Middleware(h.unreadCount()))
	add("GET", "/api/notifications/settings", a.Middleware(h.settingsGet()))
	add("PATCH", "/api/notifications/settings", a.Middleware(h.settings()))
	return routes
}

type notificationHandler struct {
	repository NotificationRepository
	now        func() time.Time
}

func (h *notificationHandler) identity(w http.ResponseWriter, r *http.Request) (account.User, bool) {
	user, ok := auth.IdentityFrom(r.Context())
	if !ok {
		writeLocalizedError(w, r, http.StatusUnauthorized, "UNAUTHENTICATED", "no active session")
		return account.User{}, false
	}
	return user, true
}

func notificationRow(n authsession.Notification) map[string]any {
	row := map[string]any{
		"id":        n.ID,
		"event":     n.Event,
		"title":     n.Title,
		"body":      n.Body,
		"read":      false,
		"createdAt": n.CreatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
	}
	if n.ReadAt != nil {
		row["read"] = true
		row["readAt"] = n.ReadAt.UTC().Format("2006-01-02T15:04:05.000Z07:00")
	}
	return row
}

func (h *notificationHandler) list() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := h.identity(w, r)
		if !ok {
			return
		}
		page, ok := intParam(r.URL.Query().Get("page"), 1)
		if !ok {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PAGE", "page must be a positive integer")
			return
		}
		pageSize, ok := intParam(r.URL.Query().Get("pageSize"), 20)
		if !ok || pageSize > maxPageSize {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PAGE_SIZE", "pageSize must be a positive integer not exceeding 100")
			return
		}
		// T-02 (GOAL-013 D-003): q (title/body keyword) + read state
		// (read=unread | read=read) joins the legacy unreadOnly param.
		var readState *bool
		switch strings.ToLower(strings.TrimSpace(r.URL.Query().Get("read"))) {
		case "unread":
			v := false
			readState = &v
		case "read":
			v := true
			readState = &v
		}
		items, total, err := h.repository.ListNotifications(user.ID, authsession.NotificationFilter{
			UnreadOnly: strings.EqualFold(r.URL.Query().Get("unreadOnly"), "true"),
			Q:          strings.TrimSpace(r.URL.Query().Get("q")),
			Read:       readState,
			Page:       page,
			PageSize:   pageSize,
		})
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not list notifications")
			return
		}
		rows := make([]map[string]any, 0, len(items))
		for _, n := range items {
			rows = append(rows, notificationRow(n))
		}
		writeJSON(w, http.StatusOK, resourceList{Items: rows, Total: total, Page: page, PageSize: pageSize})
	})
}

func (h *notificationHandler) read() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := h.identity(w, r)
		if !ok {
			return
		}
		if err := h.repository.MarkNotificationRead(r.PathValue("id"), user.ID, h.now().UTC()); err != nil {
			if errors.Is(err, authsession.ErrNotFound) {
				writeLocalizedError(w, r, http.StatusNotFound, "NOTIFICATION_NOT_FOUND", "no notification with that id")
				return
			}
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not mark notification read")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}

func (h *notificationHandler) readAll() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := h.identity(w, r)
		if !ok {
			return
		}
		updated, err := h.repository.MarkAllNotificationsRead(user.ID, h.now().UTC())
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not mark notifications read")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"updated": updated})
	})
}

func (h *notificationHandler) unreadCount() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := h.identity(w, r)
		if !ok {
			return
		}
		count, err := h.repository.UnreadNotificationCount(user.ID)
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not count notifications")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"unread": count})
	})
}

// settingsGet returns the current switch as a form-facing string value
// ("true"/"false") so the schema select control can prefill (F-001).
func (h *notificationHandler) settingsGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := h.identity(w, r)
		if !ok {
			return
		}
		enabled, err := h.repository.NotificationsEnabledFor(user.ID)
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not read notification settings")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"enabled": boolStringValue(enabled)})
	})
}

// settings accepts the switch as a JSON bool or a "true"/"false" string (the
// schema select control submits strings; F-001).
func (h *notificationHandler) settings() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := h.identity(w, r)
		if !ok {
			return
		}
		var body struct {
			Enabled json.RawMessage `json:"enabled"`
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxResourceBodyBytes)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || len(body.Enabled) == 0 {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SETTINGS_BODY", "body must be JSON with enabled")
			return
		}
		enabled, err := parseBoolValue(body.Enabled)
		if err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SETTINGS_BODY", "enabled must be a boolean or \"true\"/\"false\"")
			return
		}
		if err := h.repository.SetNotificationsEnabled(user.ID, enabled, h.now().UTC()); err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not update notification settings")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"enabled": boolStringValue(enabled)})
	})
}

// boolStringValue renders the form-facing string form of the switch.
func boolStringValue(enabled bool) string {
	if enabled {
		return "true"
	}
	return "false"
}

// parseBoolValue accepts a JSON bool or the strings "true"/"false".
func parseBoolValue(raw json.RawMessage) (bool, error) {
	var b bool
	if err := json.Unmarshal(raw, &b); err == nil {
		return b, nil
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		switch strings.ToLower(strings.TrimSpace(s)) {
		case "true":
			return true, nil
		case "false":
			return false, nil
		}
	}
	return false, errors.New("invalid boolean value")
}

// --- best-effort system-event hooks ---

// notifyAccountEvent produces one system notification if the user's master
// switch is on. Best-effort: any failure is logged by the caller's convention
// and never blocks the business response (D-002 `3).
func NotifyAccountEvent(repository NotifyRepository, userID, event string, now time.Time) {
	if repository == nil || userID == "" {
		return
	}
	enabled, err := repository.NotificationsEnabledFor(userID)
	if err != nil || !enabled {
		return
	}
	var title, body string
	switch event {
	case "account.locked":
		title, body = "Account locked", "Your account was temporarily locked after repeated failed sign-in attempts."
	case "account.disabled":
		title, body = "Account disabled", "Your account was disabled by an administrator."
	case "account.unlocked":
		title, body = "Account unlocked", "Your account was unlocked by an administrator."
	case "account.password-changed":
		title, body = "Password changed", "Your account password was changed."
	default:
		return
	}
	id, err := newNotificationID()
	if err != nil {
		slog.Error("notification id generation failed", "event", event, "err", err)
		return
	}
	if err := repository.CreateNotification(authsession.Notification{
		ID: id, UserID: userID, Event: event, Title: title, Body: body,
	}, now); err != nil {
		// F-003: best-effort failures are logged, never block the business path.
		slog.Error("notification produce failed", "event", event, "user_id", userID, "err", err)
	}
}

// newNotificationID returns a random 128-bit hex id for notification rows.
func newNotificationID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return "ntf-" + hex.EncodeToString(b[:]), nil
}