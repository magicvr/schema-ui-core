// workspace-018 R3 (GOAL-004 A-001 F-001): the managed email prefill must
// actually traverse the users resource factory — `email` travels as a raw
// string field so an explicit "" clears the address, and non-string values
// are rejected with EMAIL_INVALID.
package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

type noopMailSender struct{}

func (noopMailSender) Send(context.Context, kernel.MailMessage) error { return nil }

func TestUsersPatchEmailPrefillFlows(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	now := time.Now().UTC()

	req := bearer(t, token, http.MethodPost, "/api/users",
		`{"username":"bob","name":"Bob","password":"secret123","roles":["viewer"]}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("create status = %d: %s", rr.Code, rr.Body.String())
	}
	var created map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&created)
	id, _ := created["id"].(string)
	if id == "" {
		t.Fatalf("created user missing id: %v", created)
	}

	// Prefill → pending.
	req = bearer(t, token, http.MethodPatch, "/api/users/"+id, `{"email":"Bob@Example.com"}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("prefill patch status = %d: %s", rr.Code, rr.Body.String())
	}
	email, status, err := env.authRepository.EmailIdentityState(id)
	if err != nil || email == nil || *email != "Bob@Example.com" || status == nil || *status != "pending" {
		t.Fatalf("state after prefill = (%v, %v) err %v, want (Bob@Example.com, pending)", email, status, err)
	}

	// Clear → unbound (the "" case PatchFields would have rejected).
	req = bearer(t, token, http.MethodPatch, "/api/users/"+id, `{"email":""}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("clear patch status = %d: %s", rr.Code, rr.Body.String())
	}
	email, status, err = env.authRepository.EmailIdentityState(id)
	if err != nil || email != nil || status != nil {
		t.Fatalf("state after clear = (%v, %v) err %v, want (nil, nil)", email, status, err)
	}

	// Non-string value rejected before touching storage.
	req = bearer(t, token, http.MethodPatch, "/api/users/"+id, `{"email":12345}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("non-string email status = %d, want 400", rr.Code)
	}

	// Cross-account conflict: bind the seed admin first, then prefill bob
	// with the same address case-insensitively → 409 EMAIL_TAKEN.
	if err := env.authRepository.BindEmail("user-admin", "taken@example.com", noopMailSender{}, now); err != nil {
		t.Fatalf("bind seed admin: %v", err)
	}
	req = bearer(t, token, http.MethodPatch, "/api/users/"+id, `{"email":"TAKEN@example.com"}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusConflict {
		t.Fatalf("conflict prefill status = %d, want 409: %s", rr.Code, rr.Body.String())
	}
	var errBody map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&errBody)
	if errBody["error"] != "EMAIL_TAKEN" {
		t.Fatalf("conflict error = %v, want EMAIL_TAKEN", errBody["error"])
	}
}

// W26 (GOAL-038 D-001 §1 · C1): the managed email identity rides the users
// read faces — list and detail return email/emailStatus plus the derived
// emailStatusStyle badge preset (verified→success, pending→warning,
// unbound→""). Same-query projection: no N+1.
func TestUsersReadFacesCarryEmailIdentity(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)

	req := bearer(t, token, http.MethodPost, "/api/users",
		`{"username":"carol","name":"Carol","password":"secret123","roles":["viewer"]}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("create status = %d: %s", rr.Code, rr.Body.String())
	}
	var created map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&created)
	id, _ := created["id"].(string)

	assertIdentity := func(row map[string]any, wantEmail any, wantStatus any, wantStyle string) {
		t.Helper()
		if row["email"] != wantEmail || row["emailStatus"] != wantStatus || row["emailStatusStyle"] != wantStyle {
			t.Fatalf("identity = (%v, %v, %v), want (%v, %v, %q)", row["email"], row["emailStatus"], row["emailStatusStyle"], wantEmail, wantStatus, wantStyle)
		}
	}

	// Unbound: null/null + empty style preset.
	assertIdentity(created, nil, nil, "")

	// Prefill → pending + warning style on list AND detail.
	req = bearer(t, token, http.MethodPatch, "/api/users/"+id, `{"email":"carol@example.com"}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("prefill patch status = %d: %s", rr.Code, rr.Body.String())
	}
	var patched map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&patched)
	assertIdentity(patched, "carol@example.com", "pending", "warning")

	list := func() map[string]any {
		t.Helper()
		lreq := bearer(t, token, http.MethodGet, "/api/users?q=carol", "")
		lrr := httptest.NewRecorder()
		env.mux.ServeHTTP(lrr, lreq)
		if lrr.Code != http.StatusOK {
			t.Fatalf("list status = %d: %s", lrr.Code, lrr.Body.String())
		}
		var out struct {
			Items []map[string]any `json:"items"`
		}
		if err := json.NewDecoder(lrr.Body).Decode(&out); err != nil {
			t.Fatalf("decode list: %v", err)
		}
		if len(out.Items) != 1 {
			t.Fatalf("list items = %d, want 1", len(out.Items))
		}
		return out.Items[0]
	}
	assertIdentity(list(), "carol@example.com", "pending", "warning")

	detail := func() map[string]any {
		t.Helper()
		dreq := bearer(t, token, http.MethodGet, "/api/users/"+id, "")
		drr := httptest.NewRecorder()
		env.mux.ServeHTTP(drr, dreq)
		if drr.Code != http.StatusOK {
			t.Fatalf("detail status = %d: %s", drr.Code, drr.Body.String())
		}
		var row map[string]any
		_ = json.NewDecoder(drr.Body).Decode(&row)
		return row
	}
	assertIdentity(detail(), "carol@example.com", "pending", "warning")
}
