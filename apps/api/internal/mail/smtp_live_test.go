package mail

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// Live round-trip against a real implicit-TLS (465) SMTP submission endpoint.
// Skipped cleanly unless every MAIL_SMTP_TEST_* variable is set, so a plain
// go test ./... stays offline — mirroring the S3_TEST_* / pgtest precedents.
// This is the "与生产合同等价的 live harness" acceptance face of VP-017 R4:
// the loopback TLS harness proves the protocol shape offline; this test proves
// the same code path against a real peer when operators opt in.
func TestSMTPLiveDelivery(t *testing.T) {
	host := os.Getenv("MAIL_SMTP_TEST_HOST")
	port := os.Getenv("MAIL_SMTP_TEST_PORT")
	username := os.Getenv("MAIL_SMTP_TEST_USERNAME")
	password := os.Getenv("MAIL_SMTP_TEST_PASSWORD")
	from := os.Getenv("MAIL_SMTP_TEST_FROM")
	to := os.Getenv("MAIL_SMTP_TEST_TO")
	if host == "" || port == "" || username == "" || password == "" || from == "" || to == "" {
		t.Skip("MAIL_SMTP_TEST_HOST/PORT/USERNAME/PASSWORD/FROM/TO not set; skipping live SMTP delivery test")
	}
	portNum, err := strconv.Atoi(port)
	if err != nil || portNum < 1 || portNum > 65535 {
		t.Fatalf("MAIL_SMTP_TEST_PORT must be numeric, got %q", port)
	}
	s, err := NewSMTP(SMTPOptions{Host: host, Port: portNum, Username: username, Password: password, From: from})
	if err != nil {
		t.Fatalf("NewSMTP: %v", err)
	}
	ctx := context.Background()
	if err := s.Ping(ctx); err != nil {
		t.Fatalf("ping (endpoint must be reachable over implicit TLS): %v", err)
	}
	msg := kernel.MailMessage{
		To:       to,
		Subject:  "VP-017 R4 live delivery evidence",
		TextBody: "One verifiable delivery through the frozen implicit-TLS path.",
	}
	if err := s.Send(ctx, msg); err != nil {
		t.Fatalf("live Send: %v", err)
	}
}
