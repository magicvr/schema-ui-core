package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/ratelimit"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
)

// VP-014 GOAL-003 D-001: an explicitly configured object backend extends
// readyz with extra probes; a failing probe keeps the whole gate unavailable.
func TestReadyzObjectProbes(t *testing.T) {
	build := func(probes ...func(context.Context) error) *http.ServeMux {
		env := newAuthTestEnv(t)
		mux := http.NewServeMux()
		RegisterWithMFAProbes(mux, env.a, env.st, operationlog.Recorder(nil), kernel.Plan{}, nil, ratelimit.NewProvider(), nil, nil, probes...)
		_ = auth.IdentityFrom // keep import stable if helpers change
		return mux
	}

	t.Run("passing probe keeps readyz ok", func(t *testing.T) {
		mux := build(func(context.Context) error { return nil })
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/readyz", nil))
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rr.Code)
		}
	})

	t.Run("failing probe flips readyz unavailable", func(t *testing.T) {
		mux := build(func(context.Context) error { return errors.New("head bucket: refused") })
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/readyz", nil))
		if rr.Code != http.StatusServiceUnavailable {
			t.Fatalf("status = %d, want 503", rr.Code)
		}
	})

	t.Run("nil probe entries are ignored (composition passes typed-nil)", func(t *testing.T) {
		mux := build(nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/readyz", nil))
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rr.Code)
		}
	})
}
