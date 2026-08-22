package mail

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"math/big"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// The fake SMTP peer speaks the frozen dial shape (workspace-017 GOAL-003
// D-001): implicit TLS from the first byte, EHLO/AUTH PLAIN/MAIL/RCPT/DATA/
// QUIT. It captures enough of the session to assert envelope and payload.
type capturedSession struct {
	mailFrom string
	rcptTo   []string
	data     string
	authID   string
}

// testTLS returns a loopback server certificate signed by a throwaway CA,
// plus the CA pool the client must trust. Verification itself is never
// disabled — only WHICH anchor is trusted changes.
func testTLS(t *testing.T) (serverCert tls.Certificate, pool *x509.CertPool) {
	t.Helper()
	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate ca key: %v", err)
	}
	caTmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "loopback test ca"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	caDER, err := x509.CreateCertificate(rand.Reader, caTmpl, caTmpl, &caKey.PublicKey, caKey)
	if err != nil {
		t.Fatalf("create ca: %v", err)
	}
	caCert, err := x509.ParseCertificate(caDER)
	if err != nil {
		t.Fatalf("parse ca: %v", err)
	}

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject:      pkix.Name{CommonName: "127.0.0.1"},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1").To4()},
		DNSNames:     []string{"localhost"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, caCert, &key.PublicKey, caKey)
	if err != nil {
		t.Fatalf("create cert: %v", err)
	}
	leaf, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("parse cert: %v", err)
	}
	pool = x509.NewCertPool()
	pool.AddCert(caCert)
	return tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key, Leaf: leaf}, pool
}

// startFakeSMTP serves exactly one TLS session on a fresh loopback port.
func startFakeSMTP(t *testing.T) (addr string, got *capturedSession, pool *x509.CertPool) {
	t.Helper()
	serverCert, pool := testTLS(t)
	ln, err := tls.Listen("tcp", "127.0.0.1:0", &tls.Config{Certificates: []tls.Certificate{serverCert}})
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	got = &capturedSession{}
	served := make(chan struct{})
	go func() {
		defer close(served)
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		r := bufio.NewReader(conn)
		w := bufio.NewWriter(conn)
		write := func(s string) { w.WriteString(s + "\r\n"); _ = w.Flush() }
		write("220 mail.test ESMTP")
		for {
			line, err := r.ReadString('\n')
			if err != nil {
				return
			}
			line = strings.TrimRight(line, "\r\n")
			cmd := strings.ToUpper(line)
			switch {
			case strings.HasPrefix(cmd, "EHLO"), strings.HasPrefix(cmd, "HELO"):
				w.WriteString("250-mail.test\r\n250-AUTH PLAIN\r\n250 8BITMIME\r\n")
				_ = w.Flush()
			case strings.HasPrefix(cmd, "AUTH"):
				parts := strings.SplitN(line, " ", 3)
				if len(parts) == 3 && strings.EqualFold(parts[1], "PLAIN") {
					if dec, decErr := base64.StdEncoding.DecodeString(parts[2]); decErr == nil {
						fields := strings.Split(string(dec), "\x00")
						if len(fields) == 3 {
							got.authID = fields[1]
						}
					}
				}
				write("235 ok")
			case strings.HasPrefix(cmd, "MAIL FROM:"):
				got.mailFrom = extractAngleAddr(line)
				write("250 ok")
			case strings.HasPrefix(cmd, "RCPT TO:"):
				got.rcptTo = append(got.rcptTo, extractAngleAddr(line))
				write("250 ok")
			case cmd == "DATA":
				write("354 go")
				var sb strings.Builder
				for {
					dl, err := r.ReadString('\n')
					if err != nil {
						return
					}
					if dl == ".\r\n" || dl == ".\n" {
						break
					}
					sb.WriteString(dl)
				}
				got.data = sb.String()
				write("250 ok")
			case cmd == "QUIT":
				write("221 bye")
				return
			default:
				write("250 ok")
			}
		}
	}()
	addr = ln.Addr().String()
	t.Cleanup(func() {
		_ = ln.Close()
		select {
		case <-served:
		case <-time.After(2 * time.Second):
		}
	})
	return addr, got, pool
}

func extractAngleAddr(line string) string {
	i := strings.Index(line, "<")
	j := strings.LastIndex(line, ">")
	if i >= 0 && j > i {
		return line[i+1 : j]
	}
	return ""
}

func TestNewSMTPValidation(t *testing.T) {
	valid := SMTPOptions{Host: "smtp.example.com", Username: "u", Password: "p", From: "f@example.com"}
	tests := []struct {
		name string
		opts SMTPOptions
		want string
	}{
		{"missing host", SMTPOptions{Username: "u", Password: "p", From: "f@example.com"}, "host"},
		{"missing username", SMTPOptions{Host: "h", Password: "p", From: "f@example.com"}, "username"},
		{"missing password", SMTPOptions{Host: "h", Username: "u", From: "f@example.com"}, "password"},
		{"missing from", SMTPOptions{Host: "h", Username: "u", Password: "p"}, "from"},
		{"display-name from", SMTPOptions{Host: "h", Username: "u", Password: "p", From: "Ops <ops@example.com>"}, "bare address"},
		{"port out of range", SMTPOptions{Host: "h", Port: 70000, Username: "u", Password: "p", From: "f@example.com"}, "port"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := NewSMTP(tt.opts); err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("NewSMTP(%+v) error = %v, want containing %q", tt.opts, err, tt.want)
			}
		})
	}
	t.Run("zero port selects frozen 465", func(t *testing.T) {
		s, err := NewSMTP(valid)
		if err != nil {
			t.Fatalf("NewSMTP: %v", err)
		}
		if s.port != DefaultSMTPPort || DefaultSMTPPort != 465 {
			t.Fatalf("default port = %d, want %d", s.port, DefaultSMTPPort)
		}
	})
}

