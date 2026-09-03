package telegram

import (
	"context"
	"sync"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// CaptureSender is an in-memory test and fallback implementation of kernel.TelegramSender.
type CaptureSender struct {
	mu       sync.RWMutex
	messages []kernel.TelegramMessage
}

var _ kernel.TelegramSender = (*CaptureSender)(nil)

// NewCaptureSender constructs an empty CaptureSender.
func NewCaptureSender() *CaptureSender {
	return &CaptureSender{
		messages: make([]kernel.TelegramMessage, 0),
	}
}

// Send validates the message and records it in memory.
func (c *CaptureSender) Send(ctx context.Context, msg kernel.TelegramMessage) error {
	if err := msg.Validate(); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.messages = append(c.messages, msg)
	return nil
}

// Messages returns a copy of all captured messages.
func (c *CaptureSender) Messages() []kernel.TelegramMessage {
	c.mu.RLock()
	defer c.mu.RUnlock()

	out := make([]kernel.TelegramMessage, len(c.messages))
	copy(out, c.messages)
	return out
}

// Last returns the most recently captured message, or nil if none.
func (c *CaptureSender) Last() *kernel.TelegramMessage {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.messages) == 0 {
		return nil
	}
	msg := c.messages[len(c.messages)-1]
	return &msg
}

// Reset clears captured messages.
func (c *CaptureSender) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.messages = c.messages[:0]
}
