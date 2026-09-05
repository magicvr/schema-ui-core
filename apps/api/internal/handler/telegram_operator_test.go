package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	telegramstore "github.com/magicvr/schema-ui-core/apps/api/modules/channel/telegram/store"
)

type telegramOperatorTestSender struct {
	mu     sync.Mutex
	calls  []kernel.TelegramMessage
	err    error
	before func(context.Context, kernel.TelegramMessage) error
}

type telegramOperatorCapabilityStub struct {
	allowed bool
	err     error
	calls   []struct {
		botID  int64
		chatID int64
		force  bool
	}
}

func (s *telegramOperatorCapabilityStub) Check(_ context.Context, botID, chatID int64, force bool) (bool, error) {
	s.calls = append(s.calls, struct {
		botID  int64
		chatID int64
		force  bool
	}{botID: botID, chatID: chatID, force: force})
	return s.allowed, s.err
}

func (s *telegramOperatorTestSender) Send(ctx context.Context, message kernel.TelegramMessage) error {
	s.mu.Lock()
	s.calls = append(s.calls, message)
	err := s.err
	before := s.before
	s.mu.Unlock()
	if before != nil {
		if err := before(ctx, message); err != nil {
			return err
		}
	}
	return err
}

func (s *telegramOperatorTestSender) setError(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.err = err
}

func (s *telegramOperatorTestSender) callCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.calls)
}

type telegramOperatorTestFixture struct {
	repository *telegramstore.Repository
	sender     *telegramOperatorTestSender
	handler    *TelegramOperatorHandler
	state      *string
	business   *bool
}

func newTelegramOperatorTestFixture(t *testing.T) telegramOperatorTestFixture {
	t.Helper()
	st, err := testsupport.OpenStore(filepath.Join(t.TempDir(), "telegram-operator.db"), "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	repository := telegramstore.NewRepository(st)
	if inserted, err := repository.RecordInbound(context.Background(), telegramstore.InboundMessage{
		BotID: 101, UpdateID: 9201, ChatID: 8001, ChatType: "private",
		ChatTitle: "Ada Lovelace", ChatUsername: "ada", UserID: 7001,
		MessageID: 6001, MessageKind: "text", Text: "hello", SenderUsername: "ada",
		ReceivedAt: time.Date(2026, 9, 4, 15, 0, 0, 0, time.UTC),
	}); err != nil || !inserted {
		t.Fatalf("seed Telegram session = %v, %v; want true, nil", inserted, err)
	}
	if inserted, err := repository.RecordInbound(context.Background(), telegramstore.InboundMessage{
		BotID: 101, UpdateID: 9202, ChatID: 8002, ChatType: "group",
		MessageKind: "callback", CallbackQueryID: "callback-only", CallbackData: "approve",
		ReceivedAt: time.Date(2026, 9, 4, 15, 1, 0, 0, time.UTC),
	}); err != nil || !inserted {
		t.Fatalf("seed callback-only Telegram row = %v, %v; want true, nil", inserted, err)
	}
	state := "running"
	business := false
	sender := &telegramOperatorTestSender{}
	handler := NewTelegramOperatorHandler(
		func() (int64, string, string, string) { return 101, state, "polling", "bot-token" },
		func() bool { return business },
		repository,
		sender,
	)
	return telegramOperatorTestFixture{repository: repository, sender: sender, handler: handler, state: &state, business: &business}
}

func telegramOperatorTestRequest(t *testing.T, method, path, body string, permissions ...string) *http.Request {
	t.Helper()
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	return req.WithContext(auth.WithIdentity(req.Context(), account.User{
		ID: "operator-test", Name: "Operator Test", Permissions: permissions,
	}))
}

func serveTelegramOperator(t *testing.T, handler http.Handler, req *http.Request) (*httptest.ResponseRecorder, map[string]any) {
	t.Helper()
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	var body map[string]any
	if recorder.Body.Len() > 0 {
		if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
			t.Fatalf("decode Telegram operator response %q: %v", recorder.Body.String(), err)
		}
	}
	return recorder, body
}

