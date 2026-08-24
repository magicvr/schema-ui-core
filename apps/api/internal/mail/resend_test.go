package mail

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// R6 C3 (GOAL-006 D-002 §4): the Resend adapter speaks the frozen HTTP shape
// against a harness-equivalent server; failures map to port errors without
// leaking the key.
func TestResendSendRequestShape(t *testing.T) {
	var gotPath, gotAuth, gotContentType string
	var gotBody resendEmailRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotAuth = r.Header.Get("Authorization")
		gotContentType = r.Header.Get("Content-Type")
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Errorf("decode request body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"resend-1"}`))
	}))
	defer srv.Close()

	adapter, err := NewResend(ResendOptions{APIKey: "re-secret", From: "no-reply@example.com", BaseURL: srv.URL})
	if err != nil {
		t.Fatalf("NewResend: %v", err)
	}
	if err := adapter.Send(context.Background(), msg("user@example.com", "hello", "body")); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if gotPath != "/emails" {
		t.Fatalf("path = %q, want /emails", gotPath)
	}
	if gotAuth != "Bearer re-secret" {
		t.Fatalf("auth header = %q", gotAuth)
	}
	if !strings.HasPrefix(gotContentType, "application/json") {
		t.Fatalf("content type = %q, want application/json", gotContentType)
	}
	if gotBody.From != "no-reply@example.com" || len(gotBody.To) != 1 || gotBody.To[0] != "user@example.com" ||
		gotBody.Subject != "hello" || gotBody.Text != "body" {
		t.Fatalf("request body = %+v", gotBody)
	}
}

func TestResendNon2xxFailsWithStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"name":"validation_error","message":"invalid api key"}`))
	}))
	defer srv.Close()

	adapter, err := NewResend(ResendOptions{APIKey: "re-secret", From: "no-reply@example.com", BaseURL: srv.URL})
	if err != nil {
		t.Fatalf("NewResend: %v", err)
	}
	err = adapter.Send(context.Background(), msg("user@example.com", "s", "b"))
	if err == nil || !strings.Contains(err.Error(), "401") {
		t.Fatalf("non-2xx must fail with status, got %v", err)
	}
}

func TestResendConstructorFailsClosed(t *testing.T) {
	if _, err := NewResend(ResendOptions{From: "f@example.com"}); err == nil {
		t.Fatal("missing api key must fail at construction")
	}
	if _, err := NewResend(ResendOptions{APIKey: "k"}); err == nil {
		t.Fatal("missing from must fail at construction")
	}
	if _, err := NewResend(ResendOptions{APIKey: "k", From: "Ops <ops@example.com>"}); err == nil {
		t.Fatal("display-name from must fail at construction")
	}
}

func TestResendValidatesMessage(t *testing.T) {
	adapter, err := NewResend(ResendOptions{APIKey: "k", From: "f@example.com"})
	if err != nil {
		t.Fatal(err)
	}
	if err := adapter.Send(context.Background(), msg("bad", "s", "b")); err == nil {
		t.Fatal("invalid recipient must fail closed before any HTTP traffic")
	}
}
