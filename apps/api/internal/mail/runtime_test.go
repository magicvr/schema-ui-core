package mail

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/compiled"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
)

func TestSecretsRoundTrip(t *testing.T) {
	key := []byte(strings.Repeat("k", 32))
	enc, err := EncryptSecret(key, "re-secret-value")
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	if strings.Contains(enc, "re-secret-value") || enc == "" {
		t.Fatal("ciphertext must be present and must not contain plaintext")
	}
	got, err := DecryptSecret(key, enc)
	if err != nil || got != "re-secret-value" {
		t.Fatalf("decrypt = %q, %v", got, err)
	}
	if _, err := DecryptSecret([]byte(strings.Repeat("x", 32)), enc); err == nil {
		t.Fatal("wrong key must fail closed")
	}
	if empty, err := EncryptSecret(key, ""); err != nil || empty != "" {
		t.Fatalf("empty secret = %q, %v; want untouched empty", empty, err)
	}
	if got, err := DecryptSecret(key, ""); err != nil || got != "" {
		t.Fatalf("empty decrypt = %q, %v", got, err)
	}
}

func TestLoadOrCreateMasterKey(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "mail-master.key")

	t.Run("env wins and is deterministic", func(t *testing.T) {
		a, err := LoadOrCreateMasterKey("operator-passphrase", filepath.Join(dir, "unused.key"))
		if err != nil {
			t.Fatal(err)
		}
		b, _ := LoadOrCreateMasterKey("operator-passphrase", filepath.Join(dir, "unused.key"))
		if len(a) != 32 || string(a) != string(b) {
			t.Fatalf("env-derived key len=%d stable=%v", len(a), string(a) == string(b))
		}
	})

	t.Run("file is created once and reused", func(t *testing.T) {
		a, err := LoadOrCreateMasterKey("", path)
		if err != nil {
			t.Fatal(err)
		}
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("key file missing: %v", err)
		}
		if runtime.GOOS != "windows" {
			if perms := info.Mode().Perm(); perms&0o077 != 0 {
				t.Fatalf("key file perms too open: %v", perms)
			}
		}
		b, err := LoadOrCreateMasterKey("", path)
		if err != nil || string(a) != string(b) {
			t.Fatalf("key not reused: %v", err)
		}
	})
}

// R7 C1 (Root D-007): seed-once semantics, hot switch applies to subsequent
// sends, failed switches keep the previous channel, and no read face leaks
// secret values.
func TestSwitcherHotSwitchSemantics(t *testing.T) {
	st := openOutboxStore(t)
	key := []byte(strings.Repeat("m", 32))

	var resendHits int
	resendSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resendHits++
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "resend-1"})
	}))
	defer resendSrv.Close()

	sw, err := NewSwitcher(st, key, SeedConfig{Channel: RuntimeChannelMock}, testSwitchLogger())
	if err != nil {
		t.Fatalf("NewSwitcher: %v", err)
	}
	sw.resendBaseURL = resendSrv.URL

	t.Run("default channel is mock and persists records", func(t *testing.T) {
		if err := sw.Send(context.Background(), msg("u@example.com", "s1", "b1")); err != nil {
			t.Fatalf("send via mock: %v", err)
		}
		view, err := sw.PublicView()
		if err != nil || view.Channel != RuntimeChannelMock {
			t.Fatalf("view = %+v, %v", view, err)
		}
	})

	t.Run("switch to resend applies to subsequent sends", func(t *testing.T) {
		fullKey := "re-key"
		from := "no-reply@example.com"
		view, err := sw.Update(context.Background(), UpdateRequest{
			Channel:      RuntimeChannelResend,
			ResendAPIKey: &fullKey,
			ResendFrom:   &from,
		})
		if err != nil {
			t.Fatalf("update: %v", err)
		}
		if view.Channel != RuntimeChannelResend {
			t.Fatalf("channel = %q, want resend", view.Channel)
		}
		if err := sw.Send(context.Background(), msg("user@example.com", "s2", "b2")); err != nil {
			t.Fatalf("send via resend: %v", err)
		}
		if resendHits != 1 {
			t.Fatalf("resend endpoint hits = %d, want 1", resendHits)
		}
	})

	t.Run("public view never carries secret values", func(t *testing.T) {
		raw, _ := json.Marshal(mustView(t, sw))
		if strings.Contains(string(raw), "re-key") {
			t.Fatal("api key leaked through public view")
		}
		if !mustView(t, sw).Secrets.ResendAPIKeySet {
			t.Fatal("secret presence flag must be true after storing a key")
		}
	})

	t.Run("failed switch keeps previous channel", func(t *testing.T) {
		badFrom := "not-an-address"
		if _, err := sw.Update(context.Background(), UpdateRequest{
			Channel:      RuntimeChannelResend,
			ResendAPIKey: strPtr("re-other"),
			ResendFrom:   &badFrom,
		}); err == nil {
			t.Fatal("invalid candidate must fail closed")
		}
		if err := sw.Send(context.Background(), msg("user@example.com", "s3", "b3")); err != nil {
			t.Fatalf("previous channel must keep serving: %v", err)
		}
		if resendHits != 2 {
			t.Fatalf("hits = %d, want 2 (old resend still active)", resendHits)
		}
	})

	t.Run("unknown channel rejected", func(t *testing.T) {
		if _, err := sw.Update(context.Background(), UpdateRequest{Channel: "sendgrid"}); !errors.Is(err, ErrUnknownChannel) {
			t.Fatalf("err = %v, want ErrUnknownChannel", err)
		}
	})
}