func TestTelegramOperatorHandlerSessionsTimelineSendAndRetry(t *testing.T) {
	fixture := newTelegramOperatorTestFixture(t)
	read := telegramOperatorReadPermission
	write := telegramOperatorWritePermission

	if recorder, body := serveTelegramOperator(t, fixture.handler, httptest.NewRequest(http.MethodGet, "/api/channel/telegram/operator/sessions", nil)); recorder.Code != http.StatusUnauthorized || body["error"] != "UNAUTHENTICATED" {
		t.Fatalf("anonymous sessions response = %d %v, want 401 UNAUTHENTICATED", recorder.Code, body)
	}
	if recorder, body := serveTelegramOperator(t, fixture.handler, telegramOperatorTestRequest(t, http.MethodGet, "/api/channel/telegram/operator/sessions", "")); recorder.Code != http.StatusForbidden || body["error"] != "FORBIDDEN" {
		t.Fatalf("unpermissioned sessions response = %d %v, want 403 FORBIDDEN", recorder.Code, body)
	}
	if recorder, body := serveTelegramOperator(t, fixture.handler, telegramOperatorTestRequest(t, http.MethodPost, "/api/channel/telegram/operator/sessions/8001/messages", `{"requestId":"no-write","text":"blocked"}`, read)); recorder.Code != http.StatusForbidden || body["error"] != "FORBIDDEN" {
		t.Fatalf("read-only send response = %d %v, want 403 FORBIDDEN", recorder.Code, body)
	}

	recorder, body := serveTelegramOperator(t, fixture.handler, telegramOperatorTestRequest(t, http.MethodGet, "/api/channel/telegram/operator/sessions", "", read))
	if recorder.Code != http.StatusOK || body["total"] != float64(1) {
		t.Fatalf("sessions response = %d %v, want one visible session", recorder.Code, body)
	}
	sessions, ok := body["items"].([]any)
	if !ok || len(sessions) != 1 || sessions[0].(map[string]any)["chatId"] != "8001" {
		t.Fatalf("sessions items = %v, want chatId string 8001 only", body["items"])
	}

	recorder, body = serveTelegramOperator(t, fixture.handler, telegramOperatorTestRequest(t, http.MethodGet, "/api/channel/telegram/operator/sessions/8001/messages", "", read))
	if recorder.Code != http.StatusOK || body["total"] != float64(1) {
		t.Fatalf("initial timeline response = %d %v, want one inbound row", recorder.Code, body)
	}
	initialItems := body["items"].([]any)
	initial := initialItems[0].(map[string]any)
	if initial["direction"] != "inbound" || initial["status"] != "received" || initial["updateId"] != "9201" || initial["text"] != "hello" {
		t.Fatalf("initial timeline item = %v, want string ids and inbound text", initial)
	}

	var beforeSendErr error
	fixture.sender.before = func(ctx context.Context, message kernel.TelegramMessage) error {
		pending, err := fixture.repository.GetOutbound(ctx, 101, 8001, "root-1")
		if err != nil {
			beforeSendErr = err
			return nil
		}
		if pending.Status != "pending" {
			beforeSendErr = errors.New("outbound row was not pending before external send")
		}
		if message.ChatID != "8001" || message.Text != "reply" {
			beforeSendErr = errors.New("sender received an unexpected Telegram message")
		}
		return nil
	}
	recorder, body = serveTelegramOperator(t, fixture.handler, telegramOperatorTestRequest(t, http.MethodPost, "/api/channel/telegram/operator/sessions/8001/messages", `{"requestId":"root-1","text":"reply"}`, write))
	if beforeSendErr != nil {
		t.Fatal(beforeSendErr)
	}
	if recorder.Code != http.StatusCreated || body["status"] != "sent" || body["requestId"] != "root-1" || fixture.sender.callCount() != 1 {
		t.Fatalf("first send response = %d %v calls=%d, want 201 sent and one send", recorder.Code, body, fixture.sender.callCount())
	}

	// A terminal replay is an idempotent read and never calls the external
	// sender again; a different payload under the same request id conflicts.
	fixture.sender.before = nil
	recorder, body = serveTelegramOperator(t, fixture.handler, telegramOperatorTestRequest(t, http.MethodPost, "/api/channel/telegram/operator/sessions/8001/messages", `{"requestId":"root-1","text":"reply"}`, write))
	if recorder.Code != http.StatusOK || body["status"] != "sent" || fixture.sender.callCount() != 1 {
		t.Fatalf("terminal replay = %d %v calls=%d, want 200 sent and no second send", recorder.Code, body, fixture.sender.callCount())
	}
	recorder, body = serveTelegramOperator(t, fixture.handler, telegramOperatorTestRequest(t, http.MethodPost, "/api/channel/telegram/operator/sessions/8001/messages", `{"requestId":"root-1","text":"changed"}`, write))
	if recorder.Code != http.StatusConflict || body["error"] != "TELEGRAM_REQUEST_CONFLICT" || fixture.sender.callCount() != 1 {
		t.Fatalf("mismatched replay = %d %v calls=%d, want conflict and no send", recorder.Code, body, fixture.sender.callCount())
	}

	fixture.sender.setError(errors.New("bot api unavailable"))
	recorder, body = serveTelegramOperator(t, fixture.handler, telegramOperatorTestRequest(t, http.MethodPost, "/api/channel/telegram/operator/sessions/8001/messages", `{"requestId":"failed-1","text":"will fail"}`, write))
	if recorder.Code != http.StatusBadGateway || body["error"] != "TELEGRAM_SEND_FAILED" || fixture.sender.callCount() != 2 {
		t.Fatalf("failed send = %d %v calls=%d, want 502 cataloged failure", recorder.Code, body, fixture.sender.callCount())
	}
	failedItem := body["item"].(map[string]any)
	if failedItem["status"] != "failed" || failedItem["requestId"] != "failed-1" {
		t.Fatalf("failed item = %v, want durable failed state", failedItem)
	}

	fixture.sender.setError(nil)
	recorder, body = serveTelegramOperator(t, fixture.handler, telegramOperatorTestRequest(t, http.MethodPost, "/api/channel/telegram/operator/sessions/8001/messages/failed-1/retry", `{"requestId":"retry-1"}`, write))
	if recorder.Code != http.StatusCreated || body["status"] != "sent" || body["retryOf"] != "failed-1" || fixture.sender.callCount() != 3 {
		t.Fatalf("retry response = %d %v calls=%d, want 201 sent rooted at failed-1", recorder.Code, body, fixture.sender.callCount())
	}
	recorder, body = serveTelegramOperator(t, fixture.handler, telegramOperatorTestRequest(t, http.MethodPost, "/api/channel/telegram/operator/sessions/8001/messages/failed-1/retry", `{"requestId":"retry-2"}`, write))
	if recorder.Code != http.StatusConflict || body["error"] != "TELEGRAM_RETRY_NOT_ALLOWED" || fixture.sender.callCount() != 3 {
		t.Fatalf("retry after success = %d %v calls=%d, want retry-not-allowed", recorder.Code, body, fixture.sender.callCount())
	}

	if _, created, err := fixture.repository.CreatePending(context.Background(), 101, 8001, "pending-1", "already pending"); err != nil || !created {
		t.Fatalf("seed pending request = created %v err %v, want created", created, err)
	}
	recorder, body = serveTelegramOperator(t, fixture.handler, telegramOperatorTestRequest(t, http.MethodPost, "/api/channel/telegram/operator/sessions/8001/messages", `{"requestId":"pending-1","text":"already pending"}`, write))
	if recorder.Code != http.StatusConflict || body["error"] != "TELEGRAM_REQUEST_IN_PROGRESS" || fixture.sender.callCount() != 3 {
		t.Fatalf("pending replay = %d %v calls=%d, want in-progress and no send", recorder.Code, body, fixture.sender.callCount())
	}
	if err := fixture.repository.MarkFailed(context.Background(), 101, "pending-1", "test cleanup"); err != nil {
		t.Fatalf("cleanup pending request: %v", err)
	}

	recorder, body = serveTelegramOperator(t, fixture.handler, telegramOperatorTestRequest(t, http.MethodPost, "/api/channel/telegram/operator/sessions/8999/messages", `{"requestId":"unknown-chat","text":"blocked"}`, write))
	if recorder.Code != http.StatusNotFound || body["error"] != "TELEGRAM_CHAT_NOT_FOUND" {
		t.Fatalf("unknown chat send = %d %v, want 404 Telegram chat not found", recorder.Code, body)
	}
	recorder, body = serveTelegramOperator(t, fixture.handler, telegramOperatorTestRequest(t, http.MethodPost, "/api/channel/telegram/operator/sessions/8001/messages", `{"requestId":"bad+id","text":"blocked"}`, write))
	if recorder.Code != http.StatusBadRequest || body["error"] != "INVALID_BODY" {
		t.Fatalf("unsafe request id = %d %v, want 400 INVALID_BODY", recorder.Code, body)
	}

	recorder, body = serveTelegramOperator(t, fixture.handler, telegramOperatorTestRequest(t, http.MethodGet, "/api/channel/telegram/operator/sessions/8001/messages", "", read))
	if recorder.Code != http.StatusOK || body["total"] != float64(5) {
		t.Fatalf("final timeline response = %d %v, want inbound + 4 outbound attempts", recorder.Code, body)
	}
	finalItems := body["items"].([]any)
	if len(finalItems) != 5 {
		t.Fatalf("final timeline items = %d, want 5", len(finalItems))
	}

	callCount := fixture.sender.callCount()
	*fixture.state = "idle"
	recorder, body = serveTelegramOperator(t, fixture.handler, telegramOperatorTestRequest(t, http.MethodPost, "/api/channel/telegram/operator/sessions/8001/messages", `{"requestId":"runtime-down","text":"blocked"}`, write))
	if recorder.Code != http.StatusConflict || body["error"] != "TELEGRAM_OPERATOR_UNAVAILABLE" || fixture.sender.callCount() != callCount {
		t.Fatalf("runtime unavailable = %d %v calls=%d, want 409 and no sender", recorder.Code, body, fixture.sender.callCount())
	}
	*fixture.state = "running"
	*fixture.business = true
	recorder, body = serveTelegramOperator(t, fixture.handler, telegramOperatorTestRequest(t, http.MethodGet, "/api/channel/telegram/operator/sessions", "", read))
	if recorder.Code != http.StatusConflict || body["error"] != "TELEGRAM_OPERATOR_UNAVAILABLE" {
		t.Fatalf("business-handler runtime gate = %d %v, want 409 unavailable", recorder.Code, body)
	}
}

