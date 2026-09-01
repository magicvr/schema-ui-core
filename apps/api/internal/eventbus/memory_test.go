package eventbus

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

var _ kernel.EventBus = (*Memory)(nil)

type testEvent struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type otherEvent struct {
	Value string `json:"value"`
}

func TestNewMemoryBufferFallback(t *testing.T) {
	tests := []struct {
		input int
		want  int
	}{
		{0, kernel.DefaultEventBusBuffer},
		{-1, kernel.DefaultEventBusBuffer},
		{32, 32},
		{128, 128},
	}
	for _, tt := range tests {
		m := NewMemory(tt.input, nil)
		if m.BufferSize() != tt.want {
			t.Errorf("NewMemory(%d).BufferSize() = %d, want %d", tt.input, m.BufferSize(), tt.want)
		}
	}
}

func TestRegisterValidation(t *testing.T) {
	ctx := context.Background()
	m := NewMemory(64, slog.Default())

	t.Run("invalid topic shape fails", func(t *testing.T) {
		for _, bad := range []kernel.EventTopic{"", "Foo", "foo-", "-foo", "foo..bar", "foo.Bar"} {
			if err := m.Register(ctx, bad, testEvent{}); !errors.Is(err, kernel.ErrInvalidEventTopic) {
				t.Errorf("Register(%q) = %v, want ErrInvalidEventTopic", bad, err)
			}
		}
	})

	t.Run("nil sample fails", func(t *testing.T) {
		if err := m.Register(ctx, "test.event", nil); !errors.Is(err, kernel.ErrInvalidEventPayload) {
			t.Errorf("Register(nil) = %v, want ErrInvalidEventPayload", err)
		}
	})

	t.Run("not JSON-serializable fails", func(t *testing.T) {
		bad := make(chan int)
		if err := m.Register(ctx, "test.event", bad); !errors.Is(err, kernel.ErrEventNotSerializable) {
			t.Errorf("Register(chan) = %v, want ErrEventNotSerializable", err)
		}
	})

	t.Run("valid registration succeeds", func(t *testing.T) {
		if err := m.Register(ctx, "test.event", testEvent{}); err != nil {
			t.Fatalf("Register: %v", err)
		}
	})

	t.Run("same topic same type is idempotent", func(t *testing.T) {
		if err := m.Register(ctx, "test.event", testEvent{}); err != nil {
			t.Errorf("second Register same type: %v", err)
		}
	})

	t.Run("same topic different type fails", func(t *testing.T) {
		if err := m.Register(ctx, "test.event", otherEvent{}); !errors.Is(err, kernel.ErrEventTypeConflict) {
			t.Errorf("Register different type = %v, want ErrEventTypeConflict", err)
		}
	})
}

func TestPublishValidation(t *testing.T) {
	ctx := context.Background()
	m := NewMemory(64, slog.Default())
	_ = m.Register(ctx, "test.event", testEvent{})

	t.Run("invalid topic shape fails", func(t *testing.T) {
		if err := m.Publish(ctx, "Bad", testEvent{}); !errors.Is(err, kernel.ErrInvalidEventTopic) {
			t.Errorf("Publish(Bad) = %v, want ErrInvalidEventTopic", err)
		}
	})

	t.Run("unregistered topic fails", func(t *testing.T) {
		if err := m.Publish(ctx, "other.event", otherEvent{}); !errors.Is(err, kernel.ErrEventTopicNotRegistered) {
			t.Errorf("Publish(unregistered) = %v, want ErrEventTopicNotRegistered", err)
		}
	})

	t.Run("nil event fails", func(t *testing.T) {
		if err := m.Publish(ctx, "test.event", nil); !errors.Is(err, kernel.ErrInvalidEventPayload) {
			t.Errorf("Publish(nil) = %v, want ErrInvalidEventPayload", err)
		}
	})

	t.Run("type mismatch fails", func(t *testing.T) {
		if err := m.Publish(ctx, "test.event", otherEvent{}); !errors.Is(err, kernel.ErrEventTypeMismatch) {
			t.Errorf("Publish(wrong type) = %v, want ErrEventTypeMismatch", err)
		}
	})

	t.Run("not serializable fails", func(t *testing.T) {
		type goodEvent struct {
			ID int `json:"id"`
		}
		type badEvent struct {
			ID int         `json:"id"`
			Fn func() int `json:"fn"`
		}
		// Register with goodEvent succeeds
		if err := m.Register(ctx, "bad.event", goodEvent{ID: 1}); err != nil {
			t.Fatalf("Register goodEvent: %v", err)
		}
		// Publishing badEvent (different type) fails with type mismatch first
		if err := m.Publish(ctx, "bad.event", badEvent{ID: 1, Fn: func() int { return 0 }}); !errors.Is(err, kernel.ErrEventTypeMismatch) {
			t.Errorf("Publish(different type) = %v, want ErrEventTypeMismatch", err)
		}
		// For actual not-serializable test, register a type that CAN'T be marshaled
		type notSerializable struct {
			Fn func() `json:"fn"`
		}
		// Even registering should fail
		if err := m.Register(ctx, "truly.bad", notSerializable{}); !errors.Is(err, kernel.ErrEventNotSerializable) {
			t.Errorf("Register(not serializable) = %v, want ErrEventNotSerializable", err)
		}
	})
}

