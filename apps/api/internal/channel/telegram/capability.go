package telegram

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

const (
	// CapabilityCacheNamespace isolates Telegram operator capability decisions
	// from every other cache consumer.
	CapabilityCacheNamespace kernel.CacheNamespace = "telegram-operator-capability"
	// CapabilityCacheTTL is the absolute lifetime of both allow and deny
	// decisions. The service never refreshes it on read.
	CapabilityCacheTTL = 60 * time.Second
)

// ErrTelegramCapabilityUnavailable is intentionally stable at the service
// boundary; the HTTP handler maps it to the cataloged 502 response without
// exposing Telegram diagnostics.
var ErrTelegramCapabilityUnavailable = errors.New("telegram: capability is unavailable")

// ChatMemberGetter is the narrow Bot API seam required by the capability
// service. It keeps the cache and membership policy independent from HTTP.
type ChatMemberGetter interface {
	GetChatMember(context.Context, int64, int64) (ChatMember, error)
}

// CapabilityInvalidator is the exact-key invalidation seam used after a real
// sendMessage 403.
type CapabilityInvalidator interface {
	Invalidate(context.Context, int64, int64) error
}

type capabilityFlight struct {
	done    chan struct{}
	allowed bool
	err     error
	waiters int
}

// CapabilityService owns the server-side Telegram capability cache and the
// member-status policy. Its flight map is process-local, matching the
// process-local cache provider and preventing duplicate same-key probes.
type CapabilityService struct {
	api        ChatMemberGetter
	cache      kernel.CacheView
	initErr    error
	mu         sync.Mutex
	flights    map[string]*capabilityFlight
	generation map[string]uint64
}

// NewCapabilityService constructs the channel-owned capability service. A
// missing cache or Bot API seam is retained as an initialization error so the
// service fails closed at request time while keeping test-only composition
// seams able to mount unrelated Telegram routes.
func NewCapabilityService(api ChatMemberGetter, cachePort kernel.Cache) *CapabilityService {
	service := &CapabilityService{
		api:        api,
		flights:    make(map[string]*capabilityFlight),
		generation: make(map[string]uint64),
	}
	if api == nil {
		service.initErr = fmt.Errorf("%w: Bot API client is unavailable", ErrTelegramCapabilityUnavailable)
		return service
	}
	if cachePort == nil {
		service.initErr = fmt.Errorf("%w: cache port is unavailable", ErrTelegramCapabilityUnavailable)
		return service
	}
	view, err := cachePort.Namespace(CapabilityCacheNamespace)
	if err != nil {
		service.initErr = fmt.Errorf("%w: initialize cache namespace: %v", ErrTelegramCapabilityUnavailable, err)
		return service
	}
	service.cache = view
	return service
}

// Check returns the current bot's ability to send to chatID. force bypasses a
// cached value but still joins an existing same-key probe. The cache lookup
// and flight creation are serialized together so a normal request arriving
// during a forced probe cannot return the stale value.
func (s *CapabilityService) Check(ctx context.Context, botID, chatID int64, force bool) (bool, error) {
	if s == nil {
		return false, ErrTelegramCapabilityUnavailable
	}
	if s.initErr != nil {
		return false, s.initErr
	}
	if botID <= 0 || chatID == 0 {
		return false, fmt.Errorf("%w: invalid bot or chat id", ErrTelegramCapabilityUnavailable)
	}
	key := capabilityCacheKey(botID, chatID)

	s.mu.Lock()
	if flight, ok := s.flights[key]; ok {
		flight.waiters++
		s.mu.Unlock()
		return waitCapabilityFlight(ctx, flight)
	}
	if !force {
		if allowed, ok := decodeCapabilityValue(s.cache.Get(ctx, key)); ok {
			s.mu.Unlock()
			return allowed, nil
		}
	}
	flight := &capabilityFlight{done: make(chan struct{})}
	s.flights[key] = flight
	generation := s.generation[key]
	s.mu.Unlock()

	allowed, err := s.probe(ctx, botID, chatID)
	s.mu.Lock()
	if err == nil && generation == s.generation[key] {
		if cacheErr := s.cache.Set(ctx, key, encodeCapabilityValue(allowed), absoluteCapabilityExpiry{}); cacheErr != nil {
			allowed = false
			err = fmt.Errorf("%w: cache result: %v", ErrTelegramCapabilityUnavailable, cacheErr)
		}
	}
	flight.allowed = allowed
	flight.err = err
	delete(s.flights, key)
	close(flight.done)
	s.mu.Unlock()
	return allowed, err
}

