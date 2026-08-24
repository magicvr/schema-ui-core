package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/mail"
)

// R7 C1 (workspace-017 GOAL-008; Root D-007): the mail admin surface —
// read (settings.read) never returns secret values; write (settings.write)
// hot-switches the channel; test-send goes through the SAME Switcher and is
// audited.
func TestMailAdminSurface(t *testing.T) {
	env := newAuthTestEnv(t)
	key := []byte(strings.Repeat("m", 32))
	sw, err := mail.NewSwitcher(env.st, key, mail.SeedConfig{Channel: mail.RuntimeChannelMock}, nil)
	if err != nil {
		t.Fatalf("NewSwitcher: %v", err)
	}
	RegisterMailAdmin(env.mux, env.a, sw, env.operations)
	RegisterMailOutbox(env.mux, env.a, mail.NewOutboxSink(env.st, 0))

	token := adminToken(t, env)
	authed := func(method, path, body string) (int, map[string]any) {
		t.Helper()
		req := bearer(t, token, method, path, body)
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, req)
		var out map[string]any
		if rr.Body.Len() > 0 {
			_ = json.NewDecoder(rr.Body).Decode(&out)
		}
		return rr.Code, out
	}

	t.Run("config read is authenticated", func(t *testing.T) {
		code, _ := sendJSON(t, env.mux, http.MethodGet, "/api/mail/config", "")
		if code != http.StatusUnauthorized {
			t.Fatalf("anonymous status = %d, want 401", code)
		}
	})

	t.Run("get returns public view without secrets", func(t *testing.T) {
		code, out := authed(http.MethodGet, "/api/mail/config", "")
		if code != http.StatusOK {
			t.Fatalf("status = %d, want 200: %v", code, out)
		}
		if out["channel"] != mail.RuntimeChannelMock {
			t.Fatalf("channel = %v, want mock", out["channel"])
		}
		secrets, _ := out["secrets"].(map[string]any)
		if secrets == nil || secrets["resendApiKeySet"] != false || secrets["smtpPasswordSet"] != false {
			t.Fatalf("secret presence flags wrong: %v", out["secrets"])
		}
		raw, _ := json.Marshal(out)
		if strings.Contains(string(raw), `"apiKey":"`) || strings.Contains(string(raw), `"password":"`) {
			t.Fatalf("secret values leaked through public view: %s", raw)
		}
	})

	t.Run("put switches channel and applies to sends", func(t *testing.T) {
		body := `{"channel":"mock","mockRetention":42}`
		code, out := authed(http.MethodPut, "/api/mail/config", body)
		if code != http.StatusOK || out["mockRetention"].(float64) != 42 {
			t.Fatalf("put = %d %v", code, out)
		}
		code, sent := authed(http.MethodPost, "/api/mail/test-send", `{"to":"user@example.com","subject":"custom subject","body":"custom body"}`)
		if code != http.StatusOK || sent["sent"] != true {
			t.Fatalf("test-send = %d %v", code, sent)
		}
		reader := mail.NewOutboxSink(env.st, 0)
		if n, err := reader.Count(context.Background()); err != nil || n != 1 {
			t.Fatalf("outbox count = %d, %v; want the test message in the mock channel", n, err)
		}
		rows, err := reader.List(context.Background(), 10, 0)
		if err != nil || len(rows) != 1 || rows[0].Subject != "custom subject" {
			t.Fatalf("stored record = %+v, %v; want custom subject persisted", rows, err)
		}
		full, err := reader.Get(context.Background(), rows[0].ID)
		if err != nil || full.Body != "custom body" {
			t.Fatalf("stored body = %+v, %v; want custom body persisted", full, err)
		}
	})

	t.Run("rejected switch keeps previous channel and reports conflict", func(t *testing.T) {
		body := `{"channel":"resend","resendFrom":"not-an-address"}`
		code, _ := authed(http.MethodPut, "/api/mail/config", body)
		if code != http.StatusConflict {
			t.Fatalf("invalid candidate status = %d, want 409", code)
		}
		code, view := authed(http.MethodGet, "/api/mail/config", "")
		if code != http.StatusOK || view["channel"] != mail.RuntimeChannelMock {
			t.Fatalf("channel after rejected switch = %d %v, want mock kept", code, view["channel"])
		}
	})

	t.Run("test send with invalid recipient fails 502", func(t *testing.T) {
		code, _ := authed(http.MethodPost, "/api/mail/test-send", `{"to":"nope"}`)
		if code != http.StatusBadGateway {
			t.Fatalf("invalid recipient status = %d, want 502", code)
		}
	})

	t.Run("audit events recorded", func(t *testing.T) {
		code, list := authed(http.MethodGet, "/api/mail/outbox?limit=10", "")
		if code != http.StatusOK {
			t.Fatalf("outbox list = %d", code)
		}
		_ = list
		code, ops := authed(http.MethodGet, "/api/operations?pageSize=50", "")
		if code != http.StatusOK {
			t.Fatalf("operations = %d", code)
		}
		items, _ := ops["items"].([]any)
		events := map[string]bool{}
		for _, it := range items {
			m, _ := it.(map[string]any)
			if e, ok := m["event"].(string); ok {
				events[e] = true
			}
		}
		if !events["mail.channel-update"] || !events["mail.test-send"] {
			t.Fatalf("mail audit events missing, got %v", events)
		}
	})
}
