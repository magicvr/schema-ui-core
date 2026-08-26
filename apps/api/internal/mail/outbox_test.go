package mail

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/compiled"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
)

func openOutboxStore(t *testing.T) *store.Store {
	t.Helper()
	catalog, err := compiled.PersistenceCatalog()
	if err != nil {
		t.Fatal(err)
	}
	st, err := store.OpenWithCatalog(filepath.Join(t.TempDir(), "outbox.db"), catalog)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return st
}

func msg(to, subject, body string) kernel.MailMessage {
	return kernel.MailMessage{To: to, Subject: subject, TextBody: body}
}

// R6 C2 (GOAL-006 D-002 §3): mock publishes durable admin-inspectable records
// with bounded retention; retrieval faces list newest-first and detail by id.
func TestOutboxSinkPublishesAndLists(t *testing.T) {
	st := openOutboxStore(t)
	sink := NewOutboxSink(st, 0)

	for i := 1; i <= 3; i++ {
		if err := sink.Send(context.Background(), msg("u@example.com", subjectN(i), "b")); err != nil {
			t.Fatalf("send 1..3 failed at %d: %v", i, err)
		}
	}
	if n, err := sink.Count(context.Background()); err != nil || n != 3 {
		t.Fatalf("count = %d, %v; want 3", n, err)
	}
	records, err := sink.List(context.Background(), 10, 0)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(records) != 3 {
		t.Fatalf("list len = %d, want 3", len(records))
	}
	if records[0].Subject != "s3" || records[2].Subject != "s1" {
		t.Fatalf("list must be newest-first, got %+v", subjects(records))
	}
	// W26 (GOAL-038 D-001 §2.1): the list carries the full record — channel,
	// delivery status and body ride every row (contract revision; the
	// declarative recordView drawer renders detail from the selected row).
	for _, rec := range records {
		if rec.Body != "b" || rec.Channel != ChannelMock || rec.DeliveryStatus != DeliveryDelivered {
			t.Fatalf("list row = %+v, want body %q + mock/delivered", rec, "b")
		}
	}
	rec, err := sink.Get(context.Background(), records[0].ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if rec.Body != "b" || rec.To != "u@example.com" {
		t.Fatalf("detail record = %+v", rec)
	}
	if _, err := sink.Get(context.Background(), "outbox-missing"); !errors.Is(err, ErrOutboxRecordNotFound) {
		t.Fatalf("unknown id must return ErrOutboxRecordNotFound, got %v", err)
	}
}

// D-002 §3.5: port validation rules still gate the mock channel.
func TestOutboxSinkValidatesMessage(t *testing.T) {
	st := openOutboxStore(t)
	sink := NewOutboxSink(st, 0)
	if err := sink.Send(context.Background(), msg("not-an-address", "s", "b")); err == nil {
		t.Fatal("invalid recipient must fail closed")
	}
	if n, _ := sink.Count(context.Background()); n != 0 {
		t.Fatalf("invalid send must not persist, count = %d", n)
	}
}

// D-002 §3.3: bounded retention keeps the most recent cap records.
func TestOutboxSinkEvictsOldestBeyondCap(t *testing.T) {
	st := openOutboxStore(t)
	sink := NewOutboxSink(st, 5)
	for i := 0; i < 8; i++ {
		if err := sink.Send(t.Context(), msg("u@example.com", subjectN(i), "b")); err != nil {
			t.Fatalf("send %d: %v", i, err)
		}
	}
	if n, err := sink.Count(t.Context()); err != nil || n != 5 {
		t.Fatalf("count after eviction = %d, %v; want 5", n, err)
	}
	records, err := sink.List(t.Context(), 10, 0)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if records[0].Subject != subjectN(7) || records[4].Subject != subjectN(3) {
		t.Fatalf("eviction must keep newest five (s7..s3), got %+v", subjects(records))
	}
}

// Records survive a full close/reopen of the database (D-002 §3.1 durability).
func TestOutboxSinkSurvivesRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "restart-outbox.db")
	catalog, err := compiled.PersistenceCatalog()
	if err != nil {
		t.Fatal(err)
	}
	st, err := store.OpenWithCatalog(path, catalog)
	if err != nil {
		t.Fatal(err)
	}
	sink := NewOutboxSink(st, 0)
	if err := sink.Send(context.Background(), msg("u@example.com", "persist-me", "body")); err != nil {
		t.Fatalf("send: %v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	st2, err := store.OpenWithCatalog(path, catalog)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer st2.Close()
	sink2 := NewOutboxSink(st2, 0)
	records, err := sink2.List(context.Background(), 10, 0)
	if err != nil || len(records) != 1 {
		t.Fatalf("after restart list = %v, %v; want exactly one persisted record", records, err)
	}
	if records[0].Subject != "persist-me" {
		t.Fatalf("record subject = %q, want persist-me", records[0].Subject)
	}
}

func subjectN(i int) string { return "s" + string(rune('0'+i)) }

func subjects(records []OutboxRecord) []string {
	out := make([]string, 0, len(records))
	for _, r := range records {
		out = append(out, r.Subject)
	}
	return out
}
