// A-012 F-008 closure: Telegram-enabled composition-root acceptance for
// biz.digital-offer. The plan enables channel.telegram + biz.digital-offer and
// the app is assembled through the real newAppWithOptions/Fx graph using the
// existing fake Bot-API HTTP seam (telegramRuntimeOptions) — zero external
// Bot API dependency. The test proves:
//
//  1. the §6 price/buy/entitlements commands are registered on THE dispatcher
//     the webhook mounts (the Fx-injected *TelegramRuntime instance);
//  2. at least one command is reachable through the real webhook surface —
//     secret verification, bot identity, subject mapping, inbound persistence,
//     dispatcher routing and the outbound reply all execute in-process;
//  3. the aggregated manifest carries structured page route/schemaUrl entries
//     and sidebar navigation refs, while the schema/HTTP fail-closed
//     semantics are retained;
//  4. there is no physical DELETE route for offers (D-002 §2 negative bound).
package composition

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"go.uber.org/fx"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	telegraminternal "github.com/magicvr/schema-ui-core/apps/api/internal/channel/telegram"
	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// digitalOfferTelegramClient is the fake Bot-API seam for the Telegram-enabled
// composition test: it answers identity/webhook bookkeeping, captures every
// outbound sendMessage payload, and fails closed on any unexpected path.
type digitalOfferTelegramClient struct {
	mu  sync.Mutex
	msg []telegramCompositionMessage
}

type telegramCompositionMessage struct {
	ChatID string
	Text   string
}

func (c *digitalOfferTelegramClient) RoundTrip(r *http.Request) (*http.Response, error) {
	switch r.URL.Path {
	case "/botlive-bot-token/getMe":
		return telegramLifecycleJSONResponse(`{"ok":true,"result":{"id":201,"is_bot":true,"username":"f008_bot"}}`), nil
	case "/botlive-bot-token/deleteWebhook":
		return telegramLifecycleJSONResponse(`{"ok":true,"result":true}`), nil
	case "/botlive-bot-token/getUpdates":
		// A-012 F-008 regression: polling demand (digital-offer business
		// handlers registered before Start) starts the receiver during
		// reconcileStarted. Without a getUpdates seam the fake's default
		// branch answered ok:false, flipping the connection state to error
		// before the Fx Ready hook ran — a timing race that failed
		// TestDigitalOfferTelegramCompositionRoot on loaded runners
		// (LIFECYCLE_READY_FAILED). Block until the receiver is drained,
		// matching composition_telegram_lifecycle_test.go.
		<-r.Context().Done()
		return nil, r.Context().Err()
	case "/botlive-bot-token/sendMessage":
		var payload struct {
			ChatID string `json:"chat_id"`
			Text   string `json:"text"`
		}
		if body, err := io.ReadAll(r.Body); err == nil {
			_ = json.Unmarshal(body, &payload)
		}
		c.mu.Lock()
		c.msg = append(c.msg, telegramCompositionMessage{ChatID: payload.ChatID, Text: payload.Text})
		c.mu.Unlock()
		return telegramLifecycleJSONResponse(`{"ok":true,"result":{"message_id":1}}`), nil
	default:
		return telegramLifecycleJSONResponse(`{"ok":false,"description":"unexpected path"}`), nil
	}
}

func (c *digitalOfferTelegramClient) captured(chatID, wantSubstr string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, m := range c.msg {
		if m.ChatID == chatID && strings.Contains(m.Text, wantSubstr) {
			return true
		}
	}
	return false
}

