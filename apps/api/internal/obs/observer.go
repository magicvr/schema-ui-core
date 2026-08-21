// Package obs owns the kernel observability export face (VP-015 /
// workspace-015 GOAL-002 D-001): a dedicated Prometheus exposition listener,
// the frozen kernel series with a closed label allow-list, and HTTP
// instrumentation that tags every registered route with its owning module_id.
//
// Layering: obs depends on prometheus/client_golang and pkg/version only —
// never on config, kernel or composition — so it stays unit-testable and the
// composition root remains the single wiring point.
package obs

import (
	"crypto/subtle"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/magicvr/schema-ui-core/apps/api/pkg/version"
)

// ModuleIDCore labels centrally-registered routes (health, upload, branding,
// manifest, schema, auth, ...) that have no owning module provider.
const ModuleIDCore = "core"

// MetricsPath is the fixed exposition path on the dedicated listener.
const MetricsPath = "/metrics"

// BuildInfo carries the static suc_build_info label set (R1 D-001 §3).
type BuildInfo struct {
	Version   string
	Commit    string
	GoVersion string
	Profile   string
}

// BuildInfoFromVersion derives the build identity from the link-time version
// package plus the resolved profile name.
func BuildInfoFromVersion(profile string) BuildInfo {
	return BuildInfo{
		Version:   version.Version,
		Commit:    version.Commit,
		GoVersion: runtime.Version(),
		Profile:   profile,
	}
}

// Observer owns the private metrics registry and every frozen kernel series.
// The zero usable value is created by NewObserver; all exported methods are
// nil-receiver safe so composition can wire a disabled surface without nil
// checks at call sites.
type Observer struct {
	registry             *prometheus.Registry
	httpRequestsTotal    *prometheus.CounterVec
	httpRequestDuration  *prometheus.HistogramVec
	kernelModulesEnabled *prometheus.GaugeVec
}

// NewObserver builds the registry with exactly the R1 contract series. No
// other series may be added without a new D record revising the contract.
func NewObserver(build BuildInfo) *Observer {
	reg := prometheus.NewRegistry()

	buildInfo := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "suc_build_info",
		Help: "Build identity of this process; value is always 1.",
		ConstLabels: prometheus.Labels{
			"version":    build.Version,
			"commit":     build.Commit,
			"go_version": build.GoVersion,
			"profile":    build.Profile,
		},
	})
	buildInfo.Set(1)

	httpRequestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "suc_http_requests_total",
			Help: "Requests served by registered routes, partitioned by owning module.",
		},
		[]string{"module_id", "method", "route", "status"},
	)
	httpRequestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "suc_http_request_duration_seconds",
			Help:    "Handler latency for registered routes (fixed default buckets).",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"module_id", "method", "route"},
	)
	kernelModulesEnabled := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "suc_kernel_modules_enabled",
			Help: "Enabled kernel modules at assembly time; value is always 1.",
		},
		[]string{"module_id"},
	)

	reg.MustRegister(buildInfo, httpRequestsTotal, httpRequestDuration, kernelModulesEnabled)
	// Standard bounded Go-runtime / process collectors (R1 D-001 §3).
	reg.MustRegister(collectors.NewGoCollector())
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	return &Observer{
		registry:             reg,
		httpRequestsTotal:    httpRequestsTotal,
		httpRequestDuration:  httpRequestDuration,
		kernelModulesEnabled: kernelModulesEnabled,
	}
}

// RegisterModule records one enabled module in suc_kernel_modules_enabled.
func (o *Observer) RegisterModule(moduleID string) {
	if o == nil || moduleID == "" {
		return
	}
	o.kernelModulesEnabled.WithLabelValues(moduleID).Set(1)
}

// Handler returns the exposition handler served on the dedicated listener. A
// non-empty token requires "Authorization: Bearer <token>" (constant-time
// comparison); an empty token leaves the endpoint open on its bound interface.
func (o *Observer) Handler(authToken string) http.Handler {
	if o == nil {
		return http.NotFoundHandler()
	}
	h := promhttp.HandlerFor(o.registry, promhttp.HandlerOpts{})
	if authToken == "" {
		return h
	}
	expected := []byte(authToken)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const prefix = "Bearer "
		auth := r.Header.Get("Authorization")
		if len(auth) <= len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) ||
			subtle.ConstantTimeCompare([]byte(auth[len(prefix):]), expected) != 1 {
			w.Header().Set("WWW-Authenticate", `Bearer realm="metrics"`)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// splitPattern derives the method/route labels from a ServeMux registration
// pattern ("GET /users/{id}" -> GET, /users/{id}). Patterns without a method
// prefix report an empty method. Route labels are ALWAYS the registration
// pattern, never the request path (R1 D-001 §4 cardinality + secret rules).
func splitPattern(pattern string) (method, route string) {
	if i := strings.IndexByte(pattern, ' '); i >= 0 {
		return pattern[:i], strings.TrimSpace(pattern[i+1:])
	}
	return "", strings.TrimSpace(pattern)
}

// statusRecorder captures the response code for the metrics labels. An
// implicit body write without WriteHeader counts as 200.
type statusRecorder struct {
	http.ResponseWriter
	status int
	wrote  bool
}

func (r *statusRecorder) WriteHeader(code int) {
	if !r.wrote {
		r.status = code
		r.wrote = true
	}
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if !r.wrote {
		r.status = http.StatusOK
		r.wrote = true
	}
	return r.ResponseWriter.Write(b)
}

// Wrap instruments one registered handler: requests and handler latency are
// counted under (owner, method, route[, status]). Nil-safe: a disabled
// Observer returns the next handler unchanged.
func (o *Observer) Wrap(pattern, owner string, next http.Handler) http.Handler {
	if o == nil || next == nil {
		return next
	}
	if owner == "" {
		owner = ModuleIDCore
	}
	method, route := splitPattern(pattern)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		start := time.Now()
		next.ServeHTTP(rec, r)
		o.httpRequestsTotal.WithLabelValues(owner, method, route, strconv.Itoa(rec.status)).Inc()
		o.httpRequestDuration.WithLabelValues(owner, method, route).Observe(time.Since(start).Seconds())
	})
}
