package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/ratelimit"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/modules/wallet/subject"
)

func newTestSubjectStore(t *testing.T) *subject.Store {
	t.Helper()
	st, err := testsupport.OpenStore(":memory:", "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	return subject.NewStore(st)
}

func TestWebhook_UnconfiguredToken_Returns503(t *testing.T) {
	h := NewWebhookHandler(HandlerConfig{
		TokenGetter:  func() string { return "" },
		SecretGetter: func() string { return "secret123" },
	})

	req := httptest.NewRequest(http.MethodPost, "/api/channel/telegram/webhook", bytes.NewReader([]byte(`{}`)))
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 Service Unavailable, got %d", w.Code)
	}
}

func TestWebhook_SecretValidation_FailClosed(t *testing.T) {
	h := NewWebhookHandler(HandlerConfig{
		TokenGetter:  func() string { return "bot-token-123" },
		SecretGetter: func() string { return "correct-secret" },
		RateLimiters: ratelimit.NewProvider(),
	})

	// Missing header
	req1 := httptest.NewRequest(http.MethodPost, "/api/channel/telegram/webhook", bytes.NewReader([]byte(`{}`)))
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, req1)
	if w1.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for missing secret header, got %d", w1.Code)
	}

	// Incorrect secret
	req2 := httptest.NewRequest(http.MethodPost, "/api/channel/telegram/webhook", bytes.NewReader([]byte(`{}`)))
	req2.Header.Set(HeaderTelegramSecretToken, "wrong-secret")
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, req2)
	if w2.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for wrong secret header, got %d", w2.Code)
	}

	// Unconfigured secret on handler side -> fail closed 401
	hNoSecret := NewWebhookHandler(HandlerConfig{
		TokenGetter:  func() string { return "bot-token-123" },
		SecretGetter: func() string { return "" },
		RateLimiters: ratelimit.NewProvider(),
	})
	req3 := httptest.NewRequest(http.MethodPost, "/api/channel/telegram/webhook", bytes.NewReader([]byte(`{}`)))
	req3.Header.Set(HeaderTelegramSecretToken, "some-secret")
	w3 := httptest.NewRecorder()
	hNoSecret.ServeHTTP(w3, req3)
	if w3.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 when secret is unconfigured, got %d", w3.Code)
	}
}

func TestWebhook_MalformedJSON_Returns400(t *testing.T) {
	h := NewWebhookHandler(HandlerConfig{
		TokenGetter:  func() string { return "bot-token-123" },
		SecretGetter: func() string { return "correct-secret" },
		RateLimiters: ratelimit.NewProvider(),
	})

	req := httptest.NewRequest(http.MethodPost, "/api/channel/telegram/webhook", bytes.NewReader([]byte(`{invalid json`)))
	req.Header.Set(HeaderTelegramSecretToken, "correct-secret")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", w.Code)
	}
}

func TestWebhook_CommandDispatch_AndSubjectMapping(t *testing.T) {
	subStore := newTestSubjectStore(t)
	disp := NewDispatcher()
	sender := NewCaptureSender()

	var dispatchedUpdate kernel.TelegramUpdate
	commandCalled := false
	_ = disp.RegisterCommand("help", func(ctx context.Context, upd kernel.TelegramUpdate) error {
		commandCalled = true
		dispatchedUpdate = upd
		return nil
	})

	h := NewWebhookHandler(HandlerConfig{
		TokenGetter:  func() string { return "bot-token-123" },
		SecretGetter: func() string { return "my-secret" },
		RateLimiters: ratelimit.NewProvider(),
		SubjectStore: subStore,
		Dispatcher:   disp,
		Sender:       sender,
	})

	updateBody := UpdatePayload{
		UpdateID: 1001,
		Message: &MessagePayload{
			MessageID: 555,
			From: &UserPayload{
				ID:        987654,
				FirstName: "Alice",
			},
			Chat: &ChatPayload{
				ID:   123456,
				Type: "private",
			},
			Text: "/help@MyBot arg1 arg2",
		},
	}
	bodyBytes, _ := json.Marshal(updateBody)

	req := httptest.NewRequest(http.MethodPost, "/api/channel/telegram/webhook", bytes.NewReader(bodyBytes))
	req.Header.Set(HeaderTelegramSecretToken, "my-secret")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}
	if !commandCalled {
		t.Fatalf("expected command handler to be called")
	}
	if dispatchedUpdate.ChatID != "123456" || dispatchedUpdate.UserID != "987654" {
		t.Fatalf("unexpected ChatID or UserID: %+v", dispatchedUpdate)
	}
	if dispatchedUpdate.SubjectID == "" {
		t.Fatalf("expected SubjectID to be mapped via subject.Store")
	}

	// Verify subject in store
	sub, err := subStore.GetSubjectByExternalID(context.Background(), "telegram", "987654")
	if err != nil || sub.ID != dispatchedUpdate.SubjectID {
		t.Fatalf("subject store mismatch: %v, %+v", err, sub)
	}
}

