package telegram

import (
	"context"
	"sync"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// Dispatcher implements kernel.TelegramDispatcher in-memory with thread-safe routing.
type Dispatcher struct {
	mu        sync.RWMutex
	commands  map[string]kernel.TelegramHandler
	callbacks map[string]kernel.TelegramHandler
}

var _ kernel.TelegramDispatcher = (*Dispatcher)(nil)

// NewDispatcher constructs an empty Dispatcher.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		commands:  make(map[string]kernel.TelegramHandler),
		callbacks: make(map[string]kernel.TelegramHandler),
	}
}

// HasBusinessHandlers reports whether the dispatcher is occupied by at least
// one business command or callback. It is a concrete-package probe rather than
// a kernel interface method so the public channel contract stays unchanged.
func (d *Dispatcher) HasBusinessHandlers() bool {
	if d == nil {
		return false
	}
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.commands) > 0 || len(d.callbacks) > 0
}

// RegisterCommand registers a handler for the normalized command name.
// Strips leading slash and optional @BotName. Returns error on conflict or invalid name.
func (d *Dispatcher) RegisterCommand(name string, h kernel.TelegramHandler) error {
	norm, err := kernel.NormalizeTelegramCommand(name)
	if err != nil {
		return err
	}
	if h == nil {
		return kernel.ErrTelegramHandlerNil
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.commands[norm]; exists {
		return kernel.ErrTelegramCommandConflict
	}
	d.commands[norm] = h
	return nil
}

// UnregisterCommand removes a registered command handler.
func (d *Dispatcher) UnregisterCommand(name string) {
	norm, err := kernel.NormalizeTelegramCommand(name)
	if err != nil {
		return
	}

	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.commands, norm)
}

// RegisterCallback registers a handler for exact callback_data.
// Returns error on conflict, empty data, or data exceeding 64 bytes.
func (d *Dispatcher) RegisterCallback(data string, h kernel.TelegramHandler) error {
	if err := kernel.ValidateTelegramCallback(data); err != nil {
		return err
	}
	if h == nil {
		return kernel.ErrTelegramHandlerNil
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.callbacks[data]; exists {
		return kernel.ErrTelegramCallbackConflict
	}
	d.callbacks[data] = h
	return nil
}

// UnregisterCallback removes a registered callback handler.
func (d *Dispatcher) UnregisterCallback(data string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.callbacks, data)
}

// Dispatch executes the matched command or callback handler for an inbound Update.
// If an incoming command is not registered, it sends DefaultTelegramUnknownCommandText
// to the chat using sender (if non-nil). Unknown callbacks are silently ignored.
func (d *Dispatcher) Dispatch(ctx context.Context, upd kernel.TelegramUpdate, sender kernel.TelegramSender) error {
	if upd.Command != "" {
		norm, err := kernel.NormalizeTelegramCommand(upd.Command)
		if err != nil {
			return err
		}

		d.mu.RLock()
		h, ok := d.commands[norm]
		d.mu.RUnlock()

		if ok {
			return h(ctx, upd)
		}

		// Unknown command fallback: send frozen fallback text
		if sender != nil && upd.ChatID != "" {
			_ = sender.Send(ctx, kernel.TelegramMessage{
				ChatID: upd.ChatID,
				Text:   kernel.DefaultTelegramUnknownCommandText,
			})
		}
		return nil
	}

	if upd.CallbackData != "" {
		d.mu.RLock()
		h, ok := d.callbacks[upd.CallbackData]
		d.mu.RUnlock()

		if ok {
			return h(ctx, upd)
		}
		return nil
	}

	return nil
}
