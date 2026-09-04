package store

import (
	"context"
	"errors"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

func TestRepositoryOperatorProjectionAndOutboundStateSQLite(t *testing.T) {
	ctx := context.Background()
	st, err := testsupport.OpenStore(filepath.Join(t.TempDir(), "telegram-operator.db"), "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	repository := NewRepository(st)
	baseAt := time.Date(2026, 9, 4, 15, 0, 0, 0, time.UTC)

	// Text activity creates a visible session; callback-only activity remains
	// persisted for ingress auditability but must not become an operator chat.
	for _, message := range []InboundMessage{
		{
			BotID: 101, UpdateID: 9201, ChatID: 8001, ChatType: "private",
			ChatTitle: "Ada Lovelace", ChatUsername: "ada", UserID: 7001,
			MessageID: 6001, MessageKind: "text", Text: "hello",
			SenderUsername: "ada", ReceivedAt: baseAt,
		},
		{
			BotID: 101, UpdateID: 9202, ChatID: 8002, ChatType: "group",
			UserID: 7002, MessageID: 6002, CallbackQueryID: "callback-operator",
			MessageKind: "callback", CallbackData: "approve", ReceivedAt: baseAt.Add(time.Minute),
		},
	} {
		inserted, recordErr := repository.RecordInbound(ctx, message)
		if recordErr != nil || !inserted {
			t.Fatalf("RecordInbound(%d) = %v, %v; want true, nil", message.UpdateID, inserted, recordErr)
		}
	}

	sessions, total, err := repository.ListSessions(ctx, 101, 1, 20)
	if err != nil {
		t.Fatalf("ListSessions: %v", err)
	}
	if total != 1 || len(sessions) != 1 || sessions[0].ChatID != 8001 {
		t.Fatalf("sessions = %+v total=%d, want only text chat 8001", sessions, total)
	}
	if exists, err := repository.SessionExists(ctx, 101, 8001); err != nil || !exists {
		t.Fatalf("text session exists = %v, %v; want true, nil", exists, err)
	}
	if exists, err := repository.SessionExists(ctx, 101, 8002); err != nil || exists {
		t.Fatalf("callback-only session exists = %v, %v; want false, nil", exists, err)
	}

	entries, total, err := repository.ListTimeline(ctx, 101, 8001, 1, 20)
	if err != nil {
		t.Fatalf("initial ListTimeline: %v", err)
	}
	if total != 1 || len(entries) != 1 || entries[0].Direction != "inbound" || entries[0].Text != "hello" {
		t.Fatalf("initial timeline = %+v total=%d, want one inbound text row", entries, total)
	}

	first, created, err := repository.CreatePending(ctx, 101, 8001, "root-1", "reply")
	if err != nil || !created {
		t.Fatalf("CreatePending first = %+v, %v, %v; want created", first, created, err)
	}
	if first.Status != "pending" || first.RetryRoot != "root-1" || first.RetryOf != "" {
		t.Fatalf("first outbound = %+v, want root pending without retry_of", first)
	}

	if _, created, err := repository.CreatePending(ctx, 101, 8001, "root-1", "reply"); !errors.Is(err, ErrRequestInProgress) || created {
		t.Fatalf("pending same-request replay = created %v err %v, want false ErrRequestInProgress", created, err)
	}
	if _, _, err := repository.CreatePending(ctx, 101, 8001, "root-1", "different"); !errors.Is(err, ErrRequestConflict) {
		t.Fatalf("pending mismatched replay err = %v, want ErrRequestConflict", err)
	}

	if err := repository.MarkFailed(ctx, 101, "root-1", "  telegram\n send failed "); err != nil {
		t.Fatalf("MarkFailed root: %v", err)
	}
	failed, created, err := repository.CreatePending(ctx, 101, 8001, "root-1", "reply")
	if err != nil || created || failed.Status != "failed" {
		t.Fatalf("terminal same-request replay = %+v created=%v err=%v, want failed replay", failed, created, err)
	}
	if failed.ErrorMessage != "telegram send failed" {
		t.Fatalf("failed diagnostic = %q, want normalized bounded diagnostic", failed.ErrorMessage)
	}

	retry, created, err := repository.CreateRetry(ctx, 101, 8001, "root-1", "retry-1")
	if err != nil || !created {
		t.Fatalf("CreateRetry = %+v, %v, %v; want created", retry, created, err)
	}
	if retry.Status != "pending" || retry.RetryRoot != "root-1" || retry.RetryOf != "root-1" || retry.Text != "reply" {
		t.Fatalf("retry outbound = %+v, want pending root-1 retry", retry)
	}
	if _, _, err := repository.CreateRetry(ctx, 101, 8001, "root-1", "retry-2"); !errors.Is(err, ErrRequestInProgress) {
		t.Fatalf("second retry while pending err = %v, want ErrRequestInProgress", err)
	}
	if err := repository.MarkSent(ctx, 101, "retry-1"); err != nil {
		t.Fatalf("MarkSent retry: %v", err)
	}
	if _, _, err := repository.CreateRetry(ctx, 101, 8001, "root-1", "retry-3"); !errors.Is(err, ErrRetryNotAllowed) {
		t.Fatalf("retry after root sent err = %v, want ErrRetryNotAllowed", err)
	}
	if err := repository.MarkSent(ctx, 101, "retry-1"); !errors.Is(err, ErrOutboundNotPending) {
		t.Fatalf("second MarkSent err = %v, want ErrOutboundNotPending", err)
	}
	if inserted, err := repository.RecordInbound(ctx, InboundMessage{
		BotID: 101, UpdateID: 9203, ChatID: 8003, ChatType: "private",
		MessageKind: "text", Text: "second chat", ReceivedAt: baseAt.Add(2 * time.Minute),
	}); err != nil || !inserted {
		t.Fatalf("seed second Telegram session = %v, %v; want true, nil", inserted, err)
	}
	if _, _, err := repository.CreatePending(ctx, 101, 8003, "root-1", "different chat"); !errors.Is(err, ErrRequestConflict) {
		t.Fatalf("same request id in another chat err = %v, want ErrRequestConflict", err)
	}

	if _, created, err := repository.CreatePending(ctx, 101, 8001, "secret-test", "safe text"); err != nil || !created {
		t.Fatalf("secret diagnostic fixture = created %v err %v, want created", created, err)
	}
	if err := repository.MarkFailed(ctx, 101, "secret-test", "bot token=super-secret raw downstream response"); err != nil {
		t.Fatalf("redact downstream diagnostic: %v", err)
	}
	secretRow, err := repository.GetOutbound(ctx, 101, 8001, "secret-test")
	if err != nil {
		t.Fatalf("read redacted diagnostic: %v", err)
	}
	if secretRow.ErrorMessage != "telegram send failed" {
		t.Fatalf("persisted downstream diagnostic = %q, want fixed generic category", secretRow.ErrorMessage)
	}

	entries, total, err = repository.ListTimeline(ctx, 101, 8001, 1, 20)
	if err != nil {
		t.Fatalf("final ListTimeline: %v", err)
	}
	if total != 4 || len(entries) != 4 {
		t.Fatalf("final timeline = %+v total=%d, want inbound + root + retry + diagnostic fixture", entries, total)
	}
	foundRetry := false
	for _, entry := range entries {
		if entry.Direction == "outbound" && entry.RequestID == "retry-1" {
			foundRetry = entry.Status == "sent" && entry.RetryOf == "root-1"
		}
	}
	if !foundRetry {
		t.Fatalf("final timeline did not contain sent retry rooted at root-1: %+v", entries)
	}

	if _, _, err := repository.CreatePending(ctx, 102, 8001, "other-bot", "no leak"); !errors.Is(err, ErrChatNotFound) {
		t.Fatalf("cross-bot CreatePending err = %v, want ErrChatNotFound", err)
	}
	if _, err := repository.GetOutbound(ctx, 102, 8001, "root-1"); !errors.Is(err, ErrRequestNotFound) {
		t.Fatalf("cross-bot GetOutbound err = %v, want ErrRequestNotFound", err)
	}
}

func TestRepositoryOutboundRejectsUnsafeRequestIDSQLite(t *testing.T) {
	st, err := testsupport.OpenStore(filepath.Join(t.TempDir(), "telegram-request-id.db"), "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	repository := NewRepository(st)
	for _, requestID := range []string{"../escape", "contains space", "plus+sign", ""} {
		if _, _, err := repository.CreatePending(context.Background(), 101, 8001, requestID, "reply"); err == nil {
			t.Fatalf("unsafe request id %q was accepted", requestID)
		}
	}
}

func TestRepositoryCreatePendingConcurrentRequestSQLite(t *testing.T) {
	ctx := context.Background()
	st, err := testsupport.OpenStore(filepath.Join(t.TempDir(), "telegram-concurrent.db"), "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	repository := NewRepository(st)
	if inserted, err := repository.RecordInbound(ctx, InboundMessage{
		BotID: 101, UpdateID: 9301, ChatID: 8001, ChatType: "private",
		MessageKind: "text", Text: "hello", ReceivedAt: time.Date(2026, 9, 4, 16, 0, 0, 0, time.UTC),
	}); err != nil || !inserted {
		t.Fatalf("seed concurrent Telegram session = %v, %v; want true, nil", inserted, err)
	}

	const competitors = 8
	start := make(chan struct{})
	var wg sync.WaitGroup
	created := 0
	var mu sync.Mutex
	var unexpected []error
	for i := 0; i < competitors; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_, didCreate, callErr := repository.CreatePending(ctx, 101, 8001, "concurrent-root", "reply")
			mu.Lock()
			defer mu.Unlock()
			if didCreate {
				created++
				return
			}
			if !errors.Is(callErr, ErrRequestInProgress) {
				unexpected = append(unexpected, callErr)
			}
		}()
	}
	close(start)
	wg.Wait()
	if created != 1 {
		t.Fatalf("concurrent CreatePending created=%d, want exactly 1", created)
	}
	if len(unexpected) != 0 {
		t.Fatalf("concurrent CreatePending unexpected errors = %v", unexpected)
	}
}