func TestWebhook_UnknownCommand_SendsFallbackMessage(t *testing.T) {
	disp := NewDispatcher()
	sender := NewCaptureSender()

	h := NewWebhookHandler(HandlerConfig{
		TokenGetter:  func() string { return "bot-token-123" },
		SecretGetter: func() string { return "my-secret" },
		RateLimiters: ratelimit.NewProvider(),
		Dispatcher:   disp,
		Sender:       sender,
	})

	updateBody := UpdatePayload{
		UpdateID: 1002,
		Message: &MessagePayload{
			MessageID: 556,
			From: &UserPayload{
				ID: 111,
			},
			Chat: &ChatPayload{
				ID: 222,
			},
			Text: "/unregistered_cmd",
		},
	}
	bodyBytes, _ := json.Marshal(updateBody)

	req := httptest.NewRequest(http.MethodPost, "/api/channel/telegram/webhook", bytes.NewReader(bodyBytes))
	req.Header.Set(HeaderTelegramSecretToken, "my-secret")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}

	last := sender.Last()
	if last == nil {
		t.Fatalf("expected fallback message to be sent")
	}
	if last.ChatID != "222" || last.Text != kernel.DefaultTelegramUnknownCommandText {
		t.Fatalf("unexpected fallback message: %+v", last)
	}
}

func TestWebhook_CallbackDispatch(t *testing.T) {
	disp := NewDispatcher()
	callbackCalled := false
	var cbUpdate kernel.TelegramUpdate

	_ = disp.RegisterCallback("btn_confirm", func(ctx context.Context, upd kernel.TelegramUpdate) error {
		callbackCalled = true
		cbUpdate = upd
		return nil
	})

	h := NewWebhookHandler(HandlerConfig{
		TokenGetter:  func() string { return "bot-token-123" },
		SecretGetter: func() string { return "my-secret" },
		RateLimiters: ratelimit.NewProvider(),
		Dispatcher:   disp,
	})

	updateBody := UpdatePayload{
		UpdateID: 1003,
		CallbackQuery: &CallbackQueryPayload{
			ID:   "cb_1",
			Data: "btn_confirm",
			From: &UserPayload{
				ID: 333,
			},
			Message: &MessagePayload{
				Chat: &ChatPayload{
					ID: 444,
				},
			},
		},
	}
	bodyBytes, _ := json.Marshal(updateBody)

	req := httptest.NewRequest(http.MethodPost, "/api/channel/telegram/webhook", bytes.NewReader(bodyBytes))
	req.Header.Set(HeaderTelegramSecretToken, "my-secret")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}
	if !callbackCalled || cbUpdate.CallbackData != "btn_confirm" || cbUpdate.UserID != "333" || cbUpdate.ChatID != "444" {
		t.Fatalf("unexpected callback update: %+v", cbUpdate)
	}
}

