package store

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	telegramstore "github.com/magicvr/schema-ui-core/apps/api/modules/channel/telegram/store"
)

func TestPostgresTelegramIngressRepositoryIdempotency(t *testing.T) {
	ctx := context.Background()
	st := postgresScratchDB(t, "r3telegramingress")
	repository := telegramstore.NewRepository(st)
	firstAt := time.Date(2026, 9, 4, 14, 0, 0, 0, time.UTC)
	first := telegramstore.InboundMessage{
		BotID:          101,
		UpdateID:       9101,
		ChatID:         8101,
		ChatType:       "private",
		ChatTitle:      "Ada Lovelace",
		ChatUsername:   "ada",
		UserID:         7101,
		MessageID:      6101,
		MessageKind:    "text",
		Text:           "hello from postgres",
		SenderUsername: "ada",
		ReceivedAt:     firstAt,
	}
	inserted, err := repository.RecordInbound(ctx, first)
	if err != nil || !inserted {
		t.Fatalf("first postgres RecordInbound = %v, %v; want true, nil", inserted, err)
	}

	duplicate := first
	duplicate.ChatTitle = "duplicate should not overwrite"
	duplicate.ChatUsername = "changed"
	duplicate.ReceivedAt = firstAt.Add(time.Hour)
	inserted, err = repository.RecordInbound(ctx, duplicate)
	if err != nil || inserted {
		t.Fatalf("duplicate postgres RecordInbound = %v, %v; want false, nil", inserted, err)
	}

	const competitors = 8
	var wg sync.WaitGroup
	errors := make(chan error, competitors)
	for i := 0; i < competitors; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			concurrent := first
			concurrent.UpdateID = 9102
			concurrent.MessageID = 6102
			concurrent.Text = "concurrent winner"
			concurrent.ReceivedAt = firstAt.Add(2 * time.Hour)
			gotNew, callErr := repository.RecordInbound(ctx, concurrent)
			if callErr != nil {
				errors <- callErr
				return
			}
			if gotNew {
				return
			}
		}()
	}
	wg.Wait()
	close(errors)
	for callErr := range errors {
		if callErr != nil {
			t.Fatalf("concurrent postgres RecordInbound: %v", callErr)
		}
	}

	var messageCount, concurrentCount, sessionCount int
	var title, username string
	var lastMessageAt int64
	err = st.Run(ctx, func(tx kernel.Tx) error {
		if err := tx.QueryRow(ctx, `SELECT COUNT(*) FROM telegram_inbound_messages WHERE bot_id = ? AND update_id = ?`, first.BotID, first.UpdateID).Scan(&messageCount); err != nil {
			return err
		}
		if err := tx.QueryRow(ctx, `SELECT COUNT(*) FROM telegram_inbound_messages WHERE bot_id = ? AND update_id = ?`, first.BotID, 9102).Scan(&concurrentCount); err != nil {
			return err
		}
		if err := tx.QueryRow(ctx, `SELECT COUNT(*) FROM telegram_sessions WHERE bot_id = ? AND chat_id = ?`, first.BotID, first.ChatID).Scan(&sessionCount); err != nil {
			return err
		}
		return tx.QueryRow(ctx, `SELECT title, username, last_message_at FROM telegram_sessions WHERE bot_id = ? AND chat_id = ?`, first.BotID, first.ChatID).Scan(&title, &username, &lastMessageAt)
	})
	if err != nil {
		t.Fatalf("read postgres ingress rows: %v", err)
	}
	if messageCount != 1 || concurrentCount != 1 || sessionCount != 1 {
		t.Fatalf("postgres counts first=%d concurrent=%d sessions=%d, want 1/1/1", messageCount, concurrentCount, sessionCount)
	}
	if title != first.ChatTitle || username != first.ChatUsername || lastMessageAt != firstAt.Add(2*time.Hour).Unix() {
		t.Fatalf("postgres duplicate changed session = title %q username %q last_message_at %d", title, username, lastMessageAt)
	}
}
