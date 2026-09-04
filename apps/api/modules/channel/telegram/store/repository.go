// Package store owns the channel.telegram R3 session and inbound receipt
// persistence. It deliberately depends only on the kernel transaction port.
package store

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// TxRunner is the smallest transaction boundary needed by this module.
type TxRunner interface {
	Run(context.Context, func(kernel.Tx) error) error
}

// InboundMessage is the normalized data extracted from a Telegram update.
// Zero-valued optional IDs/strings are stored as SQL NULLs.
type InboundMessage struct {
	BotID           int64
	UpdateID        int64
	ChatID          int64
	ChatType        string
	ChatTitle       string
	ChatUsername    string
	UserID          int64
	MessageID       int64
	CallbackQueryID string
	MessageKind     string
	Text            string
	CallbackData    string
	SenderUsername  string
	ReceivedAt      time.Time
}

// Repository persists Telegram sessions and inbound idempotency receipts.
type Repository struct {
	runner TxRunner
}

// NewRepository constructs a Telegram R3 repository over a kernel store.
func NewRepository(runner TxRunner) *Repository {
	return &Repository{runner: runner}
}

// RecordInbound atomically records an inbound update and its session activity.
// It returns true only for a new receipt. Duplicate (bot_id, update_id) values
// are successful no-ops and do not update the session or trigger dispatch.
func (r *Repository) RecordInbound(ctx context.Context, msg InboundMessage) (bool, error) {
	if r == nil || r.runner == nil {
		return false, fmt.Errorf("telegram store: repository is unavailable")
	}
	if err := msg.validate(); err != nil {
		return false, err
	}

	inserted := false
	err := r.runner.Run(ctx, func(tx kernel.Tx) error {
		result, err := tx.Exec(ctx, `INSERT INTO telegram_inbound_messages (
  bot_id, update_id, chat_id, user_id, message_id, callback_query_id,
  direction, message_kind, text, callback_data, sender_username, received_at
) VALUES (?, ?, ?, ?, ?, ?, 'inbound', ?, ?, ?, ?, ?)
ON CONFLICT (bot_id, update_id) DO NOTHING`,
			msg.BotID,
			msg.UpdateID,
			msg.ChatID,
			nullableInt64(msg.UserID),
			nullableInt64(msg.MessageID),
			nullableString(msg.CallbackQueryID),
			msg.MessageKind,
			nullableString(msg.Text),
			nullableString(msg.CallbackData),
			nullableString(msg.SenderUsername),
			msg.ReceivedAt.Unix(),
		)
		if err != nil {
			return fmt.Errorf("insert telegram inbound receipt: %w", err)
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("read telegram inbound receipt result: %w", err)
		}
		switch affected {
		case 0:
			// Keep the transaction alive after a PostgreSQL conflict and avoid
			// touching the session row on the duplicate path.
			return nil
		case 1:
			inserted = true
		default:
			return fmt.Errorf("insert telegram inbound receipt affected %d rows", affected)
		}

		if _, err := tx.Exec(ctx, `INSERT INTO telegram_sessions (
  bot_id, chat_id, chat_type, title, username, last_message_at, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT (bot_id, chat_id) DO UPDATE SET
  chat_type = CASE WHEN excluded.chat_type <> '' THEN excluded.chat_type ELSE telegram_sessions.chat_type END,
  title = CASE WHEN excluded.title <> '' THEN excluded.title ELSE telegram_sessions.title END,
  username = CASE WHEN excluded.username <> '' THEN excluded.username ELSE telegram_sessions.username END,
  last_message_at = CASE WHEN excluded.last_message_at > telegram_sessions.last_message_at
    THEN excluded.last_message_at ELSE telegram_sessions.last_message_at END,
  updated_at = CASE WHEN excluded.updated_at > telegram_sessions.updated_at
    THEN excluded.updated_at ELSE telegram_sessions.updated_at END`,
			msg.BotID,
			msg.ChatID,
			msg.ChatType,
			msg.ChatTitle,
			msg.ChatUsername,
			msg.ReceivedAt.Unix(),
			msg.ReceivedAt.Unix(),
			msg.ReceivedAt.Unix(),
		); err != nil {
			return fmt.Errorf("upsert telegram session: %w", err)
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return inserted, nil
}

func (m InboundMessage) validate() error {
	if m.BotID <= 0 {
		return fmt.Errorf("telegram store: bot id must be positive")
	}
	if m.UpdateID <= 0 {
		return fmt.Errorf("telegram store: update id must be positive")
	}
	if m.ChatID == 0 {
		return fmt.Errorf("telegram store: chat id is required")
	}
	if m.ReceivedAt.IsZero() {
		return fmt.Errorf("telegram store: received time is required")
	}
	switch m.MessageKind {
	case "text", "command", "callback":
	default:
		return fmt.Errorf("telegram store: unsupported message kind %q", m.MessageKind)
	}
	if (m.MessageKind == "text" || m.MessageKind == "command") && strings.TrimSpace(m.Text) == "" {
		return fmt.Errorf("telegram store: text is required for %s", m.MessageKind)
	}
	return nil
}

func nullableInt64(value int64) any {
	if value == 0 {
		return nil
	}
	return value
}

func nullableString(value string) any {
	if value == "" {
		return nil
	}
	return value
}
