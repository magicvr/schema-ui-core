package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/mail"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// workspace-040 R3-A regression: three public wire time fields were serialized by
// encoding/json's default time.Time marshaller, which strips trailing zeros from
// the fractional part. That emitted values such as "2026-09-20T12:57:15Z" or
// "2026-09-20T12:57:15.9Z" — variable width, so a violation of the frozen public
// output contract (Root D-003: exactly six fractional digits, UTC, "Z" only).
//
// The inventory of the Go output surface missed these three fields: they are not
// inline `Format(...)` layouts but struct/DTO fields handed to writeJSON. Each
// test below forces a trailing-zero microsecond instant, which is exactly the
// case the default marshaller renders short.
const (
	trailingZeroInstant = "2026-09-20T12:57:15.900000Z" // bare .9Z before the fix
	wholeSecondInstant  = "2026-09-20T12:57:15.000000Z" // bare ...15Z before the fix
)

var wireInstantPattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{6}Z$`)

// TestDefaultTimeMarshallerBreaksTheWireContract makes the premise of the three
// regression tests above executable: the instants they seed really are rendered
// short by encoding/json's default time.Time marshaller, so those assertions
// cannot pass vacuously. If a future Go release changes the default, this test
// fails first and the comment trail stays honest.
func TestDefaultTimeMarshallerBreaksTheWireContract(t *testing.T) {
	for _, tc := range []struct{ instant, shortForm string }{
		{trailingZeroInstant, `"2026-09-20T12:57:15.9Z"`},
		{wholeSecondInstant, `"2026-09-20T12:57:15Z"`},
	} {
		parsed, err := time.Parse("2006-01-02T15:04:05.000000Z", tc.instant)
		if err != nil {
			t.Fatalf("parse %q: %v", tc.instant, err)
		}
		raw, err := json.Marshal(parsed)
		if err != nil {
			t.Fatalf("marshal %q: %v", tc.instant, err)
		}
		if string(raw) != tc.shortForm {
			t.Fatalf("default marshalling of %s = %s, want the documented short form %s",
				tc.instant, raw, tc.shortForm)
		}
		if got := FormatWireTime(parsed); got != tc.instant {
			t.Fatalf("FormatWireTime(%s) = %q, want the canonical %q", tc.instant, got, tc.instant)
		}
	}
}

func decodeObject(t *testing.T, raw []byte) map[string]json.RawMessage {
	t.Helper()
	var out map[string]json.RawMessage
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("decode %q: %v", raw, err)
	}
	return out
}

// TestHealthProbeTimestampIsCanonicalWire covers /healthz and /readyz, whose
// timestamp used to be a time.Time field on the response struct.
func TestHealthProbeTimestampIsCanonicalWire(t *testing.T) {
	env := newAuthTestEnv(t)

	for _, path := range []string{"/healthz", "/readyz"} {
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, path, nil))
		if rr.Code != http.StatusOK {
			t.Fatalf("GET %s = %d body=%s", path, rr.Code, rr.Body.String())
		}
		body := decodeObject(t, rr.Body.Bytes())
		rawTimestamp, ok := body["timestamp"]
		if !ok {
			t.Fatalf("GET %s has no timestamp: %s", path, rr.Body.String())
		}
		var value any
		if err := json.Unmarshal(rawTimestamp, &value); err != nil {
			t.Fatalf("timestamp is not a JSON scalar on %s: %s", path, rawTimestamp)
		}
		text, isString := value.(string)
		if !isString {
			t.Fatalf("GET %s timestamp = %v (%T), want a fixed-6 UTC string", path, value, value)
		}
		if !wireInstantPattern.MatchString(text) {
			t.Fatalf("GET %s timestamp = %q, want fixed-6 UTC (YYYY-MM-DDTHH:MM:SS.ffffffZ)", path, text)
		}
	}
}

// TestMailOutboxCreatedAtIsCanonicalWire covers both outbox paths: the list
// envelope item projection and the detail body (which used to write the raw
// mail.OutboxRecord struct).
func TestMailOutboxCreatedAtIsCanonicalWire(t *testing.T) {
	env := newAuthTestEnv(t)
	sink := mail.NewOutboxSink(env.st, 0)
	RegisterMailOutbox(env.mux, env.a, sink)

	// Seed one row with a trailing-zero microsecond instant so the pre-fix
	// default marshaller would have printed ".9Z" and the assertion below is
	// decisive rather than incidentally passing.
	if err := env.st.Run(context.Background(), func(tx kernel.Tx) error {
		_, err := tx.Exec(context.Background(),
			`INSERT INTO mail_outbox (id, to_addr, subject, body, channel, delivery_status, created_at)
			 VALUES ('outbox-w040', 'w040@example.com', 'trailing zero', 'body', 'mock', 'delivered', ?)`,
			trailingZeroInstant)
		return err
	}); err != nil {
		t.Fatalf("seed outbox row: %v", err)
	}

	token := adminToken(t, env)
	get := func(path string) map[string]any {
		t.Helper()
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, bearer(t, token, http.MethodGet, path, ""))
		if rr.Code != http.StatusOK {
			t.Fatalf("GET %s = %d body=%s", path, rr.Code, rr.Body.String())
		}
		var out map[string]any
		if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
			t.Fatalf("decode %s: %v", path, err)
		}
		return out
	}

	list := get("/api/mail/outbox?pageSize=10")
	items, _ := list["items"].([]any)
	if len(items) != 1 {
		t.Fatalf("list items = %v, want the seeded row", list["items"])
	}
	row, _ := items[0].(map[string]any)
	if got := row["created_at"]; got != trailingZeroInstant {
		t.Fatalf("list created_at = %v, want %q", got, trailingZeroInstant)
	}

	detail := get("/api/mail/outbox/outbox-w040")
	if got := detail["created_at"]; got != trailingZeroInstant {
		t.Fatalf("detail created_at = %v, want %q", got, trailingZeroInstant)
	}
	// Non-time fields must keep their shape on both paths.
	if detail["subject"] != "trailing zero" || detail["body"] != "body" || detail["delivery_status"] != "delivered" {
		t.Fatalf("detail row = %v, want the seeded fields unchanged", detail)
	}
}

// TestMailConfigWireProjectsExactlyOneUpdatedAt pins the two properties of the
// handler-side projection that a reviewer must not have to take on trust:
//
//  1. A nil view stays nil, so the endpoint keeps encoding JSON null exactly as
//     `writeJSON(w, http.StatusOK, view)` did before the projection existed.
//  2. The embedded mail.PublicView is shadowed by the shallower UpdatedAt, so the
//     encoded object carries exactly ONE "updated_at" key — the canonical string.
func TestMailConfigWireProjectsExactlyOneUpdatedAt(t *testing.T) {
	if got := mailConfigWire(nil); got != nil {
		t.Fatalf("mailConfigWire(nil) = %+v, want nil (JSON null on the wire)", got)
	}
	nilRaw, err := json.Marshal(mailConfigWire(nil))
	if err != nil {
		t.Fatalf("marshal nil projection: %v", err)
	}
	if string(nilRaw) != "null" {
		t.Fatalf("nil projection encodes as %s, want null", nilRaw)
	}

	instant := time.Date(2026, 9, 20, 12, 57, 15, 900_000_000, time.UTC)
	raw, err := json.Marshal(mailConfigWire(&mail.PublicView{
		Channel: mail.RuntimeChannelMock, UpdatedAt: &instant,
	}))
	if err != nil {
		t.Fatalf("marshal projection: %v", err)
	}
	if count := strings.Count(string(raw), `"updated_at"`); count != 1 {
		t.Fatalf("encoded projection has %d \"updated_at\" keys, want exactly 1: %s", count, raw)
	}
	if !strings.Contains(string(raw), `"updated_at":"`+trailingZeroInstant+`"`) {
		t.Fatalf("encoded projection = %s, want the canonical fixed-6 updated_at", raw)
	}
	// The embedded view's other fields are still forwarded.
	if !strings.Contains(string(raw), `"channel":"`+mail.RuntimeChannelMock+`"`) {
		t.Fatalf("encoded projection dropped the embedded channel: %s", raw)
	}
}

// TestMailConfigUpdatedAtIsCanonicalWire covers GET /api/mail/config, whose
// updated_at used to be a *time.Time field on mail.PublicView.
//
// It also pins the R2 ruling (GOAL-004 A-003 F-I-002 user-overruled) that this
// projection must not undo: NULL stays JSON null, never a fabricated instant.
func TestMailConfigUpdatedAtIsCanonicalWire(t *testing.T) {
	env := newAuthTestEnv(t)
	key := []byte(strings.Repeat("m", 32))
	sw, err := mail.NewSwitcher(env.st, key, mail.SeedConfig{Channel: mail.RuntimeChannelMock}, nil)
	if err != nil {
		t.Fatalf("NewSwitcher: %v", err)
	}
	RegisterMailAdmin(env.mux, env.a, sw, env.operations)

	token := adminToken(t, env)
	rawGet := func() map[string]json.RawMessage {
		t.Helper()
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, bearer(t, token, http.MethodGet, "/api/mail/config", ""))
		if rr.Code != http.StatusOK {
			t.Fatalf("GET /api/mail/config = %d body=%s", rr.Code, rr.Body.String())
		}
		return decodeObject(t, rr.Body.Bytes())
	}
	setUpdatedAt := func(value string) {
		t.Helper()
		if err := env.st.Run(context.Background(), func(tx kernel.Tx) error {
			_, err := tx.Exec(context.Background(),
				`UPDATE mail_config SET updated_at = ? WHERE id = 1`, value)
			return err
		}); err != nil {
			t.Fatalf("set mail_config.updated_at: %v", err)
		}
	}

	for _, instant := range []string{trailingZeroInstant, wholeSecondInstant} {
		setUpdatedAt(instant)
		body := rawGet()
		var got string
		if err := json.Unmarshal(body["updated_at"], &got); err != nil {
			t.Fatalf("updated_at is not a string: %s", body["updated_at"])
		}
		if got != instant {
			t.Fatalf("updated_at = %q, want %q", got, instant)
		}
		if !wireInstantPattern.MatchString(got) {
			t.Fatalf("updated_at = %q, want fixed-6 UTC", got)
		}
		// Drift guard: the embedded PublicView must keep forwarding every other
		// field of the model (only updated_at is rewritten).
		for _, field := range []string{"channel", "mockRetention", "resend", "smtp", "secrets"} {
			if _, ok := body[field]; !ok {
				t.Fatalf("GET /api/mail/config dropped %q: %v", field, body)
			}
		}
	}

	// NULL = "never configured" must stay JSON null (R2 F-I-002 user-overruled).
	setUpdatedAtNull := func() {
		t.Helper()
		if err := env.st.Run(context.Background(), func(tx kernel.Tx) error {
			_, err := tx.Exec(context.Background(),
				`UPDATE mail_config SET updated_at = NULL WHERE id = 1`)
			return err
		}); err != nil {
			t.Fatalf("null mail_config.updated_at: %v", err)
		}
	}
	setUpdatedAtNull()
	body := rawGet()
	if got := string(body["updated_at"]); got != "null" {
		t.Fatalf("updated_at = %s, want JSON null when the row has no instant", got)
	}
	if _, ok := body["channel"]; !ok {
		t.Fatalf("config body lost channel after null updated_at: %v", body)
	}
}