func TestSMTPSendDeliversOverImplicitTLS(t *testing.T) {
	addr, got, pool := startFakeSMTP(t)
	host, portStr, _ := net.SplitHostPort(addr)
	port, _ := strconv.Atoi(portStr)
	s, err := NewSMTP(SMTPOptions{Host: host, Port: port, Username: "mailer", Password: "secret", From: "no-reply@example.com", rootCAs: pool})
	if err != nil {
		t.Fatalf("NewSMTP: %v", err)
	}
	msg := kernel.MailMessage{To: "dest@example.com", Subject: "Verify your account", TextBody: "Your code is 123456.\r\nSecond line."}
	if err := s.Send(context.Background(), msg); err != nil {
		t.Fatalf("Send over implicit TLS: %v", err)
	}
	if got.authID != "mailer" {
		t.Fatalf("AUTH PLAIN identity = %q, want mailer", got.authID)
	}
	if got.mailFrom != "no-reply@example.com" {
		t.Fatalf("MAIL FROM = %q, want configured sender", got.mailFrom)
	}
	if len(got.rcptTo) != 1 || got.rcptTo[0] != "dest@example.com" {
		t.Fatalf("RCPT TO = %v, want [dest@example.com]", got.rcptTo)
	}
	for _, want := range []string{
		"From: no-reply@example.com\r\n",
		"To: dest@example.com\r\n",
		"Subject: Verify your account\r\n",
		"MIME-Version: 1.0\r\n",
		"text/plain; charset=utf-8\r\n",
		"Your code is 123456.\r\nSecond line.",
	} {
		if !strings.Contains(got.data, want) {
			t.Fatalf("DATA missing %q, got:\n%q", want, got.data)
		}
	}
}

// A plaintext listener must never receive a session: the frozen path has
// no STARTTLS/plain fallback, so the TLS handshake itself must fail.
func TestSMTPTLSIsRequired(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		_, _ = conn.Write([]byte("220 plaintext-not-tls\r\n"))
		buf := make([]byte, 512)
		for {
			if _, err := conn.Read(buf); err != nil {
				break
			}
		}
	}()
	host, portStr, _ := net.SplitHostPort(ln.Addr().String())
	port, _ := strconv.Atoi(portStr)
	s, err := NewSMTP(SMTPOptions{Host: host, Port: port, Username: "u", Password: "p", From: "f@example.com"})
	if err != nil {
		t.Fatalf("NewSMTP: %v", err)
	}
	if serr := s.Send(context.Background(), kernel.MailMessage{To: "a@b.com", Subject: "s", TextBody: "b"}); serr == nil {
		t.Fatal("Send against a plaintext endpoint must fail (implicit TLS only)")
	}
}

func TestSMTPSubjectHeaderInjectionRejected(t *testing.T) {
	s, err := NewSMTP(SMTPOptions{Host: "smtp.example.com", Username: "u", Password: "p", From: "f@example.com"})
	if err != nil {
		t.Fatalf("NewSMTP: %v", err)
	}
	msg := kernel.MailMessage{To: "a@b.com", Subject: "Hi\r\nBcc: victim@example.com", TextBody: "b"}
	if err := s.Send(context.Background(), msg); err == nil || !strings.Contains(err.Error(), "control characters") {
		t.Fatalf("header-injection subject must be rejected, got %v", err)
	}
}

func TestSMTPContextCancellation(t *testing.T) {
	s, err := NewSMTP(SMTPOptions{Host: "10.255.255.1", Port: 465, Username: "u", Password: "p", From: "f@example.com"})
	if err != nil {
		t.Fatalf("NewSMTP: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := s.Send(ctx, kernel.MailMessage{To: "a@b.com", Subject: "s", TextBody: "b"}); err == nil {
		t.Fatal("Send with cancelled context must fail")
	}
}

func TestSMTPPingAgainstImplicitTLSPeer(t *testing.T) {
	addr, _, pool := startFakeSMTP(t)
	host, portStr, _ := net.SplitHostPort(addr)
	port, _ := strconv.Atoi(portStr)
	s, err := NewSMTP(SMTPOptions{Host: host, Port: port, Username: "u", Password: "p", From: "f@example.com", rootCAs: pool})
	if err != nil {
		t.Fatalf("NewSMTP: %v", err)
	}
	if err := s.Ping(context.Background()); err != nil {
		t.Fatalf("Ping over implicit TLS: %v", err)
	}
}

func TestSMTPPingFailsAgainstPlaintextPeer(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		_, _ = conn.Write([]byte("220 plaintext-not-tls\r\n"))
		buf := make([]byte, 512)
		for {
			if _, err := conn.Read(buf); err != nil {
				break
			}
		}
	}()
	host, portStr, _ := net.SplitHostPort(ln.Addr().String())
	port, _ := strconv.Atoi(portStr)
	s, err := NewSMTP(SMTPOptions{Host: host, Port: port, Username: "u", Password: "p", From: "f@example.com"})
	if err != nil {
		t.Fatalf("NewSMTP: %v", err)
	}
	if err := s.Ping(context.Background()); err == nil {
		t.Fatal("Ping against a plaintext endpoint must fail (implicit TLS only)")
	}
}
