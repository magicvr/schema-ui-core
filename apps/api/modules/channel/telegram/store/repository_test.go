package store

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

func TestRepositoryRecordInboundSQLiteIdempotency(t *testing.T) {
	path := filepath.Join(t.TempDir(), "telegram-ingress.db")
	st, err := testsupport.OpenStore(path, "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	repository := NewRepository(st)
	firstAt := time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC)
	first := InboundMessage{
		BotID:          101,
		UpdateID:       9001,
		ChatID:         8001,
		ChatType:       "private",
		ChatTitle:      "Ada Lovelace",
		ChatUsername:   "ada",
		UserID:         7001,
		MessageID:      6001,
		MessageKind:    "text",
		Text:           "hello",
		SenderUsername: "ada",
		ReceivedAt:     firstAt,
	}
	inserted, err := repository.RecordInbound(context.Background(), first)
	if err != nil {
		t.Fatalf("first RecordInbound: %v", err)
	}
	if !inserted {
		t.Fatal("first RecordInbound = false, want true")
	}

	duplicate := first
	duplicate.ChatTitle = "Changed by duplicate"
	duplicate.ChatUsername = "changed"
	duplicate.ReceivedAt = firstAt.Add(time.Hour)
	inserted, err = repository.RecordInbound(context.Background(), duplicate)
	if err != nil {
		t.Fatalf("duplicate RecordInbound: %v", err)
	}
	if inserted {
		t.Fatal("duplicate RecordInbound = true, want false")
	}
	older := first
	older.UpdateID = 9003
	older.MessageID = 6003
	older.ChatTitle = ""
	older.ChatUsername = ""
	older.ReceivedAt = firstAt.Add(-time.Hour)
	older.Text = "older message"
	inserted, err = repository.RecordInbound(context.Background(), older)
	if err != nil || !inserted {
		t.Fatalf("older RecordInbound = %v, %v; want true, nil", inserted, err)
	}

	var messageCount, sessionCount int
	var direction, kind, text, title, username string
	var lastMessageAt int64
	err = st.Run(context.Background(), func(tx kernel.Tx) error {
		if err := tx.QueryRow(context.Background(), `SELECT COUNT(*) FROM telegram_inbound_messages`).Scan(&messageCount); err != nil {
			return err
		}
		if err := tx.QueryRow(context.Background(), `SELECT direction, message_kind, text FROM telegram_inbound_messages WHERE bot_id = ? AND update_id = ?`, first.BotID, first.UpdateID).Scan(&direction, &kind, &text); err != nil {
			return err
		}
		if err := tx.QueryRow(context.Background(), `SELECT COUNT(*) FROM telegram_sessions`).Scan(&sessionCount); err != nil {
			return err
		}
		return tx.QueryRow(context.Background(), `SELECT title, username, last_message_at FROM telegram_sessions WHERE bot_id = ? AND chat_id = ?`, first.BotID, first.ChatID).Scan(&title, &username, &lastMessageAt)
	})
	if err != nil {
		t.Fatalf("read persisted ingress: %v", err)
	}
	if messageCount != 2 || sessionCount != 1 {
		t.Fatalf("persisted counts messages=%d sessions=%d, want 2/1", messageCount, sessionCount)
	}
	if direction != "inbound" || kind != "text" || text != first.Text {
		t.Fatalf("persisted message = direction %q kind %q text %q", direction, kind, text)
	}
	if title != first.ChatTitle || username != first.ChatUsername || lastMessageAt != firstAt.Unix() {
		t.Fatalf("duplicate changed session = title %q username %q last_message_at %d", title, username, lastMessageAt)
	}
}

func TestRepositoryRecordInboundSQLiteCallback(t *testing.T) {
	path := filepath.Join(t.TempDir(), "telegram-callback.db")
	st, err := testsupport.OpenStore(path, "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	callback := InboundMessage{
		BotID:           101,
		UpdateID:        9002,
		ChatID:          8002,
		ChatType:        "group",
		UserID:          7002,
		MessageID:       6002,
		CallbackQueryID: "callback-1",
		MessageKind:     "callback",
		CallbackData:    "approve",
		ReceivedAt:      time.Date(2026, 9, 4, 13, 0, 0, 0, time.UTC),
	}
	inserted, err := NewRepository(st).RecordInbound(context.Background(), callback)
	if err != nil || !inserted {
		t.Fatalf("callback RecordInbound = %v, %v; want true, nil", inserted, err)
	}
	var kind, callbackID, callbackData string
	if err := st.Run(context.Background(), func(tx kernel.Tx) error {
		return tx.QueryRow(context.Background(), `SELECT message_kind, callback_query_id, callback_data FROM telegram_inbound_messages WHERE bot_id = ? AND update_id = ?`, callback.BotID, callback.UpdateID).Scan(&kind, &callbackID, &callbackData)
	}); err != nil {
		t.Fatal(err)
	}
	if kind != "callback" || callbackID != callback.CallbackQueryID || callbackData != callback.CallbackData {
		t.Fatalf("callback row = kind %q id %q data %q", kind, callbackID, callbackData)
	}
}

func TestRepositoryRecordInboundRejectsNonPositiveUpdateID(t *testing.T) {
	path := filepath.Join(t.TempDir(), "telegram-invalid-update.db")
	st, err := testsupport.OpenStore(path, "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	repository := NewRepository(st)
	base := InboundMessage{
		BotID:       101,
		UpdateID:    1,
		ChatID:      8003,
		ChatType:    "private",
		MessageKind: "text",
		Text:        "hello",
		ReceivedAt:  time.Date(2026, 9, 4, 14, 0, 0, 0, time.UTC),
	}
	for _, updateID := range []int64{0, -1} {
		message := base
		message.UpdateID = updateID
		inserted, err := repository.RecordInbound(context.Background(), message)
		if err == nil {
			t.Fatalf("update_id %d returned nil error", updateID)
		}
		if inserted {
			t.Fatalf("update_id %d returned inserted=true on validation failure", updateID)
		}
	}

	var messageCount, sessionCount int
	if err := st.Run(context.Background(), func(tx kernel.Tx) error {
		if err := tx.QueryRow(context.Background(), `SELECT COUNT(*) FROM telegram_inbound_messages`).Scan(&messageCount); err != nil {
			return err
		}
		return tx.QueryRow(context.Background(), `SELECT COUNT(*) FROM telegram_sessions`).Scan(&sessionCount)
	}); err != nil {
		t.Fatal(err)
	}
	if messageCount != 0 || sessionCount != 0 {
		t.Fatalf("invalid update ids persisted messages=%d sessions=%d", messageCount, sessionCount)
	}
}