func TestTelegramOperatorCapabilityRouteUsesReadPermissionRefreshAndFailClosed(t *testing.T) {
	fixture := newTelegramOperatorTestFixture(t)
	capability := &telegramOperatorCapabilityStub{allowed: true}
	handler := NewTelegramOperatorHandler(
		func() (int64, string, string, string) { return 101, *fixture.state, "polling", "bot-token" },
		func() bool { return *fixture.business },
		fixture.repository,
		fixture.sender,
		capability,
	)
	path := "/api/channel/telegram/operator/sessions/8001/capability"

	if recorder, body := serveTelegramOperator(t, handler, httptest.NewRequest(http.MethodGet, path, nil)); recorder.Code != http.StatusUnauthorized || body["error"] != "UNAUTHENTICATED" {
		t.Fatalf("anonymous capability response = %d %v, want 401", recorder.Code, body)
	}
	if recorder, body := serveTelegramOperator(t, handler, telegramOperatorTestRequest(t, http.MethodGet, path, "")); recorder.Code != http.StatusForbidden || body["error"] != "FORBIDDEN" {
		t.Fatalf("unpermissioned capability response = %d %v, want 403", recorder.Code, body)
	}
	recorder, body := serveTelegramOperator(t, handler, telegramOperatorTestRequest(t, http.MethodGet, path, "", telegramOperatorReadPermission))
	if recorder.Code != http.StatusOK || body["chatId"] != "8001" || body["canSend"] != true {
		t.Fatalf("capability response = %d %v, want allowed 8001", recorder.Code, body)
	}
	recorder, body = serveTelegramOperator(t, handler, telegramOperatorTestRequest(t, http.MethodGet, path+"?refresh=1", "", telegramOperatorReadPermission))
	if recorder.Code != http.StatusOK || body["canSend"] != true || len(capability.calls) != 2 || !capability.calls[1].force {
		t.Fatalf("forced capability response = %d %v calls=%v, want force=true", recorder.Code, body, capability.calls)
	}
	recorder, body = serveTelegramOperator(t, handler, telegramOperatorTestRequest(t, http.MethodGet, path+"?refresh=0", "", telegramOperatorReadPermission))
	if recorder.Code != http.StatusBadRequest || body["error"] != "INVALID_BODY" || len(capability.calls) != 2 {
		t.Fatalf("invalid refresh response = %d %v calls=%v, want 400 and no capability call", recorder.Code, body, capability.calls)
	}

	capability.allowed = false
	recorder, body = serveTelegramOperator(t, handler, telegramOperatorTestRequest(t, http.MethodGet, path, "", telegramOperatorReadPermission))
	if recorder.Code != http.StatusOK || body["canSend"] != false {
		t.Fatalf("denied capability response = %d %v, want 200 false", recorder.Code, body)
	}
	capability.err = errors.New("probe failed")
	recorder, body = serveTelegramOperator(t, handler, telegramOperatorTestRequest(t, http.MethodGet, path+"?refresh=1", "", telegramOperatorReadPermission))
	if recorder.Code != http.StatusBadGateway || body["error"] != "TELEGRAM_CAPABILITY_UNAVAILABLE" {
		t.Fatalf("capability error response = %d %v, want 502 cataloged unavailable", recorder.Code, body)
	}
}

