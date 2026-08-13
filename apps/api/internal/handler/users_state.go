// F-03 admin enable/disable/unlock endpoints (GOAL-005 D-002 `3): mounted
// under /api/users/{id}/… with the users.enable / users.disable permission
// keys. Disabling is fail-closed: bumps token_version and revokes every live
// refresh token; guards forbid disabling self or the last admin.
package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authsession "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
)

// UserStateRepository is the persistence surface for enable/disable/unlock.
type UserStateRepository interface {
	GetUser(string) (*authsession.User, error)
	SetUserEnabled(string, bool, string, time.Time) (*authsession.User, error)
	UnlockUser(string, time.Time) (*authsession.User, error)
}

// UserStateRoutes returns the admin enable/disable/unlock route contributions
// (admin.account module; keys users.enable / users.disable).
func UserStateRoutes(a *auth.Authenticator, repository UserStateRepository, operations operationlog.Recorder, moduleID string, notifier ...NotifyRepository) []kernel.RouteContribution {
	h := &userStateHandler{repository: repository, operations: operations, now: time.Now}
	if len(notifier) > 0 {
		h.notifier = notifier[0]
	}
	var routes []kernel.RouteContribution
	add := func(method, pattern string, handler http.Handler) {
		routes = append(routes, kernel.RouteContribution{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey(method, pattern)},
			Method:               method,
			Pattern:              pattern,
			Handler:              handler,
		})
	}
	add("POST", "/api/users/{id}/enable", a.Middleware(h.toggle("users.enable", true)))
	add("POST", "/api/users/{id}/disable", a.Middleware(h.toggle("users.disable", false)))
	add("POST", "/api/users/{id}/unlock", a.Middleware(h.unlock()))
	return routes
}

type userStateHandler struct {
	repository UserStateRepository
	operations operationlog.Recorder
	now        func() time.Time
	notifier   NotifyRepository
}

func (h *userStateHandler) toggle(permission string, enabled bool) http.Handler {
	event := operationlog.EventUserEnable
	if !enabled {
		event = operationlog.EventUserDisable
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requirePermission(w, r, permission)
		if !ok {
			return
		}
		id := r.PathValue("id")
		// F-002: only a real enabled→disabled transition produces the
		// notification (re-disabling an already-disabled account is a no-op).
		transition := false
		if !enabled {
			if before, err := h.repository.GetUser(id); err == nil && before.Enabled {
				transition = true
			}
		}
		u, err := h.repository.SetUserEnabled(id, enabled, user.ID, h.now().UTC())
		if err != nil {
			if mapped := mapUserStoreError(err); mapped != err {
				var de *DomainError
				if errors.As(mapped, &de) {
					writeLocalizedError(w, r, de.Status, de.Code, de.Message)
					return
				}
				if errors.Is(mapped, errResourceNotFound) {
					writeLocalizedError(w, r, http.StatusNotFound, "USER_NOT_FOUND", "no user with that id")
					return
				}
			}
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not update user state")
			return
		}
		h.record(event, user, id, u.Username)
		// F-04 system event: notify only on a real disable transition (F-002).
		if transition {
			NotifyAccountEvent(h.notifier, id, "account.disabled", h.now().UTC())
		}
		w.WriteHeader(http.StatusNoContent)
	})
}

func (h *userStateHandler) unlock() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requirePermission(w, r, "users.enable")
		if !ok {
			return
		}
		id := r.PathValue("id")
		// F-002: only a real lock-clear produces the notification (unlocking an
		// already-unlocked account is a no-op).
		wasLocked := false
		if before, err := h.repository.GetUser(id); err == nil {
			wasLocked = before.LockedUntil > h.now().UTC().Unix() || before.FailedLoginCount > 0
		}
		u, err := h.repository.UnlockUser(id, h.now().UTC())
		if err != nil {
			if mapped := mapUserStoreError(err); mapped != err {
				var de *DomainError
				if errors.As(mapped, &de) {
					writeLocalizedError(w, r, de.Status, de.Code, de.Message)
					return
				}
				if errors.Is(mapped, errResourceNotFound) {
					writeLocalizedError(w, r, http.StatusNotFound, "USER_NOT_FOUND", "no user with that id")
					return
				}
			}
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not unlock user")
			return
		}
		h.record(operationlog.EventUserUnlock, user, id, u.Username)
		if wasLocked {
			NotifyAccountEvent(h.notifier, id, "account.unlocked", h.now().UTC())
		}
		w.WriteHeader(http.StatusNoContent)
	})
}

func (h *userStateHandler) record(event string, actor account.User, recordID, username string) {
	op := operationlog.Operation{
		ID:        newOperationID(),
		Event:     event,
		ActorID:   actor.ID,
		ActorName: actor.Name,
		RecordID:  &recordID,
		CreatedAt: h.now().UTC(),
	}
	detail := `{"username":` + jsonQuote(username) + `}`
	op.Detail = &detail
	if h.operations == nil {
		return
	}
	if err := h.operations.RecordOperation(op); err != nil {
		slog.Error("operation log write failed", "event", event, "err", err)
	}
}