func TestSubscribeValidation(t *testing.T) {
	ctx := context.Background()
	m := NewMemory(64, slog.Default())
	_ = m.Register(ctx, "test.event", testEvent{})

	t.Run("invalid topic shape fails", func(t *testing.T) {
		h := func(context.Context, []byte) {}
		if _, err := m.Subscribe(ctx, "Bad", h); !errors.Is(err, kernel.ErrInvalidEventTopic) {
			t.Errorf("Subscribe(Bad) = %v, want ErrInvalidEventTopic", err)
		}
	})

	t.Run("unregistered topic fails", func(t *testing.T) {
		h := func(context.Context, []byte) {}
		if _, err := m.Subscribe(ctx, "other.event", h); !errors.Is(err, kernel.ErrEventTopicNotRegistered) {
			t.Errorf("Subscribe(unregistered) = %v, want ErrEventTopicNotRegistered", err)
		}
	})

	t.Run("nil handler fails", func(t *testing.T) {
		if _, err := m.Subscribe(ctx, "test.event", nil); !errors.Is(err, kernel.ErrEventHandlerNil) {
			t.Errorf("Subscribe(nil handler) = %v, want ErrEventHandlerNil", err)
		}
	})
}

func TestPublishSubscribeDelivery(t *testing.T) {
	ctx := context.Background()
	m := NewMemory(64, slog.Default())
	_ = m.Register(ctx, "test.event", testEvent{})

	var received []testEvent
	var mu sync.Mutex
	handler := func(_ context.Context, payload []byte) {
		var e testEvent
		if err := json.Unmarshal(payload, &e); err != nil {
			t.Errorf("unmarshal: %v", err)
			return
		}
		mu.Lock()
		received = append(received, e)
		mu.Unlock()
	}

	sub, err := m.Subscribe(ctx, "test.event", handler)
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}
	defer sub.Unsubscribe()

	events := []testEvent{{ID: 1, Name: "a"}, {ID: 2, Name: "b"}, {ID: 3, Name: "c"}}
	for _, e := range events {
		if err := m.Publish(ctx, "test.event", e); err != nil {
			t.Fatalf("Publish: %v", err)
		}
	}

	time.Sleep(50 * time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	if len(received) != len(events) {
		t.Fatalf("received %d events, want %d", len(received), len(events))
	}
	for i, want := range events {
		if received[i] != want {
			t.Errorf("event[%d] = %+v, want %+v", i, received[i], want)
		}
	}
}

func TestMultipleSubscribers(t *testing.T) {
	ctx := context.Background()
	m := NewMemory(64, slog.Default())
	_ = m.Register(ctx, "test.event", testEvent{})

	var count1, count2 atomic.Int32
	h1 := func(context.Context, []byte) { count1.Add(1) }
	h2 := func(context.Context, []byte) { count2.Add(1) }

	sub1, _ := m.Subscribe(ctx, "test.event", h1)
	sub2, _ := m.Subscribe(ctx, "test.event", h2)
	defer sub1.Unsubscribe()
	defer sub2.Unsubscribe()

	for i := 0; i < 5; i++ {
		_ = m.Publish(ctx, "test.event", testEvent{ID: i})
	}

	time.Sleep(50 * time.Millisecond)
	if got := count1.Load(); got != 5 {
		t.Errorf("subscriber1 received %d, want 5", got)
	}
	if got := count2.Load(); got != 5 {
		t.Errorf("subscriber2 received %d, want 5", got)
	}
}

func TestUnsubscribeStopsDelivery(t *testing.T) {
	ctx := context.Background()
	m := NewMemory(64, slog.Default())
	_ = m.Register(ctx, "test.event", testEvent{})

	var count atomic.Int32
	h := func(context.Context, []byte) { count.Add(1) }
	sub, _ := m.Subscribe(ctx, "test.event", h)

	_ = m.Publish(ctx, "test.event", testEvent{ID: 1})
	time.Sleep(20 * time.Millisecond)
	sub.Unsubscribe()
	_ = m.Publish(ctx, "test.event", testEvent{ID: 2})
	time.Sleep(20 * time.Millisecond)

	if got := count.Load(); got != 1 {
		t.Errorf("after unsubscribe received %d events, want 1", got)
	}

	// Unsubscribe is idempotent
	sub.Unsubscribe()
}