type telegramOperatorMarkSentFailureRepository struct {
	*telegramstore.Repository
	err error
}

func (r *telegramOperatorMarkSentFailureRepository) MarkSent(context.Context, int64, string) error {
	return r.err
}

func TestTelegramOperatorHandlerKeepsPendingWhenSentFinalizationFails(t *testing.T) {
	fixture := newTelegramOperatorTestFixture(t)
	failingRepository := &telegramOperatorMarkSentFailureRepository{
		Repository: fixture.repository,
		err:        errors.New("forced post-send state failure"),
	}
	var beforeSendErr error
	fixture.sender.before = func(ctx context.Context, _ kernel.TelegramMessage) error {
		message, err := fixture.repository.GetOutbound(ctx, 101, 8001, "state-failure")
		if err != nil {
			beforeSendErr = err
			return nil
		}
		if message.Status != "pending" {
			beforeSendErr = errors.New("post-send test did not observe pending row")
		}
		return nil
	}
	handler := NewTelegramOperatorHandler(
		func() (int64, string, string, string) { return 101, *fixture.state, "polling", "bot-token" },
		func() bool { return *fixture.business },
		failingRepository,
		fixture.sender,
	)
	recorder, body := serveTelegramOperator(t, handler, telegramOperatorTestRequest(t, http.MethodPost, "/api/channel/telegram/operator/sessions/8001/messages", `{"requestId":"state-failure","text":"reply once"}`, telegramOperatorWritePermission))
	if beforeSendErr != nil {
		t.Fatal(beforeSendErr)
	}
	if recorder.Code != http.StatusInternalServerError || body["error"] != "INTERNAL" || fixture.sender.callCount() != 1 {
		t.Fatalf("post-send state failure = %d %v calls=%d, want 500 INTERNAL and one external send", recorder.Code, body, fixture.sender.callCount())
	}
	message, err := fixture.repository.GetOutbound(context.Background(), 101, 8001, "state-failure")
	if err != nil || message.Status != "pending" {
		t.Fatalf("durable post-send state = %+v err=%v, want pending", message, err)
	}

	// Replaying the request is blocked by the durable pending state. The handler
	// cannot accidentally send a second time while the first result is unknown.
	recorder, body = serveTelegramOperator(t, handler, telegramOperatorTestRequest(t, http.MethodPost, "/api/channel/telegram/operator/sessions/8001/messages", `{"requestId":"state-failure","text":"reply once"}`, telegramOperatorWritePermission))
	if recorder.Code != http.StatusConflict || body["error"] != "TELEGRAM_REQUEST_IN_PROGRESS" || fixture.sender.callCount() != 1 {
		t.Fatalf("post-send replay = %d %v calls=%d, want 409 in-progress and no duplicate send", recorder.Code, body, fixture.sender.callCount())
	}
}