func (s *CapabilityService) probe(ctx context.Context, botID, chatID int64) (bool, error) {
	member, err := s.api.GetChatMember(ctx, chatID, botID)
	if err != nil {
		if IsTelegramForbidden(err) {
			return false, nil
		}
		return false, fmt.Errorf("%w: getChatMember failed", ErrTelegramCapabilityUnavailable)
	}
	return chatMemberCanSend(member), nil
}

// Invalidate removes only the current bot/chat entry. Deleting a missing key
// is intentionally idempotent under the kernel cache contract.
func (s *CapabilityService) Invalidate(ctx context.Context, botID, chatID int64) error {
	if s == nil || s.initErr != nil || s.cache == nil || botID <= 0 || chatID == 0 {
		if s != nil && s.initErr != nil {
			return s.initErr
		}
		return ErrTelegramCapabilityUnavailable
	}
	key := capabilityCacheKey(botID, chatID)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.generation[key]++
	return s.cache.Delete(ctx, key)
}

func waitCapabilityFlight(ctx context.Context, flight *capabilityFlight) (bool, error) {
	select {
	case <-flight.done:
		return flight.allowed, flight.err
	case <-ctx.Done():
		return false, ctx.Err()
	}
}

func capabilityCacheKey(botID, chatID int64) string {
	return fmt.Sprintf("bot:%d/chat:%d", botID, chatID)
}

func encodeCapabilityValue(allowed bool) []byte {
	if allowed {
		return []byte{'1'}
	}
	return []byte{'0'}
}

func decodeCapabilityValue(value []byte, ok bool) (bool, bool) {
	if !ok || len(value) != 1 {
		return false, false
	}
	switch value[0] {
	case '0':
		return false, true
	case '1':
		return true, true
	default:
		return false, false
	}
}

func chatMemberCanSend(member ChatMember) bool {
	if member.CanSendMessages != nil && !*member.CanSendMessages {
		return false
	}
	if member.CanPostMessages != nil && !*member.CanPostMessages {
		return false
	}
	switch member.Status {
	case "creator", "member", "administrator":
		return true
	case "restricted":
		return member.CanSendMessages != nil && *member.CanSendMessages
	default:
		return false
	}
}

type absoluteCapabilityExpiry struct{}

func (absoluteCapabilityExpiry) ExpireAt(now time.Time) time.Time {
	return now.Add(CapabilityCacheTTL)
}

func (absoluteCapabilityExpiry) Refresh(_ time.Time, previous time.Time) (time.Time, bool) {
	return previous, false
}

// CapabilityInvalidatingSender preserves the kernel sender contract while
// invalidating the exact capability entry when the external send proves the
// cached decision stale.
type CapabilityInvalidatingSender struct {
	delegate    kernel.TelegramSender
	invalidator CapabilityInvalidator
	botID       func() int64
}

func NewCapabilityInvalidatingSender(delegate kernel.TelegramSender, invalidator CapabilityInvalidator, botID func() int64) *CapabilityInvalidatingSender {
	return &CapabilityInvalidatingSender{delegate: delegate, invalidator: invalidator, botID: botID}
}

func (s *CapabilityInvalidatingSender) Send(ctx context.Context, msg kernel.TelegramMessage) error {
	if s == nil || s.delegate == nil {
		return errors.New("telegram: sender is unavailable")
	}
	err := s.delegate.Send(ctx, msg)
	if err == nil || !IsTelegramForbidden(err) || s.invalidator == nil || s.botID == nil {
		return err
	}
	chatID, parseErr := strconv.ParseInt(msg.ChatID, 10, 64)
	if parseErr != nil || chatID == 0 {
		return err
	}
	if invalidateErr := s.invalidator.Invalidate(ctx, s.botID(), chatID); invalidateErr != nil {
		return errors.Join(err, fmt.Errorf("telegram: invalidate capability cache: %w", invalidateErr))
	}
	return err
}
