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