func TestHandlerPanicIsolation(t *testing.T) {
	ctx := context.Background()
	m := NewMemory(64, slog.Default())
	_ = m.Register(ctx, "test.event", testEvent{})

	var ok1, ok2 atomic.Bool
	h1 := func(context.Context, []byte) {
		ok1.Store(true)
		panic("oops")
	}
	h2 := func(context.Context, []byte) {
		ok2.Store(true)
	}

	sub1, _ := m.Subscribe(ctx, "test.event", h1)
	sub2, _ := m.Subscribe(ctx, "test.event", h2)
	defer sub1.Unsubscribe()
	defer sub2.Unsubscribe()

	_ = m.Publish(ctx, "test.event", testEvent{})
	time.Sleep(50 * time.Millisecond)

	if !ok1.Load() {
		t.Error("panicking handler was not called")
	}
	if !ok2.Load() {
		t.Error("second handler was not called after first panicked")
	}
}

func TestStopDrainsAndRejects(t *testing.T) {
	ctx := context.Background()
	m := NewMemory(64, slog.Default())
	_ = m.Register(ctx, "test.event", testEvent{})

	var received atomic.Int32
	blockCh := make(chan struct{})
	h := func(context.Context, []byte) {
		<-blockCh
		received.Add(1)
	}

	sub, _ := m.Subscribe(ctx, "test.event", h)
	defer sub.Unsubscribe()

	// Publish 3 events (they'll block in handler)
	for i := 0; i < 3; i++ {
		_ = m.Publish(ctx, "test.event", testEvent{ID: i})
	}
	time.Sleep(10 * time.Millisecond)

	// Stop in a goroutine
	stopDone := make(chan error, 1)
	go func() {
		stopDone <- m.Stop(context.Background())
	}()

	// Unblock handlers
	close(blockCh)

	// Stop should complete
	if err := <-stopDone; err != nil {
		t.Fatalf("Stop: %v", err)
	}

	if got := received.Load(); got != 3 {
		t.Errorf("Stop drained %d events, want 3", got)
	}

	// After Stop, new operations fail
	if err := m.Register(ctx, "new.event", testEvent{}); !errors.Is(err, kernel.ErrEventBusStopped) {
		t.Errorf("Register after Stop = %v, want ErrEventBusStopped", err)
	}
	if err := m.Publish(ctx, "test.event", testEvent{}); !errors.Is(err, kernel.ErrEventBusStopped) {
		t.Errorf("Publish after Stop = %v, want ErrEventBusStopped", err)
	}
	if _, err := m.Subscribe(ctx, "test.event", func(context.Context, []byte) {}); !errors.Is(err, kernel.ErrEventBusStopped) {
		t.Errorf("Subscribe after Stop = %v, want ErrEventBusStopped", err)
	}

	// Stop is idempotent
	if err := m.Stop(ctx); err != nil {
		t.Errorf("second Stop: %v", err)
	}
}

func TestStopRespectsContext(t *testing.T) {
	ctx := context.Background()
	m := NewMemory(1, slog.Default())
	_ = m.Register(ctx, "test.event", testEvent{})

	blockCh := make(chan struct{})
	h := func(context.Context, []byte) { <-blockCh }
	sub, _ := m.Subscribe(ctx, "test.event", h)
	defer sub.Unsubscribe()

	_ = m.Publish(ctx, "test.event", testEvent{})
	time.Sleep(10 * time.Millisecond)

	stopCtx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := m.Stop(stopCtx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Stop with timeout = %v, want DeadlineExceeded", err)
	}

	close(blockCh)
}

