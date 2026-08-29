package kernel

// Kernel outbound-mail send port (VP-017 / workspace-017 GOAL-002 D-001, R1).
//
// The port is the only mail contract for the kernel and every module: a
// synchronous Send of one plain-text message to exactly one recipient.
// Public types carry no SMTP client types — handlers and module providers
// consume MailSender / MailMessage only, while dialing details stay inside
// the internal/mail adapters (capture/log sink default, explicit SMTP). This
// mirrors the kernel.ObjectStore storage-port precedent: domain code consumes
// the port, never a backend handle.
//
// Contract frozen by workspace-017 GOAL-002 D-001:
//
//   - Send is synchronous: no queue, no outbox, no retry. Failures return as
//     errors for the caller to handle.
//   - The From header never travels on MailMessage: the default sender comes
//     from configuration and adapters stamp it at delivery time, so callers
//     cannot forge the sender identity.
//   - To is a single RFC 5322 addr-spec (Validate rejects anything else).
//     Subject/TextBody pass through verbatim — product policy belongs to
//     consumers, not to the port.
//   - HTML/MIME bodies are out of scope for this wave (VP-017 I-017-005).

import (
	"context"
	"fmt"
	"net/mail"
	"strings"
)

// MailMessage is one outbound plain-text message addressed to a single
// recipient. There is deliberately no From field: the configured default
// sender is stamped by the active adapter at delivery time.
type MailMessage struct {
	// To is the sole recipient address (RFC 5322 addr-spec). Required; the
	// single-recipient cardinality keeps the forwarding surface minimal
	// (workspace-017 GOAL-002 D-001 §3).
	To string
	// Subject passes through verbatim; adapters must not mutate it.
	Subject string
	// TextBody is the plain-text body. HTML/MIME is not part of this wave.
	TextBody string
}

// Validate enforces the contract-level rules shared by every adapter: To must
// be present and parse as a single RFC 5322 address. Subject and TextBody are
// passed through without policy checks — message-shape decisions belong to
// consumers. Adapters must call this before touching their backend so a
// malformed message fails closed at the port instead of surfacing as a
// backend-specific error.
func (m MailMessage) Validate() error {
	addr := strings.TrimSpace(m.To)
	if addr == "" {
		return fmt.Errorf("kernel: mail message missing recipient (To)")
	}
	parsed, err := mail.ParseAddress(addr)
	if err != nil {
		return fmt.Errorf("kernel: mail recipient %q is not a valid address: %v", addr, err)
	}
	if parsed.Address != addr {
		return fmt.Errorf("kernel: mail recipient %q must be a bare address (got %q form)", addr, parsed.Address)
	}
	return nil
}

// MailSender is the kernel outbound-mail port (R1). Implementations live in
// internal/mail (capture/log sink when SMTP is unconfigured, SMTP adapter for
// explicit configuration) and must validate the message via MailMessage
//.Validate before delivering, failing closed on violations. Send is
// synchronous: it returns after the adapter accepted the message (captured or
// handed to the SMTP peer), with delivery errors returned to the caller.
type MailSender interface {
	Send(ctx context.Context, msg MailMessage) error
}
