package requestid

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type failingReader struct{}

func (failingReader) Read([]byte) (int, error) { return 0, errRandomUnavailable }

var errRandomUnavailable = errors.New("random unavailable")

func TestValid(t *testing.T) {
	for _, tc := range []struct {
		id    string
		valid bool
	}{
		{"abc-123", true},
		{" trace.id ", true},
		{"", false},
		{"../secret", false},
		{"a b", false},
		{"x" + string(make([]byte, 128)), false},
	} {
		if got := Valid(tc.id); got != tc.valid {
			t.Fatalf("Valid(%q) = %v, want %v", tc.id, got, tc.valid)
		}
	}
}

func TestMiddlewarePropagatesAndGeneratesIDs(t *testing.T) {
	h := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := FromContext(r.Context()); got == "" {
			t.Fatal("missing context correlation id")
		}
	}))

	t.Run("valid incoming id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
		req.Header.Set(HeaderName, "client-123")
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)
		if rr.Header().Get(HeaderName) != "client-123" {
			t.Fatalf("response id = %q", rr.Header().Get(HeaderName))
		}
	})

	t.Run("invalid incoming id is replaced", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
		req.Header.Set(HeaderName, "bad id")
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)
		if got := rr.Header().Get(HeaderName); !Valid(got) || got == "bad id" {
			t.Fatalf("generated response id = %q", got)
		}
	})
}

func TestNewFallbackRemainsUniqueWhenRandomUnavailable(t *testing.T) {
	original := readRandom
	readRandom = failingReader{}.Read
	defer func() { readRandom = original }()

	first, second := New(), New()
	if first == second || !Valid(first) || !Valid(second) {
		t.Fatalf("fallback IDs must be unique and valid: first=%q second=%q", first, second)
	}
}
