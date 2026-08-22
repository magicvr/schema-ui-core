// Package mail holds the outbound-mail adapters behind kernel.MailSender
// (VP-017 / workspace-017 GOAL-002 D-001): the embedded capture/log sink
// default and the explicit SMTP delivery adapter. Dialing details never
// escape this package — handlers and module providers consume the kernel
// port only.
package mail

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// DefaultSMTPPort is the frozen dial port (workspace-017 GOAL-003 D-001):
// implicit TLS submission per RFC 8314. Config may override the number, but
// there is exactly ONE dial shape: TLS from the first byte, certificate
// verification always on, STARTTLS never attempted.
const DefaultSMTPPort = 465

// SMTPOptions carries the explicitly configured endpoint. Every field except
// Port is required (enforced fail-closed by NewSMTP); a zero Port selects
// DefaultSMTPPort. Credentials travel only over the TLS session.
//
// rootCAs is an unexported trust-anchor override used by this package's loop
//back harness: it swaps WHICH anchors verify the peer certificate, never
// WHETHER verification happens (nil = system roots; InsecureSkipVerify has no
// code path).
type SMTPOptions struct {
	Host     string
	Port     int // 0 = unset -> DefaultSMTPPort
	Username string
	Password string
	From     string // default sender stamped on every message (bare addr-spec)

	rootCAs *x509.CertPool
}

// SMTP is the production kernel.MailSender adapter: one frozen dial path
// (implicit TLS), PlainAuth over TLS, synchronous delivery. It is safe for
// concurrent use: every Send opens its own connection (no shared mutable
// state), matching the synchronous no-queue contract.
type SMTP struct {
	host     string
	port     int
	username string
	password string
	from     string
	rootCAs  *x509.CertPool // nil = system roots; verification always ON
}

// NewSMTP validates options fail-closed and returns the adapter. Construction
// is the last line before startup wiring: a misconfigured endpoint must die
// here, not at first Send.
func NewSMTP(opts SMTPOptions) (*SMTP, error) {
	for _, pair := range []struct{ name, value string }{
		{"host", opts.Host},
		{"username", opts.Username},
		{"password", opts.Password},
		{"from", opts.From},
	} {
		if strings.TrimSpace(pair.value) == "" {
			return nil, fmt.Errorf("mail: explicit SMTP requires %s", pair.name)
		}
	}
	if opts.Port < 0 || opts.Port > 65535 {
		return nil, fmt.Errorf("mail: SMTP port must be between 1 and 65535")
	}
	from := strings.TrimSpace(opts.From)
	parsed, err := mail.ParseAddress(from)
	if err != nil || parsed.Address != from {
		return nil, fmt.Errorf("mail: SMTP from must be a bare address")
	}
	port := opts.Port
	if port == 0 {
		port = DefaultSMTPPort
	}
	return &SMTP{
		host:     strings.TrimSpace(opts.Host),
		port:     port,
		username: opts.Username,
		password: opts.Password,
		from:     from,
		rootCAs:  opts.rootCAs,
	}, nil
}

// Send delivers one message synchronously over the frozen implicit-TLS path:
//
//	tls.Dial (ServerName=host, MinVersion=TLS1.2, verification ON)
//	  -> smtp.NewClient -> AUTH PLAIN (only when advertised)
//	  -> MAIL FROM <configured from> -> RCPT TO <msg.To> -> DATA -> QUIT
//
// Transport guards enforced here (adapter-level, not product policy):
//   - msg.Validate() must pass (single bare-address recipient);
//   - Subject must be free of CR/LF (header-injection guard — such messages
//     are rejected instead of being silently mutated).
func (s *SMTP) Send(ctx context.Context, msg kernel.MailMessage) error {
	if err := msg.Validate(); err != nil {
		return fmt.Errorf("mail: %v", err)
	}
	subject := msg.Subject
	if strings.ContainsAny(subject, "\r\n") || strings.ContainsFunc(subject, func(r rune) bool { return r < 0x20 && r != '\t' }) {
		return fmt.Errorf("mail: subject contains control characters; refusing to deliver (header-injection guard)")
	}

	dialer := &tls.Dialer{
		Config: &tls.Config{
			ServerName: s.host,
			MinVersion: tls.VersionTLS12,
			RootCAs:    s.rootCAs, // nil = system roots; verification stays ON
		},
	}
	conn, err := dialer.DialContext(ctx, "tcp", net.JoinHostPort(s.host, strconv.Itoa(s.port)))
	if err != nil {
		return fmt.Errorf("mail: dial smtp %s:%d: %w", s.host, s.port, err)
	}
	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		_ = conn.Close()
		return fmt.Errorf("mail: smtp handshake %s: %w", s.host, err)
	}
	defer func() {
		_ = client.Close()
	}()

	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	if ok, _ := client.Extension("AUTH"); ok {
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("mail: smtp auth %s: %w", s.host, err)
		}
	}
	if err := client.Mail(s.from); err != nil {
		return fmt.Errorf("mail: smtp MAIL FROM rejected: %w", err)
	}
	if err := client.Rcpt(msg.To); err != nil {
		return fmt.Errorf("mail: smtp RCPT TO %s rejected: %w", msg.To, err)
	}
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("mail: smtp DATA rejected: %w", err)
	}
	if _, err := w.Write(buildRFC5322(s.from, msg.To, subject, msg.TextBody)); err != nil {
		_ = w.Close()
		return fmt.Errorf("mail: smtp DATA write failed: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("mail: smtp DATA close failed: %w", err)
	}
	return client.Quit()
}

// buildRFC5322 renders the wire form: From/To/Subject headers plus a UTF-8
// plain-text body. Body line endings are canonicalized to CRLF (transport
// encoding — content bytes are unchanged otherwise); dot-stuffing of leading
// dots is handled by the net/smtp DATA writer itself.
func buildRFC5322(from, to, subject, body string) []byte {
	var b strings.Builder
	b.WriteString("From: " + from + "\r\n")
	b.WriteString("To: " + to + "\r\n")
	b.WriteString("Subject: " + subject + "\r\n")
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
	b.WriteString("\r\n")
	body = strings.ReplaceAll(body, "\r\n", "\n")
	body = strings.ReplaceAll(body, "\n", "\r\n")
	b.WriteString(body)
	if !strings.HasSuffix(body, "\r\n") {
		b.WriteString("\r\n")
	}
	return []byte(b.String())
}
