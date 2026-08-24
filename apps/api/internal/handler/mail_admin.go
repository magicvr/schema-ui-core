// Outbound-mail admin configuration surface (VP-017 R7 / workspace-017
// GOAL-008; semantics frozen by Root D-007 over the GOAL-006 D-002 contract):
// read the current channel + non-secret parameters, save a new channel
// selection / configuration (hot switch; secrets write-only), and send one
// test message through the CURRENT kernel.MailSender (no bypass path).
// Independent of /api/settings/* per Root I-012. Secrets NEVER appear in any
// response — only *Set booleans.
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/internal/mail"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
)

// MailAdminService is the runtime face the handlers consume (satisfied by
// *mail.Switcher).
type MailAdminService interface {
	PublicView() (*mail.PublicView, error)
	Update(ctx context.Context, req mail.UpdateRequest) (*mail.PublicView, error)
	Send(ctx context.Context, msg kernel.MailMessage) error
}

// RegisterMailAdmin mounts:
//
//	GET  /api/mail/config      settings.read   current channel + params (+secret presence)
//	PUT  /api/mail/config      settings.write  hot-switch channel / save config
//	POST /api/mail/test-send   settings.write  send one test message via the current channel
func RegisterMailAdmin(mux routeRegistrar, a authMiddleware, svc MailAdminService, operations operationlog.Recorder) {
	mux.Handle("GET /api/mail/config", a.Middleware(mailPermRead(mailConfigGet(svc))))
	mux.Handle("PUT /api/mail/config", a.Middleware(mailPermWrite(mailConfigPut(svc, operations))))
	mux.Handle("POST /api/mail/test-send", a.Middleware(mailPermWrite(mailTestSend(svc, operations))))
}

func mailPermRead(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r, "settings.read"); !ok {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func mailPermWrite(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r, "settings.write"); !ok {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func mailConfigGet(svc MailAdminService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		view, err := svc.PublicView()
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not load mail configuration")
			return
		}
		writeJSON(w, http.StatusOK, view)
	})
}

// mailConfigPut decodes a partial update. Secret fields left empty keep the
// stored value; non-empty secrets replace them after candidate validation.
func mailConfigPut(svc MailAdminService, operations operationlog.Recorder) http.Handler {
	type putBody struct {
		Channel       string `json:"channel"`
		MockRetention *int   `json:"mockRetention"`
		Resend        *struct {
			APIKey *string `json:"apiKey"`
			From   *string `json:"from"`
		} `json:"resend"`
		SMTP *struct {
			Host     *string `json:"host"`
			Port     *int    `json:"port"`
			Username *string `json:"username"`
			Password *string `json:"password"`
			From     *string `json:"from"`
		} `json:"smtp"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, _ := auth.IdentityFrom(r.Context())
		var body putBody
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10)).Decode(&body); err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_BODY", "expected a JSON mail config body")
			return
		}
		req := mail.UpdateRequest{Channel: strings.TrimSpace(body.Channel)}
		req.MockRetention = body.MockRetention
		if body.Resend != nil {
			req.ResendFrom = body.Resend.From
			req.ResendAPIKey = body.Resend.APIKey
		}
		if body.SMTP != nil {
			req.SMTPHost = body.SMTP.Host
			req.SMTPPort = body.SMTP.Port
			req.SMTPUsername = body.SMTP.Username
			req.SMTPPassword = body.SMTP.Password
			req.SMTPFrom = body.SMTP.From
		}
		view, err := svc.Update(r.Context(), req)
		if err != nil && errors.Is(err, mail.ErrUnknownChannel) {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_MAIL_CONFIG", firstLine(err))
			return
		}
		if err != nil && strings.Contains(err.Error(), "retention") {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_MAIL_CONFIG", firstLine(err))
			return
		}
		if err != nil {
			// Candidate validation failure: previous channel keeps serving.
			writeLocalizedError(w, r, http.StatusConflict, "MAIL_SWITCH_REJECTED", firstLine(err))
			return
		}
		recordAudit(operations, user, operationlog.EventMailChannelUpdate, "", auditDetail("channel-update", map[string]any{"channel": view.Channel}), time.Now().UTC(), r.Context())
		writeJSON(w, http.StatusOK, view)
	})
}

// mailTestSend sends ONE test message through the current channel — the same
// Switcher every consumer uses, never an adapter bypass.
func mailTestSend(svc MailAdminService, operations operationlog.Recorder) http.Handler {
	type testBody struct {
		To string `json:"to"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, _ := auth.IdentityFrom(r.Context())
		var body testBody
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<10)).Decode(&body); err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_BODY", "expected a JSON test-send body")
			return
		}
		msg := kernel.MailMessage{To: strings.TrimSpace(body.To), Subject: "Schema UI test mail",
			TextBody: "This is a test message sent from the Schema UI outbound-mail console."}
		if err := svc.Send(r.Context(), msg); err != nil {
			writeLocalizedError(w, r, http.StatusBadGateway, "MAIL_SEND_FAILED", firstLine(err))
			return
		}
		channel := ""
		if view, err := svc.PublicView(); err == nil {
			channel = view.Channel
		}
		recordAudit(operations, user, operationlog.EventMailTestSend, "", auditDetail("test-send", map[string]any{"to": msg.To, "channel": channel}), time.Now().UTC(), r.Context())
		writeJSON(w, http.StatusOK, map[string]any{"sent": true, "channel": channel})
	})
}

func firstLine(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	if i := strings.IndexByte(msg, '\n'); i >= 0 {
		msg = msg[:i]
	}
	return msg
}
