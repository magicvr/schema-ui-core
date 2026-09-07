package telegram

import (
	"context"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// DisabledSender is the fail-closed kernel.TelegramSender provided when channel.telegram is disabled.
type DisabledSender struct{}

var _ kernel.TelegramSender = (*DisabledSender)(nil)

// NewDisabledSender constructs a DisabledSender.
func NewDisabledSender() *DisabledSender {
	return &DisabledSender{}
}

// Send returns ErrTelegramDisabled unconditionally.
func (d *DisabledSender) Send(ctx context.Context, msg kernel.TelegramMessage) error {
	return kernel.ErrTelegramDisabled
}

// DisabledDispatcher is the no-op kernel.TelegramDispatcher provided when channel.telegram is disabled.
type DisabledDispatcher struct{}

var _ kernel.TelegramDispatcher = (*DisabledDispatcher)(nil)

// NewDisabledDispatcher constructs a DisabledDispatcher.
func NewDisabledDispatcher() *DisabledDispatcher {
	return &DisabledDispatcher{}
}

// RegisterCommand is a successful no-op when the channel is disabled (D-002 §1).
func (d *DisabledDispatcher) RegisterCommand(name string, h kernel.TelegramHandler) error {
	return nil
}

// UnregisterCommand is a no-op.
func (d *DisabledDispatcher) UnregisterCommand(name string) {}

// RegisterCallback is a successful no-op when the channel is disabled.
func (d *DisabledDispatcher) RegisterCallback(data string, h kernel.TelegramHandler) error {
	return nil
}

// UnregisterCallback is a no-op.
func (d *DisabledDispatcher) UnregisterCallback(data string) {}