func TestPublishBlocksOnFullBuffer(t *testing.T) {
	ctx := context.Background()
	m := NewMemory(2, slog.Default())
	_ = m.Register(ctx, "test.event", testEvent{})

	blockCh := make(chan struct{})
	started := make(chan struct{})
	h := func(context.Context, []byte) {
		started <- struct{}{}
		<-blockCh
	}
	sub, _ := m.Subscribe(ctx, "test.event", h)
	defer func() {
		close(blockCh)
		sub.Unsubscribe()
	}()

	// Publish first event and wait for handler to start (so it blocks)
	_ = m.Publish(ctx, "test.event", testEvent{ID: 1})
	<-started

	// Now buffer has one processing (blocked in handler) 
	// Fill the channel buffer (capacity 2)
	_ = m.Publish(ctx, "test.event", testEvent{ID: 2})
	_ = m.Publish(ctx, "test.event", testEvent{ID: 3})

	// Fourth publish should block (handler processing one, two in buffer)
	publishDone := make(chan error, 1)
	go func() {
		publishDone <- m.Publish(ctx, "test.event", testEvent{ID: 4})
	}()

	select {
	case <-publishDone:
		t.Fatal("Publish returned immediately, expected to block on full buffer")
	case <-time.After(50 * time.Millisecond):
		// Expected: still blocked
	}

	// Cancel context should unblock a new publish
	cancelCtx, cancel := context.WithCancel(ctx)
	cancel() // Cancel immediately
	err := m.Publish(cancelCtx, "test.event", testEvent{ID: 5})
	if !errors.Is(err, context.Canceled) {
		t.Errorf("Publish with cancelled ctx = %v, want Canceled", err)
	}
}

func TestConcurrentOps(t *testing.T) {
	ctx := context.Background()
	m := NewMemory(64, slog.Default())

	var wg sync.WaitGroup
	topics := []kernel.EventTopic{"topic.a", "topic.b", "topic.c"}

	// Concurrent Register
	for _, topic := range topics {
		wg.Add(1)
		go func(t kernel.EventTopic) {
			defer wg.Done()
			_ = m.Register(ctx, t, testEvent{})
		}(topic)
	}
	wg.Wait()

	var count atomic.Int32
	h := func(context.Context, []byte) { count.Add(1) }

	// Concurrent Subscribe
	var subs []kernel.Subscription
	var subsMu sync.Mutex
	for _, topic := range topics {
		for i := 0; i < 3; i++ {
			wg.Add(1)
			go func(t kernel.EventTopic) {
				defer wg.Done()
				sub, _ := m.Subscribe(ctx, t, h)
				subsMu.Lock()
				subs = append(subs, sub)
				subsMu.Unlock()
			}(topic)
		}
	}
	wg.Wait()

	// Concurrent Publish
	for i := 0; i < 10; i++ {
		for _, topic := range topics {
			wg.Add(1)
			go func(t kernel.EventTopic, id int) {
				defer wg.Done()
				_ = m.Publish(ctx, t, testEvent{ID: id})
			}(topic, i)
		}
	}
	wg.Wait()

	time.Sleep(100 * time.Millisecond)
	want := 3 * 3 * 10 // 3 topics × 3 subs × 10 events
	if got := count.Load(); got != int32(want) {
		t.Errorf("concurrent received %d events, want %d", got, want)
	}

	// Concurrent Unsubscribe
	for _, sub := range subs {
		wg.Add(1)
		go func(s kernel.Subscription) {
			defer wg.Done()
			s.Unsubscribe()
		}(sub)
	}
	wg.Wait()

	_ = m.Stop(ctx)
}

func TestPayloadIsolation(t *testing.T) {
	ctx := context.Background()
	m := NewMemory(64, slog.Default())
	_ = m.Register(ctx, "test.event", testEvent{})

	var mu sync.Mutex
	var payload1, payload2 []byte
	var payload1FirstByte, payload2FirstByte byte
	done1 := make(chan struct{})
	done2 := make(chan struct{})
	
	h1 := func(_ context.Context, p []byte) {
		mu.Lock()
		payload1 = append([]byte(nil), p...)
		payload1FirstByte = p[0]
		p[0] = 99 // Mutate the payload
		mu.Unlock()
		close(done1)
	}
	h2 := func(_ context.Context, p []byte) {
		mu.Lock()
		payload2 = append([]byte(nil), p...)
		payload2FirstByte = p[0]
		mu.Unlock()
		close(done2)
	}

	sub1, _ := m.Subscribe(ctx, "test.event", h1)
	sub2, _ := m.Subscribe(ctx, "test.event", h2)
	defer sub1.Unsubscribe()
	defer sub2.Unsubscribe()

	_ = m.Publish(ctx, "test.event", testEvent{ID: 42, Name: "test"})
	<-done1
	<-done2

	mu.Lock()
	defer mu.Unlock()
	if len(payload1) == 0 || len(payload2) == 0 {
		t.Fatal("handlers did not receive payloads")
	}
	// Both should have received the same original first byte (123 = '{')
	if payload1FirstByte != payload2FirstByte {
		t.Errorf("handlers received different original payloads: h1[0]=%d h2[0]=%d", payload1FirstByte, payload2FirstByte)
	}
	// h1's mutation shouldn't affect h2's copy (each got independent []byte)
	if payload1FirstByte == 99 {
		t.Error("h1 saw its own mutation in the original read, test is broken")
	}
}