func TestTelegramOperatorHandlerKeepsPendingWhenTokenDisappearsAfterSend(t *testing.T) {
	fixture := newTelegramOperatorTestFixture(t)
	token := "bot-token"
	sender := &telegramOperatorTestSender{
		before: func(context.Context, kernel.TelegramMessage) error {
			token = ""
			return nil
		},
	}
	handler := NewTelegramOperatorHandler(
		func() (int64, string, string, string) { return 101, *fixture.state, "polling", token },
		func() bool { return *fixture.business },
		fixture.repository,
		sender,
	)

	recorder, body := serveTelegramOperator(t, handler, telegramOperatorTestRequest(t, http.MethodPost, "/api/channel/telegram/operator/sessions/8001/messages", `{"requestId":"token-race","text":"uncertain"}`, telegramOperatorWritePermission))
	if recorder.Code != http.StatusConflict || body["error"] != "TELEGRAM_OPERATOR_UNAVAILABLE" || sender.callCount() != 1 {
		t.Fatalf("token race response = %d %v calls=%d, want 409 unavailable and one send", recorder.Code, body, sender.callCount())
	}
	message, err := fixture.repository.GetOutbound(context.Background(), 101, 8001, "token-race")
	if err != nil || message.Status != "pending" {
		t.Fatalf("token race durable state = %+v err=%v, want pending", message, err)
	}

	// Restore runtime availability only to prove the durable pending row, not
	// the transient token state, prevents a duplicate external send.
	token = "bot-token"
	recorder, body = serveTelegramOperator(t, handler, telegramOperatorTestRequest(t, http.MethodPost, "/api/channel/telegram/operator/sessions/8001/messages", `{"requestId":"token-race","text":"uncertain"}`, telegramOperatorWritePermission))
	if recorder.Code != http.StatusConflict || body["error"] != "TELEGRAM_REQUEST_IN_PROGRESS" || sender.callCount() != 1 {
		t.Fatalf("token race replay = %d %v calls=%d, want in-progress and no duplicate send", recorder.Code, body, sender.callCount())
	}
}

