package mail

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// httptestNewResendServer returns a minimal endpoint answering every request
// with the given status (probe + send share the base URL seam).
func httptestNewResendServer(t *testing.T, status int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(`{}`))
	}))
}

// Live round-trip against the real Resend HTTP API. Skipped cleanly unless
// MAIL_RESEND_TEST_API_KEY and MAIL_RESEND_TEST_TO are set (optional base URL
// override via MAIL_RESEND_TEST_BASE_URL for staging endpoints), so a plain
// go test ./... stays offline — mirroring the MAIL_SMTP_TEST_* precedent.
// This is the opt-in "live" face of VP-017 exit criterion 3 ("live 或等价
// harness"): the httptest harness proves the request shape offline; this test
// proves at least one verifiable delivery through the real peer when
// operators supply credentials.
func TestResendLiveDelivery(t *testing.T) {
	apiKey := os.Getenv("MAIL_RESEND_TEST_API_KEY")
	to := os.Getenv("MAIL_RESEND_TEST_TO")
	from := os.Getenv("MAIL_RESEND_TEST_FROM")
	baseURL := os.Getenv("MAIL_RESEND_TEST_BASE_URL")
	if apiKey == "" || to == "" || from == "" {
		t.Skip("MAIL_RESEND_TEST_API_KEY/FROM/TO not set; skipping live Resend delivery test")
	}
	r, err := NewResend(ResendOptions{APIKey: apiKey, From: from, BaseURL: baseURL})
	if err != nil {
		t.Fatalf("NewResend: %v", err)
	}
	ctx := context.Background()
	if err := r.Ping(ctx); err != nil {
		t.Fatalf("probe (account endpoint must accept the key): %v", err)
	}
	msg := kernel.MailMessage{
		To:       to,
		Subject:  "VP-017 R8 live delivery evidence",
		TextBody: "One verifiable delivery through the explicit Resend channel.",
	}
	if err := r.Send(ctx, msg); err != nil {
		t.Fatalf("live Send: %v", err)
	}
}

// R8 C1: the availability probe accepts a healthy endpoint and rejects
// non-2xx / transport failures with status-only errors.
func TestResendPing(t *testing.T) {
	ok := httptestNewResendServer(t, http.StatusOK)
	defer ok.Close()
	r, err := NewResend(ResendOptions{APIKey: "k", From: "f@example.com", BaseURL: ok.URL})
	if err != nil {
		t.Fatal(err)
	}
	if err := r.Ping(context.Background()); err != nil {
		t.Fatalf("healthy probe must pass, got %v", err)
	}

	bad := httptestNewResendServer(t, http.StatusUnauthorized)
	defer bad.Close()
	r2, _ := NewResend(ResendOptions{APIKey: "k", From: "f@example.com", BaseURL: bad.URL})
	if err := r2.Ping(context.Background()); err == nil || !strings.Contains(err.Error(), "401") {
		t.Fatalf("non-2xx probe must fail naming status, got %v", err)
	}
}
