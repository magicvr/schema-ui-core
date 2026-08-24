// Account email identity HTTP surface (workspace-018 R3 · GOAL-004 D-001
// §3): self-service bind / verify / resend for the authenticated user's own
// email. Identity-only endpoints (no permission key), like the rest of the
// account self-service. Delivery goes through the ONE composed
// kernel.MailSender; domain sentinels map to catalog codes so zh-CN/en-US
// localization applies.
package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authsession "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
)

// EmailIdentityRepository is the persistence surface consumed by the email
// identity endpoints.
type EmailIdentityRepository interface {
	BindEmail(userID, email string, sender kernel.MailSender, now time.Time) error
	VerifyEmail(userID, code string, now time.Time) error
	ResendEmailCode(userID string, sender kernel.MailSender, now time.Time) error
}

// EmailIdentityRoutes returns the three self-service route contributions.
func EmailIdentityRoutes(a *auth.Authenticator, repo EmailIdentityRepository, sender kernel.MailSender, moduleID string) []kernel.RouteContribution {
	h := &emailIdentityHandler{auth: a, repo: repo, sender: sender, now: time.Now}
	var routes []kernel.RouteContribution
	add := func(method, pattern string, handler http.Handler) {
		routes = append(routes, kernel.RouteContribution{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey(method, pattern)},
			Method:               method,
			Pattern:              pattern,
			Handler:              handler,
		})
	}
	add("POST", "/api/account/email/bind", a.Middleware(h.bind()))
	add("POST", "/api/account/email/verify", a.Middleware(h.verify()))
	add("POST", "/api/account/email/resend", a.Middleware(h.resend()))
	return routes
}

type emailIdentityHandler struct {
	auth   *auth.Authenticator
	repo   EmailIdentityRepository
	sender kernel.MailSender
	now    func() time.Time
}

func (h *emailIdentityHandler) identity(w http.ResponseWriter, r *http.Request) (account.User, bool) {
	user, ok := auth.UserIdentityFrom(r.Context())
	if !ok {
		writeLocalizedError(w, r, http.StatusUnauthorized, "UNAUTHENTICATED", "no active session")
		return account.User{}, false
	}
	return user, true
}

func (h *emailIdentityHandler) bind() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := h.identity(w, r)
		if !ok {
			return
		}
		var body struct {
			Email string `json:"email"`
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxResourceBodyBytes)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_BODY", "body must be JSON with an email field")
			return
		}
		if err := h.repo.BindEmail(user.ID, body.Email, h.sender, h.now()); err != nil {
			h.writeDomainError(w, r, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "pending"})
	})
}

func (h *emailIdentityHandler) verify() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := h.identity(w, r)
		if !ok {
			return
		}
		var body struct {
			Code string `json:"code"`
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxResourceBodyBytes)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_BODY", "body must be JSON with a code field")
			return
		}
		if err := h.repo.VerifyEmail(user.ID, body.Code, h.now()); err != nil {
			h.writeDomainError(w, r, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "verified"})
	})
}

func (h *emailIdentityHandler) resend() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := h.identity(w, r)
		if !ok {
			return
		}
		if err := h.repo.ResendEmailCode(user.ID, h.sender, h.now()); err != nil {
			h.writeDomainError(w, r, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "pending"})
	})
}

// writeDomainError maps repository sentinel errors to their catalog codes.
func (h *emailIdentityHandler) writeDomainError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, authsession.ErrEmailInvalid):
		writeLocalizedError(w, r, http.StatusBadRequest, "EMAIL_INVALID", "invalid email address")
	case errors.Is(err, authsession.ErrEmailTaken):
		writeLocalizedError(w, r, http.StatusConflict, "EMAIL_TAKEN", "email already bound or pending on another account")
	case errors.Is(err, authsession.ErrEmailNotPending):
		writeLocalizedError(w, r, http.StatusConflict, "EMAIL_NOT_PENDING", "no pending email verification for this account")
	case errors.Is(err, authsession.ErrEmailCodeInvalid):
		writeLocalizedError(w, r, http.StatusBadRequest, "EMAIL_CODE_INVALID", "verification code is invalid")
	case errors.Is(err, authsession.ErrEmailCodeExpired):
		writeLocalizedError(w, r, http.StatusBadRequest, "EMAIL_CODE_EXPIRED", "verification code expired; request a new one")
	case errors.Is(err, authsession.ErrEmailResendCooldown):
		writeLocalizedError(w, r, http.StatusTooManyRequests, "EMAIL_RESEND_COOLDOWN", "please wait before requesting another code")
	case errors.Is(err, authsession.ErrEmailSendFailed):
		writeLocalizedError(w, r, http.StatusBadGateway, "EMAIL_SEND_FAILED", "the verification email could not be sent")
	default:
		writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not process the email identity operation")
	}
}
