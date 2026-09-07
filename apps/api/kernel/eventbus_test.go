package kernel

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
)

// Compile-time port-surface guard (D-002 §10): the frozen signatures must
// stay implementable by any provider — a stub suffices to lock the method sets.
type stubEventBus struct{}

func (stubEventBus) Register(context.Context, EventTopic, any) error { return nil }
func (stubEventBus) Publish(context.Context, EventTopic, any) error  { return nil }
func (stubEventBus) Subscribe(context.Context, EventTopic, EventHandler) (Subscription, error) {
	return stubSubscription{}, nil
}
func (stubEventBus) Stop(context.Context) error { return nil }

type stubSubscription struct{}

func (stubSubscription) Unsubscribe() {}

var (
	_ EventBus     = stubEventBus{}
	_ Subscription = stubSubscription{}
)

func TestDefaultEventBusBuffer(t *testing.T) {
	if got, want := DefaultEventBusBuffer, 64; got != want {
		t.Fatalf("DefaultEventBusBuffer = %d, want %d", got, want)
	}
}

func TestValidEventTopic(t *testing.T) {
	tests := []struct {
		name  string
		topic EventTopic
		want  bool
	}{
		{"single letter", EventTopic("a"), true},
		{"single digit", EventTopic("0"), true},
		{"dotted two", EventTopic("account.created"), true},
		{"dotted three", EventTopic("iam.mfa.enrolled"), true},
		{"digits inside", EventTopic("job2.done"), true},
		{"max length", EventTopic(strings.Repeat("a", EventTopicMaxLen)), true},

		{"empty", EventTopic(""), false},
		{"uppercase", EventTopic("Account"), false},
		{"leading dot", EventTopic(".created"), false},
		{"trailing dot", EventTopic("account."), false},
		{"double dot", EventTopic("account..created"), false},
		{"hyphen", EventTopic("account-created"), false},
		{"underscore", EventTopic("account_created"), false},
		{"slash", EventTopic("account/created"), false},
		{"space", EventTopic("account created"), false},
		{"unicode", EventTopic("账户.创建"), false},
		{"over max length", EventTopic(strings.Repeat("a", EventTopicMaxLen+1)), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidEventTopic(tt.topic); got != tt.want {
				t.Errorf("ValidEventTopic(%q) = %v, want %v", tt.topic, got, tt.want)
			}
		})
	}
}

type sampleEvent struct {
	Name string `json:"name"`
}

func TestEventMustMarshalJSON(t *testing.T) {
	t.Run("struct", func(t *testing.T) {
		b, err := EventMustMarshalJSON(sampleEvent{Name: "ok"})
		if err != nil {
			t.Fatalf("EventMustMarshalJSON(struct) error = %v", err)
		}
		if string(b) != `{"name":"ok"}` {
			t.Fatalf("payload = %s", b)
		}
	})
	t.Run("nil", func(t *testing.T) {
		_, err := EventMustMarshalJSON(nil)
		if !errors.Is(err, ErrInvalidEventPayload) {
			t.Fatalf("nil payload error = %v, want ErrInvalidEventPayload", err)
		}
	})
	t.Run("chan not serializable", func(t *testing.T) {
		ch := make(chan int)
		_, err := EventMustMarshalJSON(ch)
		if !errors.Is(err, ErrEventNotSerializable) {
			t.Fatalf("chan error = %v, want ErrEventNotSerializable", err)
		}
	})
}

func TestValidateEventRegister(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		if err := ValidateEventRegister(EventTopic("account.created"), sampleEvent{}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	t.Run("bad topic first", func(t *testing.T) {
		err := ValidateEventRegister(EventTopic("Account"), sampleEvent{})
		if !errors.Is(err, ErrInvalidEventTopic) {
			t.Fatalf("error = %v, want ErrInvalidEventTopic", err)
		}
	})
	t.Run("nil sample", func(t *testing.T) {
		err := ValidateEventRegister(EventTopic("account.created"), nil)
		if !errors.Is(err, ErrInvalidEventPayload) {
			t.Fatalf("error = %v, want ErrInvalidEventPayload", err)
		}
	})
	t.Run("not serializable", func(t *testing.T) {
		err := ValidateEventRegister(EventTopic("account.created"), make(chan int))
		if !errors.Is(err, ErrEventNotSerializable) {
			t.Fatalf("error = %v, want ErrEventNotSerializable", err)
		}
	})
}

func TestValidateEventPublish(t *testing.T) {
	registered := reflect.TypeOf(sampleEvent{})
	t.Run("ok", func(t *testing.T) {
		if err := ValidateEventPublish(registered, sampleEvent{Name: "x"}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	t.Run("unregistered first", func(t *testing.T) {
		err := ValidateEventPublish(nil, sampleEvent{})
		if !errors.Is(err, ErrEventTopicNotRegistered) {
			t.Fatalf("error = %v, want ErrEventTopicNotRegistered", err)
		}
	})
	t.Run("nil event", func(t *testing.T) {
		err := ValidateEventPublish(registered, nil)
		if !errors.Is(err, ErrInvalidEventPayload) {
			t.Fatalf("error = %v, want ErrInvalidEventPayload", err)
		}
	})
	t.Run("type mismatch pointer", func(t *testing.T) {
		err := ValidateEventPublish(registered, &sampleEvent{Name: "x"})
		if !errors.Is(err, ErrEventTypeMismatch) {
			t.Fatalf("error = %v, want ErrEventTypeMismatch", err)
		}
	})
	t.Run("type mismatch other struct", func(t *testing.T) {
		err := ValidateEventPublish(registered, struct{ Name string }{Name: "x"})
		if !errors.Is(err, ErrEventTypeMismatch) {
			t.Fatalf("error = %v, want ErrEventTypeMismatch", err)
		}
	})
	t.Run("matching type not serializable", func(t *testing.T) {
		type nestedChan struct {
			Ch chan int `json:"ch"`
		}
		err := ValidateEventPublish(reflect.TypeOf(nestedChan{}), nestedChan{Ch: make(chan int)})
		if !errors.Is(err, ErrEventNotSerializable) {
			t.Fatalf("error = %v, want ErrEventNotSerializable", err)
		}
	})
}

func TestValidateEventSubscribe(t *testing.T) {
	handler := func(context.Context, []byte) {}
	t.Run("ok", func(t *testing.T) {
		if err := ValidateEventSubscribe(EventTopic("account.created"), handler); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	t.Run("bad topic first", func(t *testing.T) {
		err := ValidateEventSubscribe(EventTopic(""), handler)
		if !errors.Is(err, ErrInvalidEventTopic) {
			t.Fatalf("error = %v, want ErrInvalidEventTopic", err)
		}
	})
	t.Run("nil handler", func(t *testing.T) {
		err := ValidateEventSubscribe(EventTopic("account.created"), nil)
		if !errors.Is(err, ErrEventHandlerNil) {
			t.Fatalf("error = %v, want ErrEventHandlerNil", err)
		}
	})
}

func TestEventSentinelsWrap(t *testing.T) {
	_, err := EventMustMarshalJSON(make(chan int))
	if !errors.Is(err, ErrEventNotSerializable) {
		t.Fatalf("wrapped marshal must satisfy errors.Is(ErrEventNotSerializable), got %v", err)
	}
}
