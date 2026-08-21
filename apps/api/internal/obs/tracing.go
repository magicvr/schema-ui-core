package obs

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

// TracingOptions carries the raw observability.traces config values plus the
// resource identity (obs does not import config; composition maps fields).
type TracingOptions struct {
	Enabled     bool
	Endpoint    string
	SampleRatio float64
	ServiceName string
	Version     string
	Environment string
}

// Tracing owns the (optional) OpenTelemetry tracer path (GOAL-004 D-001 §3):
// disabled means a pure no-op — no provider, no exporter, no background
// batch goroutine — so the mvp/dev default carries zero cost. When enabled,
// spans are exported asynchronously by a BatchSpanProcessor; export failures
// only log (never crash the process, never touch the request path).
type Tracing struct {
	enabled   bool
	provider  *sdktrace.TracerProvider // nil when disabled
	tracer    trace.Tracer
	extractor propagation.TextMapPropagator
}

// NewTracing builds the tracer path. Endpoint shape was already validated
// fail-closed at config load; a residual exporter-construction error (e.g.
// an unresolvable endpoint) degrades to no-op with a loud log rather than
// blocking startup.
func NewTracing(opts TracingOptions, logger *slog.Logger) *Tracing {
	if logger == nil {
		logger = slog.Default()
	}
	prop := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	if !opts.Enabled {
		return &Tracing{tracer: noop.NewTracerProvider().Tracer(""), extractor: prop}
	}
	if opts.SampleRatio <= 0 || opts.SampleRatio > 1 {
		opts.SampleRatio = 1.0 // defensive: validated upstream; 0 means "not set" on hand-built configs
	}
	exporter, err := otlptracehttp.New(context.Background(), otlptracehttp.WithEndpointURL(opts.Endpoint))
	if err != nil {
		logger.Error("observability: tracing exporter init failed; traces disabled", "err", err)
		return &Tracing{tracer: noop.NewTracerProvider().Tracer(""), extractor: prop}
	}
	res, err := resource.Merge(resource.Default(), resource.NewSchemaless(
		attribute.String("service.name", opts.ServiceName),
		attribute.String("service.version", opts.Version),
		attribute.String("deployment.environment", opts.Environment),
	))
	if err != nil {
		res = resource.Default()
	}
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(opts.SampleRatio))),
		sdktrace.WithResource(res),
	)
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		logger.Warn("observability: otel export issue", "err", err)
	}))
	return &Tracing{
		enabled:   true,
		provider:  provider,
		tracer:    provider.Tracer("github.com/magicvr/schema-ui-core/apps/api/internal/obs"),
		extractor: prop,
	}
}

// Enabled reports whether spans will actually be created.
func (t *Tracing) Enabled() bool { return t != nil && t.enabled }

// Shutdown flushes and stops the provider (joined into the app stop chain);
// disabled or never-started tracers return nil.
func (t *Tracing) Shutdown(ctx context.Context) error {
	if t == nil || t.provider == nil {
		return nil
	}
	err := t.provider.Shutdown(ctx)
	t.provider = nil
	return err
}

// serverSpan extracts an incoming W3C trace context (joining an upstream
// trace) and starts a SERVER span for one registered route. Attributes stay
// bounded per GOAL-004 D-001 §5: method/route come from the registration
// pattern and url.path carries no query string (credentials never enter).
func (t *Tracing) serverSpan(r *http.Request, method, route string) (context.Context, trace.Span) {
	if !t.Enabled() {
		return r.Context(), nil
	}
	ctx := t.extractor.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
	name := route
	if method != "" {
		name = method + " " + route
	}
	return t.tracer.Start(ctx, name,
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(
			attribute.String("http.request.method", method),
			attribute.String("http.route", route),
			attribute.String("url.path", r.URL.Path),
		),
	)
}

// finishServerSpan stamps the response status and ends the span. Status >=
// 500 marks the span as errored; the observation stays a bypass — nothing
// here can fail the request.
func finishServerSpan(span trace.Span, status int) {
	if span == nil {
		return
	}
	span.SetAttributes(attribute.Int("http.response.status_code", status))
	if status >= 500 {
		span.SetStatus(codes.Error, "HTTP "+strconv.Itoa(status))
	}
	span.End()
}