func TestWebhook_RateLimiting_IPBucket(t *testing.T) {
	rlProvider := ratelimit.NewProvider()
	h := NewWebhookHandler(HandlerConfig{
		TokenGetter:  func() string { return "bot-token" },
		SecretGetter: func() string { return "sec" },
		RateLimiters: rlProvider,
	})

	// Fire 60 requests -> all should pass IP check
	for i := 0; i < RateLimitMaxIP; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/channel/telegram/webhook", bytes.NewReader([]byte(`{"update_id":1}`)))
		req.Header.Set(HeaderTelegramSecretToken, "sec")
		req.RemoteAddr = "198.51.100.1:12345"
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("request %d failed: status %d", i+1, w.Code)
		}
	}

	// 61st request -> 429
	req := httptest.NewRequest(http.MethodPost, "/api/channel/telegram/webhook", bytes.NewReader([]byte(`{"update_id":1}`)))
	req.Header.Set(HeaderTelegramSecretToken, "sec")
	req.RemoteAddr = "198.51.100.1:12345"
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 Too Many Requests, got %d", w.Code)
	}
	if retryAfter := w.Header().Get("Retry-After"); retryAfter == "" {
		t.Fatalf("expected Retry-After header on 429 response")
	}
}

func TestWebhook_RateLimiting_UserBucket(t *testing.T) {
	rlProvider := ratelimit.NewProvider()
	h := NewWebhookHandler(HandlerConfig{
		TokenGetter:  func() string { return "bot-token" },
		SecretGetter: func() string { return "sec" },
		RateLimiters: rlProvider,
	})

	// User limit is 20/min
	for i := 0; i < RateLimitMaxUser; i++ {
		updateBody := UpdatePayload{
			UpdateID: int64(100 + i),
			Message: &MessagePayload{
				From: &UserPayload{ID: 777},
				Chat: &ChatPayload{ID: int64(1000 + i)}, // different chat to avoid chat limit
			},
		}
		bodyBytes, _ := json.Marshal(updateBody)
		req := httptest.NewRequest(http.MethodPost, "/api/channel/telegram/webhook", bytes.NewReader(bodyBytes))
		req.Header.Set(HeaderTelegramSecretToken, "sec")
		req.RemoteAddr = "198.51.100." + strconv.Itoa(i+10) + ":12345" // different IP
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("request %d failed: status %d", i+1, w.Code)
		}
	}

	// 21st request for user 777 -> 429
	updateBody := UpdatePayload{
		UpdateID: 999,
		Message: &MessagePayload{
			From: &UserPayload{ID: 777},
			Chat: &ChatPayload{ID: 8888},
		},
	}
	bodyBytes, _ := json.Marshal(updateBody)
	req := httptest.NewRequest(http.MethodPost, "/api/channel/telegram/webhook", bytes.NewReader(bodyBytes))
	req.Header.Set(HeaderTelegramSecretToken, "sec")
	req.RemoteAddr = "198.51.100.99:12345"
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 on user rate limit, got %d", w.Code)
	}
}

func TestWebhook_RateLimiting_IPRecordsOnSecretFailure(t *testing.T) {
	rlProvider := ratelimit.NewProvider()
	h := NewWebhookHandler(HandlerConfig{
		TokenGetter:  func() string { return "bot-token" },
		SecretGetter: func() string { return "correct-sec" },
		RateLimiters: rlProvider,
	})

	// 60 requests with wrong secret
	for i := 0; i < RateLimitMaxIP; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/channel/telegram/webhook", bytes.NewReader([]byte(`{}`)))
		req.Header.Set(HeaderTelegramSecretToken, "wrong-sec")
		req.RemoteAddr = "203.0.113.50:5555"
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401 for wrong secret, got %d", w.Code)
		}
	}

	// 61st request with correct secret should now hit 429 IP rate limit
	req := httptest.NewRequest(http.MethodPost, "/api/channel/telegram/webhook", bytes.NewReader([]byte(`{}`)))
	req.Header.Set(HeaderTelegramSecretToken, "correct-sec")
	req.RemoteAddr = "203.0.113.50:5555"
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 IP limit exceeded after failed secret attempts, got %d", w.Code)
	}
}

