// Invitation HTTP surface (workspace-019 R3 · GOAL-004 D-001 §3): the
// admin management quartet rides the admin.users module (permission
// users.invite, checked against the actor's resolved permission set), while
// token redemption is a CENTRAL pre-auth route next to login/recovery.
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authsession "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
)

// InviteRepository is the persistence surface consumed by the invitation
// endpoints (admin management + public acceptance).
type InviteRepository interface {
	CreateInvite(invitedBy string, roles []string, email string, ttl time.Duration, now time.Time) (string, *authsession.Invite, error)
	ListInvites(page, pageSize int) ([]authsession.Invite, int, error)
	RevokeInvite(id string, now time.Time) error
	ResendInvite(id string, ttl time.Duration, now time.Time) (string, *authsession.Invite, error)
	AcceptInvite(rawToken, username, name, passwordHash string, now time.Time) (*authsession.User, error)
	ValidateNewPassword(userID, plain string) error
	UserByID(id string) (*authsession.User, error)
	PermissionsForUser(userID string) ([]string, error)
}

func inviteToMap(inv *authsession.Invite) map[string]any {
	out := map[string]any{
		"id":        inv.ID,
		"roles":     inv.Roles,
		"invitedBy": inv.InvitedBy,
		"email":     "",
		"expiresAt": inv.ExpiresAt.UTC().Format(time.RFC3339),
		"createdAt": inv.CreatedAt.UTC().Format(time.RFC3339),
		"status":    "pending",
	}
	if inv.Email != nil {
		out["email"] = *inv.Email
	}
	if inv.ConsumedAt != nil {
		out["status"] = "consumed"
	} else if inv.RevokedAt != nil {
		out["status"] = "revoked"
	} else if !time.Now().UTC().Before(inv.ExpiresAt) {
		out["status"] = "expired"
	}
	return out
}

func writeInviteDomainError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, authsession.ErrInviteNotFound):
		writeLocalizedError(w, r, http.StatusNotFound, "INVITE_INVALID", "no such invitation")
	case errors.Is(err, authsession.ErrInviteInvalid):
		writeLocalizedError(w, r, http.StatusBadRequest, "INVITE_INVALID", "invitation is unknown, expired, already used or revoked")
	case errors.Is(err, authsession.ErrInviteCooldown):
		writeLocalizedError(w, r, http.StatusTooManyRequests, "EMAIL_RESEND_COOLDOWN", "please wait before resending")
	case errors.Is(err, authsession.ErrInviteRoleGone):
		writeLocalizedError(w, r, http.StatusConflict, "INVITE_ROLE_GONE", "invited roles changed; reissue the invitation")
	case errors.Is(err, authsession.ErrUsernameTaken):
		writeLocalizedError(w, r, http.StatusConflict, "USERNAME_TAKEN", "username already exists")
	case errors.Is(err, authsession.ErrPasswordPolicyViolation):
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PASSWORD", "new password violates the active password policy")
	default:
		writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not process the invitation operation")
	}
}

// InviteAdminRoutes returns the four management route contributions. Every
// request must carry the users.invite permission — resolved from the actor's
// own permission set so no new middleware machinery is needed.
func InviteAdminRoutes(a *auth.Authenticator, repo InviteRepository, sender kernel.MailSender, operations operationlog.Recorder, moduleID string) []kernel.RouteContribution {
	h := &inviteAdminHandler{auth: a, repo: repo, sender: sender, operations: operations, now: time.Now}
	var routes []kernel.RouteContribution
	add := func(method, pattern string, handler http.Handler) {
		routes = append(routes, kernel.RouteContribution{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey(method, pattern)},
			Method:               method,
			Pattern:              pattern,
			Handler:              handler,
		})
	}
	add("POST", "/api/users/invites", a.Middleware(h.create()))
	add("GET", "/api/users/invites", a.Middleware(h.list()))
	add("DELETE", "/api/users/invites/{id}", a.Middleware(h.revoke()))
	add("POST", "/api/users/invites/{id}/resend", a.Middleware(h.resend()))
	return routes
}

type inviteAdminHandler struct {
	auth       *auth.Authenticator
	repo       InviteRepository
	sender     kernel.MailSender
	operations operationlog.Recorder
	now        func() time.Time
}

func (h *inviteAdminHandler) create() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r, "users.invite"); !ok {
			return
		}

		actor, _ := auth.UserIdentityFrom(r.Context())
		var body struct {
			Email         string   `json:"email"`
			Roles         []string `json:"roles"`
			ExpiresInDays int      `json:"expiresInDays"`
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxResourceBodyBytes)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_INVITE_BODY", "body must be JSON")
			return
		}
		ttl := defaultTTL(body.ExpiresInDays)
		raw, inv, err := h.repo.CreateInvite(actor.ID, body.Roles, body.Email, ttl, h.now().UTC())
		if err != nil {
			writeInviteDomainError(w, r, err)
			return
		}
		link := inviteLink(r, raw)
		if strings.TrimSpace(body.Email) != "" {
			if serr := sendInviteMail(h.sender, *inv.Email, link, inv.ExpiresAt); serr != nil {
				// Compensate: an unsendable invite must not linger as live
				// stock the admin cannot see a link for.
				_ = h.repo.RevokeInvite(inv.ID, h.now().UTC())
				writeLocalizedError(w, r, http.StatusBadGateway, "EMAIL_SEND_FAILED", "the invitation email could not be sent")
				return
			}
		}
		out := inviteToMap(inv)
		out["token"] = raw // one-time disclosure
		out["link"] = link
		writeJSON(w, http.StatusCreated, out)
	})
}

