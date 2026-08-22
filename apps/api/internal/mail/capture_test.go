package mail

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// R3 evidence (workspace-017 GOAL-004 D-001): the embedded default sink keeps
// the process SMTP-free while tests can retrieve the last message.

func newTestCaptureSink() (*CaptureSink, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	logger := slog.New(slog.NewTextHandler(buf, nil))
	return NewCaptureSink(logger), buf
}

func TestCaptureSinkStoresLastMessageOnly(t *testing.T) {
	sink, _ := newTestCaptureSink()
	if _, has := sink.Last(); has {
		t.Fatal("fresh sink must report no captured message")
	}
	first := kernel.MailMessage{To: "a@example.com", Subject: "first", TextBody: "1"}
	if err := sink.Send(context.Background(), first); err != nil {
		t.Fatalf("Send: %v", err)
	}
	got, has := sink.Last()
	if !has || got != first {
		t.Fatalf("Last = %+v has=%v, want %+v", got, has, first)
	}
	second := kernel.MailMessage{To: "b@example.com", Subject: "second", TextBody: "2"}
	if err := sink.Send(context.Background(), second); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if got, _ := sink.Last(); got != second {
		t.Fatalf("capacity-one ring must keep only the last message, got %+v", got)
	}
}

func TestCaptureSinkValidateFailsClosed(t *testing.T) {
	sink, _ := newTestCaptureSink()
	if err := sink.Send(context.Background(), kernel.MailMessage{Subject: "s"}); err == nil {
		t.Fatal("message without recipient must fail closed at the port contract")
	}
	if _, has := sink.Last(); has {
		t.Fatal("rejected message must not be captured")
	}
}

func TestCaptureSinkReset(t *testing.T) {
	sink, _ := newTestCaptureSink()
	_ = sink.Send(context.Background(), kernel.MailMessage{To: "a@example.com"})
	sink.Reset()
	if _, has := sink.Last(); has {
		t.Fatal("Reset must clear the captured message")
	}
}

func TestCaptureSinkLogsStructuredLine(t *testing.T) {
	sink, buf := newTestCaptureSink()
	msg := kernel.MailMessage{To: "user@example.com", Subject: "Hello there", TextBody: "body"}
	if err := sink.Send(context.Background(), msg); err != nil {
		t.Fatalf("Send: %v", err)
	}
	line := buf.String()
	for _, want := range []string{"outbound mail captured", "to=user@example.com", `subject="Hello there"`, "bytes=4"} {
		if !strings.Contains(line, want) {
			t.Fatalf("log line missing %q, got: %s", want, line)
		}
	}
}

func TestCaptureSinkNilLoggerDefaults(t *testing.T) {
	sink := NewCaptureSink(nil)
	if err := sink.Send(context.Background(), kernel.MailMessage{To: "a@example.com"}); err != nil {
		t.Fatalf("nil logger must fall back to slog.Default, got error %v", err)
	}
}