func TestWebhook_RateLimiting_ChatBucket(t *testing.T) {
	rlProvider := ratelimit.NewProvider()
	h := NewWebhookHandler(HandlerConfig{
		TokenGetter:  func() string { return "bot-token" },
		SecretGetter: func() string { return "sec" },
		RateLimiters: rlProvider,
	})

	// Chat limit is 30/min
	for i := 0; i < RateLimitMaxChat; i++ {
		updateBody := UpdatePayload{
			UpdateID: int64(200 + i),
			Message: &MessagePayload{
				From: &UserPayload{ID: int64(5000 + i)}, // different users
				Chat: &ChatPayload{ID: 9999},            // same chat
			},
		}
		bodyBytes, _ := json.Marshal(updateBody)
		req := httptest.NewRequest(http.MethodPost, "/api/channel/telegram/webhook", bytes.NewReader(bodyBytes))
		req.Header.Set(HeaderTelegramSecretToken, "sec")
		req.RemoteAddr = "198.51.100." + strconv.Itoa(i+1) + ":12345" // different IPs
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("request %d failed: status %d", i+1, w.Code)
		}
	}

	// 31st request for same chat 9999 -> 429
	updateBody := UpdatePayload{
		UpdateID: 888,
		Message: &MessagePayload{
			From: &UserPayload{ID: 99999},
			Chat: &ChatPayload{ID: 9999},
		},
	}
	bodyBytes, _ := json.Marshal(updateBody)
	req := httptest.NewRequest(http.MethodPost, "/api/channel/telegram/webhook", bytes.NewReader(bodyBytes))
	req.Header.Set(HeaderTelegramSecretToken, "sec")
	req.RemoteAddr = "198.51.100.99:12345"
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 on chat rate limit, got %d", w.Code)
	}
	if retryAfter := w.Header().Get("Retry-After"); retryAfter == "" {
		t.Fatalf("expected Retry-After header on chat 429 response")
	}
}

func TestWebhook_SubjectMappingIdempotency(t *testing.T) {
	subStore := newTestSubjectStore(t)
	disp := NewDispatcher()

	var capturedSubjectIDs []string
	_ = disp.RegisterCommand("ping", func(ctx context.Context, upd kernel.TelegramUpdate) error {
		capturedSubjectIDs = append(capturedSubjectIDs, upd.SubjectID)
		return nil
	})

	h := NewWebhookHandler(HandlerConfig{
		TokenGetter:  func() string { return "bot-token-123" },
		SecretGetter: func() string { return "my-secret" },
		RateLimiters: ratelimit.NewProvider(),
		SubjectStore: subStore,
		Dispatcher:   disp,
	})

	updateBody := UpdatePayload{
		UpdateID: 1005,
		Message: &MessagePayload{
			From: &UserPayload{ID: 444555},
			Chat: &ChatPayload{ID: 123},
			Text: "/ping",
		},
	}
	bodyBytes, _ := json.Marshal(updateBody)

	// First call -> creates subject
	req1 := httptest.NewRequest(http.MethodPost, "/api/channel/telegram/webhook", bytes.NewReader(bodyBytes))
	req1.Header.Set(HeaderTelegramSecretToken, "my-secret")
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, req1)
	if w1.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w1.Code)
	}

	// Second call -> returns identical subject
	req2 := httptest.NewRequest(http.MethodPost, "/api/channel/telegram/webhook", bytes.NewReader(bodyBytes))
	req2.Header.Set(HeaderTelegramSecretToken, "my-secret")
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w2.Code)
	}

	if len(capturedSubjectIDs) != 2 {
		t.Fatalf("expected 2 captured calls, got %d", len(capturedSubjectIDs))
	}
	if capturedSubjectIDs[0] == "" || capturedSubjectIDs[0] != capturedSubjectIDs[1] {
		t.Fatalf("expected identical non-empty SubjectID on repeated webhook, got %v", capturedSubjectIDs)
	}
}