func TestTelegramOperatorHandlerPaginationRuntimeAndRetryGuards(t *testing.T) {
	fixture := newTelegramOperatorTestFixture(t)
	state := "running"
	receiver := "polling"
	botID := int64(101)
	token := "bot-token"
	business := false
	handler := NewTelegramOperatorHandler(
		func() (int64, string, string, string) { return botID, state, receiver, token },
		func() bool { return business },
		fixture.repository,
		fixture.sender,
	)

	for _, test := range []struct {
		path string
		code string
	}{
		{path: "/api/channel/telegram/operator/sessions?page=0", code: "INVALID_PAGE"},
		{path: "/api/channel/telegram/operator/sessions?pageSize=101", code: "INVALID_PAGE_SIZE"},
		{path: "/api/channel/telegram/operator/sessions/8001/messages?pageSize=not-a-number", code: "INVALID_PAGE_SIZE"},
	} {
		recorder, body := serveTelegramOperator(t, handler, telegramOperatorTestRequest(t, http.MethodGet, test.path, "", telegramOperatorReadPermission))
		if recorder.Code != http.StatusBadRequest || body["error"] != test.code {
			t.Fatalf("pagination %s = %d %v, want 400 %s", test.path, recorder.Code, body, test.code)
		}
	}

	receiver = "webhook"
	recorder, body := serveTelegramOperator(t, handler, telegramOperatorTestRequest(t, http.MethodGet, "/api/channel/telegram/operator/sessions", "", telegramOperatorReadPermission))
	if recorder.Code != http.StatusOK || body["total"] != float64(1) {
		t.Fatalf("webhook runtime = %d %v, want available sessions", recorder.Code, body)
	}
	receiver = "none"
	recorder, body = serveTelegramOperator(t, handler, telegramOperatorTestRequest(t, http.MethodGet, "/api/channel/telegram/operator/sessions", "", telegramOperatorReadPermission))
	if recorder.Code != http.StatusConflict || body["error"] != "TELEGRAM_OPERATOR_UNAVAILABLE" {
		t.Fatalf("unknown receiver runtime = %d %v, want unavailable", recorder.Code, body)
	}
	state = "unconfigured"
	receiver = "polling"
	recorder, body = serveTelegramOperator(t, handler, telegramOperatorTestRequest(t, http.MethodGet, "/api/channel/telegram/operator/sessions", "", telegramOperatorReadPermission))
	if recorder.Code != http.StatusConflict || body["error"] != "TELEGRAM_OPERATOR_UNAVAILABLE" {
		t.Fatalf("unconfigured runtime = %d %v, want unavailable", recorder.Code, body)
	}
	state = "running"
	receiver = "polling"
	botID = 0
	recorder, body = serveTelegramOperator(t, handler, telegramOperatorTestRequest(t, http.MethodGet, "/api/channel/telegram/operator/sessions", "", telegramOperatorReadPermission))
	if recorder.Code != http.StatusConflict || body["error"] != "TELEGRAM_OPERATOR_UNAVAILABLE" {
		t.Fatalf("invalid bot id runtime = %d %v, want unavailable", recorder.Code, body)
	}
	botID = 101
	token = ""
	callCount := fixture.sender.callCount()
	recorder, body = serveTelegramOperator(t, handler, telegramOperatorTestRequest(t, http.MethodPost, "/api/channel/telegram/operator/sessions/8001/messages", `{"requestId":"missing-token","text":"blocked"}`, telegramOperatorWritePermission))
	if recorder.Code != http.StatusConflict || body["error"] != "TELEGRAM_OPERATOR_UNAVAILABLE" || fixture.sender.callCount() != callCount {
		t.Fatalf("empty token runtime = %d %v calls=%d, want unavailable and no sender", recorder.Code, body, fixture.sender.callCount())
	}
	failed, err := fixture.repository.GetOutbound(context.Background(), 101, 8001, "missing-token")
	if err != nil || failed.Status != "failed" || failed.ErrorMessage != "telegram operator became unavailable" {
		t.Fatalf("empty token durable state = %+v err=%v, want sanitized failed row", failed, err)
	}
	token = "bot-token"
	recorder, body = serveTelegramOperator(t, handler, telegramOperatorTestRequest(t, http.MethodPost, "/api/channel/telegram/operator/sessions/8001/messages/missing-request/retry", `{"requestId":"retry-missing"}`, telegramOperatorWritePermission))
	if recorder.Code != http.StatusNotFound || body["error"] != "TELEGRAM_REQUEST_NOT_FOUND" || fixture.sender.callCount() != callCount {
		t.Fatalf("unknown retry request = %d %v calls=%d, want 404 and no sender", recorder.Code, body, fixture.sender.callCount())
	}
}