func TestDigitalOfferTelegramCompositionRoot(t *testing.T) {
	client := &digitalOfferTelegramClient{}
	cfg := &config.Config{
		ProfileName:           string(kernel.ProfileCustom),
		ModulesEnabled:        append(append([]string(nil), digitalOfferCoreModules...), "channel.telegram", "biz.digital-offer"),
		AuthAccessTTL:         15 * time.Minute,
		AuthRefreshTTL:        30 * 24 * time.Hour,
		HTTPAddr:              "127.0.0.1:0",
		DBPath:                filepath.Join(t.TempDir(), "digitaloffer_telegram_composition.db"),
		TelegramBotToken:      "live-bot-token",
		TelegramWebhookSecret: "correct-secret",
		TelegramMasterKey:     "test-master-key",
	}
	seedHash, err := auth.HashPassword("admin-password", 4)
	if err != nil {
		t.Fatal(err)
	}

	var tr *TelegramRuntime
	var mux *http.ServeMux
	app, err := newAppWithOptions(
		cfg,
		"test-secret",
		seedHash,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		fx.Populate(&tr, &mux),
		fx.Supply(&telegramRuntimeOptions{APIBaseURL: "https://telegram.test", HTTPClient: &http.Client{Transport: client}}),
	)
	if err != nil {
		t.Fatalf("newAppWithOptions: %v", err)
	}
	startCtx, startCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer startCancel()
	if err := app.Start(startCtx); err != nil {
		t.Fatalf("app.Start: %v", err)
	}
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer stopCancel()
	defer func() { _ = app.Stop(stopCtx) }()

	// 1. The Fx-injected runtime carries a LIVE dispatcher (not the disabled
	// no-op), and the digital-offer §6 commands are registered on it — the
	// same instance the webhook mounts (A-012 F-008 finding 2).
	if tr == nil || tr.DispatcherState == nil {
		t.Fatal("expected a live Telegram dispatcher on the Fx-injected runtime")
	}
	if tr.Dispatcher != tr.DispatcherState {
		t.Fatal("runtime must expose the SAME dispatcher instance to the digital-offer wiring and the kernel port (A-014 F-001)")
	}
	if !tr.DispatcherState.HasBusinessHandlers() {
		t.Fatal("digital-offer commands must be registered on the live dispatcher (price/buy/entitlements)")
	}

	// 2. A command is reachable through the REAL webhook surface: secret
	// verified, bot identity resolved, subject mapped, inbound persisted, the
	// /price handler dispatched on the same dispatcher, and the reply sent
	// out through the same-process sender (captured by the fake client).
	webhookPayload := `{"update_id":1,"message":{"message_id":10,"from":{"id":678,"is_bot":false,"first_name":"Alice"},"chat":{"id":12345,"type":"private"},"text":"/price"}}`
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/channel/telegram/webhook", strings.NewReader(webhookPayload))
	req.Header.Set(telegraminternal.HeaderTelegramSecretToken, "correct-secret")
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("webhook /price = %d %s, want 200", rr.Code, rr.Body.String())
	}
	if !client.captured("12345", "在售数字服务") {
		t.Fatalf("webhook /price did not reach the digital-offer handler: captured %+v", client.snapshot())
	}

	// 3. The remaining frozen commands are registered and reachable on the
	// same dispatcher (direct drive; replies captured through the same seam).
	ctx := context.Background()
	if err := tr.DispatcherState.Dispatch(ctx, kernel.TelegramUpdate{
		UpdateID: 2, ChatID: "111", SubjectID: "sub-f008", Command: "entitlements", Text: "/entitlements",
	}, tr.Sender); err != nil {
		t.Fatalf("dispatch /entitlements: %v", err)
	}
	if !client.captured("111", "您当前没有有效权益") {
		t.Fatalf("/entitlements handler did not reply: captured %+v", client.snapshot())
	}
	if err := tr.DispatcherState.Dispatch(ctx, kernel.TelegramUpdate{
		UpdateID: 3, ChatID: "222", SubjectID: "sub-f008", Command: "buy", Text: "/buy",
	}, tr.Sender); err != nil {
		t.Fatalf("dispatch /buy: %v", err)
	}
	if !client.captured("222", "用法：/buy") {
		t.Fatalf("/buy handler did not reply with usage: captured %+v", client.snapshot())
	}

	// 3b. Explicit dispatcher identity by observable effect (A-014 F-001): the
	// webhook must dispatch through the very instance we drive directly.
	// Register a probe command on the Fx-injected dispatcher AFTER startup,
	// then drive the real webhook with it — the probe handler executing proves
	// pointer identity end-to-end through the webhook surface.
	probeRan := make(chan struct{})
	if err := tr.DispatcherState.RegisterCommand("probe_f008", func(ctx context.Context, upd kernel.TelegramUpdate) error {
		close(probeRan)
		return nil
	}); err != nil {
		t.Fatalf("register probe command: %v", err)
	}
	probePayload := `{"update_id":4,"message":{"message_id":14,"from":{"id":678,"is_bot":false,"first_name":"Alice"},"chat":{"id":12345,"type":"private"},"text":"/probe_f008"}}`
	rr = httptest.NewRecorder()
	probeReq := httptest.NewRequest(http.MethodPost, "/api/channel/telegram/webhook", strings.NewReader(probePayload))
	probeReq.Header.Set(telegraminternal.HeaderTelegramSecretToken, "correct-secret")
	mux.ServeHTTP(rr, probeReq)
	if rr.Code != http.StatusOK {
		t.Fatalf("webhook /probe_f008 = %d %s, want 200", rr.Code, rr.Body.String())
	}
	select {
	case <-probeRan:
	case <-time.After(5 * time.Second):
		t.Fatal("webhook did not dispatch through the Fx-injected dispatcher (probe command never ran)")
	}

	// 4. Structured manifest: parse the aggregated document and assert page
	// route/schemaUrl and sidebar navigation refs (A-012 F-008 finding 2).
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/.well-known/schema-ui/app-manifest.json", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("manifest = %d", rr.Code)
	}
	var doc struct {
		Pages []struct {
			PageID    string `json:"pageId"`
			Route     string `json:"route"`
			SchemaURL string `json:"schemaUrl"`
		} `json:"pages"`
		Navigation struct {
			Sidebar []struct {
				PageRef string `json:"pageRef"`
			} `json:"sidebar"`
		} `json:"navigation"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &doc); err != nil {
		t.Fatalf("manifest is not valid JSON: %v", err)
	}
	pages := map[string]struct{ route, schemaURL string }{}
	for _, p := range doc.Pages {
		pages[p.PageID] = struct{ route, schemaURL string }{p.Route, p.SchemaURL}
	}
	for _, want := range []struct {
		id, route, schemaURL string
	}{
		{"digitaloffer-offers", "/digitaloffer-offers", "/api/schema/digitaloffer-offers"},
		{"digitaloffer-entitlements", "/digitaloffer-entitlements", "/api/schema/digitaloffer-entitlements"},
		{"digitaloffer-purchases", "/digitaloffer-purchases", "/api/schema/digitaloffer-purchases"},
	} {
		got, ok := pages[want.id]
		if !ok || got.route != want.route || got.schemaURL != want.schemaURL {
			t.Fatalf("manifest page %s = %+v, want route %s schemaUrl %s", want.id, got, want.route, want.schemaURL)
		}
	}
	navRefs := map[string]bool{}
	for _, n := range doc.Navigation.Sidebar {
		navRefs[n.PageRef] = true
	}
	for _, ref := range []string{"digitaloffer-offers", "digitaloffer-entitlements", "digitaloffer-purchases"} {
		if !navRefs[ref] {
			t.Fatalf("manifest sidebar missing pageRef %s (refs %v)", ref, navRefs)
		}
	}

	// 5. Schema / HTTP fail-closed semantics retained with Telegram enabled
	// (unchanged from the disabled-surface acceptance in A-006 F-006).
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/schema/digitaloffer-offers", nil))
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated schema = %d, want 401", rr.Code)
	}
	token := loginCompositionAdmin(t, mux)
	for _, pageID := range []string{"digitaloffer-offers", "digitaloffer-entitlements", "digitaloffer-purchases"} {
		schemaReq := httptest.NewRequest(http.MethodGet, "/api/schema/"+pageID, nil)
		schemaReq.Header.Set("Authorization", "Bearer "+token)
		rr = httptest.NewRecorder()
		mux.ServeHTTP(rr, schemaReq)
		if rr.Code != http.StatusOK {
			t.Fatalf("schema %s = %d %s, want 200", pageID, rr.Code, rr.Body.String())
		}
	}
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/schema/digitaloffer-unknown", nil))
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("unknown schema anonymous = %d, want 401", rr.Code)
	}
	unknownReq := httptest.NewRequest(http.MethodGet, "/api/schema/digitaloffer-unknown", nil)
	unknownReq.Header.Set("Authorization", "Bearer "+token)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, unknownReq)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("unknown schema authenticated = %d, want 404", rr.Code)
	}

	// 6. Admin surface still permission-gated and public catalog open.
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/digitaloffer/offers", nil))
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated admin offers = %d, want 401", rr.Code)
	}
	adminReq := httptest.NewRequest(http.MethodGet, "/api/digitaloffer/offers", nil)
	adminReq.Header.Set("Authorization", "Bearer "+token)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, adminReq)
	if rr.Code != http.StatusOK {
		t.Fatalf("admin offers = %d %s, want 200", rr.Code, rr.Body.String())
	}
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/biz/offers", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("public catalog = %d %s, want 200", rr.Code, rr.Body.String())
	}

	// 7. Negative bound: no physical delete route for offers (D-002 §2 /
	// A-012 F-004 suggestion). The PATCH-only {id} pattern must answer
	// 405/404, never a success.
	del := httptest.NewRequest(http.MethodDelete, "/api/digitaloffer/offers/offer-x", nil)
	del.Header.Set("Authorization", "Bearer "+token)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, del)
	if rr.Code != http.StatusMethodNotAllowed && rr.Code != http.StatusNotFound {
		t.Fatalf("DELETE /api/digitaloffer/offers/{id} = %d, want 405/404 (no physical delete)", rr.Code)
	}
}

func (c *digitalOfferTelegramClient) snapshot() []telegramCompositionMessage {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]telegramCompositionMessage(nil), c.msg...)
}