func defaultTTL(days int) time.Duration {
	if days < 1 || days > 30 {
		days = 7
	}
	return time.Duration(days) * 24 * time.Hour
}

func inviteLink(r *http.Request, rawToken string) string {
	scheme := "http"
	if proto := strings.TrimSpace(r.Header.Get("X-Forwarded-Proto")); proto != "" {
		scheme = proto
	}
	return scheme + "://" + r.Host + "/invite/accept?token=" + rawToken
}

// sendInviteMail dispatches the invitation letter through the ONE composed
// MailSender (GOAL-002 D-001 §3: 只经 kernel.MailSender).
func sendInviteMail(sender kernel.MailSender, to, link string, expires time.Time) error {
	body := "您收到一份账号邀请。请在此链接完成激活（设置用户名与密码）：\n" + link +
		"\n\n有效期至 " + expires.Format(time.RFC3339) + "。\n\n" +
		"You have been invited. Activate your account via this link:\n" + link +
		"\n\nValid until " + expires.Format(time.RFC3339) + ".\n"
	msg := kernel.MailMessage{
		To:       to,
		Subject:  "账号邀请 · Account invitation",
		TextBody: body,
	}
	if err := msg.Validate(); err != nil {
		return err
	}
	return sender.Send(context.Background(), msg)
}

func queryInt(r *http.Request, key string, fallback int) int {
	raw := strings.TrimSpace(r.URL.Query().Get(key))
	if raw == "" {
		return fallback
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 {
		return fallback
	}
	return n
}

func (h *inviteAdminHandler) list() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r, "users.invite"); !ok {
			return
		}
		page := queryInt(r, "page", 1)
		pageSize := queryInt(r, "pageSize", 20)
		if page < 1 {
			page = 1
		}
		if pageSize < 1 || pageSize > 100 {
			pageSize = 20
		}
		invites, total, err := h.repo.ListInvites(page, pageSize)
		if err != nil {
			writeInviteDomainError(w, r, err)
			return
		}
		items := make([]map[string]any, 0, len(invites))
		for i := range invites {
			items = append(items, inviteToMap(&invites[i]))
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": items, "total": total, "page": page, "pageSize": pageSize})
	})
}

func (h *inviteAdminHandler) revoke() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r, "users.invite"); !ok {
			return
		}
		if err := h.repo.RevokeInvite(r.PathValue("id"), h.now().UTC()); err != nil {
			writeInviteDomainError(w, r, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}

func (h *inviteAdminHandler) resend() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r, "users.invite"); !ok {
			return
		}
		raw, inv, err := h.repo.ResendInvite(r.PathValue("id"), 7*24*time.Hour, h.now().UTC())
		if err != nil {
			writeInviteDomainError(w, r, err)
			return
		}
		link := inviteLink(r, raw)
		// A-001 F-003: a failed dispatch surfaces (502) instead of being
		// swallowed — the invite stays live with its rotated token, so the
		// admin retries resend after the cooldown.
		if inv.Email != nil {
			if serr := sendInviteMail(h.sender, *inv.Email, link, inv.ExpiresAt); serr != nil {
				writeLocalizedError(w, r, http.StatusBadGateway, "EMAIL_SEND_FAILED", "the invitation email could not be sent")
				return
			}
		}
		out := inviteToMap(inv)
		out["token"] = raw
		out["link"] = link
		writeJSON(w, http.StatusOK, out)
	})
}

// --- public acceptance (central pre-auth surface) ---

type InviteAcceptRepository interface {
	AcceptInvite(rawToken, username, name, passwordHash string, now time.Time) (*authsession.User, error)
	ValidateNewPassword(userID, plain string) error
}

// RegisterInviteAccept mounts POST /api/auth/invite/accept on the central mux
// (same pre-auth layer as login/recovery; GOAL-004 D-001 §3). Success answers
// 204 WITHOUT tokens — the new user signs in (GOAL-002 D-001 §4 projection).
func RegisterInviteAccept(mux routeRegistrar, repo InviteAcceptRepository) {
	mux.HandleFunc("POST /api/auth/invite/accept", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Token    string `json:"token"`
			Username string `json:"username"`
			Name     string `json:"name"`
			Password string `json:"password"`
		}
		r.Body = http.MaxBytesReader(w, r.Body, 8<<10)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_INVITE_BODY", "body must be JSON with token, username and password")
			return
		}
		if strings.TrimSpace(body.Token) == "" || strings.TrimSpace(body.Username) == "" {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_INVITE_BODY", "token, username and password are required")
			return
		}
		length := len([]byte(body.Password))
		if length < minPasswordBytes || length > maxPasswordBytes || strings.TrimSpace(body.Password) == "" {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PASSWORD", "password must be a non-whitespace string of 8 to 72 bytes")
			return
		}
		if err := repo.ValidateNewPassword("", body.Password); err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PASSWORD", "password violates the active password policy")
			return
		}
		hash, herr := auth.HashPassword(body.Password, passwordHashCost)
		if herr != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not hash password")
			return
		}
		if _, aerr := repo.AcceptInvite(strings.TrimSpace(body.Token), body.Username, body.Name, hash, time.Now().UTC()); aerr != nil {
			writeInviteDomainError(w, r, aerr)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}
