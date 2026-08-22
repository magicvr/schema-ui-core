package mail

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// CaptureSink is the embedded default kernel.MailSender for unconfigured SMTP
// (workspace-017 GOAL-002 D-001 / GOAL-004 D-001): a capacity-one in-process
// ring plus one structured log line per accepted message. Tests retrieve the
// last captured message via Last(); production callers only ever see the
// kernel port. It never fails for configuration reasons and never blocks on
// network I/O — mvp/dev startup and quick tests stay SMTP-free by contract.
type CaptureSink struct {
	mu     sync.Mutex
	last   kernel.MailMessage
	has    bool
	logger *slog.Logger
}

// NewCaptureSink returns the embedded default sink. A nil logger falls back
// to slog.Default().
func NewCaptureSink(logger *slog.Logger) *CaptureSink {
	if logger == nil {
		logger = slog.Default()
	}
	return &CaptureSink{logger: logger}
}

// Send validates the message against the frozen port contract, then captures
// it (overwriting any previous message — capacity one) and writes one
// structured log line. Errors mirror the port rules only; there is no
// delivery to fail.
func (s *CaptureSink) Send(_ context.Context, msg kernel.MailMessage) error {
	if err := msg.Validate(); err != nil {
		return fmt.Errorf("mail: %v", err)
	}
	s.mu.Lock()
	s.last = msg
	s.has = true
	s.mu.Unlock()
	s.logger.Info("outbound mail captured",
		"to", msg.To,
		"subject", msg.Subject,
		"bytes", len(msg.TextBody),
	)
	return nil
}

// Last returns the most recently captured message and whether one exists.
// This is the test/dev retrieval surface frozen by workspace-017 GOAL-002
// D-001 §2; it deliberately lives on the concrete adapter, not on the
// kernel port.
func (s *CaptureSink) Last() (kernel.MailMessage, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.last, s.has
}

// Reset clears the captured message (test isolation between cases).
func (s *CaptureSink) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.last = kernel.MailMessage{}
	s.has = false
}
