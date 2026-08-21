// Package requestid owns the request correlation contract shared by the HTTP
// boundary and handlers. Incoming IDs are accepted only when they are short,
// printable, and header-safe; invalid values are replaced with a fresh ID.
package requestid

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync/atomic"
	"time"
)

const (
	HeaderName = "X-Request-ID"
	BodyName   = "correlation_id"
	maxLength  = 128
)

type contextKey struct{}

var validPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)
var readRandom = rand.Read

// Valid reports whether id is safe to propagate as a request correlation ID.
func Valid(id string) bool {
	return validPattern.MatchString(strings.TrimSpace(id))
}

// FromContext returns the correlation ID attached by Middleware, or an empty
// string when the handler is being exercised without the HTTP boundary.
func FromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	id, _ := ctx.Value(contextKey{}).(string)
	return id
}

// WithContext attaches a validated ID to ctx.
func WithContext(ctx context.Context, id string) context.Context {
	if !Valid(id) {
		return ctx
	}
	return context.WithValue(ctx, contextKey{}, strings.TrimSpace(id))
}

// New returns a fresh, opaque printable ID. The timestamp fallback preserves
// uniqueness if the process cannot access the system CSPRNG.
func New() string {
	var raw [16]byte
	if _, err := readRandom(raw[:]); err == nil {
		return hex.EncodeToString(raw[:])
	}
	return fmt.Sprintf("%x-%x", time.Now().UTC().UnixNano(), fallbackSequence.Add(1))
}

var fallbackSequence atomic.Uint64

// Middleware establishes one correlation ID for every request and emits it on
// every response, including route/auth errors produced by downstream layers.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.Header.Get(HeaderName))
		if !Valid(id) {
			id = New()
		}
		w.Header().Set(HeaderName, id)
		next.ServeHTTP(w, r.WithContext(WithContext(r.Context(), id)))
	})
}
