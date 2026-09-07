package kernel

// Kernel event-bus port (VP-028 / workspace-028 GOAL-002 D-002, R1).
//
// The port is the only in-process event transport contract for the kernel and
// every module: a topic→type registry, JSON-serializable payloads, asynchronous
// per-subscription delivery with block-on-full buffers, and an explicit Stop
// drain (VP-021). Public types carry neither provider handles nor broker
// clients. Domain code consumes EventBus, never a channel implementation.
//
// Contract frozen by workspace-028 GOAL-002 D-002:
//
//   - Non-generic port. Register freezes topic → reflect.Type via a sample
//     value that must json.Marshal; Publish requires an exact type match and
//     a successful json.Marshal; handlers receive the JSON []byte copy.
//   - Topic shape is fail-closed (ValidEventTopic). Unknown topics are a
//     programming error at Publish/Subscribe (ErrEventTopicNotRegistered).
//   - Delivery is asynchronous: Publish returns after enqueueing to every
//     current subscriber buffer and does not wait for handlers. A full
//     buffer blocks the publisher (ctx cancel / Stop abort).
//   - Handlers have no error return (swallow is structural). Providers MUST
//     recover panics and log. Stop drains buffers, waits in-flight handlers
//     against ctx, then rejects further Register/Publish/Subscribe.
//   - All EventBus methods are safe for concurrent use.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
)

// EventTopic names a typed event stream. The set is OPEN (future business
// modules create their own topics) but every value must pass ValidEventTopic
// or adapters fail closed (GOAL-002 D-002 §2).
type EventTopic string

const (
	// EventTopicMaxLen is the upper bound on topic length (in bytes).
	EventTopicMaxLen = 128
	// DefaultEventBusBuffer is the fallback per-subscription buffer size
	// when eventbus.buffer_size is unset or <= 0 (D-002 §3).
	DefaultEventBusBuffer = 64
)

// eventTopicPattern is the only topic shape the port accepts: one or more
// lower-case alphanumeric segments joined by single dots — must not start/end
// with a dot and must not contain consecutive dots. Enforcing it at the port
// keeps crafted topics from colliding with future outbox routing keys.
var eventTopicPattern = regexp.MustCompile(`^[a-z0-9]+(\.[a-z0-9]+)*$`)

// ValidEventTopic reports whether topic is a well-formed event topic.
// Unknown topics are NOT accepted here as a "soft miss" — invalid shapes
// are rejected before touching the provider.
func ValidEventTopic(topic EventTopic) bool {
	return len(topic) > 0 && len(topic) <= EventTopicMaxLen && eventTopicPattern.MatchString(string(topic))
}

// Sentinel errors (GOAL-002 D-002 §7). Callers use errors.Is.
var (
	ErrInvalidEventTopic       = errors.New("kernel: invalid event topic")
	ErrInvalidEventPayload     = errors.New("kernel: invalid event payload")
	ErrEventNotSerializable    = errors.New("kernel: event payload is not JSON-serializable")
	ErrEventTypeMismatch       = errors.New("kernel: event type does not match registration")
	ErrEventTypeConflict       = errors.New("kernel: event topic already registered with a different type")
	ErrEventTopicNotRegistered = errors.New("kernel: event topic is not registered")
	ErrEventHandlerNil         = errors.New("kernel: event handler is nil")
	ErrEventBusStopped         = errors.New("kernel: event bus is stopped")
)

// EventHandler is the consume callback (D-002 §1/§4). There is no error
// return: handler failure cannot propagate to Publish. payload is a JSON
// copy produced at Publish; the handler may retain or mutate it.
type EventHandler func(ctx context.Context, payload []byte)

// Subscription is the lifecycle handle for one Subscribe (D-002 §1/§4).
// Unsubscribe is idempotent; an in-flight handler is allowed to finish.
type Subscription interface {
	Unsubscribe()
}

// EventBus is the kernel event-transport port (R1). Implementations MUST be
// safe for concurrent use, MUST call the validation helpers below before
// touching internal state, and MUST recover handler panics.
type EventBus interface {
	// Register freezes topic → type from sample. Same topic+type is
	// idempotent; same topic+different type is ErrEventTypeConflict.
	Register(ctx context.Context, topic EventTopic, sample any) error
	// Publish type-checks and JSON-marshals event, then enqueues a copy
	// onto every current subscriber buffer. It does not wait for handlers.
	// A full buffer blocks until space, ctx cancel, or Stop.
	Publish(ctx context.Context, topic EventTopic, event any) error
	// Subscribe appends a handler on an already-registered topic.
	Subscribe(ctx context.Context, topic EventTopic, handler EventHandler) (Subscription, error)
	// Stop rejects further Register/Publish/Subscribe, unblocks publishers
	// waiting on a full buffer, drains queued events, waits in-flight
	// handlers against ctx, then cancels every subscription. Idempotent.
	Stop(ctx context.Context) error
}

// EventMustMarshalJSON is the executable serializability authority
// (D-002 §2): nil → ErrInvalidEventPayload; json.Marshal failure wraps
// ErrEventNotSerializable. Providers MUST use it for Register and Publish.
func EventMustMarshalJSON(v any) ([]byte, error) {
	if v == nil {
		return nil, ErrInvalidEventPayload
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrEventNotSerializable, err)
	}
	return b, nil
}

// ValidateEventRegister is the executable fail-closed entry for Register
// (D-002 §2): topic shape, then JSON serializability. Providers MUST call
// it before recording a type. Type-conflict detection is stateful and stays
// with the provider (ErrEventTypeConflict).
func ValidateEventRegister(topic EventTopic, sample any) error {
	if !ValidEventTopic(topic) {
		return ErrInvalidEventTopic
	}
	_, err := EventMustMarshalJSON(sample)
	return err
}

// ValidateEventPublish is the executable fail-closed entry for Publish
// (D-002 §2): registered type present, event non-nil, exact type match,
// then JSON serializability — in that order. Providers MUST call it before
// enqueueing. A nil registered type means the topic was never registered.
func ValidateEventPublish(registered reflect.Type, event any) error {
	if registered == nil {
		return ErrEventTopicNotRegistered
	}
	if event == nil {
		return ErrInvalidEventPayload
	}
	if reflect.TypeOf(event) != registered {
		return ErrEventTypeMismatch
	}
	_, err := EventMustMarshalJSON(event)
	return err
}

// ValidateEventSubscribe is the executable fail-closed entry for Subscribe
// (D-002 §2/§4): topic shape, then non-nil handler. Whether the topic is
// registered is stateful and stays with the provider
// (ErrEventTopicNotRegistered).
func ValidateEventSubscribe(topic EventTopic, handler EventHandler) error {
	if !ValidEventTopic(topic) {
		return ErrInvalidEventTopic
	}
	if handler == nil {
		return ErrEventHandlerNil
	}
	return nil
}
