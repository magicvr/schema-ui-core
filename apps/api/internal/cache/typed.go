package cache

// Typed convenience layer over the kernel cache port (VP-026 / workspace-026
// GOAL-002 D-002 §1): the port itself stays []byte and provider-agnostic;
// modules that want compile-time types wrap a kernel.CacheView with Typed.

import (
	"context"
	"encoding/json"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// Codec is the injectable value codec for Typed.
type Codec[T any] interface {
	Encode(T) ([]byte, error)
	Decode([]byte) (T, error)
}

// JSONCodec is the default codec: a plain encoding/json round trip.
type JSONCodec[T any] struct{}

// Encode implements Codec.
func (JSONCodec[T]) Encode(v T) ([]byte, error) { return json.Marshal(v) }

// Decode implements Codec.
func (JSONCodec[T]) Decode(b []byte) (T, error) {
	var v T
	if err := json.Unmarshal(b, &v); err != nil {
		return v, err
	}
	return v, nil
}

// Typed is a type-safe view over a kernel.CacheView. Build with NewTyped; the
// zero value is not ready for use.
type Typed[T any] struct {
	view  kernel.CacheView
	codec Codec[T]
}

// NewTyped wraps view with a codec, defaulting to JSONCodec[T] when none is
// supplied.
func NewTyped[T any](view kernel.CacheView, codec ...Codec[T]) *Typed[T] {
	c := Codec[T](JSONCodec[T]{})
	if len(codec) > 0 && codec[0] != nil {
		c = codec[0]
	}
	return &Typed[T]{view: view, codec: c}
}

// Get decodes a hit into T. A miss returns (zero, false, nil); a decode
// failure returns (zero, true, err) — corruption is never disguised as a
// miss.
func (t *Typed[T]) Get(ctx context.Context, key string) (T, bool, error) {
	var zero T
	b, ok := t.view.Get(ctx, key)
	if !ok {
		return zero, false, nil
	}
	v, err := t.codec.Decode(b)
	if err != nil {
		return zero, true, err
	}
	return v, true, nil
}

// Set encodes and stores value under key with the given expiry policy.
func (t *Typed[T]) Set(ctx context.Context, key string, value T, policy kernel.ExpiryPolicy) error {
	b, err := t.codec.Encode(value)
	if err != nil {
		return err
	}
	return t.view.Set(ctx, key, b, policy)
}

// Delete removes key (idempotent at the port).
func (t *Typed[T]) Delete(ctx context.Context, key string) error {
	return t.view.Delete(ctx, key)
}