// W26 (GOAL-038 D-001 §2.1): every channel's outbound mail lands in
// mail_outbox. Mock records delivered exactly once (no double-write); a real
// channel records sent on success and failed on adapter error; a record
// write never changes the send result.
func TestSwitcherRecordsAllChannelOutbound(t *testing.T) {
	st := openOutboxStore(t)

	failNext := false
	resendSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if failNext {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]any{"message": "boom"})
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "resend-1"})
	}))
	defer resendSrv.Close()

	sw, err := NewSwitcher(st, masterKeyFixed(), SeedConfig{Channel: RuntimeChannelMock}, testSwitchLogger())
	if err != nil {
		t.Fatalf("NewSwitcher: %v", err)
	}
	sw.resendBaseURL = resendSrv.URL

	sink := NewOutboxSink(st, 0)
	if err := sw.Send(context.Background(), msg("u@example.com", "mock-mail", "b")); err != nil {
		t.Fatalf("mock send: %v", err)
	}

	from := "no-reply@example.com"
	if _, err := sw.Update(context.Background(), UpdateRequest{
		Channel:      RuntimeChannelResend,
		ResendAPIKey: strPtr("re-key"),
		ResendFrom:   &from,
	}); err != nil {
		t.Fatalf("switch to resend: %v", err)
	}
	if err := sw.Send(context.Background(), msg("u@example.com", "resend-ok", "b")); err != nil {
		t.Fatalf("resend send: %v", err)
	}
	failNext = true
	if err := sw.Send(context.Background(), msg("u@example.com", "resend-fail", "b")); err == nil {
		t.Fatal("failed resend must surface its error")
	}

	records, _, err := sink.List(context.Background(), OutboxListQuery{PageSize: 10})
	if err != nil {
		t.Fatalf("list records: %v", err)
	}
	statuses := map[string]string{}
	for _, rec := range records {
		statuses[rec.Subject] = rec.Channel + "/" + rec.DeliveryStatus
	}
	want := map[string]string{
		"mock-mail":   ChannelMock + "/" + DeliveryDelivered,
		"resend-ok":   ChannelResend + "/" + DeliverySent,
		"resend-fail": ChannelResend + "/" + DeliveryFailed,
	}
	for subject, expect := range want {
		if statuses[subject] != expect {
			t.Fatalf("record %q = %q, want %q", subject, statuses[subject], expect)
		}
	}
	if len(records) != 3 {
		t.Fatalf("records = %d, want exactly one per send (no mock double-write)", len(records))
	}
}

// D-007: admin intent survives restarts (DB wins over the file seed).
func TestSwitcherConfigSurvivesRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "runtime.db")
	catalog := catalogFor(t)

	st, err := newCatalogStore(path, catalog)
	if err != nil {
		t.Fatal(err)
	}
	sw, err := NewSwitcher(st, masterKeyFixed(), SeedConfig{Channel: RuntimeChannelMock}, testSwitchLogger())
	if err != nil {
		t.Fatal(err)
	}
	retention := 42
	if _, err := sw.Update(context.Background(), UpdateRequest{MockRetention: &retention}); err != nil {
		t.Fatalf("update retention: %v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	st2, err := newCatalogStore(path, catalog)
	if err != nil {
		t.Fatal(err)
	}
	defer st2.Close()
	sw2, err := NewSwitcher(st2, masterKeyFixed(), SeedConfig{Channel: RuntimeChannelMock}, testSwitchLogger())
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	view := mustView(t, sw2)
	if view.MockRetention != 42 {
		t.Fatalf("retention after restart = %d, want 42 (DB wins over seed)", view.MockRetention)
	}
}

func strPtr(s string) *string { return &s }

func testSwitchLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func mustView(t *testing.T, sw *Switcher) *PublicView {
	t.Helper()
	view, err := sw.PublicView()
	if err != nil {
		t.Fatalf("PublicView: %v", err)
	}
	return view
}

func masterKeyFixed() []byte { return []byte(strings.Repeat("m", 32)) }

func catalogFor(t *testing.T) []kernel.MigrationContribution {
	t.Helper()
	catalog, err := compiled.PersistenceCatalog()
	if err != nil {
		t.Fatal(err)
	}
	return catalog
}

func newCatalogStore(path string, catalog []kernel.MigrationContribution) (*store.Store, error) {
	return store.OpenWithCatalog(path, catalog)
}
