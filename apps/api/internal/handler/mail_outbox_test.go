package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/internal/mail"
)

// R6 C3 (workspace-017 GOAL-007; GOAL-006 D-002 §3): the mock outbox read
// surface is authenticated, settings.read-gated, and independent of
// /api/settings/*. The list envelope follows the unified
// {items,total,page,pageSize} contract; detail exposes the full body.
func TestMailOutboxSurface(t *testing.T) {
	env := newAuthTestEnv(t)
	sink := mail.NewOutboxSink(env.st, 0)
	RegisterMailOutbox(env.mux, env.a, sink)

	t.Run("unauthenticated list is rejected", func(t *testing.T) {
		code, _ := sendJSON(t, env.mux, http.MethodGet, "/api/mail/outbox", "")
		if code != http.StatusUnauthorized {
			t.Fatalf("anonymous list status = %d, want 401", code)
		}
	})

	token := adminToken(t, env)
	authedGet := func(path string) (int, map[string]any) {
		t.Helper()
		req := bearer(t, token, http.MethodGet, path, "")
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, req)
		var body map[string]any
		if rr.Body.Len() > 0 {
			if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
				t.Fatalf("decode %q: %v", rr.Body.String(), err)
			}
		}
		return rr.Code, body
	}

	if err := sink.Send(context.Background(), kernel.MailMessage{To: "a@example.com", Subject: "first", TextBody: "b1"}); err != nil {
		t.Fatalf("seed send: %v", err)
	}
	if err := sink.Send(context.Background(), kernel.MailMessage{To: "b@example.com", Subject: "second", TextBody: "b2"}); err != nil {
		t.Fatalf("seed send: %v", err)
	}

	t.Run("list returns newest-first envelope without bodies", func(t *testing.T) {
		code, out := authedGet("/api/mail/outbox")
		if code != http.StatusOK {
			t.Fatalf("list status = %d, want 200: %v", code, out)
		}
		items, _ := out["items"].([]any)
		if len(items) != 2 {
			t.Fatalf("items = %v, want 2 rows", out["items"])
		}
		first, _ := items[0].(map[string]any)
		if first["subject"] != "second" {
			t.Fatalf("newest-first violated: %v", items)
		}
		if _, hasBody := first["body"]; hasBody {
			t.Fatal("list rows must not carry bodies")
		}
		if total, _ := out["total"].(float64); int(total) != 2 {
			t.Fatalf("total = %v, want 2", out["total"])
		}
	})

	t.Run("detail returns the full body", func(t *testing.T) {
		code, list := authedGet("/api/mail/outbox")
		if code != http.StatusOK {
			t.Fatalf("list status = %d", code)
		}
		items, _ := list["items"].([]any)
		first, _ := items[0].(map[string]any)
		id, _ := first["id"].(string)
		code, detail := authedGet("/api/mail/outbox/" + id)
		if code != http.StatusOK {
			t.Fatalf("detail status = %d, want 200: %v", code, detail)
		}
		if detail["body"] != "b2" || detail["to"] != "b@example.com" {
			t.Fatalf("detail = %v, want full record with body", detail)
		}
	})

	t.Run("unknown id is 404", func(t *testing.T) {
		code, _ := authedGet("/api/mail/outbox/outbox-missing")
		if code != http.StatusNotFound {
			t.Fatalf("unknown id status = %d, want 404", code)
		}
	})

	t.Run("limit pagination caps at 200 and pages via offset", func(t *testing.T) {
		code, out := authedGet("/api/mail/outbox?limit=1&offset=1")
		if code != http.StatusOK {
			t.Fatalf("paged list status = %d, want 200: %v", code, out)
		}
		items, _ := out["items"].([]any)
		if len(items) != 1 {
			t.Fatalf("paged items = %v, want exactly 1", items)
		}
		first, _ := items[0].(map[string]any)
		if first["subject"] != "first" {
			t.Fatalf("offset=1 must skip newest, got %v", first)
		}
		if page, _ := out["page"].(float64); int(page) != 2 {
			t.Fatalf("page = %v, want 2", out["page"])
		}
	})

	t.Run("viewer without settings.read is forbidden", func(t *testing.T) {
		env.addUser(t, "viewer9", "viewer-pass-123", []string{"viewer"})
		viewerToken := env.login(t, "viewer9", "viewer-pass-123")
		req := bearer(t, viewerToken, http.MethodGet, "/api/mail/outbox", "")
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusForbidden {
			t.Fatalf("viewer status = %d, want 403", rr.Code)
		}
	})
}
