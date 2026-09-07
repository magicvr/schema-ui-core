package cache

import (
	"errors"
	"fmt"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

type userProfile struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
}

func TestTypedJSONRoundTrip(t *testing.T) {
	m, _ := newTestMemory(t, 10)
	view := mustView(t, m, "wallet")
	tc := NewTyped[userProfile](view)

	u := userProfile{ID: 42, Name: "alice", Email: "alice@example.com"}
	if err := tc.Set(testCtx, "u:42", u, AbsoluteExpiry{TTL: 0}); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, ok, err := tc.Get(testCtx, "u:42")
	if err != nil || !ok {
		t.Fatalf("Get: ok=%v err=%v", ok, err)
	}
	if got != u {
		t.Fatalf("round trip mismatch: %+v != %+v", got, u)
	}
	// Miss semantics stay distinct.
	if _, ok, err := tc.Get(testCtx, "missing"); ok || err != nil {
		t.Fatalf("miss: ok=%v err=%v", ok, err)
	}
	// Delete through the typed layer.
	if err := tc.Delete(testCtx, "u:42"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, ok, _ := tc.Get(testCtx, "u:42"); ok {
		t.Fatal("value must be gone after Delete")
	}
}

func TestTypedDecodeErrorIsNotMiss(t *testing.T) {
	m, _ := newTestMemory(t, 10)
	view := mustView(t, m, "wallet")
	tc := NewTyped[userProfile](view)

	// Corrupt bytes stored directly through the port (the typed layer only
	// sees the value once decoded).
	if err := view.Set(testCtx, "corrupt", []byte("{not json"), AbsoluteExpiry{TTL: 0}); err != nil {
		t.Fatal(err)
	}
	v, ok, err := tc.Get(testCtx, "corrupt")
	if err == nil {
		t.Fatal("decode failure must surface as an error")
	}
	if !ok {
		t.Fatal("a hit with decode failure must report ok=true (not a miss)")
	}
	if v != (userProfile{}) {
		t.Fatalf("zero value expected on decode failure, got %+v", v)
	}
}

// customCodec proves the codec seam is injectable (D-001): a plain-string
// codec with a deterministic prefix.
type prefixCodec struct{ prefix string }

func (c prefixCodec) Encode(v string) ([]byte, error) { return []byte(c.prefix + v), nil }

func (c prefixCodec) Decode(b []byte) (string, error) {
	if len(b) < len(c.prefix) || string(b[:len(c.prefix)]) != c.prefix {
		return "", fmt.Errorf("missing prefix %q", c.prefix)
	}
	return string(b[len(c.prefix):]), nil
}

func TestTypedCustomCodec(t *testing.T) {
	m, _ := newTestMemory(t, 10)
	view := mustView(t, m, "wallet")
	tc := NewTyped[string](view, prefixCodec{prefix: "v1:"})

	if err := tc.Set(testCtx, "k", "hello", AbsoluteExpiry{TTL: 0}); err != nil {
		t.Fatal(err)
	}
	got, ok, err := tc.Get(testCtx, "k")
	if err != nil || !ok || got != "hello" {
		t.Fatalf("custom codec: got=%q ok=%v err=%v", got, ok, err)
	}
	// The raw port sees the encoded form (v1:hello) — proving the codec ran.
	raw, ok := view.Get(testCtx, "k")
	if !ok || string(raw) != "v1:hello" {
		t.Fatalf("raw port value = %q ok=%v, want v1:hello", raw, ok)
	}
}

func TestTypedPortSurface(t *testing.T) {
	// errors from the port pass through unchanged (validation still happens).
	m, _ := newTestMemory(t, 10)
	view := mustView(t, m, "wallet")
	tc := NewTyped[string](view)
	if err := tc.Set(testCtx, "", "v", AbsoluteExpiry{}); !errors.Is(err, kernel.ErrInvalidCacheKey) {
		t.Errorf("typed Set must propagate port validation, got %v", err)
	}
}
