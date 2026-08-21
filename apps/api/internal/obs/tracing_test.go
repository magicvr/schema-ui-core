package obs

import (
	"context"
	"encoding/hex"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

// recorderTracing builds an enabled Tracing backed by an in-memory span
// recorder (no network), for span-shape assertions.
func recorderTracing(t *testing.T) (*Tracing, *tracetest.SpanRecorder) {
	t.Helper()
	rec := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(rec))
	t.Cleanup(func() { _ = tp.Shutdown(context.Background()) })
	return &Tracing{
		enabled:   true,
		tracer:    tp.Tracer("test"),
		extractor: propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	}, rec
}

// TestTracingDisabledIsNoOp pins GOAL-004 D-001 搂3: default path creates no
// provider/backing machinery.
func TestTracingDisabledIsNoOp(t *testing.T) {
	tr := NewTracing(TracingOptions{Enabled: false, Endpoint: ""}, slog.Default())
	if tr.Enabled() {
		t.Fatal("disabled tracing reports Enabled")
	}
	if err := tr.Shutdown(context.Background()); err != nil {
		t.Fatalf("disabled Shutdown must be no-op, got %v", err)
	}
}

// TestServerSpanShapeAndStatusError verifies span name/kind/attributes and
// the >=500 error mapping through the observer Wrap path.
func TestServerSpanShapeAndStatusError(t *testing.T) {
	tr, rec := recorderTracing(t)
	o := NewObserver(BuildInfo{Version: "t", Commit: "t", GoVersion: "t", Profile: "test"})
	o.SetTracing(tr)

	h := o.Wrap("GET /api/files/{id}", "admin.file-library", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	mux := http.NewServeMux()
	mux.Handle("GET /api/files/{id}", h)
	req := httptest.NewRequest(http.MethodGet, "/api/files/abc?download=1", nil)
	mux.ServeHTTP(httptest.NewRecorder(), req)

	spans := rec.Ended()
	if len(spans) != 1 {
		t.Fatalf("ended spans = %d, want 1", len(spans))
	}
	s := spans[0]
	if s.Name() != "GET /api/files/{id}" {
		t.Errorf("span name = %q, want registration pattern", s.Name())
	}
	if s.SpanKind() != trace.SpanKindServer {
		t.Errorf("span kind = %v, want server", s.SpanKind())
	}
	attrs := map[string]string{}
	for _, kv := range s.Attributes() {
		switch v := kv.Value.AsInterface().(type) {
		case string:
			attrs[string(kv.Key)] = v
		case int64:
			attrs[string(kv.Key)] = strconv.FormatInt(v, 10)
		}
	}
	if attrs["http.route"] != "/api/files/{id}" || attrs["http.request.method"] != "GET" {
		t.Errorf("route/method attrs = %v", attrs)
	}
	if attrs["url.path"] != "/api/files/abc" {
		t.Errorf("url.path = %q, want path without query", attrs["url.path"])
	}
	if attrs["http.response.status_code"] != "503" {
		t.Errorf("status attr = %q, want 503", attrs["http.response.status_code"])
	}
	if s.Status().Code != codes.Error {
		t.Errorf("span status = %v, want error for 503", s.Status().Code)
	}
}

// TestServerSpanJoinsUpstreamTrace verifies W3C TraceContext extraction:
// an incoming traceparent makes the span share the caller's trace id.
func TestServerSpanJoinsUpstreamTrace(t *testing.T) {
	tr, rec := recorderTracing(t)
	o := NewObserver(BuildInfo{Version: "t", Commit: "t", GoVersion: "t", Profile: "test"})
	o.SetTracing(tr)

	h := o.Wrap("GET /healthz", "", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }))
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	traceID, _ := hex.DecodeString("4bf92f3577b34da6a3ce929d0e0e4736")
	req.Header.Set("traceparent", "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01")
	h.ServeHTTP(httptest.NewRecorder(), req)

	spans := rec.Ended()
	if len(spans) != 1 {
		t.Fatalf("ended spans = %d, want 1", len(spans))
	}
	got := spans[0].SpanContext().TraceID()
	if len(got) != len(traceID) {
		t.Fatalf("trace id length mismatch")
	}
	for i := range got {
		if got[i] != traceID[i] {
			t.Fatalf("span trace id %x, want upstream %x", got, traceID)
		}
	}
}

// TestTracingExporterDeliversToOTLPSink proves the export path end to end:
// a real OTLP/HTTP exporter + BatchSpanProcessor posts to the configured
// endpoint, and Shutdown force-flushes pending spans.
func TestTracingExporterDeliversToOTLPSink(t *testing.T) {
	var posts atomic.Int64
	sink := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The exporter's URL/path wrapping is SDK-owned; the assertion here is
		// DELIVERY 鈥?a POST carrying protobuf reaches the configured endpoint.
		if r.Method == http.MethodPost {
			_, _ = io.Copy(io.Discard, r.Body)
			posts.Add(1)
			w.WriteHeader(http.StatusOK)
			return
		}
		http.NotFound(w, r)
	}))
	t.Cleanup(sink.Close)

	tr := NewTracing(TracingOptions{
		Enabled:     true,
		Endpoint:    sink.URL,
		SampleRatio: 1.0,
		ServiceName: "test",
	}, slog.Default())
	if !tr.Enabled() {
		t.Fatal("enabled tracing reports disabled")
	}
	_, span := tr.tracer.Start(context.Background(), "smoke", trace.WithSpanKind(trace.SpanKindServer))
	span.End()
	if err := tr.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown: %v", err)
	}

	deadline := time.Now().Add(5 * time.Second)
	for posts.Load() == 0 && time.Now().Before(deadline) {
		time.Sleep(50 * time.Millisecond)
	}
	if posts.Load() != 1 {
		t.Fatalf("OTLP sink received %d posts, want 1", posts.Load())
	}
}
