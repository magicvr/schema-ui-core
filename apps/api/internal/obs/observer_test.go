package obs

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func scrape(t *testing.T, o *Observer, token string) string {
	t.Helper()
	srv := httptest.NewServer(o.Handler(token))
	t.Cleanup(srv.Close)
	resp, err := http.Get(srv.URL + MetricsPath)
	if err != nil {
		t.Fatalf("GET /metrics: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /metrics status = %d, want 200", resp.StatusCode)
	}
	var sb strings.Builder
	buf := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(buf)
		sb.Write(buf[:n])
		if err != nil {
			break
		}
	}
	return sb.String()
}

// TestObserverExposesKernelSeries verifies the frozen R1 series are present:
// build info, Go runtime collector and per-module gauges.
func TestObserverExposesKernelSeries(t *testing.T) {
	o := NewObserver(BuildInfo{Version: "9.9.9", Commit: "abc1234", GoVersion: "go1.26.0", Profile: "mvp"})
	o.RegisterModule("admin.users")
	o.RegisterModule("core.auth-session")
	body := scrape(t, o, "")

	if !strings.Contains(body, "suc_build_info{commit=\"abc1234\",go_version=\"go1.26.0\",profile=\"mvp\",version=\"9.9.9\"} 1") {
		t.Errorf("suc_build_info missing or wrong labels in:\n%s", body)
	}
	if !strings.Contains(body, `suc_kernel_modules_enabled{module_id="admin.users"} 1`) {
		t.Errorf("suc_kernel_modules_enabled missing admin.users in:\n%s", body)
	}
	if !strings.Contains(body, `suc_kernel_modules_enabled{module_id="core.auth-session"} 1`) {
		t.Errorf("suc_kernel_modules_enabled missing core.auth-session in:\n%s", body)
	}
	if !strings.Contains(body, "go_goroutines") {
		t.Errorf("Go runtime collector series missing in:\n%s", body)
	}
}

// TestWrapLabelsFollowRegistrationPattern pins the R1 D-001 §4 rule: route
// labels come from the registration pattern (wildcards intact), never from
// the request path.
func TestWrapLabelsFollowRegistrationPattern(t *testing.T) {
	o := NewObserver(BuildInfo{Version: "t", Commit: "t", GoVersion: "t", Profile: "t"})
	h := o.Wrap("POST /api/files/{id}", "admin.file-library", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))
	mux := http.NewServeMux()
	mux.Handle("POST /api/files/{id}", h)
	req := httptest.NewRequest(http.MethodPost, "/api/files/deadbeef-42?download=1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	body := scrape(t, o, "")
	want := `suc_http_requests_total{method="POST",module_id="admin.file-library",route="/api/files/{id}",status="201"} 1`
	if !strings.Contains(body, want) {
		t.Errorf("missing %s in:\n%s", want, body)
	}
	if strings.Contains(body, "deadbeef-42") {
		t.Errorf("raw request path leaked into metrics:\n%s", body)
	}
	if !strings.Contains(body, `suc_http_request_duration_seconds_count{method="POST",module_id="admin.file-library",route="/api/files/{id}"}`) {
		t.Errorf("duration histogram series missing in:\n%s", body)
	}
}

// TestWrapDefaultsOwnerToCore and implicit-200 recording.
func TestWrapDefaultsOwnerToCore(t *testing.T) {
	o := NewObserver(BuildInfo{Version: "t", Commit: "t", GoVersion: "t", Profile: "t"})
	h := o.Wrap("GET /healthz", "", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok")) // no explicit WriteHeader -> implicit 200
	}))
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	body := scrape(t, o, "")
	want := `suc_http_requests_total{method="GET",module_id="core",route="/healthz",status="200"} 1`
	if !strings.Contains(body, want) {
		t.Errorf("missing %s in:\n%s", want, body)
	}
}

// TestNilObserverWrapIsPassthrough keeps disabled surfaces zero-cost.
func TestNilObserverWrapIsPassthrough(t *testing.T) {
	var o *Observer
	called := false
	h := o.Wrap("GET /x", "core", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true }))
	h.ServeHTTP(httestRecorder(), httptest.NewRequest(http.MethodGet, "/x", nil))
	if !called {
		t.Fatal("nil observer must pass the handler through")
	}
}

func httestRecorder() *httptest.ResponseRecorder { return httptest.NewRecorder() }

// TestHandlerBearerToken enforces the R1 D-001 §2 guard semantics.
func TestHandlerBearerToken(t *testing.T) {
	const token = "tok-1234567890abcdef"
	o := NewObserver(BuildInfo{Version: "t", Commit: "t", GoVersion: "t", Profile: "t"})

	t.Run("no token configured stays open", func(t *testing.T) {
		body := scrape(t, o, "")
		if !strings.Contains(body, "suc_build_info") {
			t.Error("open handler must expose series")
		}
	})

	t.Run("missing wrong and correct tokens", func(t *testing.T) {
		srv := httptest.NewServer(o.Handler(token))
		t.Cleanup(srv.Close)

		resp, err := http.Get(srv.URL + MetricsPath)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("missing token status = %d, want 401", resp.StatusCode)
		}

		req, _ := http.NewRequest(http.MethodGet, srv.URL+MetricsPath, nil)
		req.Header.Set("Authorization", "Bearer wrong-token-value-999")
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("wrong token status = %d, want 401", resp.StatusCode)
		}

		req.Header.Set("Authorization", "bearer "+token) // scheme case-insensitive
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("correct token status = %d, want 200", resp.StatusCode)
		}
	})
}

// TestSplitPattern covers method-prefix derivation.
func TestSplitPattern(t *testing.T) {
	cases := []struct {
		pattern, method, route string
	}{
		{"GET /users/{id}", "GET", "/users/{id}"},
		{"/healthz", "", "/healthz"},
		{"POST /api/auth/login", "POST", "/api/auth/login"},
	}
	for _, c := range cases {
		m, r := splitPattern(c.pattern)
		if m != c.method || r != c.route {
			t.Errorf("splitPattern(%q) = (%q,%q), want (%q,%q)", c.pattern, m, r, c.method, c.route)
		}
	}
}
