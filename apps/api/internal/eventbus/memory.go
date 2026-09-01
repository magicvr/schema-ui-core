package eventbus

// In-memory provider for the kernel event-bus port (VP-028 / workspace-028
// GOAL-003 D-001, R2): topic→type registry, JSON payloads, per-subscription
// bounded channels with block-on-full, panic-isolated handlers, and Stop
// drain (D-002 §3/§4/§5). Safe for concurrent use.

import (
	"context"
	"log/slog"
	"reflect"
	"sync"
	"sync/atomic"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

var _ kernel.EventBus = (*Memory)(nil)

// Memory implements kernel.EventBus over in-process channels.
type Memory struct {
	buffer  int
	logger  *slog.Logger
	mu      sync.Mutex
	types   map[kernel.EventTopic]reflect.Type
	subs    map[kernel.EventTopic][]*subscription
	stopped bool
	stopCh  chan struct{}
	drainCtx context.Context
	inflight sync.WaitGroup
	nextID  atomic.Uint64
}

type subscription struct {
	id      uint64
	topic   kernel.EventTopic
	handler kernel.EventHandler
	ch      chan []byte
	stop    chan struct{}
	once    sync.Once
	bus     *Memory
}

func (s *subscription) Unsubscribe() {
	s.once.Do(func() {
		close(s.stop)
		s.bus.remove(s)
	})
}

// NewMemory builds the in-process bus. buffer <= 0 falls back to
// kernel.DefaultEventBusBuffer (D-002 §3). logger nil uses slog.Default().
func NewMemory(buffer int, logger *slog.Logger) *Memory {
	if buffer <= 0 {
		buffer = kernel.DefaultEventBusBuffer
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &Memory{
		buffer: buffer,
		logger: logger,
		types:  make(map[kernel.EventTopic]reflect.Type),
		subs:   make(map[kernel.EventTopic][]*subscription),
		stopCh: make(chan struct{}),
	}
}

// BufferSize reports the per-subscription channel capacity.
func (m *Memory) BufferSize() int { return m.buffer }

func (m *Memory) Register(ctx context.Context, topic kernel.EventTopic, sample any) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := kernel.ValidateEventRegister(topic, sample); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.stopped {
		return kernel.ErrEventBusStopped
	}
	got := reflect.TypeOf(sample)
	if existing, ok := m.types[topic]; ok {
		if existing != got {
			return kernel.ErrEventTypeConflict
		}
		return nil
	}
	m.types[topic] = got
	return nil
}

func (m *Memory) Publish(ctx context.Context, topic kernel.EventTopic, event any) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	// F-001: shape first, then lookup, then ValidateEventPublish, then enqueue.
	if !kernel.ValidEventTopic(topic) {
		return kernel.ErrInvalidEventTopic
	}
	m.mu.Lock()
	if m.stopped {
		m.mu.Unlock()
		return kernel.ErrEventBusStopped
	}
	registered := m.types[topic]
	snapshot := append([]*subscription(nil), m.subs[topic]...)
	m.mu.Unlock()
	if err := kernel.ValidateEventPublish(registered, event); err != nil {
		return err
	}
	payload, err := kernel.EventMustMarshalJSON(event)
	if err != nil {
		return err
	}
	for _, sub := range snapshot {
		cp := make([]byte, len(payload))
		copy(cp, payload)
		if err := m.enqueue(ctx, sub, cp); err != nil {
			return err
		}
	}
	return nil
}

func (m *Memory) enqueue(ctx context.Context, sub *subscription, payload []byte) error {
	// Priority check: if already stopped, fail immediately without sending.
	select {
	case <-m.stopCh:
		return kernel.ErrEventBusStopped
	default:
	}
	select {
	case sub.ch <- payload:
		return nil
	case <-sub.stop:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-m.stopCh:
		return kernel.ErrEventBusStopped
	}
}

func (m *Memory) Subscribe(ctx context.Context, topic kernel.EventTopic, handler kernel.EventHandler) (kernel.Subscription, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := kernel.ValidateEventSubscribe(topic, handler); err != nil {
		return nil, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.stopped {
		return nil, kernel.ErrEventBusStopped
	}
	if _, ok := m.types[topic]; !ok {
		return nil, kernel.ErrEventTopicNotRegistered
	}
	sub := &subscription{
		id:      m.nextID.Add(1),
		topic:   topic,
		handler: handler,
		ch:      make(chan []byte, m.buffer),
		stop:    make(chan struct{}),
		bus:     m,
	}
	m.subs[topic] = append(m.subs[topic], sub)
	m.inflight.Add(1)
	go sub.loop()
	return sub, nil
}

func (m *Memory) Stop(ctx context.Context) error {
	m.mu.Lock()
	if m.stopped {
		m.mu.Unlock()
		return nil
	}
	m.stopped = true
	m.drainCtx = ctx
	close(m.stopCh)
	var all []*subscription
	for _, list := range m.subs {
		all = append(all, list...)
	}
	m.mu.Unlock()
	for _, s := range all {
		s.Unsubscribe()
	}
	done := make(chan struct{})
	go func() {
		m.inflight.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (m *Memory) remove(s *subscription) {
	m.mu.Lock()
	defer m.mu.Unlock()
	list := m.subs[s.topic]
	out := list[:0]
	for _, item := range list {
		if item.id != s.id {
			out = append(out, item)
		}
	}
	if len(out) == 0 {
		delete(m.subs, s.topic)
	} else {
		m.subs[s.topic] = out
	}
}

func (s *subscription) loop() {
	defer s.bus.inflight.Done()
	for {
		// Priority check: if stopped, drain and exit (D-002 §5).
		select {
		case <-s.stop:
			s.drain()
			return
		default:
		}
		select {
		case payload := <-s.ch:
			s.run(payload)
		case <-s.stop:
			s.drain()
			return
		}
	}
}

func (s *subscription) drain() {
	for {
		select {
		case payload := <-s.ch:
			s.run(payload)
		default:
			return
		}
	}
}

func (s *subscription) run(payload []byte) {
	defer func() {
		if rec := recover(); rec != nil {
			s.bus.logger.Error("eventbus handler panic", "topic", string(s.topic), "panic", rec)
		}
	}()
	s.handler(s.bus.handlerCtx(), payload)
}

func (m *Memory) handlerCtx() context.Context {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.drainCtx != nil {
		return m.drainCtx
	}
	return context.Background()
}