func TestTelegramOperatorMiddlewareRejectsServiceCredentialWithoutOperatorScope(t *testing.T) {
	env := newAuthTestEnv(t)
	fixture := newTelegramOperatorTestFixture(t)
	env.mux.Handle("GET /api/channel/telegram/operator/sessions", env.a.Middleware(fixture.handler))

	expiresAt := time.Now().UTC().Add(time.Hour).Format(time.RFC3339)
	admin := adminToken(t, env)
	create := httptest.NewRecorder()
	env.mux.ServeHTTP(create, bearer(t, admin, http.MethodPost, "/api/service-credentials", `{"name":"Telegram Readless","scopes":["users.read"],"expiresAt":"`+expiresAt+`"}`))
	if create.Code != http.StatusCreated {
		t.Fatalf("service credential create = %d %s", create.Code, create.Body.String())
	}
	var created map[string]any
	if err := json.Unmarshal(create.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	secret, _ := created["secret"].(string)
	if secret == "" {
		t.Fatalf("service credential secret missing from create response: %v", created)
	}

	response := httptest.NewRecorder()
	env.mux.ServeHTTP(response, bearer(t, secret, http.MethodGet, "/api/channel/telegram/operator/sessions", ""))
	var body map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if response.Code != http.StatusForbidden || body["error"] != "FORBIDDEN" {
		t.Fatalf("service credential missing operator scope = %d %v, want 403 FORBIDDEN", response.Code, body)
	}
}

func TestTelegramOperatorHandlerKeepsPendingWhenRetryTokenDisappearsAfterSend(t *testing.T) {
	fixture := newTelegramOperatorTestFixture(t)
	if _, created, err := fixture.repository.CreatePending(context.Background(), 101, 8001, "retry-token-source", "retry me"); err != nil || !created {
		t.Fatalf("seed retry source = created %v err %v, want created", created, err)
	}
	if err := fixture.repository.MarkFailed(context.Background(), 101, "retry-token-source", "telegram send failed"); err != nil {
		t.Fatalf("mark retry source failed: %v", err)
	}

	token := "bot-token"
	sender := &telegramOperatorTestSender{
		before: func(context.Context, kernel.TelegramMessage) error {
			token = ""
			return nil
		},
	}
	handler := NewTelegramOperatorHandler(
		func() (int64, string, string, string) { return 101, *fixture.state, "polling", token },
		func() bool { return *fixture.business },
		fixture.repository,
		sender,
	)

	recorder, body := serveTelegramOperator(t, handler, telegramOperatorTestRequest(t, http.MethodPost, "/api/channel/telegram/operator/sessions/8001/messages/retry-token-source/retry", `{"requestId":"retry-token-race"}`, telegramOperatorWritePermission))
	if recorder.Code != http.StatusConflict || body["error"] != "TELEGRAM_OPERATOR_UNAVAILABLE" || sender.callCount() != 1 {
		t.Fatalf("retry token race response = %d %v calls=%d, want 409 unavailable and one send", recorder.Code, body, sender.callCount())
	}
	pending, err := fixture.repository.GetOutbound(context.Background(), 101, 8001, "retry-token-race")
	if err != nil || pending.Status != "pending" || pending.RetryOf != "retry-token-source" {
		t.Fatalf("retry token race durable state = %+v err=%v, want pending retry row", pending, err)
	}

	token = "bot-token"
	recorder, body = serveTelegramOperator(t, handler, telegramOperatorTestRequest(t, http.MethodPost, "/api/channel/telegram/operator/sessions/8001/messages/retry-token-source/retry", `{"requestId":"retry-token-race"}`, telegramOperatorWritePermission))
	if recorder.Code != http.StatusConflict || body["error"] != "TELEGRAM_REQUEST_IN_PROGRESS" || sender.callCount() != 1 {
		t.Fatalf("retry token race replay = %d %v calls=%d, want in-progress and no duplicate send", recorder.Code, body, sender.callCount())
	}
}
